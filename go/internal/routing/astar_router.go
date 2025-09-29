// Package routing contains routing algorithms used for simulation
package routing

import (
	"errors"
	"math"
	"sort"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

// AStarRouter implements the A* pathfinding algorithm between nodes.
type AStarRouter struct {
	self  types.Node   // the node this router is mounted to
	nodes []types.Node // cached list of all reachable nodes
}

// NewAStarRouter creates a new AStarRouter instance.
func NewAStarRouter() *AStarRouter {
	return &AStarRouter{}
}

// Mount binds the router to a node. This method satisfies the IRouter interface.
func (r *AStarRouter) Mount(n types.Node) error {
	if r.self != nil {
		return errors.New("router already mounted")
	}
	r.self = n
	return nil
}

// CanPreRouteCalc returns false for A* since no precomputation is used. This method satisfies the IRouter interface.
func (r *AStarRouter) CanPreRouteCalc() bool { return false }

// CanOnRouteCalc returns true for A* which calculates route on demand. This method satisfies the IRouter interface.
func (r *AStarRouter) CanOnRouteCalc() bool { return true }

// CalculateRoutingTableAsync is not applicable to A* as it is a reactive algorithm. This method is a placeholder to satisfy the IRouter interface.
func (r *AStarRouter) CalculateRoutingTableAsync() error {
	// A* doesn't need a pre-calculated table, so this method is effectively a no-op.
	return nil
}

// AdvertiseNewServiceAsync advertises a new service. A* doesn't use this, but it's implemented to satisfy the IRouter interface.
func (r *AStarRouter) AdvertiseNewServiceAsync(serviceName string) error {
	// Placeholder for service advertisement logic (future use)
	return nil
}

// ReceiveServiceAdvertismentsAsync updates the service routes. A* doesn't use this, but it's implemented to satisfy the IRouter interface.
func (r *AStarRouter) ReceiveServiceAdvertismentsAsync(serviceName string, outlink types.Link, route types.RouteResult) error {
	// Placeholder for receiving service advertisement (future use)
	return nil
}

// RouteAsync finds the nearest node that hosts the service and routes to it.
func (r *AStarRouter) RouteAsync(serviceName string, payload types.Payload) (types.RouteResult, error) {
	if r.self == nil {
		return nil, errors.New("router not mounted")
	}

	var candidates []types.Node
	for _, n := range r.getNeighbourhood() {
		if n.GetComputing().HostsService(serviceName) {
			candidates = append(candidates, n)
		}
	}

	if len(candidates) == 0 {
		return UnreachableRouteResultInstance, nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return r.self.DistanceTo(candidates[i]) < r.self.DistanceTo(candidates[j])
	})
	return r.RouteTo(candidates[0], payload)
}

// RouteAsyncToNode is used to route to a specific node. This method satisfies the IRouter interface.
func (r *AStarRouter) RouteAsyncToNode(target types.Node, payload types.Payload) (types.RouteResult, error) {
	if r.self == nil {
		return nil, errors.New("router not mounted")
	}
	return r.RouteTo(target, payload)
}

// RouteTo executes A* from the mounted node to the given target.
func (r *AStarRouter) RouteTo(target types.Node, payload types.Payload) (types.RouteResult, error) {
	if r.self == nil {
		return nil, errors.New("router not mounted")
	}

	openset := make(map[types.Node]float64)
	gScore := map[types.Node]float64{r.self: 0}
	fScore := map[types.Node]float64{r.self: heuristic(r.self, target)}
	openset[r.self] = fScore[r.self]

	for len(openset) > 0 {
		// Find node in openset with lowest fScore
		var current types.Node
		minScore := math.MaxFloat64
		for n, score := range openset {
			if score < minScore {
				current = n
				minScore = score
			}
		}
		delete(openset, current)

		if current == target {
			return NewOnRouteResult(int(gScore[current]), 0), nil
		}

		for _, l := range current.GetLinkNodeProtocol().Established() {
			neighbor := l.GetOther(current)
			alt := gScore[current] + l.Latency()
			if prev, ok := gScore[neighbor]; !ok || alt < prev {
				gScore[neighbor] = alt
				fScore[neighbor] = alt + heuristic(neighbor, target)
				openset[neighbor] = fScore[neighbor]
			}
		}
	}
	return UnreachableRouteResultInstance, nil
}

// getNeighbourhood returns all nodes connected to the current node (BFS).
func (r *AStarRouter) getNeighbourhood() []types.Node {
	visited := map[types.Node]bool{}
	queue := []types.Node{r.self}
	var result []types.Node

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if visited[n] {
			continue
		}
		visited[n] = true
		result = append(result, n)
		for _, l := range n.GetLinkNodeProtocol().Established() {
			other := l.GetOther(n)
			if !visited[other] {
				queue = append(queue, other)
			}
		}
	}
	return result
}

// heuristic estimates the distance from node a to node b (in ms).
func heuristic(a, b types.Node) float64 {
	d := a.DistanceTo(b)
	return d / configs.SpeedOfLight * 1000
}
