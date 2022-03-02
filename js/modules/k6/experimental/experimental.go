// Package experimental includes experimental module features
package experimental

import (
	"errors"
	"time"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type (
	// RootModule is the root experimental module
	RootModule struct***REMOVED******REMOVED***
	// ModuleInstance represents an instance of the experimental module
	ModuleInstance struct ***REMOVED***
		vu modules.VU
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &ModuleInstance***REMOVED******REMOVED***
)

// NewModuleInstance implements modules.Module interface
func (*RootModule) NewModuleInstance(m modules.VU) modules.Instance ***REMOVED***
	return &ModuleInstance***REMOVED***vu: m***REMOVED***
***REMOVED***

// New returns a new RootModule.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// Exports returns the exports of the experimental module
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"setTimeout": mi.setTimeout,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (mi *ModuleInstance) setTimeout(f goja.Callable, t float64) ***REMOVED***
	if f == nil ***REMOVED***
		common.Throw(mi.vu.Runtime(), errors.New("setTimeout requires a function as first argument"))
	***REMOVED***
	// TODO maybe really return something to use with `clearTimeout
	// TODO support arguments ... maybe
	runOnLoop := mi.vu.RegisterCallback()
	go func() ***REMOVED***
		timer := time.NewTimer(time.Duration(t * float64(time.Millisecond)))
		select ***REMOVED***
		case <-timer.C:
			runOnLoop(func() error ***REMOVED***
				_, err := f(goja.Undefined())
				return err
			***REMOVED***)
		case <-mi.vu.Context().Done():
			// TODO log something?

			timer.Stop()
			runOnLoop(func() error ***REMOVED*** return nil ***REMOVED***)
		***REMOVED***
	***REMOVED***()
***REMOVED***
