// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// weakFields adds methods to the exported WeakFields type for internal use.
//
// The exported type is an alias to an unnamed type, so methods can't be
// defined directly on it.
type weakFields WeakFields

func (w weakFields) get(num protoreflect.FieldNumber) (protoreflect.ProtoMessage, bool) ***REMOVED***
	m, ok := w[int32(num)]
	return m, ok
***REMOVED***

func (w *weakFields) set(num protoreflect.FieldNumber, m protoreflect.ProtoMessage) ***REMOVED***
	if *w == nil ***REMOVED***
		*w = make(weakFields)
	***REMOVED***
	(*w)[int32(num)] = m
***REMOVED***

func (w *weakFields) clear(num protoreflect.FieldNumber) ***REMOVED***
	delete(*w, int32(num))
***REMOVED***

func (Export) HasWeak(w WeakFields, num protoreflect.FieldNumber) bool ***REMOVED***
	_, ok := w[int32(num)]
	return ok
***REMOVED***

func (Export) ClearWeak(w *WeakFields, num protoreflect.FieldNumber) ***REMOVED***
	delete(*w, int32(num))
***REMOVED***

func (Export) GetWeak(w WeakFields, num protoreflect.FieldNumber, name protoreflect.FullName) protoreflect.ProtoMessage ***REMOVED***
	if m, ok := w[int32(num)]; ok ***REMOVED***
		return m
	***REMOVED***
	mt, _ := protoregistry.GlobalTypes.FindMessageByName(name)
	if mt == nil ***REMOVED***
		panic(fmt.Sprintf("message %v for weak field is not linked in", name))
	***REMOVED***
	return mt.Zero().Interface()
***REMOVED***

func (Export) SetWeak(w *WeakFields, num protoreflect.FieldNumber, name protoreflect.FullName, m protoreflect.ProtoMessage) ***REMOVED***
	if m != nil ***REMOVED***
		mt, _ := protoregistry.GlobalTypes.FindMessageByName(name)
		if mt == nil ***REMOVED***
			panic(fmt.Sprintf("message %v for weak field is not linked in", name))
		***REMOVED***
		if mt != m.ProtoReflect().Type() ***REMOVED***
			panic(fmt.Sprintf("invalid message type for weak field: got %T, want %T", m, mt.Zero().Interface()))
		***REMOVED***
	***REMOVED***
	if m == nil || !m.ProtoReflect().IsValid() ***REMOVED***
		delete(*w, int32(num))
		return
	***REMOVED***
	if *w == nil ***REMOVED***
		*w = make(weakFields)
	***REMOVED***
	(*w)[int32(num)] = m
***REMOVED***
