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

package js

import (
	"sync/atomic"
	"time"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/robertkrimen/otto"
)

type JSAPI struct ***REMOVED***
	vu *VU
***REMOVED***

func (a JSAPI) Sleep(secs float64) ***REMOVED***
	d := time.Duration(secs * float64(time.Second))
	t := time.NewTimer(d)
	select ***REMOVED***
	case <-t.C:
	case <-a.vu.ctx.Done():
	***REMOVED***
	t.Stop()
***REMOVED***

func (a JSAPI) DoGroup(call otto.FunctionCall) otto.Value ***REMOVED***
	name := call.Argument(0).String()
	group, err := a.vu.group.Group(name)
	if err != nil ***REMOVED***
		throw(call.Otto, err)
	***REMOVED***
	a.vu.group = group
	defer func() ***REMOVED*** a.vu.group = group.Parent ***REMOVED***()

	fn := call.Argument(1)
	if !fn.IsFunction() ***REMOVED***
		panic(call.Otto.MakeSyntaxError("fn must be a function"))
	***REMOVED***

	val, err := fn.Call(call.This)
	if err != nil ***REMOVED***
		throw(call.Otto, err)
	***REMOVED***

	if val.IsUndefined() ***REMOVED***
		return otto.TrueValue()
	***REMOVED***
	return val
***REMOVED***

func (a JSAPI) DoCheck(obj otto.Value, conds map[string]otto.Value, extraTags map[string]string) bool ***REMOVED***
	t := time.Now()
	success := true
	for name, cond := range conds ***REMOVED***
		check, err := a.vu.group.Check(name)
		if err != nil ***REMOVED***
			throw(a.vu.vm, err)
		***REMOVED***

		result, err := Check(cond, obj)
		if err != nil ***REMOVED***
			throw(a.vu.vm, err)
		***REMOVED***

		tags := map[string]string***REMOVED***
			"group": check.Group.Path,
			"check": check.Name,
		***REMOVED***
		for k, v := range extraTags ***REMOVED***
			tags[k] = v
		***REMOVED***

		if result ***REMOVED***
			atomic.AddInt64(&check.Passes, 1)
			a.vu.Samples = append(a.vu.Samples,
				stats.Sample***REMOVED***Time: t, Metric: metrics.Checks, Tags: tags, Value: 1***REMOVED***,
			)
		***REMOVED*** else ***REMOVED***
			success = false
			atomic.AddInt64(&check.Fails, 1)
			a.vu.Samples = append(a.vu.Samples,
				stats.Sample***REMOVED***Time: t, Metric: metrics.Checks, Tags: tags, Value: 0***REMOVED***,
			)
		***REMOVED***
	***REMOVED***

	return success
***REMOVED***

func (a JSAPI) ElapsedMs() float64 ***REMOVED***
	return float64(time.Since(a.vu.started)) / float64(time.Millisecond)
***REMOVED***
