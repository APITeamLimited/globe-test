package flate

import "fmt"

// fastGen maintains the table for matches,
// and the previous byte block for level 2.
// This is the generic implementation.
type fastEncL2 struct ***REMOVED***
	fastGen
	table [bTableSize]tableEntry
***REMOVED***

// EncodeL2 uses a similar algorithm to level 1, but is capable
// of matching across blocks giving better compression at a small slowdown.
func (e *fastEncL2) Encode(dst *tokens, src []byte) ***REMOVED***
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
		// When should we start skipping if we haven't found matches in a long while.
		const skipLog = 5
		const doEvery = 2

		nextS := s
		var candidate tableEntry
		for ***REMOVED***
			nextHash := hash4u(cv, bTableBits)
			s = nextS
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit ***REMOVED***
				goto emitRemainder
			***REMOVED***
			candidate = e.table[nextHash]
			now := load6432(src, nextS)
			e.table[nextHash] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***
			nextHash = hash4u(uint32(now), bTableBits)

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
				break
			***REMOVED***
			cv = uint32(now)
		***REMOVED***

		// A 4-byte match has been found. We'll later see if more than 4 bytes
		// match. But, prior to the match, src[nextEmit:s] are unmatched. Emit
		// them as literal bytes.

		// Call emitCopy, and then see if another emitCopy could be our next
		// move. Repeat until we find no match for the input immediately after
		// what was consumed by the last emitCopy call.
		//
		// If we exit this loop normally then we need to call emitLiteral next,
		// though we don't yet know how big the literal will be. We handle that
		// by proceeding to the next iteration of the main loop. We also can
		// exit this loop via goto if we get close to exhausting the input.
		for ***REMOVED***
			// Invariant: we have a 4-byte match at s, and no need to emit any
			// literal bytes prior to s.

			// Extend the 4-byte match as long as possible.
			t := candidate.offset - e.cur
			l := e.matchlenLong(s+4, t+4, src) + 4

			// Extend backwards
			for t > 0 && s > nextEmit && src[t-1] == src[s-1] ***REMOVED***
				s--
				t--
				l++
			***REMOVED***
			if nextEmit < s ***REMOVED***
				emitLiteral(dst, src[nextEmit:s])
			***REMOVED***

			dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
			s += l
			nextEmit = s
			if nextS >= s ***REMOVED***
				s = nextS + 1
			***REMOVED***

			if s >= sLimit ***REMOVED***
				// Index first pair after match end.
				if int(s+l+4) < len(src) ***REMOVED***
					cv := load3232(src, s)
					e.table[hash4u(cv, bTableBits)] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***
				***REMOVED***
				goto emitRemainder
			***REMOVED***

			// Store every second hash in-between, but offset by 1.
			for i := s - l + 2; i < s-5; i += 7 ***REMOVED***
				x := load6432(src, i)
				nextHash := hash4u(uint32(x), bTableBits)
				e.table[nextHash] = tableEntry***REMOVED***offset: e.cur + i***REMOVED***
				// Skip one
				x >>= 16
				nextHash = hash4u(uint32(x), bTableBits)
				e.table[nextHash] = tableEntry***REMOVED***offset: e.cur + i + 2***REMOVED***
				// Skip one
				x >>= 16
				nextHash = hash4u(uint32(x), bTableBits)
				e.table[nextHash] = tableEntry***REMOVED***offset: e.cur + i + 4***REMOVED***
			***REMOVED***

			// We could immediately start working at s now, but to improve
			// compression we first update the hash table at s-2 to s. If
			// another emitCopy is not our next move, also calculate nextHash
			// at s+1. At least on GOARCH=amd64, these three hash calculations
			// are faster as one load64 call (with some shifts) instead of
			// three load32 calls.
			x := load6432(src, s-2)
			o := e.cur + s - 2
			prevHash := hash4u(uint32(x), bTableBits)
			prevHash2 := hash4u(uint32(x>>8), bTableBits)
			e.table[prevHash] = tableEntry***REMOVED***offset: o***REMOVED***
			e.table[prevHash2] = tableEntry***REMOVED***offset: o + 1***REMOVED***
			currHash := hash4u(uint32(x>>16), bTableBits)
			candidate = e.table[currHash]
			e.table[currHash] = tableEntry***REMOVED***offset: o + 2***REMOVED***

			offset := s - (candidate.offset - e.cur)
			if offset > maxMatchOffset || uint32(x>>16) != load3232(src, candidate.offset-e.cur) ***REMOVED***
				cv = uint32(x >> 24)
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
