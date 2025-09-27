package links

import (
	"errors"
	"log"
	"sync"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.InterSatelliteLinkProtocol = (*LinkFilterProtocol)(nil)

// LinkFilterProtocol wraps another IInterSatelliteLinkProtocol and filters links
// so that only those involving the mounted node are retained and processed.
type LinkFilterProtocol struct {
	inner       types.InterSatelliteLinkProtocol
	node        types.Node
	links       map[types.Link]struct{}
	established map[types.Link]struct{}
	out         map[types.Link]struct{}
	mu          sync.Mutex
}

// NewLinkFilterProtocol initializes the filter protocol
func NewLinkFilterProtocol(inner types.InterSatelliteLinkProtocol) *LinkFilterProtocol {
	return &LinkFilterProtocol{
		inner:       inner,
		links:       make(map[types.Link]struct{}),
		established: make(map[types.Link]struct{}),
		out:         make(map[types.Link]struct{}),
	}
}

func (p *LinkFilterProtocol) Mount(s types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.node != nil {
		log.Fatalln("LinkFilterProtocol should not be mounted to a node twice")
	}
	p.node = s
	p.inner.Mount(s)
}

func (p *LinkFilterProtocol) AddLink(link types.Link) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if involves(link, p.node) {
		p.links[link] = struct{}{}
	}
	p.inner.AddLink(link)
}

func (p *LinkFilterProtocol) ConnectLink(link types.Link) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.established[link] = struct{}{}
	return p.inner.ConnectLink(link)
}

func (p *LinkFilterProtocol) DisconnectLink(link types.Link) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.established, link)
	return p.inner.DisconnectLink(link)
}

func (p *LinkFilterProtocol) UpdateLinks() ([]types.Link, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.node == nil {
		return nil, errors.New("not mounted")
	}

	all, err := p.inner.UpdateLinks()
	if err != nil {
		return nil, err
	}

	oldOut := p.out
	p.out = make(map[types.Link]struct{})
	filtered := make([]types.Link, 0, len(all))
	for _, link := range all {
		if involves(link, p.node) {
			filtered = append(filtered, link)
			p.established[link] = struct{}{}
			p.out[link] = struct{}{}
			delete(oldOut, link)
		}
	}

	for l := range oldOut {
		delete(p.established, l)
	}

	return filtered, nil
}

func (p *LinkFilterProtocol) Links() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]types.Link, 0, len(p.links))
	for l := range p.links {
		out = append(out, l)
	}
	return out
}

func (p *LinkFilterProtocol) Established() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]types.Link, 0, len(p.established))
	for l := range p.established {
		out = append(out, l)
	}
	return out
}

func involves(link types.Link, node types.Node) bool {
	n1, n2 := link.Nodes()
	return n1 == node || n2 == node
}
