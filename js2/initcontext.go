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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js2/common"
	"github.com/loadimpact/k6/js2/modules"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Provides APIs for use in the init context.
type InitContext struct ***REMOVED***
	// Bound runtime; used to instantiate objects.
	runtime *goja.Runtime

	// Index of all loaded modules.
	Modules map[string]*common.Module `js:"-"`

	// Filesystem to load files and scripts from.
	fs  afero.Fs
	pwd string

	// Console object.
	Console *Console
***REMOVED***

func NewInitContext(rt *goja.Runtime, fs afero.Fs, pwd string) *InitContext ***REMOVED***
	return &InitContext***REMOVED***
		runtime: rt,
		fs:      fs,
		pwd:     pwd,

		Modules: make(map[string]*common.Module),

		Console: NewConsole(),
	***REMOVED***
***REMOVED***

func (i *InitContext) Require(arg string) goja.Value ***REMOVED***
	switch ***REMOVED***
	case arg == "k6", strings.HasPrefix(arg, "k6/"):
		// Builtin modules ("k6" or "k6/...") are handled specially, as they don't exist on the
		// filesystem. This intentionally shadows attempts to name your own modules this.
		v, err := i.requireModule(arg)
		if err != nil ***REMOVED***
			panic(i.runtime.NewGoError(err))
		***REMOVED***
		return v
	***REMOVED***
	return goja.Undefined()
***REMOVED***

func (i *InitContext) requireModule(name string) (goja.Value, error) ***REMOVED***
	log.WithField("name", name).Info("require module")
	mod, ok := i.Modules[name]
	if !ok ***REMOVED***
		mod_, ok := modules.Index[name]
		if !ok ***REMOVED***
			panic(i.runtime.NewGoError(errors.Errorf("unknown builtin module: %s", name)))
		***REMOVED***
		mod = &mod_
		i.Modules[name] = mod
	***REMOVED***
	return mod.Export(i.runtime), nil
***REMOVED***

// func (i *InitContext) requireProgram(pgm *goja.Program) (goja.Value, error) ***REMOVED***
// 	// Switch out the 'exports' global for a module-specific one.
// 	oldExports := i.runtime.Get("exports")
// 	i.runtime.Set("exports", i.runtime.NewObject())
// 	defer i.runtime.Set("exports", oldExports)

// 	// Run the program, this will populate the swapped-in exports.
// 	if _, err := i.runtime.RunProgram(pgm); err != nil ***REMOVED***
// 		log.WithError(err).Error("couldn't run module program")
// 		return goja.Undefined(), err
// 	***REMOVED***

// 	// Return the current exports, before the defer'd Set swaps it back.
// 	return i.runtime.Get("exports"), nil
// ***REMOVED***
