// Javascript api for interracting with APITeam, for more information, see https://apiteam.cloud
package apiteam

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
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
			"mark": mi.Mark,
		},
	}
}

// Returns a marked value to the orchestrator
func (mi *APITeam) Mark(mark string, markedObject *goja.Object) error {
	workerInfo := mi.vu.InitEnv().WorkerInfo
	rt := mi.vu.Runtime()

	// Ensure no ':' in the tag
	if strings.Contains(mark, ":") {
		return common.NewInitContextError(fmt.Sprintf("Mark tag cannot contain ':' character: %s", mark))
	}

	exportedResponse := map[string]interface{}{}
	err := rt.ExportTo(markedObject, &exportedResponse)
	if err != nil {
		common.Throw(rt, err)
	}

	markedMessage := libWorker.MarkMessage{
		Mark:    mark,
		Message: exportedResponse,
	}

	// Loop over marked response message and delete any function calls
	for key, value := range markedMessage.Message {
		if reflect.TypeOf(value).Kind() == reflect.Func {
			delete(markedMessage.Message, key)
		}
	}

	marshalledMarkedMessage, err := json.Marshal(markedMessage)
	if err != nil {
		return err
	}

	libWorker.DispatchMessage(*workerInfo.Gs, string(marshalledMarkedMessage), "MARK")

	return nil
}
