// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protojson

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"google.golang.org/protobuf/internal/encoding/json"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/set"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Unmarshal reads the given []byte into the given proto.Message.
func Unmarshal(b []byte, m proto.Message) error ***REMOVED***
	return UnmarshalOptions***REMOVED******REMOVED***.Unmarshal(b, m)
***REMOVED***

// UnmarshalOptions is a configurable JSON format parser.
type UnmarshalOptions struct ***REMOVED***
	pragma.NoUnkeyedLiterals

	// If AllowPartial is set, input for messages that will result in missing
	// required fields will not return an error.
	AllowPartial bool

	// If DiscardUnknown is set, unknown fields are ignored.
	DiscardUnknown bool

	// Resolver is used for looking up types when unmarshaling
	// google.protobuf.Any messages or extension fields.
	// If nil, this defaults to using protoregistry.GlobalTypes.
	Resolver interface ***REMOVED***
		protoregistry.MessageTypeResolver
		protoregistry.ExtensionTypeResolver
	***REMOVED***
***REMOVED***

// Unmarshal reads the given []byte and populates the given proto.Message using
// options in UnmarshalOptions object. It will clear the message first before
// setting the fields. If it returns an error, the given message may be
// partially set.
func (o UnmarshalOptions) Unmarshal(b []byte, m proto.Message) error ***REMOVED***
	return o.unmarshal(b, m)
***REMOVED***

// unmarshal is a centralized function that all unmarshal operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for unmarshal that do not go through this.
func (o UnmarshalOptions) unmarshal(b []byte, m proto.Message) error ***REMOVED***
	proto.Reset(m)

	if o.Resolver == nil ***REMOVED***
		o.Resolver = protoregistry.GlobalTypes
	***REMOVED***

	dec := decoder***REMOVED***json.NewDecoder(b), o***REMOVED***
	if err := dec.unmarshalMessage(m.ProtoReflect(), false); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check for EOF.
	tok, err := dec.Read()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tok.Kind() != json.EOF ***REMOVED***
		return dec.unexpectedTokenError(tok)
	***REMOVED***

	if o.AllowPartial ***REMOVED***
		return nil
	***REMOVED***
	return proto.CheckInitialized(m)
***REMOVED***

type decoder struct ***REMOVED***
	*json.Decoder
	opts UnmarshalOptions
***REMOVED***

// newError returns an error object with position info.
func (d decoder) newError(pos int, f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	line, column := d.Position(pos)
	head := fmt.Sprintf("(line %d:%d): ", line, column)
	return errors.New(head+f, x...)
***REMOVED***

// unexpectedTokenError returns a syntax error for the given unexpected token.
func (d decoder) unexpectedTokenError(tok json.Token) error ***REMOVED***
	return d.syntaxError(tok.Pos(), "unexpected token %s", tok.RawString())
***REMOVED***

// syntaxError returns a syntax error for given position.
func (d decoder) syntaxError(pos int, f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	line, column := d.Position(pos)
	head := fmt.Sprintf("syntax error (line %d:%d): ", line, column)
	return errors.New(head+f, x...)
***REMOVED***

// unmarshalMessage unmarshals a message into the given protoreflect.Message.
func (d decoder) unmarshalMessage(m pref.Message, skipTypeURL bool) error ***REMOVED***
	if unmarshal := wellKnownTypeUnmarshaler(m.Descriptor().FullName()); unmarshal != nil ***REMOVED***
		return unmarshal(d, m)
	***REMOVED***

	tok, err := d.Read()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tok.Kind() != json.ObjectOpen ***REMOVED***
		return d.unexpectedTokenError(tok)
	***REMOVED***

	if err := d.unmarshalFields(m, skipTypeURL); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// unmarshalFields unmarshals the fields into the given protoreflect.Message.
func (d decoder) unmarshalFields(m pref.Message, skipTypeURL bool) error ***REMOVED***
	messageDesc := m.Descriptor()
	if !flags.ProtoLegacy && messageset.IsMessageSet(messageDesc) ***REMOVED***
		return errors.New("no support for proto1 MessageSets")
	***REMOVED***

	var seenNums set.Ints
	var seenOneofs set.Ints
	fieldDescs := messageDesc.Fields()
	for ***REMOVED***
		// Read field name.
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch tok.Kind() ***REMOVED***
		default:
			return d.unexpectedTokenError(tok)
		case json.ObjectClose:
			return nil
		case json.Name:
			// Continue below.
		***REMOVED***

		name := tok.Name()
		// Unmarshaling a non-custom embedded message in Any will contain the
		// JSON field "@type" which should be skipped because it is not a field
		// of the embedded message, but simply an artifact of the Any format.
		if skipTypeURL && name == "@type" ***REMOVED***
			d.Read()
			continue
		***REMOVED***

		// Get the FieldDescriptor.
		var fd pref.FieldDescriptor
		if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") ***REMOVED***
			// Only extension names are in [name] format.
			extName := pref.FullName(name[1 : len(name)-1])
			extType, err := d.findExtension(extName)
			if err != nil && err != protoregistry.NotFound ***REMOVED***
				return d.newError(tok.Pos(), "unable to resolve %s: %v", tok.RawString(), err)
			***REMOVED***
			if extType != nil ***REMOVED***
				fd = extType.TypeDescriptor()
				if !messageDesc.ExtensionRanges().Has(fd.Number()) || fd.ContainingMessage().FullName() != messageDesc.FullName() ***REMOVED***
					return d.newError(tok.Pos(), "message %v cannot be extended by %v", messageDesc.FullName(), fd.FullName())
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// The name can either be the JSON name or the proto field name.
			fd = fieldDescs.ByJSONName(name)
			if fd == nil ***REMOVED***
				fd = fieldDescs.ByName(pref.Name(name))
				if fd == nil ***REMOVED***
					// The proto name of a group field is in all lowercase,
					// while the textual field name is the group message name.
					gd := fieldDescs.ByName(pref.Name(strings.ToLower(name)))
					if gd != nil && gd.Kind() == pref.GroupKind && gd.Message().Name() == pref.Name(name) ***REMOVED***
						fd = gd
					***REMOVED***
				***REMOVED*** else if fd.Kind() == pref.GroupKind && fd.Message().Name() != pref.Name(name) ***REMOVED***
					fd = nil // reset since field name is actually the message name
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if flags.ProtoLegacy ***REMOVED***
			if fd != nil && fd.IsWeak() && fd.Message().IsPlaceholder() ***REMOVED***
				fd = nil // reset since the weak reference is not linked in
			***REMOVED***
		***REMOVED***

		if fd == nil ***REMOVED***
			// Field is unknown.
			if d.opts.DiscardUnknown ***REMOVED***
				if err := d.skipJSONValue(); err != nil ***REMOVED***
					return err
				***REMOVED***
				continue
			***REMOVED***
			return d.newError(tok.Pos(), "unknown field %v", tok.RawString())
		***REMOVED***

		// Do not allow duplicate fields.
		num := uint64(fd.Number())
		if seenNums.Has(num) ***REMOVED***
			return d.newError(tok.Pos(), "duplicate field %v", tok.RawString())
		***REMOVED***
		seenNums.Set(num)

		// No need to set values for JSON null unless the field type is
		// google.protobuf.Value or google.protobuf.NullValue.
		if tok, _ := d.Peek(); tok.Kind() == json.Null && !isKnownValue(fd) && !isNullValue(fd) ***REMOVED***
			d.Read()
			continue
		***REMOVED***

		switch ***REMOVED***
		case fd.IsList():
			list := m.Mutable(fd).List()
			if err := d.unmarshalList(list, fd); err != nil ***REMOVED***
				return err
			***REMOVED***
		case fd.IsMap():
			mmap := m.Mutable(fd).Map()
			if err := d.unmarshalMap(mmap, fd); err != nil ***REMOVED***
				return err
			***REMOVED***
		default:
			// If field is a oneof, check if it has already been set.
			if od := fd.ContainingOneof(); od != nil ***REMOVED***
				idx := uint64(od.Index())
				if seenOneofs.Has(idx) ***REMOVED***
					return d.newError(tok.Pos(), "error parsing %s, oneof %v is already set", tok.RawString(), od.FullName())
				***REMOVED***
				seenOneofs.Set(idx)
			***REMOVED***

			// Required or optional fields.
			if err := d.unmarshalSingular(m, fd); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// findExtension returns protoreflect.ExtensionType from the resolver if found.
func (d decoder) findExtension(xtName pref.FullName) (pref.ExtensionType, error) ***REMOVED***
	xt, err := d.opts.Resolver.FindExtensionByName(xtName)
	if err == nil ***REMOVED***
		return xt, nil
	***REMOVED***
	return messageset.FindMessageSetExtension(d.opts.Resolver, xtName)
***REMOVED***

func isKnownValue(fd pref.FieldDescriptor) bool ***REMOVED***
	md := fd.Message()
	return md != nil && md.FullName() == genid.Value_message_fullname
***REMOVED***

func isNullValue(fd pref.FieldDescriptor) bool ***REMOVED***
	ed := fd.Enum()
	return ed != nil && ed.FullName() == genid.NullValue_enum_fullname
***REMOVED***

// unmarshalSingular unmarshals to the non-repeated field specified
// by the given FieldDescriptor.
func (d decoder) unmarshalSingular(m pref.Message, fd pref.FieldDescriptor) error ***REMOVED***
	var val pref.Value
	var err error
	switch fd.Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind:
		val = m.NewField(fd)
		err = d.unmarshalMessage(val.Message(), false)
	default:
		val, err = d.unmarshalScalar(fd)
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m.Set(fd, val)
	return nil
***REMOVED***

// unmarshalScalar unmarshals to a scalar/enum protoreflect.Value specified by
// the given FieldDescriptor.
func (d decoder) unmarshalScalar(fd pref.FieldDescriptor) (pref.Value, error) ***REMOVED***
	const b32 int = 32
	const b64 int = 64

	tok, err := d.Read()
	if err != nil ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, err
	***REMOVED***

	kind := fd.Kind()
	switch kind ***REMOVED***
	case pref.BoolKind:
		if tok.Kind() == json.Bool ***REMOVED***
			return pref.ValueOfBool(tok.Bool()), nil
		***REMOVED***

	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		if v, ok := unmarshalInt(tok, b32); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		if v, ok := unmarshalInt(tok, b64); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.Uint32Kind, pref.Fixed32Kind:
		if v, ok := unmarshalUint(tok, b32); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.Uint64Kind, pref.Fixed64Kind:
		if v, ok := unmarshalUint(tok, b64); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.FloatKind:
		if v, ok := unmarshalFloat(tok, b32); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.DoubleKind:
		if v, ok := unmarshalFloat(tok, b64); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.StringKind:
		if tok.Kind() == json.String ***REMOVED***
			return pref.ValueOfString(tok.ParsedString()), nil
		***REMOVED***

	case pref.BytesKind:
		if v, ok := unmarshalBytes(tok); ok ***REMOVED***
			return v, nil
		***REMOVED***

	case pref.EnumKind:
		if v, ok := unmarshalEnum(tok, fd); ok ***REMOVED***
			return v, nil
		***REMOVED***

	default:
		panic(fmt.Sprintf("unmarshalScalar: invalid scalar kind %v", kind))
	***REMOVED***

	return pref.Value***REMOVED******REMOVED***, d.newError(tok.Pos(), "invalid value for %v type: %v", kind, tok.RawString())
***REMOVED***

func unmarshalInt(tok json.Token, bitSize int) (pref.Value, bool) ***REMOVED***
	switch tok.Kind() ***REMOVED***
	case json.Number:
		return getInt(tok, bitSize)

	case json.String:
		// Decode number from string.
		s := strings.TrimSpace(tok.ParsedString())
		if len(s) != len(tok.ParsedString()) ***REMOVED***
			return pref.Value***REMOVED******REMOVED***, false
		***REMOVED***
		dec := json.NewDecoder([]byte(s))
		tok, err := dec.Read()
		if err != nil ***REMOVED***
			return pref.Value***REMOVED******REMOVED***, false
		***REMOVED***
		return getInt(tok, bitSize)
	***REMOVED***
	return pref.Value***REMOVED******REMOVED***, false
***REMOVED***

func getInt(tok json.Token, bitSize int) (pref.Value, bool) ***REMOVED***
	n, ok := tok.Int(bitSize)
	if !ok ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, false
	***REMOVED***
	if bitSize == 32 ***REMOVED***
		return pref.ValueOfInt32(int32(n)), true
	***REMOVED***
	return pref.ValueOfInt64(n), true
***REMOVED***

func unmarshalUint(tok json.Token, bitSize int) (pref.Value, bool) ***REMOVED***
	switch tok.Kind() ***REMOVED***
	case json.Number:
		return getUint(tok, bitSize)

	case json.String:
		// Decode number from string.
		s := strings.TrimSpace(tok.ParsedString())
		if len(s) != len(tok.ParsedString()) ***REMOVED***
			return pref.Value***REMOVED******REMOVED***, false
		***REMOVED***
		dec := json.NewDecoder([]byte(s))
		tok, err := dec.Read()
		if err != nil ***REMOVED***
			return pref.Value***REMOVED******REMOVED***, false
		***REMOVED***
		return getUint(tok, bitSize)
	***REMOVED***
	return pref.Value***REMOVED******REMOVED***, false
***REMOVED***

func getUint(tok json.Token, bitSize int) (pref.Value, bool) ***REMOVED***
	n, ok := tok.Uint(bitSize)
	if !ok ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, false
	***REMOVED***
	if bitSize == 32 ***REMOVED***
		return pref.ValueOfUint32(uint32(n)), true
	***REMOVED***
	return pref.ValueOfUint64(n), true
***REMOVED***

func unmarshalFloat(tok json.Token, bitSize int) (pref.Value, bool) ***REMOVED***
	switch tok.Kind() ***REMOVED***
	case json.Number:
		return getFloat(tok, bitSize)

	case json.String:
		s := tok.ParsedString()
		switch s ***REMOVED***
		case "NaN":
			if bitSize == 32 ***REMOVED***
				return pref.ValueOfFloat32(float32(math.NaN())), true
			***REMOVED***
			return pref.ValueOfFloat64(math.NaN()), true
		case "Infinity":
			if bitSize == 32 ***REMOVED***
				return pref.ValueOfFloat32(float32(math.Inf(+1))), true
			***REMOVED***
			return pref.ValueOfFloat64(math.Inf(+1)), true
		case "-Infinity":
			if bitSize == 32 ***REMOVED***
				return pref.ValueOfFloat32(float32(math.Inf(-1))), true
			***REMOVED***
			return pref.ValueOfFloat64(math.Inf(-1)), true
		***REMOVED***

		// Decode number from string.
		if len(s) != len(strings.TrimSpace(s)) ***REMOVED***
			return pref.Value***REMOVED******REMOVED***, false
		***REMOVED***
		dec := json.NewDecoder([]byte(s))
		tok, err := dec.Read()
		if err != nil ***REMOVED***
			return pref.Value***REMOVED******REMOVED***, false
		***REMOVED***
		return getFloat(tok, bitSize)
	***REMOVED***
	return pref.Value***REMOVED******REMOVED***, false
***REMOVED***

func getFloat(tok json.Token, bitSize int) (pref.Value, bool) ***REMOVED***
	n, ok := tok.Float(bitSize)
	if !ok ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, false
	***REMOVED***
	if bitSize == 32 ***REMOVED***
		return pref.ValueOfFloat32(float32(n)), true
	***REMOVED***
	return pref.ValueOfFloat64(n), true
***REMOVED***

func unmarshalBytes(tok json.Token) (pref.Value, bool) ***REMOVED***
	if tok.Kind() != json.String ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, false
	***REMOVED***

	s := tok.ParsedString()
	enc := base64.StdEncoding
	if strings.ContainsAny(s, "-_") ***REMOVED***
		enc = base64.URLEncoding
	***REMOVED***
	if len(s)%4 != 0 ***REMOVED***
		enc = enc.WithPadding(base64.NoPadding)
	***REMOVED***
	b, err := enc.DecodeString(s)
	if err != nil ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, false
	***REMOVED***
	return pref.ValueOfBytes(b), true
***REMOVED***

func unmarshalEnum(tok json.Token, fd pref.FieldDescriptor) (pref.Value, bool) ***REMOVED***
	switch tok.Kind() ***REMOVED***
	case json.String:
		// Lookup EnumNumber based on name.
		s := tok.ParsedString()
		if enumVal := fd.Enum().Values().ByName(pref.Name(s)); enumVal != nil ***REMOVED***
			return pref.ValueOfEnum(enumVal.Number()), true
		***REMOVED***

	case json.Number:
		if n, ok := tok.Int(32); ok ***REMOVED***
			return pref.ValueOfEnum(pref.EnumNumber(n)), true
		***REMOVED***

	case json.Null:
		// This is only valid for google.protobuf.NullValue.
		if isNullValue(fd) ***REMOVED***
			return pref.ValueOfEnum(0), true
		***REMOVED***
	***REMOVED***

	return pref.Value***REMOVED******REMOVED***, false
***REMOVED***

func (d decoder) unmarshalList(list pref.List, fd pref.FieldDescriptor) error ***REMOVED***
	tok, err := d.Read()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tok.Kind() != json.ArrayOpen ***REMOVED***
		return d.unexpectedTokenError(tok)
	***REMOVED***

	switch fd.Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind:
		for ***REMOVED***
			tok, err := d.Peek()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if tok.Kind() == json.ArrayClose ***REMOVED***
				d.Read()
				return nil
			***REMOVED***

			val := list.NewElement()
			if err := d.unmarshalMessage(val.Message(), false); err != nil ***REMOVED***
				return err
			***REMOVED***
			list.Append(val)
		***REMOVED***
	default:
		for ***REMOVED***
			tok, err := d.Peek()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if tok.Kind() == json.ArrayClose ***REMOVED***
				d.Read()
				return nil
			***REMOVED***

			val, err := d.unmarshalScalar(fd)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			list.Append(val)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d decoder) unmarshalMap(mmap pref.Map, fd pref.FieldDescriptor) error ***REMOVED***
	tok, err := d.Read()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tok.Kind() != json.ObjectOpen ***REMOVED***
		return d.unexpectedTokenError(tok)
	***REMOVED***

	// Determine ahead whether map entry is a scalar type or a message type in
	// order to call the appropriate unmarshalMapValue func inside the for loop
	// below.
	var unmarshalMapValue func() (pref.Value, error)
	switch fd.MapValue().Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind:
		unmarshalMapValue = func() (pref.Value, error) ***REMOVED***
			val := mmap.NewValue()
			if err := d.unmarshalMessage(val.Message(), false); err != nil ***REMOVED***
				return pref.Value***REMOVED******REMOVED***, err
			***REMOVED***
			return val, nil
		***REMOVED***
	default:
		unmarshalMapValue = func() (pref.Value, error) ***REMOVED***
			return d.unmarshalScalar(fd.MapValue())
		***REMOVED***
	***REMOVED***

Loop:
	for ***REMOVED***
		// Read field name.
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch tok.Kind() ***REMOVED***
		default:
			return d.unexpectedTokenError(tok)
		case json.ObjectClose:
			break Loop
		case json.Name:
			// Continue.
		***REMOVED***

		// Unmarshal field name.
		pkey, err := d.unmarshalMapKey(tok, fd.MapKey())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Check for duplicate field name.
		if mmap.Has(pkey) ***REMOVED***
			return d.newError(tok.Pos(), "duplicate map key %v", tok.RawString())
		***REMOVED***

		// Read and unmarshal field value.
		pval, err := unmarshalMapValue()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		mmap.Set(pkey, pval)
	***REMOVED***

	return nil
***REMOVED***

// unmarshalMapKey converts given token of Name kind into a protoreflect.MapKey.
// A map key type is any integral or string type.
func (d decoder) unmarshalMapKey(tok json.Token, fd pref.FieldDescriptor) (pref.MapKey, error) ***REMOVED***
	const b32 = 32
	const b64 = 64
	const base10 = 10

	name := tok.Name()
	kind := fd.Kind()
	switch kind ***REMOVED***
	case pref.StringKind:
		return pref.ValueOfString(name).MapKey(), nil

	case pref.BoolKind:
		switch name ***REMOVED***
		case "true":
			return pref.ValueOfBool(true).MapKey(), nil
		case "false":
			return pref.ValueOfBool(false).MapKey(), nil
		***REMOVED***

	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		if n, err := strconv.ParseInt(name, base10, b32); err == nil ***REMOVED***
			return pref.ValueOfInt32(int32(n)).MapKey(), nil
		***REMOVED***

	case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		if n, err := strconv.ParseInt(name, base10, b64); err == nil ***REMOVED***
			return pref.ValueOfInt64(int64(n)).MapKey(), nil
		***REMOVED***

	case pref.Uint32Kind, pref.Fixed32Kind:
		if n, err := strconv.ParseUint(name, base10, b32); err == nil ***REMOVED***
			return pref.ValueOfUint32(uint32(n)).MapKey(), nil
		***REMOVED***

	case pref.Uint64Kind, pref.Fixed64Kind:
		if n, err := strconv.ParseUint(name, base10, b64); err == nil ***REMOVED***
			return pref.ValueOfUint64(uint64(n)).MapKey(), nil
		***REMOVED***

	default:
		panic(fmt.Sprintf("invalid kind for map key: %v", kind))
	***REMOVED***

	return pref.MapKey***REMOVED******REMOVED***, d.newError(tok.Pos(), "invalid value for %v key: %s", kind, tok.RawString())
***REMOVED***
