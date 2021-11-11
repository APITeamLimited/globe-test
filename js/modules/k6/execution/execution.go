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
	"fmt"
	"reflect"
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
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &ModuleInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
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

// Exports returns the exports of the execution module.
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
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

	o, err := newInfoObj(rt, vi)
	if err != nil ***REMOVED***
		return o, err
	***REMOVED***

	err = o.Set("tags", rt.NewDynamicObject(&tagsDynamicObject***REMOVED***
		Runtime: rt,
		State:   vuState,
	***REMOVED***))
	return o, err
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

type tagsDynamicObject struct ***REMOVED***
	Runtime *goja.Runtime
	State   *lib.State
***REMOVED***

// Get a property value for the key. May return nil if the property does not exist.
func (o *tagsDynamicObject) Get(key string) goja.Value ***REMOVED***
	tag, ok := o.State.Tags.Get(key)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return o.Runtime.ToValue(tag)
***REMOVED***

// Set a property value for the key. It returns true if succeed.
// String, Boolean and Number types are implicitly converted
// to the goja's relative string representation.
// In any other case, if the Throw option is set then an error is raised
// otherwise just a Warning is written.
func (o *tagsDynamicObject) Set(key string, val goja.Value) bool ***REMOVED***
	switch val.ExportType().Kind() ***REMOVED*** //nolint:exhaustive
	case
		reflect.String,
		reflect.Bool,
		reflect.Int64,
		reflect.Float64:

		o.State.Tags.Set(key, val.String())
		return true
	default:
		err := fmt.Errorf("only String, Boolean and Number types are accepted as a Tag value")
		if o.State.Options.Throw.Bool ***REMOVED***
			common.Throw(o.Runtime, err)
			return false
		***REMOVED***
		o.State.Logger.Warnf("the execution.vu.tags.Set('%s') operation has been discarded because %s", key, err.Error())
		return false
	***REMOVED***
***REMOVED***

// Has returns true if the property exists.
func (o *tagsDynamicObject) Has(key string) bool ***REMOVED***
	_, ok := o.State.Tags.Get(key)
	return ok
***REMOVED***

// Delete deletes the property for the key. It returns true on success (note, that includes missing property).
func (o *tagsDynamicObject) Delete(key string) bool ***REMOVED***
	o.State.Tags.Delete(key)
	return true
***REMOVED***

// Keys returns a slice with all existing property keys. The order is not deterministic.
func (o *tagsDynamicObject) Keys() []string ***REMOVED***
	if o.State.Tags.Len() < 1 ***REMOVED***
		return nil
	***REMOVED***

	tags := o.State.Tags.Clone()
	keys := make([]string, 0, len(tags))
	for k := range tags ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	return keys
***REMOVED***
