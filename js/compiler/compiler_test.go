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
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/stretchr/testify/assert"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
)

func TestTransform(t *testing.T) ***REMOVED***
	c := New(testutils.NewLogger(t))
	t.Run("blank", func(t *testing.T) ***REMOVED***
		src, _, err := c.Transform("", "test.js")
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";`, src)
		// assert.Equal(t, 3, srcmap.Version)
		// assert.Equal(t, "test.js", srcmap.File)
		// assert.Equal(t, "", srcmap.Mappings)
	***REMOVED***)
	t.Run("double-arrow", func(t *testing.T) ***REMOVED***
		src, _, err := c.Transform("()=> true", "test.js")
		assert.NoError(t, err)
		assert.Equal(t, `"use strict";(function () ***REMOVED***return true;***REMOVED***);`, src)
		// assert.Equal(t, 3, srcmap.Version)
		// assert.Equal(t, "test.js", srcmap.File)
		// assert.Equal(t, "aAAA,qBAAK,IAAL", srcmap.Mappings)
	***REMOVED***)
	t.Run("longer", func(t *testing.T) ***REMOVED***
		src, _, err := c.Transform(strings.Join([]string***REMOVED***
			`function add(a, b) ***REMOVED***`,
			`    return a + b;`,
			`***REMOVED***;`,
			``,
			`let res = add(1, 2);`,
		***REMOVED***, "\n"), "test.js")
		assert.NoError(t, err)
		assert.Equal(t, strings.Join([]string***REMOVED***
			`"use strict";function add(a, b) ***REMOVED***`,
			`    return a + b;`,
			`***REMOVED***;`,
			``,
			`let res = add(1, 2);`,
		***REMOVED***, "\n"), src)
		// assert.Equal(t, 3, srcmap.Version)
		// assert.Equal(t, "test.js", srcmap.File)
		// assert.Equal(t, "aAAA,SAASA,GAAT,CAAaC,CAAb,EAAgBC,CAAhB,EAAmB;AACf,WAAOD,IAAIC,CAAX;AACH;;AAED,IAAIC,MAAMH,IAAI,CAAJ,EAAO,CAAP,CAAV", srcmap.Mappings)
	***REMOVED***)
***REMOVED***

func TestCompile(t *testing.T) ***REMOVED***
	c := New(testutils.NewLogger(t))
	t.Run("ES5", func(t *testing.T) ***REMOVED***
		src := `1+(function() ***REMOVED*** return 2; ***REMOVED***)()`
		pgm, code, err := c.Compile(src, "script.js", "", "", true, lib.CompatibilityModeBase)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, src, code)
		v, err := goja.New().RunProgram(pgm)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, int64(3), v.Export())
		***REMOVED***

		t.Run("Wrap", func(t *testing.T) ***REMOVED***
			pgm, code, err := c.Compile(src, "script.js",
				"(function()***REMOVED***return ", "***REMOVED***)", true, lib.CompatibilityModeBase)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			assert.Equal(t, `(function()***REMOVED***return 1+(function() ***REMOVED*** return 2; ***REMOVED***)()***REMOVED***)`, code)
			v, err := goja.New().RunProgram(pgm)
			if assert.NoError(t, err) ***REMOVED***
				fn, ok := goja.AssertFunction(v)
				if assert.True(t, ok, "not a function") ***REMOVED***
					v, err := fn(goja.Undefined())
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(3), v.Export())
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***)

		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			src := `1+(function() ***REMOVED*** return 2; )()`
			_, _, err := c.Compile(src, "script.js", "", "", true, lib.CompatibilityModeExtended)
			assert.IsType(t, &goja.Exception***REMOVED******REMOVED***, err)
			assert.Contains(t, err.Error(), `SyntaxError: script.js: Unexpected token (1:26)
> 1 | 1+(function() ***REMOVED*** return 2; )()`)
		***REMOVED***)
	***REMOVED***)
	t.Run("ES6", func(t *testing.T) ***REMOVED***
		pgm, code, err := c.Compile(`1+(()=>2)()`, "script.js", "", "", true, lib.CompatibilityModeExtended)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, `"use strict";1 + function () ***REMOVED***return 2;***REMOVED***();`, code)
		v, err := goja.New().RunProgram(pgm)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, int64(3), v.Export())
		***REMOVED***

		t.Run("Wrap", func(t *testing.T) ***REMOVED***
			pgm, code, err := c.Compile(`fn(1+(()=>2)())`, "script.js", "(function(fn)***REMOVED***", "***REMOVED***)", true, lib.CompatibilityModeExtended)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			assert.Equal(t, `(function(fn)***REMOVED***"use strict";fn(1 + function () ***REMOVED***return 2;***REMOVED***());***REMOVED***)`, code)
			rt := goja.New()
			v, err := rt.RunProgram(pgm)
			if assert.NoError(t, err) ***REMOVED***
				fn, ok := goja.AssertFunction(v)
				if assert.True(t, ok, "not a function") ***REMOVED***
					var out interface***REMOVED******REMOVED***
					_, err := fn(goja.Undefined(), rt.ToValue(func(v goja.Value) ***REMOVED***
						out = v.Export()
					***REMOVED***))
					assert.NoError(t, err)
					assert.Equal(t, int64(3), out)
				***REMOVED***
			***REMOVED***
		***REMOVED***)

		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			_, _, err := c.Compile(`1+(=>2)()`, "script.js", "", "", true, lib.CompatibilityModeExtended)
			assert.IsType(t, &goja.Exception***REMOVED******REMOVED***, err)
			assert.Contains(t, err.Error(), `SyntaxError: script.js: Unexpected token (1:3)
> 1 | 1+(=>2)()`)
		***REMOVED***)

		t.Run("Invalid for goja but not babel", func(t *testing.T) ***REMOVED***
			t.Skip("Find something else that breaks this as this was fixed in goja :(")
			ch := make(chan struct***REMOVED******REMOVED***)
			go func() ***REMOVED***
				defer close(ch)
				// This is a string with U+2029 Paragraph separator in it
				// the important part is that goja won't parse it but babel will transform it but still
				// goja won't be able to parse the result it is actually "\<U+2029>"
				_, _, err := c.Compile(string([]byte***REMOVED***0x22, 0x5c, 0xe2, 0x80, 0xa9, 0x22***REMOVED***), "script.js", "", "", true, lib.CompatibilityModeExtended)
				assert.IsType(t, parser.ErrorList***REMOVED******REMOVED***, err)
				assert.Contains(t, err.Error(), ` Unexpected token ILLEGAL`)
			***REMOVED***()

			select ***REMOVED***
			case <-ch:
				// everything is fine
			case <-time.After(time.Second):
				// it took too long
				t.Fatal("takes too long")
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***
