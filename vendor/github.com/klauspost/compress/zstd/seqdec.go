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
	hist         []byte
	literals     []byte
	out          []byte
	maxBits      uint8
***REMOVED***

// initialize all 3 decoders from the stream input.
func (s *sequenceDecs) initialize(br *bitReader, hist *history, literals, out []byte) error ***REMOVED***
	if err := s.litLengths.init(br); err != nil ***REMOVED***
		return errors.New("litLengths:" + err.Error())
	***REMOVED***
	if err := s.offsets.init(br); err != nil ***REMOVED***
		return errors.New("offsets:" + err.Error())
	***REMOVED***
	if err := s.matchLengths.init(br); err != nil ***REMOVED***
		return errors.New("matchLengths:" + err.Error())
	***REMOVED***
	s.literals = literals
	s.hist = hist.b
	s.prevOffset = hist.recentOffsets
	s.maxBits = s.litLengths.fse.maxBits + s.offsets.fse.maxBits + s.matchLengths.fse.maxBits
	s.out = out
	return nil
***REMOVED***

// decode sequences from the stream with the provided history.
func (s *sequenceDecs) decode(seqs int, br *bitReader, hist []byte) error ***REMOVED***
	startSize := len(s.out)
	for i := seqs - 1; i >= 0; i-- ***REMOVED***
		if br.overread() ***REMOVED***
			printf("reading sequence %d, exceeded available data\n", seqs-i)
			return io.ErrUnexpectedEOF
		***REMOVED***
		var litLen, matchOff, matchLen int
		if br.off > 4+((maxOffsetBits+16+16)>>3) ***REMOVED***
			litLen, matchOff, matchLen = s.nextFast(br)
			br.fillFast()
		***REMOVED*** else ***REMOVED***
			litLen, matchOff, matchLen = s.next(br)
			br.fill()
		***REMOVED***

		if debugSequences ***REMOVED***
			println("Seq", seqs-i-1, "Litlen:", litLen, "matchOff:", matchOff, "(abs) matchLen:", matchLen)
		***REMOVED***

		if litLen > len(s.literals) ***REMOVED***
			return fmt.Errorf("unexpected literal count, want %d bytes, but only %d is available", litLen, len(s.literals))
		***REMOVED***
		size := litLen + matchLen + len(s.out)
		if size-startSize > maxBlockSize ***REMOVED***
			return fmt.Errorf("output (%d) bigger than max block size", size)
		***REMOVED***
		if size > cap(s.out) ***REMOVED***
			// Not enough size, will be extremely rarely triggered,
			// but could be if destination slice is too small for sync operations.
			// We add maxBlockSize to the capacity.
			s.out = append(s.out, make([]byte, maxBlockSize)...)
			s.out = s.out[:len(s.out)-maxBlockSize]
		***REMOVED***
		if matchLen > maxMatchLen ***REMOVED***
			return fmt.Errorf("match len (%d) bigger than max allowed length", matchLen)
		***REMOVED***
		if matchOff > len(s.out)+len(hist)+litLen ***REMOVED***
			return fmt.Errorf("match offset (%d) bigger than current history (%d)", matchOff, len(s.out)+len(hist)+litLen)
		***REMOVED***
		if matchOff == 0 && matchLen > 0 ***REMOVED***
			return fmt.Errorf("zero matchoff and matchlen > 0")
		***REMOVED***

		s.out = append(s.out, s.literals[:litLen]...)
		s.literals = s.literals[litLen:]
		out := s.out

		// Copy from history.
		// TODO: Blocks without history could be made to ignore this completely.
		if v := matchOff - len(s.out); v > 0 ***REMOVED***
			// v is the start position in history from end.
			start := len(s.hist) - v
			if matchLen > v ***REMOVED***
				// Some goes into current block.
				// Copy remainder of history
				out = append(out, s.hist[start:]...)
				matchOff -= v
				matchLen -= v
			***REMOVED*** else ***REMOVED***
				out = append(out, s.hist[start:start+matchLen]...)
				matchLen = 0
			***REMOVED***
		***REMOVED***
		// We must be in current buffer now
		if matchLen > 0 ***REMOVED***
			start := len(s.out) - matchOff
			if matchLen <= len(s.out)-start ***REMOVED***
				// No overlap
				out = append(out, s.out[start:start+matchLen]...)
			***REMOVED*** else ***REMOVED***
				// Overlapping copy
				// Extend destination slice and copy one byte at the time.
				out = out[:len(out)+matchLen]
				src := out[start : start+matchLen]
				// Destination is the space we just added.
				dst := out[len(out)-matchLen:]
				dst = dst[:len(src)]
				for i := range src ***REMOVED***
					dst[i] = src[i]
				***REMOVED***
			***REMOVED***
		***REMOVED***
		s.out = out
		if i == 0 ***REMOVED***
			// This is the last sequence, so we shouldn't update state.
			break
		***REMOVED***
		if true ***REMOVED***
			// Manually inlined, ~ 5-20% faster
			// Update all 3 states at once. Approx 20% faster.
			a, b, c := s.litLengths.state.state, s.matchLengths.state.state, s.offsets.state.state

			nBits := a.nbBits + b.nbBits + c.nbBits
			if nBits == 0 ***REMOVED***
				s.litLengths.state.state = s.litLengths.state.dt[a.newState]
				s.matchLengths.state.state = s.matchLengths.state.dt[b.newState]
				s.offsets.state.state = s.offsets.state.dt[c.newState]
			***REMOVED*** else ***REMOVED***
				bits := br.getBitsFast(nBits)
				lowBits := uint16(bits >> ((c.nbBits + b.nbBits) & 31))
				s.litLengths.state.state = s.litLengths.state.dt[a.newState+lowBits]

				lowBits = uint16(bits >> (c.nbBits & 31))
				lowBits &= bitMask[b.nbBits&15]
				s.matchLengths.state.state = s.matchLengths.state.dt[b.newState+lowBits]

				lowBits = uint16(bits) & bitMask[c.nbBits&15]
				s.offsets.state.state = s.offsets.state.dt[c.newState+lowBits]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			s.updateAlt(br)
		***REMOVED***
	***REMOVED***

	// Add final literals
	s.out = append(s.out, s.literals...)
	return nil
***REMOVED***

// update states, at least 27 bits must be available.
func (s *sequenceDecs) update(br *bitReader) ***REMOVED***
	// Max 8 bits
	s.litLengths.state.next(br)
	// Max 9 bits
	s.matchLengths.state.next(br)
	// Max 8 bits
	s.offsets.state.next(br)
***REMOVED***

var bitMask [16]uint16

func init() ***REMOVED***
	for i := range bitMask[:] ***REMOVED***
		bitMask[i] = uint16((1 << uint(i)) - 1)
	***REMOVED***
***REMOVED***

// update states, at least 27 bits must be available.
func (s *sequenceDecs) updateAlt(br *bitReader) ***REMOVED***
	// Update all 3 states at once. Approx 20% faster.
	a, b, c := s.litLengths.state.state, s.matchLengths.state.state, s.offsets.state.state

	nBits := a.nbBits + b.nbBits + c.nbBits
	if nBits == 0 ***REMOVED***
		s.litLengths.state.state = s.litLengths.state.dt[a.newState]
		s.matchLengths.state.state = s.matchLengths.state.dt[b.newState]
		s.offsets.state.state = s.offsets.state.dt[c.newState]
		return
	***REMOVED***
	bits := br.getBitsFast(nBits)
	lowBits := uint16(bits >> ((c.nbBits + b.nbBits) & 31))
	s.litLengths.state.state = s.litLengths.state.dt[a.newState+lowBits]

	lowBits = uint16(bits >> (c.nbBits & 31))
	lowBits &= bitMask[b.nbBits&15]
	s.matchLengths.state.state = s.matchLengths.state.dt[b.newState+lowBits]

	lowBits = uint16(bits) & bitMask[c.nbBits&15]
	s.offsets.state.state = s.offsets.state.dt[c.newState+lowBits]
***REMOVED***

// nextFast will return new states when there are at least 4 unused bytes left on the stream when done.
func (s *sequenceDecs) nextFast(br *bitReader) (ll, mo, ml int) ***REMOVED***
	// Final will not read from stream.
	ll, llB := s.litLengths.state.final()
	ml, mlB := s.matchLengths.state.final()
	mo, moB := s.offsets.state.final()

	// extra bits are stored in reverse order.
	br.fillFast()
	if s.maxBits <= 32 ***REMOVED***
		mo += br.getBits(moB)
		ml += br.getBits(mlB)
		ll += br.getBits(llB)
	***REMOVED*** else ***REMOVED***
		mo += br.getBits(moB)
		br.fillFast()
		// matchlength+literal length, max 32 bits
		ml += br.getBits(mlB)
		ll += br.getBits(llB)
	***REMOVED***

	// mo = s.adjustOffset(mo, ll, moB)
	// Inlined for rather big speedup
	if moB > 1 ***REMOVED***
		s.prevOffset[2] = s.prevOffset[1]
		s.prevOffset[1] = s.prevOffset[0]
		s.prevOffset[0] = mo
		return
	***REMOVED***

	if ll == 0 ***REMOVED***
		// There is an exception though, when current sequence's literals_length = 0.
		// In this case, repeated offsets are shifted by one, so an offset_value of 1 means Repeated_Offset2,
		// an offset_value of 2 means Repeated_Offset3, and an offset_value of 3 means Repeated_Offset1 - 1_byte.
		mo++
	***REMOVED***

	if mo == 0 ***REMOVED***
		mo = s.prevOffset[0]
		return
	***REMOVED***
	var temp int
	if mo == 3 ***REMOVED***
		temp = s.prevOffset[0] - 1
	***REMOVED*** else ***REMOVED***
		temp = s.prevOffset[mo]
	***REMOVED***

	if temp == 0 ***REMOVED***
		// 0 is not valid; input is corrupted; force offset to 1
		println("temp was 0")
		temp = 1
	***REMOVED***

	if mo != 1 ***REMOVED***
		s.prevOffset[2] = s.prevOffset[1]
	***REMOVED***
	s.prevOffset[1] = s.prevOffset[0]
	s.prevOffset[0] = temp
	mo = temp
	return
***REMOVED***

func (s *sequenceDecs) next(br *bitReader) (ll, mo, ml int) ***REMOVED***
	// Final will not read from stream.
	ll, llB := s.litLengths.state.final()
	ml, mlB := s.matchLengths.state.final()
	mo, moB := s.offsets.state.final()

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

// mergeHistory will merge history.
func (s *sequenceDecs) mergeHistory(hist *sequenceDecs) (*sequenceDecs, error) ***REMOVED***
	for i := uint(0); i < 3; i++ ***REMOVED***
		var sNew, sHist *sequenceDec
		switch i ***REMOVED***
		default:
			// same as "case 0":
			sNew = &s.litLengths
			sHist = &hist.litLengths
		case 1:
			sNew = &s.offsets
			sHist = &hist.offsets
		case 2:
			sNew = &s.matchLengths
			sHist = &hist.matchLengths
		***REMOVED***
		if sNew.repeat ***REMOVED***
			if sHist.fse == nil ***REMOVED***
				return nil, fmt.Errorf("sequence stream %d, repeat requested, but no history", i)
			***REMOVED***
			continue
		***REMOVED***
		if sNew.fse == nil ***REMOVED***
			return nil, fmt.Errorf("sequence stream %d, no fse found", i)
		***REMOVED***
		if sHist.fse != nil && !sHist.fse.preDefined ***REMOVED***
			fseDecoderPool.Put(sHist.fse)
		***REMOVED***
		sHist.fse = sNew.fse
	***REMOVED***
	return hist, nil
***REMOVED***
