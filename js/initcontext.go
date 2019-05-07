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

package js

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/compiler"
	"github.com/loadimpact/k6/js/modules"
	"github.com/loadimpact/k6/loader"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type programWithSource struct ***REMOVED***
	pgm     *goja.Program
	src     string
	exports goja.Value
***REMOVED***

// InitContext provides APIs for use in the init context.
type InitContext struct ***REMOVED***
	// Bound runtime; used to instantiate objects.
	runtime  *goja.Runtime
	compiler *compiler.Compiler

	// Pointer to a context that bridged modules are invoked with.
	ctxPtr *context.Context

	// Filesystem to load files and scripts from.
	fs  afero.Fs
	pwd string

	// Cache of loaded programs and files.
	programs map[string]programWithSource
	files    map[string][]byte
***REMOVED***

// NewInitContext creates a new initcontext with the provided arguments
func NewInitContext(
	rt *goja.Runtime, compiler *compiler.Compiler, ctxPtr *context.Context, fs afero.Fs, pwd string,
) *InitContext ***REMOVED***
	return &InitContext***REMOVED***
		runtime:  rt,
		compiler: compiler,
		ctxPtr:   ctxPtr,
		fs:       fs,
		pwd:      filepath.ToSlash(pwd),

		programs: make(map[string]programWithSource),
		files:    make(map[string][]byte),
	***REMOVED***
***REMOVED***

func newBoundInitContext(base *InitContext, ctxPtr *context.Context, rt *goja.Runtime) *InitContext ***REMOVED***
	// we don't copy the exports as otherwise they will be shared and we don't want this.
	// this means that all the files will be executed again but once again only once per compilation
	// of the main file.
	var programs = make(map[string]programWithSource, len(base.programs))
	for key, program := range base.programs ***REMOVED***
		programs[key] = programWithSource***REMOVED***
			src: program.src,
			pgm: program.pgm,
		***REMOVED***
	***REMOVED***
	return &InitContext***REMOVED***
		runtime: rt,
		ctxPtr:  ctxPtr,

		fs:       base.fs,
		pwd:      base.pwd,
		compiler: base.compiler,

		programs: programs,
		files:    base.files,
	***REMOVED***
***REMOVED***

// Require is called when a module/file needs to be loaded by a script
func (i *InitContext) Require(arg string) goja.Value ***REMOVED***
	switch ***REMOVED***
	case arg == "k6", strings.HasPrefix(arg, "k6/"):
		// Builtin modules ("k6" or "k6/...") are handled specially, as they don't exist on the
		// filesystem. This intentionally shadows attempts to name your own modules this.
		v, err := i.requireModule(arg)
		if err != nil ***REMOVED***
			common.Throw(i.runtime, err)
		***REMOVED***
		return v
	default:
		// Fall back to loading from the filesystem.
		v, err := i.requireFile(arg)
		if err != nil ***REMOVED***
			common.Throw(i.runtime, err)
		***REMOVED***
		return v
	***REMOVED***
***REMOVED***

func (i *InitContext) requireModule(name string) (goja.Value, error) ***REMOVED***
	mod, ok := modules.Index[name]
	if !ok ***REMOVED***
		return nil, errors.Errorf("unknown builtin module: %s", name)
	***REMOVED***
	return i.runtime.ToValue(common.Bind(i.runtime, mod, i.ctxPtr)), nil
***REMOVED***

func (i *InitContext) requireFile(name string) (goja.Value, error) ***REMOVED***
	// Resolve the file path, push the target directory as pwd to make relative imports work.
	pwd := i.pwd
	filename := loader.Resolve(pwd, name)

	// First, check if we have a cached program already.
	pgm, ok := i.programs[filename]
	if !ok || pgm.exports == nil ***REMOVED***
		i.pwd = loader.Dir(filename)
		defer func() ***REMOVED*** i.pwd = pwd ***REMOVED***()

		// Swap the importing scope's exports out, then put it back again.
		oldExports := i.runtime.Get("exports")
		defer i.runtime.Set("exports", oldExports)
		oldModule := i.runtime.Get("module")
		defer i.runtime.Set("module", oldModule)
		exports := i.runtime.NewObject()
		i.runtime.Set("exports", exports)
		module := i.runtime.NewObject()
		_ = module.Set("exports", exports)
		i.runtime.Set("module", module)
		if pgm.pgm == nil ***REMOVED***
			// Load the sources; the loader takes care of remote loading, etc.
			data, err := loader.Load(i.fs, pwd, name)
			if err != nil ***REMOVED***
				return goja.Undefined(), err
			***REMOVED***
			pgm.src = string(data.Data)

			// Compile the sources; this handles ES5 vs ES6 automatically.
			pgm.pgm, err = i.compileImport(pgm.src, data.Filename)
			if err != nil ***REMOVED***
				return goja.Undefined(), err
			***REMOVED***
		***REMOVED***

		pgm.exports = module.Get("exports")
		i.programs[filename] = pgm

		// Run the program.
		if _, err := i.runtime.RunProgram(pgm.pgm); err != nil ***REMOVED***
			return goja.Undefined(), err
		***REMOVED***

		pgm.exports = module.Get("exports")
	***REMOVED***

	return pgm.exports, nil
***REMOVED***

func (i *InitContext) compileImport(src, filename string) (*goja.Program, error) ***REMOVED***
	pgm, _, err := i.compiler.Compile(src, filename, "(function()***REMOVED***\n", "\n***REMOVED***)()\n", true)
	return pgm, err
***REMOVED***

// Open implements open() in the init context and will read and return the contents of a file
func (i *InitContext) Open(filename string, args ...string) (goja.Value, error) ***REMOVED***
	if filename == "" ***REMOVED***
		return nil, errors.New("open() can't be used with an empty filename")
	***REMOVED***

	// Here IsAbs should be enough but unfortunately it doesn't handle absolute paths starting from
	// the current drive on windows like `\users\noname\...`. Also it makes it more easy to test and
	// will probably be need for archive execution under windows if always consider '/...' as an
	// absolute path.
	if filename[0] != '/' && filename[0] != '\\' && !filepath.IsAbs(filename) ***REMOVED***
		filename = filepath.Join(i.pwd, filename)
	***REMOVED***
	filename = filepath.ToSlash(filename)

	data, ok := i.files[filename]
	if !ok ***REMOVED***
		var (
			err   error
			isDir bool
		)

		// Workaround for https://github.com/spf13/afero/issues/201
		if isDir, err = afero.IsDir(i.fs, filename); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if isDir ***REMOVED***
			return nil, errors.New("open() can't be used with directories")
		***REMOVED***
		data, err = afero.ReadFile(i.fs, filename)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		i.files[filename] = data
	***REMOVED***

	if len(args) > 0 && args[0] == "b" ***REMOVED***
		return i.runtime.ToValue(data), nil
	***REMOVED***
	return i.runtime.ToValue(string(data)), nil
***REMOVED***
