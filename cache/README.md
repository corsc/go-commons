# Cache

Simple cache implementation with pluggable storage; optional logging and instrumentation.

For usage examples please refer [here](cache_examples_test.go)

## Redis storage
* This library makes no effort to ensure it does not overwrite other data in the server.  Key names should be chosen carefully

### Tests

Tests that require Redis or DDB are protected with environment variable flags.

Use `REDIS=true go test ./...` to run Redis tests.
(These tests assume redis running on `:6379`)

Use `DDB=true go test ./...` to run DDB tests.
(These tests assume redis running on `:6379`)

(These tests use [DynamoDbLocal](http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) 
and assume server running on `https://s3-ap-southeast-1.amazonaws.com/dynamodb-local-singapore/release`)

## DynamoDB storage
* TTL should be enabled on the table with attribute name `ttl` see [reference](http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/time-to-live-ttl-how-to.html)

## Notes:

### Logging
While this package does log, it only logs asynchronous errors.  Personally I would prefer not to log at all but this 
would leave async issues completely invisible.

That said, logging is optional.

### Metrics
Metrics are provided but optional.
