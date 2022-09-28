package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/APITeamLimited/k6-worker/errext"
	"github.com/APITeamLimited/k6-worker/errext/exitcodes"
	"github.com/APITeamLimited/k6-worker/lib/types"
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

func newThreshold(src string, abortOnFail bool, gracePeriod types.NullDuration) *Threshold ***REMOVED***
	return &Threshold***REMOVED***
		Source:           src,
		AbortOnFail:      abortOnFail,
		AbortGracePeriod: gracePeriod,
		parsed:           nil,
	***REMOVED***
***REMOVED***

func (t *Threshold) runNoTaint(sinks map[string]float64) (bool, error) ***REMOVED***
	// Extract the sink value for the aggregation method used in the threshold
	// expression. Considering we already validated thresholds before starting
	// the execution, we assume that a missing sink entry means that no samples
	// are available yet, and that it's safe to ignore this run.
	lhs, ok := sinks[t.parsed.SinkKey()]
	if !ok ***REMOVED***
		return true, nil
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
func NewThresholds(sources []string) Thresholds ***REMOVED***
	tcs := make([]thresholdConfig, len(sources))
	for i, source := range sources ***REMOVED***
		tcs[i].Threshold = source
	***REMOVED***

	return newThresholdsWithConfig(tcs)
***REMOVED***

func newThresholdsWithConfig(configs []thresholdConfig) Thresholds ***REMOVED***
	thresholds := make([]*Threshold, len(configs))
	sinked := make(map[string]float64)

	for i, config := range configs ***REMOVED***
		t := newThreshold(config.Threshold, config.AbortOnFail, config.AbortGracePeriod)
		thresholds[i] = t
	***REMOVED***

	return Thresholds***REMOVED***thresholds, false, sinked***REMOVED***
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

	// FIXME: Remove this comment as soon as the metrics.Sink does not expose Format anymore.
	//
	// As of December 2021, this block reproduces the behavior of the
	// metrics.Sink.Format behavior. As we intend to try to get away from it,
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
			if threshold.parsed.AggregationMethod != tokenPercentile ***REMOVED***
				continue
			***REMOVED***

			key := fmt.Sprintf("p(%g)", threshold.parsed.AggregationValue.Float64)
			ts.sinked[key] = sinkImpl.P(threshold.parsed.AggregationValue.Float64 / 100)
		***REMOVED***
	case *RateSink:
		// We want to avoid division by zero, which
		// would lead to [#2520](https://github.com/grafana/k6/issues/2520)
		if sinkImpl.Total > 0 ***REMOVED***
			ts.sinked["rate"] = float64(sinkImpl.Trues) / float64(sinkImpl.Total)
		***REMOVED***
	case DummySink:
		for k, v := range sinkImpl ***REMOVED***
			ts.sinked[k] = v
		***REMOVED***
	default:
		return false, fmt.Errorf("unable to run Thresholds; reason: unknown sink type")
	***REMOVED***

	return ts.runAll(duration)
***REMOVED***

// Parse parses the Thresholds and fills each Threshold.parsed field with the result.
// It effectively asserts they are syntaxically correct.
func (ts *Thresholds) Parse() error ***REMOVED***
	for _, t := range ts.Thresholds ***REMOVED***
		parsed, err := parseThresholdExpression(t.Source)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		t.parsed = parsed
	***REMOVED***

	return nil
***REMOVED***

// ErrInvalidThreshold indicates a threshold is not valid
var ErrInvalidThreshold = errors.New("invalid threshold")

// Validate ensures a threshold definition is consistent with the metric it applies to.
// Given a metric registry and a metric name to apply the expressions too, Validate will
// assert that each threshold expression uses an aggregation method that's supported by the
// provided metric. It returns an error otherwise.
// Note that this function expects the passed in thresholds to have been parsed already, and
// have their Parsed (ThresholdExpression) field already filled.
func (ts *Thresholds) Validate(metricName string, r *Registry) error ***REMOVED***
	parsedMetricName, _, err := ParseMetricName(metricName)
	if err != nil ***REMOVED***
		parseErr := fmt.Errorf("unable to validate threshold expressions; reason: %w", err)
		return errext.WithExitCodeIfNone(parseErr, exitcodes.InvalidConfig)
	***REMOVED***

	// Obtain the metric the thresholds apply to from the registry.
	// if the metric doesn't exist, then we return an error indicating
	// the InvalidConfig exitcode should be used.
	metric := r.Get(parsedMetricName)
	if metric == nil ***REMOVED***
		err := fmt.Errorf("%w defined on %s; reason: no metric name %q found", ErrInvalidThreshold, metricName, metricName)
		return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
	***REMOVED***

	for _, threshold := range ts.Thresholds ***REMOVED***
		// Return a digestable error if we attempt to validate a threshold
		// that hasn't been parsed yet.
		if threshold.parsed == nil ***REMOVED***
			thresholdExpression, err := parseThresholdExpression(threshold.Source)
			if err != nil ***REMOVED***
				return fmt.Errorf("unable to validate threshold %q on metric %s; reason: "+
					"parsing threshold failed %w", threshold.Source, metricName, err)
			***REMOVED***

			threshold.parsed = thresholdExpression
		***REMOVED***

		// If the threshold's expression aggregation method is not
		// supported for the metric we validate against, then we return
		// an error indicating the InvalidConfig exitcode should be used.
		if !metric.Type.supportsAggregationMethod(threshold.parsed.AggregationMethod) ***REMOVED***
			err := fmt.Errorf(
				"%w %q applied on metric %s; reason: "+
					"unsupported aggregation method %s on metric of type %s. "+
					"supported aggregation methods for this metric are: %s",
				ErrInvalidThreshold, threshold.Source, metricName,
				threshold.parsed.AggregationMethod, metric.Type,
				strings.Join(metric.Type.supportedAggregationMethods(), ", "),
			)
			return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// UnmarshalJSON is implementation of json.Unmarshaler
func (ts *Thresholds) UnmarshalJSON(data []byte) error ***REMOVED***
	var configs []thresholdConfig
	if err := json.Unmarshal(data, &configs); err != nil ***REMOVED***
		return err
	***REMOVED***

	*ts = newThresholdsWithConfig(configs)

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

var (
	_ json.Unmarshaler = &Thresholds***REMOVED******REMOVED***
	_ json.Marshaler   = &Thresholds***REMOVED******REMOVED***
)
