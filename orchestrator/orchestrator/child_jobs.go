package orchestrator

import (
	"fmt"
	"math"

	orchOptions "github.com/APITeamLimited/globe-test/orchestrator/options"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/google/uuid"
)

const maxJobSize = 500

func determineChildJobs(healthy bool, job libOrch.Job, options *libWorker.Options,
	workerClients libOrch.WorkerClients) (map[string]jobDistribution, error) {
	// Don't run if not healthy
	if !healthy {
		return nil, nil
	}

	childJobs := make(map[string]jobDistribution)

	/*childJob := libOrch.ChildJob{
		Job:               job,
		ChildJobId:        uuid.NewString(),
		Options:           *options,
		UnderlyingRequest: job.UnderlyingRequest,
		FinalRequest:      job.FinalRequest,
	}

	childJobs["portsmouth"] = jobDistribution{
		jobs:         &[]libOrch.ChildJob{childJob},
		workerClient: workerClients.DefaultClient.Client,
	}*/

	// Loop through options load distribution
	for _, loadZone := range options.LoadDistribution.Value {
		// Find worker client
		var workerClient *libOrch.NamedClient

		for _, client := range workerClients.Clients {
			if client.Name == loadZone.Location {
				workerClient = client
				break
			}
		}

		if workerClient == nil {
			return nil, fmt.Errorf("failed to find worker client %s, this is an internal error", loadZone.Location)
		}

		totalFraction := loadZone.Fraction

		var subFractions = []float32{}

		// If max possible vus is greater than 500, then split into multiple jobs
		if int(options.MaxPossibleVUs.ValueOrZero())*(totalFraction/100) <= maxJobSize {
			subFractions = append(subFractions, float32(totalFraction))
		} else {
			// Split into multiple jobs, each with a max of 500 vus and one job with the remainder
			// Floor plus one to ensure we don't lose any vus
			numJobs := int(math.Floor(float64(options.MaxPossibleVUs.ValueOrZero())/maxJobSize)) + 1

			// Calculate sub fractions
			for i := 0; i < numJobs-1; i++ {
				subFractions = append(subFractions, float32(totalFraction)/maxJobSize)
			}

			// Add remainder
			subFractions = append(subFractions, float32(totalFraction%maxJobSize)/maxJobSize)
		}

		zoneChildJobs := make([]libOrch.ChildJob, len(subFractions))

		jobNoOptions := job
		// Remove options from job
		jobNoOptions.Options = nil

		// Create child jobs
		for i, subFraction := range subFractions {
			zoneChildJobs[i] = libOrch.ChildJob{
				Job:               jobNoOptions,
				ChildJobId:        uuid.NewString(),
				Options:           orchOptions.DetermineChildDerivedOptions(loadZone, workerClient, *options, subFraction),
				UnderlyingRequest: job.UnderlyingRequest,
				FinalRequest:      job.FinalRequest,
			}
		}

		childJobs[loadZone.Location] = jobDistribution{
			jobs:         zoneChildJobs,
			workerClient: workerClient.Client,
		}
	}

	return childJobs, nil
}
