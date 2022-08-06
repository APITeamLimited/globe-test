// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"
	"sync"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
)

type errInvalidUTF8 struct***REMOVED******REMOVED***

func (errInvalidUTF8) Error() string     ***REMOVED*** return "string field contains invalid UTF-8" ***REMOVED***
func (errInvalidUTF8) InvalidUTF8() bool ***REMOVED*** return true ***REMOVED***
func (errInvalidUTF8) Unwrap() error     ***REMOVED*** return errors.Error ***REMOVED***

// initOneofFieldCoders initializes the fast-path functions for the fields in a oneof.
//
// For size, marshal, and isInit operations, functions are set only on the first field
// in the oneof. The functions are called when the oneof is non-nil, and will dispatch
// to the appropriate field-specific function as necessary.
//
// The unmarshal function is set on each field individually as usual.
func (mi *MessageInfo) initOneofFieldCoders(od protoreflect.OneofDescriptor, si structInfo) ***REMOVED***
	fs := si.oneofsByName[od.Name()]
	ft := fs.Type
	oneofFields := make(map[reflect.Type]*coderFieldInfo)
	needIsInit := false
	fields := od.Fields()
	for i, lim := 0, fields.Len(); i < lim; i++ ***REMOVED***
		fd := od.Fields().Get(i)
		num := fd.Number()
		// Make a copy of the original coderFieldInfo for use in unmarshaling.
		//
		// oneofFields[oneofType].funcs.marshal is the field-specific marshal function.
		//
		// mi.coderFields[num].marshal is set on only the first field in the oneof,
		// and dispatches to the field-specific marshaler in oneofFields.
		cf := *mi.coderFields[num]
		ot := si.oneofWrappersByNumber[num]
		cf.ft = ot.Field(0).Type
		cf.mi, cf.funcs = fieldCoder(fd, cf.ft)
		oneofFields[ot] = &cf
		if cf.funcs.isInit != nil ***REMOVED***
			needIsInit = true
		***REMOVED***
		mi.coderFields[num].funcs.unmarshal = func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
			var vw reflect.Value         // pointer to wrapper type
			vi := p.AsValueOf(ft).Elem() // oneof field value of interface kind
			if !vi.IsNil() && !vi.Elem().IsNil() && vi.Elem().Elem().Type() == ot ***REMOVED***
				vw = vi.Elem()
			***REMOVED*** else ***REMOVED***
				vw = reflect.New(ot)
			***REMOVED***
			out, err := cf.funcs.unmarshal(b, pointerOfValue(vw).Apply(zeroOffset), wtyp, &cf, opts)
			if err != nil ***REMOVED***
				return out, err
			***REMOVED***
			vi.Set(vw)
			return out, nil
		***REMOVED***
	***REMOVED***
	getInfo := func(p pointer) (pointer, *coderFieldInfo) ***REMOVED***
		v := p.AsValueOf(ft).Elem()
		if v.IsNil() ***REMOVED***
			return pointer***REMOVED******REMOVED***, nil
		***REMOVED***
		v = v.Elem() // interface -> *struct
		if v.IsNil() ***REMOVED***
			return pointer***REMOVED******REMOVED***, nil
		***REMOVED***
		return pointerOfValue(v).Apply(zeroOffset), oneofFields[v.Elem().Type()]
	***REMOVED***
	first := mi.coderFields[od.Fields().Get(0).Number()]
	first.funcs.size = func(p pointer, _ *coderFieldInfo, opts marshalOptions) int ***REMOVED***
		p, info := getInfo(p)
		if info == nil || info.funcs.size == nil ***REMOVED***
			return 0
		***REMOVED***
		return info.funcs.size(p, info, opts)
	***REMOVED***
	first.funcs.marshal = func(b []byte, p pointer, _ *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
		p, info := getInfo(p)
		if info == nil || info.funcs.marshal == nil ***REMOVED***
			return b, nil
		***REMOVED***
		return info.funcs.marshal(b, p, info, opts)
	***REMOVED***
	first.funcs.merge = func(dst, src pointer, _ *coderFieldInfo, opts mergeOptions) ***REMOVED***
		srcp, srcinfo := getInfo(src)
		if srcinfo == nil || srcinfo.funcs.merge == nil ***REMOVED***
			return
		***REMOVED***
		dstp, dstinfo := getInfo(dst)
		if dstinfo != srcinfo ***REMOVED***
			dst.AsValueOf(ft).Elem().Set(reflect.New(src.AsValueOf(ft).Elem().Elem().Elem().Type()))
			dstp = pointerOfValue(dst.AsValueOf(ft).Elem().Elem()).Apply(zeroOffset)
		***REMOVED***
		srcinfo.funcs.merge(dstp, srcp, srcinfo, opts)
	***REMOVED***
	if needIsInit ***REMOVED***
		first.funcs.isInit = func(p pointer, _ *coderFieldInfo) error ***REMOVED***
			p, info := getInfo(p)
			if info == nil || info.funcs.isInit == nil ***REMOVED***
				return nil
			***REMOVED***
			return info.funcs.isInit(p, info)
		***REMOVED***
	***REMOVED***
***REMOVED***

func makeWeakMessageFieldCoder(fd protoreflect.FieldDescriptor) pointerCoderFuncs ***REMOVED***
	var once sync.Once
	var messageType protoreflect.MessageType
	lazyInit := func() ***REMOVED***
		once.Do(func() ***REMOVED***
			messageName := fd.Message().FullName()
			messageType, _ = protoregistry.GlobalTypes.FindMessageByName(messageName)
		***REMOVED***)
	***REMOVED***

	return pointerCoderFuncs***REMOVED***
		size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
			m, ok := p.WeakFields().get(f.num)
			if !ok ***REMOVED***
				return 0
			***REMOVED***
			lazyInit()
			if messageType == nil ***REMOVED***
				panic(fmt.Sprintf("weak message %v is not linked in", fd.Message().FullName()))
			***REMOVED***
			return sizeMessage(m, f.tagsize, opts)
		***REMOVED***,
		marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
			m, ok := p.WeakFields().get(f.num)
			if !ok ***REMOVED***
				return b, nil
			***REMOVED***
			lazyInit()
			if messageType == nil ***REMOVED***
				panic(fmt.Sprintf("weak message %v is not linked in", fd.Message().FullName()))
			***REMOVED***
			return appendMessage(b, m, f.wiretag, opts)
		***REMOVED***,
		unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
			fs := p.WeakFields()
			m, ok := fs.get(f.num)
			if !ok ***REMOVED***
				lazyInit()
				if messageType == nil ***REMOVED***
					return unmarshalOutput***REMOVED******REMOVED***, errUnknown
				***REMOVED***
				m = messageType.New().Interface()
				fs.set(f.num, m)
			***REMOVED***
			return consumeMessage(b, m, wtyp, opts)
		***REMOVED***,
		isInit: func(p pointer, f *coderFieldInfo) error ***REMOVED***
			m, ok := p.WeakFields().get(f.num)
			if !ok ***REMOVED***
				return nil
			***REMOVED***
			return proto.CheckInitialized(m)
		***REMOVED***,
		merge: func(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
			sm, ok := src.WeakFields().get(f.num)
			if !ok ***REMOVED***
				return
			***REMOVED***
			dm, ok := dst.WeakFields().get(f.num)
			if !ok ***REMOVED***
				lazyInit()
				if messageType == nil ***REMOVED***
					panic(fmt.Sprintf("weak message %v is not linked in", fd.Message().FullName()))
				***REMOVED***
				dm = messageType.New().Interface()
				dst.WeakFields().set(f.num, dm)
			***REMOVED***
			opts.Merge(dm, sm)
		***REMOVED***,
	***REMOVED***
***REMOVED***

func makeMessageFieldCoder(fd protoreflect.FieldDescriptor, ft reflect.Type) pointerCoderFuncs ***REMOVED***
	if mi := getMessageInfo(ft); mi != nil ***REMOVED***
		funcs := pointerCoderFuncs***REMOVED***
			size:      sizeMessageInfo,
			marshal:   appendMessageInfo,
			unmarshal: consumeMessageInfo,
			merge:     mergeMessage,
		***REMOVED***
		if needsInitCheck(mi.Desc) ***REMOVED***
			funcs.isInit = isInitMessageInfo
		***REMOVED***
		return funcs
	***REMOVED*** else ***REMOVED***
		return pointerCoderFuncs***REMOVED***
			size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
				m := asMessage(p.AsValueOf(ft).Elem())
				return sizeMessage(m, f.tagsize, opts)
			***REMOVED***,
			marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
				m := asMessage(p.AsValueOf(ft).Elem())
				return appendMessage(b, m, f.wiretag, opts)
			***REMOVED***,
			unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
				mp := p.AsValueOf(ft).Elem()
				if mp.IsNil() ***REMOVED***
					mp.Set(reflect.New(ft.Elem()))
				***REMOVED***
				return consumeMessage(b, asMessage(mp), wtyp, opts)
			***REMOVED***,
			isInit: func(p pointer, f *coderFieldInfo) error ***REMOVED***
				m := asMessage(p.AsValueOf(ft).Elem())
				return proto.CheckInitialized(m)
			***REMOVED***,
			merge: mergeMessage,
		***REMOVED***
	***REMOVED***
***REMOVED***

func sizeMessageInfo(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
	return protowire.SizeBytes(f.mi.sizePointer(p.Elem(), opts)) + f.tagsize
***REMOVED***

func appendMessageInfo(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(f.mi.sizePointer(p.Elem(), opts)))
	return f.mi.marshalAppendPointer(b, p.Elem(), opts)
***REMOVED***

func consumeMessageInfo(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	if p.Elem().IsNil() ***REMOVED***
		p.SetPointer(pointerOfValue(reflect.New(f.mi.GoReflectType.Elem())))
	***REMOVED***
	o, err := f.mi.unmarshalPointer(v, p.Elem(), 0, opts)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	out.n = n
	out.initialized = o.initialized
	return out, nil
***REMOVED***

func isInitMessageInfo(p pointer, f *coderFieldInfo) error ***REMOVED***
	return f.mi.checkInitializedPointer(p.Elem())
***REMOVED***

func sizeMessage(m proto.Message, tagsize int, _ marshalOptions) int ***REMOVED***
	return protowire.SizeBytes(proto.Size(m)) + tagsize
***REMOVED***

func appendMessage(b []byte, m proto.Message, wiretag uint64, opts marshalOptions) ([]byte, error) ***REMOVED***
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, uint64(proto.Size(m)))
	return opts.Options().MarshalAppend(b, m)
***REMOVED***

func consumeMessage(b []byte, m proto.Message, wtyp protowire.Type, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	o, err := opts.Options().UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     v,
		Message: m.ProtoReflect(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	out.n = n
	out.initialized = o.Flags&protoiface.UnmarshalInitialized != 0
	return out, nil
***REMOVED***

func sizeMessageValue(v protoreflect.Value, tagsize int, opts marshalOptions) int ***REMOVED***
	m := v.Message().Interface()
	return sizeMessage(m, tagsize, opts)
***REMOVED***

func appendMessageValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) ***REMOVED***
	m := v.Message().Interface()
	return appendMessage(b, m, wiretag, opts)
***REMOVED***

func consumeMessageValue(b []byte, v protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (protoreflect.Value, unmarshalOutput, error) ***REMOVED***
	m := v.Message().Interface()
	out, err := consumeMessage(b, m, wtyp, opts)
	return v, out, err
***REMOVED***

func isInitMessageValue(v protoreflect.Value) error ***REMOVED***
	m := v.Message().Interface()
	return proto.CheckInitialized(m)
***REMOVED***

var coderMessageValue = valueCoderFuncs***REMOVED***
	size:      sizeMessageValue,
	marshal:   appendMessageValue,
	unmarshal: consumeMessageValue,
	isInit:    isInitMessageValue,
	merge:     mergeMessageValue,
***REMOVED***

func sizeGroupValue(v protoreflect.Value, tagsize int, opts marshalOptions) int ***REMOVED***
	m := v.Message().Interface()
	return sizeGroup(m, tagsize, opts)
***REMOVED***

func appendGroupValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) ***REMOVED***
	m := v.Message().Interface()
	return appendGroup(b, m, wiretag, opts)
***REMOVED***

func consumeGroupValue(b []byte, v protoreflect.Value, num protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (protoreflect.Value, unmarshalOutput, error) ***REMOVED***
	m := v.Message().Interface()
	out, err := consumeGroup(b, m, num, wtyp, opts)
	return v, out, err
***REMOVED***

var coderGroupValue = valueCoderFuncs***REMOVED***
	size:      sizeGroupValue,
	marshal:   appendGroupValue,
	unmarshal: consumeGroupValue,
	isInit:    isInitMessageValue,
	merge:     mergeMessageValue,
***REMOVED***

func makeGroupFieldCoder(fd protoreflect.FieldDescriptor, ft reflect.Type) pointerCoderFuncs ***REMOVED***
	num := fd.Number()
	if mi := getMessageInfo(ft); mi != nil ***REMOVED***
		funcs := pointerCoderFuncs***REMOVED***
			size:      sizeGroupType,
			marshal:   appendGroupType,
			unmarshal: consumeGroupType,
			merge:     mergeMessage,
		***REMOVED***
		if needsInitCheck(mi.Desc) ***REMOVED***
			funcs.isInit = isInitMessageInfo
		***REMOVED***
		return funcs
	***REMOVED*** else ***REMOVED***
		return pointerCoderFuncs***REMOVED***
			size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
				m := asMessage(p.AsValueOf(ft).Elem())
				return sizeGroup(m, f.tagsize, opts)
			***REMOVED***,
			marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
				m := asMessage(p.AsValueOf(ft).Elem())
				return appendGroup(b, m, f.wiretag, opts)
			***REMOVED***,
			unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
				mp := p.AsValueOf(ft).Elem()
				if mp.IsNil() ***REMOVED***
					mp.Set(reflect.New(ft.Elem()))
				***REMOVED***
				return consumeGroup(b, asMessage(mp), num, wtyp, opts)
			***REMOVED***,
			isInit: func(p pointer, f *coderFieldInfo) error ***REMOVED***
				m := asMessage(p.AsValueOf(ft).Elem())
				return proto.CheckInitialized(m)
			***REMOVED***,
			merge: mergeMessage,
		***REMOVED***
	***REMOVED***
***REMOVED***

func sizeGroupType(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
	return 2*f.tagsize + f.mi.sizePointer(p.Elem(), opts)
***REMOVED***

func appendGroupType(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	b = protowire.AppendVarint(b, f.wiretag) // start group
	b, err := f.mi.marshalAppendPointer(b, p.Elem(), opts)
	b = protowire.AppendVarint(b, f.wiretag+1) // end group
	return b, err
***REMOVED***

func consumeGroupType(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.StartGroupType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	if p.Elem().IsNil() ***REMOVED***
		p.SetPointer(pointerOfValue(reflect.New(f.mi.GoReflectType.Elem())))
	***REMOVED***
	return f.mi.unmarshalPointer(b, p.Elem(), f.num, opts)
***REMOVED***

func sizeGroup(m proto.Message, tagsize int, _ marshalOptions) int ***REMOVED***
	return 2*tagsize + proto.Size(m)
***REMOVED***

func appendGroup(b []byte, m proto.Message, wiretag uint64, opts marshalOptions) ([]byte, error) ***REMOVED***
	b = protowire.AppendVarint(b, wiretag) // start group
	b, err := opts.Options().MarshalAppend(b, m)
	b = protowire.AppendVarint(b, wiretag+1) // end group
	return b, err
***REMOVED***

func consumeGroup(b []byte, m proto.Message, num protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.StartGroupType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	b, n := protowire.ConsumeGroup(num, b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	o, err := opts.Options().UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     b,
		Message: m.ProtoReflect(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	out.n = n
	out.initialized = o.Flags&protoiface.UnmarshalInitialized != 0
	return out, nil
***REMOVED***

func makeMessageSliceFieldCoder(fd protoreflect.FieldDescriptor, ft reflect.Type) pointerCoderFuncs ***REMOVED***
	if mi := getMessageInfo(ft); mi != nil ***REMOVED***
		funcs := pointerCoderFuncs***REMOVED***
			size:      sizeMessageSliceInfo,
			marshal:   appendMessageSliceInfo,
			unmarshal: consumeMessageSliceInfo,
			merge:     mergeMessageSlice,
		***REMOVED***
		if needsInitCheck(mi.Desc) ***REMOVED***
			funcs.isInit = isInitMessageSliceInfo
		***REMOVED***
		return funcs
	***REMOVED***
	return pointerCoderFuncs***REMOVED***
		size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
			return sizeMessageSlice(p, ft, f.tagsize, opts)
		***REMOVED***,
		marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
			return appendMessageSlice(b, p, f.wiretag, ft, opts)
		***REMOVED***,
		unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
			return consumeMessageSlice(b, p, ft, wtyp, opts)
		***REMOVED***,
		isInit: func(p pointer, f *coderFieldInfo) error ***REMOVED***
			return isInitMessageSlice(p, ft)
		***REMOVED***,
		merge: mergeMessageSlice,
	***REMOVED***
***REMOVED***

func sizeMessageSliceInfo(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
	s := p.PointerSlice()
	n := 0
	for _, v := range s ***REMOVED***
		n += protowire.SizeBytes(f.mi.sizePointer(v, opts)) + f.tagsize
	***REMOVED***
	return n
***REMOVED***

func appendMessageSliceInfo(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	s := p.PointerSlice()
	var err error
	for _, v := range s ***REMOVED***
		b = protowire.AppendVarint(b, f.wiretag)
		siz := f.mi.sizePointer(v, opts)
		b = protowire.AppendVarint(b, uint64(siz))
		b, err = f.mi.marshalAppendPointer(b, v, opts)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func consumeMessageSliceInfo(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	m := reflect.New(f.mi.GoReflectType.Elem()).Interface()
	mp := pointerOfIface(m)
	o, err := f.mi.unmarshalPointer(v, mp, 0, opts)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	p.AppendPointerSlice(mp)
	out.n = n
	out.initialized = o.initialized
	return out, nil
***REMOVED***

func isInitMessageSliceInfo(p pointer, f *coderFieldInfo) error ***REMOVED***
	s := p.PointerSlice()
	for _, v := range s ***REMOVED***
		if err := f.mi.checkInitializedPointer(v); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func sizeMessageSlice(p pointer, goType reflect.Type, tagsize int, _ marshalOptions) int ***REMOVED***
	s := p.PointerSlice()
	n := 0
	for _, v := range s ***REMOVED***
		m := asMessage(v.AsValueOf(goType.Elem()))
		n += protowire.SizeBytes(proto.Size(m)) + tagsize
	***REMOVED***
	return n
***REMOVED***

func appendMessageSlice(b []byte, p pointer, wiretag uint64, goType reflect.Type, opts marshalOptions) ([]byte, error) ***REMOVED***
	s := p.PointerSlice()
	var err error
	for _, v := range s ***REMOVED***
		m := asMessage(v.AsValueOf(goType.Elem()))
		b = protowire.AppendVarint(b, wiretag)
		siz := proto.Size(m)
		b = protowire.AppendVarint(b, uint64(siz))
		b, err = opts.Options().MarshalAppend(b, m)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func consumeMessageSlice(b []byte, p pointer, goType reflect.Type, wtyp protowire.Type, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	mp := reflect.New(goType.Elem())
	o, err := opts.Options().UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     v,
		Message: asMessage(mp).ProtoReflect(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	p.AppendPointerSlice(pointerOfValue(mp))
	out.n = n
	out.initialized = o.Flags&protoiface.UnmarshalInitialized != 0
	return out, nil
***REMOVED***

func isInitMessageSlice(p pointer, goType reflect.Type) error ***REMOVED***
	s := p.PointerSlice()
	for _, v := range s ***REMOVED***
		m := asMessage(v.AsValueOf(goType.Elem()))
		if err := proto.CheckInitialized(m); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Slices of messages

func sizeMessageSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) int ***REMOVED***
	list := listv.List()
	n := 0
	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		m := list.Get(i).Message().Interface()
		n += protowire.SizeBytes(proto.Size(m)) + tagsize
	***REMOVED***
	return n
***REMOVED***

func appendMessageSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) ***REMOVED***
	list := listv.List()
	mopts := opts.Options()
	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		m := list.Get(i).Message().Interface()
		b = protowire.AppendVarint(b, wiretag)
		siz := proto.Size(m)
		b = protowire.AppendVarint(b, uint64(siz))
		var err error
		b, err = mopts.MarshalAppend(b, m)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func consumeMessageSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) ***REMOVED***
	list := listv.List()
	if wtyp != protowire.BytesType ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***, out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***, out, errDecode
	***REMOVED***
	m := list.NewElement()
	o, err := opts.Options().UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     v,
		Message: m.Message(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***, out, err
	***REMOVED***
	list.Append(m)
	out.n = n
	out.initialized = o.Flags&protoiface.UnmarshalInitialized != 0
	return listv, out, nil
***REMOVED***

func isInitMessageSliceValue(listv protoreflect.Value) error ***REMOVED***
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		m := list.Get(i).Message().Interface()
		if err := proto.CheckInitialized(m); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var coderMessageSliceValue = valueCoderFuncs***REMOVED***
	size:      sizeMessageSliceValue,
	marshal:   appendMessageSliceValue,
	unmarshal: consumeMessageSliceValue,
	isInit:    isInitMessageSliceValue,
	merge:     mergeMessageListValue,
***REMOVED***

func sizeGroupSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) int ***REMOVED***
	list := listv.List()
	n := 0
	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		m := list.Get(i).Message().Interface()
		n += 2*tagsize + proto.Size(m)
	***REMOVED***
	return n
***REMOVED***

func appendGroupSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) ***REMOVED***
	list := listv.List()
	mopts := opts.Options()
	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		m := list.Get(i).Message().Interface()
		b = protowire.AppendVarint(b, wiretag) // start group
		var err error
		b, err = mopts.MarshalAppend(b, m)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
		b = protowire.AppendVarint(b, wiretag+1) // end group
	***REMOVED***
	return b, nil
***REMOVED***

func consumeGroupSliceValue(b []byte, listv protoreflect.Value, num protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) ***REMOVED***
	list := listv.List()
	if wtyp != protowire.StartGroupType ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***, out, errUnknown
	***REMOVED***
	b, n := protowire.ConsumeGroup(num, b)
	if n < 0 ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***, out, errDecode
	***REMOVED***
	m := list.NewElement()
	o, err := opts.Options().UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     b,
		Message: m.Message(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***, out, err
	***REMOVED***
	list.Append(m)
	out.n = n
	out.initialized = o.Flags&protoiface.UnmarshalInitialized != 0
	return listv, out, nil
***REMOVED***

var coderGroupSliceValue = valueCoderFuncs***REMOVED***
	size:      sizeGroupSliceValue,
	marshal:   appendGroupSliceValue,
	unmarshal: consumeGroupSliceValue,
	isInit:    isInitMessageSliceValue,
	merge:     mergeMessageListValue,
***REMOVED***

func makeGroupSliceFieldCoder(fd protoreflect.FieldDescriptor, ft reflect.Type) pointerCoderFuncs ***REMOVED***
	num := fd.Number()
	if mi := getMessageInfo(ft); mi != nil ***REMOVED***
		funcs := pointerCoderFuncs***REMOVED***
			size:      sizeGroupSliceInfo,
			marshal:   appendGroupSliceInfo,
			unmarshal: consumeGroupSliceInfo,
			merge:     mergeMessageSlice,
		***REMOVED***
		if needsInitCheck(mi.Desc) ***REMOVED***
			funcs.isInit = isInitMessageSliceInfo
		***REMOVED***
		return funcs
	***REMOVED***
	return pointerCoderFuncs***REMOVED***
		size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
			return sizeGroupSlice(p, ft, f.tagsize, opts)
		***REMOVED***,
		marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
			return appendGroupSlice(b, p, f.wiretag, ft, opts)
		***REMOVED***,
		unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
			return consumeGroupSlice(b, p, num, wtyp, ft, opts)
		***REMOVED***,
		isInit: func(p pointer, f *coderFieldInfo) error ***REMOVED***
			return isInitMessageSlice(p, ft)
		***REMOVED***,
		merge: mergeMessageSlice,
	***REMOVED***
***REMOVED***

func sizeGroupSlice(p pointer, messageType reflect.Type, tagsize int, _ marshalOptions) int ***REMOVED***
	s := p.PointerSlice()
	n := 0
	for _, v := range s ***REMOVED***
		m := asMessage(v.AsValueOf(messageType.Elem()))
		n += 2*tagsize + proto.Size(m)
	***REMOVED***
	return n
***REMOVED***

func appendGroupSlice(b []byte, p pointer, wiretag uint64, messageType reflect.Type, opts marshalOptions) ([]byte, error) ***REMOVED***
	s := p.PointerSlice()
	var err error
	for _, v := range s ***REMOVED***
		m := asMessage(v.AsValueOf(messageType.Elem()))
		b = protowire.AppendVarint(b, wiretag) // start group
		b, err = opts.Options().MarshalAppend(b, m)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
		b = protowire.AppendVarint(b, wiretag+1) // end group
	***REMOVED***
	return b, nil
***REMOVED***

func consumeGroupSlice(b []byte, p pointer, num protowire.Number, wtyp protowire.Type, goType reflect.Type, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.StartGroupType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	b, n := protowire.ConsumeGroup(num, b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	mp := reflect.New(goType.Elem())
	o, err := opts.Options().UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     b,
		Message: asMessage(mp).ProtoReflect(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	p.AppendPointerSlice(pointerOfValue(mp))
	out.n = n
	out.initialized = o.Flags&protoiface.UnmarshalInitialized != 0
	return out, nil
***REMOVED***

func sizeGroupSliceInfo(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
	s := p.PointerSlice()
	n := 0
	for _, v := range s ***REMOVED***
		n += 2*f.tagsize + f.mi.sizePointer(v, opts)
	***REMOVED***
	return n
***REMOVED***

func appendGroupSliceInfo(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	s := p.PointerSlice()
	var err error
	for _, v := range s ***REMOVED***
		b = protowire.AppendVarint(b, f.wiretag) // start group
		b, err = f.mi.marshalAppendPointer(b, v, opts)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
		b = protowire.AppendVarint(b, f.wiretag+1) // end group
	***REMOVED***
	return b, nil
***REMOVED***

func consumeGroupSliceInfo(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
	if wtyp != protowire.StartGroupType ***REMOVED***
		return unmarshalOutput***REMOVED******REMOVED***, errUnknown
	***REMOVED***
	m := reflect.New(f.mi.GoReflectType.Elem()).Interface()
	mp := pointerOfIface(m)
	out, err := f.mi.unmarshalPointer(b, mp, f.num, opts)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	p.AppendPointerSlice(mp)
	return out, nil
***REMOVED***

func asMessage(v reflect.Value) protoreflect.ProtoMessage ***REMOVED***
	if m, ok := v.Interface().(protoreflect.ProtoMessage); ok ***REMOVED***
		return m
	***REMOVED***
	return legacyWrapMessage(v).Interface()
***REMOVED***
