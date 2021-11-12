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

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/stats"
)

func TestFail(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	require.NoError(t, rt.Set("k6", common.Bind(rt, New(), nil)))
	_, err := rt.RunString(`k6.fail("blah")`)
	assert.Contains(t, err.Error(), "blah")
***REMOVED***

func TestSleep(t *testing.T) ***REMOVED***
	t.Parallel()

	testdata := map[string]time.Duration***REMOVED***
		"1":   1 * time.Second,
		"1.0": 1 * time.Second,
		"0.5": 500 * time.Millisecond,
	***REMOVED***
	for name, d := range testdata ***REMOVED***
		d := d
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			rt := goja.New()
			ctx := context.Background()
			require.NoError(t, rt.Set("k6", common.Bind(rt, New(), &ctx)))
			startTime := time.Now()
			_, err := rt.RunString(`k6.sleep(1)`)
			endTime := time.Now()
			assert.NoError(t, err)
			assert.True(t, endTime.Sub(startTime) > d, "did not sleep long enough")
		***REMOVED***)
	***REMOVED***

	t.Run("Cancel", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := goja.New()
		ctx, cancel := context.WithCancel(context.Background())
		require.NoError(t, rt.Set("k6", common.Bind(rt, New(), &ctx)))
		dch := make(chan time.Duration)
		go func() ***REMOVED***
			startTime := time.Now()
			_, err := rt.RunString(`k6.sleep(10)`)
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

func TestRandSeed(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()

	ctx := context.Background()
	ctx = common.WithRuntime(ctx, rt)

	require.NoError(t, rt.Set("k6", common.Bind(rt, New(), &ctx)))

	rand := 0.8487305991992138
	_, err := rt.RunString(fmt.Sprintf(`
		var rnd = Math.random();
		if (rnd == %.16f) ***REMOVED*** throw new Error("wrong random: " + rnd); ***REMOVED***
	`, rand))
	assert.NoError(t, err)

	_, err = rt.RunString(fmt.Sprintf(`
		k6.randomSeed(12345)
		var rnd = Math.random();
		if (rnd != %.16f) ***REMOVED*** throw new Error("wrong random: " + rnd); ***REMOVED***
	`, rand))
	assert.NoError(t, err)
***REMOVED***

func TestGroup(t *testing.T) ***REMOVED***
	t.Parallel()
	setupGroupTest := func() (*goja.Runtime, *lib.State, *lib.Group) ***REMOVED***
		root, err := lib.NewGroup("", nil)
		assert.NoError(t, err)

		rt := goja.New()
		state := &lib.State***REMOVED***
			Group:   root,
			Samples: make(chan stats.SampleContainer, 1000),
			Tags:    lib.NewTagMap(nil),
			Options: lib.Options***REMOVED***
				SystemTags: stats.NewSystemTagSet(stats.TagGroup),
			***REMOVED***,
		***REMOVED***
		ctx := context.Background()
		ctx = lib.WithState(ctx, state)
		ctx = common.WithRuntime(ctx, rt)
		state.BuiltinMetrics = metrics.RegisterBuiltinMetrics(metrics.NewRegistry())
		require.NoError(t, rt.Set("k6", common.Bind(rt, New(), &ctx)))
		return rt, state, root
	***REMOVED***

	t.Run("Valid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt, state, root := setupGroupTest()
		assert.Equal(t, state.Group, root)
		require.NoError(t, rt.Set("fn", func() ***REMOVED***
			groupTag, ok := state.Tags.Get("group")
			require.True(t, ok)
			assert.Equal(t, groupTag, "::my group")
			assert.Equal(t, state.Group.Name, "my group")
			assert.Equal(t, state.Group.Parent, root)
		***REMOVED***))
		_, err := rt.RunString(`k6.group("my group", fn)`)
		assert.NoError(t, err)
		assert.Equal(t, state.Group, root)
		groupTag, ok := state.Tags.Get("group")
		require.True(t, ok)
		assert.Equal(t, groupTag, root.Name)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt, _, _ := setupGroupTest()
		_, err := rt.RunString(`k6.group("::", function() ***REMOVED*** throw new Error("nooo") ***REMOVED***)`)
		assert.Contains(t, err.Error(), "group and check names may not contain '::'")
	***REMOVED***)
***REMOVED***

func checkTestRuntime(t testing.TB, ctxs ...*context.Context) (
	*goja.Runtime, chan stats.SampleContainer, *metrics.BuiltinMetrics,
) ***REMOVED***
	rt := goja.New()

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)
	samples := make(chan stats.SampleContainer, 1000)
	state := &lib.State***REMOVED***
		Group: root,
		Options: lib.Options***REMOVED***
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		Samples: samples,
		Tags: lib.NewTagMap(map[string]string***REMOVED***
			"group": root.Path,
		***REMOVED***),
	***REMOVED***
	ctx := context.Background()
	if len(ctxs) == 1 ***REMOVED*** // hacks
		ctx = *ctxs[0]
	***REMOVED***
	ctx = common.WithRuntime(ctx, rt)
	ctx = lib.WithState(ctx, state)
	state.BuiltinMetrics = metrics.RegisterBuiltinMetrics(metrics.NewRegistry())
	require.NoError(t, rt.Set("k6", common.Bind(rt, New(), &ctx)))
	if len(ctxs) == 1 ***REMOVED*** // hacks
		*ctxs[0] = ctx
	***REMOVED***
	return rt, samples, state.BuiltinMetrics
***REMOVED***

func TestCheckObject(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, samples, builtinMetrics := checkTestRuntime(t)

	_, err := rt.RunString(`k6.check(null, ***REMOVED*** "check": true ***REMOVED***)`)
	assert.NoError(t, err)

	bufSamples := stats.GetBufferedSamples(samples)
	if assert.Len(t, bufSamples, 1) ***REMOVED***
		sample, ok := bufSamples[0].(stats.Sample)
		require.True(t, ok)

		assert.NotZero(t, sample.Time)
		assert.Equal(t, builtinMetrics.Checks, sample.Metric)
		assert.Equal(t, float64(1), sample.Value)
		assert.Equal(t, map[string]string***REMOVED***
			"group": "",
			"check": "check",
		***REMOVED***, sample.Tags.CloneTags())
	***REMOVED***

	t.Run("Multiple", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt, samples, _ := checkTestRuntime(t)

		_, err := rt.RunString(`k6.check(null, ***REMOVED*** "a": true, "b": false ***REMOVED***)`)
		assert.NoError(t, err)

		bufSamples := stats.GetBufferedSamples(samples)
		assert.Len(t, bufSamples, 2)
		var foundA, foundB bool
		for _, sampleC := range bufSamples ***REMOVED***
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
		t.Parallel()
		rt, _, _ := checkTestRuntime(t)
		_, err := rt.RunString(`k6.check(null, ***REMOVED*** "::": true ***REMOVED***)`)
		assert.Contains(t, err.Error(), "group and check names may not contain '::'")
	***REMOVED***)
***REMOVED***

func TestCheckArray(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, samples, builtinMetrics := checkTestRuntime(t)

	_, err := rt.RunString(`k6.check(null, [ true ])`)
	assert.NoError(t, err)

	bufSamples := stats.GetBufferedSamples(samples)
	if assert.Len(t, bufSamples, 1) ***REMOVED***
		sample, ok := bufSamples[0].(stats.Sample)
		require.True(t, ok)

		assert.NotZero(t, sample.Time)
		assert.Equal(t, builtinMetrics.Checks, sample.Metric)
		assert.Equal(t, float64(1), sample.Value)
		assert.Equal(t, map[string]string***REMOVED***
			"group": "",
			"check": "0",
		***REMOVED***, sample.Tags.CloneTags())
	***REMOVED***
***REMOVED***

func TestCheckLiteral(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, samples, _ := checkTestRuntime(t)

	_, err := rt.RunString(`k6.check(null, 12345)`)
	assert.NoError(t, err)
	assert.Len(t, stats.GetBufferedSamples(samples), 0)
***REMOVED***

func TestCheckThrows(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, samples, builtinMetrics := checkTestRuntime(t)
	_, err := rt.RunString(`
		k6.check(null, ***REMOVED***
			"a": function() ***REMOVED*** throw new Error("error A") ***REMOVED***,
			"b": function() ***REMOVED*** throw new Error("error B") ***REMOVED***,
		***REMOVED***)
		`)
	assert.EqualError(t, err, "Error: error A at a (<eval>:3:28(4))")

	bufSamples := stats.GetBufferedSamples(samples)
	if assert.Len(t, bufSamples, 1) ***REMOVED***
		sample, ok := bufSamples[0].(stats.Sample)
		require.True(t, ok)

		assert.NotZero(t, sample.Time)
		assert.Equal(t, builtinMetrics.Checks, sample.Metric)
		assert.Equal(t, float64(0), sample.Value)
		assert.Equal(t, map[string]string***REMOVED***
			"group": "",
			"check": "a",
		***REMOVED***, sample.Tags.CloneTags())
	***REMOVED***
***REMOVED***

func TestCheckTypes(t *testing.T) ***REMOVED***
	t.Parallel()
	templates := map[string]string***REMOVED***
		"Literal":      `k6.check(null,***REMOVED***"check": %s***REMOVED***)`,
		"Callable":     `k6.check(null,***REMOVED***"check": function() ***REMOVED*** return %s; ***REMOVED******REMOVED***)`,
		"Callable/Arg": `k6.check(%s,***REMOVED***"check": function(v) ***REMOVED***return v; ***REMOVED******REMOVED***)`,
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
		name, tpl := name, tpl
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for value, succ := range testdata ***REMOVED***
				value, succ := value, succ
				t.Run(value, func(t *testing.T) ***REMOVED***
					t.Parallel()
					rt, samples, builtinMetrics := checkTestRuntime(t)

					v, err := rt.RunString(fmt.Sprintf(tpl, value))
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, succ, v.Export())
					***REMOVED***

					bufSamples := stats.GetBufferedSamples(samples)
					if assert.Len(t, bufSamples, 1) ***REMOVED***
						sample, ok := bufSamples[0].(stats.Sample)
						require.True(t, ok)

						assert.NotZero(t, sample.Time)
						assert.Equal(t, builtinMetrics.Checks, sample.Metric)
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
***REMOVED***

func TestCheckContextExpiry(t *testing.T) ***REMOVED***
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	rt, _, _ := checkTestRuntime(t, &ctx)
	root := lib.GetState(ctx).Group

	v, err := rt.RunString(`k6.check(null, ***REMOVED*** "check": true ***REMOVED***)`)
	if assert.NoError(t, err) ***REMOVED***
		assert.Equal(t, true, v.Export())
	***REMOVED***

	check, _ := root.Check("check")
	assert.Equal(t, int64(1), check.Passes)
	assert.Equal(t, int64(0), check.Fails)

	cancel()

	v, err = rt.RunString(`k6.check(null, ***REMOVED*** "check": true ***REMOVED***)`)
	if assert.NoError(t, err) ***REMOVED***
		assert.Equal(t, true, v.Export())
	***REMOVED***

	assert.Equal(t, int64(1), check.Passes)
	assert.Equal(t, int64(0), check.Fails)
***REMOVED***

func TestCheckTags(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, samples, builtinMetrics := checkTestRuntime(t)

	v, err := rt.RunString(`k6.check(null, ***REMOVED***"check": true***REMOVED***, ***REMOVED***a: 1, b: "2"***REMOVED***)`)
	if assert.NoError(t, err) ***REMOVED***
		assert.Equal(t, true, v.Export())
	***REMOVED***

	bufSamples := stats.GetBufferedSamples(samples)
	if assert.Len(t, bufSamples, 1) ***REMOVED***
		sample, ok := bufSamples[0].(stats.Sample)
		require.True(t, ok)

		assert.NotZero(t, sample.Time)
		assert.Equal(t, builtinMetrics.Checks, sample.Metric)
		assert.Equal(t, float64(1), sample.Value)
		assert.Equal(t, map[string]string***REMOVED***
			"group": "",
			"check": "check",
			"a":     "1",
			"b":     "2",
		***REMOVED***, sample.Tags.CloneTags())
	***REMOVED***
***REMOVED***
