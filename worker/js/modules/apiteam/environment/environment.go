package environment

import (
	"sync"

	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create environment module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct {
		env sharedEnvironment
	}

	// Environment represents an instance of the environment module.
	Environment struct {
		vu          modules.VU
		environment *sharedEnvironment
	}

	sharedEnvironment struct {
		isEnabled bool
		data      map[string]libWorker.KeyValueItem
		mu        sync.RWMutex
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Environment{}
)

// New returns a pointer to a new RootModule instance.
func New(workerInfo *libWorker.WorkerInfo) *RootModule {
	// Check environment actually exists
	if workerInfo.Environment != nil {
		return &RootModule{
			env: sharedEnvironment{
				isEnabled: true,
				data:      make(map[string]libWorker.KeyValueItem),
			},
		}
	}

	return &RootModule{
		env: sharedEnvironment{
			isEnabled: false,
		},
	}
}

// NewModuleInstance returns an environment module instance for each VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Environment{
		vu:          vu,
		environment: &rm.env,
	}
}

// Exports returns the exports of the environment module.
func (mi *Environment) Exports() modules.Exports {
	if !mi.environment.isEnabled {
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
func (mi *Environment) isEnabled() bool {
	return mi.environment.isEnabled
}

// set sets a key-value pair in the environment.
func (mi *Environment) set(key string, value string) bool {
	mi.environment.mu.Lock()
	defer mi.environment.mu.Unlock()

	// Overwrite existing value if key already exists
	mi.environment.data[key] = libWorker.KeyValueItem{
		Key:   key,
		Value: value,
	}

	return true
}

// get gets a value from the environment.
func (mi *Environment) get(key string) string {
	mi.environment.mu.RLock()
	defer mi.environment.mu.RUnlock()

	if value, ok := mi.environment.data[key]; ok {
		return value.Value
	}

	return ""
}

// has checks if a key exists in the environment.
func (mi *Environment) has(key string) bool {
	// Check if environment is enabled
	if !mi.environment.isEnabled {
		return false
	}

	mi.environment.mu.RLock()
	defer mi.environment.mu.RUnlock()

	_, ok := mi.environment.data[key]
	return ok
}

// unset removes a key-value pair from the environment.
func (mi *Environment) unset(key string) bool {
	mi.environment.mu.Lock()
	defer mi.environment.mu.Unlock()

	if _, ok := mi.environment.data[key]; ok {
		delete(mi.environment.data, key)
		return true
	}

	return false
}

// clear removes all key-value pairs from the environment.
func (mi *Environment) clear() bool {
	mi.environment.mu.Lock()
	defer mi.environment.mu.Unlock()

	mi.environment.data = make(map[string]libWorker.KeyValueItem)
	return true
}

// list returns a list of all key-value pairs in the environment.
func (mi *Environment) list() []libWorker.KeyValueItem {
	mi.environment.mu.RLock()
	defer mi.environment.mu.RUnlock()

	list := make([]libWorker.KeyValueItem, 0, len(mi.environment.data))
	for _, item := range mi.environment.data {
		list = append(list, item)
	}

	return list
}
