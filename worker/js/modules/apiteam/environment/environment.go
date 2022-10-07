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
	EnvironmentModule struct {
		sharedEnvironment sharedEnvironment
	}

	// EnvironmentInstance represents an instance of the environment module.
	EnvironmentInstance struct {
		vu     modules.VU
		module *EnvironmentModule
	}

	sharedEnvironment struct {
		isEnabled bool
		data      map[string]string

		mu *sync.RWMutex
	}
)

var (
	_ modules.Module   = &EnvironmentModule{}
	_ modules.Instance = &EnvironmentInstance{}
)

// New returns a pointer to a new EnvironmentModule instance.
func New(workerInfo *libWorker.WorkerInfo) *EnvironmentModule {
	marshalled, _ := json.Marshal(workerInfo)
	fmt.Println("new environment module", string(marshalled))

	// Check environment actually exists
	if workerInfo.Environment != nil {
		return &EnvironmentModule{
			sharedEnvironment: sharedEnvironment{
				isEnabled: true,
				data:      workerInfo.Environment,
			},
		}
	}

	return &EnvironmentModule{
		sharedEnvironment: sharedEnvironment{
			isEnabled: false,
		},
	}
}

// NewModuleInstance returns an environment module instance for each VU.
func (module *EnvironmentModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &EnvironmentInstance{
		vu:     vu,
		module: module,
	}
}

// Exports returns the exports of the environment module.
func (mi *EnvironmentInstance) Exports() modules.Exports {
	if !mi.module.sharedEnvironment.isEnabled {
		return modules.Exports{
			Named: map[string]interface{}{
				"isEnabled": mi.isEnabled,
			},
		}
	}

	return modules.Exports{
		Named: map[string]interface{}{
			"isEnabled": mi.isEnabled,
			"set":       mi.set,
			"get":       mi.get,
			"has":       mi.has,
			"unset":     mi.unset,
			"clear":     mi.clear,
			"list":      mi.list,
		},
	}
}

// isEnabled is a getter for the isEnabled property.
func (mi *EnvironmentInstance) isEnabled() bool {
	return mi.module.sharedEnvironment.isEnabled
}

// set sets a key-value pair in the environment.
func (mi *EnvironmentInstance) set(key string, value string) bool {
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.module.sharedEnvironment.data[key] = value
	return true
}

// get gets a value from the environment.
func (mi *EnvironmentInstance) get(key string) string {
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	if value, ok := mi.module.sharedEnvironment.data[key]; ok {
		return value
	}

	return ""
}

// has checks if a key exists in the environment.
func (mi *EnvironmentInstance) has(key string) bool {
	// Check if environment is enabled
	if !mi.module.sharedEnvironment.isEnabled {
		return false
	}

	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	_, ok := mi.module.sharedEnvironment.data[key]
	return ok
}

// unset removes a key-value pair from the environment.
func (mi *EnvironmentInstance) unset(key string) bool {
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	if _, ok := mi.module.sharedEnvironment.data[key]; ok {
		delete(mi.module.sharedEnvironment.data, key)
		return true
	}

	return false
}

// clear removes all key-value pairs from the environment.
func (mi *EnvironmentInstance) clear() bool {
	mi.module.sharedEnvironment.mu.Lock()
	defer mi.module.sharedEnvironment.mu.Unlock()

	mi.module.sharedEnvironment.data = make(map[string]string)
	return true
}

// list returns a list of all key-value pairs in the environment.
func (mi *EnvironmentInstance) list() []string {
	mi.module.sharedEnvironment.mu.RLock()
	defer mi.module.sharedEnvironment.mu.RUnlock()

	list := make([]string, 0, len(mi.module.sharedEnvironment.data))
	for _, item := range mi.module.sharedEnvironment.data {
		list = append(list, item)
	}

	return list
}
