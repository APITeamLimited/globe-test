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
	"context"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/compiler"
	"github.com/loadimpact/k6/js2/common"
	"github.com/loadimpact/k6/js2/modules"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Provides APIs for use in the init context.
type InitContext struct ***REMOVED***
	// Bound runtime; used to instantiate objects.
	runtime *goja.Runtime

	// Pointer to a context that bridged modules are invoked with.
	ctxPtr *context.Context

	// Filesystem to load files and scripts from.
	fs  afero.Fs
	pwd string

	// Cache of loaded programs.
	programs map[string]*goja.Program `js:"-"`

	// Console object.
	Console *Console
***REMOVED***

func NewInitContext(rt *goja.Runtime, ctxPtr *context.Context, fs afero.Fs, pwd string) *InitContext ***REMOVED***
	return &InitContext***REMOVED***
		runtime: rt,
		ctxPtr:  ctxPtr,
		fs:      fs,
		pwd:     pwd,

		programs: make(map[string]*goja.Program),

		Console: NewConsole(),
	***REMOVED***
***REMOVED***

func newBoundInitContext(base *InitContext, ctxPtr *context.Context, rt *goja.Runtime) *InitContext ***REMOVED***
	return NewInitContext(rt, ctxPtr, base.fs, base.pwd)
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
	default:
		// Fall back to loading from the filesystem.
		v, err := i.requireFile(arg)
		if err != nil ***REMOVED***
			panic(i.runtime.NewGoError(err))
		***REMOVED***
		return v
	***REMOVED***
***REMOVED***

func (i *InitContext) requireModule(name string) (goja.Value, error) ***REMOVED***
	mod, ok := modules.Index[name]
	if !ok ***REMOVED***
		panic(i.runtime.NewGoError(errors.Errorf("unknown builtin module: %s", name)))
	***REMOVED***
	return i.runtime.ToValue(common.Bind(i.runtime, mod, i.ctxPtr)), nil
***REMOVED***

func (i *InitContext) requireFile(name string) (goja.Value, error) ***REMOVED***
	// Resolve the file path, push the target directory as pwd to make relative imports work.
	pwd := i.pwd
	filename := filepath.Join(pwd, name)
	i.pwd = filepath.Dir(filename)
	defer func() ***REMOVED*** i.pwd = pwd ***REMOVED***()

	// Swap the importing scope's imports out, then put it back again.
	oldExports := i.runtime.Get("exports")
	i.runtime.Set("exports", i.runtime.NewObject())
	defer i.runtime.Set("exports", oldExports)

	// Read sources, transform into ES6 and cache the compiled program.
	pgm, ok := i.programs[filename]
	if !ok ***REMOVED***
		data, err := afero.ReadFile(i.fs, filename)
		if err != nil ***REMOVED***
			return goja.Undefined(), err
		***REMOVED***
		src, _, err := compiler.Transform(string(data), filename)
		if err != nil ***REMOVED***
			return goja.Undefined(), err
		***REMOVED***
		pgm_, err := goja.Compile(filename, src, true)
		if err != nil ***REMOVED***
			return goja.Undefined(), err
		***REMOVED***
		i.programs[filename] = pgm_
		pgm = pgm_
	***REMOVED***

	// Execute the program to populate exports. You may notice that this theoretically allows an
	// imported file to access or overwrite globals defined outside of it. Please don't do anything
	// stupid with this, consider *any* use of it undefined behavior >_>;;
	if _, err := i.runtime.RunProgram(pgm); err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***

	exports := i.runtime.Get("exports")
	return exports, nil
***REMOVED***
