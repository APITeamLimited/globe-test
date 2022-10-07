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
	ContextModule struct {
		workerInfo *libWorker.WorkerInfo
	}

	// Context represents an instance of the context module.
	ContextInstance struct {
		vu modules.VU

		contextModule *ContextModule
	}
)

var (
	_ modules.Module   = &ContextModule{}
	_ modules.Instance = &ContextInstance{}
)

// New returns a pointer to a new ContextModule instance.
func New(workerInfo *libWorker.WorkerInfo) *ContextModule {
	// Check environment actually exists
	contextModule := &ContextModule{
		workerInfo: workerInfo,
	}

	return contextModule
}

// NewModuleInstance returns an environment module instance for each VU.
func (rm *ContextModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ContextInstance{
		vu:            vu,
		contextModule: rm,
	}
}

// Exports returns the module's exports.
func (ci *ContextInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"environment": environment.New(ci.contextModule.workerInfo),
			"collection":  collection.New(ci.contextModule.workerInfo),
		},
	}
}
