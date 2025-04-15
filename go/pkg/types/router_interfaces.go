package types

// IPayload represents a simulation payload interface.
type IPayload interface{}

// IRouteResult represents the result of a routing calculation.
type IRouteResult interface {
	Reachable() bool
	Latency() int
	AddCalculationDuration(duration int) IRouteResult
	WaitLatencyAsync() error
}

// IRouter represents a router capable of resolving routes to services or nodes.
type IRouter interface {
	CanPreRouteCalc() bool
	CanOnRouteCalc() bool
	Mount(node INode) error
	CalculateRoutingTableAsync() error
	AdvertiseNewServiceAsync(serviceName string) error
	ReceiveServiceAdvertismentsAsync(serviceName string, outlink ILink, route IRouteResult) error
	RouteAsyncToNode(target INode, payload IPayload) (IRouteResult, error)
	RouteAsync(serviceName string, payload IPayload) (IRouteResult, error)
}

// Route describes a hop and cost to reach a destination.
type Route struct {
	Target  INode
	NextHop INode
	Metric  float64
}

// RouteAdvertisment packages an outbound link with a list of routes for advertisement.
type RouteAdvertisment struct {
	Link   ILink
	Routes []Route
}
