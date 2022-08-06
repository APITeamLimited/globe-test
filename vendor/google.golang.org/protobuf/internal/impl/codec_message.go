// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"
	"sort"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// coderMessageInfo contains per-message information used by the fast-path functions.
// This is a different type from MessageInfo to keep MessageInfo as general-purpose as
// possible.
type coderMessageInfo struct ***REMOVED***
	methods protoiface.Methods

	orderedCoderFields []*coderFieldInfo
	denseCoderFields   []*coderFieldInfo
	coderFields        map[protowire.Number]*coderFieldInfo
	sizecacheOffset    offset
	unknownOffset      offset
	unknownPtrKind     bool
	extensionOffset    offset
	needsInitCheck     bool
	isMessageSet       bool
	numRequiredFields  uint8
***REMOVED***

type coderFieldInfo struct ***REMOVED***
	funcs      pointerCoderFuncs // fast-path per-field functions
	mi         *MessageInfo      // field's message
	ft         reflect.Type
	validation validationInfo           // information used by message validation
	num        protoreflect.FieldNumber // field number
	offset     offset                   // struct field offset
	wiretag    uint64                   // field tag (number + wire type)
	tagsize    int                      // size of the varint-encoded tag
	isPointer  bool                     // true if IsNil may be called on the struct field
	isRequired bool                     // true if field is required
***REMOVED***

func (mi *MessageInfo) makeCoderMethods(t reflect.Type, si structInfo) ***REMOVED***
	mi.sizecacheOffset = invalidOffset
	mi.unknownOffset = invalidOffset
	mi.extensionOffset = invalidOffset

	if si.sizecacheOffset.IsValid() && si.sizecacheType == sizecacheType ***REMOVED***
		mi.sizecacheOffset = si.sizecacheOffset
	***REMOVED***
	if si.unknownOffset.IsValid() && (si.unknownType == unknownFieldsAType || si.unknownType == unknownFieldsBType) ***REMOVED***
		mi.unknownOffset = si.unknownOffset
		mi.unknownPtrKind = si.unknownType.Kind() == reflect.Ptr
	***REMOVED***
	if si.extensionOffset.IsValid() && si.extensionType == extensionFieldsType ***REMOVED***
		mi.extensionOffset = si.extensionOffset
	***REMOVED***

	mi.coderFields = make(map[protowire.Number]*coderFieldInfo)
	fields := mi.Desc.Fields()
	preallocFields := make([]coderFieldInfo, fields.Len())
	for i := 0; i < fields.Len(); i++ ***REMOVED***
		fd := fields.Get(i)

		fs := si.fieldsByNumber[fd.Number()]
		isOneof := fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic()
		if isOneof ***REMOVED***
			fs = si.oneofsByName[fd.ContainingOneof().Name()]
		***REMOVED***
		ft := fs.Type
		var wiretag uint64
		if !fd.IsPacked() ***REMOVED***
			wiretag = protowire.EncodeTag(fd.Number(), wireTypes[fd.Kind()])
		***REMOVED*** else ***REMOVED***
			wiretag = protowire.EncodeTag(fd.Number(), protowire.BytesType)
		***REMOVED***
		var fieldOffset offset
		var funcs pointerCoderFuncs
		var childMessage *MessageInfo
		switch ***REMOVED***
		case ft == nil:
			// This never occurs for generated message types.
			// It implies that a hand-crafted type has missing Go fields
			// for specific protobuf message fields.
			funcs = pointerCoderFuncs***REMOVED***
				size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
					return 0
				***REMOVED***,
				marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
					return nil, nil
				***REMOVED***,
				unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
					panic("missing Go struct field for " + string(fd.FullName()))
				***REMOVED***,
				isInit: func(p pointer, f *coderFieldInfo) error ***REMOVED***
					panic("missing Go struct field for " + string(fd.FullName()))
				***REMOVED***,
				merge: func(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
					panic("missing Go struct field for " + string(fd.FullName()))
				***REMOVED***,
			***REMOVED***
		case isOneof:
			fieldOffset = offsetOf(fs, mi.Exporter)
		case fd.IsWeak():
			fieldOffset = si.weakOffset
			funcs = makeWeakMessageFieldCoder(fd)
		default:
			fieldOffset = offsetOf(fs, mi.Exporter)
			childMessage, funcs = fieldCoder(fd, ft)
		***REMOVED***
		cf := &preallocFields[i]
		*cf = coderFieldInfo***REMOVED***
			num:        fd.Number(),
			offset:     fieldOffset,
			wiretag:    wiretag,
			ft:         ft,
			tagsize:    protowire.SizeVarint(wiretag),
			funcs:      funcs,
			mi:         childMessage,
			validation: newFieldValidationInfo(mi, si, fd, ft),
			isPointer:  fd.Cardinality() == protoreflect.Repeated || fd.HasPresence(),
			isRequired: fd.Cardinality() == protoreflect.Required,
		***REMOVED***
		mi.orderedCoderFields = append(mi.orderedCoderFields, cf)
		mi.coderFields[cf.num] = cf
	***REMOVED***
	for i, oneofs := 0, mi.Desc.Oneofs(); i < oneofs.Len(); i++ ***REMOVED***
		if od := oneofs.Get(i); !od.IsSynthetic() ***REMOVED***
			mi.initOneofFieldCoders(od, si)
		***REMOVED***
	***REMOVED***
	if messageset.IsMessageSet(mi.Desc) ***REMOVED***
		if !mi.extensionOffset.IsValid() ***REMOVED***
			panic(fmt.Sprintf("%v: MessageSet with no extensions field", mi.Desc.FullName()))
		***REMOVED***
		if !mi.unknownOffset.IsValid() ***REMOVED***
			panic(fmt.Sprintf("%v: MessageSet with no unknown field", mi.Desc.FullName()))
		***REMOVED***
		mi.isMessageSet = true
	***REMOVED***
	sort.Slice(mi.orderedCoderFields, func(i, j int) bool ***REMOVED***
		return mi.orderedCoderFields[i].num < mi.orderedCoderFields[j].num
	***REMOVED***)

	var maxDense protoreflect.FieldNumber
	for _, cf := range mi.orderedCoderFields ***REMOVED***
		if cf.num >= 16 && cf.num >= 2*maxDense ***REMOVED***
			break
		***REMOVED***
		maxDense = cf.num
	***REMOVED***
	mi.denseCoderFields = make([]*coderFieldInfo, maxDense+1)
	for _, cf := range mi.orderedCoderFields ***REMOVED***
		if int(cf.num) >= len(mi.denseCoderFields) ***REMOVED***
			break
		***REMOVED***
		mi.denseCoderFields[cf.num] = cf
	***REMOVED***

	// To preserve compatibility with historic wire output, marshal oneofs last.
	if mi.Desc.Oneofs().Len() > 0 ***REMOVED***
		sort.Slice(mi.orderedCoderFields, func(i, j int) bool ***REMOVED***
			fi := fields.ByNumber(mi.orderedCoderFields[i].num)
			fj := fields.ByNumber(mi.orderedCoderFields[j].num)
			return order.LegacyFieldOrder(fi, fj)
		***REMOVED***)
	***REMOVED***

	mi.needsInitCheck = needsInitCheck(mi.Desc)
	if mi.methods.Marshal == nil && mi.methods.Size == nil ***REMOVED***
		mi.methods.Flags |= protoiface.SupportMarshalDeterministic
		mi.methods.Marshal = mi.marshal
		mi.methods.Size = mi.size
	***REMOVED***
	if mi.methods.Unmarshal == nil ***REMOVED***
		mi.methods.Flags |= protoiface.SupportUnmarshalDiscardUnknown
		mi.methods.Unmarshal = mi.unmarshal
	***REMOVED***
	if mi.methods.CheckInitialized == nil ***REMOVED***
		mi.methods.CheckInitialized = mi.checkInitialized
	***REMOVED***
	if mi.methods.Merge == nil ***REMOVED***
		mi.methods.Merge = mi.merge
	***REMOVED***
***REMOVED***

// getUnknownBytes returns a *[]byte for the unknown fields.
// It is the caller's responsibility to check whether the pointer is nil.
// This function is specially designed to be inlineable.
func (mi *MessageInfo) getUnknownBytes(p pointer) *[]byte ***REMOVED***
	if mi.unknownPtrKind ***REMOVED***
		return *p.Apply(mi.unknownOffset).BytesPtr()
	***REMOVED*** else ***REMOVED***
		return p.Apply(mi.unknownOffset).Bytes()
	***REMOVED***
***REMOVED***

// mutableUnknownBytes returns a *[]byte for the unknown fields.
// The returned pointer is guaranteed to not be nil.
func (mi *MessageInfo) mutableUnknownBytes(p pointer) *[]byte ***REMOVED***
	if mi.unknownPtrKind ***REMOVED***
		bp := p.Apply(mi.unknownOffset).BytesPtr()
		if *bp == nil ***REMOVED***
			*bp = new([]byte)
		***REMOVED***
		return *bp
	***REMOVED*** else ***REMOVED***
		return p.Apply(mi.unknownOffset).Bytes()
	***REMOVED***
***REMOVED***
