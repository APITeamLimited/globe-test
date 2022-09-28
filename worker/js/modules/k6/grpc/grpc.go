package grpc

import (
	"github.com/APITeamLimited/k6-worker/js/modules"
	"github.com/dop251/goja"
	"google.golang.org/grpc/codes"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct***REMOVED******REMOVED***

	// ModuleInstance represents an instance of the GRPC module for every VU.
	ModuleInstance struct ***REMOVED***
		vu      modules.VU
		exports map[string]interface***REMOVED******REMOVED***
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &ModuleInstance***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	mi := &ModuleInstance***REMOVED***
		vu:      vu,
		exports: make(map[string]interface***REMOVED******REMOVED***),
	***REMOVED***

	mi.exports["Client"] = mi.NewClient
	mi.defineConstants()
	return mi
***REMOVED***

// NewClient is the JS constructor for the grpc Client.
func (mi *ModuleInstance) NewClient(call goja.ConstructorCall) *goja.Object ***REMOVED***
	rt := mi.vu.Runtime()
	return rt.ToValue(&Client***REMOVED***vu: mi.vu***REMOVED***).ToObject(rt)
***REMOVED***

// defineConstants defines the constant variables of the module.
func (mi *ModuleInstance) defineConstants() ***REMOVED***
	rt := mi.vu.Runtime()
	mustAddCode := func(name string, code codes.Code) ***REMOVED***
		mi.exports[name] = rt.ToValue(code)
	***REMOVED***

	mustAddCode("StatusOK", codes.OK)
	mustAddCode("StatusCanceled", codes.Canceled)
	mustAddCode("StatusUnknown", codes.Unknown)
	mustAddCode("StatusInvalidArgument", codes.InvalidArgument)
	mustAddCode("StatusDeadlineExceeded", codes.DeadlineExceeded)
	mustAddCode("StatusNotFound", codes.NotFound)
	mustAddCode("StatusAlreadyExists", codes.AlreadyExists)
	mustAddCode("StatusPermissionDenied", codes.PermissionDenied)
	mustAddCode("StatusResourceExhausted", codes.ResourceExhausted)
	mustAddCode("StatusFailedPrecondition", codes.FailedPrecondition)
	mustAddCode("StatusAborted", codes.Aborted)
	mustAddCode("StatusOutOfRange", codes.OutOfRange)
	mustAddCode("StatusUnimplemented", codes.Unimplemented)
	mustAddCode("StatusInternal", codes.Internal)
	mustAddCode("StatusUnavailable", codes.Unavailable)
	mustAddCode("StatusDataLoss", codes.DataLoss)
	mustAddCode("StatusUnauthenticated", codes.Unauthenticated)
***REMOVED***

// Exports returns the exports of the grpc module.
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: mi.exports,
	***REMOVED***
***REMOVED***
