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

package js2

import (
	"github.com/dop251/goja"
	"reflect"
	"strings"
)

// The field name mapper translates Go symbol names for bridging to JS.
type FieldNameMapper struct***REMOVED******REMOVED***

// Bridge exported fields, camelCasing their names. A `js:"name"` tag overrides, `js:"-"` hides.
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
	return strings.ToLower(f.Name[0:1]) + f.Name[1:]
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
	mapper := FieldNameMapper***REMOVED******REMOVED***
	keys := []string***REMOVED******REMOVED***

	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr ***REMOVED***
		val = val.Elem()
		typ = val.Type()
	***REMOVED***
	for i := 0; i < typ.NumField(); i++ ***REMOVED***
		f := typ.Field(i)
		k := mapper.FieldName(typ, f)
		if k != "" ***REMOVED***
			v := val.Field(i).Interface()
			keys = append(keys, k)
			rt.Set(k, v)
		***REMOVED***
	***REMOVED***
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		m := typ.Method(i)
		k := mapper.MethodName(typ, m)
		if k != "" ***REMOVED***
			v := val.Method(i).Interface()
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
