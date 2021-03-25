// comment this out // + build ignore

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// Code generated from gen-helper.go.tmpl - DO NOT EDIT.

package codec

import "encoding"

// GenVersion is the current version of codecgen.
const GenVersion = 16

// This file is used to generate helper code for codecgen.
// The values here i.e. genHelper(En|De)coder are not to be used directly by
// library users. They WILL change continuously and without notice.

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

type genHelperDecDriver struct ***REMOVED***
	decDriver
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
type genHelperEncoder struct ***REMOVED***
	M must
	F fastpathT
	e *Encoder
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
type genHelperDecoder struct ***REMOVED***
	C checkOverflow
	F fastpathT
	d *Decoder
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
	// f.e.encodeI(iv, false, false)
	f.e.encodeValue(rv4i(iv), nil)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncTextMarshal(iv encoding.TextMarshaler) ***REMOVED***
	bs, fnerr := iv.MarshalText()
	f.e.marshalUtf8(bs, fnerr)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncJSONMarshal(iv jsonMarshaler) ***REMOVED***
	bs, fnerr := iv.MarshalJSON()
	f.e.marshalAsis(bs, fnerr)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncBinaryMarshal(iv encoding.BinaryMarshaler) ***REMOVED***
	bs, fnerr := iv.MarshalBinary()
	f.e.marshalRaw(bs, fnerr)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncRaw(iv Raw) ***REMOVED*** f.e.rawBytes(iv) ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) I2Rtid(v interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return i2rtid(v)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) Extension(rtid uintptr) (xfn *extTypeTagFn) ***REMOVED***
	return f.e.h.getExt(rtid, true)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncExtension(v interface***REMOVED******REMOVED***, xfFn *extTypeTagFn) ***REMOVED***
	f.e.e.EncodeExt(v, xfFn.tag, xfFn.ext)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) WriteStr(s string) ***REMOVED***
	f.e.w().writestr(s)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) BytesView(v string) []byte ***REMOVED*** return bytesView(v) ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapStart(length int) ***REMOVED*** f.e.mapStart(length) ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapEnd() ***REMOVED*** f.e.mapEnd() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteArrayStart(length int) ***REMOVED*** f.e.arrayStart(length) ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteArrayEnd() ***REMOVED*** f.e.arrayEnd() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteArrayElem() ***REMOVED*** f.e.arrayElem() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapElemKey() ***REMOVED*** f.e.mapElemKey() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperEncoder) EncWriteMapElemValue() ***REMOVED*** f.e.mapElemValue() ***REMOVED***

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
	rv := rv4i(iv)
	if chkPtr ***REMOVED***
		f.d.ensureDecodeable(rv)
	***REMOVED***
	f.d.decodeValue(rv, nil)
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
	if fnerr := tm.UnmarshalText(f.d.d.DecodeStringAsBytes()); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecJSONUnmarshal(tm jsonUnmarshaler) ***REMOVED***
	// bs := f.dd.DecodeStringAsBytes()
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	if fnerr := tm.UnmarshalJSON(f.d.nextValueBytes()); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecBinaryUnmarshal(bm encoding.BinaryUnmarshaler) ***REMOVED***
	if fnerr := bm.UnmarshalBinary(f.d.d.DecodeBytes(nil, true)); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecRaw() []byte ***REMOVED*** return f.d.rawBytes() ***REMOVED***

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
	return f.d.h.getExt(rtid, true)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecExtension(v interface***REMOVED******REMOVED***, xfFn *extTypeTagFn) ***REMOVED***
	f.d.d.DecodeExt(v, xfFn.tag, xfFn.ext)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecInferLen(clen, maxlen, unit int) (rvlen int) ***REMOVED***
	return decInferLen(clen, maxlen, unit)
***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) StringView(v []byte) string ***REMOVED*** return stringView(v) ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapStart() int ***REMOVED*** return f.d.mapStart() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapEnd() ***REMOVED*** f.d.mapEnd() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadArrayStart() int ***REMOVED*** return f.d.arrayStart() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadArrayEnd() ***REMOVED*** f.d.arrayEnd() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadArrayElem() ***REMOVED*** f.d.arrayElem() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapElemKey() ***REMOVED*** f.d.mapElemKey() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecReadMapElemValue() ***REMOVED*** f.d.mapElemValue() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecDecodeFloat32() float32 ***REMOVED*** return f.d.decodeFloat32() ***REMOVED***

// FOR USE BY CODECGEN ONLY. IT *WILL* CHANGE WITHOUT NOTICE. *DO NOT USE*
func (f genHelperDecoder) DecCheckBreak() bool ***REMOVED*** return f.d.checkBreak() ***REMOVED***
