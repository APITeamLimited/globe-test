package fse

import (
	"errors"
	"fmt"
)

const (
	tablelogAbsoluteMax = 15
)

// Decompress a block of data.
// You can provide a scratch buffer to avoid allocations.
// If nil is provided a temporary one will be allocated.
// It is possible, but by no way guaranteed that corrupt data will
// return an error.
// It is up to the caller to verify integrity of the returned data.
// Use a predefined Scrach to set maximum acceptable output size.
func Decompress(b []byte, s *Scratch) ([]byte, error) ***REMOVED***
	s, err := s.prepare(b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s.Out = s.Out[:0]
	err = s.readNCount()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = s.buildDtable()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = s.decompress()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return s.Out, nil
***REMOVED***

// readNCount will read the symbol distribution so decoding tables can be constructed.
func (s *Scratch) readNCount() error ***REMOVED***
	var (
		charnum   uint16
		previous0 bool
		b         = &s.br
	)
	iend := b.remain()
	if iend < 4 ***REMOVED***
		return errors.New("input too small")
	***REMOVED***
	bitStream := b.Uint32()
	nbBits := uint((bitStream & 0xF) + minTablelog) // extract tableLog
	if nbBits > tablelogAbsoluteMax ***REMOVED***
		return errors.New("tableLog too large")
	***REMOVED***
	bitStream >>= 4
	bitCount := uint(4)

	s.actualTableLog = uint8(nbBits)
	remaining := int32((1 << nbBits) + 1)
	threshold := int32(1 << nbBits)
	gotTotal := int32(0)
	nbBits++

	for remaining > 1 ***REMOVED***
		if previous0 ***REMOVED***
			n0 := charnum
			for (bitStream & 0xFFFF) == 0xFFFF ***REMOVED***
				n0 += 24
				if b.off < iend-5 ***REMOVED***
					b.advance(2)
					bitStream = b.Uint32() >> bitCount
				***REMOVED*** else ***REMOVED***
					bitStream >>= 16
					bitCount += 16
				***REMOVED***
			***REMOVED***
			for (bitStream & 3) == 3 ***REMOVED***
				n0 += 3
				bitStream >>= 2
				bitCount += 2
			***REMOVED***
			n0 += uint16(bitStream & 3)
			bitCount += 2
			if n0 > maxSymbolValue ***REMOVED***
				return errors.New("maxSymbolValue too small")
			***REMOVED***
			for charnum < n0 ***REMOVED***
				s.norm[charnum&0xff] = 0
				charnum++
			***REMOVED***

			if b.off <= iend-7 || b.off+int(bitCount>>3) <= iend-4 ***REMOVED***
				b.advance(bitCount >> 3)
				bitCount &= 7
				bitStream = b.Uint32() >> bitCount
			***REMOVED*** else ***REMOVED***
				bitStream >>= 2
			***REMOVED***
		***REMOVED***

		max := (2*(threshold) - 1) - (remaining)
		var count int32

		if (int32(bitStream) & (threshold - 1)) < max ***REMOVED***
			count = int32(bitStream) & (threshold - 1)
			bitCount += nbBits - 1
		***REMOVED*** else ***REMOVED***
			count = int32(bitStream) & (2*threshold - 1)
			if count >= threshold ***REMOVED***
				count -= max
			***REMOVED***
			bitCount += nbBits
		***REMOVED***

		count-- // extra accuracy
		if count < 0 ***REMOVED***
			// -1 means +1
			remaining += count
			gotTotal -= count
		***REMOVED*** else ***REMOVED***
			remaining -= count
			gotTotal += count
		***REMOVED***
		s.norm[charnum&0xff] = int16(count)
		charnum++
		previous0 = count == 0
		for remaining < threshold ***REMOVED***
			nbBits--
			threshold >>= 1
		***REMOVED***
		if b.off <= iend-7 || b.off+int(bitCount>>3) <= iend-4 ***REMOVED***
			b.advance(bitCount >> 3)
			bitCount &= 7
		***REMOVED*** else ***REMOVED***
			bitCount -= (uint)(8 * (len(b.b) - 4 - b.off))
			b.off = len(b.b) - 4
		***REMOVED***
		bitStream = b.Uint32() >> (bitCount & 31)
	***REMOVED***
	s.symbolLen = charnum

	if s.symbolLen <= 1 ***REMOVED***
		return fmt.Errorf("symbolLen (%d) too small", s.symbolLen)
	***REMOVED***
	if s.symbolLen > maxSymbolValue+1 ***REMOVED***
		return fmt.Errorf("symbolLen (%d) too big", s.symbolLen)
	***REMOVED***
	if remaining != 1 ***REMOVED***
		return fmt.Errorf("corruption detected (remaining %d != 1)", remaining)
	***REMOVED***
	if bitCount > 32 ***REMOVED***
		return fmt.Errorf("corruption detected (bitCount %d > 32)", bitCount)
	***REMOVED***
	if gotTotal != 1<<s.actualTableLog ***REMOVED***
		return fmt.Errorf("corruption detected (total %d != %d)", gotTotal, 1<<s.actualTableLog)
	***REMOVED***
	b.advance((bitCount + 7) >> 3)
	return nil
***REMOVED***

// decSymbol contains information about a state entry,
// Including the state offset base, the output symbol and
// the number of bits to read for the low part of the destination state.
type decSymbol struct ***REMOVED***
	newState uint16
	symbol   uint8
	nbBits   uint8
***REMOVED***

// allocDtable will allocate decoding tables if they are not big enough.
func (s *Scratch) allocDtable() ***REMOVED***
	tableSize := 1 << s.actualTableLog
	if cap(s.decTable) < int(tableSize) ***REMOVED***
		s.decTable = make([]decSymbol, tableSize)
	***REMOVED***
	s.decTable = s.decTable[:tableSize]

	if cap(s.ct.tableSymbol) < 256 ***REMOVED***
		s.ct.tableSymbol = make([]byte, 256)
	***REMOVED***
	s.ct.tableSymbol = s.ct.tableSymbol[:256]

	if cap(s.ct.stateTable) < 256 ***REMOVED***
		s.ct.stateTable = make([]uint16, 256)
	***REMOVED***
	s.ct.stateTable = s.ct.stateTable[:256]
***REMOVED***

// buildDtable will build the decoding table.
func (s *Scratch) buildDtable() error ***REMOVED***
	tableSize := uint32(1 << s.actualTableLog)
	highThreshold := tableSize - 1
	s.allocDtable()
	symbolNext := s.ct.stateTable[:256]

	// Init, lay down lowprob symbols
	s.zeroBits = false
	***REMOVED***
		largeLimit := int16(1 << (s.actualTableLog - 1))
		for i, v := range s.norm[:s.symbolLen] ***REMOVED***
			if v == -1 ***REMOVED***
				s.decTable[highThreshold].symbol = uint8(i)
				highThreshold--
				symbolNext[i] = 1
			***REMOVED*** else ***REMOVED***
				if v >= largeLimit ***REMOVED***
					s.zeroBits = true
				***REMOVED***
				symbolNext[i] = uint16(v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Spread symbols
	***REMOVED***
		tableMask := tableSize - 1
		step := tableStep(tableSize)
		position := uint32(0)
		for ss, v := range s.norm[:s.symbolLen] ***REMOVED***
			for i := 0; i < int(v); i++ ***REMOVED***
				s.decTable[position].symbol = uint8(ss)
				position = (position + step) & tableMask
				for position > highThreshold ***REMOVED***
					// lowprob area
					position = (position + step) & tableMask
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if position != 0 ***REMOVED***
			// position must reach all cells once, otherwise normalizedCounter is incorrect
			return errors.New("corrupted input (position != 0)")
		***REMOVED***
	***REMOVED***

	// Build Decoding table
	***REMOVED***
		tableSize := uint16(1 << s.actualTableLog)
		for u, v := range s.decTable ***REMOVED***
			symbol := v.symbol
			nextState := symbolNext[symbol]
			symbolNext[symbol] = nextState + 1
			nBits := s.actualTableLog - byte(highBits(uint32(nextState)))
			s.decTable[u].nbBits = nBits
			newState := (nextState << nBits) - tableSize
			if newState > tableSize ***REMOVED***
				return fmt.Errorf("newState (%d) outside table size (%d)", newState, tableSize)
			***REMOVED***
			if newState == uint16(u) && nBits == 0 ***REMOVED***
				// Seems weird that this is possible with nbits > 0.
				return fmt.Errorf("newState (%d) == oldState (%d) and no bits", newState, u)
			***REMOVED***
			s.decTable[u].newState = newState
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// decompress will decompress the bitstream.
// If the buffer is over-read an error is returned.
func (s *Scratch) decompress() error ***REMOVED***
	br := &s.bits
	br.init(s.br.unread())

	var s1, s2 decoder
	// Initialize and decode first state and symbol.
	s1.init(br, s.decTable, s.actualTableLog)
	s2.init(br, s.decTable, s.actualTableLog)

	// Use temp table to avoid bound checks/append penalty.
	var tmp = s.ct.tableSymbol[:256]
	var off uint8

	// Main part
	if !s.zeroBits ***REMOVED***
		for br.off >= 8 ***REMOVED***
			br.fillFast()
			tmp[off+0] = s1.nextFast()
			tmp[off+1] = s2.nextFast()
			br.fillFast()
			tmp[off+2] = s1.nextFast()
			tmp[off+3] = s2.nextFast()
			off += 4
			if off == 0 ***REMOVED***
				s.Out = append(s.Out, tmp...)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for br.off >= 8 ***REMOVED***
			br.fillFast()
			tmp[off+0] = s1.next()
			tmp[off+1] = s2.next()
			br.fillFast()
			tmp[off+2] = s1.next()
			tmp[off+3] = s2.next()
			off += 4
			if off == 0 ***REMOVED***
				s.Out = append(s.Out, tmp...)
				off = 0
				if len(s.Out) >= s.DecompressLimit ***REMOVED***
					return fmt.Errorf("output size (%d) > DecompressLimit (%d)", len(s.Out), s.DecompressLimit)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	s.Out = append(s.Out, tmp[:off]...)

	// Final bits, a bit more expensive check
	for ***REMOVED***
		if s1.finished() ***REMOVED***
			s.Out = append(s.Out, s1.final(), s2.final())
			break
		***REMOVED***
		br.fill()
		s.Out = append(s.Out, s1.next())
		if s2.finished() ***REMOVED***
			s.Out = append(s.Out, s2.final(), s1.final())
			break
		***REMOVED***
		s.Out = append(s.Out, s2.next())
		if len(s.Out) >= s.DecompressLimit ***REMOVED***
			return fmt.Errorf("output size (%d) > DecompressLimit (%d)", len(s.Out), s.DecompressLimit)
		***REMOVED***
	***REMOVED***
	return br.close()
***REMOVED***

// decoder keeps track of the current state and updates it from the bitstream.
type decoder struct ***REMOVED***
	state uint16
	br    *bitReader
	dt    []decSymbol
***REMOVED***

// init will initialize the decoder and read the first state from the stream.
func (d *decoder) init(in *bitReader, dt []decSymbol, tableLog uint8) ***REMOVED***
	d.dt = dt
	d.br = in
	d.state = uint16(in.getBits(tableLog))
***REMOVED***

// next returns the next symbol and sets the next state.
// At least tablelog bits must be available in the bit reader.
func (d *decoder) next() uint8 ***REMOVED***
	n := &d.dt[d.state]
	lowBits := d.br.getBits(n.nbBits)
	d.state = n.newState + lowBits
	return n.symbol
***REMOVED***

// finished returns true if all bits have been read from the bitstream
// and the next state would require reading bits from the input.
func (d *decoder) finished() bool ***REMOVED***
	return d.br.finished() && d.dt[d.state].nbBits > 0
***REMOVED***

// final returns the current state symbol without decoding the next.
func (d *decoder) final() uint8 ***REMOVED***
	return d.dt[d.state].symbol
***REMOVED***

// nextFast returns the next symbol and sets the next state.
// This can only be used if no symbols are 0 bits.
// At least tablelog bits must be available in the bit reader.
func (d *decoder) nextFast() uint8 ***REMOVED***
	n := d.dt[d.state]
	lowBits := d.br.getBitsFast(n.nbBits)
	d.state = n.newState + lowBits
	return n.symbol
***REMOVED***
