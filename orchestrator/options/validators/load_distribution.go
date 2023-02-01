package validators

import (
	"fmt"

	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func LoadDistribution(options *libWorker.Options, gs libOrch.BaseGlobalState, job libOrch.Job) error {
	allowedLoadZones := gs.LoadZones()

	if gs.Standalone() {
		// In case not speified in job, use all load zones
		perittedLoadZones := allowedLoadZones

		if job.PermittedLoadZones != nil && len(job.PermittedLoadZones) > 0 {
			// Filter out worker clients that are not in permitted load zones#
			perittedLoadZones = filterPermittedLoadZones(allowedLoadZones, job.PermittedLoadZones)
		}

		err := cloudLoadDistribution(options, perittedLoadZones)
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
	}

	return localLoadDistribution(options)
}

// Filters out worker clients that are not in permitted load zones
func filterPermittedLoadZones(desiredLoadZones []string, permittedLoadZones []string) []string {
	permittedZones := make([]string, 0)

	for _, desiredLoadZone := range desiredLoadZones {
		for _, permittedLoadZone := range permittedLoadZones {
			if desiredLoadZone == permittedLoadZone {
				permittedZones = append(permittedZones, desiredLoadZone)
			}
		}
	}

	return permittedZones
}

// For each location in load distribution checks if cloud functions are available
func ensureCloudFunctionsAvailable(options *libWorker.Options, authClient libOrch.RunAuthClient) error {
	for _, location := range options.LoadDistribution.Value {
		// If cloud functions are not available for this location
		err := authClient.CheckServiceAvailability(location.Location)
		if err != nil {
			return err
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

func cloudLoadDistribution(options *libWorker.Options, permittedLoadZones []string) error {
	// In case user wants equal distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.GlobalName {
		// Check fraction is valid and set to 100
		if options.LoadDistribution.Value[0].Fraction != 100 {
			return fmt.Errorf("fraction must be 100 when using global distribution")
		}

		// Set load distribution to equal distribution across all worker clients
		locationCount := len(permittedLoadZones)
		fractionSize := 100 / float64(locationCount)

		options.LoadDistribution.Value = make([]types.LoadZone, locationCount)

		for index, loadZone := range permittedLoadZones {
			options.LoadDistribution.Value[index] = types.LoadZone{
				Location: loadZone,
				Fraction: fractionSize,
			}
		}
	}

	// In case user just wants default load distribution
	if len(options.LoadDistribution.Value) == 1 && options.LoadDistribution.Value[0].Location == libOrch.DefaultName {
		options.LoadDistribution = types.NullLoadDistribution{
			Valid: true,
			Value: []types.LoadZone{{
				Location: permittedLoadZones[0],
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
					Location: permittedLoadZones[0],
					Fraction: 100,
				}},
			}

			return nil
		}

		return checkSingleCloudLD(options, permittedLoadZones)
	} else if options.ExecutionMode.ValueOrZero() == types.HTTPMultipleExecutionMode {
		return checkMultiCloudLD(options, permittedLoadZones)
	}

	return fmt.Errorf("invalid execution mode %s", options.ExecutionMode.ValueOrZero())
}

func checkSingleCloudLD(options *libWorker.Options, permittedLoadZones []string) error {
	if len(options.LoadDistribution.Value) != 1 {
		return fmt.Errorf("load distribution must be a single zone when execution mode is %s", types.HTTPSingleExecutionMode)
	}

	if options.LoadDistribution.Value[0].Fraction != 100 {
		return fmt.Errorf("load distribution fraction must be 100 when execution mode is %s", types.HTTPSingleExecutionMode)
	}

	// Check valid location
	for _, loadZone := range permittedLoadZones {
		if options.LoadDistribution.Value[0].Location == loadZone {
			return nil
		}
	}

	return fmt.Errorf("invalid location %s", options.LoadDistribution.Value[0].Location)
}

func checkMultiCloudLD(options *libWorker.Options, permittedLoadZones []string) error {
	if !options.LoadDistribution.Valid {
		return fmt.Errorf("load distribution must be set when execution mode is %s", types.HTTPMultipleExecutionMode)
	}

	// Check all names valid and fractions add up to 100
	var totalFraction float64

	for _, loadZone := range options.LoadDistribution.Value {
		// Check valid location
		validLoadZone := false
		for _, permittedLoadZone := range permittedLoadZones {
			if loadZone.Location == permittedLoadZone {
				validLoadZone = true
				break
			}
		}

		if !validLoadZone {
			return fmt.Errorf("invalid location %s", loadZone.Location)
		}

		if loadZone.Fraction < 1 || loadZone.Fraction > 100 {
			return fmt.Errorf("invalid fraction %f", loadZone.Fraction)
		}

		totalFraction += loadZone.Fraction
	}

	// Allow for some rounding error
	if totalFraction < 99.99 || totalFraction > 100.01 {
		return fmt.Errorf("total fraction must be 100, got %f", totalFraction)
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

	if options.LoadDistribution.Value[0].Location == libOrch.GlobalName {
		options.LoadDistribution.Value[0].Location = agent.AgentWorkerName
	}

	if options.LoadDistribution.Value[0].Location != agent.AgentWorkerName {
		return fmt.Errorf("load distribution location must be %s when running locally", agent.AgentWorkerName)
	}

	return nil
}
