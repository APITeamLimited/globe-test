package validators

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/lib/types"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

const setupTimeoutSeconds = 60 * 1

func SetupTimeout(options *libWorker.Options) error {
	// Ensure that user duration is within the allowed range

	if !options.SetupTimeout.Valid {
		options.SetupTimeout = types.NullDurationFrom(time.Duration(setupTimeoutSeconds * time.Second))
	}

	durationMilliseconds := int64(options.SetupTimeout.TimeDuration().Milliseconds())

	if durationMilliseconds > (setupTimeoutSeconds * 1000) {
		// Round duration to nearest 0.1 second
		return fmt.Errorf("setupTimeout duration is too large, the maximum is %ds", setupTimeoutSeconds)
	}

	return nil
}
