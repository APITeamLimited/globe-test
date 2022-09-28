package validators

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"gopkg.in/guregu/null.v3"
)

func InsecureSkipTLSVerify(options *libWorker.Options) ***REMOVED***
	// Ensure that user duration is within the allowed range

	if !options.InsecureSkipTLSVerify.Valid ***REMOVED***
		options.InsecureSkipTLSVerify = null.BoolFrom(false)
	***REMOVED***
***REMOVED***
