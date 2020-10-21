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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const wrapJSONMarshalV2 = false

// Marshaler is a configurable object for marshaling protocol buffer messages
// to the specified JSON representation.
type Marshaler struct ***REMOVED***
	// OrigName specifies whether to use the original protobuf name for fields.
	OrigName bool

	// EnumsAsInts specifies whether to render enum values as integers,
	// as opposed to string values.
	EnumsAsInts bool

	// EmitDefaults specifies whether to render fields with zero values.
	EmitDefaults bool

	// Indent controls whether the output is compact or not.
	// If empty, the output is compact JSON. Otherwise, every JSON object
	// entry and JSON array value will be on its own line.
	// Each line will be preceded by repeated copies of Indent, where the
	// number of copies is the current indentation depth.
	Indent string

	// AnyResolver is used to resolve the google.protobuf.Any well-known type.
	// If unset, the global registry is used by default.
	AnyResolver AnyResolver
***REMOVED***

// JSONPBMarshaler is implemented by protobuf messages that customize the
// way they are marshaled to JSON. Messages that implement this should also
// implement JSONPBUnmarshaler so that the custom format can be parsed.
//
// The JSON marshaling must follow the proto to JSON specification:
//	https://developers.google.com/protocol-buffers/docs/proto3#json
//
// Deprecated: Custom types should implement protobuf reflection instead.
type JSONPBMarshaler interface ***REMOVED***
	MarshalJSONPB(*Marshaler) ([]byte, error)
***REMOVED***

// Marshal serializes a protobuf message as JSON into w.
func (jm *Marshaler) Marshal(w io.Writer, m proto.Message) error ***REMOVED***
	b, err := jm.marshal(m)
	if len(b) > 0 ***REMOVED***
		if _, err := w.Write(b); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

// MarshalToString serializes a protobuf message as JSON in string form.
func (jm *Marshaler) MarshalToString(m proto.Message) (string, error) ***REMOVED***
	b, err := jm.marshal(m)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(b), nil
***REMOVED***

func (jm *Marshaler) marshal(m proto.Message) ([]byte, error) ***REMOVED***
	v := reflect.ValueOf(m)
	if m == nil || (v.Kind() == reflect.Ptr && v.IsNil()) ***REMOVED***
		return nil, errors.New("Marshal called with nil")
	***REMOVED***

	// Check for custom marshalers first since they may not properly
	// implement protobuf reflection that the logic below relies on.
	if jsm, ok := m.(JSONPBMarshaler); ok ***REMOVED***
		return jsm.MarshalJSONPB(jm)
	***REMOVED***

	if wrapJSONMarshalV2 ***REMOVED***
		opts := protojson.MarshalOptions***REMOVED***
			UseProtoNames:   jm.OrigName,
			UseEnumNumbers:  jm.EnumsAsInts,
			EmitUnpopulated: jm.EmitDefaults,
			Indent:          jm.Indent,
		***REMOVED***
		if jm.AnyResolver != nil ***REMOVED***
			opts.Resolver = anyResolver***REMOVED***jm.AnyResolver***REMOVED***
		***REMOVED***
		return opts.Marshal(proto.MessageReflect(m).Interface())
	***REMOVED*** else ***REMOVED***
		// Check for unpopulated required fields first.
		m2 := proto.MessageReflect(m)
		if err := protoV2.CheckInitialized(m2.Interface()); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		w := jsonWriter***REMOVED***Marshaler: jm***REMOVED***
		err := w.marshalMessage(m2, "", "")
		return w.buf, err
	***REMOVED***
***REMOVED***

type jsonWriter struct ***REMOVED***
	*Marshaler
	buf []byte
***REMOVED***

func (w *jsonWriter) write(s string) ***REMOVED***
	w.buf = append(w.buf, s...)
***REMOVED***

func (w *jsonWriter) marshalMessage(m protoreflect.Message, indent, typeURL string) error ***REMOVED***
	if jsm, ok := proto.MessageV1(m.Interface()).(JSONPBMarshaler); ok ***REMOVED***
		b, err := jsm.MarshalJSONPB(w.Marshaler)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if typeURL != "" ***REMOVED***
			// we are marshaling this object to an Any type
			var js map[string]*json.RawMessage
			if err = json.Unmarshal(b, &js); err != nil ***REMOVED***
				return fmt.Errorf("type %T produced invalid JSON: %v", m.Interface(), err)
			***REMOVED***
			turl, err := json.Marshal(typeURL)
			if err != nil ***REMOVED***
				return fmt.Errorf("failed to marshal type URL %q to JSON: %v", typeURL, err)
			***REMOVED***
			js["@type"] = (*json.RawMessage)(&turl)
			if b, err = json.Marshal(js); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		w.write(string(b))
		return nil
	***REMOVED***

	md := m.Descriptor()
	fds := md.Fields()

	// Handle well-known types.
	const secondInNanos = int64(time.Second / time.Nanosecond)
	switch wellKnownType(md.FullName()) ***REMOVED***
	case "Any":
		return w.marshalAny(m, indent)
	case "BoolValue", "BytesValue", "StringValue",
		"Int32Value", "UInt32Value", "FloatValue",
		"Int64Value", "UInt64Value", "DoubleValue":
		fd := fds.ByNumber(1)
		return w.marshalValue(fd, m.Get(fd), indent)
	case "Duration":
		// "Generated output always contains 0, 3, 6, or 9 fractional digits,
		//  depending on required precision."
		s := m.Get(fds.ByNumber(1)).Int()
		ns := m.Get(fds.ByNumber(2)).Int()
		if ns <= -secondInNanos || ns >= secondInNanos ***REMOVED***
			return fmt.Errorf("ns out of range (%v, %v)", -secondInNanos, secondInNanos)
		***REMOVED***
		if (s > 0 && ns < 0) || (s < 0 && ns > 0) ***REMOVED***
			return errors.New("signs of seconds and nanos do not match")
		***REMOVED***
		if s < 0 ***REMOVED***
			ns = -ns
		***REMOVED***
		x := fmt.Sprintf("%d.%09d", s, ns)
		x = strings.TrimSuffix(x, "000")
		x = strings.TrimSuffix(x, "000")
		x = strings.TrimSuffix(x, ".000")
		w.write(fmt.Sprintf(`"%vs"`, x))
		return nil
	case "Timestamp":
		// "RFC 3339, where generated output will always be Z-normalized
		//  and uses 0, 3, 6 or 9 fractional digits."
		s := m.Get(fds.ByNumber(1)).Int()
		ns := m.Get(fds.ByNumber(2)).Int()
		if ns < 0 || ns >= secondInNanos ***REMOVED***
			return fmt.Errorf("ns out of range [0, %v)", secondInNanos)
		***REMOVED***
		t := time.Unix(s, ns).UTC()
		// time.RFC3339Nano isn't exactly right (we need to get 3/6/9 fractional digits).
		x := t.Format("2006-01-02T15:04:05.000000000")
		x = strings.TrimSuffix(x, "000")
		x = strings.TrimSuffix(x, "000")
		x = strings.TrimSuffix(x, ".000")
		w.write(fmt.Sprintf(`"%vZ"`, x))
		return nil
	case "Value":
		// JSON value; which is a null, number, string, bool, object, or array.
		od := md.Oneofs().Get(0)
		fd := m.WhichOneof(od)
		if fd == nil ***REMOVED***
			return errors.New("nil Value")
		***REMOVED***
		return w.marshalValue(fd, m.Get(fd), indent)
	case "Struct", "ListValue":
		// JSON object or array.
		fd := fds.ByNumber(1)
		return w.marshalValue(fd, m.Get(fd), indent)
	***REMOVED***

	w.write("***REMOVED***")
	if w.Indent != "" ***REMOVED***
		w.write("\n")
	***REMOVED***

	firstField := true
	if typeURL != "" ***REMOVED***
		if err := w.marshalTypeURL(indent, typeURL); err != nil ***REMOVED***
			return err
		***REMOVED***
		firstField = false
	***REMOVED***

	for i := 0; i < fds.Len(); ***REMOVED***
		fd := fds.Get(i)
		if od := fd.ContainingOneof(); od != nil ***REMOVED***
			fd = m.WhichOneof(od)
			i += od.Fields().Len()
			if fd == nil ***REMOVED***
				continue
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			i++
		***REMOVED***

		v := m.Get(fd)

		if !m.Has(fd) ***REMOVED***
			if !w.EmitDefaults || fd.ContainingOneof() != nil ***REMOVED***
				continue
			***REMOVED***
			if fd.Cardinality() != protoreflect.Repeated && (fd.Message() != nil || fd.Syntax() == protoreflect.Proto2) ***REMOVED***
				v = protoreflect.Value***REMOVED******REMOVED*** // use "null" for singular messages or proto2 scalars
			***REMOVED***
		***REMOVED***

		if !firstField ***REMOVED***
			w.writeComma()
		***REMOVED***
		if err := w.marshalField(fd, v, indent); err != nil ***REMOVED***
			return err
		***REMOVED***
		firstField = false
	***REMOVED***

	// Handle proto2 extensions.
	if md.ExtensionRanges().Len() > 0 ***REMOVED***
		// Collect a sorted list of all extension descriptor and values.
		type ext struct ***REMOVED***
			desc protoreflect.FieldDescriptor
			val  protoreflect.Value
		***REMOVED***
		var exts []ext
		m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
			if fd.IsExtension() ***REMOVED***
				exts = append(exts, ext***REMOVED***fd, v***REMOVED***)
			***REMOVED***
			return true
		***REMOVED***)
		sort.Slice(exts, func(i, j int) bool ***REMOVED***
			return exts[i].desc.Number() < exts[j].desc.Number()
		***REMOVED***)

		for _, ext := range exts ***REMOVED***
			if !firstField ***REMOVED***
				w.writeComma()
			***REMOVED***
			if err := w.marshalField(ext.desc, ext.val, indent); err != nil ***REMOVED***
				return err
			***REMOVED***
			firstField = false
		***REMOVED***
	***REMOVED***

	if w.Indent != "" ***REMOVED***
		w.write("\n")
		w.write(indent)
	***REMOVED***
	w.write("***REMOVED***")
	return nil
***REMOVED***

func (w *jsonWriter) writeComma() ***REMOVED***
	if w.Indent != "" ***REMOVED***
		w.write(",\n")
	***REMOVED*** else ***REMOVED***
		w.write(",")
	***REMOVED***
***REMOVED***

func (w *jsonWriter) marshalAny(m protoreflect.Message, indent string) error ***REMOVED***
	// "If the Any contains a value that has a special JSON mapping,
	//  it will be converted as follows: ***REMOVED***"@type": xxx, "value": yyy***REMOVED***.
	//  Otherwise, the value will be converted into a JSON object,
	//  and the "@type" field will be inserted to indicate the actual data type."
	md := m.Descriptor()
	typeURL := m.Get(md.Fields().ByNumber(1)).String()
	rawVal := m.Get(md.Fields().ByNumber(2)).Bytes()

	var m2 protoreflect.Message
	if w.AnyResolver != nil ***REMOVED***
		mi, err := w.AnyResolver.Resolve(typeURL)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m2 = proto.MessageReflect(mi)
	***REMOVED*** else ***REMOVED***
		mt, err := protoregistry.GlobalTypes.FindMessageByURL(typeURL)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m2 = mt.New()
	***REMOVED***

	if err := protoV2.Unmarshal(rawVal, m2.Interface()); err != nil ***REMOVED***
		return err
	***REMOVED***

	if wellKnownType(m2.Descriptor().FullName()) == "" ***REMOVED***
		return w.marshalMessage(m2, indent, typeURL)
	***REMOVED***

	w.write("***REMOVED***")
	if w.Indent != "" ***REMOVED***
		w.write("\n")
	***REMOVED***
	if err := w.marshalTypeURL(indent, typeURL); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.writeComma()
	if w.Indent != "" ***REMOVED***
		w.write(indent)
		w.write(w.Indent)
		w.write(`"value": `)
	***REMOVED*** else ***REMOVED***
		w.write(`"value":`)
	***REMOVED***
	if err := w.marshalMessage(m2, indent+w.Indent, ""); err != nil ***REMOVED***
		return err
	***REMOVED***
	if w.Indent != "" ***REMOVED***
		w.write("\n")
		w.write(indent)
	***REMOVED***
	w.write("***REMOVED***")
	return nil
***REMOVED***

func (w *jsonWriter) marshalTypeURL(indent, typeURL string) error ***REMOVED***
	if w.Indent != "" ***REMOVED***
		w.write(indent)
		w.write(w.Indent)
	***REMOVED***
	w.write(`"@type":`)
	if w.Indent != "" ***REMOVED***
		w.write(" ")
	***REMOVED***
	b, err := json.Marshal(typeURL)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.write(string(b))
	return nil
***REMOVED***

// marshalField writes field description and value to the Writer.
func (w *jsonWriter) marshalField(fd protoreflect.FieldDescriptor, v protoreflect.Value, indent string) error ***REMOVED***
	if w.Indent != "" ***REMOVED***
		w.write(indent)
		w.write(w.Indent)
	***REMOVED***
	w.write(`"`)
	switch ***REMOVED***
	case fd.IsExtension():
		// For message set, use the fname of the message as the extension name.
		name := string(fd.FullName())
		if isMessageSet(fd.ContainingMessage()) ***REMOVED***
			name = strings.TrimSuffix(name, ".message_set_extension")
		***REMOVED***

		w.write("[" + name + "]")
	case w.OrigName:
		name := string(fd.Name())
		if fd.Kind() == protoreflect.GroupKind ***REMOVED***
			name = string(fd.Message().Name())
		***REMOVED***
		w.write(name)
	default:
		w.write(string(fd.JSONName()))
	***REMOVED***
	w.write(`":`)
	if w.Indent != "" ***REMOVED***
		w.write(" ")
	***REMOVED***
	return w.marshalValue(fd, v, indent)
***REMOVED***

func (w *jsonWriter) marshalValue(fd protoreflect.FieldDescriptor, v protoreflect.Value, indent string) error ***REMOVED***
	switch ***REMOVED***
	case fd.IsList():
		w.write("[")
		comma := ""
		lv := v.List()
		for i := 0; i < lv.Len(); i++ ***REMOVED***
			w.write(comma)
			if w.Indent != "" ***REMOVED***
				w.write("\n")
				w.write(indent)
				w.write(w.Indent)
				w.write(w.Indent)
			***REMOVED***
			if err := w.marshalSingularValue(fd, lv.Get(i), indent+w.Indent); err != nil ***REMOVED***
				return err
			***REMOVED***
			comma = ","
		***REMOVED***
		if w.Indent != "" ***REMOVED***
			w.write("\n")
			w.write(indent)
			w.write(w.Indent)
		***REMOVED***
		w.write("]")
		return nil
	case fd.IsMap():
		kfd := fd.MapKey()
		vfd := fd.MapValue()
		mv := v.Map()

		// Collect a sorted list of all map keys and values.
		type entry struct***REMOVED*** key, val protoreflect.Value ***REMOVED***
		var entries []entry
		mv.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
			entries = append(entries, entry***REMOVED***k.Value(), v***REMOVED***)
			return true
		***REMOVED***)
		sort.Slice(entries, func(i, j int) bool ***REMOVED***
			switch kfd.Kind() ***REMOVED***
			case protoreflect.BoolKind:
				return !entries[i].key.Bool() && entries[j].key.Bool()
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				return entries[i].key.Int() < entries[j].key.Int()
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				return entries[i].key.Uint() < entries[j].key.Uint()
			case protoreflect.StringKind:
				return entries[i].key.String() < entries[j].key.String()
			default:
				panic("invalid kind")
			***REMOVED***
		***REMOVED***)

		w.write(`***REMOVED***`)
		comma := ""
		for _, entry := range entries ***REMOVED***
			w.write(comma)
			if w.Indent != "" ***REMOVED***
				w.write("\n")
				w.write(indent)
				w.write(w.Indent)
				w.write(w.Indent)
			***REMOVED***

			s := fmt.Sprint(entry.key.Interface())
			b, err := json.Marshal(s)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			w.write(string(b))

			w.write(`:`)
			if w.Indent != "" ***REMOVED***
				w.write(` `)
			***REMOVED***

			if err := w.marshalSingularValue(vfd, entry.val, indent+w.Indent); err != nil ***REMOVED***
				return err
			***REMOVED***
			comma = ","
		***REMOVED***
		if w.Indent != "" ***REMOVED***
			w.write("\n")
			w.write(indent)
			w.write(w.Indent)
		***REMOVED***
		w.write(`***REMOVED***`)
		return nil
	default:
		return w.marshalSingularValue(fd, v, indent)
	***REMOVED***
***REMOVED***

func (w *jsonWriter) marshalSingularValue(fd protoreflect.FieldDescriptor, v protoreflect.Value, indent string) error ***REMOVED***
	switch ***REMOVED***
	case !v.IsValid():
		w.write("null")
		return nil
	case fd.Message() != nil:
		return w.marshalMessage(v.Message(), indent+w.Indent, "")
	case fd.Enum() != nil:
		if fd.Enum().FullName() == "google.protobuf.NullValue" ***REMOVED***
			w.write("null")
			return nil
		***REMOVED***

		vd := fd.Enum().Values().ByNumber(v.Enum())
		if vd == nil || w.EnumsAsInts ***REMOVED***
			w.write(strconv.Itoa(int(v.Enum())))
		***REMOVED*** else ***REMOVED***
			w.write(`"` + string(vd.Name()) + `"`)
		***REMOVED***
		return nil
	default:
		switch v.Interface().(type) ***REMOVED***
		case float32, float64:
			switch ***REMOVED***
			case math.IsInf(v.Float(), +1):
				w.write(`"Infinity"`)
				return nil
			case math.IsInf(v.Float(), -1):
				w.write(`"-Infinity"`)
				return nil
			case math.IsNaN(v.Float()):
				w.write(`"NaN"`)
				return nil
			***REMOVED***
		case int64, uint64:
			w.write(fmt.Sprintf(`"%d"`, v.Interface()))
			return nil
		***REMOVED***

		b, err := json.Marshal(v.Interface())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		w.write(string(b))
		return nil
	***REMOVED***
***REMOVED***
