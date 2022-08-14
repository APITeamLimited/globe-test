// Package apiteam implements the module imported as 'apiteam' from inside k6.
package stores

import (
	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("apiteam/x/stores", new(Stores))
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU
		// comparator is the exported type
		store *Stores
	}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// Compare is the type for our custom API.
type Stores struct {
	vu modules.VU // provides methods for accessing internal k6 objects

}

// IsGreater returns true if a is greater than b, or false otherwise, setting textual result message.
func (s *Stores) StoreResult(key, value string) bool {
	// Store the result in the global state.

	//stringedMessage =
	//
	//s.vu.State().Logger().Info(
}

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		//Default: s.stores
	}
}
