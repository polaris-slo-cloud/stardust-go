package types

// RouteResult represents the result of a routing calculation.
type RouteResult interface {
	Reachable() bool
	Latency() int
	AddCalculationDuration(duration int) RouteResult
	WaitLatencyAsync() error
}
