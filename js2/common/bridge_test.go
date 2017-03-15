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

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

type bridgeTestType struct ***REMOVED***
	Exported      string
	ExportedTag   string `js:"renamed"`
	unexported    string
	unexportedTag string `js:"unexported"`
***REMOVED***

func (bridgeTestType) ExportedFn()   ***REMOVED******REMOVED***
func (bridgeTestType) unexportedFn() ***REMOVED******REMOVED***

func (*bridgeTestType) ExportedPtrFn()   ***REMOVED******REMOVED***
func (*bridgeTestType) unexportedPtrFn() ***REMOVED******REMOVED***

func TestFieldNameMapper(t *testing.T) ***REMOVED***
	typ := reflect.TypeOf(bridgeTestType***REMOVED******REMOVED***)
	t.Run("Fields", func(t *testing.T) ***REMOVED***
		names := map[string]string***REMOVED***
			"Exported":      "exported",
			"ExportedTag":   "renamed",
			"unexported":    "",
			"unexportedTag": "",
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
	t.Run("Exported", func(t *testing.T) ***REMOVED***
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
