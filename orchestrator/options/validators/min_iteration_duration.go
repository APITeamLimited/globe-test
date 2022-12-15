package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const minIterationDurationMilliseconds = 1000

func MinIterationDuration(options *libWorker.Options) error {
	// Ensure that user duration is within the allowed range

	// Don't enforce minimum iteration duration for HTTPSingleExecutionMode
	if options.ExecutionMode.Valid && options.ExecutionMode.Value == types.HTTPSingleExecutionMode {
		return nil
	}

	if !options.MinIterationDuration.Valid {
		// Set default
		options.MinIterationDuration = types.NullDurationFrom(time.Duration(minIterationDurationMilliseconds * time.Millisecond))
	}

	durationMilliseconds := int64(options.MinIterationDuration.TimeDuration().Milliseconds())

	// Allow errors of 1 millisecond for floating point errors
	if durationMilliseconds < (minIterationDurationMilliseconds - 1) {
		return fmt.Errorf("duration %dms is too small, the minimum is %dms", durationMilliseconds, minIterationDurationMilliseconds)
	}

	return nil
}
