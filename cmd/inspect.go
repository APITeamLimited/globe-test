package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/types"
)

// TODO: split apart like `k6 run` and `k6 archive`
func getCmdInspect(gs *globalState) *cobra.Command ***REMOVED***
	var addExecReqs bool

	// inspectCmd represents the inspect command
	inspectCmd := &cobra.Command***REMOVED***
		Use:   "inspect [file]",
		Short: "Inspect a script or archive",
		Long:  `Inspect a script or archive.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			test, err := loadTest(gs, cmd, args)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			// At the moment, `k6 inspect` output can take 2 forms: standard
			// (equal to the lib.Options struct) and extended, with additional
			// fields with execution requirements.
			var inspectOutput interface***REMOVED******REMOVED***
			if addExecReqs ***REMOVED***
				inspectOutput, err = inspectOutputWithExecRequirements(gs, cmd, test)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				inspectOutput = test.initRunner.GetOptions()
			***REMOVED***

			data, err := json.MarshalIndent(inspectOutput, "", "  ")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			printToStdout(gs, string(data))

			return nil
		***REMOVED***,
	***REMOVED***

	inspectCmd.Flags().SortFlags = false
	inspectCmd.Flags().AddFlagSet(runtimeOptionFlagSet(false))
	inspectCmd.Flags().BoolVar(&addExecReqs,
		"execution-requirements",
		false,
		"include calculations of execution requirements for the test")

	return inspectCmd
***REMOVED***

// If --execution-requirements is enabled, this will consolidate the config,
// derive the value of `scenarios` and calculate the max test duration and VUs.
func inspectOutputWithExecRequirements(gs *globalState, cmd *cobra.Command, test *loadedTest) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// we don't actually support CLI flags here, so we pass nil as the getter
	configuredTest, err := test.consolidateDeriveAndValidateConfig(gs, cmd, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	et, err := lib.NewExecutionTuple(
		configuredTest.derivedConfig.ExecutionSegment,
		configuredTest.derivedConfig.ExecutionSegmentSequence,
	)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	executionPlan := configuredTest.derivedConfig.Scenarios.GetFullExecutionRequirements(et)
	duration, _ := lib.GetEndOffset(executionPlan)

	return struct ***REMOVED***
		lib.Options
		TotalDuration types.NullDuration `json:"totalDuration"`
		MaxVUs        uint64             `json:"maxVUs"`
	***REMOVED******REMOVED***
		configuredTest.derivedConfig.Options,
		types.NewNullDuration(duration, true),
		lib.GetMaxPossibleVUs(executionPlan),
	***REMOVED***, nil
***REMOVED***
