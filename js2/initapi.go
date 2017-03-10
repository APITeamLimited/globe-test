/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package js2

import (
	"github.com/dop251/goja"
	"github.com/spf13/afero"
)

// Provides APIs for use in the init context.
type InitContext struct ***REMOVED***
	// Filesystem to load files and scripts from.
	Fs  afero.Fs
	Pwd string

	// Cache of loaded modules.
	Modules map[string]*goja.Program
***REMOVED***

func (i *InitContext) Require(mod string) goja.Value ***REMOVED***
	return goja.Undefined()
***REMOVED***
