package validators

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"gopkg.in/guregu/null.v3"
)

func InsecureSkipTLSVerify(options *libWorker.Options) {
	// Ensure that user duration is within the allowed range

	if !options.InsecureSkipTLSVerify.Valid {
		options.InsecureSkipTLSVerify = null.BoolFrom(false)
	}
}
