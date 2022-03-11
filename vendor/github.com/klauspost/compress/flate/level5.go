package flate

import "fmt"

type fastEncL5 struct ***REMOVED***
	fastGen
	table  [tableSize]tableEntry
	bTable [tableSize]tableEntryPrev
***REMOVED***

func (e *fastEncL5) Encode(dst *tokens, src []byte) ***REMOVED***
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
				e.bTable[i] = tableEntryPrev***REMOVED******REMOVED***
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
			v := e.bTable[i]
			if v.Cur.offset <= minOff ***REMOVED***
				v.Cur.offset = 0
				v.Prev.offset = 0
			***REMOVED*** else ***REMOVED***
				v.Cur.offset = v.Cur.offset - e.cur + maxMatchOffset
				if v.Prev.offset <= minOff ***REMOVED***
					v.Prev.offset = 0
				***REMOVED*** else ***REMOVED***
					v.Prev.offset = v.Prev.offset - e.cur + maxMatchOffset
				***REMOVED***
			***REMOVED***
			e.bTable[i] = v
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
		var l int32
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
			eLong := &e.bTable[nextHashL]
			eLong.Cur, eLong.Prev = entry, eLong.Cur

			nextHashS = hash4x64(next, tableBits)
			nextHashL = hash7(next, tableBits)

			t = lCandidate.Cur.offset - e.cur
			if s-t < maxMatchOffset ***REMOVED***
				if uint32(cv) == load3232(src, lCandidate.Cur.offset-e.cur) ***REMOVED***
					// Store the next match
					e.table[nextHashS] = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***
					eLong := &e.bTable[nextHashL]
					eLong.Cur, eLong.Prev = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***, eLong.Cur

					t2 := lCandidate.Prev.offset - e.cur
					if s-t2 < maxMatchOffset && uint32(cv) == load3232(src, lCandidate.Prev.offset-e.cur) ***REMOVED***
						l = e.matchlen(s+4, t+4, src) + 4
						ml1 := e.matchlen(s+4, t2+4, src) + 4
						if ml1 > l ***REMOVED***
							t = t2
							l = ml1
							break
						***REMOVED***
					***REMOVED***
					break
				***REMOVED***
				t = lCandidate.Prev.offset - e.cur
				if s-t < maxMatchOffset && uint32(cv) == load3232(src, lCandidate.Prev.offset-e.cur) ***REMOVED***
					// Store the next match
					e.table[nextHashS] = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***
					eLong := &e.bTable[nextHashL]
					eLong.Cur, eLong.Prev = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***, eLong.Cur
					break
				***REMOVED***
			***REMOVED***

			t = sCandidate.offset - e.cur
			if s-t < maxMatchOffset && uint32(cv) == load3232(src, sCandidate.offset-e.cur) ***REMOVED***
				// Found a 4 match...
				l = e.matchlen(s+4, t+4, src) + 4
				lCandidate = e.bTable[nextHashL]
				// Store the next match

				e.table[nextHashS] = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***
				eLong := &e.bTable[nextHashL]
				eLong.Cur, eLong.Prev = tableEntry***REMOVED***offset: nextS + e.cur***REMOVED***, eLong.Cur

				// If the next long is a candidate, use that...
				t2 := lCandidate.Cur.offset - e.cur
				if nextS-t2 < maxMatchOffset ***REMOVED***
					if load3232(src, lCandidate.Cur.offset-e.cur) == uint32(next) ***REMOVED***
						ml := e.matchlen(nextS+4, t2+4, src) + 4
						if ml > l ***REMOVED***
							t = t2
							s = nextS
							l = ml
							break
						***REMOVED***
					***REMOVED***
					// If the previous long is a candidate, use that...
					t2 = lCandidate.Prev.offset - e.cur
					if nextS-t2 < maxMatchOffset && load3232(src, lCandidate.Prev.offset-e.cur) == uint32(next) ***REMOVED***
						ml := e.matchlen(nextS+4, t2+4, src) + 4
						if ml > l ***REMOVED***
							t = t2
							s = nextS
							l = ml
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***
			cv = next
		***REMOVED***

		// A 4-byte match has been found. We'll later see if more than 4 bytes
		// match. But, prior to the match, src[nextEmit:s] are unmatched. Emit
		// them as literal bytes.

		if l == 0 ***REMOVED***
			// Extend the 4-byte match as long as possible.
			l = e.matchlenLong(s+4, t+4, src) + 4
		***REMOVED*** else if l == maxMatchLength ***REMOVED***
			l += e.matchlenLong(s+l, t+l, src)
		***REMOVED***

		// Try to locate a better match by checking the end of best match...
		if sAt := s + l; l < 30 && sAt < sLimit ***REMOVED***
			eLong := e.bTable[hash7(load6432(src, sAt), tableBits)].Cur.offset
			// Test current
			t2 := eLong - e.cur - l
			off := s - t2
			if t2 >= 0 && off < maxMatchOffset && off > 0 ***REMOVED***
				if l2 := e.matchlenLong(s, t2, src); l2 > l ***REMOVED***
					t = t2
					l = l2
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
			emitLiteral(dst, src[nextEmit:s])
		***REMOVED***
		if debugDeflate ***REMOVED***
			if t >= s ***REMOVED***
				panic(fmt.Sprintln("s-t", s, t))
			***REMOVED***
			if (s - t) > maxMatchOffset ***REMOVED***
				panic(fmt.Sprintln("mmo", s-t))
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
			goto emitRemainder
		***REMOVED***

		// Store every 3rd hash in-between.
		if true ***REMOVED***
			const hashEvery = 3
			i := s - l + 1
			if i < s-1 ***REMOVED***
				cv := load6432(src, i)
				t := tableEntry***REMOVED***offset: i + e.cur***REMOVED***
				e.table[hash4x64(cv, tableBits)] = t
				eLong := &e.bTable[hash7(cv, tableBits)]
				eLong.Cur, eLong.Prev = t, eLong.Cur

				// Do an long at i+1
				cv >>= 8
				t = tableEntry***REMOVED***offset: t.offset + 1***REMOVED***
				eLong = &e.bTable[hash7(cv, tableBits)]
				eLong.Cur, eLong.Prev = t, eLong.Cur

				// We only have enough bits for a short entry at i+2
				cv >>= 8
				t = tableEntry***REMOVED***offset: t.offset + 1***REMOVED***
				e.table[hash4x64(cv, tableBits)] = t

				// Skip one - otherwise we risk hitting 's'
				i += 4
				for ; i < s-1; i += hashEvery ***REMOVED***
					cv := load6432(src, i)
					t := tableEntry***REMOVED***offset: i + e.cur***REMOVED***
					t2 := tableEntry***REMOVED***offset: t.offset + 1***REMOVED***
					eLong := &e.bTable[hash7(cv, tableBits)]
					eLong.Cur, eLong.Prev = t, eLong.Cur
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
		eLong := &e.bTable[prevHashL]
		eLong.Cur, eLong.Prev = tableEntry***REMOVED***offset: o***REMOVED***, eLong.Cur
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
