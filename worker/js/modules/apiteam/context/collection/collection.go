package collection

import (
	"fmt"
	"sync"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
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
	CollectionModule struct ***REMOVED***
		isEnabled        bool
		setEnabled       bool
		sharedCollection sharedCollection
	***REMOVED***

	// CollectionInstance represents an instance of the collection module.
	CollectionInstance struct ***REMOVED***
		vu            modules.VU
		module        *CollectionModule
		defaultExport *goja.Object
	***REMOVED***

	sharedCollection struct ***REMOVED***
		data *libWorker.Collection

		mu *sync.RWMutex
	***REMOVED***
)

var (
	_ modules.Module   = &CollectionModule***REMOVED******REMOVED***
	_ modules.Instance = &CollectionInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new CollectionModule instance.
func New(workerInfo *libWorker.WorkerInfo) *CollectionModule ***REMOVED***
	// Check collection actually exists
	if workerInfo.Collection != nil ***REMOVED***
		return &CollectionModule***REMOVED***
			isEnabled: true,
			sharedCollection: sharedCollection***REMOVED***
				data: workerInfo.Collection,
				mu:   &sync.RWMutex***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return &CollectionModule***REMOVED***
			isEnabled:        false,
			setEnabled:       false,
			sharedCollection: sharedCollection***REMOVED******REMOVED***,
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an collection module instance for each VU.
func (module *CollectionModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	rt := vu.Runtime()

	mi := &CollectionInstance***REMOVED***
		vu:            vu,
		module:        module,
		defaultExport: rt.NewObject(),
	***REMOVED***

	mi.defaultExport.DefineDataProperty(
		"enabled", rt.ToValue(module.isEnabled), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)

	if module.isEnabled ***REMOVED***
		mi.defaultExport.DefineDataProperty(
			"name", rt.ToValue(module.sharedCollection.data.Name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)

		if err := mi.defaultExport.Set("variables", mi.getVariablesObject()); err != nil ***REMOVED***
			fmt.Println("Error setting collection variables object: ", err)
			common.Throw(rt, err)
		***REMOVED***
	***REMOVED***

	return mi
***REMOVED***

// Exports returns the exports of the collection module.
func (mi *CollectionInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Default: mi.defaultExport,
	***REMOVED***
***REMOVED***
