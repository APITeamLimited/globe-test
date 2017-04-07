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
	"strings"

	"github.com/dop251/goja"
	"github.com/serenize/snaker"
)

var (
	ctxPtrT = reflect.TypeOf((*context.Context)(nil))
	ctxT    = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorT  = reflect.TypeOf((*error)(nil)).Elem()
)

// Returns the JS name for an exported struct field. The name is snake_cased, with respect for
// certain common initialisms (URL, ID, HTTP, etc).
func FieldName(t reflect.Type, f reflect.StructField) string ***REMOVED***
	// PkgPath is non-empty for unexported fields.
	if f.PkgPath != "" ***REMOVED***
		return ""
	***REMOVED***

	// Allow a `js:"name"` tag to override the default name.
	if tag := f.Tag.Get("js"); tag != "" ***REMOVED***
		// Matching encoding/json, `js:"-"` hides a field.
		if tag == "-" ***REMOVED***
			return ""
		***REMOVED***
		return tag
	***REMOVED***

	// Default to lowercasing the first character of the field name.
	return snaker.CamelToSnake(f.Name)
***REMOVED***

// Returns the JS name for an exported method. The first letter of the method's name is
// lowercased, otherwise it is unaltered.
func MethodName(t reflect.Type, m reflect.Method) string ***REMOVED***
	// PkgPath is non-empty for unexported methods.
	if m.PkgPath != "" ***REMOVED***
		return ""
	***REMOVED***

	// Lowercase the first character of the method name.
	return strings.ToLower(m.Name[0:1]) + m.Name[1:]
***REMOVED***

// FieldNameMapper for goja.Runtime.SetFieldNameMapper()
type FieldNameMapper struct***REMOVED******REMOVED***

func (FieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string ***REMOVED*** return FieldName(t, f) ***REMOVED***

func (FieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string ***REMOVED*** return MethodName(t, m) ***REMOVED***

// Binds an object's members to the global scope. Returns a function that un-binds them.
// Note that this will panic if passed something that isn't a struct; please don't do that.
func BindToGlobal(rt *goja.Runtime, data map[string]interface***REMOVED******REMOVED***) func() ***REMOVED***
	keys := make([]string, len(data))
	i := 0
	for k, v := range data ***REMOVED***
		rt.Set(k, v)
		keys[i] = k
		i++
	***REMOVED***

	return func() ***REMOVED***
		for _, k := range keys ***REMOVED***
			rt.Set(k, goja.Undefined())
		***REMOVED***
	***REMOVED***
***REMOVED***

func Bind(rt *goja.Runtime, v interface***REMOVED******REMOVED***, ctxPtr *context.Context) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	exports := make(map[string]interface***REMOVED******REMOVED***)

	val := reflect.ValueOf(v)
	typ := val.Type()
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		meth := typ.Method(i)
		name := MethodName(typ, meth)
		if name == "" ***REMOVED***
			continue
		***REMOVED***
		fn := val.Method(i)

		// Figure out if we want to do any wrapping of it.
		fnT := fn.Type()
		numIn := fnT.NumIn()
		numOut := fnT.NumOut()
		hasError := (numOut > 1 && fnT.Out(1) == errorT)
		wantsContext := false
		wantsContextPtr := false
		if numIn > 0 ***REMOVED***
			in0 := fnT.In(0)
			switch in0 ***REMOVED***
			case ctxT:
				wantsContext = true
			case ctxPtrT:
				wantsContextPtr = true
			***REMOVED***
		***REMOVED***
		if hasError || wantsContext || wantsContextPtr ***REMOVED***
			// Varadic functions are called a bit differently.
			varadic := fnT.IsVariadic()

			// Collect input types, but skip the context (if any).
			var in []reflect.Type
			if numIn > 0 ***REMOVED***
				inOffset := 0
				if wantsContext || wantsContextPtr ***REMOVED***
					inOffset = 1
				***REMOVED***
				in = make([]reflect.Type, numIn-inOffset)
				for i := inOffset; i < numIn; i++ ***REMOVED***
					in[i-inOffset] = fnT.In(i)
				***REMOVED***
			***REMOVED***

			// Collect the output type (if any). JS functions can only return a single value, but
			// allow returning an error, which will be thrown as a JS exception.
			var out []reflect.Type
			if numOut != 0 ***REMOVED***
				out = []reflect.Type***REMOVED***fnT.Out(0)***REMOVED***
			***REMOVED***

			wrappedFn := fn
			fn = reflect.MakeFunc(
				reflect.FuncOf(in, out, varadic),
				func(args []reflect.Value) []reflect.Value ***REMOVED***
					if wantsContext ***REMOVED***
						if ctxPtr == nil || *ctxPtr == nil ***REMOVED***
							Throw(rt, errors.New(fmt.Sprintf("%s needs a valid VU context", meth.Name)))
						***REMOVED***
						args = append([]reflect.Value***REMOVED***reflect.ValueOf(*ctxPtr)***REMOVED***, args...)
					***REMOVED*** else if wantsContextPtr ***REMOVED***
						args = append([]reflect.Value***REMOVED***reflect.ValueOf(ctxPtr)***REMOVED***, args...)
					***REMOVED***

					var res []reflect.Value
					if varadic ***REMOVED***
						res = wrappedFn.CallSlice(args)
					***REMOVED*** else ***REMOVED***
						res = wrappedFn.Call(args)
					***REMOVED***

					if hasError ***REMOVED***
						if !res[1].IsNil() ***REMOVED***
							Throw(rt, res[1].Interface().(error))
						***REMOVED***
						res = res[:1]
					***REMOVED***

					return res
				***REMOVED***,
			)
		***REMOVED***

		exports[name] = fn.Interface()
	***REMOVED***

	// If v is a pointer, we need to indirect it to access fields.
	if typ.Kind() == reflect.Ptr ***REMOVED***
		val = val.Elem()
		typ = val.Type()
	***REMOVED***
	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		field := typ.Field(i)
		name := FieldName(typ, field)
		if name != "" ***REMOVED***
			exports[name] = val.Field(i).Interface()
		***REMOVED***
	***REMOVED***

	return exports
***REMOVED***
