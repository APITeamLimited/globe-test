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
func (FieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string ***REMOVED***
	return FieldName(t, f)
***REMOVED***

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
