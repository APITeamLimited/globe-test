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
	"reflect"
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

					if typ.FnName != "" ***REMOVED***
						b.Run("Call", func(b *testing.B) ***REMOVED***
							rt := goja.New()
							rt.SetFieldNameMapper(FieldNameMapper***REMOVED******REMOVED***)
							fn := func() ***REMOVED******REMOVED***
							typ.Fn(b, fn)
						***REMOVED***)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
