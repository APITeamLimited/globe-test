// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build purego appengine

package impl

import (
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
)

func sizeEnum(p pointer, f *coderFieldInfo, _ marshalOptions) (size int) ***REMOVED***
	v := p.v.Elem().Int()
	return f.tagsize + protowire.SizeVarint(uint64(v))
***REMOVED***

func appendEnum(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	v := p.v.Elem().Int()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
***REMOVED***

func consumeEnum(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, _ unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.VarintType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeVarint(b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	p.v.Elem().SetInt(int64(v))
	out.n = n
	return out, nil
***REMOVED***

func mergeEnum(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	dst.v.Elem().Set(src.v.Elem())
***REMOVED***

var coderEnum = pointerCoderFuncs***REMOVED***
	size:      sizeEnum,
	marshal:   appendEnum,
	unmarshal: consumeEnum,
	merge:     mergeEnum,
***REMOVED***

func sizeEnumNoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) ***REMOVED***
	if p.v.Elem().Int() == 0 ***REMOVED***
		return 0
	***REMOVED***
	return sizeEnum(p, f, opts)
***REMOVED***

func appendEnumNoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	if p.v.Elem().Int() == 0 ***REMOVED***
		return b, nil
	***REMOVED***
	return appendEnum(b, p, f, opts)
***REMOVED***

func mergeEnumNoZero(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	if src.v.Elem().Int() != 0 ***REMOVED***
		dst.v.Elem().Set(src.v.Elem())
	***REMOVED***
***REMOVED***

var coderEnumNoZero = pointerCoderFuncs***REMOVED***
	size:      sizeEnumNoZero,
	marshal:   appendEnumNoZero,
	unmarshal: consumeEnum,
	merge:     mergeEnumNoZero,
***REMOVED***

func sizeEnumPtr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) ***REMOVED***
	return sizeEnum(pointer***REMOVED***p.v.Elem()***REMOVED***, f, opts)
***REMOVED***

func appendEnumPtr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	return appendEnum(b, pointer***REMOVED***p.v.Elem()***REMOVED***, f, opts)
***REMOVED***

func consumeEnumPtr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.VarintType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	if p.v.Elem().IsNil() ***REMOVED***
		p.v.Elem().Set(reflect.New(p.v.Elem().Type().Elem()))
	***REMOVED***
	return consumeEnum(b, pointer***REMOVED***p.v.Elem()***REMOVED***, wtyp, f, opts)
***REMOVED***

func mergeEnumPtr(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	if !src.v.Elem().IsNil() ***REMOVED***
		v := reflect.New(dst.v.Type().Elem().Elem())
		v.Elem().Set(src.v.Elem().Elem())
		dst.v.Elem().Set(v)
	***REMOVED***
***REMOVED***

var coderEnumPtr = pointerCoderFuncs***REMOVED***
	size:      sizeEnumPtr,
	marshal:   appendEnumPtr,
	unmarshal: consumeEnumPtr,
	merge:     mergeEnumPtr,
***REMOVED***

func sizeEnumSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) ***REMOVED***
	s := p.v.Elem()
	for i, llen := 0, s.Len(); i < llen; i++ ***REMOVED***
		size += protowire.SizeVarint(uint64(s.Index(i).Int())) + f.tagsize
	***REMOVED***
	return size
***REMOVED***

func appendEnumSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	s := p.v.Elem()
	for i, llen := 0, s.Len(); i < llen; i++ ***REMOVED***
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, uint64(s.Index(i).Int()))
	***REMOVED***
	return b, nil
***REMOVED***

func consumeEnumSlice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	s := p.v.Elem()
	if wtyp == protowire.BytesType ***REMOVED***
		b, n := protowire.ConsumeBytes(b)
		if n < 0 ***REMOVED***
			return out, errDecode
		***REMOVED***
		for len(b) > 0 ***REMOVED***
			v, n := protowire.ConsumeVarint(b)
			if n < 0 ***REMOVED***
				return out, errDecode
			***REMOVED***
			rv := reflect.New(s.Type().Elem()).Elem()
			rv.SetInt(int64(v))
			s.Set(reflect.Append(s, rv))
			b = b[n:]
		***REMOVED***
		out.n = n
		return out, nil
	***REMOVED***
	if wtyp != protowire.VarintType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	v, n := protowire.ConsumeVarint(b)
	if n < 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	rv := reflect.New(s.Type().Elem()).Elem()
	rv.SetInt(int64(v))
	s.Set(reflect.Append(s, rv))
	out.n = n
	return out, nil
***REMOVED***

func mergeEnumSlice(dst, src pointer, _ *coderFieldInfo, _ mergeOptions) ***REMOVED***
	dst.v.Elem().Set(reflect.AppendSlice(dst.v.Elem(), src.v.Elem()))
***REMOVED***

var coderEnumSlice = pointerCoderFuncs***REMOVED***
	size:      sizeEnumSlice,
	marshal:   appendEnumSlice,
	unmarshal: consumeEnumSlice,
	merge:     mergeEnumSlice,
***REMOVED***

func sizeEnumPackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) ***REMOVED***
	s := p.v.Elem()
	llen := s.Len()
	if llen == 0 ***REMOVED***
		return 0
	***REMOVED***
	n := 0
	for i := 0; i < llen; i++ ***REMOVED***
		n += protowire.SizeVarint(uint64(s.Index(i).Int()))
	***REMOVED***
	return f.tagsize + protowire.SizeBytes(n)
***REMOVED***

func appendEnumPackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	s := p.v.Elem()
	llen := s.Len()
	if llen == 0 ***REMOVED***
		return b, nil
	***REMOVED***
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for i := 0; i < llen; i++ ***REMOVED***
		n += protowire.SizeVarint(uint64(s.Index(i).Int()))
	***REMOVED***
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ ***REMOVED***
		b = protowire.AppendVarint(b, uint64(s.Index(i).Int()))
	***REMOVED***
	return b, nil
***REMOVED***

var coderEnumPackedSlice = pointerCoderFuncs***REMOVED***
	size:      sizeEnumPackedSlice,
	marshal:   appendEnumPackedSlice,
	unmarshal: consumeEnumSlice,
	merge:     mergeEnumSlice,
***REMOVED***
