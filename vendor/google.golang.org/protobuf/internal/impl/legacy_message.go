// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/internal/descopts"
	ptag "google.golang.org/protobuf/internal/encoding/tag"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

// legacyWrapMessage wraps v as a protoreflect.Message,
// where v must be a *struct kind and not implement the v2 API already.
func legacyWrapMessage(v reflect.Value) pref.Message ***REMOVED***
	typ := v.Type()
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct ***REMOVED***
		return aberrantMessage***REMOVED***v: v***REMOVED***
	***REMOVED***
	mt := legacyLoadMessageInfo(typ, "")
	return mt.MessageOf(v.Interface())
***REMOVED***

// legacyLoadMessageType dynamically loads a protoreflect.Type for t,
// where t must be not implement the v2 API already.
// The provided name is used if it cannot be determined from the message.
func legacyLoadMessageType(t reflect.Type, name pref.FullName) protoreflect.MessageType ***REMOVED***
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct ***REMOVED***
		return aberrantMessageType***REMOVED***t***REMOVED***
	***REMOVED***
	return legacyLoadMessageInfo(t, name)
***REMOVED***

var legacyMessageTypeCache sync.Map // map[reflect.Type]*MessageInfo

// legacyLoadMessageInfo dynamically loads a *MessageInfo for t,
// where t must be a *struct kind and not implement the v2 API already.
// The provided name is used if it cannot be determined from the message.
func legacyLoadMessageInfo(t reflect.Type, name pref.FullName) *MessageInfo ***REMOVED***
	// Fast-path: check if a MessageInfo is cached for this concrete type.
	if mt, ok := legacyMessageTypeCache.Load(t); ok ***REMOVED***
		return mt.(*MessageInfo)
	***REMOVED***

	// Slow-path: derive message descriptor and initialize MessageInfo.
	mi := &MessageInfo***REMOVED***
		Desc:          legacyLoadMessageDesc(t, name),
		GoReflectType: t,
	***REMOVED***

	v := reflect.Zero(t).Interface()
	if _, ok := v.(legacyMarshaler); ok ***REMOVED***
		mi.methods.Marshal = legacyMarshal

		// We have no way to tell whether the type's Marshal method
		// supports deterministic serialization or not, but this
		// preserves the v1 implementation's behavior of always
		// calling Marshal methods when present.
		mi.methods.Flags |= piface.SupportMarshalDeterministic
	***REMOVED***
	if _, ok := v.(legacyUnmarshaler); ok ***REMOVED***
		mi.methods.Unmarshal = legacyUnmarshal
	***REMOVED***
	if _, ok := v.(legacyMerger); ok ***REMOVED***
		mi.methods.Merge = legacyMerge
	***REMOVED***

	if mi, ok := legacyMessageTypeCache.LoadOrStore(t, mi); ok ***REMOVED***
		return mi.(*MessageInfo)
	***REMOVED***
	return mi
***REMOVED***

var legacyMessageDescCache sync.Map // map[reflect.Type]protoreflect.MessageDescriptor

// LegacyLoadMessageDesc returns an MessageDescriptor derived from the Go type,
// which must be a *struct kind and not implement the v2 API already.
//
// This is exported for testing purposes.
func LegacyLoadMessageDesc(t reflect.Type) pref.MessageDescriptor ***REMOVED***
	return legacyLoadMessageDesc(t, "")
***REMOVED***
func legacyLoadMessageDesc(t reflect.Type, name pref.FullName) pref.MessageDescriptor ***REMOVED***
	// Fast-path: check if a MessageDescriptor is cached for this concrete type.
	if mi, ok := legacyMessageDescCache.Load(t); ok ***REMOVED***
		return mi.(pref.MessageDescriptor)
	***REMOVED***

	// Slow-path: initialize MessageDescriptor from the raw descriptor.
	mv := reflect.Zero(t).Interface()
	if _, ok := mv.(pref.ProtoMessage); ok ***REMOVED***
		panic(fmt.Sprintf("%v already implements proto.Message", t))
	***REMOVED***
	mdV1, ok := mv.(messageV1)
	if !ok ***REMOVED***
		return aberrantLoadMessageDesc(t, name)
	***REMOVED***

	// If this is a dynamic message type where there isn't a 1-1 mapping between
	// Go and protobuf types, calling the Descriptor method on the zero value of
	// the message type isn't likely to work. If it panics, swallow the panic and
	// continue as if the Descriptor method wasn't present.
	b, idxs := func() ([]byte, []int) ***REMOVED***
		defer func() ***REMOVED***
			recover()
		***REMOVED***()
		return mdV1.Descriptor()
	***REMOVED***()
	if b == nil ***REMOVED***
		return aberrantLoadMessageDesc(t, name)
	***REMOVED***

	// If the Go type has no fields, then this might be a proto3 empty message
	// from before the size cache was added. If there are any fields, check to
	// see that at least one of them looks like something we generated.
	if nfield := t.Elem().NumField(); nfield > 0 ***REMOVED***
		hasProtoField := false
		for i := 0; i < nfield; i++ ***REMOVED***
			f := t.Elem().Field(i)
			if f.Tag.Get("protobuf") != "" || f.Tag.Get("protobuf_oneof") != "" || strings.HasPrefix(f.Name, "XXX_") ***REMOVED***
				hasProtoField = true
				break
			***REMOVED***
		***REMOVED***
		if !hasProtoField ***REMOVED***
			return aberrantLoadMessageDesc(t, name)
		***REMOVED***
	***REMOVED***

	md := legacyLoadFileDesc(b).Messages().Get(idxs[0])
	for _, i := range idxs[1:] ***REMOVED***
		md = md.Messages().Get(i)
	***REMOVED***
	if name != "" && md.FullName() != name ***REMOVED***
		panic(fmt.Sprintf("mismatching message name: got %v, want %v", md.FullName(), name))
	***REMOVED***
	if md, ok := legacyMessageDescCache.LoadOrStore(t, md); ok ***REMOVED***
		return md.(protoreflect.MessageDescriptor)
	***REMOVED***
	return md
***REMOVED***

var (
	aberrantMessageDescLock  sync.Mutex
	aberrantMessageDescCache map[reflect.Type]protoreflect.MessageDescriptor
)

// aberrantLoadMessageDesc returns an MessageDescriptor derived from the Go type,
// which must not implement protoreflect.ProtoMessage or messageV1.
//
// This is a best-effort derivation of the message descriptor using the protobuf
// tags on the struct fields.
func aberrantLoadMessageDesc(t reflect.Type, name pref.FullName) pref.MessageDescriptor ***REMOVED***
	aberrantMessageDescLock.Lock()
	defer aberrantMessageDescLock.Unlock()
	if aberrantMessageDescCache == nil ***REMOVED***
		aberrantMessageDescCache = make(map[reflect.Type]protoreflect.MessageDescriptor)
	***REMOVED***
	return aberrantLoadMessageDescReentrant(t, name)
***REMOVED***
func aberrantLoadMessageDescReentrant(t reflect.Type, name pref.FullName) pref.MessageDescriptor ***REMOVED***
	// Fast-path: check if an MessageDescriptor is cached for this concrete type.
	if md, ok := aberrantMessageDescCache[t]; ok ***REMOVED***
		return md
	***REMOVED***

	// Slow-path: construct a descriptor from the Go struct type (best-effort).
	// Cache the MessageDescriptor early on so that we can resolve internal
	// cyclic references.
	md := &filedesc.Message***REMOVED***L2: new(filedesc.MessageL2)***REMOVED***
	md.L0.FullName = aberrantDeriveMessageName(t, name)
	md.L0.ParentFile = filedesc.SurrogateProto2
	aberrantMessageDescCache[t] = md

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct ***REMOVED***
		return md
	***REMOVED***

	// Try to determine if the message is using proto3 by checking scalars.
	for i := 0; i < t.Elem().NumField(); i++ ***REMOVED***
		f := t.Elem().Field(i)
		if tag := f.Tag.Get("protobuf"); tag != "" ***REMOVED***
			switch f.Type.Kind() ***REMOVED***
			case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
				md.L0.ParentFile = filedesc.SurrogateProto3
			***REMOVED***
			for _, s := range strings.Split(tag, ",") ***REMOVED***
				if s == "proto3" ***REMOVED***
					md.L0.ParentFile = filedesc.SurrogateProto3
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Obtain a list of oneof wrapper types.
	var oneofWrappers []reflect.Type
	for _, method := range []string***REMOVED***"XXX_OneofFuncs", "XXX_OneofWrappers"***REMOVED*** ***REMOVED***
		if fn, ok := t.MethodByName(method); ok ***REMOVED***
			for _, v := range fn.Func.Call([]reflect.Value***REMOVED***reflect.Zero(fn.Type.In(0))***REMOVED***) ***REMOVED***
				if vs, ok := v.Interface().([]interface***REMOVED******REMOVED***); ok ***REMOVED***
					for _, v := range vs ***REMOVED***
						oneofWrappers = append(oneofWrappers, reflect.TypeOf(v))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Obtain a list of the extension ranges.
	if fn, ok := t.MethodByName("ExtensionRangeArray"); ok ***REMOVED***
		vs := fn.Func.Call([]reflect.Value***REMOVED***reflect.Zero(fn.Type.In(0))***REMOVED***)[0]
		for i := 0; i < vs.Len(); i++ ***REMOVED***
			v := vs.Index(i)
			md.L2.ExtensionRanges.List = append(md.L2.ExtensionRanges.List, [2]pref.FieldNumber***REMOVED***
				pref.FieldNumber(v.FieldByName("Start").Int()),
				pref.FieldNumber(v.FieldByName("End").Int() + 1),
			***REMOVED***)
			md.L2.ExtensionRangeOptions = append(md.L2.ExtensionRangeOptions, nil)
		***REMOVED***
	***REMOVED***

	// Derive the message fields by inspecting the struct fields.
	for i := 0; i < t.Elem().NumField(); i++ ***REMOVED***
		f := t.Elem().Field(i)
		if tag := f.Tag.Get("protobuf"); tag != "" ***REMOVED***
			tagKey := f.Tag.Get("protobuf_key")
			tagVal := f.Tag.Get("protobuf_val")
			aberrantAppendField(md, f.Type, tag, tagKey, tagVal)
		***REMOVED***
		if tag := f.Tag.Get("protobuf_oneof"); tag != "" ***REMOVED***
			n := len(md.L2.Oneofs.List)
			md.L2.Oneofs.List = append(md.L2.Oneofs.List, filedesc.Oneof***REMOVED******REMOVED***)
			od := &md.L2.Oneofs.List[n]
			od.L0.FullName = md.FullName().Append(pref.Name(tag))
			od.L0.ParentFile = md.L0.ParentFile
			od.L0.Parent = md
			od.L0.Index = n

			for _, t := range oneofWrappers ***REMOVED***
				if t.Implements(f.Type) ***REMOVED***
					f := t.Elem().Field(0)
					if tag := f.Tag.Get("protobuf"); tag != "" ***REMOVED***
						aberrantAppendField(md, f.Type, tag, "", "")
						fd := &md.L2.Fields.List[len(md.L2.Fields.List)-1]
						fd.L1.ContainingOneof = od
						od.L1.Fields.List = append(od.L1.Fields.List, fd)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return md
***REMOVED***

func aberrantDeriveMessageName(t reflect.Type, name pref.FullName) pref.FullName ***REMOVED***
	if name.IsValid() ***REMOVED***
		return name
	***REMOVED***
	func() ***REMOVED***
		defer func() ***REMOVED*** recover() ***REMOVED***() // swallow possible nil panics
		if m, ok := reflect.Zero(t).Interface().(interface***REMOVED*** XXX_MessageName() string ***REMOVED***); ok ***REMOVED***
			name = pref.FullName(m.XXX_MessageName())
		***REMOVED***
	***REMOVED***()
	if name.IsValid() ***REMOVED***
		return name
	***REMOVED***
	if t.Kind() == reflect.Ptr ***REMOVED***
		t = t.Elem()
	***REMOVED***
	return AberrantDeriveFullName(t)
***REMOVED***

func aberrantAppendField(md *filedesc.Message, goType reflect.Type, tag, tagKey, tagVal string) ***REMOVED***
	t := goType
	isOptional := t.Kind() == reflect.Ptr && t.Elem().Kind() != reflect.Struct
	isRepeated := t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
	if isOptional || isRepeated ***REMOVED***
		t = t.Elem()
	***REMOVED***
	fd := ptag.Unmarshal(tag, t, placeholderEnumValues***REMOVED******REMOVED***).(*filedesc.Field)

	// Append field descriptor to the message.
	n := len(md.L2.Fields.List)
	md.L2.Fields.List = append(md.L2.Fields.List, *fd)
	fd = &md.L2.Fields.List[n]
	fd.L0.FullName = md.FullName().Append(fd.Name())
	fd.L0.ParentFile = md.L0.ParentFile
	fd.L0.Parent = md
	fd.L0.Index = n

	if fd.L1.IsWeak || fd.L1.HasPacked ***REMOVED***
		fd.L1.Options = func() pref.ProtoMessage ***REMOVED***
			opts := descopts.Field.ProtoReflect().New()
			if fd.L1.IsWeak ***REMOVED***
				opts.Set(opts.Descriptor().Fields().ByName("weak"), protoreflect.ValueOfBool(true))
			***REMOVED***
			if fd.L1.HasPacked ***REMOVED***
				opts.Set(opts.Descriptor().Fields().ByName("packed"), protoreflect.ValueOfBool(fd.L1.IsPacked))
			***REMOVED***
			return opts.Interface()
		***REMOVED***
	***REMOVED***

	// Populate Enum and Message.
	if fd.Enum() == nil && fd.Kind() == pref.EnumKind ***REMOVED***
		switch v := reflect.Zero(t).Interface().(type) ***REMOVED***
		case pref.Enum:
			fd.L1.Enum = v.Descriptor()
		default:
			fd.L1.Enum = LegacyLoadEnumDesc(t)
		***REMOVED***
	***REMOVED***
	if fd.Message() == nil && (fd.Kind() == pref.MessageKind || fd.Kind() == pref.GroupKind) ***REMOVED***
		switch v := reflect.Zero(t).Interface().(type) ***REMOVED***
		case pref.ProtoMessage:
			fd.L1.Message = v.ProtoReflect().Descriptor()
		case messageV1:
			fd.L1.Message = LegacyLoadMessageDesc(t)
		default:
			if t.Kind() == reflect.Map ***REMOVED***
				n := len(md.L1.Messages.List)
				md.L1.Messages.List = append(md.L1.Messages.List, filedesc.Message***REMOVED***L2: new(filedesc.MessageL2)***REMOVED***)
				md2 := &md.L1.Messages.List[n]
				md2.L0.FullName = md.FullName().Append(pref.Name(strs.MapEntryName(string(fd.Name()))))
				md2.L0.ParentFile = md.L0.ParentFile
				md2.L0.Parent = md
				md2.L0.Index = n

				md2.L1.IsMapEntry = true
				md2.L2.Options = func() pref.ProtoMessage ***REMOVED***
					opts := descopts.Message.ProtoReflect().New()
					opts.Set(opts.Descriptor().Fields().ByName("map_entry"), protoreflect.ValueOfBool(true))
					return opts.Interface()
				***REMOVED***

				aberrantAppendField(md2, t.Key(), tagKey, "", "")
				aberrantAppendField(md2, t.Elem(), tagVal, "", "")

				fd.L1.Message = md2
				break
			***REMOVED***
			fd.L1.Message = aberrantLoadMessageDescReentrant(t, "")
		***REMOVED***
	***REMOVED***
***REMOVED***

type placeholderEnumValues struct ***REMOVED***
	protoreflect.EnumValueDescriptors
***REMOVED***

func (placeholderEnumValues) ByNumber(n pref.EnumNumber) pref.EnumValueDescriptor ***REMOVED***
	return filedesc.PlaceholderEnumValue(pref.FullName(fmt.Sprintf("UNKNOWN_%d", n)))
***REMOVED***

// legacyMarshaler is the proto.Marshaler interface superseded by protoiface.Methoder.
type legacyMarshaler interface ***REMOVED***
	Marshal() ([]byte, error)
***REMOVED***

// legacyUnmarshaler is the proto.Unmarshaler interface superseded by protoiface.Methoder.
type legacyUnmarshaler interface ***REMOVED***
	Unmarshal([]byte) error
***REMOVED***

// legacyMerger is the proto.Merger interface superseded by protoiface.Methoder.
type legacyMerger interface ***REMOVED***
	Merge(protoiface.MessageV1)
***REMOVED***

var aberrantProtoMethods = &piface.Methods***REMOVED***
	Marshal:   legacyMarshal,
	Unmarshal: legacyUnmarshal,
	Merge:     legacyMerge,

	// We have no way to tell whether the type's Marshal method
	// supports deterministic serialization or not, but this
	// preserves the v1 implementation's behavior of always
	// calling Marshal methods when present.
	Flags: piface.SupportMarshalDeterministic,
***REMOVED***

func legacyMarshal(in piface.MarshalInput) (piface.MarshalOutput, error) ***REMOVED***
	v := in.Message.(unwrapper).protoUnwrap()
	marshaler, ok := v.(legacyMarshaler)
	if !ok ***REMOVED***
		return piface.MarshalOutput***REMOVED******REMOVED***, errors.New("%T does not implement Marshal", v)
	***REMOVED***
	out, err := marshaler.Marshal()
	if in.Buf != nil ***REMOVED***
		out = append(in.Buf, out...)
	***REMOVED***
	return piface.MarshalOutput***REMOVED***
		Buf: out,
	***REMOVED***, err
***REMOVED***

func legacyUnmarshal(in piface.UnmarshalInput) (piface.UnmarshalOutput, error) ***REMOVED***
	v := in.Message.(unwrapper).protoUnwrap()
	unmarshaler, ok := v.(legacyUnmarshaler)
	if !ok ***REMOVED***
		return piface.UnmarshalOutput***REMOVED******REMOVED***, errors.New("%T does not implement Marshal", v)
	***REMOVED***
	return piface.UnmarshalOutput***REMOVED******REMOVED***, unmarshaler.Unmarshal(in.Buf)
***REMOVED***

func legacyMerge(in piface.MergeInput) piface.MergeOutput ***REMOVED***
	dstv := in.Destination.(unwrapper).protoUnwrap()
	merger, ok := dstv.(legacyMerger)
	if !ok ***REMOVED***
		return piface.MergeOutput***REMOVED******REMOVED***
	***REMOVED***
	merger.Merge(Export***REMOVED******REMOVED***.ProtoMessageV1Of(in.Source))
	return piface.MergeOutput***REMOVED***Flags: piface.MergeComplete***REMOVED***
***REMOVED***

// aberrantMessageType implements MessageType for all types other than pointer-to-struct.
type aberrantMessageType struct ***REMOVED***
	t reflect.Type
***REMOVED***

func (mt aberrantMessageType) New() pref.Message ***REMOVED***
	return aberrantMessage***REMOVED***reflect.Zero(mt.t)***REMOVED***
***REMOVED***
func (mt aberrantMessageType) Zero() pref.Message ***REMOVED***
	return aberrantMessage***REMOVED***reflect.Zero(mt.t)***REMOVED***
***REMOVED***
func (mt aberrantMessageType) GoType() reflect.Type ***REMOVED***
	return mt.t
***REMOVED***
func (mt aberrantMessageType) Descriptor() pref.MessageDescriptor ***REMOVED***
	return LegacyLoadMessageDesc(mt.t)
***REMOVED***

// aberrantMessage implements Message for all types other than pointer-to-struct.
//
// When the underlying type implements legacyMarshaler or legacyUnmarshaler,
// the aberrant Message can be marshaled or unmarshaled. Otherwise, there is
// not much that can be done with values of this type.
type aberrantMessage struct ***REMOVED***
	v reflect.Value
***REMOVED***

func (m aberrantMessage) ProtoReflect() pref.Message ***REMOVED***
	return m
***REMOVED***

func (m aberrantMessage) Descriptor() pref.MessageDescriptor ***REMOVED***
	return LegacyLoadMessageDesc(m.v.Type())
***REMOVED***
func (m aberrantMessage) Type() pref.MessageType ***REMOVED***
	return aberrantMessageType***REMOVED***m.v.Type()***REMOVED***
***REMOVED***
func (m aberrantMessage) New() pref.Message ***REMOVED***
	return aberrantMessage***REMOVED***reflect.Zero(m.v.Type())***REMOVED***
***REMOVED***
func (m aberrantMessage) Interface() pref.ProtoMessage ***REMOVED***
	return m
***REMOVED***
func (m aberrantMessage) Range(f func(pref.FieldDescriptor, pref.Value) bool) ***REMOVED***
***REMOVED***
func (m aberrantMessage) Has(pref.FieldDescriptor) bool ***REMOVED***
	panic("invalid field descriptor")
***REMOVED***
func (m aberrantMessage) Clear(pref.FieldDescriptor) ***REMOVED***
	panic("invalid field descriptor")
***REMOVED***
func (m aberrantMessage) Get(pref.FieldDescriptor) pref.Value ***REMOVED***
	panic("invalid field descriptor")
***REMOVED***
func (m aberrantMessage) Set(pref.FieldDescriptor, pref.Value) ***REMOVED***
	panic("invalid field descriptor")
***REMOVED***
func (m aberrantMessage) Mutable(pref.FieldDescriptor) pref.Value ***REMOVED***
	panic("invalid field descriptor")
***REMOVED***
func (m aberrantMessage) NewField(pref.FieldDescriptor) pref.Value ***REMOVED***
	panic("invalid field descriptor")
***REMOVED***
func (m aberrantMessage) WhichOneof(pref.OneofDescriptor) pref.FieldDescriptor ***REMOVED***
	panic("invalid oneof descriptor")
***REMOVED***
func (m aberrantMessage) GetUnknown() pref.RawFields ***REMOVED***
	return nil
***REMOVED***
func (m aberrantMessage) SetUnknown(pref.RawFields) ***REMOVED***
	// SetUnknown discards its input on messages which don't support unknown field storage.
***REMOVED***
func (m aberrantMessage) IsValid() bool ***REMOVED***
	// An invalid message is a read-only, empty message. Since we don't know anything
	// about the alleged contents of this message, we can't say with confidence that
	// it is invalid in this sense. Therefore, report it as valid.
	return true
***REMOVED***
func (m aberrantMessage) ProtoMethods() *piface.Methods ***REMOVED***
	return aberrantProtoMethods
***REMOVED***
func (m aberrantMessage) protoUnwrap() interface***REMOVED******REMOVED*** ***REMOVED***
	return m.v.Interface()
***REMOVED***
