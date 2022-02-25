/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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
package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.k6.io/k6/lib/types"
)

// Threshold is a representation of a single threshold for a single metric
type Threshold struct ***REMOVED***
	// Source is the text based source of the threshold
	Source string
	// LastFailed is a marker if the last testing of this threshold failed
	LastFailed bool
	// AbortOnFail marks if a given threshold fails that the whole test should be aborted
	AbortOnFail bool
	// AbortGracePeriod is a the minimum amount of time a test should be running before a failing
	// this threshold will abort the test
	AbortGracePeriod types.NullDuration
	// parsed is the threshold expression parsed from the Source
	parsed *thresholdExpression
***REMOVED***

func newThreshold(src string, abortOnFail bool, gracePeriod types.NullDuration) (*Threshold, error) ***REMOVED***
	parsedExpression, err := parseThresholdExpression(src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Threshold***REMOVED***
		Source:           src,
		AbortOnFail:      abortOnFail,
		AbortGracePeriod: gracePeriod,
		parsed:           parsedExpression,
	***REMOVED***, nil
***REMOVED***

func (t *Threshold) runNoTaint(sinks map[string]float64) (bool, error) ***REMOVED***
	// Extract the sink value for the aggregation method used in the threshold
	// expression
	lhs, ok := sinks[t.parsed.AggregationMethod]
	if !ok ***REMOVED***
		return false, fmt.Errorf("unable to apply threshold %s over metrics; reason: "+
			"no metric supporting the %s aggregation method found",
			t.Source,
			t.parsed.AggregationMethod)
	***REMOVED***

	// Apply the threshold expression operator to the left and
	// right hand side values
	var passes bool
	switch t.parsed.Operator ***REMOVED***
	case ">":
		passes = lhs > t.parsed.Value
	case ">=":
		passes = lhs >= t.parsed.Value
	case "<=":
		passes = lhs <= t.parsed.Value
	case "<":
		passes = lhs < t.parsed.Value
	case "==", "===":
		// Considering a sink always maps to float64 values,
		// strictly equal is equivalent to loosely equal
		passes = lhs == t.parsed.Value
	case "!=":
		passes = lhs != t.parsed.Value
	default:
		// The parseThresholdExpression function should ensure that no invalid
		// operator gets through, but let's protect our future selves anyhow.
		return false, fmt.Errorf("unable to apply threshold %s over metrics; "+
			"reason: %s is an invalid operator",
			t.Source,
			t.parsed.Operator,
		)
	***REMOVED***

	// Perform the actual threshold verification
	return passes, nil
***REMOVED***

func (t *Threshold) run(sinks map[string]float64) (bool, error) ***REMOVED***
	passes, err := t.runNoTaint(sinks)
	t.LastFailed = !passes
	return passes, err
***REMOVED***

type thresholdConfig struct ***REMOVED***
	Threshold        string             `json:"threshold"`
	AbortOnFail      bool               `json:"abortOnFail"`
	AbortGracePeriod types.NullDuration `json:"delayAbortEval"`
***REMOVED***

// used internally for JSON marshalling
type rawThresholdConfig thresholdConfig

func (tc *thresholdConfig) UnmarshalJSON(data []byte) error ***REMOVED***
	// shortcircuit unmarshalling for simple string format
	if err := json.Unmarshal(data, &tc.Threshold); err == nil ***REMOVED***
		return nil
	***REMOVED***

	rawConfig := (*rawThresholdConfig)(tc)
	return json.Unmarshal(data, rawConfig)
***REMOVED***

func (tc thresholdConfig) MarshalJSON() ([]byte, error) ***REMOVED***
	var data interface***REMOVED******REMOVED*** = tc.Threshold
	if tc.AbortOnFail ***REMOVED***
		data = rawThresholdConfig(tc)
	***REMOVED***

	return MarshalJSONWithoutHTMLEscape(data)
***REMOVED***

// Thresholds is the combination of all Thresholds for a given metric
type Thresholds struct ***REMOVED***
	Thresholds []*Threshold
	Abort      bool
	sinked     map[string]float64
***REMOVED***

// NewThresholds returns Thresholds objects representing the provided source strings
func NewThresholds(sources []string) (Thresholds, error) ***REMOVED***
	tcs := make([]thresholdConfig, len(sources))
	for i, source := range sources ***REMOVED***
		tcs[i].Threshold = source
	***REMOVED***

	return newThresholdsWithConfig(tcs)
***REMOVED***

func newThresholdsWithConfig(configs []thresholdConfig) (Thresholds, error) ***REMOVED***
	thresholds := make([]*Threshold, len(configs))
	sinked := make(map[string]float64)

	for i, config := range configs ***REMOVED***
		t, err := newThreshold(config.Threshold, config.AbortOnFail, config.AbortGracePeriod)
		if err != nil ***REMOVED***
			return Thresholds***REMOVED******REMOVED***, fmt.Errorf("threshold %d error: %w", i, err)
		***REMOVED***
		thresholds[i] = t
	***REMOVED***

	return Thresholds***REMOVED***thresholds, false, sinked***REMOVED***, nil
***REMOVED***

func (ts *Thresholds) runAll(timeSpentInTest time.Duration) (bool, error) ***REMOVED***
	succeeded := true
	for i, threshold := range ts.Thresholds ***REMOVED***
		b, err := threshold.run(ts.sinked)
		if err != nil ***REMOVED***
			return false, fmt.Errorf("threshold %d run error: %w", i, err)
		***REMOVED***

		if !b ***REMOVED***
			succeeded = false

			if ts.Abort || !threshold.AbortOnFail ***REMOVED***
				continue
			***REMOVED***

			ts.Abort = !threshold.AbortGracePeriod.Valid ||
				threshold.AbortGracePeriod.Duration < types.Duration(timeSpentInTest)
		***REMOVED***
	***REMOVED***

	return succeeded, nil
***REMOVED***

// Run processes all the thresholds with the provided Sink at the provided time and returns if any
// of them fails
func (ts *Thresholds) Run(sink Sink, duration time.Duration) (bool, error) ***REMOVED***
	// Initialize the sinks store
	ts.sinked = make(map[string]float64)

	// FIXME: Remove this comment as soon as the stats.Sink does not expose Format anymore.
	//
	// As of December 2021, this block reproduces the behavior of the
	// stats.Sink.Format behavior. As we intend to try to get away from it,
	// we instead implement the behavior directly here.
	//
	// For more details, see https://github.com/grafana/k6/issues/2320
	switch sinkImpl := sink.(type) ***REMOVED***
	case *CounterSink:
		ts.sinked["count"] = sinkImpl.Value
		ts.sinked["rate"] = sinkImpl.Value / (float64(duration) / float64(time.Second))
	case *GaugeSink:
		ts.sinked["value"] = sinkImpl.Value
	case *TrendSink:
		ts.sinked["min"] = sinkImpl.Min
		ts.sinked["max"] = sinkImpl.Max
		ts.sinked["avg"] = sinkImpl.Avg
		ts.sinked["med"] = sinkImpl.Med

		// Parse the percentile thresholds and insert them in
		// the sinks mapping.
		for _, threshold := range ts.Thresholds ***REMOVED***
			if !strings.HasPrefix(threshold.parsed.AggregationMethod, "p(") ***REMOVED***
				continue
			***REMOVED***

			ts.sinked[threshold.parsed.AggregationMethod] = sinkImpl.P(threshold.parsed.AggregationValue.Float64 / 100)
		***REMOVED***
	case *RateSink:
		ts.sinked["rate"] = float64(sinkImpl.Trues) / float64(sinkImpl.Total)
	case DummySink:
		for k, v := range sinkImpl ***REMOVED***
			ts.sinked[k] = v
		***REMOVED***
	default:
		return false, fmt.Errorf("unable to run Thresholds; reason: unknown sink type")
	***REMOVED***

	return ts.runAll(duration)
***REMOVED***

// UnmarshalJSON is implementation of json.Unmarshaler
func (ts *Thresholds) UnmarshalJSON(data []byte) error ***REMOVED***
	var configs []thresholdConfig
	if err := json.Unmarshal(data, &configs); err != nil ***REMOVED***
		return err
	***REMOVED***
	newts, err := newThresholdsWithConfig(configs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*ts = newts
	return nil
***REMOVED***

// MarshalJSON is implementation of json.Marshaler
func (ts Thresholds) MarshalJSON() ([]byte, error) ***REMOVED***
	configs := make([]thresholdConfig, len(ts.Thresholds))
	for i, t := range ts.Thresholds ***REMOVED***
		configs[i].Threshold = t.Source
		configs[i].AbortOnFail = t.AbortOnFail
		configs[i].AbortGracePeriod = t.AbortGracePeriod
	***REMOVED***

	return MarshalJSONWithoutHTMLEscape(configs)
***REMOVED***

// MarshalJSONWithoutHTMLEscape marshals t to JSON without escaping characters
// for safe use in HTML.
func MarshalJSONWithoutHTMLEscape(t interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	buffer := &bytes.Buffer***REMOVED******REMOVED***
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	bytes := buffer.Bytes()
	if err == nil && len(bytes) > 0 ***REMOVED***
		// Remove the newline appended by Encode() :-/
		// See https://github.com/golang/go/issues/37083
		bytes = bytes[:len(bytes)-1]
	***REMOVED***
	return bytes, err
***REMOVED***

var _ json.Unmarshaler = &Thresholds***REMOVED******REMOVED***
var _ json.Marshaler = &Thresholds***REMOVED******REMOVED***
