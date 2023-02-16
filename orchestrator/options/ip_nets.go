package options

import (
	"net"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

// Returns internal ip nets that are banned
func getBannedIPNets(gs libOrch.BaseGlobalState) *[]*net.IPNet {
	// If we are not running in standalone mode, we don't need to ban any ip ranges
	// as this is a local test
	if !gs.Standalone() {
		return &[]*net.IPNet{}
	}

	return generateBannedIPNets()
}

func generateBannedIPNets() *[]*net.IPNet {
	localhostIPNets := make([]*net.IPNet, 6)

	// Private IPv4 addresses
	localhostIPNets[0] = &net.IPNet{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)}
	localhostIPNets[1] = &net.IPNet{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)}
	localhostIPNets[2] = &net.IPNet{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)}
	// Private IPv6 addresses
	localhostIPNets[3] = &net.IPNet{IP: net.ParseIP("fc00::"), Mask: net.CIDRMask(7, 128)}
	localhostIPNets[4] = &net.IPNet{IP: net.ParseIP("fe80::"), Mask: net.CIDRMask(10, 128)}
	localhostIPNets[5] = &net.IPNet{IP: net.ParseIP("::1"), Mask: net.CIDRMask(128, 128)}

	return &localhostIPNets
}
