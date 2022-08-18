// Package apiteam implements the module imported as 'apiteam' from inside k6.
package apiteam

import (
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

var (
	// ErrGroupInInitContext is returned when group() are using in the init context.
	ErrGroupInInitContext = common.NewInitContextError("Using group() in the init context is not supported")

	// ErrCheckInInitContext is returned when check() are using in the init context.
	ErrCheckInInitContext = common.NewInitContextError("Using check() in the init context is not supported")
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// K6 represents an instance of the k6 module.
	APITeam struct {
		vu modules.VU
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &APITeam{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &APITeam{vu: vu}
}

// Exports returns the exports of the k6 module.
func (mi *APITeam) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"info": mi.Info,
		},
	}
}

// Info returns current info about the APITeam Execution Context.
func (mi *APITeam) Info(secs float64) *goja.Object {
	workerInfo := mi.vu.InitEnv().WorkerInfo
	workerId := workerInfo.WorkerId

	newObject := mi.vu.Runtime().CreateObject(&goja.Object{})
	newObject.Set("workerId", workerId)

	return newObject
}
