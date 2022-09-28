package validators

import (
	"fmt"
	"net"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

func Hosts(options *libWorker.Options, bannedIPNets *[]*net.IPNet) error ***REMOVED***
	for _, ip := range options.Hosts ***REMOVED***
		for _, localhostIPNet := range *bannedIPNets ***REMOVED***
			netIp := net.ParseIP(string(ip.IP))
			if netIp != nil && localhostIPNet.Contains(netIp) ***REMOVED***
				return fmt.Errorf("host %s is banned", ip)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
