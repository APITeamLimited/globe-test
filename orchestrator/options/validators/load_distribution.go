package validators

import (
	"fmt"
	"math"

	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func LoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients,
	gs libOrch.BaseGlobalState, job libOrch.Job) error {
	if gs.FuncMode() {
		permittedWorkerClients := workerClients

		if job.PermittedLoadZones != nil && len(job.PermittedLoadZones) > 0 {
			// Filter out worker clients that are not in permitted load zones#
			permittedWorkerClients = filterPermittedLoadZones(workerClients, job.PermittedLoadZones)
		}

		err := cloudLoadDistribution(options, permittedWorkerClients)
		if err != nil {
			return err
		}

		err = ensureCloudFunctionsAvailable(options, gs.FuncAuthClient())
		if err != nil {
			return err
		}

		if job.PermittedLoadZones != nil && len(job.PermittedLoadZones) > 0 {
			err = ensurePermittedLoadZones(options, job.PermittedLoadZones)
			if err != nil {
				return err
			}
		}

		return nil
	} else if gs.Standalone() {
		return cloudLoadDistribution(options, workerClients)
	}
	return localLoadDistribution(options)
}

// Filters out worker clients that are not in permitted load zones
func filterPermittedLoadZones(workerClients libOrch.WorkerClients, permittedLoadZones []string) libOrch.WorkerClients {
	permittedWorkerClients := libOrch.WorkerClients{
		Clients:       make(map[string]*libOrch.NamedClient),
		DefaultClient: workerClients.DefaultClient,
	}

	// Filter out worker clients that are not in permitted load zones
	for _, workerClient := range workerClients.Clients {
		for _, permittedLoadZone := range permittedLoadZones {
			if workerClient.Name == permittedLoadZone {
				permittedWorkerClients.Clients[workerClient.Name] = workerClient
			}
		}
	}

	return permittedWorkerClients
}

// For each location in load distribution checks if cloud functions are available
func ensureCloudFunctionsAvailable(options *libWorker.Options, authClient libOrch.FunctionAuthClient) error {
	for _, location := range options.LoadDistribution.Value {
		// If cloud functions are not available for this location
		error := authClient.CheckFunctionAvailability(location.Location)
		if error != nil {
			return error
		}
	}

	return nil
}

func ensurePermittedLoadZones(options *libWorker.Options, loadZones []string) error {
	// Only allow load zones named in loadZones
	for _, loadZone := range options.LoadDistribution.Value {
		validLoadZone := false
		for _, allowedLoadZone := range loadZones {
			if loadZone.Location == allowedLoadZone {
				validLoadZone = true
				break
			}
		}

		if !validLoadZone {
			return fmt.Errorf("invalid location %s is not permitted", loadZone.Location)
		}
	}

	return nil
}

func cloudLoadDistribution(options *libWorker.Options, workerClients libOrch.WorkerClients) error {
	// In case user wants equal distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.GlobalName {
		// Check fraction is valid and set to 100
		if options.LoadDistribution.Value[0].Fraction != 100 {
			return fmt.Errorf("fraction must be 100 when using global distribution")
		}

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
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.DefaultName {
		options.LoadDistribution = types.NullLoadDistribution{
			Valid: true,
			Value: []types.LoadZone{{
				Location: workerClients.DefaultClient.Name,
				Fraction: 100,
			}},
		}
	}

	// Check execution mode scenarios
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

		return checkSingleCloudLD(options, workerClients)
	} else if options.ExecutionMode.ValueOrZero() == types.HTTPMultipleExecutionMode {
		return checkMultiCloudLD(options, workerClients)
	}

	return fmt.Errorf("invalid execution mode %s", options.ExecutionMode.ValueOrZero())
}

func checkSingleCloudLD(options *libWorker.Options, workerClients libOrch.WorkerClients) error {
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

func checkMultiCloudLD(options *libWorker.Options, workerClients libOrch.WorkerClients) error {
	if !options.LoadDistribution.Valid {
		return fmt.Errorf("load distribution must be set when execution mode is %s", types.HTTPMultipleExecutionMode)
	}

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

func localLoadDistribution(options *libWorker.Options) error {
	// Check single load zone

	// If no config, just apply default
	if !options.LoadDistribution.Valid {
		options.LoadDistribution = types.NullLoadDistribution{
			Valid: true,
			Value: []types.LoadZone{{
				Location: agent.AgentWorkerName,
				Fraction: 100,
			}},
		}

		return nil
	}

	// In case user just wants default load distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.DefaultName {
		options.LoadDistribution = types.NullLoadDistribution{
			Valid: true,
			Value: []types.LoadZone{{
				Location: agent.AgentWorkerName,
				Fraction: 100,
			}},
		}
	}

	if len(options.LoadDistribution.Value) != 1 {
		return fmt.Errorf("load distribution must be a single zone when running locally")
	}

	if options.LoadDistribution.Value[0].Fraction != 100 {
		return fmt.Errorf("load distribution fraction must be 100 when running locally")
	}

	if options.LoadDistribution.Value[0].Location != agent.AgentWorkerName {
		return fmt.Errorf("load distribution location must be %s when running locally", agent.AgentWorkerName)
	}

	return nil
}
