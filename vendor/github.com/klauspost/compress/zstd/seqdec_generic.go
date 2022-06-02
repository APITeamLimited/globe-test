//go:build !amd64 || appengine || !gc || noasm
// +build !amd64 appengine !gc noasm

package zstd

import (
	"fmt"
	"io"
)

// decode sequences from the stream with the provided history but without dictionary.
func (s *sequenceDecs) decodeSyncSimple(hist []byte) (bool, error) ***REMOVED***
	return false, nil
***REMOVED***

// decode sequences from the stream without the provided history.
func (s *sequenceDecs) decode(seqs []seqVals) error ***REMOVED***
	br := s.br

	// Grab full sizes tables, to avoid bounds checks.
	llTable, mlTable, ofTable := s.litLengths.fse.dt[:maxTablesize], s.matchLengths.fse.dt[:maxTablesize], s.offsets.fse.dt[:maxTablesize]
	llState, mlState, ofState := s.litLengths.state.state, s.matchLengths.state.state, s.offsets.state.state
	s.seqSize = 0
	litRemain := len(s.literals)

	maxBlockSize := maxCompressedBlockSize
	if s.windowSize < maxBlockSize ***REMOVED***
		maxBlockSize = s.windowSize
	***REMOVED***
	for i := range seqs ***REMOVED***
		var ll, mo, ml int
		if br.off > 4+((maxOffsetBits+16+16)>>3) ***REMOVED***
			// inlined function:
			// ll, mo, ml = s.nextFast(br, llState, mlState, ofState)

			// Final will not read from stream.
			var llB, mlB, moB uint8
			ll, llB = llState.final()
			ml, mlB = mlState.final()
			mo, moB = ofState.final()

			// extra bits are stored in reverse order.
			br.fillFast()
			mo += br.getBits(moB)
			if s.maxBits > 32 ***REMOVED***
				br.fillFast()
			***REMOVED***
			ml += br.getBits(mlB)
			ll += br.getBits(llB)

			if moB > 1 ***REMOVED***
				s.prevOffset[2] = s.prevOffset[1]
				s.prevOffset[1] = s.prevOffset[0]
				s.prevOffset[0] = mo
			***REMOVED*** else ***REMOVED***
				// mo = s.adjustOffset(mo, ll, moB)
				// Inlined for rather big speedup
				if ll == 0 ***REMOVED***
					// There is an exception though, when current sequence's literals_length = 0.
					// In this case, repeated offsets are shifted by one, so an offset_value of 1 means Repeated_Offset2,
					// an offset_value of 2 means Repeated_Offset3, and an offset_value of 3 means Repeated_Offset1 - 1_byte.
					mo++
				***REMOVED***

				if mo == 0 ***REMOVED***
					mo = s.prevOffset[0]
				***REMOVED*** else ***REMOVED***
					var temp int
					if mo == 3 ***REMOVED***
						temp = s.prevOffset[0] - 1
					***REMOVED*** else ***REMOVED***
						temp = s.prevOffset[mo]
					***REMOVED***

					if temp == 0 ***REMOVED***
						// 0 is not valid; input is corrupted; force offset to 1
						println("WARNING: temp was 0")
						temp = 1
					***REMOVED***

					if mo != 1 ***REMOVED***
						s.prevOffset[2] = s.prevOffset[1]
					***REMOVED***
					s.prevOffset[1] = s.prevOffset[0]
					s.prevOffset[0] = temp
					mo = temp
				***REMOVED***
			***REMOVED***
			br.fillFast()
		***REMOVED*** else ***REMOVED***
			if br.overread() ***REMOVED***
				if debugDecoder ***REMOVED***
					printf("reading sequence %d, exceeded available data\n", i)
				***REMOVED***
				return io.ErrUnexpectedEOF
			***REMOVED***
			ll, mo, ml = s.next(br, llState, mlState, ofState)
			br.fill()
		***REMOVED***

		if debugSequences ***REMOVED***
			println("Seq", i, "Litlen:", ll, "mo:", mo, "(abs) ml:", ml)
		***REMOVED***
		// Evaluate.
		// We might be doing this async, so do it early.
		if mo == 0 && ml > 0 ***REMOVED***
			return fmt.Errorf("zero matchoff and matchlen (%d) > 0", ml)
		***REMOVED***
		if ml > maxMatchLen ***REMOVED***
			return fmt.Errorf("match len (%d) bigger than max allowed length", ml)
		***REMOVED***
		s.seqSize += ll + ml
		if s.seqSize > maxBlockSize ***REMOVED***
			return fmt.Errorf("output (%d) bigger than max block size (%d)", s.seqSize, maxBlockSize)
		***REMOVED***
		litRemain -= ll
		if litRemain < 0 ***REMOVED***
			return fmt.Errorf("unexpected literal count, want %d bytes, but only %d is available", ll, litRemain+ll)
		***REMOVED***
		seqs[i] = seqVals***REMOVED***
			ll: ll,
			ml: ml,
			mo: mo,
		***REMOVED***
		if i == len(seqs)-1 ***REMOVED***
			// This is the last sequence, so we shouldn't update state.
			break
		***REMOVED***

		// Manually inlined, ~ 5-20% faster
		// Update all 3 states at once. Approx 20% faster.
		nBits := llState.nbBits() + mlState.nbBits() + ofState.nbBits()
		if nBits == 0 ***REMOVED***
			llState = llTable[llState.newState()&maxTableMask]
			mlState = mlTable[mlState.newState()&maxTableMask]
			ofState = ofTable[ofState.newState()&maxTableMask]
		***REMOVED*** else ***REMOVED***
			bits := br.get32BitsFast(nBits)
			lowBits := uint16(bits >> ((ofState.nbBits() + mlState.nbBits()) & 31))
			llState = llTable[(llState.newState()+lowBits)&maxTableMask]

			lowBits = uint16(bits >> (ofState.nbBits() & 31))
			lowBits &= bitMask[mlState.nbBits()&15]
			mlState = mlTable[(mlState.newState()+lowBits)&maxTableMask]

			lowBits = uint16(bits) & bitMask[ofState.nbBits()&15]
			ofState = ofTable[(ofState.newState()+lowBits)&maxTableMask]
		***REMOVED***
	***REMOVED***
	s.seqSize += litRemain
	if s.seqSize > maxBlockSize ***REMOVED***
		return fmt.Errorf("output (%d) bigger than max block size (%d)", s.seqSize, maxBlockSize)
	***REMOVED***
	err := br.close()
	if err != nil ***REMOVED***
		printf("Closing sequences: %v, %+v\n", err, *br)
	***REMOVED***
	return err
***REMOVED***

// executeSimple handles cases when a dictionary is not used.
func (s *sequenceDecs) executeSimple(seqs []seqVals, hist []byte) error ***REMOVED***
	// Ensure we have enough output size...
	if len(s.out)+s.seqSize > cap(s.out) ***REMOVED***
		addBytes := s.seqSize + len(s.out)
		s.out = append(s.out, make([]byte, addBytes)...)
		s.out = s.out[:len(s.out)-addBytes]
	***REMOVED***

	if debugDecoder ***REMOVED***
		printf("Execute %d seqs with literals: %d into %d bytes\n", len(seqs), len(s.literals), s.seqSize)
	***REMOVED***

	var t = len(s.out)
	out := s.out[:t+s.seqSize]

	for _, seq := range seqs ***REMOVED***
		// Add literals
		copy(out[t:], s.literals[:seq.ll])
		t += seq.ll
		s.literals = s.literals[seq.ll:]

		// Malformed input
		if seq.mo > t+len(hist) || seq.mo > s.windowSize ***REMOVED***
			return fmt.Errorf("match offset (%d) bigger than current history (%d)", seq.mo, t+len(hist))
		***REMOVED***

		// Copy from history.
		if v := seq.mo - t; v > 0 ***REMOVED***
			// v is the start position in history from end.
			start := len(hist) - v
			if seq.ml > v ***REMOVED***
				// Some goes into the current block.
				// Copy remainder of history
				copy(out[t:], hist[start:])
				t += v
				seq.ml -= v
			***REMOVED*** else ***REMOVED***
				copy(out[t:], hist[start:start+seq.ml])
				t += seq.ml
				continue
			***REMOVED***
		***REMOVED***

		// We must be in the current buffer now
		if seq.ml > 0 ***REMOVED***
			start := t - seq.mo
			if seq.ml <= t-start ***REMOVED***
				// No overlap
				copy(out[t:], out[start:start+seq.ml])
				t += seq.ml
			***REMOVED*** else ***REMOVED***
				// Overlapping copy
				// Extend destination slice and copy one byte at the time.
				src := out[start : start+seq.ml]
				dst := out[t:]
				dst = dst[:len(src)]
				t += len(src)
				// Destination is the space we just added.
				for i := range src ***REMOVED***
					dst[i] = src[i]
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Add final literals
	copy(out[t:], s.literals)
	if debugDecoder ***REMOVED***
		t += len(s.literals)
		if t != len(out) ***REMOVED***
			panic(fmt.Errorf("length mismatch, want %d, got %d, ss: %d", len(out), t, s.seqSize))
		***REMOVED***
	***REMOVED***
	s.out = out

	return nil
***REMOVED***
