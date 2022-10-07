package environment

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create environment module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	EnvironmentModule struct ***REMOVED***
		sharedEnvironment sharedEnvironment
	***REMOVED***

	// EnvironmentInstance represents an instance of the environment module.
	EnvironmentInstance struct ***REMOVED***
		vu     modules.VU
		module *EnvironmentModule
	***REMOVED***

	sharedEnvironment struct ***REMOVED***
		isEnabled bool
		data      map[string]string

		mu *sync.RWMutex
	***REMOVED***
)

var (
	_ modules.Module   = &EnvironmentModule***REMOVED******REMOVED***
	_ modules.Instance = &EnvironmentInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new EnvironmentModule instance.
func New(workerInfo *libWorker.WorkerInfo) *EnvironmentModule ***REMOVED***
	marshalled, _ := json.Marshal(workerInfo)
	fmt.Println("new environment module", string(marshalled))

	// Check environment actually exists
	if workerInfo.Environment != nil ***REMOVED***
		return &EnvironmentModule***REMOVED***
			sharedEnvironment: sharedEnvironment***REMOVED***
				isEnabled: true,
				data:      workerInfo.Environment,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return &EnvironmentModule***REMOVED***
		sharedEnvironment: sharedEnvironment***REMOVED***
			isEnabled: false,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an environment module instance for each VU.
func (module *EnvironmentModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &EnvironmentInstance***REMOVED***
		vu:     vu,
		module: module,
	***REMOVED***
***REMOVED***

// Exports returns the exports of the environment module.
func (mi *EnvironmentInstance) Exports() modules.Exports ***REMOVED***
	if !mi.module.sharedEnvironment.isEnabled ***REMOVED***
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
func (mi *EnvironmentInstance) isEnabled() bool ***REMOVED***
	return mi.module.sharedEnvironment.isEnabled
***REMOVED***

// set sets a key-value pair in the environment.
func (mi *EnvironmentInstance) set(key string, value string) bool ***REMOVED***
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.module.sharedEnvironment.data[key] = value
	return true
***REMOVED***

// get gets a value from the environment.
func (mi *EnvironmentInstance) get(key string) string ***REMOVED***
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	if value, ok := mi.module.sharedEnvironment.data[key]; ok ***REMOVED***
		return value
	***REMOVED***

	return ""
***REMOVED***

// has checks if a key exists in the environment.
func (mi *EnvironmentInstance) has(key string) bool ***REMOVED***
	// Check if environment is enabled
	if !mi.module.sharedEnvironment.isEnabled ***REMOVED***
		return false
	***REMOVED***

	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	_, ok := mi.module.sharedEnvironment.data[key]
	return ok
***REMOVED***

// unset removes a key-value pair from the environment.
func (mi *EnvironmentInstance) unset(key string) bool ***REMOVED***
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	if _, ok := mi.module.sharedEnvironment.data[key]; ok ***REMOVED***
		delete(mi.module.sharedEnvironment.data, key)
		return true
	***REMOVED***

	return false
***REMOVED***

// clear removes all key-value pairs from the environment.
func (mi *EnvironmentInstance) clear() bool ***REMOVED***
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	mi.module.sharedEnvironment.data = make(map[string]string)
	return true
***REMOVED***

// list returns a list of all key-value pairs in the environment.
func (mi *EnvironmentInstance) list() []string ***REMOVED***
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	list := make([]string, 0, len(mi.module.sharedEnvironment.data))
	for _, item := range mi.module.sharedEnvironment.data ***REMOVED***
		list = append(list, item)
	***REMOVED***

	return list
***REMOVED***
