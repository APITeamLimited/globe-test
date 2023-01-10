package options

import (
	"net"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

// Returns internal ip nets that are banned
func generateBannedIPNets(gs libOrch.BaseGlobalState) *[]*net.IPNet {
	// If we are not running in standalone mode, we don't need to ban any ip ranges
	// as this is a local test
	if !gs.Standalone() {
		return &[]*net.IPNet{}
	}

	localhostIPNets := make([]*net.IPNet, 0, 4)

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv4(10, 0, 0, 0),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	})

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv4(172, 16, 0, 0),
		Mask: net.IPv4Mask(255, 240, 0, 0),
	})

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv4(192, 168, 0, 0),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	})

	localhostIPNets = append(localhostIPNets, &net.IPNet{
		IP:   net.IPv6loopback,
		Mask: net.CIDRMask(128, 128),
	})

	return &localhostIPNets
}
