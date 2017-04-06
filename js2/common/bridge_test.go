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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

type bridgeTestType struct ***REMOVED***
	Exported      string
	ExportedTag   string `js:"renamed"`
	unexported    string
	unexportedTag string `js:"unexported"`

	TwoWords string
	URL      string

	Counter int
***REMOVED***

func (bridgeTestType) ExportedFn()   ***REMOVED******REMOVED***
func (bridgeTestType) unexportedFn() ***REMOVED******REMOVED***

func (*bridgeTestType) ExportedPtrFn()   ***REMOVED******REMOVED***
func (*bridgeTestType) unexportedPtrFn() ***REMOVED******REMOVED***

func (bridgeTestType) Error() error ***REMOVED*** return errors.New("error") ***REMOVED***

func (bridgeTestType) Add(a, b int) int ***REMOVED*** return a + b ***REMOVED***

func (bridgeTestType) AddWithError(a, b int) (int, error) ***REMOVED***
	res := a + b
	if res < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return res, nil
***REMOVED***

func (bridgeTestType) Context(ctx context.Context) ***REMOVED******REMOVED***

func (bridgeTestType) ContextAdd(ctx context.Context, a, b int) int ***REMOVED***
	return a + b
***REMOVED***

func (bridgeTestType) ContextAddWithError(ctx context.Context, a, b int) (int, error) ***REMOVED***
	res := a + b
	if res < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return res, nil
***REMOVED***

func (m *bridgeTestType) Count() int ***REMOVED***
	m.Counter++
	return m.Counter
***REMOVED***

func (bridgeTestType) Sum(nums ...int) int ***REMOVED***
	sum := 0
	for v := range nums ***REMOVED***
		sum += v
	***REMOVED***
	return sum
***REMOVED***

func (m bridgeTestType) SumWithContext(ctx context.Context, nums ...int) int ***REMOVED***
	return m.Sum(nums...)
***REMOVED***

func (m bridgeTestType) SumWithError(nums ...int) (int, error) ***REMOVED***
	sum := m.Sum(nums...)
	if sum < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return sum, nil
***REMOVED***

func (m bridgeTestType) SumWithContextAndError(ctx context.Context, nums ...int) (int, error) ***REMOVED***
	return m.SumWithError(nums...)
***REMOVED***

func TestFieldNameMapper(t *testing.T) ***REMOVED***
	typ := reflect.TypeOf(bridgeTestType***REMOVED******REMOVED***)
	t.Run("Fields", func(t *testing.T) ***REMOVED***
		names := map[string]string***REMOVED***
			"Exported":      "exported",
			"ExportedTag":   "renamed",
			"unexported":    "",
			"unexportedTag": "",
			"TwoWords":      "two_words",
			"URL":           "url",
		***REMOVED***
		for name, result := range names ***REMOVED***
			t.Run(name, func(t *testing.T) ***REMOVED***
				f, ok := typ.FieldByName(name)
				if assert.True(t, ok) ***REMOVED***
					assert.Equal(t, result, (FieldNameMapper***REMOVED******REMOVED***).FieldName(typ, f))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Methods", func(t *testing.T) ***REMOVED***
		t.Run("ExportedFn", func(t *testing.T) ***REMOVED***
			m, ok := typ.MethodByName("ExportedFn")
			if assert.True(t, ok) ***REMOVED***
				assert.Equal(t, "exportedFn", (FieldNameMapper***REMOVED******REMOVED***).MethodName(typ, m))
			***REMOVED***
		***REMOVED***)
		t.Run("unexportedFn", func(t *testing.T) ***REMOVED***
			_, ok := typ.MethodByName("unexportedFn")
			assert.False(t, ok)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestBindToGlobal(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Obj  interface***REMOVED******REMOVED***
		Keys []string
		Not  []string
	***REMOVED******REMOVED***
		"Value": ***REMOVED***
			bridgeTestType***REMOVED******REMOVED***,
			[]string***REMOVED***"exported", "renamed", "exportedFn"***REMOVED***,
			[]string***REMOVED***"exportedPtrFn"***REMOVED***,
		***REMOVED***,
		"Pointer": ***REMOVED***
			&bridgeTestType***REMOVED******REMOVED***,
			[]string***REMOVED***"exported", "renamed", "exportedFn", "exportedPtrFn"***REMOVED***,
			[]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			rt := goja.New()
			unbind := BindToGlobal(rt, data.Obj)
			for _, k := range data.Keys ***REMOVED***
				t.Run(k, func(t *testing.T) ***REMOVED***
					v := rt.Get(k)
					if assert.NotNil(t, v) ***REMOVED***
						assert.False(t, goja.IsUndefined(v), "value is undefined")
					***REMOVED***
				***REMOVED***)
			***REMOVED***
			for _, k := range data.Not ***REMOVED***
				t.Run(k, func(t *testing.T) ***REMOVED***
					assert.Nil(t, rt.Get(k), "unexpected member bridged")
				***REMOVED***)
			***REMOVED***

			t.Run("Unbind", func(t *testing.T) ***REMOVED***
				unbind()
				for _, k := range data.Keys ***REMOVED***
					t.Run(k, func(t *testing.T) ***REMOVED***
						v := rt.Get(k)
						assert.True(t, goja.IsUndefined(v), "value is not undefined")
					***REMOVED***)
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBind(t *testing.T) ***REMOVED***
	template := bridgeTestType***REMOVED***
		Exported:      "a",
		ExportedTag:   "b",
		unexported:    "c",
		unexportedTag: "d",
	***REMOVED***
	testdata := map[string]func() interface***REMOVED******REMOVED******REMOVED***
		"Value":   func() interface***REMOVED******REMOVED*** ***REMOVED*** return template ***REMOVED***,
		"Pointer": func() interface***REMOVED******REMOVED*** ***REMOVED*** tmp := template; return &tmp ***REMOVED***,
	***REMOVED***
	for vtype, vfn := range testdata ***REMOVED***
		t.Run(vtype, func(t *testing.T) ***REMOVED***
			rt := goja.New()
			v := vfn()
			ctx := new(context.Context)
			rt.Set("obj", Bind(rt, v, ctx))

			t.Run("unexportedFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.unexportedFn()`)
				assert.EqualError(t, err, "TypeError: Object has no member 'unexportedFn'")
			***REMOVED***)
			t.Run("ExportedFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.exportedFn()`)
				assert.NoError(t, err)
			***REMOVED***)
			t.Run("unexportedPtrFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.unexportedPtrFn()`)
				assert.EqualError(t, err, "TypeError: Object has no member 'unexportedPtrFn'")
			***REMOVED***)
			t.Run("ExportedPtrFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.exportedPtrFn()`)
				if vtype == "Pointer" ***REMOVED***
					assert.NoError(t, err)
				***REMOVED*** else ***REMOVED***
					assert.EqualError(t, err, "TypeError: Object has no member 'exportedPtrFn'")
				***REMOVED***
			***REMOVED***)
			t.Run("Error", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.error()`)
				assert.EqualError(t, err, "GoError: error")
			***REMOVED***)
			t.Run("Add", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.add(1, 2)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(3), v.Export())
				***REMOVED***
			***REMOVED***)
			t.Run("AddWithError", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.addWithError(1, 2)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(3), v.Export())
				***REMOVED***

				t.Run("Negative", func(t *testing.T) ***REMOVED***
					_, err := RunString(rt, `obj.addWithError(0, -1)`)
					assert.EqualError(t, err, "GoError: answer is negative")
				***REMOVED***)
			***REMOVED***)
			t.Run("Context", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.context()`)
				assert.EqualError(t, err, "GoError: Context needs a valid VU context")

				t.Run("Valid", func(t *testing.T) ***REMOVED***
					*ctx = context.Background()
					defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

					_, err := RunString(rt, `obj.context()`)
					assert.NoError(t, err)
				***REMOVED***)

				t.Run("Expired", func(t *testing.T) ***REMOVED***
					ctx_, cancel := context.WithCancel(context.Background())
					cancel()
					*ctx = ctx_
					defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

					_, err := RunString(rt, `obj.context()`)
					assert.EqualError(t, err, "GoError: test has ended")
				***REMOVED***)
			***REMOVED***)
			t.Run("ContextAdd", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.contextAdd(1, 2)`)
				assert.EqualError(t, err, "GoError: ContextAdd needs a valid VU context")

				t.Run("Valid", func(t *testing.T) ***REMOVED***
					*ctx = context.Background()
					defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

					v, err := RunString(rt, `obj.contextAdd(1, 2)`)
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(3), v.Export())
					***REMOVED***
				***REMOVED***)

				t.Run("Expired", func(t *testing.T) ***REMOVED***
					ctx_, cancel := context.WithCancel(context.Background())
					cancel()
					*ctx = ctx_
					defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

					_, err := RunString(rt, `obj.contextAdd(1, 2)`)
					assert.EqualError(t, err, "GoError: test has ended")
				***REMOVED***)
			***REMOVED***)
			t.Run("ContextAddWithError", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.contextAddWithError(1, 2)`)
				assert.EqualError(t, err, "GoError: ContextAddWithError needs a valid VU context")

				t.Run("Valid", func(t *testing.T) ***REMOVED***
					*ctx = context.Background()
					defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

					v, err := RunString(rt, `obj.contextAddWithError(1, 2)`)
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(3), v.Export())
					***REMOVED***

					t.Run("Negative", func(t *testing.T) ***REMOVED***
						_, err := RunString(rt, `obj.contextAddWithError(0, -1)`)
						assert.EqualError(t, err, "GoError: answer is negative")
					***REMOVED***)
				***REMOVED***)
				t.Run("Expired", func(t *testing.T) ***REMOVED***
					ctx_, cancel := context.WithCancel(context.Background())
					cancel()
					*ctx = ctx_
					defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

					_, err := RunString(rt, `obj.contextAddWithError(1, 2)`)
					assert.EqualError(t, err, "GoError: test has ended")
				***REMOVED***)
			***REMOVED***)
			if impl, ok := v.(*bridgeTestType); ok ***REMOVED***
				t.Run("Count", func(t *testing.T) ***REMOVED***
					for i := 0; i < 10; i++ ***REMOVED***
						t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
							v, err := RunString(rt, `obj.count()`)
							if assert.NoError(t, err) ***REMOVED***
								assert.Equal(t, int64(i+1), v.Export())
								assert.Equal(t, i+1, impl.Counter)
							***REMOVED***
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				t.Run("Count", func(t *testing.T) ***REMOVED***
					_, err := RunString(rt, `obj.count()`)
					assert.EqualError(t, err, "TypeError: Object has no member 'count'")
				***REMOVED***)
			***REMOVED***
			for name, fname := range map[string]string***REMOVED***
				"Sum":                    "sum",
				"SumWithContext":         "sumWithContext",
				"SumWithError":           "sumWithError",
				"SumWithContextAndError": "sumWithContextAndError",
			***REMOVED*** ***REMOVED***
				*ctx = context.Background()
				defer func() ***REMOVED*** *ctx = nil ***REMOVED***()

				t.Run(name, func(t *testing.T) ***REMOVED***
					sum := 0
					args := []string***REMOVED******REMOVED***
					for i := 0; i < 10; i++ ***REMOVED***
						args = append(args, strconv.Itoa(i))
						sum += i
						t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
							code := fmt.Sprintf(`obj.%s(%s)`, fname, strings.Join(args, ", "))
							v, err := RunString(rt, code)
							if assert.NoError(t, err) ***REMOVED***
								assert.Equal(t, int64(sum), v.Export())
							***REMOVED***
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***

			t.Run("Exported", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.exported`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "a", v.Export())
				***REMOVED***
			***REMOVED***)
			t.Run("ExportedTag", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.renamed`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "b", v.Export())
				***REMOVED***
			***REMOVED***)
			t.Run("unexported", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.unexported`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, nil, v.Export())
				***REMOVED***
			***REMOVED***)
			t.Run("unexportedTag", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.unexportedTag`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, nil, v.Export())
				***REMOVED***
			***REMOVED***)

			t.Run("Counter", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.counter`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(0), v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkProxy(b *testing.B) ***REMOVED***
	var v bridgeTestType
	rt := goja.New()
	b.ResetTimer()

	b.Run("ToValue", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			_ = rt.ToValue(v)
		***REMOVED***
	***REMOVED***)

	b.Run("Bind", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			_ = Bind(rt, v, nil)
		***REMOVED***
	***REMOVED***)
***REMOVED***
