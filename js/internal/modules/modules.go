/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

package modules

import (
	"fmt"
	"sync"
)

//nolint:gochecknoglobals
var (
	modules = make(map[string]interface***REMOVED******REMOVED***)
	mx      sync.RWMutex
)

// Get returns the module registered with name.
func Get(name string) interface***REMOVED******REMOVED*** ***REMOVED***
	mx.RLock()
	defer mx.RUnlock()
	mod := modules[name]
	if i, ok := mod.(HasModuleInstancePerVU); ok ***REMOVED***
		return i.NewModuleInstancePerVU()
	***REMOVED***
	return mod
***REMOVED***

// HasModuleInstancePerVU should be implemented by all native Golang modules that
// would require per-VU state. k6 will call their NewModuleInstancePerVU() methods
// every time a VU imports the module and use its result as the returned object.
type HasModuleInstancePerVU interface ***REMOVED***
	NewModuleInstancePerVU() interface***REMOVED******REMOVED***
***REMOVED***

// Register the given mod as a JavaScript module, available
// for import from JS scripts by name.
// This function panics if a module with the same name is already registered.
func Register(name string, mod interface***REMOVED******REMOVED***) ***REMOVED***
	mx.Lock()
	defer mx.Unlock()

	if _, ok := modules[name]; ok ***REMOVED***
		panic(fmt.Sprintf("module already registered: %s", name))
	***REMOVED***
	modules[name] = mod
***REMOVED***
