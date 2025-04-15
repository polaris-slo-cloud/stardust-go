package links

import (
	"errors"
	linkmod "github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/types"
	"sort"
	"sync"
)

// IslPstProtocol implements a link selection strategy inspired by partial spanning trees (PST).
// It connects satellites by preferring links with the lowest latency while maintaining a maximum
// number of links per satellite. This helps form a distributed yet connected network with minimal overhead.
type IslPstProtocol struct {
	setLink     map[*linkmod.IslLink]struct{} // All candidate links seen by the protocol
	established []*linkmod.IslLink            // Currently active links

	satellite      types.INode                 // The satellite this protocol is mounted to
	satellites     []types.NodeWithISL         // All reachable satellites
	representative map[types.INode]types.INode // Union-find mapping for MST cycles

	position types.Vector  // Last position we calculated for
	mu       sync.Mutex    // Protects concurrent access
	readyCh  chan struct{} // Notifies when UpdateLinks finishes (mimics ManualResetEvent)
}

// NewIslPstProtocol initializes the protocol instance
func NewIslPstProtocol() *IslPstProtocol {
	return &IslPstProtocol{
		setLink:        make(map[*linkmod.IslLink]struct{}),
		established:    []*linkmod.IslLink{},
		representative: make(map[types.INode]types.INode),
		readyCh:        make(chan struct{}, 1),
	}
}

// Mount assigns this protocol to a given satellite
func (p *IslPstProtocol) Mount(s types.INode) {
	if p.satellite == nil {
		p.satellite = s
	}
}

// AddLink registers a new candidate link to the protocol's pool
func (p *IslPstProtocol) AddLink(link types.ILink) {
	if isl, ok := link.(*linkmod.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.setLink[isl] = struct{}{}
	}
}

// ConnectLink adds a link to the active set if not already connected
func (p *IslPstProtocol) ConnectLink(link types.ILink) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if isl, ok := link.(*linkmod.IslLink); ok {
		for _, l := range p.established {
			if l == isl {
				return nil
			}
		}
		p.established = append(p.established, isl)
	}
	return nil
}

// DisconnectLink removes a link from the active set
func (p *IslPstProtocol) DisconnectLink(link types.ILink) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if isl, ok := link.(*linkmod.IslLink); ok {
		for i, l := range p.established {
			if l == isl {
				p.established = append(p.established[:i], p.established[i+1:]...)
				break
			}
		}
	}
	return nil
}

// ConnectSatellite is a stub
func (p *IslPstProtocol) ConnectSatellite(s types.INode) error {
	return errors.New("ConnectSatellite not implemented")
}

// DisconnectSatellite is a stub
func (p *IslPstProtocol) DisconnectSatellite(s types.INode) error {
	return errors.New("DisconnectSatellite not implemented")
}

// UpdateLinks recalculates which links to establish using a distributed MST-style heuristic
func (p *IslPstProtocol) UpdateLinks() ([]types.ILink, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.satellite == nil {
		return nil, errors.New("satellite not mounted")
	}

	if p.position.Equals(p.satellite.PositionVector()) {
		select {
		case <-p.readyCh:
		default:
		}
		links := make([]types.ILink, len(p.established))
		for i, l := range p.established {
			links[i] = l
		}
		return links, nil
	}
	p.position = p.satellite.PositionVector()

	satSet := make(map[types.INode]bool)
	for l := range p.setLink {
		satSet[l.Node1] = true
		satSet[l.Node2] = true
	}
	p.satellites = []types.NodeWithISL{}
	for s := range satSet {
		if node, ok := s.(types.NodeWithISL); ok {
			p.satellites = append(p.satellites, node)
			p.representative[node] = node
		}
	}

	maxLinks := 4
	nodes := map[types.INode]int{}
	mstLinks := []*linkmod.IslLink{}

	for _, sat := range p.satellites {
		links := sat.InterSatelliteLinkProtocol().Links()
		valid := []*linkmod.IslLink{}
		for _, l := range links {
			if isl, ok := l.(*linkmod.IslLink); ok && isl.IsReachable() {
				valid = append(valid, isl)
			}
		}
		sort.Slice(valid, func(i, j int) bool {
			return valid[i].Latency() < valid[j].Latency()
		})

		for _, link := range valid {
			other := link.GetOther(sat)
			rep1 := getRepresentative(p.representative, sat)
			rep2 := getRepresentative(p.representative, other)
			if rep1 == rep2 || nodes[sat] >= maxLinks || nodes[other] >= maxLinks {
				continue
			}
			nodes[sat]++
			nodes[other]++
			p.representative[rep2] = rep1
			p.representative[other] = rep1
			p.representative[sat] = rep1
			mstLinks = append(mstLinks, link)
			break
		}
	}

	estSet := make(map[*linkmod.IslLink]bool)
	for _, l := range mstLinks {
		l.SetEstablished(true)
		estSet[l] = true
	}
	for _, l := range p.established {
		if !estSet[l] {
			l.SetEstablished(false)
		}
	}
	p.established = mstLinks

	select {
	case p.readyCh <- struct{}{}:
	default:
	}

	links := make([]types.ILink, len(mstLinks))
	for i, l := range mstLinks {
		links[i] = l
	}
	return links, nil
}

// Links returns all registered candidate links
func (p *IslPstProtocol) Links() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	res := []types.ILink{}
	for l := range p.setLink {
		res = append(res, l)
	}
	return res
}

// Established returns the list of currently active links
func (p *IslPstProtocol) Established() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	res := make([]types.ILink, len(p.established))
	for i, l := range p.established {
		res[i] = l
	}
	return res
}

// getRepresentative returns the root representative for a satellite using path compression
func getRepresentative(reps map[types.INode]types.INode, sat types.INode) types.INode {
	cur := sat
	for reps[cur] != cur {
		cur = reps[cur]
	}
	reps[sat] = cur
	return cur
}
