/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

package data

import (
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
)

// TODO fix it not working really well with setupData or just make it more broken
// TODO fix it working with console.log
type sharedArray struct ***REMOVED***
	arr []string
***REMOVED***

type wrappedSharedArray struct ***REMOVED***
	sharedArray

	rt       *goja.Runtime
	freeze   goja.Callable
	isFrozen goja.Callable
	parse    goja.Callable
***REMOVED***

func (s sharedArray) wrap(rt *goja.Runtime) goja.Value ***REMOVED***
	freeze, _ := goja.AssertFunction(rt.GlobalObject().Get("Object").ToObject(rt).Get("freeze"))
	isFrozen, _ := goja.AssertFunction(rt.GlobalObject().Get("Object").ToObject(rt).Get("isFrozen"))
	parse, _ := goja.AssertFunction(rt.GlobalObject().Get("JSON").ToObject(rt).Get("parse"))
	return rt.NewDynamicArray(wrappedSharedArray***REMOVED***
		sharedArray: s,
		rt:          rt,
		freeze:      freeze,
		isFrozen:    isFrozen,
		parse:       parse,
	***REMOVED***)
***REMOVED***

func (s wrappedSharedArray) Set(index int, val goja.Value) bool ***REMOVED***
	panic(s.rt.NewTypeError("SharedArray is immutable")) // this is specifically a type error
***REMOVED***

func (s wrappedSharedArray) SetLen(len int) bool ***REMOVED***
	panic(s.rt.NewTypeError("SharedArray is immutable")) // this is specifically a type error
***REMOVED***

func (s wrappedSharedArray) Get(index int) goja.Value ***REMOVED***
	if index < 0 || index >= len(s.arr) ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	val, err := s.parse(goja.Undefined(), s.rt.ToValue(s.arr[index]))
	if err != nil ***REMOVED***
		common.Throw(s.rt, err)
	***REMOVED***
	err = s.deepFreeze(s.rt, val)
	if err != nil ***REMOVED***
		common.Throw(s.rt, err)
	***REMOVED***

	return val
***REMOVED***

func (s wrappedSharedArray) Len() int ***REMOVED***
	return len(s.arr)
***REMOVED***

func (s wrappedSharedArray) deepFreeze(rt *goja.Runtime, val goja.Value) error ***REMOVED***
	if val != nil && goja.IsNull(val) ***REMOVED***
		return nil
	***REMOVED***

	_, err := s.freeze(goja.Undefined(), val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	o := val.ToObject(rt)
	if o == nil ***REMOVED***
		return nil
	***REMOVED***
	for _, key := range o.Keys() ***REMOVED***
		prop := o.Get(key)
		if prop != nil ***REMOVED***
			// isFrozen returns true for all non objects so it we don't need to check that
			frozen, err := s.isFrozen(goja.Undefined(), prop)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if !frozen.ToBoolean() ***REMOVED*** // prevent cycles
				if err = s.deepFreeze(rt, prop); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
