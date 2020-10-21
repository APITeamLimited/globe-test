// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"reflect"

	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type EnumInfo struct ***REMOVED***
	GoReflectType reflect.Type // int32 kind
	Desc          pref.EnumDescriptor
***REMOVED***

func (t *EnumInfo) New(n pref.EnumNumber) pref.Enum ***REMOVED***
	return reflect.ValueOf(n).Convert(t.GoReflectType).Interface().(pref.Enum)
***REMOVED***
func (t *EnumInfo) Descriptor() pref.EnumDescriptor ***REMOVED*** return t.Desc ***REMOVED***
