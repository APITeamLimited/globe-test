package validators

import (
	"net"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

func BlacklistIPs(options *libWorker.Options, bannedIPNets *[]*net.IPNet) ***REMOVED***
	for _, ipNet := range *bannedIPNets ***REMOVED***
		options.BlacklistIPs = append(options.BlacklistIPs, &libWorker.IPNet***REMOVED***
			IPNet: *ipNet,
		***REMOVED***)
	***REMOVED***
***REMOVED***
