package environment

import (
	"sync"

	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/APITeamLimited/globe-test/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create environment module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	EnvironmentModule struct {
		isEnabled         bool
		sharedEnvironment sharedEnvironment
	}

	// EnvironmentInstance represents an instance of the environment module.
	EnvironmentInstance struct {
		vu            modules.VU
		module        *EnvironmentModule
		defaultExport *goja.Object
	}

	sharedEnvironment struct {
		data *libWorker.Environment

		mu *sync.RWMutex
	}
)

var (
	_ modules.Module   = &EnvironmentModule{}
	_ modules.Instance = &EnvironmentInstance{}
)

// New returns a pointer to a new EnvironmentModule instance.
func New(workerInfo *libWorker.WorkerInfo) *EnvironmentModule {
	// Check environment actually exists
	if workerInfo.Environment != nil {
		return &EnvironmentModule{
			isEnabled: true,
			sharedEnvironment: sharedEnvironment{
				data: workerInfo.Environment,
				mu:   &sync.RWMutex{},
			},
		}
	}

	return &EnvironmentModule{
		isEnabled:         false,
		sharedEnvironment: sharedEnvironment{},
	}
}

// NewModuleInstance returns an environment module instance for each VU.
func (module *EnvironmentModule) NewModuleInstance(vu modules.VU) modules.Instance {
	rt := vu.Runtime()

	mi := &EnvironmentInstance{
		vu:            vu,
		module:        module,
		defaultExport: rt.NewObject(),
	}

	mi.defaultExport.DefineDataProperty(
		"enabled", rt.ToValue(module.isEnabled), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
	)

	if mi.module.isEnabled {
		mi.defaultExport.DefineDataProperty(
			"name", rt.ToValue(module.sharedEnvironment.data.Name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)

		if err := mi.defaultExport.Set("variables", mi.getVariablesObject()); err != nil {
			common.Throw(rt, err)
		}
	}

	return mi
}

// Exports returns the JS values this module exports.
func (mi *EnvironmentInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.defaultExport,
	}
}
