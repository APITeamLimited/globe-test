// Package apiteam implements the module imported as 'apiteam' from inside k6.
package stores

import (
	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() ***REMOVED***
	modules.Register("apiteam/x/stores", new(Stores))
***REMOVED***

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct***REMOVED******REMOVED***

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct ***REMOVED***
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU
		// comparator is the exported type
		store *Stores
	***REMOVED***
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance***REMOVED******REMOVED***
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// Compare is the type for our custom API.
type Stores struct ***REMOVED***
	vu modules.VU // provides methods for accessing internal k6 objects

***REMOVED***

// IsGreater returns true if a is greater than b, or false otherwise, setting textual result message.
func (s *Stores) StoreResult(key, value string) bool ***REMOVED***
	// Store the result in the global state.

	//stringedMessage =
	//
	//s.vu.State().Logger().Info(
***REMOVED***

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		//Default: s.stores
	***REMOVED***
***REMOVED***
