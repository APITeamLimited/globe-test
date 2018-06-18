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
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
)

type K6 struct***REMOVED******REMOVED***

func New() *K6 ***REMOVED***
	return &K6***REMOVED******REMOVED***
***REMOVED***

func (*K6) Fail(msg string) (goja.Value, error) ***REMOVED***
	return goja.Undefined(), errors.New(msg)
***REMOVED***

func (*K6) Sleep(ctx context.Context, secs float64) ***REMOVED***
	timer := time.NewTimer(time.Duration(secs * float64(time.Second)))
	select ***REMOVED***
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
	***REMOVED***
***REMOVED***

func (*K6) Group(ctx context.Context, name string, fn goja.Callable) (goja.Value, error) ***REMOVED***
	state := common.GetState(ctx)

	g, err := state.Group.Group(name)
	if err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***

	old := state.Group
	state.Group = g
	defer func() ***REMOVED*** state.Group = old ***REMOVED***()

	startTime := time.Now()
	ret, err := fn(goja.Undefined())
	t := time.Now()

	tags := state.Options.RunTags.CloneTags()
	if state.Options.SystemTags["group"] ***REMOVED***
		tags["group"] = g.Path
	***REMOVED***
	if state.Options.SystemTags["vu"] ***REMOVED***
		tags["vu"] = strconv.FormatInt(state.Vu, 10)
	***REMOVED***
	if state.Options.SystemTags["iter"] ***REMOVED***
		tags["iter"] = strconv.FormatInt(state.Iteration, 10)
	***REMOVED***

	state.Samples <- stats.Sample***REMOVED***
		Time:   t,
		Metric: metrics.GroupDuration,
		Tags:   stats.IntoSampleTags(&tags),
		Value:  stats.D(t.Sub(startTime)),
	***REMOVED***
	return ret, err
***REMOVED***

func (*K6) Check(ctx context.Context, arg0, checks goja.Value, extras ...goja.Value) (bool, error) ***REMOVED***
	state := common.GetState(ctx)
	rt := common.GetRuntime(ctx)
	t := time.Now()

	// Prepare tags, make sure the `group` tag can't be overwritten.
	commonTags := state.Options.RunTags.CloneTags()
	if state.Options.SystemTags["group"] ***REMOVED***
		commonTags["group"] = state.Group.Path
	***REMOVED***
	if len(extras) > 0 ***REMOVED***
		obj := extras[0].ToObject(rt)
		for _, k := range obj.Keys() ***REMOVED***
			commonTags[k] = obj.Get(k).String()
		***REMOVED***
	***REMOVED***
	if state.Options.SystemTags["vu"] ***REMOVED***
		commonTags["vu"] = strconv.FormatInt(state.Vu, 10)
	***REMOVED***
	if state.Options.SystemTags["iter"] ***REMOVED***
		commonTags["iter"] = strconv.FormatInt(state.Iteration, 10)
	***REMOVED***

	succ := true
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
		if state.Options.SystemTags["check"] ***REMOVED***
			tags["check"] = check.Name
		***REMOVED***

		// Resolve callables into values.
		fn, ok := goja.AssertFunction(val)
		if ok ***REMOVED***
			tmpVal, err := fn(goja.Undefined(), arg0)
			if err != nil ***REMOVED***
				return false, err
			***REMOVED***
			val = tmpVal
		***REMOVED***

		sampleTags := stats.IntoSampleTags(&tags)

		// Emit! (But only if we have a valid context.)
		select ***REMOVED***
		case <-ctx.Done():
		default:
			if val.ToBoolean() ***REMOVED***
				atomic.AddInt64(&check.Passes, 1)
				state.Samples <- stats.Sample***REMOVED***Time: t, Metric: metrics.Checks, Tags: sampleTags, Value: 1***REMOVED***
			***REMOVED*** else ***REMOVED***
				atomic.AddInt64(&check.Fails, 1)
				state.Samples <- stats.Sample***REMOVED***Time: t, Metric: metrics.Checks, Tags: sampleTags, Value: 0***REMOVED***
				// A single failure makes the return value false.
				succ = false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return succ, nil
***REMOVED***
