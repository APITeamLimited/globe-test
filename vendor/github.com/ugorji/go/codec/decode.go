// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Some tagging information for error messages.
const (
	msgBadDesc            = "Unrecognized descriptor byte"
	msgDecCannotExpandArr = "cannot expand go array from %v to stream length: %v"
)

var (
	errstrOnlyMapOrArrayCanDecodeIntoStruct = "only encoded map or array can be decoded into a struct"
	errstrCannotDecodeIntoNil               = "cannot decode into nil"

	errmsgExpandSliceOverflow     = "expand slice: slice overflow"
	errmsgExpandSliceCannotChange = "expand slice: cannot change"

	errDecoderNotInitialized = errors.New("Decoder not initialized")

	errDecUnreadByteNothingToRead   = errors.New("cannot unread - nothing has been read")
	errDecUnreadByteLastByteNotRead = errors.New("cannot unread - last byte has not been read")
	errDecUnreadByteUnknown         = errors.New("cannot unread - reason unknown")
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface ***REMOVED***
	unreadn1()

	// readx will use the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// just return a view of the []byte being decoded from.
	// Ensure you call detachZeroCopyBytes later if this needs to be sent outside codec control.
	readx(n int) []byte
	readb([]byte)
	readn1() uint8
	numread() int // number of bytes read
	track()
	stopTrack() []byte

	// skip will skip any byte that matches, and return the first non-matching byte
	skip(accept *bitset256) (token byte)
	// readTo will read any byte that matches, stopping once no-longer matching.
	readTo(in []byte, accept *bitset256) (out []byte)
	// readUntil will read, only stopping once it matches the 'stop' byte.
	readUntil(in []byte, stop byte) (out []byte)
***REMOVED***

type decDriver interface ***REMOVED***
	// this will check if the next token is a break.
	CheckBreak() bool
	// Note: TryDecodeAsNil should be careful not to share any temporary []byte with
	// the rest of the decDriver. This is because sometimes, we optimize by holding onto
	// a transient []byte, and ensuring the only other call we make to the decDriver
	// during that time is maybe a TryDecodeAsNil() call.
	TryDecodeAsNil() bool
	// vt is one of: Bytes, String, Nil, Slice or Map. Return unSet if not known.
	ContainerType() (vt valueType)
	// IsBuiltinType(rt uintptr) bool

	// DecodeNaked will decode primitives (number, bool, string, []byte) and RawExt.
	// For maps and arrays, it will not do the decoding in-band, but will signal
	// the decoder, so that is done later, by setting the decNaked.valueType field.
	//
	// Note: Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	// for extensions, DecodeNaked must read the tag and the []byte if it exists.
	// if the []byte is not read, then kInterfaceNaked will treat it as a Handle
	// that stores the subsequent value in-band, and complete reading the RawExt.
	//
	// extensions should also use readx to decode them, for efficiency.
	// kInterface will extract the detached byte slice if it has to pass it outside its realm.
	DecodeNaked()

	// Deprecated: use DecodeInt64 and DecodeUint64 instead
	// DecodeInt(bitsize uint8) (i int64)
	// DecodeUint(bitsize uint8) (ui uint64)

	DecodeInt64() (i int64)
	DecodeUint64() (ui uint64)

	DecodeFloat64() (f float64)
	DecodeBool() (b bool)
	// DecodeString can also decode symbols.
	// It looks redundant as DecodeBytes is available.
	// However, some codecs (e.g. binc) support symbols and can
	// return a pre-stored string value, meaning that it can bypass
	// the cost of []byte->string conversion.
	DecodeString() (s string)
	DecodeStringAsBytes() (v []byte)

	// DecodeBytes may be called directly, without going through reflection.
	// Consequently, it must be designed to handle possible nil.
	DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte)
	// DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte)

	// decodeExt will decode into a *RawExt or into an extension.
	DecodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64)
	// decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)

	DecodeTime() (t time.Time)

	ReadArrayStart() int
	ReadArrayElem()
	ReadArrayEnd()
	ReadMapStart() int
	ReadMapElemKey()
	ReadMapElemValue()
	ReadMapEnd()

	reset()
	uncacheRead()
***REMOVED***

type decDriverNoopContainerReader struct***REMOVED******REMOVED***

func (x decDriverNoopContainerReader) ReadArrayStart() (v int) ***REMOVED*** return ***REMOVED***
func (x decDriverNoopContainerReader) ReadArrayElem()          ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) ReadArrayEnd()           ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) ReadMapStart() (v int)   ***REMOVED*** return ***REMOVED***
func (x decDriverNoopContainerReader) ReadMapElemKey()         ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) ReadMapElemValue()       ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) ReadMapEnd()             ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) CheckBreak() (v bool)    ***REMOVED*** return ***REMOVED***

// func (x decNoSeparator) uncacheRead() ***REMOVED******REMOVED***

// DecodeOptions captures configuration options during decode.
type DecodeOptions struct ***REMOVED***
	// MapType specifies type to use during schema-less decoding of a map in the stream.
	// If nil (unset), we default to map[string]interface***REMOVED******REMOVED*** iff json handle and MapStringAsKey=true,
	// else map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***.
	MapType reflect.Type

	// SliceType specifies type to use during schema-less decoding of an array in the stream.
	// If nil (unset), we default to []interface***REMOVED******REMOVED*** for all formats.
	SliceType reflect.Type

	// MaxInitLen defines the maxinum initial length that we "make" a collection
	// (string, slice, map, chan). If 0 or negative, we default to a sensible value
	// based on the size of an element in the collection.
	//
	// For example, when decoding, a stream may say that it has 2^64 elements.
	// We should not auto-matically provision a slice of that size, to prevent Out-Of-Memory crash.
	// Instead, we provision up to MaxInitLen, fill that up, and start appending after that.
	MaxInitLen int

	// ReaderBufferSize is the size of the buffer used when reading.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	ReaderBufferSize int

	// If ErrorIfNoField, return an error when decoding a map
	// from a codec stream into a struct, and no matching struct field is found.
	ErrorIfNoField bool

	// If ErrorIfNoArrayExpand, return an error when decoding a slice/array that cannot be expanded.
	// For example, the stream contains an array of 8 items, but you are decoding into a [4]T array,
	// or you are decoding into a slice of length 4 which is non-addressable (and so cannot be set).
	ErrorIfNoArrayExpand bool

	// If SignedInteger, use the int64 during schema-less decoding of unsigned values (not uint64).
	SignedInteger bool

	// MapValueReset controls how we decode into a map value.
	//
	// By default, we MAY retrieve the mapping for a key, and then decode into that.
	// However, especially with big maps, that retrieval may be expensive and unnecessary
	// if the stream already contains all that is necessary to recreate the value.
	//
	// If true, we will never retrieve the previous mapping,
	// but rather decode into a new value and set that in the map.
	//
	// If false, we will retrieve the previous mapping if necessary e.g.
	// the previous mapping is a pointer, or is a struct or array with pre-set state,
	// or is an interface.
	MapValueReset bool

	// SliceElementReset: on decoding a slice, reset the element to a zero value first.
	//
	// concern: if the slice already contained some garbage, we will decode into that garbage.
	SliceElementReset bool

	// InterfaceReset controls how we decode into an interface.
	//
	// By default, when we see a field that is an interface***REMOVED***...***REMOVED***,
	// or a map with interface***REMOVED***...***REMOVED*** value, we will attempt decoding into the
	// "contained" value.
	//
	// However, this prevents us from reading a string into an interface***REMOVED******REMOVED***
	// that formerly contained a number.
	//
	// If true, we will decode into a new "blank" value, and set that in the interface.
	// If false, we will decode into whatever is contained in the interface.
	InterfaceReset bool

	// InternString controls interning of strings during decoding.
	//
	// Some handles, e.g. json, typically will read map keys as strings.
	// If the set of keys are finite, it may help reduce allocation to
	// look them up from a map (than to allocate them afresh).
	//
	// Note: Handles will be smart when using the intern functionality.
	// Every string should not be interned.
	// An excellent use-case for interning is struct field names,
	// or map keys where key type is string.
	InternString bool

	// PreferArrayOverSlice controls whether to decode to an array or a slice.
	//
	// This only impacts decoding into a nil interface***REMOVED******REMOVED***.
	// Consequently, it has no effect on codecgen.
	//
	// *Note*: This only applies if using go1.5 and above,
	// as it requires reflect.ArrayOf support which was absent before go1.5.
	PreferArrayOverSlice bool

	// DeleteOnNilMapValue controls how to decode a nil value in the stream.
	//
	// If true, we will delete the mapping of the key.
	// Else, just set the mapping to the zero value of the type.
	DeleteOnNilMapValue bool
***REMOVED***

// ------------------------------------

type bufioDecReader struct ***REMOVED***
	buf []byte
	r   io.Reader

	c   int // cursor
	n   int // num read
	err error

	tr  []byte
	trb bool
	b   [4]byte
***REMOVED***

func (z *bufioDecReader) reset(r io.Reader) ***REMOVED***
	z.r, z.c, z.n, z.err, z.trb = r, 0, 0, nil, false
	if z.tr != nil ***REMOVED***
		z.tr = z.tr[:0]
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) Read(p []byte) (n int, err error) ***REMOVED***
	if z.err != nil ***REMOVED***
		return 0, z.err
	***REMOVED***
	p0 := p
	n = copy(p, z.buf[z.c:])
	z.c += n
	if z.c == len(z.buf) ***REMOVED***
		z.c = 0
	***REMOVED***
	z.n += n
	if len(p) == n ***REMOVED***
		if z.c == 0 ***REMOVED***
			z.buf = z.buf[:1]
			z.buf[0] = p[len(p)-1]
			z.c = 1
		***REMOVED***
		if z.trb ***REMOVED***
			z.tr = append(z.tr, p0[:n]...)
		***REMOVED***
		return
	***REMOVED***
	p = p[n:]
	var n2 int
	// if we are here, then z.buf is all read
	if len(p) > len(z.buf) ***REMOVED***
		n2, err = decReadFull(z.r, p)
		n += n2
		z.n += n2
		z.err = err
		// don't return EOF if some bytes were read. keep for next time.
		if n > 0 && err == io.EOF ***REMOVED***
			err = nil
		***REMOVED***
		// always keep last byte in z.buf
		z.buf = z.buf[:1]
		z.buf[0] = p[len(p)-1]
		z.c = 1
		if z.trb ***REMOVED***
			z.tr = append(z.tr, p0[:n]...)
		***REMOVED***
		return
	***REMOVED***
	// z.c is now 0, and len(p) <= len(z.buf)
	for len(p) > 0 && z.err == nil ***REMOVED***
		// println("len(p) loop starting ... ")
		z.c = 0
		z.buf = z.buf[0:cap(z.buf)]
		n2, err = z.r.Read(z.buf)
		if n2 > 0 ***REMOVED***
			if err == io.EOF ***REMOVED***
				err = nil
			***REMOVED***
			z.buf = z.buf[:n2]
			n2 = copy(p, z.buf)
			z.c = n2
			n += n2
			z.n += n2
			p = p[n2:]
		***REMOVED***
		z.err = err
		// println("... len(p) loop done")
	***REMOVED***
	if z.c == 0 ***REMOVED***
		z.buf = z.buf[:1]
		z.buf[0] = p[len(p)-1]
		z.c = 1
	***REMOVED***
	if z.trb ***REMOVED***
		z.tr = append(z.tr, p0[:n]...)
	***REMOVED***
	return
***REMOVED***

func (z *bufioDecReader) ReadByte() (b byte, err error) ***REMOVED***
	z.b[0] = 0
	_, err = z.Read(z.b[:1])
	b = z.b[0]
	return
***REMOVED***

func (z *bufioDecReader) UnreadByte() (err error) ***REMOVED***
	if z.err != nil ***REMOVED***
		return z.err
	***REMOVED***
	if z.c > 0 ***REMOVED***
		z.c--
		z.n--
		if z.trb ***REMOVED***
			z.tr = z.tr[:len(z.tr)-1]
		***REMOVED***
		return
	***REMOVED***
	return errDecUnreadByteNothingToRead
***REMOVED***

func (z *bufioDecReader) numread() int ***REMOVED***
	return z.n
***REMOVED***

func (z *bufioDecReader) readx(n int) (bs []byte) ***REMOVED***
	if n <= 0 || z.err != nil ***REMOVED***
		return
	***REMOVED***
	if z.c+n <= len(z.buf) ***REMOVED***
		bs = z.buf[z.c : z.c+n]
		z.n += n
		z.c += n
		if z.trb ***REMOVED***
			z.tr = append(z.tr, bs...)
		***REMOVED***
		return
	***REMOVED***
	bs = make([]byte, n)
	_, err := z.Read(bs)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***

func (z *bufioDecReader) readb(bs []byte) ***REMOVED***
	_, err := z.Read(bs)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// func (z *bufioDecReader) readn1eof() (b uint8, eof bool) ***REMOVED***
// 	b, err := z.ReadByte()
// 	if err != nil ***REMOVED***
// 		if err == io.EOF ***REMOVED***
// 			eof = true
// 		***REMOVED*** else ***REMOVED***
// 			panic(err)
// 		***REMOVED***
// 	***REMOVED***
// 	return
// ***REMOVED***

func (z *bufioDecReader) readn1() (b uint8) ***REMOVED***
	b, err := z.ReadByte()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***

func (z *bufioDecReader) search(in []byte, accept *bitset256, stop, flag uint8) (token byte, out []byte) ***REMOVED***
	// flag: 1 (skip), 2 (readTo), 4 (readUntil)
	if flag == 4 ***REMOVED***
		for i := z.c; i < len(z.buf); i++ ***REMOVED***
			if z.buf[i] == stop ***REMOVED***
				token = z.buf[i]
				z.n = z.n + (i - z.c) - 1
				i++
				out = z.buf[z.c:i]
				if z.trb ***REMOVED***
					z.tr = append(z.tr, z.buf[z.c:i]...)
				***REMOVED***
				z.c = i
				return
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := z.c; i < len(z.buf); i++ ***REMOVED***
			if !accept.isset(z.buf[i]) ***REMOVED***
				token = z.buf[i]
				z.n = z.n + (i - z.c) - 1
				if flag == 1 ***REMOVED***
					i++
				***REMOVED*** else ***REMOVED***
					out = z.buf[z.c:i]
				***REMOVED***
				if z.trb ***REMOVED***
					z.tr = append(z.tr, z.buf[z.c:i]...)
				***REMOVED***
				z.c = i
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
	z.n += len(z.buf) - z.c
	if flag != 1 ***REMOVED***
		out = append(in, z.buf[z.c:]...)
	***REMOVED***
	if z.trb ***REMOVED***
		z.tr = append(z.tr, z.buf[z.c:]...)
	***REMOVED***
	var n2 int
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	for ***REMOVED***
		z.c = 0
		z.buf = z.buf[0:cap(z.buf)]
		n2, z.err = z.r.Read(z.buf)
		if n2 > 0 && z.err != nil ***REMOVED***
			z.err = nil
		***REMOVED***
		z.buf = z.buf[:n2]
		if flag == 4 ***REMOVED***
			for i := 0; i < n2; i++ ***REMOVED***
				if z.buf[i] == stop ***REMOVED***
					token = z.buf[i]
					z.n += i - 1
					i++
					out = append(out, z.buf[z.c:i]...)
					if z.trb ***REMOVED***
						z.tr = append(z.tr, z.buf[z.c:i]...)
					***REMOVED***
					z.c = i
					return
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := 0; i < n2; i++ ***REMOVED***
				if !accept.isset(z.buf[i]) ***REMOVED***
					token = z.buf[i]
					z.n += i - 1
					if flag == 1 ***REMOVED***
						i++
					***REMOVED***
					if flag != 1 ***REMOVED***
						out = append(out, z.buf[z.c:i]...)
					***REMOVED***
					if z.trb ***REMOVED***
						z.tr = append(z.tr, z.buf[z.c:i]...)
					***REMOVED***
					z.c = i
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if flag != 1 ***REMOVED***
			out = append(out, z.buf[:n2]...)
		***REMOVED***
		z.n += n2
		if z.err != nil ***REMOVED***
			return
		***REMOVED***
		if z.trb ***REMOVED***
			z.tr = append(z.tr, z.buf[:n2]...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) skip(accept *bitset256) (token byte) ***REMOVED***
	token, _ = z.search(nil, accept, 0, 1)
	return
***REMOVED***

func (z *bufioDecReader) readTo(in []byte, accept *bitset256) (out []byte) ***REMOVED***
	_, out = z.search(in, accept, 0, 2)
	return
***REMOVED***

func (z *bufioDecReader) readUntil(in []byte, stop byte) (out []byte) ***REMOVED***
	_, out = z.search(in, nil, stop, 4)
	return
***REMOVED***

func (z *bufioDecReader) unreadn1() ***REMOVED***
	err := z.UnreadByte()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) track() ***REMOVED***
	if z.tr != nil ***REMOVED***
		z.tr = z.tr[:0]
	***REMOVED***
	z.trb = true
***REMOVED***

func (z *bufioDecReader) stopTrack() (bs []byte) ***REMOVED***
	z.trb = false
	return z.tr
***REMOVED***

// ioDecReader is a decReader that reads off an io.Reader.
//
// It also has a fallback implementation of ByteScanner if needed.
type ioDecReader struct ***REMOVED***
	r io.Reader // the reader passed in

	rr io.Reader
	br io.ByteScanner

	l   byte // last byte
	ls  byte // last byte status. 0: init-canDoNothing, 1: canRead, 2: canUnread
	trb bool // tracking bytes turned on
	_   bool
	b   [4]byte // tiny buffer for reading single bytes

	x  [scratchByteArrayLen]byte // for: get struct field name, swallow valueTypeBytes, etc
	n  int                       // num read
	tr []byte                    // tracking bytes read
***REMOVED***

func (z *ioDecReader) reset(r io.Reader) ***REMOVED***
	z.r = r
	z.rr = r
	z.l, z.ls, z.n, z.trb = 0, 0, 0, false
	if z.tr != nil ***REMOVED***
		z.tr = z.tr[:0]
	***REMOVED***
	var ok bool
	if z.br, ok = r.(io.ByteScanner); !ok ***REMOVED***
		z.br = z
		z.rr = z
	***REMOVED***
***REMOVED***

func (z *ioDecReader) Read(p []byte) (n int, err error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return
	***REMOVED***
	var firstByte bool
	if z.ls == 1 ***REMOVED***
		z.ls = 2
		p[0] = z.l
		if len(p) == 1 ***REMOVED***
			n = 1
			return
		***REMOVED***
		firstByte = true
		p = p[1:]
	***REMOVED***
	n, err = z.r.Read(p)
	if n > 0 ***REMOVED***
		if err == io.EOF && n == len(p) ***REMOVED***
			err = nil // read was successful, so postpone EOF (till next time)
		***REMOVED***
		z.l = p[n-1]
		z.ls = 2
	***REMOVED***
	if firstByte ***REMOVED***
		n++
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) ReadByte() (c byte, err error) ***REMOVED***
	n, err := z.Read(z.b[:1])
	if n == 1 ***REMOVED***
		c = z.b[0]
		if err == io.EOF ***REMOVED***
			err = nil // read was successful, so postpone EOF (till next time)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) UnreadByte() (err error) ***REMOVED***
	switch z.ls ***REMOVED***
	case 2:
		z.ls = 1
	case 0:
		err = errDecUnreadByteNothingToRead
	case 1:
		err = errDecUnreadByteLastByteNotRead
	default:
		err = errDecUnreadByteUnknown
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) numread() int ***REMOVED***
	return z.n
***REMOVED***

func (z *ioDecReader) readx(n int) (bs []byte) ***REMOVED***
	if n <= 0 ***REMOVED***
		return
	***REMOVED***
	if n < len(z.x) ***REMOVED***
		bs = z.x[:n]
	***REMOVED*** else ***REMOVED***
		bs = make([]byte, n)
	***REMOVED***
	if _, err := decReadFull(z.rr, bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n += len(bs)
	if z.trb ***REMOVED***
		z.tr = append(z.tr, bs...)
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readb(bs []byte) ***REMOVED***
	// if len(bs) == 0 ***REMOVED***
	// 	return
	// ***REMOVED***
	if _, err := decReadFull(z.rr, bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n += len(bs)
	if z.trb ***REMOVED***
		z.tr = append(z.tr, bs...)
	***REMOVED***
***REMOVED***

func (z *ioDecReader) readn1eof() (b uint8, eof bool) ***REMOVED***
	b, err := z.br.ReadByte()
	if err == nil ***REMOVED***
		z.n++
		if z.trb ***REMOVED***
			z.tr = append(z.tr, b)
		***REMOVED***
	***REMOVED*** else if err == io.EOF ***REMOVED***
		eof = true
	***REMOVED*** else ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readn1() (b uint8) ***REMOVED***
	var err error
	if b, err = z.br.ReadByte(); err == nil ***REMOVED***
		z.n++
		if z.trb ***REMOVED***
			z.tr = append(z.tr, b)
		***REMOVED***
		return
	***REMOVED***
	panic(err)
***REMOVED***

func (z *ioDecReader) skip(accept *bitset256) (token byte) ***REMOVED***
	for ***REMOVED***
		var eof bool
		token, eof = z.readn1eof()
		if eof ***REMOVED***
			return
		***REMOVED***
		if accept.isset(token) ***REMOVED***
			continue
		***REMOVED***
		return
	***REMOVED***
***REMOVED***

func (z *ioDecReader) readTo(in []byte, accept *bitset256) (out []byte) ***REMOVED***
	out = in
	for ***REMOVED***
		token, eof := z.readn1eof()
		if eof ***REMOVED***
			return
		***REMOVED***
		if accept.isset(token) ***REMOVED***
			out = append(out, token)
		***REMOVED*** else ***REMOVED***
			z.unreadn1()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *ioDecReader) readUntil(in []byte, stop byte) (out []byte) ***REMOVED***
	out = in
	for ***REMOVED***
		token, eof := z.readn1eof()
		if eof ***REMOVED***
			panic(io.EOF)
		***REMOVED***
		out = append(out, token)
		if token == stop ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *ioDecReader) unreadn1() ***REMOVED***
	err := z.br.UnreadByte()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n--
	if z.trb ***REMOVED***
		if l := len(z.tr) - 1; l >= 0 ***REMOVED***
			z.tr = z.tr[:l]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *ioDecReader) track() ***REMOVED***
	if z.tr != nil ***REMOVED***
		z.tr = z.tr[:0]
	***REMOVED***
	z.trb = true
***REMOVED***

func (z *ioDecReader) stopTrack() (bs []byte) ***REMOVED***
	z.trb = false
	return z.tr
***REMOVED***

// ------------------------------------

var errBytesDecReaderCannotUnread = errors.New("cannot unread last byte read")

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct ***REMOVED***
	b []byte // data
	c int    // cursor
	a int    // available
	t int    // track start
***REMOVED***

func (z *bytesDecReader) reset(in []byte) ***REMOVED***
	z.b = in
	z.a = len(in)
	z.c = 0
	z.t = 0
***REMOVED***

func (z *bytesDecReader) numread() int ***REMOVED***
	return z.c
***REMOVED***

func (z *bytesDecReader) unreadn1() ***REMOVED***
	if z.c == 0 || len(z.b) == 0 ***REMOVED***
		panic(errBytesDecReaderCannotUnread)
	***REMOVED***
	z.c--
	z.a++
	return
***REMOVED***

func (z *bytesDecReader) readx(n int) (bs []byte) ***REMOVED***
	// slicing from a non-constant start position is more expensive,
	// as more computation is required to decipher the pointer start position.
	// However, we do it only once, and it's better than reslicing both z.b and return value.

	if n <= 0 ***REMOVED***
	***REMOVED*** else if z.a == 0 ***REMOVED***
		panic(io.EOF)
	***REMOVED*** else if n > z.a ***REMOVED***
		panic(io.ErrUnexpectedEOF)
	***REMOVED*** else ***REMOVED***
		c0 := z.c
		z.c = c0 + n
		z.a = z.a - n
		bs = z.b[c0:z.c]
	***REMOVED***
	return
***REMOVED***

func (z *bytesDecReader) readb(bs []byte) ***REMOVED***
	copy(bs, z.readx(len(bs)))
***REMOVED***

func (z *bytesDecReader) readn1() (v uint8) ***REMOVED***
	if z.a == 0 ***REMOVED***
		panic(io.EOF)
	***REMOVED***
	v = z.b[z.c]
	z.c++
	z.a--
	return
***REMOVED***

// func (z *bytesDecReader) readn1eof() (v uint8, eof bool) ***REMOVED***
// 	if z.a == 0 ***REMOVED***
// 		eof = true
// 		return
// 	***REMOVED***
// 	v = z.b[z.c]
// 	z.c++
// 	z.a--
// 	return
// ***REMOVED***

func (z *bytesDecReader) skip(accept *bitset256) (token byte) ***REMOVED***
	if z.a == 0 ***REMOVED***
		return
	***REMOVED***
	blen := len(z.b)
	for i := z.c; i < blen; i++ ***REMOVED***
		if !accept.isset(z.b[i]) ***REMOVED***
			token = z.b[i]
			i++
			z.a -= (i - z.c)
			z.c = i
			return
		***REMOVED***
	***REMOVED***
	z.a, z.c = 0, blen
	return
***REMOVED***

func (z *bytesDecReader) readTo(_ []byte, accept *bitset256) (out []byte) ***REMOVED***
	if z.a == 0 ***REMOVED***
		return
	***REMOVED***
	blen := len(z.b)
	for i := z.c; i < blen; i++ ***REMOVED***
		if !accept.isset(z.b[i]) ***REMOVED***
			out = z.b[z.c:i]
			z.a -= (i - z.c)
			z.c = i
			return
		***REMOVED***
	***REMOVED***
	out = z.b[z.c:]
	z.a, z.c = 0, blen
	return
***REMOVED***

func (z *bytesDecReader) readUntil(_ []byte, stop byte) (out []byte) ***REMOVED***
	if z.a == 0 ***REMOVED***
		panic(io.EOF)
	***REMOVED***
	blen := len(z.b)
	for i := z.c; i < blen; i++ ***REMOVED***
		if z.b[i] == stop ***REMOVED***
			i++
			out = z.b[z.c:i]
			z.a -= (i - z.c)
			z.c = i
			return
		***REMOVED***
	***REMOVED***
	z.a, z.c = 0, blen
	panic(io.EOF)
***REMOVED***

func (z *bytesDecReader) track() ***REMOVED***
	z.t = z.c
***REMOVED***

func (z *bytesDecReader) stopTrack() (bs []byte) ***REMOVED***
	return z.b[z.t:z.c]
***REMOVED***

// ----------------------------------------

// func (d *Decoder) builtin(f *codecFnInfo, rv reflect.Value) ***REMOVED***
// 	d.d.DecodeBuiltin(f.ti.rtid, rv2i(rv))
// ***REMOVED***

func (d *Decoder) rawExt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	d.d.DecodeExt(rv2i(rv), 0, nil)
***REMOVED***

func (d *Decoder) ext(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	d.d.DecodeExt(rv2i(rv), f.xfTag, f.xfFn)
***REMOVED***

func (d *Decoder) selferUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv2i(rv).(Selfer).CodecDecodeSelf(d)
***REMOVED***

func (d *Decoder) binaryUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil, true)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (d *Decoder) textUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (d *Decoder) jsonUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	tm := rv2i(rv).(jsonUnmarshaler)
	// bs := d.d.DecodeBytes(d.b[:], true, true)
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	fnerr := tm.UnmarshalJSON(d.nextValueBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (d *Decoder) kErr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	d.errorf("no decoding function defined for kind %v", rv.Kind())
***REMOVED***

// var kIntfCtr uint64

func (d *Decoder) kInterfaceNaked(f *codecFnInfo) (rvn reflect.Value) ***REMOVED***
	// nil interface:
	// use some hieristics to decode it appropriately
	// based on the detected next value in the stream.
	n := d.naked()
	d.d.DecodeNaked()
	if n.v == valueTypeNil ***REMOVED***
		return
	***REMOVED***
	// We cannot decode non-nil stream value into nil interface with methods (e.g. io.Reader).
	if f.ti.numMeth > 0 ***REMOVED***
		d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
		return
	***REMOVED***
	// var useRvn bool
	switch n.v ***REMOVED***
	case valueTypeMap:
		// if json, default to a map type with string keys
		mtid := d.mtid
		if mtid == 0 ***REMOVED***
			if d.jsms ***REMOVED***
				mtid = mapStrIntfTypId
			***REMOVED*** else ***REMOVED***
				mtid = mapIntfIntfTypId
			***REMOVED***
		***REMOVED***
		if mtid == mapIntfIntfTypId ***REMOVED***
			n.initContainers()
			if n.lm < arrayCacheLen ***REMOVED***
				n.ma[n.lm] = nil
				rvn = n.rma[n.lm]
				n.lm++
				d.decode(&n.ma[n.lm-1])
				n.lm--
			***REMOVED*** else ***REMOVED***
				var v2 map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
				d.decode(&v2)
				rvn = reflect.ValueOf(&v2).Elem()
			***REMOVED***
		***REMOVED*** else if mtid == mapStrIntfTypId ***REMOVED*** // for json performance
			n.initContainers()
			if n.ln < arrayCacheLen ***REMOVED***
				n.na[n.ln] = nil
				rvn = n.rna[n.ln]
				n.ln++
				d.decode(&n.na[n.ln-1])
				n.ln--
			***REMOVED*** else ***REMOVED***
				var v2 map[string]interface***REMOVED******REMOVED***
				d.decode(&v2)
				rvn = reflect.ValueOf(&v2).Elem()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if d.mtr ***REMOVED***
				rvn = reflect.New(d.h.MapType)
				d.decode(rv2i(rvn))
				rvn = rvn.Elem()
			***REMOVED*** else ***REMOVED***
				rvn = reflect.New(d.h.MapType).Elem()
				d.decodeValue(rvn, nil, true)
			***REMOVED***
		***REMOVED***
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId ***REMOVED***
			n.initContainers()
			if n.ls < arrayCacheLen ***REMOVED***
				n.sa[n.ls] = nil
				rvn = n.rsa[n.ls]
				n.ls++
				d.decode(&n.sa[n.ls-1])
				n.ls--
			***REMOVED*** else ***REMOVED***
				var v2 []interface***REMOVED******REMOVED***
				d.decode(&v2)
				rvn = reflect.ValueOf(&v2).Elem()
			***REMOVED***
			if reflectArrayOfSupported && d.stid == 0 && d.h.PreferArrayOverSlice ***REMOVED***
				rvn2 := reflect.New(reflectArrayOf(rvn.Len(), intfTyp)).Elem()
				reflect.Copy(rvn2, rvn)
				rvn = rvn2
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if d.str ***REMOVED***
				rvn = reflect.New(d.h.SliceType)
				d.decode(rv2i(rvn))
				rvn = rvn.Elem()
			***REMOVED*** else ***REMOVED***
				rvn = reflect.New(d.h.SliceType).Elem()
				d.decodeValue(rvn, nil, true)
			***REMOVED***
		***REMOVED***
	case valueTypeExt:
		var v interface***REMOVED******REMOVED***
		tag, bytes := n.u, n.l // calling decode below might taint the values
		if bytes == nil ***REMOVED***
			n.initContainers()
			if n.li < arrayCacheLen ***REMOVED***
				n.ia[n.li] = nil
				n.li++
				d.decode(&n.ia[n.li-1])
				// v = *(&n.ia[l])
				n.li--
				v = n.ia[n.li]
				n.ia[n.li] = nil
			***REMOVED*** else ***REMOVED***
				d.decode(&v)
			***REMOVED***
		***REMOVED***
		bfn := d.h.getExtForTag(tag)
		if bfn == nil ***REMOVED***
			var re RawExt
			re.Tag = tag
			re.Data = detachZeroCopyBytes(d.bytes, nil, bytes)
			re.Value = v
			rvn = reflect.ValueOf(&re).Elem()
		***REMOVED*** else ***REMOVED***
			rvnA := reflect.New(bfn.rt)
			if bytes != nil ***REMOVED***
				bfn.ext.ReadExt(rv2i(rvnA), bytes)
			***REMOVED*** else ***REMOVED***
				bfn.ext.UpdateExt(rv2i(rvnA), v)
			***REMOVED***
			rvn = rvnA.Elem()
		***REMOVED***
	case valueTypeNil:
		// no-op
	case valueTypeInt:
		rvn = n.ri
	case valueTypeUint:
		rvn = n.ru
	case valueTypeFloat:
		rvn = n.rf
	case valueTypeBool:
		rvn = n.rb
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs
	case valueTypeBytes:
		rvn = n.rl
	case valueTypeTime:
		rvn = n.rt
	default:
		panicv.errorf("kInterfaceNaked: unexpected valueType: %d", n.v)
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) kInterface(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	// Note:
	// A consequence of how kInterface works, is that
	// if an interface already contains something, we try
	// to decode into what was there before.
	// We do not replace with a generic value (as got from decodeNaked).

	// every interface passed here MUST be settable.
	var rvn reflect.Value
	if rv.IsNil() || d.h.InterfaceReset ***REMOVED***
		// check if mapping to a type: if so, initialize it and move on
		rvn = d.h.intf2impl(f.ti.rtid)
		if rvn.IsValid() ***REMOVED***
			rv.Set(rvn)
		***REMOVED*** else ***REMOVED***
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() ***REMOVED***
				rv.Set(rvn)
			***REMOVED*** else if d.h.InterfaceReset ***REMOVED***
				// reset to zero value based on current type in there.
				rv.Set(reflect.Zero(rv.Elem().Type()))
			***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// now we have a non-nil interface value, meaning it contains a type
		rvn = rv.Elem()
	***REMOVED***
	if d.d.TryDecodeAsNil() ***REMOVED***
		rv.Set(reflect.Zero(rvn.Type()))
		return
	***REMOVED***

	// Note: interface***REMOVED******REMOVED*** is settable, but underlying type may not be.
	// Consequently, we MAY have to create a decodable value out of the underlying value,
	// decode into it, and reset the interface itself.
	// fmt.Printf(">>>> kInterface: rvn type: %v, rv type: %v\n", rvn.Type(), rv.Type())

	rvn2, canDecode := isDecodeable(rvn)
	if canDecode ***REMOVED***
		d.decodeValue(rvn2, nil, true)
		return
	***REMOVED***

	rvn2 = reflect.New(rvn.Type()).Elem()
	rvn2.Set(rvn)
	d.decodeValue(rvn2, nil, true)
	rv.Set(rvn2)
***REMOVED***

func decStructFieldKey(dd decDriver, keyType valueType, b *[decScratchByteArrayLen]byte) (rvkencname []byte) ***REMOVED***
	// use if-else-if, not switch (which compiles to binary-search)
	// since keyType is typically valueTypeString, branch prediction is pretty good.

	if keyType == valueTypeString ***REMOVED***
		rvkencname = dd.DecodeStringAsBytes()
	***REMOVED*** else if keyType == valueTypeInt ***REMOVED***
		rvkencname = strconv.AppendInt(b[:0], dd.DecodeInt64(), 10)
	***REMOVED*** else if keyType == valueTypeUint ***REMOVED***
		rvkencname = strconv.AppendUint(b[:0], dd.DecodeUint64(), 10)
	***REMOVED*** else if keyType == valueTypeFloat ***REMOVED***
		rvkencname = strconv.AppendFloat(b[:0], dd.DecodeFloat64(), 'f', -1, 64)
	***REMOVED*** else ***REMOVED***
		rvkencname = dd.DecodeStringAsBytes()
	***REMOVED***
	return rvkencname
***REMOVED***

func (d *Decoder) kStruct(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	fti := f.ti
	dd := d.d
	elemsep := d.esep
	sfn := structFieldNode***REMOVED***v: rv, update: true***REMOVED***
	ctyp := dd.ContainerType()
	if ctyp == valueTypeMap ***REMOVED***
		containerLen := dd.ReadMapStart()
		if containerLen == 0 ***REMOVED***
			dd.ReadMapEnd()
			return
		***REMOVED***
		tisfi := fti.sfiSort
		hasLen := containerLen >= 0

		var rvkencname []byte
		for j := 0; (hasLen && j < containerLen) || !(hasLen || dd.CheckBreak()); j++ ***REMOVED***
			if elemsep ***REMOVED***
				dd.ReadMapElemKey()
			***REMOVED***
			rvkencname = decStructFieldKey(dd, fti.keyType, &d.b)
			if elemsep ***REMOVED***
				dd.ReadMapElemValue()
			***REMOVED***
			if k := fti.indexForEncName(rvkencname); k > -1 ***REMOVED***
				si := tisfi[k]
				if dd.TryDecodeAsNil() ***REMOVED***
					si.setToZeroValue(rv)
				***REMOVED*** else ***REMOVED***
					d.decodeValue(sfn.field(si), nil, true)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				d.structFieldNotFound(-1, stringView(rvkencname))
			***REMOVED***
			// keepAlive4StringView(rvkencnameB) // not needed, as reference is outside loop
		***REMOVED***
		dd.ReadMapEnd()
	***REMOVED*** else if ctyp == valueTypeArray ***REMOVED***
		containerLen := dd.ReadArrayStart()
		if containerLen == 0 ***REMOVED***
			dd.ReadArrayEnd()
			return
		***REMOVED***
		// Not much gain from doing it two ways for array.
		// Arrays are not used as much for structs.
		hasLen := containerLen >= 0
		for j, si := range fti.sfiSrc ***REMOVED***
			if (hasLen && j == containerLen) || (!hasLen && dd.CheckBreak()) ***REMOVED***
				break
			***REMOVED***
			if elemsep ***REMOVED***
				dd.ReadArrayElem()
			***REMOVED***
			if dd.TryDecodeAsNil() ***REMOVED***
				si.setToZeroValue(rv)
			***REMOVED*** else ***REMOVED***
				d.decodeValue(sfn.field(si), nil, true)
			***REMOVED***
		***REMOVED***
		if containerLen > len(fti.sfiSrc) ***REMOVED***
			// read remaining values and throw away
			for j := len(fti.sfiSrc); j < containerLen; j++ ***REMOVED***
				if elemsep ***REMOVED***
					dd.ReadArrayElem()
				***REMOVED***
				d.structFieldNotFound(j, "")
			***REMOVED***
		***REMOVED***
		dd.ReadArrayEnd()
	***REMOVED*** else ***REMOVED***
		d.errorstr(errstrOnlyMapOrArrayCanDecodeIntoStruct)
		return
	***REMOVED***
***REMOVED***

func (d *Decoder) kSlice(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).
	ti := f.ti
	if f.seq == seqTypeChan && ti.chandir&uint8(reflect.SendDir) == 0 ***REMOVED***
		d.errorf("receive-only channel cannot be used for sending byte(s)")
	***REMOVED***
	dd := d.d
	rtelem0 := ti.elem
	ctyp := dd.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString ***REMOVED***
		// you can only decode bytes or string in the stream into a slice or array of bytes
		if !(ti.rtid == uint8SliceTypId || rtelem0.Kind() == reflect.Uint8) ***REMOVED***
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		***REMOVED***
		if f.seq == seqTypeChan ***REMOVED***
			bs2 := dd.DecodeBytes(nil, true)
			irv := rv2i(rv)
			ch, ok := irv.(chan<- byte)
			if !ok ***REMOVED***
				ch = irv.(chan byte)
			***REMOVED***
			for _, b := range bs2 ***REMOVED***
				ch <- b
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			rvbs := rv.Bytes()
			bs2 := dd.DecodeBytes(rvbs, false)
			// if rvbs == nil && bs2 != nil || rvbs != nil && bs2 == nil || len(bs2) != len(rvbs) ***REMOVED***
			if !(len(bs2) > 0 && len(bs2) == len(rvbs) && &bs2[0] == &rvbs[0]) ***REMOVED***
				if rv.CanSet() ***REMOVED***
					rv.SetBytes(bs2)
				***REMOVED*** else if len(rvbs) > 0 && len(bs2) > 0 ***REMOVED***
					copy(rvbs, bs2)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	// array := f.seq == seqTypeChan

	slh, containerLenS := d.decSliceHelperStart() // only expects valueType(Array|Map)

	// an array can never return a nil slice. so no need to check f.array here.
	if containerLenS == 0 ***REMOVED***
		if rv.CanSet() ***REMOVED***
			if f.seq == seqTypeSlice ***REMOVED***
				if rv.IsNil() ***REMOVED***
					rv.Set(reflect.MakeSlice(ti.rt, 0, 0))
				***REMOVED*** else ***REMOVED***
					rv.SetLen(0)
				***REMOVED***
			***REMOVED*** else if f.seq == seqTypeChan ***REMOVED***
				if rv.IsNil() ***REMOVED***
					rv.Set(reflect.MakeChan(ti.rt, 0))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		slh.End()
		return
	***REMOVED***

	rtelem0Size := int(rtelem0.Size())
	rtElem0Kind := rtelem0.Kind()
	rtelem0Mut := !isImmutableKind(rtElem0Kind)
	rtelem := rtelem0
	rtelemkind := rtelem.Kind()
	for rtelemkind == reflect.Ptr ***REMOVED***
		rtelem = rtelem.Elem()
		rtelemkind = rtelem.Kind()
	***REMOVED***

	var fn *codecFn

	var rvCanset = rv.CanSet()
	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	rvlen := rv.Len()
	rvcap := rv.Cap()
	hasLen := containerLenS > 0
	if hasLen && f.seq == seqTypeSlice ***REMOVED***
		if containerLenS > rvcap ***REMOVED***
			oldRvlenGtZero := rvlen > 0
			rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(rtelem0.Size()))
			if rvlen <= rvcap ***REMOVED***
				if rvCanset ***REMOVED***
					rv.SetLen(rvlen)
				***REMOVED***
			***REMOVED*** else if rvCanset ***REMOVED***
				rv = reflect.MakeSlice(ti.rt, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = true
			***REMOVED*** else ***REMOVED***
				d.errorf("cannot decode into non-settable slice")
			***REMOVED***
			if rvChanged && oldRvlenGtZero && !isImmutableKind(rtelem0.Kind()) ***REMOVED***
				reflect.Copy(rv, rv0) // only copy up to length NOT cap i.e. rv0.Slice(0, rvcap)
			***REMOVED***
		***REMOVED*** else if containerLenS != rvlen ***REMOVED***
			rvlen = containerLenS
			if rvCanset ***REMOVED***
				rv.SetLen(rvlen)
			***REMOVED***
			// else ***REMOVED***
			// rv = rv.Slice(0, rvlen)
			// rvChanged = true
			// d.errorf("cannot decode into non-settable slice")
			// ***REMOVED***
		***REMOVED***
	***REMOVED***

	// consider creating new element once, and just decoding into it.
	var rtelem0Zero reflect.Value
	var rtelem0ZeroValid bool
	var decodeAsNil bool
	var j int
	d.cfer()
	for ; (hasLen && j < containerLenS) || !(hasLen || dd.CheckBreak()); j++ ***REMOVED***
		if j == 0 && (f.seq == seqTypeSlice || f.seq == seqTypeChan) && rv.IsNil() ***REMOVED***
			if hasLen ***REMOVED***
				rvlen = decInferLen(containerLenS, d.h.MaxInitLen, rtelem0Size)
			***REMOVED*** else ***REMOVED***
				rvlen = 8
			***REMOVED***
			if rvCanset ***REMOVED***
				if f.seq == seqTypeSlice ***REMOVED***
					rv = reflect.MakeSlice(ti.rt, rvlen, rvlen)
					rvChanged = true
				***REMOVED*** else ***REMOVED*** // chan
					rv = reflect.MakeChan(ti.rt, rvlen)
					rvChanged = true
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				d.errorf("cannot decode into non-settable slice")
			***REMOVED***
		***REMOVED***
		slh.ElemContainerState(j)
		decodeAsNil = dd.TryDecodeAsNil()
		if f.seq == seqTypeChan ***REMOVED***
			if decodeAsNil ***REMOVED***
				rv.Send(reflect.Zero(rtelem0))
				continue
			***REMOVED***
			if rtelem0Mut || !rv9.IsValid() ***REMOVED*** // || (rtElem0Kind == reflect.Ptr && rv9.IsNil()) ***REMOVED***
				rv9 = reflect.New(rtelem0).Elem()
			***REMOVED***
			if fn == nil ***REMOVED***
				fn = d.cf.get(rtelem, true, true)
			***REMOVED***
			d.decodeValue(rv9, fn, true)
			rv.Send(rv9)
		***REMOVED*** else ***REMOVED***
			// if indefinite, etc, then expand the slice if necessary
			var decodeIntoBlank bool
			if j >= rvlen ***REMOVED***
				if f.seq == seqTypeArray ***REMOVED***
					d.arrayCannotExpand(rvlen, j+1)
					decodeIntoBlank = true
				***REMOVED*** else ***REMOVED*** // if f.seq == seqTypeSlice
					// rv = reflect.Append(rv, reflect.Zero(rtelem0)) // append logic + varargs
					var rvcap2 int
					var rvErrmsg2 string
					rv9, rvcap2, rvChanged, rvErrmsg2 =
						expandSliceRV(rv, ti.rt, rvCanset, rtelem0Size, 1, rvlen, rvcap)
					if rvErrmsg2 != "" ***REMOVED***
						d.errorf(rvErrmsg2)
					***REMOVED***
					rvlen++
					if rvChanged ***REMOVED***
						rv = rv9
						rvcap = rvcap2
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if decodeIntoBlank ***REMOVED***
				if !decodeAsNil ***REMOVED***
					d.swallow()
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				rv9 = rv.Index(j)
				if d.h.SliceElementReset || decodeAsNil ***REMOVED***
					if !rtelem0ZeroValid ***REMOVED***
						rtelem0ZeroValid = true
						rtelem0Zero = reflect.Zero(rtelem0)
					***REMOVED***
					rv9.Set(rtelem0Zero)
				***REMOVED***
				if decodeAsNil ***REMOVED***
					continue
				***REMOVED***

				if fn == nil ***REMOVED***
					fn = d.cf.get(rtelem, true, true)
				***REMOVED***
				d.decodeValue(rv9, fn, true)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if f.seq == seqTypeSlice ***REMOVED***
		if j < rvlen ***REMOVED***
			if rv.CanSet() ***REMOVED***
				rv.SetLen(j)
			***REMOVED*** else if rvCanset ***REMOVED***
				rv = rv.Slice(0, j)
				rvChanged = true
			***REMOVED*** // else ***REMOVED*** d.errorf("kSlice: cannot change non-settable slice") ***REMOVED***
			rvlen = j
		***REMOVED*** else if j == 0 && rv.IsNil() ***REMOVED***
			if rvCanset ***REMOVED***
				rv = reflect.MakeSlice(ti.rt, 0, 0)
				rvChanged = true
			***REMOVED*** // else ***REMOVED*** d.errorf("kSlice: cannot change non-settable slice") ***REMOVED***
		***REMOVED***
	***REMOVED***
	slh.End()

	if rvChanged ***REMOVED*** // infers rvCanset=true, so it can be reset
		rv0.Set(rv)
	***REMOVED***
***REMOVED***

// func (d *Decoder) kArray(f *codecFnInfo, rv reflect.Value) ***REMOVED***
// 	// d.decodeValueFn(rv.Slice(0, rv.Len()))
// 	f.kSlice(rv.Slice(0, rv.Len()))
// ***REMOVED***

func (d *Decoder) kMap(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	dd := d.d
	containerLen := dd.ReadMapStart()
	elemsep := d.esep
	ti := f.ti
	if rv.IsNil() ***REMOVED***
		rv.Set(makeMapReflect(ti.rt, containerLen))
	***REMOVED***

	if containerLen == 0 ***REMOVED***
		dd.ReadMapEnd()
		return
	***REMOVED***

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := vtype.Kind()

	var keyFn, valFn *codecFn
	var ktypeLo, vtypeLo reflect.Type

	for ktypeLo = ktype; ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() ***REMOVED***
	***REMOVED***

	for vtypeLo = vtype; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() ***REMOVED***
	***REMOVED***

	var mapGet, mapSet bool
	rvvImmut := isImmutableKind(vtypeKind)
	if !d.h.MapValueReset ***REMOVED***
		// if pointer, mapGet = true
		// if interface, mapGet = true if !DecodeNakedAlways (else false)
		// if builtin, mapGet = false
		// else mapGet = true
		if vtypeKind == reflect.Ptr ***REMOVED***
			mapGet = true
		***REMOVED*** else if vtypeKind == reflect.Interface ***REMOVED***
			if !d.h.InterfaceReset ***REMOVED***
				mapGet = true
			***REMOVED***
		***REMOVED*** else if !rvvImmut ***REMOVED***
			mapGet = true
		***REMOVED***
	***REMOVED***

	var rvk, rvkp, rvv, rvz reflect.Value
	rvkMut := !isImmutableKind(ktype.Kind()) // if ktype is immutable, then re-use the same rvk.
	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId
	hasLen := containerLen > 0
	var kstrbs []byte
	d.cfer()
	for j := 0; (hasLen && j < containerLen) || !(hasLen || dd.CheckBreak()); j++ ***REMOVED***
		if rvkMut || !rvkp.IsValid() ***REMOVED***
			rvkp = reflect.New(ktype)
			rvk = rvkp.Elem()
		***REMOVED***
		if elemsep ***REMOVED***
			dd.ReadMapElemKey()
		***REMOVED***
		if false && dd.TryDecodeAsNil() ***REMOVED*** // nil cannot be a map key, so disregard this block
			// Previously, if a nil key, we just ignored the mapped value and continued.
			// However, that makes the result of encoding and then decoding map[intf]intf***REMOVED***nil:nil***REMOVED***
			// to be an empty map.
			// Instead, we treat a nil key as the zero value of the type.
			rvk.Set(reflect.Zero(ktype))
		***REMOVED*** else if ktypeIsString ***REMOVED***
			kstrbs = dd.DecodeStringAsBytes()
			rvk.SetString(stringView(kstrbs))
			// NOTE: if doing an insert, you MUST use a real string (not stringview)
		***REMOVED*** else ***REMOVED***
			if keyFn == nil ***REMOVED***
				keyFn = d.cf.get(ktypeLo, true, true)
			***REMOVED***
			d.decodeValue(rvk, keyFn, true)
		***REMOVED***
		// special case if a byte array.
		if ktypeIsIntf ***REMOVED***
			if rvk2 := rvk.Elem(); rvk2.IsValid() ***REMOVED***
				if rvk2.Type() == uint8SliceTyp ***REMOVED***
					rvk = reflect.ValueOf(d.string(rvk2.Bytes()))
				***REMOVED*** else ***REMOVED***
					rvk = rvk2
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if elemsep ***REMOVED***
			dd.ReadMapElemValue()
		***REMOVED***

		// Brittle, but OK per TryDecodeAsNil() contract.
		// i.e. TryDecodeAsNil never shares slices with other decDriver procedures
		if dd.TryDecodeAsNil() ***REMOVED***
			if ktypeIsString ***REMOVED***
				rvk.SetString(d.string(kstrbs))
			***REMOVED***
			if d.h.DeleteOnNilMapValue ***REMOVED***
				rv.SetMapIndex(rvk, reflect.Value***REMOVED******REMOVED***)
			***REMOVED*** else ***REMOVED***
				rv.SetMapIndex(rvk, reflect.Zero(vtype))
			***REMOVED***
			continue
		***REMOVED***

		mapSet = true // set to false if u do a get, and its a non-nil pointer
		if mapGet ***REMOVED***
			// mapGet true only in case where kind=Ptr|Interface or kind is otherwise mutable.
			rvv = rv.MapIndex(rvk)
			if !rvv.IsValid() ***REMOVED***
				rvv = reflect.New(vtype).Elem()
			***REMOVED*** else if vtypeKind == reflect.Ptr ***REMOVED***
				if rvv.IsNil() ***REMOVED***
					rvv = reflect.New(vtype).Elem()
				***REMOVED*** else ***REMOVED***
					mapSet = false
				***REMOVED***
			***REMOVED*** else if vtypeKind == reflect.Interface ***REMOVED***
				// not addressable, and thus not settable.
				// e MUST create a settable/addressable variant
				rvv2 := reflect.New(rvv.Type()).Elem()
				if !rvv.IsNil() ***REMOVED***
					rvv2.Set(rvv)
				***REMOVED***
				rvv = rvv2
			***REMOVED***
			// else it is ~mutable, and we can just decode into it directly
		***REMOVED*** else if rvvImmut ***REMOVED***
			if !rvz.IsValid() ***REMOVED***
				rvz = reflect.New(vtype).Elem()
			***REMOVED***
			rvv = rvz
		***REMOVED*** else ***REMOVED***
			rvv = reflect.New(vtype).Elem()
		***REMOVED***

		// We MUST be done with the stringview of the key, before decoding the value
		// so that we don't bastardize the reused byte array.
		if mapSet && ktypeIsString ***REMOVED***
			rvk.SetString(d.string(kstrbs))
		***REMOVED***
		if valFn == nil ***REMOVED***
			valFn = d.cf.get(vtypeLo, true, true)
		***REMOVED***
		d.decodeValue(rvv, valFn, true)
		// d.decodeValueFn(rvv, valFn)
		if mapSet ***REMOVED***
			rv.SetMapIndex(rvk, rvv)
		***REMOVED***
		// if ktypeIsString ***REMOVED***
		// 	// keepAlive4StringView(kstrbs) // not needed, as reference is outside loop
		// ***REMOVED***
	***REMOVED***

	dd.ReadMapEnd()
***REMOVED***

// decNaked is used to keep track of the primitives decoded.
// Without it, we would have to decode each primitive and wrap it
// in an interface***REMOVED******REMOVED***, causing an allocation.
// In this model, the primitives are decoded in a "pseudo-atomic" fashion,
// so we can rest assured that no other decoding happens while these
// primitives are being decoded.
//
// maps and arrays are not handled by this mechanism.
// However, RawExt is, and we accommodate for extensions that decode
// RawExt from DecodeNaked, but need to decode the value subsequently.
// kInterfaceNaked and swallow, which call DecodeNaked, handle this caveat.
//
// However, decNaked also keeps some arrays of default maps and slices
// used in DecodeNaked. This way, we can get a pointer to it
// without causing a new heap allocation.
//
// kInterfaceNaked will ensure that there is no allocation for the common
// uses.

type decNakedContainers struct ***REMOVED***
	// array/stacks for reducing allocation
	// keep arrays at the bottom? Chance is that they are not used much.
	ia [arrayCacheLen]interface***REMOVED******REMOVED***
	ma [arrayCacheLen]map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	na [arrayCacheLen]map[string]interface***REMOVED******REMOVED***
	sa [arrayCacheLen][]interface***REMOVED******REMOVED***

	// ria [arrayCacheLen]reflect.Value // not needed, as we decode directly into &ia[n]
	rma, rna, rsa [arrayCacheLen]reflect.Value // reflect.Value mapping to above
***REMOVED***

func (n *decNakedContainers) init() ***REMOVED***
	for i := 0; i < arrayCacheLen; i++ ***REMOVED***
		// n.ria[i] = reflect.ValueOf(&(n.ia[i])).Elem()
		n.rma[i] = reflect.ValueOf(&(n.ma[i])).Elem()
		n.rna[i] = reflect.ValueOf(&(n.na[i])).Elem()
		n.rsa[i] = reflect.ValueOf(&(n.sa[i])).Elem()
	***REMOVED***
***REMOVED***

type decNaked struct ***REMOVED***
	// r RawExt // used for RawExt, uint, []byte.

	// primitives below
	u uint64
	i int64
	f float64
	l []byte
	s string

	// ---- cpu cache line boundary?
	t time.Time
	b bool

	// state
	v              valueType
	li, lm, ln, ls int8
	inited         bool

	*decNakedContainers

	ru, ri, rf, rl, rs, rb, rt reflect.Value // mapping to the primitives above

	// _ [6]uint64 // padding // no padding - rt goes into next cache line
***REMOVED***

func (n *decNaked) init() ***REMOVED***
	if n.inited ***REMOVED***
		return
	***REMOVED***
	n.ru = reflect.ValueOf(&n.u).Elem()
	n.ri = reflect.ValueOf(&n.i).Elem()
	n.rf = reflect.ValueOf(&n.f).Elem()
	n.rl = reflect.ValueOf(&n.l).Elem()
	n.rs = reflect.ValueOf(&n.s).Elem()
	n.rt = reflect.ValueOf(&n.t).Elem()
	n.rb = reflect.ValueOf(&n.b).Elem()

	n.inited = true
	// n.rr[] = reflect.ValueOf(&n.)
***REMOVED***

func (n *decNaked) initContainers() ***REMOVED***
	if n.decNakedContainers == nil ***REMOVED***
		n.decNakedContainers = new(decNakedContainers)
		n.decNakedContainers.init()
	***REMOVED***
***REMOVED***

func (n *decNaked) reset() ***REMOVED***
	if n == nil ***REMOVED***
		return
	***REMOVED***
	n.li, n.lm, n.ln, n.ls = 0, 0, 0, 0
***REMOVED***

type rtid2rv struct ***REMOVED***
	rtid uintptr
	rv   reflect.Value
***REMOVED***

// --------------

type decReaderSwitch struct ***REMOVED***
	rb bytesDecReader
	// ---- cpu cache line boundary?
	ri       *ioDecReader
	mtr, str bool // whether maptype or slicetype are known types

	be    bool // is binary encoding
	bytes bool // is bytes reader
	js    bool // is json handle
	jsms  bool // is json handle, and MapKeyAsString
	esep  bool // has elem separators
***REMOVED***

// TODO: Uncomment after mid-stack inlining enabled in go 1.10
//
// func (z *decReaderSwitch) unreadn1() ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		z.rb.unreadn1()
// 	***REMOVED*** else ***REMOVED***
// 		z.ri.unreadn1()
// 	***REMOVED***
// ***REMOVED***
// func (z *decReaderSwitch) readx(n int) []byte ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.readx(n)
// 	***REMOVED***
// 	return z.ri.readx(n)
// ***REMOVED***
// func (z *decReaderSwitch) readb(s []byte) ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		z.rb.readb(s)
// 	***REMOVED*** else ***REMOVED***
// 		z.ri.readb(s)
// 	***REMOVED***
// ***REMOVED***
// func (z *decReaderSwitch) readn1() uint8 ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.readn1()
// 	***REMOVED***
// 	return z.ri.readn1()
// ***REMOVED***
// func (z *decReaderSwitch) numread() int ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.numread()
// 	***REMOVED***
// 	return z.ri.numread()
// ***REMOVED***
// func (z *decReaderSwitch) track() ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		z.rb.track()
// 	***REMOVED*** else ***REMOVED***
// 		z.ri.track()
// 	***REMOVED***
// ***REMOVED***
// func (z *decReaderSwitch) stopTrack() []byte ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.stopTrack()
// 	***REMOVED***
// 	return z.ri.stopTrack()
// ***REMOVED***
// func (z *decReaderSwitch) skip(accept *bitset256) (token byte) ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.skip(accept)
// 	***REMOVED***
// 	return z.ri.skip(accept)
// ***REMOVED***
// func (z *decReaderSwitch) readTo(in []byte, accept *bitset256) (out []byte) ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.readTo(in, accept)
// 	***REMOVED***
// 	return z.ri.readTo(in, accept)
// ***REMOVED***
// func (z *decReaderSwitch) readUntil(in []byte, stop byte) (out []byte) ***REMOVED***
// 	if z.bytes ***REMOVED***
// 		return z.rb.readUntil(in, stop)
// 	***REMOVED***
// 	return z.ri.readUntil(in, stop)
// ***REMOVED***

const decScratchByteArrayLen = cacheLineSize - 8

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct ***REMOVED***
	panicHdl
	// hopefully, reduce derefencing cost by laying the decReader inside the Decoder.
	// Try to put things that go together to fit within a cache line (8 words).

	d decDriver
	// NOTE: Decoder shouldn't call it's read methods,
	// as the handler MAY need to do some coordination.
	r  decReader
	h  *BasicHandle
	bi *bufioDecReader
	// cache the mapTypeId and sliceTypeId for faster comparisons
	mtid uintptr
	stid uintptr

	// ---- cpu cache line boundary?
	decReaderSwitch

	// ---- cpu cache line boundary?
	codecFnPooler
	// cr containerStateRecv
	n   *decNaked
	nsp *sync.Pool
	err error

	// ---- cpu cache line boundary?
	b  [decScratchByteArrayLen]byte // scratch buffer, used by Decoder and xxxEncDrivers
	is map[string]string            // used for interning strings

	// padding - false sharing help // modify 232 if Decoder struct changes.
	// _ [cacheLineSize - 232%cacheLineSize]byte
***REMOVED***

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to pass in a memory buffered reader
// (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder ***REMOVED***
	d := newDecoder(h)
	d.Reset(r)
	return d
***REMOVED***

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder ***REMOVED***
	d := newDecoder(h)
	d.ResetBytes(in)
	return d
***REMOVED***

var defaultDecNaked decNaked

func newDecoder(h Handle) *Decoder ***REMOVED***
	d := &Decoder***REMOVED***h: h.getBasicHandle(), err: errDecoderNotInitialized***REMOVED***
	d.hh = h
	d.be = h.isBinary()
	// NOTE: do not initialize d.n here. It is lazily initialized in d.naked()
	var jh *JsonHandle
	jh, d.js = h.(*JsonHandle)
	if d.js ***REMOVED***
		d.jsms = jh.MapKeyAsString
	***REMOVED***
	d.esep = d.hh.hasElemSeparators()
	if d.h.InternString ***REMOVED***
		d.is = make(map[string]string, 32)
	***REMOVED***
	d.d = h.newDecDriver(d)
	// d.cr, _ = d.d.(containerStateRecv)
	return d
***REMOVED***

func (d *Decoder) resetCommon() ***REMOVED***
	d.n.reset()
	d.d.reset()
	d.err = nil
	// reset all things which were cached from the Handle, but could change
	d.mtid, d.stid = 0, 0
	d.mtr, d.str = false, false
	if d.h.MapType != nil ***REMOVED***
		d.mtid = rt2id(d.h.MapType)
		d.mtr = fastpathAV.index(d.mtid) != -1
	***REMOVED***
	if d.h.SliceType != nil ***REMOVED***
		d.stid = rt2id(d.h.SliceType)
		d.str = fastpathAV.index(d.stid) != -1
	***REMOVED***
***REMOVED***

// Reset the Decoder with a new Reader to decode from,
// clearing all state from last run(s).
func (d *Decoder) Reset(r io.Reader) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	if d.bi == nil ***REMOVED***
		d.bi = new(bufioDecReader)
	***REMOVED***
	d.bytes = false
	if d.h.ReaderBufferSize > 0 ***REMOVED***
		d.bi.buf = make([]byte, 0, d.h.ReaderBufferSize)
		d.bi.reset(r)
		d.r = d.bi
	***REMOVED*** else ***REMOVED***
		// d.ri.x = &d.b
		// d.s = d.sa[:0]
		if d.ri == nil ***REMOVED***
			d.ri = new(ioDecReader)
		***REMOVED***
		d.ri.reset(r)
		d.r = d.ri
	***REMOVED***
	d.resetCommon()
***REMOVED***

// ResetBytes resets the Decoder with a new []byte to decode from,
// clearing all state from last run(s).
func (d *Decoder) ResetBytes(in []byte) ***REMOVED***
	if in == nil ***REMOVED***
		return
	***REMOVED***
	d.bytes = true
	d.rb.reset(in)
	d.r = &d.rb
	d.resetCommon()
***REMOVED***

// naked must be called before each call to .DecodeNaked,
// as they will use it.
func (d *Decoder) naked() *decNaked ***REMOVED***
	if d.n == nil ***REMOVED***
		// consider one of:
		//   - get from sync.Pool  (if GC is frequent, there's no value here)
		//   - new alloc           (safest. only init'ed if it a naked decode will be done)
		//   - field in Decoder    (makes the Decoder struct very big)
		// To support using a decoder where a DecodeNaked is not needed,
		// we prefer #1 or #2.
		// d.n = new(decNaked) // &d.nv // new(decNaked) // grab from a sync.Pool
		// d.n.init()
		var v interface***REMOVED******REMOVED***
		d.nsp, v = pool.decNaked()
		d.n = v.(*decNaked)
	***REMOVED***
	return d.n
***REMOVED***

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
//
// Sample usages:
//   // Decoding into a non-nil typed value
//   var f float32
//   err = codec.NewDecoder(r, handle).Decode(&f)
//
//   // Decoding into nil interface
//   var v interface***REMOVED******REMOVED***
//   dec := codec.NewDecoder(r, handle)
//   err = dec.Decode(&v)
//
// When decoding into a nil interface***REMOVED******REMOVED***, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64.
//   - Other values are decoded appropriately depending on the type:
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
// Configurations exist on the Handle to override defaults
// (e.g. for MapType, SliceType and how to decode raw bytes).
//
// When decoding into a non-nil interface***REMOVED******REMOVED*** value, the mode of encoding is based on the
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
//
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// When decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
// Note: we allow nil values in the stream anywhere except for map keys.
// A nil value in the encoded stream where a map key is expected is treated as an error.
func (d *Decoder) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	// need to call defer directly, else it seems the recover is not fully handled
	defer panicToErrs2(d, &d.err, &err)
	defer d.alwaysAtEnd()
	d.MustDecode(v)
	return
***REMOVED***

// MustDecode is like Decode, but panics if unable to Decode.
// This provides insight to the code location that triggered the error.
func (d *Decoder) MustDecode(v interface***REMOVED******REMOVED***) ***REMOVED***
	// TODO: Top-level: ensure that v is a pointer and not nil.
	if d.err != nil ***REMOVED***
		panic(d.err)
	***REMOVED***
	if d.d.TryDecodeAsNil() ***REMOVED***
		setZero(v)
	***REMOVED*** else ***REMOVED***
		d.decode(v)
	***REMOVED***
	d.alwaysAtEnd()
	// xprintf(">>>>>>>> >>>>>>>> num decFns: %v\n", d.cf.sn)
***REMOVED***

// // this is not a smart swallow, as it allocates objects and does unnecessary work.
// func (d *Decoder) swallowViaHammer() ***REMOVED***
// 	var blank interface***REMOVED******REMOVED***
// 	d.decodeValueNoFn(reflect.ValueOf(&blank).Elem())
// ***REMOVED***

func (d *Decoder) alwaysAtEnd() ***REMOVED***
	if d.n != nil ***REMOVED***
		// if n != nil, then nsp != nil (they are always set together)
		d.nsp.Put(d.n)
		d.n, d.nsp = nil, nil
	***REMOVED***
	d.codecFnPooler.alwaysAtEnd()
***REMOVED***

func (d *Decoder) swallow() ***REMOVED***
	// smarter decode that just swallows the content
	dd := d.d
	if dd.TryDecodeAsNil() ***REMOVED***
		return
	***REMOVED***
	elemsep := d.esep
	switch dd.ContainerType() ***REMOVED***
	case valueTypeMap:
		containerLen := dd.ReadMapStart()
		hasLen := containerLen >= 0
		for j := 0; (hasLen && j < containerLen) || !(hasLen || dd.CheckBreak()); j++ ***REMOVED***
			// if clenGtEqualZero ***REMOVED***if j >= containerLen ***REMOVED***break***REMOVED*** ***REMOVED*** else if dd.CheckBreak() ***REMOVED***break***REMOVED***
			if elemsep ***REMOVED***
				dd.ReadMapElemKey()
			***REMOVED***
			d.swallow()
			if elemsep ***REMOVED***
				dd.ReadMapElemValue()
			***REMOVED***
			d.swallow()
		***REMOVED***
		dd.ReadMapEnd()
	case valueTypeArray:
		containerLen := dd.ReadArrayStart()
		hasLen := containerLen >= 0
		for j := 0; (hasLen && j < containerLen) || !(hasLen || dd.CheckBreak()); j++ ***REMOVED***
			if elemsep ***REMOVED***
				dd.ReadArrayElem()
			***REMOVED***
			d.swallow()
		***REMOVED***
		dd.ReadArrayEnd()
	case valueTypeBytes:
		dd.DecodeBytes(d.b[:], true)
	case valueTypeString:
		dd.DecodeStringAsBytes()
	default:
		// these are all primitives, which we can get from decodeNaked
		// if RawExt using Value, complete the processing.
		n := d.naked()
		dd.DecodeNaked()
		if n.v == valueTypeExt && n.l == nil ***REMOVED***
			n.initContainers()
			if n.li < arrayCacheLen ***REMOVED***
				n.ia[n.li] = nil
				n.li++
				d.decode(&n.ia[n.li-1])
				n.ia[n.li-1] = nil
				n.li--
			***REMOVED*** else ***REMOVED***
				var v2 interface***REMOVED******REMOVED***
				d.decode(&v2)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func setZero(iv interface***REMOVED******REMOVED***) ***REMOVED***
	if iv == nil || definitelyNil(iv) ***REMOVED***
		return
	***REMOVED***
	var canDecode bool
	switch v := iv.(type) ***REMOVED***
	case *string:
		*v = ""
	case *bool:
		*v = false
	case *int:
		*v = 0
	case *int8:
		*v = 0
	case *int16:
		*v = 0
	case *int32:
		*v = 0
	case *int64:
		*v = 0
	case *uint:
		*v = 0
	case *uint8:
		*v = 0
	case *uint16:
		*v = 0
	case *uint32:
		*v = 0
	case *uint64:
		*v = 0
	case *float32:
		*v = 0
	case *float64:
		*v = 0
	case *[]uint8:
		*v = nil
	case *Raw:
		*v = nil
	case *time.Time:
		*v = time.Time***REMOVED******REMOVED***
	case reflect.Value:
		if v, canDecode = isDecodeable(v); canDecode && v.CanSet() ***REMOVED***
			v.Set(reflect.Zero(v.Type()))
		***REMOVED*** // TODO: else drain if chan, clear if map, set all to nil if slice???
	default:
		if !fastpathDecodeSetZeroTypeSwitch(iv) ***REMOVED***
			v := reflect.ValueOf(iv)
			if v, canDecode = isDecodeable(v); canDecode && v.CanSet() ***REMOVED***
				v.Set(reflect.Zero(v.Type()))
			***REMOVED*** // TODO: else drain if chan, clear if map, set all to nil if slice???
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *Decoder) decode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	// check nil and interfaces explicitly,
	// so that type switches just have a run of constant non-interface types.
	if iv == nil ***REMOVED***
		d.errorstr(errstrCannotDecodeIntoNil)
		return
	***REMOVED***
	if v, ok := iv.(Selfer); ok ***REMOVED***
		v.CodecDecodeSelf(d)
		return
	***REMOVED***

	switch v := iv.(type) ***REMOVED***
	// case nil:
	// case Selfer:

	case reflect.Value:
		v = d.ensureDecodeable(v)
		d.decodeValue(v, nil, true)

	case *string:
		*v = d.d.DecodeString()
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *float32:
		f64 := d.d.DecodeFloat64()
		if chkOvf.Float32(f64) ***REMOVED***
			d.errorf("float32 overflow: %v", f64)
		***REMOVED***
		*v = float32(f64)
	case *float64:
		*v = d.d.DecodeFloat64()
	case *[]uint8:
		*v = d.d.DecodeBytes(*v, false)
	case []uint8:
		b := d.d.DecodeBytes(v, false)
		if !(len(b) > 0 && len(b) == len(v) && &b[0] == &v[0]) ***REMOVED***
			copy(v, b)
		***REMOVED***
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface***REMOVED******REMOVED***:
		d.decodeValue(reflect.ValueOf(iv).Elem(), nil, true)
		// d.decodeValueNotNil(reflect.ValueOf(iv).Elem())

	default:
		if !fastpathDecodeTypeSwitch(iv, d) ***REMOVED***
			v := reflect.ValueOf(iv)
			v = d.ensureDecodeable(v)
			d.decodeValue(v, nil, false)
			// d.decodeValueFallback(v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *Decoder) decodeValue(rv reflect.Value, fn *codecFn, chkAll bool) ***REMOVED***
	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
	var rvp reflect.Value
	var rvpValid bool
	if rv.Kind() == reflect.Ptr ***REMOVED***
		rvpValid = true
		for ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rv.Set(reflect.New(rv.Type().Elem()))
			***REMOVED***
			rvp = rv
			rv = rv.Elem()
			if rv.Kind() != reflect.Ptr ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if fn == nil ***REMOVED***
		// always pass checkCodecSelfer=true, in case T or ****T is passed, where *T is a Selfer
		fn = d.cfer().get(rv.Type(), chkAll, true) // chkAll, chkAll)
	***REMOVED***
	if fn.i.addrD ***REMOVED***
		if rvpValid ***REMOVED***
			fn.fd(d, &fn.i, rvp)
		***REMOVED*** else if rv.CanAddr() ***REMOVED***
			fn.fd(d, &fn.i, rv.Addr())
		***REMOVED*** else if !fn.i.addrF ***REMOVED***
			fn.fd(d, &fn.i, rv)
		***REMOVED*** else ***REMOVED***
			d.errorf("cannot decode into a non-pointer value")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fn.fd(d, &fn.i, rv)
	***REMOVED***
	// return rv
***REMOVED***

func (d *Decoder) structFieldNotFound(index int, rvkencname string) ***REMOVED***
	// NOTE: rvkencname may be a stringView, so don't pass it to another function.
	if d.h.ErrorIfNoField ***REMOVED***
		if index >= 0 ***REMOVED***
			d.errorf("no matching struct field found when decoding stream array at index %v", index)
			return
		***REMOVED*** else if rvkencname != "" ***REMOVED***
			d.errorf("no matching struct field found when decoding stream map with key " + rvkencname)
			return
		***REMOVED***
	***REMOVED***
	d.swallow()
***REMOVED***

func (d *Decoder) arrayCannotExpand(sliceLen, streamLen int) ***REMOVED***
	if d.h.ErrorIfNoArrayExpand ***REMOVED***
		d.errorf("cannot expand array len during decode from %v to %v", sliceLen, streamLen)
	***REMOVED***
***REMOVED***

func isDecodeable(rv reflect.Value) (rv2 reflect.Value, canDecode bool) ***REMOVED***
	switch rv.Kind() ***REMOVED***
	case reflect.Array:
		return rv, true
	case reflect.Ptr:
		if !rv.IsNil() ***REMOVED***
			return rv.Elem(), true
		***REMOVED***
	case reflect.Slice, reflect.Chan, reflect.Map:
		if !rv.IsNil() ***REMOVED***
			return rv, true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) ensureDecodeable(rv reflect.Value) (rv2 reflect.Value) ***REMOVED***
	// decode can take any reflect.Value that is a inherently addressable i.e.
	//   - array
	//   - non-nil chan    (we will SEND to it)
	//   - non-nil slice   (we will set its elements)
	//   - non-nil map     (we will put into it)
	//   - non-nil pointer (we can "update" it)
	rv2, canDecode := isDecodeable(rv)
	if canDecode ***REMOVED***
		return
	***REMOVED***
	if !rv.IsValid() ***REMOVED***
		d.errorstr(errstrCannotDecodeIntoNil)
		return
	***REMOVED***
	if !rv.CanInterface() ***REMOVED***
		d.errorf("cannot decode into a value without an interface: %v", rv)
		return
	***REMOVED***
	rvi := rv2i(rv)
	rvk := rv.Kind()
	d.errorf("cannot decode into value of kind: %v, type: %T, %v", rvk, rvi, rvi)
	return
***REMOVED***

// Possibly get an interned version of a string
//
// This should mostly be used for map keys, where the key type is string.
// This is because keys of a map/struct are typically reused across many objects.
func (d *Decoder) string(v []byte) (s string) ***REMOVED***
	if d.is == nil ***REMOVED***
		return string(v) // don't return stringView, as we need a real string here.
	***REMOVED***
	s, ok := d.is[string(v)] // no allocation here, per go implementation
	if !ok ***REMOVED***
		s = string(v) // new allocation here
		d.is[s] = s
	***REMOVED***
	return s
***REMOVED***

// nextValueBytes returns the next value in the stream as a set of bytes.
func (d *Decoder) nextValueBytes() (bs []byte) ***REMOVED***
	d.d.uncacheRead()
	d.r.track()
	d.swallow()
	bs = d.r.stopTrack()
	return
***REMOVED***

func (d *Decoder) rawBytes() []byte ***REMOVED***
	// ensure that this is not a view into the bytes
	// i.e. make new copy always.
	bs := d.nextValueBytes()
	bs2 := make([]byte, len(bs))
	copy(bs2, bs)
	return bs2
***REMOVED***

func (d *Decoder) wrapErrstr(v interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	*err = fmt.Errorf("%s decode error [pos %d]: %v", d.hh.Name(), d.r.numread(), v)
***REMOVED***

// --------------------------------------------------

// decSliceHelper assists when decoding into a slice, from a map or an array in the stream.
// A slice can be set from a map or array in stream. This supports the MapBySlice interface.
type decSliceHelper struct ***REMOVED***
	d *Decoder
	// ct valueType
	array bool
***REMOVED***

func (d *Decoder) decSliceHelperStart() (x decSliceHelper, clen int) ***REMOVED***
	dd := d.d
	ctyp := dd.ContainerType()
	switch ctyp ***REMOVED***
	case valueTypeArray:
		x.array = true
		clen = dd.ReadArrayStart()
	case valueTypeMap:
		clen = dd.ReadMapStart() * 2
	default:
		d.errorf("only encoded map or array can be decoded into a slice (%d)", ctyp)
	***REMOVED***
	// x.ct = ctyp
	x.d = d
	return
***REMOVED***

func (x decSliceHelper) End() ***REMOVED***
	if x.array ***REMOVED***
		x.d.d.ReadArrayEnd()
	***REMOVED*** else ***REMOVED***
		x.d.d.ReadMapEnd()
	***REMOVED***
***REMOVED***

func (x decSliceHelper) ElemContainerState(index int) ***REMOVED***
	if x.array ***REMOVED***
		x.d.d.ReadArrayElem()
	***REMOVED*** else if index%2 == 0 ***REMOVED***
		x.d.d.ReadMapElemKey()
	***REMOVED*** else ***REMOVED***
		x.d.d.ReadMapElemValue()
	***REMOVED***
***REMOVED***

func decByteSlice(r decReader, clen, maxInitLen int, bs []byte) (bsOut []byte) ***REMOVED***
	if clen == 0 ***REMOVED***
		return zeroByteSlice
	***REMOVED***
	if len(bs) == clen ***REMOVED***
		bsOut = bs
		r.readb(bsOut)
	***REMOVED*** else if cap(bs) >= clen ***REMOVED***
		bsOut = bs[:clen]
		r.readb(bsOut)
	***REMOVED*** else ***REMOVED***
		// bsOut = make([]byte, clen)
		len2 := decInferLen(clen, maxInitLen, 1)
		bsOut = make([]byte, len2)
		r.readb(bsOut)
		for len2 < clen ***REMOVED***
			len3 := decInferLen(clen-len2, maxInitLen, 1)
			bs3 := bsOut
			bsOut = make([]byte, len2+len3)
			copy(bsOut, bs3)
			r.readb(bsOut[len2:])
			len2 += len3
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func detachZeroCopyBytes(isBytesReader bool, dest []byte, in []byte) (out []byte) ***REMOVED***
	if xlen := len(in); xlen > 0 ***REMOVED***
		if isBytesReader || xlen <= scratchByteArrayLen ***REMOVED***
			if cap(dest) >= xlen ***REMOVED***
				out = dest[:xlen]
			***REMOVED*** else ***REMOVED***
				out = make([]byte, xlen)
			***REMOVED***
			copy(out, in)
			return
		***REMOVED***
	***REMOVED***
	return in
***REMOVED***

// decInferLen will infer a sensible length, given the following:
//    - clen: length wanted.
//    - maxlen: max length to be returned.
//      if <= 0, it is unset, and we infer it based on the unit size
//    - unit: number of bytes for each element of the collection
func decInferLen(clen, maxlen, unit int) (rvlen int) ***REMOVED***
	// handle when maxlen is not set i.e. <= 0
	if clen <= 0 ***REMOVED***
		return
	***REMOVED***
	if unit == 0 ***REMOVED***
		return clen
	***REMOVED***
	if maxlen <= 0 ***REMOVED***
		// no maxlen defined. Use maximum of 256K memory, with a floor of 4K items.
		// maxlen = 256 * 1024 / unit
		// if maxlen < (4 * 1024) ***REMOVED***
		// 	maxlen = 4 * 1024
		// ***REMOVED***
		if unit < (256 / 4) ***REMOVED***
			maxlen = 256 * 1024 / unit
		***REMOVED*** else ***REMOVED***
			maxlen = 4 * 1024
		***REMOVED***
	***REMOVED***
	if clen > maxlen ***REMOVED***
		rvlen = maxlen
	***REMOVED*** else ***REMOVED***
		rvlen = clen
	***REMOVED***
	return
***REMOVED***

func expandSliceRV(s reflect.Value, st reflect.Type, canChange bool, stElemSize, num, slen, scap int) (
	s2 reflect.Value, scap2 int, changed bool, err string) ***REMOVED***
	l1 := slen + num // new slice length
	if l1 < slen ***REMOVED***
		err = errmsgExpandSliceOverflow
		return
	***REMOVED***
	if l1 <= scap ***REMOVED***
		if s.CanSet() ***REMOVED***
			s.SetLen(l1)
		***REMOVED*** else if canChange ***REMOVED***
			s2 = s.Slice(0, l1)
			scap2 = scap
			changed = true
		***REMOVED*** else ***REMOVED***
			err = errmsgExpandSliceCannotChange
			return
		***REMOVED***
		return
	***REMOVED***
	if !canChange ***REMOVED***
		err = errmsgExpandSliceCannotChange
		return
	***REMOVED***
	scap2 = growCap(scap, stElemSize, num)
	s2 = reflect.MakeSlice(st, l1, scap2)
	changed = true
	reflect.Copy(s2, s)
	return
***REMOVED***

func decReadFull(r io.Reader, bs []byte) (n int, err error) ***REMOVED***
	var nn int
	for n < len(bs) && err == nil ***REMOVED***
		nn, err = r.Read(bs[n:])
		if nn > 0 ***REMOVED***
			if err == io.EOF ***REMOVED***
				// leave EOF for next time
				err = nil
			***REMOVED***
			n += nn
		***REMOVED***
	***REMOVED***

	// do not do this - it serves no purpose
	// if n != len(bs) && err == io.EOF ***REMOVED*** err = io.ErrUnexpectedEOF ***REMOVED***
	return
***REMOVED***
