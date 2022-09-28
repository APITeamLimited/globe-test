package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const maxDurationMinutes = 10

func Duration(options *libWorker.Options) error {
	// Ensure that user duration is within the allowed range

	// TODO: in future validate based off user plan

	if (!options.Duration.Valid) || (options.Duration.TimeDuration() < 0) {
		options.Duration = types.NullDurationFrom(time.Duration(maxDurationMinutes * time.Minute))
	}

	durationMinutes := int64(options.Duration.TimeDuration().Minutes())

	if durationMinutes > maxDurationMinutes {
		return fmt.Errorf("duration %d is too large, the maximum is %d", durationMinutes, maxDurationMinutes)
	}

	return nil
}
