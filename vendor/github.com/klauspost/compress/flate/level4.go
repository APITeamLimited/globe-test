package flate

import "fmt"

type fastEncL4 struct ***REMOVED***
	fastGen
	table  [tableSize]tableEntry
	bTable [tableSize]tableEntry
***REMOVED***

func (e *fastEncL4) Encode(dst *tokens, src []byte) ***REMOVED***
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
			for i := range e.bTable[:] ***REMOVED***
				e.bTable[i] = tableEntry***REMOVED******REMOVED***
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
		for i := range e.bTable[:] ***REMOVED***
			v := e.bTable[i].offset
			if v <= minOff ***REMOVED***
				v = 0
			***REMOVED*** else ***REMOVED***
				v = v - e.cur + maxMatchOffset
			***REMOVED***
			e.bTable[i].offset = v
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
	cv := load6432(src, s)
	for ***REMOVED***
		const skipLog = 6
		const doEvery = 1

		nextS := s
		var t int32
		for ***REMOVED***
			nextHashS := hash4x64(cv, tableBits)
			nextHashL := hash7(cv, tableBits)

			s = nextS
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit ***REMOVED***
				goto emitRemainder
			***REMOVED***
			// Fetch a short+long candidate
			sCandidate := e.table[nextHashS]
			lCandidate := e.bTable[nextHashL]
			next := load6432(src, nextS)
			entry := tableEntry***REMOVED***offset: s + e.cur***REMOVED***
			e.table[nextHashS] = entry
			e.bTable[nextHashL] = entry

			t = lCandidate.offset - e.cur
			if s-t < maxMatchOffset && uint32(cv) == load3232(src, lCandidate.offset-e.cur) ***REMOVED***
				// We got a long match. Use that.
				break
			***REMOVED***

			t = sCandidate.offset - e.cur
			if s-t < maxMatchOffset && uint32(cv) == load3232(src, sCandidate.offset-e.cur) ***REMOVED***
				// Found a 4 match...
				lCandidate = e.bTable[hash7(next, tableBits)]

				// If the next long is a candidate, check if we should use that instead...
				lOff := nextS - (lCandidate.offset - e.cur)
				if lOff < maxMatchOffset && load3232(src, lCandidate.offset-e.cur) == uint32(next) ***REMOVED***
					l1, l2 := matchLen(src[s+4:], src[t+4:]), matchLen(src[nextS+4:], src[nextS-lOff+4:])
					if l2 > l1 ***REMOVED***
						s = nextS
						t = lCandidate.offset - e.cur
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***
			cv = next
		***REMOVED***

		// A 4-byte match has been found. We'll later see if more than 4 bytes
		// match. But, prior to the match, src[nextEmit:s] are unmatched. Emit
		// them as literal bytes.

		// Extend the 4-byte match as long as possible.
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
		if debugDeflate ***REMOVED***
			if t >= s ***REMOVED***
				panic("s-t")
			***REMOVED***
			if (s - t) > maxMatchOffset ***REMOVED***
				panic(fmt.Sprintln("mmo", t))
			***REMOVED***
			if l < baseMatchLength ***REMOVED***
				panic("bml")
			***REMOVED***
		***REMOVED***

		dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
		s += l
		nextEmit = s
		if nextS >= s ***REMOVED***
			s = nextS + 1
		***REMOVED***

		if s >= sLimit ***REMOVED***
			// Index first pair after match end.
			if int(s+8) < len(src) ***REMOVED***
				cv := load6432(src, s)
				e.table[hash4x64(cv, tableBits)] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***
				e.bTable[hash7(cv, tableBits)] = tableEntry***REMOVED***offset: s + e.cur***REMOVED***
			***REMOVED***
			goto emitRemainder
		***REMOVED***

		// Store every 3rd hash in-between
		if true ***REMOVED***
			i := nextS
			if i < s-1 ***REMOVED***
				cv := load6432(src, i)
				t := tableEntry***REMOVED***offset: i + e.cur***REMOVED***
				t2 := tableEntry***REMOVED***offset: t.offset + 1***REMOVED***
				e.bTable[hash7(cv, tableBits)] = t
				e.bTable[hash7(cv>>8, tableBits)] = t2
				e.table[hash4u(uint32(cv>>8), tableBits)] = t2

				i += 3
				for ; i < s-1; i += 3 ***REMOVED***
					cv := load6432(src, i)
					t := tableEntry***REMOVED***offset: i + e.cur***REMOVED***
					t2 := tableEntry***REMOVED***offset: t.offset + 1***REMOVED***
					e.bTable[hash7(cv, tableBits)] = t
					e.bTable[hash7(cv>>8, tableBits)] = t2
					e.table[hash4u(uint32(cv>>8), tableBits)] = t2
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// We could immediately start working at s now, but to improve
		// compression we first update the hash table at s-1 and at s.
		x := load6432(src, s-1)
		o := e.cur + s - 1
		prevHashS := hash4x64(x, tableBits)
		prevHashL := hash7(x, tableBits)
		e.table[prevHashS] = tableEntry***REMOVED***offset: o***REMOVED***
		e.bTable[prevHashL] = tableEntry***REMOVED***offset: o***REMOVED***
		cv = x >> 8
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
