// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"fmt"
	"math"
	"math/bits"
)

const (
	tableBits      = 15                               // Bits used in the table
	tableSize      = 1 << tableBits                   // Size of the table
	tableShardCnt  = 1 << (tableBits - dictShardBits) // Number of shards in the table
	tableShardSize = tableSize / tableShardCnt        // Size of an individual shard
	tableMask      = tableSize - 1                    // Mask for table indices. Redundant, but can eliminate bounds checks.
	maxMatchLength = 131074
)

type tableEntry struct ***REMOVED***
	val    uint32
	offset int32
***REMOVED***

type fastEncoder struct ***REMOVED***
	fastBase
	table [tableSize]tableEntry
***REMOVED***

type fastEncoderDict struct ***REMOVED***
	fastEncoder
	dictTable       []tableEntry
	tableShardDirty [tableShardCnt]bool
	allDirty        bool
***REMOVED***

// Encode mimmics functionality in zstd_fast.c
func (e *fastEncoder) Encode(blk *blockEnc, src []byte) ***REMOVED***
	const (
		inputMargin            = 8
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	// Protect against e.cur wraparound.
	for e.cur >= bufferReset ***REMOVED***
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
		break
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
	const stepSize = 2

	// TEMPLATE
	const hashLog = tableBits
	// seems global, but would be nice to tweak.
	const kSearchStrength = 7

	// nextEmit is where in src the next emitLiteral should start from.
	nextEmit := s
	cv := load6432(src, s)

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
			if debugAsserts && canRepeat && offset1 == 0 ***REMOVED***
				panic("offset0 was 0")
			***REMOVED***

			nextHash := hash6(cv, hashLog)
			nextHash2 := hash6(cv>>8, hashLog)
			candidate := e.table[nextHash]
			candidate2 := e.table[nextHash2]
			repIndex := s - offset1 + 2

			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			e.table[nextHash2] = tableEntry***REMOVED***offset: s + e.cur + 1, val: uint32(cv >> 8)***REMOVED***

			if canRepeat && repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>16) ***REMOVED***
				// Consider history as well.
				var seq seq
				var length int32
				// length = 4 + e.matchlen(s+6, repIndex+4, src)
				***REMOVED***
					a := src[s+6:]
					b := src[repIndex+4:]
					endI := len(a) & (math.MaxInt32 - 7)
					length = int32(endI) + 4
					for i := 0; i < endI; i += 8 ***REMOVED***
						if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
							length = int32(i+bits.TrailingZeros64(diff)>>3) + 4
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***

				seq.matchLen = uint32(length - zstdMinMatch)

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
				s += length + 2
				nextEmit = s
				if s >= sLimit ***REMOVED***
					if debug ***REMOVED***
						println("repeat ended", s, length)

					***REMOVED***
					break encodeLoop
				***REMOVED***
				cv = load6432(src, s)
				continue
			***REMOVED***
			coffset0 := s - (candidate.offset - e.cur)
			coffset1 := s - (candidate2.offset - e.cur) + 1
			if coffset0 < e.maxMatchOff && uint32(cv) == candidate.val ***REMOVED***
				// found a regular match
				t = candidate.offset - e.cur
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				break
			***REMOVED***

			if coffset1 < e.maxMatchOff && uint32(cv>>8) == candidate2.val ***REMOVED***
				// found a regular match
				t = candidate2.offset - e.cur
				s++
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugAsserts && t < 0 ***REMOVED***
					panic("t<0")
				***REMOVED***
				break
			***REMOVED***
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***
		// A 4-byte match has been found. We'll later see if more than 4 bytes.
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t ***REMOVED***
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		***REMOVED***

		if debugAsserts && canRepeat && int(offset1) > len(src) ***REMOVED***
			panic("invalid offset")
		***REMOVED***

		// Extend the 4-byte match as long as possible.
		//l := e.matchlen(s+4, t+4, src) + 4
		var l int32
		***REMOVED***
			a := src[s+4:]
			b := src[t+4:]
			endI := len(a) & (math.MaxInt32 - 7)
			l = int32(endI) + 4
			for i := 0; i < endI; i += 8 ***REMOVED***
				if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
					l = int32(i+bits.TrailingZeros64(diff)>>3) + 4
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***

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

		// Check offset 2
		if o2 := s - offset2; canRepeat && load3232(src, o2) == uint32(cv) ***REMOVED***
			// We have at least 4 byte match.
			// No need to check backwards. We come straight from a match
			//l := 4 + e.matchlen(s+4, o2+4, src)
			var l int32
			***REMOVED***
				a := src[s+4:]
				b := src[o2+4:]
				endI := len(a) & (math.MaxInt32 - 7)
				l = int32(endI) + 4
				for i := 0; i < endI; i += 8 ***REMOVED***
					if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
						l = int32(i+bits.TrailingZeros64(diff)>>3) + 4
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Store this, since we have it.
			nextHash := hash6(cv, hashLog)
			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
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

// EncodeNoHist will encode a block with no history and no following blocks.
// Most notable difference is that src will not be copied for history and
// we do not need to check for max match length.
func (e *fastEncoder) EncodeNoHist(blk *blockEnc, src []byte) ***REMOVED***
	const (
		inputMargin            = 8
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)
	if debug ***REMOVED***
		if len(src) > maxBlockSize ***REMOVED***
			panic("src too big")
		***REMOVED***
	***REMOVED***

	// Protect against e.cur wraparound.
	if e.cur >= bufferReset ***REMOVED***
		for i := range e.table[:] ***REMOVED***
			e.table[i] = tableEntry***REMOVED******REMOVED***
		***REMOVED***
		e.cur = e.maxMatchOff
	***REMOVED***

	s := int32(0)
	blk.size = len(src)
	if len(src) < minNonLiteralBlockSize ***REMOVED***
		blk.extraLits = len(src)
		blk.literals = blk.literals[:len(src)]
		copy(blk.literals, src)
		return
	***REMOVED***

	sLimit := int32(len(src)) - inputMargin
	// stepSize is the number of bytes to skip on every main loop iteration.
	// It should be >= 2.
	const stepSize = 2

	// TEMPLATE
	const hashLog = tableBits
	// seems global, but would be nice to tweak.
	const kSearchStrength = 8

	// nextEmit is where in src the next emitLiteral should start from.
	nextEmit := s
	cv := load6432(src, s)

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

		for ***REMOVED***
			nextHash := hash6(cv, hashLog)
			nextHash2 := hash6(cv>>8, hashLog)
			candidate := e.table[nextHash]
			candidate2 := e.table[nextHash2]
			repIndex := s - offset1 + 2

			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			e.table[nextHash2] = tableEntry***REMOVED***offset: s + e.cur + 1, val: uint32(cv >> 8)***REMOVED***

			if len(blk.sequences) > 2 && load3232(src, repIndex) == uint32(cv>>16) ***REMOVED***
				// Consider history as well.
				var seq seq
				// length := 4 + e.matchlen(s+6, repIndex+4, src)
				// length := 4 + int32(matchLen(src[s+6:], src[repIndex+4:]))
				var length int32
				***REMOVED***
					a := src[s+6:]
					b := src[repIndex+4:]
					endI := len(a) & (math.MaxInt32 - 7)
					length = int32(endI) + 4
					for i := 0; i < endI; i += 8 ***REMOVED***
						if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
							length = int32(i+bits.TrailingZeros64(diff)>>3) + 4
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***

				seq.matchLen = uint32(length - zstdMinMatch)

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
				for repIndex > sMin && start > startLimit && src[repIndex-1] == src[start-1] ***REMOVED***
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
				s += length + 2
				nextEmit = s
				if s >= sLimit ***REMOVED***
					if debug ***REMOVED***
						println("repeat ended", s, length)

					***REMOVED***
					break encodeLoop
				***REMOVED***
				cv = load6432(src, s)
				continue
			***REMOVED***
			coffset0 := s - (candidate.offset - e.cur)
			coffset1 := s - (candidate2.offset - e.cur) + 1
			if coffset0 < e.maxMatchOff && uint32(cv) == candidate.val ***REMOVED***
				// found a regular match
				t = candidate.offset - e.cur
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugAsserts && t < 0 ***REMOVED***
					panic(fmt.Sprintf("t (%d) < 0, candidate.offset: %d, e.cur: %d, coffset0: %d, e.maxMatchOff: %d", t, candidate.offset, e.cur, coffset0, e.maxMatchOff))
				***REMOVED***
				break
			***REMOVED***

			if coffset1 < e.maxMatchOff && uint32(cv>>8) == candidate2.val ***REMOVED***
				// found a regular match
				t = candidate2.offset - e.cur
				s++
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugAsserts && t < 0 ***REMOVED***
					panic("t<0")
				***REMOVED***
				break
			***REMOVED***
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***
		// A 4-byte match has been found. We'll later see if more than 4 bytes.
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t ***REMOVED***
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		***REMOVED***

		if debugAsserts && t < 0 ***REMOVED***
			panic(fmt.Sprintf("t (%d) < 0 ", t))
		***REMOVED***
		// Extend the 4-byte match as long as possible.
		//l := e.matchlenNoHist(s+4, t+4, src) + 4
		// l := int32(matchLen(src[s+4:], src[t+4:])) + 4
		var l int32
		***REMOVED***
			a := src[s+4:]
			b := src[t+4:]
			endI := len(a) & (math.MaxInt32 - 7)
			l = int32(endI) + 4
			for i := 0; i < endI; i += 8 ***REMOVED***
				if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
					l = int32(i+bits.TrailingZeros64(diff)>>3) + 4
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Extend backwards
		tMin := s - e.maxMatchOff
		if tMin < 0 ***REMOVED***
			tMin = 0
		***REMOVED***
		for t > tMin && s > nextEmit && src[t-1] == src[s-1] ***REMOVED***
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

		// Check offset 2
		if o2 := s - offset2; len(blk.sequences) > 2 && load3232(src, o2) == uint32(cv) ***REMOVED***
			// We have at least 4 byte match.
			// No need to check backwards. We come straight from a match
			//l := 4 + e.matchlenNoHist(s+4, o2+4, src)
			// l := 4 + int32(matchLen(src[s+4:], src[o2+4:]))
			var l int32
			***REMOVED***
				a := src[s+4:]
				b := src[o2+4:]
				endI := len(a) & (math.MaxInt32 - 7)
				l = int32(endI) + 4
				for i := 0; i < endI; i += 8 ***REMOVED***
					if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
						l = int32(i+bits.TrailingZeros64(diff)>>3) + 4
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Store this, since we have it.
			nextHash := hash6(cv, hashLog)
			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
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
		***REMOVED***
	***REMOVED***

	if int(nextEmit) < len(src) ***REMOVED***
		blk.literals = append(blk.literals, src[nextEmit:]...)
		blk.extraLits = len(src) - int(nextEmit)
	***REMOVED***
	if debug ***REMOVED***
		println("returning, recent offsets:", blk.recentOffsets, "extra literals:", blk.extraLits)
	***REMOVED***
	// We do not store history, so we must offset e.cur to avoid false matches for next user.
	if e.cur < bufferReset ***REMOVED***
		e.cur += int32(len(src))
	***REMOVED***
***REMOVED***

// Encode will encode the content, with a dictionary if initialized for it.
func (e *fastEncoderDict) Encode(blk *blockEnc, src []byte) ***REMOVED***
	const (
		inputMargin            = 8
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)
	if e.allDirty || len(src) > 32<<10 ***REMOVED***
		e.fastEncoder.Encode(blk, src)
		e.allDirty = true
		return
	***REMOVED***
	// Protect against e.cur wraparound.
	for e.cur >= bufferReset ***REMOVED***
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
		break
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
	const stepSize = 2

	// TEMPLATE
	const hashLog = tableBits
	// seems global, but would be nice to tweak.
	const kSearchStrength = 7

	// nextEmit is where in src the next emitLiteral should start from.
	nextEmit := s
	cv := load6432(src, s)

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
			if debugAsserts && canRepeat && offset1 == 0 ***REMOVED***
				panic("offset0 was 0")
			***REMOVED***

			nextHash := hash6(cv, hashLog)
			nextHash2 := hash6(cv>>8, hashLog)
			candidate := e.table[nextHash]
			candidate2 := e.table[nextHash2]
			repIndex := s - offset1 + 2

			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			e.markShardDirty(nextHash)
			e.table[nextHash2] = tableEntry***REMOVED***offset: s + e.cur + 1, val: uint32(cv >> 8)***REMOVED***
			e.markShardDirty(nextHash2)

			if canRepeat && repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>16) ***REMOVED***
				// Consider history as well.
				var seq seq
				var length int32
				// length = 4 + e.matchlen(s+6, repIndex+4, src)
				***REMOVED***
					a := src[s+6:]
					b := src[repIndex+4:]
					endI := len(a) & (math.MaxInt32 - 7)
					length = int32(endI) + 4
					for i := 0; i < endI; i += 8 ***REMOVED***
						if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
							length = int32(i+bits.TrailingZeros64(diff)>>3) + 4
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***

				seq.matchLen = uint32(length - zstdMinMatch)

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
				s += length + 2
				nextEmit = s
				if s >= sLimit ***REMOVED***
					if debug ***REMOVED***
						println("repeat ended", s, length)

					***REMOVED***
					break encodeLoop
				***REMOVED***
				cv = load6432(src, s)
				continue
			***REMOVED***
			coffset0 := s - (candidate.offset - e.cur)
			coffset1 := s - (candidate2.offset - e.cur) + 1
			if coffset0 < e.maxMatchOff && uint32(cv) == candidate.val ***REMOVED***
				// found a regular match
				t = candidate.offset - e.cur
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				break
			***REMOVED***

			if coffset1 < e.maxMatchOff && uint32(cv>>8) == candidate2.val ***REMOVED***
				// found a regular match
				t = candidate2.offset - e.cur
				s++
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugAsserts && t < 0 ***REMOVED***
					panic("t<0")
				***REMOVED***
				break
			***REMOVED***
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***
		// A 4-byte match has been found. We'll later see if more than 4 bytes.
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t ***REMOVED***
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		***REMOVED***

		if debugAsserts && canRepeat && int(offset1) > len(src) ***REMOVED***
			panic("invalid offset")
		***REMOVED***

		// Extend the 4-byte match as long as possible.
		//l := e.matchlen(s+4, t+4, src) + 4
		var l int32
		***REMOVED***
			a := src[s+4:]
			b := src[t+4:]
			endI := len(a) & (math.MaxInt32 - 7)
			l = int32(endI) + 4
			for i := 0; i < endI; i += 8 ***REMOVED***
				if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
					l = int32(i+bits.TrailingZeros64(diff)>>3) + 4
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***

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

		// Check offset 2
		if o2 := s - offset2; canRepeat && load3232(src, o2) == uint32(cv) ***REMOVED***
			// We have at least 4 byte match.
			// No need to check backwards. We come straight from a match
			//l := 4 + e.matchlen(s+4, o2+4, src)
			var l int32
			***REMOVED***
				a := src[s+4:]
				b := src[o2+4:]
				endI := len(a) & (math.MaxInt32 - 7)
				l = int32(endI) + 4
				for i := 0; i < endI; i += 8 ***REMOVED***
					if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
						l = int32(i+bits.TrailingZeros64(diff)>>3) + 4
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Store this, since we have it.
			nextHash := hash6(cv, hashLog)
			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			e.markShardDirty(nextHash)
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

// ResetDict will reset and set a dictionary if not nil
func (e *fastEncoder) Reset(d *dict, singleBlock bool) ***REMOVED***
	e.resetBase(d, singleBlock)
	if d != nil ***REMOVED***
		panic("fastEncoder: Reset with dict")
	***REMOVED***
***REMOVED***

// ResetDict will reset and set a dictionary if not nil
func (e *fastEncoderDict) Reset(d *dict, singleBlock bool) ***REMOVED***
	e.resetBase(d, singleBlock)
	if d == nil ***REMOVED***
		return
	***REMOVED***

	// Init or copy dict table
	if len(e.dictTable) != len(e.table) || d.id != e.lastDictID ***REMOVED***
		if len(e.dictTable) != len(e.table) ***REMOVED***
			e.dictTable = make([]tableEntry, len(e.table))
		***REMOVED***
		if true ***REMOVED***
			end := e.maxMatchOff + int32(len(d.content)) - 8
			for i := e.maxMatchOff; i < end; i += 3 ***REMOVED***
				const hashLog = tableBits

				cv := load6432(d.content, i-e.maxMatchOff)
				nextHash := hash6(cv, hashLog)      // 0 -> 5
				nextHash1 := hash6(cv>>8, hashLog)  // 1 -> 6
				nextHash2 := hash6(cv>>16, hashLog) // 2 -> 7
				e.dictTable[nextHash] = tableEntry***REMOVED***
					val:    uint32(cv),
					offset: i,
				***REMOVED***
				e.dictTable[nextHash1] = tableEntry***REMOVED***
					val:    uint32(cv >> 8),
					offset: i + 1,
				***REMOVED***
				e.dictTable[nextHash2] = tableEntry***REMOVED***
					val:    uint32(cv >> 16),
					offset: i + 2,
				***REMOVED***
			***REMOVED***
		***REMOVED***
		e.lastDictID = d.id
		e.allDirty = true
	***REMOVED***

	e.cur = e.maxMatchOff
	dirtyShardCnt := 0
	if !e.allDirty ***REMOVED***
		for i := range e.tableShardDirty ***REMOVED***
			if e.tableShardDirty[i] ***REMOVED***
				dirtyShardCnt++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	const shardCnt = tableShardCnt
	const shardSize = tableShardSize
	if e.allDirty || dirtyShardCnt > shardCnt*4/6 ***REMOVED***
		copy(e.table[:], e.dictTable)
		for i := range e.tableShardDirty ***REMOVED***
			e.tableShardDirty[i] = false
		***REMOVED***
		e.allDirty = false
		return
	***REMOVED***
	for i := range e.tableShardDirty ***REMOVED***
		if !e.tableShardDirty[i] ***REMOVED***
			continue
		***REMOVED***

		copy(e.table[i*shardSize:(i+1)*shardSize], e.dictTable[i*shardSize:(i+1)*shardSize])
		e.tableShardDirty[i] = false
	***REMOVED***
	e.allDirty = false
***REMOVED***

func (e *fastEncoderDict) markAllShardsDirty() ***REMOVED***
	e.allDirty = true
***REMOVED***

func (e *fastEncoderDict) markShardDirty(entryNum uint32) ***REMOVED***
	e.tableShardDirty[entryNum/tableShardSize] = true
***REMOVED***
