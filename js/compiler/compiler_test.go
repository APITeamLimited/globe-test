/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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
package compiler

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
)

func TestTransform(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("blank", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		src, _, err := c.Transform("", "test.js", nil)
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";`, src)
	***REMOVED***)
	t.Run("double-arrow", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		src, _, err := c.Transform("()=> true", "test.js", nil)
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";() => true;`, src)
	***REMOVED***)
	t.Run("longer", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		src, _, err := c.Transform(strings.Join([]string***REMOVED***
			`function add(a, b) ***REMOVED***`,
			`    return a + b;`,
			`***REMOVED***;`,
			``,
			`let res = add(1, 2);`,
		***REMOVED***, "\n"), "test.js", nil)
		assert.NoError(t, err)
		assert.Equal(t, strings.Join([]string***REMOVED***
			`"use strict";function add(a, b) ***REMOVED***`,
			`    return a + b;`,
			`***REMOVED***;`,
			``,
			`let res = add(1, 2);`,
		***REMOVED***, "\n"), src)
	***REMOVED***)

	t.Run("double-arrow with sourceMap", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		c.Options.SourceMapLoader = func(string) ([]byte, error) ***REMOVED*** return nil, errors.New("shouldn't be called") ***REMOVED***
		src, _, err := c.Transform("()=> true", "test.js", nil)
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";

() => true;
//# sourceMappingURL=k6://internal-should-not-leak/file.map`, src)
	***REMOVED***)
***REMOVED***

func TestCompile(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("ES5", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		src := `1+(function() ***REMOVED*** return 2; ***REMOVED***)()`
		pgm, code, err := c.Compile(src, "script.js", true)
		require.NoError(t, err)
		assert.Equal(t, src, code)
		v, err := goja.New().RunProgram(pgm)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, int64(3), v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("ES5 Wrap", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		src := `exports.d=1+(function() ***REMOVED*** return 2; ***REMOVED***)()`
		pgm, code, err := c.Compile(src, "script.js", false)
		require.NoError(t, err)
		assert.Equal(t, "(function(module, exports)***REMOVED***\nexports.d=1+(function() ***REMOVED*** return 2; ***REMOVED***)()\n***REMOVED***)\n", code)
		rt := goja.New()
		v, err := rt.RunProgram(pgm)
		if assert.NoError(t, err) ***REMOVED***
			fn, ok := goja.AssertFunction(v)
			if assert.True(t, ok, "not a function") ***REMOVED***
				exp := make(map[string]goja.Value)
				_, err := fn(goja.Undefined(), goja.Undefined(), rt.ToValue(exp))
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(3), exp["d"].Export())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	t.Run("ES5 Invalid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		src := `1+(function() ***REMOVED*** return 2; )()`
		c.Options.CompatibilityMode = lib.CompatibilityModeExtended
		_, _, err := c.Compile(src, "script.js", false)
		assert.IsType(t, &goja.Exception***REMOVED******REMOVED***, err)
		assert.Contains(t, err.Error(), `SyntaxError: script.js: Unexpected token (1:26)
> 1 | 1+(function() ***REMOVED*** return 2; )()`)
	***REMOVED***)
	t.Run("ES6", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		c.Options.CompatibilityMode = lib.CompatibilityModeExtended
		pgm, code, err := c.Compile(`class A ***REMOVED***nine()***REMOVED***return 9***REMOVED******REMOVED***; new A().nine()`, "script.js", true)
		require.NoError(t, err)
		assert.Equal(t, `"use strict";var _createClass = function () ***REMOVED***function defineProperties(target, props) ***REMOVED***for (var i = 0; i < props.length; i++) ***REMOVED***var descriptor = props[i];descriptor.enumerable = descriptor.enumerable || false;descriptor.configurable = true;if ("value" in descriptor) descriptor.writable = true;Object.defineProperty(target, descriptor.key, descriptor);***REMOVED******REMOVED***return function (Constructor, protoProps, staticProps) ***REMOVED***if (protoProps) defineProperties(Constructor.prototype, protoProps);if (staticProps) defineProperties(Constructor, staticProps);return Constructor;***REMOVED***;***REMOVED***();function _classCallCheck(instance, Constructor) ***REMOVED***if (!(instance instanceof Constructor)) ***REMOVED***throw new TypeError("Cannot call a class as a function");***REMOVED******REMOVED***let A = function () ***REMOVED***function A() ***REMOVED***_classCallCheck(this, A);***REMOVED***_createClass(A, [***REMOVED*** key: "nine", value: function nine() ***REMOVED***return 9;***REMOVED*** ***REMOVED***]);return A;***REMOVED***();;new A().nine();`,
			code)
		v, err := goja.New().RunProgram(pgm)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, int64(9), v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Wrap", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		c.Options.CompatibilityMode = lib.CompatibilityModeExtended
		pgm, code, err := c.Compile(`class A ***REMOVED***nine()***REMOVED***return 9***REMOVED******REMOVED***; exports.fn(new A().nine())`, "script.js", false)
		require.NoError(t, err)
		assert.Equal(t, `(function(module, exports)***REMOVED***
"use strict";var _createClass = function () ***REMOVED***function defineProperties(target, props) ***REMOVED***for (var i = 0; i < props.length; i++) ***REMOVED***var descriptor = props[i];descriptor.enumerable = descriptor.enumerable || false;descriptor.configurable = true;if ("value" in descriptor) descriptor.writable = true;Object.defineProperty(target, descriptor.key, descriptor);***REMOVED******REMOVED***return function (Constructor, protoProps, staticProps) ***REMOVED***if (protoProps) defineProperties(Constructor.prototype, protoProps);if (staticProps) defineProperties(Constructor, staticProps);return Constructor;***REMOVED***;***REMOVED***();function _classCallCheck(instance, Constructor) ***REMOVED***if (!(instance instanceof Constructor)) ***REMOVED***throw new TypeError("Cannot call a class as a function");***REMOVED******REMOVED***let A = function () ***REMOVED***function A() ***REMOVED***_classCallCheck(this, A);***REMOVED***_createClass(A, [***REMOVED*** key: "nine", value: function nine() ***REMOVED***return 9;***REMOVED*** ***REMOVED***]);return A;***REMOVED***();;exports.fn(new A().nine());
***REMOVED***)
`, code)
		rt := goja.New()
		v, err := rt.RunProgram(pgm)
		if assert.NoError(t, err) ***REMOVED***
			fn, ok := goja.AssertFunction(v)
			if assert.True(t, ok, "not a function") ***REMOVED***
				exp := make(map[string]goja.Value)
				var out interface***REMOVED******REMOVED***
				exp["fn"] = rt.ToValue(func(v goja.Value) ***REMOVED***
					out = v.Export()
				***REMOVED***)
				_, err := fn(goja.Undefined(), goja.Undefined(), rt.ToValue(exp))
				assert.NoError(t, err)
				assert.Equal(t, int64(9), out)
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		c := New(testutils.NewLogger(t))
		c.Options.CompatibilityMode = lib.CompatibilityModeExtended
		_, _, err := c.Compile(`1+(=>2)()`, "script.js", true)
		assert.IsType(t, &goja.Exception***REMOVED******REMOVED***, err)
		assert.Contains(t, err.Error(), `SyntaxError: script.js: Unexpected token (1:3)
> 1 | 1+(=>2)()`)
	***REMOVED***)
***REMOVED***

func TestCorruptSourceMap(t *testing.T) ***REMOVED***
	t.Parallel()
	corruptSourceMap := []byte(`***REMOVED***"mappings": 12***REMOVED***`) // 12 is a number not a string

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = ioutil.Discard
	hook := testutils.SimpleLogrusHook***REMOVED***
		HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel, logrus.WarnLevel***REMOVED***,
	***REMOVED***
	logger.AddHook(&hook)

	compiler := New(logger)
	compiler.Options = Options***REMOVED***
		Strict: true,
		SourceMapLoader: func(string) ([]byte, error) ***REMOVED***
			return corruptSourceMap, nil
		***REMOVED***,
	***REMOVED***
	_, _, err := compiler.Compile("var s = 5;\n//# sourceMappingURL=somefile", "somefile", false)
	require.NoError(t, err)
	entries := hook.Drain()
	require.Len(t, entries, 1)
	msg, err := entries[0].String() // we need this in order to get the field error
	require.NoError(t, err)

	require.Contains(t, msg, `Couldn't load source map for somefile`)
	require.Contains(t, msg, `json: cannot unmarshal number into Go struct field v3.mappings of type string`)
***REMOVED***

func TestCorruptSourceMapOnlyForBabel(t *testing.T) ***REMOVED***
	t.Parallel()
	// this a valid source map for the go implementation but babel doesn't like it
	corruptSourceMap := []byte(`***REMOVED***"mappings": ";"***REMOVED***`)

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = ioutil.Discard
	hook := testutils.SimpleLogrusHook***REMOVED***
		HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel, logrus.WarnLevel***REMOVED***,
	***REMOVED***
	logger.AddHook(&hook)

	compiler := New(logger)
	compiler.Options = Options***REMOVED***
		CompatibilityMode: lib.CompatibilityModeExtended,
		Strict:            true,
		SourceMapLoader: func(string) ([]byte, error) ***REMOVED***
			return corruptSourceMap, nil
		***REMOVED***,
	***REMOVED***
	_, _, err := compiler.Compile("class s ***REMOVED******REMOVED***;\n//# sourceMappingURL=somefile", "somefile", false)
	require.NoError(t, err)
	entries := hook.Drain()
	require.Len(t, entries, 1)
	msg, err := entries[0].String() // we need this in order to get the field error
	require.NoError(t, err)

	require.Contains(t, msg, `needs to be transpiled by Babel, but its source map will not be accepted by Babel`)
	require.Contains(t, msg, `source map missing required 'version' field`)
***REMOVED***

func TestMinimalSourceMap(t *testing.T) ***REMOVED***
	t.Parallel()
	// this is the minimal sourcemap valid for both go and babel implementations
	corruptSourceMap := []byte(`***REMOVED***"version":3,"mappings":";","sources":[]***REMOVED***`)

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = ioutil.Discard
	hook := testutils.SimpleLogrusHook***REMOVED***
		HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel, logrus.WarnLevel***REMOVED***,
	***REMOVED***
	logger.AddHook(&hook)

	compiler := New(logger)
	compiler.Options = Options***REMOVED***
		CompatibilityMode: lib.CompatibilityModeExtended,
		Strict:            true,
		SourceMapLoader: func(string) ([]byte, error) ***REMOVED***
			return corruptSourceMap, nil
		***REMOVED***,
	***REMOVED***
	_, _, err := compiler.Compile("class s ***REMOVED******REMOVED***;\n//# sourceMappingURL=somefile", "somefile", false)
	require.NoError(t, err)
	require.Empty(t, hook.Drain())
***REMOVED***
