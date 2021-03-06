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
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	defaultMaxAttempts = 3
	defaultBaseDelay   = 10 * time.Millisecond
	defaultMaxDelay    = 1 * time.Second
)

// ErrAttemptsExceeded is returned when we exceeded the max attempts without succeeding
var ErrAttemptsExceeded = errors.New("exceeded max attempts")

// MetricsClient defines the metrics interface used by this package
type MetricsClient interface {
	Incr(key string, tags ...string)
	Duration(key string, start time.Time, tags ...string)
}

// Client implements the Exponential Backoff
type Client struct {
	// MaxAttempts is the maximum number of retry attempts before giving up. (default: 3)
	MaxAttempts int

	// BaseDelay is the base amount of time between attempts (default: 10 ms)
	BaseDelay time.Duration

	// MaxDelay is the maximum possible delay (default: 1 second)
	MaxDelay time.Duration

	// CanRetry allows for selectively shortcutting the retries.
	// Useful for cases were retrying would never work. (default: retry always)
	CanRetry func(error) bool

	// IsIgnored allows for selectively ignoring some errors (from tracking) but still returning the error.
	// Useful for things like ignoring a HTTP404 when that is not an error
	IsIgnored func(error) bool

	// MetricsClient allows this package to emit metrics (default: no metrics)
	Metrics MetricsClient
}

// Do executes the lambda until success, context is cancelled, attempts are exceeded or a fatal error.
func (r *Client) Do(ctx context.Context, metricKey string, do func() error) error {
	for attempt := 0; attempt < r.getMaxAttempts(); attempt++ {
		if ctx.Err() != nil {
			return r.returnContextError(ctx, metricKey)
		}

		// wrap lamba so we can "quit early" with using context
		doChan := make(chan error, 1)
		go func() {
			defer close(doChan)
			defer r.trackDoDuration(metricKey, time.Now(), attempt)

			err := do()
			if err != nil {
				doChan <- err
			}
		}()

		// wait for lambda or context to complete
		select {
		case err := <-doChan:
			if err == nil {
				r.getMetrics().Incr(metricKey, "type:success", r.buildAttemptTag(attempt))
				return nil
			}

			isFatal := r.handleError(metricKey, attempt, err)
			if isFatal {
				return err
			}

		case <-ctx.Done():
			return r.returnContextError(ctx, metricKey)
		}

		// sleep before trying again
		sleep := r.getSleep(attempt)
		select {
		case <-time.After(sleep):
			// nothing

		case <-ctx.Done():
			r.getMetrics().Incr(metricKey, "type:error", "cause:context", r.buildAttemptTag(attempt))
			return ctx.Err()
		}
	}

	// give up
	r.getMetrics().Incr(metricKey, "type:error", "cause:attempts", r.buildAttemptTag(r.getMaxAttempts()))
	return ErrAttemptsExceeded
}

func (r *Client) handleError(metricKey string, attempt int, err error) bool {
	if r.isIgnored(err) {
		r.getMetrics().Incr(metricKey, "type:ignored", r.buildAttemptTag(attempt))
		return true
	}

	if !r.canRetry(err) {
		r.getMetrics().Incr(metricKey, "type:error", "cause:fatal", r.buildAttemptTag(attempt))
		return true
	}

	r.getMetrics().Incr(metricKey, "type:error", "cause:lambda", r.buildAttemptTag(attempt))
	return false
}

func (r *Client) trackDoDuration(metricKey string, start time.Time, attempt int) {
	r.getMetrics().Duration(metricKey+".do", start, r.buildAttemptTag(attempt))
}

func (r *Client) getSleep(attempt int) time.Duration {
	maxDelayFloat := float64(r.getMaxDelay())

	delayByAttempt := float64(r.getBaseDelay()) * 2 * math.Exp2(float64(attempt))
	temp := int64(math.Min(maxDelayFloat, delayByAttempt))
	sleep := temp/2 + rand.Int63n(temp/2)

	randComp := float64(rand.Int63n((sleep*3)-r.getBaseDelay()) + r.getBaseDelay())
	sleep = int64(math.Min(maxDelayFloat, randComp))
	return time.Duration(sleep)
}

func (r *Client) getMaxAttempts() int {
	if r.MaxAttempts > 0 {
		return r.MaxAttempts
	}

	return defaultMaxAttempts
}

func (r *Client) getBaseDelay() int64 {
	if int64(r.BaseDelay) > int64(0) {
		return int64(r.BaseDelay)
	}

	return int64(defaultBaseDelay)
}

func (r *Client) getMaxDelay() int64 {
	if int64(r.MaxDelay) > int64(0) {
		return int64(r.MaxDelay)
	}

	return int64(defaultMaxDelay)
}

func (r *Client) canRetry(err error) bool {
	if r.CanRetry != nil {
		return r.CanRetry(err)
	}

	return true
}

func (r *Client) isIgnored(err error) bool {
	if r.IsIgnored != nil {
		return r.IsIgnored(err)
	}

	return false
}

func (r *Client) getMetrics() MetricsClient {
	if r.Metrics != nil {
		return r.Metrics
	}
	return &noopMetrics{}
}

func (r *Client) buildAttemptTag(attempt int) string {
	return "attempt:" + strconv.Itoa(attempt)
}

func (r *Client) returnContextError(ctx context.Context, metricKey string) error {
	r.getMetrics().Incr(metricKey, "type:error", "cause:context")
	return ctx.Err()
}

type noopMetrics struct{}

func (n *noopMetrics) Incr(key string, tags ...string) {
	// intentionally does nothing
}

func (n *noopMetrics) Duration(key string, start time.Time, tags ...string) {
	// intentionally does nothing
}
