// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note: the file data_test.go that is generated should not be checked in.
//go:generate go run maketables.go triegen.go
//go:generate go test -tags test

// Package norm contains types and functions for normalizing Unicode strings.
package norm // import "golang.org/x/text/unicode/norm"

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// A Form denotes a canonical representation of Unicode code points.
// The Unicode-defined normalization and equivalence forms are:
//
//   NFC   Unicode Normalization Form C
//   NFD   Unicode Normalization Form D
//   NFKC  Unicode Normalization Form KC
//   NFKD  Unicode Normalization Form KD
//
// For a Form f, this documentation uses the notation f(x) to mean
// the bytes or string x converted to the given form.
// A position n in x is called a boundary if conversion to the form can
// proceed independently on both sides:
//   f(x) == append(f(x[0:n]), f(x[n:])...)
//
// References: https://unicode.org/reports/tr15/ and
// https://unicode.org/notes/tn5/.
type Form int

const (
	NFC Form = iota
	NFD
	NFKC
	NFKD
)

// Bytes returns f(b). May return b if f(b) = b.
func (f Form) Bytes(b []byte) []byte ***REMOVED***
	src := inputBytes(b)
	ft := formTable[f]
	n, ok := ft.quickSpan(src, 0, len(b), true)
	if ok ***REMOVED***
		return b
	***REMOVED***
	out := make([]byte, n, len(b))
	copy(out, b[0:n])
	rb := reorderBuffer***REMOVED***f: *ft, src: src, nsrc: len(b), out: out, flushF: appendFlush***REMOVED***
	return doAppendInner(&rb, n)
***REMOVED***

// String returns f(s).
func (f Form) String(s string) string ***REMOVED***
	src := inputString(s)
	ft := formTable[f]
	n, ok := ft.quickSpan(src, 0, len(s), true)
	if ok ***REMOVED***
		return s
	***REMOVED***
	out := make([]byte, n, len(s))
	copy(out, s[0:n])
	rb := reorderBuffer***REMOVED***f: *ft, src: src, nsrc: len(s), out: out, flushF: appendFlush***REMOVED***
	return string(doAppendInner(&rb, n))
***REMOVED***

// IsNormal returns true if b == f(b).
func (f Form) IsNormal(b []byte) bool ***REMOVED***
	src := inputBytes(b)
	ft := formTable[f]
	bp, ok := ft.quickSpan(src, 0, len(b), true)
	if ok ***REMOVED***
		return true
	***REMOVED***
	rb := reorderBuffer***REMOVED***f: *ft, src: src, nsrc: len(b)***REMOVED***
	rb.setFlusher(nil, cmpNormalBytes)
	for bp < len(b) ***REMOVED***
		rb.out = b[bp:]
		if bp = decomposeSegment(&rb, bp, true); bp < 0 ***REMOVED***
			return false
		***REMOVED***
		bp, _ = rb.f.quickSpan(rb.src, bp, len(b), true)
	***REMOVED***
	return true
***REMOVED***

func cmpNormalBytes(rb *reorderBuffer) bool ***REMOVED***
	b := rb.out
	for i := 0; i < rb.nrune; i++ ***REMOVED***
		info := rb.rune[i]
		if int(info.size) > len(b) ***REMOVED***
			return false
		***REMOVED***
		p := info.pos
		pe := p + info.size
		for ; p < pe; p++ ***REMOVED***
			if b[0] != rb.byte[p] ***REMOVED***
				return false
			***REMOVED***
			b = b[1:]
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// IsNormalString returns true if s == f(s).
func (f Form) IsNormalString(s string) bool ***REMOVED***
	src := inputString(s)
	ft := formTable[f]
	bp, ok := ft.quickSpan(src, 0, len(s), true)
	if ok ***REMOVED***
		return true
	***REMOVED***
	rb := reorderBuffer***REMOVED***f: *ft, src: src, nsrc: len(s)***REMOVED***
	rb.setFlusher(nil, func(rb *reorderBuffer) bool ***REMOVED***
		for i := 0; i < rb.nrune; i++ ***REMOVED***
			info := rb.rune[i]
			if bp+int(info.size) > len(s) ***REMOVED***
				return false
			***REMOVED***
			p := info.pos
			pe := p + info.size
			for ; p < pe; p++ ***REMOVED***
				if s[bp] != rb.byte[p] ***REMOVED***
					return false
				***REMOVED***
				bp++
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***)
	for bp < len(s) ***REMOVED***
		if bp = decomposeSegment(&rb, bp, true); bp < 0 ***REMOVED***
			return false
		***REMOVED***
		bp, _ = rb.f.quickSpan(rb.src, bp, len(s), true)
	***REMOVED***
	return true
***REMOVED***

// patchTail fixes a case where a rune may be incorrectly normalized
// if it is followed by illegal continuation bytes. It returns the
// patched buffer and whether the decomposition is still in progress.
func patchTail(rb *reorderBuffer) bool ***REMOVED***
	info, p := lastRuneStart(&rb.f, rb.out)
	if p == -1 || info.size == 0 ***REMOVED***
		return true
	***REMOVED***
	end := p + int(info.size)
	extra := len(rb.out) - end
	if extra > 0 ***REMOVED***
		// Potentially allocating memory. However, this only
		// happens with ill-formed UTF-8.
		x := make([]byte, 0)
		x = append(x, rb.out[len(rb.out)-extra:]...)
		rb.out = rb.out[:end]
		decomposeToLastBoundary(rb)
		rb.doFlush()
		rb.out = append(rb.out, x...)
		return false
	***REMOVED***
	buf := rb.out[p:]
	rb.out = rb.out[:p]
	decomposeToLastBoundary(rb)
	if s := rb.ss.next(info); s == ssStarter ***REMOVED***
		rb.doFlush()
		rb.ss.first(info)
	***REMOVED*** else if s == ssOverflow ***REMOVED***
		rb.doFlush()
		rb.insertCGJ()
		rb.ss = 0
	***REMOVED***
	rb.insertUnsafe(inputBytes(buf), 0, info)
	return true
***REMOVED***

func appendQuick(rb *reorderBuffer, i int) int ***REMOVED***
	if rb.nsrc == i ***REMOVED***
		return i
	***REMOVED***
	end, _ := rb.f.quickSpan(rb.src, i, rb.nsrc, true)
	rb.out = rb.src.appendSlice(rb.out, i, end)
	return end
***REMOVED***

// Append returns f(append(out, b...)).
// The buffer out must be nil, empty, or equal to f(out).
func (f Form) Append(out []byte, src ...byte) []byte ***REMOVED***
	return f.doAppend(out, inputBytes(src), len(src))
***REMOVED***

func (f Form) doAppend(out []byte, src input, n int) []byte ***REMOVED***
	if n == 0 ***REMOVED***
		return out
	***REMOVED***
	ft := formTable[f]
	// Attempt to do a quickSpan first so we can avoid initializing the reorderBuffer.
	if len(out) == 0 ***REMOVED***
		p, _ := ft.quickSpan(src, 0, n, true)
		out = src.appendSlice(out, 0, p)
		if p == n ***REMOVED***
			return out
		***REMOVED***
		rb := reorderBuffer***REMOVED***f: *ft, src: src, nsrc: n, out: out, flushF: appendFlush***REMOVED***
		return doAppendInner(&rb, p)
	***REMOVED***
	rb := reorderBuffer***REMOVED***f: *ft, src: src, nsrc: n***REMOVED***
	return doAppend(&rb, out, 0)
***REMOVED***

func doAppend(rb *reorderBuffer, out []byte, p int) []byte ***REMOVED***
	rb.setFlusher(out, appendFlush)
	src, n := rb.src, rb.nsrc
	doMerge := len(out) > 0
	if q := src.skipContinuationBytes(p); q > p ***REMOVED***
		// Move leading non-starters to destination.
		rb.out = src.appendSlice(rb.out, p, q)
		p = q
		doMerge = patchTail(rb)
	***REMOVED***
	fd := &rb.f
	if doMerge ***REMOVED***
		var info Properties
		if p < n ***REMOVED***
			info = fd.info(src, p)
			if !info.BoundaryBefore() || info.nLeadingNonStarters() > 0 ***REMOVED***
				if p == 0 ***REMOVED***
					decomposeToLastBoundary(rb)
				***REMOVED***
				p = decomposeSegment(rb, p, true)
			***REMOVED***
		***REMOVED***
		if info.size == 0 ***REMOVED***
			rb.doFlush()
			// Append incomplete UTF-8 encoding.
			return src.appendSlice(rb.out, p, n)
		***REMOVED***
		if rb.nrune > 0 ***REMOVED***
			return doAppendInner(rb, p)
		***REMOVED***
	***REMOVED***
	p = appendQuick(rb, p)
	return doAppendInner(rb, p)
***REMOVED***

func doAppendInner(rb *reorderBuffer, p int) []byte ***REMOVED***
	for n := rb.nsrc; p < n; ***REMOVED***
		p = decomposeSegment(rb, p, true)
		p = appendQuick(rb, p)
	***REMOVED***
	return rb.out
***REMOVED***

// AppendString returns f(append(out, []byte(s))).
// The buffer out must be nil, empty, or equal to f(out).
func (f Form) AppendString(out []byte, src string) []byte ***REMOVED***
	return f.doAppend(out, inputString(src), len(src))
***REMOVED***

// QuickSpan returns a boundary n such that b[0:n] == f(b[0:n]).
// It is not guaranteed to return the largest such n.
func (f Form) QuickSpan(b []byte) int ***REMOVED***
	n, _ := formTable[f].quickSpan(inputBytes(b), 0, len(b), true)
	return n
***REMOVED***

// Span implements transform.SpanningTransformer. It returns a boundary n such
// that b[0:n] == f(b[0:n]). It is not guaranteed to return the largest such n.
func (f Form) Span(b []byte, atEOF bool) (n int, err error) ***REMOVED***
	n, ok := formTable[f].quickSpan(inputBytes(b), 0, len(b), atEOF)
	if n < len(b) ***REMOVED***
		if !ok ***REMOVED***
			err = transform.ErrEndOfSpan
		***REMOVED*** else ***REMOVED***
			err = transform.ErrShortSrc
		***REMOVED***
	***REMOVED***
	return n, err
***REMOVED***

// SpanString returns a boundary n such that s[0:n] == f(s[0:n]).
// It is not guaranteed to return the largest such n.
func (f Form) SpanString(s string, atEOF bool) (n int, err error) ***REMOVED***
	n, ok := formTable[f].quickSpan(inputString(s), 0, len(s), atEOF)
	if n < len(s) ***REMOVED***
		if !ok ***REMOVED***
			err = transform.ErrEndOfSpan
		***REMOVED*** else ***REMOVED***
			err = transform.ErrShortSrc
		***REMOVED***
	***REMOVED***
	return n, err
***REMOVED***

// quickSpan returns a boundary n such that src[0:n] == f(src[0:n]) and
// whether any non-normalized parts were found. If atEOF is false, n will
// not point past the last segment if this segment might be become
// non-normalized by appending other runes.
func (f *formInfo) quickSpan(src input, i, end int, atEOF bool) (n int, ok bool) ***REMOVED***
	var lastCC uint8
	ss := streamSafe(0)
	lastSegStart := i
	for n = end; i < n; ***REMOVED***
		if j := src.skipASCII(i, n); i != j ***REMOVED***
			i = j
			lastSegStart = i - 1
			lastCC = 0
			ss = 0
			continue
		***REMOVED***
		info := f.info(src, i)
		if info.size == 0 ***REMOVED***
			if atEOF ***REMOVED***
				// include incomplete runes
				return n, true
			***REMOVED***
			return lastSegStart, true
		***REMOVED***
		// This block needs to be before the next, because it is possible to
		// have an overflow for runes that are starters (e.g. with U+FF9E).
		switch ss.next(info) ***REMOVED***
		case ssStarter:
			lastSegStart = i
		case ssOverflow:
			return lastSegStart, false
		case ssSuccess:
			if lastCC > info.ccc ***REMOVED***
				return lastSegStart, false
			***REMOVED***
		***REMOVED***
		if f.composing ***REMOVED***
			if !info.isYesC() ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if !info.isYesD() ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		lastCC = info.ccc
		i += int(info.size)
	***REMOVED***
	if i == n ***REMOVED***
		if !atEOF ***REMOVED***
			n = lastSegStart
		***REMOVED***
		return n, true
	***REMOVED***
	return lastSegStart, false
***REMOVED***

// QuickSpanString returns a boundary n such that s[0:n] == f(s[0:n]).
// It is not guaranteed to return the largest such n.
func (f Form) QuickSpanString(s string) int ***REMOVED***
	n, _ := formTable[f].quickSpan(inputString(s), 0, len(s), true)
	return n
***REMOVED***

// FirstBoundary returns the position i of the first boundary in b
// or -1 if b contains no boundary.
func (f Form) FirstBoundary(b []byte) int ***REMOVED***
	return f.firstBoundary(inputBytes(b), len(b))
***REMOVED***

func (f Form) firstBoundary(src input, nsrc int) int ***REMOVED***
	i := src.skipContinuationBytes(0)
	if i >= nsrc ***REMOVED***
		return -1
	***REMOVED***
	fd := formTable[f]
	ss := streamSafe(0)
	// We should call ss.first here, but we can't as the first rune is
	// skipped already. This means FirstBoundary can't really determine
	// CGJ insertion points correctly. Luckily it doesn't have to.
	for ***REMOVED***
		info := fd.info(src, i)
		if info.size == 0 ***REMOVED***
			return -1
		***REMOVED***
		if s := ss.next(info); s != ssSuccess ***REMOVED***
			return i
		***REMOVED***
		i += int(info.size)
		if i >= nsrc ***REMOVED***
			if !info.BoundaryAfter() && !ss.isMax() ***REMOVED***
				return -1
			***REMOVED***
			return nsrc
		***REMOVED***
	***REMOVED***
***REMOVED***

// FirstBoundaryInString returns the position i of the first boundary in s
// or -1 if s contains no boundary.
func (f Form) FirstBoundaryInString(s string) int ***REMOVED***
	return f.firstBoundary(inputString(s), len(s))
***REMOVED***

// NextBoundary reports the index of the boundary between the first and next
// segment in b or -1 if atEOF is false and there are not enough bytes to
// determine this boundary.
func (f Form) NextBoundary(b []byte, atEOF bool) int ***REMOVED***
	return f.nextBoundary(inputBytes(b), len(b), atEOF)
***REMOVED***

// NextBoundaryInString reports the index of the boundary between the first and
// next segment in b or -1 if atEOF is false and there are not enough bytes to
// determine this boundary.
func (f Form) NextBoundaryInString(s string, atEOF bool) int ***REMOVED***
	return f.nextBoundary(inputString(s), len(s), atEOF)
***REMOVED***

func (f Form) nextBoundary(src input, nsrc int, atEOF bool) int ***REMOVED***
	if nsrc == 0 ***REMOVED***
		if atEOF ***REMOVED***
			return 0
		***REMOVED***
		return -1
	***REMOVED***
	fd := formTable[f]
	info := fd.info(src, 0)
	if info.size == 0 ***REMOVED***
		if atEOF ***REMOVED***
			return 1
		***REMOVED***
		return -1
	***REMOVED***
	ss := streamSafe(0)
	ss.first(info)

	for i := int(info.size); i < nsrc; i += int(info.size) ***REMOVED***
		info = fd.info(src, i)
		if info.size == 0 ***REMOVED***
			if atEOF ***REMOVED***
				return i
			***REMOVED***
			return -1
		***REMOVED***
		// TODO: Using streamSafe to determine the boundary isn't the same as
		// using BoundaryBefore. Determine which should be used.
		if s := ss.next(info); s != ssSuccess ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	if !atEOF && !info.BoundaryAfter() && !ss.isMax() ***REMOVED***
		return -1
	***REMOVED***
	return nsrc
***REMOVED***

// LastBoundary returns the position i of the last boundary in b
// or -1 if b contains no boundary.
func (f Form) LastBoundary(b []byte) int ***REMOVED***
	return lastBoundary(formTable[f], b)
***REMOVED***

func lastBoundary(fd *formInfo, b []byte) int ***REMOVED***
	i := len(b)
	info, p := lastRuneStart(fd, b)
	if p == -1 ***REMOVED***
		return -1
	***REMOVED***
	if info.size == 0 ***REMOVED*** // ends with incomplete rune
		if p == 0 ***REMOVED*** // starts with incomplete rune
			return -1
		***REMOVED***
		i = p
		info, p = lastRuneStart(fd, b[:i])
		if p == -1 ***REMOVED*** // incomplete UTF-8 encoding or non-starter bytes without a starter
			return i
		***REMOVED***
	***REMOVED***
	if p+int(info.size) != i ***REMOVED*** // trailing non-starter bytes: illegal UTF-8
		return i
	***REMOVED***
	if info.BoundaryAfter() ***REMOVED***
		return i
	***REMOVED***
	ss := streamSafe(0)
	v := ss.backwards(info)
	for i = p; i >= 0 && v != ssStarter; i = p ***REMOVED***
		info, p = lastRuneStart(fd, b[:i])
		if v = ss.backwards(info); v == ssOverflow ***REMOVED***
			break
		***REMOVED***
		if p+int(info.size) != i ***REMOVED***
			if p == -1 ***REMOVED*** // no boundary found
				return -1
			***REMOVED***
			return i // boundary after an illegal UTF-8 encoding
		***REMOVED***
	***REMOVED***
	return i
***REMOVED***

// decomposeSegment scans the first segment in src into rb. It inserts 0x034f
// (Grapheme Joiner) when it encounters a sequence of more than 30 non-starters
// and returns the number of bytes consumed from src or iShortDst or iShortSrc.
func decomposeSegment(rb *reorderBuffer, sp int, atEOF bool) int ***REMOVED***
	// Force one character to be consumed.
	info := rb.f.info(rb.src, sp)
	if info.size == 0 ***REMOVED***
		return 0
	***REMOVED***
	if s := rb.ss.next(info); s == ssStarter ***REMOVED***
		// TODO: this could be removed if we don't support merging.
		if rb.nrune > 0 ***REMOVED***
			goto end
		***REMOVED***
	***REMOVED*** else if s == ssOverflow ***REMOVED***
		rb.insertCGJ()
		goto end
	***REMOVED***
	if err := rb.insertFlush(rb.src, sp, info); err != iSuccess ***REMOVED***
		return int(err)
	***REMOVED***
	for ***REMOVED***
		sp += int(info.size)
		if sp >= rb.nsrc ***REMOVED***
			if !atEOF && !info.BoundaryAfter() ***REMOVED***
				return int(iShortSrc)
			***REMOVED***
			break
		***REMOVED***
		info = rb.f.info(rb.src, sp)
		if info.size == 0 ***REMOVED***
			if !atEOF ***REMOVED***
				return int(iShortSrc)
			***REMOVED***
			break
		***REMOVED***
		if s := rb.ss.next(info); s == ssStarter ***REMOVED***
			break
		***REMOVED*** else if s == ssOverflow ***REMOVED***
			rb.insertCGJ()
			break
		***REMOVED***
		if err := rb.insertFlush(rb.src, sp, info); err != iSuccess ***REMOVED***
			return int(err)
		***REMOVED***
	***REMOVED***
end:
	if !rb.doFlush() ***REMOVED***
		return int(iShortDst)
	***REMOVED***
	return sp
***REMOVED***

// lastRuneStart returns the runeInfo and position of the last
// rune in buf or the zero runeInfo and -1 if no rune was found.
func lastRuneStart(fd *formInfo, buf []byte) (Properties, int) ***REMOVED***
	p := len(buf) - 1
	for ; p >= 0 && !utf8.RuneStart(buf[p]); p-- ***REMOVED***
	***REMOVED***
	if p < 0 ***REMOVED***
		return Properties***REMOVED******REMOVED***, -1
	***REMOVED***
	return fd.info(inputBytes(buf), p), p
***REMOVED***

// decomposeToLastBoundary finds an open segment at the end of the buffer
// and scans it into rb. Returns the buffer minus the last segment.
func decomposeToLastBoundary(rb *reorderBuffer) ***REMOVED***
	fd := &rb.f
	info, i := lastRuneStart(fd, rb.out)
	if int(info.size) != len(rb.out)-i ***REMOVED***
		// illegal trailing continuation bytes
		return
	***REMOVED***
	if info.BoundaryAfter() ***REMOVED***
		return
	***REMOVED***
	var add [maxNonStarters + 1]Properties // stores runeInfo in reverse order
	padd := 0
	ss := streamSafe(0)
	p := len(rb.out)
	for ***REMOVED***
		add[padd] = info
		v := ss.backwards(info)
		if v == ssOverflow ***REMOVED***
			// Note that if we have an overflow, it the string we are appending to
			// is not correctly normalized. In this case the behavior is undefined.
			break
		***REMOVED***
		padd++
		p -= int(info.size)
		if v == ssStarter || p < 0 ***REMOVED***
			break
		***REMOVED***
		info, i = lastRuneStart(fd, rb.out[:p])
		if int(info.size) != p-i ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	rb.ss = ss
	// Copy bytes for insertion as we may need to overwrite rb.out.
	var buf [maxBufferSize * utf8.UTFMax]byte
	cp := buf[:copy(buf[:], rb.out[p:])]
	rb.out = rb.out[:p]
	for padd--; padd >= 0; padd-- ***REMOVED***
		info = add[padd]
		rb.insertUnsafe(inputBytes(cp), 0, info)
		cp = cp[info.size:]
	***REMOVED***
***REMOVED***
