package codec

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc"
)

var varintTypes = map[descriptor.FieldDescriptorProto_Type]bool***REMOVED******REMOVED***
var fixed32Types = map[descriptor.FieldDescriptorProto_Type]bool***REMOVED******REMOVED***
var fixed64Types = map[descriptor.FieldDescriptorProto_Type]bool***REMOVED******REMOVED***

func init() ***REMOVED***
	varintTypes[descriptor.FieldDescriptorProto_TYPE_BOOL] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_INT32] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_INT64] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_UINT32] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_UINT64] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_SINT32] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_SINT64] = true
	varintTypes[descriptor.FieldDescriptorProto_TYPE_ENUM] = true

	fixed32Types[descriptor.FieldDescriptorProto_TYPE_FIXED32] = true
	fixed32Types[descriptor.FieldDescriptorProto_TYPE_SFIXED32] = true
	fixed32Types[descriptor.FieldDescriptorProto_TYPE_FLOAT] = true

	fixed64Types[descriptor.FieldDescriptorProto_TYPE_FIXED64] = true
	fixed64Types[descriptor.FieldDescriptorProto_TYPE_SFIXED64] = true
	fixed64Types[descriptor.FieldDescriptorProto_TYPE_DOUBLE] = true
***REMOVED***

// ErrWireTypeEndGroup is returned from DecodeFieldValue if the tag and wire-type
// it reads indicates an end-group marker.
var ErrWireTypeEndGroup = errors.New("unexpected wire type: end group")

// MessageFactory is used to instantiate messages when DecodeFieldValue needs to
// decode a message value.
//
// Also see MessageFactory in "github.com/jhump/protoreflect/dynamic", which
// implements this interface.
type MessageFactory interface ***REMOVED***
	NewMessage(md *desc.MessageDescriptor) proto.Message
***REMOVED***

// UnknownField represents a field that was parsed from the binary wire
// format for a message, but was not a recognized field number. Enough
// information is preserved so that re-serializing the message won't lose
// any of the unrecognized data.
type UnknownField struct ***REMOVED***
	// The tag number for the unrecognized field.
	Tag int32

	// Encoding indicates how the unknown field was encoded on the wire. If it
	// is proto.WireBytes or proto.WireGroupStart then Contents will be set to
	// the raw bytes. If it is proto.WireTypeFixed32 then the data is in the least
	// significant 32 bits of Value. Otherwise, the data is in all 64 bits of
	// Value.
	Encoding int8
	Contents []byte
	Value    uint64
***REMOVED***

// DecodeZigZag32 decodes a signed 32-bit integer from the given
// zig-zag encoded value.
func DecodeZigZag32(v uint64) int32 ***REMOVED***
	return int32((uint32(v) >> 1) ^ uint32((int32(v&1)<<31)>>31))
***REMOVED***

// DecodeZigZag64 decodes a signed 64-bit integer from the given
// zig-zag encoded value.
func DecodeZigZag64(v uint64) int64 ***REMOVED***
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63))
***REMOVED***

// DecodeFieldValue will read a field value from the buffer and return its
// value and the corresponding field descriptor. The given function is used
// to lookup a field descriptor by tag number. The given factory is used to
// instantiate a message if the field value is (or contains) a message value.
//
// On error, the field descriptor and value are typically nil. However, if the
// error returned is ErrWireTypeEndGroup, the returned value will indicate any
// tag number encoded in the end-group marker.
//
// If the field descriptor returned is nil, that means that the given function
// returned nil. This is expected to happen for unrecognized tag numbers. In
// that case, no error is returned, and the value will be an UnknownField.
func (cb *Buffer) DecodeFieldValue(fieldFinder func(int32) *desc.FieldDescriptor, fact MessageFactory) (*desc.FieldDescriptor, interface***REMOVED******REMOVED***, error) ***REMOVED***
	if cb.EOF() ***REMOVED***
		return nil, nil, io.EOF
	***REMOVED***
	tagNumber, wireType, err := cb.DecodeTagAndWireType()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	if wireType == proto.WireEndGroup ***REMOVED***
		return nil, tagNumber, ErrWireTypeEndGroup
	***REMOVED***
	fd := fieldFinder(tagNumber)
	if fd == nil ***REMOVED***
		val, err := cb.decodeUnknownField(tagNumber, wireType)
		return nil, val, err
	***REMOVED***
	val, err := cb.decodeKnownField(fd, wireType, fact)
	return fd, val, err
***REMOVED***

// DecodeScalarField extracts a properly-typed value from v. The returned value's
// type depends on the given field descriptor type. It will be the same type as
// generated structs use for the field descriptor's type. Enum types will return
// an int32. If the given field type uses length-delimited encoding (nested
// messages, bytes, and strings), an error is returned.
func DecodeScalarField(fd *desc.FieldDescriptor, v uint64) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch fd.GetType() ***REMOVED***
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return v != 0, nil
	case descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32:
		if v > math.MaxUint32 ***REMOVED***
			return nil, ErrOverflow
		***REMOVED***
		return uint32(v), nil

	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM:
		s := int64(v)
		if s > math.MaxInt32 || s < math.MinInt32 ***REMOVED***
			return nil, ErrOverflow
		***REMOVED***
		return int32(s), nil

	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		if v > math.MaxUint32 ***REMOVED***
			return nil, ErrOverflow
		***REMOVED***
		return int32(v), nil

	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		if v > math.MaxUint32 ***REMOVED***
			return nil, ErrOverflow
		***REMOVED***
		return DecodeZigZag32(v), nil

	case descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return v, nil

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return int64(v), nil

	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return DecodeZigZag64(v), nil

	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if v > math.MaxUint32 ***REMOVED***
			return nil, ErrOverflow
		***REMOVED***
		return math.Float32frombits(uint32(v)), nil

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return math.Float64frombits(v), nil

	default:
		// bytes, string, message, and group cannot be represented as a simple numeric value
		return nil, fmt.Errorf("bad input; field %s requires length-delimited wire type", fd.GetFullyQualifiedName())
	***REMOVED***
***REMOVED***

// DecodeLengthDelimitedField extracts a properly-typed value from bytes. The
// returned value's type will usually be []byte, string, or, for nested messages,
// the type returned from the given message factory. However, since repeated
// scalar fields can be length-delimited, when they used packed encoding, it can
// also return an []interface***REMOVED******REMOVED***, where each element is a scalar value. Furthermore,
// it could return a scalar type, not in a slice, if the given field descriptor is
// not repeated. This is to support cases where a field is changed from optional
// to repeated. New code may emit a packed repeated representation, but old code
// still expects a single scalar value. In this case, if the actual data in bytes
// contains multiple values, only the last value is returned.
func DecodeLengthDelimitedField(fd *desc.FieldDescriptor, bytes []byte, mf MessageFactory) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch ***REMOVED***
	case fd.GetType() == descriptor.FieldDescriptorProto_TYPE_BYTES:
		return bytes, nil

	case fd.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING:
		return string(bytes), nil

	case fd.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE ||
		fd.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP:
		msg := mf.NewMessage(fd.GetMessageType())
		err := proto.Unmarshal(bytes, msg)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else ***REMOVED***
			return msg, nil
		***REMOVED***

	default:
		// even if the field is not repeated or not packed, we still parse it as such for
		// backwards compatibility (e.g. message we are de-serializing could have been both
		// repeated and packed at the time of serialization)
		packedBuf := NewBuffer(bytes)
		var slice []interface***REMOVED******REMOVED***
		var val interface***REMOVED******REMOVED***
		for !packedBuf.EOF() ***REMOVED***
			var v uint64
			var err error
			if varintTypes[fd.GetType()] ***REMOVED***
				v, err = packedBuf.DecodeVarint()
			***REMOVED*** else if fixed32Types[fd.GetType()] ***REMOVED***
				v, err = packedBuf.DecodeFixed32()
			***REMOVED*** else if fixed64Types[fd.GetType()] ***REMOVED***
				v, err = packedBuf.DecodeFixed64()
			***REMOVED*** else ***REMOVED***
				return nil, fmt.Errorf("bad input; cannot parse length-delimited wire type for field %s", fd.GetFullyQualifiedName())
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			val, err = DecodeScalarField(fd, v)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if fd.IsRepeated() ***REMOVED***
				slice = append(slice, val)
			***REMOVED***
		***REMOVED***
		if fd.IsRepeated() ***REMOVED***
			return slice, nil
		***REMOVED*** else ***REMOVED***
			// if not a repeated field, last value wins
			return val, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (b *Buffer) decodeKnownField(fd *desc.FieldDescriptor, encoding int8, fact MessageFactory) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	var val interface***REMOVED******REMOVED***
	var err error
	switch encoding ***REMOVED***
	case proto.WireFixed32:
		var num uint64
		num, err = b.DecodeFixed32()
		if err == nil ***REMOVED***
			val, err = DecodeScalarField(fd, num)
		***REMOVED***
	case proto.WireFixed64:
		var num uint64
		num, err = b.DecodeFixed64()
		if err == nil ***REMOVED***
			val, err = DecodeScalarField(fd, num)
		***REMOVED***
	case proto.WireVarint:
		var num uint64
		num, err = b.DecodeVarint()
		if err == nil ***REMOVED***
			val, err = DecodeScalarField(fd, num)
		***REMOVED***

	case proto.WireBytes:
		alloc := fd.GetType() == descriptor.FieldDescriptorProto_TYPE_BYTES
		var raw []byte
		raw, err = b.DecodeRawBytes(alloc)
		if err == nil ***REMOVED***
			val, err = DecodeLengthDelimitedField(fd, raw, fact)
		***REMOVED***

	case proto.WireStartGroup:
		if fd.GetMessageType() == nil ***REMOVED***
			return nil, fmt.Errorf("cannot parse field %s from group-encoded wire type", fd.GetFullyQualifiedName())
		***REMOVED***
		msg := fact.NewMessage(fd.GetMessageType())
		var data []byte
		data, err = b.ReadGroup(false)
		if err == nil ***REMOVED***
			err = proto.Unmarshal(data, msg)
			if err == nil ***REMOVED***
				val = msg
			***REMOVED***
		***REMOVED***

	default:
		return nil, ErrBadWireType
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return val, nil
***REMOVED***

func (b *Buffer) decodeUnknownField(tagNumber int32, encoding int8) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	u := UnknownField***REMOVED***Tag: tagNumber, Encoding: encoding***REMOVED***
	var err error
	switch encoding ***REMOVED***
	case proto.WireFixed32:
		u.Value, err = b.DecodeFixed32()
	case proto.WireFixed64:
		u.Value, err = b.DecodeFixed64()
	case proto.WireVarint:
		u.Value, err = b.DecodeVarint()
	case proto.WireBytes:
		u.Contents, err = b.DecodeRawBytes(true)
	case proto.WireStartGroup:
		u.Contents, err = b.ReadGroup(true)
	default:
		err = ErrBadWireType
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return u, nil
***REMOVED***
