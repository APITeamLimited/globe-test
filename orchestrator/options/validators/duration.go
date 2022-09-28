package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const maxDurationMinutes = 10

func Duration(options *libWorker.Options) error ***REMOVED***
	// Ensure that user duration is within the allowed range

	// TODO: in future validate based off user plan

	if (!options.Duration.Valid) || (options.Duration.TimeDuration() < 0) ***REMOVED***
		options.Duration = types.NullDurationFrom(time.Duration(maxDurationMinutes * time.Minute))
	***REMOVED***

	durationMinutes := int64(options.Duration.TimeDuration().Minutes())

	if durationMinutes > maxDurationMinutes ***REMOVED***
		return fmt.Errorf("duration %d is too large, the maximum is %d", durationMinutes, maxDurationMinutes)
	***REMOVED***

	return nil
***REMOVED***
