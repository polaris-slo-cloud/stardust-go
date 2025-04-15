package links

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
	"log"
)

// IslProtocolBuilder constructs inter-satellite link protocols based on config
// It wraps MST, PST, and smart loop strategies with filtering or enhancements as needed
// Available protocols: mst, pst, mst_loop, pst_loop, mst_smart_loop, pst_smart_loop, other_mst, other_mst_loop, other_mst_smart_loop, nearest

type IslProtocolBuilder struct {
	config   configs.InterSatelliteLinkConfig
	mst      *IslMstProtocol
	pst      *IslPstProtocol
	otherMst *IslSatelliteCentricMstProtocol
}

// NewIslProtocolBuilder initializes a protocol builder instance
func NewIslProtocolBuilder(cfg configs.InterSatelliteLinkConfig) *IslProtocolBuilder {
	return &IslProtocolBuilder{config: cfg}
}

// Build selects and wraps the desired link protocol
func (b *IslProtocolBuilder) Build() types.IInterSatelliteLinkProtocol {
	switch b.config.Protocol {
	case "mst":
		return NewIslFilterProtocol(b.getMst())
	case "pst":
		return NewIslFilterProtocol(b.getPst())
	case "mst_loop":
		return NewIslAddLoopProtocol(NewIslFilterProtocol(b.getMst()), b.config)
	case "pst_loop":
		return NewIslAddLoopProtocol(NewIslFilterProtocol(b.getPst()), b.config)
	case "mst_smart_loop":
		return NewIslFilterProtocol(NewIslAddSmartLoopProtocol(b.getMst(), b.config))
	case "pst_smart_loop":
		return NewIslFilterProtocol(NewIslAddSmartLoopProtocol(b.getPst(), b.config))
	case "other_mst":
		return NewIslFilterProtocol(b.getOtherMst())
	case "other_mst_loop":
		return NewIslAddLoopProtocol(NewIslFilterProtocol(b.getOtherMst()), b.config)
	case "other_mst_smart_loop":
		return NewIslFilterProtocol(NewIslAddSmartLoopProtocol(b.getOtherMst(), b.config))
	case "nearest":
		return NewIslNearestProtocol(b.config)
	default:
		log.Printf("[WARN] Unknown ISL protocol '%s', falling back to 'nearest'", b.config.Protocol)
		return NewIslNearestProtocol(b.config)
	}
}

func (b *IslProtocolBuilder) getMst() *IslMstProtocol {
	if b.mst == nil {
		b.mst = NewIslMstProtocol()
	}
	return b.mst
}

func (b *IslProtocolBuilder) getPst() *IslPstProtocol {
	if b.pst == nil {
		b.pst = NewIslPstProtocol()
	}
	return b.pst
}

func (b *IslProtocolBuilder) getOtherMst() *IslSatelliteCentricMstProtocol {
	if b.otherMst == nil {
		b.otherMst = NewIslSatelliteCentricMstProtocol()
	}
	return b.otherMst
}
