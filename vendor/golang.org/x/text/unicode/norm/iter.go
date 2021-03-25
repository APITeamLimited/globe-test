// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"fmt"
	"unicode/utf8"
)

// MaxSegmentSize is the maximum size of a byte buffer needed to consider any
// sequence of starter and non-starter runes for the purpose of normalization.
const MaxSegmentSize = maxByteBufferSize

// An Iter iterates over a string or byte slice, while normalizing it
// to a given Form.
type Iter struct ***REMOVED***
	rb     reorderBuffer
	buf    [maxByteBufferSize]byte
	info   Properties // first character saved from previous iteration
	next   iterFunc   // implementation of next depends on form
	asciiF iterFunc

	p        int    // current position in input source
	multiSeg []byte // remainder of multi-segment decomposition
***REMOVED***

type iterFunc func(*Iter) []byte

// Init initializes i to iterate over src after normalizing it to Form f.
func (i *Iter) Init(f Form, src []byte) ***REMOVED***
	i.p = 0
	if len(src) == 0 ***REMOVED***
		i.setDone()
		i.rb.nsrc = 0
		return
	***REMOVED***
	i.multiSeg = nil
	i.rb.init(f, src)
	i.next = i.rb.f.nextMain
	i.asciiF = nextASCIIBytes
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
***REMOVED***

// InitString initializes i to iterate over src after normalizing it to Form f.
func (i *Iter) InitString(f Form, src string) ***REMOVED***
	i.p = 0
	if len(src) == 0 ***REMOVED***
		i.setDone()
		i.rb.nsrc = 0
		return
	***REMOVED***
	i.multiSeg = nil
	i.rb.initString(f, src)
	i.next = i.rb.f.nextMain
	i.asciiF = nextASCIIString
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
***REMOVED***

// Seek sets the segment to be returned by the next call to Next to start
// at position p.  It is the responsibility of the caller to set p to the
// start of a segment.
func (i *Iter) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	var abs int64
	switch whence ***REMOVED***
	case 0:
		abs = offset
	case 1:
		abs = int64(i.p) + offset
	case 2:
		abs = int64(i.rb.nsrc) + offset
	default:
		return 0, fmt.Errorf("norm: invalid whence")
	***REMOVED***
	if abs < 0 ***REMOVED***
		return 0, fmt.Errorf("norm: negative position")
	***REMOVED***
	if int(abs) >= i.rb.nsrc ***REMOVED***
		i.setDone()
		return int64(i.p), nil
	***REMOVED***
	i.p = int(abs)
	i.multiSeg = nil
	i.next = i.rb.f.nextMain
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
	return abs, nil
***REMOVED***

// returnSlice returns a slice of the underlying input type as a byte slice.
// If the underlying is of type []byte, it will simply return a slice.
// If the underlying is of type string, it will copy the slice to the buffer
// and return that.
func (i *Iter) returnSlice(a, b int) []byte ***REMOVED***
	if i.rb.src.bytes == nil ***REMOVED***
		return i.buf[:copy(i.buf[:], i.rb.src.str[a:b])]
	***REMOVED***
	return i.rb.src.bytes[a:b]
***REMOVED***

// Pos returns the byte position at which the next call to Next will commence processing.
func (i *Iter) Pos() int ***REMOVED***
	return i.p
***REMOVED***

func (i *Iter) setDone() ***REMOVED***
	i.next = nextDone
	i.p = i.rb.nsrc
***REMOVED***

// Done returns true if there is no more input to process.
func (i *Iter) Done() bool ***REMOVED***
	return i.p >= i.rb.nsrc
***REMOVED***

// Next returns f(i.input[i.Pos():n]), where n is a boundary of i.input.
// For any input a and b for which f(a) == f(b), subsequent calls
// to Next will return the same segments.
// Modifying runes are grouped together with the preceding starter, if such a starter exists.
// Although not guaranteed, n will typically be the smallest possible n.
func (i *Iter) Next() []byte ***REMOVED***
	return i.next(i)
***REMOVED***

func nextASCIIBytes(i *Iter) []byte ***REMOVED***
	p := i.p + 1
	if p >= i.rb.nsrc ***REMOVED***
		p0 := i.p
		i.setDone()
		return i.rb.src.bytes[p0:p]
	***REMOVED***
	if i.rb.src.bytes[p] < utf8.RuneSelf ***REMOVED***
		p0 := i.p
		i.p = p
		return i.rb.src.bytes[p0:p]
	***REMOVED***
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.next = i.rb.f.nextMain
	return i.next(i)
***REMOVED***

func nextASCIIString(i *Iter) []byte ***REMOVED***
	p := i.p + 1
	if p >= i.rb.nsrc ***REMOVED***
		i.buf[0] = i.rb.src.str[i.p]
		i.setDone()
		return i.buf[:1]
	***REMOVED***
	if i.rb.src.str[p] < utf8.RuneSelf ***REMOVED***
		i.buf[0] = i.rb.src.str[i.p]
		i.p = p
		return i.buf[:1]
	***REMOVED***
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.next = i.rb.f.nextMain
	return i.next(i)
***REMOVED***

func nextHangul(i *Iter) []byte ***REMOVED***
	p := i.p
	next := p + hangulUTF8Size
	if next >= i.rb.nsrc ***REMOVED***
		i.setDone()
	***REMOVED*** else if i.rb.src.hangul(next) == 0 ***REMOVED***
		i.rb.ss.next(i.info)
		i.info = i.rb.f.info(i.rb.src, i.p)
		i.next = i.rb.f.nextMain
		return i.next(i)
	***REMOVED***
	i.p = next
	return i.buf[:decomposeHangul(i.buf[:], i.rb.src.hangul(p))]
***REMOVED***

func nextDone(i *Iter) []byte ***REMOVED***
	return nil
***REMOVED***

// nextMulti is used for iterating over multi-segment decompositions
// for decomposing normal forms.
func nextMulti(i *Iter) []byte ***REMOVED***
	j := 0
	d := i.multiSeg
	// skip first rune
	for j = 1; j < len(d) && !utf8.RuneStart(d[j]); j++ ***REMOVED***
	***REMOVED***
	for j < len(d) ***REMOVED***
		info := i.rb.f.info(input***REMOVED***bytes: d***REMOVED***, j)
		if info.BoundaryBefore() ***REMOVED***
			i.multiSeg = d[j:]
			return d[:j]
		***REMOVED***
		j += int(info.size)
	***REMOVED***
	// treat last segment as normal decomposition
	i.next = i.rb.f.nextMain
	return i.next(i)
***REMOVED***

// nextMultiNorm is used for iterating over multi-segment decompositions
// for composing normal forms.
func nextMultiNorm(i *Iter) []byte ***REMOVED***
	j := 0
	d := i.multiSeg
	for j < len(d) ***REMOVED***
		info := i.rb.f.info(input***REMOVED***bytes: d***REMOVED***, j)
		if info.BoundaryBefore() ***REMOVED***
			i.rb.compose()
			seg := i.buf[:i.rb.flushCopy(i.buf[:])]
			i.rb.insertUnsafe(input***REMOVED***bytes: d***REMOVED***, j, info)
			i.multiSeg = d[j+int(info.size):]
			return seg
		***REMOVED***
		i.rb.insertUnsafe(input***REMOVED***bytes: d***REMOVED***, j, info)
		j += int(info.size)
	***REMOVED***
	i.multiSeg = nil
	i.next = nextComposed
	return doNormComposed(i)
***REMOVED***

// nextDecomposed is the implementation of Next for forms NFD and NFKD.
func nextDecomposed(i *Iter) (next []byte) ***REMOVED***
	outp := 0
	inCopyStart, outCopyStart := i.p, 0
	for ***REMOVED***
		if sz := int(i.info.size); sz <= 1 ***REMOVED***
			i.rb.ss = 0
			p := i.p
			i.p++ // ASCII or illegal byte.  Either way, advance by 1.
			if i.p >= i.rb.nsrc ***REMOVED***
				i.setDone()
				return i.returnSlice(p, i.p)
			***REMOVED*** else if i.rb.src._byte(i.p) < utf8.RuneSelf ***REMOVED***
				i.next = i.asciiF
				return i.returnSlice(p, i.p)
			***REMOVED***
			outp++
		***REMOVED*** else if d := i.info.Decomposition(); d != nil ***REMOVED***
			// Note: If leading CCC != 0, then len(d) == 2 and last is also non-zero.
			// Case 1: there is a leftover to copy.  In this case the decomposition
			// must begin with a modifier and should always be appended.
			// Case 2: no leftover. Simply return d if followed by a ccc == 0 value.
			p := outp + len(d)
			if outp > 0 ***REMOVED***
				i.rb.src.copySlice(i.buf[outCopyStart:], inCopyStart, i.p)
				// TODO: this condition should not be possible, but we leave it
				// in for defensive purposes.
				if p > len(i.buf) ***REMOVED***
					return i.buf[:outp]
				***REMOVED***
			***REMOVED*** else if i.info.multiSegment() ***REMOVED***
				// outp must be 0 as multi-segment decompositions always
				// start a new segment.
				if i.multiSeg == nil ***REMOVED***
					i.multiSeg = d
					i.next = nextMulti
					return nextMulti(i)
				***REMOVED***
				// We are in the last segment.  Treat as normal decomposition.
				d = i.multiSeg
				i.multiSeg = nil
				p = len(d)
			***REMOVED***
			prevCC := i.info.tccc
			if i.p += sz; i.p >= i.rb.nsrc ***REMOVED***
				i.setDone()
				i.info = Properties***REMOVED******REMOVED*** // Force BoundaryBefore to succeed.
			***REMOVED*** else ***REMOVED***
				i.info = i.rb.f.info(i.rb.src, i.p)
			***REMOVED***
			switch i.rb.ss.next(i.info) ***REMOVED***
			case ssOverflow:
				i.next = nextCGJDecompose
				fallthrough
			case ssStarter:
				if outp > 0 ***REMOVED***
					copy(i.buf[outp:], d)
					return i.buf[:p]
				***REMOVED***
				return d
			***REMOVED***
			copy(i.buf[outp:], d)
			outp = p
			inCopyStart, outCopyStart = i.p, outp
			if i.info.ccc < prevCC ***REMOVED***
				goto doNorm
			***REMOVED***
			continue
		***REMOVED*** else if r := i.rb.src.hangul(i.p); r != 0 ***REMOVED***
			outp = decomposeHangul(i.buf[:], r)
			i.p += hangulUTF8Size
			inCopyStart, outCopyStart = i.p, outp
			if i.p >= i.rb.nsrc ***REMOVED***
				i.setDone()
				break
			***REMOVED*** else if i.rb.src.hangul(i.p) != 0 ***REMOVED***
				i.next = nextHangul
				return i.buf[:outp]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			p := outp + sz
			if p > len(i.buf) ***REMOVED***
				break
			***REMOVED***
			outp = p
			i.p += sz
		***REMOVED***
		if i.p >= i.rb.nsrc ***REMOVED***
			i.setDone()
			break
		***REMOVED***
		prevCC := i.info.tccc
		i.info = i.rb.f.info(i.rb.src, i.p)
		if v := i.rb.ss.next(i.info); v == ssStarter ***REMOVED***
			break
		***REMOVED*** else if v == ssOverflow ***REMOVED***
			i.next = nextCGJDecompose
			break
		***REMOVED***
		if i.info.ccc < prevCC ***REMOVED***
			goto doNorm
		***REMOVED***
	***REMOVED***
	if outCopyStart == 0 ***REMOVED***
		return i.returnSlice(inCopyStart, i.p)
	***REMOVED*** else if inCopyStart < i.p ***REMOVED***
		i.rb.src.copySlice(i.buf[outCopyStart:], inCopyStart, i.p)
	***REMOVED***
	return i.buf[:outp]
doNorm:
	// Insert what we have decomposed so far in the reorderBuffer.
	// As we will only reorder, there will always be enough room.
	i.rb.src.copySlice(i.buf[outCopyStart:], inCopyStart, i.p)
	i.rb.insertDecomposed(i.buf[0:outp])
	return doNormDecomposed(i)
***REMOVED***

func doNormDecomposed(i *Iter) []byte ***REMOVED***
	for ***REMOVED***
		i.rb.insertUnsafe(i.rb.src, i.p, i.info)
		if i.p += int(i.info.size); i.p >= i.rb.nsrc ***REMOVED***
			i.setDone()
			break
		***REMOVED***
		i.info = i.rb.f.info(i.rb.src, i.p)
		if i.info.ccc == 0 ***REMOVED***
			break
		***REMOVED***
		if s := i.rb.ss.next(i.info); s == ssOverflow ***REMOVED***
			i.next = nextCGJDecompose
			break
		***REMOVED***
	***REMOVED***
	// new segment or too many combining characters: exit normalization
	return i.buf[:i.rb.flushCopy(i.buf[:])]
***REMOVED***

func nextCGJDecompose(i *Iter) []byte ***REMOVED***
	i.rb.ss = 0
	i.rb.insertCGJ()
	i.next = nextDecomposed
	i.rb.ss.first(i.info)
	buf := doNormDecomposed(i)
	return buf
***REMOVED***

// nextComposed is the implementation of Next for forms NFC and NFKC.
func nextComposed(i *Iter) []byte ***REMOVED***
	outp, startp := 0, i.p
	var prevCC uint8
	for ***REMOVED***
		if !i.info.isYesC() ***REMOVED***
			goto doNorm
		***REMOVED***
		prevCC = i.info.tccc
		sz := int(i.info.size)
		if sz == 0 ***REMOVED***
			sz = 1 // illegal rune: copy byte-by-byte
		***REMOVED***
		p := outp + sz
		if p > len(i.buf) ***REMOVED***
			break
		***REMOVED***
		outp = p
		i.p += sz
		if i.p >= i.rb.nsrc ***REMOVED***
			i.setDone()
			break
		***REMOVED*** else if i.rb.src._byte(i.p) < utf8.RuneSelf ***REMOVED***
			i.rb.ss = 0
			i.next = i.asciiF
			break
		***REMOVED***
		i.info = i.rb.f.info(i.rb.src, i.p)
		if v := i.rb.ss.next(i.info); v == ssStarter ***REMOVED***
			break
		***REMOVED*** else if v == ssOverflow ***REMOVED***
			i.next = nextCGJCompose
			break
		***REMOVED***
		if i.info.ccc < prevCC ***REMOVED***
			goto doNorm
		***REMOVED***
	***REMOVED***
	return i.returnSlice(startp, i.p)
doNorm:
	// reset to start position
	i.p = startp
	i.info = i.rb.f.info(i.rb.src, i.p)
	i.rb.ss.first(i.info)
	if i.info.multiSegment() ***REMOVED***
		d := i.info.Decomposition()
		info := i.rb.f.info(input***REMOVED***bytes: d***REMOVED***, 0)
		i.rb.insertUnsafe(input***REMOVED***bytes: d***REMOVED***, 0, info)
		i.multiSeg = d[int(info.size):]
		i.next = nextMultiNorm
		return nextMultiNorm(i)
	***REMOVED***
	i.rb.ss.first(i.info)
	i.rb.insertUnsafe(i.rb.src, i.p, i.info)
	return doNormComposed(i)
***REMOVED***

func doNormComposed(i *Iter) []byte ***REMOVED***
	// First rune should already be inserted.
	for ***REMOVED***
		if i.p += int(i.info.size); i.p >= i.rb.nsrc ***REMOVED***
			i.setDone()
			break
		***REMOVED***
		i.info = i.rb.f.info(i.rb.src, i.p)
		if s := i.rb.ss.next(i.info); s == ssStarter ***REMOVED***
			break
		***REMOVED*** else if s == ssOverflow ***REMOVED***
			i.next = nextCGJCompose
			break
		***REMOVED***
		i.rb.insertUnsafe(i.rb.src, i.p, i.info)
	***REMOVED***
	i.rb.compose()
	seg := i.buf[:i.rb.flushCopy(i.buf[:])]
	return seg
***REMOVED***

func nextCGJCompose(i *Iter) []byte ***REMOVED***
	i.rb.ss = 0 // instead of first
	i.rb.insertCGJ()
	i.next = nextComposed
	// Note that we treat any rune with nLeadingNonStarters > 0 as a non-starter,
	// even if they are not. This is particularly dubious for U+FF9E and UFF9A.
	// If we ever change that, insert a check here.
	i.rb.ss.first(i.info)
	i.rb.insertUnsafe(i.rb.src, i.p, i.info)
	return doNormComposed(i)
***REMOVED***
