package types

// Route describes a hop and cost to reach a destination.
type Route struct {
	Target  Node
	NextHop Node
	Metric  float64
}
