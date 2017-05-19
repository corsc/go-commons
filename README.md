# Go Commons

Find yourself needing the same packages over and over?

Me too.

Instead of re-writting them for each project/client/company I decided it was time to do it one more time in public.

## Getting Started

To use these packages simply:
```
go get github.com/corsc/go-commons/
```
(vendoring recommended)

## Packages

* [**Cache**](cache/) - A simple cache with pluggable storage (currently includes Redis and DynamoDb storage)
* [**Concurrency**](concurrency/) - Packages related to concurrency
    * [**Concurrent Map**](concurrency/cmap) - A concurrent map implementations with pluggable sharding implementations 
* [**HTTP**](http/) - Packages related to serving or consuming HTTP
    * [**Middleware**](http/middleware) - Middleware to decorate your HTTP handlers with additional features/functionality. Including
        * [Instrumentation](http/middleware/metrics_examples_test.go)
        * [JSON Request processing](http/middleware/input_body_examples_test.go)
        * [Output JSON](http/middleware/output_json.go)
        * [Panic catch and log](http/middleware/panic_examples_test.go)
        * Security Related Headers:
            * [Content-Security-Policy](http/middleware/csp_examples_test.go)
            * [HTTP Strict Transport Security](http/middleware/hsts_examples_test.go)
            * [X-Content-Type-Options](http/middleware/content_no_sniff_examples_test.go)
            * [X-XSS-Protection](http/middleware/xxss_examples_test.go)
        * Version Header
* [**I/O Closer** - a convenience function for closing and optionally logging io.Closers in 1 line (useful for defer calls)
* [**Resilience** - Packages related to resilience
    * [**Retry** - Retry with Exponential Backoff & Decorrelated Jitter Algorithm described in [here](https://www.awsarchitectureblog.com/2015/03/backoff.html)

### Prerequisites

* Go 1.8
* (optional) [GoMetaLinter](https://github.com/alecthomas/gometalinter)
* (optional) [My GoMetaLinter Config](https://raw.githubusercontent.com/corsc/PersonalTools/master/go-scripts/gometa-config.json)

## Contributions
Contributions, suggestions, bug request, etc are all welcome.  Please use Github [Issues](https://github.com/corsc/go-commons/issues) and [Pull Requests](https://github.com/corsc/go-commons/pulls).

### Running the tests

Nothing special, standard `go test ./...` will get the job done.

If me, you want the fastest possible tests, I would skip vendor by using:
```
go test $(go list ./... | grep -v /vendor)
```

### Lint checking contributions

Please check all PRs before sending them using the following settings

```
gometalinter --config=gometa-config.json ./...
```

## Authors

* **Corey Scott** - [corsc](https://github.com/corsc)

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* **Chao Gao** - [reterVision](https://github.com/reterVision) - Code Reviews and Advice
* **Ryan Cumming** - [etaoins](https://github.com/etaoins) - Code Reviews and Advice
