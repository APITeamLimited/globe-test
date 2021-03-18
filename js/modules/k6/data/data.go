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
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
)

type data struct ***REMOVED***
	shared sharedArrays
***REMOVED***

type sharedArrays struct ***REMOVED***
	data map[string]sharedArray
	mu   sync.RWMutex
***REMOVED***

func (s *sharedArrays) get(rt *goja.Runtime, name string, call goja.Callable) sharedArray ***REMOVED***
	s.mu.RLock()
	array, ok := s.data[name]
	s.mu.RUnlock()
	if !ok ***REMOVED***
		s.mu.Lock()
		array, ok = s.data[name]
		if !ok ***REMOVED***
			func() ***REMOVED*** // this is done for the defer below
				defer s.mu.Unlock()
				array = getShareArrayFromCall(rt, call)
				s.data[name] = array
			***REMOVED***()
		***REMOVED***
	***REMOVED***

	return array
***REMOVED***

// New return a new Module instance
func New() interface***REMOVED******REMOVED*** ***REMOVED***
	return &data***REMOVED***
		shared: sharedArrays***REMOVED***
			data: make(map[string]sharedArray),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// XSharedArray is a constructor returning a shareable read-only array
// indentified by the name and having their contents be whatever the call returns
func (d *data) XSharedArray(ctx context.Context, name string, call goja.Callable) (goja.Value, error) ***REMOVED***
	if lib.GetState(ctx) != nil ***REMOVED***
		return nil, errors.New("new SharedArray must be called in the init context")
	***REMOVED***

	if len(name) == 0 ***REMOVED***
		return nil, errors.New("empty name provided to SharedArray's constructor")
	***REMOVED***

	rt := common.GetRuntime(ctx)
	array := d.shared.get(rt, name, call)

	return array.wrap(rt), nil
***REMOVED***

func getShareArrayFromCall(rt *goja.Runtime, call goja.Callable) sharedArray ***REMOVED***
	gojaValue, err := call(goja.Undefined())
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	obj := gojaValue.ToObject(rt)
	if obj.ClassName() != "Array" ***REMOVED***
		common.Throw(rt, errors.New("only arrays can be made into SharedArray")) // TODO better error
	***REMOVED***
	arr := make([]string, obj.Get("length").ToInteger())

	stringify, _ := goja.AssertFunction(rt.GlobalObject().Get("JSON").ToObject(rt).Get("stringify"))
	var val goja.Value
	for i := range arr ***REMOVED***
		val, err = stringify(goja.Undefined(), obj.Get(strconv.Itoa(i)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		arr[i] = val.String()
	***REMOVED***

	return sharedArray***REMOVED***arr: arr***REMOVED***
***REMOVED***
