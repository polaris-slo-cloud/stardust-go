package routing

import (
	"stardustGo/pkg/types"
)

type PreRouteResult struct {
	latency int
}

// NewPreRouteResult creates a new PreRouteResult with the specified latency
func NewPreRouteResult(latency int) types.IRouteResult {
	if latency < 0 {
		return nil // Return nil if latency is negative
	}
	return &PreRouteResult{latency: latency}
}

// Reachable returns whether the route is reachable (PreRoute is always reachable)
func (r *PreRouteResult) Reachable() bool {
	return true
}

// Latency returns the calculated latency for the route
func (r *PreRouteResult) Latency() int {
	return r.latency
}

// WaitLatencyAsync simulates waiting for the latency (asynchronous operation)
func (r *PreRouteResult) WaitLatencyAsync() error {
	return delayMilliseconds(r.latency)
}

// AddCalculationDuration adds additional calculation duration to the route and returns the updated result
func (r *PreRouteResult) AddCalculationDuration(calculationDuration int) types.IRouteResult {
	return NewOnRouteResult(r.latency, calculationDuration)
}
