package validators

import (
	"net"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

func BlacklistIPs(options *libWorker.Options, bannedIPNets *[]*net.IPNet) {
	for _, ipNet := range *bannedIPNets {
		options.BlacklistIPs = append(options.BlacklistIPs, &libWorker.IPNet{
			IPNet: *ipNet,
		})
	}
}
