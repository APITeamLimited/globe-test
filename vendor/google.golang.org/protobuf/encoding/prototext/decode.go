// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prototext

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/encoding/text"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/fieldnum"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/set"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Unmarshal reads the given []byte into the given proto.Message.
func Unmarshal(b []byte, m proto.Message) error ***REMOVED***
	return UnmarshalOptions***REMOVED******REMOVED***.Unmarshal(b, m)
***REMOVED***

// UnmarshalOptions is a configurable textproto format unmarshaler.
type UnmarshalOptions struct ***REMOVED***
	pragma.NoUnkeyedLiterals

	// AllowPartial accepts input for messages that will result in missing
	// required fields. If AllowPartial is false (the default), Unmarshal will
	// return error if there are any missing required fields.
	AllowPartial bool

	// DiscardUnknown specifies whether to ignore unknown fields when parsing.
	// An unknown field is any field whose field name or field number does not
	// resolve to any known or extension field in the message.
	// By default, unmarshal rejects unknown fields as an error.
	DiscardUnknown bool

	// Resolver is used for looking up types when unmarshaling
	// google.protobuf.Any messages or extension fields.
	// If nil, this defaults to using protoregistry.GlobalTypes.
	Resolver interface ***REMOVED***
		protoregistry.MessageTypeResolver
		protoregistry.ExtensionTypeResolver
	***REMOVED***
***REMOVED***

// Unmarshal reads the given []byte and populates the given proto.Message using options in
// UnmarshalOptions object.
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

	dec := decoder***REMOVED***text.NewDecoder(b), o***REMOVED***
	if err := dec.unmarshalMessage(m.ProtoReflect(), false); err != nil ***REMOVED***
		return err
	***REMOVED***
	if o.AllowPartial ***REMOVED***
		return nil
	***REMOVED***
	return proto.CheckInitialized(m)
***REMOVED***

type decoder struct ***REMOVED***
	*text.Decoder
	opts UnmarshalOptions
***REMOVED***

// newError returns an error object with position info.
func (d decoder) newError(pos int, f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	line, column := d.Position(pos)
	head := fmt.Sprintf("(line %d:%d): ", line, column)
	return errors.New(head+f, x...)
***REMOVED***

// unexpectedTokenError returns a syntax error for the given unexpected token.
func (d decoder) unexpectedTokenError(tok text.Token) error ***REMOVED***
	return d.syntaxError(tok.Pos(), "unexpected token: %s", tok.RawString())
***REMOVED***

// syntaxError returns a syntax error for given position.
func (d decoder) syntaxError(pos int, f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	line, column := d.Position(pos)
	head := fmt.Sprintf("syntax error (line %d:%d): ", line, column)
	return errors.New(head+f, x...)
***REMOVED***

// unmarshalMessage unmarshals into the given protoreflect.Message.
func (d decoder) unmarshalMessage(m pref.Message, checkDelims bool) error ***REMOVED***
	messageDesc := m.Descriptor()
	if !flags.ProtoLegacy && messageset.IsMessageSet(messageDesc) ***REMOVED***
		return errors.New("no support for proto1 MessageSets")
	***REMOVED***

	if messageDesc.FullName() == "google.protobuf.Any" ***REMOVED***
		return d.unmarshalAny(m, checkDelims)
	***REMOVED***

	if checkDelims ***REMOVED***
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if tok.Kind() != text.MessageOpen ***REMOVED***
			return d.unexpectedTokenError(tok)
		***REMOVED***
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
		switch typ := tok.Kind(); typ ***REMOVED***
		case text.Name:
			// Continue below.
		case text.EOF:
			if checkDelims ***REMOVED***
				return text.ErrUnexpectedEOF
			***REMOVED***
			return nil
		default:
			if checkDelims && typ == text.MessageClose ***REMOVED***
				return nil
			***REMOVED***
			return d.unexpectedTokenError(tok)
		***REMOVED***

		// Resolve the field descriptor.
		var name pref.Name
		var fd pref.FieldDescriptor
		var xt pref.ExtensionType
		var xtErr error
		var isFieldNumberName bool

		switch tok.NameKind() ***REMOVED***
		case text.IdentName:
			name = pref.Name(tok.IdentName())
			fd = fieldDescs.ByName(name)
			if fd == nil ***REMOVED***
				// The proto name of a group field is in all lowercase,
				// while the textproto field name is the group message name.
				gd := fieldDescs.ByName(pref.Name(strings.ToLower(string(name))))
				if gd != nil && gd.Kind() == pref.GroupKind && gd.Message().Name() == name ***REMOVED***
					fd = gd
				***REMOVED***
			***REMOVED*** else if fd.Kind() == pref.GroupKind && fd.Message().Name() != name ***REMOVED***
				fd = nil // reset since field name is actually the message name
			***REMOVED***

		case text.TypeName:
			// Handle extensions only. This code path is not for Any.
			xt, xtErr = d.findExtension(pref.FullName(tok.TypeName()))

		case text.FieldNumber:
			isFieldNumberName = true
			num := pref.FieldNumber(tok.FieldNumber())
			if !num.IsValid() ***REMOVED***
				return d.newError(tok.Pos(), "invalid field number: %d", num)
			***REMOVED***
			fd = fieldDescs.ByNumber(num)
			if fd == nil ***REMOVED***
				xt, xtErr = d.opts.Resolver.FindExtensionByNumber(messageDesc.FullName(), num)
			***REMOVED***
		***REMOVED***

		if xt != nil ***REMOVED***
			fd = xt.TypeDescriptor()
			if !messageDesc.ExtensionRanges().Has(fd.Number()) || fd.ContainingMessage().FullName() != messageDesc.FullName() ***REMOVED***
				return d.newError(tok.Pos(), "message %v cannot be extended by %v", messageDesc.FullName(), fd.FullName())
			***REMOVED***
		***REMOVED*** else if xtErr != nil && xtErr != protoregistry.NotFound ***REMOVED***
			return d.newError(tok.Pos(), "unable to resolve [%s]: %v", tok.RawString(), xtErr)
		***REMOVED***
		if flags.ProtoLegacy ***REMOVED***
			if fd != nil && fd.IsWeak() && fd.Message().IsPlaceholder() ***REMOVED***
				fd = nil // reset since the weak reference is not linked in
			***REMOVED***
		***REMOVED***

		// Handle unknown fields.
		if fd == nil ***REMOVED***
			if d.opts.DiscardUnknown || messageDesc.ReservedNames().Has(name) ***REMOVED***
				d.skipValue()
				continue
			***REMOVED***
			return d.newError(tok.Pos(), "unknown field: %v", tok.RawString())
		***REMOVED***

		// Handle fields identified by field number.
		if isFieldNumberName ***REMOVED***
			// TODO: Add an option to permit parsing field numbers.
			//
			// This requires careful thought as the MarshalOptions.EmitUnknown
			// option allows formatting unknown fields as the field number and the
			// best-effort textual representation of the field value.  In that case,
			// it may not be possible to unmarshal the value from a parser that does
			// have information about the unknown field.
			return d.newError(tok.Pos(), "cannot specify field by number: %v", tok.RawString())
		***REMOVED***

		switch ***REMOVED***
		case fd.IsList():
			kind := fd.Kind()
			if kind != pref.MessageKind && kind != pref.GroupKind && !tok.HasSeparator() ***REMOVED***
				return d.syntaxError(tok.Pos(), "missing field separator :")
			***REMOVED***

			list := m.Mutable(fd).List()
			if err := d.unmarshalList(fd, list); err != nil ***REMOVED***
				return err
			***REMOVED***

		case fd.IsMap():
			mmap := m.Mutable(fd).Map()
			if err := d.unmarshalMap(fd, mmap); err != nil ***REMOVED***
				return err
			***REMOVED***

		default:
			kind := fd.Kind()
			if kind != pref.MessageKind && kind != pref.GroupKind && !tok.HasSeparator() ***REMOVED***
				return d.syntaxError(tok.Pos(), "missing field separator :")
			***REMOVED***

			// If field is a oneof, check if it has already been set.
			if od := fd.ContainingOneof(); od != nil ***REMOVED***
				idx := uint64(od.Index())
				if seenOneofs.Has(idx) ***REMOVED***
					return d.newError(tok.Pos(), "error parsing %q, oneof %v is already set", tok.RawString(), od.FullName())
				***REMOVED***
				seenOneofs.Set(idx)
			***REMOVED***

			num := uint64(fd.Number())
			if seenNums.Has(num) ***REMOVED***
				return d.newError(tok.Pos(), "non-repeated field %q is repeated", tok.RawString())
			***REMOVED***

			if err := d.unmarshalSingular(fd, m); err != nil ***REMOVED***
				return err
			***REMOVED***
			seenNums.Set(num)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// findExtension returns protoreflect.ExtensionType from the Resolver if found.
func (d decoder) findExtension(xtName pref.FullName) (pref.ExtensionType, error) ***REMOVED***
	xt, err := d.opts.Resolver.FindExtensionByName(xtName)
	if err == nil ***REMOVED***
		return xt, nil
	***REMOVED***
	return messageset.FindMessageSetExtension(d.opts.Resolver, xtName)
***REMOVED***

// unmarshalSingular unmarshals a non-repeated field value specified by the
// given FieldDescriptor.
func (d decoder) unmarshalSingular(fd pref.FieldDescriptor, m pref.Message) error ***REMOVED***
	var val pref.Value
	var err error
	switch fd.Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind:
		val = m.NewField(fd)
		err = d.unmarshalMessage(val.Message(), true)
	default:
		val, err = d.unmarshalScalar(fd)
	***REMOVED***
	if err == nil ***REMOVED***
		m.Set(fd, val)
	***REMOVED***
	return err
***REMOVED***

// unmarshalScalar unmarshals a scalar/enum protoreflect.Value specified by the
// given FieldDescriptor.
func (d decoder) unmarshalScalar(fd pref.FieldDescriptor) (pref.Value, error) ***REMOVED***
	tok, err := d.Read()
	if err != nil ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, err
	***REMOVED***

	if tok.Kind() != text.Scalar ***REMOVED***
		return pref.Value***REMOVED******REMOVED***, d.unexpectedTokenError(tok)
	***REMOVED***

	kind := fd.Kind()
	switch kind ***REMOVED***
	case pref.BoolKind:
		if b, ok := tok.Bool(); ok ***REMOVED***
			return pref.ValueOfBool(b), nil
		***REMOVED***

	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		if n, ok := tok.Int32(); ok ***REMOVED***
			return pref.ValueOfInt32(n), nil
		***REMOVED***

	case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		if n, ok := tok.Int64(); ok ***REMOVED***
			return pref.ValueOfInt64(n), nil
		***REMOVED***

	case pref.Uint32Kind, pref.Fixed32Kind:
		if n, ok := tok.Uint32(); ok ***REMOVED***
			return pref.ValueOfUint32(n), nil
		***REMOVED***

	case pref.Uint64Kind, pref.Fixed64Kind:
		if n, ok := tok.Uint64(); ok ***REMOVED***
			return pref.ValueOfUint64(n), nil
		***REMOVED***

	case pref.FloatKind:
		if n, ok := tok.Float32(); ok ***REMOVED***
			return pref.ValueOfFloat32(n), nil
		***REMOVED***

	case pref.DoubleKind:
		if n, ok := tok.Float64(); ok ***REMOVED***
			return pref.ValueOfFloat64(n), nil
		***REMOVED***

	case pref.StringKind:
		if s, ok := tok.String(); ok ***REMOVED***
			if strs.EnforceUTF8(fd) && !utf8.ValidString(s) ***REMOVED***
				return pref.Value***REMOVED******REMOVED***, d.newError(tok.Pos(), "contains invalid UTF-8")
			***REMOVED***
			return pref.ValueOfString(s), nil
		***REMOVED***

	case pref.BytesKind:
		if b, ok := tok.String(); ok ***REMOVED***
			return pref.ValueOfBytes([]byte(b)), nil
		***REMOVED***

	case pref.EnumKind:
		if lit, ok := tok.Enum(); ok ***REMOVED***
			// Lookup EnumNumber based on name.
			if enumVal := fd.Enum().Values().ByName(pref.Name(lit)); enumVal != nil ***REMOVED***
				return pref.ValueOfEnum(enumVal.Number()), nil
			***REMOVED***
		***REMOVED***
		if num, ok := tok.Int32(); ok ***REMOVED***
			return pref.ValueOfEnum(pref.EnumNumber(num)), nil
		***REMOVED***

	default:
		panic(fmt.Sprintf("invalid scalar kind %v", kind))
	***REMOVED***

	return pref.Value***REMOVED******REMOVED***, d.newError(tok.Pos(), "invalid value for %v type: %v", kind, tok.RawString())
***REMOVED***

// unmarshalList unmarshals into given protoreflect.List. A list value can
// either be in [] syntax or simply just a single scalar/message value.
func (d decoder) unmarshalList(fd pref.FieldDescriptor, list pref.List) error ***REMOVED***
	tok, err := d.Peek()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch fd.Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind:
		switch tok.Kind() ***REMOVED***
		case text.ListOpen:
			d.Read()
			for ***REMOVED***
				tok, err := d.Peek()
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				switch tok.Kind() ***REMOVED***
				case text.ListClose:
					d.Read()
					return nil
				case text.MessageOpen:
					pval := list.NewElement()
					if err := d.unmarshalMessage(pval.Message(), true); err != nil ***REMOVED***
						return err
					***REMOVED***
					list.Append(pval)
				default:
					return d.unexpectedTokenError(tok)
				***REMOVED***
			***REMOVED***

		case text.MessageOpen:
			pval := list.NewElement()
			if err := d.unmarshalMessage(pval.Message(), true); err != nil ***REMOVED***
				return err
			***REMOVED***
			list.Append(pval)
			return nil
		***REMOVED***

	default:
		switch tok.Kind() ***REMOVED***
		case text.ListOpen:
			d.Read()
			for ***REMOVED***
				tok, err := d.Peek()
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				switch tok.Kind() ***REMOVED***
				case text.ListClose:
					d.Read()
					return nil
				case text.Scalar:
					pval, err := d.unmarshalScalar(fd)
					if err != nil ***REMOVED***
						return err
					***REMOVED***
					list.Append(pval)
				default:
					return d.unexpectedTokenError(tok)
				***REMOVED***
			***REMOVED***

		case text.Scalar:
			pval, err := d.unmarshalScalar(fd)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			list.Append(pval)
			return nil
		***REMOVED***
	***REMOVED***

	return d.unexpectedTokenError(tok)
***REMOVED***

// unmarshalMap unmarshals into given protoreflect.Map. A map value is a
// textproto message containing ***REMOVED***key: <kvalue>, value: <mvalue>***REMOVED***.
func (d decoder) unmarshalMap(fd pref.FieldDescriptor, mmap pref.Map) error ***REMOVED***
	// Determine ahead whether map entry is a scalar type or a message type in
	// order to call the appropriate unmarshalMapValue func inside
	// unmarshalMapEntry.
	var unmarshalMapValue func() (pref.Value, error)
	switch fd.MapValue().Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind:
		unmarshalMapValue = func() (pref.Value, error) ***REMOVED***
			pval := mmap.NewValue()
			if err := d.unmarshalMessage(pval.Message(), true); err != nil ***REMOVED***
				return pref.Value***REMOVED******REMOVED***, err
			***REMOVED***
			return pval, nil
		***REMOVED***
	default:
		unmarshalMapValue = func() (pref.Value, error) ***REMOVED***
			return d.unmarshalScalar(fd.MapValue())
		***REMOVED***
	***REMOVED***

	tok, err := d.Read()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch tok.Kind() ***REMOVED***
	case text.MessageOpen:
		return d.unmarshalMapEntry(fd, mmap, unmarshalMapValue)

	case text.ListOpen:
		for ***REMOVED***
			tok, err := d.Read()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			switch tok.Kind() ***REMOVED***
			case text.ListClose:
				return nil
			case text.MessageOpen:
				if err := d.unmarshalMapEntry(fd, mmap, unmarshalMapValue); err != nil ***REMOVED***
					return err
				***REMOVED***
			default:
				return d.unexpectedTokenError(tok)
			***REMOVED***
		***REMOVED***

	default:
		return d.unexpectedTokenError(tok)
	***REMOVED***
***REMOVED***

// unmarshalMap unmarshals into given protoreflect.Map. A map value is a
// textproto message containing ***REMOVED***key: <kvalue>, value: <mvalue>***REMOVED***.
func (d decoder) unmarshalMapEntry(fd pref.FieldDescriptor, mmap pref.Map, unmarshalMapValue func() (pref.Value, error)) error ***REMOVED***
	var key pref.MapKey
	var pval pref.Value
Loop:
	for ***REMOVED***
		// Read field name.
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch tok.Kind() ***REMOVED***
		case text.Name:
			if tok.NameKind() != text.IdentName ***REMOVED***
				if !d.opts.DiscardUnknown ***REMOVED***
					return d.newError(tok.Pos(), "unknown map entry field %q", tok.RawString())
				***REMOVED***
				d.skipValue()
				continue Loop
			***REMOVED***
			// Continue below.
		case text.MessageClose:
			break Loop
		default:
			return d.unexpectedTokenError(tok)
		***REMOVED***

		name := tok.IdentName()
		switch name ***REMOVED***
		case "key":
			if !tok.HasSeparator() ***REMOVED***
				return d.syntaxError(tok.Pos(), "missing field separator :")
			***REMOVED***
			if key.IsValid() ***REMOVED***
				return d.newError(tok.Pos(), `map entry "key" cannot be repeated`)
			***REMOVED***
			val, err := d.unmarshalScalar(fd.MapKey())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			key = val.MapKey()

		case "value":
			if kind := fd.MapValue().Kind(); (kind != pref.MessageKind) && (kind != pref.GroupKind) ***REMOVED***
				if !tok.HasSeparator() ***REMOVED***
					return d.syntaxError(tok.Pos(), "missing field separator :")
				***REMOVED***
			***REMOVED***
			if pval.IsValid() ***REMOVED***
				return d.newError(tok.Pos(), `map entry "value" cannot be repeated`)
			***REMOVED***
			pval, err = unmarshalMapValue()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

		default:
			if !d.opts.DiscardUnknown ***REMOVED***
				return d.newError(tok.Pos(), "unknown map entry field %q", name)
			***REMOVED***
			d.skipValue()
		***REMOVED***
	***REMOVED***

	if !key.IsValid() ***REMOVED***
		key = fd.MapKey().Default().MapKey()
	***REMOVED***
	if !pval.IsValid() ***REMOVED***
		switch fd.MapValue().Kind() ***REMOVED***
		case pref.MessageKind, pref.GroupKind:
			// If value field is not set for message/group types, construct an
			// empty one as default.
			pval = mmap.NewValue()
		default:
			pval = fd.MapValue().Default()
		***REMOVED***
	***REMOVED***
	mmap.Set(key, pval)
	return nil
***REMOVED***

// unmarshalAny unmarshals an Any textproto. It can either be in expanded form
// or non-expanded form.
func (d decoder) unmarshalAny(m pref.Message, checkDelims bool) error ***REMOVED***
	var typeURL string
	var bValue []byte

	// hasFields tracks which valid fields have been seen in the loop below in
	// order to flag an error if there are duplicates or conflicts. It may
	// contain the strings "type_url", "value" and "expanded".  The literal
	// "expanded" is used to indicate that the expanded form has been
	// encountered already.
	hasFields := map[string]bool***REMOVED******REMOVED***

	if checkDelims ***REMOVED***
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if tok.Kind() != text.MessageOpen ***REMOVED***
			return d.unexpectedTokenError(tok)
		***REMOVED***
	***REMOVED***

Loop:
	for ***REMOVED***
		// Read field name. Can only have 3 possible field names, i.e. type_url,
		// value and type URL name inside [].
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if typ := tok.Kind(); typ != text.Name ***REMOVED***
			if checkDelims ***REMOVED***
				if typ == text.MessageClose ***REMOVED***
					break Loop
				***REMOVED***
			***REMOVED*** else if typ == text.EOF ***REMOVED***
				break Loop
			***REMOVED***
			return d.unexpectedTokenError(tok)
		***REMOVED***

		switch tok.NameKind() ***REMOVED***
		case text.IdentName:
			// Both type_url and value fields require field separator :.
			if !tok.HasSeparator() ***REMOVED***
				return d.syntaxError(tok.Pos(), "missing field separator :")
			***REMOVED***

			switch tok.IdentName() ***REMOVED***
			case "type_url":
				if hasFields["type_url"] ***REMOVED***
					return d.newError(tok.Pos(), "duplicate Any type_url field")
				***REMOVED***
				if hasFields["expanded"] ***REMOVED***
					return d.newError(tok.Pos(), "conflict with [%s] field", typeURL)
				***REMOVED***
				tok, err := d.Read()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				var ok bool
				typeURL, ok = tok.String()
				if !ok ***REMOVED***
					return d.newError(tok.Pos(), "invalid Any type_url: %v", tok.RawString())
				***REMOVED***
				hasFields["type_url"] = true

			case "value":
				if hasFields["value"] ***REMOVED***
					return d.newError(tok.Pos(), "duplicate Any value field")
				***REMOVED***
				if hasFields["expanded"] ***REMOVED***
					return d.newError(tok.Pos(), "conflict with [%s] field", typeURL)
				***REMOVED***
				tok, err := d.Read()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				s, ok := tok.String()
				if !ok ***REMOVED***
					return d.newError(tok.Pos(), "invalid Any value: %v", tok.RawString())
				***REMOVED***
				bValue = []byte(s)
				hasFields["value"] = true

			default:
				if !d.opts.DiscardUnknown ***REMOVED***
					return d.newError(tok.Pos(), "invalid field name %q in google.protobuf.Any message", tok.RawString())
				***REMOVED***
			***REMOVED***

		case text.TypeName:
			if hasFields["expanded"] ***REMOVED***
				return d.newError(tok.Pos(), "cannot have more than one type")
			***REMOVED***
			if hasFields["type_url"] ***REMOVED***
				return d.newError(tok.Pos(), "conflict with type_url field")
			***REMOVED***
			typeURL = tok.TypeName()
			var err error
			bValue, err = d.unmarshalExpandedAny(typeURL, tok.Pos())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hasFields["expanded"] = true

		default:
			if !d.opts.DiscardUnknown ***REMOVED***
				return d.newError(tok.Pos(), "invalid field name %q in google.protobuf.Any message", tok.RawString())
			***REMOVED***
		***REMOVED***
	***REMOVED***

	fds := m.Descriptor().Fields()
	if len(typeURL) > 0 ***REMOVED***
		m.Set(fds.ByNumber(fieldnum.Any_TypeUrl), pref.ValueOfString(typeURL))
	***REMOVED***
	if len(bValue) > 0 ***REMOVED***
		m.Set(fds.ByNumber(fieldnum.Any_Value), pref.ValueOfBytes(bValue))
	***REMOVED***
	return nil
***REMOVED***

func (d decoder) unmarshalExpandedAny(typeURL string, pos int) ([]byte, error) ***REMOVED***
	mt, err := d.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil ***REMOVED***
		return nil, d.newError(pos, "unable to resolve message [%v]: %v", typeURL, err)
	***REMOVED***
	// Create new message for the embedded message type and unmarshal the value
	// field into it.
	m := mt.New()
	if err := d.unmarshalMessage(m, true); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Serialize the embedded message and return the resulting bytes.
	b, err := proto.MarshalOptions***REMOVED***
		AllowPartial:  true, // Never check required fields inside an Any.
		Deterministic: true,
	***REMOVED***.Marshal(m.Interface())
	if err != nil ***REMOVED***
		return nil, d.newError(pos, "error in marshaling message into Any.value: %v", err)
	***REMOVED***
	return b, nil
***REMOVED***

// skipValue makes the decoder parse a field value in order to advance the read
// to the next field. It relies on Read returning an error if the types are not
// in valid sequence.
func (d decoder) skipValue() error ***REMOVED***
	tok, err := d.Read()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Only need to continue reading for messages and lists.
	switch tok.Kind() ***REMOVED***
	case text.MessageOpen:
		return d.skipMessageValue()

	case text.ListOpen:
		for ***REMOVED***
			tok, err := d.Read()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			switch tok.Kind() ***REMOVED***
			case text.ListClose:
				return nil
			case text.MessageOpen:
				return d.skipMessageValue()
			default:
				// Skip items. This will not validate whether skipped values are
				// of the same type or not, same behavior as C++
				// TextFormat::Parser::AllowUnknownField(true) version 3.8.0.
				if err := d.skipValue(); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// skipMessageValue makes the decoder parse and skip over all fields in a
// message. It assumes that the previous read type is MessageOpen.
func (d decoder) skipMessageValue() error ***REMOVED***
	for ***REMOVED***
		tok, err := d.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch tok.Kind() ***REMOVED***
		case text.MessageClose:
			return nil
		case text.Name:
			if err := d.skipValue(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
