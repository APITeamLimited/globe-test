package context

import (
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/js/modules/apiteam/collection"
	"github.com/APITeamLimited/globe-test/worker/js/modules/apiteam/environment"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	ContextModule struct ***REMOVED***
		workerInfo *libWorker.WorkerInfo
	***REMOVED***

	// Context represents an instance of the context module.
	ContextInstance struct ***REMOVED***
		vu modules.VU

		contextModule *ContextModule
	***REMOVED***
)

var (
	_ modules.Module   = &ContextModule***REMOVED******REMOVED***
	_ modules.Instance = &ContextInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new ContextModule instance.
func New(workerInfo *libWorker.WorkerInfo) *ContextModule ***REMOVED***
	// Check environment actually exists
	contextModule := &ContextModule***REMOVED***
		workerInfo: workerInfo,
	***REMOVED***

	return contextModule
***REMOVED***

// NewModuleInstance returns an environment module instance for each VU.
func (rm *ContextModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &ContextInstance***REMOVED***
		vu:            vu,
		contextModule: rm,
	***REMOVED***
***REMOVED***

// Exports returns the module's exports.
func (ci *ContextInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"environment": environment.New(ci.contextModule.workerInfo),
			"collection":  collection.New(ci.contextModule.workerInfo),
		***REMOVED***,
	***REMOVED***
***REMOVED***
