// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package defval marshals and unmarshals textual forms of default values.
//
// This package handles both the form historically used in Go struct field tags
// and also the form used by google.protobuf.FieldDescriptorProto.default_value
// since they differ in superficial ways.
package defval

import (
	"fmt"
	"math"
	"strconv"

	ptext "google.golang.org/protobuf/internal/encoding/text"
	errors "google.golang.org/protobuf/internal/errors"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// Format is the serialization format used to represent the default value.
type Format int

const (
	_ Format = iota

	// Descriptor uses the serialization format that protoc uses with the
	// google.protobuf.FieldDescriptorProto.default_value field.
	Descriptor

	// GoTag uses the historical serialization format in Go struct field tags.
	GoTag
)

// Unmarshal deserializes the default string s according to the given kind k.
// When k is an enum, a list of enum value descriptors must be provided.
func Unmarshal(s string, k pref.Kind, evs pref.EnumValueDescriptors, f Format) (pref.Value, pref.EnumValueDescriptor, error) ***REMOVED***
	switch k ***REMOVED***
	case pref.BoolKind:
		if f == GoTag ***REMOVED***
			switch s ***REMOVED***
			case "1":
				return pref.ValueOfBool(true), nil, nil
			case "0":
				return pref.ValueOfBool(false), nil, nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch s ***REMOVED***
			case "true":
				return pref.ValueOfBool(true), nil, nil
			case "false":
				return pref.ValueOfBool(false), nil, nil
			***REMOVED***
		***REMOVED***
	case pref.EnumKind:
		if f == GoTag ***REMOVED***
			// Go tags use the numeric form of the enum value.
			if n, err := strconv.ParseInt(s, 10, 32); err == nil ***REMOVED***
				if ev := evs.ByNumber(pref.EnumNumber(n)); ev != nil ***REMOVED***
					return pref.ValueOfEnum(ev.Number()), ev, nil
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Descriptor default_value use the enum identifier.
			ev := evs.ByName(pref.Name(s))
			if ev != nil ***REMOVED***
				return pref.ValueOfEnum(ev.Number()), ev, nil
			***REMOVED***
		***REMOVED***
	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		if v, err := strconv.ParseInt(s, 10, 32); err == nil ***REMOVED***
			return pref.ValueOfInt32(int32(v)), nil, nil
		***REMOVED***
	case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		if v, err := strconv.ParseInt(s, 10, 64); err == nil ***REMOVED***
			return pref.ValueOfInt64(int64(v)), nil, nil
		***REMOVED***
	case pref.Uint32Kind, pref.Fixed32Kind:
		if v, err := strconv.ParseUint(s, 10, 32); err == nil ***REMOVED***
			return pref.ValueOfUint32(uint32(v)), nil, nil
		***REMOVED***
	case pref.Uint64Kind, pref.Fixed64Kind:
		if v, err := strconv.ParseUint(s, 10, 64); err == nil ***REMOVED***
			return pref.ValueOfUint64(uint64(v)), nil, nil
		***REMOVED***
	case pref.FloatKind, pref.DoubleKind:
		var v float64
		var err error
		switch s ***REMOVED***
		case "-inf":
			v = math.Inf(-1)
		case "inf":
			v = math.Inf(+1)
		case "nan":
			v = math.NaN()
		default:
			v, err = strconv.ParseFloat(s, 64)
		***REMOVED***
		if err == nil ***REMOVED***
			if k == pref.FloatKind ***REMOVED***
				return pref.ValueOfFloat32(float32(v)), nil, nil
			***REMOVED*** else ***REMOVED***
				return pref.ValueOfFloat64(float64(v)), nil, nil
			***REMOVED***
		***REMOVED***
	case pref.StringKind:
		// String values are already unescaped and can be used as is.
		return pref.ValueOfString(s), nil, nil
	case pref.BytesKind:
		if b, ok := unmarshalBytes(s); ok ***REMOVED***
			return pref.ValueOfBytes(b), nil, nil
		***REMOVED***
	***REMOVED***
	return pref.Value***REMOVED******REMOVED***, nil, errors.New("could not parse value for %v: %q", k, s)
***REMOVED***

// Marshal serializes v as the default string according to the given kind k.
// When specifying the Descriptor format for an enum kind, the associated
// enum value descriptor must be provided.
func Marshal(v pref.Value, ev pref.EnumValueDescriptor, k pref.Kind, f Format) (string, error) ***REMOVED***
	switch k ***REMOVED***
	case pref.BoolKind:
		if f == GoTag ***REMOVED***
			if v.Bool() ***REMOVED***
				return "1", nil
			***REMOVED*** else ***REMOVED***
				return "0", nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if v.Bool() ***REMOVED***
				return "true", nil
			***REMOVED*** else ***REMOVED***
				return "false", nil
			***REMOVED***
		***REMOVED***
	case pref.EnumKind:
		if f == GoTag ***REMOVED***
			return strconv.FormatInt(int64(v.Enum()), 10), nil
		***REMOVED*** else ***REMOVED***
			return string(ev.Name()), nil
		***REMOVED***
	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind, pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		return strconv.FormatInt(v.Int(), 10), nil
	case pref.Uint32Kind, pref.Fixed32Kind, pref.Uint64Kind, pref.Fixed64Kind:
		return strconv.FormatUint(v.Uint(), 10), nil
	case pref.FloatKind, pref.DoubleKind:
		f := v.Float()
		switch ***REMOVED***
		case math.IsInf(f, -1):
			return "-inf", nil
		case math.IsInf(f, +1):
			return "inf", nil
		case math.IsNaN(f):
			return "nan", nil
		default:
			if k == pref.FloatKind ***REMOVED***
				return strconv.FormatFloat(f, 'g', -1, 32), nil
			***REMOVED*** else ***REMOVED***
				return strconv.FormatFloat(f, 'g', -1, 64), nil
			***REMOVED***
		***REMOVED***
	case pref.StringKind:
		// String values are serialized as is without any escaping.
		return v.String(), nil
	case pref.BytesKind:
		if s, ok := marshalBytes(v.Bytes()); ok ***REMOVED***
			return s, nil
		***REMOVED***
	***REMOVED***
	return "", errors.New("could not format value for %v: %v", k, v)
***REMOVED***

// unmarshalBytes deserializes bytes by applying C unescaping.
func unmarshalBytes(s string) ([]byte, bool) ***REMOVED***
	// Bytes values use the same escaping as the text format,
	// however they lack the surrounding double quotes.
	v, err := ptext.UnmarshalString(`"` + s + `"`)
	if err != nil ***REMOVED***
		return nil, false
	***REMOVED***
	return []byte(v), true
***REMOVED***

// marshalBytes serializes bytes by using C escaping.
// To match the exact output of protoc, this is identical to the
// CEscape function in strutil.cc of the protoc source code.
func marshalBytes(b []byte) (string, bool) ***REMOVED***
	var s []byte
	for _, c := range b ***REMOVED***
		switch c ***REMOVED***
		case '\n':
			s = append(s, `\n`...)
		case '\r':
			s = append(s, `\r`...)
		case '\t':
			s = append(s, `\t`...)
		case '"':
			s = append(s, `\"`...)
		case '\'':
			s = append(s, `\'`...)
		case '\\':
			s = append(s, `\\`...)
		default:
			if printableASCII := c >= 0x20 && c <= 0x7e; printableASCII ***REMOVED***
				s = append(s, c)
			***REMOVED*** else ***REMOVED***
				s = append(s, fmt.Sprintf(`\%03o`, c)...)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return string(s), true
***REMOVED***
