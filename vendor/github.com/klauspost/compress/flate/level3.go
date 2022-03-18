package flate

import "fmt"

// fastEncL3
type fastEncL3 struct ***REMOVED***
	fastGen
	table [1 << 16]tableEntryPrev
***REMOVED***

// Encode uses a similar algorithm to level 2, will check up to two candidates.
func (e *fastEncL3) Encode(dst *tokens, src []byte) ***REMOVED***
	const (
		inputMargin            = 8 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
		tableBits              = 16
		tableSize              = 1 << tableBits
	)

	if debugDeflate && e.cur < 0 ***REMOVED***
		panic(fmt.Sprint("e.cur < 0: ", e.cur))
	***REMOVED***

	// Protect against e.cur wraparound.
	for e.cur >= bufferReset ***REMOVED***
		if len(e.hist) == 0 ***REMOVED***
			for i := range e.table[:] ***REMOVED***
				e.table[i] = tableEntryPrev***REMOVED******REMOVED***
			***REMOVED***
			e.cur = maxMatchOffset
			break
		***REMOVED***
		// Shift down everything in the table that isn't already too far away.
		minOff := e.cur + int32(len(e.hist)) - maxMatchOffset
		for i := range e.table[:] ***REMOVED***
			v := e.table[i]
			if v.Cur.offset <= minOff ***REMOVED***
				v.Cur.offset = 0
			***REMOVED*** else ***REMOVED***
				v.Cur.offset = v.Cur.offset - e.cur + maxMatchOffset
			***REMOVED***
			if v.Prev.offset <= minOff ***REMOVED***
				v.Prev.offset = 0
			***REMOVED*** else ***REMOVED***
				v.Prev.offset = v.Prev.offset - e.cur + maxMatchOffset
			***REMOVED***
			e.table[i] = v
		***REMOVED***
		e.cur = maxMatchOffset
	***REMOVED***

	s := e.addBlock(src)

	// Skip if too small.
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
		const skipLog = 6
		nextS := s
		var candidate tableEntry
		for ***REMOVED***
			nextHash := hash4u(cv, tableBits)
			s = nextS
			nextS = s + 1 + (s-nextEmit)>>skipLog
			if nextS > sLimit ***REMOVED***
				goto emitRemainder
			***REMOVED***
			candidates := e.table[nextHash]
			now := load3232(src, nextS)

			// Safe offset distance until s + 4...
			minOffset := e.cur + s - (maxMatchOffset - 4)
			e.table[nextHash] = tableEntryPrev***REMOVED***Prev: candidates.Cur, Cur: tableEntry***REMOVED***offset: s + e.cur***REMOVED******REMOVED***

			// Check both candidates
			candidate = candidates.Cur
			if candidate.offset < minOffset ***REMOVED***
				cv = now
				// Previous will also be invalid, we have nothing.
				continue
			***REMOVED***

			if cv == load3232(src, candidate.offset-e.cur) ***REMOVED***
				if candidates.Prev.offset < minOffset || cv != load3232(src, candidates.Prev.offset-e.cur) ***REMOVED***
					break
				***REMOVED***
				// Both match and are valid, pick longest.
				offset := s - (candidate.offset - e.cur)
				o2 := s - (candidates.Prev.offset - e.cur)
				l1, l2 := matchLen(src[s+4:], src[s-offset+4:]), matchLen(src[s+4:], src[s-o2+4:])
				if l2 > l1 ***REMOVED***
					candidate = candidates.Prev
				***REMOVED***
				break
			***REMOVED*** else ***REMOVED***
				// We only check if value mismatches.
				// Offset will always be invalid in other cases.
				candidate = candidates.Prev
				if candidate.offset > minOffset && cv == load3232(src, candidate.offset-e.cur) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			cv = now
		***REMOVED***

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
			//
			t := candidate.offset - e.cur
			l := e.matchlenLong(s+4, t+4, src) + 4

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

			dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
			s += l
			nextEmit = s
			if nextS >= s ***REMOVED***
				s = nextS + 1
			***REMOVED***

			if s >= sLimit ***REMOVED***
				t += l
				// Index first pair after match end.
				if int(t+4) < len(src) && t > 0 ***REMOVED***
					cv := load3232(src, t)
					nextHash := hash4u(cv, tableBits)
					e.table[nextHash] = tableEntryPrev***REMOVED***
						Prev: e.table[nextHash].Cur,
						Cur:  tableEntry***REMOVED***offset: e.cur + t***REMOVED***,
					***REMOVED***
				***REMOVED***
				goto emitRemainder
			***REMOVED***

			// Store every 5th hash in-between.
			for i := s - l + 2; i < s-5; i += 5 ***REMOVED***
				nextHash := hash4u(load3232(src, i), tableBits)
				e.table[nextHash] = tableEntryPrev***REMOVED***
					Prev: e.table[nextHash].Cur,
					Cur:  tableEntry***REMOVED***offset: e.cur + i***REMOVED******REMOVED***
			***REMOVED***
			// We could immediately start working at s now, but to improve
			// compression we first update the hash table at s-2 to s.
			x := load6432(src, s-2)
			prevHash := hash4u(uint32(x), tableBits)

			e.table[prevHash] = tableEntryPrev***REMOVED***
				Prev: e.table[prevHash].Cur,
				Cur:  tableEntry***REMOVED***offset: e.cur + s - 2***REMOVED***,
			***REMOVED***
			x >>= 8
			prevHash = hash4u(uint32(x), tableBits)

			e.table[prevHash] = tableEntryPrev***REMOVED***
				Prev: e.table[prevHash].Cur,
				Cur:  tableEntry***REMOVED***offset: e.cur + s - 1***REMOVED***,
			***REMOVED***
			x >>= 8
			currHash := hash4u(uint32(x), tableBits)
			candidates := e.table[currHash]
			cv = uint32(x)
			e.table[currHash] = tableEntryPrev***REMOVED***
				Prev: candidates.Cur,
				Cur:  tableEntry***REMOVED***offset: s + e.cur***REMOVED***,
			***REMOVED***

			// Check both candidates
			candidate = candidates.Cur
			minOffset := e.cur + s - (maxMatchOffset - 4)

			if candidate.offset > minOffset ***REMOVED***
				if cv == load3232(src, candidate.offset-e.cur) ***REMOVED***
					// Found a match...
					continue
				***REMOVED***
				candidate = candidates.Prev
				if candidate.offset > minOffset && cv == load3232(src, candidate.offset-e.cur) ***REMOVED***
					// Match at prev...
					continue
				***REMOVED***
			***REMOVED***
			cv = uint32(x >> 8)
			s++
			break
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
