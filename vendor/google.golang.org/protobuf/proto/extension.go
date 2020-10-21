// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// HasExtension reports whether an extension field is populated.
// It returns false if m is invalid or if xt does not extend m.
func HasExtension(m Message, xt protoreflect.ExtensionType) bool ***REMOVED***
	// Treat nil message interface as an empty message; no populated fields.
	if m == nil ***REMOVED***
		return false
	***REMOVED***

	// As a special-case, we reports invalid or mismatching descriptors
	// as always not being populated (since they aren't).
	if xt == nil || m.ProtoReflect().Descriptor() != xt.TypeDescriptor().ContainingMessage() ***REMOVED***
		return false
	***REMOVED***

	return m.ProtoReflect().Has(xt.TypeDescriptor())
***REMOVED***

// ClearExtension clears an extension field such that subsequent
// HasExtension calls return false.
// It panics if m is invalid or if xt does not extend m.
func ClearExtension(m Message, xt protoreflect.ExtensionType) ***REMOVED***
	m.ProtoReflect().Clear(xt.TypeDescriptor())
***REMOVED***

// GetExtension retrieves the value for an extension field.
// If the field is unpopulated, it returns the default value for
// scalars and an immutable, empty value for lists or messages.
// It panics if xt does not extend m.
func GetExtension(m Message, xt protoreflect.ExtensionType) interface***REMOVED******REMOVED*** ***REMOVED***
	// Treat nil message interface as an empty message; return the default.
	if m == nil ***REMOVED***
		return xt.InterfaceOf(xt.Zero())
	***REMOVED***

	return xt.InterfaceOf(m.ProtoReflect().Get(xt.TypeDescriptor()))
***REMOVED***

// SetExtension stores the value of an extension field.
// It panics if m is invalid, xt does not extend m, or if type of v
// is invalid for the specified extension field.
func SetExtension(m Message, xt protoreflect.ExtensionType, v interface***REMOVED******REMOVED***) ***REMOVED***
	xd := xt.TypeDescriptor()
	pv := xt.ValueOf(v)

	// Specially treat an invalid list, map, or message as clear.
	isValid := true
	switch ***REMOVED***
	case xd.IsList():
		isValid = pv.List().IsValid()
	case xd.IsMap():
		isValid = pv.Map().IsValid()
	case xd.Message() != nil:
		isValid = pv.Message().IsValid()
	***REMOVED***
	if !isValid ***REMOVED***
		m.ProtoReflect().Clear(xd)
		return
	***REMOVED***

	m.ProtoReflect().Set(xd, pv)
***REMOVED***

// RangeExtensions iterates over every populated extension field in m in an
// undefined order, calling f for each extension type and value encountered.
// It returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current extension field.
func RangeExtensions(m Message, f func(protoreflect.ExtensionType, interface***REMOVED******REMOVED***) bool) ***REMOVED***
	// Treat nil message interface as an empty message; nothing to range over.
	if m == nil ***REMOVED***
		return
	***REMOVED***

	m.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		if fd.IsExtension() ***REMOVED***
			xt := fd.(protoreflect.ExtensionTypeDescriptor).Type()
			vi := xt.InterfaceOf(v)
			return f(xt, vi)
		***REMOVED***
		return true
	***REMOVED***)
***REMOVED***
