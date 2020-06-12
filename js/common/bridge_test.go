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

type bridgeTestFieldsType struct ***REMOVED***
	Exported       string
	ExportedTag    string `js:"renamed"`
	ExportedHidden string `js:"-"`
	unexported     string
	unexportedTag  string `js:"unexported"`
***REMOVED***

type bridgeTestMethodsType struct***REMOVED******REMOVED***

func (bridgeTestMethodsType) ExportedFn() ***REMOVED******REMOVED***

//lint:ignore U1000 needed for the actual test to check that it won't be seen
func (bridgeTestMethodsType) unexportedFn() ***REMOVED******REMOVED***

func (*bridgeTestMethodsType) ExportedPtrFn() ***REMOVED******REMOVED***

//lint:ignore U1000 needed for the actual test to check that it won't be seen
func (*bridgeTestMethodsType) unexportedPtrFn() ***REMOVED******REMOVED***

type bridgeTestOddFieldsType struct ***REMOVED***
	TwoWords string
	URL      string
***REMOVED***

type bridgeTestErrorType struct***REMOVED******REMOVED***

func (bridgeTestErrorType) Error() error ***REMOVED*** return errors.New("error") ***REMOVED***

type bridgeTestJSValueType struct***REMOVED******REMOVED***

func (bridgeTestJSValueType) Func(arg goja.Value) goja.Value ***REMOVED*** return arg ***REMOVED***

type bridgeTestJSValueErrorType struct***REMOVED******REMOVED***

func (bridgeTestJSValueErrorType) Func(arg goja.Value) (goja.Value, error) ***REMOVED***
	if goja.IsUndefined(arg) ***REMOVED***
		return goja.Undefined(), errors.New("missing argument")
	***REMOVED***
	return arg, nil
***REMOVED***

type bridgeTestJSValueContextType struct***REMOVED******REMOVED***

func (bridgeTestJSValueContextType) Func(ctx context.Context, arg goja.Value) goja.Value ***REMOVED***
	return arg
***REMOVED***

type bridgeTestJSValueContextErrorType struct***REMOVED******REMOVED***

func (bridgeTestJSValueContextErrorType) Func(ctx context.Context, arg goja.Value) (goja.Value, error) ***REMOVED***
	if goja.IsUndefined(arg) ***REMOVED***
		return goja.Undefined(), errors.New("missing argument")
	***REMOVED***
	return arg, nil
***REMOVED***

type bridgeTestNativeFunctionType struct***REMOVED******REMOVED***

func (bridgeTestNativeFunctionType) Func(call goja.FunctionCall) goja.Value ***REMOVED***
	return call.Argument(0)
***REMOVED***

type bridgeTestNativeFunctionErrorType struct***REMOVED******REMOVED***

func (bridgeTestNativeFunctionErrorType) Func(call goja.FunctionCall) (goja.Value, error) ***REMOVED***
	arg := call.Argument(0)
	if goja.IsUndefined(arg) ***REMOVED***
		return goja.Undefined(), errors.New("missing argument")
	***REMOVED***
	return arg, nil
***REMOVED***

type bridgeTestNativeFunctionContextType struct***REMOVED******REMOVED***

func (bridgeTestNativeFunctionContextType) Func(ctx context.Context, call goja.FunctionCall) goja.Value ***REMOVED***
	return call.Argument(0)
***REMOVED***

type bridgeTestNativeFunctionContextErrorType struct***REMOVED******REMOVED***

func (bridgeTestNativeFunctionContextErrorType) Func(ctx context.Context, call goja.FunctionCall) (goja.Value, error) ***REMOVED***
	arg := call.Argument(0)
	if goja.IsUndefined(arg) ***REMOVED***
		return goja.Undefined(), errors.New("missing argument")
	***REMOVED***
	return arg, nil
***REMOVED***

type bridgeTestAddType struct***REMOVED******REMOVED***

func (bridgeTestAddType) Add(a, b int) int ***REMOVED*** return a + b ***REMOVED***

type bridgeTestAddWithErrorType struct***REMOVED******REMOVED***

func (bridgeTestAddWithErrorType) AddWithError(a, b int) (int, error) ***REMOVED***
	res := a + b
	if res < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return res, nil
***REMOVED***

type bridgeTestContextType struct***REMOVED******REMOVED***

func (bridgeTestContextType) Context(ctx context.Context) ***REMOVED******REMOVED***

type bridgeTestContextAddType struct***REMOVED******REMOVED***

func (bridgeTestContextAddType) ContextAdd(ctx context.Context, a, b int) int ***REMOVED***
	return a + b
***REMOVED***

type bridgeTestContextAddWithErrorType struct***REMOVED******REMOVED***

func (bridgeTestContextAddWithErrorType) ContextAddWithError(ctx context.Context, a, b int) (int, error) ***REMOVED***
	res := a + b
	if res < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return res, nil
***REMOVED***

type bridgeTestContextInjectType struct ***REMOVED***
	ctx context.Context
***REMOVED***

func (t *bridgeTestContextInjectType) ContextInject(ctx context.Context) ***REMOVED*** t.ctx = ctx ***REMOVED***

type bridgeTestContextInjectPtrType struct ***REMOVED***
	ctxPtr *context.Context
***REMOVED***

func (t *bridgeTestContextInjectPtrType) ContextInjectPtr(ctxPtr *context.Context) ***REMOVED*** t.ctxPtr = ctxPtr ***REMOVED***

type bridgeTestSumType struct***REMOVED******REMOVED***

func (bridgeTestSumType) Sum(nums ...int) int ***REMOVED***
	sum := 0
	for v := range nums ***REMOVED***
		sum += v
	***REMOVED***
	return sum
***REMOVED***

type bridgeTestSumWithContextType struct***REMOVED******REMOVED***

func (bridgeTestSumWithContextType) SumWithContext(ctx context.Context, nums ...int) int ***REMOVED***
	sum := 0
	for v := range nums ***REMOVED***
		sum += v
	***REMOVED***
	return sum
***REMOVED***

type bridgeTestSumWithErrorType struct***REMOVED******REMOVED***

func (bridgeTestSumWithErrorType) SumWithError(nums ...int) (int, error) ***REMOVED***
	sum := 0
	for v := range nums ***REMOVED***
		sum += v
	***REMOVED***
	if sum < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return sum, nil
***REMOVED***

type bridgeTestSumWithContextAndErrorType struct***REMOVED******REMOVED***

func (m bridgeTestSumWithContextAndErrorType) SumWithContextAndError(ctx context.Context, nums ...int) (int, error) ***REMOVED***
	sum := 0
	for v := range nums ***REMOVED***
		sum += v
	***REMOVED***
	if sum < 0 ***REMOVED***
		return 0, errors.New("answer is negative")
	***REMOVED***
	return sum, nil
***REMOVED***

type bridgeTestCounterType struct ***REMOVED***
	Counter int
***REMOVED***

func (m *bridgeTestCounterType) Count() int ***REMOVED***
	m.Counter++
	return m.Counter
***REMOVED***

type bridgeTestConstructorType struct***REMOVED******REMOVED***

type bridgeTestConstructorSpawnedType struct***REMOVED******REMOVED***

func (bridgeTestConstructorType) XConstructor() bridgeTestConstructorSpawnedType ***REMOVED***
	return bridgeTestConstructorSpawnedType***REMOVED******REMOVED***
***REMOVED***

func TestFieldNameMapper(t *testing.T) ***REMOVED***
	testdata := []struct ***REMOVED***
		Typ     reflect.Type
		Fields  map[string]string
		Methods map[string]string
	***REMOVED******REMOVED***
		***REMOVED***reflect.TypeOf(bridgeTestFieldsType***REMOVED******REMOVED***), map[string]string***REMOVED***
			"Exported":       "exported",
			"ExportedTag":    "renamed",
			"ExportedHidden": "",
			"unexported":     "",
			"unexportedTag":  "",
		***REMOVED***, nil***REMOVED***,
		***REMOVED***reflect.TypeOf(bridgeTestMethodsType***REMOVED******REMOVED***), nil, map[string]string***REMOVED***
			"ExportedFn":   "exportedFn",
			"unexportedFn": "",
		***REMOVED******REMOVED***,
		***REMOVED***reflect.TypeOf(bridgeTestOddFieldsType***REMOVED******REMOVED***), map[string]string***REMOVED***
			"TwoWords": "two_words",
			"URL":      "url",
		***REMOVED***, nil***REMOVED***,
		***REMOVED***reflect.TypeOf(bridgeTestConstructorType***REMOVED******REMOVED***), nil, map[string]string***REMOVED***
			"XConstructor": "Constructor",
		***REMOVED******REMOVED***,
	***REMOVED***
	for _, data := range testdata ***REMOVED***
		for field, name := range data.Fields ***REMOVED***
			t.Run(field, func(t *testing.T) ***REMOVED***
				f, ok := data.Typ.FieldByName(field)
				if assert.True(t, ok, "no such field") ***REMOVED***
					assert.Equal(t, name, (FieldNameMapper***REMOVED******REMOVED***).FieldName(data.Typ, f))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		for meth, name := range data.Methods ***REMOVED***
			t.Run(meth, func(t *testing.T) ***REMOVED***
				m, ok := data.Typ.MethodByName(meth)
				if name != "" ***REMOVED***
					if assert.True(t, ok, "no such method") ***REMOVED***
						assert.Equal(t, name, (FieldNameMapper***REMOVED******REMOVED***).MethodName(data.Typ, m))
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					assert.False(t, ok, "exported by accident")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBindToGlobal(t *testing.T) ***REMOVED***
	rt := goja.New()
	unbind := BindToGlobal(rt, map[string]interface***REMOVED******REMOVED******REMOVED***"a": 1***REMOVED***)
	assert.Equal(t, int64(1), rt.Get("a").Export())
	unbind()
	assert.Nil(t, rt.Get("a").Export())
***REMOVED***

func TestBind(t *testing.T) ***REMOVED***
	ctxPtr := new(context.Context)
	testdata := []struct ***REMOVED***
		Name string
		V    interface***REMOVED******REMOVED***
		Fn   func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime)
	***REMOVED******REMOVED***
		***REMOVED***"Fields", bridgeTestFieldsType***REMOVED***"a", "b", "c", "d", "e"***REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
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
		***REMOVED******REMOVED***,
		***REMOVED***"Methods", bridgeTestMethodsType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			t.Run("unexportedFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.unexportedFn()`)
				assert.EqualError(t, err, "TypeError: Object has no member 'unexportedFn' at <eval>:1:17(3)")
			***REMOVED***)
			t.Run("ExportedFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.exportedFn()`)
				assert.NoError(t, err)
			***REMOVED***)
			t.Run("unexportedPtrFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.unexportedPtrFn()`)
				assert.EqualError(t, err, "TypeError: Object has no member 'unexportedPtrFn' at <eval>:1:20(3)")
			***REMOVED***)
			t.Run("ExportedPtrFn", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.exportedPtrFn()`)
				switch obj.(type) ***REMOVED***
				case *bridgeTestMethodsType:
					assert.NoError(t, err)
				case bridgeTestMethodsType:
					assert.EqualError(t, err, "TypeError: Object has no member 'exportedPtrFn' at <eval>:1:18(3)")
				default:
					assert.Fail(t, "INVALID TYPE")
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"Error", bridgeTestErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.error()`)
			assert.EqualError(t, err, "GoError: error")
		***REMOVED******REMOVED***,
		***REMOVED***"JSValue", bridgeTestJSValueType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			v, err := RunString(rt, `obj.func(1234)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1234), v.Export())
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"JSValueError", bridgeTestJSValueErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.func()`)
			assert.EqualError(t, err, "GoError: missing argument")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.func(1234)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(1234), v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"JSValueContext", bridgeTestJSValueContextType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.func()`)
			assert.EqualError(t, err, "GoError: func() can only be called from within default()")

			t.Run("Context", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				v, err := RunString(rt, `obj.func(1234)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(1234), v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"JSValueContextError", bridgeTestJSValueContextErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.func()`)
			assert.EqualError(t, err, "GoError: func() can only be called from within default()")

			t.Run("Context", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				_, err := RunString(rt, `obj.func()`)
				assert.EqualError(t, err, "GoError: missing argument")

				t.Run("Valid", func(t *testing.T) ***REMOVED***
					v, err := RunString(rt, `obj.func(1234)`)
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(1234), v.Export())
					***REMOVED***
				***REMOVED***)
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"NativeFunction", bridgeTestNativeFunctionType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			v, err := RunString(rt, `obj.func(1234)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1234), v.Export())
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"NativeFunctionError", bridgeTestNativeFunctionErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.func()`)
			assert.EqualError(t, err, "GoError: missing argument")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				v, err := RunString(rt, `obj.func(1234)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(1234), v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"NativeFunctionContext", bridgeTestNativeFunctionContextType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.func()`)
			assert.EqualError(t, err, "GoError: func() can only be called from within default()")

			t.Run("Context", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				v, err := RunString(rt, `obj.func(1234)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(1234), v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"NativeFunctionContextError", bridgeTestNativeFunctionContextErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.func()`)
			assert.EqualError(t, err, "GoError: func() can only be called from within default()")

			t.Run("Context", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				_, err := RunString(rt, `obj.func()`)
				assert.EqualError(t, err, "GoError: missing argument")

				t.Run("Valid", func(t *testing.T) ***REMOVED***
					v, err := RunString(rt, `obj.func(1234)`)
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(1234), v.Export())
					***REMOVED***
				***REMOVED***)
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"Add", bridgeTestAddType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			v, err := RunString(rt, `obj.add(1, 2)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"AddWithError", bridgeTestAddWithErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			v, err := RunString(rt, `obj.addWithError(1, 2)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***

			t.Run("Negative", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.addWithError(0, -1)`)
				assert.EqualError(t, err, "GoError: answer is negative")
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"AddWithError", bridgeTestAddWithErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			v, err := RunString(rt, `obj.addWithError(1, 2)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***

			t.Run("Negative", func(t *testing.T) ***REMOVED***
				_, err := RunString(rt, `obj.addWithError(0, -1)`)
				assert.EqualError(t, err, "GoError: answer is negative")
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"Context", bridgeTestContextType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.context()`)
			assert.EqualError(t, err, "GoError: context() can only be called from within default()")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				_, err := RunString(rt, `obj.context()`)
				assert.NoError(t, err)
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"ContextAdd", bridgeTestContextAddType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.contextAdd(1, 2)`)
			assert.EqualError(t, err, "GoError: contextAdd() can only be called from within default()")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				v, err := RunString(rt, `obj.contextAdd(1, 2)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(3), v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"ContextAddWithError", bridgeTestContextAddWithErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.contextAddWithError(1, 2)`)
			assert.EqualError(t, err, "GoError: contextAddWithError() can only be called from within default()")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				v, err := RunString(rt, `obj.contextAddWithError(1, 2)`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, int64(3), v.Export())
				***REMOVED***

				t.Run("Negative", func(t *testing.T) ***REMOVED***
					_, err := RunString(rt, `obj.contextAddWithError(0, -1)`)
					assert.EqualError(t, err, "GoError: answer is negative")
				***REMOVED***)
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"ContextInject", bridgeTestContextInjectType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.contextInject()`)
			switch impl := obj.(type) ***REMOVED***
			case bridgeTestContextInjectType:
				assert.EqualError(t, err, "TypeError: Object has no member 'contextInject' at <eval>:1:18(3)")
			case *bridgeTestContextInjectType:
				assert.EqualError(t, err, "GoError: contextInject() can only be called from within default()")
				assert.Equal(t, nil, impl.ctx)

				t.Run("Valid", func(t *testing.T) ***REMOVED***
					*ctxPtr = context.Background()
					defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

					_, err := RunString(rt, `obj.contextInject()`)
					assert.NoError(t, err)
					assert.Equal(t, *ctxPtr, impl.ctx)
				***REMOVED***)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"ContextInjectPtr", bridgeTestContextInjectPtrType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.contextInjectPtr()`)
			switch impl := obj.(type) ***REMOVED***
			case bridgeTestContextInjectPtrType:
				assert.EqualError(t, err, "TypeError: Object has no member 'contextInjectPtr' at <eval>:1:21(3)")
			case *bridgeTestContextInjectPtrType:
				assert.NoError(t, err)
				assert.Equal(t, ctxPtr, impl.ctxPtr)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"Count", bridgeTestCounterType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			switch impl := obj.(type) ***REMOVED***
			case *bridgeTestCounterType:
				for i := 0; i < 10; i++ ***REMOVED***
					t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
						v, err := RunString(rt, `obj.count()`)
						if assert.NoError(t, err) ***REMOVED***
							assert.Equal(t, int64(i+1), v.Export())
							assert.Equal(t, i+1, impl.Counter)
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			case bridgeTestCounterType:
				_, err := RunString(rt, `obj.count()`)
				assert.EqualError(t, err, "TypeError: Object has no member 'count' at <eval>:1:10(3)")
			default:
				assert.Fail(t, "UNKNOWN TYPE")
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"Sum", bridgeTestSumType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			sum := 0
			args := []string***REMOVED******REMOVED***
			for i := 0; i < 10; i++ ***REMOVED***
				args = append(args, strconv.Itoa(i))
				sum += i
				t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
					code := fmt.Sprintf(`obj.sum(%s)`, strings.Join(args, ", "))
					v, err := RunString(rt, code)
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(sum), v.Export())
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"SumWithContext", bridgeTestSumWithContextType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.sumWithContext(1, 2)`)
			assert.EqualError(t, err, "GoError: sumWithContext() can only be called from within default()")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				sum := 0
				args := []string***REMOVED******REMOVED***
				for i := 0; i < 10; i++ ***REMOVED***
					args = append(args, strconv.Itoa(i))
					sum += i
					t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
						code := fmt.Sprintf(`obj.sumWithContext(%s)`, strings.Join(args, ", "))
						v, err := RunString(rt, code)
						if assert.NoError(t, err) ***REMOVED***
							assert.Equal(t, int64(sum), v.Export())
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"SumWithError", bridgeTestSumWithErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			sum := 0
			args := []string***REMOVED******REMOVED***
			for i := 0; i < 10; i++ ***REMOVED***
				args = append(args, strconv.Itoa(i))
				sum += i
				t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
					code := fmt.Sprintf(`obj.sumWithError(%s)`, strings.Join(args, ", "))
					v, err := RunString(rt, code)
					if assert.NoError(t, err) ***REMOVED***
						assert.Equal(t, int64(sum), v.Export())
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"SumWithContextAndError", bridgeTestSumWithContextAndErrorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			_, err := RunString(rt, `obj.sumWithContextAndError(1, 2)`)
			assert.EqualError(t, err, "GoError: sumWithContextAndError() can only be called from within default()")

			t.Run("Valid", func(t *testing.T) ***REMOVED***
				*ctxPtr = context.Background()
				defer func() ***REMOVED*** *ctxPtr = nil ***REMOVED***()

				sum := 0
				args := []string***REMOVED******REMOVED***
				for i := 0; i < 10; i++ ***REMOVED***
					args = append(args, strconv.Itoa(i))
					sum += i
					t.Run(strconv.Itoa(i), func(t *testing.T) ***REMOVED***
						code := fmt.Sprintf(`obj.sumWithContextAndError(%s)`, strings.Join(args, ", "))
						v, err := RunString(rt, code)
						if assert.NoError(t, err) ***REMOVED***
							assert.Equal(t, int64(sum), v.Export())
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			***REMOVED***)
		***REMOVED******REMOVED***,
		***REMOVED***"Constructor", bridgeTestConstructorType***REMOVED******REMOVED***, func(t *testing.T, obj interface***REMOVED******REMOVED***, rt *goja.Runtime) ***REMOVED***
			v, err := RunString(rt, `new obj.Constructor()`)
			assert.NoError(t, err)
			assert.IsType(t, bridgeTestConstructorSpawnedType***REMOVED******REMOVED***, v.Export())
		***REMOVED******REMOVED***,
	***REMOVED***

	vfns := map[string]func(interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED******REMOVED***
		"Value": func(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED*** return v ***REMOVED***,
		"Pointer": func(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
			val := reflect.ValueOf(v)
			ptr := reflect.New(val.Type())
			ptr.Elem().Set(val)
			return ptr.Interface()
		***REMOVED***,
	***REMOVED***

	for name, vfn := range vfns ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			for _, data := range testdata ***REMOVED***
				t.Run(data.Name, func(t *testing.T) ***REMOVED***
					rt := goja.New()
					rt.SetFieldNameMapper(FieldNameMapper***REMOVED******REMOVED***)
					obj := vfn(data.V)
					rt.Set("obj", Bind(rt, obj, ctxPtr))
					data.Fn(t, obj, rt)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkProxy(b *testing.B) ***REMOVED***
	types := []struct ***REMOVED***
		Name, FnName string
		Value        interface***REMOVED******REMOVED***
		Fn           func(b *testing.B, fn interface***REMOVED******REMOVED***)
	***REMOVED******REMOVED***
		***REMOVED***"Fields", "", bridgeTestFieldsType***REMOVED******REMOVED***, nil***REMOVED***,
		***REMOVED***"Methods", "exportedFn", bridgeTestMethodsType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			f := fn.(func())
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f()
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"Error", "", bridgeTestErrorType***REMOVED******REMOVED***, nil***REMOVED***,
		***REMOVED***"Add", "add", bridgeTestAddType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			f := fn.(func(int, int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"AddError", "addWithError", bridgeTestAddWithErrorType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func(int, int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"Context", "context", bridgeTestContextType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func())
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f()
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"ContextAdd", "contextAdd", bridgeTestContextAddType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func(int, int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"ContextAddError", "contextAddWithError", bridgeTestContextAddWithErrorType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func(int, int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"Sum", "sum", bridgeTestSumType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			f := fn.(func(...int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2, 3)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"SumContext", "sumWithContext", bridgeTestSumWithContextType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func(...int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2, 3)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"SumError", "sumWithError", bridgeTestSumWithErrorType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func(...int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2, 3)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"SumContextError", "sumWithContextAndError", bridgeTestSumWithContextAndErrorType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			b.Skip()
			f := fn.(func(...int) int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				f(1, 2, 3)
			***REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***"Constructor", "Constructor", bridgeTestConstructorType***REMOVED******REMOVED***, func(b *testing.B, fn interface***REMOVED******REMOVED***) ***REMOVED***
			f, _ := goja.AssertFunction(fn.(goja.Value))
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				_, _ = f(goja.Undefined())
			***REMOVED***
		***REMOVED******REMOVED***,
	***REMOVED***
	vfns := []struct ***REMOVED***
		Name string
		Fn   func(interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***"Value", func(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED*** return v ***REMOVED******REMOVED***,
		***REMOVED***"Pointer", func(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
			val := reflect.ValueOf(v)
			ptr := reflect.New(val.Type())
			ptr.Elem().Set(val)
			return ptr.Interface()
		***REMOVED******REMOVED***,
	***REMOVED***

	for _, vfn := range vfns ***REMOVED***
		b.Run(vfn.Name, func(b *testing.B) ***REMOVED***
			for _, typ := range types ***REMOVED***
				b.Run(typ.Name, func(b *testing.B) ***REMOVED***
					v := vfn.Fn(typ.Value)

					b.Run("ToValue", func(b *testing.B) ***REMOVED***
						rt := goja.New()
						rt.SetFieldNameMapper(FieldNameMapper***REMOVED******REMOVED***)
						b.ResetTimer()

						for i := 0; i < b.N; i++ ***REMOVED***
							rt.ToValue(v)
						***REMOVED***
					***REMOVED***)

					b.Run("Bridge", func(b *testing.B) ***REMOVED***
						rt := goja.New()
						rt.SetFieldNameMapper(FieldNameMapper***REMOVED******REMOVED***)
						ctx := context.Background()
						b.ResetTimer()

						for i := 0; i < b.N; i++ ***REMOVED***
							Bind(rt, v, &ctx)
						***REMOVED***
					***REMOVED***)

					if typ.FnName != "" ***REMOVED***
						b.Run("Call", func(b *testing.B) ***REMOVED***
							rt := goja.New()
							rt.SetFieldNameMapper(FieldNameMapper***REMOVED******REMOVED***)
							ctx := context.Background()
							fn := Bind(rt, v, &ctx)[typ.FnName]
							typ.Fn(b, fn)
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
