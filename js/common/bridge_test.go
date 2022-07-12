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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type bridgeTestFieldsType struct ***REMOVED***
	Exported       string
	ExportedTag    string `js:"renamed"`
	ExportedHidden string `js:"-"`
	unexported     string //nolint:structcheck,unused // actually checked in the test
	unexportedTag  string `js:"unexported"` //nolint:structcheck,unused // actually checked in the test
***REMOVED***

type bridgeTestMethodsType struct***REMOVED******REMOVED***

func (bridgeTestMethodsType) ExportedFn() ***REMOVED******REMOVED***

//nolint:unused // needed for the actual test to check that it won't be seen
func (bridgeTestMethodsType) unexportedFn() ***REMOVED******REMOVED***

func (*bridgeTestMethodsType) ExportedPtrFn() ***REMOVED******REMOVED***

//nolint:unused // needed for the actual test to check that it won't be seen
func (*bridgeTestMethodsType) unexportedPtrFn() ***REMOVED******REMOVED***

type bridgeTestOddFieldsType struct ***REMOVED***
	TwoWords string
	URL      string
***REMOVED***

type bridgeTestConstructorType struct***REMOVED******REMOVED***

type bridgeTestConstructorSpawnedType struct***REMOVED******REMOVED***

func (bridgeTestConstructorType) XConstructor() bridgeTestConstructorSpawnedType ***REMOVED***
	return bridgeTestConstructorSpawnedType***REMOVED******REMOVED***
***REMOVED***

func TestFieldNameMapper(t *testing.T) ***REMOVED***
	t.Parallel()
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
		data := data
		for field, name := range data.Fields ***REMOVED***
			field, name := field, name
			t.Run(field, func(t *testing.T) ***REMOVED***
				t.Parallel()
				f, ok := data.Typ.FieldByName(field)
				if assert.True(t, ok, "no such field") ***REMOVED***
					assert.Equal(t, name, (FieldNameMapper***REMOVED******REMOVED***).FieldName(data.Typ, f))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
		for meth, name := range data.Methods ***REMOVED***
			meth, name := meth, name
			t.Run(meth, func(t *testing.T) ***REMOVED***
				t.Parallel()
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
