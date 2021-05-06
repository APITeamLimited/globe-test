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

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

	"github.com/loadimpact/k6/lib/types"
)

func TestNewThreshold(t *testing.T) ***REMOVED***
	src := `1+1==2`
	rt := goja.New()
	abortOnFail := false
	gracePeriod := types.NullDurationFrom(2 * time.Second)
	th, err := newThreshold(src, rt, abortOnFail, gracePeriod)
	assert.NoError(t, err)

	assert.Equal(t, src, th.Source)
	assert.False(t, th.LastFailed)
	assert.NotNil(t, th.pgm)
	assert.Equal(t, rt, th.rt)
	assert.Equal(t, abortOnFail, th.AbortOnFail)
	assert.Equal(t, gracePeriod, th.AbortGracePeriod)
***REMOVED***

func TestThresholdRun(t *testing.T) ***REMOVED***
	t.Run("true", func(t *testing.T) ***REMOVED***
		th, err := newThreshold(`1+1==2`, goja.New(), false, types.NullDuration***REMOVED******REMOVED***)
		assert.NoError(t, err)

		t.Run("no taint", func(t *testing.T) ***REMOVED***
			b, err := th.runNoTaint()
			assert.NoError(t, err)
			assert.True(t, b)
			assert.False(t, th.LastFailed)
		***REMOVED***)

		t.Run("taint", func(t *testing.T) ***REMOVED***
			b, err := th.run()
			assert.NoError(t, err)
			assert.True(t, b)
			assert.False(t, th.LastFailed)
		***REMOVED***)
	***REMOVED***)

	t.Run("false", func(t *testing.T) ***REMOVED***
		th, err := newThreshold(`1+1==4`, goja.New(), false, types.NullDuration***REMOVED******REMOVED***)
		assert.NoError(t, err)

		t.Run("no taint", func(t *testing.T) ***REMOVED***
			b, err := th.runNoTaint()
			assert.NoError(t, err)
			assert.False(t, b)
			assert.False(t, th.LastFailed)
		***REMOVED***)

		t.Run("taint", func(t *testing.T) ***REMOVED***
			b, err := th.run()
			assert.NoError(t, err)
			assert.False(t, b)
			assert.True(t, th.LastFailed)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestNewThresholds(t *testing.T) ***REMOVED***
	t.Run("empty", func(t *testing.T) ***REMOVED***
		ts, err := NewThresholds([]string***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Len(t, ts.Thresholds, 0)
	***REMOVED***)
	t.Run("two", func(t *testing.T) ***REMOVED***
		sources := []string***REMOVED***`1+1==2`, `1+1==4`***REMOVED***
		ts, err := NewThresholds(sources)
		assert.NoError(t, err)
		assert.Len(t, ts.Thresholds, 2)
		for i, th := range ts.Thresholds ***REMOVED***
			assert.Equal(t, sources[i], th.Source)
			assert.False(t, th.LastFailed)
			assert.False(t, th.AbortOnFail)
			assert.NotNil(t, th.pgm)
			assert.Equal(t, ts.Runtime, th.rt)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestNewThresholdsWithConfig(t *testing.T) ***REMOVED***
	t.Run("empty", func(t *testing.T) ***REMOVED***
		ts, err := NewThresholds([]string***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Len(t, ts.Thresholds, 0)
	***REMOVED***)
	t.Run("two", func(t *testing.T) ***REMOVED***
		configs := []thresholdConfig***REMOVED***
			***REMOVED***`1+1==2`, false, types.NullDuration***REMOVED******REMOVED******REMOVED***,
			***REMOVED***`1+1==4`, true, types.NullDuration***REMOVED******REMOVED******REMOVED***,
		***REMOVED***
		ts, err := newThresholdsWithConfig(configs)
		assert.NoError(t, err)
		assert.Len(t, ts.Thresholds, 2)
		for i, th := range ts.Thresholds ***REMOVED***
			assert.Equal(t, configs[i].Threshold, th.Source)
			assert.False(t, th.LastFailed)
			assert.Equal(t, configs[i].AbortOnFail, th.AbortOnFail)
			assert.NotNil(t, th.pgm)
			assert.Equal(t, ts.Runtime, th.rt)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestThresholdsUpdateVM(t *testing.T) ***REMOVED***
	ts, err := NewThresholds(nil)
	assert.NoError(t, err)
	assert.NoError(t, ts.updateVM(DummySink***REMOVED***"a": 1234.5***REMOVED***, 0))
	assert.Equal(t, 1234.5, ts.Runtime.Get("a").ToFloat())
***REMOVED***

func TestThresholdsRunAll(t *testing.T) ***REMOVED***
	zero := types.NullDuration***REMOVED******REMOVED***
	oneSec := types.NullDurationFrom(time.Second)
	twoSec := types.NullDurationFrom(2 * time.Second)
	testdata := map[string]struct ***REMOVED***
		succ  bool
		err   bool
		abort bool
		grace types.NullDuration
		srcs  []string
	***REMOVED******REMOVED***
		"one passing":                ***REMOVED***true, false, false, zero, []string***REMOVED***`1+1==2`***REMOVED******REMOVED***,
		"one failing":                ***REMOVED***false, false, false, zero, []string***REMOVED***`1+1==4`***REMOVED******REMOVED***,
		"two passing":                ***REMOVED***true, false, false, zero, []string***REMOVED***`1+1==2`, `2+2==4`***REMOVED******REMOVED***,
		"two failing":                ***REMOVED***false, false, false, zero, []string***REMOVED***`1+1==4`, `2+2==2`***REMOVED******REMOVED***,
		"two mixed":                  ***REMOVED***false, false, false, zero, []string***REMOVED***`1+1==2`, `1+1==4`***REMOVED******REMOVED***,
		"one erroring":               ***REMOVED***false, true, false, zero, []string***REMOVED***`throw new Error('?!');`***REMOVED******REMOVED***,
		"one aborting":               ***REMOVED***false, false, true, zero, []string***REMOVED***`1+1==4`***REMOVED******REMOVED***,
		"abort with grace period":    ***REMOVED***false, false, true, oneSec, []string***REMOVED***`1+1==4`***REMOVED******REMOVED***,
		"no abort with grace period": ***REMOVED***false, false, true, twoSec, []string***REMOVED***`1+1==4`***REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			ts, err := NewThresholds(data.srcs)
			assert.Nil(t, err)
			ts.Thresholds[0].AbortOnFail = data.abort
			ts.Thresholds[0].AbortGracePeriod = data.grace

			runDuration := 1500 * time.Millisecond

			assert.NoError(t, err)

			b, err := ts.runAll(runDuration)

			if data.err ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
			***REMOVED***

			if data.succ ***REMOVED***
				assert.True(t, b)
			***REMOVED*** else ***REMOVED***
				assert.False(t, b)
			***REMOVED***

			if data.abort && data.grace.Duration < types.Duration(runDuration) ***REMOVED***
				assert.True(t, ts.Abort)
			***REMOVED*** else ***REMOVED***
				assert.False(t, ts.Abort)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestThresholdsRun(t *testing.T) ***REMOVED***
	ts, err := NewThresholds([]string***REMOVED***"a>0"***REMOVED***)
	assert.NoError(t, err)

	t.Run("error", func(t *testing.T) ***REMOVED***
		b, err := ts.Run(DummySink***REMOVED******REMOVED***, 0)
		assert.Error(t, err)
		assert.False(t, b)
	***REMOVED***)

	t.Run("pass", func(t *testing.T) ***REMOVED***
		b, err := ts.Run(DummySink***REMOVED***"a": 1234.5***REMOVED***, 0)
		assert.NoError(t, err)
		assert.True(t, b)
	***REMOVED***)

	t.Run("fail", func(t *testing.T) ***REMOVED***
		b, err := ts.Run(DummySink***REMOVED***"a": 0***REMOVED***, 0)
		assert.NoError(t, err)
		assert.False(t, b)
	***REMOVED***)
***REMOVED***

func TestThresholdsJSON(t *testing.T) ***REMOVED***
	var testdata = []struct ***REMOVED***
		JSON        string
		srcs        []string
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
			`["1+1==2"]`,
			[]string***REMOVED***"1+1==2"***REMOVED***,
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
			`["1+1==2","1+1==3"]`,
			[]string***REMOVED***"1+1==2", "1+1==3"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"1+1==2"***REMOVED***]`,
			[]string***REMOVED***"1+1==2"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["1+1==2"]`,
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"1+1==2","abortOnFail":true,"delayAbortEval":null***REMOVED***]`,
			[]string***REMOVED***"1+1==2"***REMOVED***,
			true,
			types.NullDuration***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"1+1==2","abortOnFail":true,"delayAbortEval":"2s"***REMOVED***]`,
			[]string***REMOVED***"1+1==2"***REMOVED***,
			true,
			types.NullDurationFrom(2 * time.Second),
			"",
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"1+1==2","abortOnFail":false***REMOVED***]`,
			[]string***REMOVED***"1+1==2"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["1+1==2"]`,
		***REMOVED***,
		***REMOVED***
			`[***REMOVED***"threshold":"1+1==2"***REMOVED***, "1+1==3"]`,
			[]string***REMOVED***"1+1==2", "1+1==3"***REMOVED***,
			false,
			types.NullDuration***REMOVED******REMOVED***,
			`["1+1==2","1+1==3"]`,
		***REMOVED***,
	***REMOVED***

	for _, data := range testdata ***REMOVED***
		t.Run(data.JSON, func(t *testing.T) ***REMOVED***
			var ts Thresholds
			assert.NoError(t, json.Unmarshal([]byte(data.JSON), &ts))
			assert.Equal(t, len(data.srcs), len(ts.Thresholds))
			for i, src := range data.srcs ***REMOVED***
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
		var ts Thresholds
		assert.Error(t, json.Unmarshal([]byte("42"), &ts))
		assert.Nil(t, ts.Thresholds)
		assert.Nil(t, ts.Runtime)
		assert.False(t, ts.Abort)
	***REMOVED***)

	t.Run("bad source", func(t *testing.T) ***REMOVED***
		var ts Thresholds
		assert.Error(t, json.Unmarshal([]byte(`["="]`), &ts))
		assert.Nil(t, ts.Thresholds)
		assert.Nil(t, ts.Runtime)
		assert.False(t, ts.Abort)
	***REMOVED***)
***REMOVED***
