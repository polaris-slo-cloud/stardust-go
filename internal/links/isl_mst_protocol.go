package links

import (
	"errors"
	"sort"
	"sync"

	"stardustGo/configs"
	"stardustGo/internal/links/linktypes"
	"stardustGo/pkg/types"
)

// IslMstProtocol builds a global minimum spanning tree (MST) of ISL links.
// It uses Kruskal’s algorithm with a union-find structure over node names.
type IslMstProtocol struct {
	setLink         map[*linktypes.IslLink]bool // All candidate links
	established     []*linktypes.IslLink        // Currently active MST links
	satellite       types.INode                 // Local satellite
	satellites      []types.INode               // All satellites involved
	representatives map[string]string           // Disjoint-set forest (by node name)

	position types.Vector // Last position when links were updated
	mu       sync.Mutex   // Protects concurrent access
	readyCh  chan struct{}
}

// NewIslMstProtocol initializes an empty protocol instance.
func NewIslMstProtocol() *IslMstProtocol {
	return &IslMstProtocol{
		setLink:         make(map[*linktypes.IslLink]bool),
		established:     []*linktypes.IslLink{},
		representatives: make(map[string]string),
		readyCh:         make(chan struct{}, 1),
	}
}

// Mount assigns this protocol to a local satellite.
func (p *IslMstProtocol) Mount(s types.INode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.satellite == nil {
		p.satellite = s
	}
}

// AddLink registers a new candidate link.
func (p *IslMstProtocol) AddLink(link types.ILink) {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.setLink[isl] = true
	}
}

// ConnectLink adds a link to the established set if not already present.
func (p *IslMstProtocol) ConnectLink(link types.ILink) error {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		for _, l := range p.established {
			if l == isl {
				return nil
			}
		}
		p.established = append(p.established, isl)
	}
	return nil
}

// DisconnectLink removes a link from the established set.
func (p *IslMstProtocol) DisconnectLink(link types.ILink) error {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		for i, l := range p.established {
			if l == isl {
				p.established = append(p.established[:i], p.established[i+1:]...)
				break
			}
		}
	}
	return nil
}

// UpdateLinks builds a new MST if the satellite's position has changed.
func (p *IslMstProtocol) UpdateLinks() ([]types.ILink, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.satellite == nil {
		return nil, errors.New("satellite not mounted")
	}

	// Return cached if position hasn't changed
	if p.position.Equals(p.satellite.PositionVector()) {
		select {
		case <-p.readyCh:
		default:
		}
		result := make([]types.ILink, len(p.established))
		for i, l := range p.established {
			result[i] = l
		}
		return result, nil
	}
	p.position = p.satellite.PositionVector()

	// Collect all satellites from links
	satMap := map[string]types.INode{}
	for l := range p.setLink {
		satMap[l.Node1.GetName()] = l.Node1
		satMap[l.Node2.GetName()] = l.Node2
	}
	p.satellites = make([]types.INode, 0, len(satMap))
	for _, sat := range satMap {
		p.satellites = append(p.satellites, sat)
		p.representatives[sat.GetName()] = sat.GetName()
	}

	// Build list of eligible links
	type edge struct {
		dist float64
		link *linktypes.IslLink
	}
	var edges []edge
	for l := range p.setLink {
		if l.Distance() <= configs.MaxISLDistance {
			edges = append(edges, edge{l.Distance(), l})
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].dist < edges[j].dist
	})

	// Apply Kruskal's algorithm
	mst := []*linktypes.IslLink{}
	for _, e := range edges {
		n1 := e.link.Node1.GetName()
		n2 := e.link.Node2.GetName()

		rep1 := p.getRepresentative(n1)
		rep2 := p.getRepresentative(n2)

		if rep1 == rep2 || !e.link.IsReachable() {
			continue
		}

		mst = append(mst, e.link)
		p.representatives[rep2] = rep1

		if len(mst) == len(p.satellites)-1 {
			break
		}
	}

	// Update established set
	estSet := make(map[*linktypes.IslLink]bool)
	for _, l := range mst {
		l.SetEstablished(true)
		estSet[l] = true
	}
	for _, l := range p.established {
		if !estSet[l] {
			l.SetEstablished(false)
		}
	}
	p.established = mst

	// Signal readiness for reuse
	select {
	case p.readyCh <- struct{}{}:
	default:
	}

	result := make([]types.ILink, len(mst))
	for i, l := range mst {
		result[i] = l
	}
	return result, nil
}

// getRepresentative finds the disjoint-set root for a node.
func (p *IslMstProtocol) getRepresentative(name string) string {
	for p.representatives[name] != name {
		name = p.representatives[name]
	}
	return name
}

// ConnectSatellite is not implemented in this strategy.
func (p *IslMstProtocol) ConnectSatellite(types.INode) error {
	return errors.New("ConnectSatellite not implemented")
}

// DisconnectSatellite is not implemented in this strategy.
func (p *IslMstProtocol) DisconnectSatellite(types.INode) error {
	return errors.New("DisconnectSatellite not implemented")
}

// Links returns a snapshot of all known candidate links.
func (p *IslMstProtocol) Links() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make([]types.ILink, 0, len(p.setLink))
	for l := range p.setLink {
		result = append(result, l)
	}
	return result
}

// Established returns all currently active MST links.
func (p *IslMstProtocol) Established() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make([]types.ILink, len(p.established))
	for i, l := range p.established {
		result[i] = l
	}
	return result
}
