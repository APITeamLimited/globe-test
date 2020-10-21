package codec

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc"
)

// EncodeZigZag64 does zig-zag encoding to convert the given
// signed 64-bit integer into a form that can be expressed
// efficiently as a varint, even for negative values.
func EncodeZigZag64(v int64) uint64 ***REMOVED***
	return (uint64(v) << 1) ^ uint64(v>>63)
***REMOVED***

// EncodeZigZag32 does zig-zag encoding to convert the given
// signed 32-bit integer into a form that can be expressed
// efficiently as a varint, even for negative values.
func EncodeZigZag32(v int32) uint64 ***REMOVED***
	return uint64((uint32(v) << 1) ^ uint32((v >> 31)))
***REMOVED***

func (cb *Buffer) EncodeFieldValue(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if fd.IsMap() ***REMOVED***
		mp := val.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
		entryType := fd.GetMessageType()
		keyType := entryType.FindFieldByNumber(1)
		valType := entryType.FindFieldByNumber(2)
		var entryBuffer Buffer
		if cb.IsDeterministic() ***REMOVED***
			keys := make([]interface***REMOVED******REMOVED***, 0, len(mp))
			for k := range mp ***REMOVED***
				keys = append(keys, k)
			***REMOVED***
			sort.Sort(sortable(keys))
			for _, k := range keys ***REMOVED***
				v := mp[k]
				entryBuffer.Reset()
				if err := entryBuffer.encodeFieldElement(keyType, k); err != nil ***REMOVED***
					return err
				***REMOVED***
				rv := reflect.ValueOf(v)
				if rv.Kind() != reflect.Ptr || !rv.IsNil() ***REMOVED***
					if err := entryBuffer.encodeFieldElement(valType, v); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				if err := cb.EncodeTagAndWireType(fd.GetNumber(), proto.WireBytes); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := cb.EncodeRawBytes(entryBuffer.Bytes()); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for k, v := range mp ***REMOVED***
				entryBuffer.Reset()
				if err := entryBuffer.encodeFieldElement(keyType, k); err != nil ***REMOVED***
					return err
				***REMOVED***
				rv := reflect.ValueOf(v)
				if rv.Kind() != reflect.Ptr || !rv.IsNil() ***REMOVED***
					if err := entryBuffer.encodeFieldElement(valType, v); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				if err := cb.EncodeTagAndWireType(fd.GetNumber(), proto.WireBytes); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := cb.EncodeRawBytes(entryBuffer.Bytes()); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED*** else if fd.IsRepeated() ***REMOVED***
		sl := val.([]interface***REMOVED******REMOVED***)
		wt, err := getWireType(fd.GetType())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if isPacked(fd) && len(sl) > 0 &&
			(wt == proto.WireVarint || wt == proto.WireFixed32 || wt == proto.WireFixed64) ***REMOVED***
			// packed repeated field
			var packedBuffer Buffer
			for _, v := range sl ***REMOVED***
				if err := packedBuffer.encodeFieldValue(fd, v); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			if err := cb.EncodeTagAndWireType(fd.GetNumber(), proto.WireBytes); err != nil ***REMOVED***
				return err
			***REMOVED***
			return cb.EncodeRawBytes(packedBuffer.Bytes())
		***REMOVED*** else ***REMOVED***
			// non-packed repeated field
			for _, v := range sl ***REMOVED***
				if err := cb.encodeFieldElement(fd, v); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return cb.encodeFieldElement(fd, val)
	***REMOVED***
***REMOVED***

func isPacked(fd *desc.FieldDescriptor) bool ***REMOVED***
	opts := fd.AsFieldDescriptorProto().GetOptions()
	// if set, use that value
	if opts != nil && opts.Packed != nil ***REMOVED***
		return opts.GetPacked()
	***REMOVED***
	// if unset: proto2 defaults to false, proto3 to true
	return fd.GetFile().IsProto3()
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

func (b *Buffer) encodeFieldElement(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	wt, err := getWireType(fd.GetType())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.EncodeTagAndWireType(fd.GetNumber(), wt); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.encodeFieldValue(fd, val); err != nil ***REMOVED***
		return err
	***REMOVED***
	if wt == proto.WireStartGroup ***REMOVED***
		return b.EncodeTagAndWireType(fd.GetNumber(), proto.WireEndGroup)
	***REMOVED***
	return nil
***REMOVED***

func (b *Buffer) encodeFieldValue(fd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	switch fd.GetType() ***REMOVED***
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		v := val.(bool)
		if v ***REMOVED***
			return b.EncodeVarint(1)
		***REMOVED***
		return b.EncodeVarint(0)

	case descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_INT32:
		v := val.(int32)
		return b.EncodeVarint(uint64(v))

	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		v := val.(int32)
		return b.EncodeFixed32(uint64(v))

	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		v := val.(int32)
		return b.EncodeVarint(EncodeZigZag32(v))

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		v := val.(uint32)
		return b.EncodeVarint(uint64(v))

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		v := val.(uint32)
		return b.EncodeFixed32(uint64(v))

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		v := val.(int64)
		return b.EncodeVarint(uint64(v))

	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		v := val.(int64)
		return b.EncodeFixed64(uint64(v))

	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		v := val.(int64)
		return b.EncodeVarint(EncodeZigZag64(v))

	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		v := val.(uint64)
		return b.EncodeVarint(v)

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		v := val.(uint64)
		return b.EncodeFixed64(v)

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		v := val.(float64)
		return b.EncodeFixed64(math.Float64bits(v))

	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		v := val.(float32)
		return b.EncodeFixed32(uint64(math.Float32bits(v)))

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		v := val.([]byte)
		return b.EncodeRawBytes(v)

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		v := val.(string)
		return b.EncodeRawBytes(([]byte)(v))

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return b.EncodeDelimitedMessage(val.(proto.Message))

	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		// just append the nested message to this buffer
		return b.EncodeMessage(val.(proto.Message))
		// whosoever writeth start-group tag (e.g. caller) is responsible for writing end-group tag

	default:
		return fmt.Errorf("unrecognized field type: %v", fd.GetType())
	***REMOVED***
***REMOVED***

func getWireType(t descriptor.FieldDescriptorProto_Type) (int8, error) ***REMOVED***
	switch t ***REMOVED***
	case descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return proto.WireVarint, nil

	case descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return proto.WireFixed32, nil

	case descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return proto.WireFixed64, nil

	case descriptor.FieldDescriptorProto_TYPE_BYTES,
		descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return proto.WireBytes, nil

	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return proto.WireStartGroup, nil

	default:
		return 0, ErrBadWireType
	***REMOVED***
***REMOVED***
