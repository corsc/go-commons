# Middleware

A collections of mux / router agnostic middlewares.

All middleware takes and returns `http.Handler`; you can of course use `http.HandlerFunc` instead by casting your
`http.Handler` with `http.HandlerFunc(myHandler)`.

For usage examples please refer to:

* [Instrumentation](metrics_examples_test.go)
* [JSON Request processing](input_body_examples_test.go)
* [Output JSON](output_json.go)
* [Panic catch and log](panic_examples_test.go)
* [Response Cache](response_cache_examples_test.go)
* Security Related Headers:
    * [Content-Security-Policy](csp_examples_test.go)
    * [HTTP Strict Transport Security](hsts_examples_test.go)
    * [X-Content-Type-Options](content_no_sniff_examples_test.go)
    * [X-XSS-Protection](xxss_examples_test.go)
* [Version Header](version_examples_test.go)
