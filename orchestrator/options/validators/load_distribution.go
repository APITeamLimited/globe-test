package validators

import (
	"fmt"
	"math"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func LoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error ***REMOVED***
	// In case user wants equal distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.GlobalName ***REMOVED***
		clientCount := len(workerClients.Clients)

		// Set load distribution to equal distribution across all worker clients
		// Or as close as possible
		floorSize := int(math.Floor(float64(100 / float64(clientCount))))
		remainderSize := 100 - (floorSize * (clientCount - 1))

		options.LoadDistribution.Value = make([]types.LoadZone, clientCount)

		currentIndex := 0

		for _, workerClient := range workerClients.Clients ***REMOVED***
			if currentIndex == clientCount-1 ***REMOVED***
				options.LoadDistribution.Value[currentIndex] = types.LoadZone***REMOVED***
					Location: workerClient.Name,
					Fraction: remainderSize,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				options.LoadDistribution.Value[currentIndex] = types.LoadZone***REMOVED***
					Location: workerClient.Name,
					Fraction: floorSize,
				***REMOVED***
				currentIndex++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// In case user just wants default load distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == "Default" ***REMOVED***
		options.LoadDistribution = types.NullLoadDistribution***REMOVED***
			Valid: true,
			Value: []types.LoadZone***REMOVED******REMOVED***
				Location: workerClients.DefaultClient.Name,
				Fraction: 100,
			***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***

	if options.ExecutionMode.Value == types.HTTPSingleExecutionMode ***REMOVED***
		if !options.LoadDistribution.Valid ***REMOVED***
			options.LoadDistribution = types.NullLoadDistribution***REMOVED***
				Valid: true,
				Value: []types.LoadZone***REMOVED******REMOVED***
					Location: workerClients.DefaultClient.Name,
					Fraction: 100,
				***REMOVED******REMOVED***,
			***REMOVED***

			return nil
		***REMOVED***

		return checkSingleLoadDistribution(options, workerClients)
	***REMOVED*** else if options.ExecutionMode.ValueOrZero() == types.HTTPMultipleExecutionMode ***REMOVED***
		if !options.LoadDistribution.Valid ***REMOVED***
			options.LoadDistribution = types.NullLoadDistribution***REMOVED***
				Valid: true,
				Value: []types.LoadZone***REMOVED******REMOVED***
					Location: workerClients.DefaultClient.Name,
					Fraction: 100,
				***REMOVED******REMOVED***,
			***REMOVED***

			return nil
		***REMOVED***

		return checkMultiLoadDistribution(options, workerClients)
	***REMOVED***

	return fmt.Errorf("invalid execution mode %s", options.ExecutionMode.ValueOrZero())
***REMOVED***

func checkSingleLoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error ***REMOVED***
	if len(options.LoadDistribution.Value) != 1 ***REMOVED***
		return fmt.Errorf("load distribution must be a single zone when execution mode is %s", types.HTTPSingleExecutionMode)
	***REMOVED***

	if options.LoadDistribution.Value[0].Fraction != 100 ***REMOVED***
		return fmt.Errorf("load distribution fraction must be 100 when execution mode is %s", types.HTTPSingleExecutionMode)
	***REMOVED***

	// Check valid location
	for _, workerClient := range workerClients.Clients ***REMOVED***
		if options.LoadDistribution.Value[0].Location == workerClient.Name ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return fmt.Errorf("invalid location %s", options.LoadDistribution.Value[0].Location)
***REMOVED***

func checkMultiLoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error ***REMOVED***
	// Check all names valid and fractions add up to 100
	var totalFraction int

	for _, loadZone := range options.LoadDistribution.Value ***REMOVED***
		// Check valid location
		validLoadZone := false
		for _, workerClient := range workerClients.Clients ***REMOVED***
			if loadZone.Location == workerClient.Name ***REMOVED***
				validLoadZone = true
				break
			***REMOVED***
		***REMOVED***

		if !validLoadZone ***REMOVED***
			return fmt.Errorf("invalid location %s", loadZone.Location)
		***REMOVED***

		if loadZone.Fraction < 1 || loadZone.Fraction > 100 ***REMOVED***
			return fmt.Errorf("invalid fraction %d", loadZone.Fraction)
		***REMOVED***

		totalFraction += loadZone.Fraction
	***REMOVED***

	if totalFraction != 100 ***REMOVED***
		return fmt.Errorf("total fraction must be 100, got %d", totalFraction)
	***REMOVED***

	return nil
***REMOVED***
