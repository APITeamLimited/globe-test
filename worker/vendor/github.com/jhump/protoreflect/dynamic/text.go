package dynamic

// Marshalling and unmarshalling of dynamic messages to/from proto's standard text format

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/codec"
	"github.com/jhump/protoreflect/desc"
)

// MarshalText serializes this message to bytes in the standard text format,
// returning an error if the operation fails. The resulting bytes will be a
// valid UTF8 string.
//
// This method uses a compact form: no newlines, and spaces between field
// identifiers and values are elided.
func (m *Message) MarshalText() ([]byte, error) ***REMOVED***
	var b indentBuffer
	b.indentCount = -1 // no indentation
	if err := m.marshalText(&b); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.Bytes(), nil
***REMOVED***

// MarshalTextIndent serializes this message to bytes in the standard text
// format, returning an error if the operation fails. The resulting bytes will
// be a valid UTF8 string.
//
// This method uses a "pretty-printed" form, with each field on its own line and
// spaces between field identifiers and values.
func (m *Message) MarshalTextIndent() ([]byte, error) ***REMOVED***
	var b indentBuffer
	b.indent = "  " // TODO: option for indent?
	if err := m.marshalText(&b); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.Bytes(), nil
***REMOVED***

func (m *Message) marshalText(b *indentBuffer) error ***REMOVED***
	// TODO: option for emitting extended Any format?
	first := true
	// first the known fields
	for _, tag := range m.knownFieldTags() ***REMOVED***
		itag := int32(tag)
		v := m.values[itag]
		fd := m.FindFieldDescriptor(itag)
		if fd.IsMap() ***REMOVED***
			md := fd.GetMessageType()
			kfd := md.FindFieldByNumber(1)
			vfd := md.FindFieldByNumber(2)
			mp := v.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
			keys := make([]interface***REMOVED******REMOVED***, 0, len(mp))
			for k := range mp ***REMOVED***
				keys = append(keys, k)
			***REMOVED***
			sort.Sort(sortable(keys))
			for _, mk := range keys ***REMOVED***
				mv := mp[mk]
				err := b.maybeNext(&first)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = marshalKnownFieldMapEntryText(b, fd, kfd, mk, vfd, mv)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if fd.IsRepeated() ***REMOVED***
			sl := v.([]interface***REMOVED******REMOVED***)
			for _, slv := range sl ***REMOVED***
				err := b.maybeNext(&first)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = marshalKnownFieldText(b, fd, slv)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err := b.maybeNext(&first)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = marshalKnownFieldText(b, fd, v)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// then the unknown fields
	for _, tag := range m.unknownFieldTags() ***REMOVED***
		itag := int32(tag)
		ufs := m.unknownFields[itag]
		for _, uf := range ufs ***REMOVED***
			err := b.maybeNext(&first)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			_, err = fmt.Fprintf(b, "%d", tag)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if uf.Encoding == proto.WireStartGroup ***REMOVED***
				err = b.WriteByte('***REMOVED***')
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = b.start()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				in := codec.NewBuffer(uf.Contents)
				err = marshalUnknownGroupText(b, in, true)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = b.end()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = b.WriteByte('***REMOVED***')
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				err = b.sep()
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if uf.Encoding == proto.WireBytes ***REMOVED***
					err = writeString(b, string(uf.Contents))
					if err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					_, err = b.WriteString(strconv.FormatUint(uf.Value, 10))
					if err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func marshalKnownFieldMapEntryText(b *indentBuffer, fd *desc.FieldDescriptor, kfd *desc.FieldDescriptor, mk interface***REMOVED******REMOVED***, vfd *desc.FieldDescriptor, mv interface***REMOVED******REMOVED***) error ***REMOVED***
	var name string
	if fd.IsExtension() ***REMOVED***
		name = fmt.Sprintf("[%s]", fd.GetFullyQualifiedName())
	***REMOVED*** else ***REMOVED***
		name = fd.GetName()
	***REMOVED***
	_, err := b.WriteString(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.sep()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = b.WriteByte('<')
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.start()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = marshalKnownFieldText(b, kfd, mk)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = b.next()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !isNil(mv) ***REMOVED***
		err = marshalKnownFieldText(b, vfd, mv)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	err = b.end()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return b.WriteByte('>')
***REMOVED***

func marshalKnownFieldText(b *indentBuffer, fd *desc.FieldDescriptor, v interface***REMOVED******REMOVED***) error ***REMOVED***
	group := fd.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP
	if group ***REMOVED***
		var name string
		if fd.IsExtension() ***REMOVED***
			name = fmt.Sprintf("[%s]", fd.GetMessageType().GetFullyQualifiedName())
		***REMOVED*** else ***REMOVED***
			name = fd.GetMessageType().GetName()
		***REMOVED***
		_, err := b.WriteString(name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var name string
		if fd.IsExtension() ***REMOVED***
			name = fmt.Sprintf("[%s]", fd.GetFullyQualifiedName())
		***REMOVED*** else ***REMOVED***
			name = fd.GetName()
		***REMOVED***
		_, err := b.WriteString(name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = b.sep()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	rv := reflect.ValueOf(v)
	switch rv.Kind() ***REMOVED***
	case reflect.Int32, reflect.Int64:
		ed := fd.GetEnumType()
		if ed != nil ***REMOVED***
			n := int32(rv.Int())
			vd := ed.FindValueByNumber(n)
			if vd == nil ***REMOVED***
				_, err := b.WriteString(strconv.FormatInt(rv.Int(), 10))
				return err
			***REMOVED*** else ***REMOVED***
				_, err := b.WriteString(vd.GetName())
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			_, err := b.WriteString(strconv.FormatInt(rv.Int(), 10))
			return err
		***REMOVED***
	case reflect.Uint32, reflect.Uint64:
		_, err := b.WriteString(strconv.FormatUint(rv.Uint(), 10))
		return err
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		var str string
		if math.IsNaN(f) ***REMOVED***
			str = "nan"
		***REMOVED*** else if math.IsInf(f, 1) ***REMOVED***
			str = "inf"
		***REMOVED*** else if math.IsInf(f, -1) ***REMOVED***
			str = "-inf"
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
		return writeString(b, string(rv.Bytes()))
	case reflect.String:
		return writeString(b, rv.String())
	default:
		var err error
		if group ***REMOVED***
			err = b.WriteByte('***REMOVED***')
		***REMOVED*** else ***REMOVED***
			err = b.WriteByte('<')
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = b.start()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// must be a message
		if dm, ok := v.(*Message); ok ***REMOVED***
			err = dm.marshalText(b)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = proto.CompactText(b, v.(proto.Message))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		err = b.end()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if group ***REMOVED***
			return b.WriteByte('***REMOVED***')
		***REMOVED*** else ***REMOVED***
			return b.WriteByte('>')
		***REMOVED***
	***REMOVED***
***REMOVED***

// writeString writes a string in the protocol buffer text format.
// It is similar to strconv.Quote except we don't use Go escape sequences,
// we treat the string as a byte sequence, and we use octal escapes.
// These differences are to maintain interoperability with the other
// languages' implementations of the text format.
func writeString(b *indentBuffer, s string) error ***REMOVED***
	// use WriteByte here to get any needed indent
	if err := b.WriteByte('"'); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Loop over the bytes, not the runes.
	for i := 0; i < len(s); i++ ***REMOVED***
		var err error
		// Divergence from C++: we don't escape apostrophes.
		// There's no need to escape them, and the C++ parser
		// copes with a naked apostrophe.
		switch c := s[i]; c ***REMOVED***
		case '\n':
			_, err = b.WriteString("\\n")
		case '\r':
			_, err = b.WriteString("\\r")
		case '\t':
			_, err = b.WriteString("\\t")
		case '"':
			_, err = b.WriteString("\\\"")
		case '\\':
			_, err = b.WriteString("\\\\")
		default:
			if c >= 0x20 && c < 0x7f ***REMOVED***
				err = b.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				_, err = fmt.Fprintf(b, "\\%03o", c)
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return b.WriteByte('"')
***REMOVED***

func marshalUnknownGroupText(b *indentBuffer, in *codec.Buffer, topLevel bool) error ***REMOVED***
	first := true
	for ***REMOVED***
		if in.EOF() ***REMOVED***
			if topLevel ***REMOVED***
				return nil
			***REMOVED***
			// this is a nested message: we are expecting an end-group tag, not EOF!
			return io.ErrUnexpectedEOF
		***REMOVED***
		tag, wireType, err := in.DecodeTagAndWireType()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if wireType == proto.WireEndGroup ***REMOVED***
			return nil
		***REMOVED***
		err = b.maybeNext(&first)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		_, err = fmt.Fprintf(b, "%d", tag)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if wireType == proto.WireStartGroup ***REMOVED***
			err = b.WriteByte('***REMOVED***')
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = b.start()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = marshalUnknownGroupText(b, in, false)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = b.end()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = b.WriteByte('***REMOVED***')
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED*** else ***REMOVED***
			err = b.sep()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if wireType == proto.WireBytes ***REMOVED***
				contents, err := in.DecodeRawBytes(false)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = writeString(b, string(contents))
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				var v uint64
				switch wireType ***REMOVED***
				case proto.WireVarint:
					v, err = in.DecodeVarint()
				case proto.WireFixed32:
					v, err = in.DecodeFixed32()
				case proto.WireFixed64:
					v, err = in.DecodeFixed64()
				default:
					return proto.ErrInternalBadWireType
				***REMOVED***
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				_, err = b.WriteString(strconv.FormatUint(v, 10))
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// UnmarshalText de-serializes the message that is present, in text format, in
// the given bytes into this message. It first resets the current message. It
// returns an error if the given bytes do not contain a valid encoding of this
// message type in the standard text format
func (m *Message) UnmarshalText(text []byte) error ***REMOVED***
	m.Reset()
	if err := m.UnmarshalMergeText(text); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.Validate()
***REMOVED***

// UnmarshalMergeText de-serializes the message that is present, in text format,
// in the given bytes into this message. Unlike UnmarshalText, it does not first
// reset the message, instead merging the data in the given bytes into the
// existing data in this message.
func (m *Message) UnmarshalMergeText(text []byte) error ***REMOVED***
	return m.unmarshalText(newReader(text), tokenEOF)
***REMOVED***

func (m *Message) unmarshalText(tr *txtReader, end tokenType) error ***REMOVED***
	for ***REMOVED***
		tok := tr.next()
		if tok.tokTyp == end ***REMOVED***
			return nil
		***REMOVED***
		if tok.tokTyp == tokenEOF ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		var fd *desc.FieldDescriptor
		var extendedAnyType *desc.MessageDescriptor
		if tok.tokTyp == tokenInt ***REMOVED***
			// tag number (indicates unknown field)
			tag, err := strconv.ParseInt(tok.val.(string), 10, 32)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			itag := int32(tag)
			fd = m.FindFieldDescriptor(itag)
			if fd == nil ***REMOVED***
				// can't parse the value w/out field descriptor, so skip it
				tok = tr.next()
				if tok.tokTyp == tokenEOF ***REMOVED***
					return io.ErrUnexpectedEOF
				***REMOVED*** else if tok.tokTyp == tokenOpenBrace ***REMOVED***
					if err := skipMessageText(tr, true); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED*** else if tok.tokTyp == tokenColon ***REMOVED***
					if err := skipFieldValueText(tr); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					return textError(tok, "Expecting a colon ':' or brace '***REMOVED***'; instead got %q", tok.txt)
				***REMOVED***
				tok = tr.peek()
				if tok.tokTyp.IsSep() ***REMOVED***
					tr.next() // consume separator
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			fieldName, err := unmarshalFieldNameText(tr, tok)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			fd = m.FindFieldDescriptorByName(fieldName)
			if fd == nil ***REMOVED***
				// See if it's a group name
				for _, field := range m.md.GetFields() ***REMOVED***
					if field.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP && field.GetMessageType().GetName() == fieldName ***REMOVED***
						fd = field
						break
					***REMOVED***
				***REMOVED***
				if fd == nil ***REMOVED***
					// maybe this is an extended Any
					if m.md.GetFullyQualifiedName() == "google.protobuf.Any" && fieldName[0] == '[' && strings.Contains(fieldName, "/") ***REMOVED***
						// strip surrounding "[" and "]" and extract type name from URL
						typeUrl := fieldName[1 : len(fieldName)-1]
						mname := typeUrl
						if slash := strings.LastIndex(mname, "/"); slash >= 0 ***REMOVED***
							mname = mname[slash+1:]
						***REMOVED***
						// TODO: add a way to weave an AnyResolver to this point
						extendedAnyType = findMessageDescriptor(mname, m.md.GetFile())
						if extendedAnyType == nil ***REMOVED***
							return textError(tok, "could not parse Any with unknown type URL %q", fieldName)
						***REMOVED***
						// field 1 is "type_url"
						typeUrlField := m.md.FindFieldByNumber(1)
						if err := m.TrySetField(typeUrlField, typeUrl); err != nil ***REMOVED***
							return err
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						// TODO: add a flag to just ignore unrecognized field names
						return textError(tok, "%q is not a recognized field name of %q", fieldName, m.md.GetFullyQualifiedName())
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		tok = tr.next()
		if tok.tokTyp == tokenEOF ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		if extendedAnyType != nil ***REMOVED***
			// consume optional colon; make sure this is a "start message" token
			if tok.tokTyp == tokenColon ***REMOVED***
				tok = tr.next()
				if tok.tokTyp == tokenEOF ***REMOVED***
					return io.ErrUnexpectedEOF
				***REMOVED***
			***REMOVED***
			if tok.tokTyp.EndToken() == tokenError ***REMOVED***
				return textError(tok, "Expecting a '<' or '***REMOVED***'; instead got %q", tok.txt)
			***REMOVED***

			// TODO: use mf.NewMessage and, if not a dynamic message, use proto.UnmarshalText to unmarshal it
			g := m.mf.NewDynamicMessage(extendedAnyType)
			if err := g.unmarshalText(tr, tok.tokTyp.EndToken()); err != nil ***REMOVED***
				return err
			***REMOVED***
			// now we marshal the message to bytes and store in the Any
			b, err := g.Marshal()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			// field 2 is "value"
			anyValueField := m.md.FindFieldByNumber(2)
			if err := m.TrySetField(anyValueField, b); err != nil ***REMOVED***
				return err
			***REMOVED***

		***REMOVED*** else if (fd.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP ||
			fd.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE) &&
			tok.tokTyp.EndToken() != tokenError ***REMOVED***

			// TODO: use mf.NewMessage and, if not a dynamic message, use proto.UnmarshalText to unmarshal it
			g := m.mf.NewDynamicMessage(fd.GetMessageType())
			if err := g.unmarshalText(tr, tok.tokTyp.EndToken()); err != nil ***REMOVED***
				return err
			***REMOVED***
			if fd.IsRepeated() ***REMOVED***
				if err := m.TryAddRepeatedField(fd, g); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if err := m.TrySetField(fd, g); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if tok.tokTyp != tokenColon ***REMOVED***
				return textError(tok, "Expecting a colon ':'; instead got %q", tok.txt)
			***REMOVED***
			if err := m.unmarshalFieldValueText(fd, tr); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		tok = tr.peek()
		if tok.tokTyp.IsSep() ***REMOVED***
			tr.next() // consume separator
		***REMOVED***
	***REMOVED***
***REMOVED***
func findMessageDescriptor(name string, fd *desc.FileDescriptor) *desc.MessageDescriptor ***REMOVED***
	md := findMessageInTransitiveDeps(name, fd, map[*desc.FileDescriptor]struct***REMOVED******REMOVED******REMOVED******REMOVED***)
	if md == nil ***REMOVED***
		// couldn't find it; see if we have this message linked in
		md, _ = desc.LoadMessageDescriptor(name)
	***REMOVED***
	return md
***REMOVED***

func findMessageInTransitiveDeps(name string, fd *desc.FileDescriptor, seen map[*desc.FileDescriptor]struct***REMOVED******REMOVED***) *desc.MessageDescriptor ***REMOVED***
	if _, ok := seen[fd]; ok ***REMOVED***
		// already checked this file
		return nil
	***REMOVED***
	seen[fd] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	md := fd.FindMessage(name)
	if md != nil ***REMOVED***
		return md
	***REMOVED***
	// not in this file so recursively search its deps
	for _, dep := range fd.GetDependencies() ***REMOVED***
		md = findMessageInTransitiveDeps(name, dep, seen)
		if md != nil ***REMOVED***
			return md
		***REMOVED***
	***REMOVED***
	// couldn't find it
	return nil
***REMOVED***

func textError(tok *token, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	var msg string
	if tok.tokTyp == tokenError ***REMOVED***
		msg = tok.val.(error).Error()
	***REMOVED*** else ***REMOVED***
		msg = fmt.Sprintf(format, args...)
	***REMOVED***
	return fmt.Errorf("line %d, col %d: %s", tok.pos.Line, tok.pos.Column, msg)
***REMOVED***

type setFunction func(*Message, *desc.FieldDescriptor, interface***REMOVED******REMOVED***) error

func (m *Message) unmarshalFieldValueText(fd *desc.FieldDescriptor, tr *txtReader) error ***REMOVED***
	var set setFunction
	if fd.IsRepeated() ***REMOVED***
		set = (*Message).addRepeatedField
	***REMOVED*** else ***REMOVED***
		set = mergeField
	***REMOVED***
	tok := tr.peek()
	if tok.tokTyp == tokenOpenBracket ***REMOVED***
		tr.next() // consume tok
		for ***REMOVED***
			if err := m.unmarshalFieldElementText(fd, tr, set); err != nil ***REMOVED***
				return err
			***REMOVED***
			tok = tr.peek()
			if tok.tokTyp == tokenCloseBracket ***REMOVED***
				tr.next() // consume tok
				return nil
			***REMOVED*** else if tok.tokTyp.IsSep() ***REMOVED***
				tr.next() // consume separator
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return m.unmarshalFieldElementText(fd, tr, set)
***REMOVED***

func (m *Message) unmarshalFieldElementText(fd *desc.FieldDescriptor, tr *txtReader, set setFunction) error ***REMOVED***
	tok := tr.next()
	if tok.tokTyp == tokenEOF ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***

	var expected string
	switch fd.GetType() ***REMOVED***
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if tok.tokTyp == tokenIdent ***REMOVED***
			if tok.val.(string) == "true" ***REMOVED***
				return set(m, fd, true)
			***REMOVED*** else if tok.val.(string) == "false" ***REMOVED***
				return set(m, fd, false)
			***REMOVED***
		***REMOVED***
		expected = "boolean value"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if tok.tokTyp == tokenString ***REMOVED***
			return set(m, fd, []byte(tok.val.(string)))
		***REMOVED***
		expected = "bytes string value"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		if tok.tokTyp == tokenString ***REMOVED***
			return set(m, fd, tok.val)
		***REMOVED***
		expected = "string value"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		switch tok.tokTyp ***REMOVED***
		case tokenFloat:
			return set(m, fd, float32(tok.val.(float64)))
		case tokenInt:
			if f, err := strconv.ParseFloat(tok.val.(string), 32); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, float32(f))
			***REMOVED***
		case tokenIdent:
			ident := strings.ToLower(tok.val.(string))
			if ident == "inf" ***REMOVED***
				return set(m, fd, float32(math.Inf(1)))
			***REMOVED*** else if ident == "nan" ***REMOVED***
				return set(m, fd, float32(math.NaN()))
			***REMOVED***
		case tokenMinus:
			peeked := tr.peek()
			if peeked.tokTyp == tokenIdent ***REMOVED***
				ident := strings.ToLower(peeked.val.(string))
				if ident == "inf" ***REMOVED***
					tr.next() // consume peeked token
					return set(m, fd, float32(math.Inf(-1)))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		expected = "float value"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		switch tok.tokTyp ***REMOVED***
		case tokenFloat:
			return set(m, fd, tok.val)
		case tokenInt:
			if f, err := strconv.ParseFloat(tok.val.(string), 64); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, f)
			***REMOVED***
		case tokenIdent:
			ident := strings.ToLower(tok.val.(string))
			if ident == "inf" ***REMOVED***
				return set(m, fd, math.Inf(1))
			***REMOVED*** else if ident == "nan" ***REMOVED***
				return set(m, fd, math.NaN())
			***REMOVED***
		case tokenMinus:
			peeked := tr.peek()
			if peeked.tokTyp == tokenIdent ***REMOVED***
				ident := strings.ToLower(peeked.val.(string))
				if ident == "inf" ***REMOVED***
					tr.next() // consume peeked token
					return set(m, fd, math.Inf(-1))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		expected = "float value"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		if tok.tokTyp == tokenInt ***REMOVED***
			if i, err := strconv.ParseInt(tok.val.(string), 10, 32); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, int32(i))
			***REMOVED***
		***REMOVED***
		expected = "int value"
	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		if tok.tokTyp == tokenInt ***REMOVED***
			if i, err := strconv.ParseInt(tok.val.(string), 10, 64); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, i)
			***REMOVED***
		***REMOVED***
		expected = "int value"
	case descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32:
		if tok.tokTyp == tokenInt ***REMOVED***
			if i, err := strconv.ParseUint(tok.val.(string), 10, 32); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, uint32(i))
			***REMOVED***
		***REMOVED***
		expected = "unsigned int value"
	case descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64:
		if tok.tokTyp == tokenInt ***REMOVED***
			if i, err := strconv.ParseUint(tok.val.(string), 10, 64); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, i)
			***REMOVED***
		***REMOVED***
		expected = "unsigned int value"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		if tok.tokTyp == tokenIdent ***REMOVED***
			// TODO: add a flag to just ignore unrecognized enum value names?
			vd := fd.GetEnumType().FindValueByName(tok.val.(string))
			if vd != nil ***REMOVED***
				return set(m, fd, vd.GetNumber())
			***REMOVED***
		***REMOVED*** else if tok.tokTyp == tokenInt ***REMOVED***
			if i, err := strconv.ParseInt(tok.val.(string), 10, 32); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				return set(m, fd, int32(i))
			***REMOVED***
		***REMOVED***
		expected = fmt.Sprintf("enum %s value", fd.GetEnumType().GetFullyQualifiedName())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_GROUP:

		endTok := tok.tokTyp.EndToken()
		if endTok != tokenError ***REMOVED***
			dm := m.mf.NewDynamicMessage(fd.GetMessageType())
			if err := dm.unmarshalText(tr, endTok); err != nil ***REMOVED***
				return err
			***REMOVED***
			// TODO: ideally we would use mf.NewMessage and, if not a dynamic message, use
			// proto package to unmarshal it. But the text parser isn't particularly amenable
			// to that, so we instead convert a dynamic message to a generated one if the
			// known-type registry knows about the generated type...
			var ktr *KnownTypeRegistry
			if m.mf != nil ***REMOVED***
				ktr = m.mf.ktr
			***REMOVED***
			pm := ktr.CreateIfKnown(fd.GetMessageType().GetFullyQualifiedName())
			if pm != nil ***REMOVED***
				if err := dm.ConvertTo(pm); err != nil ***REMOVED***
					return set(m, fd, pm)
				***REMOVED***
			***REMOVED***
			return set(m, fd, dm)
		***REMOVED***
		expected = fmt.Sprintf("message %s value", fd.GetMessageType().GetFullyQualifiedName())
	default:
		return fmt.Errorf("field %q of message %q has unrecognized type: %v", fd.GetFullyQualifiedName(), m.md.GetFullyQualifiedName(), fd.GetType())
	***REMOVED***

	// if we get here, token was wrong type; create error message
	var article string
	if strings.Contains("aieou", expected[0:1]) ***REMOVED***
		article = "an"
	***REMOVED*** else ***REMOVED***
		article = "a"
	***REMOVED***
	return textError(tok, "Expecting %s %s; got %q", article, expected, tok.txt)
***REMOVED***

func unmarshalFieldNameText(tr *txtReader, tok *token) (string, error) ***REMOVED***
	if tok.tokTyp == tokenOpenBracket || tok.tokTyp == tokenOpenParen ***REMOVED***
		// extension name
		var closeType tokenType
		var closeChar string
		if tok.tokTyp == tokenOpenBracket ***REMOVED***
			closeType = tokenCloseBracket
			closeChar = "close bracket ']'"
		***REMOVED*** else ***REMOVED***
			closeType = tokenCloseParen
			closeChar = "close paren ')'"
		***REMOVED***
		// must be followed by an identifier
		idents := make([]string, 0, 1)
		for ***REMOVED***
			tok = tr.next()
			if tok.tokTyp == tokenEOF ***REMOVED***
				return "", io.ErrUnexpectedEOF
			***REMOVED*** else if tok.tokTyp != tokenIdent ***REMOVED***
				return "", textError(tok, "Expecting an identifier; instead got %q", tok.txt)
			***REMOVED***
			idents = append(idents, tok.val.(string))
			// and then close bracket/paren, or "/" to keep adding URL elements to name
			tok = tr.next()
			if tok.tokTyp == tokenEOF ***REMOVED***
				return "", io.ErrUnexpectedEOF
			***REMOVED*** else if tok.tokTyp == closeType ***REMOVED***
				break
			***REMOVED*** else if tok.tokTyp != tokenSlash ***REMOVED***
				return "", textError(tok, "Expecting a %s; instead got %q", closeChar, tok.txt)
			***REMOVED***
		***REMOVED***
		return "[" + strings.Join(idents, "/") + "]", nil
	***REMOVED*** else if tok.tokTyp == tokenIdent ***REMOVED***
		// normal field name
		return tok.val.(string), nil
	***REMOVED*** else ***REMOVED***
		return "", textError(tok, "Expecting an identifier or tag number; instead got %q", tok.txt)
	***REMOVED***
***REMOVED***

func skipFieldNameText(tr *txtReader) error ***REMOVED***
	tok := tr.next()
	if tok.tokTyp == tokenEOF ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED*** else if tok.tokTyp == tokenInt || tok.tokTyp == tokenIdent ***REMOVED***
		return nil
	***REMOVED*** else ***REMOVED***
		_, err := unmarshalFieldNameText(tr, tok)
		return err
	***REMOVED***
***REMOVED***

func skipFieldValueText(tr *txtReader) error ***REMOVED***
	tok := tr.peek()
	if tok.tokTyp == tokenOpenBracket ***REMOVED***
		tr.next() // consume tok
		for ***REMOVED***
			if err := skipFieldElementText(tr); err != nil ***REMOVED***
				return err
			***REMOVED***
			tok = tr.peek()
			if tok.tokTyp == tokenCloseBracket ***REMOVED***
				tr.next() // consume tok
				return nil
			***REMOVED*** else if tok.tokTyp.IsSep() ***REMOVED***
				tr.next() // consume separator
			***REMOVED***

		***REMOVED***
	***REMOVED***
	return skipFieldElementText(tr)
***REMOVED***

func skipFieldElementText(tr *txtReader) error ***REMOVED***
	tok := tr.next()
	switch tok.tokTyp ***REMOVED***
	case tokenEOF:
		return io.ErrUnexpectedEOF
	case tokenInt, tokenFloat, tokenString, tokenIdent:
		return nil
	case tokenOpenAngle:
		return skipMessageText(tr, false)
	default:
		return textError(tok, "Expecting an angle bracket '<' or a value; instead got %q", tok.txt)
	***REMOVED***
***REMOVED***

func skipMessageText(tr *txtReader, isGroup bool) error ***REMOVED***
	for ***REMOVED***
		tok := tr.peek()
		if tok.tokTyp == tokenEOF ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED*** else if isGroup && tok.tokTyp == tokenCloseBrace ***REMOVED***
			return nil
		***REMOVED*** else if !isGroup && tok.tokTyp == tokenCloseAngle ***REMOVED***
			return nil
		***REMOVED***

		// field name or tag
		if err := skipFieldNameText(tr); err != nil ***REMOVED***
			return err
		***REMOVED***

		// field value
		tok = tr.next()
		if tok.tokTyp == tokenEOF ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED*** else if tok.tokTyp == tokenOpenBrace ***REMOVED***
			if err := skipMessageText(tr, true); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if tok.tokTyp == tokenColon ***REMOVED***
			if err := skipFieldValueText(tr); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return textError(tok, "Expecting a colon ':' or brace '***REMOVED***'; instead got %q", tok.txt)
		***REMOVED***

		tok = tr.peek()
		if tok.tokTyp.IsSep() ***REMOVED***
			tr.next() // consume separator
		***REMOVED***
	***REMOVED***
***REMOVED***

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenIdent
	tokenString
	tokenInt
	tokenFloat
	tokenColon
	tokenComma
	tokenSemiColon
	tokenOpenBrace
	tokenCloseBrace
	tokenOpenBracket
	tokenCloseBracket
	tokenOpenAngle
	tokenCloseAngle
	tokenOpenParen
	tokenCloseParen
	tokenSlash
	tokenMinus
)

func (t tokenType) IsSep() bool ***REMOVED***
	return t == tokenComma || t == tokenSemiColon
***REMOVED***

func (t tokenType) EndToken() tokenType ***REMOVED***
	switch t ***REMOVED***
	case tokenOpenAngle:
		return tokenCloseAngle
	case tokenOpenBrace:
		return tokenCloseBrace
	default:
		return tokenError
	***REMOVED***
***REMOVED***

type token struct ***REMOVED***
	tokTyp tokenType
	val    interface***REMOVED******REMOVED***
	txt    string
	pos    scanner.Position
***REMOVED***

type txtReader struct ***REMOVED***
	scanner    scanner.Scanner
	peeked     token
	havePeeked bool
***REMOVED***

func newReader(text []byte) *txtReader ***REMOVED***
	sc := scanner.Scanner***REMOVED******REMOVED***
	sc.Init(bytes.NewReader(text))
	sc.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanChars |
		scanner.ScanStrings | scanner.ScanComments | scanner.SkipComments
	// identifiers are same restrictions as Go identifiers, except we also allow dots since
	// we accept fully-qualified names
	sc.IsIdentRune = func(ch rune, i int) bool ***REMOVED***
		return ch == '_' || unicode.IsLetter(ch) ||
			(i > 0 && unicode.IsDigit(ch)) ||
			(i > 0 && ch == '.')
	***REMOVED***
	// ignore errors; we handle them if/when we see malformed tokens
	sc.Error = func(s *scanner.Scanner, msg string) ***REMOVED******REMOVED***
	return &txtReader***REMOVED***scanner: sc***REMOVED***
***REMOVED***

func (p *txtReader) peek() *token ***REMOVED***
	if p.havePeeked ***REMOVED***
		return &p.peeked
	***REMOVED***
	t := p.scanner.Scan()
	if t == scanner.EOF ***REMOVED***
		p.peeked.tokTyp = tokenEOF
		p.peeked.val = nil
		p.peeked.txt = ""
		p.peeked.pos = p.scanner.Position
	***REMOVED*** else if err := p.processToken(t, p.scanner.TokenText(), p.scanner.Position); err != nil ***REMOVED***
		p.peeked.tokTyp = tokenError
		p.peeked.val = err
	***REMOVED***
	p.havePeeked = true
	return &p.peeked
***REMOVED***

func (p *txtReader) processToken(t rune, text string, pos scanner.Position) error ***REMOVED***
	p.peeked.pos = pos
	p.peeked.txt = text
	switch t ***REMOVED***
	case scanner.Ident:
		p.peeked.tokTyp = tokenIdent
		p.peeked.val = text
	case scanner.Int:
		p.peeked.tokTyp = tokenInt
		p.peeked.val = text // can't parse the number because we don't know if it's signed or unsigned
	case scanner.Float:
		p.peeked.tokTyp = tokenFloat
		var err error
		if p.peeked.val, err = strconv.ParseFloat(text, 64); err != nil ***REMOVED***
			return err
		***REMOVED***
	case scanner.Char, scanner.String:
		p.peeked.tokTyp = tokenString
		var err error
		if p.peeked.val, err = strconv.Unquote(text); err != nil ***REMOVED***
			return err
		***REMOVED***
	case '-': // unary minus, for negative ints and floats
		ch := p.scanner.Peek()
		if ch < '0' || ch > '9' ***REMOVED***
			p.peeked.tokTyp = tokenMinus
			p.peeked.val = '-'
		***REMOVED*** else ***REMOVED***
			t := p.scanner.Scan()
			if t == scanner.EOF ***REMOVED***
				return io.ErrUnexpectedEOF
			***REMOVED*** else if t == scanner.Float ***REMOVED***
				p.peeked.tokTyp = tokenFloat
				text += p.scanner.TokenText()
				p.peeked.txt = text
				var err error
				if p.peeked.val, err = strconv.ParseFloat(text, 64); err != nil ***REMOVED***
					p.peeked.pos = p.scanner.Position
					return err
				***REMOVED***
			***REMOVED*** else if t == scanner.Int ***REMOVED***
				p.peeked.tokTyp = tokenInt
				text += p.scanner.TokenText()
				p.peeked.txt = text
				p.peeked.val = text // can't parse the number because we don't know if it's signed or unsigned
			***REMOVED*** else ***REMOVED***
				p.peeked.pos = p.scanner.Position
				return fmt.Errorf("expecting an int or float but got %q", p.scanner.TokenText())
			***REMOVED***
		***REMOVED***
	case ':':
		p.peeked.tokTyp = tokenColon
		p.peeked.val = ':'
	case ',':
		p.peeked.tokTyp = tokenComma
		p.peeked.val = ','
	case ';':
		p.peeked.tokTyp = tokenSemiColon
		p.peeked.val = ';'
	case '***REMOVED***':
		p.peeked.tokTyp = tokenOpenBrace
		p.peeked.val = '***REMOVED***'
	case '***REMOVED***':
		p.peeked.tokTyp = tokenCloseBrace
		p.peeked.val = '***REMOVED***'
	case '<':
		p.peeked.tokTyp = tokenOpenAngle
		p.peeked.val = '<'
	case '>':
		p.peeked.tokTyp = tokenCloseAngle
		p.peeked.val = '>'
	case '[':
		p.peeked.tokTyp = tokenOpenBracket
		p.peeked.val = '['
	case ']':
		p.peeked.tokTyp = tokenCloseBracket
		p.peeked.val = ']'
	case '(':
		p.peeked.tokTyp = tokenOpenParen
		p.peeked.val = '('
	case ')':
		p.peeked.tokTyp = tokenCloseParen
		p.peeked.val = ')'
	case '/':
		// only allowed to separate URL components in expanded Any format
		p.peeked.tokTyp = tokenSlash
		p.peeked.val = '/'
	default:
		return fmt.Errorf("invalid character: %c", t)
	***REMOVED***
	return nil
***REMOVED***

func (p *txtReader) next() *token ***REMOVED***
	t := p.peek()
	if t.tokTyp != tokenEOF && t.tokTyp != tokenError ***REMOVED***
		p.havePeeked = false
	***REMOVED***
	return t
***REMOVED***
