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
	"strconv"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DynamoDbStorage implements Storage
//
// It is strongly recommended that users customize the circuit breaker settings with a call similar to:
//
//    hystrix.ConfigureCommand(cache.CbDynamoDbStorage, hystrix.CommandConfig{
//        Timeout: 1 * 1000,
//        MaxConcurrentRequests: 1000,
//        ErrorPercentThreshold: 50,
//        })
//
type DynamoDbStorage struct {
	// Service is the AWS DDB Client instance
	Service dynamodbiface.DynamoDBAPI

	// TableName is the AWS DDB Table name
	TableName string

	// TTL is the max TTL for cache items (required)
	TTL time.Duration
}

// Get implements Storage
func (r *DynamoDbStorage) Get(ctx context.Context, key string) ([]byte, error) {
	resultCh := make(chan []byte, 1)
	errorCh := hystrix.Go(CbDynamoDbStorage, func() error {
		params := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				ddbKey: {
					S: aws.String(key),
				},
			},
			TableName: aws.String(r.TableName),
			AttributesToGet: []*string{
				aws.String(ddbData),
			},
		}

		resp, err := r.Service.GetItemWithContext(ctx, params)
		if err != nil {
			return err
		}

		if len(resp.Item) == 0 {
			// cache miss (cannot be returned as error or the CB will track it)
			resultCh <- nil
			return nil
		}

		resultCh <- resp.Item[ddbData].B
		return nil
	}, nil)

	select {
	case result := <-resultCh:
		if result == nil {
			return nil, ErrCacheMiss
		}
		// success
		return result, nil

	case <-ctx.Done():
		// timeout/context cancelled
		return nil, ctx.Err()

	case err := <-errorCh:
		// failure
		return nil, err
	}
}

// Set implements Storage
func (r *DynamoDbStorage) Set(ctx context.Context, key string, bytes []byte) error {
	resultCh := make(chan struct{}, 1)
	errorCh := hystrix.Go(CbDynamoDbStorage, func() error {
		defer close(resultCh)

		timestamp := time.Now().Add(r.TTL).Unix()

		params := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				ddbKey: {
					S: aws.String(key),
				},
				ddbData: {
					B: bytes,
				},
				ddbTTL: {
					N: aws.String(strconv.FormatInt(timestamp, 10)),
				},
			},
			TableName: aws.String(r.TableName),
		}

		_, err := r.Service.PutItemWithContext(ctx, params)
		return err
	}, nil)

	select {
	case <-resultCh:
		// success
		return nil

	case <-ctx.Done():
		// timeout/context cancelled
		return ctx.Err()

	case err := <-errorCh:
		// failure
		return err
	}
}

// Invalidate implements Storage
func (r *DynamoDbStorage) Invalidate(ctx context.Context, key string) error {
	resultCh := make(chan struct{}, 1)
	errorCh := hystrix.Go(CbDynamoDbStorage, func() error {
		defer close(resultCh)

		params := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				ddbKey: {
					S: aws.String(key),
				},
			},
			TableName: aws.String(r.TableName),
		}

		_, err := r.Service.DeleteItemWithContext(ctx, params)
		return err
	}, nil)

	select {
	case <-resultCh:
		// success
		return nil

	case <-ctx.Done():
		// timeout/context cancelled
		return ctx.Err()

	case err := <-errorCh:
		// failure
		return err
	}
}
