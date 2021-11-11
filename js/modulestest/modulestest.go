/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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
	CtxField     context.Context
	InitEnvField *common.InitEnvironment
	StateField   *lib.State
	RuntimeField *goja.Runtime
***REMOVED***

// Context returns internally set field to conform to modules.VU interface
func (m *VU) Context() context.Context ***REMOVED***
	return m.CtxField
***REMOVED***

// InitEnv returns internally set field to conform to modules.VU interface
func (m *VU) InitEnv() *common.InitEnvironment ***REMOVED***
	return m.InitEnvField
***REMOVED***

// State returns internally set field to conform to modules.VU interface
func (m *VU) State() *lib.State ***REMOVED***
	return m.StateField
***REMOVED***

// Runtime returns internally set field to conform to modules.VU interface
func (m *VU) Runtime() *goja.Runtime ***REMOVED***
	return m.RuntimeField
***REMOVED***
