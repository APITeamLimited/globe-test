package collection

import (
	"sync"

	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

// CollectionModule is the global module object type. It is instantiated once per test
// run and will be used to create collection module instances for each VU.
//
// TODO: add sync.Once for all of the deprecation warnings we might want to do
// for the old k6/http APIs here, so they are shown only once in a test run.
type (
	// CollectionModule is the global module instance that will create module
	// instances for each VU.
	CollectionModule struct {
		sharedCollection sharedCollection
	}

	// CollectionInstance represents an instance of the collection module.
	CollectionInstance struct {
		vu     modules.VU
		module *CollectionModule
	}

	sharedCollection struct {
		isEnabled bool
		data      libWorker.Collection

		mu *sync.RWMutex
	}
)

var (
	_ modules.Module   = &CollectionModule{}
	_ modules.Instance = &CollectionInstance{}
)

// New returns a pointer to a new CollectionModule instance.
func New(workerInfo *libWorker.WorkerInfo) *CollectionModule {
	// Check collection actually exists
	if workerInfo.Collection != nil {
		return &CollectionModule{
			sharedCollection: sharedCollection{
				isEnabled: true,
				data:      *workerInfo.Collection,
			},
		}
	} else {
		return &CollectionModule{
			sharedCollection: sharedCollection{
				isEnabled: false,
			},
		}
	}
}

// NewModuleInstance returns an collection module instance for each VU.
func (module *CollectionModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &CollectionInstance{
		vu,
		module,
	}
}

// Exports returns the exports of the collection module.
func (mi *CollectionInstance) Exports() modules.Exports {
	if !mi.module.sharedCollection.isEnabled {
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
func (mi *CollectionInstance) isEnabled() bool {
	return mi.module.sharedCollection.isEnabled
}
