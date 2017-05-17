package cache

// Metrics allows for instrumenting the cache for hit/miss and errors
//go:generate mockery -name Metrics -inpkg -testonly -case underscore
type Metrics interface {
	// Tracks
	Track(event Event)
}

// MetricsFunc implements Metrics as a function
type MetricsFunc func(event Event)

// Track implements Metrics
func (m MetricsFunc) Track(event Event) {
	m(event)
}

// No op implementation of metrics
var noopMetrics = MetricsFunc(func(event Event) {
	// intentionally do nothing
})
