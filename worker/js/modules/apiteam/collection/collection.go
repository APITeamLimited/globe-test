package collection

import (
	"sync"

	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create collection module instances for each VU.
//
// TODO: add sync.Once for all of the deprecation warnings we might want to do
// for the old k6/http APIs here, so they are shown only once in a test run.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct {
		col sharedCollection
	}

	// Collection represents an instance of the collection module.
	Collection struct {
		vu         modules.VU
		collection *sharedCollection
	}

	sharedCollection struct {
		isEnabled bool
		variables map[string]libWorker.KeyValueItem
		mu        sync.RWMutex
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Collection{}
)

// New returns a pointer to a new RootModule instance.
func New(workerInfo *libWorker.WorkerInfo) *RootModule {
	// Check collection actually exists
	if workerInfo.Collection != nil {
		return &RootModule{
			col: sharedCollection{
				isEnabled: true,
				variables: *workerInfo.Collection.Variables,
			},
		}
	} else {
		return &RootModule{
			col: sharedCollection{
				isEnabled: false,
			},
		}
	}
}

// NewModuleInstance returns an collection module instance for each VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Collection{
		vu:         vu,
		collection: &rm.col,
	}
}

// Exports returns the exports of the collection module.
func (mi *Collection) Exports() modules.Exports {
	if !mi.collection.isEnabled {
		return modules.Exports{
			Named: map[string]interface{}{
				"isEnabled": mi.isEnabled,
			},
		}
	}

	return modules.Exports{
		Named: map[string]interface{}{
			"isEnabled": mi.isEnabled,
			"variables": mi.getVariables,
		},
	}
}

// isEnabled returns whether the collection is enabled or not.
func (mi *Collection) isEnabled() bool {
	return mi.collection.isEnabled
}
