package flate

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

// fastGen maintains the table for matches,
// and the previous byte block for level 2.
// This is the generic implementation.
type fastEncL1 struct ***REMOVED***
	fastGen
	table [tableSize]tableEntry
***REMOVED***

// EncodeL1 uses a similar algorithm to level 1
func (e *fastEncL1) Encode(dst *tokens, src []byte) ***REMOVED***
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)
	if debugDeflate && e.cur < 0 ***REMOVED***
		panic(fmt.Sprint("e.cur < 0: ", e.cur))
	***REMOVED***

	// Protect against e.cur wraparound.
	for e.cur >= bufferReset ***REMOVED***
		if len(e.hist) == 0 ***REMOVED***
			for i := range e.table[:] ***REMOVED***
				e.table[i] = tableEntry***REMOVED******REMOVED***
			***REMOVED***
			e.cur = maxMatchOffset
			break
		***REMOVED***
		// Shift down everything in the table that isn't already too far away.
		minOff := e.cur + int32(len(e.hist)) - maxMatchOffset
		for i := range e.table[:] ***REMOVED***
			v := e.table[i].offset
			if v <= minOff ***REMOVED***
				v = 0
			***REMOVED*** else ***REMOVED***
				v = v - e.cur + maxMatchOffset
			***REMOVED***
			e.table[i].offset = v
		***REMOVED***
		e.cur = maxMatchOffset
	***REMOVED***

	s := e.addBlock(src)

	// This check isn't in the Snappy implementation, but there, the caller
	// instead of the callee handles this case.
	if len(src) < minNonLiteralBlockSize ***REMOVED***
		// We do not fill the token table.
		// This will be picked up by caller.
		dst.n = uint16(len(src))
		return
	***REMOVED***

	// Override src
	src = e.hist
	nextEmit := s

	// sLimit is when to stop looking for offset/length copies. The inputMargin
	// lets us use a fast path for emitLiteral in the main loop, while we are
	// looking for copies.
	sLimit := int32(len(src) - inputMargin)

	// nextEmit is where in src the next emitLiteral should start from.
	cv := load3232(src, s)

	for ***REMOVED***
		const skipLog = 5
		const doEvery = 2

		nextS := s
		var candidate tableEntry
		for ***REMOVED***
			nextHash := hash(cv)
			candidate = e.table[nextHash]
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit ***REMOVED***
				goto emitRemainder
			***REMOVED***

			now := load6432(src, nextS)
			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***
			nextHash = hash(uint32(now))

			offset := s - (candidate.offset - e.cur)
			if offset < maxMatchOffset && cv == load3232(src, candidate.offset-e.cur) ***REMOVED***
				e.table[nextHash] = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***
				break
			***REMOVED***

			// Do one right away...
			cv = uint32(now)
			s = nextS
			nextS++
			candidate = e.table[nextHash]
			now >>= 8
			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***

			offset = s - (candidate.offset - e.cur)
			if offset < maxMatchOffset && cv == load3232(src, candidate.offset-e.cur) ***REMOVED***
				e.table[nextHash] = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***
				break
			***REMOVED***
			cv = uint32(now)
			s = nextS
		***REMOVED***

		// A 4-byte match has been found. We'll later see if more than 4 bytes
		// match. But, prior to the match, src[nextEmit:s] are unmatched. Emit
		// them as literal bytes.
		for ***REMOVED***
			// Invariant: we have a 4-byte match at s, and no need to emit any
			// literal bytes prior to s.

			// Extend the 4-byte match as long as possible.
			t := candidate.offset - e.cur
			var l = int32(4)
			if false ***REMOVED***
				l = e.matchlenLong(s+4, t+4, src) + 4
			***REMOVED*** else ***REMOVED***
				// inlined:
				a := src[s+4:]
				b := src[t+4:]
				for len(a) >= 8 ***REMOVED***
					if diff := binary.LittleEndian.Uint64(a) ^ binary.LittleEndian.Uint64(b); diff != 0 ***REMOVED***
						l += int32(bits.TrailingZeros64(diff) >> 3)
						break
					***REMOVED***
					l += 8
					a = a[8:]
					b = b[8:]
				***REMOVED***
				if len(a) < 8 ***REMOVED***
					b = b[:len(a)]
					for i := range a ***REMOVED***
						if a[i] != b[i] ***REMOVED***
							break
						***REMOVED***
						l++
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Extend backwards
			for t > 0 && s > nextEmit && src[t-1] == src[s-1] ***REMOVED***
				s--
				t--
				l++
			***REMOVED***
			if nextEmit < s ***REMOVED***
				if false ***REMOVED***
					emitLiteral(dst, src[nextEmit:s])
				***REMOVED*** else ***REMOVED***
					for _, v := range src[nextEmit:s] ***REMOVED***
						dst.tokens[dst.n] = token(v)
						dst.litHist[v]++
						dst.n++
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Save the match found
			if false ***REMOVED***
				dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
			***REMOVED*** else ***REMOVED***
				// Inlined...
				xoffset := uint32(s - t - baseMatchOffset)
				xlength := l
				oc := offsetCode(xoffset)
				xoffset |= oc << 16
				for xlength > 0 ***REMOVED***
					xl := xlength
					if xl > 258 ***REMOVED***
						if xl > 258+baseMatchLength ***REMOVED***
							xl = 258
						***REMOVED*** else ***REMOVED***
							xl = 258 - baseMatchLength
						***REMOVED***
					***REMOVED***
					xlength -= xl
					xl -= baseMatchLength
					dst.extraHist[lengthCodes1[uint8(xl)]]++
					dst.offHist[oc]++
					dst.tokens[dst.n] = token(matchType | uint32(xl)<<lengthShift | xoffset)
					dst.n++
				***REMOVED***
			***REMOVED***
			s += l
			nextEmit = s
			if nextS >= s ***REMOVED***
				s = nextS + 1
			***REMOVED***
			if s >= sLimit ***REMOVED***
				// Index first pair after match end.
				if int(s+l+4) < len(src) ***REMOVED***
					cv := load3232(src, s)
					e.table[hash(cv)] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***
				***REMOVED***
				goto emitRemainder
			***REMOVED***

			// We could immediately start working at s now, but to improve
			// compression we first update the hash table at s-2 and at s. If
			// another emitCopy is not our next move, also calculate nextHash
			// at s+1. At least on GOARCH=amd64, these three hash calculations
			// are faster as one load64 call (with some shifts) instead of
			// three load32 calls.
			x := load6432(src, s-2)
			o := e.cur + s - 2
			prevHash := hash(uint32(x))
			e.table[prevHash] = tableEntry***REMOVED***offset: o***REMOVED***
			x >>= 16
			currHash := hash(uint32(x))
			candidate = e.table[currHash]
			e.table[currHash] = tableEntry***REMOVED***offset: o + 2***REMOVED***

			offset := s - (candidate.offset - e.cur)
			if offset > maxMatchOffset || uint32(x) != load3232(src, candidate.offset-e.cur) ***REMOVED***
				cv = uint32(x >> 8)
				s++
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

emitRemainder:
	if int(nextEmit) < len(src) ***REMOVED***
		// If nothing was added, don't encode literals.
		if dst.n == 0 ***REMOVED***
			return
		***REMOVED***
		emitLiteral(dst, src[nextEmit:])
	***REMOVED***
***REMOVED***
