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

package common

import (
	"context"
	"strconv"
	"testing"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testModule struct ***REMOVED***
	Counter int
***REMOVED***

func (testModule) unexported() bool ***REMOVED*** return true ***REMOVED***

func (testModule) Func() ***REMOVED******REMOVED***

func (testModule) Error() error ***REMOVED*** return errors.New("error") ***REMOVED***

func (testModule) Add(a, b int) int ***REMOVED*** return a + b ***REMOVED***

func (testModule) AddWithError(a, b int) (int, error) ***REMOVED***
	res := a + b
	if res < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return res, nil
***REMOVED***

func (testModule) Context(ctx context.Context) ***REMOVED******REMOVED***

func (testModule) ContextAdd(ctx context.Context, a, b int) int ***REMOVED***
	return a + b
***REMOVED***

func (testModule) ContextAddWithError(ctx context.Context, a, b int) (int, error) ***REMOVED***
	res := a + b
	if res < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return res, nil
***REMOVED***

func (m *testModule) Count() int ***REMOVED***
	m.Counter++
	return m.Counter
***REMOVED***

func TestModuleProxy(t *testing.T) ***REMOVED***
	testdata := map[string]func() interface***REMOVED******REMOVED******REMOVED***
		"Value":   func() interface***REMOVED******REMOVED*** ***REMOVED*** return testModule***REMOVED******REMOVED*** ***REMOVED***,
		"Pointer": func() interface***REMOVED******REMOVED*** ***REMOVED*** return &testModule***REMOVED******REMOVED*** ***REMOVED***,
	***REMOVED***
	for vtype, vfn := range testdata ***REMOVED***
		t.Run(vtype, func(t *testing.T) ***REMOVED***
			mod := Module***REMOVED******REMOVED***
			rt := goja.New()

			fns := map[string]func(interface***REMOVED******REMOVED***)***REMOVED***
				"Export": func(v interface***REMOVED******REMOVED***) ***REMOVED*** mod.Impl = v; rt.Set("mod", mod.Export(rt)) ***REMOVED***,
				"Proxy":  func(v interface***REMOVED******REMOVED***) ***REMOVED*** rt.Set("mod", mod.Proxy(rt, v)) ***REMOVED***,
			***REMOVED***
			for name, fn := range fns ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					v := vfn()
					fn(v)

					t.Run("unexported", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `mod.unexported()`)
						assert.EqualError(t, err, "TypeError: Object has no member 'unexported'")
					***REMOVED***)
					t.Run("Func", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `mod.func()`)
						assert.NoError(t, err)
					***REMOVED***)
					t.Run("Error", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `mod.error()`)
						assert.EqualError(t, err, "GoError: error")
					***REMOVED***)
					t.Run("Add", func(t *testing.T) ***REMOVED***
						v, err := RunString(rt, `mod.add(1, 2)`)
						assert.NoError(t, err)
						assert.Equal(t, int64(3), v.Export())
					***REMOVED***)
					t.Run("AddWithError", func(t *testing.T) ***REMOVED***
						v, err := RunString(rt, `mod.addWithError(1, 2)`)
						assert.NoError(t, err)
						assert.Equal(t, int64(3), v.Export())

						t.Run("Negative", func(t *testing.T) ***REMOVED***
							_, err := RunString(rt, `mod.addWithError(0, -1)`)
							assert.EqualError(t, err, "GoError: answer is negative")
						***REMOVED***)
					***REMOVED***)
					t.Run("Context", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `mod.context()`)
						assert.EqualError(t, err, "GoError: Context needs a valid VU context")

						t.Run("Valid", func(t *testing.T) ***REMOVED***
							mod.Context = context.Background()
							defer func() ***REMOVED*** mod.Context = nil ***REMOVED***()

							_, err := RunString(rt, `mod.context()`)
							assert.NoError(t, err)
						***REMOVED***)

						t.Run("Expired", func(t *testing.T) ***REMOVED***
							ctx, cancel := context.WithCancel(context.Background())
							cancel()
							mod.Context = ctx
							defer func() ***REMOVED*** mod.Context = nil ***REMOVED***()

							_, err := RunString(rt, `mod.context()`)
							assert.EqualError(t, err, "GoError: test has ended")
						***REMOVED***)
					***REMOVED***)
					t.Run("ContextAdd", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `mod.contextAdd(1, 2)`)
						assert.EqualError(t, err, "GoError: ContextAdd needs a valid VU context")

						t.Run("Valid", func(t *testing.T) ***REMOVED***
							mod.Context = context.Background()
							defer func() ***REMOVED*** mod.Context = nil ***REMOVED***()

							v, err := RunString(rt, `mod.contextAdd(1, 2)`)
							assert.NoError(t, err)
							assert.Equal(t, int64(3), v.Export())
						***REMOVED***)

						t.Run("Expired", func(t *testing.T) ***REMOVED***
							ctx, cancel := context.WithCancel(context.Background())
							cancel()
							mod.Context = ctx
							defer func() ***REMOVED*** mod.Context = nil ***REMOVED***()

							_, err := RunString(rt, `mod.contextAdd(1, 2)`)
							assert.EqualError(t, err, "GoError: test has ended")
						***REMOVED***)
					***REMOVED***)
					t.Run("ContextAddWithError", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `mod.contextAddWithError(1, 2)`)
						assert.EqualError(t, err, "GoError: ContextAddWithError needs a valid VU context")

						t.Run("Valid", func(t *testing.T) ***REMOVED***
							mod.Context = context.Background()
							defer func() ***REMOVED*** mod.Context = nil ***REMOVED***()

							v, err := RunString(rt, `mod.contextAddWithError(1, 2)`)
							assert.NoError(t, err)
							assert.Equal(t, int64(3), v.Export())

							t.Run("Negative", func(t *testing.T) ***REMOVED***
								_, err := RunString(rt, `mod.contextAddWithError(0, -1)`)
								assert.EqualError(t, err, "GoError: answer is negative")
							***REMOVED***)
						***REMOVED***)
						t.Run("Expired", func(t *testing.T) ***REMOVED***
							ctx, cancel := context.WithCancel(context.Background())
							cancel()
							mod.Context = ctx
							defer func() ***REMOVED*** mod.Context = nil ***REMOVED***()

							_, err := RunString(rt, `mod.contextAddWithError(1, 2)`)
							assert.EqualError(t, err, "GoError: test has ended")
						***REMOVED***)
					***REMOVED***)
					if impl, ok := v.(*testModule); ok ***REMOVED***
						t.Run("Count", func(t *testing.T) ***REMOVED***
							for i := 0; i < 10; i++ ***REMOVED***
								t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
									v, err := RunString(rt, `mod.count()`)
									assert.NoError(t, err)
									assert.Equal(t, int64(i+1), v.Export())
									assert.Equal(t, i+1, impl.Counter)
								***REMOVED***)
							***REMOVED***
						***REMOVED***)
					***REMOVED*** else ***REMOVED***
						t.Run("Count", func(t *testing.T) ***REMOVED***
							_, err := RunString(rt, `mod.count()`)
							assert.EqualError(t, err, "TypeError: Object has no member 'count'")
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
