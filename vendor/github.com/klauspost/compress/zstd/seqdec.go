// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"fmt"
	"io"
)

type seq struct ***REMOVED***
	litLen   uint32
	matchLen uint32
	offset   uint32

	// Codes are stored here for the encoder
	// so they only have to be looked up once.
	llCode, mlCode, ofCode uint8
***REMOVED***

type seqVals struct ***REMOVED***
	ll, ml, mo int
***REMOVED***

func (s seq) String() string ***REMOVED***
	if s.offset <= 3 ***REMOVED***
		if s.offset == 0 ***REMOVED***
			return fmt.Sprint("litLen:", s.litLen, ", matchLen:", s.matchLen+zstdMinMatch, ", offset: INVALID (0)")
		***REMOVED***
		return fmt.Sprint("litLen:", s.litLen, ", matchLen:", s.matchLen+zstdMinMatch, ", offset:", s.offset, " (repeat)")
	***REMOVED***
	return fmt.Sprint("litLen:", s.litLen, ", matchLen:", s.matchLen+zstdMinMatch, ", offset:", s.offset-3, " (new)")
***REMOVED***

type seqCompMode uint8

const (
	compModePredefined seqCompMode = iota
	compModeRLE
	compModeFSE
	compModeRepeat
)

type sequenceDec struct ***REMOVED***
	// decoder keeps track of the current state and updates it from the bitstream.
	fse    *fseDecoder
	state  fseState
	repeat bool
***REMOVED***

// init the state of the decoder with input from stream.
func (s *sequenceDec) init(br *bitReader) error ***REMOVED***
	if s.fse == nil ***REMOVED***
		return errors.New("sequence decoder not defined")
	***REMOVED***
	s.state.init(br, s.fse.actualTableLog, s.fse.dt[:1<<s.fse.actualTableLog])
	return nil
***REMOVED***

// sequenceDecs contains all 3 sequence decoders and their state.
type sequenceDecs struct ***REMOVED***
	litLengths   sequenceDec
	offsets      sequenceDec
	matchLengths sequenceDec
	prevOffset   [3]int
	dict         []byte
	literals     []byte
	out          []byte
	nSeqs        int
	br           *bitReader
	seqSize      int
	windowSize   int
	maxBits      uint8
	maxSyncLen   uint64
***REMOVED***

// initialize all 3 decoders from the stream input.
func (s *sequenceDecs) initialize(br *bitReader, hist *history, out []byte) error ***REMOVED***
	if err := s.litLengths.init(br); err != nil ***REMOVED***
		return errors.New("litLengths:" + err.Error())
	***REMOVED***
	if err := s.offsets.init(br); err != nil ***REMOVED***
		return errors.New("offsets:" + err.Error())
	***REMOVED***
	if err := s.matchLengths.init(br); err != nil ***REMOVED***
		return errors.New("matchLengths:" + err.Error())
	***REMOVED***
	s.br = br
	s.prevOffset = hist.recentOffsets
	s.maxBits = s.litLengths.fse.maxBits + s.offsets.fse.maxBits + s.matchLengths.fse.maxBits
	s.windowSize = hist.windowSize
	s.out = out
	s.dict = nil
	if hist.dict != nil ***REMOVED***
		s.dict = hist.dict.content
	***REMOVED***
	return nil
***REMOVED***

// execute will execute the decoded sequence with the provided history.
// The sequence must be evaluated before being sent.
func (s *sequenceDecs) execute(seqs []seqVals, hist []byte) error ***REMOVED***
	if len(s.dict) == 0 ***REMOVED***
		return s.executeSimple(seqs, hist)
	***REMOVED***

	// Ensure we have enough output size...
	if len(s.out)+s.seqSize > cap(s.out) ***REMOVED***
		addBytes := s.seqSize + len(s.out)
		s.out = append(s.out, make([]byte, addBytes)...)
		s.out = s.out[:len(s.out)-addBytes]
	***REMOVED***

	if debugDecoder ***REMOVED***
		printf("Execute %d seqs with hist %d, dict %d, literals: %d into %d bytes\n", len(seqs), len(hist), len(s.dict), len(s.literals), s.seqSize)
	***REMOVED***

	var t = len(s.out)
	out := s.out[:t+s.seqSize]

	for _, seq := range seqs ***REMOVED***
		// Add literals
		copy(out[t:], s.literals[:seq.ll])
		t += seq.ll
		s.literals = s.literals[seq.ll:]

		// Copy from dictionary...
		if seq.mo > t+len(hist) || seq.mo > s.windowSize ***REMOVED***
			if len(s.dict) == 0 ***REMOVED***
				return fmt.Errorf("match offset (%d) bigger than current history (%d)", seq.mo, t+len(hist))
			***REMOVED***

			// we may be in dictionary.
			dictO := len(s.dict) - (seq.mo - (t + len(hist)))
			if dictO < 0 || dictO >= len(s.dict) ***REMOVED***
				return fmt.Errorf("match offset (%d) bigger than current history+dict (%d)", seq.mo, t+len(hist)+len(s.dict))
			***REMOVED***
			end := dictO + seq.ml
			if end > len(s.dict) ***REMOVED***
				n := len(s.dict) - dictO
				copy(out[t:], s.dict[dictO:])
				t += n
				seq.ml -= n
			***REMOVED*** else ***REMOVED***
				copy(out[t:], s.dict[dictO:end])
				t += end - dictO
				continue
			***REMOVED***
		***REMOVED***

		// Copy from history.
		if v := seq.mo - t; v > 0 ***REMOVED***
			// v is the start position in history from end.
			start := len(hist) - v
			if seq.ml > v ***REMOVED***
				// Some goes into current block.
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
		// We must be in current buffer now
		if seq.ml > 0 ***REMOVED***
			start := t - seq.mo
			if seq.ml <= t-start ***REMOVED***
				// No overlap
				copy(out[t:], out[start:start+seq.ml])
				t += seq.ml
				continue
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

// decode sequences from the stream with the provided history.
func (s *sequenceDecs) decodeSync(hist []byte) error ***REMOVED***
	supported, err := s.decodeSyncSimple(hist)
	if supported ***REMOVED***
		return err
	***REMOVED***

	br := s.br
	seqs := s.nSeqs
	startSize := len(s.out)
	// Grab full sizes tables, to avoid bounds checks.
	llTable, mlTable, ofTable := s.litLengths.fse.dt[:maxTablesize], s.matchLengths.fse.dt[:maxTablesize], s.offsets.fse.dt[:maxTablesize]
	llState, mlState, ofState := s.litLengths.state.state, s.matchLengths.state.state, s.offsets.state.state
	out := s.out
	maxBlockSize := maxCompressedBlockSize
	if s.windowSize < maxBlockSize ***REMOVED***
		maxBlockSize = s.windowSize
	***REMOVED***

	for i := seqs - 1; i >= 0; i-- ***REMOVED***
		if br.overread() ***REMOVED***
			printf("reading sequence %d, exceeded available data\n", seqs-i)
			return io.ErrUnexpectedEOF
		***REMOVED***
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
			ll, mo, ml = s.next(br, llState, mlState, ofState)
			br.fill()
		***REMOVED***

		if debugSequences ***REMOVED***
			println("Seq", seqs-i-1, "Litlen:", ll, "mo:", mo, "(abs) ml:", ml)
		***REMOVED***

		if ll > len(s.literals) ***REMOVED***
			return fmt.Errorf("unexpected literal count, want %d bytes, but only %d is available", ll, len(s.literals))
		***REMOVED***
		size := ll + ml + len(out)
		if size-startSize > maxBlockSize ***REMOVED***
			return fmt.Errorf("output (%d) bigger than max block size (%d)", size-startSize, maxBlockSize)
		***REMOVED***
		if size > cap(out) ***REMOVED***
			// Not enough size, which can happen under high volume block streaming conditions
			// but could be if destination slice is too small for sync operations.
			// over-allocating here can create a large amount of GC pressure so we try to keep
			// it as contained as possible
			used := len(out) - startSize
			addBytes := 256 + ll + ml + used>>2
			// Clamp to max block size.
			if used+addBytes > maxBlockSize ***REMOVED***
				addBytes = maxBlockSize - used
			***REMOVED***
			out = append(out, make([]byte, addBytes)...)
			out = out[:len(out)-addBytes]
		***REMOVED***
		if ml > maxMatchLen ***REMOVED***
			return fmt.Errorf("match len (%d) bigger than max allowed length", ml)
		***REMOVED***

		// Add literals
		out = append(out, s.literals[:ll]...)
		s.literals = s.literals[ll:]

		if mo == 0 && ml > 0 ***REMOVED***
			return fmt.Errorf("zero matchoff and matchlen (%d) > 0", ml)
		***REMOVED***

		if mo > len(out)+len(hist) || mo > s.windowSize ***REMOVED***
			if len(s.dict) == 0 ***REMOVED***
				return fmt.Errorf("match offset (%d) bigger than current history (%d)", mo, len(out)+len(hist)-startSize)
			***REMOVED***

			// we may be in dictionary.
			dictO := len(s.dict) - (mo - (len(out) + len(hist)))
			if dictO < 0 || dictO >= len(s.dict) ***REMOVED***
				return fmt.Errorf("match offset (%d) bigger than current history (%d)", mo, len(out)+len(hist)-startSize)
			***REMOVED***
			end := dictO + ml
			if end > len(s.dict) ***REMOVED***
				out = append(out, s.dict[dictO:]...)
				ml -= len(s.dict) - dictO
			***REMOVED*** else ***REMOVED***
				out = append(out, s.dict[dictO:end]...)
				mo = 0
				ml = 0
			***REMOVED***
		***REMOVED***

		// Copy from history.
		// TODO: Blocks without history could be made to ignore this completely.
		if v := mo - len(out); v > 0 ***REMOVED***
			// v is the start position in history from end.
			start := len(hist) - v
			if ml > v ***REMOVED***
				// Some goes into current block.
				// Copy remainder of history
				out = append(out, hist[start:]...)
				ml -= v
			***REMOVED*** else ***REMOVED***
				out = append(out, hist[start:start+ml]...)
				ml = 0
			***REMOVED***
		***REMOVED***
		// We must be in current buffer now
		if ml > 0 ***REMOVED***
			start := len(out) - mo
			if ml <= len(out)-start ***REMOVED***
				// No overlap
				out = append(out, out[start:start+ml]...)
			***REMOVED*** else ***REMOVED***
				// Overlapping copy
				// Extend destination slice and copy one byte at the time.
				out = out[:len(out)+ml]
				src := out[start : start+ml]
				// Destination is the space we just added.
				dst := out[len(out)-ml:]
				dst = dst[:len(src)]
				for i := range src ***REMOVED***
					dst[i] = src[i]
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if i == 0 ***REMOVED***
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

	// Check if space for literals
	if size := len(s.literals) + len(s.out) - startSize; size > maxBlockSize ***REMOVED***
		return fmt.Errorf("output (%d) bigger than max block size (%d)", size, maxBlockSize)
	***REMOVED***

	// Add final literals
	s.out = append(out, s.literals...)
	return br.close()
***REMOVED***

var bitMask [16]uint16

func init() ***REMOVED***
	for i := range bitMask[:] ***REMOVED***
		bitMask[i] = uint16((1 << uint(i)) - 1)
	***REMOVED***
***REMOVED***

func (s *sequenceDecs) next(br *bitReader, llState, mlState, ofState decSymbol) (ll, mo, ml int) ***REMOVED***
	// Final will not read from stream.
	ll, llB := llState.final()
	ml, mlB := mlState.final()
	mo, moB := ofState.final()

	// extra bits are stored in reverse order.
	br.fill()
	if s.maxBits <= 32 ***REMOVED***
		mo += br.getBits(moB)
		ml += br.getBits(mlB)
		ll += br.getBits(llB)
	***REMOVED*** else ***REMOVED***
		mo += br.getBits(moB)
		br.fill()
		// matchlength+literal length, max 32 bits
		ml += br.getBits(mlB)
		ll += br.getBits(llB)

	***REMOVED***
	mo = s.adjustOffset(mo, ll, moB)
	return
***REMOVED***

func (s *sequenceDecs) adjustOffset(offset, litLen int, offsetB uint8) int ***REMOVED***
	if offsetB > 1 ***REMOVED***
		s.prevOffset[2] = s.prevOffset[1]
		s.prevOffset[1] = s.prevOffset[0]
		s.prevOffset[0] = offset
		return offset
	***REMOVED***

	if litLen == 0 ***REMOVED***
		// There is an exception though, when current sequence's literals_length = 0.
		// In this case, repeated offsets are shifted by one, so an offset_value of 1 means Repeated_Offset2,
		// an offset_value of 2 means Repeated_Offset3, and an offset_value of 3 means Repeated_Offset1 - 1_byte.
		offset++
	***REMOVED***

	if offset == 0 ***REMOVED***
		return s.prevOffset[0]
	***REMOVED***
	var temp int
	if offset == 3 ***REMOVED***
		temp = s.prevOffset[0] - 1
	***REMOVED*** else ***REMOVED***
		temp = s.prevOffset[offset]
	***REMOVED***

	if temp == 0 ***REMOVED***
		// 0 is not valid; input is corrupted; force offset to 1
		println("temp was 0")
		temp = 1
	***REMOVED***

	if offset != 1 ***REMOVED***
		s.prevOffset[2] = s.prevOffset[1]
	***REMOVED***
	s.prevOffset[1] = s.prevOffset[0]
	s.prevOffset[0] = temp
	return temp
***REMOVED***
