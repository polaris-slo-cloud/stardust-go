package types

// Router represents a router capable of resolving routes to services or nodes.
type Router interface {
	// CanPreRouteCalc indicates if the protocol can pre-calculate a routing table
	CanPreRouteCalc() bool

	// CanOnRouteCalc indicates if the protocol can calculate routes on demand
	CanOnRouteCalc() bool

	// Mount associates this router with the given node
	Mount(node Node) error

	// CalculateRoutingTable calculates the routing table
	CalculateRoutingTable() error

	// AdvertiseNewServiceAsync sends a service advertisement to other routers
	AdvertiseNewServiceAsync(serviceName string) error

	// ReceiveServiceAdvertismentsAsync receives the service advertisements to
	// i.e. fill into routing table
	ReceiveServiceAdvertismentsAsync(serviceName string, outlink Link, route RouteResult) error

	// RouteToNode returns a route from mounted node to target node
	// i.e. read from routing table or calculate on demand
	RouteToNode(target Node, payload Payload) (RouteResult, error)

	// RouteToService returns a route from mounted node to target service
	RouteToService(serviceName string, payload Payload) (RouteResult, error)
}
