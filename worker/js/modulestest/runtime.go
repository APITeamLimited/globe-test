// Package modulestest contains helpers to test js modules
package modulestest

import (
	"context"
	"testing"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/eventloop"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/dop251/goja"
)

// Runtime is a helper struct that contains what is needed to run a (simple) module test
type Runtime struct ***REMOVED***
	VU             *VU
	EventLoop      *eventloop.EventLoop
	CancelContext  func()
	BuiltinMetrics *workerMetrics.BuiltinMetrics
***REMOVED***

// NewRuntime will create a new test runtime and will cancel the context on test/benchmark end
func NewRuntime(t testing.TB) *Runtime ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	vu := &VU***REMOVED***
		CtxField:     ctx,
		RuntimeField: goja.New(),
	***REMOVED***
	vu.RuntimeField.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	vu.InitEnvField = &common.InitEnvironment***REMOVED***
		Logger:   testutils.NewLogger(t),
		Registry: workerMetrics.NewRegistry(),
	***REMOVED***

	eventloop := eventloop.New(vu)
	vu.RegisterCallbackField = eventloop.RegisterCallback
	result := &Runtime***REMOVED***
		VU:             vu,
		EventLoop:      eventloop,
		CancelContext:  cancel,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(vu.InitEnvField.Registry),
	***REMOVED***
	// let's cancel again in case it has changed
	t.Cleanup(func() ***REMOVED*** result.CancelContext() ***REMOVED***)
	return result
***REMOVED***

// MoveToVUContext will set the state and nil the InitEnv just as a real VU
func (r *Runtime) MoveToVUContext(state *libWorker.State) ***REMOVED***
	r.VU.InitEnvField = nil
	r.VU.StateField = state
***REMOVED***
