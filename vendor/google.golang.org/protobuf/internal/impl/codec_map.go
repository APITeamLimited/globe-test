// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"errors"
	"reflect"
	"sort"

	"google.golang.org/protobuf/encoding/protowire"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type mapInfo struct ***REMOVED***
	goType     reflect.Type
	keyWiretag uint64
	valWiretag uint64
	keyFuncs   valueCoderFuncs
	valFuncs   valueCoderFuncs
	keyZero    pref.Value
	keyKind    pref.Kind
	conv       *mapConverter
***REMOVED***

func encoderFuncsForMap(fd pref.FieldDescriptor, ft reflect.Type) (valueMessage *MessageInfo, funcs pointerCoderFuncs) ***REMOVED***
	// TODO: Consider generating specialized map coders.
	keyField := fd.MapKey()
	valField := fd.MapValue()
	keyWiretag := protowire.EncodeTag(1, wireTypes[keyField.Kind()])
	valWiretag := protowire.EncodeTag(2, wireTypes[valField.Kind()])
	keyFuncs := encoderFuncsForValue(keyField)
	valFuncs := encoderFuncsForValue(valField)
	conv := newMapConverter(ft, fd)

	mapi := &mapInfo***REMOVED***
		goType:     ft,
		keyWiretag: keyWiretag,
		valWiretag: valWiretag,
		keyFuncs:   keyFuncs,
		valFuncs:   valFuncs,
		keyZero:    keyField.Default(),
		keyKind:    keyField.Kind(),
		conv:       conv,
	***REMOVED***
	if valField.Kind() == pref.MessageKind ***REMOVED***
		valueMessage = getMessageInfo(ft.Elem())
	***REMOVED***

	funcs = pointerCoderFuncs***REMOVED***
		size: func(p pointer, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
			return sizeMap(p.AsValueOf(ft).Elem(), mapi, f, opts)
		***REMOVED***,
		marshal: func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
			return appendMap(b, p.AsValueOf(ft).Elem(), mapi, f, opts)
		***REMOVED***,
		unmarshal: func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error) ***REMOVED***
			mp := p.AsValueOf(ft)
			if mp.Elem().IsNil() ***REMOVED***
				mp.Elem().Set(reflect.MakeMap(mapi.goType))
			***REMOVED***
			if f.mi == nil ***REMOVED***
				return consumeMap(b, mp.Elem(), wtyp, mapi, f, opts)
			***REMOVED*** else ***REMOVED***
				return consumeMapOfMessage(b, mp.Elem(), wtyp, mapi, f, opts)
			***REMOVED***
		***REMOVED***,
	***REMOVED***
	switch valField.Kind() ***REMOVED***
	case pref.MessageKind:
		funcs.merge = mergeMapOfMessage
	case pref.BytesKind:
		funcs.merge = mergeMapOfBytes
	default:
		funcs.merge = mergeMap
	***REMOVED***
	if valFuncs.isInit != nil ***REMOVED***
		funcs.isInit = func(p pointer, f *coderFieldInfo) error ***REMOVED***
			return isInitMap(p.AsValueOf(ft).Elem(), mapi, f)
		***REMOVED***
	***REMOVED***
	return valueMessage, funcs
***REMOVED***

const (
	mapKeyTagSize = 1 // field 1, tag size 1.
	mapValTagSize = 1 // field 2, tag size 2.
)

func sizeMap(mapv reflect.Value, mapi *mapInfo, f *coderFieldInfo, opts marshalOptions) int ***REMOVED***
	if mapv.Len() == 0 ***REMOVED***
		return 0
	***REMOVED***
	n := 0
	iter := mapRange(mapv)
	for iter.Next() ***REMOVED***
		key := mapi.conv.keyConv.PBValueOf(iter.Key()).MapKey()
		keySize := mapi.keyFuncs.size(key.Value(), mapKeyTagSize, opts)
		var valSize int
		value := mapi.conv.valConv.PBValueOf(iter.Value())
		if f.mi == nil ***REMOVED***
			valSize = mapi.valFuncs.size(value, mapValTagSize, opts)
		***REMOVED*** else ***REMOVED***
			p := pointerOfValue(iter.Value())
			valSize += mapValTagSize
			valSize += protowire.SizeBytes(f.mi.sizePointer(p, opts))
		***REMOVED***
		n += f.tagsize + protowire.SizeBytes(keySize+valSize)
	***REMOVED***
	return n
***REMOVED***

func consumeMap(b []byte, mapv reflect.Value, wtyp protowire.Type, mapi *mapInfo, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	b, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return out, protowire.ParseError(n)
	***REMOVED***
	var (
		key = mapi.keyZero
		val = mapi.conv.valConv.New()
	)
	for len(b) > 0 ***REMOVED***
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 ***REMOVED***
			return out, protowire.ParseError(n)
		***REMOVED***
		if num > protowire.MaxValidNumber ***REMOVED***
			return out, errors.New("invalid field number")
		***REMOVED***
		b = b[n:]
		err := errUnknown
		switch num ***REMOVED***
		case 1:
			var v pref.Value
			var o unmarshalOutput
			v, o, err = mapi.keyFuncs.unmarshal(b, key, num, wtyp, opts)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			key = v
			n = o.n
		case 2:
			var v pref.Value
			var o unmarshalOutput
			v, o, err = mapi.valFuncs.unmarshal(b, val, num, wtyp, opts)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			val = v
			n = o.n
		***REMOVED***
		if err == errUnknown ***REMOVED***
			n = protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 ***REMOVED***
				return out, protowire.ParseError(n)
			***REMOVED***
		***REMOVED*** else if err != nil ***REMOVED***
			return out, err
		***REMOVED***
		b = b[n:]
	***REMOVED***
	mapv.SetMapIndex(mapi.conv.keyConv.GoValueOf(key), mapi.conv.valConv.GoValueOf(val))
	out.n = n
	return out, nil
***REMOVED***

func consumeMapOfMessage(b []byte, mapv reflect.Value, wtyp protowire.Type, mapi *mapInfo, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return out, errUnknown
	***REMOVED***
	b, n := protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return out, protowire.ParseError(n)
	***REMOVED***
	var (
		key = mapi.keyZero
		val = reflect.New(f.mi.GoReflectType.Elem())
	)
	for len(b) > 0 ***REMOVED***
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 ***REMOVED***
			return out, protowire.ParseError(n)
		***REMOVED***
		if num > protowire.MaxValidNumber ***REMOVED***
			return out, errors.New("invalid field number")
		***REMOVED***
		b = b[n:]
		err := errUnknown
		switch num ***REMOVED***
		case 1:
			var v pref.Value
			var o unmarshalOutput
			v, o, err = mapi.keyFuncs.unmarshal(b, key, num, wtyp, opts)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			key = v
			n = o.n
		case 2:
			if wtyp != protowire.BytesType ***REMOVED***
				break
			***REMOVED***
			var v []byte
			v, n = protowire.ConsumeBytes(b)
			if n < 0 ***REMOVED***
				return out, protowire.ParseError(n)
			***REMOVED***
			var o unmarshalOutput
			o, err = f.mi.unmarshalPointer(v, pointerOfValue(val), 0, opts)
			if o.initialized ***REMOVED***
				// Consider this map item initialized so long as we see
				// an initialized value.
				out.initialized = true
			***REMOVED***
		***REMOVED***
		if err == errUnknown ***REMOVED***
			n = protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 ***REMOVED***
				return out, protowire.ParseError(n)
			***REMOVED***
		***REMOVED*** else if err != nil ***REMOVED***
			return out, err
		***REMOVED***
		b = b[n:]
	***REMOVED***
	mapv.SetMapIndex(mapi.conv.keyConv.GoValueOf(key), val)
	out.n = n
	return out, nil
***REMOVED***

func appendMapItem(b []byte, keyrv, valrv reflect.Value, mapi *mapInfo, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	if f.mi == nil ***REMOVED***
		key := mapi.conv.keyConv.PBValueOf(keyrv).MapKey()
		val := mapi.conv.valConv.PBValueOf(valrv)
		size := 0
		size += mapi.keyFuncs.size(key.Value(), mapKeyTagSize, opts)
		size += mapi.valFuncs.size(val, mapValTagSize, opts)
		b = protowire.AppendVarint(b, uint64(size))
		b, err := mapi.keyFuncs.marshal(b, key.Value(), mapi.keyWiretag, opts)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return mapi.valFuncs.marshal(b, val, mapi.valWiretag, opts)
	***REMOVED*** else ***REMOVED***
		key := mapi.conv.keyConv.PBValueOf(keyrv).MapKey()
		val := pointerOfValue(valrv)
		valSize := f.mi.sizePointer(val, opts)
		size := 0
		size += mapi.keyFuncs.size(key.Value(), mapKeyTagSize, opts)
		size += mapValTagSize + protowire.SizeBytes(valSize)
		b = protowire.AppendVarint(b, uint64(size))
		b, err := mapi.keyFuncs.marshal(b, key.Value(), mapi.keyWiretag, opts)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		b = protowire.AppendVarint(b, mapi.valWiretag)
		b = protowire.AppendVarint(b, uint64(valSize))
		return f.mi.marshalAppendPointer(b, val, opts)
	***REMOVED***
***REMOVED***

func appendMap(b []byte, mapv reflect.Value, mapi *mapInfo, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	if mapv.Len() == 0 ***REMOVED***
		return b, nil
	***REMOVED***
	if opts.Deterministic() ***REMOVED***
		return appendMapDeterministic(b, mapv, mapi, f, opts)
	***REMOVED***
	iter := mapRange(mapv)
	for iter.Next() ***REMOVED***
		var err error
		b = protowire.AppendVarint(b, f.wiretag)
		b, err = appendMapItem(b, iter.Key(), iter.Value(), mapi, f, opts)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func appendMapDeterministic(b []byte, mapv reflect.Value, mapi *mapInfo, f *coderFieldInfo, opts marshalOptions) ([]byte, error) ***REMOVED***
	keys := mapv.MapKeys()
	sort.Slice(keys, func(i, j int) bool ***REMOVED***
		switch keys[i].Kind() ***REMOVED***
		case reflect.Bool:
			return !keys[i].Bool() && keys[j].Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return keys[i].Int() < keys[j].Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return keys[i].Uint() < keys[j].Uint()
		case reflect.Float32, reflect.Float64:
			return keys[i].Float() < keys[j].Float()
		case reflect.String:
			return keys[i].String() < keys[j].String()
		default:
			panic("invalid kind: " + keys[i].Kind().String())
		***REMOVED***
	***REMOVED***)
	for _, key := range keys ***REMOVED***
		var err error
		b = protowire.AppendVarint(b, f.wiretag)
		b, err = appendMapItem(b, key, mapv.MapIndex(key), mapi, f, opts)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func isInitMap(mapv reflect.Value, mapi *mapInfo, f *coderFieldInfo) error ***REMOVED***
	if mi := f.mi; mi != nil ***REMOVED***
		mi.init()
		if !mi.needsInitCheck ***REMOVED***
			return nil
		***REMOVED***
		iter := mapRange(mapv)
		for iter.Next() ***REMOVED***
			val := pointerOfValue(iter.Value())
			if err := mi.checkInitializedPointer(val); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		iter := mapRange(mapv)
		for iter.Next() ***REMOVED***
			val := mapi.conv.valConv.PBValueOf(iter.Value())
			if err := mapi.valFuncs.isInit(val); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func mergeMap(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
	dstm := dst.AsValueOf(f.ft).Elem()
	srcm := src.AsValueOf(f.ft).Elem()
	if srcm.Len() == 0 ***REMOVED***
		return
	***REMOVED***
	if dstm.IsNil() ***REMOVED***
		dstm.Set(reflect.MakeMap(f.ft))
	***REMOVED***
	iter := mapRange(srcm)
	for iter.Next() ***REMOVED***
		dstm.SetMapIndex(iter.Key(), iter.Value())
	***REMOVED***
***REMOVED***

func mergeMapOfBytes(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
	dstm := dst.AsValueOf(f.ft).Elem()
	srcm := src.AsValueOf(f.ft).Elem()
	if srcm.Len() == 0 ***REMOVED***
		return
	***REMOVED***
	if dstm.IsNil() ***REMOVED***
		dstm.Set(reflect.MakeMap(f.ft))
	***REMOVED***
	iter := mapRange(srcm)
	for iter.Next() ***REMOVED***
		dstm.SetMapIndex(iter.Key(), reflect.ValueOf(append(emptyBuf[:], iter.Value().Bytes()...)))
	***REMOVED***
***REMOVED***

func mergeMapOfMessage(dst, src pointer, f *coderFieldInfo, opts mergeOptions) ***REMOVED***
	dstm := dst.AsValueOf(f.ft).Elem()
	srcm := src.AsValueOf(f.ft).Elem()
	if srcm.Len() == 0 ***REMOVED***
		return
	***REMOVED***
	if dstm.IsNil() ***REMOVED***
		dstm.Set(reflect.MakeMap(f.ft))
	***REMOVED***
	iter := mapRange(srcm)
	for iter.Next() ***REMOVED***
		val := reflect.New(f.ft.Elem().Elem())
		if f.mi != nil ***REMOVED***
			f.mi.mergePointer(pointerOfValue(val), pointerOfValue(iter.Value()), opts)
		***REMOVED*** else ***REMOVED***
			opts.Merge(asMessage(val), asMessage(iter.Value()))
		***REMOVED***
		dstm.SetMapIndex(iter.Key(), val)
	***REMOVED***
***REMOVED***
