package context

import (
	"github.com/APITeamLimited/globe-test/js/modules"
	"github.com/APITeamLimited/globe-test/js/modules/apiteam/context/collection"
	"github.com/APITeamLimited/globe-test/js/modules/apiteam/context/environment"
	"github.com/APITeamLimited/globe-test/js/modules/apiteam/context/lifecycle"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	ContextModule struct {
		workerInfo *libWorker.WorkerInfo

		environmentModule *environment.EnvironmentModule
		collectionModule  *collection.CollectionModule
		lifecycleModlule  *lifecycle.LifecycleModule
	}

	// Context represents an instance of the context module.
	ContextInstance struct {
		vu modules.VU

		contextModule *ContextModule
		exports       modules.Exports
	}
)

var (
	_ modules.Module   = &ContextModule{}
	_ modules.Instance = &ContextInstance{}
)

// New returns a pointer to a new ContextModule instance.
func New(workerInfo *libWorker.WorkerInfo) *ContextModule {
	// Will determine if each submodule is enabled enabled in the modules
	contextModule := &ContextModule{
		workerInfo:        workerInfo,
		environmentModule: environment.New(workerInfo),
		collectionModule:  collection.New(workerInfo),
		lifecycleModlule:  lifecycle.New(workerInfo),
	}

	return contextModule
}

// NewModuleInstance returns an context module instance for each VU.
func (rm *ContextModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ContextInstance{
		vu:            vu,
		contextModule: rm,
		exports: modules.Exports{
			Named: map[string]interface{}{
				"environment": rm.environmentModule.NewModuleInstance(vu).Exports().Default,
				"collection":  rm.collectionModule.NewModuleInstance(vu).Exports().Default,
				"lifecycle":   rm.lifecycleModlule.NewModuleInstance(vu).Exports().Default,
			},
		},
	}
}

// Exports returns the module's exports.
func (ci *ContextInstance) Exports() modules.Exports {
	return ci.exports
}
