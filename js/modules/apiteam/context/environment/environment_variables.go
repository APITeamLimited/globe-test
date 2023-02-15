package environment

import (
	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/dop251/goja"
)

func (mi *EnvironmentInstance) getVariablesObject() *goja.Object {
	// Make sure variables are enabled before calling

	rt := mi.vu.Runtime()
	variablesObject := rt.NewObject()

	mustExport := func(name string, value interface{}) {
		if err := variablesObject.Set(name, value); err != nil {
			common.Throw(rt, err)
		}
	}

	mustExport("set", mi.set)
	mustExport("get", mi.get)
	mustExport("has", mi.has)
	mustExport("unset", mi.unset)
	mustExport("clear", mi.clear)
	mustExport("list", mi.list)

	return variablesObject
}

// set sets a key-value pair in the environment.
func (mi *EnvironmentInstance) set(key string, value string) bool {
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.module.sharedEnvironment.data.Variables[key] = value

	return true
}

// get gets a value from the environment.
func (mi *EnvironmentInstance) get(key string) string {
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	return mi.module.sharedEnvironment.data.Variables[key]
}

// has checks if a key exists in the environment.
func (mi *EnvironmentInstance) has(key string) bool {
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	_, ok := mi.module.sharedEnvironment.data.Variables[key]

	return ok
}

// unset removes a key-value pair from the environment.
func (mi *EnvironmentInstance) unset(key string) bool {
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	if _, ok := mi.module.sharedEnvironment.data.Variables[key]; ok {
		delete(mi.module.sharedEnvironment.data.Variables, key)
		return true
	}

	return false
}

// clear removes all key-value pairs from the environment.
func (mi *EnvironmentInstance) clear() bool {
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	mi.module.sharedEnvironment.data.Variables = make(map[string]string)

	return true
}

// list returns a dictionary of all key-value pairs in the environment.
func (mi *EnvironmentInstance) list() map[string]string {
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	return mi.module.sharedEnvironment.data.Variables
}
