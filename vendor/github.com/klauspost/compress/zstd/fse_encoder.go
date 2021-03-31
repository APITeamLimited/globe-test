// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"fmt"
	"math"
)

const (
	// For encoding we only support up to
	maxEncTableLog    = 8
	maxEncTablesize   = 1 << maxTableLog
	maxEncTableMask   = (1 << maxTableLog) - 1
	minEncTablelog    = 5
	maxEncSymbolValue = maxMatchLengthSymbol
)

// Scratch provides temporary storage for compression and decompression.
type fseEncoder struct ***REMOVED***
	symbolLen      uint16 // Length of active part of the symbol table.
	actualTableLog uint8  // Selected tablelog.
	ct             cTable // Compression tables.
	maxCount       int    // count of the most probable symbol
	zeroBits       bool   // no bits has prob > 50%.
	clearCount     bool   // clear count
	useRLE         bool   // This encoder is for RLE
	preDefined     bool   // This encoder is predefined.
	reUsed         bool   // Set to know when the encoder has been reused.
	rleVal         uint8  // RLE Symbol
	maxBits        uint8  // Maximum output bits after transform.

	// TODO: Technically zstd should be fine with 64 bytes.
	count [256]uint32
	norm  [256]int16
***REMOVED***

// cTable contains tables used for compression.
type cTable struct ***REMOVED***
	tableSymbol []byte
	stateTable  []uint16
	symbolTT    []symbolTransform
***REMOVED***

// symbolTransform contains the state transform for a symbol.
type symbolTransform struct ***REMOVED***
	deltaNbBits    uint32
	deltaFindState int16
	outBits        uint8
***REMOVED***

// String prints values as a human readable string.
func (s symbolTransform) String() string ***REMOVED***
	return fmt.Sprintf("***REMOVED***deltabits: %08x, findstate:%d outbits:%d***REMOVED***", s.deltaNbBits, s.deltaFindState, s.outBits)
***REMOVED***

// Histogram allows to populate the histogram and skip that step in the compression,
// It otherwise allows to inspect the histogram when compression is done.
// To indicate that you have populated the histogram call HistogramFinished
// with the value of the highest populated symbol, as well as the number of entries
// in the most populated entry. These are accepted at face value.
// The returned slice will always be length 256.
func (s *fseEncoder) Histogram() []uint32 ***REMOVED***
	return s.count[:]
***REMOVED***

// HistogramFinished can be called to indicate that the histogram has been populated.
// maxSymbol is the index of the highest set symbol of the next data segment.
// maxCount is the number of entries in the most populated entry.
// These are accepted at face value.
func (s *fseEncoder) HistogramFinished(maxSymbol uint8, maxCount int) ***REMOVED***
	s.maxCount = maxCount
	s.symbolLen = uint16(maxSymbol) + 1
	s.clearCount = maxCount != 0
***REMOVED***

// prepare will prepare and allocate scratch tables used for both compression and decompression.
func (s *fseEncoder) prepare() (*fseEncoder, error) ***REMOVED***
	if s == nil ***REMOVED***
		s = &fseEncoder***REMOVED******REMOVED***
	***REMOVED***
	s.useRLE = false
	if s.clearCount && s.maxCount == 0 ***REMOVED***
		for i := range s.count ***REMOVED***
			s.count[i] = 0
		***REMOVED***
		s.clearCount = false
	***REMOVED***
	return s, nil
***REMOVED***

// allocCtable will allocate tables needed for compression.
// If existing tables a re big enough, they are simply re-used.
func (s *fseEncoder) allocCtable() ***REMOVED***
	tableSize := 1 << s.actualTableLog
	// get tableSymbol that is big enough.
	if cap(s.ct.tableSymbol) < tableSize ***REMOVED***
		s.ct.tableSymbol = make([]byte, tableSize)
	***REMOVED***
	s.ct.tableSymbol = s.ct.tableSymbol[:tableSize]

	ctSize := tableSize
	if cap(s.ct.stateTable) < ctSize ***REMOVED***
		s.ct.stateTable = make([]uint16, ctSize)
	***REMOVED***
	s.ct.stateTable = s.ct.stateTable[:ctSize]

	if cap(s.ct.symbolTT) < 256 ***REMOVED***
		s.ct.symbolTT = make([]symbolTransform, 256)
	***REMOVED***
	s.ct.symbolTT = s.ct.symbolTT[:256]
***REMOVED***

// buildCTable will populate the compression table so it is ready to be used.
func (s *fseEncoder) buildCTable() error ***REMOVED***
	tableSize := uint32(1 << s.actualTableLog)
	highThreshold := tableSize - 1
	var cumul [256]int16

	s.allocCtable()
	tableSymbol := s.ct.tableSymbol[:tableSize]
	// symbol start positions
	***REMOVED***
		cumul[0] = 0
		for ui, v := range s.norm[:s.symbolLen-1] ***REMOVED***
			u := byte(ui) // one less than reference
			if v == -1 ***REMOVED***
				// Low proba symbol
				cumul[u+1] = cumul[u] + 1
				tableSymbol[highThreshold] = u
				highThreshold--
			***REMOVED*** else ***REMOVED***
				cumul[u+1] = cumul[u] + v
			***REMOVED***
		***REMOVED***
		// Encode last symbol separately to avoid overflowing u
		u := int(s.symbolLen - 1)
		v := s.norm[s.symbolLen-1]
		if v == -1 ***REMOVED***
			// Low proba symbol
			cumul[u+1] = cumul[u] + 1
			tableSymbol[highThreshold] = byte(u)
			highThreshold--
		***REMOVED*** else ***REMOVED***
			cumul[u+1] = cumul[u] + v
		***REMOVED***
		if uint32(cumul[s.symbolLen]) != tableSize ***REMOVED***
			return fmt.Errorf("internal error: expected cumul[s.symbolLen] (%d) == tableSize (%d)", cumul[s.symbolLen], tableSize)
		***REMOVED***
		cumul[s.symbolLen] = int16(tableSize) + 1
	***REMOVED***
	// Spread symbols
	s.zeroBits = false
	***REMOVED***
		step := tableStep(tableSize)
		tableMask := tableSize - 1
		var position uint32
		// if any symbol > largeLimit, we may have 0 bits output.
		largeLimit := int16(1 << (s.actualTableLog - 1))
		for ui, v := range s.norm[:s.symbolLen] ***REMOVED***
			symbol := byte(ui)
			if v > largeLimit ***REMOVED***
				s.zeroBits = true
			***REMOVED***
			for nbOccurrences := int16(0); nbOccurrences < v; nbOccurrences++ ***REMOVED***
				tableSymbol[position] = symbol
				position = (position + step) & tableMask
				for position > highThreshold ***REMOVED***
					position = (position + step) & tableMask
				***REMOVED*** /* Low proba area */
			***REMOVED***
		***REMOVED***

		// Check if we have gone through all positions
		if position != 0 ***REMOVED***
			return errors.New("position!=0")
		***REMOVED***
	***REMOVED***

	// Build table
	table := s.ct.stateTable
	***REMOVED***
		tsi := int(tableSize)
		for u, v := range tableSymbol ***REMOVED***
			// TableU16 : sorted by symbol order; gives next state value
			table[cumul[v]] = uint16(tsi + u)
			cumul[v]++
		***REMOVED***
	***REMOVED***

	// Build Symbol Transformation Table
	***REMOVED***
		total := int16(0)
		symbolTT := s.ct.symbolTT[:s.symbolLen]
		tableLog := s.actualTableLog
		tl := (uint32(tableLog) << 16) - (1 << tableLog)
		for i, v := range s.norm[:s.symbolLen] ***REMOVED***
			switch v ***REMOVED***
			case 0:
			case -1, 1:
				symbolTT[i].deltaNbBits = tl
				symbolTT[i].deltaFindState = total - 1
				total++
			default:
				maxBitsOut := uint32(tableLog) - highBit(uint32(v-1))
				minStatePlus := uint32(v) << maxBitsOut
				symbolTT[i].deltaNbBits = (maxBitsOut << 16) - minStatePlus
				symbolTT[i].deltaFindState = total - v
				total += v
			***REMOVED***
		***REMOVED***
		if total != int16(tableSize) ***REMOVED***
			return fmt.Errorf("total mismatch %d (got) != %d (want)", total, tableSize)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var rtbTable = [...]uint32***REMOVED***0, 473195, 504333, 520860, 550000, 700000, 750000, 830000***REMOVED***

func (s *fseEncoder) setRLE(val byte) ***REMOVED***
	s.allocCtable()
	s.actualTableLog = 0
	s.ct.stateTable = s.ct.stateTable[:1]
	s.ct.symbolTT[val] = symbolTransform***REMOVED***
		deltaFindState: 0,
		deltaNbBits:    0,
	***REMOVED***
	if debug ***REMOVED***
		println("setRLE: val", val, "symbolTT", s.ct.symbolTT[val])
	***REMOVED***
	s.rleVal = val
	s.useRLE = true
***REMOVED***

// setBits will set output bits for the transform.
// if nil is provided, the number of bits is equal to the index.
func (s *fseEncoder) setBits(transform []byte) ***REMOVED***
	if s.reUsed || s.preDefined ***REMOVED***
		return
	***REMOVED***
	if s.useRLE ***REMOVED***
		if transform == nil ***REMOVED***
			s.ct.symbolTT[s.rleVal].outBits = s.rleVal
			s.maxBits = s.rleVal
			return
		***REMOVED***
		s.maxBits = transform[s.rleVal]
		s.ct.symbolTT[s.rleVal].outBits = s.maxBits
		return
	***REMOVED***
	if transform == nil ***REMOVED***
		for i := range s.ct.symbolTT[:s.symbolLen] ***REMOVED***
			s.ct.symbolTT[i].outBits = uint8(i)
		***REMOVED***
		s.maxBits = uint8(s.symbolLen - 1)
		return
	***REMOVED***
	s.maxBits = 0
	for i, v := range transform[:s.symbolLen] ***REMOVED***
		s.ct.symbolTT[i].outBits = v
		if v > s.maxBits ***REMOVED***
			// We could assume bits always going up, but we play safe.
			s.maxBits = v
		***REMOVED***
	***REMOVED***
***REMOVED***

// normalizeCount will normalize the count of the symbols so
// the total is equal to the table size.
// If successful, compression tables will also be made ready.
func (s *fseEncoder) normalizeCount(length int) error ***REMOVED***
	if s.reUsed ***REMOVED***
		return nil
	***REMOVED***
	s.optimalTableLog(length)
	var (
		tableLog          = s.actualTableLog
		scale             = 62 - uint64(tableLog)
		step              = (1 << 62) / uint64(length)
		vStep             = uint64(1) << (scale - 20)
		stillToDistribute = int16(1 << tableLog)
		largest           int
		largestP          int16
		lowThreshold      = (uint32)(length >> tableLog)
	)
	if s.maxCount == length ***REMOVED***
		s.useRLE = true
		return nil
	***REMOVED***
	s.useRLE = false
	for i, cnt := range s.count[:s.symbolLen] ***REMOVED***
		// already handled
		// if (count[s] == s.length) return 0;   /* rle special case */

		if cnt == 0 ***REMOVED***
			s.norm[i] = 0
			continue
		***REMOVED***
		if cnt <= lowThreshold ***REMOVED***
			s.norm[i] = -1
			stillToDistribute--
		***REMOVED*** else ***REMOVED***
			proba := (int16)((uint64(cnt) * step) >> scale)
			if proba < 8 ***REMOVED***
				restToBeat := vStep * uint64(rtbTable[proba])
				v := uint64(cnt)*step - (uint64(proba) << scale)
				if v > restToBeat ***REMOVED***
					proba++
				***REMOVED***
			***REMOVED***
			if proba > largestP ***REMOVED***
				largestP = proba
				largest = i
			***REMOVED***
			s.norm[i] = proba
			stillToDistribute -= proba
		***REMOVED***
	***REMOVED***

	if -stillToDistribute >= (s.norm[largest] >> 1) ***REMOVED***
		// corner case, need another normalization method
		err := s.normalizeCount2(length)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if debugAsserts ***REMOVED***
			err = s.validateNorm()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return s.buildCTable()
	***REMOVED***
	s.norm[largest] += stillToDistribute
	if debugAsserts ***REMOVED***
		err := s.validateNorm()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return s.buildCTable()
***REMOVED***

// Secondary normalization method.
// To be used when primary method fails.
func (s *fseEncoder) normalizeCount2(length int) error ***REMOVED***
	const notYetAssigned = -2
	var (
		distributed  uint32
		total        = uint32(length)
		tableLog     = s.actualTableLog
		lowThreshold = total >> tableLog
		lowOne       = (total * 3) >> (tableLog + 1)
	)
	for i, cnt := range s.count[:s.symbolLen] ***REMOVED***
		if cnt == 0 ***REMOVED***
			s.norm[i] = 0
			continue
		***REMOVED***
		if cnt <= lowThreshold ***REMOVED***
			s.norm[i] = -1
			distributed++
			total -= cnt
			continue
		***REMOVED***
		if cnt <= lowOne ***REMOVED***
			s.norm[i] = 1
			distributed++
			total -= cnt
			continue
		***REMOVED***
		s.norm[i] = notYetAssigned
	***REMOVED***
	toDistribute := (1 << tableLog) - distributed

	if (total / toDistribute) > lowOne ***REMOVED***
		// risk of rounding to zero
		lowOne = (total * 3) / (toDistribute * 2)
		for i, cnt := range s.count[:s.symbolLen] ***REMOVED***
			if (s.norm[i] == notYetAssigned) && (cnt <= lowOne) ***REMOVED***
				s.norm[i] = 1
				distributed++
				total -= cnt
				continue
			***REMOVED***
		***REMOVED***
		toDistribute = (1 << tableLog) - distributed
	***REMOVED***
	if distributed == uint32(s.symbolLen)+1 ***REMOVED***
		// all values are pretty poor;
		//   probably incompressible data (should have already been detected);
		//   find max, then give all remaining points to max
		var maxV int
		var maxC uint32
		for i, cnt := range s.count[:s.symbolLen] ***REMOVED***
			if cnt > maxC ***REMOVED***
				maxV = i
				maxC = cnt
			***REMOVED***
		***REMOVED***
		s.norm[maxV] += int16(toDistribute)
		return nil
	***REMOVED***

	if total == 0 ***REMOVED***
		// all of the symbols were low enough for the lowOne or lowThreshold
		for i := uint32(0); toDistribute > 0; i = (i + 1) % (uint32(s.symbolLen)) ***REMOVED***
			if s.norm[i] > 0 ***REMOVED***
				toDistribute--
				s.norm[i]++
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	var (
		vStepLog = 62 - uint64(tableLog)
		mid      = uint64((1 << (vStepLog - 1)) - 1)
		rStep    = (((1 << vStepLog) * uint64(toDistribute)) + mid) / uint64(total) // scale on remaining
		tmpTotal = mid
	)
	for i, cnt := range s.count[:s.symbolLen] ***REMOVED***
		if s.norm[i] == notYetAssigned ***REMOVED***
			var (
				end    = tmpTotal + uint64(cnt)*rStep
				sStart = uint32(tmpTotal >> vStepLog)
				sEnd   = uint32(end >> vStepLog)
				weight = sEnd - sStart
			)
			if weight < 1 ***REMOVED***
				return errors.New("weight < 1")
			***REMOVED***
			s.norm[i] = int16(weight)
			tmpTotal = end
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// optimalTableLog calculates and sets the optimal tableLog in s.actualTableLog
func (s *fseEncoder) optimalTableLog(length int) ***REMOVED***
	tableLog := uint8(maxEncTableLog)
	minBitsSrc := highBit(uint32(length)) + 1
	minBitsSymbols := highBit(uint32(s.symbolLen-1)) + 2
	minBits := uint8(minBitsSymbols)
	if minBitsSrc < minBitsSymbols ***REMOVED***
		minBits = uint8(minBitsSrc)
	***REMOVED***

	maxBitsSrc := uint8(highBit(uint32(length-1))) - 2
	if maxBitsSrc < tableLog ***REMOVED***
		// Accuracy can be reduced
		tableLog = maxBitsSrc
	***REMOVED***
	if minBits > tableLog ***REMOVED***
		tableLog = minBits
	***REMOVED***
	// Need a minimum to safely represent all symbol values
	if tableLog < minEncTablelog ***REMOVED***
		tableLog = minEncTablelog
	***REMOVED***
	if tableLog > maxEncTableLog ***REMOVED***
		tableLog = maxEncTableLog
	***REMOVED***
	s.actualTableLog = tableLog
***REMOVED***

// validateNorm validates the normalized histogram table.
func (s *fseEncoder) validateNorm() (err error) ***REMOVED***
	var total int
	for _, v := range s.norm[:s.symbolLen] ***REMOVED***
		if v >= 0 ***REMOVED***
			total += int(v)
		***REMOVED*** else ***REMOVED***
			total -= int(v)
		***REMOVED***
	***REMOVED***
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			return
		***REMOVED***
		fmt.Printf("selected TableLog: %d, Symbol length: %d\n", s.actualTableLog, s.symbolLen)
		for i, v := range s.norm[:s.symbolLen] ***REMOVED***
			fmt.Printf("%3d: %5d -> %4d \n", i, s.count[i], v)
		***REMOVED***
	***REMOVED***()
	if total != (1 << s.actualTableLog) ***REMOVED***
		return fmt.Errorf("warning: Total == %d != %d", total, 1<<s.actualTableLog)
	***REMOVED***
	for i, v := range s.count[s.symbolLen:] ***REMOVED***
		if v != 0 ***REMOVED***
			return fmt.Errorf("warning: Found symbol out of range, %d after cut", i)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// writeCount will write the normalized histogram count to header.
// This is read back by readNCount.
func (s *fseEncoder) writeCount(out []byte) ([]byte, error) ***REMOVED***
	if s.useRLE ***REMOVED***
		return append(out, s.rleVal), nil
	***REMOVED***
	if s.preDefined || s.reUsed ***REMOVED***
		// Never write predefined.
		return out, nil
	***REMOVED***

	var (
		tableLog  = s.actualTableLog
		tableSize = 1 << tableLog
		previous0 bool
		charnum   uint16

		// maximum header size plus 2 extra bytes for final output if bitCount == 0.
		maxHeaderSize = ((int(s.symbolLen) * int(tableLog)) >> 3) + 3 + 2

		// Write Table Size
		bitStream = uint32(tableLog - minEncTablelog)
		bitCount  = uint(4)
		remaining = int16(tableSize + 1) /* +1 for extra accuracy */
		threshold = int16(tableSize)
		nbBits    = uint(tableLog + 1)
		outP      = len(out)
	)
	if cap(out) < outP+maxHeaderSize ***REMOVED***
		out = append(out, make([]byte, maxHeaderSize*3)...)
		out = out[:len(out)-maxHeaderSize*3]
	***REMOVED***
	out = out[:outP+maxHeaderSize]

	// stops at 1
	for remaining > 1 ***REMOVED***
		if previous0 ***REMOVED***
			start := charnum
			for s.norm[charnum] == 0 ***REMOVED***
				charnum++
			***REMOVED***
			for charnum >= start+24 ***REMOVED***
				start += 24
				bitStream += uint32(0xFFFF) << bitCount
				out[outP] = byte(bitStream)
				out[outP+1] = byte(bitStream >> 8)
				outP += 2
				bitStream >>= 16
			***REMOVED***
			for charnum >= start+3 ***REMOVED***
				start += 3
				bitStream += 3 << bitCount
				bitCount += 2
			***REMOVED***
			bitStream += uint32(charnum-start) << bitCount
			bitCount += 2
			if bitCount > 16 ***REMOVED***
				out[outP] = byte(bitStream)
				out[outP+1] = byte(bitStream >> 8)
				outP += 2
				bitStream >>= 16
				bitCount -= 16
			***REMOVED***
		***REMOVED***

		count := s.norm[charnum]
		charnum++
		max := (2*threshold - 1) - remaining
		if count < 0 ***REMOVED***
			remaining += count
		***REMOVED*** else ***REMOVED***
			remaining -= count
		***REMOVED***
		count++ // +1 for extra accuracy
		if count >= threshold ***REMOVED***
			count += max // [0..max[ [max..threshold[ (...) [threshold+max 2*threshold[
		***REMOVED***
		bitStream += uint32(count) << bitCount
		bitCount += nbBits
		if count < max ***REMOVED***
			bitCount--
		***REMOVED***

		previous0 = count == 1
		if remaining < 1 ***REMOVED***
			return nil, errors.New("internal error: remaining < 1")
		***REMOVED***
		for remaining < threshold ***REMOVED***
			nbBits--
			threshold >>= 1
		***REMOVED***

		if bitCount > 16 ***REMOVED***
			out[outP] = byte(bitStream)
			out[outP+1] = byte(bitStream >> 8)
			outP += 2
			bitStream >>= 16
			bitCount -= 16
		***REMOVED***
	***REMOVED***

	if outP+2 > len(out) ***REMOVED***
		return nil, fmt.Errorf("internal error: %d > %d, maxheader: %d, sl: %d, tl: %d, normcount: %v", outP+2, len(out), maxHeaderSize, s.symbolLen, int(tableLog), s.norm[:s.symbolLen])
	***REMOVED***
	out[outP] = byte(bitStream)
	out[outP+1] = byte(bitStream >> 8)
	outP += int((bitCount + 7) / 8)

	if charnum > s.symbolLen ***REMOVED***
		return nil, errors.New("internal error: charnum > s.symbolLen")
	***REMOVED***
	return out[:outP], nil
***REMOVED***

// Approximate symbol cost, as fractional value, using fixed-point format (accuracyLog fractional bits)
// note 1 : assume symbolValue is valid (<= maxSymbolValue)
// note 2 : if freq[symbolValue]==0, @return a fake cost of tableLog+1 bits *
func (s *fseEncoder) bitCost(symbolValue uint8, accuracyLog uint32) uint32 ***REMOVED***
	minNbBits := s.ct.symbolTT[symbolValue].deltaNbBits >> 16
	threshold := (minNbBits + 1) << 16
	if debugAsserts ***REMOVED***
		if !(s.actualTableLog < 16) ***REMOVED***
			panic("!s.actualTableLog < 16")
		***REMOVED***
		// ensure enough room for renormalization double shift
		if !(uint8(accuracyLog) < 31-s.actualTableLog) ***REMOVED***
			panic("!uint8(accuracyLog) < 31-s.actualTableLog")
		***REMOVED***
	***REMOVED***
	tableSize := uint32(1) << s.actualTableLog
	deltaFromThreshold := threshold - (s.ct.symbolTT[symbolValue].deltaNbBits + tableSize)
	// linear interpolation (very approximate)
	normalizedDeltaFromThreshold := (deltaFromThreshold << accuracyLog) >> s.actualTableLog
	bitMultiplier := uint32(1) << accuracyLog
	if debugAsserts ***REMOVED***
		if s.ct.symbolTT[symbolValue].deltaNbBits+tableSize > threshold ***REMOVED***
			panic("s.ct.symbolTT[symbolValue].deltaNbBits+tableSize > threshold")
		***REMOVED***
		if normalizedDeltaFromThreshold > bitMultiplier ***REMOVED***
			panic("normalizedDeltaFromThreshold > bitMultiplier")
		***REMOVED***
	***REMOVED***
	return (minNbBits+1)*bitMultiplier - normalizedDeltaFromThreshold
***REMOVED***

// Returns the cost in bits of encoding the distribution in count using ctable.
// Histogram should only be up to the last non-zero symbol.
// Returns an -1 if ctable cannot represent all the symbols in count.
func (s *fseEncoder) approxSize(hist []uint32) uint32 ***REMOVED***
	if int(s.symbolLen) < len(hist) ***REMOVED***
		// More symbols than we have.
		return math.MaxUint32
	***REMOVED***
	if s.useRLE ***REMOVED***
		// We will never reuse RLE encoders.
		return math.MaxUint32
	***REMOVED***
	const kAccuracyLog = 8
	badCost := (uint32(s.actualTableLog) + 1) << kAccuracyLog
	var cost uint32
	for i, v := range hist ***REMOVED***
		if v == 0 ***REMOVED***
			continue
		***REMOVED***
		if s.norm[i] == 0 ***REMOVED***
			return math.MaxUint32
		***REMOVED***
		bitCost := s.bitCost(uint8(i), kAccuracyLog)
		if bitCost > badCost ***REMOVED***
			return math.MaxUint32
		***REMOVED***
		cost += v * bitCost
	***REMOVED***
	return cost >> kAccuracyLog
***REMOVED***

// maxHeaderSize returns the maximum header size in bits.
// This is not exact size, but we want a penalty for new tables anyway.
func (s *fseEncoder) maxHeaderSize() uint32 ***REMOVED***
	if s.preDefined ***REMOVED***
		return 0
	***REMOVED***
	if s.useRLE ***REMOVED***
		return 8
	***REMOVED***
	return (((uint32(s.symbolLen) * uint32(s.actualTableLog)) >> 3) + 3) * 8
***REMOVED***

// cState contains the compression state of a stream.
type cState struct ***REMOVED***
	bw         *bitWriter
	stateTable []uint16
	state      uint16
***REMOVED***

// init will initialize the compression state to the first symbol of the stream.
func (c *cState) init(bw *bitWriter, ct *cTable, first symbolTransform) ***REMOVED***
	c.bw = bw
	c.stateTable = ct.stateTable
	if len(c.stateTable) == 1 ***REMOVED***
		// RLE
		c.stateTable[0] = uint16(0)
		c.state = 0
		return
	***REMOVED***
	nbBitsOut := (first.deltaNbBits + (1 << 15)) >> 16
	im := int32((nbBitsOut << 16) - first.deltaNbBits)
	lu := (im >> nbBitsOut) + int32(first.deltaFindState)
	c.state = c.stateTable[lu]
	return
***REMOVED***

// encode the output symbol provided and write it to the bitstream.
func (c *cState) encode(symbolTT symbolTransform) ***REMOVED***
	nbBitsOut := (uint32(c.state) + symbolTT.deltaNbBits) >> 16
	dstState := int32(c.state>>(nbBitsOut&15)) + int32(symbolTT.deltaFindState)
	c.bw.addBits16NC(c.state, uint8(nbBitsOut))
	c.state = c.stateTable[dstState]
***REMOVED***

// flush will write the tablelog to the output and flush the remaining full bytes.
func (c *cState) flush(tableLog uint8) ***REMOVED***
	c.bw.flush32()
	c.bw.addBits16NC(c.state, tableLog)
***REMOVED***
