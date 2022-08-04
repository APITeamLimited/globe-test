package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/output"
	"go.k6.io/k6/output/cloud"
	"go.k6.io/k6/output/csv"
	"go.k6.io/k6/output/influxdb"
	"go.k6.io/k6/output/json"
	"go.k6.io/k6/output/statsd"
)

// TODO: move this to an output sub-module after we get rid of the old collectors?
func getAllOutputConstructors() (map[string]func(output.Params) (output.Output, error), error) ***REMOVED***
	// Start with the built-in outputs
	result := map[string]func(output.Params) (output.Output, error)***REMOVED***
		"json":     json.New,
		"cloud":    cloud.New,
		"influxdb": influxdb.New,
		"kafka": func(params output.Params) (output.Output, error) ***REMOVED***
			return nil, errors.New("the kafka output was deprecated in k6 v0.32.0 and removed in k6 v0.34.0, " +
				"please use the new xk6 kafka output extension instead - https://github.com/k6io/xk6-output-kafka")
		***REMOVED***,
		"statsd": statsd.New,
		"datadog": func(params output.Params) (output.Output, error) ***REMOVED***
			return nil, errors.New("the datadog output was deprecated in k6 v0.32.0 and removed in k6 v0.34.0, " +
				"please use the statsd output with env. variable K6_STATSD_ENABLE_TAGS=true instead")
		***REMOVED***,
		"csv": csv.New,
	***REMOVED***

	exts := output.GetExtensions()
	for k, v := range exts ***REMOVED***
		if _, ok := result[k]; ok ***REMOVED***
			return nil, fmt.Errorf("invalid output extension %s, built-in output with the same type already exists", k)
		***REMOVED***
		result[k] = v
	***REMOVED***

	return result, nil
***REMOVED***

func getPossibleIDList(constrs map[string]func(output.Params) (output.Output, error)) string ***REMOVED***
	res := make([]string, 0, len(constrs))
	for k := range constrs ***REMOVED***
		if k == "kafka" || k == "datadog" ***REMOVED***
			continue
		***REMOVED***
		res = append(res, k)
	***REMOVED***
	sort.Strings(res)
	return strings.Join(res, ", ")
***REMOVED***

func createOutputs(
	gs *globalState, test *loadedAndConfiguredTest, executionPlan []lib.ExecutionStep,
) ([]output.Output, error) ***REMOVED***
	outputConstructors, err := getAllOutputConstructors()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	baseParams := output.Params***REMOVED***
		ScriptPath:     test.source.URL,
		Logger:         gs.logger,
		Environment:    gs.envVars,
		StdOut:         gs.stdOut,
		StdErr:         gs.stdErr,
		FS:             gs.fs,
		ScriptOptions:  test.derivedConfig.Options,
		RuntimeOptions: test.preInitState.RuntimeOptions,
		ExecutionPlan:  executionPlan,
	***REMOVED***
	result := make([]output.Output, 0, len(test.derivedConfig.Out))

	for _, outputFullArg := range test.derivedConfig.Out ***REMOVED***
		outputType, outputArg := parseOutputArgument(outputFullArg)
		outputConstructor, ok := outputConstructors[outputType]
		if !ok ***REMOVED***
			return nil, fmt.Errorf(
				"invalid output type '%s', available types are: %s",
				outputType, getPossibleIDList(outputConstructors),
			)
		***REMOVED***

		params := baseParams
		params.OutputType = outputType
		params.ConfigArgument = outputArg
		params.JSONConfig = test.derivedConfig.Collectors[outputType]

		out, err := outputConstructor(params)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("could not create the '%s' output: %w", outputType, err)
		***REMOVED***

		if thresholdOut, ok := out.(output.WithThresholds); ok ***REMOVED***
			thresholdOut.SetThresholds(test.derivedConfig.Thresholds)
		***REMOVED***

		if builtinMetricOut, ok := out.(output.WithBuiltinMetrics); ok ***REMOVED***
			builtinMetricOut.SetBuiltinMetrics(test.preInitState.BuiltinMetrics)
		***REMOVED***

		result = append(result, out)
	***REMOVED***

	return result, nil
***REMOVED***

func parseOutputArgument(s string) (t, arg string) ***REMOVED***
	parts := strings.SplitN(s, "=", 2)
	switch len(parts) ***REMOVED***
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	***REMOVED***
***REMOVED***
