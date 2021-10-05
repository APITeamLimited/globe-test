// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import "fmt"

const (
	betterLongTableBits = 19                       // Bits used in the long match table
	betterLongTableSize = 1 << betterLongTableBits // Size of the table
	betterLongLen       = 8                        // Bytes used for table hash

	// Note: Increasing the short table bits or making the hash shorter
	// can actually lead to compression degradation since it will 'steal' more from the
	// long match table and match offsets are quite big.
	// This greatly depends on the type of input.
	betterShortTableBits = 13                        // Bits used in the short match table
	betterShortTableSize = 1 << betterShortTableBits // Size of the table
	betterShortLen       = 5                         // Bytes used for table hash

	betterLongTableShardCnt  = 1 << (betterLongTableBits - dictShardBits)    // Number of shards in the table
	betterLongTableShardSize = betterLongTableSize / betterLongTableShardCnt // Size of an individual shard

	betterShortTableShardCnt  = 1 << (betterShortTableBits - dictShardBits)     // Number of shards in the table
	betterShortTableShardSize = betterShortTableSize / betterShortTableShardCnt // Size of an individual shard
)

type prevEntry struct ***REMOVED***
	offset int32
	prev   int32
***REMOVED***

// betterFastEncoder uses 2 tables, one for short matches (5 bytes) and one for long matches.
// The long match table contains the previous entry with the same hash,
// effectively making it a "chain" of length 2.
// When we find a long match we choose between the two values and select the longest.
// When we find a short match, after checking the long, we check if we can find a long at n+1
// and that it is longer (lazy matching).
type betterFastEncoder struct ***REMOVED***
	fastBase
	table     [betterShortTableSize]tableEntry
	longTable [betterLongTableSize]prevEntry
***REMOVED***

type betterFastEncoderDict struct ***REMOVED***
	betterFastEncoder
	dictTable            []tableEntry
	dictLongTable        []prevEntry
	shortTableShardDirty [betterShortTableShardCnt]bool
	longTableShardDirty  [betterLongTableShardCnt]bool
	allDirty             bool
***REMOVED***

// Encode improves compression...
func (e *betterFastEncoder) Encode(blk *blockEnc, src []byte) ***REMOVED***
	const (
		// Input margin is the number of bytes we read (8)
		// and the maximum we will read ahead (2)
		inputMargin            = 8 + 2
		minNonLiteralBlockSize = 16
	)

	// Protect against e.cur wraparound.
	for e.cur >= bufferReset ***REMOVED***
		if len(e.hist) == 0 ***REMOVED***
			for i := range e.table[:] ***REMOVED***
				e.table[i] = tableEntry***REMOVED******REMOVED***
			***REMOVED***
			for i := range e.longTable[:] ***REMOVED***
				e.longTable[i] = prevEntry***REMOVED******REMOVED***
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
		for i := range e.longTable[:] ***REMOVED***
			v := e.longTable[i].offset
			v2 := e.longTable[i].prev
			if v < minOff ***REMOVED***
				v = 0
				v2 = 0
			***REMOVED*** else ***REMOVED***
				v = v - e.cur + e.maxMatchOff
				if v2 < minOff ***REMOVED***
					v2 = 0
				***REMOVED*** else ***REMOVED***
					v2 = v2 - e.cur + e.maxMatchOff
				***REMOVED***
			***REMOVED***
			e.longTable[i] = prevEntry***REMOVED***
				offset: v,
				prev:   v2,
			***REMOVED***
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
	// It should be >= 1.
	const stepSize = 1

	const kSearchStrength = 9

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
	if debugEncoder ***REMOVED***
		println("recent offsets:", blk.recentOffsets)
	***REMOVED***

encodeLoop:
	for ***REMOVED***
		var t int32
		// We allow the encoder to optionally turn off repeat offsets across blocks
		canRepeat := len(blk.sequences) > 2
		var matched int32

		for ***REMOVED***
			if debugAsserts && canRepeat && offset1 == 0 ***REMOVED***
				panic("offset0 was 0")
			***REMOVED***

			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)
			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)
			candidateL := e.longTable[nextHashL]
			candidateS := e.table[nextHashS]

			const repOff = 1
			repIndex := s - offset1 + repOff
			off := s + e.cur
			e.longTable[nextHashL] = prevEntry***REMOVED***offset: off, prev: candidateL.offset***REMOVED***
			e.table[nextHashS] = tableEntry***REMOVED***offset: off, val: uint32(cv)***REMOVED***

			if canRepeat ***REMOVED***
				if repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>(repOff*8)) ***REMOVED***
					// Consider history as well.
					var seq seq
					lenght := 4 + e.matchlen(s+4+repOff, repIndex+4, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					// We might be able to match backwards.
					// Extend as long as we can.
					start := s + repOff
					// We end the search early, so we don't risk 0 literals
					// and have to do special offset treatment.
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 ***REMOVED***
						tMin = 0
					***REMOVED***
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 ***REMOVED***
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

					// Index match start+1 (long) -> s - 1
					index0 := s + repOff
					s += lenght + repOff

					nextEmit = s
					if s >= sLimit ***REMOVED***
						if debugEncoder ***REMOVED***
							println("repeat ended", s, lenght)

						***REMOVED***
						break encodeLoop
					***REMOVED***
					// Index skipped...
					for index0 < s-1 ***REMOVED***
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry***REMOVED***offset: off, prev: e.longTable[h0].offset***REMOVED***
						e.table[hashLen(cv1, betterShortTableBits, betterShortLen)] = tableEntry***REMOVED***offset: off + 1, val: uint32(cv1)***REMOVED***
						index0 += 2
					***REMOVED***
					cv = load6432(src, s)
					continue
				***REMOVED***
				const repOff2 = 1

				// We deviate from the reference encoder and also check offset 2.
				// Still slower and not much better, so disabled.
				// repIndex = s - offset2 + repOff2
				if false && repIndex >= 0 && load6432(src, repIndex) == load6432(src, s+repOff) ***REMOVED***
					// Consider history as well.
					var seq seq
					lenght := 8 + e.matchlen(s+8+repOff2, repIndex+8, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					// We might be able to match backwards.
					// Extend as long as we can.
					start := s + repOff2
					// We end the search early, so we don't risk 0 literals
					// and have to do special offset treatment.
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 ***REMOVED***
						tMin = 0
					***REMOVED***
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 ***REMOVED***
						repIndex--
						start--
						seq.matchLen++
					***REMOVED***
					addLiterals(&seq, start)

					// rep 2
					seq.offset = 2
					if debugSequences ***REMOVED***
						println("repeat sequence 2", seq, "next s:", s)
					***REMOVED***
					blk.sequences = append(blk.sequences, seq)

					index0 := s + repOff2
					s += lenght + repOff2
					nextEmit = s
					if s >= sLimit ***REMOVED***
						if debugEncoder ***REMOVED***
							println("repeat ended", s, lenght)

						***REMOVED***
						break encodeLoop
					***REMOVED***

					// Index skipped...
					for index0 < s-1 ***REMOVED***
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry***REMOVED***offset: off, prev: e.longTable[h0].offset***REMOVED***
						e.table[hashLen(cv1, betterShortTableBits, betterShortLen)] = tableEntry***REMOVED***offset: off + 1, val: uint32(cv1)***REMOVED***
						index0 += 2
					***REMOVED***
					cv = load6432(src, s)
					// Swap offsets
					offset1, offset2 = offset2, offset1
					continue
				***REMOVED***
			***REMOVED***
			// Find the offsets of our two matches.
			coffsetL := candidateL.offset - e.cur
			coffsetLP := candidateL.prev - e.cur

			// Check if we have a long match.
			if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) ***REMOVED***
				// Found a long match, at least 8 bytes.
				matched = e.matchlen(s+8, coffsetL+8, src) + 8
				t = coffsetL
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugMatches ***REMOVED***
					println("long match")
				***REMOVED***

				if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) ***REMOVED***
					// Found a long match, at least 8 bytes.
					prevMatch := e.matchlen(s+8, coffsetLP+8, src) + 8
					if prevMatch > matched ***REMOVED***
						matched = prevMatch
						t = coffsetLP
					***REMOVED***
					if debugAsserts && s <= t ***REMOVED***
						panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
					***REMOVED***
					if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
						panic("s - t >e.maxMatchOff")
					***REMOVED***
					if debugMatches ***REMOVED***
						println("long match")
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***

			// Check if we have a long match on prev.
			if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) ***REMOVED***
				// Found a long match, at least 8 bytes.
				matched = e.matchlen(s+8, coffsetLP+8, src) + 8
				t = coffsetLP
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugMatches ***REMOVED***
					println("long match")
				***REMOVED***
				break
			***REMOVED***

			coffsetS := candidateS.offset - e.cur

			// Check if we have a short match.
			if s-coffsetS < e.maxMatchOff && uint32(cv) == candidateS.val ***REMOVED***
				// found a regular match
				matched = e.matchlen(s+4, coffsetS+4, src) + 4

				// See if we can find a long match at s+1
				const checkAt = 1
				cv := load6432(src, s+checkAt)
				nextHashL = hashLen(cv, betterLongTableBits, betterLongLen)
				candidateL = e.longTable[nextHashL]
				coffsetL = candidateL.offset - e.cur

				// We can store it, since we have at least a 4 byte match.
				e.longTable[nextHashL] = prevEntry***REMOVED***offset: s + checkAt + e.cur, prev: candidateL.offset***REMOVED***
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) ***REMOVED***
					// Found a long match, at least 8 bytes.
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched ***REMOVED***
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches ***REMOVED***
							println("long match (after short)")
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***

				// Check prev long...
				coffsetL = candidateL.prev - e.cur
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) ***REMOVED***
					// Found a long match, at least 8 bytes.
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched ***REMOVED***
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches ***REMOVED***
							println("prev long match (after short)")
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				t = coffsetS
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugAsserts && t < 0 ***REMOVED***
					panic("t<0")
				***REMOVED***
				if debugMatches ***REMOVED***
					println("short match")
				***REMOVED***
				break
			***REMOVED***

			// No match found, move forward in input.
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***

		// Try to find a better match by searching for a long match at the end of the current best match
		if s+matched < sLimit ***REMOVED***
			nextHashL := hashLen(load6432(src, s+matched), betterLongTableBits, betterLongLen)
			cv := load3232(src, s)
			candidateL := e.longTable[nextHashL]
			coffsetL := candidateL.offset - e.cur - matched
			if coffsetL >= 0 && coffsetL < s && s-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) ***REMOVED***
				// Found a long match, at least 4 bytes.
				matchedNext := e.matchlen(s+4, coffsetL+4, src) + 4
				if matchedNext > matched ***REMOVED***
					t = coffsetL
					matched = matchedNext
					if debugMatches ***REMOVED***
						println("long match at end-of-match")
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Check prev long...
			if true ***REMOVED***
				coffsetL = candidateL.prev - e.cur - matched
				if coffsetL >= 0 && coffsetL < s && s-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) ***REMOVED***
					// Found a long match, at least 4 bytes.
					matchedNext := e.matchlen(s+4, coffsetL+4, src) + 4
					if matchedNext > matched ***REMOVED***
						t = coffsetL
						matched = matchedNext
						if debugMatches ***REMOVED***
							println("prev long match at end-of-match")
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// A match has been found. Update recent offsets.
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t ***REMOVED***
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		***REMOVED***

		if debugAsserts && canRepeat && int(offset1) > len(src) ***REMOVED***
			panic("invalid offset")
		***REMOVED***

		// Extend the n-byte match as long as possible.
		l := matched

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

		// Write our sequence
		var seq seq
		seq.litLen = uint32(s - nextEmit)
		seq.matchLen = uint32(l - zstdMinMatch)
		if seq.litLen > 0 ***REMOVED***
			blk.literals = append(blk.literals, src[nextEmit:s]...)
		***REMOVED***
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

		// Index match start+1 (long) -> s - 1
		index0 := s - l + 1
		for index0 < s-1 ***REMOVED***
			cv0 := load6432(src, index0)
			cv1 := cv0 >> 8
			h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
			off := index0 + e.cur
			e.longTable[h0] = prevEntry***REMOVED***offset: off, prev: e.longTable[h0].offset***REMOVED***
			e.table[hashLen(cv1, betterShortTableBits, betterShortLen)] = tableEntry***REMOVED***offset: off + 1, val: uint32(cv1)***REMOVED***
			index0 += 2
		***REMOVED***

		cv = load6432(src, s)
		if !canRepeat ***REMOVED***
			continue
		***REMOVED***

		// Check offset 2
		for ***REMOVED***
			o2 := s - offset2
			if load3232(src, o2) != uint32(cv) ***REMOVED***
				// Do regular search
				break
			***REMOVED***

			// Store this, since we have it.
			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)
			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)

			// We have at least 4 byte match.
			// No need to check backwards. We come straight from a match
			l := 4 + e.matchlen(s+4, o2+4, src)

			e.longTable[nextHashL] = prevEntry***REMOVED***offset: s + e.cur, prev: e.longTable[nextHashL].offset***REMOVED***
			e.table[nextHashS] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
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
				// Finished
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***
	***REMOVED***

	if int(nextEmit) < len(src) ***REMOVED***
		blk.literals = append(blk.literals, src[nextEmit:]...)
		blk.extraLits = len(src) - int(nextEmit)
	***REMOVED***
	blk.recentOffsets[0] = uint32(offset1)
	blk.recentOffsets[1] = uint32(offset2)
	if debugEncoder ***REMOVED***
		println("returning, recent offsets:", blk.recentOffsets, "extra literals:", blk.extraLits)
	***REMOVED***
***REMOVED***

// EncodeNoHist will encode a block with no history and no following blocks.
// Most notable difference is that src will not be copied for history and
// we do not need to check for max match length.
func (e *betterFastEncoder) EncodeNoHist(blk *blockEnc, src []byte) ***REMOVED***
	e.ensureHist(len(src))
	e.Encode(blk, src)
***REMOVED***

// Encode improves compression...
func (e *betterFastEncoderDict) Encode(blk *blockEnc, src []byte) ***REMOVED***
	const (
		// Input margin is the number of bytes we read (8)
		// and the maximum we will read ahead (2)
		inputMargin            = 8 + 2
		minNonLiteralBlockSize = 16
	)

	// Protect against e.cur wraparound.
	for e.cur >= bufferReset ***REMOVED***
		if len(e.hist) == 0 ***REMOVED***
			for i := range e.table[:] ***REMOVED***
				e.table[i] = tableEntry***REMOVED******REMOVED***
			***REMOVED***
			for i := range e.longTable[:] ***REMOVED***
				e.longTable[i] = prevEntry***REMOVED******REMOVED***
			***REMOVED***
			e.cur = e.maxMatchOff
			e.allDirty = true
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
		for i := range e.longTable[:] ***REMOVED***
			v := e.longTable[i].offset
			v2 := e.longTable[i].prev
			if v < minOff ***REMOVED***
				v = 0
				v2 = 0
			***REMOVED*** else ***REMOVED***
				v = v - e.cur + e.maxMatchOff
				if v2 < minOff ***REMOVED***
					v2 = 0
				***REMOVED*** else ***REMOVED***
					v2 = v2 - e.cur + e.maxMatchOff
				***REMOVED***
			***REMOVED***
			e.longTable[i] = prevEntry***REMOVED***
				offset: v,
				prev:   v2,
			***REMOVED***
		***REMOVED***
		e.allDirty = true
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
	// It should be >= 1.
	const stepSize = 1

	const kSearchStrength = 9

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
	if debugEncoder ***REMOVED***
		println("recent offsets:", blk.recentOffsets)
	***REMOVED***

encodeLoop:
	for ***REMOVED***
		var t int32
		// We allow the encoder to optionally turn off repeat offsets across blocks
		canRepeat := len(blk.sequences) > 2
		var matched int32

		for ***REMOVED***
			if debugAsserts && canRepeat && offset1 == 0 ***REMOVED***
				panic("offset0 was 0")
			***REMOVED***

			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)
			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)
			candidateL := e.longTable[nextHashL]
			candidateS := e.table[nextHashS]

			const repOff = 1
			repIndex := s - offset1 + repOff
			off := s + e.cur
			e.longTable[nextHashL] = prevEntry***REMOVED***offset: off, prev: candidateL.offset***REMOVED***
			e.markLongShardDirty(nextHashL)
			e.table[nextHashS] = tableEntry***REMOVED***offset: off, val: uint32(cv)***REMOVED***
			e.markShortShardDirty(nextHashS)

			if canRepeat ***REMOVED***
				if repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>(repOff*8)) ***REMOVED***
					// Consider history as well.
					var seq seq
					lenght := 4 + e.matchlen(s+4+repOff, repIndex+4, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					// We might be able to match backwards.
					// Extend as long as we can.
					start := s + repOff
					// We end the search early, so we don't risk 0 literals
					// and have to do special offset treatment.
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 ***REMOVED***
						tMin = 0
					***REMOVED***
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 ***REMOVED***
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

					// Index match start+1 (long) -> s - 1
					index0 := s + repOff
					s += lenght + repOff

					nextEmit = s
					if s >= sLimit ***REMOVED***
						if debugEncoder ***REMOVED***
							println("repeat ended", s, lenght)

						***REMOVED***
						break encodeLoop
					***REMOVED***
					// Index skipped...
					for index0 < s-1 ***REMOVED***
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry***REMOVED***offset: off, prev: e.longTable[h0].offset***REMOVED***
						e.markLongShardDirty(h0)
						h1 := hashLen(cv1, betterShortTableBits, betterShortLen)
						e.table[h1] = tableEntry***REMOVED***offset: off + 1, val: uint32(cv1)***REMOVED***
						e.markShortShardDirty(h1)
						index0 += 2
					***REMOVED***
					cv = load6432(src, s)
					continue
				***REMOVED***
				const repOff2 = 1

				// We deviate from the reference encoder and also check offset 2.
				// Still slower and not much better, so disabled.
				// repIndex = s - offset2 + repOff2
				if false && repIndex >= 0 && load6432(src, repIndex) == load6432(src, s+repOff) ***REMOVED***
					// Consider history as well.
					var seq seq
					lenght := 8 + e.matchlen(s+8+repOff2, repIndex+8, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					// We might be able to match backwards.
					// Extend as long as we can.
					start := s + repOff2
					// We end the search early, so we don't risk 0 literals
					// and have to do special offset treatment.
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 ***REMOVED***
						tMin = 0
					***REMOVED***
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 ***REMOVED***
						repIndex--
						start--
						seq.matchLen++
					***REMOVED***
					addLiterals(&seq, start)

					// rep 2
					seq.offset = 2
					if debugSequences ***REMOVED***
						println("repeat sequence 2", seq, "next s:", s)
					***REMOVED***
					blk.sequences = append(blk.sequences, seq)

					index0 := s + repOff2
					s += lenght + repOff2
					nextEmit = s
					if s >= sLimit ***REMOVED***
						if debugEncoder ***REMOVED***
							println("repeat ended", s, lenght)

						***REMOVED***
						break encodeLoop
					***REMOVED***

					// Index skipped...
					for index0 < s-1 ***REMOVED***
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry***REMOVED***offset: off, prev: e.longTable[h0].offset***REMOVED***
						e.markLongShardDirty(h0)
						h1 := hashLen(cv1, betterShortTableBits, betterShortLen)
						e.table[h1] = tableEntry***REMOVED***offset: off + 1, val: uint32(cv1)***REMOVED***
						e.markShortShardDirty(h1)
						index0 += 2
					***REMOVED***
					cv = load6432(src, s)
					// Swap offsets
					offset1, offset2 = offset2, offset1
					continue
				***REMOVED***
			***REMOVED***
			// Find the offsets of our two matches.
			coffsetL := candidateL.offset - e.cur
			coffsetLP := candidateL.prev - e.cur

			// Check if we have a long match.
			if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) ***REMOVED***
				// Found a long match, at least 8 bytes.
				matched = e.matchlen(s+8, coffsetL+8, src) + 8
				t = coffsetL
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugMatches ***REMOVED***
					println("long match")
				***REMOVED***

				if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) ***REMOVED***
					// Found a long match, at least 8 bytes.
					prevMatch := e.matchlen(s+8, coffsetLP+8, src) + 8
					if prevMatch > matched ***REMOVED***
						matched = prevMatch
						t = coffsetLP
					***REMOVED***
					if debugAsserts && s <= t ***REMOVED***
						panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
					***REMOVED***
					if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
						panic("s - t >e.maxMatchOff")
					***REMOVED***
					if debugMatches ***REMOVED***
						println("long match")
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***

			// Check if we have a long match on prev.
			if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) ***REMOVED***
				// Found a long match, at least 8 bytes.
				matched = e.matchlen(s+8, coffsetLP+8, src) + 8
				t = coffsetLP
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugMatches ***REMOVED***
					println("long match")
				***REMOVED***
				break
			***REMOVED***

			coffsetS := candidateS.offset - e.cur

			// Check if we have a short match.
			if s-coffsetS < e.maxMatchOff && uint32(cv) == candidateS.val ***REMOVED***
				// found a regular match
				matched = e.matchlen(s+4, coffsetS+4, src) + 4

				// See if we can find a long match at s+1
				const checkAt = 1
				cv := load6432(src, s+checkAt)
				nextHashL = hashLen(cv, betterLongTableBits, betterLongLen)
				candidateL = e.longTable[nextHashL]
				coffsetL = candidateL.offset - e.cur

				// We can store it, since we have at least a 4 byte match.
				e.longTable[nextHashL] = prevEntry***REMOVED***offset: s + checkAt + e.cur, prev: candidateL.offset***REMOVED***
				e.markLongShardDirty(nextHashL)
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) ***REMOVED***
					// Found a long match, at least 8 bytes.
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched ***REMOVED***
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches ***REMOVED***
							println("long match (after short)")
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***

				// Check prev long...
				coffsetL = candidateL.prev - e.cur
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) ***REMOVED***
					// Found a long match, at least 8 bytes.
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched ***REMOVED***
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches ***REMOVED***
							println("prev long match (after short)")
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				t = coffsetS
				if debugAsserts && s <= t ***REMOVED***
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				***REMOVED***
				if debugAsserts && s-t > e.maxMatchOff ***REMOVED***
					panic("s - t >e.maxMatchOff")
				***REMOVED***
				if debugAsserts && t < 0 ***REMOVED***
					panic("t<0")
				***REMOVED***
				if debugMatches ***REMOVED***
					println("short match")
				***REMOVED***
				break
			***REMOVED***

			// No match found, move forward in input.
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit ***REMOVED***
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***
		// Try to find a better match by searching for a long match at the end of the current best match
		if s+matched < sLimit ***REMOVED***
			nextHashL := hashLen(load6432(src, s+matched), betterLongTableBits, betterLongLen)
			cv := load3232(src, s)
			candidateL := e.longTable[nextHashL]
			coffsetL := candidateL.offset - e.cur - matched
			if coffsetL >= 0 && coffsetL < s && s-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) ***REMOVED***
				// Found a long match, at least 4 bytes.
				matchedNext := e.matchlen(s+4, coffsetL+4, src) + 4
				if matchedNext > matched ***REMOVED***
					t = coffsetL
					matched = matchedNext
					if debugMatches ***REMOVED***
						println("long match at end-of-match")
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Check prev long...
			if true ***REMOVED***
				coffsetL = candidateL.prev - e.cur - matched
				if coffsetL >= 0 && coffsetL < s && s-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) ***REMOVED***
					// Found a long match, at least 4 bytes.
					matchedNext := e.matchlen(s+4, coffsetL+4, src) + 4
					if matchedNext > matched ***REMOVED***
						t = coffsetL
						matched = matchedNext
						if debugMatches ***REMOVED***
							println("prev long match at end-of-match")
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// A match has been found. Update recent offsets.
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t ***REMOVED***
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		***REMOVED***

		if debugAsserts && canRepeat && int(offset1) > len(src) ***REMOVED***
			panic("invalid offset")
		***REMOVED***

		// Extend the n-byte match as long as possible.
		l := matched

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

		// Write our sequence
		var seq seq
		seq.litLen = uint32(s - nextEmit)
		seq.matchLen = uint32(l - zstdMinMatch)
		if seq.litLen > 0 ***REMOVED***
			blk.literals = append(blk.literals, src[nextEmit:s]...)
		***REMOVED***
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

		// Index match start+1 (long) -> s - 1
		index0 := s - l + 1
		for index0 < s-1 ***REMOVED***
			cv0 := load6432(src, index0)
			cv1 := cv0 >> 8
			h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
			off := index0 + e.cur
			e.longTable[h0] = prevEntry***REMOVED***offset: off, prev: e.longTable[h0].offset***REMOVED***
			e.markLongShardDirty(h0)
			h1 := hashLen(cv1, betterShortTableBits, betterShortLen)
			e.table[h1] = tableEntry***REMOVED***offset: off + 1, val: uint32(cv1)***REMOVED***
			e.markShortShardDirty(h1)
			index0 += 2
		***REMOVED***

		cv = load6432(src, s)
		if !canRepeat ***REMOVED***
			continue
		***REMOVED***

		// Check offset 2
		for ***REMOVED***
			o2 := s - offset2
			if load3232(src, o2) != uint32(cv) ***REMOVED***
				// Do regular search
				break
			***REMOVED***

			// Store this, since we have it.
			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)
			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)

			// We have at least 4 byte match.
			// No need to check backwards. We come straight from a match
			l := 4 + e.matchlen(s+4, o2+4, src)

			e.longTable[nextHashL] = prevEntry***REMOVED***offset: s + e.cur, prev: e.longTable[nextHashL].offset***REMOVED***
			e.markLongShardDirty(nextHashL)
			e.table[nextHashS] = tableEntry***REMOVED***offset: s + e.cur, val: uint32(cv)***REMOVED***
			e.markShortShardDirty(nextHashS)
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
				// Finished
				break encodeLoop
			***REMOVED***
			cv = load6432(src, s)
		***REMOVED***
	***REMOVED***

	if int(nextEmit) < len(src) ***REMOVED***
		blk.literals = append(blk.literals, src[nextEmit:]...)
		blk.extraLits = len(src) - int(nextEmit)
	***REMOVED***
	blk.recentOffsets[0] = uint32(offset1)
	blk.recentOffsets[1] = uint32(offset2)
	if debugEncoder ***REMOVED***
		println("returning, recent offsets:", blk.recentOffsets, "extra literals:", blk.extraLits)
	***REMOVED***
***REMOVED***

// ResetDict will reset and set a dictionary if not nil
func (e *betterFastEncoder) Reset(d *dict, singleBlock bool) ***REMOVED***
	e.resetBase(d, singleBlock)
	if d != nil ***REMOVED***
		panic("betterFastEncoder: Reset with dict")
	***REMOVED***
***REMOVED***

// ResetDict will reset and set a dictionary if not nil
func (e *betterFastEncoderDict) Reset(d *dict, singleBlock bool) ***REMOVED***
	e.resetBase(d, singleBlock)
	if d == nil ***REMOVED***
		return
	***REMOVED***
	// Init or copy dict table
	if len(e.dictTable) != len(e.table) || d.id != e.lastDictID ***REMOVED***
		if len(e.dictTable) != len(e.table) ***REMOVED***
			e.dictTable = make([]tableEntry, len(e.table))
		***REMOVED***
		end := int32(len(d.content)) - 8 + e.maxMatchOff
		for i := e.maxMatchOff; i < end; i += 4 ***REMOVED***
			const hashLog = betterShortTableBits

			cv := load6432(d.content, i-e.maxMatchOff)
			nextHash := hashLen(cv, hashLog, betterShortLen)      // 0 -> 4
			nextHash1 := hashLen(cv>>8, hashLog, betterShortLen)  // 1 -> 5
			nextHash2 := hashLen(cv>>16, hashLog, betterShortLen) // 2 -> 6
			nextHash3 := hashLen(cv>>24, hashLog, betterShortLen) // 3 -> 7
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
			e.dictTable[nextHash3] = tableEntry***REMOVED***
				val:    uint32(cv >> 24),
				offset: i + 3,
			***REMOVED***
		***REMOVED***
		e.lastDictID = d.id
		e.allDirty = true
	***REMOVED***

	// Init or copy dict table
	if len(e.dictLongTable) != len(e.longTable) || d.id != e.lastDictID ***REMOVED***
		if len(e.dictLongTable) != len(e.longTable) ***REMOVED***
			e.dictLongTable = make([]prevEntry, len(e.longTable))
		***REMOVED***
		if len(d.content) >= 8 ***REMOVED***
			cv := load6432(d.content, 0)
			h := hashLen(cv, betterLongTableBits, betterLongLen)
			e.dictLongTable[h] = prevEntry***REMOVED***
				offset: e.maxMatchOff,
				prev:   e.dictLongTable[h].offset,
			***REMOVED***

			end := int32(len(d.content)) - 8 + e.maxMatchOff
			off := 8 // First to read
			for i := e.maxMatchOff + 1; i < end; i++ ***REMOVED***
				cv = cv>>8 | (uint64(d.content[off]) << 56)
				h := hashLen(cv, betterLongTableBits, betterLongLen)
				e.dictLongTable[h] = prevEntry***REMOVED***
					offset: i,
					prev:   e.dictLongTable[h].offset,
				***REMOVED***
				off++
			***REMOVED***
		***REMOVED***
		e.lastDictID = d.id
		e.allDirty = true
	***REMOVED***

	// Reset table to initial state
	***REMOVED***
		dirtyShardCnt := 0
		if !e.allDirty ***REMOVED***
			for i := range e.shortTableShardDirty ***REMOVED***
				if e.shortTableShardDirty[i] ***REMOVED***
					dirtyShardCnt++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		const shardCnt = betterShortTableShardCnt
		const shardSize = betterShortTableShardSize
		if e.allDirty || dirtyShardCnt > shardCnt*4/6 ***REMOVED***
			copy(e.table[:], e.dictTable)
			for i := range e.shortTableShardDirty ***REMOVED***
				e.shortTableShardDirty[i] = false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := range e.shortTableShardDirty ***REMOVED***
				if !e.shortTableShardDirty[i] ***REMOVED***
					continue
				***REMOVED***

				copy(e.table[i*shardSize:(i+1)*shardSize], e.dictTable[i*shardSize:(i+1)*shardSize])
				e.shortTableShardDirty[i] = false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	***REMOVED***
		dirtyShardCnt := 0
		if !e.allDirty ***REMOVED***
			for i := range e.shortTableShardDirty ***REMOVED***
				if e.shortTableShardDirty[i] ***REMOVED***
					dirtyShardCnt++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		const shardCnt = betterLongTableShardCnt
		const shardSize = betterLongTableShardSize
		if e.allDirty || dirtyShardCnt > shardCnt*4/6 ***REMOVED***
			copy(e.longTable[:], e.dictLongTable)
			for i := range e.longTableShardDirty ***REMOVED***
				e.longTableShardDirty[i] = false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := range e.longTableShardDirty ***REMOVED***
				if !e.longTableShardDirty[i] ***REMOVED***
					continue
				***REMOVED***

				copy(e.longTable[i*shardSize:(i+1)*shardSize], e.dictLongTable[i*shardSize:(i+1)*shardSize])
				e.longTableShardDirty[i] = false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	e.cur = e.maxMatchOff
	e.allDirty = false
***REMOVED***

func (e *betterFastEncoderDict) markLongShardDirty(entryNum uint32) ***REMOVED***
	e.longTableShardDirty[entryNum/betterLongTableShardSize] = true
***REMOVED***

func (e *betterFastEncoderDict) markShortShardDirty(entryNum uint32) ***REMOVED***
	e.shortTableShardDirty[entryNum/betterShortTableShardSize] = true
***REMOVED***
