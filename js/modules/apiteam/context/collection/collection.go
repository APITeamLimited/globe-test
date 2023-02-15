package collection

import (
	"fmt"
	"sync"

	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/APITeamLimited/globe-test/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
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
		isEnabled        bool
		setEnabled       bool
		sharedCollection sharedCollection
	}

	// CollectionInstance represents an instance of the collection module.
	CollectionInstance struct {
		vu            modules.VU
		module        *CollectionModule
		defaultExport *goja.Object
	}

	sharedCollection struct {
		data *libWorker.Collection

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
			isEnabled: true,
			sharedCollection: sharedCollection{
				data: workerInfo.Collection,
				mu:   &sync.RWMutex{},
			},
		}
	}

	return &CollectionModule{
		isEnabled:        false,
		setEnabled:       false,
		sharedCollection: sharedCollection{},
	}
}

// NewModuleInstance returns an collection module instance for each VU.
func (module *CollectionModule) NewModuleInstance(vu modules.VU) modules.Instance {
	rt := vu.Runtime()

	mi := &CollectionInstance{
		vu:            vu,
		module:        module,
		defaultExport: rt.NewObject(),
	}

	mi.defaultExport.DefineDataProperty(
		"enabled", rt.ToValue(module.isEnabled), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)

	if module.isEnabled {
		mi.defaultExport.DefineDataProperty(
			"name", rt.ToValue(module.sharedCollection.data.Name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)

		if err := mi.defaultExport.Set("variables", mi.getVariablesObject()); err != nil {
			fmt.Println("Error setting collection variables object: ", err)
			common.Throw(rt, err)
		}
	}

	return mi
}

// Exports returns the exports of the collection module.
func (mi *CollectionInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.defaultExport,
	}
}
