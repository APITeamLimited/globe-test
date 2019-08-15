// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"math/bits"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

const (
	tableBits      = 15             // Bits used in the table
	tableSize      = 1 << tableBits // Size of the table
	tableMask      = tableSize - 1  // Mask for table indices. Redundant, but can eliminate bounds checks.
	maxMatchLength = 131074
)

type tableEntry struct ***REMOVED***
	val    uint32
	offset int32
***REMOVED***

type fastEncoder struct ***REMOVED***
	o encParams
	// cur is the offset at the start of hist
	cur int32
	// maximum offset. Should be at least 2x block size.
	maxMatchOff int32
	hist        []byte
	crc         *xxhash.Digest
	table       [tableSize]tableEntry
	tmp         [8]byte
	blk         *blockEnc
***REMOVED***

// CRC returns the underlying CRC writer.
func (e *fastEncoder) CRC() *xxhash.Digest ***REMOVED***
	return e.crc
***REMOVED***

// AppendCRC will append the CRC to the destination slice and return it.
func (e *fastEncoder) AppendCRC(dst []byte) []byte ***REMOVED***
	crc := e.crc.Sum(e.tmp[:0])
	dst = append(dst, crc[7], crc[6], crc[5], crc[4])
	return dst
***REMOVED***

// WindowSize returns the window size of the encoder,
// or a window size small enough to contain the input size, if > 0.
func (e *fastEncoder) WindowSize(size int) int32 ***REMOVED***
	if size > 0 && size < int(e.maxMatchOff) ***REMOVED***
		b := int32(1) << uint(bits.Len(uint(size)))
		// Keep minimum window.
		if b < 1024 ***REMOVED***
			b = 1024
		***REMOVED***
		return b
	***REMOVED***
	return e.maxMatchOff
***REMOVED***

// Block returns the current block.
func (e *fastEncoder) Block() *blockEnc ***REMOVED***
	return e.blk
***REMOVED***

// Encode mimmics functionality in zstd_fast.c
func (e *fastEncoder) Encode(blk *blockEnc, src []byte) ***REMOVED***
	const (
		inputMargin            = 8
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	// Protect against e.cur wraparound.
	for e.cur > (1<<30)+e.maxMatchOff ***REMOVED***
		if len(e.hist) == 0 ***REMOVED***
			for i := range e.table[:] ***REMOVED***
				e.table[i] = tableEntry***REMOVED******REMOVED***
			***REMOVED***
			e.cur = e.maxMatchOff
			break
		***REMOVED***
		// Shift down everything in the table that isn't already too far away.
		minOff := e.cur + int32(len(e.hist)) - e.maxMatchOff
		for i := range e.table[:] ***REMOVED***
			v := e.table[i].offset
			if v < minOff ***REMOVED***
				v = 0
			***REMOVED*** else ***REMOVED***
				v = v - e.cur + e.maxMatchOff
			***REMOVED***
			e.table[i].offset = v
		***REMOVED***
		e.cur = e.maxMatchOff
	***REMOVED***

	s := e.addBlock(src)
	blk.size = len(src)
	if len(src) < minNonLiteralBlockSize ***REMOVED***
		blk.extraLits = len(src)
		blk.literals = blk.literals[:len(src)]
		copy(blk.literals, src)
		return
	***REMOVED***

	// Override src
	src = e.hist
	sLimit := int32(len(src)) - inputMargin
	// stepSize is the number of bytes to skip on every main loop iteration.
	// It should be >= 2.
	stepSize := int32(e.o.targetLength)
	if stepSize == 0 ***REMOVED***
		stepSize++
	***REMOVED***
	stepSize++

	// TEMPLATE
	const hashLog = tableBits
	// seems global, but would be nice to tweak.
	const kSearchStrength = 8

	// nextEmit is where in src the next emitLiteral should start from.
	nextEmit := s
	cv := load6432(src, s)
	// nextHash is the hash at s
	nextHash := hash6(cv, hashLog)

	// Relative offsets
	offset1 := int32(blk.recentOffsets[0])
	offset2 := int32(blk.recentOffsets[1])

	addLiterals := func(s *seq, until int32) ***REMOVED***
		if until == nextEmit ***REMOVED***
			return
		***REMOVED***
		blk.literals = append(blk.literals, src[nextEmit:until]...)
		s.litLen = uint32(until - nextEmit)
	***REMOVED***
	if debug ***REMOVED***
		println("recent offsets:", blk.recentOffsets)
	***REMOVED***

encodeLoop:
	for ***REMOVED***
		// t will contain the match offset when we find one.
		// When existing the search loop, we have already checked 4 bytes.
		var t int32

		// We will not use repeat offsets across blocks.
		// By not using them for the first 3 matches
		canRepeat := len(blk.sequences) > 2

		for ***REMOVED***
			if debug && canRepeat && offset1 == 0 ***REMOVED***
				panic("offset0 was 0")
			***REMOVED***

			nextHash2 := hash6(cv>>8, hashLog) & tableMask
			nextHash = nextHash & tableMask
			candidate := e.table[nextHash]
			candidate2 := e.table[nextHash2]
			repIndex := s - offset1 + 2

			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			e.table[nextHash2] = tableEntry***REMOVED***offset: s + e.cur + 1, val: uint32(cv >> 8)***REMOVED***

			if canRepeat && repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>16) ***REMOVED***
				// Consider history as well.
				var seq seq
				lenght := 4 + e.matchlen(s+6, repIndex+4, src)

				seq.matchLen = uint32(lenght - zstdMinMatch)

				// We might be able to match backwards.
				// Extend as long as we can.
				start := s + 2
				// We end the search early, so we don't risk 0 literals
				// and have to do special offset treatment.
				startLimit := nextEmit + 1

				sMin := s - e.maxMatchOff
				if sMin < 0 ***REMOVED***
					sMin = 0
				***REMOVED***
				for repIndex > sMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch ***REMOVED***
					repIndex--
					start--
					seq.matchLen++
				***REMOVED***
				addLiterals(&seq, start)

				// rep 0
				seq.offset = 1
				if debugSequences ***REMOVED***
					println("repeat sequence", seq, "next s:", s)
				***REMOVED***
				blk.sequences = append(blk.sequences, seq)
				s += lenght + 2
				nextEmit = s
				if s >= sLimit ***REMOVED***
					if debug ***REMOVED***
						println("repeat ended", s, lenght)

					***REMOVED***
					break encodeLoop
				***REMOVED***
				cv = load6432(src, s)
				//nextHash = hashLen(cv, hashLog, mls)
				nextHash = hash6(cv, hashLog)
				continue
			***REMOVED***
			coffset0 := s - (candidate.offset - e.cur)
			coffset1 := s - (candidate2.offset - e.cur) + 1
			if coffset0 < e.maxMatchOff && uint32(cv) == candidate.val ***REMOVED***
				// found a regular match
				t = candidate.offset - e.cur
				if debug && s <= t ***REMOVED***
					panic("s <= t")
				***REMOVED***
				if debug && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				break
			***REMOVED***

			if coffset1 < e.maxMatchOff && uint32(cv>>8) == candidate2.val ***REMOVED***
				// found a regular match
				t = candidate2.offset - e.cur
				s++
				if debug && s <= t ***REMOVED***
					panic("s <= t")
				***REMOVED***
				if debug && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debug && t < 0 ***REMOVED***
					panic("t<0")
				***REMOVED***
				break
			***REMOVED***
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
			nextHash = hash6(cv, hashLog)
		***REMOVED***
		// A 4-byte match has been found. We'll later see if more than 4 bytes.
		offset2 = offset1
		offset1 = s - t

		if debug && s <= t ***REMOVED***
			panic("s <= t")
		***REMOVED***

		if debug && canRepeat && int(offset1) > len(src) ***REMOVED***
			panic("invalid offset")
		***REMOVED***

		// Extend the 4-byte match as long as possible.
		l := e.matchlen(s+4, t+4, src) + 4

		// Extend backwards
		tMin := s - e.maxMatchOff
		if tMin < 0 ***REMOVED***
			tMin = 0
		***REMOVED***
		for t > tMin && s > nextEmit && src[t-1] == src[s-1] && l < maxMatchLength ***REMOVED***
			s--
			t--
			l++
		***REMOVED***

		// Write our sequence.
		var seq seq
		seq.litLen = uint32(s - nextEmit)
		seq.matchLen = uint32(l - zstdMinMatch)
		if seq.litLen > 0 ***REMOVED***
			blk.literals = append(blk.literals, src[nextEmit:s]...)
		***REMOVED***
		// Don't use repeat offsets
		seq.offset = uint32(s-t) + 3
		s += l
		if debugSequences ***REMOVED***
			println("sequence", seq, "next s:", s)
		***REMOVED***
		blk.sequences = append(blk.sequences, seq)
		nextEmit = s
		if s >= sLimit ***REMOVED***
			break encodeLoop
		***REMOVED***
		cv = load6432(src, s)
		nextHash = hash6(cv, hashLog)

		// Check offset 2
		if o2 := s - offset2; canRepeat && o2 > 0 && load3232(src, o2) == uint32(cv) ***REMOVED***
			// We have at least 4 byte match.
			// No need to check backwards. We come straight from a match
			l := 4 + e.matchlen(s+4, o2+4, src)
			// Store this, since we have it.
			e.table[nextHash&tableMask] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			seq.matchLen = uint32(l) - zstdMinMatch
			seq.litLen = 0
			// Since litlen is always 0, this is offset 1.
			seq.offset = 1
			s += l
			nextEmit = s
			if debugSequences ***REMOVED***
				println("sequence", seq, "next s:", s)
			***REMOVED***
			blk.sequences = append(blk.sequences, seq)

			// Swap offset 1 and 2.
			offset1, offset2 = offset2, offset1
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			// Prepare next loop.
			cv = load6432(src, s)
			nextHash = hash6(cv, hashLog)
		***REMOVED***
	***REMOVED***

	if int(nextEmit) < len(src) ***REMOVED***
		blk.literals = append(blk.literals, src[nextEmit:]...)
		blk.extraLits = len(src) - int(nextEmit)
	***REMOVED***
	blk.recentOffsets[0] = uint32(offset1)
	blk.recentOffsets[1] = uint32(offset2)
	if debug ***REMOVED***
		println("returning, recent offsets:", blk.recentOffsets, "extra literals:", blk.extraLits)
	***REMOVED***
***REMOVED***

func (e *fastEncoder) addBlock(src []byte) int32 ***REMOVED***
	// check if we have space already
	if len(e.hist)+len(src) > cap(e.hist) ***REMOVED***
		if cap(e.hist) == 0 ***REMOVED***
			l := e.maxMatchOff * 2
			// Make it at least 1MB.
			if l < 1<<20 ***REMOVED***
				l = 1 << 20
			***REMOVED***
			e.hist = make([]byte, 0, l)
		***REMOVED*** else ***REMOVED***
			if cap(e.hist) < int(e.maxMatchOff*2) ***REMOVED***
				panic("unexpected buffer size")
			***REMOVED***
			// Move down
			offset := int32(len(e.hist)) - e.maxMatchOff
			copy(e.hist[0:e.maxMatchOff], e.hist[offset:])
			e.cur += offset
			e.hist = e.hist[:e.maxMatchOff]
		***REMOVED***
	***REMOVED***
	s := int32(len(e.hist))
	e.hist = append(e.hist, src...)
	return s
***REMOVED***

// useBlock will replace the block with the provided one,
// but transfer recent offsets from the previous.
func (e *fastEncoder) UseBlock(enc *blockEnc) ***REMOVED***
	enc.reset(e.blk)
	e.blk = enc
***REMOVED***

func (e *fastEncoder) matchlen(s, t int32, src []byte) int32 ***REMOVED***
	if debug ***REMOVED***
		if s < 0 ***REMOVED***
			panic("s<0")
		***REMOVED***
		if t < 0 ***REMOVED***
			panic("t<0")
		***REMOVED***
		if s-t > e.maxMatchOff ***REMOVED***
			panic(s - t)
		***REMOVED***
	***REMOVED***
	s1 := int(s) + maxMatchLength - 4
	if s1 > len(src) ***REMOVED***
		s1 = len(src)
	***REMOVED***

	// Extend the match to be as long as possible.
	return int32(matchLen(src[s:s1], src[t:]))
***REMOVED***

// Reset the encoding table.
func (e *fastEncoder) Reset() ***REMOVED***
	if e.blk == nil ***REMOVED***
		e.blk = &blockEnc***REMOVED******REMOVED***
		e.blk.init()
	***REMOVED*** else ***REMOVED***
		e.blk.reset(nil)
	***REMOVED***
	e.blk.initNewEncode()
	if e.crc == nil ***REMOVED***
		e.crc = xxhash.New()
	***REMOVED*** else ***REMOVED***
		e.crc.Reset()
	***REMOVED***
	if cap(e.hist) < int(e.maxMatchOff*2) ***REMOVED***
		l := e.maxMatchOff * 2
		// Make it at least 1MB.
		if l < 1<<20 ***REMOVED***
			l = 1 << 20
		***REMOVED***
		e.hist = make([]byte, 0, l)
	***REMOVED***
	// We offset current position so everything will be out of reach
	e.cur += e.maxMatchOff + int32(len(e.hist))
	e.hist = e.hist[:0]
***REMOVED***
