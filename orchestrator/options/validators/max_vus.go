package validators

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

func MaxVUs(maxVUsCount int64, job libOrch.Job) error {
	if job.MaxSimulatedUsers != 0 && maxVUsCount > job.MaxSimulatedUsers {
		return fmt.Errorf("max VU count is limited to %d", job.MaxSimulatedUsers)
	}

	return nil
}
