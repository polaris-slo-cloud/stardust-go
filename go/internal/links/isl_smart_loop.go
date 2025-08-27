package links

import (
	"errors"
	"sync"

	"github.com/keniack/stardustGo/configs"
	linkmod "github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/types"
)

// IslAddSmartLoopProtocol wraps another IInterSatelliteLinkProtocol and augments it
// by adding extra links to form smart loops (e.g., when a satellite has only 1 connection).
type IslAddSmartLoopProtocol struct {
	inner       types.IInterSatelliteLinkProtocol
	config      configs.InterSatelliteLinkConfig
	satellite   types.Node
	position    types.Vector
	mu          sync.Mutex
	readyCh     chan struct{}
	lastUpdated []types.Link
}

// NewIslAddSmartLoopProtocol creates a new smart loop-enhancing protocol
func NewIslAddSmartLoopProtocol(inner types.IInterSatelliteLinkProtocol, cfg configs.InterSatelliteLinkConfig) *IslAddSmartLoopProtocol {
	return &IslAddSmartLoopProtocol{
		inner:   inner,
		config:  cfg,
		readyCh: make(chan struct{}, 1),
	}
}

// Mount binds the protocol to a satellite node
func (p *IslAddSmartLoopProtocol) Mount(s types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.satellite == nil {
		p.satellite = s
		p.inner.Mount(s)
	}
}

// AddLink registers a link to the underlying protocol
func (p *IslAddSmartLoopProtocol) AddLink(link types.Link) {
	p.inner.AddLink(link)
}

// ConnectLink forwards link connection to the inner protocol
func (p *IslAddSmartLoopProtocol) ConnectLink(link types.Link) error {
	return p.inner.ConnectLink(link)
}

// DisconnectLink forwards link disconnection to the inner protocol
func (p *IslAddSmartLoopProtocol) DisconnectLink(link types.Link) error {
	return p.inner.DisconnectLink(link)
}

// ConnectSatellite forwards satellite connection to the inner protocol
func (p *IslAddSmartLoopProtocol) ConnectSatellite(s types.Node) error {
	return p.inner.ConnectSatellite(s)
}

// DisconnectSatellite forwards satellite disconnection to the inner protocol
func (p *IslAddSmartLoopProtocol) DisconnectSatellite(s types.Node) error {
	return p.inner.DisconnectSatellite(s)
}

// Links returns all candidate links from the inner protocol
func (p *IslAddSmartLoopProtocol) Links() []types.Link {
	return p.inner.Links()
}

// Established returns all currently established links
func (p *IslAddSmartLoopProtocol) Established() []types.Link {
	return p.inner.Established()
}

// UpdateLinks enhances the underlying protocol's MST with smart loops
func (p *IslAddSmartLoopProtocol) UpdateLinks() ([]types.Link, error) {
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
		return p.lastUpdated, nil
	}
	p.position = p.satellite.PositionVector()

	mstLinks, err := p.inner.UpdateLinks()
	if err != nil {
		return nil, err
	}

	satToLinks := make(map[types.Node]int)
	uniqueLinks := make(map[types.Node]types.Link)

	// Count how many times each satellite appears in links
	for _, link := range mstLinks {
		n1 := link.GetOther(nil)
		n2 := link.GetOther(n1)

		for _, n := range []types.Node{n1, n2} {
			satToLinks[n]++
			if satToLinks[n] == 1 {
				uniqueLinks[n] = link
			} else {
				delete(uniqueLinks, n)
			}
		}
	}

	var additions []types.Link
	for s := range uniqueLinks {
		for _, candidate := range p.inner.Links() {
			if isl, ok := candidate.(*linkmod.IslLink); ok && isl.Distance() <= configs.MaxISLDistance {
				if isl.Node1 != s && isl.Node2 != s {
					continue
				}
				other := isl.GetOther(s)
				if other == nil || satToLinks[s] >= p.config.Neighbours || satToLinks[other] >= p.config.Neighbours {
					continue
				}
				satToLinks[s]++
				satToLinks[other]++
				isl.SetEstablished(true)
				additions = append(additions, isl)
			}
		}
	}

	p.lastUpdated = append(mstLinks, additions...)
	select {
	case p.readyCh <- struct{}{}:
	default:
	}
	return p.lastUpdated, nil
}
