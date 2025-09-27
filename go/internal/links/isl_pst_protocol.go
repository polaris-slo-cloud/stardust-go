package links

import (
	"errors"
	"sort"
	"sync"

	"github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/helper"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.InterSatelliteLinkProtocol = (*IslPstProtocol)(nil)

// IslPstProtocol implements a link selection strategy inspired by partial spanning trees (PST).
// It connects satellites by preferring links with the lowest latency while maintaining a maximum
// number of links per satellite. This helps form a distributed yet connected network with minimal overhead.
type IslPstProtocol struct {
	setLink     map[*linktypes.IslLink]struct{} // All candidate links seen by the protocol
	established []*linktypes.IslLink            // Currently active links
	resultCache []types.Link                    // Cached result of last UpdateLinks

	satellite      types.Node                // The satellite this protocol is mounted to
	satellites     []types.Node              // All reachable satellites
	representative map[types.Node]types.Node // Union-find mapping for MST cycles

	position   types.Vector             // Last position we calculated for
	mu         sync.Mutex               // Protects concurrent access
	resetEvent *helper.ManualResetEvent // Notifies when UpdateLinks finishes
}

// NewIslPstProtocol initializes the protocol instance
func NewIslPstProtocol() *IslPstProtocol {
	return &IslPstProtocol{
		setLink:        make(map[*linktypes.IslLink]struct{}),
		established:    []*linktypes.IslLink{},
		representative: make(map[types.Node]types.Node),
		resetEvent:     helper.NewManualResetEvent(true),
	}
}

// Mount assigns this protocol to a given satellite
func (p *IslPstProtocol) Mount(s types.Node) {
	if p.satellite == nil {
		p.satellite = s
	}
}

// AddLink registers a new candidate link to the protocol's pool
func (p *IslPstProtocol) AddLink(link types.Link) {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.setLink[isl] = struct{}{}
	}
}

// ConnectLink adds a link to the active set if not already connected
func (p *IslPstProtocol) ConnectLink(link types.Link) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if isl, ok := link.(*linktypes.IslLink); ok {
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
func (p *IslPstProtocol) DisconnectLink(link types.Link) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if isl, ok := link.(*linktypes.IslLink); ok {
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
func (p *IslPstProtocol) ConnectSatellite(s types.Node) error {
	return errors.New("ConnectSatellite not implemented")
}

// DisconnectSatellite is a stub
func (p *IslPstProtocol) DisconnectSatellite(s types.Node) error {
	return errors.New("DisconnectSatellite not implemented")
}

// UpdateLinks recalculates which links to establish using a distributed MST-style heuristic
func (p *IslPstProtocol) UpdateLinks() ([]types.Link, error) {
	if p.satellite == nil {
		return nil, errors.New("satellite not mounted")
	}

	p.mu.Lock()
	if p.position.Equals(p.satellite.GetPosition()) {
		p.mu.Unlock()
		p.resetEvent.Wait()
		return p.resultCache, nil
	}
	p.position = p.satellite.GetPosition()
	p.resetEvent.Reset()
	p.mu.Unlock()

	satSet := make(map[types.Node]bool)
	for l := range p.setLink {
		satSet[l.Node1] = true
		satSet[l.Node2] = true
	}
	p.satellites = []types.Node{}
	for s := range satSet {
		p.satellites = append(p.satellites, s)
		p.representative[s] = s
	}

	maxLinks := 4
	nodes := map[types.Node]int{}
	mstLinks := []*linktypes.IslLink{}

	for _, sat := range p.satellites {
		links := sat.GetLinkNodeProtocol().Links()
		valid := []*linktypes.IslLink{}
		for _, l := range links {
			if isl, ok := l.(*linktypes.IslLink); ok && isl.IsReachable() {
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

	estSet := make(map[*linktypes.IslLink]bool)
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

	p.resultCache = make([]types.Link, len(mstLinks))
	for i, l := range mstLinks {
		p.resultCache[i] = l
	}
	p.resetEvent.Set()
	return p.resultCache, nil
}

// Links returns all registered candidate links
func (p *IslPstProtocol) Links() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	res := []types.Link{}
	for l := range p.setLink {
		res = append(res, l)
	}
	return res
}

// Established returns the list of currently active links
func (p *IslPstProtocol) Established() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	res := make([]types.Link, len(p.established))
	for i, l := range p.established {
		res[i] = l
	}
	return res
}

// getRepresentative returns the root representative for a satellite using path compression
func getRepresentative(reps map[types.Node]types.Node, sat types.Node) types.Node {
	cur := sat
	for reps[cur] != cur {
		cur = reps[cur]
	}
	reps[sat] = cur
	return cur
}
