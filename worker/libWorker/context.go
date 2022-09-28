package libWorker

import (
	"context"
)

type ctxKey int

const (
	ctxKeyExecState ctxKey = iota
	ctxKeyScenario
)

// WithExecutionState embeds an ExecutionState in ctx.
func WithExecutionState(ctx context.Context, s *ExecutionState) context.Context ***REMOVED***
	return context.WithValue(ctx, ctxKeyExecState, s)
***REMOVED***

// GetExecutionState returns an ExecutionState from ctx.
func GetExecutionState(ctx context.Context) *ExecutionState ***REMOVED***
	v := ctx.Value(ctxKeyExecState)
	if v == nil ***REMOVED***
		return nil
	***REMOVED***
	return v.(*ExecutionState)
***REMOVED***

// WithScenarioState embeds a ScenarioState in ctx.
func WithScenarioState(ctx context.Context, s *ScenarioState) context.Context ***REMOVED***
	return context.WithValue(ctx, ctxKeyScenario, s)
***REMOVED***

// GetScenarioState returns a ScenarioState from ctx.
func GetScenarioState(ctx context.Context) *ScenarioState ***REMOVED***
	v := ctx.Value(ctxKeyScenario)
	if v == nil ***REMOVED***
		return nil
	***REMOVED***
	return v.(*ScenarioState)
***REMOVED***
