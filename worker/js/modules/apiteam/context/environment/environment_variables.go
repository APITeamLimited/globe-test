package environment

import (
	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/dop251/goja"
)

func (mi *EnvironmentInstance) getVariablesObject() *goja.Object ***REMOVED***
	// Make sure variables are enabled before calling

	rt := mi.vu.Runtime()
	variablesObject := rt.NewObject()

	mustExport := func(name string, value interface***REMOVED******REMOVED***) ***REMOVED***
		if err := variablesObject.Set(name, value); err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
	***REMOVED***

	mustExport("set", mi.set)
	mustExport("get", mi.get)
	mustExport("has", mi.has)
	mustExport("unset", mi.unset)
	mustExport("clear", mi.clear)
	mustExport("list", mi.list)

	return variablesObject
***REMOVED***

// set sets a key-value pair in the environment.
func (mi *EnvironmentInstance) set(key string, value string) bool ***REMOVED***
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.module.sharedEnvironment.data.Variables[key] = value

	return true
***REMOVED***

// get gets a value from the environment.
func (mi *EnvironmentInstance) get(key string) string ***REMOVED***
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	return mi.module.sharedEnvironment.data.Variables[key]
***REMOVED***

// has checks if a key exists in the environment.
func (mi *EnvironmentInstance) has(key string) bool ***REMOVED***
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	_, ok := mi.module.sharedEnvironment.data.Variables[key]

	return ok
***REMOVED***

// unset removes a key-value pair from the environment.
func (mi *EnvironmentInstance) unset(key string) bool ***REMOVED***
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	if _, ok := mi.module.sharedEnvironment.data.Variables[key]; ok ***REMOVED***
		delete(mi.module.sharedEnvironment.data.Variables, key)
		return true
	***REMOVED***

	return false
***REMOVED***

// clear removes all key-value pairs from the environment.
func (mi *EnvironmentInstance) clear() bool ***REMOVED***
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	mi.module.sharedEnvironment.data.Variables = make(map[string]string)

	return true
***REMOVED***

// list returns a dictionary of all key-value pairs in the environment.
func (mi *EnvironmentInstance) list() map[string]string ***REMOVED***
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	return mi.module.sharedEnvironment.data.Variables
***REMOVED***
