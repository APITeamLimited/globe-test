package validators

import (
	"database/sql"
	"fmt"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"gopkg.in/guregu/null.v3"
)

func MaxVUs(checkedOptions *libWorker.Options, job libOrch.Job) error {
	// Find max possible VU count
	maxVUsCount := int64(0)
	for _, scenario := range checkedOptions.Scenarios {
		maxVUsCount += scenario.GetMaxExecutorVUs()
	}

	checkedOptions.MaxPossibleVUs = null.Int{
		NullInt64: sql.NullInt64{
			Int64: maxVUsCount,
			Valid: true,
		},
	}

	if job.MaxSimulatedUsers != 0 && maxVUsCount > job.MaxSimulatedUsers {
		return fmt.Errorf("max VU count is limited to %d", job.MaxSimulatedUsers)
	}

	return nil
}
