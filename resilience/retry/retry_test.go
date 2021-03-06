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

package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Do_happyPath(t *testing.T) {
	callsChan := make(chan struct{}, defaultMaxAttempts)
	happyLambda := func() error {
		callsChan <- struct{}{}
		return nil
	}

	// create a retry client with all the defaults
	retry := &Client{}

	resultErr := retry.Do(context.Background(), "foo", happyLambda)
	assert.Nil(t, resultErr)
	assert.True(t, len(callsChan) == 1)
}

func TestClient_Do_happyPathMetrics(t *testing.T) {
	mockMetrics := &mockMetricsClient{}
	mockMetrics.On("Incr", "foo", mock.Anything).Once()
	mockMetrics.On("Duration", "foo.do", mock.Anything, mock.Anything).Once()

	happyLambda := func() error {
		return nil
	}

	// create a retry client with defaults and mock metrics
	retry := &Client{
		Metrics: mockMetrics,
	}

	resultErr := retry.Do(context.Background(), "foo", happyLambda)
	assert.Nil(t, resultErr)
	assert.True(t, mockMetrics.AssertExpectations(t))
}

func TestClient_Do_error(t *testing.T) {
	callsChan := make(chan struct{}, defaultMaxAttempts)
	sadLambda := func() error {
		callsChan <- struct{}{}
		return errors.New("something broke")
	}

	// create a retry client with all the defaults
	retry := &Client{}

	resultErr := retry.Do(context.Background(), "foo", sadLambda)
	assert.Equal(t, ErrAttemptsExceeded, resultErr)
	assert.Equal(t, defaultMaxAttempts, len(callsChan))
}

func TestClient_Do_fatalError(t *testing.T) {
	callsChan := make(chan struct{}, defaultMaxAttempts)
	sadLambda := func() error {
		callsChan <- struct{}{}
		return errors.New("something broke")
	}

	// create a retry client with defaults and custom fatal error detection
	retry := &Client{
		CanRetry: func(err error) bool {
			return false
		},
	}

	resultErr := retry.Do(context.Background(), "foo", sadLambda)
	assert.NotNil(t, resultErr)
	assert.Equal(t, 1, len(callsChan))
}

func TestClient_Do_ignoringErrors(t *testing.T) {
	metricKey := "foo"
	knownErr := errors.New("something broke")

	callsChan := make(chan struct{}, defaultMaxAttempts)
	sadLambda := func() error {
		callsChan <- struct{}{}
		return knownErr
	}

	mockMetrics := &mockMetricsClient{}
	mockMetrics.On("Incr", metricKey, []string{"type:ignored", "attempt:0"}).Once()
	mockMetrics.On("Duration", metricKey+".do", mock.Anything, mock.Anything).Once()

	// create a retry client with defaults and custom fatal error detection
	retry := &Client{
		Metrics: mockMetrics,
		IsIgnored: func(err error) bool {
			return err == knownErr
		},
	}

	resultErr := retry.Do(context.Background(), metricKey, sadLambda)
	assert.NotNil(t, resultErr)
	assert.Equal(t, 1, len(callsChan))
}

func TestClient_Do_sleep(t *testing.T) {
	scenarios := []struct {
		desc        string
		attempt     int
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{
			desc:        "defaults + attempt 1",
			attempt:     1,
			expectedMin: time.Duration(10 * time.Millisecond),
			expectedMax: time.Duration(180 * time.Millisecond),
		},
		{
			desc:        "defaults + attempt 2",
			attempt:     1,
			expectedMin: time.Duration(10 * time.Millisecond),
			expectedMax: time.Duration(360 * time.Millisecond),
		},
		{
			desc:        "defaults + attempt 3",
			attempt:     1,
			expectedMin: time.Duration(10 * time.Millisecond),
			expectedMax: time.Duration(720 * time.Millisecond),
		},
		{
			desc:        "defaults + attempt 4 (exceeds max)",
			attempt:     1,
			expectedMin: time.Duration(10 * time.Millisecond),
			expectedMax: time.Duration(1000 * time.Millisecond),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			retry := &Client{}

			result := retry.getSleep(scenario.attempt)

			assert.True(t, int64(result) >= int64(scenario.expectedMin), scenario.desc)
			assert.True(t, int64(result) <= int64(scenario.expectedMax), scenario.desc)
		})
	}
}

func TestClient_Do_contextAlreadyDone(t *testing.T) {
	callsChan := make(chan struct{}, defaultMaxAttempts)
	sadLambda := func() error {
		callsChan <- struct{}{}
		return errors.New("something broke")
	}

	// create a retry client with defaults and custom fatal error detection
	retry := &Client{}

	// create closed context
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	// call should immediately fail without trying
	resultErr := retry.Do(ctx, "foo", sadLambda)
	assert.NotNil(t, resultErr)
	assert.Equal(t, 0, len(callsChan))
}

func TestClient_Do_contextTimeoutDuringAttempts(t *testing.T) {
	callsChan := make(chan struct{}, defaultMaxAttempts)
	sadLambda := func() error {
		callsChan <- struct{}{}
		return errors.New("something broke")
	}

	// create a retry client with silly settings
	retry := &Client{
		MaxAttempts: 100,
		BaseDelay:   1 * time.Second,
		MaxDelay:    10 * time.Second,
	}

	// create closed context
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancelFn()

	// call should timeout before we run out of attempts
	resultErr := retry.Do(ctx, "foo", sadLambda)
	assert.NotNil(t, resultErr)
	assert.True(t, len(callsChan) >= 1)
	assert.True(t, len(callsChan) < 100)
}

func TestClient_Do_contextTimeoutSlowLambda(t *testing.T) {
	callsChan := make(chan struct{}, defaultMaxAttempts)
	sadLambda := func() error {
		callsChan <- struct{}{}
		<-time.After(10 * time.Second)
		return errors.New("something broke")
	}

	// create a retry client with silly settings
	retry := &Client{}

	// create closed context
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancelFn()

	// call should fail before the first attempt is complete
	resultErr := retry.Do(ctx, "foo", sadLambda)
	assert.NotNil(t, resultErr)
	assert.True(t, len(callsChan) == 1)
}

type mockMetricsClient struct {
	mock.Mock
}

// Incr implements MetricsClient
func (m *mockMetricsClient) Incr(key string, tags ...string) {
	m.Called(key, tags)
}

// Duration implements MetricsClient
func (m *mockMetricsClient) Duration(key string, start time.Time, tags ...string) {
	m.Called(key, start, tags)
}
