package validators

import (
	"fmt"
	"net"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

func Hosts(options *libWorker.Options, bannedIPNets *[]*net.IPNet) error {
	for _, ip := range options.Hosts {
		for _, localhostIPNet := range *bannedIPNets {
			netIp := net.ParseIP(string(ip.IP))
			if netIp != nil && localhostIPNet.Contains(netIp) {
				return fmt.Errorf("host %s is banned", ip)
			}
		}
	}

	return nil
}
