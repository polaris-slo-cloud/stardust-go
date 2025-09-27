package links

import (
	"sync"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.InterSatelliteLinkProtocol = (*SimulatedLinkProtocol)(nil)

type SimulatedLinkProtocol struct {
	node        types.Node
	links       []types.Link
	established [][]types.Link
	currentIx   int

	position types.Vector
	mu       sync.Mutex
}

func NewSimulatedLinkProtocol() *SimulatedLinkProtocol {
	return &SimulatedLinkProtocol{
		links:     []types.Link{},
		currentIx: -1,
		position:  types.NewVector(-1, -1, -1),
	}
}

func (p *SimulatedLinkProtocol) InjectEstablishedLinks(links [][]types.Link) {
	p.established = links
}

func (p *SimulatedLinkProtocol) Mount(s types.Node) {
	if p.node == nil {
		p.node = s
	}
}

func (p *SimulatedLinkProtocol) ConnectLink(link types.Link) error {
	return nil
}

func (p *SimulatedLinkProtocol) DisconnectLink(link types.Link) error {
	return nil
}

// UpdateLinks updates the list of active links.
func (p *SimulatedLinkProtocol) UpdateLinks() ([]types.Link, error) {
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
func (p *SimulatedLinkProtocol) Established() []types.Link {
	return p.established[p.currentIx]
}

func (p *SimulatedLinkProtocol) Links() []types.Link {
	return p.links
}

func (p *SimulatedLinkProtocol) ConnectSatellite(s types.Node) error {
	return nil
}

func (p *SimulatedLinkProtocol) DisconnectSatellite(s types.Node) error {
	return nil
}

func (p *SimulatedLinkProtocol) AddLink(link types.Link) {
	p.links = append(p.links, link)
}
