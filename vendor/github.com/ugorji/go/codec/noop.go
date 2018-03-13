// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build ignore

package codec

import (
	"math/rand"
	"time"
)

// NoopHandle returns a no-op handle. It basically does nothing.
// It is only useful for benchmarking, as it gives an idea of the
// overhead from the codec framework.
//
// LIBRARY USERS: *** DO NOT USE ***
func NoopHandle(slen int) *noopHandle ***REMOVED***
	h := noopHandle***REMOVED******REMOVED***
	h.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	h.B = make([][]byte, slen)
	h.S = make([]string, slen)
	for i := 0; i < len(h.S); i++ ***REMOVED***
		b := make([]byte, i+1)
		for j := 0; j < len(b); j++ ***REMOVED***
			b[j] = 'a' + byte(i)
		***REMOVED***
		h.B[i] = b
		h.S[i] = string(b)
	***REMOVED***
	return &h
***REMOVED***

// noopHandle does nothing.
// It is used to simulate the overhead of the codec framework.
type noopHandle struct ***REMOVED***
	BasicHandle
	binaryEncodingType
	noopDrv // noopDrv is unexported here, so we can get a copy of it when needed.
***REMOVED***

type noopDrv struct ***REMOVED***
	d    *Decoder
	e    *Encoder
	i    int
	S    []string
	B    [][]byte
	mks  []bool    // stack. if map (true), else if array (false)
	mk   bool      // top of stack. what container are we on? map or array?
	ct   valueType // last response for IsContainerType.
	cb   int       // counter for ContainerType
	rand *rand.Rand
***REMOVED***

func (h *noopDrv) r(v int) int ***REMOVED*** return h.rand.Intn(v) ***REMOVED***
func (h *noopDrv) m(v int) int ***REMOVED*** h.i++; return h.i % v ***REMOVED***

func (h *noopDrv) newEncDriver(e *Encoder) encDriver ***REMOVED*** h.e = e; return h ***REMOVED***
func (h *noopDrv) newDecDriver(d *Decoder) decDriver ***REMOVED*** h.d = d; return h ***REMOVED***

func (h *noopDrv) reset()       ***REMOVED******REMOVED***
func (h *noopDrv) uncacheRead() ***REMOVED******REMOVED***

// --- encDriver

// stack functions (for map and array)
func (h *noopDrv) start(b bool) ***REMOVED***
	// println("start", len(h.mks)+1)
	h.mks = append(h.mks, b)
	h.mk = b
***REMOVED***
func (h *noopDrv) end() ***REMOVED***
	// println("end: ", len(h.mks)-1)
	h.mks = h.mks[:len(h.mks)-1]
	if len(h.mks) > 0 ***REMOVED***
		h.mk = h.mks[len(h.mks)-1]
	***REMOVED*** else ***REMOVED***
		h.mk = false
	***REMOVED***
***REMOVED***

func (h *noopDrv) EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
func (h *noopDrv) EncodeNil()                              ***REMOVED******REMOVED***
func (h *noopDrv) EncodeInt(i int64)                       ***REMOVED******REMOVED***
func (h *noopDrv) EncodeUint(i uint64)                     ***REMOVED******REMOVED***
func (h *noopDrv) EncodeBool(b bool)                       ***REMOVED******REMOVED***
func (h *noopDrv) EncodeFloat32(f float32)                 ***REMOVED******REMOVED***
func (h *noopDrv) EncodeFloat64(f float64)                 ***REMOVED******REMOVED***
func (h *noopDrv) EncodeRawExt(re *RawExt, e *Encoder)     ***REMOVED******REMOVED***
func (h *noopDrv) EncodeArrayStart(length int)             ***REMOVED*** h.start(true) ***REMOVED***
func (h *noopDrv) EncodeMapStart(length int)               ***REMOVED*** h.start(false) ***REMOVED***
func (h *noopDrv) EncodeEnd()                              ***REMOVED*** h.end() ***REMOVED***

func (h *noopDrv) EncodeString(c charEncoding, v string) ***REMOVED******REMOVED***

// func (h *noopDrv) EncodeSymbol(v string)                      ***REMOVED******REMOVED***
func (h *noopDrv) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED******REMOVED***

func (h *noopDrv) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext, e *Encoder) ***REMOVED******REMOVED***

// ---- decDriver
func (h *noopDrv) initReadNext()                              ***REMOVED******REMOVED***
func (h *noopDrv) CheckBreak() bool                           ***REMOVED*** return false ***REMOVED***
func (h *noopDrv) IsBuiltinType(rt uintptr) bool              ***REMOVED*** return false ***REMOVED***
func (h *noopDrv) DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***)    ***REMOVED******REMOVED***
func (h *noopDrv) DecodeInt(bitsize uint8) (i int64)          ***REMOVED*** return int64(h.m(15)) ***REMOVED***
func (h *noopDrv) DecodeUint(bitsize uint8) (ui uint64)       ***REMOVED*** return uint64(h.m(35)) ***REMOVED***
func (h *noopDrv) DecodeFloat(chkOverflow32 bool) (f float64) ***REMOVED*** return float64(h.m(95)) ***REMOVED***
func (h *noopDrv) DecodeBool() (b bool)                       ***REMOVED*** return h.m(2) == 0 ***REMOVED***
func (h *noopDrv) DecodeString() (s string)                   ***REMOVED*** return h.S[h.m(8)] ***REMOVED***
func (h *noopDrv) DecodeStringAsBytes() []byte                ***REMOVED*** return h.DecodeBytes(nil, true) ***REMOVED***

func (h *noopDrv) DecodeBytes(bs []byte, zerocopy bool) []byte ***REMOVED*** return h.B[h.m(len(h.B))] ***REMOVED***

func (h *noopDrv) ReadEnd() ***REMOVED*** h.end() ***REMOVED***

// toggle map/slice
func (h *noopDrv) ReadMapStart() int   ***REMOVED*** h.start(true); return h.m(10) ***REMOVED***
func (h *noopDrv) ReadArrayStart() int ***REMOVED*** h.start(false); return h.m(10) ***REMOVED***

func (h *noopDrv) ContainerType() (vt valueType) ***REMOVED***
	// return h.m(2) == 0
	// handle kStruct, which will bomb is it calls this and
	// doesn't get back a map or array.
	// consequently, if the return value is not map or array,
	// reset it to one of them based on h.m(7) % 2
	// for kstruct: at least one out of every 2 times,
	// return one of valueTypeMap or Array (else kstruct bombs)
	// however, every 10th time it is called, we just return something else.
	var vals = [...]valueType***REMOVED***valueTypeArray, valueTypeMap***REMOVED***
	//  ------------ TAKE ------------
	// if h.cb%2 == 0 ***REMOVED***
	// 	if h.ct == valueTypeMap || h.ct == valueTypeArray ***REMOVED***
	// 	***REMOVED*** else ***REMOVED***
	// 		h.ct = vals[h.m(2)]
	// 	***REMOVED***
	// ***REMOVED*** else if h.cb%5 == 0 ***REMOVED***
	// 	h.ct = valueType(h.m(8))
	// ***REMOVED*** else ***REMOVED***
	// 	h.ct = vals[h.m(2)]
	// ***REMOVED***
	//  ------------ TAKE ------------
	// if h.cb%16 == 0 ***REMOVED***
	// 	h.ct = valueType(h.cb % 8)
	// ***REMOVED*** else ***REMOVED***
	// 	h.ct = vals[h.cb%2]
	// ***REMOVED***
	h.ct = vals[h.cb%2]
	h.cb++
	return h.ct

	// if h.ct == valueTypeNil || h.ct == valueTypeString || h.ct == valueTypeBytes ***REMOVED***
	// 	return h.ct
	// ***REMOVED***
	// return valueTypeUnset
	// TODO: may need to tweak this so it works.
	// if h.ct == valueTypeMap && vt == valueTypeArray ||
	// 	h.ct == valueTypeArray && vt == valueTypeMap ***REMOVED***
	// 	h.cb = !h.cb
	// 	h.ct = vt
	// 	return h.cb
	// ***REMOVED***
	// // go in a loop and check it.
	// h.ct = vt
	// h.cb = h.m(7) == 0
	// return h.cb
***REMOVED***
func (h *noopDrv) TryDecodeAsNil() bool ***REMOVED***
	if h.mk ***REMOVED***
		return false
	***REMOVED*** else ***REMOVED***
		return h.m(8) == 0
	***REMOVED***
***REMOVED***
func (h *noopDrv) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) uint64 ***REMOVED***
	return 0
***REMOVED***

func (h *noopDrv) DecodeNaked() ***REMOVED***
	// use h.r (random) not h.m() because h.m() could cause the same value to be given.
	var sk int
	if h.mk ***REMOVED***
		// if mapkey, do not support values of nil OR bytes, array, map or rawext
		sk = h.r(7) + 1
	***REMOVED*** else ***REMOVED***
		sk = h.r(12)
	***REMOVED***
	n := &h.d.n
	switch sk ***REMOVED***
	case 0:
		n.v = valueTypeNil
	case 1:
		n.v, n.b = valueTypeBool, false
	case 2:
		n.v, n.b = valueTypeBool, true
	case 3:
		n.v, n.i = valueTypeInt, h.DecodeInt(64)
	case 4:
		n.v, n.u = valueTypeUint, h.DecodeUint(64)
	case 5:
		n.v, n.f = valueTypeFloat, h.DecodeFloat(true)
	case 6:
		n.v, n.f = valueTypeFloat, h.DecodeFloat(false)
	case 7:
		n.v, n.s = valueTypeString, h.DecodeString()
	case 8:
		n.v, n.l = valueTypeBytes, h.B[h.m(len(h.B))]
	case 9:
		n.v = valueTypeArray
	case 10:
		n.v = valueTypeMap
	default:
		n.v = valueTypeExt
		n.u = h.DecodeUint(64)
		n.l = h.B[h.m(len(h.B))]
	***REMOVED***
	h.ct = n.v
	return
***REMOVED***
