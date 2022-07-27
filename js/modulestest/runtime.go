// Package modulestest contains helpers to test js modules
package modulestest

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/eventloop"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"
)

// Runtime is a helper struct that contains what is needed to run a (simple) module test
type Runtime struct ***REMOVED***
	VU             *VU
	EventLoop      *eventloop.EventLoop
	CancelContext  func()
	BuiltinMetrics *metrics.BuiltinMetrics
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
		Registry: metrics.NewRegistry(),
	***REMOVED***

	eventloop := eventloop.New(vu)
	vu.RegisterCallbackField = eventloop.RegisterCallback
	result := &Runtime***REMOVED***
		VU:             vu,
		EventLoop:      eventloop,
		CancelContext:  cancel,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(vu.InitEnvField.Registry),
	***REMOVED***
	// let's cancel again in case it has changed
	t.Cleanup(func() ***REMOVED*** result.CancelContext() ***REMOVED***)
	return result
***REMOVED***

// MoveToVUContext will set the state and nil the InitEnv just as a real VU
func (r *Runtime) MoveToVUContext(state *lib.State) ***REMOVED***
	r.VU.InitEnvField = nil
	r.VU.StateField = state
***REMOVED***