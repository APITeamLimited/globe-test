// Package apiteam implements the module imported as 'apiteam' from inside k6.
package apiteam

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
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

	// APITeam represents an instance of the k6 module.
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

// Exports returns the exports of the apiteam module.
func (mi *APITeam) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"context": mi.Context,
			"mark":    mi.Mark,
		},
	}
}

// Info returns current info about the APITeam Execution Context.
func (mi *APITeam) Context() *libWorker.WorkerInfo {
	workerInfo := mi.vu.InitEnv().WorkerInfo

	return workerInfo
}

type markMessage struct {
	Mark    string      `json:"mark"`
	Message interface{} `json:"message"`
}

// Returns a marked value to the orchestrator
func (mi *APITeam) Mark(mark string, value interface{}) error {
	workerInfo := mi.vu.InitEnv().WorkerInfo

	// Ensure no ':' in the tag
	if strings.Contains(mark, ":") {
		return fmt.Errorf("filename cannot contain ':'")
	}

	markMessage := markMessage{
		Mark:    mark,
		Message: value,
	}

	marshalled, err := json.Marshal(markMessage)
	if err != nil {
		return err
	}

	libWorker.DispatchMessage(workerInfo.Ctx, workerInfo.Client, workerInfo.JobId, workerInfo.WorkerId, string(marshalled), "MARK")

	return nil
}
