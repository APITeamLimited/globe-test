// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonpb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const wrapJSONUnmarshalV2 = false

// UnmarshalNext unmarshals the next JSON object from d into m.
func UnmarshalNext(d *json.Decoder, m proto.Message) error ***REMOVED***
	return new(Unmarshaler).UnmarshalNext(d, m)
***REMOVED***

// Unmarshal unmarshals a JSON object from r into m.
func Unmarshal(r io.Reader, m proto.Message) error ***REMOVED***
	return new(Unmarshaler).Unmarshal(r, m)
***REMOVED***

// UnmarshalString unmarshals a JSON object from s into m.
func UnmarshalString(s string, m proto.Message) error ***REMOVED***
	return new(Unmarshaler).Unmarshal(strings.NewReader(s), m)
***REMOVED***

// Unmarshaler is a configurable object for converting from a JSON
// representation to a protocol buffer object.
type Unmarshaler struct ***REMOVED***
	// AllowUnknownFields specifies whether to allow messages to contain
	// unknown JSON fields, as opposed to failing to unmarshal.
	AllowUnknownFields bool

	// AnyResolver is used to resolve the google.protobuf.Any well-known type.
	// If unset, the global registry is used by default.
	AnyResolver AnyResolver
***REMOVED***

// JSONPBUnmarshaler is implemented by protobuf messages that customize the way
// they are unmarshaled from JSON. Messages that implement this should also
// implement JSONPBMarshaler so that the custom format can be produced.
//
// The JSON unmarshaling must follow the JSON to proto specification:
//	https://developers.google.com/protocol-buffers/docs/proto3#json
//
// Deprecated: Custom types should implement protobuf reflection instead.
type JSONPBUnmarshaler interface ***REMOVED***
	UnmarshalJSONPB(*Unmarshaler, []byte) error
***REMOVED***

// Unmarshal unmarshals a JSON object from r into m.
func (u *Unmarshaler) Unmarshal(r io.Reader, m proto.Message) error ***REMOVED***
	return u.UnmarshalNext(json.NewDecoder(r), m)
***REMOVED***

// UnmarshalNext unmarshals the next JSON object from d into m.
func (u *Unmarshaler) UnmarshalNext(d *json.Decoder, m proto.Message) error ***REMOVED***
	if m == nil ***REMOVED***
		return errors.New("invalid nil message")
	***REMOVED***

	// Parse the next JSON object from the stream.
	raw := json.RawMessage***REMOVED******REMOVED***
	if err := d.Decode(&raw); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check for custom unmarshalers first since they may not properly
	// implement protobuf reflection that the logic below relies on.
	if jsu, ok := m.(JSONPBUnmarshaler); ok ***REMOVED***
		return jsu.UnmarshalJSONPB(u, raw)
	***REMOVED***

	mr := proto.MessageReflect(m)

	// NOTE: For historical reasons, a top-level null is treated as a noop.
	// This is incorrect, but kept for compatibility.
	if string(raw) == "null" && mr.Descriptor().FullName() != "google.protobuf.Value" ***REMOVED***
		return nil
	***REMOVED***

	if wrapJSONUnmarshalV2 ***REMOVED***
		// NOTE: If input message is non-empty, we need to preserve merge semantics
		// of the old jsonpb implementation. These semantics are not supported by
		// the protobuf JSON specification.
		isEmpty := true
		mr.Range(func(protoreflect.FieldDescriptor, protoreflect.Value) bool ***REMOVED***
			isEmpty = false // at least one iteration implies non-empty
			return false
		***REMOVED***)
		if !isEmpty ***REMOVED***
			// Perform unmarshaling into a newly allocated, empty message.
			mr = mr.New()

			// Use a defer to copy all unmarshaled fields into the original message.
			dst := proto.MessageReflect(m)
			defer mr.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
				dst.Set(fd, v)
				return true
			***REMOVED***)
		***REMOVED***

		// Unmarshal using the v2 JSON unmarshaler.
		opts := protojson.UnmarshalOptions***REMOVED***
			DiscardUnknown: u.AllowUnknownFields,
		***REMOVED***
		if u.AnyResolver != nil ***REMOVED***
			opts.Resolver = anyResolver***REMOVED***u.AnyResolver***REMOVED***
		***REMOVED***
		return opts.Unmarshal(raw, mr.Interface())
	***REMOVED*** else ***REMOVED***
		if err := u.unmarshalMessage(mr, raw); err != nil ***REMOVED***
			return err
		***REMOVED***
		return protoV2.CheckInitialized(mr.Interface())
	***REMOVED***
***REMOVED***

func (u *Unmarshaler) unmarshalMessage(m protoreflect.Message, in []byte) error ***REMOVED***
	md := m.Descriptor()
	fds := md.Fields()

	if jsu, ok := proto.MessageV1(m.Interface()).(JSONPBUnmarshaler); ok ***REMOVED***
		return jsu.UnmarshalJSONPB(u, in)
	***REMOVED***

	if string(in) == "null" && md.FullName() != "google.protobuf.Value" ***REMOVED***
		return nil
	***REMOVED***

	switch wellKnownType(md.FullName()) ***REMOVED***
	case "Any":
		var jsonObject map[string]json.RawMessage
		if err := json.Unmarshal(in, &jsonObject); err != nil ***REMOVED***
			return err
		***REMOVED***

		rawTypeURL, ok := jsonObject["@type"]
		if !ok ***REMOVED***
			return errors.New("Any JSON doesn't have '@type'")
		***REMOVED***
		typeURL, err := unquoteString(string(rawTypeURL))
		if err != nil ***REMOVED***
			return fmt.Errorf("can't unmarshal Any's '@type': %q", rawTypeURL)
		***REMOVED***
		m.Set(fds.ByNumber(1), protoreflect.ValueOfString(typeURL))

		var m2 protoreflect.Message
		if u.AnyResolver != nil ***REMOVED***
			mi, err := u.AnyResolver.Resolve(typeURL)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			m2 = proto.MessageReflect(mi)
		***REMOVED*** else ***REMOVED***
			mt, err := protoregistry.GlobalTypes.FindMessageByURL(typeURL)
			if err != nil ***REMOVED***
				if err == protoregistry.NotFound ***REMOVED***
					return fmt.Errorf("could not resolve Any message type: %v", typeURL)
				***REMOVED***
				return err
			***REMOVED***
			m2 = mt.New()
		***REMOVED***

		if wellKnownType(m2.Descriptor().FullName()) != "" ***REMOVED***
			rawValue, ok := jsonObject["value"]
			if !ok ***REMOVED***
				return errors.New("Any JSON doesn't have 'value'")
			***REMOVED***
			if err := u.unmarshalMessage(m2, rawValue); err != nil ***REMOVED***
				return fmt.Errorf("can't unmarshal Any nested proto %v: %v", typeURL, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			delete(jsonObject, "@type")
			rawJSON, err := json.Marshal(jsonObject)
			if err != nil ***REMOVED***
				return fmt.Errorf("can't generate JSON for Any's nested proto to be unmarshaled: %v", err)
			***REMOVED***
			if err = u.unmarshalMessage(m2, rawJSON); err != nil ***REMOVED***
				return fmt.Errorf("can't unmarshal Any nested proto %v: %v", typeURL, err)
			***REMOVED***
		***REMOVED***

		rawWire, err := protoV2.Marshal(m2.Interface())
		if err != nil ***REMOVED***
			return fmt.Errorf("can't marshal proto %v into Any.Value: %v", typeURL, err)
		***REMOVED***
		m.Set(fds.ByNumber(2), protoreflect.ValueOfBytes(rawWire))
		return nil
	case "BoolValue", "BytesValue", "StringValue",
		"Int32Value", "UInt32Value", "FloatValue",
		"Int64Value", "UInt64Value", "DoubleValue":
		fd := fds.ByNumber(1)
		v, err := u.unmarshalValue(m.NewField(fd), in, fd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m.Set(fd, v)
		return nil
	case "Duration":
		v, err := unquoteString(string(in))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		d, err := time.ParseDuration(v)
		if err != nil ***REMOVED***
			return fmt.Errorf("bad Duration: %v", err)
		***REMOVED***

		sec := d.Nanoseconds() / 1e9
		nsec := d.Nanoseconds() % 1e9
		m.Set(fds.ByNumber(1), protoreflect.ValueOfInt64(int64(sec)))
		m.Set(fds.ByNumber(2), protoreflect.ValueOfInt32(int32(nsec)))
		return nil
	case "Timestamp":
		v, err := unquoteString(string(in))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		t, err := time.Parse(time.RFC3339Nano, v)
		if err != nil ***REMOVED***
			return fmt.Errorf("bad Timestamp: %v", err)
		***REMOVED***

		sec := t.Unix()
		nsec := t.Nanosecond()
		m.Set(fds.ByNumber(1), protoreflect.ValueOfInt64(int64(sec)))
		m.Set(fds.ByNumber(2), protoreflect.ValueOfInt32(int32(nsec)))
		return nil
	case "Value":
		switch ***REMOVED***
		case string(in) == "null":
			m.Set(fds.ByNumber(1), protoreflect.ValueOfEnum(0))
		case string(in) == "true":
			m.Set(fds.ByNumber(4), protoreflect.ValueOfBool(true))
		case string(in) == "false":
			m.Set(fds.ByNumber(4), protoreflect.ValueOfBool(false))
		case hasPrefixAndSuffix('"', in, '"'):
			s, err := unquoteString(string(in))
			if err != nil ***REMOVED***
				return fmt.Errorf("unrecognized type for Value %q", in)
			***REMOVED***
			m.Set(fds.ByNumber(3), protoreflect.ValueOfString(s))
		case hasPrefixAndSuffix('[', in, ']'):
			v := m.Mutable(fds.ByNumber(6))
			return u.unmarshalMessage(v.Message(), in)
		case hasPrefixAndSuffix('***REMOVED***', in, '***REMOVED***'):
			v := m.Mutable(fds.ByNumber(5))
			return u.unmarshalMessage(v.Message(), in)
		default:
			f, err := strconv.ParseFloat(string(in), 0)
			if err != nil ***REMOVED***
				return fmt.Errorf("unrecognized type for Value %q", in)
			***REMOVED***
			m.Set(fds.ByNumber(2), protoreflect.ValueOfFloat64(f))
		***REMOVED***
		return nil
	case "ListValue":
		var jsonArray []json.RawMessage
		if err := json.Unmarshal(in, &jsonArray); err != nil ***REMOVED***
			return fmt.Errorf("bad ListValue: %v", err)
		***REMOVED***

		lv := m.Mutable(fds.ByNumber(1)).List()
		for _, raw := range jsonArray ***REMOVED***
			ve := lv.NewElement()
			if err := u.unmarshalMessage(ve.Message(), raw); err != nil ***REMOVED***
				return err
			***REMOVED***
			lv.Append(ve)
		***REMOVED***
		return nil
	case "Struct":
		var jsonObject map[string]json.RawMessage
		if err := json.Unmarshal(in, &jsonObject); err != nil ***REMOVED***
			return fmt.Errorf("bad StructValue: %v", err)
		***REMOVED***

		mv := m.Mutable(fds.ByNumber(1)).Map()
		for key, raw := range jsonObject ***REMOVED***
			kv := protoreflect.ValueOf(key).MapKey()
			vv := mv.NewValue()
			if err := u.unmarshalMessage(vv.Message(), raw); err != nil ***REMOVED***
				return fmt.Errorf("bad value in StructValue for key %q: %v", key, err)
			***REMOVED***
			mv.Set(kv, vv)
		***REMOVED***
		return nil
	***REMOVED***

	var jsonObject map[string]json.RawMessage
	if err := json.Unmarshal(in, &jsonObject); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Handle known fields.
	for i := 0; i < fds.Len(); i++ ***REMOVED***
		fd := fds.Get(i)
		if fd.IsWeak() && fd.Message().IsPlaceholder() ***REMOVED***
			continue //  weak reference is not linked in
		***REMOVED***

		// Search for any raw JSON value associated with this field.
		var raw json.RawMessage
		name := string(fd.Name())
		if fd.Kind() == protoreflect.GroupKind ***REMOVED***
			name = string(fd.Message().Name())
		***REMOVED***
		if v, ok := jsonObject[name]; ok ***REMOVED***
			delete(jsonObject, name)
			raw = v
		***REMOVED***
		name = string(fd.JSONName())
		if v, ok := jsonObject[name]; ok ***REMOVED***
			delete(jsonObject, name)
			raw = v
		***REMOVED***

		field := m.NewField(fd)
		// Unmarshal the field value.
		if raw == nil || (string(raw) == "null" && !isSingularWellKnownValue(fd) && !isSingularJSONPBUnmarshaler(field, fd)) ***REMOVED***
			continue
		***REMOVED***
		v, err := u.unmarshalValue(field, raw, fd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m.Set(fd, v)
	***REMOVED***

	// Handle extension fields.
	for name, raw := range jsonObject ***REMOVED***
		if !strings.HasPrefix(name, "[") || !strings.HasSuffix(name, "]") ***REMOVED***
			continue
		***REMOVED***

		// Resolve the extension field by name.
		xname := protoreflect.FullName(name[len("[") : len(name)-len("]")])
		xt, _ := protoregistry.GlobalTypes.FindExtensionByName(xname)
		if xt == nil && isMessageSet(md) ***REMOVED***
			xt, _ = protoregistry.GlobalTypes.FindExtensionByName(xname.Append("message_set_extension"))
		***REMOVED***
		if xt == nil ***REMOVED***
			continue
		***REMOVED***
		delete(jsonObject, name)
		fd := xt.TypeDescriptor()
		if fd.ContainingMessage().FullName() != m.Descriptor().FullName() ***REMOVED***
			return fmt.Errorf("extension field %q does not extend message %q", xname, m.Descriptor().FullName())
		***REMOVED***

		field := m.NewField(fd)
		// Unmarshal the field value.
		if raw == nil || (string(raw) == "null" && !isSingularWellKnownValue(fd) && !isSingularJSONPBUnmarshaler(field, fd)) ***REMOVED***
			continue
		***REMOVED***
		v, err := u.unmarshalValue(field, raw, fd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m.Set(fd, v)
	***REMOVED***

	if !u.AllowUnknownFields && len(jsonObject) > 0 ***REMOVED***
		for name := range jsonObject ***REMOVED***
			return fmt.Errorf("unknown field %q in %v", name, md.FullName())
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func isSingularWellKnownValue(fd protoreflect.FieldDescriptor) bool ***REMOVED***
	if md := fd.Message(); md != nil ***REMOVED***
		return md.FullName() == "google.protobuf.Value" && fd.Cardinality() != protoreflect.Repeated
	***REMOVED***
	return false
***REMOVED***

func isSingularJSONPBUnmarshaler(v protoreflect.Value, fd protoreflect.FieldDescriptor) bool ***REMOVED***
	if fd.Message() != nil && fd.Cardinality() != protoreflect.Repeated ***REMOVED***
		_, ok := proto.MessageV1(v.Interface()).(JSONPBUnmarshaler)
		return ok
	***REMOVED***
	return false
***REMOVED***

func (u *Unmarshaler) unmarshalValue(v protoreflect.Value, in []byte, fd protoreflect.FieldDescriptor) (protoreflect.Value, error) ***REMOVED***
	switch ***REMOVED***
	case fd.IsList():
		var jsonArray []json.RawMessage
		if err := json.Unmarshal(in, &jsonArray); err != nil ***REMOVED***
			return v, err
		***REMOVED***
		lv := v.List()
		for _, raw := range jsonArray ***REMOVED***
			ve, err := u.unmarshalSingularValue(lv.NewElement(), raw, fd)
			if err != nil ***REMOVED***
				return v, err
			***REMOVED***
			lv.Append(ve)
		***REMOVED***
		return v, nil
	case fd.IsMap():
		var jsonObject map[string]json.RawMessage
		if err := json.Unmarshal(in, &jsonObject); err != nil ***REMOVED***
			return v, err
		***REMOVED***
		kfd := fd.MapKey()
		vfd := fd.MapValue()
		mv := v.Map()
		for key, raw := range jsonObject ***REMOVED***
			var kv protoreflect.MapKey
			if kfd.Kind() == protoreflect.StringKind ***REMOVED***
				kv = protoreflect.ValueOf(key).MapKey()
			***REMOVED*** else ***REMOVED***
				v, err := u.unmarshalSingularValue(kfd.Default(), []byte(key), kfd)
				if err != nil ***REMOVED***
					return v, err
				***REMOVED***
				kv = v.MapKey()
			***REMOVED***

			vv, err := u.unmarshalSingularValue(mv.NewValue(), raw, vfd)
			if err != nil ***REMOVED***
				return v, err
			***REMOVED***
			mv.Set(kv, vv)
		***REMOVED***
		return v, nil
	default:
		return u.unmarshalSingularValue(v, in, fd)
	***REMOVED***
***REMOVED***

var nonFinite = map[string]float64***REMOVED***
	`"NaN"`:       math.NaN(),
	`"Infinity"`:  math.Inf(+1),
	`"-Infinity"`: math.Inf(-1),
***REMOVED***

func (u *Unmarshaler) unmarshalSingularValue(v protoreflect.Value, in []byte, fd protoreflect.FieldDescriptor) (protoreflect.Value, error) ***REMOVED***
	switch fd.Kind() ***REMOVED***
	case protoreflect.BoolKind:
		return unmarshalValue(in, new(bool))
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return unmarshalValue(trimQuote(in), new(int32))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return unmarshalValue(trimQuote(in), new(int64))
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return unmarshalValue(trimQuote(in), new(uint32))
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return unmarshalValue(trimQuote(in), new(uint64))
	case protoreflect.FloatKind:
		if f, ok := nonFinite[string(in)]; ok ***REMOVED***
			return protoreflect.ValueOfFloat32(float32(f)), nil
		***REMOVED***
		return unmarshalValue(trimQuote(in), new(float32))
	case protoreflect.DoubleKind:
		if f, ok := nonFinite[string(in)]; ok ***REMOVED***
			return protoreflect.ValueOfFloat64(float64(f)), nil
		***REMOVED***
		return unmarshalValue(trimQuote(in), new(float64))
	case protoreflect.StringKind:
		return unmarshalValue(in, new(string))
	case protoreflect.BytesKind:
		return unmarshalValue(in, new([]byte))
	case protoreflect.EnumKind:
		if hasPrefixAndSuffix('"', in, '"') ***REMOVED***
			vd := fd.Enum().Values().ByName(protoreflect.Name(trimQuote(in)))
			if vd == nil ***REMOVED***
				return v, fmt.Errorf("unknown value %q for enum %s", in, fd.Enum().FullName())
			***REMOVED***
			return protoreflect.ValueOfEnum(vd.Number()), nil
		***REMOVED***
		return unmarshalValue(in, new(protoreflect.EnumNumber))
	case protoreflect.MessageKind, protoreflect.GroupKind:
		err := u.unmarshalMessage(v.Message(), in)
		return v, err
	default:
		panic(fmt.Sprintf("invalid kind %v", fd.Kind()))
	***REMOVED***
***REMOVED***

func unmarshalValue(in []byte, v interface***REMOVED******REMOVED***) (protoreflect.Value, error) ***REMOVED***
	err := json.Unmarshal(in, v)
	return protoreflect.ValueOf(reflect.ValueOf(v).Elem().Interface()), err
***REMOVED***

func unquoteString(in string) (out string, err error) ***REMOVED***
	err = json.Unmarshal([]byte(in), &out)
	return out, err
***REMOVED***

func hasPrefixAndSuffix(prefix byte, in []byte, suffix byte) bool ***REMOVED***
	if len(in) >= 2 && in[0] == prefix && in[len(in)-1] == suffix ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// trimQuote is like unquoteString but simply strips surrounding quotes.
// This is incorrect, but is behavior done by the legacy implementation.
func trimQuote(in []byte) []byte ***REMOVED***
	if len(in) >= 2 && in[0] == '"' && in[len(in)-1] == '"' ***REMOVED***
		in = in[1 : len(in)-1]
	***REMOVED***
	return in
***REMOVED***
