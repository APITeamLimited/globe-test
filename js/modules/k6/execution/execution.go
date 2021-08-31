/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package execution

import (
	"errors"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct***REMOVED******REMOVED***

	// ModuleInstance represents an instance of the execution module.
	ModuleInstance struct ***REMOVED***
		modules.InstanceCore
		obj *goja.Object
	***REMOVED***
)

var (
	_ modules.IsModuleV2 = &RootModule***REMOVED******REMOVED***
	_ modules.Instance   = &ModuleInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.IsModuleV2 interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(m modules.InstanceCore) modules.Instance ***REMOVED***
	mi := &ModuleInstance***REMOVED***InstanceCore: m***REMOVED***
	rt := m.GetRuntime()
	o := rt.NewObject()
	defProp := func(name string, newInfo func() (*goja.Object, error)) ***REMOVED***
		err := o.DefineAccessorProperty(name, rt.ToValue(func() goja.Value ***REMOVED***
			obj, err := newInfo()
			if err != nil ***REMOVED***
				common.Throw(rt, err)
			***REMOVED***
			return obj
		***REMOVED***), nil, goja.FLAG_FALSE, goja.FLAG_TRUE)
		if err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
	***REMOVED***
	defProp("scenario", mi.newScenarioInfo)
	defProp("instance", mi.newInstanceInfo)
	defProp("vu", mi.newVUInfo)

	mi.obj = o

	return mi
***REMOVED***

// GetExports returns the exports of the execution module.
func (mi *ModuleInstance) GetExports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***Default: mi.obj***REMOVED***
***REMOVED***

// newScenarioInfo returns a goja.Object with property accessors to retrieve
// information about the scenario the current VU is running in.
func (mi *ModuleInstance) newScenarioInfo() (*goja.Object, error) ***REMOVED***
	ctx := mi.GetContext()
	rt := common.GetRuntime(ctx)
	vuState := mi.GetState()
	if vuState == nil ***REMOVED***
		return nil, errors.New("getting scenario information in the init context is not supported")
	***REMOVED***
	if rt == nil ***REMOVED***
		return nil, errors.New("goja runtime is nil in context")
	***REMOVED***
	getScenarioState := func() *lib.ScenarioState ***REMOVED***
		ss := lib.GetScenarioState(mi.GetContext())
		if ss == nil ***REMOVED***
			common.Throw(rt, errors.New("getting scenario information in the init context is not supported"))
		***REMOVED***
		return ss
	***REMOVED***

	si := map[string]func() interface***REMOVED******REMOVED******REMOVED***
		"name": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return getScenarioState().Name
		***REMOVED***,
		"executor": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return getScenarioState().Executor
		***REMOVED***,
		"startTime": func() interface***REMOVED******REMOVED*** ***REMOVED***
			//nolint:lll
			// Return the timestamp in milliseconds, since that's how JS
			// timestamps usually are:
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/Date#time_value_or_timestamp_number
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/now#return_value
			return getScenarioState().StartTime.UnixNano() / int64(time.Millisecond)
		***REMOVED***,
		"progress": func() interface***REMOVED******REMOVED*** ***REMOVED***
			p, _ := getScenarioState().ProgressFn()
			return p
		***REMOVED***,
		"iterationInInstance": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return vuState.GetScenarioLocalVUIter()
		***REMOVED***,
		"iterationInTest": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return vuState.GetScenarioGlobalVUIter()
		***REMOVED***,
	***REMOVED***

	return newInfoObj(rt, si)
***REMOVED***

// newInstanceInfo returns a goja.Object with property accessors to retrieve
// information about the local instance stats.
func (mi *ModuleInstance) newInstanceInfo() (*goja.Object, error) ***REMOVED***
	ctx := mi.GetContext()
	es := lib.GetExecutionState(ctx)
	if es == nil ***REMOVED***
		return nil, errors.New("getting instance information in the init context is not supported")
	***REMOVED***

	rt := common.GetRuntime(ctx)
	if rt == nil ***REMOVED***
		return nil, errors.New("goja runtime is nil in context")
	***REMOVED***

	ti := map[string]func() interface***REMOVED******REMOVED******REMOVED***
		"currentTestRunDuration": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return float64(es.GetCurrentTestRunDuration()) / float64(time.Millisecond)
		***REMOVED***,
		"iterationsCompleted": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return es.GetFullIterationCount()
		***REMOVED***,
		"iterationsInterrupted": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return es.GetPartialIterationCount()
		***REMOVED***,
		"vusActive": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return es.GetCurrentlyActiveVUsCount()
		***REMOVED***,
		"vusInitialized": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return es.GetInitializedVUsCount()
		***REMOVED***,
	***REMOVED***

	return newInfoObj(rt, ti)
***REMOVED***

// newVUInfo returns a goja.Object with property accessors to retrieve
// information about the currently executing VU.
func (mi *ModuleInstance) newVUInfo() (*goja.Object, error) ***REMOVED***
	ctx := mi.GetContext()
	vuState := lib.GetState(ctx)
	if vuState == nil ***REMOVED***
		return nil, errors.New("getting VU information in the init context is not supported")
	***REMOVED***

	rt := common.GetRuntime(ctx)
	if rt == nil ***REMOVED***
		return nil, errors.New("goja runtime is nil in context")
	***REMOVED***

	vi := map[string]func() interface***REMOVED******REMOVED******REMOVED***
		"idInInstance":        func() interface***REMOVED******REMOVED*** ***REMOVED*** return vuState.VUID ***REMOVED***,
		"idInTest":            func() interface***REMOVED******REMOVED*** ***REMOVED*** return vuState.VUIDGlobal ***REMOVED***,
		"iterationInInstance": func() interface***REMOVED******REMOVED*** ***REMOVED*** return vuState.Iteration ***REMOVED***,
		"iterationInScenario": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return vuState.GetScenarioVUIter()
		***REMOVED***,
	***REMOVED***

	return newInfoObj(rt, vi)
***REMOVED***

func newInfoObj(rt *goja.Runtime, props map[string]func() interface***REMOVED******REMOVED***) (*goja.Object, error) ***REMOVED***
	o := rt.NewObject()

	for p, get := range props ***REMOVED***
		err := o.DefineAccessorProperty(p, rt.ToValue(get), nil, goja.FLAG_FALSE, goja.FLAG_TRUE)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return o, nil
***REMOVED***
