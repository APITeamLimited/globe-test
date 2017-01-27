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
	log "github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"
	"strconv"
	"sync/atomic"
	"time"
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

func (a JSAPI) Log(level int, msg string, args []otto.Value) ***REMOVED***
	fields := make(log.Fields, len(args))
	for i, arg := range args ***REMOVED***
		if arg.IsObject() ***REMOVED***
			obj := arg.Object()
			for _, key := range obj.Keys() ***REMOVED***
				v, err := obj.Get(key)
				if err != nil ***REMOVED***
					throw(a.vu.vm, err)
				***REMOVED***
				fields[key] = v.String()
			***REMOVED***
			continue
		***REMOVED***
		fields["arg"+strconv.Itoa(i)] = arg.String()
	***REMOVED***

	entry := log.WithFields(fields)
	switch level ***REMOVED***
	case 0:
		entry.Debug(msg)
	case 1:
		entry.Info(msg)
	case 2:
		entry.Warn(msg)
	case 3:
		entry.Error(msg)
	***REMOVED***
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

func (a JSAPI) DoCheck(call otto.FunctionCall) otto.Value ***REMOVED***
	if len(call.ArgumentList) < 2 ***REMOVED***
		return otto.UndefinedValue()
	***REMOVED***

	success := true
	arg0 := call.Argument(0)
	for _, v := range call.ArgumentList[1:] ***REMOVED***
		obj := v.Object()
		if obj == nil ***REMOVED***
			panic(call.Otto.MakeTypeError("checks must be objects"))
		***REMOVED***
		for _, name := range obj.Keys() ***REMOVED***
			val, err := obj.Get(name)
			if err != nil ***REMOVED***
				throw(call.Otto, err)
			***REMOVED***

			result, err := Check(val, arg0)
			if err != nil ***REMOVED***
				throw(call.Otto, err)
			***REMOVED***

			check, err := a.vu.group.Check(name)
			if err != nil ***REMOVED***
				throw(call.Otto, err)
			***REMOVED***
			if result ***REMOVED***
				atomic.AddInt64(&(check.Passes), 1)
			***REMOVED*** else ***REMOVED***
				atomic.AddInt64(&(check.Fails), 1)
				success = false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !success ***REMOVED***
		a.vu.Taint = true
		return otto.FalseValue()
	***REMOVED***
	return otto.TrueValue()
***REMOVED***

func (a JSAPI) Taint() ***REMOVED***
	a.vu.Taint = true
***REMOVED***

func (a JSAPI) ElapsedMs() float64 ***REMOVED***
	return float64(time.Since(a.vu.started)) / float64(time.Millisecond)
***REMOVED***
