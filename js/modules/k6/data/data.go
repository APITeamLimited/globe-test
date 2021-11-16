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
	"errors"
	"strconv"
	"sync"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct ***REMOVED***
		shared sharedArrays
	***REMOVED***

	// Data represents an instance of the data module.
	Data struct ***REMOVED***
		vu     modules.VU
		shared *sharedArrays
	***REMOVED***

	sharedArrays struct ***REMOVED***
		data map[string]sharedArray
		mu   sync.RWMutex
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &Data***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED***
		shared: sharedArrays***REMOVED***
			data: make(map[string]sharedArray),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &Data***REMOVED***
		vu:     vu,
		shared: &rm.shared,
	***REMOVED***
***REMOVED***

// Exports returns the exports of the data module.
func (d *Data) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"SharedArray": d.sharedArray,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// sharedArray is a constructor returning a shareable read-only array
// indentified by the name and having their contents be whatever the call returns
func (d *Data) sharedArray(call goja.ConstructorCall) *goja.Object ***REMOVED***
	rt := d.vu.Runtime()

	if d.vu.State() != nil ***REMOVED***
		common.Throw(rt, errors.New("new SharedArray must be called in the init context"))
	***REMOVED***

	name := call.Argument(0).String()
	if name == "" ***REMOVED***
		common.Throw(rt, errors.New("empty name provided to SharedArray's constructor"))
	***REMOVED***

	fn, ok := goja.AssertFunction(call.Argument(1))
	if !ok ***REMOVED***
		common.Throw(rt, errors.New("a function is expected as the second argument of SharedArray's constructor"))
	***REMOVED***

	array := d.shared.get(rt, name, fn)
	return array.wrap(rt).ToObject(rt)
***REMOVED***

func (s *sharedArrays) get(rt *goja.Runtime, name string, call goja.Callable) sharedArray ***REMOVED***
	s.mu.RLock()
	array, ok := s.data[name]
	s.mu.RUnlock()
	if !ok ***REMOVED***
		s.mu.Lock()
		defer s.mu.Unlock()
		array, ok = s.data[name]
		if !ok ***REMOVED***
			array = getShareArrayFromCall(rt, call)
			s.data[name] = array
		***REMOVED***
	***REMOVED***

	return array
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
