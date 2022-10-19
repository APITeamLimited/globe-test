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
	workerClients libOrch.WorkerClients) (map[string]jobDistribution, error) ***REMOVED***
	// Don't run if not healthy
	if !healthy ***REMOVED***
		return nil, nil
	***REMOVED***

	childJobs := make(map[string]jobDistribution)

	/*childJob := libOrch.ChildJob***REMOVED***
		Job:               job,
		ChildJobId:        uuid.NewString(),
		Options:           *options,
		UnderlyingRequest: job.UnderlyingRequest,
		FinalRequest:      job.FinalRequest,
	***REMOVED***

	childJobs["portsmouth"] = jobDistribution***REMOVED***
		jobs:         &[]libOrch.ChildJob***REMOVED***childJob***REMOVED***,
		workerClient: workerClients.DefaultClient.Client,
	***REMOVED****/

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

		totalFraction := loadZone.Fraction

		var subFractions = []float32***REMOVED******REMOVED***

		// If max possible vus is greater than 500, then split into multiple jobs
		if int(options.MaxPossibleVUs.ValueOrZero())*(totalFraction/100) <= maxJobSize ***REMOVED***
			subFractions = append(subFractions, float32(totalFraction))
		***REMOVED*** else ***REMOVED***
			// Split into multiple jobs, each with a max of 500 vus and one job with the remainder
			// Floor plus one to ensure we don't lose any vus
			numJobs := int(math.Floor(float64(options.MaxPossibleVUs.ValueOrZero())/maxJobSize)) + 1

			// Calculate sub fractions
			for i := 0; i < numJobs-1; i++ ***REMOVED***
				subFractions = append(subFractions, float32(totalFraction)/maxJobSize)
			***REMOVED***

			// Add remainder
			subFractions = append(subFractions, float32(totalFraction%maxJobSize)/maxJobSize)
		***REMOVED***

		zoneChildJobs := make([]libOrch.ChildJob, len(subFractions))

		jobNoOptions := job
		// Remove options from job
		jobNoOptions.Options = nil

		// Create child jobs
		for i, subFraction := range subFractions ***REMOVED***
			zoneChildJobs[i] = libOrch.ChildJob***REMOVED***
				Job:               jobNoOptions,
				ChildJobId:        uuid.NewString(),
				Options:           orchOptions.DetermineChildDerivedOptions(loadZone, workerClient, *options, subFraction),
				UnderlyingRequest: job.UnderlyingRequest,
				FinalRequest:      job.FinalRequest,
			***REMOVED***
		***REMOVED***

		childJobs[loadZone.Location] = jobDistribution***REMOVED***
			jobs:         zoneChildJobs,
			workerClient: workerClient.Client,
		***REMOVED***
	***REMOVED***

	return childJobs, nil
***REMOVED***
