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
	"reflect"
	"strings"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
)

var (
	ctxPtrT = reflect.TypeOf((*context.Context)(nil))
	ctxT    = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorT  = reflect.TypeOf((*error)(nil)).Elem()
	jsValT  = reflect.TypeOf((*goja.Value)(nil)).Elem()
	fnCallT = reflect.TypeOf((*goja.FunctionCall)(nil)).Elem()

	constructWrap = goja.MustCompile(
		"__constructor__",
		`(function(impl) ***REMOVED*** return function() ***REMOVED*** return impl.apply(this, arguments); ***REMOVED*** ***REMOVED***)`,
		true,
	)
)

// if a fieldName is the key of this map exactly than the value for the given key should be used as
// the name of the field in js
//nolint: gochecknoglobals
var fieldNameExceptions = map[string]string***REMOVED***
	"OCSP": "ocsp",
***REMOVED***

// FieldName Returns the JS name for an exported struct field. The name is snake_cased, with respect for
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

	if exception, ok := fieldNameExceptions[f.Name]; ok ***REMOVED***
		return exception
	***REMOVED***

	// Default to lowercasing the first character of the field name.
	return snaker.CamelToSnake(f.Name)
***REMOVED***

// if a methodName is the key of this map exactly than the value for the given key should be used as
// the name of the method in js
//nolint: gochecknoglobals
var methodNameExceptions = map[string]string***REMOVED***
	"JSON": "json",
	"HTML": "html",
	"URL":  "url",
	"OCSP": "ocsp",
***REMOVED***

// MethodName Returns the JS name for an exported method. The first letter of the method's name is
// lowercased, otherwise it is unaltered.
func MethodName(t reflect.Type, m reflect.Method) string ***REMOVED***
	// A field with a name beginning with an X is a constructor, and just gets the prefix stripped.
	// Note: They also get some special treatment from Bridge(), see further down.
	if m.Name[0] == 'X' ***REMOVED***
		return m.Name[1:]
	***REMOVED***

	if exception, ok := methodNameExceptions[m.Name]; ok ***REMOVED***
		return exception
	***REMOVED***
	// Lowercase the first character of the method name.
	return strings.ToLower(m.Name[0:1]) + m.Name[1:]
***REMOVED***

// FieldNameMapper for goja.Runtime.SetFieldNameMapper()
type FieldNameMapper struct***REMOVED******REMOVED***

// FieldName is part of the goja.FieldNameMapper interface
// https://godoc.org/github.com/dop251/goja#FieldNameMapper
func (FieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string ***REMOVED*** return FieldName(t, f) ***REMOVED***

// MethodName is part of the goja.FieldNameMapper interface
// https://godoc.org/github.com/dop251/goja#FieldNameMapper
func (FieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string ***REMOVED*** return MethodName(t, m) ***REMOVED***

// BindToGlobal Binds an object's members to the global scope. Returns a function that un-binds them.
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

// Bind the provided value v to the provided runtime
func Bind(rt *goja.Runtime, v interface***REMOVED******REMOVED***, ctxPtr *context.Context) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	exports := make(map[string]interface***REMOVED******REMOVED***)

	val := reflect.ValueOf(v)
	typ := val.Type()
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		meth := typ.Method(i)
		name := MethodName(typ, meth)
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
			isVariadic := fnT.IsVariadic()
			realFn := fn
			fn = reflect.ValueOf(func(call goja.FunctionCall) goja.Value ***REMOVED***
				// Number of arguments: the higher number between the function's required arguments
				// and the number of arguments actually given.
				args := make([]reflect.Value, numIn)

				// Inject any requested parameters, and reserve them to offset user args.
				reservedArgs := 0
				if wantsContext ***REMOVED***
					if ctxPtr == nil || *ctxPtr == nil ***REMOVED***
						Throw(rt, errors.Errorf("%s() can only be called from within default()", name))
					***REMOVED***
					args[0] = reflect.ValueOf(*ctxPtr)
					reservedArgs++
				***REMOVED*** else if wantsContextPtr ***REMOVED***
					args[0] = reflect.ValueOf(ctxPtr)
					reservedArgs++
				***REMOVED***

				// Copy over arguments.
				for i := 0; i < numIn; i++ ***REMOVED***
					if i < reservedArgs ***REMOVED***
						continue
					***REMOVED***

					T := fnT.In(i)

					// A function that takes a goja.FunctionCall takes only that arg (+ injected).
					if T == fnCallT ***REMOVED***
						args[i] = reflect.ValueOf(call)
						break
					***REMOVED***

					// The last arg to a varadic function is a slice of the remainder.
					if isVariadic && i == numIn-1 ***REMOVED***
						varArgsLen := len(call.Arguments) - (i - reservedArgs)
						if varArgsLen <= 0 ***REMOVED***
							args[i] = reflect.Zero(T)
							break
						***REMOVED***
						varArgs := reflect.MakeSlice(T, varArgsLen, varArgsLen)
						emT := T.Elem()
						for j := 0; j < varArgsLen; j++ ***REMOVED***
							arg := call.Arguments[i+j-reservedArgs]
							v := reflect.New(emT)
							if err := rt.ExportTo(arg, v.Interface()); err != nil ***REMOVED***
								Throw(rt, err)
							***REMOVED***
							varArgs.Index(j).Set(v.Elem())
						***REMOVED***
						args[i] = varArgs
						break
					***REMOVED***

					arg := call.Argument(i - reservedArgs)

					// Optimization: no need to allocate a pointer and export for a zero value.
					if goja.IsUndefined(arg) ***REMOVED***
						if T == jsValT ***REMOVED***
							args[i] = reflect.ValueOf(goja.Undefined())
							continue
						***REMOVED***
						args[i] = reflect.Zero(T)
						continue
					***REMOVED***

					// Allocate a T* and export the JS value to it.
					v := reflect.New(T)
					if err := rt.ExportTo(arg, v.Interface()); err != nil ***REMOVED***
						Throw(rt, err)
					***REMOVED***
					args[i] = v.Elem()
				***REMOVED***

				var ret []reflect.Value
				if isVariadic ***REMOVED***
					ret = realFn.CallSlice(args)
				***REMOVED*** else ***REMOVED***
					ret = realFn.Call(args)
				***REMOVED***

				if len(ret) > 0 ***REMOVED***
					if hasError && !ret[1].IsNil() ***REMOVED***
						Throw(rt, ret[1].Interface().(error))
					***REMOVED***
					return rt.ToValue(ret[0].Interface())
				***REMOVED***
				return goja.Undefined()
			***REMOVED***)
		***REMOVED***

		// X-Prefixed methods are assumed to be constructors; use a closure to wrap them in a
		// pure-JS function to allow them to be `new`d. (This is an awful hack...)
		if meth.Name[0] == 'X' ***REMOVED***
			wrapperV, _ := rt.RunProgram(constructWrap)
			wrapper, _ := goja.AssertFunction(wrapperV)
			v, _ := wrapper(goja.Undefined(), rt.ToValue(fn.Interface()))
			exports[name] = v
		***REMOVED*** else ***REMOVED***
			exports[name] = fn.Interface()
		***REMOVED***
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
