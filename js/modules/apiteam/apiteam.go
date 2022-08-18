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
	RootModule struct***REMOVED******REMOVED***

	// K6 represents an instance of the k6 module.
	APITeam struct ***REMOVED***
		vu modules.VU
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &APITeam***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &APITeam***REMOVED***vu: vu***REMOVED***
***REMOVED***

// Exports returns the exports of the k6 module.
func (mi *APITeam) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"info": mi.Info,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Info returns current info about the APITeam Execution Context.
func (mi *APITeam) Info(secs float64) *goja.Object ***REMOVED***
	workerInfo := mi.vu.InitEnv().WorkerInfo
	workerId := workerInfo.WorkerId

	newObject := mi.vu.Runtime().CreateObject(&goja.Object***REMOVED******REMOVED***)
	newObject.Set("workerId", workerId)

	return newObject
***REMOVED***
