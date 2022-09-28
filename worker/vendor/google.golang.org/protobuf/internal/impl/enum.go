// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type EnumInfo struct ***REMOVED***
	GoReflectType reflect.Type // int32 kind
	Desc          protoreflect.EnumDescriptor
***REMOVED***

func (t *EnumInfo) New(n protoreflect.EnumNumber) protoreflect.Enum ***REMOVED***
	return reflect.ValueOf(n).Convert(t.GoReflectType).Interface().(protoreflect.Enum)
***REMOVED***
func (t *EnumInfo) Descriptor() protoreflect.EnumDescriptor ***REMOVED*** return t.Desc ***REMOVED***
