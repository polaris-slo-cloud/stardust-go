package links

import (
	"sync"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.InterSatelliteLinkProtocol = (*PrecomputedLinkProtocol)(nil)

type PrecomputedLinkProtocol struct {
	node        types.Node
	links       []types.Link
	established [][]types.Link
	currentIx   int

	position types.Vector
	mu       sync.Mutex
}

func NewSimulatedLinkProtocol() *PrecomputedLinkProtocol {
	return &PrecomputedLinkProtocol{
		links:     []types.Link{},
		currentIx: -1,
		position:  types.NewVector(-1, -1, -1),
	}
}

func (p *PrecomputedLinkProtocol) InjectEstablishedLinks(links [][]types.Link) {
	p.established = links
}

func (p *PrecomputedLinkProtocol) Mount(s types.Node) {
	if p.node == nil {
		p.node = s
	}
}

func (p *PrecomputedLinkProtocol) ConnectLink(link types.Link) error {
	return nil
}

func (p *PrecomputedLinkProtocol) DisconnectLink(link types.Link) error {
	return nil
}

// UpdateLinks updates the list of active links.
func (p *PrecomputedLinkProtocol) UpdateLinks() ([]types.Link, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.position == p.node.GetPosition() {
		return p.established[p.currentIx], nil
	}

	p.currentIx++
	p.position = p.node.GetPosition()
	return p.established[p.currentIx], nil
}

// Established returns the list of established links.
func (p *PrecomputedLinkProtocol) Established() []types.Link {
	curr := p.established[p.currentIx]
	out := make([]types.Link, 0, len(curr))
	out = append(out, curr...)
	return out
}

func (p *PrecomputedLinkProtocol) Links() []types.Link {
	return p.links
}

func (p *PrecomputedLinkProtocol) AddLink(link types.Link) {
	p.links = append(p.links, link)
}
