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
	RootModule struct***REMOVED******REMOVED***

	// APITeam represents an instance of the k6 module.
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

// Exports returns the exports of the apiteam module.
func (mi *APITeam) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"mark": mi.Mark,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Returns a marked value to the orchestrator
func (mi *APITeam) Mark(mark string, markedObject *goja.Object) error ***REMOVED***
	workerInfo := mi.vu.InitEnv().WorkerInfo
	rt := mi.vu.Runtime()

	// Ensure no ':' in the tag
	if strings.Contains(mark, ":") ***REMOVED***
		return common.NewInitContextError(fmt.Sprintf("Mark tag cannot contain ':' character: %s", mark))
	***REMOVED***

	exportedResponse := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err := rt.ExportTo(markedObject, &exportedResponse)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	markedMessage := libWorker.MarkMessage***REMOVED***
		Mark:    mark,
		Message: exportedResponse,
	***REMOVED***

	// Loop over marked response message and delete any function calls
	for key, value := range markedMessage.Message ***REMOVED***
		if reflect.TypeOf(value).Kind() == reflect.Func ***REMOVED***
			delete(markedMessage.Message, key)
		***REMOVED***
	***REMOVED***

	marshalledMarkedMessage, err := json.Marshal(markedMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	libWorker.DispatchMessage(*workerInfo.Gs, string(marshalledMarkedMessage), "MARK")

	return nil
***REMOVED***
