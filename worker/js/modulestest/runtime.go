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
type Runtime struct {
	VU             *VU
	EventLoop      *eventloop.EventLoop
	CancelContext  func()
	BuiltinMetrics *workerMetrics.BuiltinMetrics
}

// NewRuntime will create a new test runtime and will cancel the context on test/benchmark end
func NewRuntime(t testing.TB) *Runtime {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	vu := &VU{
		CtxField:     ctx,
		RuntimeField: goja.New(),
	}
	vu.RuntimeField.SetFieldNameMapper(common.FieldNameMapper{})
	vu.InitEnvField = &common.InitEnvironment{
		Logger:   testutils.NewLogger(t),
		Registry: workerMetrics.NewRegistry(),
	}

	eventloop := eventloop.New(vu)
	vu.RegisterCallbackField = eventloop.RegisterCallback
	result := &Runtime{
		VU:             vu,
		EventLoop:      eventloop,
		CancelContext:  cancel,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(vu.InitEnvField.Registry),
	}
	// let's cancel again in case it has changed
	t.Cleanup(func() { result.CancelContext() })
	return result
}

// MoveToVUContext will set the state and nil the InitEnv just as a real VU
func (r *Runtime) MoveToVUContext(state *libWorker.State) {
	r.VU.InitEnvField = nil
	r.VU.StateField = state
}
