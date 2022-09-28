package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const teardownTimeoutSeconds = 10

func TeardownTimeout(options *libWorker.Options) error ***REMOVED***
	// Ensure that user duration is within the allowed range

	if !options.TeardownTimeout.Valid ***REMOVED***
		options.TeardownTimeout = types.NullDurationFrom(time.Duration(teardownTimeoutSeconds * time.Second))
	***REMOVED***

	durationMilliseconds := int64(options.TeardownTimeout.TimeDuration().Milliseconds())

	if durationMilliseconds > (teardownTimeoutSeconds * 1000) ***REMOVED***
		// Round duration to nearest 0.1 second
		return fmt.Errorf("teardownTimeout duration is too large, the maximum is %ds", teardownTimeoutSeconds)
	***REMOVED***

	return nil
***REMOVED***
