package dynamic

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/codec"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/internal"
)

// ErrUnknownTagNumber is an error that is returned when an operation refers
// to an unknown tag number.
var ErrUnknownTagNumber = errors.New("unknown tag number")

// UnknownTagNumberError is the same as ErrUnknownTagNumber.
// Deprecated: use ErrUnknownTagNumber
var UnknownTagNumberError = ErrUnknownTagNumber

// ErrUnknownFieldName is an error that is returned when an operation refers
// to an unknown field name.
var ErrUnknownFieldName = errors.New("unknown field name")

// UnknownFieldNameError is the same as ErrUnknownFieldName.
// Deprecated: use ErrUnknownFieldName
var UnknownFieldNameError = ErrUnknownFieldName

// ErrFieldIsNotMap is an error that is returned when map-related operations
// are attempted with fields that are not maps.
var ErrFieldIsNotMap = errors.New("field is not a map type")

// FieldIsNotMapError is the same as ErrFieldIsNotMap.
// Deprecated: use ErrFieldIsNotMap
var FieldIsNotMapError = ErrFieldIsNotMap

// ErrFieldIsNotRepeated is an error that is returned when repeated field
// operations are attempted with fields that are not repeated.
var ErrFieldIsNotRepeated = errors.New("field is not repeated")

// FieldIsNotRepeatedError is the same as ErrFieldIsNotRepeated.
// Deprecated: use ErrFieldIsNotRepeated
var FieldIsNotRepeatedError = ErrFieldIsNotRepeated

// ErrIndexOutOfRange is an error that is returned when an invalid index is
// provided when access a single element of a repeated field.
var ErrIndexOutOfRange = errors.New("index is out of range")

// IndexOutOfRangeError is the same as ErrIndexOutOfRange.
// Deprecated: use ErrIndexOutOfRange
var IndexOutOfRangeError = ErrIndexOutOfRange

// ErrNumericOverflow is an error returned by operations that encounter a
// numeric value that is too large, for example de-serializing a value into an
// int32 field when the value is larger that can fit into a 32-bit value.
var ErrNumericOverflow = errors.New("numeric value is out of range")

// NumericOverflowError is the same as ErrNumericOverflow.
// Deprecated: use ErrNumericOverflow
var NumericOverflowError = ErrNumericOverflow

var typeOfProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()
var typeOfDynamicMessage = reflect.TypeOf((*Message)(nil))
var typeOfBytes = reflect.TypeOf(([]byte)(nil))

// Message is a dynamic protobuf message. Instead of a generated struct,
// like most protobuf messages, this is a map of field number to values and
// a message descriptor, which is used to validate the field values and
// also to de-serialize messages (from the standard binary format, as well
// as from the text format and from JSON).
type Message struct ***REMOVED***
	md            *desc.MessageDescriptor
	er            *ExtensionRegistry
	mf            *MessageFactory
	extraFields   map[int32]*desc.FieldDescriptor
	values        map[int32]interface***REMOVED******REMOVED***
	unknownFields map[int32][]UnknownField
***REMOVED***

// UnknownField represents a field that was parsed from the binary wire
// format for a message, but was not a recognized field number. Enough
// information is preserved so that re-serializing the message won't lose
// any of the unrecognized data.
type UnknownField struct ***REMOVED***
	// Encoding indicates how the unknown field was encoded on the wire. If it
	// is proto.WireBytes or proto.WireGroupStart then Contents will be set to
	// the raw bytes. If it is proto.WireTypeFixed32 then the data is in the least
	// significant 32 bits of Value. Otherwise, the data is in all 64 bits of
	// Value.
	Encoding int8
	Contents []byte
	Value    uint64
***REMOVED***

// NewMessage creates a new dynamic message for the type represented by the given
// message descriptor. During de-serialization, a default MessageFactory is used to
// instantiate any nested message fields and no extension fields will be parsed. To
// use a custom MessageFactory or ExtensionRegistry, use MessageFactory.NewMessage.
func NewMessage(md *desc.MessageDescriptor) *Message ***REMOVED***
	return NewMessageWithMessageFactory(md, nil)
***REMOVED***

// NewMessageWithExtensionRegistry creates a new dynamic message for the type
// represented by the given message descriptor. During de-serialization, the given
// ExtensionRegistry is used to parse extension fields and nested messages will be
// instantiated using dynamic.NewMessageFactoryWithExtensionRegistry(er).
func NewMessageWithExtensionRegistry(md *desc.MessageDescriptor, er *ExtensionRegistry) *Message ***REMOVED***
	mf := NewMessageFactoryWithExtensionRegistry(er)
	return NewMessageWithMessageFactory(md, mf)
***REMOVED***

// NewMessageWithMessageFactory creates a new dynamic message for the type
// represented by the given message descriptor. During de-serialization, the given
// MessageFactory is used to instantiate nested messages.
func NewMessageWithMessageFactory(md *desc.MessageDescriptor, mf *MessageFactory) *Message ***REMOVED***
	var er *ExtensionRegistry
	if mf != nil ***REMOVED***
		er = mf.er
	***REMOVED***
	return &Message***REMOVED***
		md: md,
		mf: mf,
		er: er,
	***REMOVED***
***REMOVED***

// AsDynamicMessage converts the given message to a dynamic message. If the
// given message is dynamic, it is returned. Otherwise, a dynamic message is
// created using NewMessage.
func AsDynamicMessage(msg proto.Message) (*Message, error) ***REMOVED***
	return AsDynamicMessageWithMessageFactory(msg, nil)
***REMOVED***

// AsDynamicMessageWithExtensionRegistry converts the given message to a dynamic
// message. If the given message is dynamic, it is returned. Otherwise, a
// dynamic message is created using NewMessageWithExtensionRegistry.
func AsDynamicMessageWithExtensionRegistry(msg proto.Message, er *ExtensionRegistry) (*Message, error) ***REMOVED***
	mf := NewMessageFactoryWithExtensionRegistry(er)
	return AsDynamicMessageWithMessageFactory(msg, mf)
***REMOVED***

// AsDynamicMessageWithMessageFactory converts the given message to a dynamic
// message. If the given message is dynamic, it is returned. Otherwise, a
// dynamic message is created using NewMessageWithMessageFactory.
func AsDynamicMessageWithMessageFactory(msg proto.Message, mf *MessageFactory) (*Message, error) ***REMOVED***
	if dm, ok := msg.(*Message); ok ***REMOVED***
		return dm, nil
	***REMOVED***
	md, err := desc.LoadMessageDescriptorForMessage(msg)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	dm := NewMessageWithMessageFactory(md, mf)
	err = dm.mergeFrom(msg)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return dm, nil
***REMOVED***

// GetMessageDescriptor returns a descriptor for this message's type.
func (m *Message) GetMessageDescriptor() *desc.MessageDescriptor ***REMOVED***
	return m.md
***REMOVED***

// GetKnownFields returns a slice of descriptors for all known fields. The
// fields will not be in any defined order.
func (m *Message) GetKnownFields() []*desc.FieldDescriptor ***REMOVED***
	if len(m.extraFields) == 0 ***REMOVED***
		return m.md.GetFields()
	***REMOVED***
	flds := make([]*desc.FieldDescriptor, len(m.md.GetFields()), len(m.md.GetFields())+len(m.extraFields))
	copy(flds, m.md.GetFields())
	for _, fld := range m.extraFields ***REMOVED***
		if !fld.IsExtension() ***REMOVED***
			flds = append(flds, fld)
		***REMOVED***
	***REMOVED***
	return flds
***REMOVED***

// GetKnownExtensions returns a slice of descriptors for all extensions known by
// the message's extension registry. The fields will not be in any defined order.
func (m *Message) GetKnownExtensions() []*desc.FieldDescriptor ***REMOVED***
	if !m.md.IsExtendable() ***REMOVED***
		return nil
	***REMOVED***
	exts := m.er.AllExtensionsForType(m.md.GetFullyQualifiedName())
	for _, fld := range m.extraFields ***REMOVED***
		if fld.IsExtension() ***REMOVED***
			exts = append(exts, fld)
		***REMOVED***
	***REMOVED***
	return exts
***REMOVED***

// GetUnknownFields returns a slice of tag numbers for all unknown fields that
// this message contains. The tags will not be in any defined order.
func (m *Message) GetUnknownFields() []int32 ***REMOVED***
	flds := make([]int32, 0, len(m.unknownFields))
	for tag := range m.unknownFields ***REMOVED***
		flds = append(flds, tag)
	***REMOVED***
	return flds
***REMOVED***

// Descriptor returns the serialized form of the file descriptor in which the
// message was defined and a path to the message type therein. This mimics the
// method of the same name on message types generated by protoc.
func (m *Message) Descriptor() ([]byte, []int) ***REMOVED***
	// get encoded file descriptor
	b, err := proto.Marshal(m.md.GetFile().AsProto())
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("failed to get encoded descriptor for %s: %v", m.md.GetFile().GetName(), err))
	***REMOVED***
	var zippedBytes bytes.Buffer
	w := gzip.NewWriter(&zippedBytes)
	if _, err := w.Write(b); err != nil ***REMOVED***
		panic(fmt.Sprintf("failed to get encoded descriptor for %s: %v", m.md.GetFile().GetName(), err))
	***REMOVED***
	if err := w.Close(); err != nil ***REMOVED***
		panic(fmt.Sprintf("failed to get an encoded descriptor for %s: %v", m.md.GetFile().GetName(), err))
	***REMOVED***

	// and path to message
	path := []int***REMOVED******REMOVED***
	var d desc.Descriptor
	name := m.md.GetFullyQualifiedName()
	for d = m.md.GetParent(); d != nil; name, d = d.GetFullyQualifiedName(), d.GetParent() ***REMOVED***
		found := false
		switch d := d.(type) ***REMOVED***
		case (*desc.FileDescriptor):
			for i, md := range d.GetMessageTypes() ***REMOVED***
				if md.GetFullyQualifiedName() == name ***REMOVED***
					found = true
					path = append(path, i)
				***REMOVED***
			***REMOVED***
		case (*desc.MessageDescriptor):
			for i, md := range d.GetNestedMessageTypes() ***REMOVED***
				if md.GetFullyQualifiedName() == name ***REMOVED***
					found = true
					path = append(path, i)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			panic(fmt.Sprintf("failed to compute descriptor path for %s", m.md.GetFullyQualifiedName()))
		***REMOVED***
	***REMOVED***
	// reverse the path
	i := 0
	j := len(path) - 1
	for i < j ***REMOVED***
		path[i], path[j] = path[j], path[i]
		i++
		j--
	***REMOVED***

	return zippedBytes.Bytes(), path
***REMOVED***

// XXX_MessageName returns the fully qualified name of this message's type. This
// allows dynamic messages to be used with proto.MessageName.
func (m *Message) XXX_MessageName() string ***REMOVED***
	return m.md.GetFullyQualifiedName()
***REMOVED***

// FindFieldDescriptor returns a field descriptor for the given tag number. This
// searches known fields in the descriptor, known fields discovered during calls
// to GetField or SetField, and extension fields known by the message's extension
// registry. It returns nil if the tag is unknown.
func (m *Message) FindFieldDescriptor(tagNumber int32) *desc.FieldDescriptor ***REMOVED***
	fd := m.md.FindFieldByNumber(tagNumber)
	if fd != nil ***REMOVED***
		return fd
	***REMOVED***
	fd = m.er.FindExtension(m.md.GetFullyQualifiedName(), tagNumber)
	if fd != nil ***REMOVED***
		return fd
	***REMOVED***
	return m.extraFields[tagNumber]
***REMOVED***

// FindFieldDescriptorByName returns a field descriptor for the given field
// name. This searches known fields in the descriptor, known fields discovered
// during calls to GetField or SetField, and extension fields known by the
// message's extension registry. It returns nil if the name is unknown. If the
// given name refers to an extension, it should be fully qualified and may be
// optionally enclosed in parentheses or brackets.
func (m *Message) FindFieldDescriptorByName(name string) *desc.FieldDescriptor ***REMOVED***
	if name == "" ***REMOVED***
		return nil
	***REMOVED***
	fd := m.md.FindFieldByName(name)
	if fd != nil ***REMOVED***
		return fd
	***REMOVED***
	mustBeExt := false
	if name[0] == '(' ***REMOVED***
		if name[len(name)-1] != ')' ***REMOVED***
			// malformed name
			return nil
		***REMOVED***
		mustBeExt = true
		name = name[1 : len(name)-1]
	***REMOVED*** else if name[0] == '[' ***REMOVED***
		if name[len(name)-1] != ']' ***REMOVED***
			// malformed name
			return nil
		***REMOVED***
		mustBeExt = true
		name = name[1 : len(name)-1]
	***REMOVED***
	fd = m.er.FindExtensionByName(m.md.GetFullyQualifiedName(), name)
	if fd != nil ***REMOVED***
		return fd
	***REMOVED***
	for _, fd := range m.extraFields ***REMOVED***
		if fd.IsExtension() && name == fd.GetFullyQualifiedName() ***REMOVED***
			return fd
		***REMOVED*** else if !mustBeExt && !fd.IsExtension() && name == fd.GetName() ***REMOVED***
			return fd
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// FindFieldDescriptorByJSONName returns a field descriptor for the given JSON
// name. This searches known fields in the descriptor, known fields discovered
// during calls to GetField or SetField, and extension fields known by the
// message's extension registry. If no field matches the given JSON name, it
// will fall back to searching field names (e.g. FindFieldDescriptorByName). If
// this also yields no match, nil is returned.
func (m *Message) FindFieldDescriptorByJSONName(name string) *desc.FieldDescriptor ***REMOVED***
	if name == "" ***REMOVED***
		return nil
	***REMOVED***
	fd := m.md.FindFieldByJSONName(name)
	if fd != nil ***REMOVED***
		return fd
	***REMOVED***
	mustBeExt := false
	if name[0] == '(' ***REMOVED***
		if name[len(name)-1] != ')' ***REMOVED***
			// malformed name
			return nil
		***REMOVED***
		mustBeExt = true
		name = name[1 : len(name)-1]
	***REMOVED*** else if name[0] == '[' ***REMOVED***
		if name[len(name)-1] != ']' ***REMOVED***
			// malformed name
			return nil
		***REMOVED***
		mustBeExt = true
		name = name[1 : len(name)-1]
	***REMOVED***
	fd = m.er.FindExtensionByJSONName(m.md.GetFullyQualifiedName(), name)
	if fd != nil ***REMOVED***
		return fd
	***REMOVED***
	for _, fd := range m.extraFields ***REMOVED***
		if fd.IsExtension() && name == fd.GetFullyQualifiedJSONName() ***REMOVED***
			return fd
		***REMOVED*** else if !mustBeExt && !fd.IsExtension() && name == fd.GetJSONName() ***REMOVED***
			return fd
		***REMOVED***
	***REMOVED***

	// try non-JSON names
	return m.FindFieldDescriptorByName(name)
***REMOVED***

func (m *Message) checkField(fd *desc.FieldDescriptor) error ***REMOVED***
	return checkField(fd, m.md)
***REMOVED***

func checkField(fd *desc.FieldDescriptor, md *desc.MessageDescriptor) error ***REMOVED***
	if fd.GetOwner().GetFullyQualifiedName() != md.GetFullyQualifiedName() ***REMOVED***
		return fmt.Errorf("given field, %s, is for wrong message type: %s; expecting %s", fd.GetName(), fd.GetOwner().GetFullyQualifiedName(), md.GetFullyQualifiedName())
	***REMOVED***
	if fd.IsExtension() && !md.IsExtension(fd.GetNumber()) ***REMOVED***
		return fmt.Errorf("given field, %s, is an extension but is not in message extension range: %v", fd.GetFullyQualifiedName(), md.GetExtensionRanges())
	***REMOVED***
	return nil
***REMOVED***

// GetField returns the value for the given field descriptor. It panics if an
// error is encountered. See TryGetField.
func (m *Message) GetField(fd *desc.FieldDescriptor) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetField(fd); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetField returns the value for the given field descriptor. An error is
// returned if the given field descriptor does not belong to the right message
// type.
//
// The Go type of the returned value, for scalar fields, is the same as protoc
// would generate for the field (in a non-dynamic message). The table below
// lists the scalar types and the corresponding Go types.
//  +-------------------------+-----------+
//  |       Declared Type     |  Go Type  |
//  +-------------------------+-----------+
//  | int32, sint32, sfixed32 | int32     |
//  | int64, sint64, sfixed64 | int64     |
//  | uint32, fixed32         | uint32    |
//  | uint64, fixed64         | uint64    |
//  | float                   | float32   |
//  | double                  | double32  |
//  | bool                    | bool      |
//  | string                  | string    |
//  | bytes                   | []byte    |
//  +-------------------------+-----------+
//
// Values for enum fields will always be int32 values. You can use the enum
// descriptor associated with the field to lookup value names with those values.
// Values for message type fields may be an instance of the generated type *or*
// may be another *dynamic.Message that represents the type.
//
// If the given field is a map field, the returned type will be
// map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***. The actual concrete types of keys and values is
// as described above. If the given field is a (non-map) repeated field, the
// returned type is always []interface***REMOVED******REMOVED***; the type of the actual elements is as
// described above.
//
// If this message has no value for the given field, its default value is
// returned. If the message is defined in a file with "proto3" syntax, the
// default is always the zero value for the field. The default value for map and
// repeated fields is a nil map or slice (respectively). For field's whose types
// is a message, the default value is an empty message for "proto2" syntax or a
// nil message for "proto3" syntax. Note that the in the latter case, a non-nil
// interface with a nil pointer is returned, not a nil interface. Also note that
// whether the returned value is an empty message or nil depends on if *this*
// message was defined as "proto3" syntax, not the message type referred to by
// the field's type.
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) but corresponds to an unknown field, the unknown value will be
// parsed and become known. The parsed value will be returned, or an error will
// be returned if the unknown value cannot be parsed according to the field
// descriptor's type information.
func (m *Message) TryGetField(fd *desc.FieldDescriptor) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m.getField(fd)
***REMOVED***

// GetFieldByName returns the value for the field with the given name. It panics
// if an error is encountered. See TryGetFieldByName.
func (m *Message) GetFieldByName(name string) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetFieldByName(name); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetFieldByName returns the value for the field with the given name. An
// error is returned if the given name is unknown. If the given name refers to
// an extension field, it should be fully qualified and optionally enclosed in
// parenthesis or brackets.
//
// If this message has no value for the given field, its default value is
// returned. (See TryGetField for more info on types and default field values.)
func (m *Message) TryGetFieldByName(name string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return nil, UnknownFieldNameError
	***REMOVED***
	return m.getField(fd)
***REMOVED***

// GetFieldByNumber returns the value for the field with the given tag number.
// It panics if an error is encountered. See TryGetFieldByNumber.
func (m *Message) GetFieldByNumber(tagNumber int) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetFieldByNumber(tagNumber); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetFieldByNumber returns the value for the field with the given tag
// number. An error is returned if the given tag is unknown.
//
// If this message has no value for the given field, its default value is
// returned. (See TryGetField for more info on types and default field values.)
func (m *Message) TryGetFieldByNumber(tagNumber int) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return nil, UnknownTagNumberError
	***REMOVED***
	return m.getField(fd)
***REMOVED***

func (m *Message) getField(fd *desc.FieldDescriptor) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return m.doGetField(fd, false)
***REMOVED***

func (m *Message) doGetField(fd *desc.FieldDescriptor, nilIfAbsent bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	res := m.values[fd.GetNumber()]
	if res == nil ***REMOVED***
		var err error
		if res, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if res == nil ***REMOVED***
			if nilIfAbsent ***REMOVED***
				return nil, nil
			***REMOVED*** else ***REMOVED***
				def := fd.GetDefaultValue()
				if def != nil ***REMOVED***
					return def, nil
				***REMOVED***
				// GetDefaultValue only returns nil for message types
				md := fd.GetMessageType()
				if m.md.IsProto3() ***REMOVED***
					return nilMessage(md), nil
				***REMOVED*** else ***REMOVED***
					// for proto2, return default instance of message
					return m.mf.NewMessage(md), nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	rt := reflect.TypeOf(res)
	if rt.Kind() == reflect.Map ***REMOVED***
		// make defensive copies to prevent caller from storing illegal keys and values
		m := res.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
		res := map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for k, v := range m ***REMOVED***
			res[k] = v
		***REMOVED***
		return res, nil
	***REMOVED*** else if rt.Kind() == reflect.Slice && rt != typeOfBytes ***REMOVED***
		// make defensive copies to prevent caller from storing illegal elements
		sl := res.([]interface***REMOVED******REMOVED***)
		res := make([]interface***REMOVED******REMOVED***, len(sl))
		copy(res, sl)
		return res, nil
	***REMOVED***
	return res, nil
***REMOVED***

func nilMessage(md *desc.MessageDescriptor) interface***REMOVED******REMOVED*** ***REMOVED***
	// try to return a proper nil pointer
	msgType := proto.MessageType(md.GetFullyQualifiedName())
	if msgType != nil && msgType.Implements(typeOfProtoMessage) ***REMOVED***
		return reflect.Zero(msgType).Interface().(proto.Message)
	***REMOVED***
	// fallback to nil dynamic message pointer
	return (*Message)(nil)
***REMOVED***

// HasField returns true if this message has a value for the given field. If the
// given field is not valid (e.g. belongs to a different message type), false is
// returned. If this message is defined in a file with "proto3" syntax, this
// will return false even if a field was explicitly assigned its zero value (the
// zero values for a field are intentionally indistinguishable from absent).
func (m *Message) HasField(fd *desc.FieldDescriptor) bool ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return false
	***REMOVED***
	return m.HasFieldNumber(int(fd.GetNumber()))
***REMOVED***

// HasFieldName returns true if this message has a value for a field with the
// given name. If the given name is unknown, this returns false.
func (m *Message) HasFieldName(name string) bool ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return false
	***REMOVED***
	return m.HasFieldNumber(int(fd.GetNumber()))
***REMOVED***

// HasFieldNumber returns true if this message has a value for a field with the
// given tag number. If the given tag is unknown, this returns false.
func (m *Message) HasFieldNumber(tagNumber int) bool ***REMOVED***
	if _, ok := m.values[int32(tagNumber)]; ok ***REMOVED***
		return true
	***REMOVED***
	_, ok := m.unknownFields[int32(tagNumber)]
	return ok
***REMOVED***

// SetField sets the value for the given field descriptor to the given value. It
// panics if an error is encountered. See TrySetField.
func (m *Message) SetField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TrySetField(fd, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TrySetField sets the value for the given field descriptor to the given value.
// An error is returned if the given field descriptor does not belong to the
// right message type or if the given value is not a correct/compatible type for
// the given field.
//
// The Go type expected for a field  is the same as TryGetField would return for
// the field. So message values can be supplied as either the correct generated
// message type or as a *dynamic.Message.
//
// Since it is cumbersome to work with dynamic messages, some concessions are
// made to simplify usage regarding types:
//
//  1. If a numeric type is provided that can be converted *without loss or
//     overflow*, it is accepted. This allows for setting int64 fields using int
//     or int32 values. Similarly for uint64 with uint and uint32 values and for
//     float64 fields with float32 values.
//  2. The value can be a named type, as long as its underlying type is correct.
//  3. Map and repeated fields can be set using any kind of concrete map or
//     slice type, as long as the values within are all of the correct type. So
//     a field defined as a 'map<string, int32>` can be set using a
//     map[string]int32, a map[string]interface***REMOVED******REMOVED***, or even a
//     map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***.
//  4. Finally, dynamic code that chooses to not treat maps as a special-case
//     find that they can set map fields using a slice where each element is a
//     message that matches the implicit map-entry field message type.
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) it will become known. Subsequent operations using tag numbers or
// names will be able to resolve the newly-known type. If the message has a
// value for the unknown value, it is cleared, replaced by the given known
// value.
func (m *Message) TrySetField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.setField(fd, val)
***REMOVED***

// SetFieldByName sets the value for the field with the given name to the given
// value. It panics if an error is encountered. See TrySetFieldByName.
func (m *Message) SetFieldByName(name string, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TrySetFieldByName(name, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TrySetFieldByName sets the value for the field with the given name to the
// given value. An error is returned if the given name is unknown or if the
// given value has an incorrect type. If the given name refers to an extension
// field, it should be fully qualified and optionally enclosed in parenthesis or
// brackets.
//
// (See TrySetField for more info on types.)
func (m *Message) TrySetFieldByName(name string, val interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	return m.setField(fd, val)
***REMOVED***

// SetFieldByNumber sets the value for the field with the given tag number to
// the given value. It panics if an error is encountered. See
// TrySetFieldByNumber.
func (m *Message) SetFieldByNumber(tagNumber int, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TrySetFieldByNumber(tagNumber, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TrySetFieldByNumber sets the value for the field with the given tag number to
// the given value. An error is returned if the given tag is unknown or if the
// given value has an incorrect type.
//
// (See TrySetField for more info on types.)
func (m *Message) TrySetFieldByNumber(tagNumber int, val interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	return m.setField(fd, val)
***REMOVED***

func (m *Message) setField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error
	if val, err = validFieldValue(fd, val); err != nil ***REMOVED***
		return err
	***REMOVED***
	m.internalSetField(fd, val)
	return nil
***REMOVED***

func (m *Message) internalSetField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) ***REMOVED***
	if fd.IsRepeated() ***REMOVED***
		// Unset fields and zero-length fields are indistinguishable, in both
		// proto2 and proto3 syntax
		if reflect.ValueOf(val).Len() == 0 ***REMOVED***
			if m.values != nil ***REMOVED***
				delete(m.values, fd.GetNumber())
			***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else if m.md.IsProto3() && fd.GetOneOf() == nil ***REMOVED***
		// proto3 considers fields that are set to their zero value as unset
		// (we already handled repeated fields above)
		var equal bool
		if b, ok := val.([]byte); ok ***REMOVED***
			// can't compare slices, so we have to special-case []byte values
			equal = ok && bytes.Equal(b, fd.GetDefaultValue().([]byte))
		***REMOVED*** else ***REMOVED***
			defVal := fd.GetDefaultValue()
			equal = defVal == val
			if !equal && defVal == nil ***REMOVED***
				// above just checks if value is the nil interface,
				// but we should also test if the given value is a
				// nil pointer
				rv := reflect.ValueOf(val)
				if rv.Kind() == reflect.Ptr && rv.IsNil() ***REMOVED***
					equal = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if equal ***REMOVED***
			if m.values != nil ***REMOVED***
				delete(m.values, fd.GetNumber())
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if m.values == nil ***REMOVED***
		m.values = map[int32]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	m.values[fd.GetNumber()] = val
	// if this field is part of a one-of, make sure all other one-of choices are cleared
	od := fd.GetOneOf()
	if od != nil ***REMOVED***
		for _, other := range od.GetChoices() ***REMOVED***
			if other.GetNumber() != fd.GetNumber() ***REMOVED***
				delete(m.values, other.GetNumber())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// also clear any unknown fields
	if m.unknownFields != nil ***REMOVED***
		delete(m.unknownFields, fd.GetNumber())
	***REMOVED***
	// and add this field if it was previously unknown
	if existing := m.FindFieldDescriptor(fd.GetNumber()); existing == nil ***REMOVED***
		m.addField(fd)
	***REMOVED***
***REMOVED***

func (m *Message) addField(fd *desc.FieldDescriptor) ***REMOVED***
	if m.extraFields == nil ***REMOVED***
		m.extraFields = map[int32]*desc.FieldDescriptor***REMOVED******REMOVED***
	***REMOVED***
	m.extraFields[fd.GetNumber()] = fd
***REMOVED***

// ClearField removes any value for the given field. It panics if an error is
// encountered. See TryClearField.
func (m *Message) ClearField(fd *desc.FieldDescriptor) ***REMOVED***
	if err := m.TryClearField(fd); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryClearField removes any value for the given field. An error is returned if
// the given field descriptor does not belong to the right message type.
func (m *Message) TryClearField(fd *desc.FieldDescriptor) error ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	m.clearField(fd)
	return nil
***REMOVED***

// ClearFieldByName removes any value for the field with the given name. It
// panics if an error is encountered. See TryClearFieldByName.
func (m *Message) ClearFieldByName(name string) ***REMOVED***
	if err := m.TryClearFieldByName(name); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryClearFieldByName removes any value for the field with the given name. An
// error is returned if the given name is unknown. If the given name refers to
// an extension field, it should be fully qualified and optionally enclosed in
// parenthesis or brackets.
func (m *Message) TryClearFieldByName(name string) error ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	m.clearField(fd)
	return nil
***REMOVED***

// ClearFieldByNumber removes any value for the field with the given tag number.
// It panics if an error is encountered. See TryClearFieldByNumber.
func (m *Message) ClearFieldByNumber(tagNumber int) ***REMOVED***
	if err := m.TryClearFieldByNumber(tagNumber); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryClearFieldByNumber removes any value for the field with the given tag
// number. An error is returned if the given tag is unknown.
func (m *Message) TryClearFieldByNumber(tagNumber int) error ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	m.clearField(fd)
	return nil
***REMOVED***

func (m *Message) clearField(fd *desc.FieldDescriptor) ***REMOVED***
	// clear value
	if m.values != nil ***REMOVED***
		delete(m.values, fd.GetNumber())
	***REMOVED***
	// also clear any unknown fields
	if m.unknownFields != nil ***REMOVED***
		delete(m.unknownFields, fd.GetNumber())
	***REMOVED***
	// and add this field if it was previously unknown
	if existing := m.FindFieldDescriptor(fd.GetNumber()); existing == nil ***REMOVED***
		m.addField(fd)
	***REMOVED***
***REMOVED***

// GetOneOfField returns which of the given one-of's fields is set and the
// corresponding value. It panics if an error is encountered. See
// TryGetOneOfField.
func (m *Message) GetOneOfField(od *desc.OneOfDescriptor) (*desc.FieldDescriptor, interface***REMOVED******REMOVED***) ***REMOVED***
	if fd, val, err := m.TryGetOneOfField(od); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return fd, val
	***REMOVED***
***REMOVED***

// TryGetOneOfField returns which of the given one-of's fields is set and the
// corresponding value. An error is returned if the given one-of belongs to the
// wrong message type. If the given one-of has no field set, this method will
// return nil, nil.
//
// The type of the value, if one is set, is the same as would be returned by
// TryGetField using the returned field descriptor.
//
// Like with TryGetField, if the given one-of contains any fields that are not
// known (e.g. not present in this message's descriptor), they will become known
// and any unknown value will be parsed (and become a known value on success).
func (m *Message) TryGetOneOfField(od *desc.OneOfDescriptor) (*desc.FieldDescriptor, interface***REMOVED******REMOVED***, error) ***REMOVED***
	if od.GetOwner().GetFullyQualifiedName() != m.md.GetFullyQualifiedName() ***REMOVED***
		return nil, nil, fmt.Errorf("given one-of, %s, is for wrong message type: %s; expecting %s", od.GetName(), od.GetOwner().GetFullyQualifiedName(), m.md.GetFullyQualifiedName())
	***REMOVED***
	for _, fd := range od.GetChoices() ***REMOVED***
		val, err := m.doGetField(fd, true)
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
		if val != nil ***REMOVED***
			return fd, val, nil
		***REMOVED***
	***REMOVED***
	return nil, nil, nil
***REMOVED***

// ClearOneOfField removes any value for any of the given one-of's fields. It
// panics if an error is encountered. See TryClearOneOfField.
func (m *Message) ClearOneOfField(od *desc.OneOfDescriptor) ***REMOVED***
	if err := m.TryClearOneOfField(od); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryClearOneOfField removes any value for any of the given one-of's fields. An
// error is returned if the given one-of descriptor does not belong to the right
// message type.
func (m *Message) TryClearOneOfField(od *desc.OneOfDescriptor) error ***REMOVED***
	if od.GetOwner().GetFullyQualifiedName() != m.md.GetFullyQualifiedName() ***REMOVED***
		return fmt.Errorf("given one-of, %s, is for wrong message type: %s; expecting %s", od.GetName(), od.GetOwner().GetFullyQualifiedName(), m.md.GetFullyQualifiedName())
	***REMOVED***
	for _, fd := range od.GetChoices() ***REMOVED***
		m.clearField(fd)
	***REMOVED***
	return nil
***REMOVED***

// GetMapField returns the value for the given map field descriptor and given
// key. It panics if an error is encountered. See TryGetMapField.
func (m *Message) GetMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetMapField(fd, key); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetMapField returns the value for the given map field descriptor and given
// key. An error is returned if the given field descriptor does not belong to
// the right message type or if it is not a map field.
//
// If the map field does not contain the requested key, this method returns
// nil, nil. The Go type of the value returned mirrors the type that protoc
// would generate for the field. (See TryGetField for more details on types).
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) but corresponds to an unknown field, the unknown value will be
// parsed and become known. The parsed value will be searched for the requested
// key and any value returned. An error will be returned if the unknown value
// cannot be parsed according to the field descriptor's type information.
func (m *Message) TryGetMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m.getMapField(fd, key)
***REMOVED***

// GetMapFieldByName returns the value for the map field with the given name and
// given key. It panics if an error is encountered. See TryGetMapFieldByName.
func (m *Message) GetMapFieldByName(name string, key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetMapFieldByName(name, key); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetMapFieldByName returns the value for the map field with the given name
// and given key. An error is returned if the given name is unknown or if it
// names a field that is not a map field.
//
// If this message has no value for the given field or the value has no value
// for the requested key, then this method returns nil, nil.
//
// (See TryGetField for more info on types.)
func (m *Message) TryGetMapFieldByName(name string, key interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return nil, UnknownFieldNameError
	***REMOVED***
	return m.getMapField(fd, key)
***REMOVED***

// GetMapFieldByNumber returns the value for the map field with the given tag
// number and given key. It panics if an error is encountered. See
// TryGetMapFieldByNumber.
func (m *Message) GetMapFieldByNumber(tagNumber int, key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetMapFieldByNumber(tagNumber, key); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetMapFieldByNumber returns the value for the map field with the given tag
// number and given key. An error is returned if the given tag is unknown or if
// it indicates a field that is not a map field.
//
// If this message has no value for the given field or the value has no value
// for the requested key, then this method returns nil, nil.
//
// (See TryGetField for more info on types.)
func (m *Message) TryGetMapFieldByNumber(tagNumber int, key interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return nil, UnknownTagNumberError
	***REMOVED***
	return m.getMapField(fd, key)
***REMOVED***

func (m *Message) getMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if !fd.IsMap() ***REMOVED***
		return nil, FieldIsNotMapError
	***REMOVED***
	kfd := fd.GetMessageType().GetFields()[0]
	ki, err := validElementFieldValue(kfd, key, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mp := m.values[fd.GetNumber()]
	if mp == nil ***REMOVED***
		if mp, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if mp == nil ***REMOVED***
			return nil, nil
		***REMOVED***
	***REMOVED***
	return mp.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)[ki], nil
***REMOVED***

// ForEachMapFieldEntry executes the given function for each entry in the map
// value for the given field descriptor. It stops iteration if the function
// returns false. It panics if an error is encountered. See
// TryForEachMapFieldEntry.
func (m *Message) ForEachMapFieldEntry(fd *desc.FieldDescriptor, fn func(key, val interface***REMOVED******REMOVED***) bool) ***REMOVED***
	if err := m.TryForEachMapFieldEntry(fd, fn); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryForEachMapFieldEntry executes the given function for each entry in the map
// value for the given field descriptor. An error is returned if the given field
// descriptor does not belong to the right message type or if it is not a  map
// field.
//
// Iteration ends either when all entries have been examined or when the given
// function returns false. So the function is expected to return true for normal
// iteration and false to break out. If this message has no value for the given
// field, it returns without invoking the given function.
//
// The Go type of the key and value supplied to the function mirrors the type
// that protoc would generate for the field. (See TryGetField for more details
// on types).
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) but corresponds to an unknown field, the unknown value will be
// parsed and become known. The parsed value will be searched for the requested
// key and any value returned. An error will be returned if the unknown value
// cannot be parsed according to the field descriptor's type information.
func (m *Message) TryForEachMapFieldEntry(fd *desc.FieldDescriptor, fn func(key, val interface***REMOVED******REMOVED***) bool) error ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.forEachMapFieldEntry(fd, fn)
***REMOVED***

// ForEachMapFieldEntryByName executes the given function for each entry in the
// map value for the field with the given name. It stops iteration if the
// function returns false. It panics if an error is encountered. See
// TryForEachMapFieldEntryByName.
func (m *Message) ForEachMapFieldEntryByName(name string, fn func(key, val interface***REMOVED******REMOVED***) bool) ***REMOVED***
	if err := m.TryForEachMapFieldEntryByName(name, fn); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryForEachMapFieldEntryByName executes the given function for each entry in
// the map value for the field with the given name. It stops iteration if the
// function returns false. An error is returned if the given name is unknown or
// if it names a field that is not a map field.
//
// If this message has no value for the given field, it returns without ever
// invoking the given function.
//
// (See TryGetField for more info on types supplied to the function.)
func (m *Message) TryForEachMapFieldEntryByName(name string, fn func(key, val interface***REMOVED******REMOVED***) bool) error ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	return m.forEachMapFieldEntry(fd, fn)
***REMOVED***

// ForEachMapFieldEntryByNumber executes the given function for each entry in
// the map value for the field with the given tag number. It stops iteration if
// the function returns false. It panics if an error is encountered. See
// TryForEachMapFieldEntryByNumber.
func (m *Message) ForEachMapFieldEntryByNumber(tagNumber int, fn func(key, val interface***REMOVED******REMOVED***) bool) ***REMOVED***
	if err := m.TryForEachMapFieldEntryByNumber(tagNumber, fn); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryForEachMapFieldEntryByNumber executes the given function for each entry in
// the map value for the field with the given tag number. It stops iteration if
// the function returns false. An error is returned if the given tag is unknown
// or if it indicates a field that is not a map field.
//
// If this message has no value for the given field, it returns without ever
// invoking the given function.
//
// (See TryGetField for more info on types supplied to the function.)
func (m *Message) TryForEachMapFieldEntryByNumber(tagNumber int, fn func(key, val interface***REMOVED******REMOVED***) bool) error ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	return m.forEachMapFieldEntry(fd, fn)
***REMOVED***

func (m *Message) forEachMapFieldEntry(fd *desc.FieldDescriptor, fn func(key, val interface***REMOVED******REMOVED***) bool) error ***REMOVED***
	if !fd.IsMap() ***REMOVED***
		return FieldIsNotMapError
	***REMOVED***
	mp := m.values[fd.GetNumber()]
	if mp == nil ***REMOVED***
		if mp, err := m.parseUnknownField(fd); err != nil ***REMOVED***
			return err
		***REMOVED*** else if mp == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	for k, v := range mp.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***) ***REMOVED***
		if !fn(k, v) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// PutMapField sets the value for the given map field descriptor and given key
// to the given value. It panics if an error is encountered. See TryPutMapField.
func (m *Message) PutMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryPutMapField(fd, key, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryPutMapField sets the value for the given map field descriptor and given
// key to the given value. An error is returned if the given field descriptor
// does not belong to the right message type, if the given field is not a map
// field, or if the given value is not a correct/compatible type for the given
// field.
//
// The Go type expected for a field  is the same as required by TrySetField for
// a field with the same type as the map's value type.
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) it will become known. Subsequent operations using tag numbers or
// names will be able to resolve the newly-known type. If the message has a
// value for the unknown value, it is cleared, replaced by the given known
// value.
func (m *Message) TryPutMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.putMapField(fd, key, val)
***REMOVED***

// PutMapFieldByName sets the value for the map field with the given name and
// given key to the given value. It panics if an error is encountered. See
// TryPutMapFieldByName.
func (m *Message) PutMapFieldByName(name string, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryPutMapFieldByName(name, key, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryPutMapFieldByName sets the value for the map field with the given name and
// the given key to the given value. An error is returned if the given name is
// unknown, if it names a field that is not a map, or if the given value has an
// incorrect type.
//
// (See TrySetField for more info on types.)
func (m *Message) TryPutMapFieldByName(name string, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	return m.putMapField(fd, key, val)
***REMOVED***

// PutMapFieldByNumber sets the value for the map field with the given tag
// number and given key to the given value. It panics if an error is
// encountered. See TryPutMapFieldByNumber.
func (m *Message) PutMapFieldByNumber(tagNumber int, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryPutMapFieldByNumber(tagNumber, key, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryPutMapFieldByNumber sets the value for the map field with the given tag
// number and the given key to the given value. An error is returned if the
// given tag is unknown, if it indicates a field that is not a map, or if the
// given value has an incorrect type.
//
// (See TrySetField for more info on types.)
func (m *Message) TryPutMapFieldByNumber(tagNumber int, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	return m.putMapField(fd, key, val)
***REMOVED***

func (m *Message) putMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if !fd.IsMap() ***REMOVED***
		return FieldIsNotMapError
	***REMOVED***
	kfd := fd.GetMessageType().GetFields()[0]
	ki, err := validElementFieldValue(kfd, key, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	vfd := fd.GetMessageType().GetFields()[1]
	vi, err := validElementFieldValue(vfd, val, true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mp := m.values[fd.GetNumber()]
	if mp == nil ***REMOVED***
		if mp, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return err
		***REMOVED*** else if mp == nil ***REMOVED***
			m.internalSetField(fd, map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***ki: vi***REMOVED***)
			return nil
		***REMOVED***
	***REMOVED***
	mp.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)[ki] = vi
	return nil
***REMOVED***

// RemoveMapField changes the value for the given field descriptor by removing
// any value associated with the given key. It panics if an error is
// encountered. See TryRemoveMapField.
func (m *Message) RemoveMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryRemoveMapField(fd, key); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryRemoveMapField changes the value for the given field descriptor by
// removing any value associated with the given key. An error is returned if the
// given field descriptor does not belong to the right message type or if the
// given field is not a map field.
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) it will become known. Subsequent operations using tag numbers or
// names will be able to resolve the newly-known type. If the message has a
// value for the unknown value, it is parsed and any value for the given key
// removed.
func (m *Message) TryRemoveMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.removeMapField(fd, key)
***REMOVED***

// RemoveMapFieldByName changes the value for the field with the given name by
// removing any value associated with the given key. It panics if an error is
// encountered. See TryRemoveMapFieldByName.
func (m *Message) RemoveMapFieldByName(name string, key interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryRemoveMapFieldByName(name, key); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryRemoveMapFieldByName changes the value for the field with the given name
// by removing any value associated with the given key. An error is returned if
// the given name is unknown or if it names a field that is not a map.
func (m *Message) TryRemoveMapFieldByName(name string, key interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	return m.removeMapField(fd, key)
***REMOVED***

// RemoveMapFieldByNumber changes the value for the field with the given tag
// number by removing any value associated with the given key. It panics if an
// error is encountered. See TryRemoveMapFieldByNumber.
func (m *Message) RemoveMapFieldByNumber(tagNumber int, key interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryRemoveMapFieldByNumber(tagNumber, key); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryRemoveMapFieldByNumber changes the value for the field with the given tag
// number by removing any value associated with the given key. An error is
// returned if the given tag is unknown or if it indicates a field that is not
// a map.
func (m *Message) TryRemoveMapFieldByNumber(tagNumber int, key interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	return m.removeMapField(fd, key)
***REMOVED***

func (m *Message) removeMapField(fd *desc.FieldDescriptor, key interface***REMOVED******REMOVED***) error ***REMOVED***
	if !fd.IsMap() ***REMOVED***
		return FieldIsNotMapError
	***REMOVED***
	kfd := fd.GetMessageType().GetFields()[0]
	ki, err := validElementFieldValue(kfd, key, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mp := m.values[fd.GetNumber()]
	if mp == nil ***REMOVED***
		if mp, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return err
		***REMOVED*** else if mp == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	res := mp.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
	delete(res, ki)
	if len(res) == 0 ***REMOVED***
		delete(m.values, fd.GetNumber())
	***REMOVED***
	return nil
***REMOVED***

// FieldLength returns the number of elements in this message for the given
// field descriptor. It panics if an error is encountered. See TryFieldLength.
func (m *Message) FieldLength(fd *desc.FieldDescriptor) int ***REMOVED***
	l, err := m.TryFieldLength(fd)
	if err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
	return l
***REMOVED***

// TryFieldLength returns the number of elements in this message for the given
// field descriptor. An error is returned if the given field descriptor does not
// belong to the right message type or if it is neither a map field nor a
// repeated field.
func (m *Message) TryFieldLength(fd *desc.FieldDescriptor) (int, error) ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return m.fieldLength(fd)
***REMOVED***

// FieldLengthByName returns the number of elements in this message for the
// field with the given name. It panics if an error is encountered. See
// TryFieldLengthByName.
func (m *Message) FieldLengthByName(name string) int ***REMOVED***
	l, err := m.TryFieldLengthByName(name)
	if err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
	return l
***REMOVED***

// TryFieldLengthByName returns the number of elements in this message for the
// field with the given name. An error is returned if the given name is unknown
// or if the named field is neither a map field nor a repeated field.
func (m *Message) TryFieldLengthByName(name string) (int, error) ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return 0, UnknownFieldNameError
	***REMOVED***
	return m.fieldLength(fd)
***REMOVED***

// FieldLengthByNumber returns the number of elements in this message for the
// field with the given tag number. It panics if an error is encountered. See
// TryFieldLengthByNumber.
func (m *Message) FieldLengthByNumber(tagNumber int32) int ***REMOVED***
	l, err := m.TryFieldLengthByNumber(tagNumber)
	if err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
	return l
***REMOVED***

// TryFieldLengthByNumber returns the number of elements in this message for the
// field with the given tag number. An error is returned if the given tag is
// unknown or if the named field is neither a map field nor a repeated field.
func (m *Message) TryFieldLengthByNumber(tagNumber int32) (int, error) ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return 0, UnknownTagNumberError
	***REMOVED***
	return m.fieldLength(fd)
***REMOVED***

func (m *Message) fieldLength(fd *desc.FieldDescriptor) (int, error) ***REMOVED***
	if !fd.IsRepeated() ***REMOVED***
		return 0, FieldIsNotRepeatedError
	***REMOVED***
	val := m.values[fd.GetNumber()]
	if val == nil ***REMOVED***
		var err error
		if val, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return 0, err
		***REMOVED*** else if val == nil ***REMOVED***
			return 0, nil
		***REMOVED***
	***REMOVED***
	if sl, ok := val.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return len(sl), nil
	***REMOVED*** else if mp, ok := val.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***); ok ***REMOVED***
		return len(mp), nil
	***REMOVED***
	return 0, nil
***REMOVED***

// GetRepeatedField returns the value for the given repeated field descriptor at
// the given index. It panics if an error is encountered. See
// TryGetRepeatedField.
func (m *Message) GetRepeatedField(fd *desc.FieldDescriptor, index int) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetRepeatedField(fd, index); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetRepeatedField returns the value for the given repeated field descriptor
// at the given index. An error is returned if the given field descriptor does
// not belong to the right message type, if it is not a repeated field, or if
// the given index is out of range (less than zero or greater than or equal to
// the length of the repeated field). Also, even though map fields technically
// are repeated fields, if the given field is a map field an error will result:
// map representation does not lend itself to random access by index.
//
// The Go type of the value returned mirrors the type that protoc would generate
// for the field's element type. (See TryGetField for more details on types).
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) but corresponds to an unknown field, the unknown value will be
// parsed and become known. The value at the given index in the parsed value
// will be returned. An error will be returned if the unknown value cannot be
// parsed according to the field descriptor's type information.
func (m *Message) TryGetRepeatedField(fd *desc.FieldDescriptor, index int) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if index < 0 ***REMOVED***
		return nil, IndexOutOfRangeError
	***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return m.getRepeatedField(fd, index)
***REMOVED***

// GetRepeatedFieldByName returns the value for the repeated field with the
// given name at the given index. It panics if an error is encountered. See
// TryGetRepeatedFieldByName.
func (m *Message) GetRepeatedFieldByName(name string, index int) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetRepeatedFieldByName(name, index); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetRepeatedFieldByName returns the value for the repeated field with the
// given name at the given index. An error is returned if the given name is
// unknown, if it names a field that is not a repeated field (or is a map
// field), or if the given index is out of range (less than zero or greater
// than or equal to the length of the repeated field).
//
// (See TryGetField for more info on types.)
func (m *Message) TryGetRepeatedFieldByName(name string, index int) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if index < 0 ***REMOVED***
		return nil, IndexOutOfRangeError
	***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return nil, UnknownFieldNameError
	***REMOVED***
	return m.getRepeatedField(fd, index)
***REMOVED***

// GetRepeatedFieldByNumber returns the value for the repeated field with the
// given tag number at the given index. It panics if an error is encountered.
// See TryGetRepeatedFieldByNumber.
func (m *Message) GetRepeatedFieldByNumber(tagNumber int, index int) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, err := m.TryGetRepeatedFieldByNumber(tagNumber, index); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED*** else ***REMOVED***
		return v
	***REMOVED***
***REMOVED***

// TryGetRepeatedFieldByNumber returns the value for the repeated field with the
// given tag number at the given index. An error is returned if the given tag is
// unknown, if it indicates a field that is not a repeated field (or is a map
// field), or if the given index is out of range (less than zero or greater than
// or equal to the length of the repeated field).
//
// (See TryGetField for more info on types.)
func (m *Message) TryGetRepeatedFieldByNumber(tagNumber int, index int) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if index < 0 ***REMOVED***
		return nil, IndexOutOfRangeError
	***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return nil, UnknownTagNumberError
	***REMOVED***
	return m.getRepeatedField(fd, index)
***REMOVED***

func (m *Message) getRepeatedField(fd *desc.FieldDescriptor, index int) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if fd.IsMap() || !fd.IsRepeated() ***REMOVED***
		return nil, FieldIsNotRepeatedError
	***REMOVED***
	sl := m.values[fd.GetNumber()]
	if sl == nil ***REMOVED***
		var err error
		if sl, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if sl == nil ***REMOVED***
			return nil, IndexOutOfRangeError
		***REMOVED***
	***REMOVED***
	res := sl.([]interface***REMOVED******REMOVED***)
	if index >= len(res) ***REMOVED***
		return nil, IndexOutOfRangeError
	***REMOVED***
	return res[index], nil
***REMOVED***

// AddRepeatedField appends the given value to the given repeated field. It
// panics if an error is encountered. See TryAddRepeatedField.
func (m *Message) AddRepeatedField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryAddRepeatedField(fd, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryAddRepeatedField appends the given value to the given repeated field. An
// error is returned if the given field descriptor does not belong to the right
// message type, if the given field is not repeated, or if the given value is
// not a correct/compatible type for the given field. If the given field is a
// map field, the call will succeed if the given value is an instance of the
// map's entry message type.
//
// The Go type expected for a field  is the same as required by TrySetField for
// a non-repeated field of the same type.
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) it will become known. Subsequent operations using tag numbers or
// names will be able to resolve the newly-known type. If the message has a
// value for the unknown value, it is parsed and the given value is appended to
// it.
func (m *Message) TryAddRepeatedField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.addRepeatedField(fd, val)
***REMOVED***

// AddRepeatedFieldByName appends the given value to the repeated field with the
// given name. It panics if an error is encountered. See
// TryAddRepeatedFieldByName.
func (m *Message) AddRepeatedFieldByName(name string, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryAddRepeatedFieldByName(name, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryAddRepeatedFieldByName appends the given value to the repeated field with
// the given name. An error is returned if the given name is unknown, if it
// names a field that is not repeated, or if the given value has an incorrect
// type.
//
// (See TrySetField for more info on types.)
func (m *Message) TryAddRepeatedFieldByName(name string, val interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	return m.addRepeatedField(fd, val)
***REMOVED***

// AddRepeatedFieldByNumber appends the given value to the repeated field with
// the given tag number. It panics if an error is encountered. See
// TryAddRepeatedFieldByNumber.
func (m *Message) AddRepeatedFieldByNumber(tagNumber int, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TryAddRepeatedFieldByNumber(tagNumber, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TryAddRepeatedFieldByNumber appends the given value to the repeated field
// with the given tag number. An error is returned if the given tag is unknown,
// if it indicates a field that is not repeated, or if the given value has an
// incorrect type.
//
// (See TrySetField for more info on types.)
func (m *Message) TryAddRepeatedFieldByNumber(tagNumber int, val interface***REMOVED******REMOVED***) error ***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	return m.addRepeatedField(fd, val)
***REMOVED***

func (m *Message) addRepeatedField(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if !fd.IsRepeated() ***REMOVED***
		return FieldIsNotRepeatedError
	***REMOVED***
	val, err := validElementFieldValue(fd, val, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if fd.IsMap() ***REMOVED***
		// We're lenient. Just as we allow setting a map field to a slice of entry messages, we also allow
		// adding entries one at a time (as if the field were a normal repeated field).
		msg := val.(proto.Message)
		dm, err := asDynamicMessage(msg, fd.GetMessageType(), m.mf)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		k, err := dm.TryGetFieldByNumber(1)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v, err := dm.TryGetFieldByNumber(2)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return m.putMapField(fd, k, v)
	***REMOVED***

	sl := m.values[fd.GetNumber()]
	if sl == nil ***REMOVED***
		if sl, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return err
		***REMOVED*** else if sl == nil ***REMOVED***
			sl = []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	res := sl.([]interface***REMOVED******REMOVED***)
	res = append(res, val)
	m.internalSetField(fd, res)
	return nil
***REMOVED***

// SetRepeatedField sets the value for the given repeated field descriptor and
// given index to the given value. It panics if an error is encountered. See
// SetRepeatedField.
func (m *Message) SetRepeatedField(fd *desc.FieldDescriptor, index int, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TrySetRepeatedField(fd, index, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TrySetRepeatedField sets the value for the given repeated field descriptor
// and given index to the given value. An error is returned if the given field
// descriptor does not belong to the right message type, if the given field is
// not repeated, or if the given value is not a correct/compatible type for the
// given field. Also, even though map fields technically are repeated fields, if
// the given field is a map field an error will result: map representation does
// not lend itself to random access by index.
//
// The Go type expected for a field  is the same as required by TrySetField for
// a non-repeated field of the same type.
//
// If the given field descriptor is not known (e.g. not present in the message
// descriptor) it will become known. Subsequent operations using tag numbers or
// names will be able to resolve the newly-known type. If the message has a
// value for the unknown value, it is parsed and the element at the given index
// is replaced with the given value.
func (m *Message) TrySetRepeatedField(fd *desc.FieldDescriptor, index int, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if index < 0 ***REMOVED***
		return IndexOutOfRangeError
	***REMOVED***
	if err := m.checkField(fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.setRepeatedField(fd, index, val)
***REMOVED***

// SetRepeatedFieldByName sets the value for the repeated field with the given
// name and given index to the given value. It panics if an error is
// encountered. See TrySetRepeatedFieldByName.
func (m *Message) SetRepeatedFieldByName(name string, index int, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TrySetRepeatedFieldByName(name, index, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TrySetRepeatedFieldByName sets the value for the repeated field with the
// given name and the given index to the given value. An error is returned if
// the given name is unknown, if it names a field that is not repeated (or is a
// map field), or if the given value has an incorrect type.
//
// (See TrySetField for more info on types.)
func (m *Message) TrySetRepeatedFieldByName(name string, index int, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if index < 0 ***REMOVED***
		return IndexOutOfRangeError
	***REMOVED***
	fd := m.FindFieldDescriptorByName(name)
	if fd == nil ***REMOVED***
		return UnknownFieldNameError
	***REMOVED***
	return m.setRepeatedField(fd, index, val)
***REMOVED***

// SetRepeatedFieldByNumber sets the value for the repeated field with the given
// tag number and given index to the given value. It panics if an error is
// encountered. See TrySetRepeatedFieldByNumber.
func (m *Message) SetRepeatedFieldByNumber(tagNumber int, index int, val interface***REMOVED******REMOVED***) ***REMOVED***
	if err := m.TrySetRepeatedFieldByNumber(tagNumber, index, val); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

// TrySetRepeatedFieldByNumber sets the value for the repeated field with the
// given tag number and the given index to the given value. An error is returned
// if the given tag is unknown, if it indicates a field that is not repeated (or
// is a map field), or if the given value has an incorrect type.
//
// (See TrySetField for more info on types.)
func (m *Message) TrySetRepeatedFieldByNumber(tagNumber int, index int, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if index < 0 ***REMOVED***
		return IndexOutOfRangeError
	***REMOVED***
	fd := m.FindFieldDescriptor(int32(tagNumber))
	if fd == nil ***REMOVED***
		return UnknownTagNumberError
	***REMOVED***
	return m.setRepeatedField(fd, index, val)
***REMOVED***

func (m *Message) setRepeatedField(fd *desc.FieldDescriptor, index int, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if fd.IsMap() || !fd.IsRepeated() ***REMOVED***
		return FieldIsNotRepeatedError
	***REMOVED***
	val, err := validElementFieldValue(fd, val, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sl := m.values[fd.GetNumber()]
	if sl == nil ***REMOVED***
		if sl, err = m.parseUnknownField(fd); err != nil ***REMOVED***
			return err
		***REMOVED*** else if sl == nil ***REMOVED***
			return IndexOutOfRangeError
		***REMOVED***
	***REMOVED***
	res := sl.([]interface***REMOVED******REMOVED***)
	if index >= len(res) ***REMOVED***
		return IndexOutOfRangeError
	***REMOVED***
	res[index] = val
	return nil
***REMOVED***

// GetUnknownField gets the value(s) for the given unknown tag number. If this
// message has no unknown fields with the given tag, nil is returned.
func (m *Message) GetUnknownField(tagNumber int32) []UnknownField ***REMOVED***
	if u, ok := m.unknownFields[tagNumber]; ok ***REMOVED***
		return u
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

func (m *Message) parseUnknownField(fd *desc.FieldDescriptor) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	unks, ok := m.unknownFields[fd.GetNumber()]
	if !ok ***REMOVED***
		return nil, nil
	***REMOVED***
	var v interface***REMOVED******REMOVED***
	var sl []interface***REMOVED******REMOVED***
	var mp map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	if fd.IsMap() ***REMOVED***
		mp = map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	var err error
	for _, unk := range unks ***REMOVED***
		var val interface***REMOVED******REMOVED***
		if unk.Encoding == proto.WireBytes || unk.Encoding == proto.WireStartGroup ***REMOVED***
			val, err = codec.DecodeLengthDelimitedField(fd, unk.Contents, m.mf)
		***REMOVED*** else ***REMOVED***
			val, err = codec.DecodeScalarField(fd, unk.Value)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if fd.IsMap() ***REMOVED***
			newEntry := val.(*Message)
			kk, err := newEntry.TryGetFieldByNumber(1)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			vv, err := newEntry.TryGetFieldByNumber(2)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			mp[kk] = vv
			v = mp
		***REMOVED*** else if fd.IsRepeated() ***REMOVED***
			t := reflect.TypeOf(val)
			if t.Kind() == reflect.Slice && t != typeOfBytes ***REMOVED***
				// append slices if we unmarshalled a packed repeated field
				newVals := val.([]interface***REMOVED******REMOVED***)
				sl = append(sl, newVals...)
			***REMOVED*** else ***REMOVED***
				sl = append(sl, val)
			***REMOVED***
			v = sl
		***REMOVED*** else ***REMOVED***
			v = val
		***REMOVED***
	***REMOVED***
	m.internalSetField(fd, v)
	return v, nil
***REMOVED***

func validFieldValue(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return validFieldValueForRv(fd, reflect.ValueOf(val))
***REMOVED***

func validFieldValueForRv(fd *desc.FieldDescriptor, val reflect.Value) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if fd.IsMap() && val.Kind() == reflect.Map ***REMOVED***
		return validFieldValueForMapField(fd, val)
	***REMOVED***

	if fd.IsRepeated() ***REMOVED*** // this will also catch map fields where given value was not a map
		if val.Kind() != reflect.Array && val.Kind() != reflect.Slice ***REMOVED***
			if fd.IsMap() ***REMOVED***
				return nil, fmt.Errorf("value for map field must be a map; instead was %v", val.Type())
			***REMOVED*** else ***REMOVED***
				return nil, fmt.Errorf("value for repeated field must be a slice; instead was %v", val.Type())
			***REMOVED***
		***REMOVED***

		if fd.IsMap() ***REMOVED***
			// value should be a slice of entry messages that we need convert into a map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
			m := map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED******REMOVED***
			for i := 0; i < val.Len(); i++ ***REMOVED***
				e, err := validElementFieldValue(fd, val.Index(i).Interface(), false)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				msg := e.(proto.Message)
				dm, err := asDynamicMessage(msg, fd.GetMessageType(), nil)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				k, err := dm.TryGetFieldByNumber(1)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				v, err := dm.TryGetFieldByNumber(2)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				m[k] = v
			***REMOVED***
			return m, nil
		***REMOVED***

		// make a defensive copy while checking contents (also converts to []interface***REMOVED******REMOVED***)
		s := make([]interface***REMOVED******REMOVED***, val.Len())
		for i := 0; i < val.Len(); i++ ***REMOVED***
			ev := val.Index(i)
			if ev.Kind() == reflect.Interface ***REMOVED***
				// unwrap it
				ev = reflect.ValueOf(ev.Interface())
			***REMOVED***
			e, err := validElementFieldValueForRv(fd, ev, false)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			s[i] = e
		***REMOVED***

		return s, nil
	***REMOVED***

	return validElementFieldValueForRv(fd, val, false)
***REMOVED***

func asDynamicMessage(m proto.Message, md *desc.MessageDescriptor, mf *MessageFactory) (*Message, error) ***REMOVED***
	if dm, ok := m.(*Message); ok ***REMOVED***
		return dm, nil
	***REMOVED***
	dm := NewMessageWithMessageFactory(md, mf)
	if err := dm.mergeFrom(m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return dm, nil
***REMOVED***

func validElementFieldValue(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***, allowNilMessage bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return validElementFieldValueForRv(fd, reflect.ValueOf(val), allowNilMessage)
***REMOVED***

func validElementFieldValueForRv(fd *desc.FieldDescriptor, val reflect.Value, allowNilMessage bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	t := fd.GetType()
	if !val.IsValid() ***REMOVED***
		return nil, typeError(fd, nil)
	***REMOVED***

	switch t ***REMOVED***
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		return toInt32(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return toInt64(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32:
		return toUint32(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return toUint64(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return toFloat32(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return toFloat64(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return toBool(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return toBytes(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return toString(reflect.Indirect(val), fd)

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_GROUP:
		m, err := asMessage(val, fd.GetFullyQualifiedName())
		// check that message is correct type
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var msgType string
		if dm, ok := m.(*Message); ok ***REMOVED***
			if allowNilMessage && dm == nil ***REMOVED***
				// if dm == nil, we'll panic below, so early out if that is allowed
				// (only allowed for map values, to indicate an entry w/ no value)
				return m, nil
			***REMOVED***
			msgType = dm.GetMessageDescriptor().GetFullyQualifiedName()
		***REMOVED*** else ***REMOVED***
			msgType = proto.MessageName(m)
		***REMOVED***
		if msgType != fd.GetMessageType().GetFullyQualifiedName() ***REMOVED***
			return nil, fmt.Errorf("message field %s requires value of type %s; received %s", fd.GetFullyQualifiedName(), fd.GetMessageType().GetFullyQualifiedName(), msgType)
		***REMOVED***
		return m, nil

	default:
		return nil, fmt.Errorf("unable to handle unrecognized field type: %v", fd.GetType())
	***REMOVED***
***REMOVED***

func toInt32(v reflect.Value, fd *desc.FieldDescriptor) (int32, error) ***REMOVED***
	if v.Kind() == reflect.Int32 ***REMOVED***
		return int32(v.Int()), nil
	***REMOVED***
	return 0, typeError(fd, v.Type())
***REMOVED***

func toUint32(v reflect.Value, fd *desc.FieldDescriptor) (uint32, error) ***REMOVED***
	if v.Kind() == reflect.Uint32 ***REMOVED***
		return uint32(v.Uint()), nil
	***REMOVED***
	return 0, typeError(fd, v.Type())
***REMOVED***

func toFloat32(v reflect.Value, fd *desc.FieldDescriptor) (float32, error) ***REMOVED***
	if v.Kind() == reflect.Float32 ***REMOVED***
		return float32(v.Float()), nil
	***REMOVED***
	return 0, typeError(fd, v.Type())
***REMOVED***

func toInt64(v reflect.Value, fd *desc.FieldDescriptor) (int64, error) ***REMOVED***
	if v.Kind() == reflect.Int64 || v.Kind() == reflect.Int || v.Kind() == reflect.Int32 ***REMOVED***
		return v.Int(), nil
	***REMOVED***
	return 0, typeError(fd, v.Type())
***REMOVED***

func toUint64(v reflect.Value, fd *desc.FieldDescriptor) (uint64, error) ***REMOVED***
	if v.Kind() == reflect.Uint64 || v.Kind() == reflect.Uint || v.Kind() == reflect.Uint32 ***REMOVED***
		return v.Uint(), nil
	***REMOVED***
	return 0, typeError(fd, v.Type())
***REMOVED***

func toFloat64(v reflect.Value, fd *desc.FieldDescriptor) (float64, error) ***REMOVED***
	if v.Kind() == reflect.Float64 || v.Kind() == reflect.Float32 ***REMOVED***
		return v.Float(), nil
	***REMOVED***
	return 0, typeError(fd, v.Type())
***REMOVED***

func toBool(v reflect.Value, fd *desc.FieldDescriptor) (bool, error) ***REMOVED***
	if v.Kind() == reflect.Bool ***REMOVED***
		return v.Bool(), nil
	***REMOVED***
	return false, typeError(fd, v.Type())
***REMOVED***

func toBytes(v reflect.Value, fd *desc.FieldDescriptor) ([]byte, error) ***REMOVED***
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 ***REMOVED***
		return v.Bytes(), nil
	***REMOVED***
	return nil, typeError(fd, v.Type())
***REMOVED***

func toString(v reflect.Value, fd *desc.FieldDescriptor) (string, error) ***REMOVED***
	if v.Kind() == reflect.String ***REMOVED***
		return v.String(), nil
	***REMOVED***
	return "", typeError(fd, v.Type())
***REMOVED***

func typeError(fd *desc.FieldDescriptor, t reflect.Type) error ***REMOVED***
	return fmt.Errorf(
		"%s field %s is not compatible with value of type %v",
		getTypeString(fd), fd.GetFullyQualifiedName(), t)
***REMOVED***

func getTypeString(fd *desc.FieldDescriptor) string ***REMOVED***
	return strings.ToLower(fd.GetType().String())
***REMOVED***

func asMessage(v reflect.Value, fieldName string) (proto.Message, error) ***REMOVED***
	t := v.Type()
	// we need a pointer to a struct that implements proto.Message
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct || !t.Implements(typeOfProtoMessage) ***REMOVED***
		return nil, fmt.Errorf("message field %s requires is not compatible with value of type %v", fieldName, v.Type())
	***REMOVED***
	return v.Interface().(proto.Message), nil
***REMOVED***

// Reset resets this message to an empty message. It removes all values set in
// the message.
func (m *Message) Reset() ***REMOVED***
	for k := range m.values ***REMOVED***
		delete(m.values, k)
	***REMOVED***
	for k := range m.unknownFields ***REMOVED***
		delete(m.unknownFields, k)
	***REMOVED***
***REMOVED***

// String returns this message rendered in compact text format.
func (m *Message) String() string ***REMOVED***
	b, err := m.MarshalText()
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("Failed to create string representation of message: %s", err.Error()))
	***REMOVED***
	return string(b)
***REMOVED***

// ProtoMessage is present to satisfy the proto.Message interface.
func (m *Message) ProtoMessage() ***REMOVED***
***REMOVED***

// ConvertTo converts this dynamic message into the given message. This is
// shorthand for resetting then merging:
//   target.Reset()
//   m.MergeInto(target)
func (m *Message) ConvertTo(target proto.Message) error ***REMOVED***
	if err := m.checkType(target); err != nil ***REMOVED***
		return err
	***REMOVED***

	target.Reset()
	return m.mergeInto(target, defaultDeterminism)
***REMOVED***

// ConvertToDeterministic converts this dynamic message into the given message.
// It is just like ConvertTo, but it attempts to produce deterministic results.
// That means that if the target is a generated message (not another dynamic
// message) and the current runtime is unaware of any fields or extensions that
// are present in m, they will be serialized into the target's unrecognized
// fields deterministically.
func (m *Message) ConvertToDeterministic(target proto.Message) error ***REMOVED***
	if err := m.checkType(target); err != nil ***REMOVED***
		return err
	***REMOVED***

	target.Reset()
	return m.mergeInto(target, true)
***REMOVED***

// ConvertFrom converts the given message into this dynamic message. This is
// shorthand for resetting then merging:
//   m.Reset()
//   m.MergeFrom(target)
func (m *Message) ConvertFrom(target proto.Message) error ***REMOVED***
	if err := m.checkType(target); err != nil ***REMOVED***
		return err
	***REMOVED***

	m.Reset()
	return m.mergeFrom(target)
***REMOVED***

// MergeInto merges this dynamic message into the given message. All field
// values in this message will be set on the given message. For map fields,
// entries are added to the given message (if the given message has existing
// values for like keys, they are overwritten). For slice fields, elements are
// added.
//
// If the given message has a different set of known fields, it is possible for
// some known fields in this message to be represented as unknown fields in the
// given message after merging, and vice versa.
func (m *Message) MergeInto(target proto.Message) error ***REMOVED***
	if err := m.checkType(target); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.mergeInto(target, defaultDeterminism)
***REMOVED***

// MergeIntoDeterministic merges this dynamic message into the given message.
// It is just like MergeInto, but it attempts to produce deterministic results.
// That means that if the target is a generated message (not another dynamic
// message) and the current runtime is unaware of any fields or extensions that
// are present in m, they will be serialized into the target's unrecognized
// fields deterministically.
func (m *Message) MergeIntoDeterministic(target proto.Message) error ***REMOVED***
	if err := m.checkType(target); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.mergeInto(target, true)
***REMOVED***

// MergeFrom merges the given message into this dynamic message. All field
// values in the given message will be set on this message. For map fields,
// entries are added to this message (if this message has existing values for
// like keys, they are overwritten). For slice fields, elements are added.
//
// If the given message has a different set of known fields, it is possible for
// some known fields in that message to be represented as unknown fields in this
// message after merging, and vice versa.
func (m *Message) MergeFrom(source proto.Message) error ***REMOVED***
	if err := m.checkType(source); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.mergeFrom(source)
***REMOVED***

// Merge implements the proto.Merger interface so that dynamic messages are
// compatible with the proto.Merge function. It delegates to MergeFrom but will
// panic on error as the proto.Merger interface doesn't allow for returning an
// error.
//
// Unlike nearly all other methods, this method can work if this message's type
// is not defined (such as instantiating the message without using NewMessage).
// This is strictly so that dynamic message's are compatible with the
// proto.Clone function, which instantiates a new message via reflection (thus
// its message descriptor will not be set) and than calls Merge.
func (m *Message) Merge(source proto.Message) ***REMOVED***
	if m.md == nil ***REMOVED***
		// To support proto.Clone, initialize the descriptor from the source.
		if dm, ok := source.(*Message); ok ***REMOVED***
			m.md = dm.md
			// also make sure the clone uses the same message factory and
			// extensions and also knows about the same extra fields (if any)
			m.mf = dm.mf
			m.er = dm.er
			m.extraFields = dm.extraFields
		***REMOVED*** else if md, err := desc.LoadMessageDescriptorForMessage(source); err != nil ***REMOVED***
			panic(err.Error())
		***REMOVED*** else ***REMOVED***
			m.md = md
		***REMOVED***
	***REMOVED***

	if err := m.MergeFrom(source); err != nil ***REMOVED***
		panic(err.Error())
	***REMOVED***
***REMOVED***

func (m *Message) checkType(target proto.Message) error ***REMOVED***
	if dm, ok := target.(*Message); ok ***REMOVED***
		if dm.md.GetFullyQualifiedName() != m.md.GetFullyQualifiedName() ***REMOVED***
			return fmt.Errorf("given message has wrong type: %q; expecting %q", dm.md.GetFullyQualifiedName(), m.md.GetFullyQualifiedName())
		***REMOVED***
		return nil
	***REMOVED***

	msgName := proto.MessageName(target)
	if msgName != m.md.GetFullyQualifiedName() ***REMOVED***
		return fmt.Errorf("given message has wrong type: %q; expecting %q", msgName, m.md.GetFullyQualifiedName())
	***REMOVED***
	return nil
***REMOVED***

func (m *Message) mergeInto(pm proto.Message, deterministic bool) error ***REMOVED***
	if dm, ok := pm.(*Message); ok ***REMOVED***
		return dm.mergeFrom(m)
	***REMOVED***

	target := reflect.ValueOf(pm)
	if target.Kind() == reflect.Ptr ***REMOVED***
		target = target.Elem()
	***REMOVED***

	// track tags for which the dynamic message has data but the given
	// message doesn't know about it
	unknownTags := map[int32]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for tag := range m.values ***REMOVED***
		unknownTags[tag] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	// check that we can successfully do the merge
	structProps := proto.GetProperties(reflect.TypeOf(pm).Elem())
	for _, prop := range structProps.Prop ***REMOVED***
		if prop.Tag == 0 ***REMOVED***
			continue // one-of or special field (such as XXX_unrecognized, etc.)
		***REMOVED***
		tag := int32(prop.Tag)
		v, ok := m.values[tag]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if unknownTags != nil ***REMOVED***
			delete(unknownTags, tag)
		***REMOVED***
		f := target.FieldByName(prop.Name)
		ft := f.Type()
		val := reflect.ValueOf(v)
		if !canConvert(val, ft) ***REMOVED***
			return fmt.Errorf("cannot convert %v to %v", val.Type(), ft)
		***REMOVED***
	***REMOVED***
	// check one-of fields
	for _, oop := range structProps.OneofTypes ***REMOVED***
		prop := oop.Prop
		tag := int32(prop.Tag)
		v, ok := m.values[tag]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if unknownTags != nil ***REMOVED***
			delete(unknownTags, tag)
		***REMOVED***
		stf, ok := oop.Type.Elem().FieldByName(prop.Name)
		if !ok ***REMOVED***
			return fmt.Errorf("one-of field indicates struct field name %s, but type %v has no such field", prop.Name, oop.Type.Elem())
		***REMOVED***
		ft := stf.Type
		val := reflect.ValueOf(v)
		if !canConvert(val, ft) ***REMOVED***
			return fmt.Errorf("cannot convert %v to %v", val.Type(), ft)
		***REMOVED***
	***REMOVED***
	// and check extensions, too
	for tag, ext := range proto.RegisteredExtensions(pm) ***REMOVED***
		v, ok := m.values[tag]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if unknownTags != nil ***REMOVED***
			delete(unknownTags, tag)
		***REMOVED***
		ft := reflect.TypeOf(ext.ExtensionType)
		val := reflect.ValueOf(v)
		if !canConvert(val, ft) ***REMOVED***
			return fmt.Errorf("cannot convert %v to %v", val.Type(), ft)
		***REMOVED***
	***REMOVED***

	// now actually perform the merge
	for _, prop := range structProps.Prop ***REMOVED***
		v, ok := m.values[int32(prop.Tag)]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		f := target.FieldByName(prop.Name)
		if err := mergeVal(reflect.ValueOf(v), f, deterministic); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// merge one-ofs
	for _, oop := range structProps.OneofTypes ***REMOVED***
		prop := oop.Prop
		tag := int32(prop.Tag)
		v, ok := m.values[tag]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		oov := reflect.New(oop.Type.Elem())
		f := oov.Elem().FieldByName(prop.Name)
		if err := mergeVal(reflect.ValueOf(v), f, deterministic); err != nil ***REMOVED***
			return err
		***REMOVED***
		target.Field(oop.Field).Set(oov)
	***REMOVED***
	// merge extensions, too
	for tag, ext := range proto.RegisteredExtensions(pm) ***REMOVED***
		v, ok := m.values[tag]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		e := reflect.New(reflect.TypeOf(ext.ExtensionType)).Elem()
		if err := mergeVal(reflect.ValueOf(v), e, deterministic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := proto.SetExtension(pm, ext, e.Interface()); err != nil ***REMOVED***
			// shouldn't happen since we already checked that the extension type was compatible above
			return err
		***REMOVED***
	***REMOVED***

	// if we have fields that the given message doesn't know about, add to its unknown fields
	if len(unknownTags) > 0 ***REMOVED***
		var b codec.Buffer
		b.SetDeterministic(deterministic)
		if deterministic ***REMOVED***
			// if we need to emit things deterministically, sort the
			// extensions by their tag number
			sortedUnknownTags := make([]int32, 0, len(unknownTags))
			for tag := range unknownTags ***REMOVED***
				sortedUnknownTags = append(sortedUnknownTags, tag)
			***REMOVED***
			sort.Slice(sortedUnknownTags, func(i, j int) bool ***REMOVED***
				return sortedUnknownTags[i] < sortedUnknownTags[j]
			***REMOVED***)
			for _, tag := range sortedUnknownTags ***REMOVED***
				fd := m.FindFieldDescriptor(tag)
				if err := b.EncodeFieldValue(fd, m.values[tag]); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for tag := range unknownTags ***REMOVED***
				fd := m.FindFieldDescriptor(tag)
				if err := b.EncodeFieldValue(fd, m.values[tag]); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		internal.SetUnrecognized(pm, b.Bytes())
	***REMOVED***

	// finally, convey unknown fields into the given message by letting it unmarshal them
	// (this will append to its unknown fields if not known; if somehow the given message recognizes
	// a field even though the dynamic message did not, it will get correctly unmarshalled)
	if unknownTags != nil && len(m.unknownFields) > 0 ***REMOVED***
		var b codec.Buffer
		_ = m.marshalUnknownFields(&b)
		_ = proto.UnmarshalMerge(b.Bytes(), pm)
	***REMOVED***

	return nil
***REMOVED***

func canConvert(src reflect.Value, target reflect.Type) bool ***REMOVED***
	if src.Kind() == reflect.Interface ***REMOVED***
		src = reflect.ValueOf(src.Interface())
	***REMOVED***
	srcType := src.Type()
	// we allow convertible types instead of requiring exact types so that calling
	// code can, for example, assign an enum constant to an enum field. In that case,
	// one type is the enum type (a sub-type of int32) and the other may be the int32
	// type. So we automatically do the conversion in that case.
	if srcType.ConvertibleTo(target) ***REMOVED***
		return true
	***REMOVED*** else if target.Kind() == reflect.Ptr && srcType.ConvertibleTo(target.Elem()) ***REMOVED***
		return true
	***REMOVED*** else if target.Kind() == reflect.Slice ***REMOVED***
		if srcType.Kind() != reflect.Slice ***REMOVED***
			return false
		***REMOVED***
		et := target.Elem()
		for i := 0; i < src.Len(); i++ ***REMOVED***
			if !canConvert(src.Index(i), et) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED*** else if target.Kind() == reflect.Map ***REMOVED***
		if srcType.Kind() != reflect.Map ***REMOVED***
			return false
		***REMOVED***
		return canConvertMap(src, target)
	***REMOVED*** else if srcType == typeOfDynamicMessage && target.Implements(typeOfProtoMessage) ***REMOVED***
		z := reflect.Zero(target).Interface()
		msgType := proto.MessageName(z.(proto.Message))
		return msgType == src.Interface().(*Message).GetMessageDescriptor().GetFullyQualifiedName()
	***REMOVED*** else ***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func mergeVal(src, target reflect.Value, deterministic bool) error ***REMOVED***
	if src.Kind() == reflect.Interface && !src.IsNil() ***REMOVED***
		src = src.Elem()
	***REMOVED***
	srcType := src.Type()
	targetType := target.Type()
	if srcType.ConvertibleTo(targetType) ***REMOVED***
		if targetType.Implements(typeOfProtoMessage) && !target.IsNil() ***REMOVED***
			Merge(target.Interface().(proto.Message), src.Convert(targetType).Interface().(proto.Message))
		***REMOVED*** else ***REMOVED***
			target.Set(src.Convert(targetType))
		***REMOVED***
	***REMOVED*** else if targetType.Kind() == reflect.Ptr && srcType.ConvertibleTo(targetType.Elem()) ***REMOVED***
		if !src.CanAddr() ***REMOVED***
			target.Set(reflect.New(targetType.Elem()))
			target.Elem().Set(src.Convert(targetType.Elem()))
		***REMOVED*** else ***REMOVED***
			target.Set(src.Addr().Convert(targetType))
		***REMOVED***
	***REMOVED*** else if targetType.Kind() == reflect.Slice ***REMOVED***
		l := target.Len()
		newL := l + src.Len()
		if target.Cap() < newL ***REMOVED***
			// expand capacity of the slice and copy
			newSl := reflect.MakeSlice(targetType, newL, newL)
			for i := 0; i < target.Len(); i++ ***REMOVED***
				newSl.Index(i).Set(target.Index(i))
			***REMOVED***
			target.Set(newSl)
		***REMOVED*** else ***REMOVED***
			target.SetLen(newL)
		***REMOVED***
		for i := 0; i < src.Len(); i++ ***REMOVED***
			dest := target.Index(l + i)
			if dest.Kind() == reflect.Ptr ***REMOVED***
				dest.Set(reflect.New(dest.Type().Elem()))
			***REMOVED***
			if err := mergeVal(src.Index(i), dest, deterministic); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if targetType.Kind() == reflect.Map ***REMOVED***
		return mergeMapVal(src, target, targetType, deterministic)
	***REMOVED*** else if srcType == typeOfDynamicMessage && targetType.Implements(typeOfProtoMessage) ***REMOVED***
		dm := src.Interface().(*Message)
		if target.IsNil() ***REMOVED***
			target.Set(reflect.New(targetType.Elem()))
		***REMOVED***
		m := target.Interface().(proto.Message)
		if err := dm.mergeInto(m, deterministic); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return fmt.Errorf("cannot convert %v to %v", srcType, targetType)
	***REMOVED***
	return nil
***REMOVED***

func (m *Message) mergeFrom(pm proto.Message) error ***REMOVED***
	if dm, ok := pm.(*Message); ok ***REMOVED***
		// if given message is also a dynamic message, we merge differently
		for tag, v := range dm.values ***REMOVED***
			fd := m.FindFieldDescriptor(tag)
			if fd == nil ***REMOVED***
				fd = dm.FindFieldDescriptor(tag)
			***REMOVED***
			if err := mergeField(m, fd, v); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	pmrv := reflect.ValueOf(pm)
	if pmrv.IsNil() ***REMOVED***
		// nil is an empty message, so nothing to do
		return nil
	***REMOVED***

	// check that we can successfully do the merge
	src := pmrv.Elem()
	values := map[*desc.FieldDescriptor]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	props := proto.GetProperties(reflect.TypeOf(pm).Elem())
	if props == nil ***REMOVED***
		return fmt.Errorf("could not determine message properties to merge for %v", reflect.TypeOf(pm).Elem())
	***REMOVED***

	// regular fields
	for _, prop := range props.Prop ***REMOVED***
		if prop.Tag == 0 ***REMOVED***
			continue // one-of or special field (such as XXX_unrecognized, etc.)
		***REMOVED***
		fd := m.FindFieldDescriptor(int32(prop.Tag))
		if fd == nil ***REMOVED***
			// Our descriptor has different fields than this message object. So
			// try to reflect on the message object's fields.
			md, err := desc.LoadMessageDescriptorForMessage(pm)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			fd = md.FindFieldByNumber(int32(prop.Tag))
			if fd == nil ***REMOVED***
				return fmt.Errorf("message descriptor %q did not contain field for tag %d (%q)", md.GetFullyQualifiedName(), prop.Tag, prop.Name)
			***REMOVED***
		***REMOVED***
		rv := src.FieldByName(prop.Name)
		if (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Slice) && rv.IsNil() ***REMOVED***
			continue
		***REMOVED***
		if v, err := validFieldValueForRv(fd, rv); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			values[fd] = v
		***REMOVED***
	***REMOVED***

	// one-of fields
	for _, oop := range props.OneofTypes ***REMOVED***
		oov := src.Field(oop.Field).Elem()
		if !oov.IsValid() || oov.Type() != oop.Type ***REMOVED***
			// this field is unset (in other words, one-of message field is not currently set to this option)
			continue
		***REMOVED***
		prop := oop.Prop
		rv := oov.Elem().FieldByName(prop.Name)
		fd := m.FindFieldDescriptor(int32(prop.Tag))
		if fd == nil ***REMOVED***
			// Our descriptor has different fields than this message object. So
			// try to reflect on the message object's fields.
			md, err := desc.LoadMessageDescriptorForMessage(pm)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			fd = md.FindFieldByNumber(int32(prop.Tag))
			if fd == nil ***REMOVED***
				return fmt.Errorf("message descriptor %q did not contain field for tag %d (%q in one-of %q)", md.GetFullyQualifiedName(), prop.Tag, prop.Name, src.Type().Field(oop.Field).Name)
			***REMOVED***
		***REMOVED***
		if v, err := validFieldValueForRv(fd, rv); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			values[fd] = v
		***REMOVED***
	***REMOVED***

	// extension fields
	rexts, _ := proto.ExtensionDescs(pm)
	for _, ed := range rexts ***REMOVED***
		v, _ := proto.GetExtension(pm, ed)
		if v == nil ***REMOVED***
			continue
		***REMOVED***
		if ed.ExtensionType == nil ***REMOVED***
			// unrecognized extension: we'll handle that below when we
			// handle other unrecognized fields
			continue
		***REMOVED***
		fd := m.er.FindExtension(m.md.GetFullyQualifiedName(), ed.Field)
		if fd == nil ***REMOVED***
			var err error
			if fd, err = desc.LoadFieldDescriptorForExtension(ed); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if v, err := validFieldValue(fd, v); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			values[fd] = v
		***REMOVED***
	***REMOVED***

	// now actually perform the merge
	for fd, v := range values ***REMOVED***
		if err := mergeField(m, fd, v); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	data := internal.GetUnrecognized(pm)
	if len(data) > 0 ***REMOVED***
		// ignore any error returned: pulling in unknown fields is best-effort
		_ = m.UnmarshalMerge(data)
	***REMOVED***

	return nil
***REMOVED***

// Validate checks that all required fields are present. It returns an error if any are absent.
func (m *Message) Validate() error ***REMOVED***
	missingFields := m.findMissingFields()
	if len(missingFields) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return fmt.Errorf("some required fields missing: %v", strings.Join(missingFields, ", "))
***REMOVED***

func (m *Message) findMissingFields() []string ***REMOVED***
	if m.md.IsProto3() ***REMOVED***
		// proto3 does not allow required fields
		return nil
	***REMOVED***
	var missingFields []string
	for _, fd := range m.md.GetFields() ***REMOVED***
		if fd.IsRequired() ***REMOVED***
			if _, ok := m.values[fd.GetNumber()]; !ok ***REMOVED***
				missingFields = append(missingFields, fd.GetName())
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return missingFields
***REMOVED***

// ValidateRecursive checks that all required fields are present and also
// recursively validates all fields who are also messages. It returns an error
// if any required fields, in this message or nested within, are absent.
func (m *Message) ValidateRecursive() error ***REMOVED***
	return m.validateRecursive("")
***REMOVED***

func (m *Message) validateRecursive(prefix string) error ***REMOVED***
	if missingFields := m.findMissingFields(); len(missingFields) > 0 ***REMOVED***
		for i := range missingFields ***REMOVED***
			missingFields[i] = fmt.Sprintf("%s%s", prefix, missingFields[i])
		***REMOVED***
		return fmt.Errorf("some required fields missing: %v", strings.Join(missingFields, ", "))
	***REMOVED***

	for tag, fld := range m.values ***REMOVED***
		fd := m.FindFieldDescriptor(tag)
		var chprefix string
		var md *desc.MessageDescriptor
		checkMsg := func(pm proto.Message) error ***REMOVED***
			var dm *Message
			if d, ok := pm.(*Message); ok ***REMOVED***
				dm = d
			***REMOVED*** else if pm != nil ***REMOVED***
				dm = m.mf.NewDynamicMessage(md)
				if err := dm.ConvertFrom(pm); err != nil ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
			if dm == nil ***REMOVED***
				return nil
			***REMOVED***
			if err := dm.validateRecursive(chprefix); err != nil ***REMOVED***
				return err
			***REMOVED***
			return nil
		***REMOVED***
		isMap := fd.IsMap()
		if isMap && fd.GetMapValueType().GetMessageType() != nil ***REMOVED***
			md = fd.GetMapValueType().GetMessageType()
			mp := fld.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
			for k, v := range mp ***REMOVED***
				chprefix = fmt.Sprintf("%s%s[%v].", prefix, getName(fd), k)
				if err := checkMsg(v.(proto.Message)); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if !isMap && fd.GetMessageType() != nil ***REMOVED***
			md = fd.GetMessageType()
			if fd.IsRepeated() ***REMOVED***
				sl := fld.([]interface***REMOVED******REMOVED***)
				for i, v := range sl ***REMOVED***
					chprefix = fmt.Sprintf("%s%s[%d].", prefix, getName(fd), i)
					if err := checkMsg(v.(proto.Message)); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				chprefix = fmt.Sprintf("%s%s.", prefix, getName(fd))
				if err := checkMsg(fld.(proto.Message)); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func getName(fd *desc.FieldDescriptor) string ***REMOVED***
	if fd.IsExtension() ***REMOVED***
		return fmt.Sprintf("(%s)", fd.GetFullyQualifiedName())
	***REMOVED*** else ***REMOVED***
		return fd.GetName()
	***REMOVED***
***REMOVED***

// knownFieldTags return tags of present and recognized fields, in sorted order.
func (m *Message) knownFieldTags() []int ***REMOVED***
	if len(m.values) == 0 ***REMOVED***
		return []int(nil)
	***REMOVED***

	keys := make([]int, len(m.values))
	i := 0
	for k := range m.values ***REMOVED***
		keys[i] = int(k)
		i++
	***REMOVED***

	sort.Ints(keys)
	return keys
***REMOVED***

// allKnownFieldTags return tags of present and recognized fields, including
// those that are unset, in sorted order. This only includes extensions that are
// present. Known but not-present extensions are not included in the returned
// set of tags.
func (m *Message) allKnownFieldTags() []int ***REMOVED***
	fds := m.md.GetFields()
	keys := make([]int, 0, len(fds)+len(m.extraFields))

	for k := range m.values ***REMOVED***
		keys = append(keys, int(k))
	***REMOVED***

	// also include known fields that are not present
	for _, fd := range fds ***REMOVED***
		if _, ok := m.values[fd.GetNumber()]; !ok ***REMOVED***
			keys = append(keys, int(fd.GetNumber()))
		***REMOVED***
	***REMOVED***
	for _, fd := range m.extraFields ***REMOVED***
		if !fd.IsExtension() ***REMOVED*** // skip extensions that are not present
			if _, ok := m.values[fd.GetNumber()]; !ok ***REMOVED***
				keys = append(keys, int(fd.GetNumber()))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	sort.Ints(keys)
	return keys
***REMOVED***

// unknownFieldTags return tags of present but unrecognized fields, in sorted order.
func (m *Message) unknownFieldTags() []int ***REMOVED***
	if len(m.unknownFields) == 0 ***REMOVED***
		return []int(nil)
	***REMOVED***
	keys := make([]int, len(m.unknownFields))
	i := 0
	for k := range m.unknownFields ***REMOVED***
		keys[i] = int(k)
		i++
	***REMOVED***
	sort.Ints(keys)
	return keys
***REMOVED***
