package dynamic

// JSON marshalling and unmarshalling for dynamic messages

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	// link in the well-known-types that have a special JSON format
	_ "github.com/golang/protobuf/ptypes/any"
	_ "github.com/golang/protobuf/ptypes/duration"
	_ "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/golang/protobuf/ptypes/struct"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/golang/protobuf/ptypes/wrappers"

	"github.com/jhump/protoreflect/desc"
)

var wellKnownTypeNames = map[string]struct***REMOVED******REMOVED******REMOVED***
	"google.protobuf.Any":       ***REMOVED******REMOVED***,
	"google.protobuf.Empty":     ***REMOVED******REMOVED***,
	"google.protobuf.Duration":  ***REMOVED******REMOVED***,
	"google.protobuf.Timestamp": ***REMOVED******REMOVED***,
	// struct.proto
	"google.protobuf.Struct":    ***REMOVED******REMOVED***,
	"google.protobuf.Value":     ***REMOVED******REMOVED***,
	"google.protobuf.ListValue": ***REMOVED******REMOVED***,
	// wrappers.proto
	"google.protobuf.DoubleValue": ***REMOVED******REMOVED***,
	"google.protobuf.FloatValue":  ***REMOVED******REMOVED***,
	"google.protobuf.Int64Value":  ***REMOVED******REMOVED***,
	"google.protobuf.UInt64Value": ***REMOVED******REMOVED***,
	"google.protobuf.Int32Value":  ***REMOVED******REMOVED***,
	"google.protobuf.UInt32Value": ***REMOVED******REMOVED***,
	"google.protobuf.BoolValue":   ***REMOVED******REMOVED***,
	"google.protobuf.StringValue": ***REMOVED******REMOVED***,
	"google.protobuf.BytesValue":  ***REMOVED******REMOVED***,
***REMOVED***

// MarshalJSON serializes this message to bytes in JSON format, returning an
// error if the operation fails. The resulting bytes will be a valid UTF8
// string.
//
// This method uses a compact form: no newlines, and spaces between fields and
// between field identifiers and values are elided.
//
// This method is convenient shorthand for invoking MarshalJSONPB with a default
// (zero value) marshaler:
//
//    m.MarshalJSONPB(&jsonpb.Marshaler***REMOVED******REMOVED***)
//
// So enums are serialized using enum value name strings, and values that are
// not present (including those with default/zero value for messages defined in
// "proto3" syntax) are omitted.
func (m *Message) MarshalJSON() ([]byte, error) ***REMOVED***
	return m.MarshalJSONPB(&jsonpb.Marshaler***REMOVED******REMOVED***)
***REMOVED***

// MarshalJSONIndent serializes this message to bytes in JSON format, returning
// an error if the operation fails. The resulting bytes will be a valid UTF8
// string.
//
// This method uses a "pretty-printed" form, with each field on its own line and
// spaces between field identifiers and values. Indentation of two spaces is
// used.
//
// This method is convenient shorthand for invoking MarshalJSONPB with a default
// (zero value) marshaler:
//
//    m.MarshalJSONPB(&jsonpb.Marshaler***REMOVED***Indent: "  "***REMOVED***)
//
// So enums are serialized using enum value name strings, and values that are
// not present (including those with default/zero value for messages defined in
// "proto3" syntax) are omitted.
func (m *Message) MarshalJSONIndent() ([]byte, error) ***REMOVED***
	return m.MarshalJSONPB(&jsonpb.Marshaler***REMOVED***Indent: "  "***REMOVED***)
***REMOVED***

// MarshalJSONPB serializes this message to bytes in JSON format, returning an
// error if the operation fails. The resulting bytes will be a valid UTF8
// string. The given marshaler is used to convey options used during marshaling.
//
// If this message contains nested messages that are generated message types (as
// opposed to dynamic messages), the given marshaler is used to marshal it.
//
// When marshaling any nested messages, any jsonpb.AnyResolver configured in the
// given marshaler is augmented with knowledge of message types known to this
// message's descriptor (and its enclosing file and set of transitive
// dependencies).
func (m *Message) MarshalJSONPB(opts *jsonpb.Marshaler) ([]byte, error) ***REMOVED***
	var b indentBuffer
	b.indent = opts.Indent
	if len(opts.Indent) == 0 ***REMOVED***
		b.indentCount = -1
	***REMOVED***
	b.comma = true
	if err := m.marshalJSON(&b, opts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.Bytes(), nil
***REMOVED***

func (m *Message) marshalJSON(b *indentBuffer, opts *jsonpb.Marshaler) error ***REMOVED***
	if m == nil ***REMOVED***
		_, err := b.WriteString("null")
		return err
	***REMOVED***
	if r, changed := wrapResolver(opts.AnyResolver, m.mf, m.md.GetFile()); changed ***REMOVED***
		newOpts := *opts
		newOpts.AnyResolver = r
		opts = &newOpts
	***REMOVED***

	if ok, err := marshalWellKnownType(m, b, opts); ok ***REMOVED***
		return err
	***REMOVED***

	err := b.WriteByte('***REMOVED***')
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.start()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var tags []int
	if opts.EmitDefaults ***REMOVED***
		tags = m.allKnownFieldTags()
	***REMOVED*** else ***REMOVED***
		tags = m.knownFieldTags()
	***REMOVED***

	first := true

	for _, tag := range tags ***REMOVED***
		itag := int32(tag)
		fd := m.FindFieldDescriptor(itag)

		v, ok := m.values[itag]
		if !ok ***REMOVED***
			if fd.GetOneOf() != nil ***REMOVED***
				// don't print defaults for fields in a oneof
				continue
			***REMOVED***
			v = fd.GetDefaultValue()
		***REMOVED***

		err := b.maybeNext(&first)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = marshalKnownFieldJSON(b, fd, v, opts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	err = b.end()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.WriteByte('***REMOVED***')
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func marshalWellKnownType(m *Message, b *indentBuffer, opts *jsonpb.Marshaler) (bool, error) ***REMOVED***
	fqn := m.md.GetFullyQualifiedName()
	if _, ok := wellKnownTypeNames[fqn]; !ok ***REMOVED***
		return false, nil
	***REMOVED***

	msgType := proto.MessageType(fqn)
	if msgType == nil ***REMOVED***
		// wtf?
		panic(fmt.Sprintf("could not find registered message type for %q", fqn))
	***REMOVED***

	// convert dynamic message to well-known type and let jsonpb marshal it
	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
	if err := m.MergeInto(msg); err != nil ***REMOVED***
		return true, err
	***REMOVED***
	return true, opts.Marshal(b, msg)
***REMOVED***

func marshalKnownFieldJSON(b *indentBuffer, fd *desc.FieldDescriptor, v interface***REMOVED******REMOVED***, opts *jsonpb.Marshaler) error ***REMOVED***
	var jsonName string
	if opts.OrigName ***REMOVED***
		jsonName = fd.GetName()
	***REMOVED*** else ***REMOVED***
		jsonName = fd.AsFieldDescriptorProto().GetJsonName()
		if jsonName == "" ***REMOVED***
			jsonName = fd.GetName()
		***REMOVED***
	***REMOVED***
	if fd.IsExtension() ***REMOVED***
		var scope string
		switch parent := fd.GetParent().(type) ***REMOVED***
		case *desc.FileDescriptor:
			scope = parent.GetPackage()
		default:
			scope = parent.GetFullyQualifiedName()
		***REMOVED***
		if scope == "" ***REMOVED***
			jsonName = fmt.Sprintf("[%s]", jsonName)
		***REMOVED*** else ***REMOVED***
			jsonName = fmt.Sprintf("[%s.%s]", scope, jsonName)
		***REMOVED***
	***REMOVED***
	err := writeJsonString(b, jsonName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.sep()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if isNil(v) ***REMOVED***
		_, err := b.WriteString("null")
		return err
	***REMOVED***

	if fd.IsMap() ***REMOVED***
		err = b.WriteByte('***REMOVED***')
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = b.start()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		md := fd.GetMessageType()
		vfd := md.FindFieldByNumber(2)

		mp := v.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
		keys := make([]interface***REMOVED******REMOVED***, 0, len(mp))
		for k := range mp ***REMOVED***
			keys = append(keys, k)
		***REMOVED***
		sort.Sort(sortable(keys))
		first := true
		for _, mk := range keys ***REMOVED***
			mv := mp[mk]
			err := b.maybeNext(&first)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			err = marshalKnownFieldMapEntryJSON(b, mk, vfd, mv, opts)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		err = b.end()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return b.WriteByte('***REMOVED***')

	***REMOVED*** else if fd.IsRepeated() ***REMOVED***
		err = b.WriteByte('[')
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = b.start()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		sl := v.([]interface***REMOVED******REMOVED***)
		first := true
		for _, slv := range sl ***REMOVED***
			err := b.maybeNext(&first)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = marshalKnownFieldValueJSON(b, fd, slv, opts)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		err = b.end()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return b.WriteByte(']')

	***REMOVED*** else ***REMOVED***
		return marshalKnownFieldValueJSON(b, fd, v, opts)
	***REMOVED***
***REMOVED***

// sortable is used to sort map keys. Values will be integers (int32, int64, uint32, and uint64),
// bools, or strings.
type sortable []interface***REMOVED******REMOVED***

func (s sortable) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s sortable) Less(i, j int) bool ***REMOVED***
	vi := s[i]
	vj := s[j]
	switch reflect.TypeOf(vi).Kind() ***REMOVED***
	case reflect.Int32:
		return vi.(int32) < vj.(int32)
	case reflect.Int64:
		return vi.(int64) < vj.(int64)
	case reflect.Uint32:
		return vi.(uint32) < vj.(uint32)
	case reflect.Uint64:
		return vi.(uint64) < vj.(uint64)
	case reflect.String:
		return vi.(string) < vj.(string)
	case reflect.Bool:
		return !vi.(bool) && vj.(bool)
	default:
		panic(fmt.Sprintf("cannot compare keys of type %v", reflect.TypeOf(vi)))
	***REMOVED***
***REMOVED***

func (s sortable) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func isNil(v interface***REMOVED******REMOVED***) bool ***REMOVED***
	if v == nil ***REMOVED***
		return true
	***REMOVED***
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
***REMOVED***

func marshalKnownFieldMapEntryJSON(b *indentBuffer, mk interface***REMOVED******REMOVED***, vfd *desc.FieldDescriptor, mv interface***REMOVED******REMOVED***, opts *jsonpb.Marshaler) error ***REMOVED***
	rk := reflect.ValueOf(mk)
	var strkey string
	switch rk.Kind() ***REMOVED***
	case reflect.Bool:
		strkey = strconv.FormatBool(rk.Bool())
	case reflect.Int32, reflect.Int64:
		strkey = strconv.FormatInt(rk.Int(), 10)
	case reflect.Uint32, reflect.Uint64:
		strkey = strconv.FormatUint(rk.Uint(), 10)
	case reflect.String:
		strkey = rk.String()
	default:
		return fmt.Errorf("invalid map key value: %v (%v)", mk, rk.Type())
	***REMOVED***
	err := writeString(b, strkey)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.sep()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return marshalKnownFieldValueJSON(b, vfd, mv, opts)
***REMOVED***

func marshalKnownFieldValueJSON(b *indentBuffer, fd *desc.FieldDescriptor, v interface***REMOVED******REMOVED***, opts *jsonpb.Marshaler) error ***REMOVED***
	rv := reflect.ValueOf(v)
	switch rv.Kind() ***REMOVED***
	case reflect.Int64:
		return writeJsonString(b, strconv.FormatInt(rv.Int(), 10))
	case reflect.Int32:
		ed := fd.GetEnumType()
		if !opts.EnumsAsInts && ed != nil ***REMOVED***
			n := int32(rv.Int())
			vd := ed.FindValueByNumber(n)
			if vd == nil ***REMOVED***
				_, err := b.WriteString(strconv.FormatInt(rv.Int(), 10))
				return err
			***REMOVED*** else ***REMOVED***
				return writeJsonString(b, vd.GetName())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			_, err := b.WriteString(strconv.FormatInt(rv.Int(), 10))
			return err
		***REMOVED***
	case reflect.Uint64:
		return writeJsonString(b, strconv.FormatUint(rv.Uint(), 10))
	case reflect.Uint32:
		_, err := b.WriteString(strconv.FormatUint(rv.Uint(), 10))
		return err
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		var str string
		if math.IsNaN(f) ***REMOVED***
			str = `"NaN"`
		***REMOVED*** else if math.IsInf(f, 1) ***REMOVED***
			str = `"Infinity"`
		***REMOVED*** else if math.IsInf(f, -1) ***REMOVED***
			str = `"-Infinity"`
		***REMOVED*** else ***REMOVED***
			var bits int
			if rv.Kind() == reflect.Float32 ***REMOVED***
				bits = 32
			***REMOVED*** else ***REMOVED***
				bits = 64
			***REMOVED***
			str = strconv.FormatFloat(rv.Float(), 'g', -1, bits)
		***REMOVED***
		_, err := b.WriteString(str)
		return err
	case reflect.Bool:
		_, err := b.WriteString(strconv.FormatBool(rv.Bool()))
		return err
	case reflect.Slice:
		bstr := base64.StdEncoding.EncodeToString(rv.Bytes())
		return writeJsonString(b, bstr)
	case reflect.String:
		return writeJsonString(b, rv.String())
	default:
		// must be a message
		if isNil(v) ***REMOVED***
			_, err := b.WriteString("null")
			return err
		***REMOVED***

		if dm, ok := v.(*Message); ok ***REMOVED***
			return dm.marshalJSON(b, opts)
		***REMOVED***

		var err error
		if b.indentCount <= 0 || len(b.indent) == 0 ***REMOVED***
			err = opts.Marshal(b, v.(proto.Message))
		***REMOVED*** else ***REMOVED***
			str, err := opts.MarshalToString(v.(proto.Message))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			indent := strings.Repeat(b.indent, b.indentCount)
			pos := 0
			// add indention prefix to each line
			for pos < len(str) ***REMOVED***
				start := pos
				nextPos := strings.Index(str[pos:], "\n")
				if nextPos == -1 ***REMOVED***
					nextPos = len(str)
				***REMOVED*** else ***REMOVED***
					nextPos = pos + nextPos + 1 // include newline
				***REMOVED***
				line := str[start:nextPos]
				if pos > 0 ***REMOVED***
					_, err = b.WriteString(indent)
					if err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				_, err = b.WriteString(line)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				pos = nextPos
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

func writeJsonString(b *indentBuffer, s string) error ***REMOVED***
	if sbytes, err := json.Marshal(s); err != nil ***REMOVED***
		return err
	***REMOVED*** else ***REMOVED***
		_, err := b.Write(sbytes)
		return err
	***REMOVED***
***REMOVED***

// UnmarshalJSON de-serializes the message that is present, in JSON format, in
// the given bytes into this message. It first resets the current message. It
// returns an error if the given bytes do not contain a valid encoding of this
// message type in JSON format.
//
// This method is shorthand for invoking UnmarshalJSONPB with a default (zero
// value) unmarshaler:
//
//    m.UnmarshalMergeJSONPB(&jsonpb.Unmarshaler***REMOVED******REMOVED***, js)
//
// So unknown fields will result in an error, and no provided jsonpb.AnyResolver
// will be used when parsing google.protobuf.Any messages.
func (m *Message) UnmarshalJSON(js []byte) error ***REMOVED***
	return m.UnmarshalJSONPB(&jsonpb.Unmarshaler***REMOVED******REMOVED***, js)
***REMOVED***

// UnmarshalMergeJSON de-serializes the message that is present, in JSON format,
// in the given bytes into this message. Unlike UnmarshalJSON, it does not first
// reset the message, instead merging the data in the given bytes into the
// existing data in this message.
func (m *Message) UnmarshalMergeJSON(js []byte) error ***REMOVED***
	return m.UnmarshalMergeJSONPB(&jsonpb.Unmarshaler***REMOVED******REMOVED***, js)
***REMOVED***

// UnmarshalJSONPB de-serializes the message that is present, in JSON format, in
// the given bytes into this message. The given unmarshaler conveys options used
// when parsing the JSON. This function first resets the current message. It
// returns an error if the given bytes do not contain a valid encoding of this
// message type in JSON format.
//
// The decoding is lenient:
//  1. The JSON can refer to fields either by their JSON name or by their
//     declared name.
//  2. The JSON can use either numeric values or string names for enum values.
//
// When instantiating nested messages, if this message's associated factory
// returns a generated message type (as opposed to a dynamic message), the given
// unmarshaler is used to unmarshal it.
//
// When unmarshaling any nested messages, any jsonpb.AnyResolver configured in
// the given unmarshaler is augmented with knowledge of message types known to
// this message's descriptor (and its enclosing file and set of transitive
// dependencies).
func (m *Message) UnmarshalJSONPB(opts *jsonpb.Unmarshaler, js []byte) error ***REMOVED***
	m.Reset()
	if err := m.UnmarshalMergeJSONPB(opts, js); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.Validate()
***REMOVED***

// UnmarshalMergeJSONPB de-serializes the message that is present, in JSON
// format, in the given bytes into this message. The given unmarshaler conveys
// options used when parsing the JSON. Unlike UnmarshalJSONPB, it does not first
// reset the message, instead merging the data in the given bytes into the
// existing data in this message.
func (m *Message) UnmarshalMergeJSONPB(opts *jsonpb.Unmarshaler, js []byte) error ***REMOVED***
	r := newJsReader(js)
	err := m.unmarshalJson(r, opts)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if t, err := r.poll(); err != io.EOF ***REMOVED***
		b, _ := ioutil.ReadAll(r.unread())
		s := fmt.Sprintf("%v%s", t, string(b))
		return fmt.Errorf("superfluous data found after JSON object: %q", s)
	***REMOVED***
	return nil
***REMOVED***

func unmarshalWellKnownType(m *Message, r *jsReader, opts *jsonpb.Unmarshaler) (bool, error) ***REMOVED***
	fqn := m.md.GetFullyQualifiedName()
	if _, ok := wellKnownTypeNames[fqn]; !ok ***REMOVED***
		return false, nil
	***REMOVED***

	msgType := proto.MessageType(fqn)
	if msgType == nil ***REMOVED***
		// wtf?
		panic(fmt.Sprintf("could not find registered message type for %q", fqn))
	***REMOVED***

	// extract json value from r
	var js json.RawMessage
	if err := json.NewDecoder(r.unread()).Decode(&js); err != nil ***REMOVED***
		return true, err
	***REMOVED***
	if err := r.skip(); err != nil ***REMOVED***
		return true, err
	***REMOVED***

	// unmarshal into well-known type and then convert to dynamic message
	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
	if err := opts.Unmarshal(bytes.NewReader(js), msg); err != nil ***REMOVED***
		return true, err
	***REMOVED***
	return true, m.MergeFrom(msg)
***REMOVED***

func (m *Message) unmarshalJson(r *jsReader, opts *jsonpb.Unmarshaler) error ***REMOVED***
	if r, changed := wrapResolver(opts.AnyResolver, m.mf, m.md.GetFile()); changed ***REMOVED***
		newOpts := *opts
		newOpts.AnyResolver = r
		opts = &newOpts
	***REMOVED***

	if ok, err := unmarshalWellKnownType(m, r, opts); ok ***REMOVED***
		return err
	***REMOVED***

	t, err := r.peek()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if t == nil ***REMOVED***
		// if json is simply "null" we do nothing
		r.poll()
		return nil
	***REMOVED***

	if err := r.beginObject(); err != nil ***REMOVED***
		return err
	***REMOVED***

	for r.hasNext() ***REMOVED***
		f, err := r.nextObjectKey()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		fd := m.FindFieldDescriptorByJSONName(f)
		if fd == nil ***REMOVED***
			if opts.AllowUnknownFields ***REMOVED***
				r.skip()
				continue
			***REMOVED***
			return fmt.Errorf("message type %s has no known field named %s", m.md.GetFullyQualifiedName(), f)
		***REMOVED***
		v, err := unmarshalJsField(fd, r, m.mf, opts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if v != nil ***REMOVED***
			if err := mergeField(m, fd, v); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if fd.GetOneOf() != nil ***REMOVED***
			// preserve explicit null for oneof fields (this is a little odd but
			// mimics the behavior of jsonpb with oneofs in generated message types)
			if fd.GetMessageType() != nil ***REMOVED***
				typ := m.mf.GetKnownTypeRegistry().GetKnownType(fd.GetMessageType().GetFullyQualifiedName())
				if typ != nil ***REMOVED***
					// typed nil
					if typ.Kind() != reflect.Ptr ***REMOVED***
						typ = reflect.PtrTo(typ)
					***REMOVED***
					v = reflect.Zero(typ).Interface()
				***REMOVED*** else ***REMOVED***
					// can't use nil dynamic message, so we just use empty one instead
					v = m.mf.NewDynamicMessage(fd.GetMessageType())
				***REMOVED***
				if err := m.setField(fd, v); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// not a message... explicit null makes no sense
				return fmt.Errorf("message type %s cannot set field %s to null: it is not a message type", m.md.GetFullyQualifiedName(), f)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			m.clearField(fd)
		***REMOVED***
	***REMOVED***

	if err := r.endObject(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func isWellKnownValue(fd *desc.FieldDescriptor) bool ***REMOVED***
	return !fd.IsRepeated() && fd.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE &&
		fd.GetMessageType().GetFullyQualifiedName() == "google.protobuf.Value"
***REMOVED***

func isWellKnownListValue(fd *desc.FieldDescriptor) bool ***REMOVED***
	return !fd.IsRepeated() && fd.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE &&
		fd.GetMessageType().GetFullyQualifiedName() == "google.protobuf.ListValue"
***REMOVED***

func unmarshalJsField(fd *desc.FieldDescriptor, r *jsReader, mf *MessageFactory, opts *jsonpb.Unmarshaler) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	t, err := r.peek()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if t == nil && !isWellKnownValue(fd) ***REMOVED***
		// if value is null, just return nil
		// (unless field is google.protobuf.Value, in which case
		// we fall through to parse it as an instance where its
		// underlying value is set to a NullValue)
		r.poll()
		return nil, nil
	***REMOVED***

	if t == json.Delim('***REMOVED***') && fd.IsMap() ***REMOVED***
		entryType := fd.GetMessageType()
		keyType := entryType.FindFieldByNumber(1)
		valueType := entryType.FindFieldByNumber(2)
		mp := map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***

		// TODO: if there are just two map keys "key" and "value" and they have the right type of values,
		// treat this JSON object as a single map entry message. (In keeping with support of map fields as
		// if they were normal repeated field of entry messages as well as supporting a transition from
		// optional to repeated...)

		if err := r.beginObject(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for r.hasNext() ***REMOVED***
			kk, err := unmarshalJsFieldElement(keyType, r, mf, opts, false)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			vv, err := unmarshalJsFieldElement(valueType, r, mf, opts, true)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			mp[kk] = vv
		***REMOVED***
		if err := r.endObject(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return mp, nil
	***REMOVED*** else if t == json.Delim('[') && !isWellKnownListValue(fd) ***REMOVED***
		// We support parsing an array, even if field is not repeated, to mimic support in proto
		// binary wire format that supports changing an optional field to repeated and vice versa.
		// If the field is not repeated, we only keep the last value in the array.

		if err := r.beginArray(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var sl []interface***REMOVED******REMOVED***
		var v interface***REMOVED******REMOVED***
		for r.hasNext() ***REMOVED***
			var err error
			v, err = unmarshalJsFieldElement(fd, r, mf, opts, false)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if fd.IsRepeated() && v != nil ***REMOVED***
				sl = append(sl, v)
			***REMOVED***
		***REMOVED***
		if err := r.endArray(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if fd.IsMap() ***REMOVED***
			mp := map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***
			for _, m := range sl ***REMOVED***
				msg := m.(*Message)
				kk, err := msg.TryGetFieldByNumber(1)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				vv, err := msg.TryGetFieldByNumber(2)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				mp[kk] = vv
			***REMOVED***
			return mp, nil
		***REMOVED*** else if fd.IsRepeated() ***REMOVED***
			return sl, nil
		***REMOVED*** else ***REMOVED***
			return v, nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// We support parsing a singular value, even if field is repeated, to mimic support in proto
		// binary wire format that supports changing an optional field to repeated and vice versa.
		// If the field is repeated, we store value as singleton slice of that one value.

		v, err := unmarshalJsFieldElement(fd, r, mf, opts, false)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if v == nil ***REMOVED***
			return nil, nil
		***REMOVED***
		if fd.IsRepeated() ***REMOVED***
			return []interface***REMOVED******REMOVED******REMOVED***v***REMOVED***, nil
		***REMOVED*** else ***REMOVED***
			return v, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func unmarshalJsFieldElement(fd *desc.FieldDescriptor, r *jsReader, mf *MessageFactory, opts *jsonpb.Unmarshaler, allowNilMessage bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	t, err := r.peek()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch fd.GetType() ***REMOVED***
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_GROUP:

		if t == nil && allowNilMessage ***REMOVED***
			// if json is simply "null" return a nil pointer
			r.poll()
			return nilMessage(fd.GetMessageType()), nil
		***REMOVED***

		m := mf.NewMessage(fd.GetMessageType())
		if dm, ok := m.(*Message); ok ***REMOVED***
			if err := dm.unmarshalJson(r, opts); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			var msg json.RawMessage
			if err := json.NewDecoder(r.unread()).Decode(&msg); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err := r.skip(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err := opts.Unmarshal(bytes.NewReader([]byte(msg)), m); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return m, nil

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		if e, err := r.nextNumber(); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else ***REMOVED***
			// value could be string or number
			if i, err := e.Int64(); err != nil ***REMOVED***
				// number cannot be parsed, so see if it's an enum value name
				vd := fd.GetEnumType().FindValueByName(string(e))
				if vd != nil ***REMOVED***
					return vd.GetNumber(), nil
				***REMOVED*** else ***REMOVED***
					return nil, fmt.Errorf("enum %q does not have value named %q", fd.GetEnumType().GetFullyQualifiedName(), e)
				***REMOVED***
			***REMOVED*** else if i > math.MaxInt32 || i < math.MinInt32 ***REMOVED***
				return nil, NumericOverflowError
			***REMOVED*** else ***REMOVED***
				return int32(i), err
			***REMOVED***
		***REMOVED***

	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		if i, err := r.nextInt(); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if i > math.MaxInt32 || i < math.MinInt32 ***REMOVED***
			return nil, NumericOverflowError
		***REMOVED*** else ***REMOVED***
			return int32(i), err
		***REMOVED***

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return r.nextInt()

	case descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32:
		if i, err := r.nextUint(); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if i > math.MaxUint32 ***REMOVED***
			return nil, NumericOverflowError
		***REMOVED*** else ***REMOVED***
			return uint32(i), err
		***REMOVED***

	case descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return r.nextUint()

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if str, ok := t.(string); ok ***REMOVED***
			if str == "true" ***REMOVED***
				r.poll() // consume token
				return true, err
			***REMOVED*** else if str == "false" ***REMOVED***
				r.poll() // consume token
				return false, err
			***REMOVED***
		***REMOVED***
		return r.nextBool()

	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if f, err := r.nextFloat(); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else ***REMOVED***
			return float32(f), nil
		***REMOVED***

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return r.nextFloat()

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return r.nextBytes()

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return r.nextString()

	default:
		return nil, fmt.Errorf("unknown field type: %v", fd.GetType())
	***REMOVED***
***REMOVED***

type jsReader struct ***REMOVED***
	reader  *bytes.Reader
	dec     *json.Decoder
	current json.Token
	peeked  bool
***REMOVED***

func newJsReader(b []byte) *jsReader ***REMOVED***
	reader := bytes.NewReader(b)
	dec := json.NewDecoder(reader)
	dec.UseNumber()
	return &jsReader***REMOVED***reader: reader, dec: dec***REMOVED***
***REMOVED***

func (r *jsReader) unread() io.Reader ***REMOVED***
	bufs := make([]io.Reader, 3)
	var peeked []byte
	if r.peeked ***REMOVED***
		if _, ok := r.current.(json.Delim); ok ***REMOVED***
			peeked = []byte(fmt.Sprintf("%v", r.current))
		***REMOVED*** else ***REMOVED***
			peeked, _ = json.Marshal(r.current)
		***REMOVED***
	***REMOVED***
	readerCopy := *r.reader
	decCopy := *r.dec

	bufs[0] = bytes.NewReader(peeked)
	bufs[1] = decCopy.Buffered()
	bufs[2] = &readerCopy
	return &concatReader***REMOVED***bufs: bufs***REMOVED***
***REMOVED***

func (r *jsReader) hasNext() bool ***REMOVED***
	return r.dec.More()
***REMOVED***

func (r *jsReader) peek() (json.Token, error) ***REMOVED***
	if r.peeked ***REMOVED***
		return r.current, nil
	***REMOVED***
	t, err := r.dec.Token()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	r.peeked = true
	r.current = t
	return t, nil
***REMOVED***

func (r *jsReader) poll() (json.Token, error) ***REMOVED***
	if r.peeked ***REMOVED***
		ret := r.current
		r.current = nil
		r.peeked = false
		return ret, nil
	***REMOVED***
	return r.dec.Token()
***REMOVED***

func (r *jsReader) beginObject() error ***REMOVED***
	_, err := r.expect(func(t json.Token) bool ***REMOVED*** return t == json.Delim('***REMOVED***') ***REMOVED***, nil, "start of JSON object: '***REMOVED***'")
	return err
***REMOVED***

func (r *jsReader) endObject() error ***REMOVED***
	_, err := r.expect(func(t json.Token) bool ***REMOVED*** return t == json.Delim('***REMOVED***') ***REMOVED***, nil, "end of JSON object: '***REMOVED***'")
	return err
***REMOVED***

func (r *jsReader) beginArray() error ***REMOVED***
	_, err := r.expect(func(t json.Token) bool ***REMOVED*** return t == json.Delim('[') ***REMOVED***, nil, "start of array: '['")
	return err
***REMOVED***

func (r *jsReader) endArray() error ***REMOVED***
	_, err := r.expect(func(t json.Token) bool ***REMOVED*** return t == json.Delim(']') ***REMOVED***, nil, "end of array: ']'")
	return err
***REMOVED***

func (r *jsReader) nextObjectKey() (string, error) ***REMOVED***
	return r.nextString()
***REMOVED***

func (r *jsReader) nextString() (string, error) ***REMOVED***
	t, err := r.expect(func(t json.Token) bool ***REMOVED*** _, ok := t.(string); return ok ***REMOVED***, "", "string")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return t.(string), nil
***REMOVED***

func (r *jsReader) nextBytes() ([]byte, error) ***REMOVED***
	str, err := r.nextString()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return base64.StdEncoding.DecodeString(str)
***REMOVED***

func (r *jsReader) nextBool() (bool, error) ***REMOVED***
	t, err := r.expect(func(t json.Token) bool ***REMOVED*** _, ok := t.(bool); return ok ***REMOVED***, false, "boolean")
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return t.(bool), nil
***REMOVED***

func (r *jsReader) nextInt() (int64, error) ***REMOVED***
	n, err := r.nextNumber()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return n.Int64()
***REMOVED***

func (r *jsReader) nextUint() (uint64, error) ***REMOVED***
	n, err := r.nextNumber()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return strconv.ParseUint(string(n), 10, 64)
***REMOVED***

func (r *jsReader) nextFloat() (float64, error) ***REMOVED***
	n, err := r.nextNumber()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return n.Float64()
***REMOVED***

func (r *jsReader) nextNumber() (json.Number, error) ***REMOVED***
	t, err := r.expect(func(t json.Token) bool ***REMOVED*** return reflect.TypeOf(t).Kind() == reflect.String ***REMOVED***, "0", "number")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	switch t := t.(type) ***REMOVED***
	case json.Number:
		return t, nil
	case string:
		return json.Number(t), nil
	***REMOVED***
	return "", fmt.Errorf("expecting a number but got %v", t)
***REMOVED***

func (r *jsReader) skip() error ***REMOVED***
	t, err := r.poll()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if t == json.Delim('[') ***REMOVED***
		if err := r.skipArray(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else if t == json.Delim('***REMOVED***') ***REMOVED***
		if err := r.skipObject(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *jsReader) skipArray() error ***REMOVED***
	for r.hasNext() ***REMOVED***
		if err := r.skip(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := r.endArray(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (r *jsReader) skipObject() error ***REMOVED***
	for r.hasNext() ***REMOVED***
		// skip object key
		if err := r.skip(); err != nil ***REMOVED***
			return err
		***REMOVED***
		// and value
		if err := r.skip(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := r.endObject(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (r *jsReader) expect(predicate func(json.Token) bool, ifNil interface***REMOVED******REMOVED***, expected string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	t, err := r.poll()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if t == nil && ifNil != nil ***REMOVED***
		return ifNil, nil
	***REMOVED***
	if !predicate(t) ***REMOVED***
		return t, fmt.Errorf("bad input: expecting %s ; instead got %v", expected, t)
	***REMOVED***
	return t, nil
***REMOVED***

type concatReader struct ***REMOVED***
	bufs []io.Reader
	curr int
***REMOVED***

func (r *concatReader) Read(p []byte) (n int, err error) ***REMOVED***
	for ***REMOVED***
		if r.curr >= len(r.bufs) ***REMOVED***
			err = io.EOF
			return
		***REMOVED***
		var c int
		c, err = r.bufs[r.curr].Read(p)
		n += c
		if err != io.EOF ***REMOVED***
			return
		***REMOVED***
		r.curr++
		p = p[c:]
	***REMOVED***
***REMOVED***

// AnyResolver returns a jsonpb.AnyResolver that uses the given file descriptors
// to resolve message names. It uses the given factory, which may be nil, to
// instantiate messages. The messages that it returns when resolving a type name
// may often be dynamic messages.
func AnyResolver(mf *MessageFactory, files ...*desc.FileDescriptor) jsonpb.AnyResolver ***REMOVED***
	return &anyResolver***REMOVED***mf: mf, files: files***REMOVED***
***REMOVED***

type anyResolver struct ***REMOVED***
	mf      *MessageFactory
	files   []*desc.FileDescriptor
	ignored map[*desc.FileDescriptor]struct***REMOVED******REMOVED***
	other   jsonpb.AnyResolver
***REMOVED***

func wrapResolver(r jsonpb.AnyResolver, mf *MessageFactory, f *desc.FileDescriptor) (jsonpb.AnyResolver, bool) ***REMOVED***
	if r, ok := r.(*anyResolver); ok ***REMOVED***
		if _, ok := r.ignored[f]; ok ***REMOVED***
			// if the current resolver is ignoring this file, it's because another
			// (upstream) resolver is already handling it, so nothing to do
			return r, false
		***REMOVED***
		for _, file := range r.files ***REMOVED***
			if file == f ***REMOVED***
				// no need to wrap!
				return r, false
			***REMOVED***
		***REMOVED***
		// ignore files that will be checked by the resolver we're wrapping
		// (we'll just delegate and let it search those files)
		ignored := map[*desc.FileDescriptor]struct***REMOVED******REMOVED******REMOVED******REMOVED***
		for i := range r.ignored ***REMOVED***
			ignored[i] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
		ignore(r.files, ignored)
		return &anyResolver***REMOVED***mf: mf, files: []*desc.FileDescriptor***REMOVED***f***REMOVED***, ignored: ignored, other: r***REMOVED***, true
	***REMOVED***
	return &anyResolver***REMOVED***mf: mf, files: []*desc.FileDescriptor***REMOVED***f***REMOVED***, other: r***REMOVED***, true
***REMOVED***

func ignore(files []*desc.FileDescriptor, ignored map[*desc.FileDescriptor]struct***REMOVED******REMOVED***) ***REMOVED***
	for _, f := range files ***REMOVED***
		if _, ok := ignored[f]; ok ***REMOVED***
			continue
		***REMOVED***
		ignored[f] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		ignore(f.GetDependencies(), ignored)
	***REMOVED***
***REMOVED***

func (r *anyResolver) Resolve(typeUrl string) (proto.Message, error) ***REMOVED***
	mname := typeUrl
	if slash := strings.LastIndex(mname, "/"); slash >= 0 ***REMOVED***
		mname = mname[slash+1:]
	***REMOVED***

	// see if the user-specified resolver is able to do the job
	if r.other != nil ***REMOVED***
		msg, err := r.other.Resolve(typeUrl)
		if err == nil ***REMOVED***
			return msg, nil
		***REMOVED***
	***REMOVED***

	// try to find the message in our known set of files
	checked := map[*desc.FileDescriptor]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, f := range r.files ***REMOVED***
		md := r.findMessage(f, mname, checked)
		if md != nil ***REMOVED***
			return r.mf.NewMessage(md), nil
		***REMOVED***
	***REMOVED***
	// failing that, see if the message factory knows about this type
	var ktr *KnownTypeRegistry
	if r.mf != nil ***REMOVED***
		ktr = r.mf.ktr
	***REMOVED*** else ***REMOVED***
		ktr = (*KnownTypeRegistry)(nil)
	***REMOVED***
	m := ktr.CreateIfKnown(mname)
	if m != nil ***REMOVED***
		return m, nil
	***REMOVED***

	// no other resolver to fallback to? mimic default behavior
	mt := proto.MessageType(mname)
	if mt == nil ***REMOVED***
		return nil, fmt.Errorf("unknown message type %q", mname)
	***REMOVED***
	return reflect.New(mt.Elem()).Interface().(proto.Message), nil
***REMOVED***

func (r *anyResolver) findMessage(fd *desc.FileDescriptor, msgName string, checked map[*desc.FileDescriptor]struct***REMOVED******REMOVED***) *desc.MessageDescriptor ***REMOVED***
	// if this is an ignored descriptor, skip
	if _, ok := r.ignored[fd]; ok ***REMOVED***
		return nil
	***REMOVED***

	// bail if we've already checked this file
	if _, ok := checked[fd]; ok ***REMOVED***
		return nil
	***REMOVED***
	checked[fd] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	// see if this file has the message
	md := fd.FindMessage(msgName)
	if md != nil ***REMOVED***
		return md
	***REMOVED***

	// if not, recursively search the file's imports
	for _, dep := range fd.GetDependencies() ***REMOVED***
		md = r.findMessage(dep, msgName, checked)
		if md != nil ***REMOVED***
			return md
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var _ jsonpb.AnyResolver = (*anyResolver)(nil)
