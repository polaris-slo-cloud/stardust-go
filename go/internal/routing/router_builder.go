package routing

import (
	"fmt"
	"strings"

	"github.com/keniack/stardustGo/pkg/types"

	"github.com/keniack/stardustGo/configs"
)

// RouterBuilder constructs routers based on configuration.
type RouterBuilder struct {
	Config configs.RouterConfig
}

// Supported routing strategies
const (
	Dijkstra = "dijkstra"
	AStar    = "a-star"
)

// NewRouterBuilder creates a new builder using the provided config.
func NewRouterBuilder(cfg configs.RouterConfig) *RouterBuilder {
	return &RouterBuilder{Config: cfg}
}

// Build creates an IRouter implementation based on config.
func (b *RouterBuilder) Build() (types.Router, error) {
	switch strings.ToLower(b.Config.Protocol) {
	case Dijkstra:
		return NewDijkstraRouter(), nil
	case AStar:
		return NewAStarRouter(), nil
	default:
		return nil, fmt.Errorf("unknown routing protocol: %s", b.Config.Protocol)
	}
}
