package lifecycle

import (
	"encoding/json"

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
	LifecycleModule struct ***REMOVED***
		isEnabled     bool
		sharedRequest sharedRequest
	***REMOVED***

	// RequestInstance represents an instance of the lifecycle module.
	LifecylcleInstance struct ***REMOVED***
		vu            modules.VU
		module        *LifecycleModule
		defaultExport *goja.Object
	***REMOVED***

	sharedRequest struct ***REMOVED***
		finalRequest      map[string]interface***REMOVED******REMOVED***
		underlyingRequest map[string]interface***REMOVED******REMOVED***
	***REMOVED***
)

var (
	_ modules.Module   = &LifecycleModule***REMOVED******REMOVED***
	_ modules.Instance = &LifecylcleInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new LifecycleModule instance.
func New(workerInfo *libWorker.WorkerInfo) *LifecycleModule ***REMOVED***
	// Only run in http_single execution mode

	// Check lifecycle actually exists
	if workerInfo.WorkerOptions.ExecutionMode.Value == types.HTTPSingleExecutionMode ***REMOVED***
		return &LifecycleModule***REMOVED***
			isEnabled: true,
			sharedRequest: sharedRequest***REMOVED***
				finalRequest:      workerInfo.FinalRequest,
				underlyingRequest: workerInfo.UnderlyingRequest,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	return &LifecycleModule***REMOVED***
		isEnabled:     false,
		sharedRequest: sharedRequest***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// NewModuleInstance returns an lifecycle module instance for each VU.
func (module *LifecycleModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	rt := vu.Runtime()

	mi := &LifecylcleInstance***REMOVED***
		vu:            vu,
		module:        module,
		defaultExport: rt.NewObject(),
	***REMOVED***

	mi.defaultExport.DefineDataProperty(
		"enabled", rt.ToValue(module.isEnabled), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
	)

	if mi.module.isEnabled ***REMOVED***
		mi.defaultExport.DefineDataProperty(
			"underlyingRequest", rt.ToValue(module.sharedRequest.underlyingRequest), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE,
		)

		mi.defaultExport.DefineDataProperty(
			"finalRequest", rt.ToValue(module.sharedRequest.finalRequest), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	***REMOVED***

	if err := mi.defaultExport.Set("markResponse", mi.markResponse); err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	return mi
***REMOVED***

// Exports returns the JS values this module exports.
func (mi *LifecylcleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Default: mi.defaultExport,
	***REMOVED***
***REMOVED***

func (mi *LifecylcleInstance) markResponse(responseObject *goja.Object) ***REMOVED***
	// Get golang value from goja object
	workerInfo := *mi.vu.InitEnv().WorkerInfo
	rt := mi.vu.Runtime()

	exportedResponse := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err := rt.ExportTo(responseObject, &exportedResponse)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	markedResponse := libWorker.MarkMessage***REMOVED***
		Mark:    "MarkedResponse",
		Message: exportedResponse,
	***REMOVED***

	// Marshal response to JSON
	marshalledMarkedResponse, err := json.Marshal(markedResponse)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	libWorker.DispatchMessage(workerInfo.Ctx, workerInfo.Client, workerInfo.JobId, workerInfo.WorkerId, string(marshalledMarkedResponse), "MARK")
***REMOVED***
