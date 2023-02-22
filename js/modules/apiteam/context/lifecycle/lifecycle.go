package lifecycle

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/APITeamLimited/globe-test/js/modules"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
)

const (
	stateUnknown  = "unknown"
	stateEnabled  = "enabled"
	stateDisabled = "disabled"
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create lifecycle module instances for each VU.
type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	LifecycleModule struct {
	}

	// RequestInstance represents an instance of the lifecycle module.
	LifecycleInstance struct {
		vu              modules.VU
		module          *LifecycleModule
		lifecycleObject *goja.Object
		// Provide cached value of enabled state so don't have to keep redefining
		// it in the runtime
		enableState string // enabled, disabled, unknown
	}
)

var (
	_ modules.Module   = &LifecycleModule{}
	_ modules.Instance = &LifecycleInstance{}
)

// New returns a pointer to a new LifecycleModule instance.
func New(workerInfo *libWorker.WorkerInfo) *LifecycleModule {
	return &LifecycleModule{}
}

// NewModuleInstance returns an lifecycle module instance for each VU.
func (module *LifecycleModule) NewModuleInstance(vu modules.VU) modules.Instance {
	rt := vu.Runtime()

	mi := &LifecycleInstance{
		vu:              vu,
		module:          module,
		lifecycleObject: rt.NewObject(),
		enableState:     stateUnknown,
	}

	if err := mi.lifecycleObject.Set("markResponse", mi.markResponse); err != nil {
		common.Throw(rt, err)
	}

	currentNodeCallable := func(call goja.FunctionCall) goja.Value {
		state := vu.State()
		return mi.getNodeObject(state.CurrentNode)
	}

	if err := mi.lifecycleObject.Set("node", currentNodeCallable); err != nil {
		common.Throw(rt, err)
	}

	return mi
}

func (mi *LifecycleInstance) getNodeObject(node libWorker.Node) *goja.Object {
	rt := mi.vu.Runtime()
	//state := mi.vu.State()

	nodeObject := rt.NewObject()

	nodeObject.DefineDataProperty("variant", rt.ToValue(node.GetVariant()), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	nodeObject.DefineDataProperty("name", rt.ToValue(node.GetName()), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	nodeObject.DefineDataProperty("id", rt.ToValue(node.GetId()), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)

	// TODO: Figure out why we can't currently use multiple scripts
	/*scriptsObject := rt.NewObject()
	for scriptName, exports := range node.GetScripts() {
		scriptObject := rt.NewObject()

		for exportKey, callable := range exports {
			// Ignore options and setup/teardown functions
			if exportKey == consts.Options || exportKey == consts.SetupFn || exportKey == consts.TeardownFn {
				continue
			}

			exportObject := rt.NewObject()

			exportObject.DefineDataProperty("name", rt.ToValue(exportKey), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
			//exportObject.Set("call", callable)
			exportObject.Set("call", func(call goja.FunctionCall) goja.Value {
				mi.nodeMutex.Lock()

				parentNode := state.CurrentNode

				// 	//fmt.Println("Parent node is", parentNode.GetName(), "current node is", node.GetName())

				// Update the context with the current node
				if node.GetId() != state.CurrentNode.GetId() {
					fmt.Println("Setting current node to", node.GetName())
					state.CurrentNode = node
				}

				// 	fmt.Println("Calling exported function", exportKey, callable, goja.IsUndefined(call.This), goja.IsNull(call.This), call.This, call.Arguments)

				mi.nodeMutex.Unlock()
				// TODO: Figure out why the runtime is panicking when we call the
				// function directly, not proud of this forcer function
				value, err := callable(goja.Undefined(), call.Arguments...)
				if err != nil {
					// Unsure how we can best handle this
					// TODO: Handle this better
					fmt.Println("Error calling exported function", err)
					panic(err)
				}

				// Reset the current node
				if node.GetId() != state.CurrentNode.GetId() {
					state.CurrentNode = parentNode
				}

				return value
			})

			scriptObject.Set(exportKey, exportObject)
		}

		scriptsObject.Set(scriptName, rt.ToValue(scriptObject))
	}

	nodeObject.Set("scripts", scriptsObject)
	*/

	nodeVariant := node.GetVariant()

	if nodeVariant == libWorker.HTTPRequestVariant {
		httpRequestNode := node.(*libWorker.HTTPRequestNode)

		nodeObject.DefineDataProperty(
			"finalRequest", rt.ToValue(httpRequestNode.FinalRequest), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)

		// Batch not yet working

		// nodeObject.Set(
		// 	"getBatchedRequests", func(call goja.FunctionCall) goja.Value {
		// 		return rt.ToValue(httpRequestNode.GetBatchedRequests(rt))
		// 	})

	} else if nodeVariant == libWorker.GroupVariant {
		groupNode := node.(*libWorker.GroupNode)

		childObjects := make([]*goja.Object, len(groupNode.Children))
		for index, child := range groupNode.Children {
			// Node on state will be different to the child node, enabling us to
			// automatically set the current node when an exported function is
			// called
			childObjects[index] = mi.getNodeObject(child)
		}

		nodeObject.DefineDataProperty("children", rt.ToValue(childObjects), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)

		// Batch not yet working

		// nodeObject.Set(
		// 	"getBatchedRequests", func(call goja.FunctionCall) goja.Value {

		// 		fmt.Println("Getting batched requests", groupNode.GetBatchedRequests(rt))

		// 		return rt.ToValue(groupNode.GetBatchedRequests(rt))
		// 	})
	}

	return nodeObject
}

// Exports returns the JS values this module exports.
func (mi *LifecycleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.lifecycleObject,
		Named: map[string]interface{}{
			"markResponse": mi.markResponse,
			"node":         mi.lifecycleObject.Get("node"),
		},
	}
}

func (mi *LifecycleInstance) markResponse(responseObject goja.Value) {
	// Get golang value from goja object
	workerInfo := *mi.vu.InitEnv().WorkerInfo
	rt := mi.vu.Runtime()

	exportedResponse := map[string]interface{}{}
	err := rt.ExportTo(responseObject, &exportedResponse)
	if err != nil {
		common.Throw(rt, err)
	}

	if exportedResponse["error"].(string) != "" {
		common.Throw(rt, errors.New(exportedResponse["error"].(string)))
		return
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
		common.Throw(rt, err)
	}

	libWorker.DispatchMessage(*workerInfo.Gs, string(marshalledMarkedResponse), "MARK")
}
