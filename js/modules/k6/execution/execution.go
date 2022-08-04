package execution

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/errext"
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
		vu  modules.VU
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
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	mi := &ModuleInstance***REMOVED***vu: vu***REMOVED***
	rt := vu.Runtime()
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
	defProp("instance", mi.newInstanceInfo)
	defProp("scenario", mi.newScenarioInfo)
	defProp("test", mi.newTestInfo)
	defProp("vu", mi.newVUInfo)

	mi.obj = o

	return mi
***REMOVED***

// Exports returns the exports of the execution module.
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***Default: mi.obj***REMOVED***
***REMOVED***

var errRunInInitContext = errors.New("getting scenario information outside of the VU context is not supported")

// newScenarioInfo returns a goja.Object with property accessors to retrieve
// information about the scenario the current VU is running in.
func (mi *ModuleInstance) newScenarioInfo() (*goja.Object, error) ***REMOVED***
	rt := mi.vu.Runtime()
	vuState := mi.vu.State()
	if vuState == nil ***REMOVED***
		return nil, errRunInInitContext
	***REMOVED***
	getScenarioState := func() *lib.ScenarioState ***REMOVED***
		ss := lib.GetScenarioState(mi.vu.Context())
		if ss == nil ***REMOVED***
			common.Throw(rt, errRunInInitContext)
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
			if vuState.GetScenarioLocalVUIter == nil ***REMOVED***
				common.Throw(rt, errRunInInitContext)
			***REMOVED***

			return vuState.GetScenarioLocalVUIter()
		***REMOVED***,
		"iterationInTest": func() interface***REMOVED******REMOVED*** ***REMOVED***
			if vuState.GetScenarioGlobalVUIter == nil ***REMOVED***
				common.Throw(rt, errRunInInitContext)
			***REMOVED***

			return vuState.GetScenarioGlobalVUIter()
		***REMOVED***,
	***REMOVED***

	return newInfoObj(rt, si)
***REMOVED***

// newInstanceInfo returns a goja.Object with property accessors to retrieve
// information about the local instance stats.
func (mi *ModuleInstance) newInstanceInfo() (*goja.Object, error) ***REMOVED***
	es := lib.GetExecutionState(mi.vu.Context())
	if es == nil ***REMOVED***
		return nil, errors.New("getting instance information in the init context is not supported")
	***REMOVED***
	rt := mi.vu.Runtime()

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

// newTestInfo returns a goja.Object with property accessors to retrieve
// information and control execution of the overall test run.
func (mi *ModuleInstance) newTestInfo() (*goja.Object, error) ***REMOVED***
	// the cache of goja.Object in the optimal parsed form
	// for the consolidated and derived lib.Options
	var optionsObject *goja.Object
	rt := mi.vu.Runtime()
	ti := map[string]func() interface***REMOVED******REMOVED******REMOVED***
		// stop the test run
		"abort": func() interface***REMOVED******REMOVED*** ***REMOVED***
			return func(msg goja.Value) ***REMOVED***
				reason := errext.AbortTest
				if msg != nil && !goja.IsUndefined(msg) ***REMOVED***
					reason = fmt.Sprintf("%s: %s", reason, msg.String())
				***REMOVED***
				rt.Interrupt(&errext.InterruptError***REMOVED***Reason: reason***REMOVED***)
			***REMOVED***
		***REMOVED***,
		"options": func() interface***REMOVED******REMOVED*** ***REMOVED***
			if optionsObject == nil ***REMOVED***
				opts, err := optionsAsObject(rt, mi.vu.State().Options)
				if err != nil ***REMOVED***
					common.Throw(rt, err)
				***REMOVED***
				optionsObject = opts
			***REMOVED***
			return optionsObject
		***REMOVED***,
	***REMOVED***

	return newInfoObj(rt, ti)
***REMOVED***

// newVUInfo returns a goja.Object with property accessors to retrieve
// information about the currently executing VU.
func (mi *ModuleInstance) newVUInfo() (*goja.Object, error) ***REMOVED***
	vuState := mi.vu.State()
	if vuState == nil ***REMOVED***
		return nil, errors.New("getting VU information in the init context is not supported")
	***REMOVED***
	rt := mi.vu.Runtime()

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

// optionsAsObject maps the lib.Options struct that contains the consolidated
// and derived options configuration in a goja.Object.
//
// When values are not set then the default value returned from JSON is used.
// Most of the lib.Options are Nullable types so they will be null on default.
func optionsAsObject(rt *goja.Runtime, options lib.Options) (*goja.Object, error) ***REMOVED***
	b, err := json.Marshal(options)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to encode the lib.Options as json: %w", err)
	***REMOVED***

	// Using the native JS parser function guarantees getting
	// the supported types for deep freezing the complex object.
	jsonParse, _ := goja.AssertFunction(rt.GlobalObject().Get("JSON").ToObject(rt).Get("parse"))
	parsed, err := jsonParse(goja.Undefined(), rt.ToValue(string(b)))
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	obj := parsed.ToObject(rt)

	mustDelete := func(prop string) ***REMOVED***
		delErr := obj.Delete(prop)
		if err != nil ***REMOVED***
			common.Throw(rt, delErr)
		***REMOVED***
	***REMOVED***
	mustSetReadOnlyProperty := func(k string, v interface***REMOVED******REMOVED***) ***REMOVED***
		defErr := obj.DefineDataProperty(k, rt.ToValue(v), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
		if err != nil ***REMOVED***
			common.Throw(rt, defErr)
		***REMOVED***
	***REMOVED***

	mustDelete("vus")
	mustDelete("iterations")
	mustDelete("duration")
	mustDelete("stages")

	consoleOutput := goja.Null()
	if options.ConsoleOutput.Valid ***REMOVED***
		consoleOutput = rt.ToValue(options.ConsoleOutput.String)
	***REMOVED***
	mustSetReadOnlyProperty("consoleOutput", consoleOutput)

	localIPs := goja.Null()
	if options.LocalIPs.Valid ***REMOVED***
		raw, marshalErr := options.LocalIPs.MarshalText()
		if err != nil ***REMOVED***
			common.Throw(rt, marshalErr)
		***REMOVED***
		localIPs = rt.ToValue(string(raw))
	***REMOVED***
	mustSetReadOnlyProperty("localIPs", localIPs)

	err = common.FreezeObject(rt, obj)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	return obj, nil
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
	kind := reflect.Invalid
	if typ := val.ExportType(); typ != nil ***REMOVED***
		kind = typ.Kind()
	***REMOVED***
	switch kind ***REMOVED***
	case
		reflect.String,
		reflect.Bool,
		reflect.Int64,
		reflect.Float64:

		o.State.Tags.Set(key, val.String())
		return true
	default:
		reason := "only String, Boolean and Number types are accepted as a Tag value"
		if o.State.Options.Throw.Bool ***REMOVED***
			panic(o.Runtime.NewTypeError(reason))
		***REMOVED***
		o.State.Logger.Warnf("the execution.vu.tags.Set('%s') operation has been discarded because %s", key, reason)
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
