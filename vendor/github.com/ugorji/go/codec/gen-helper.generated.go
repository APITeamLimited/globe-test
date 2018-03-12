/* // +build ignore */

// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from gen-helper.go.tmpl - DO NOT EDIT.

package codec

import (
	"encoding"
	"reflect"
)

// GenVersion is the current version of codecgen.
const GenVersion = 8

// This file is used to generate helper code for codecgen.
// The values here i.e. genHelper(En|De)coder are not to be used directly by
// library users. They WILL change continuously and without notice.
//
// To help enforce this, we create an unexported type with exported members.
// The only way to get the type is via the one exported type that we control (somewhat).
//
// When static codecs are created for types, they will use this value
// to perform encoding or decoding of primitives or known slice or map types.

// GenHelperEncoder is exported so that it can be used externally by codecgen.
//
// Library users: DO NOT USE IT DIRECTLY. IT WILL CHANGE CONTINOUSLY WITHOUT NOTICE.
func GenHelperEncoder(e *Encoder) (ge genHelperEncoder, ee genHelperEncDriver) ***REMOVED***
	ge = genHelperEncoder***REMOVED***e: e***REMOVED***
	ee = genHelperEncDriver***REMOVED***encDriver: e.e***REMOVED***
	return
***REMOVED***

// GenHelperDecoder is exported so that it can be used externally by codecgen.
//
// Library users: DO NOT USE IT DIRECTLY. IT WILL CHANGE CONTINOUSLY WITHOUT NOTICE.
func GenHelperDecoder(d *Decoder) (gd genHelperDecoder, dd genHelperDecDriver) ***REMOVED***
	gd = genHelperDecoder***REMOVED***d: d***REMOVED***
	dd = genHelperDecDriver***REMOVED***decDriver: d.d***REMOVED***
	return
***REMOVED***

type genHelperEncDriver struct ***REMOVED***
	encDriver
***REMOVED***

func (x genHelperEncDriver) EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
func (x genHelperEncDriver) EncStructFieldKey(keyType valueType, s string) ***REMOVED***
	encStructFieldKey(x.encDriver, keyType, s)
***REMOVED***
func (x genHelperEncDriver) EncodeSymbol(s string) ***REMOVED***
	x.encDriver.EncodeString(cUTF8, s)
***REMOVED***

type genHelperDecDriver struct ***REMOVED***
	decDriver
	C checkOverflow
***REMOVED***

func (x genHelperDecDriver) DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
func (x genHelperDecDriver) DecStructFieldKey(keyType valueType, buf *[decScratchByteArrayLen]byte) []byte ***REMOVED***
	return decStructFieldKey(x.decDriver, keyType, buf)
***REMOVED***
func (x genHelperDecDriver) DecodeInt(bitsize uint8) (i int64) ***REMOVED***
	return x.C.IntV(x.decDriver.DecodeInt64(), bitsize)
***REMOVED***
func (x genHelperDecDriver) DecodeUint(bitsize uint8) (ui uint64) ***REMOVED***
	return x.C.UintV(x.decDriver.DecodeUint64(), bitsize)
***REMOVED***
func (x genHelperDecDriver) DecodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
	f = x.DecodeFloat64()
	if chkOverflow32 && chkOvf.Float32(f) ***REMOVED***
		panicv.errorf("float32 overflow: %v", f)
	***REMOVED***
	return
***REMOVED***
func (x genHelperDecDriver) DecodeFloat32As64() (f float64) ***REMOVED***
	f = x.DecodeFloat64()
	if chkOvf.Float32(f) ***REMOVED***
		panicv.errorf("float32 overflow: %v", f)
	***REMOVED***
	return
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
type genHelperEncoder struct ***REMOVED***
	M must
	e *Encoder
	F fastpathT
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
type genHelperDecoder struct ***REMOVED***
	C checkOverflow
	d *Decoder
	F fastpathT
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBasicHandle() *BasicHandle ***REMOVED***
	return f.e.h
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBinary() bool ***REMOVED***
	return f.e.be // f.e.hh.isBinaryEncoding()
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) IsJSONHandle() bool ***REMOVED***
	return f.e.js
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncFallback(iv interface***REMOVED******REMOVED***) ***REMOVED***
	// println(">>>>>>>>> EncFallback")
	// f.e.encodeI(iv, false, false)
	f.e.encodeValue(reflect.ValueOf(iv), nil, false)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncTextMarshal(iv encoding.TextMarshaler) ***REMOVED***
	bs, fnerr := iv.MarshalText()
	f.e.marshal(bs, fnerr, false, cUTF8)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncJSONMarshal(iv jsonMarshaler) ***REMOVED***
	bs, fnerr := iv.MarshalJSON()
	f.e.marshal(bs, fnerr, true, cUTF8)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBinaryMarshal(iv encoding.BinaryMarshaler) ***REMOVED***
	bs, fnerr := iv.MarshalBinary()
	f.e.marshal(bs, fnerr, false, cRAW)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncRaw(iv Raw) ***REMOVED*** f.e.rawBytes(iv) ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: builtin no longer supported - so we make this method a no-op,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperEncoder) TimeRtidIfBinc() (v uintptr) ***REMOVED*** return ***REMOVED***

// func (f genHelperEncoder) TimeRtidIfBinc() uintptr ***REMOVED***
// 	if _, ok := f.e.hh.(*BincHandle); ok ***REMOVED***
// 		return timeTypId
// 	***REMOVED***
// ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) I2Rtid(v interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return i2rtid(v)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) Extension(rtid uintptr) (xfn *extTypeTagFn) ***REMOVED***
	return f.e.h.getExt(rtid)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncExtension(v interface***REMOVED******REMOVED***, xfFn *extTypeTagFn) ***REMOVED***
	f.e.e.EncodeExt(v, xfFn.tag, xfFn.ext, f.e)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: No longer used,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperEncoder) HasExtensions() bool ***REMOVED***
	return len(f.e.h.extHandle) != 0
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: No longer used,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperEncoder) EncExt(v interface***REMOVED******REMOVED***) (r bool) ***REMOVED***
	if xfFn := f.e.h.getExt(i2rtid(v)); xfFn != nil ***REMOVED***
		f.e.e.EncodeExt(v, xfFn.tag, xfFn.ext, f.e)
		return true
	***REMOVED***
	return false
***REMOVED***

// ---------------- DECODER FOLLOWS -----------------

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBasicHandle() *BasicHandle ***REMOVED***
	return f.d.h
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBinary() bool ***REMOVED***
	return f.d.be // f.d.hh.isBinaryEncoding()
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecSwallow() ***REMOVED*** f.d.swallow() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecScratchBuffer() []byte ***REMOVED***
	return f.d.b[:]
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecScratchArrayBuffer() *[decScratchByteArrayLen]byte ***REMOVED***
	return &f.d.b
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecFallback(iv interface***REMOVED******REMOVED***, chkPtr bool) ***REMOVED***
	// println(">>>>>>>>> DecFallback")
	rv := reflect.ValueOf(iv)
	if chkPtr ***REMOVED***
		rv = f.d.ensureDecodeable(rv)
	***REMOVED***
	f.d.decodeValue(rv, nil, false)
	// f.d.decodeValueFallback(rv)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecSliceHelperStart() (decSliceHelper, int) ***REMOVED***
	return f.d.decSliceHelperStart()
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecStructFieldNotFound(index int, name string) ***REMOVED***
	f.d.structFieldNotFound(index, name)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecArrayCannotExpand(sliceLen, streamLen int) ***REMOVED***
	f.d.arrayCannotExpand(sliceLen, streamLen)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecTextUnmarshal(tm encoding.TextUnmarshaler) ***REMOVED***
	fnerr := tm.UnmarshalText(f.d.d.DecodeStringAsBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecJSONUnmarshal(tm jsonUnmarshaler) ***REMOVED***
	// bs := f.dd.DecodeStringAsBytes()
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	fnerr := tm.UnmarshalJSON(f.d.nextValueBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBinaryUnmarshal(bm encoding.BinaryUnmarshaler) ***REMOVED***
	fnerr := bm.UnmarshalBinary(f.d.d.DecodeBytes(nil, true))
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecRaw() []byte ***REMOVED*** return f.d.rawBytes() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: builtin no longer supported - so we make this method a no-op,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperDecoder) TimeRtidIfBinc() (v uintptr) ***REMOVED*** return ***REMOVED***

// func (f genHelperDecoder) TimeRtidIfBinc() uintptr ***REMOVED***
// 	// Note: builtin is no longer supported - so make this a no-op
// 	if _, ok := f.d.hh.(*BincHandle); ok ***REMOVED***
// 		return timeTypId
// 	***REMOVED***
// 	return 0
// ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) IsJSONHandle() bool ***REMOVED***
	return f.d.js
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) I2Rtid(v interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return i2rtid(v)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) Extension(rtid uintptr) (xfn *extTypeTagFn) ***REMOVED***
	return f.d.h.getExt(rtid)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecExtension(v interface***REMOVED******REMOVED***, xfFn *extTypeTagFn) ***REMOVED***
	f.d.d.DecodeExt(v, xfFn.tag, xfFn.ext)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: No longer used,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperDecoder) HasExtensions() bool ***REMOVED***
	return len(f.d.h.extHandle) != 0
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: No longer used,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperDecoder) DecExt(v interface***REMOVED******REMOVED***) (r bool) ***REMOVED***
	if xfFn := f.d.h.getExt(i2rtid(v)); xfFn != nil ***REMOVED***
		f.d.d.DecodeExt(v, xfFn.tag, xfFn.ext)
		return true
	***REMOVED***
	return false
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecInferLen(clen, maxlen, unit int) (rvlen int) ***REMOVED***
	return decInferLen(clen, maxlen, unit)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
//
// Deprecated: no longer used,
// but leave in-place so that old generated files continue to work without regeneration.
func (f genHelperDecoder) StringView(v []byte) string ***REMOVED*** return stringView(v) ***REMOVED***
