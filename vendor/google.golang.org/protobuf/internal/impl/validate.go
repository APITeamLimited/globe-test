// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
)

// ValidationStatus is the result of validating the wire-format encoding of a message.
type ValidationStatus int

const (
	// ValidationUnknown indicates that unmarshaling the message might succeed or fail.
	// The validator was unable to render a judgement.
	//
	// The only causes of this status are an aberrant message type appearing somewhere
	// in the message or a failure in the extension resolver.
	ValidationUnknown ValidationStatus = iota + 1

	// ValidationInvalid indicates that unmarshaling the message will fail.
	ValidationInvalid

	// ValidationValid indicates that unmarshaling the message will succeed.
	ValidationValid
)

func (v ValidationStatus) String() string ***REMOVED***
	switch v ***REMOVED***
	case ValidationUnknown:
		return "ValidationUnknown"
	case ValidationInvalid:
		return "ValidationInvalid"
	case ValidationValid:
		return "ValidationValid"
	default:
		return fmt.Sprintf("ValidationStatus(%d)", int(v))
	***REMOVED***
***REMOVED***

// Validate determines whether the contents of the buffer are a valid wire encoding
// of the message type.
//
// This function is exposed for testing.
func Validate(mt protoreflect.MessageType, in protoiface.UnmarshalInput) (out protoiface.UnmarshalOutput, _ ValidationStatus) ***REMOVED***
	mi, ok := mt.(*MessageInfo)
	if !ok ***REMOVED***
		return out, ValidationUnknown
	***REMOVED***
	if in.Resolver == nil ***REMOVED***
		in.Resolver = protoregistry.GlobalTypes
	***REMOVED***
	o, st := mi.validate(in.Buf, 0, unmarshalOptions***REMOVED***
		flags:    in.Flags,
		resolver: in.Resolver,
	***REMOVED***)
	if o.initialized ***REMOVED***
		out.Flags |= protoiface.UnmarshalInitialized
	***REMOVED***
	return out, st
***REMOVED***

type validationInfo struct ***REMOVED***
	mi               *MessageInfo
	typ              validationType
	keyType, valType validationType

	// For non-required fields, requiredBit is 0.
	//
	// For required fields, requiredBit's nth bit is set, where n is a
	// unique index in the range [0, MessageInfo.numRequiredFields).
	//
	// If there are more than 64 required fields, requiredBit is 0.
	requiredBit uint64
***REMOVED***

type validationType uint8

const (
	validationTypeOther validationType = iota
	validationTypeMessage
	validationTypeGroup
	validationTypeMap
	validationTypeRepeatedVarint
	validationTypeRepeatedFixed32
	validationTypeRepeatedFixed64
	validationTypeVarint
	validationTypeFixed32
	validationTypeFixed64
	validationTypeBytes
	validationTypeUTF8String
	validationTypeMessageSetItem
)

func newFieldValidationInfo(mi *MessageInfo, si structInfo, fd protoreflect.FieldDescriptor, ft reflect.Type) validationInfo ***REMOVED***
	var vi validationInfo
	switch ***REMOVED***
	case fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic():
		switch fd.Kind() ***REMOVED***
		case protoreflect.MessageKind:
			vi.typ = validationTypeMessage
			if ot, ok := si.oneofWrappersByNumber[fd.Number()]; ok ***REMOVED***
				vi.mi = getMessageInfo(ot.Field(0).Type)
			***REMOVED***
		case protoreflect.GroupKind:
			vi.typ = validationTypeGroup
			if ot, ok := si.oneofWrappersByNumber[fd.Number()]; ok ***REMOVED***
				vi.mi = getMessageInfo(ot.Field(0).Type)
			***REMOVED***
		case protoreflect.StringKind:
			if strs.EnforceUTF8(fd) ***REMOVED***
				vi.typ = validationTypeUTF8String
			***REMOVED***
		***REMOVED***
	default:
		vi = newValidationInfo(fd, ft)
	***REMOVED***
	if fd.Cardinality() == protoreflect.Required ***REMOVED***
		// Avoid overflow. The required field check is done with a 64-bit mask, with
		// any message containing more than 64 required fields always reported as
		// potentially uninitialized, so it is not important to get a precise count
		// of the required fields past 64.
		if mi.numRequiredFields < math.MaxUint8 ***REMOVED***
			mi.numRequiredFields++
			vi.requiredBit = 1 << (mi.numRequiredFields - 1)
		***REMOVED***
	***REMOVED***
	return vi
***REMOVED***

func newValidationInfo(fd protoreflect.FieldDescriptor, ft reflect.Type) validationInfo ***REMOVED***
	var vi validationInfo
	switch ***REMOVED***
	case fd.IsList():
		switch fd.Kind() ***REMOVED***
		case protoreflect.MessageKind:
			vi.typ = validationTypeMessage
			if ft.Kind() == reflect.Slice ***REMOVED***
				vi.mi = getMessageInfo(ft.Elem())
			***REMOVED***
		case protoreflect.GroupKind:
			vi.typ = validationTypeGroup
			if ft.Kind() == reflect.Slice ***REMOVED***
				vi.mi = getMessageInfo(ft.Elem())
			***REMOVED***
		case protoreflect.StringKind:
			vi.typ = validationTypeBytes
			if strs.EnforceUTF8(fd) ***REMOVED***
				vi.typ = validationTypeUTF8String
			***REMOVED***
		default:
			switch wireTypes[fd.Kind()] ***REMOVED***
			case protowire.VarintType:
				vi.typ = validationTypeRepeatedVarint
			case protowire.Fixed32Type:
				vi.typ = validationTypeRepeatedFixed32
			case protowire.Fixed64Type:
				vi.typ = validationTypeRepeatedFixed64
			***REMOVED***
		***REMOVED***
	case fd.IsMap():
		vi.typ = validationTypeMap
		switch fd.MapKey().Kind() ***REMOVED***
		case protoreflect.StringKind:
			if strs.EnforceUTF8(fd) ***REMOVED***
				vi.keyType = validationTypeUTF8String
			***REMOVED***
		***REMOVED***
		switch fd.MapValue().Kind() ***REMOVED***
		case protoreflect.MessageKind:
			vi.valType = validationTypeMessage
			if ft.Kind() == reflect.Map ***REMOVED***
				vi.mi = getMessageInfo(ft.Elem())
			***REMOVED***
		case protoreflect.StringKind:
			if strs.EnforceUTF8(fd) ***REMOVED***
				vi.valType = validationTypeUTF8String
			***REMOVED***
		***REMOVED***
	default:
		switch fd.Kind() ***REMOVED***
		case protoreflect.MessageKind:
			vi.typ = validationTypeMessage
			if !fd.IsWeak() ***REMOVED***
				vi.mi = getMessageInfo(ft)
			***REMOVED***
		case protoreflect.GroupKind:
			vi.typ = validationTypeGroup
			vi.mi = getMessageInfo(ft)
		case protoreflect.StringKind:
			vi.typ = validationTypeBytes
			if strs.EnforceUTF8(fd) ***REMOVED***
				vi.typ = validationTypeUTF8String
			***REMOVED***
		default:
			switch wireTypes[fd.Kind()] ***REMOVED***
			case protowire.VarintType:
				vi.typ = validationTypeVarint
			case protowire.Fixed32Type:
				vi.typ = validationTypeFixed32
			case protowire.Fixed64Type:
				vi.typ = validationTypeFixed64
			case protowire.BytesType:
				vi.typ = validationTypeBytes
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return vi
***REMOVED***

func (mi *MessageInfo) validate(b []byte, groupTag protowire.Number, opts unmarshalOptions) (out unmarshalOutput, result ValidationStatus) ***REMOVED***
	mi.init()
	type validationState struct ***REMOVED***
		typ              validationType
		keyType, valType validationType
		endGroup         protowire.Number
		mi               *MessageInfo
		tail             []byte
		requiredMask     uint64
	***REMOVED***

	// Pre-allocate some slots to avoid repeated slice reallocation.
	states := make([]validationState, 0, 16)
	states = append(states, validationState***REMOVED***
		typ: validationTypeMessage,
		mi:  mi,
	***REMOVED***)
	if groupTag > 0 ***REMOVED***
		states[0].typ = validationTypeGroup
		states[0].endGroup = groupTag
	***REMOVED***
	initialized := true
	start := len(b)
State:
	for len(states) > 0 ***REMOVED***
		st := &states[len(states)-1]
		for len(b) > 0 ***REMOVED***
			// Parse the tag (field number and wire type).
			var tag uint64
			if b[0] < 0x80 ***REMOVED***
				tag = uint64(b[0])
				b = b[1:]
			***REMOVED*** else if len(b) >= 2 && b[1] < 128 ***REMOVED***
				tag = uint64(b[0]&0x7f) + uint64(b[1])<<7
				b = b[2:]
			***REMOVED*** else ***REMOVED***
				var n int
				tag, n = protowire.ConsumeVarint(b)
				if n < 0 ***REMOVED***
					return out, ValidationInvalid
				***REMOVED***
				b = b[n:]
			***REMOVED***
			var num protowire.Number
			if n := tag >> 3; n < uint64(protowire.MinValidNumber) || n > uint64(protowire.MaxValidNumber) ***REMOVED***
				return out, ValidationInvalid
			***REMOVED*** else ***REMOVED***
				num = protowire.Number(n)
			***REMOVED***
			wtyp := protowire.Type(tag & 7)

			if wtyp == protowire.EndGroupType ***REMOVED***
				if st.endGroup == num ***REMOVED***
					goto PopState
				***REMOVED***
				return out, ValidationInvalid
			***REMOVED***
			var vi validationInfo
			switch ***REMOVED***
			case st.typ == validationTypeMap:
				switch num ***REMOVED***
				case genid.MapEntry_Key_field_number:
					vi.typ = st.keyType
				case genid.MapEntry_Value_field_number:
					vi.typ = st.valType
					vi.mi = st.mi
					vi.requiredBit = 1
				***REMOVED***
			case flags.ProtoLegacy && st.mi.isMessageSet:
				switch num ***REMOVED***
				case messageset.FieldItem:
					vi.typ = validationTypeMessageSetItem
				***REMOVED***
			default:
				var f *coderFieldInfo
				if int(num) < len(st.mi.denseCoderFields) ***REMOVED***
					f = st.mi.denseCoderFields[num]
				***REMOVED*** else ***REMOVED***
					f = st.mi.coderFields[num]
				***REMOVED***
				if f != nil ***REMOVED***
					vi = f.validation
					if vi.typ == validationTypeMessage && vi.mi == nil ***REMOVED***
						// Probable weak field.
						//
						// TODO: Consider storing the results of this lookup somewhere
						// rather than recomputing it on every validation.
						fd := st.mi.Desc.Fields().ByNumber(num)
						if fd == nil || !fd.IsWeak() ***REMOVED***
							break
						***REMOVED***
						messageName := fd.Message().FullName()
						messageType, err := protoregistry.GlobalTypes.FindMessageByName(messageName)
						switch err ***REMOVED***
						case nil:
							vi.mi, _ = messageType.(*MessageInfo)
						case protoregistry.NotFound:
							vi.typ = validationTypeBytes
						default:
							return out, ValidationUnknown
						***REMOVED***
					***REMOVED***
					break
				***REMOVED***
				// Possible extension field.
				//
				// TODO: We should return ValidationUnknown when:
				//   1. The resolver is not frozen. (More extensions may be added to it.)
				//   2. The resolver returns preg.NotFound.
				// In this case, a type added to the resolver in the future could cause
				// unmarshaling to begin failing. Supporting this requires some way to
				// determine if the resolver is frozen.
				xt, err := opts.resolver.FindExtensionByNumber(st.mi.Desc.FullName(), num)
				if err != nil && err != protoregistry.NotFound ***REMOVED***
					return out, ValidationUnknown
				***REMOVED***
				if err == nil ***REMOVED***
					vi = getExtensionFieldInfo(xt).validation
				***REMOVED***
			***REMOVED***
			if vi.requiredBit != 0 ***REMOVED***
				// Check that the field has a compatible wire type.
				// We only need to consider non-repeated field types,
				// since repeated fields (and maps) can never be required.
				ok := false
				switch vi.typ ***REMOVED***
				case validationTypeVarint:
					ok = wtyp == protowire.VarintType
				case validationTypeFixed32:
					ok = wtyp == protowire.Fixed32Type
				case validationTypeFixed64:
					ok = wtyp == protowire.Fixed64Type
				case validationTypeBytes, validationTypeUTF8String, validationTypeMessage:
					ok = wtyp == protowire.BytesType
				case validationTypeGroup:
					ok = wtyp == protowire.StartGroupType
				***REMOVED***
				if ok ***REMOVED***
					st.requiredMask |= vi.requiredBit
				***REMOVED***
			***REMOVED***

			switch wtyp ***REMOVED***
			case protowire.VarintType:
				if len(b) >= 10 ***REMOVED***
					switch ***REMOVED***
					case b[0] < 0x80:
						b = b[1:]
					case b[1] < 0x80:
						b = b[2:]
					case b[2] < 0x80:
						b = b[3:]
					case b[3] < 0x80:
						b = b[4:]
					case b[4] < 0x80:
						b = b[5:]
					case b[5] < 0x80:
						b = b[6:]
					case b[6] < 0x80:
						b = b[7:]
					case b[7] < 0x80:
						b = b[8:]
					case b[8] < 0x80:
						b = b[9:]
					case b[9] < 0x80 && b[9] < 2:
						b = b[10:]
					default:
						return out, ValidationInvalid
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					switch ***REMOVED***
					case len(b) > 0 && b[0] < 0x80:
						b = b[1:]
					case len(b) > 1 && b[1] < 0x80:
						b = b[2:]
					case len(b) > 2 && b[2] < 0x80:
						b = b[3:]
					case len(b) > 3 && b[3] < 0x80:
						b = b[4:]
					case len(b) > 4 && b[4] < 0x80:
						b = b[5:]
					case len(b) > 5 && b[5] < 0x80:
						b = b[6:]
					case len(b) > 6 && b[6] < 0x80:
						b = b[7:]
					case len(b) > 7 && b[7] < 0x80:
						b = b[8:]
					case len(b) > 8 && b[8] < 0x80:
						b = b[9:]
					case len(b) > 9 && b[9] < 2:
						b = b[10:]
					default:
						return out, ValidationInvalid
					***REMOVED***
				***REMOVED***
				continue State
			case protowire.BytesType:
				var size uint64
				if len(b) >= 1 && b[0] < 0x80 ***REMOVED***
					size = uint64(b[0])
					b = b[1:]
				***REMOVED*** else if len(b) >= 2 && b[1] < 128 ***REMOVED***
					size = uint64(b[0]&0x7f) + uint64(b[1])<<7
					b = b[2:]
				***REMOVED*** else ***REMOVED***
					var n int
					size, n = protowire.ConsumeVarint(b)
					if n < 0 ***REMOVED***
						return out, ValidationInvalid
					***REMOVED***
					b = b[n:]
				***REMOVED***
				if size > uint64(len(b)) ***REMOVED***
					return out, ValidationInvalid
				***REMOVED***
				v := b[:size]
				b = b[size:]
				switch vi.typ ***REMOVED***
				case validationTypeMessage:
					if vi.mi == nil ***REMOVED***
						return out, ValidationUnknown
					***REMOVED***
					vi.mi.init()
					fallthrough
				case validationTypeMap:
					if vi.mi != nil ***REMOVED***
						vi.mi.init()
					***REMOVED***
					states = append(states, validationState***REMOVED***
						typ:     vi.typ,
						keyType: vi.keyType,
						valType: vi.valType,
						mi:      vi.mi,
						tail:    b,
					***REMOVED***)
					b = v
					continue State
				case validationTypeRepeatedVarint:
					// Packed field.
					for len(v) > 0 ***REMOVED***
						_, n := protowire.ConsumeVarint(v)
						if n < 0 ***REMOVED***
							return out, ValidationInvalid
						***REMOVED***
						v = v[n:]
					***REMOVED***
				case validationTypeRepeatedFixed32:
					// Packed field.
					if len(v)%4 != 0 ***REMOVED***
						return out, ValidationInvalid
					***REMOVED***
				case validationTypeRepeatedFixed64:
					// Packed field.
					if len(v)%8 != 0 ***REMOVED***
						return out, ValidationInvalid
					***REMOVED***
				case validationTypeUTF8String:
					if !utf8.Valid(v) ***REMOVED***
						return out, ValidationInvalid
					***REMOVED***
				***REMOVED***
			case protowire.Fixed32Type:
				if len(b) < 4 ***REMOVED***
					return out, ValidationInvalid
				***REMOVED***
				b = b[4:]
			case protowire.Fixed64Type:
				if len(b) < 8 ***REMOVED***
					return out, ValidationInvalid
				***REMOVED***
				b = b[8:]
			case protowire.StartGroupType:
				switch ***REMOVED***
				case vi.typ == validationTypeGroup:
					if vi.mi == nil ***REMOVED***
						return out, ValidationUnknown
					***REMOVED***
					vi.mi.init()
					states = append(states, validationState***REMOVED***
						typ:      validationTypeGroup,
						mi:       vi.mi,
						endGroup: num,
					***REMOVED***)
					continue State
				case flags.ProtoLegacy && vi.typ == validationTypeMessageSetItem:
					typeid, v, n, err := messageset.ConsumeFieldValue(b, false)
					if err != nil ***REMOVED***
						return out, ValidationInvalid
					***REMOVED***
					xt, err := opts.resolver.FindExtensionByNumber(st.mi.Desc.FullName(), typeid)
					switch ***REMOVED***
					case err == protoregistry.NotFound:
						b = b[n:]
					case err != nil:
						return out, ValidationUnknown
					default:
						xvi := getExtensionFieldInfo(xt).validation
						if xvi.mi != nil ***REMOVED***
							xvi.mi.init()
						***REMOVED***
						states = append(states, validationState***REMOVED***
							typ:  xvi.typ,
							mi:   xvi.mi,
							tail: b[n:],
						***REMOVED***)
						b = v
						continue State
					***REMOVED***
				default:
					n := protowire.ConsumeFieldValue(num, wtyp, b)
					if n < 0 ***REMOVED***
						return out, ValidationInvalid
					***REMOVED***
					b = b[n:]
				***REMOVED***
			default:
				return out, ValidationInvalid
			***REMOVED***
		***REMOVED***
		if st.endGroup != 0 ***REMOVED***
			return out, ValidationInvalid
		***REMOVED***
		if len(b) != 0 ***REMOVED***
			return out, ValidationInvalid
		***REMOVED***
		b = st.tail
	PopState:
		numRequiredFields := 0
		switch st.typ ***REMOVED***
		case validationTypeMessage, validationTypeGroup:
			numRequiredFields = int(st.mi.numRequiredFields)
		case validationTypeMap:
			// If this is a map field with a message value that contains
			// required fields, require that the value be present.
			if st.mi != nil && st.mi.numRequiredFields > 0 ***REMOVED***
				numRequiredFields = 1
			***REMOVED***
		***REMOVED***
		// If there are more than 64 required fields, this check will
		// always fail and we will report that the message is potentially
		// uninitialized.
		if numRequiredFields > 0 && bits.OnesCount64(st.requiredMask) != numRequiredFields ***REMOVED***
			initialized = false
		***REMOVED***
		states = states[:len(states)-1]
	***REMOVED***
	out.n = start - len(b)
	if initialized ***REMOVED***
		out.initialized = true
	***REMOVED***
	return out, ValidationValid
***REMOVED***
