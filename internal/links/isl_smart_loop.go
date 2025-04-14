package links

import (
	"errors"
	"sync"

	"stardustGo/configs"
	linkmod "stardustGo/internal/links/linktypes"
	"stardustGo/pkg/types"
)

// IslAddSmartLoopProtocol wraps another IInterSatelliteLinkProtocol and augments it
// by adding extra links to form smart loops (e.g., when a satellite has only 1 connection).
type IslAddSmartLoopProtocol struct {
	inner       types.IInterSatelliteLinkProtocol
	config      configs.InterSatelliteLinkConfig
	satellite   types.INode
	position    types.Vector
	mu          sync.Mutex
	readyCh     chan struct{}
	lastUpdated []types.ILink
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
func (p *IslAddSmartLoopProtocol) Mount(s types.INode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.satellite == nil {
		p.satellite = s
		p.inner.Mount(s)
	}
}

// AddLink registers a link to the underlying protocol
func (p *IslAddSmartLoopProtocol) AddLink(link types.ILink) {
	p.inner.AddLink(link)
}

// ConnectLink forwards link connection to the inner protocol
func (p *IslAddSmartLoopProtocol) ConnectLink(link types.ILink) error {
	return p.inner.ConnectLink(link)
}

// DisconnectLink forwards link disconnection to the inner protocol
func (p *IslAddSmartLoopProtocol) DisconnectLink(link types.ILink) error {
	return p.inner.DisconnectLink(link)
}

// ConnectSatellite forwards satellite connection to the inner protocol
func (p *IslAddSmartLoopProtocol) ConnectSatellite(s types.INode) error {
	return p.inner.ConnectSatellite(s)
}

// DisconnectSatellite forwards satellite disconnection to the inner protocol
func (p *IslAddSmartLoopProtocol) DisconnectSatellite(s types.INode) error {
	return p.inner.DisconnectSatellite(s)
}

// Links returns all candidate links from the inner protocol
func (p *IslAddSmartLoopProtocol) Links() []types.ILink {
	return p.inner.Links()
}

// Established returns all currently established links
func (p *IslAddSmartLoopProtocol) Established() []types.ILink {
	return p.inner.Established()
}

// UpdateLinks enhances the underlying protocol's MST with smart loops
func (p *IslAddSmartLoopProtocol) UpdateLinks() ([]types.ILink, error) {
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

	satToLinks := make(map[types.INode]int)
	uniqueLinks := make(map[types.INode]types.ILink)

	// Count how many times each satellite appears in links
	for _, link := range mstLinks {
		n1 := link.GetOther(nil)
		n2 := link.GetOther(n1)

		for _, n := range []types.INode{n1, n2} {
			satToLinks[n]++
			if satToLinks[n] == 1 {
				uniqueLinks[n] = link
			} else {
				delete(uniqueLinks, n)
			}
		}
	}

	var additions []types.ILink
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
