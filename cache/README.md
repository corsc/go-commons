# Cache


## Notes:

### Logging
While this package does log, it only logs asynchronous errors.  Personally I would prefer not to log at all but this 
would leave async issues completely invisible.

That said, logging is optional.

### Metrics
Metrics are provided but optional.

### Running tests with different storage
To run tests on `RedisStorage` use `go test -tags redis ./...`
(These tests assume redis running on localhost:6379)