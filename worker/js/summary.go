package js

import (
	_ "embed" // this is used to embed the contents of summary.js
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/dop251/goja"
)

// Copied from https://github.com/k6io/jslibWorker.k6.io/tree/master/lib/k6-summary
//go:embed summary.js
var jslibSummaryCode string //nolint:gochecknoglobals

//go:embed summary-wrapper.js
var summaryWrapperLambdaCode string //nolint:gochecknoglobals

// TODO: figure out something saner... refactor the sinks and how we deal with
// metrics in general... so much pain and misery... :sob:
func metricValueGetter(summaryTrendStats []string) func(workerMetrics.Sink, time.Duration) map[string]float64 ***REMOVED***
	trendResolvers, err := workerMetrics.GetResolversForTrendColumns(summaryTrendStats)
	if err != nil ***REMOVED***
		panic(err.Error()) // this should have been validated already
	***REMOVED***

	return func(sink workerMetrics.Sink, t time.Duration) (result map[string]float64) ***REMOVED***
		sink.Calc()

		switch sink := sink.(type) ***REMOVED***
		case *workerMetrics.CounterSink:
			result = sink.Format(t)
			rate := 0.0
			if t > 0 ***REMOVED***
				rate = sink.Value / (float64(t) / float64(time.Second))
			***REMOVED***
			result["rate"] = rate
		case *workerMetrics.GaugeSink:
			result = sink.Format(t)
			result["min"] = sink.Min
			result["max"] = sink.Max
		case *workerMetrics.RateSink:
			result = sink.Format(t)
			result["passes"] = float64(sink.Trues)
			result["fails"] = float64(sink.Total - sink.Trues)
		case *workerMetrics.TrendSink:
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
func summarizeMetricsToObject(data *libWorker.Summary, options libWorker.Options, setupData []byte) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***)
	m["root_group"] = exportGroup(data.RootGroup)
	m["options"] = map[string]interface***REMOVED******REMOVED******REMOVED***
		// TODO: improve when we can easily export all option values, including defaults?
		"summaryTrendStats": options.SummaryTrendStats,
		"summaryTimeUnit":   options.SummaryTimeUnit.String,
	***REMOVED***
	m["state"] = map[string]interface***REMOVED******REMOVED******REMOVED***
		"isStdOutTTY":       data.UIState.IsStdOutTTY,
		"isStdErrTTY":       data.UIState.IsStdErrTTY,
		"testRunDurationMs": float64(data.TestRunDuration) / float64(time.Millisecond),
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

	var setupDataI interface***REMOVED******REMOVED***
	if setupData != nil ***REMOVED***
		if err := json.Unmarshal(setupData, &setupDataI); err != nil ***REMOVED***
			// TODO: log the error
			return m
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		setupDataI = goja.Undefined()
	***REMOVED***

	m["setup_data"] = setupDataI

	return m
***REMOVED***

func exportGroup(group *libWorker.Group) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
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
