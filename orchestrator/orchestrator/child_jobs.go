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

type jobDistribution struct ***REMOVED***
	Jobs         []libOrch.ChildJob `json:"jobs"`
	workerClient *redis.Client
***REMOVED***

func determineChildJobs(healthy bool, job libOrch.Job, options *libWorker.Options,
	workerClients libOrch.WorkerClients) (map[string]jobDistribution, error) ***REMOVED***
	// Don't run if not healthy
	if !healthy ***REMOVED***
		return nil, nil
	***REMOVED***

	childJobs := make(map[string]jobDistribution)

	// Loop through options load distribution
	for _, loadZone := range options.LoadDistribution.Value ***REMOVED***
		// Find worker client
		var workerClient *libOrch.NamedClient

		for _, client := range workerClients.Clients ***REMOVED***
			if client.Name == loadZone.Location ***REMOVED***
				workerClient = client
				break
			***REMOVED***
		***REMOVED***

		if workerClient == nil ***REMOVED***
			return nil, fmt.Errorf("failed to find worker client %s, this is an internal error", loadZone.Location)
		***REMOVED***

		subFractions := determineSubFractions(loadZone.Fraction)

		zoneChildJobs := make([]libOrch.ChildJob, len(subFractions))

		jobNoOptions := job
		// Remove options from job
		jobNoOptions.Options = nil

		// Create child jobs
		for i, subFraction := range subFractions ***REMOVED***
			// Need to deep copy job, json only way that seems to work
			childOptions, _ := json.Marshal(job.Options)
			parsed := libWorker.Options***REMOVED******REMOVED***
			json.Unmarshal(childOptions, &parsed)

			zoneChildJobs[i] = libOrch.ChildJob***REMOVED***
				Job:               jobNoOptions,
				ChildJobId:        uuid.NewString(),
				Options:           orchOptions.DetermineChildDerivedOptions(loadZone, workerClient, parsed, subFraction),
				UnderlyingRequest: job.UnderlyingRequest,
				FinalRequest:      job.FinalRequest,
				SubFraction:       subFraction,
			***REMOVED***
		***REMOVED***

		childJobs[loadZone.Location] = jobDistribution***REMOVED***
			Jobs:         zoneChildJobs,
			workerClient: workerClient.Client,
		***REMOVED***
	***REMOVED***

	return childJobs, nil
***REMOVED***

func determineSubFractions(fraction int) []float64 ***REMOVED***
	actualFraction := float64(fraction) / 100

	if int(actualFraction*float64(maxJobSize)) <= maxJobSize ***REMOVED***
		return []float64***REMOVED***actualFraction***REMOVED***
	***REMOVED***

	// Split into multiple jobs, each with a max of 500 vus and one job with the remainder

	// Floor plus one to ensure we don't lose any vus

	numJobs := int(math.Floor(actualFraction*float64(maxJobSize))/maxJobSize) + 1

	// Calculate sub fractions
	subFractions := make([]float64, numJobs-1)

	for i := 0; i < numJobs-1; i++ ***REMOVED***
		subFractions[i] = actualFraction * float64(maxJobSize) / float64(maxJobSize)
	***REMOVED***

	// Add remainder
	remainingVUs := actualFraction*float64(maxJobSize) - float64(maxJobSize)*float64(numJobs-1)

	if remainingVUs > 0 ***REMOVED***
		subFractions = append(subFractions, actualFraction*remainingVUs/float64(maxJobSize))
	***REMOVED***

	return subFractions
***REMOVED***
