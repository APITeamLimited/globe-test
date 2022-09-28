// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"bytes"
	"math"
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Equal reports whether two messages are equal.
// If two messages marshal to the same bytes under deterministic serialization,
// then Equal is guaranteed to report true.
//
// Two messages are equal if they belong to the same message descriptor,
// have the same set of populated known and extension field values,
// and the same set of unknown fields values. If either of the top-level
// messages are invalid, then Equal reports true only if both are invalid.
//
// Scalar values are compared with the equivalent of the == operator in Go,
// except bytes values which are compared using bytes.Equal and
// floating point values which specially treat NaNs as equal.
// Message values are compared by recursively calling Equal.
// Lists are equal if each element value is also equal.
// Maps are equal if they have the same set of keys, where the pair of values
// for each key is also equal.
func Equal(x, y Message) bool ***REMOVED***
	if x == nil || y == nil ***REMOVED***
		return x == nil && y == nil
	***REMOVED***
	if reflect.TypeOf(x).Kind() == reflect.Ptr && x == y ***REMOVED***
		// Avoid an expensive comparison if both inputs are identical pointers.
		return true
	***REMOVED***
	mx := x.ProtoReflect()
	my := y.ProtoReflect()
	if mx.IsValid() != my.IsValid() ***REMOVED***
		return false
	***REMOVED***
	return equalMessage(mx, my)
***REMOVED***

// equalMessage compares two messages.
func equalMessage(mx, my protoreflect.Message) bool ***REMOVED***
	if mx.Descriptor() != my.Descriptor() ***REMOVED***
		return false
	***REMOVED***

	nx := 0
	equal := true
	mx.Range(func(fd protoreflect.FieldDescriptor, vx protoreflect.Value) bool ***REMOVED***
		nx++
		vy := my.Get(fd)
		equal = my.Has(fd) && equalField(fd, vx, vy)
		return equal
	***REMOVED***)
	if !equal ***REMOVED***
		return false
	***REMOVED***
	ny := 0
	my.Range(func(fd protoreflect.FieldDescriptor, vx protoreflect.Value) bool ***REMOVED***
		ny++
		return true
	***REMOVED***)
	if nx != ny ***REMOVED***
		return false
	***REMOVED***

	return equalUnknown(mx.GetUnknown(), my.GetUnknown())
***REMOVED***

// equalField compares two fields.
func equalField(fd protoreflect.FieldDescriptor, x, y protoreflect.Value) bool ***REMOVED***
	switch ***REMOVED***
	case fd.IsList():
		return equalList(fd, x.List(), y.List())
	case fd.IsMap():
		return equalMap(fd, x.Map(), y.Map())
	default:
		return equalValue(fd, x, y)
	***REMOVED***
***REMOVED***

// equalMap compares two maps.
func equalMap(fd protoreflect.FieldDescriptor, x, y protoreflect.Map) bool ***REMOVED***
	if x.Len() != y.Len() ***REMOVED***
		return false
	***REMOVED***
	equal := true
	x.Range(func(k protoreflect.MapKey, vx protoreflect.Value) bool ***REMOVED***
		vy := y.Get(k)
		equal = y.Has(k) && equalValue(fd.MapValue(), vx, vy)
		return equal
	***REMOVED***)
	return equal
***REMOVED***

// equalList compares two lists.
func equalList(fd protoreflect.FieldDescriptor, x, y protoreflect.List) bool ***REMOVED***
	if x.Len() != y.Len() ***REMOVED***
		return false
	***REMOVED***
	for i := x.Len() - 1; i >= 0; i-- ***REMOVED***
		if !equalValue(fd, x.Get(i), y.Get(i)) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// equalValue compares two singular values.
func equalValue(fd protoreflect.FieldDescriptor, x, y protoreflect.Value) bool ***REMOVED***
	switch fd.Kind() ***REMOVED***
	case protoreflect.BoolKind:
		return x.Bool() == y.Bool()
	case protoreflect.EnumKind:
		return x.Enum() == y.Enum()
	case protoreflect.Int32Kind, protoreflect.Sint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		return x.Int() == y.Int()
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		return x.Uint() == y.Uint()
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		fx := x.Float()
		fy := y.Float()
		if math.IsNaN(fx) || math.IsNaN(fy) ***REMOVED***
			return math.IsNaN(fx) && math.IsNaN(fy)
		***REMOVED***
		return fx == fy
	case protoreflect.StringKind:
		return x.String() == y.String()
	case protoreflect.BytesKind:
		return bytes.Equal(x.Bytes(), y.Bytes())
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return equalMessage(x.Message(), y.Message())
	default:
		return x.Interface() == y.Interface()
	***REMOVED***
***REMOVED***

// equalUnknown compares unknown fields by direct comparison on the raw bytes
// of each individual field number.
func equalUnknown(x, y protoreflect.RawFields) bool ***REMOVED***
	if len(x) != len(y) ***REMOVED***
		return false
	***REMOVED***
	if bytes.Equal([]byte(x), []byte(y)) ***REMOVED***
		return true
	***REMOVED***

	mx := make(map[protoreflect.FieldNumber]protoreflect.RawFields)
	my := make(map[protoreflect.FieldNumber]protoreflect.RawFields)
	for len(x) > 0 ***REMOVED***
		fnum, _, n := protowire.ConsumeField(x)
		mx[fnum] = append(mx[fnum], x[:n]...)
		x = x[n:]
	***REMOVED***
	for len(y) > 0 ***REMOVED***
		fnum, _, n := protowire.ConsumeField(y)
		my[fnum] = append(my[fnum], y[:n]...)
		y = y[n:]
	***REMOVED***
	return reflect.DeepEqual(mx, my)
***REMOVED***
