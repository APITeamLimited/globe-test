// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import "unicode/utf8"

const (
	maxNonStarters = 30
	// The maximum number of characters needed for a buffer is
	// maxNonStarters + 1 for the starter + 1 for the GCJ
	maxBufferSize    = maxNonStarters + 2
	maxNFCExpansion  = 3  // NFC(0x1D160)
	maxNFKCExpansion = 18 // NFKC(0xFDFA)

	maxByteBufferSize = utf8.UTFMax * maxBufferSize // 128
)

// ssState is used for reporting the segment state after inserting a rune.
// It is returned by streamSafe.next.
type ssState int

const (
	// Indicates a rune was successfully added to the segment.
	ssSuccess ssState = iota
	// Indicates a rune starts a new segment and should not be added.
	ssStarter
	// Indicates a rune caused a segment overflow and a CGJ should be inserted.
	ssOverflow
)

// streamSafe implements the policy of when a CGJ should be inserted.
type streamSafe uint8

// first inserts the first rune of a segment. It is a faster version of next if
// it is known p represents the first rune in a segment.
func (ss *streamSafe) first(p Properties) ***REMOVED***
	*ss = streamSafe(p.nTrailingNonStarters())
***REMOVED***

// insert returns a ssState value to indicate whether a rune represented by p
// can be inserted.
func (ss *streamSafe) next(p Properties) ssState ***REMOVED***
	if *ss > maxNonStarters ***REMOVED***
		panic("streamSafe was not reset")
	***REMOVED***
	n := p.nLeadingNonStarters()
	if *ss += streamSafe(n); *ss > maxNonStarters ***REMOVED***
		*ss = 0
		return ssOverflow
	***REMOVED***
	// The Stream-Safe Text Processing prescribes that the counting can stop
	// as soon as a starter is encountered. However, there are some starters,
	// like Jamo V and T, that can combine with other runes, leaving their
	// successive non-starters appended to the previous, possibly causing an
	// overflow. We will therefore consider any rune with a non-zero nLead to
	// be a non-starter. Note that it always hold that if nLead > 0 then
	// nLead == nTrail.
	if n == 0 ***REMOVED***
		*ss = streamSafe(p.nTrailingNonStarters())
		return ssStarter
	***REMOVED***
	return ssSuccess
***REMOVED***

// backwards is used for checking for overflow and segment starts
// when traversing a string backwards. Users do not need to call first
// for the first rune. The state of the streamSafe retains the count of
// the non-starters loaded.
func (ss *streamSafe) backwards(p Properties) ssState ***REMOVED***
	if *ss > maxNonStarters ***REMOVED***
		panic("streamSafe was not reset")
	***REMOVED***
	c := *ss + streamSafe(p.nTrailingNonStarters())
	if c > maxNonStarters ***REMOVED***
		return ssOverflow
	***REMOVED***
	*ss = c
	if p.nLeadingNonStarters() == 0 ***REMOVED***
		return ssStarter
	***REMOVED***
	return ssSuccess
***REMOVED***

func (ss streamSafe) isMax() bool ***REMOVED***
	return ss == maxNonStarters
***REMOVED***

// GraphemeJoiner is inserted after maxNonStarters non-starter runes.
const GraphemeJoiner = "\u034F"

// reorderBuffer is used to normalize a single segment.  Characters inserted with
// insert are decomposed and reordered based on CCC. The compose method can
// be used to recombine characters.  Note that the byte buffer does not hold
// the UTF-8 characters in order.  Only the rune array is maintained in sorted
// order. flush writes the resulting segment to a byte array.
type reorderBuffer struct ***REMOVED***
	rune  [maxBufferSize]Properties // Per character info.
	byte  [maxByteBufferSize]byte   // UTF-8 buffer. Referenced by runeInfo.pos.
	nbyte uint8                     // Number or bytes.
	ss    streamSafe                // For limiting length of non-starter sequence.
	nrune int                       // Number of runeInfos.
	f     formInfo

	src      input
	nsrc     int
	tmpBytes input

	out    []byte
	flushF func(*reorderBuffer) bool
***REMOVED***

func (rb *reorderBuffer) init(f Form, src []byte) ***REMOVED***
	rb.f = *formTable[f]
	rb.src.setBytes(src)
	rb.nsrc = len(src)
	rb.ss = 0
***REMOVED***

func (rb *reorderBuffer) initString(f Form, src string) ***REMOVED***
	rb.f = *formTable[f]
	rb.src.setString(src)
	rb.nsrc = len(src)
	rb.ss = 0
***REMOVED***

func (rb *reorderBuffer) setFlusher(out []byte, f func(*reorderBuffer) bool) ***REMOVED***
	rb.out = out
	rb.flushF = f
***REMOVED***

// reset discards all characters from the buffer.
func (rb *reorderBuffer) reset() ***REMOVED***
	rb.nrune = 0
	rb.nbyte = 0
***REMOVED***

func (rb *reorderBuffer) doFlush() bool ***REMOVED***
	if rb.f.composing ***REMOVED***
		rb.compose()
	***REMOVED***
	res := rb.flushF(rb)
	rb.reset()
	return res
***REMOVED***

// appendFlush appends the normalized segment to rb.out.
func appendFlush(rb *reorderBuffer) bool ***REMOVED***
	for i := 0; i < rb.nrune; i++ ***REMOVED***
		start := rb.rune[i].pos
		end := start + rb.rune[i].size
		rb.out = append(rb.out, rb.byte[start:end]...)
	***REMOVED***
	return true
***REMOVED***

// flush appends the normalized segment to out and resets rb.
func (rb *reorderBuffer) flush(out []byte) []byte ***REMOVED***
	for i := 0; i < rb.nrune; i++ ***REMOVED***
		start := rb.rune[i].pos
		end := start + rb.rune[i].size
		out = append(out, rb.byte[start:end]...)
	***REMOVED***
	rb.reset()
	return out
***REMOVED***

// flushCopy copies the normalized segment to buf and resets rb.
// It returns the number of bytes written to buf.
func (rb *reorderBuffer) flushCopy(buf []byte) int ***REMOVED***
	p := 0
	for i := 0; i < rb.nrune; i++ ***REMOVED***
		runep := rb.rune[i]
		p += copy(buf[p:], rb.byte[runep.pos:runep.pos+runep.size])
	***REMOVED***
	rb.reset()
	return p
***REMOVED***

// insertOrdered inserts a rune in the buffer, ordered by Canonical Combining Class.
// It returns false if the buffer is not large enough to hold the rune.
// It is used internally by insert and insertString only.
func (rb *reorderBuffer) insertOrdered(info Properties) ***REMOVED***
	n := rb.nrune
	b := rb.rune[:]
	cc := info.ccc
	if cc > 0 ***REMOVED***
		// Find insertion position + move elements to make room.
		for ; n > 0; n-- ***REMOVED***
			if b[n-1].ccc <= cc ***REMOVED***
				break
			***REMOVED***
			b[n] = b[n-1]
		***REMOVED***
	***REMOVED***
	rb.nrune += 1
	pos := uint8(rb.nbyte)
	rb.nbyte += utf8.UTFMax
	info.pos = pos
	b[n] = info
***REMOVED***

// insertErr is an error code returned by insert. Using this type instead
// of error improves performance up to 20% for many of the benchmarks.
type insertErr int

const (
	iSuccess insertErr = -iota
	iShortDst
	iShortSrc
)

// insertFlush inserts the given rune in the buffer ordered by CCC.
// If a decomposition with multiple segments are encountered, they leading
// ones are flushed.
// It returns a non-zero error code if the rune was not inserted.
func (rb *reorderBuffer) insertFlush(src input, i int, info Properties) insertErr ***REMOVED***
	if rune := src.hangul(i); rune != 0 ***REMOVED***
		rb.decomposeHangul(rune)
		return iSuccess
	***REMOVED***
	if info.hasDecomposition() ***REMOVED***
		return rb.insertDecomposed(info.Decomposition())
	***REMOVED***
	rb.insertSingle(src, i, info)
	return iSuccess
***REMOVED***

// insertUnsafe inserts the given rune in the buffer ordered by CCC.
// It is assumed there is sufficient space to hold the runes. It is the
// responsibility of the caller to ensure this. This can be done by checking
// the state returned by the streamSafe type.
func (rb *reorderBuffer) insertUnsafe(src input, i int, info Properties) ***REMOVED***
	if rune := src.hangul(i); rune != 0 ***REMOVED***
		rb.decomposeHangul(rune)
	***REMOVED***
	if info.hasDecomposition() ***REMOVED***
		// TODO: inline.
		rb.insertDecomposed(info.Decomposition())
	***REMOVED*** else ***REMOVED***
		rb.insertSingle(src, i, info)
	***REMOVED***
***REMOVED***

// insertDecomposed inserts an entry in to the reorderBuffer for each rune
// in dcomp. dcomp must be a sequence of decomposed UTF-8-encoded runes.
// It flushes the buffer on each new segment start.
func (rb *reorderBuffer) insertDecomposed(dcomp []byte) insertErr ***REMOVED***
	rb.tmpBytes.setBytes(dcomp)
	// As the streamSafe accounting already handles the counting for modifiers,
	// we don't have to call next. However, we do need to keep the accounting
	// intact when flushing the buffer.
	for i := 0; i < len(dcomp); ***REMOVED***
		info := rb.f.info(rb.tmpBytes, i)
		if info.BoundaryBefore() && rb.nrune > 0 && !rb.doFlush() ***REMOVED***
			return iShortDst
		***REMOVED***
		i += copy(rb.byte[rb.nbyte:], dcomp[i:i+int(info.size)])
		rb.insertOrdered(info)
	***REMOVED***
	return iSuccess
***REMOVED***

// insertSingle inserts an entry in the reorderBuffer for the rune at
// position i. info is the runeInfo for the rune at position i.
func (rb *reorderBuffer) insertSingle(src input, i int, info Properties) ***REMOVED***
	src.copySlice(rb.byte[rb.nbyte:], i, i+int(info.size))
	rb.insertOrdered(info)
***REMOVED***

// insertCGJ inserts a Combining Grapheme Joiner (0x034f) into rb.
func (rb *reorderBuffer) insertCGJ() ***REMOVED***
	rb.insertSingle(input***REMOVED***str: GraphemeJoiner***REMOVED***, 0, Properties***REMOVED***size: uint8(len(GraphemeJoiner))***REMOVED***)
***REMOVED***

// appendRune inserts a rune at the end of the buffer. It is used for Hangul.
func (rb *reorderBuffer) appendRune(r rune) ***REMOVED***
	bn := rb.nbyte
	sz := utf8.EncodeRune(rb.byte[bn:], rune(r))
	rb.nbyte += utf8.UTFMax
	rb.rune[rb.nrune] = Properties***REMOVED***pos: bn, size: uint8(sz)***REMOVED***
	rb.nrune++
***REMOVED***

// assignRune sets a rune at position pos. It is used for Hangul and recomposition.
func (rb *reorderBuffer) assignRune(pos int, r rune) ***REMOVED***
	bn := rb.rune[pos].pos
	sz := utf8.EncodeRune(rb.byte[bn:], rune(r))
	rb.rune[pos] = Properties***REMOVED***pos: bn, size: uint8(sz)***REMOVED***
***REMOVED***

// runeAt returns the rune at position n. It is used for Hangul and recomposition.
func (rb *reorderBuffer) runeAt(n int) rune ***REMOVED***
	inf := rb.rune[n]
	r, _ := utf8.DecodeRune(rb.byte[inf.pos : inf.pos+inf.size])
	return r
***REMOVED***

// bytesAt returns the UTF-8 encoding of the rune at position n.
// It is used for Hangul and recomposition.
func (rb *reorderBuffer) bytesAt(n int) []byte ***REMOVED***
	inf := rb.rune[n]
	return rb.byte[inf.pos : int(inf.pos)+int(inf.size)]
***REMOVED***

// For Hangul we combine algorithmically, instead of using tables.
const (
	hangulBase  = 0xAC00 // UTF-8(hangulBase) -> EA B0 80
	hangulBase0 = 0xEA
	hangulBase1 = 0xB0
	hangulBase2 = 0x80

	hangulEnd  = hangulBase + jamoLVTCount // UTF-8(0xD7A4) -> ED 9E A4
	hangulEnd0 = 0xED
	hangulEnd1 = 0x9E
	hangulEnd2 = 0xA4

	jamoLBase  = 0x1100 // UTF-8(jamoLBase) -> E1 84 00
	jamoLBase0 = 0xE1
	jamoLBase1 = 0x84
	jamoLEnd   = 0x1113
	jamoVBase  = 0x1161
	jamoVEnd   = 0x1176
	jamoTBase  = 0x11A7
	jamoTEnd   = 0x11C3

	jamoTCount   = 28
	jamoVCount   = 21
	jamoVTCount  = 21 * 28
	jamoLVTCount = 19 * 21 * 28
)

const hangulUTF8Size = 3

func isHangul(b []byte) bool ***REMOVED***
	if len(b) < hangulUTF8Size ***REMOVED***
		return false
	***REMOVED***
	b0 := b[0]
	if b0 < hangulBase0 ***REMOVED***
		return false
	***REMOVED***
	b1 := b[1]
	switch ***REMOVED***
	case b0 == hangulBase0:
		return b1 >= hangulBase1
	case b0 < hangulEnd0:
		return true
	case b0 > hangulEnd0:
		return false
	case b1 < hangulEnd1:
		return true
	***REMOVED***
	return b1 == hangulEnd1 && b[2] < hangulEnd2
***REMOVED***

func isHangulString(b string) bool ***REMOVED***
	if len(b) < hangulUTF8Size ***REMOVED***
		return false
	***REMOVED***
	b0 := b[0]
	if b0 < hangulBase0 ***REMOVED***
		return false
	***REMOVED***
	b1 := b[1]
	switch ***REMOVED***
	case b0 == hangulBase0:
		return b1 >= hangulBase1
	case b0 < hangulEnd0:
		return true
	case b0 > hangulEnd0:
		return false
	case b1 < hangulEnd1:
		return true
	***REMOVED***
	return b1 == hangulEnd1 && b[2] < hangulEnd2
***REMOVED***

// Caller must ensure len(b) >= 2.
func isJamoVT(b []byte) bool ***REMOVED***
	// True if (rune & 0xff00) == jamoLBase
	return b[0] == jamoLBase0 && (b[1]&0xFC) == jamoLBase1
***REMOVED***

func isHangulWithoutJamoT(b []byte) bool ***REMOVED***
	c, _ := utf8.DecodeRune(b)
	c -= hangulBase
	return c < jamoLVTCount && c%jamoTCount == 0
***REMOVED***

// decomposeHangul writes the decomposed Hangul to buf and returns the number
// of bytes written.  len(buf) should be at least 9.
func decomposeHangul(buf []byte, r rune) int ***REMOVED***
	const JamoUTF8Len = 3
	r -= hangulBase
	x := r % jamoTCount
	r /= jamoTCount
	utf8.EncodeRune(buf, jamoLBase+r/jamoVCount)
	utf8.EncodeRune(buf[JamoUTF8Len:], jamoVBase+r%jamoVCount)
	if x != 0 ***REMOVED***
		utf8.EncodeRune(buf[2*JamoUTF8Len:], jamoTBase+x)
		return 3 * JamoUTF8Len
	***REMOVED***
	return 2 * JamoUTF8Len
***REMOVED***

// decomposeHangul algorithmically decomposes a Hangul rune into
// its Jamo components.
// See https://unicode.org/reports/tr15/#Hangul for details on decomposing Hangul.
func (rb *reorderBuffer) decomposeHangul(r rune) ***REMOVED***
	r -= hangulBase
	x := r % jamoTCount
	r /= jamoTCount
	rb.appendRune(jamoLBase + r/jamoVCount)
	rb.appendRune(jamoVBase + r%jamoVCount)
	if x != 0 ***REMOVED***
		rb.appendRune(jamoTBase + x)
	***REMOVED***
***REMOVED***

// combineHangul algorithmically combines Jamo character components into Hangul.
// See https://unicode.org/reports/tr15/#Hangul for details on combining Hangul.
func (rb *reorderBuffer) combineHangul(s, i, k int) ***REMOVED***
	b := rb.rune[:]
	bn := rb.nrune
	for ; i < bn; i++ ***REMOVED***
		cccB := b[k-1].ccc
		cccC := b[i].ccc
		if cccB == 0 ***REMOVED***
			s = k - 1
		***REMOVED***
		if s != k-1 && cccB >= cccC ***REMOVED***
			// b[i] is blocked by greater-equal cccX below it
			b[k] = b[i]
			k++
		***REMOVED*** else ***REMOVED***
			l := rb.runeAt(s) // also used to compare to hangulBase
			v := rb.runeAt(i) // also used to compare to jamoT
			switch ***REMOVED***
			case jamoLBase <= l && l < jamoLEnd &&
				jamoVBase <= v && v < jamoVEnd:
				// 11xx plus 116x to LV
				rb.assignRune(s, hangulBase+
					(l-jamoLBase)*jamoVTCount+(v-jamoVBase)*jamoTCount)
			case hangulBase <= l && l < hangulEnd &&
				jamoTBase < v && v < jamoTEnd &&
				((l-hangulBase)%jamoTCount) == 0:
				// ACxx plus 11Ax to LVT
				rb.assignRune(s, l+v-jamoTBase)
			default:
				b[k] = b[i]
				k++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	rb.nrune = k
***REMOVED***

// compose recombines the runes in the buffer.
// It should only be used to recompose a single segment, as it will not
// handle alternations between Hangul and non-Hangul characters correctly.
func (rb *reorderBuffer) compose() ***REMOVED***
	// Lazily load the map used by the combine func below, but do
	// it outside of the loop.
	recompMapOnce.Do(buildRecompMap)

	// UAX #15, section X5 , including Corrigendum #5
	// "In any character sequence beginning with starter S, a character C is
	//  blocked from S if and only if there is some character B between S
	//  and C, and either B is a starter or it has the same or higher
	//  combining class as C."
	bn := rb.nrune
	if bn == 0 ***REMOVED***
		return
	***REMOVED***
	k := 1
	b := rb.rune[:]
	for s, i := 0, 1; i < bn; i++ ***REMOVED***
		if isJamoVT(rb.bytesAt(i)) ***REMOVED***
			// Redo from start in Hangul mode. Necessary to support
			// U+320E..U+321E in NFKC mode.
			rb.combineHangul(s, i, k)
			return
		***REMOVED***
		ii := b[i]
		// We can only use combineForward as a filter if we later
		// get the info for the combined character. This is more
		// expensive than using the filter. Using combinesBackward()
		// is safe.
		if ii.combinesBackward() ***REMOVED***
			cccB := b[k-1].ccc
			cccC := ii.ccc
			blocked := false // b[i] blocked by starter or greater or equal CCC?
			if cccB == 0 ***REMOVED***
				s = k - 1
			***REMOVED*** else ***REMOVED***
				blocked = s != k-1 && cccB >= cccC
			***REMOVED***
			if !blocked ***REMOVED***
				combined := combine(rb.runeAt(s), rb.runeAt(i))
				if combined != 0 ***REMOVED***
					rb.assignRune(s, combined)
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***
		b[k] = b[i]
		k++
	***REMOVED***
	rb.nrune = k
***REMOVED***
