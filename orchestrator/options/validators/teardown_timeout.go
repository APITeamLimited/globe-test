package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

const teardownTimeoutSeconds = 30

func TeardownTimeout(options *libWorker.Options) error {
	// Ensure that user duration is within the allowed range

	if !options.TeardownTimeout.Valid {
		options.TeardownTimeout = types.NullDurationFrom(time.Duration(teardownTimeoutSeconds * time.Second))
	}

	durationMilliseconds := int64(options.TeardownTimeout.TimeDuration().Milliseconds())

	if durationMilliseconds > (teardownTimeoutSeconds * 1000) {
		// Round duration to nearest 0.1 second
		return fmt.Errorf("teardownTimeout duration is too large, the maximum is %ds", teardownTimeoutSeconds)
	}

	return nil
}
