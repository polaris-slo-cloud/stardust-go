package routing

import (
	"sync"
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

type OnRouteResult struct {
	latency             int
	calculationDuration int
	firstRequest        bool
	lock                sync.Mutex
}

func NewOnRouteResult(latency, calculationDuration int) *OnRouteResult {
	return &OnRouteResult{
		latency:             latency,
		calculationDuration: calculationDuration,
		firstRequest:        true,
	}
}

func (r *OnRouteResult) Reachable() bool {
	return true
}

func (r *OnRouteResult) Latency() int {
	return r.latency
}

func (r *OnRouteResult) WaitLatencyAsync() error {
	wait := r.latency
	r.lock.Lock()
	if r.firstRequest {
		wait -= r.calculationDuration
		r.firstRequest = false
	}
	r.lock.Unlock()

	if wait > 0 {
		return delayMilliseconds(wait)
	}
	return nil
}

func (r *OnRouteResult) AddCalculationDuration(calculationDuration int) types.RouteResult {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.firstRequest {
		r.calculationDuration += calculationDuration
	} else {
		r.calculationDuration = calculationDuration
		r.firstRequest = true
	}
	return r
}

func delayMilliseconds(ms int) error {
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}
