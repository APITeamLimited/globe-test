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
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/lib/types"
	"gopkg.in/guregu/null.v3"
)

func TestNewThreshold(t *testing.T) ***REMOVED***
	t.Parallel()

	src := `rate<0.01`
	abortOnFail := false
	gracePeriod := types.NullDurationFrom(2 * time.Second)

	gotThreshold := newThreshold(src, abortOnFail, gracePeriod)

	assert.Equal(t, src, gotThreshold.Source)
	assert.False(t, gotThreshold.LastFailed)
	assert.Equal(t, abortOnFail, gotThreshold.AbortOnFail)
	assert.Equal(t, gracePeriod, gotThreshold.AbortGracePeriod)
	assert.Nil(t, gotThreshold.parsed)
***REMOVED***

func TestThreshold_runNoTaint(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name             string
		parsed           *thresholdExpression
		abortGracePeriod types.NullDuration
		sinks            map[string]float64
		wantOk           bool
		wantErr          bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:             "valid expression using the > operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenGreater, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 1***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using the > operator over passing threshold and defined abort grace period",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenGreater, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(2 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 1***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using the >= operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenGreaterEqual, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.01***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using the <= operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenLessEqual, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.01***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using the < operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenLess, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.00001***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using the == operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenLooselyEqual, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.01***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using the === operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenStrictlyEqual, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.01***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression using != operator over passing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenBangEqual, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.02***REMOVED***,
			wantOk:           true,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression over failing threshold",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenGreater, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.00001***REMOVED***,
			wantOk:           false,
			wantErr:          false,
		***REMOVED***,
		***REMOVED***
			name:             "valid expression over non-existing sink",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenGreater, 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"med": 27.2***REMOVED***,
			wantOk:           false,
			wantErr:          true,
		***REMOVED***,
		***REMOVED***
			// The ParseThresholdCondition constructor should ensure that no invalid
			// operator gets through, but let's protect our future selves anyhow.
			name:             "invalid expression operator",
			parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, "&", 0.01***REMOVED***,
			abortGracePeriod: types.NullDurationFrom(0 * time.Second),
			sinks:            map[string]float64***REMOVED***"rate": 0.00001***REMOVED***,
			wantOk:           false,
			wantErr:          true,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range tests ***REMOVED***
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			threshold := &Threshold***REMOVED***
				LastFailed:       false,
				AbortOnFail:      false,
				AbortGracePeriod: testCase.abortGracePeriod,
				parsed:           testCase.parsed,
			***REMOVED***

			gotOk, gotErr := threshold.runNoTaint(testCase.sinks)

			assert.Equal(t,
				testCase.wantErr,
				gotErr != nil,
				"Threshold.runNoTaint() error = %v, wantErr %v", gotErr, testCase.wantErr,
			)

			assert.Equal(t,
				testCase.wantOk,
				gotOk,
				"Threshold.runNoTaint() gotOk = %v, want %v", gotOk, testCase.wantOk,
			)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkRunNoTaint(b *testing.B) ***REMOVED***
	threshold := &Threshold***REMOVED***
		Source:           "rate>0.01",
		LastFailed:       false,
		AbortOnFail:      false,
		AbortGracePeriod: types.NullDurationFrom(2 * time.Second),
		parsed:           &thresholdExpression***REMOVED***tokenRate, null.Float***REMOVED******REMOVED***, tokenGreater, 0.01***REMOVED***,
	***REMOVED***

	sinks := map[string]float64***REMOVED***"rate": 1***REMOVED***

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		threshold.runNoTaint(sinks) // nolint
	***REMOVED***
***REMOVED***

func TestThresholdRun(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("true", func(t *testing.T) ***REMOVED***
		t.Parallel()

		sinks := map[string]float64***REMOVED***"rate": 0.0001***REMOVED***
		parsed, parseErr := parseThresholdExpression("rate<0.01")
		require.NoError(t, parseErr)
		threshold := newThreshold(`rate<0.01`, false, types.NullDuration***REMOVED******REMOVED***)
		threshold.parsed = parsed

		t.Run("no taint", func(t *testing.T) ***REMOVED***
			b, err := threshold.runNoTaint(sinks)
			assert.NoError(t, err)
			assert.True(t, b)
			assert.False(t, threshold.LastFailed)
		***REMOVED***)

		t.Run("taint", func(t *testing.T) ***REMOVED***
			t.Parallel()

			b, err := threshold.run(sinks)
			assert.NoError(t, err)
			assert.True(t, b)
			assert.False(t, threshold.LastFailed)
		***REMOVED***)
	***REMOVED***)

	t.Run("false", func(t *testing.T) ***REMOVED***
		t.Parallel()

		sinks := map[string]float64***REMOVED***"rate": 1***REMOVED***
		parsed, parseErr := parseThresholdExpression("rate<0.01")
		require.NoError(t, parseErr)
		threshold := newThreshold(`rate<0.01`, false, types.NullDuration***REMOVED******REMOVED***)
		threshold.parsed = parsed

		t.Run("no taint", func(t *testing.T) ***REMOVED***
			b, err := threshold.runNoTaint(sinks)
			assert.NoError(t, err)
			assert.False(t, b)
			assert.False(t, threshold.LastFailed)
		***REMOVED***)

		t.Run("taint", func(t *testing.T) ***REMOVED***
			b, err := threshold.run(sinks)
			assert.NoError(t, err)
			assert.False(t, b)
			assert.True(t, threshold.LastFailed)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestThresholdsParse(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("valid threshold expressions", func(t *testing.T) ***REMOVED***
		t.Parallel()

		// Prepare a Thresholds collection containing syntaxically
		// correct thresholds
		ts := Thresholds***REMOVED***
			Thresholds: []*Threshold***REMOVED***
				newThreshold("rate<1", false, types.NullDuration***REMOVED******REMOVED***),
			***REMOVED***,
		***REMOVED***

		// Collect the result of the parsing operation
		gotErr := ts.Parse()

		assert.NoError(t, gotErr, "Parse shouldn't fail parsing valid expressions")
		assert.Condition(t, func() bool ***REMOVED***
			for _, threshold := range ts.Thresholds ***REMOVED***
				if threshold.parsed == nil ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***

			return true
		***REMOVED***, "Parse did not fail, but some Thresholds' parsed field is left empty")
	***REMOVED***)

	t.Run("invalid threshold expressions", func(t *testing.T) ***REMOVED***
		t.Parallel()

		// Prepare a Thresholds collection containing syntaxically
		// correct thresholds
		ts := Thresholds***REMOVED***
			Thresholds: []*Threshold***REMOVED***
				newThreshold("foo&1", false, types.NullDuration***REMOVED******REMOVED***),
			***REMOVED***,
		***REMOVED***

		// Collect the result of the parsing operation
		gotErr := ts.Parse()

		assert.Error(t, gotErr, "Parse should fail parsing invalid expressions")
		assert.Condition(t, func() bool ***REMOVED***
			for _, threshold := range ts.Thresholds ***REMOVED***
				if threshold.parsed == nil ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***

			return false
		***REMOVED***, "Parse failed, but some Thresholds' parsed field was not empty")
	***REMOVED***)

	t.Run("mixed valid/invalid threshold expressions", func(t *testing.T) ***REMOVED***
		t.Parallel()

		// Prepare a Thresholds collection containing syntaxically
		// correct thresholds
		ts := Thresholds***REMOVED***
			Thresholds: []*Threshold***REMOVED***
				newThreshold("rate<1", false, types.NullDuration***REMOVED******REMOVED***),
				newThreshold("foo&1", false, types.NullDuration***REMOVED******REMOVED***),
			***REMOVED***,
		***REMOVED***

		// Collect the result of the parsing operation
		gotErr := ts.Parse()

		assert.Error(t, gotErr, "Parse should fail parsing invalid expressions")
		assert.Condition(t, func() bool ***REMOVED***
			for _, threshold := range ts.Thresholds ***REMOVED***
				if threshold.parsed == nil ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***

			return false
		***REMOVED***, "Parse failed, but some Thresholds' parsed field was not empty")
	***REMOVED***)
***REMOVED***

func TestNewThresholds(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("empty", func(t *testing.T) ***REMOVED***
		t.Parallel()

		ts := NewThresholds([]string***REMOVED******REMOVED***)
		assert.Len(t, ts.Thresholds, 0)
	***REMOVED***)
	t.Run("two", func(t *testing.T) ***REMOVED***
		t.Parallel()

		sources := []string***REMOVED***`rate<0.01`, `p(95)<200`***REMOVED***
		ts := NewThresholds(sources)
		assert.Len(t, ts.Thresholds, 2)
		for i, th := range ts.Thresholds ***REMOVED***
			assert.Equal(t, sources[i], th.Source)
			assert.False(t, th.LastFailed)
			assert.False(t, th.AbortOnFail)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestNewThresholdsWithConfig(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("empty", func(t *testing.T) ***REMOVED***
		t.Parallel()

		ts := NewThresholds([]string***REMOVED******REMOVED***)
		assert.Len(t, ts.Thresholds, 0)
	***REMOVED***)
	t.Run("two", func(t *testing.T) ***REMOVED***
		t.Parallel()

		configs := []thresholdConfig***REMOVED***
			***REMOVED***`rate<0.01`, false, types.NullDuration***REMOVED******REMOVED******REMOVED***,
			***REMOVED***`p(95)<200`, true, types.NullDuration***REMOVED******REMOVED******REMOVED***,
		***REMOVED***
		ts := newThresholdsWithConfig(configs)
		assert.Len(t, ts.Thresholds, 2)
		for i, th := range ts.Thresholds ***REMOVED***
			assert.Equal(t, configs[i].Threshold, th.Source)
			assert.False(t, th.LastFailed)
			assert.Equal(t, configs[i].AbortOnFail, th.AbortOnFail)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestThresholdsRunAll(t *testing.T) ***REMOVED***
	t.Parallel()

	zero := types.NullDuration***REMOVED******REMOVED***
	oneSec := types.NullDurationFrom(time.Second)
	twoSec := types.NullDurationFrom(2 * time.Second)
	testdata := map[string]struct ***REMOVED***
		succeeded bool
		err       bool
		abort     bool
		grace     types.NullDuration
		sources   []string
	***REMOVED******REMOVED***
		"one passing":                ***REMOVED***true, false, false, zero, []string***REMOVED***`rate<0.01`***REMOVED******REMOVED***,
		"one failing":                ***REMOVED***false, false, false, zero, []string***REMOVED***`p(95)<200`***REMOVED******REMOVED***,
		"two passing":                ***REMOVED***true, false, false, zero, []string***REMOVED***`rate<0.1`, `rate<0.01`***REMOVED******REMOVED***,
		"two failing":                ***REMOVED***false, false, false, zero, []string***REMOVED***`p(95)<200`, `rate<0.1`***REMOVED******REMOVED***,
		"two mixed":                  ***REMOVED***false, false, false, zero, []string***REMOVED***`rate<0.01`, `p(95)<200`***REMOVED******REMOVED***,
		"one aborting":               ***REMOVED***false, false, true, zero, []string***REMOVED***`p(95)<200`***REMOVED******REMOVED***,
		"abort with grace period":    ***REMOVED***false, false, true, oneSec, []string***REMOVED***`p(95)<200`***REMOVED******REMOVED***,
		"no abort with grace period": ***REMOVED***false, false, true, twoSec, []string***REMOVED***`p(95)<200`***REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			thresholds := NewThresholds(data.sources)
			gotParseErr := thresholds.Parse()
			require.NoError(t, gotParseErr)
			thresholds.sinked = map[string]float64***REMOVED***"rate": 0.0001, "p(95)": 500***REMOVED***
			thresholds.Thresholds[0].AbortOnFail = data.abort
			thresholds.Thresholds[0].AbortGracePeriod = data.grace

			runDuration := 1500 * time.Millisecond

			succeeded, err := thresholds.runAll(runDuration)

			if data.err ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
			***REMOVED***

			if data.succeeded ***REMOVED***
				assert.True(t, succeeded)
			***REMOVED*** else ***REMOVED***
				assert.False(t, succeeded)
			***REMOVED***

			if data.abort && data.grace.Duration < types.Duration(runDuration) ***REMOVED***
				assert.True(t, thresholds.Abort)
			***REMOVED*** else ***REMOVED***
				assert.False(t, thresholds.Abort)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestThresholds_Run(t *testing.T) ***REMOVED***
	t.Parallel()

	type args struct ***REMOVED***
		sink     Sink
		duration time.Duration
	***REMOVED***
	tests := []struct ***REMOVED***
		name    string
		args    args
		want    bool
		wantErr bool
	***REMOVED******REMOVED***
		***REMOVED***
			"Running thresholds of existing sink",
			args***REMOVED***DummySink***REMOVED***"p(95)": 1234.5***REMOVED***, 0***REMOVED***,
			true,
			false,
		***REMOVED***,
		***REMOVED***
			"Running thresholds of existing sink but failing threshold",
			args***REMOVED***DummySink***REMOVED***"p(95)": 3000***REMOVED***, 0***REMOVED***,
			false,
			false,
		***REMOVED***,
		***REMOVED***
			"Running threshold on non existing sink fails",
			args***REMOVED***DummySink***REMOVED***"dummy": 0***REMOVED***, 0***REMOVED***,
			false,
			true,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range tests ***REMOVED***
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			thresholds := NewThresholds([]string***REMOVED***"p(95)<2000"***REMOVED***)
			gotParseErr := thresholds.Parse()
			require.NoError(t, gotParseErr)

			gotOk, gotErr := thresholds.Run(testCase.args.sink, testCase.args.duration)
			assert.Equal(t, gotErr != nil, testCase.wantErr, "Thresholds.Run() error = %v, wantErr %v", gotErr, testCase.wantErr)
			assert.Equal(t, gotOk, testCase.want, "Thresholds.Run() = %v, want %v", gotOk, testCase.want)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestThresholdsJSON(t *testing.T) ***REMOVED***
	t.Parallel()

	testdata := []struct ***REMOVED***
		JSON        string
		sources     []string
		abortOnFail bool
		gracePeriod types.NullDuration
		outputJSON  string
	***REMOVED******REMOVED***
		***REMOVED***
			`[]`,
			[]string***REMOVED******REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			`["rate<0.01"]`,
			[]string***REMOVED***"rate<0.01"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			`["rate<0.01"]`,
			[]string***REMOVED***"rate<0.01"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["rate<0.01"]`,
		***REMOVED***,
		***REMOVED***
			`["rate<0.01","p(95)<200"]`,
			[]string***REMOVED***"rate<0.01", "p(95)<200"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"rate<0.01"***REMOVED***]`,
			[]string***REMOVED***"rate<0.01"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["rate<0.01"]`,
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"rate<0.01","abortOnFail":true,"delayAbortEval":null***REMOVED***]`,
			[]string***REMOVED***"rate<0.01"***REMOVED***,
			true,
			types.NullDuration***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"rate<0.01","abortOnFail":true,"delayAbortEval":"2s"***REMOVED***]`,
			[]string***REMOVED***"rate<0.01"***REMOVED***,
			true,
			types.NullDurationFrom(2 * time.Second),
			"",
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"rate<0.01","abortOnFail":false***REMOVED***]`,
			[]string***REMOVED***"rate<0.01"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["rate<0.01"]`,
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"rate<0.01"***REMOVED***, "p(95)<200"]`,
			[]string***REMOVED***"rate<0.01", "p(95)<200"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["rate<0.01","p(95)<200"]`,
		***REMOVED***,
	***REMOVED***

	for _, data := range testdata ***REMOVED***
		data := data

		t.Run(data.JSON, func(t *testing.T) ***REMOVED***
			t.Parallel()

			var ts Thresholds
			assert.NoError(t, json.Unmarshal([]byte(data.JSON), &ts))
			assert.Equal(t, len(data.sources), len(ts.Thresholds))
			for i, src := range data.sources ***REMOVED***
				assert.Equal(t, src, ts.Thresholds[i].Source)
				assert.Equal(t, data.abortOnFail, ts.Thresholds[i].AbortOnFail)
				assert.Equal(t, data.gracePeriod, ts.Thresholds[i].AbortGracePeriod)
			***REMOVED***

			t.Run("marshal", func(t *testing.T) ***REMOVED***
				data2, err := MarshalJSONWithoutHTMLEscape(ts)
				assert.NoError(t, err)
				output := data.JSON
				if data.outputJSON != "" ***REMOVED***
					output = data.outputJSON
				***REMOVED***
				assert.Equal(t, output, string(data2))
			***REMOVED***)
		***REMOVED***)
	***REMOVED***

	t.Run("bad JSON", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var ts Thresholds
		assert.Error(t, json.Unmarshal([]byte("42"), &ts))
		assert.Nil(t, ts.Thresholds)
		assert.False(t, ts.Abort)
	***REMOVED***)

	t.Run("bad source", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var ts Thresholds
		assert.Nil(t, ts.Thresholds)
		assert.False(t, ts.Abort)
	***REMOVED***)
***REMOVED***
