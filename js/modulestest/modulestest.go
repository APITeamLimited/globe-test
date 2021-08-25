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

var _ modules.InstanceCore = &InstanceCore***REMOVED******REMOVED***

// InstanceCore is a modules.InstanceCore implementation meant to be used within tests
type InstanceCore struct ***REMOVED***
	Ctx     context.Context
	InitEnv *common.InitEnvironment
	State   *lib.State
	Runtime *goja.Runtime
***REMOVED***

// GetContext returns internally set field to conform to modules.InstanceCore interface
func (m *InstanceCore) GetContext() context.Context ***REMOVED***
	return m.Ctx
***REMOVED***

// GetInitEnv returns internally set field to conform to modules.InstanceCore interface
func (m *InstanceCore) GetInitEnv() *common.InitEnvironment ***REMOVED***
	return m.InitEnv
***REMOVED***

// GetState returns internally set field to conform to modules.InstanceCore interface
func (m *InstanceCore) GetState() *lib.State ***REMOVED***
	return m.State
***REMOVED***

// GetRuntime returns internally set field to conform to modules.InstanceCore interface
func (m *InstanceCore) GetRuntime() *goja.Runtime ***REMOVED***
	return m.Runtime
***REMOVED***
