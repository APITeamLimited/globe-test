/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

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
