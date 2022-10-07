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
	CollectionModule struct ***REMOVED***
		sharedCollection sharedCollection
	***REMOVED***

	// CollectionInstance represents an instance of the collection module.
	CollectionInstance struct ***REMOVED***
		vu     modules.VU
		module *CollectionModule
	***REMOVED***

	sharedCollection struct ***REMOVED***
		isEnabled bool
		data      libWorker.Collection

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
			sharedCollection: sharedCollection***REMOVED***
				isEnabled: true,
				data:      *workerInfo.Collection,
			***REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return &CollectionModule***REMOVED***
			sharedCollection: sharedCollection***REMOVED***
				isEnabled: false,
			***REMOVED***,
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an collection module instance for each VU.
func (module *CollectionModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &CollectionInstance***REMOVED***
		vu,
		module,
	***REMOVED***
***REMOVED***

// Exports returns the exports of the collection module.
func (mi *CollectionInstance) Exports() modules.Exports ***REMOVED***
	if !mi.module.sharedCollection.isEnabled ***REMOVED***
		return modules.Exports***REMOVED***
			Named: map[string]interface***REMOVED******REMOVED******REMOVED***
				"isEnabled": mi.isEnabled,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"isEnabled": mi.isEnabled,
			"variables": mi.getVariables,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// isEnabled returns whether the collection is enabled or not.
func (mi *CollectionInstance) isEnabled() bool ***REMOVED***
	return mi.module.sharedCollection.isEnabled
***REMOVED***
