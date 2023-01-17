package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const maxDurationMinutesDefault int64 = 50

func Duration(options *libWorker.Options, job libOrch.Job) error {
	// Ensure that user duration is within the allowed range
	maxDurationMinutes := maxDurationMinutesDefault

	if job.MaxTestDurationMinutes != 0 {
		maxDurationMinutes = job.MaxTestDurationMinutes
	}

	// TODO: in future validate based off user plan

	if (!options.Duration.Valid) || (options.Duration.TimeDuration() < 0) {
		options.Duration = types.NullDurationFrom(time.Duration(maxDurationMinutes * int64(time.Minute)))
	}

	durationMinutes := int64(options.Duration.TimeDuration().Minutes())

	if durationMinutes > maxDurationMinutes {
		return fmt.Errorf("duration %d is too large, the maximum is %d", durationMinutes, maxDurationMinutes)
	}

	// TODO: enforce max duration

	return nil
}
