package types

// Router represents a router capable of resolving routes to services or nodes.
type Router interface {
	CanPreRouteCalc() bool
	CanOnRouteCalc() bool
	Mount(node Node) error
	CalculateRoutingTableAsync() error
	AdvertiseNewServiceAsync(serviceName string) error
	ReceiveServiceAdvertismentsAsync(serviceName string, outlink Link, route RouteResult) error
	RouteAsyncToNode(target Node, payload Payload) (RouteResult, error)
	RouteAsync(serviceName string, payload Payload) (RouteResult, error)
}
