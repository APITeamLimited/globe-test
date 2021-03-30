// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

type mergeOptions struct***REMOVED******REMOVED***

func (o mergeOptions) Merge(dst, src proto.Message) ***REMOVED***
	proto.Merge(dst, src)
***REMOVED***

// merge is protoreflect.Methods.Merge.
func (mi *MessageInfo) merge(in piface.MergeInput) piface.MergeOutput ***REMOVED***
	dp, ok := mi.getPointer(in.Destination)
	if !ok ***REMOVED***
		return piface.MergeOutput***REMOVED******REMOVED***
	***REMOVED***
	sp, ok := mi.getPointer(in.Source)
	if !ok ***REMOVED***
		return piface.MergeOutput***REMOVED******REMOVED***
	***REMOVED***
	mi.mergePointer(dp, sp, mergeOptions***REMOVED******REMOVED***)
	return piface.MergeOutput***REMOVED***Flags: piface.MergeComplete***REMOVED***
***REMOVED***

func (mi *MessageInfo) mergePointer(dst, src pointer, opts mergeOptions) ***REMOVED***
	mi.init()
	if dst.IsNil() ***REMOVED***
		panic(fmt.Sprintf("invalid value: merging into nil message"))
	***REMOVED***
	if src.IsNil() ***REMOVED***
		return
	***REMOVED***
	for _, f := range mi.orderedCoderFields ***REMOVED***
		if f.funcs.merge == nil ***REMOVED***
			continue
		***REMOVED***
		sfptr := src.Apply(f.offset)
		if f.isPointer && sfptr.Elem().IsNil() ***REMOVED***
			continue
		***REMOVED***
		f.funcs.merge(dst.Apply(f.offset), sfptr, f, opts)
	***REMOVED***
	if mi.extensionOffset.IsValid() ***REMOVED***
		sext := src.Apply(mi.extensionOffset).Extensions()
		dext := dst.Apply(mi.extensionOffset).Extensions()
		if *dext == nil ***REMOVED***
			*dext = make(map[int32]ExtensionField)
		***REMOVED***
		for num, sx := range *sext ***REMOVED***
			xt := sx.Type()
			xi := getExtensionFieldInfo(xt)
			if xi.funcs.merge == nil ***REMOVED***
				continue
			***REMOVED***
			dx := (*dext)[num]
			var dv pref.Value
			if dx.Type() == sx.Type() ***REMOVED***
				dv = dx.Value()
			***REMOVED***
			if !dv.IsValid() && xi.unmarshalNeedsValue ***REMOVED***
				dv = xt.New()
			***REMOVED***
			dv = xi.funcs.merge(dv, sx.Value(), opts)
			dx.Set(sx.Type(), dv)
			(*dext)[num] = dx
		***REMOVED***
	***REMOVED***
	if mi.unknownOffset.IsValid() ***REMOVED***
		su := mi.getUnknownBytes(src)
		if su != nil && len(*su) > 0 ***REMOVED***
			du := mi.mutableUnknownBytes(dst)
			*du = append(*du, *su...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func mergeScalarValue(dst, src pref.Value, opts mergeOptions) pref.Value ***REMOVED***
	return src
***REMOVED***

func mergeBytesValue(dst, src pref.Value, opts mergeOptions) pref.Value ***REMOVED***
	return pref.ValueOfBytes(append(emptyBuf[:], src.Bytes()...))
***REMOVED***

func mergeListValue(dst, src pref.Value, opts mergeOptions) pref.Value ***REMOVED***
	dstl := dst.List()
	srcl := src.List()
	for i, llen := 0, srcl.Len(); i < llen; i++ ***REMOVED***
		dstl.Append(srcl.Get(i))
	***REMOVED***
	return dst
***REMOVED***

func mergeBytesListValue(dst, src pref.Value, opts mergeOptions) pref.Value ***REMOVED***
	dstl := dst.List()
	srcl := src.List()
	for i, llen := 0, srcl.Len(); i < llen; i++ ***REMOVED***
		sb := srcl.Get(i).Bytes()
		db := append(emptyBuf[:], sb...)
		dstl.Append(pref.ValueOfBytes(db))
	***REMOVED***
	return dst
***REMOVED***

func mergeMessageListValue(dst, src pref.Value, opts mergeOptions) pref.Value ***REMOVED***
	dstl := dst.List()
	srcl := src.List()
	for i, llen := 0, srcl.Len(); i < llen; i++ ***REMOVED***
		sm := srcl.Get(i).Message()
		dm := proto.Clone(sm.Interface()).ProtoReflect()
		dstl.Append(pref.ValueOfMessage(dm))
	***REMOVED***
	return dst
***REMOVED***

func mergeMessageValue(dst, src pref.Value, opts mergeOptions) pref.Value ***REMOVED***
	opts.Merge(dst.Message().Interface(), src.Message().Interface())
	return dst
***REMOVED***

func mergeMessage(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
	if f.mi != nil ***REMOVED***
		if dst.Elem().IsNil() ***REMOVED***
			dst.SetPointer(pointerOfValue(reflect.New(f.mi.GoReflectType.Elem())))
		***REMOVED***
		f.mi.mergePointer(dst.Elem(), src.Elem(), opts)
	***REMOVED*** else ***REMOVED***
		dm := dst.AsValueOf(f.ft).Elem()
		sm := src.AsValueOf(f.ft).Elem()
		if dm.IsNil() ***REMOVED***
			dm.Set(reflect.New(f.ft.Elem()))
		***REMOVED***
		opts.Merge(asMessage(dm), asMessage(sm))
	***REMOVED***
***REMOVED***

func mergeMessageSlice(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
	for _, sp := range src.PointerSlice() ***REMOVED***
		dm := reflect.New(f.ft.Elem().Elem())
		if f.mi != nil ***REMOVED***
			f.mi.mergePointer(pointerOfValue(dm), sp, opts)
		***REMOVED*** else ***REMOVED***
			opts.Merge(asMessage(dm), asMessage(sp.AsValueOf(f.ft.Elem().Elem())))
		***REMOVED***
		dst.AppendPointerSlice(pointerOfValue(dm))
	***REMOVED***
***REMOVED***

func mergeBytes(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	*dst.Bytes() = append(emptyBuf[:], *src.Bytes()...)
***REMOVED***

func mergeBytesNoZero(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	v := *src.Bytes()
	if len(v) > 0 ***REMOVED***
		*dst.Bytes() = append(emptyBuf[:], v...)
	***REMOVED***
***REMOVED***

func mergeBytesSlice(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	ds := dst.BytesSlice()
	for _, v := range *src.BytesSlice() ***REMOVED***
		*ds = append(*ds, append(emptyBuf[:], v...))
	***REMOVED***
***REMOVED***
