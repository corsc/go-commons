# Middleware

A collections of mux / router agnostic middlewares.

All middleware takes and returns `http.Handler`; you can of course use `http.HandlerFunc` instead by casting your
`http.Handler` with `http.HandlerFunc(myHandler)`.

## Ideas / TOOD:
* Conversion from HTTP request to struct
* JWT / 2 legged oauth2
* Response Content Type (by status code)
* Debug/Log Request (2 modes: all request and only errors)
* Request ID
* Response Cache
* Validate Request Content Type
* Panic Handling
* CORS
* Request Timeout
