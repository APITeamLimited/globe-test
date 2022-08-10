package modulestest

import (
	"context"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
)

var _ modules.VU = &VU***REMOVED******REMOVED***

// VU is a modules.VU implementation meant to be used within tests
type VU struct ***REMOVED***
	CtxField              context.Context
	InitEnvField          *common.InitEnvironment
	StateField            *lib.State
	RuntimeField          *goja.Runtime
	RegisterCallbackField func() func(f func() error)
***REMOVED***

// Context returns internally set field to conform to modules.VU interface
func (m *VU) Context() context.Context ***REMOVED***
	return m.CtxField
***REMOVED***

// InitEnv returns internally set field to conform to modules.VU interface
func (m *VU) InitEnv() *common.InitEnvironment ***REMOVED***
	m.checkIntegrity()
	return m.InitEnvField
***REMOVED***

// State returns internally set field to conform to modules.VU interface
func (m *VU) State() *lib.State ***REMOVED***
	m.checkIntegrity()
	return m.StateField
***REMOVED***

// Runtime returns internally set field to conform to modules.VU interface
func (m *VU) Runtime() *goja.Runtime ***REMOVED***
	return m.RuntimeField
***REMOVED***

// RegisterCallback is not really implemented
func (m *VU) RegisterCallback() func(f func() error) ***REMOVED***
	return m.RegisterCallbackField()
***REMOVED***

func (m *VU) checkIntegrity() ***REMOVED***
	if m.InitEnvField != nil && m.StateField != nil ***REMOVED***
		panic("there is a bug in the test: InitEnvField and StateField are not allowed at the same time")
	***REMOVED***
***REMOVED***
