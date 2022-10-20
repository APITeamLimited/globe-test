package lifecycle

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/dop251/goja"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create lifecycle module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	LifecycleModule struct {
		isEnabled     bool
		sharedRequest sharedRequest
	}

	// RequestInstance represents an instance of the lifecycle module.
	LifecylcleInstance struct {
		vu            modules.VU
		module        *LifecycleModule
		defaultExport *goja.Object
	}

	sharedRequest struct {
		finalRequest      map[string]interface{}
		underlyingRequest map[string]interface{}
	}
)

var (
	_ modules.Module   = &LifecycleModule{}
	_ modules.Instance = &LifecylcleInstance{}
)

// New returns a pointer to a new LifecycleModule instance.
func New(workerInfo *libWorker.WorkerInfo) *LifecycleModule {
	// Only run in http_single execution mode

	// Check lifecycle actually exists
	if workerInfo.WorkerOptions.ExecutionMode.Value == types.HTTPSingleExecutionMode ||
		workerInfo.WorkerOptions.ExecutionMode.Value == types.HTTPMultipleExecutionMode {
		return &LifecycleModule{
			isEnabled: true,
			sharedRequest: sharedRequest{
				finalRequest:      workerInfo.FinalRequest,
				underlyingRequest: workerInfo.UnderlyingRequest,
			},
		}
	}

	return &LifecycleModule{
		isEnabled: false,
	}
}

// NewModuleInstance returns an lifecycle module instance for each VU.
func (module *LifecycleModule) NewModuleInstance(vu modules.VU) modules.Instance {
	rt := vu.Runtime()

	mi := &LifecylcleInstance{
		vu:            vu,
		module:        module,
		defaultExport: rt.NewObject(),
	}

	mi.defaultExport.DefineDataProperty(
		"enabled", rt.ToValue(module.isEnabled), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
	)

	if mi.module.isEnabled {
		mi.defaultExport.DefineDataProperty(
			"underlyingRequest", rt.ToValue(module.sharedRequest.underlyingRequest), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)

		mi.defaultExport.DefineDataProperty(
			"finalRequest", rt.ToValue(module.sharedRequest.finalRequest), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	}

	if err := mi.defaultExport.Set("markResponse", mi.markResponse); err != nil {
		common.Throw(rt, err)
	}

	return mi
}

// Exports returns the JS values this module exports.
func (mi *LifecylcleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.defaultExport,
	}
}

func (mi *LifecylcleInstance) markResponse(responseObject goja.Value) {
	// Get golang value from goja object
	workerInfo := *mi.vu.InitEnv().WorkerInfo
	rt := mi.vu.Runtime()

	exportedResponse := map[string]interface{}{}
	err := rt.ExportTo(responseObject, &exportedResponse)
	if err != nil {
		common.Throw(rt, err)
	}

	markedResponse := libWorker.MarkMessage{
		Mark:    "MarkedResponse",
		Message: exportedResponse,
	}

	// Loop over marked response message and delete any function calls
	for key, value := range markedResponse.Message {
		if reflect.TypeOf(value).Kind() == reflect.Func {
			delete(markedResponse.Message, key)
		}
	}

	// Marshal response to JSON
	marshalledMarkedResponse, err := json.Marshal(markedResponse)
	if err != nil {
		fmt.Println("Error marshalling marked response", err)
		common.Throw(rt, err)
	}

	libWorker.DispatchMessage(*workerInfo.Gs, string(marshalledMarkedResponse), "MARK")
}
