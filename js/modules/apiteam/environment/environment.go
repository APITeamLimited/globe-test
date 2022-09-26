package environment

import (
	"sync"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create environment module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct ***REMOVED***
		env sharedEnvironment
	***REMOVED***

	// Environment represents an instance of the environment module.
	Environment struct ***REMOVED***
		vu          modules.VU
		environment *sharedEnvironment
	***REMOVED***

	sharedEnvironment struct ***REMOVED***
		isEnabled bool
		data      map[string]lib.KeyValueItem
		mu        sync.RWMutex
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &Environment***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New(workerInfo *lib.WorkerInfo) *RootModule ***REMOVED***
	// Check environment actually exists
	if workerInfo.Environment != nil ***REMOVED***
		return &RootModule***REMOVED***
			env: sharedEnvironment***REMOVED***
				isEnabled: true,
				data:      make(map[string]lib.KeyValueItem),
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return &RootModule***REMOVED***
		env: sharedEnvironment***REMOVED***
			isEnabled: false,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an environment module instance for each VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &Environment***REMOVED***
		vu:          vu,
		environment: &rm.env,
	***REMOVED***
***REMOVED***

// Exports returns the exports of the environment module.
func (mi *Environment) Exports() modules.Exports ***REMOVED***
	if !mi.environment.isEnabled ***REMOVED***
		return modules.Exports***REMOVED***
			Named: map[string]interface***REMOVED******REMOVED******REMOVED***
				"isEnabled": mi.isEnabled,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"isEnabled": mi.isEnabled,
			"set":       mi.set,
			"get":       mi.get,
			"has":       mi.has,
			"unset":     mi.unset,
			"clear":     mi.clear,
			"list":      mi.list,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// isEnabled is a getter for the isEnabled property.
func (mi *Environment) isEnabled() bool ***REMOVED***
	return mi.environment.isEnabled
***REMOVED***

// set sets a key-value pair in the environment.
func (mi *Environment) set(key string, value string) bool ***REMOVED***
	mi.environment.mu.Lock()
	defer mi.environment.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.environment.data[key] = lib.KeyValueItem***REMOVED***
		Key:   key,
		Value: value,
	***REMOVED***

	return true
***REMOVED***

// get gets a value from the environment.
func (mi *Environment) get(key string) string ***REMOVED***
	mi.environment.mu.RLock()
	defer mi.environment.mu.RUnlock()

	if value, ok := mi.environment.data[key]; ok ***REMOVED***
		return value.Value
	***REMOVED***

	return ""
***REMOVED***

// has checks if a key exists in the environment.
func (mi *Environment) has(key string) bool ***REMOVED***
	// Check if environment is enabled
	if !mi.environment.isEnabled ***REMOVED***
		return false
	***REMOVED***

	mi.environment.mu.RLock()
	defer mi.environment.mu.RUnlock()

	_, ok := mi.environment.data[key]
	return ok
***REMOVED***

// unset removes a key-value pair from the environment.
func (mi *Environment) unset(key string) bool ***REMOVED***
	mi.environment.mu.Lock()
	defer mi.environment.mu.Unlock()

	if _, ok := mi.environment.data[key]; ok ***REMOVED***
		delete(mi.environment.data, key)
		return true
	***REMOVED***

	return false
***REMOVED***

// clear removes all key-value pairs from the environment.
func (mi *Environment) clear() bool ***REMOVED***
	mi.environment.mu.Lock()
	defer mi.environment.mu.Unlock()

	mi.environment.data = make(map[string]lib.KeyValueItem)
	return true
***REMOVED***

// list returns a list of all key-value pairs in the environment.
func (mi *Environment) list() []lib.KeyValueItem ***REMOVED***
	mi.environment.mu.RLock()
	defer mi.environment.mu.RUnlock()

	list := make([]lib.KeyValueItem, 0, len(mi.environment.data))
	for _, item := range mi.environment.data ***REMOVED***
		list = append(list, item)
	***REMOVED***

	return list
***REMOVED***
