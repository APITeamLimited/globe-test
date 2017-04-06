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
	ctxT   = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorT = reflect.TypeOf((*error)(nil)).Elem()
	mapper = FieldNameMapper***REMOVED******REMOVED***
)

// The field name mapper translates Go symbol names for bridging to JS.
type FieldNameMapper struct***REMOVED******REMOVED***

// Bridge exported fields, snake_casing their names. A `js:"name"` tag overrides, `js:"-"` hides.
func (FieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string ***REMOVED***
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

// Bridge exported methods, but camelCase their names.
func (FieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string ***REMOVED***
	// PkgPath is non-empty for unexported methods.
	if m.PkgPath != "" ***REMOVED***
		return ""
	***REMOVED***

	// Lowercase the first character of the method name.
	return strings.ToLower(m.Name[0:1]) + m.Name[1:]
***REMOVED***

// Binds an object's members to the global scope. Returns a function that un-binds them.
// Note that this will panic if passed something that isn't a struct; please don't do that.
func BindToGlobal(rt *goja.Runtime, v interface***REMOVED******REMOVED***) func() ***REMOVED***
	keys := []string***REMOVED******REMOVED***

	val := reflect.ValueOf(v)
	typ := val.Type()
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		m := typ.Method(i)
		k := mapper.MethodName(typ, m)
		if k != "" ***REMOVED***
			fn := val.Method(i).Interface()
			keys = append(keys, k)
			rt.Set(k, fn)
		***REMOVED***
	***REMOVED***

	elem := val
	elemTyp := typ
	if typ.Kind() == reflect.Ptr ***REMOVED***
		elem = val.Elem()
		elemTyp = elem.Type()
	***REMOVED***
	for i := 0; i < elemTyp.NumField(); i++ ***REMOVED***
		f := elemTyp.Field(i)
		k := mapper.FieldName(elemTyp, f)
		if k != "" ***REMOVED***
			v := elem.Field(i).Interface()
			keys = append(keys, k)
			rt.Set(k, v)
		***REMOVED***
	***REMOVED***

	return func() ***REMOVED***
		for _, k := range keys ***REMOVED***
			rt.Set(k, goja.Undefined())
		***REMOVED***
	***REMOVED***
***REMOVED***

func Bind(rt *goja.Runtime, v interface***REMOVED******REMOVED***, ctx *context.Context) goja.Value ***REMOVED***
	val := reflect.ValueOf(v)
	typ := val.Type()
	exports := make(map[string]interface***REMOVED******REMOVED***)
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		methT := typ.Method(i)
		name := mapper.MethodName(typ, methT)
		if name == "" ***REMOVED***
			continue
		***REMOVED***
		meth := val.Method(i)

		in := make([]reflect.Type, methT.Type.NumIn())
		for i := 0; i < len(in); i++ ***REMOVED***
			in[i] = methT.Type.In(i)
		***REMOVED***
		out := make([]reflect.Type, methT.Type.NumOut())
		for i := 0; i < len(out); i++ ***REMOVED***
			out[i] = methT.Type.Out(i)
		***REMOVED***

		// Skip over the first input arg; it'll be the bound object.
		in = in[1:]

		// If the first argument is a context.Context, inject the given context.
		// The function will error if called outside of a valid context.
		if len(in) > 0 && in[0].Implements(ctxT) ***REMOVED***
			in = in[1:]
			meth = bindContext(in, out, methT, meth, rt, ctx)
		***REMOVED***

		// If the last return value is an error, turn it into a JS throw.
		if len(out) > 0 && out[len(out)-1] == errorT ***REMOVED***
			out = out[:len(out)-1]
			meth = bindErrorHandler(in, out, methT, meth, rt)
		***REMOVED***

		exports[name] = meth.Interface()
	***REMOVED***

	elem := val
	elemTyp := typ
	if typ.Kind() == reflect.Ptr ***REMOVED***
		elem = val.Elem()
		elemTyp = elem.Type()
	***REMOVED***
	for i := 0; i < elemTyp.NumField(); i++ ***REMOVED***
		f := elemTyp.Field(i)
		k := mapper.FieldName(elemTyp, f)
		if k == "" ***REMOVED***
			continue
		***REMOVED***
		exports[k] = elem.Field(i).Interface()
	***REMOVED***

	return rt.ToValue(exports)
***REMOVED***

func bindContext(in, out []reflect.Type, methT reflect.Method, meth reflect.Value, rt *goja.Runtime, ctxPtr *context.Context) reflect.Value ***REMOVED***
	return reflect.MakeFunc(
		reflect.FuncOf(in, out, methT.Type.IsVariadic()),
		func(args []reflect.Value) []reflect.Value ***REMOVED***
			if ctxPtr == nil || *ctxPtr == nil ***REMOVED***
				Throw(rt, errors.Errorf("%s needs a valid VU context", methT.Name))
			***REMOVED***
			ctx := *ctxPtr

			select ***REMOVED***
			case <-ctx.Done():
				Throw(rt, errors.Errorf("test has ended"))
			default:
			***REMOVED***

			return callBound(methT, meth, append([]reflect.Value***REMOVED***reflect.ValueOf(ctx)***REMOVED***, args...))
		***REMOVED***,
	)
***REMOVED***

func bindErrorHandler(in, out []reflect.Type, methT reflect.Method, meth reflect.Value, rt *goja.Runtime) reflect.Value ***REMOVED***
	return reflect.MakeFunc(
		reflect.FuncOf(in, out, methT.Type.IsVariadic()),
		func(args []reflect.Value) []reflect.Value ***REMOVED***
			ret := callBound(methT, meth, args)
			err := ret[len(ret)-1]
			if !err.IsNil() ***REMOVED***
				Throw(rt, err.Interface().(error))
			***REMOVED***
			return ret[:len(ret)-1]
		***REMOVED***,
	)
***REMOVED***

func callBound(methT reflect.Method, meth reflect.Value, args []reflect.Value) []reflect.Value ***REMOVED***
	if methT.Type.IsVariadic() ***REMOVED***
		return meth.CallSlice(args)
	***REMOVED***
	return meth.Call(args)
***REMOVED***
