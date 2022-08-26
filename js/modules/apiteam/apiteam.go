// Package apiteam implements the module imported as 'apiteam' from inside k6.
package apiteam

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
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
			"context": mi.Context,
			"tag":     mi.Tag,
		},
	}
}

// Info returns current info about the APITeam Execution Context.
func (mi *APITeam) Context() *lib.WorkerInfo {
	workerInfo := mi.vu.InitEnv().WorkerInfo

	return workerInfo
}

type tagMessage struct {
	Tag     string      `json:"tag"`
	Message interface{} `json:"message"`
}

// Returns a tagged value to the orchestrator
func (mi *APITeam) Tag(tag string, value interface{}) error {
	workerInfo := mi.vu.InitEnv().WorkerInfo

	// Ensure no ':' in the tag
	if strings.Contains(tag, ":") {
		return fmt.Errorf("filename cannot contain ':'")
	}

	tagMessage := tagMessage{
		Tag:     tag,
		Message: value,
	}

	marshalled, err := json.Marshal(tagMessage)
	if err != nil {
		return err
	}

	lib.DispatchMessage(workerInfo.Ctx, workerInfo.Client, workerInfo.JobId, workerInfo.WorkerId, string(marshalled), "TAG")

	return nil
}
