/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

package ui

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

func TestSummary(t *testing.T) ***REMOVED***
	t.Run("SummarizeMetrics", func(t *testing.T) ***REMOVED***
		var (
			checksOut = "     █ child\n\n" +
				"       ✗ check1\n        ↳  33% — ✓ 5 / ✗ 10\n" +
				"       ✗ check2\n        ↳  66% — ✓ 10 / ✗ 5\n\n" +
				"   ✓ checks......: 100.00% ✓ 3   ✗ 0  \n"
			countOut = "   ✗ http_reqs...: 3       3/s\n"
			gaugeOut = "     vus.........: 1       min=1 max=1\n"
			trendOut = "   ✗ my_trend....: avg=15ms min=10ms med=15ms max=20ms p(90)=19ms " +
				"p(95)=19.5ms p(99.9)=19.99ms\n"
		)

		metrics := createTestMetrics()
		testCases := []struct ***REMOVED***
			stats    []string
			expected string
		***REMOVED******REMOVED***
			***REMOVED***[]string***REMOVED***"avg", "min", "med", "max", "p(90)", "p(95)", "p(99.9)"***REMOVED***,
				checksOut + countOut + trendOut + gaugeOut***REMOVED***,
			***REMOVED***[]string***REMOVED***"count"***REMOVED***, checksOut + countOut + "   ✗ my_trend....: count=3\n" + gaugeOut***REMOVED***,
			***REMOVED***[]string***REMOVED***"avg", "count"***REMOVED***, checksOut + countOut + "   ✗ my_trend....: avg=15ms count=3\n" + gaugeOut***REMOVED***,
		***REMOVED***

		rootG, _ := lib.NewGroup("", nil)
		childG, _ := rootG.Group("child")
		check1, _ := childG.Check("check1")
		check1.Passes = 5
		check1.Fails = 10
		check2, _ := childG.Check("check2")
		check2.Passes = 10
		check2.Fails = 5
		for _, tc := range testCases ***REMOVED***
			tc := tc
			t.Run(fmt.Sprintf("%v", tc.stats), func(t *testing.T) ***REMOVED***
				var w bytes.Buffer
				s := NewSummary(tc.stats)

				s.SummarizeMetrics(&w, " ", SummaryData***REMOVED***
					Metrics:   metrics,
					RootGroup: rootG,
					Time:      time.Second,
					TimeUnit:  "",
				***REMOVED***)
				assert.Equal(t, tc.expected, w.String())
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	t.Run("generateCustomTrendValueResolvers", func(t *testing.T) ***REMOVED***
		var customResolversTests = []struct ***REMOVED***
			stats      []string
			percentile float64
		***REMOVED******REMOVED***
			***REMOVED***[]string***REMOVED***"p(99)", "p(err)"***REMOVED***, 0.99***REMOVED***,
			***REMOVED***[]string***REMOVED***"p(none", "p(99.9)"***REMOVED***, 0.9990000000000001***REMOVED***,
			***REMOVED***[]string***REMOVED***"p(none", "p(99.99)"***REMOVED***, 0.9998999999999999***REMOVED***,
			***REMOVED***[]string***REMOVED***"p(none", "p(99.999)"***REMOVED***, 0.9999899999999999***REMOVED***,
		***REMOVED***

		sink := createTestTrendSink(100)

		for _, tc := range customResolversTests ***REMOVED***
			tc := tc
			t.Run(fmt.Sprintf("%v", tc.stats), func(t *testing.T) ***REMOVED***
				s := Summary***REMOVED***trendColumns: tc.stats***REMOVED***
				res := s.generateCustomTrendValueResolvers(tc.stats)
				assert.Len(t, res, 1)
				for k := range res ***REMOVED***
					assert.Equal(t, sink.P(tc.percentile), res[k](sink))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestValidateSummary(t *testing.T) ***REMOVED***
	var validateTests = []struct ***REMOVED***
		stats  []string
		expErr error
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED******REMOVED***, nil***REMOVED***,
		***REMOVED***[]string***REMOVED***"avg", "min", "med", "max", "p(0)", "p(99)", "p(99.999)", "count"***REMOVED***, nil***REMOVED***,
		***REMOVED***[]string***REMOVED***"avg", "p(err)"***REMOVED***, ErrInvalidStat***REMOVED***"p(err)", errPercentileStatInvalidValue***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"nil", "p(err)"***REMOVED***, ErrInvalidStat***REMOVED***"nil", errStatUnknownFormat***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"p90"***REMOVED***, ErrInvalidStat***REMOVED***"p90", errStatUnknownFormat***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"p(90"***REMOVED***, ErrInvalidStat***REMOVED***"p(90", errStatUnknownFormat***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***" avg"***REMOVED***, ErrInvalidStat***REMOVED***" avg", errStatUnknownFormat***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"avg "***REMOVED***, ErrInvalidStat***REMOVED***"avg ", errStatUnknownFormat***REMOVED******REMOVED***,
		***REMOVED***[]string***REMOVED***"", "avg "***REMOVED***, ErrInvalidStat***REMOVED***"", errStatEmptyString***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range validateTests ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%v", tc.stats), func(t *testing.T) ***REMOVED***
			err := ValidateSummary(tc.stats)
			assert.Equal(t, tc.expErr, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func createTestTrendSink(count int) *stats.TrendSink ***REMOVED***
	sink := stats.TrendSink***REMOVED******REMOVED***

	for i := 0; i < count; i++ ***REMOVED***
		sink.Add(stats.Sample***REMOVED***Value: float64(i)***REMOVED***)
	***REMOVED***

	return &sink
***REMOVED***

func createTestMetrics() map[string]*stats.Metric ***REMOVED***
	metrics := make(map[string]*stats.Metric)
	gaugeMetric := stats.New("vus", stats.Gauge)
	gaugeMetric.Sink.Add(stats.Sample***REMOVED***Value: 1***REMOVED***)

	countMetric := stats.New("http_reqs", stats.Counter)
	countMetric.Tainted = null.BoolFrom(true)
	countMetric.Thresholds = stats.Thresholds***REMOVED***Thresholds: []*stats.Threshold***REMOVED******REMOVED***Source: "rate<100"***REMOVED******REMOVED******REMOVED***

	checksMetric := stats.New("checks", stats.Rate)
	checksMetric.Tainted = null.BoolFrom(false)
	sink := &stats.TrendSink***REMOVED******REMOVED***

	samples := []float64***REMOVED***10.0, 15.0, 20.0***REMOVED***
	for _, s := range samples ***REMOVED***
		sink.Add(stats.Sample***REMOVED***Value: s***REMOVED***)
		checksMetric.Sink.Add(stats.Sample***REMOVED***Value: 1***REMOVED***)
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

	return metrics
***REMOVED***

func TestSummarizeMetricsJSON(t *testing.T) ***REMOVED***
	metrics := createTestMetrics()
	expected := `***REMOVED***
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
                        "passes": 5,
                        "fails": 10
                    ***REMOVED***
                ***REMOVED***
            ***REMOVED***
        ***REMOVED***,
        "checks": ***REMOVED******REMOVED***
    ***REMOVED***,
    "metrics": ***REMOVED***
        "checks": ***REMOVED***
            "value": 0,
            "passes": 3,
            "fails": 0
        ***REMOVED***,
        "http_reqs": ***REMOVED***
            "count": 3,
            "rate": 3,
            "thresholds": ***REMOVED***
                "rate<100": false
            ***REMOVED***
        ***REMOVED***,
        "my_trend": ***REMOVED***
            "avg": 15,
            "max": 20,
            "med": 15,
            "min": 10,
            "p(90)": 19,
            "p(95)": 19.5,
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
	rootG, _ := lib.NewGroup("", nil)
	childG, _ := rootG.Group("child")
	check, _ := lib.NewCheck("check1", childG)
	check.Passes = 5
	check.Fails = 10
	childG.Checks["check1"] = check

	s := NewSummary([]string***REMOVED***"avg", "min", "med", "max", "p(90)", "p(95)", "p(99.9)"***REMOVED***)
	data := SummaryData***REMOVED***
		Metrics:   metrics,
		RootGroup: rootG,
		Time:      time.Second,
		TimeUnit:  "",
	***REMOVED***

	var w bytes.Buffer
	err := s.SummarizeMetricsJSON(&w, data)
	require.Nil(t, err)
	require.Contains(t, w.String(), "<")
	require.JSONEq(t, expected, w.String())
***REMOVED***
