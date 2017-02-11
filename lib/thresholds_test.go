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

package lib

import (
	"encoding/json"
	"github.com/loadimpact/k6/stats"
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewThreshold(t *testing.T) ***REMOVED***
	src := `1+1==2`
	vm := otto.New()
	th, err := NewThreshold(src, vm)
	assert.NoError(t, err)

	assert.Equal(t, src, th.Source)
	assert.False(t, th.Failed)
	assert.NotNil(t, th.script)
	assert.Equal(t, vm, th.vm)
***REMOVED***

func TestThresholdRun(t *testing.T) ***REMOVED***
	t.Run("true", func(t *testing.T) ***REMOVED***
		th, err := NewThreshold(`1+1==2`, otto.New())
		assert.NoError(t, err)

		t.Run("no taint", func(t *testing.T) ***REMOVED***
			b, err := th.RunNoTaint()
			assert.NoError(t, err)
			assert.True(t, b)
			assert.False(t, th.Failed)
		***REMOVED***)

		t.Run("taint", func(t *testing.T) ***REMOVED***
			b, err := th.Run()
			assert.NoError(t, err)
			assert.True(t, b)
			assert.False(t, th.Failed)
		***REMOVED***)
	***REMOVED***)

	t.Run("false", func(t *testing.T) ***REMOVED***
		th, err := NewThreshold(`1+1==4`, otto.New())
		assert.NoError(t, err)

		t.Run("no taint", func(t *testing.T) ***REMOVED***
			b, err := th.RunNoTaint()
			assert.NoError(t, err)
			assert.False(t, b)
			assert.False(t, th.Failed)
		***REMOVED***)

		t.Run("taint", func(t *testing.T) ***REMOVED***
			b, err := th.Run()
			assert.NoError(t, err)
			assert.False(t, b)
			assert.True(t, th.Failed)
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
			assert.False(t, th.Failed)
			assert.NotNil(t, th.script)
			assert.Equal(t, ts.VM, th.vm)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestThresholdsUpdateVM(t *testing.T) ***REMOVED***
	ts, err := NewThresholds(nil)
	assert.NoError(t, err)
	assert.NoError(t, ts.UpdateVM(stats.DummySink***REMOVED***"a": 1234.5***REMOVED***))

	v, err := ts.VM.Get("a")
	assert.NoError(t, err)
	f, err := v.ToFloat()
	assert.NoError(t, err)
	assert.Equal(t, 1234.5, f)
***REMOVED***

func TestThresholdsRunAll(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		succ bool
		err  bool
		srcs []string
	***REMOVED******REMOVED***
		"one passing":  ***REMOVED***true, false, []string***REMOVED***`1+1==2`***REMOVED******REMOVED***,
		"one failing":  ***REMOVED***false, false, []string***REMOVED***`1+1==4`***REMOVED******REMOVED***,
		"two passing":  ***REMOVED***true, false, []string***REMOVED***`1+1==2`, `2+2==4`***REMOVED******REMOVED***,
		"two failing":  ***REMOVED***false, false, []string***REMOVED***`1+1==4`, `2+2==2`***REMOVED******REMOVED***,
		"two mixed":    ***REMOVED***false, false, []string***REMOVED***`1+1==2`, `1+1==4`***REMOVED******REMOVED***,
		"one erroring": ***REMOVED***false, true, []string***REMOVED***`throw new Error('?!');`***REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			ts, err := NewThresholds(data.srcs)
			assert.NoError(t, err)
			b, err := ts.RunAll()

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
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestThresholdsRun(t *testing.T) ***REMOVED***
	ts, err := NewThresholds([]string***REMOVED***"a>0"***REMOVED***)
	assert.NoError(t, err)

	t.Run("error", func(t *testing.T) ***REMOVED***
		b, err := ts.Run(stats.DummySink***REMOVED******REMOVED***)
		assert.Error(t, err)
		assert.False(t, b)
	***REMOVED***)

	t.Run("pass", func(t *testing.T) ***REMOVED***
		b, err := ts.Run(stats.DummySink***REMOVED***"a": 1234.5***REMOVED***)
		assert.NoError(t, err)
		assert.True(t, b)
	***REMOVED***)

	t.Run("fail", func(t *testing.T) ***REMOVED***
		b, err := ts.Run(stats.DummySink***REMOVED***"a": 0***REMOVED***)
		assert.NoError(t, err)
		assert.False(t, b)
	***REMOVED***)
***REMOVED***

func TestThresholdsJSON(t *testing.T) ***REMOVED***
	testdata := map[string][]string***REMOVED***
		`[]`:                  ***REMOVED******REMOVED***,
		`["1+1==2"]`:          ***REMOVED***"1+1==2"***REMOVED***,
		`["1+1==2","1+1==3"]`: ***REMOVED***"1+1==2", "1+1==3"***REMOVED***,
	***REMOVED***

	for data, srcs := range testdata ***REMOVED***
		t.Run(data, func(t *testing.T) ***REMOVED***
			var ts Thresholds
			assert.NoError(t, json.Unmarshal([]byte(data), &ts))
			assert.Equal(t, len(srcs), len(ts.Thresholds))
			for i, src := range srcs ***REMOVED***
				assert.Equal(t, src, ts.Thresholds[i].Source)
			***REMOVED***

			t.Run("marshal", func(t *testing.T) ***REMOVED***
				data2, err := json.Marshal(ts)
				assert.NoError(t, err)
				assert.Equal(t, data, string(data2))
			***REMOVED***)
		***REMOVED***)
	***REMOVED***
***REMOVED***
