// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "io"

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface ***REMOVED***
	unreadn1()
	// readx will use the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// just return a view of the []byte being decoded from.
	// Ensure you call detachZeroCopyBytes later if this needs to be sent outside codec control.
	readx(n uint) []byte
	readb([]byte)
	readn1() uint8
	// read up to 7 bytes at a time
	readn(num uint8) (v [rwNLen]byte)
	numread() uint // number of bytes read
	track()
	stopTrack() []byte

	// skip will skip any byte that matches, and return the first non-matching byte
	skip(accept *bitset256) (token byte)
	// readTo will read any byte that matches, stopping once no-longer matching.
	readTo(accept *bitset256) (out []byte)
	// readUntil will read, only stopping once it matches the 'stop' byte.
	readUntil(stop byte, includeLast bool) (out []byte)
***REMOVED***

// ------------------------------------------------

type unreadByteStatus uint8

// unreadByteStatus goes from
// undefined (when initialized) -- (read) --> canUnread -- (unread) --> canRead ...
const (
	unreadByteUndefined unreadByteStatus = iota
	unreadByteCanRead
	unreadByteCanUnread
)

// --------------------

type ioDecReaderCommon struct ***REMOVED***
	r io.Reader // the reader passed in

	n uint // num read

	l   byte             // last byte
	ls  unreadByteStatus // last byte status
	trb bool             // tracking bytes turned on
	_   bool
	b   [4]byte // tiny buffer for reading single bytes

	blist *bytesFreelist

	tr   []byte // buffer for tracking bytes
	bufr []byte // buffer for readTo/readUntil
***REMOVED***

func (z *ioDecReaderCommon) last() byte ***REMOVED***
	return z.l
***REMOVED***

func (z *ioDecReaderCommon) reset(r io.Reader, blist *bytesFreelist) ***REMOVED***
	z.blist = blist
	z.r = r
	z.ls = unreadByteUndefined
	z.l, z.n = 0, 0
	z.trb = false
***REMOVED***

func (z *ioDecReaderCommon) numread() uint ***REMOVED***
	return z.n
***REMOVED***

func (z *ioDecReaderCommon) track() ***REMOVED***
	z.tr = z.blist.check(z.tr, 256)[:0]
	z.trb = true
***REMOVED***

func (z *ioDecReaderCommon) stopTrack() (bs []byte) ***REMOVED***
	z.trb = false
	return z.tr
***REMOVED***

// ------------------------------------------

// ioDecReader is a decReader that reads off an io.Reader.
//
// It also has a fallback implementation of ByteScanner if needed.
type ioDecReader struct ***REMOVED***
	ioDecReaderCommon

	// rr io.Reader
	br io.ByteScanner

	x [64 + 16]byte // for: get struct field name, swallow valueTypeBytes, etc
	// _ [1]uint64                 // padding
***REMOVED***

func (z *ioDecReader) reset(r io.Reader, blist *bytesFreelist) ***REMOVED***
	z.ioDecReaderCommon.reset(r, blist)

	z.br, _ = r.(io.ByteScanner)
***REMOVED***

func (z *ioDecReader) Read(p []byte) (n int, err error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return
	***REMOVED***
	var firstByte bool
	if z.ls == unreadByteCanRead ***REMOVED***
		z.ls = unreadByteCanUnread
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
		z.ls = unreadByteCanUnread
	***REMOVED***
	if firstByte ***REMOVED***
		n++
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) ReadByte() (c byte, err error) ***REMOVED***
	if z.br != nil ***REMOVED***
		c, err = z.br.ReadByte()
		if err == nil ***REMOVED***
			z.l = c
			z.ls = unreadByteCanUnread
		***REMOVED***
		return
	***REMOVED***

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
	if z.br != nil ***REMOVED***
		err = z.br.UnreadByte()
		if err == nil ***REMOVED***
			z.ls = unreadByteCanRead
		***REMOVED***
		return
	***REMOVED***

	switch z.ls ***REMOVED***
	case unreadByteCanUnread:
		z.ls = unreadByteCanRead
	case unreadByteCanRead:
		err = errDecUnreadByteLastByteNotRead
	case unreadByteUndefined:
		err = errDecUnreadByteNothingToRead
	default:
		err = errDecUnreadByteUnknown
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readn(num uint8) (bs [rwNLen]byte) ***REMOVED***
	z.readb(bs[:num])
	// copy(bs[:], z.readx(uint(num)))
	return
***REMOVED***

func (z *ioDecReader) readx(n uint) (bs []byte) ***REMOVED***
	if n == 0 ***REMOVED***
		return
	***REMOVED***
	if n < uint(len(z.x)) ***REMOVED***
		bs = z.x[:n]
	***REMOVED*** else ***REMOVED***
		bs = make([]byte, n)
	***REMOVED***
	if _, err := decReadFull(z.r, bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n += uint(len(bs))
	if z.trb ***REMOVED***
		z.tr = append(z.tr, bs...)
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readb(bs []byte) ***REMOVED***
	if len(bs) == 0 ***REMOVED***
		return
	***REMOVED***
	if _, err := decReadFull(z.r, bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n += uint(len(bs))
	if z.trb ***REMOVED***
		z.tr = append(z.tr, bs...)
	***REMOVED***
***REMOVED***

func (z *ioDecReader) readn1eof() (b uint8, eof bool) ***REMOVED***
	b, err := z.ReadByte()
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
	b, err := z.ReadByte()
	if err == nil ***REMOVED***
		z.n++
		if z.trb ***REMOVED***
			z.tr = append(z.tr, b)
		***REMOVED***
		return
	***REMOVED***
	panic(err)
***REMOVED***

func (z *ioDecReader) skip(accept *bitset256) (token byte) ***REMOVED***
	var eof bool
LOOP:
	token, eof = z.readn1eof()
	if eof ***REMOVED***
		return
	***REMOVED***
	if accept.isset(token) ***REMOVED***
		goto LOOP
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readTo(accept *bitset256) []byte ***REMOVED***
	z.bufr = z.blist.check(z.bufr, 256)[:0]
LOOP:
	token, eof := z.readn1eof()
	if eof ***REMOVED***
		return z.bufr
	***REMOVED***
	if accept.isset(token) ***REMOVED***
		z.bufr = append(z.bufr, token)
		goto LOOP
	***REMOVED***
	z.unreadn1()
	return z.bufr
***REMOVED***

func (z *ioDecReader) readUntil(stop byte, includeLast bool) []byte ***REMOVED***
	z.bufr = z.blist.check(z.bufr, 256)[:0]
LOOP:
	token, eof := z.readn1eof()
	if eof ***REMOVED***
		panic(io.EOF)
	***REMOVED***
	z.bufr = append(z.bufr, token)
	if token == stop ***REMOVED***
		if includeLast ***REMOVED***
			return z.bufr
		***REMOVED***
		return z.bufr[:len(z.bufr)-1]
	***REMOVED***
	goto LOOP
***REMOVED***

//go:noinline
func (z *ioDecReader) unreadn1() ***REMOVED***
	err := z.UnreadByte()
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

// ------------------------------------

type bufioDecReader struct ***REMOVED***
	ioDecReaderCommon

	c   uint // cursor
	buf []byte
***REMOVED***

func (z *bufioDecReader) reset(r io.Reader, bufsize int, blist *bytesFreelist) ***REMOVED***
	z.ioDecReaderCommon.reset(r, blist)
	z.c = 0
	if cap(z.buf) < bufsize ***REMOVED***
		z.buf = blist.get(bufsize)
	***REMOVED***
	z.buf = z.buf[:0]
***REMOVED***

func (z *bufioDecReader) readb(p []byte) ***REMOVED***
	var n = uint(copy(p, z.buf[z.c:]))
	z.n += n
	z.c += n
	if len(p) == int(n) ***REMOVED***
		if z.trb ***REMOVED***
			z.tr = append(z.tr, p...)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		z.readbFill(p, n)
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) readbFill(p0 []byte, n uint) ***REMOVED***
	// at this point, there's nothing in z.buf to read (z.buf is fully consumed)
	p := p0[n:]
	var n2 uint
	var err error
	if len(p) > cap(z.buf) ***REMOVED***
		n2, err = decReadFull(z.r, p)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		n += n2
		z.n += n2
		// always keep last byte in z.buf
		z.buf = z.buf[:1]
		z.buf[0] = p[len(p)-1]
		z.c = 1
		if z.trb ***REMOVED***
			z.tr = append(z.tr, p0[:n]...)
		***REMOVED***
		return
	***REMOVED***
	// z.c is now 0, and len(p) <= cap(z.buf)
LOOP:
	// for len(p) > 0 && z.err == nil ***REMOVED***
	if len(p) > 0 ***REMOVED***
		z.buf = z.buf[0:cap(z.buf)]
		var n1 int
		n1, err = z.r.Read(z.buf)
		n2 = uint(n1)
		if n2 == 0 && err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		z.buf = z.buf[:n2]
		n2 = uint(copy(p, z.buf))
		z.c = n2
		n += n2
		z.n += n2
		p = p[n2:]
		goto LOOP
	***REMOVED***
	if z.c == 0 ***REMOVED***
		z.buf = z.buf[:1]
		z.buf[0] = p[len(p)-1]
		z.c = 1
	***REMOVED***
	if z.trb ***REMOVED***
		z.tr = append(z.tr, p0[:n]...)
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) last() byte ***REMOVED***
	return z.buf[z.c-1]
***REMOVED***

func (z *bufioDecReader) readn1() (b byte) ***REMOVED***
	// fast-path, so we elide calling into Read() most of the time
	if z.c < uint(len(z.buf)) ***REMOVED***
		b = z.buf[z.c]
		z.c++
		z.n++
		if z.trb ***REMOVED***
			z.tr = append(z.tr, b)
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // meaning z.c == len(z.buf) or greater ... so need to fill
		z.readbFill(z.b[:1], 0)
		b = z.b[0]
	***REMOVED***
	return
***REMOVED***

func (z *bufioDecReader) unreadn1() ***REMOVED***
	if z.c == 0 ***REMOVED***
		panic(errDecUnreadByteNothingToRead)
	***REMOVED***
	z.c--
	z.n--
	if z.trb ***REMOVED***
		z.tr = z.tr[:len(z.tr)-1]
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) readn(num uint8) (bs [rwNLen]byte) ***REMOVED***
	z.readb(bs[:num])
	// copy(bs[:], z.readx(uint(num)))
	return
***REMOVED***

func (z *bufioDecReader) readx(n uint) (bs []byte) ***REMOVED***
	if n == 0 ***REMOVED***
		// return
	***REMOVED*** else if z.c+n <= uint(len(z.buf)) ***REMOVED***
		bs = z.buf[z.c : z.c+n]
		z.n += n
		z.c += n
		if z.trb ***REMOVED***
			z.tr = append(z.tr, bs...)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		bs = make([]byte, n)
		// n no longer used - can reuse
		n = uint(copy(bs, z.buf[z.c:]))
		z.n += n
		z.c += n
		z.readbFill(bs, n)
	***REMOVED***
	return
***REMOVED***

func (z *bufioDecReader) skip(accept *bitset256) (token byte) ***REMOVED***
	i := z.c
LOOP:
	if i < uint(len(z.buf)) ***REMOVED***
		// inline z.skipLoopFn(i) and refactor, so cost is within inline budget
		token = z.buf[i]
		i++
		if accept.isset(token) ***REMOVED***
			goto LOOP
		***REMOVED***
		z.n += i - 2 - z.c
		if z.trb ***REMOVED***
			z.tr = append(z.tr, z.buf[z.c:i]...) // z.doTrack(i)
		***REMOVED***
		z.c = i
		return
	***REMOVED***
	return z.skipFill(accept)
***REMOVED***

func (z *bufioDecReader) skipFill(accept *bitset256) (token byte) ***REMOVED***
	z.n += uint(len(z.buf)) - z.c
	if z.trb ***REMOVED***
		z.tr = append(z.tr, z.buf[z.c:]...)
	***REMOVED***
	var i, n2 int
	var err error
	for ***REMOVED***
		z.c = 0
		z.buf = z.buf[0:cap(z.buf)]
		n2, err = z.r.Read(z.buf)
		if n2 == 0 && err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		z.buf = z.buf[:n2]
		for i, token = range z.buf ***REMOVED***
			// if !accept.isset(token) ***REMOVED***
			if accept.check(token) == 0 ***REMOVED***
				z.n += (uint(i) - z.c) - 1
				z.loopFn(uint(i + 1))
				return
			***REMOVED***
		***REMOVED***
		z.n += uint(n2)
		if z.trb ***REMOVED***
			z.tr = append(z.tr, z.buf...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) loopFn(i uint) ***REMOVED***
	if z.trb ***REMOVED***
		z.tr = append(z.tr, z.buf[z.c:i]...) // z.doTrack(i)
	***REMOVED***
	z.c = i
***REMOVED***

func (z *bufioDecReader) readTo(accept *bitset256) (out []byte) ***REMOVED***
	i := z.c
LOOP:
	if i < uint(len(z.buf)) ***REMOVED***
		// if !accept.isset(z.buf[i]) ***REMOVED***
		if accept.check(z.buf[i]) == 0 ***REMOVED***
			// inline readToLoopFn here (for performance)
			z.n += (i - z.c) - 1
			out = z.buf[z.c:i]
			if z.trb ***REMOVED***
				z.tr = append(z.tr, z.buf[z.c:i]...) // z.doTrack(i)
			***REMOVED***
			z.c = i
			return
		***REMOVED***
		i++
		goto LOOP
	***REMOVED***
	return z.readToFill(accept)
***REMOVED***

func (z *bufioDecReader) readToFill(accept *bitset256) []byte ***REMOVED***
	z.bufr = z.blist.check(z.bufr, 256)[:0]
	z.n += uint(len(z.buf)) - z.c
	z.bufr = append(z.bufr, z.buf[z.c:]...)
	if z.trb ***REMOVED***
		z.tr = append(z.tr, z.buf[z.c:]...)
	***REMOVED***
	var n2 int
	var err error
	for ***REMOVED***
		z.c = 0
		z.buf = z.buf[:cap(z.buf)]
		n2, err = z.r.Read(z.buf)
		if n2 == 0 && err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				return z.bufr // readTo should read until it matches or end is reached
			***REMOVED***
			panic(err)
		***REMOVED***
		z.buf = z.buf[:n2]
		for i, token := range z.buf ***REMOVED***
			// if !accept.isset(token) ***REMOVED***
			if accept.check(token) == 0 ***REMOVED***
				z.n += (uint(i) - z.c) - 1
				z.bufr = append(z.bufr, z.buf[z.c:i]...)
				z.loopFn(uint(i))
				return z.bufr
			***REMOVED***
		***REMOVED***
		z.bufr = append(z.bufr, z.buf...)
		z.n += uint(n2)
		if z.trb ***REMOVED***
			z.tr = append(z.tr, z.buf...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *bufioDecReader) readUntil(stop byte, includeLast bool) (out []byte) ***REMOVED***
	i := z.c
LOOP:
	if i < uint(len(z.buf)) ***REMOVED***
		if z.buf[i] == stop ***REMOVED***
			z.n += (i - z.c) - 1
			i++
			out = z.buf[z.c:i]
			if z.trb ***REMOVED***
				z.tr = append(z.tr, z.buf[z.c:i]...) // z.doTrack(i)
			***REMOVED***
			z.c = i
			goto FINISH
		***REMOVED***
		i++
		goto LOOP
	***REMOVED***
	out = z.readUntilFill(stop)
FINISH:
	if includeLast ***REMOVED***
		return
	***REMOVED***
	return out[:len(out)-1]
***REMOVED***

func (z *bufioDecReader) readUntilFill(stop byte) []byte ***REMOVED***
	z.bufr = z.blist.check(z.bufr, 256)[:0]
	z.n += uint(len(z.buf)) - z.c
	z.bufr = append(z.bufr, z.buf[z.c:]...)
	if z.trb ***REMOVED***
		z.tr = append(z.tr, z.buf[z.c:]...)
	***REMOVED***
	for ***REMOVED***
		z.c = 0
		z.buf = z.buf[0:cap(z.buf)]
		n1, err := z.r.Read(z.buf)
		if n1 == 0 && err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		n2 := uint(n1)
		z.buf = z.buf[:n2]
		for i, token := range z.buf ***REMOVED***
			if token == stop ***REMOVED***
				z.n += (uint(i) - z.c) - 1
				z.bufr = append(z.bufr, z.buf[z.c:i+1]...)
				z.loopFn(uint(i + 1))
				return z.bufr
			***REMOVED***
		***REMOVED***
		z.bufr = append(z.bufr, z.buf...)
		z.n += n2
		if z.trb ***REMOVED***
			z.tr = append(z.tr, z.buf...)
		***REMOVED***
	***REMOVED***
***REMOVED***

// ------------------------------------

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct ***REMOVED***
	b []byte // data
	c uint   // cursor
	t uint   // track start
	// a int    // available
***REMOVED***

func (z *bytesDecReader) reset(in []byte) ***REMOVED***
	z.b = in
	z.c = 0
	z.t = 0
***REMOVED***

func (z *bytesDecReader) numread() uint ***REMOVED***
	return z.c
***REMOVED***

func (z *bytesDecReader) last() byte ***REMOVED***
	return z.b[z.c-1]
***REMOVED***

func (z *bytesDecReader) unreadn1() ***REMOVED***
	if z.c == 0 || len(z.b) == 0 ***REMOVED***
		panic(errBytesDecReaderCannotUnread)
	***REMOVED***
	z.c--
***REMOVED***

func (z *bytesDecReader) readx(n uint) (bs []byte) ***REMOVED***
	// slicing from a non-constant start position is more expensive,
	// as more computation is required to decipher the pointer start position.
	// However, we do it only once, and it's better than reslicing both z.b and return value.

	z.c += n
	return z.b[z.c-n : z.c]
***REMOVED***

func (z *bytesDecReader) readb(bs []byte) ***REMOVED***
	copy(bs, z.readx(uint(len(bs))))
***REMOVED***

func (z *bytesDecReader) readn1() (v uint8) ***REMOVED***
	v = z.b[z.c]
	z.c++
	return
***REMOVED***

func (z *bytesDecReader) readn(num uint8) (bs [rwNLen]byte) ***REMOVED***
	// if z.c >= uint(len(z.b)) || z.c+uint(num) >= uint(len(z.b)) ***REMOVED***
	// 	panic(io.EOF)
	// ***REMOVED***

	// for bounds-check elimination, reslice z.b and ensure bs is within len
	// bb := z.b[z.c:][:num]
	bb := z.b[z.c : z.c+uint(num)]
	_ = bs[len(bb)-1]
	var i int
LOOP:
	if i < len(bb) ***REMOVED***
		bs[i] = bb[i]
		i++
		goto LOOP
	***REMOVED***

	z.c += uint(num)
	return
***REMOVED***

func (z *bytesDecReader) skip(accept *bitset256) (token byte) ***REMOVED***
	i := z.c
LOOP:
	// if i < uint(len(z.b)) ***REMOVED***
	token = z.b[i]
	i++
	if accept.isset(token) ***REMOVED***
		goto LOOP
	***REMOVED***
	z.c = i
	return
***REMOVED***

func (z *bytesDecReader) readTo(accept *bitset256) (out []byte) ***REMOVED***
	i := z.c
LOOP:
	if i < uint(len(z.b)) ***REMOVED***
		if accept.isset(z.b[i]) ***REMOVED***
			i++
			goto LOOP
		***REMOVED***
	***REMOVED***

	out = z.b[z.c:i]
	z.c = i
	return // z.b[c:i]
***REMOVED***

func (z *bytesDecReader) readUntil(stop byte, includeLast bool) (out []byte) ***REMOVED***
	i := z.c
LOOP:
	// if i < uint(len(z.b)) ***REMOVED***
	if z.b[i] == stop ***REMOVED***
		i++
		if includeLast ***REMOVED***
			out = z.b[z.c:i]
		***REMOVED*** else ***REMOVED***
			out = z.b[z.c : i-1]
		***REMOVED***
		// z.a -= (i - z.c)
		z.c = i
		return
	***REMOVED***
	i++
	goto LOOP
	// ***REMOVED***
	// panic(io.EOF)
***REMOVED***

func (z *bytesDecReader) track() ***REMOVED***
	z.t = z.c
***REMOVED***

func (z *bytesDecReader) stopTrack() (bs []byte) ***REMOVED***
	return z.b[z.t:z.c]
***REMOVED***

// --------------

type decRd struct ***REMOVED***
	mtr bool // is maptype a known type?
	str bool // is slicetype a known type?

	be   bool // is binary encoding
	js   bool // is json handle
	jsms bool // is json handle, and MapKeyAsString
	cbor bool // is cbor handle

	bytes bool // is bytes reader
	bufio bool // is this a bufioDecReader?

	rb bytesDecReader
	ri *ioDecReader
	bi *bufioDecReader
***REMOVED***

// numread, track and stopTrack are always inlined, as they just check int fields, etc.

// the if/else-if/else block is expensive to inline.
// Each node of this construct costs a lot and dominates the budget.
// Best to only do an if fast-path else block (so fast-path is inlined).
// This is irrespective of inlineExtraCallCost set in $GOROOT/src/cmd/compile/internal/gc/inl.go
//
// In decRd methods below, we delegate all IO functions into their own methods.
// This allows for the inlining of the common path when z.bytes=true.
// Go 1.12+ supports inlining methods with up to 1 inlined function (or 2 if no other constructs).
//
// However, up through Go 1.13, decRd's readXXX, skip and unreadXXX methods are not inlined.
// Consequently, there is no benefit to do the xxxIO methods for decRd at this time.
// Instead, we have a if/else-if/else block so that IO calls do not have to jump through
// a second unnecessary function call.
//
// If golang inlining gets better and bytesDecReader methods can be inlined,
// then we can revert to using these 2 functions so the bytesDecReader
// methods are inlined and the IO paths call out to a function.

func (z *decRd) numread() uint ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.numread()
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.numread()
	***REMOVED*** else ***REMOVED***
		return z.ri.numread()
	***REMOVED***
***REMOVED***
func (z *decRd) stopTrack() []byte ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.stopTrack()
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.stopTrack()
	***REMOVED*** else ***REMOVED***
		return z.ri.stopTrack()
	***REMOVED***
***REMOVED***

func (z *decRd) track() ***REMOVED***
	if z.bytes ***REMOVED***
		z.rb.track()
	***REMOVED*** else if z.bufio ***REMOVED***
		z.bi.track()
	***REMOVED*** else ***REMOVED***
		z.ri.track()
	***REMOVED***
***REMOVED***

func (z *decRd) unreadn1() ***REMOVED***
	if z.bytes ***REMOVED***
		z.rb.unreadn1()
	***REMOVED*** else if z.bufio ***REMOVED***
		z.bi.unreadn1()
	***REMOVED*** else ***REMOVED***
		z.ri.unreadn1() // not inlined
	***REMOVED***
***REMOVED***

func (z *decRd) readn(num uint8) [rwNLen]byte ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readn(num)
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.readn(num)
	***REMOVED*** else ***REMOVED***
		return z.ri.readn(num)
	***REMOVED***
***REMOVED***

func (z *decRd) readx(n uint) []byte ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readx(n)
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.readx(n)
	***REMOVED*** else ***REMOVED***
		return z.ri.readx(n)
	***REMOVED***
***REMOVED***

func (z *decRd) readb(s []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		z.rb.readb(s)
	***REMOVED*** else if z.bufio ***REMOVED***
		z.bi.readb(s)
	***REMOVED*** else ***REMOVED***
		z.ri.readb(s)
	***REMOVED***
***REMOVED***

func (z *decRd) readn1() uint8 ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readn1()
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.readn1()
	***REMOVED*** else ***REMOVED***
		return z.ri.readn1()
	***REMOVED***
***REMOVED***

func (z *decRd) skip(accept *bitset256) (token byte) ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.skip(accept)
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.skip(accept)
	***REMOVED*** else ***REMOVED***
		return z.ri.skip(accept)
	***REMOVED***
***REMOVED***

func (z *decRd) readTo(accept *bitset256) (out []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readTo(accept)
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.readTo(accept)
	***REMOVED*** else ***REMOVED***
		return z.ri.readTo(accept)
	***REMOVED***
***REMOVED***

func (z *decRd) readUntil(stop byte, includeLast bool) (out []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readUntil(stop, includeLast)
	***REMOVED*** else if z.bufio ***REMOVED***
		return z.bi.readUntil(stop, includeLast)
	***REMOVED*** else ***REMOVED***
		return z.ri.readUntil(stop, includeLast)
	***REMOVED***
***REMOVED***

/*
func (z *decRd) track() ***REMOVED***
	if z.bytes ***REMOVED***
		z.rb.track()
	***REMOVED*** else ***REMOVED***
		z.trackIO()
	***REMOVED***
***REMOVED***
func (z *decRd) trackIO() ***REMOVED***
	if z.bufio ***REMOVED***
		z.bi.track()
	***REMOVED*** else ***REMOVED***
		z.ri.track()
	***REMOVED***
***REMOVED***

func (z *decRd) unreadn1() ***REMOVED***
	if z.bytes ***REMOVED***
		z.rb.unreadn1()
	***REMOVED*** else ***REMOVED***
		z.unreadn1IO()
	***REMOVED***
***REMOVED***
func (z *decRd) unreadn1IO() ***REMOVED***
	if z.bufio ***REMOVED***
		z.bi.unreadn1()
	***REMOVED*** else ***REMOVED***
		z.ri.unreadn1()
	***REMOVED***
***REMOVED***

func (z *decRd) readn(num uint8) [rwNLen]byte ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readn(num)
	***REMOVED***
	return z.readnIO(num)
***REMOVED***
func (z *decRd) readnIO(num uint8) [rwNLen]byte ***REMOVED***
	if z.bufio ***REMOVED***
		return z.bi.readn(num)
	***REMOVED***
	return z.ri.readn(num)
***REMOVED***

func (z *decRd) readx(n uint) []byte ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readx(n)
	***REMOVED***
	return z.readxIO(n)
***REMOVED***
func (z *decRd) readxIO(n uint) []byte ***REMOVED***
	if z.bufio ***REMOVED***
		return z.bi.readx(n)
	***REMOVED***
	return z.ri.readx(n)
***REMOVED***

func (z *decRd) readb(s []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		z.rb.readb(s)
	***REMOVED*** else ***REMOVED***
		z.readbIO(s)
	***REMOVED***
***REMOVED***
func (z *decRd) readbIO(s []byte) ***REMOVED***
	if z.bufio ***REMOVED***
		z.bi.readb(s)
	***REMOVED*** else ***REMOVED***
		z.ri.readb(s)
	***REMOVED***
***REMOVED***

func (z *decRd) readn1() uint8 ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readn1()
	***REMOVED***
	return z.readn1IO()
***REMOVED***
func (z *decRd) readn1IO() uint8 ***REMOVED***
	if z.bufio ***REMOVED***
		return z.bi.readn1()
	***REMOVED***
	return z.ri.readn1()
***REMOVED***

func (z *decRd) skip(accept *bitset256) (token byte) ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.skip(accept)
	***REMOVED***
	return z.skipIO(accept)
***REMOVED***
func (z *decRd) skipIO(accept *bitset256) (token byte) ***REMOVED***
	if z.bufio ***REMOVED***
		return z.bi.skip(accept)
	***REMOVED***
	return z.ri.skip(accept)
***REMOVED***

func (z *decRd) readTo(accept *bitset256) (out []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readTo(accept)
	***REMOVED***
	return z.readToIO(accept)
***REMOVED***
func (z *decRd) readToIO(accept *bitset256) (out []byte) ***REMOVED***
	if z.bufio ***REMOVED***
		return z.bi.readTo(accept)
	***REMOVED***
	return z.ri.readTo(accept)
***REMOVED***

func (z *decRd) readUntil(stop byte, includeLast bool) (out []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		return z.rb.readUntil(stop, includeLast)
	***REMOVED***
	return z.readUntilIO(stop, includeLast)
***REMOVED***
func (z *decRd) readUntilIO(stop byte, includeLast bool) (out []byte) ***REMOVED***
	if z.bufio ***REMOVED***
		return z.bi.readUntil(stop, includeLast)
	***REMOVED***
	return z.ri.readUntil(stop, includeLast)
***REMOVED***
*/

var _ decReader = (*decRd)(nil)
