# Retry

This packages attempts to provide an easy way to retry lambas.

It is implemented with the Exponential Backoff & Decorrelated Jitter Algorithm described in [here](https://www.awsarchitectureblog.com/2015/03/backoff.html).

For usage examples please refer [here](retry_examples_test.go)