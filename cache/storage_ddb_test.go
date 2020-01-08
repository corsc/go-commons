// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/corsc/go-commons/testing/skip"
	"github.com/stretchr/testify/assert"
)

const (
	DDBTestFlag = "DDB"
)

func TestDynamoDbStorage_implements(t *testing.T) {
	assert.Implements(t, (*Storage)(nil), &DynamoDbStorage{})
}

func TestDynamoDbStorage_happyPath(t *testing.T) {
	skip.IfNotSet(t, DDBTestFlag)

	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()

	storage := getTestDynamoDbStorage()

	// get a value (should fail)
	result, resultErr := storage.Get(ctx, key)
	assert.Nil(t, result)
	assert.Equal(t, ErrCacheMiss, resultErr)

	// set a value
	data := []byte(`this is foo`)
	resultErr = storage.Set(ctx, key, data)
	assert.Nil(t, resultErr)

	// get a value
	result, resultErr = storage.Get(ctx, key)
	assert.Equal(t, data, result)
	assert.Nil(t, resultErr)
}

func TestDynamoDbStorage_Invalidate(t *testing.T) {
	skip.IfNotSet(t, DDBTestFlag)

	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()

	storage := getTestDynamoDbStorage()

	// set a value
	data := []byte(`this is foo`)
	resultErr := storage.Set(ctx, key, data)
	assert.Nil(t, resultErr)

	// get a value
	result, resultErr := storage.Get(ctx, key)
	assert.Equal(t, data, result)
	assert.Nil(t, resultErr)

	// invalidate that value
	resultErr = storage.Invalidate(ctx, key)
	assert.Nil(t, resultErr)

	// get a value (should fail)
	result, resultErr = storage.Get(ctx, key)
	assert.Nil(t, result)
	assert.Equal(t, ErrCacheMiss, resultErr)
}

func TestDynamoDbStorage_getWithCtxDone(t *testing.T) {
	skip.IfNotSet(t, DDBTestFlag)

	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	key := getTestKey()

	storage := getTestDynamoDbStorage()

	// attempt to get with a cancelled context
	cancelFn()

	result, resultErr := storage.Get(ctx, key)
	assert.Nil(t, result)
	assert.Equal(t, context.Canceled, resultErr)
}

func TestDynamoDbStorage_setWithCtxDone(t *testing.T) {
	skip.IfNotSet(t, DDBTestFlag)

	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	key := getTestKey()

	storage := getTestDynamoDbStorage()

	// attempt to get with a cancelled context
	cancelFn()

	resultErr := storage.Set(ctx, key, []byte("this is foo"))
	assert.Equal(t, context.Canceled, resultErr)
}

func getTestDynamoDbStorage() *DynamoDbStorage {
	creds := credentials.NewStaticCredentials("123", "123", "")

	svc := dynamodb.New(
		session.New(),
		aws.NewConfig().
			WithLogger(aws.LoggerFunc(log.Print)).
			WithCredentials(creds).
			WithRegion("ap-southeast-2").
			WithEndpoint("http://localhost:8000"),
	)

	tableName := "cachetest"

	createCalled := false
	timer := time.After(3 * time.Second)
	for {
		select {
		case <-timer:
			panic("failed to create table")

		default:
			if testTableExists(svc, tableName) {
				return &DynamoDbStorage{
					Service:   svc,
					TTL:       60 * time.Second,
					TableName: tableName,
				}
			}

			if !createCalled {
				createTestTable(svc, tableName)
			}
		}
	}
}

func createTestTable(svc *dynamodb.DynamoDB, tableName string) {
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(ddbKey),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(ddbKey),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})

	if err != nil {
		panic(err)
	}
}

func testTableExists(svc *dynamodb.DynamoDB, tableName string) bool {
	resp, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return false
	}

	if *(resp.Table.TableStatus) == "ACTIVE" {
		return true
	}

	return false
}
