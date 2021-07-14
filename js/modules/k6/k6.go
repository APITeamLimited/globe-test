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

// Package k6 implements the module imported as 'k6' from inside k6.
package k6

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/stats"
)

// K6 is just the module struct.
type K6 struct***REMOVED******REMOVED***

// ErrGroupInInitContext is returned when group() are using in the init context.
var ErrGroupInInitContext = common.NewInitContextError("Using group() in the init context is not supported")

// ErrCheckInInitContext is returned when check() are using in the init context.
var ErrCheckInInitContext = common.NewInitContextError("Using check() in the init context is not supported")

// New returns a new module Struct.
func New() *K6 ***REMOVED***
	return &K6***REMOVED******REMOVED***
***REMOVED***

// Fail is a fancy way of saying `throw "something"`.
func (*K6) Fail(msg string) (goja.Value, error) ***REMOVED***
	return goja.Undefined(), errors.New(msg)
***REMOVED***

// Sleep waits the provided seconds before continuing the execution.
func (*K6) Sleep(ctx context.Context, secs float64) ***REMOVED***
	timer := time.NewTimer(time.Duration(secs * float64(time.Second)))
	select ***REMOVED***
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
	***REMOVED***
***REMOVED***

// RandomSeed sets the seed to the random generator used for this VU.
func (*K6) RandomSeed(ctx context.Context, seed int64) ***REMOVED***
	randSource := rand.New(rand.NewSource(seed)).Float64 //nolint:gosec

	rt := common.GetRuntime(ctx)
	rt.SetRandSource(randSource)
***REMOVED***

// Group wraps a function call and executes it within the provided group name.
func (*K6) Group(ctx context.Context, name string, fn goja.Callable) (goja.Value, error) ***REMOVED***
	state := lib.GetState(ctx)
	if state == nil ***REMOVED***
		return nil, ErrGroupInInitContext
	***REMOVED***

	if fn == nil ***REMOVED***
		return nil, errors.New("group() requires a callback as a second argument")
	***REMOVED***

	g, err := state.Group.Group(name)
	if err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***

	old := state.Group
	state.Group = g

	shouldUpdateTag := state.Options.SystemTags.Has(stats.TagGroup)
	if shouldUpdateTag ***REMOVED***
		state.Tags["group"] = g.Path
	***REMOVED***
	defer func() ***REMOVED***
		state.Group = old
		if shouldUpdateTag ***REMOVED***
			state.Tags["group"] = old.Path
		***REMOVED***
	***REMOVED***()

	startTime := time.Now()
	ret, err := fn(goja.Undefined())
	t := time.Now()

	tags := state.CloneTags()
	stats.PushIfNotDone(ctx, state.Samples, stats.Sample***REMOVED***
		Time:   t,
		Metric: metrics.GroupDuration,
		Tags:   stats.IntoSampleTags(&tags),
		Value:  stats.D(t.Sub(startTime)),
	***REMOVED***)

	return ret, err
***REMOVED***

// Check will emit check metrics for the provided checks.
//nolint:cyclop
func (*K6) Check(ctx context.Context, arg0, checks goja.Value, extras ...goja.Value) (bool, error) ***REMOVED***
	state := lib.GetState(ctx)
	if state == nil ***REMOVED***
		return false, ErrCheckInInitContext
	***REMOVED***
	rt := common.GetRuntime(ctx)
	t := time.Now()

	// Prepare the metric tags
	commonTags := state.CloneTags()
	if len(extras) > 0 ***REMOVED***
		obj := extras[0].ToObject(rt)
		for _, k := range obj.Keys() ***REMOVED***
			commonTags[k] = obj.Get(k).String()
		***REMOVED***
	***REMOVED***

	succ := true
	var exc error
	obj := checks.ToObject(rt)
	for _, name := range obj.Keys() ***REMOVED***
		val := obj.Get(name)

		tags := make(map[string]string, len(commonTags))
		for k, v := range commonTags ***REMOVED***
			tags[k] = v
		***REMOVED***

		// Resolve the check record.
		check, err := state.Group.Check(name)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		if state.Options.SystemTags.Has(stats.TagCheck) ***REMOVED***
			tags["check"] = check.Name
		***REMOVED***

		// Resolve callables into values.
		fn, ok := goja.AssertFunction(val)
		if ok ***REMOVED***
			tmpVal, err := fn(goja.Undefined(), arg0)
			val = tmpVal
			if err != nil ***REMOVED***
				val = rt.ToValue(false)
				exc = err
			***REMOVED***
		***REMOVED***

		sampleTags := stats.IntoSampleTags(&tags)

		// Emit! (But only if we have a valid context.)
		select ***REMOVED***
		case <-ctx.Done():
		default:
			if val.ToBoolean() ***REMOVED***
				atomic.AddInt64(&check.Passes, 1)
				stats.PushIfNotDone(ctx, state.Samples, stats.Sample***REMOVED***Time: t, Metric: metrics.Checks, Tags: sampleTags, Value: 1***REMOVED***)
			***REMOVED*** else ***REMOVED***
				atomic.AddInt64(&check.Fails, 1)
				stats.PushIfNotDone(ctx, state.Samples, stats.Sample***REMOVED***Time: t, Metric: metrics.Checks, Tags: sampleTags, Value: 0***REMOVED***)
				// A single failure makes the return value false.
				succ = false
			***REMOVED***
		***REMOVED***

		if exc != nil ***REMOVED***
			return succ, exc
		***REMOVED***
	***REMOVED***

	return succ, nil
***REMOVED***
