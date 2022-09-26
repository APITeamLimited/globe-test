package collection

import (
	"sync"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create collection module instances for each VU.
//
// TODO: add sync.Once for all of the deprecation warnings we might want to do
// for the old k6/http APIs here, so they are shown only once in a test run.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct ***REMOVED***
		col sharedCollection
	***REMOVED***

	// Collection represents an instance of the collection module.
	Collection struct ***REMOVED***
		vu         modules.VU
		collection *sharedCollection
	***REMOVED***

	sharedCollection struct ***REMOVED***
		isEnabled bool
		variables map[string]lib.KeyValueItem
		mu        sync.RWMutex
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &Collection***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New(workerInfo *lib.WorkerInfo) *RootModule ***REMOVED***
	// Check collection actually exists
	if workerInfo.Collection != nil ***REMOVED***
		return &RootModule***REMOVED***
			col: sharedCollection***REMOVED***
				isEnabled: true,
				variables: *workerInfo.Collection.Variables,
			***REMOVED***,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return &RootModule***REMOVED***
			col: sharedCollection***REMOVED***
				isEnabled: false,
			***REMOVED***,
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an collection module instance for each VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &Collection***REMOVED***
		vu:         vu,
		collection: &rm.col,
	***REMOVED***
***REMOVED***

// Exports returns the exports of the collection module.
func (mi *Collection) Exports() modules.Exports ***REMOVED***
	if !mi.collection.isEnabled ***REMOVED***
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
func (mi *Collection) isEnabled() bool ***REMOVED***
	return mi.collection.isEnabled
***REMOVED***
