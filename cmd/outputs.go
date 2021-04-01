/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/output/cloud"
	"github.com/loadimpact/k6/output/json"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/csv"
	"github.com/loadimpact/k6/stats/datadog"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/loadimpact/k6/stats/statsd"

	"github.com/k6io/xk6-output-kafka/pkg/kafka"
)

// TODO: move this to an output sub-module after we get rid of the old collectors?
//nolint: funlen
func getAllOutputConstructors() (map[string]func(output.Params) (output.Output, error), error) ***REMOVED***
	// Start with the built-in outputs
	result := map[string]func(output.Params) (output.Output, error)***REMOVED***
		"json":  json.New,
		"cloud": cloud.New,

		// TODO: remove all of these
		"influxdb": func(params output.Params) (output.Output, error) ***REMOVED***
			conf, err := influxdb.GetConsolidatedConfig(params.JSONConfig, params.Environment, params.ConfigArgument)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			influxc, err := influxdb.New(params.Logger, conf)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return newCollectorAdapter(params, influxc), nil
		***REMOVED***,
		"kafka": func(params output.Params) (output.Output, error) ***REMOVED***
			params.Logger.Warn("The kafka output is deprecated, and will be removed in a future k6 version. " +
				"Please use the new xk6 kafka output extension instead. " +
				"It can be found at https://github.com/k6io/xk6-output-kafka.")
			return kafka.New(params)
		***REMOVED***,
		"statsd": func(params output.Params) (output.Output, error) ***REMOVED***
			conf, err := statsd.GetConsolidatedConfig(params.JSONConfig, params.Environment)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			statsdc, err := statsd.New(params.Logger, conf)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return newCollectorAdapter(params, statsdc), nil
		***REMOVED***,
		"datadog": func(params output.Params) (output.Output, error) ***REMOVED***
			conf, err := datadog.GetConsolidatedConfig(params.JSONConfig, params.Environment)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			datadogc, err := datadog.New(params.Logger, conf)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return newCollectorAdapter(params, datadogc), nil
		***REMOVED***,
		"csv": func(params output.Params) (output.Output, error) ***REMOVED***
			conf, err := csv.GetConsolidatedConfig(params.JSONConfig, params.Environment, params.ConfigArgument)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			csvc, err := csv.New(params.Logger, params.FS, params.ScriptOptions.SystemTags.Map(), conf)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return newCollectorAdapter(params, csvc), nil
		***REMOVED***,
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
		res = append(res, k)
	***REMOVED***
	sort.Strings(res)
	return strings.Join(res, ", ")
***REMOVED***

func createOutputs(
	outputFullArguments []string, src *loader.SourceData, conf Config, rtOpts lib.RuntimeOptions,
	executionPlan []lib.ExecutionStep, osEnvironment map[string]string, logger logrus.FieldLogger,
) ([]output.Output, error) ***REMOVED***
	outputConstructors, err := getAllOutputConstructors()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	baseParams := output.Params***REMOVED***
		ScriptPath:     src.URL,
		Logger:         logger,
		Environment:    osEnvironment,
		StdOut:         stdout,
		StdErr:         stderr,
		FS:             afero.NewOsFs(),
		ScriptOptions:  conf.Options,
		RuntimeOptions: rtOpts,
		ExecutionPlan:  executionPlan,
	***REMOVED***
	result := make([]output.Output, 0, len(outputFullArguments))

	for _, outputFullArg := range outputFullArguments ***REMOVED***
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
		params.JSONConfig = conf.Collectors[outputType]

		output, err := outputConstructor(params)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("could not create the '%s' output: %w", outputType, err)
		***REMOVED***
		result = append(result, output)
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

// TODO: remove this after we transition every collector to the output interface

func newCollectorAdapter(params output.Params, collector lib.Collector) output.Output ***REMOVED***
	return &collectorAdapter***REMOVED***
		outputType: params.OutputType,
		collector:  collector,
		stopCh:     make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// collectorAdapter is a _temporary_ fix until we move all of the old
// "collectors" to the new output interface
type collectorAdapter struct ***REMOVED***
	collector    lib.Collector
	outputType   string
	runCtx       context.Context
	runCtxCancel func()
	stopCh       chan struct***REMOVED******REMOVED***
***REMOVED***

func (ca *collectorAdapter) Description() string ***REMOVED***
	link := ca.collector.Link()
	if link != "" ***REMOVED***
		return fmt.Sprintf("%s (%s)", ca.outputType, link)
	***REMOVED***
	return ca.outputType
***REMOVED***

func (ca *collectorAdapter) Start() error ***REMOVED***
	if err := ca.collector.Init(); err != nil ***REMOVED***
		return err
	***REMOVED***
	ca.runCtx, ca.runCtxCancel = context.WithCancel(context.Background())
	go func() ***REMOVED***
		ca.collector.Run(ca.runCtx)
		close(ca.stopCh)
	***REMOVED***()
	return nil
***REMOVED***

func (ca *collectorAdapter) AddMetricSamples(samples []stats.SampleContainer) ***REMOVED***
	ca.collector.Collect(samples)
***REMOVED***

// Stop implements the new output interface.
func (ca *collectorAdapter) Stop() error ***REMOVED***
	ca.runCtxCancel()
	<-ca.stopCh
	return nil
***REMOVED***
