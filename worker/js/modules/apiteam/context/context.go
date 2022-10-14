package context

import (
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/js/modules/apiteam/context/collection"
	"github.com/APITeamLimited/globe-test/worker/js/modules/apiteam/context/environment"
	"github.com/APITeamLimited/globe-test/worker/js/modules/apiteam/context/lifecycle"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	ContextModule struct ***REMOVED***
		workerInfo *libWorker.WorkerInfo

		environmentModule *environment.EnvironmentModule
		collectionModule  *collection.CollectionModule
		lifecycleModlule  *lifecycle.LifecycleModule
	***REMOVED***

	// Context represents an instance of the context module.
	ContextInstance struct ***REMOVED***
		vu modules.VU

		contextModule *ContextModule
		exports       modules.Exports
	***REMOVED***
)

var (
	_ modules.Module   = &ContextModule***REMOVED******REMOVED***
	_ modules.Instance = &ContextInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new ContextModule instance.
func New(workerInfo *libWorker.WorkerInfo) *ContextModule ***REMOVED***
	// Will determine if each submodule is enabled enabled in the modules
	contextModule := &ContextModule***REMOVED***
		workerInfo:        workerInfo,
		environmentModule: environment.New(workerInfo),
		collectionModule:  collection.New(workerInfo),
		lifecycleModlule:  lifecycle.New(workerInfo),
	***REMOVED***

	return contextModule
***REMOVED***

// NewModuleInstance returns an context module instance for each VU.
func (rm *ContextModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &ContextInstance***REMOVED***
		vu:            vu,
		contextModule: rm,
		exports: modules.Exports***REMOVED***
			Named: map[string]interface***REMOVED******REMOVED******REMOVED***
				"environment": rm.environmentModule.NewModuleInstance(vu).Exports().Default,
				"collection":  rm.collectionModule.NewModuleInstance(vu).Exports().Default,
				"lifecycle":   rm.lifecycleModlule.NewModuleInstance(vu).Exports().Default,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Exports returns the module's exports.
func (ci *ContextInstance) Exports() modules.Exports ***REMOVED***
	return ci.exports
***REMOVED***
