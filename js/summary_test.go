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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/stats"
)

const (
	checksOut = "     █ child\n\n" +
		"       ✓ check1\n" +
		"       ✗ check3\n        ↳  66% — ✓ 10 / ✗ 5\n" +
		"       ✗ check2\n        ↳  33% — ✓ 5 / ✗ 10\n\n" +
		"   ✓ checks......: 75.00% ✓ 45  ✗ 15 \n"
	countOut = "   ✗ http_reqs...: 3      3/s\n"
	gaugeOut = "     vus.........: 1      min=1 max=1\n"
	trendOut = "   ✗ my_trend....: avg=15ms min=10ms med=15ms max=20ms p(90)=19ms " +
		"p(95)=19.5ms p(99.9)=19.99ms\n"
)

func TestTextSummary(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		stats    []string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***
			[]string***REMOVED***"avg", "min", "med", "max", "p(90)", "p(95)", "p(99.9)"***REMOVED***,
			checksOut + countOut + trendOut + gaugeOut,
		***REMOVED***,
		***REMOVED***[]string***REMOVED***"count"***REMOVED***, checksOut + countOut + "   ✗ my_trend....: count=3\n" + gaugeOut***REMOVED***,
		***REMOVED***[]string***REMOVED***"avg", "count"***REMOVED***, checksOut + countOut + "   ✗ my_trend....: avg=15ms count=3\n" + gaugeOut***REMOVED***,
	***REMOVED***

	for i, tc := range testCases ***REMOVED***
		i, tc := i, tc
		t.Run(fmt.Sprintf("%d_%v", i, tc.stats), func(t *testing.T) ***REMOVED***
			t.Parallel()
			summary := createTestSummary(t)
			trendStats, err := json.Marshal(tc.stats)
			require.NoError(t, err)
			runner, err := getSimpleRunner(
				t, "/script.js",
				fmt.Sprintf(`
					exports.options = ***REMOVED***summaryTrendStats: %s***REMOVED***;
					exports.default = function() ***REMOVED***/* we don't run this, metrics are mocked */***REMOVED***;
				`, string(trendStats)),
				lib.RuntimeOptions***REMOVED***CompatibilityMode: null.NewString("base", true)***REMOVED***,
			)
			require.NoError(t, err)

			result, err := runner.HandleSummary(context.Background(), summary)
			require.NoError(t, err)

			require.Len(t, result, 1)
			stdout := result["stdout"]
			require.NotNil(t, stdout)
			summaryOut, err := ioutil.ReadAll(stdout)
			require.NoError(t, err)
			assert.Equal(t, "\n"+tc.expected+"\n", string(summaryOut))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func createTestMetrics(t *testing.T) (map[string]*stats.Metric, *lib.Group) ***REMOVED***
	metrics := make(map[string]*stats.Metric)
	gaugeMetric := stats.New("vus", stats.Gauge)
	gaugeMetric.Sink.Add(stats.Sample***REMOVED***Value: 1***REMOVED***)

	countMetric := stats.New("http_reqs", stats.Counter)
	countMetric.Tainted = null.BoolFrom(true)
	countMetric.Thresholds = stats.Thresholds***REMOVED***Thresholds: []*stats.Threshold***REMOVED******REMOVED***Source: "rate<100", LastFailed: true***REMOVED******REMOVED******REMOVED***

	checksMetric := stats.New("checks", stats.Rate)
	checksMetric.Tainted = null.BoolFrom(false)
	checksMetric.Thresholds = stats.Thresholds***REMOVED***Thresholds: []*stats.Threshold***REMOVED******REMOVED***Source: "rate>70", LastFailed: false***REMOVED******REMOVED******REMOVED***
	sink := &stats.TrendSink***REMOVED******REMOVED***

	samples := []float64***REMOVED***10.0, 15.0, 20.0***REMOVED***
	for _, s := range samples ***REMOVED***
		sink.Add(stats.Sample***REMOVED***Value: s***REMOVED***)
		countMetric.Sink.Add(stats.Sample***REMOVED***Value: 1***REMOVED***)
	***REMOVED***

	metrics["vus"] = gaugeMetric
	metrics["http_reqs"] = countMetric
	metrics["checks"] = checksMetric
	metrics["my_trend"] = &stats.Metric***REMOVED***
		Name:     "my_trend",
		Type:     stats.Trend,
		Contains: stats.Time,
		Sink:     sink,
		Tainted:  null.BoolFrom(true),
		Thresholds: stats.Thresholds***REMOVED***
			Thresholds: []*stats.Threshold***REMOVED***
				***REMOVED***
					Source:     "my_trend<1000",
					LastFailed: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	rootG, err := lib.NewGroup("", nil)
	require.NoError(t, err)
	childG, err := rootG.Group("child")
	require.NoError(t, err)
	check1, err := childG.Check("check1")
	require.NoError(t, err)
	check1.Passes = 30

	check3, err := childG.Check("check3") // intentionally before check2
	require.NoError(t, err)
	check3.Passes = 10
	check3.Fails = 5

	check2, err := childG.Check("check2")
	require.NoError(t, err)
	check2.Passes = 5
	check2.Fails = 10

	for i := 0; i < int(check1.Passes+check2.Passes+check3.Passes); i++ ***REMOVED***
		checksMetric.Sink.Add(stats.Sample***REMOVED***Value: 1***REMOVED***)
	***REMOVED***
	for i := 0; i < int(check1.Fails+check2.Fails+check3.Fails); i++ ***REMOVED***
		checksMetric.Sink.Add(stats.Sample***REMOVED***Value: 0***REMOVED***)
	***REMOVED***

	return metrics, rootG
***REMOVED***

func createTestSummary(t *testing.T) *lib.Summary ***REMOVED***
	metrics, rootG := createTestMetrics(t)
	return &lib.Summary***REMOVED***
		Metrics:         metrics,
		RootGroup:       rootG,
		TestRunDuration: time.Second,
	***REMOVED***
***REMOVED***

const expectedOldJSONExportResult = `***REMOVED***
    "root_group": ***REMOVED***
        "name": "",
        "path": "",
        "id": "d41d8cd98f00b204e9800998ecf8427e",
        "groups": ***REMOVED***
            "child": ***REMOVED***
                "name": "child",
                "path": "::child",
                "id": "f41cbb53a398ec1c9fb3d33e20c9b040",
                "groups": ***REMOVED******REMOVED***,
                "checks": ***REMOVED***
                    "check1": ***REMOVED***
                        "name": "check1",
                        "path": "::child::check1",
                        "id": "6289a7a06253a1c3f6137dfb25695563",
                        "passes":30,
                        "fails": 0
                    ***REMOVED***,
                    "check2": ***REMOVED***
                        "name": "check2",
                        "path": "::child::check2",
                        "id": "06f5922794bef0d4584ba76a49893e1f",
                        "passes": 5,
                        "fails": 10
                    ***REMOVED***,
                    "check3": ***REMOVED***
                        "name": "check3",
                        "path": "::child::check3",
                        "id": "c7553eca92d3e034b5808332296d304a",
                        "passes": 10,
                        "fails": 5
                    ***REMOVED***
                ***REMOVED***
            ***REMOVED***
        ***REMOVED***,
        "checks": ***REMOVED******REMOVED***
    ***REMOVED***,
    "metrics": ***REMOVED***
        "checks": ***REMOVED***
            "value": 0.75,
            "passes": 45,
            "fails": 15,
            "thresholds": ***REMOVED***
                "rate>70": false
            ***REMOVED***
        ***REMOVED***,
        "http_reqs": ***REMOVED***
            "count": 3,
            "rate": 3,
            "thresholds": ***REMOVED***
                "rate<100": true
            ***REMOVED***
        ***REMOVED***,
        "my_trend": ***REMOVED***
            "avg": 15,
            "max": 20,
            "med": 15,
            "min": 10,
            "p(90)": 19,
            "p(95)": 19.5,
            "p(99)": 19.9,
			"count": 3,
            "thresholds": ***REMOVED***
                "my_trend<1000": true
            ***REMOVED***
        ***REMOVED***,
        "vus": ***REMOVED***
            "value": 1,
            "min": 1,
            "max": 1
        ***REMOVED***
    ***REMOVED***
***REMOVED***
`

func TestOldJSONExport(t *testing.T) ***REMOVED***
	t.Parallel()
	runner, err := getSimpleRunner(
		t, "/script.js",
		`
		exports.options = ***REMOVED***summaryTrendStats: ["avg", "min", "med", "max", "p(90)", "p(95)", "p(99)", "count"]***REMOVED***;
		exports.default = function() ***REMOVED***/* we don't run this, metrics are mocked */***REMOVED***;
		`,
		lib.RuntimeOptions***REMOVED***
			CompatibilityMode: null.NewString("base", true),
			SummaryExport:     null.StringFrom("result.json"),
		***REMOVED***,
	)

	require.NoError(t, err)

	summary := createTestSummary(t)
	result, err := runner.HandleSummary(context.Background(), summary)
	require.NoError(t, err)

	require.Len(t, result, 2)
	require.NotNil(t, result["stdout"])
	textSummary, err := ioutil.ReadAll(result["stdout"])
	require.NoError(t, err)
	assert.Contains(t, string(textSummary), checksOut+countOut)
	require.NotNil(t, result["result.json"])
	jsonExport, err := ioutil.ReadAll(result["result.json"])
	require.NoError(t, err)
	assert.JSONEq(t, expectedOldJSONExportResult, string(jsonExport))
***REMOVED***

const expectedHandleSummaryRawData = `
***REMOVED***
    "root_group": ***REMOVED***
        "groups": [
            ***REMOVED***
                "name": "child",
                "path": "::child",
                "id": "f41cbb53a398ec1c9fb3d33e20c9b040",
                "groups": [],
                "checks": [
                        ***REMOVED***
                            "id": "6289a7a06253a1c3f6137dfb25695563",
                            "passes": 30,
                            "fails": 0,
                            "name": "check1",
                            "path": "::child::check1"
                        ***REMOVED***,
                        ***REMOVED***
                            "fails": 5,
                            "name": "check3",
                            "path": "::child::check3",
                            "id": "c7553eca92d3e034b5808332296d304a",
                            "passes": 10
                        ***REMOVED***,
                        ***REMOVED***
                            "name": "check2",
                            "path": "::child::check2",
                            "id": "06f5922794bef0d4584ba76a49893e1f",
                            "passes": 5,
                            "fails": 10
                        ***REMOVED***
                    ]
            ***REMOVED***
        ],
        "checks": [],
        "name": "",
        "path": "",
        "id": "d41d8cd98f00b204e9800998ecf8427e"
    ***REMOVED***,
    "options": ***REMOVED***
        "summaryTrendStats": [
            "avg",
            "min",
            "med",
            "max",
            "p(90)",
            "p(95)",
            "p(99)",
            "count"
        ],
        "summaryTimeUnit": "",
		"noColor": false
    ***REMOVED***,
	"state": ***REMOVED***
		"isStdErrTTY": false,
		"isStdOutTTY": false
	***REMOVED***,
    "metrics": ***REMOVED***
        "checks": ***REMOVED***
            "contains": "default",
            "values": ***REMOVED***
                "passes": 45,
                "fails": 15,
                "rate": 0.75
            ***REMOVED***,
            "type": "rate",
            "thresholds": ***REMOVED***
                "rate>70": ***REMOVED***
                    "ok": true
                ***REMOVED***
            ***REMOVED***
        ***REMOVED***,
        "my_trend": ***REMOVED***
            "thresholds": ***REMOVED***
                "my_trend<1000": ***REMOVED***
                    "ok": false
                ***REMOVED***
            ***REMOVED***,
            "type": "trend",
            "contains": "time",
            "values": ***REMOVED***
                "max": 20,
                "p(90)": 19,
                "p(95)": 19.5,
                "p(99)": 19.9,
                "count": 3,
                "avg": 15,
                "min": 10,
                "med": 15
            ***REMOVED***
        ***REMOVED***,
        "vus": ***REMOVED***
            "contains": "default",
            "values": ***REMOVED***
                "value": 1,
                "min": 1,
                "max": 1
            ***REMOVED***,
            "type": "gauge"
        ***REMOVED***,
        "http_reqs": ***REMOVED***
            "type": "counter",
            "contains": "default",
            "values": ***REMOVED***
                "count": 3,
                "rate": 3
            ***REMOVED***,
            "thresholds": ***REMOVED***
                "rate<100": ***REMOVED***
                    "ok": false
                ***REMOVED***
            ***REMOVED***
        ***REMOVED***
    ***REMOVED***
***REMOVED***`

func TestRawHandleSummaryData(t *testing.T) ***REMOVED***
	t.Parallel()
	runner, err := getSimpleRunner(
		t, "/script.js",
		`
		exports.options = ***REMOVED***summaryTrendStats: ["avg", "min", "med", "max", "p(90)", "p(95)", "p(99)", "count"]***REMOVED***;
		exports.default = function() ***REMOVED*** /* we don't run this, metrics are mocked */ ***REMOVED***;
		exports.handleSummary = function(data) ***REMOVED***
			return ***REMOVED***'rawdata.json': JSON.stringify(data)***REMOVED***;
		***REMOVED***;
		`,
		lib.RuntimeOptions***REMOVED***
			CompatibilityMode: null.NewString("base", true),
			// we still want to check this
			SummaryExport: null.StringFrom("old-export.json"),
		***REMOVED***,
	)

	require.NoError(t, err)

	summary := createTestSummary(t)
	result, err := runner.HandleSummary(context.Background(), summary)
	require.NoError(t, err)

	require.Len(t, result, 2)
	require.Nil(t, result["stdout"])

	require.NotNil(t, result["old-export.json"])
	oldExport, err := ioutil.ReadAll(result["old-export.json"])
	require.NoError(t, err)
	assert.JSONEq(t, expectedOldJSONExportResult, string(oldExport))
	require.NotNil(t, result["rawdata.json"])
	newRawData, err := ioutil.ReadAll(result["rawdata.json"])
	require.NoError(t, err)
	assert.JSONEq(t, expectedHandleSummaryRawData, string(newRawData))
***REMOVED***

func TestWrongSummaryHandlerExportTypes(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []string***REMOVED***"***REMOVED******REMOVED***", `"foo"`, "null", "undefined", "123"***REMOVED***

	for i, tc := range testCases ***REMOVED***
		i, tc := i, tc
		t.Run(fmt.Sprintf("%d_%s", i, tc), func(t *testing.T) ***REMOVED***
			t.Parallel()
			runner, err := getSimpleRunner(t, "/script.js",
				fmt.Sprintf(`
					exports.default = function() ***REMOVED*** /* we don't run this, metrics are mocked */ ***REMOVED***;
					exports.handleSummary = %s;
				`, tc),
				lib.RuntimeOptions***REMOVED***CompatibilityMode: null.NewString("base", true)***REMOVED***,
			)
			require.NoError(t, err)

			summary := createTestSummary(t)
			_, err = runner.HandleSummary(context.Background(), summary)
			require.Error(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExceptionInHandleSummaryFallsBackToTextSummary(t *testing.T) ***REMOVED***
	t.Parallel()

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.ErrorLevel***REMOVED******REMOVED***
	logger.AddHook(&logHook)

	runner, err := getSimpleRunner(t, "/script.js", `
			exports.default = function() ***REMOVED***/* we don't run this, metrics are mocked */***REMOVED***;
			exports.handleSummary = function(data) ***REMOVED***
				throw new Error('intentional error');
			***REMOVED***;
		`, logger, lib.RuntimeOptions***REMOVED***CompatibilityMode: null.NewString("base", true)***REMOVED***,
	)

	require.NoError(t, err)

	summary := createTestSummary(t)
	result, err := runner.HandleSummary(context.Background(), summary)
	require.NoError(t, err)

	require.Len(t, result, 1)
	require.NotNil(t, result["stdout"])
	textSummary, err := ioutil.ReadAll(result["stdout"])
	require.NoError(t, err)
	assert.Contains(t, string(textSummary), checksOut+countOut)

	logErrors := logHook.Drain()
	assert.Equal(t, 1, len(logErrors))
	errMsg, err := logErrors[0].String()
	require.NoError(t, err)
	assert.Contains(t, errMsg, "intentional error")
***REMOVED***
