package validators

import (
	"fmt"
	"math"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func LoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error {
	// In case user wants equal distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.GlobalName {
		clientCount := len(workerClients.Clients)

		// Set load distribution to equal distribution across all worker clients
		// Or as close as possible
		floorSize := int(math.Floor(float64(100 / float64(clientCount))))
		remainderSize := 100 - (floorSize * (clientCount - 1))

		options.LoadDistribution.Value = make([]types.LoadZone, clientCount)

		currentIndex := 0

		for _, workerClient := range workerClients.Clients {
			if currentIndex == clientCount-1 {
				options.LoadDistribution.Value[currentIndex] = types.LoadZone{
					Location: workerClient.Name,
					Fraction: remainderSize,
				}
			} else {
				options.LoadDistribution.Value[currentIndex] = types.LoadZone{
					Location: workerClient.Name,
					Fraction: floorSize,
				}
				currentIndex++
			}
		}
	}

	// In case user just wants default load distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == "Default" {
		options.LoadDistribution = types.NullLoadDistribution{
			Valid: true,
			Value: []types.LoadZone{{
				Location: workerClients.DefaultClient.Name,
				Fraction: 100,
			}},
		}
	}

	if options.ExecutionMode.Value == types.HTTPSingleExecutionMode {
		if !options.LoadDistribution.Valid {
			options.LoadDistribution = types.NullLoadDistribution{
				Valid: true,
				Value: []types.LoadZone{{
					Location: workerClients.DefaultClient.Name,
					Fraction: 100,
				}},
			}

			return nil
		}

		return checkSingleLoadDistribution(options, workerClients)
	} else if options.ExecutionMode.ValueOrZero() == types.HTTPMultipleExecutionMode {
		if !options.LoadDistribution.Valid {
			options.LoadDistribution = types.NullLoadDistribution{
				Valid: true,
				Value: []types.LoadZone{{
					Location: workerClients.DefaultClient.Name,
					Fraction: 100,
				}},
			}

			return nil
		}

		return checkMultiLoadDistribution(options, workerClients)
	}

	return fmt.Errorf("invalid execution mode %s", options.ExecutionMode.ValueOrZero())
}

func checkSingleLoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error {
	if len(options.LoadDistribution.Value) != 1 {
		return fmt.Errorf("load distribution must be a single zone when execution mode is %s", types.HTTPSingleExecutionMode)
	}

	if options.LoadDistribution.Value[0].Fraction != 100 {
		return fmt.Errorf("load distribution fraction must be 100 when execution mode is %s", types.HTTPSingleExecutionMode)
	}

	// Check valid location
	for _, workerClient := range workerClients.Clients {
		if options.LoadDistribution.Value[0].Location == workerClient.Name {
			return nil
		}
	}

	return fmt.Errorf("invalid location %s", options.LoadDistribution.Value[0].Location)
}

func checkMultiLoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error {
	// Check all names valid and fractions add up to 100
	var totalFraction int

	for _, loadZone := range options.LoadDistribution.Value {
		// Check valid location
		validLoadZone := false
		for _, workerClient := range workerClients.Clients {
			if loadZone.Location == workerClient.Name {
				validLoadZone = true
				break
			}
		}

		if !validLoadZone {
			return fmt.Errorf("invalid location %s", loadZone.Location)
		}

		if loadZone.Fraction < 1 || loadZone.Fraction > 100 {
			return fmt.Errorf("invalid fraction %d", loadZone.Fraction)
		}

		totalFraction += loadZone.Fraction
	}

	if totalFraction != 100 {
		return fmt.Errorf("total fraction must be 100, got %d", totalFraction)
	}

	return nil
}
