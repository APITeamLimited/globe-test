package orchestrator

import (
	"encoding/json"
	"fmt"
	"math"

	orchOptions "github.com/APITeamLimited/globe-test/orchestrator/options"
	"github.com/APITeamLimited/redis/v9"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/google/uuid"
)

const maxJobSize = 500

type jobDistribution struct {
	Jobs         []libOrch.ChildJob `json:"jobs"`
	workerClient *redis.Client
}

func determineChildJobs(healthy bool, job libOrch.Job, options *libWorker.Options,
	workerClients libOrch.WorkerClients) (map[string]jobDistribution, error) {
	// Don't run if not healthy
	if !healthy {
		return nil, nil
	}

	childJobs := make(map[string]jobDistribution)

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

		maxVUs := options.MaxPossibleVUs.ValueOrZero()

		var subFractions = []float64{}

		// If max possible vus is greater than 500, then split into multiple jobs
		if int(maxVUs)*(loadZone.Fraction/100) <= maxJobSize {
			subFractions = append(subFractions, float64(loadZone.Fraction/100))
		} else {
			// Split into multiple jobs, each with a max of 500 vus and one job with the remainder
			// Floor plus one to ensure we don't lose any vus
			numJobs := int(math.Floor(float64(maxVUs)/maxJobSize)) + 1

			// Calculate sub fractions
			for i := 0; i < numJobs-1; i++ {
				subFractions = append(subFractions, float64(loadZone.Fraction/100)*float64(maxJobSize)/float64(maxVUs))
			}

			// Add remainder
			remainingVUs := float64(maxVUs) - float64(maxJobSize)*(float64(numJobs)-1)
			if remainingVUs > 0 {
				subFractions = append(subFractions, float64(loadZone.Fraction/100)*remainingVUs/float64(maxVUs))
			}
		}

		zoneChildJobs := make([]libOrch.ChildJob, len(subFractions))

		jobNoOptions := job
		// Remove options from job
		jobNoOptions.Options = nil

		// Create child jobs
		for i, subFraction := range subFractions {
			// Need to deep copy job, json only way that seems to work
			childOptions, _ := json.Marshal(job.Options)
			parsed := libWorker.Options{}
			json.Unmarshal(childOptions, &parsed)

			zoneChildJobs[i] = libOrch.ChildJob{
				Job:               jobNoOptions,
				ChildJobId:        uuid.NewString(),
				Options:           orchOptions.DetermineChildDerivedOptions(loadZone, workerClient, parsed, subFraction),
				UnderlyingRequest: job.UnderlyingRequest,
				FinalRequest:      job.FinalRequest,
				SubFraction:       subFraction,
			}
		}

		childJobs[loadZone.Location] = jobDistribution{
			Jobs:         zoneChildJobs,
			workerClient: workerClient.Client,
		}
	}

	return childJobs, nil
}
