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

package k6

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/loadimpact/k6/stats"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFail(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.Set("k6", common.Bind(rt, New(), nil))
	_, err := common.RunString(rt, `k6.fail("blah")`)
	assert.EqualError(t, err, "GoError: blah")
***REMOVED***

func TestSleep(t *testing.T) ***REMOVED***
	rt := goja.New()
	ctx, cancel := context.WithCancel(context.Background())
	rt.Set("k6", common.Bind(rt, New(), &ctx))

	testdata := map[string]time.Duration***REMOVED***
		"1":   1 * time.Second,
		"1.0": 1 * time.Second,
		"0.5": 500 * time.Millisecond,
	***REMOVED***
	for name, d := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			startTime := time.Now()
			_, err := common.RunString(rt, `k6.sleep(1)`)
			endTime := time.Now()
			assert.NoError(t, err)
			assert.True(t, endTime.Sub(startTime) > d, "did not sleep long enough")
		***REMOVED***)
	***REMOVED***

	t.Run("Cancel", func(t *testing.T) ***REMOVED***
		dch := make(chan time.Duration)
		go func() ***REMOVED***
			startTime := time.Now()
			_, err := common.RunString(rt, `k6.sleep(10)`)
			endTime := time.Now()
			assert.NoError(t, err)
			dch <- endTime.Sub(startTime)
		***REMOVED***()
		runtime.Gosched()
		time.Sleep(1 * time.Second)
		runtime.Gosched()
		cancel()
		runtime.Gosched()
		d := <-dch
		assert.True(t, d > 500*time.Millisecond, "did not sleep long enough")
		assert.True(t, d < 2*time.Second, "slept for too long!!")
	***REMOVED***)
***REMOVED***

func TestGroup(t *testing.T) ***REMOVED***
	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	rt := goja.New()
	state := &common.State***REMOVED***Group: root***REMOVED***

	ctx := context.Background()
	ctx = common.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, rt)
	rt.Set("k6", common.Bind(rt, New(), &ctx))

	t.Run("Valid", func(t *testing.T) ***REMOVED***
		assert.Equal(t, state.Group, root)
		rt.Set("fn", func() ***REMOVED***
			assert.Equal(t, state.Group.Name, "my group")
			assert.Equal(t, state.Group.Parent, root)
		***REMOVED***)
		_, err = common.RunString(rt, `k6.group("my group", fn)`)
		assert.NoError(t, err)
		assert.Equal(t, state.Group, root)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `k6.group("::", function() ***REMOVED*** throw new Error("nooo") ***REMOVED***)`)
		assert.EqualError(t, err, "GoError: group and check names may not contain '::'")
	***REMOVED***)
***REMOVED***

func TestCheck(t *testing.T) ***REMOVED***
	rt := goja.New()

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	baseCtx := common.WithRuntime(context.Background(), rt)

	ctx := new(context.Context)
	*ctx = baseCtx
	rt.Set("k6", common.Bind(rt, New(), ctx))

	getState := func() *common.State ***REMOVED***
		return &common.State***REMOVED***
			Group: root,
			Options: lib.Options***REMOVED***
				SystemTags: lib.GetTagSet(lib.DefaultSystemTagList...),
			***REMOVED***,
		***REMOVED***
	***REMOVED***
	t.Run("Object", func(t *testing.T) ***REMOVED***
		state := getState()
		*ctx = common.WithState(baseCtx, state)

		_, err := common.RunString(rt, `k6.check(null, ***REMOVED*** "check": true ***REMOVED***)`)
		assert.NoError(t, err)

		if assert.Len(t, state.Samples, 1) ***REMOVED***
			sample, ok := state.Samples[0].(stats.Sample)
			require.True(t, ok)

			assert.NotZero(t, sample.Time)
			assert.Equal(t, metrics.Checks, sample.Metric)
			assert.Equal(t, float64(1), sample.Value)
			assert.Equal(t, map[string]string***REMOVED***
				"group": "",
				"check": "check",
			***REMOVED***, sample.Tags.CloneTags())
		***REMOVED***

		t.Run("Multiple", func(t *testing.T) ***REMOVED***
			state := getState()
			*ctx = common.WithState(baseCtx, state)

			_, err := common.RunString(rt, `k6.check(null, ***REMOVED*** "a": true, "b": false ***REMOVED***)`)
			assert.NoError(t, err)

			assert.Len(t, state.Samples, 2)
			var foundA, foundB bool
			for _, sampleC := range state.Samples ***REMOVED***
				for _, sample := range sampleC.GetSamples() ***REMOVED***
					name, ok := sample.Tags.Get("check")
					assert.True(t, ok)
					switch name ***REMOVED***
					case "a":
						assert.False(t, foundA, "duplicate 'a'")
						foundA = true
					case "b":
						assert.False(t, foundB, "duplicate 'b'")
						foundB = true
					default:
						assert.Fail(t, name)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			assert.True(t, foundA, "missing 'a'")
			assert.True(t, foundB, "missing 'b'")
		***REMOVED***)

		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `k6.check(null, ***REMOVED*** "::": true ***REMOVED***)`)
			assert.EqualError(t, err, "GoError: group and check names may not contain '::'")
		***REMOVED***)
	***REMOVED***)

	t.Run("Array", func(t *testing.T) ***REMOVED***
		state := getState()
		*ctx = common.WithState(baseCtx, state)

		_, err := common.RunString(rt, `k6.check(null, [ true ])`)
		assert.NoError(t, err)

		if assert.Len(t, state.Samples, 1) ***REMOVED***
			sample, ok := state.Samples[0].(stats.Sample)
			require.True(t, ok)

			assert.NotZero(t, sample.Time)
			assert.Equal(t, metrics.Checks, sample.Metric)
			assert.Equal(t, float64(1), sample.Value)
			assert.Equal(t, map[string]string***REMOVED***
				"group": "",
				"check": "0",
			***REMOVED***, sample.Tags.CloneTags())
		***REMOVED***
	***REMOVED***)

	t.Run("Literal", func(t *testing.T) ***REMOVED***
		state := getState()
		*ctx = common.WithState(baseCtx, state)

		_, err := common.RunString(rt, `k6.check(null, 12345)`)
		assert.NoError(t, err)
		assert.Len(t, state.Samples, 0)
	***REMOVED***)

	t.Run("Throws", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, `
		k6.check(null, ***REMOVED***
			"a": function() ***REMOVED*** throw new Error("error A") ***REMOVED***,
			"b": function() ***REMOVED*** throw new Error("error B") ***REMOVED***,
		***REMOVED***)
		`)
		assert.EqualError(t, err, "Error: error A at a (<eval>:3:27(6))")
	***REMOVED***)

	t.Run("Types", func(t *testing.T) ***REMOVED***
		templates := map[string]string***REMOVED***
			"Literal":      `k6.check(null,***REMOVED***"check": %s***REMOVED***)`,
			"Callable":     `k6.check(null,***REMOVED***"check": ()=>%s***REMOVED***)`,
			"Callable/Arg": `k6.check(%s,***REMOVED***"check":(v)=>v***REMOVED***)`,
		***REMOVED***
		testdata := map[string]bool***REMOVED***
			`0`:         false,
			`1`:         true,
			`-1`:        true,
			`""`:        false,
			`"true"`:    true,
			`"false"`:   true,
			`true`:      true,
			`false`:     false,
			`null`:      false,
			`undefined`: false,
		***REMOVED***
		for name, tpl := range templates ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				for value, succ := range testdata ***REMOVED***
					t.Run(value, func(t *testing.T) ***REMOVED***
						state := getState()
						*ctx = common.WithState(baseCtx, state)

						v, err := common.RunString(rt, fmt.Sprintf(tpl, value))
						if assert.NoError(t, err) ***REMOVED***
							assert.Equal(t, succ, v.Export())
						***REMOVED***

						if assert.Len(t, state.Samples, 1) ***REMOVED***
							sample, ok := state.Samples[0].(stats.Sample)
							require.True(t, ok)

							assert.NotZero(t, sample.Time)
							assert.Equal(t, metrics.Checks, sample.Metric)
							if succ ***REMOVED***
								assert.Equal(t, float64(1), sample.Value)
							***REMOVED*** else ***REMOVED***
								assert.Equal(t, float64(0), sample.Value)
							***REMOVED***
							assert.Equal(t, map[string]string***REMOVED***
								"group": "",
								"check": "check",
							***REMOVED***, sample.Tags.CloneTags())
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			***REMOVED***)
		***REMOVED***

		t.Run("ContextExpiry", func(t *testing.T) ***REMOVED***
			root, err := lib.NewGroup("", nil)
			assert.NoError(t, err)

			state := &common.State***REMOVED***Group: root***REMOVED***
			ctx2, cancel := context.WithCancel(common.WithState(baseCtx, state))
			*ctx = ctx2

			v, err := common.RunString(rt, `k6.check(null, ***REMOVED*** "check": true ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***

			check, _ := root.Check("check")
			assert.Equal(t, int64(1), check.Passes)
			assert.Equal(t, int64(0), check.Fails)

			cancel()

			v, err = common.RunString(rt, `k6.check(null, ***REMOVED*** "check": true ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***

			assert.Equal(t, int64(1), check.Passes)
			assert.Equal(t, int64(0), check.Fails)
		***REMOVED***)
	***REMOVED***)

	t.Run("Tags", func(t *testing.T) ***REMOVED***
		state := getState()
		*ctx = common.WithState(baseCtx, state)

		v, err := common.RunString(rt, `k6.check(null, ***REMOVED***"check": true***REMOVED***, ***REMOVED***a: 1, b: "2"***REMOVED***)`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***

		if assert.Len(t, state.Samples, 1) ***REMOVED***
			sample, ok := state.Samples[0].(stats.Sample)
			require.True(t, ok)

			assert.NotZero(t, sample.Time)
			assert.Equal(t, metrics.Checks, sample.Metric)
			assert.Equal(t, float64(1), sample.Value)
			assert.Equal(t, map[string]string***REMOVED***
				"group": "",
				"check": "check",
				"a":     "1",
				"b":     "2",
			***REMOVED***, sample.Tags.CloneTags())
		***REMOVED***
	***REMOVED***)
***REMOVED***
