package routing

import (
	"errors"
	"sort"
	"stardustGo/internal/links/linktypes"
	"stardustGo/pkg/types"
)

// DijkstraRouter implements shortest-path routing using Dijkstra's algorithm
// and supports precomputed routing tables.
type DijkstraRouter struct {
	self     types.INode
	routes   map[types.INode]routeEntry
	services map[string]routeEntry
	comparer func(a, b dijkstraEntry) bool
}

type routeEntry struct {
	OutLink types.ILink
	Route   types.IRouteResult
}

type dijkstraEntry struct {
	Link    types.ILink
	Target  types.INode
	Via     types.ILink
	Latency float64
}

// NewDijkstraRouter creates a new Dijkstra-based router
func NewDijkstraRouter() *DijkstraRouter {
	return &DijkstraRouter{
		routes:   make(map[types.INode]routeEntry),
		services: make(map[string]routeEntry),
		comparer: func(a, b dijkstraEntry) bool {
			return a.Latency < b.Latency
		},
	}
}

// Mount attaches the router to a node
func (r *DijkstraRouter) Mount(n types.INode) error {
	if r.self != nil {
		return errors.New("router already mounted")
	}
	r.self = n
	return nil
}

// CanPreRouteCalc returns true (Dijkstra supports pre-calculation)
func (r *DijkstraRouter) CanPreRouteCalc() bool { return true }

// CanOnRouteCalc returns true (also usable live)
func (r *DijkstraRouter) CanOnRouteCalc() bool { return true }

// RouteAsyncToNode finds a route to a specific node
func (r *DijkstraRouter) RouteAsyncToNode(target types.INode, payload types.IPayload) (types.IRouteResult, error) {
	if r.self == nil {
		return nil, errors.New("router not mounted")
	}
	if r.self == target { // Compare values directly since self is of type types.INode
		// Return a PreRouteResult with 0 latency when self == target
		return NewPreRouteResult(0), nil // Use the function to create the PreRouteResult
	}
	if entry, ok := r.routes[target]; ok {
		return entry.Route, nil
	}
	return UnreachableRouteResultInstance, nil
}

// RouteAsync finds a route by service name
func (r *DijkstraRouter) RouteAsync(serviceName string, payload types.IPayload) (types.IRouteResult, error) {
	if r.self == nil {
		return nil, errors.New("router not mounted")
	}

	// Check if the service is hosted on this node's computing
	if r.self.GetComputing().HostsService(serviceName) {
		// Create a PreRouteResult with 0 latency if the service is hosted on this node
		return NewPreRouteResult(0), nil // Use NewPreRouteResult for flexibility
	}

	// If the service exists in the routing table, return the associated route
	if entry, ok := r.services[serviceName]; ok {
		return entry.Route, nil
	}

	// If the service is not reachable, return the UnreachableRouteResultInstance
	return UnreachableRouteResultInstance, nil
}

// CalculateRoutingTableAsync populates all shortest paths using Dijkstra
func (r *DijkstraRouter) CalculateRoutingTableAsync() error {
	if r.self == nil {
		return errors.New("router not mounted")
	}
	r.routes = make(map[types.INode]routeEntry)
	r.services = make(map[string]routeEntry)

	queue := []dijkstraEntry{}
	r.routes[r.self] = routeEntry{}

	// Initialize priority queue with links
	for _, l := range r.self.GetLinks() {
		// We assume that links are of type *linktypes.IslLink
		if islLink, ok := l.(*linktypes.IslLink); ok && islLink.Established() {
			// Only add established ISL links
			queue = append(queue, dijkstraEntry{
				Link:    l,
				Target:  l.GetOther(r.self),
				Via:     l,
				Latency: islLink.Latency(),
			})
		}
	}

	// Sort the queue based on latency
	sort.Slice(queue, func(i, j int) bool { return r.comparer(queue[i], queue[j]) })

	// Initialize visited map
	visited := map[types.INode]bool{r.self: true}

	// Process the queue
	for len(queue) > 0 {
		// Get the entry with the least latency
		entry := queue[0]
		queue = queue[1:] // Pop the first entry from the queue

		// Skip already visited targets
		if visited[entry.Target] {
			continue
		}

		// Mark the target as visited and add it to the routes
		visited[entry.Target] = true
		r.routes[entry.Target] = routeEntry{
			OutLink: entry.Via,
			Route:   NewPreRouteResult(int(entry.Latency)),
		}

		// Handle services for the target node
		r.addServicesToRoutes(entry.Target, entry.Latency)

		// Add the neighbors to the queue
		for _, link := range entry.Target.GetLinks() {
			// Only process ISL links
			if islLink, ok := link.(*linktypes.IslLink); ok {
				// Add the neighboring ISL link to the queue
				neighbor := link.GetOther(entry.Target)
				if !visited[neighbor] {
					queue = append(queue, dijkstraEntry{
						Link:    islLink,
						Target:  neighbor,
						Via:     entry.Via,
						Latency: entry.Latency + islLink.Latency(),
					})
				}
			}
		}

		// Re-sort the queue by latency
		sort.Slice(queue, func(i, j int) bool { return r.comparer(queue[i], queue[j]) })
	}

	return nil
}

// AdvertiseNewServiceAsync pushes a service to neighbors (future use)
func (r *DijkstraRouter) AdvertiseNewServiceAsync(serviceName string) error {
	return nil // placeholder for broadcast mechanism
}

// ReceiveServiceAdvertismentsAsync updates service routes
func (r *DijkstraRouter) ReceiveServiceAdvertismentsAsync(serviceName string, outlink types.ILink, route types.IRouteResult) error {
	if existing, ok := r.services[serviceName]; ok && existing.Route.Latency() <= route.Latency() {
		return nil
	}
	r.services[serviceName] = routeEntry{outlink, route}
	return nil
}

// addServicesToRoutes helps manage the services associated with a node in the routes map.
func (r *DijkstraRouter) addServicesToRoutes(target types.INode, latency float64) {
	// Loop through all services hosted on the target node
	for _, service := range target.GetComputing().GetServices() {
		// Check if the service already exists in the routing table
		if _, exists := r.services[service.GetServiceName()]; !exists {
			// Handle the case where there are no links available
			if len(target.GetLinks()) > 0 {
				// Use the first available link as the "via" link for simplicity
				r.services[service.GetServiceName()] = routeEntry{
					OutLink: target.GetLinks()[0],            // Using the first link as the "via"
					Route:   NewPreRouteResult(int(latency)), // Creating PreRouteResult with latency
				}
			} else {
				// If no links are available, we can set the route to unreachable or handle it differently
				r.services[service.GetServiceName()] = routeEntry{
					OutLink: nil,                            // No valid link
					Route:   UnreachableRouteResultInstance, // Set to unreachable
				}
			}
		}
	}
}
