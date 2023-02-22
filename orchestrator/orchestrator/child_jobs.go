package orchestrator

import (
	"encoding/json"

	orchOptions "github.com/APITeamLimited/globe-test/orchestrator/options"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/google/uuid"
)

const maxWorkerJobSize = 250

func determineChildJobs(healthy bool, job libOrch.Job, options *libWorker.Options) (map[string]libOrch.ChildJobDistribution, error) {
	// Don't run if not healthy
	if !healthy {
		return nil, nil
	}

	childJobs := make(map[string]libOrch.ChildJobDistribution)

	// Loop through options load distribution
	for _, loadZone := range options.LoadDistribution.Value {
		subFractions := determineSubFractions(loadZone.Fraction, job.Options.MaxPossibleVUs.Int64)

		zoneChildJobs := make([]*libOrch.ChildJob, len(subFractions))

		jobNoOptions := job
		// Remove options from job
		jobNoOptions.Options = nil

		// Create child jobs
		for i, subFraction := range subFractions {
			// Need to deep copy job, json only way that seems to work
			childOptions, _ := json.Marshal(job.Options)
			parsed := libWorker.Options{}
			json.Unmarshal(childOptions, &parsed)

			zoneChildJobs[i] = &libOrch.ChildJob{
				Job:              jobNoOptions,
				ChildJobId:       uuid.NewString(),
				ChildOptions:     orchOptions.DetermineChildDerivedOptions(loadZone, parsed, subFraction),
				SubFraction:      subFraction,
				Location:         loadZone.Location,
				WorkerConnection: nil,
			}
		}

		childJobs[loadZone.Location] = libOrch.ChildJobDistribution{
			ChildJobs: zoneChildJobs,
		}
	}

	return childJobs, nil
}

func determineSubFractions(zoneFraction float64, totalMaxVUs int64) []float64 {
	actualFraction := zoneFraction / 100
	zoneMaxVUsFloat := float64(totalMaxVUs) * actualFraction

	// Split into multiple jobs, each with a max of 500 vus and one job with the remainder

	childJobs := make([]float64, 0)

	for {
		if zoneMaxVUsFloat <= maxWorkerJobSize {
			childJobs = append(childJobs, zoneMaxVUsFloat)
			break
		}

		childJobs = append(childJobs, maxWorkerJobSize)
		zoneMaxVUsFloat -= maxWorkerJobSize
	}

	childSubFractions := make([]float64, len(childJobs))

	for i, childJob := range childJobs {
		childSubFractions[i] = childJob / float64(totalMaxVUs)
	}

	return childSubFractions
}
