package environment

import (
	"sync"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create environment module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	EnvironmentModule struct ***REMOVED***
		isEnabled         bool
		sharedEnvironment sharedEnvironment
	***REMOVED***

	// EnvironmentInstance represents an instance of the environment module.
	EnvironmentInstance struct ***REMOVED***
		vu            modules.VU
		module        *EnvironmentModule
		defaultExport *goja.Object
	***REMOVED***

	sharedEnvironment struct ***REMOVED***
		data *libWorker.Environment

		mu *sync.RWMutex
	***REMOVED***
)

var (
	_ modules.Module   = &EnvironmentModule***REMOVED******REMOVED***
	_ modules.Instance = &EnvironmentInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new EnvironmentModule instance.
func New(workerInfo *libWorker.WorkerInfo) *EnvironmentModule ***REMOVED***
	// Check environment actually exists
	if workerInfo.Environment != nil ***REMOVED***
		return &EnvironmentModule***REMOVED***
			isEnabled: true,
			sharedEnvironment: sharedEnvironment***REMOVED***
				data: workerInfo.Environment,
				mu:   &sync.RWMutex***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return &EnvironmentModule***REMOVED***
		isEnabled:         false,
		sharedEnvironment: sharedEnvironment***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an environment module instance for each VU.
func (module *EnvironmentModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	rt := vu.Runtime()

	mi := &EnvironmentInstance***REMOVED***
		vu:            vu,
		module:        module,
		defaultExport: rt.NewObject(),
	***REMOVED***

	mi.defaultExport.DefineDataProperty(
		"enabled", rt.ToValue(module.isEnabled), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
	)

	if mi.module.isEnabled ***REMOVED***
		mi.defaultExport.DefineDataProperty(
			"name", rt.ToValue(module.sharedEnvironment.data.Name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)

		if err := mi.defaultExport.Set("variables", mi.getVariablesObject()); err != nil ***REMOVED***
			common.Throw(rt, err)
		***REMOVED***
	***REMOVED***

	return mi
***REMOVED***

// Exports returns the JS values this module exports.
func (mi *EnvironmentInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Default: mi.defaultExport,
	***REMOVED***
***REMOVED***
