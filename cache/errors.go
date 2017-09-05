package cache

// LambdaError is the error type returned when the user fetch/build lambda failed
type LambdaError struct {
	Cause error
}

// Error implements error
func (e LambdaError) Error() string {
	return e.Cause.Error()
}
