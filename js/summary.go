/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

package js

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/ui"
)

// TODO: move this to a separate JS file and use go.rice to embed it
const summaryWrapperLambdaCode = `
(function() ***REMOVED***
	var forEach = function(obj, callback) ***REMOVED***
		for (var key in obj) ***REMOVED***
			if (obj.hasOwnProperty(key)) ***REMOVED***
				callback(key, obj[key]);
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var transformGroup = function(group) ***REMOVED***
		if (Array.isArray(group.groups)) ***REMOVED***
			var newFormatGroups = group.groups;
			group.groups = ***REMOVED******REMOVED***;
			for (var i = 0; i < newFormatGroups.length; i++) ***REMOVED***
				group.groups[newFormatGroups[i].name] = transformGroup(newFormatGroups[i]);
			***REMOVED***
		***REMOVED***
		if (Array.isArray(group.checks)) ***REMOVED***
			var newFormatChecks = group.checks;
			group.checks = ***REMOVED******REMOVED***;
			for (var i = 0; i < newFormatChecks.length; i++) ***REMOVED***
				group.checks[newFormatChecks[i].name] = newFormatChecks[i];
			***REMOVED***
		***REMOVED***
		return group;
	***REMOVED***;

	var oldJSONSummary = function(data) ***REMOVED***
		// Quick copy of the data, since it's easiest to modify it in place.
		var results = JSON.parse(JSON.stringify(data));
		delete results.options;

		forEach(results.metrics, function(metricName, metric) ***REMOVED***
			var oldFormatMetric = metric.values;
			if (metric.thresholds && Object.keys(metric.thresholds).length > 0) ***REMOVED***
				var newFormatThresholds = metric.thresholds;
				oldFormatMetric.thresholds = ***REMOVED******REMOVED***;
				forEach(newFormatThresholds, function(thresholdName, threshold) ***REMOVED***
					oldFormatMetric.thresholds[thresholdName] = !threshold.ok;
				***REMOVED***);
			***REMOVED***
			if (metric.type == 'rate' && oldFormatMetric.hasOwnProperty('rate')) ***REMOVED***
				oldFormatMetric.value = oldFormatMetric.rate; // sigh...
				delete oldFormatMetric.rate;
			***REMOVED***
			results.metrics[metricName] = oldFormatMetric;
		***REMOVED***);

		results.root_group = transformGroup(results.root_group);

		return JSON.stringify(results, null, 4);
	***REMOVED***;

	var oldTextSummary = function(data) ***REMOVED***
		// TODO: implement something like the current end of test summary
	***REMOVED***;

	return function(exportedSummaryCallback, jsonSummaryPath, data, oldCallback) ***REMOVED***
		var result = ***REMOVED******REMOVED***;
		if (exportedSummaryCallback) ***REMOVED***
			try ***REMOVED***
				result = exportedSummaryCallback(data, oldCallback);
			***REMOVED*** catch (e) ***REMOVED***
				console.error('handleSummary() failed with error "' + e + '", falling back to the default summary');
				//result["stdout"] = oldTextSummary(data);
				result["stdout"] = oldCallback(); // TODO: delete
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// result["stdout"] = oldTextSummary(data);
			result["stdout"] = oldCallback(); // TODO: delete
		***REMOVED***

		// TODO: ensure we're returning a map of strings or null/undefined...
		// and if not, log an error and generate the default summary?

		if (jsonSummaryPath != '') ***REMOVED***
			result[jsonSummaryPath] = oldJSONSummary(data);
		***REMOVED***

		return result;
	***REMOVED***;
***REMOVED***)();
`

// TODO: figure out something saner... refactor the sinks and how we deal with
// metrics in general... so much pain and misery... :sob:
func metricValueGetter(summaryTrendStats []string) func(stats.Sink, time.Duration) map[string]float64 ***REMOVED***
	trendResolvers, err := stats.GetResolversForTrendColumns(summaryTrendStats)
	if err != nil ***REMOVED***
		panic(err.Error()) // this should have been validated already
	***REMOVED***

	return func(sink stats.Sink, t time.Duration) (result map[string]float64) ***REMOVED***
		sink.Calc()

		switch sink := sink.(type) ***REMOVED***
		case *stats.CounterSink:
			result = sink.Format(t)
			rate := 0.0
			if t > 0 ***REMOVED***
				rate = sink.Value / (float64(t) / float64(time.Second))
			***REMOVED***
			result["rate"] = rate
		case *stats.GaugeSink:
			result = sink.Format(t)
			result["min"] = sink.Min
			result["max"] = sink.Max
		case *stats.RateSink:
			result = sink.Format(t)
			result["passes"] = float64(sink.Trues)
			result["fails"] = float64(sink.Total - sink.Trues)
		case *stats.TrendSink:
			result = make(map[string]float64, len(summaryTrendStats))
			for _, col := range summaryTrendStats ***REMOVED***
				result[col] = trendResolvers[col](sink)
			***REMOVED***
		***REMOVED***

		return result
	***REMOVED***
***REMOVED***

// summarizeMetricsToObject transforms the summary objects in a way that's
// suitable to pass to the JS runtime or export to JSON.
func summarizeMetricsToObject(data *lib.Summary, options lib.Options) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***)
	m["root_group"] = exportGroup(data.RootGroup)
	m["options"] = map[string]interface***REMOVED******REMOVED******REMOVED***
		// TODO: improve when we can easily export all option values, including defaults?
		"summaryTrendStats": options.SummaryTrendStats,
		"summaryTimeUnit":   options.SummaryTimeUnit.String,
	***REMOVED***

	getMetricValues := metricValueGetter(options.SummaryTrendStats)

	metricsData := make(map[string]interface***REMOVED******REMOVED***)
	for name, m := range data.Metrics ***REMOVED***
		metricData := map[string]interface***REMOVED******REMOVED******REMOVED***
			"type":     m.Type.String(),
			"contains": m.Contains.String(),
			"values":   getMetricValues(m.Sink, data.TestRunDuration),
		***REMOVED***

		if len(m.Thresholds.Thresholds) > 0 ***REMOVED***
			thresholds := make(map[string]interface***REMOVED******REMOVED***)
			for _, threshold := range m.Thresholds.Thresholds ***REMOVED***
				thresholds[threshold.Source] = map[string]interface***REMOVED******REMOVED******REMOVED***
					"ok": !threshold.LastFailed,
				***REMOVED***
			***REMOVED***
			metricData["thresholds"] = thresholds
		***REMOVED***
		metricsData[name] = metricData
	***REMOVED***
	m["metrics"] = metricsData

	return m
***REMOVED***

func exportGroup(group *lib.Group) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	subGroups := make([]map[string]interface***REMOVED******REMOVED***, len(group.OrderedGroups))
	for i, subGroup := range group.OrderedGroups ***REMOVED***
		subGroups[i] = exportGroup(subGroup)
	***REMOVED***

	checks := make([]map[string]interface***REMOVED******REMOVED***, len(group.OrderedChecks))
	for i, check := range group.OrderedChecks ***REMOVED***
		checks[i] = map[string]interface***REMOVED******REMOVED******REMOVED***
			"name":   check.Name,
			"path":   check.Path,
			"id":     check.ID,
			"passes": check.Passes,
			"fails":  check.Fails,
		***REMOVED***
	***REMOVED***

	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"name":   group.Name,
		"path":   group.Path,
		"id":     group.ID,
		"groups": subGroups,
		"checks": checks,
	***REMOVED***
***REMOVED***

func getSummaryResult(rawResult goja.Value) (map[string]io.Reader, error) ***REMOVED***
	if goja.IsNull(rawResult) || goja.IsUndefined(rawResult) ***REMOVED***
		return nil, nil
	***REMOVED***

	rawResultMap, ok := rawResult.Export().(map[string]interface***REMOVED******REMOVED***)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("handleSummary() should return a map with string keys")
	***REMOVED***

	results := make(map[string]io.Reader, len(rawResultMap))
	for path, val := range rawResultMap ***REMOVED***
		readerVal, err := common.GetReader(val)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("error handling summary object %s: %w", path, err)
		***REMOVED***
		results[path] = readerVal
	***REMOVED***

	return results, nil
***REMOVED***

// TODO: remove this after the JS alternative is written
func getOldTextSummaryFunc(summary *lib.Summary, options lib.Options) func() string ***REMOVED***
	data := ui.SummaryData***REMOVED***
		Metrics:   summary.Metrics,
		RootGroup: summary.RootGroup,
		Time:      summary.TestRunDuration,
		TimeUnit:  options.SummaryTimeUnit.String,
	***REMOVED***

	return func() string ***REMOVED***
		buffer := bytes.NewBuffer(nil)
		_ = buffer.WriteByte('\n')

		s := ui.NewSummary(options.SummaryTrendStats)
		s.SummarizeMetrics(buffer, " ", data)

		_ = buffer.WriteByte('\n')

		return buffer.String()
	***REMOVED***
***REMOVED***
