package huff0

import (
	"fmt"
	"runtime"
	"sync"
)

// Compress1X will compress the input.
// The output can be decoded using Decompress1X.
// Supply a Scratch object. The scratch object contains state about re-use,
// So when sharing across independent encodes, be sure to set the re-use policy.
func Compress1X(in []byte, s *Scratch) (out []byte, reUsed bool, err error) ***REMOVED***
	s, err = s.prepare(in)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	return compress(in, s, s.compress1X)
***REMOVED***

// Compress4X will compress the input. The input is split into 4 independent blocks
// and compressed similar to Compress1X.
// The output can be decoded using Decompress4X.
// Supply a Scratch object. The scratch object contains state about re-use,
// So when sharing across independent encodes, be sure to set the re-use policy.
func Compress4X(in []byte, s *Scratch) (out []byte, reUsed bool, err error) ***REMOVED***
	s, err = s.prepare(in)
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***
	if false ***REMOVED***
		// TODO: compress4Xp only slightly faster.
		const parallelThreshold = 8 << 10
		if len(in) < parallelThreshold || runtime.GOMAXPROCS(0) == 1 ***REMOVED***
			return compress(in, s, s.compress4X)
		***REMOVED***
		return compress(in, s, s.compress4Xp)
	***REMOVED***
	return compress(in, s, s.compress4X)
***REMOVED***

func compress(in []byte, s *Scratch, compressor func(src []byte) ([]byte, error)) (out []byte, reUsed bool, err error) ***REMOVED***
	// Nuke previous table if we cannot reuse anyway.
	if s.Reuse == ReusePolicyNone ***REMOVED***
		s.prevTable = s.prevTable[:0]
	***REMOVED***

	// Create histogram, if none was provided.
	maxCount := s.maxCount
	var canReuse = false
	if maxCount == 0 ***REMOVED***
		maxCount, canReuse = s.countSimple(in)
	***REMOVED*** else ***REMOVED***
		canReuse = s.canUseTable(s.prevTable)
	***REMOVED***

	// Reset for next run.
	s.clearCount = true
	s.maxCount = 0
	if maxCount >= len(in) ***REMOVED***
		if maxCount > len(in) ***REMOVED***
			return nil, false, fmt.Errorf("maxCount (%d) > length (%d)", maxCount, len(in))
		***REMOVED***
		if len(in) == 1 ***REMOVED***
			return nil, false, ErrIncompressible
		***REMOVED***
		// One symbol, use RLE
		return nil, false, ErrUseRLE
	***REMOVED***
	if maxCount == 1 || maxCount < (len(in)>>7) ***REMOVED***
		// Each symbol present maximum once or too well distributed.
		return nil, false, ErrIncompressible
	***REMOVED***

	if s.Reuse == ReusePolicyPrefer && canReuse ***REMOVED***
		keepTable := s.cTable
		s.cTable = s.prevTable
		s.Out, err = compressor(in)
		s.cTable = keepTable
		if err == nil && len(s.Out) < len(in) ***REMOVED***
			s.OutData = s.Out
			return s.Out, true, nil
		***REMOVED***
		// Do not attempt to re-use later.
		s.prevTable = s.prevTable[:0]
	***REMOVED***

	// Calculate new table.
	s.optimalTableLog()
	err = s.buildCTable()
	if err != nil ***REMOVED***
		return nil, false, err
	***REMOVED***

	if false && !s.canUseTable(s.cTable) ***REMOVED***
		panic("invalid table generated")
	***REMOVED***

	if s.Reuse == ReusePolicyAllow && canReuse ***REMOVED***
		hSize := len(s.Out)
		oldSize := s.prevTable.estimateSize(s.count[:s.symbolLen])
		newSize := s.cTable.estimateSize(s.count[:s.symbolLen])
		if oldSize <= hSize+newSize || hSize+12 >= len(in) ***REMOVED***
			// Retain cTable even if we re-use.
			keepTable := s.cTable
			s.cTable = s.prevTable
			s.Out, err = compressor(in)
			s.cTable = keepTable
			if len(s.Out) >= len(in) ***REMOVED***
				return nil, false, ErrIncompressible
			***REMOVED***
			s.OutData = s.Out
			return s.Out, true, nil
		***REMOVED***
	***REMOVED***

	// Use new table
	err = s.cTable.write(s)
	if err != nil ***REMOVED***
		s.OutTable = nil
		return nil, false, err
	***REMOVED***
	s.OutTable = s.Out

	// Compress using new table
	s.Out, err = compressor(in)
	if err != nil ***REMOVED***
		s.OutTable = nil
		return nil, false, err
	***REMOVED***
	if len(s.Out) >= len(in) ***REMOVED***
		s.OutTable = nil
		return nil, false, ErrIncompressible
	***REMOVED***
	// Move current table into previous.
	s.prevTable, s.cTable = s.cTable, s.prevTable[:0]
	s.OutData = s.Out[len(s.OutTable):]
	return s.Out, false, nil
***REMOVED***

func (s *Scratch) compress1X(src []byte) ([]byte, error) ***REMOVED***
	return s.compress1xDo(s.Out, src)
***REMOVED***

func (s *Scratch) compress1xDo(dst, src []byte) ([]byte, error) ***REMOVED***
	var bw = bitWriter***REMOVED***out: dst***REMOVED***

	// N is length divisible by 4.
	n := len(src)
	n -= n & 3
	cTable := s.cTable[:256]

	// Encode last bytes.
	for i := len(src) & 3; i > 0; i-- ***REMOVED***
		bw.encSymbol(cTable, src[n+i-1])
	***REMOVED***
	if s.actualTableLog <= 8 ***REMOVED***
		n -= 4
		for ; n >= 0; n -= 4 ***REMOVED***
			tmp := src[n : n+4]
			// tmp should be len 4
			bw.flush32()
			bw.encSymbol(cTable, tmp[3])
			bw.encSymbol(cTable, tmp[2])
			bw.encSymbol(cTable, tmp[1])
			bw.encSymbol(cTable, tmp[0])
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		n -= 4
		for ; n >= 0; n -= 4 ***REMOVED***
			tmp := src[n : n+4]
			// tmp should be len 4
			bw.flush32()
			bw.encSymbol(cTable, tmp[3])
			bw.encSymbol(cTable, tmp[2])
			bw.flush32()
			bw.encSymbol(cTable, tmp[1])
			bw.encSymbol(cTable, tmp[0])
		***REMOVED***
	***REMOVED***
	err := bw.close()
	return bw.out, err
***REMOVED***

var sixZeros [6]byte

func (s *Scratch) compress4X(src []byte) ([]byte, error) ***REMOVED***
	if len(src) < 12 ***REMOVED***
		return nil, ErrIncompressible
	***REMOVED***
	segmentSize := (len(src) + 3) / 4

	// Add placeholder for output length
	offsetIdx := len(s.Out)
	s.Out = append(s.Out, sixZeros[:]...)

	for i := 0; i < 4; i++ ***REMOVED***
		toDo := src
		if len(toDo) > segmentSize ***REMOVED***
			toDo = toDo[:segmentSize]
		***REMOVED***
		src = src[len(toDo):]

		var err error
		idx := len(s.Out)
		s.Out, err = s.compress1xDo(s.Out, toDo)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// Write compressed length as little endian before block.
		if i < 3 ***REMOVED***
			// Last length is not written.
			length := len(s.Out) - idx
			s.Out[i*2+offsetIdx] = byte(length)
			s.Out[i*2+offsetIdx+1] = byte(length >> 8)
		***REMOVED***
	***REMOVED***

	return s.Out, nil
***REMOVED***

// compress4Xp will compress 4 streams using separate goroutines.
func (s *Scratch) compress4Xp(src []byte) ([]byte, error) ***REMOVED***
	if len(src) < 12 ***REMOVED***
		return nil, ErrIncompressible
	***REMOVED***
	// Add placeholder for output length
	s.Out = s.Out[:6]

	segmentSize := (len(src) + 3) / 4
	var wg sync.WaitGroup
	var errs [4]error
	wg.Add(4)
	for i := 0; i < 4; i++ ***REMOVED***
		toDo := src
		if len(toDo) > segmentSize ***REMOVED***
			toDo = toDo[:segmentSize]
		***REMOVED***
		src = src[len(toDo):]

		// Separate goroutine for each block.
		go func(i int) ***REMOVED***
			s.tmpOut[i], errs[i] = s.compress1xDo(s.tmpOut[i][:0], toDo)
			wg.Done()
		***REMOVED***(i)
	***REMOVED***
	wg.Wait()
	for i := 0; i < 4; i++ ***REMOVED***
		if errs[i] != nil ***REMOVED***
			return nil, errs[i]
		***REMOVED***
		o := s.tmpOut[i]
		// Write compressed length as little endian before block.
		if i < 3 ***REMOVED***
			// Last length is not written.
			s.Out[i*2] = byte(len(o))
			s.Out[i*2+1] = byte(len(o) >> 8)
		***REMOVED***

		// Write output.
		s.Out = append(s.Out, o...)
	***REMOVED***
	return s.Out, nil
***REMOVED***

// countSimple will create a simple histogram in s.count.
// Returns the biggest count.
// Does not update s.clearCount.
func (s *Scratch) countSimple(in []byte) (max int, reuse bool) ***REMOVED***
	reuse = true
	for _, v := range in ***REMOVED***
		s.count[v]++
	***REMOVED***
	m := uint32(0)
	if len(s.prevTable) > 0 ***REMOVED***
		for i, v := range s.count[:] ***REMOVED***
			if v > m ***REMOVED***
				m = v
			***REMOVED***
			if v > 0 ***REMOVED***
				s.symbolLen = uint16(i) + 1
				if i >= len(s.prevTable) ***REMOVED***
					reuse = false
				***REMOVED*** else ***REMOVED***
					if s.prevTable[i].nBits == 0 ***REMOVED***
						reuse = false
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return int(m), reuse
	***REMOVED***
	for i, v := range s.count[:] ***REMOVED***
		if v > m ***REMOVED***
			m = v
		***REMOVED***
		if v > 0 ***REMOVED***
			s.symbolLen = uint16(i) + 1
		***REMOVED***
	***REMOVED***
	return int(m), false
***REMOVED***

func (s *Scratch) canUseTable(c cTable) bool ***REMOVED***
	if len(c) < int(s.symbolLen) ***REMOVED***
		return false
	***REMOVED***
	for i, v := range s.count[:s.symbolLen] ***REMOVED***
		if v != 0 && c[i].nBits == 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// minTableLog provides the minimum logSize to safely represent a distribution.
func (s *Scratch) minTableLog() uint8 ***REMOVED***
	minBitsSrc := highBit32(uint32(s.br.remain()-1)) + 1
	minBitsSymbols := highBit32(uint32(s.symbolLen-1)) + 2
	if minBitsSrc < minBitsSymbols ***REMOVED***
		return uint8(minBitsSrc)
	***REMOVED***
	return uint8(minBitsSymbols)
***REMOVED***

// optimalTableLog calculates and sets the optimal tableLog in s.actualTableLog
func (s *Scratch) optimalTableLog() ***REMOVED***
	tableLog := s.TableLog
	minBits := s.minTableLog()
	maxBitsSrc := uint8(highBit32(uint32(s.br.remain()-1))) - 2
	if maxBitsSrc < tableLog ***REMOVED***
		// Accuracy can be reduced
		tableLog = maxBitsSrc
	***REMOVED***
	if minBits > tableLog ***REMOVED***
		tableLog = minBits
	***REMOVED***
	// Need a minimum to safely represent all symbol values
	if tableLog < minTablelog ***REMOVED***
		tableLog = minTablelog
	***REMOVED***
	if tableLog > tableLogMax ***REMOVED***
		tableLog = tableLogMax
	***REMOVED***
	s.actualTableLog = tableLog
***REMOVED***

type cTableEntry struct ***REMOVED***
	val   uint16
	nBits uint8
	// We have 8 bits extra
***REMOVED***

const huffNodesMask = huffNodesLen - 1

func (s *Scratch) buildCTable() error ***REMOVED***
	s.huffSort()
	if cap(s.cTable) < maxSymbolValue+1 ***REMOVED***
		s.cTable = make([]cTableEntry, s.symbolLen, maxSymbolValue+1)
	***REMOVED*** else ***REMOVED***
		s.cTable = s.cTable[:s.symbolLen]
		for i := range s.cTable ***REMOVED***
			s.cTable[i] = cTableEntry***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	var startNode = int16(s.symbolLen)
	nonNullRank := s.symbolLen - 1

	nodeNb := int16(startNode)
	huffNode := s.nodes[1 : huffNodesLen+1]

	// This overlays the slice above, but allows "-1" index lookups.
	// Different from reference implementation.
	huffNode0 := s.nodes[0 : huffNodesLen+1]

	for huffNode[nonNullRank].count == 0 ***REMOVED***
		nonNullRank--
	***REMOVED***

	lowS := int16(nonNullRank)
	nodeRoot := nodeNb + lowS - 1
	lowN := nodeNb
	huffNode[nodeNb].count = huffNode[lowS].count + huffNode[lowS-1].count
	huffNode[lowS].parent, huffNode[lowS-1].parent = uint16(nodeNb), uint16(nodeNb)
	nodeNb++
	lowS -= 2
	for n := nodeNb; n <= nodeRoot; n++ ***REMOVED***
		huffNode[n].count = 1 << 30
	***REMOVED***
	// fake entry, strong barrier
	huffNode0[0].count = 1 << 31

	// create parents
	for nodeNb <= nodeRoot ***REMOVED***
		var n1, n2 int16
		if huffNode0[lowS+1].count < huffNode0[lowN+1].count ***REMOVED***
			n1 = lowS
			lowS--
		***REMOVED*** else ***REMOVED***
			n1 = lowN
			lowN++
		***REMOVED***
		if huffNode0[lowS+1].count < huffNode0[lowN+1].count ***REMOVED***
			n2 = lowS
			lowS--
		***REMOVED*** else ***REMOVED***
			n2 = lowN
			lowN++
		***REMOVED***

		huffNode[nodeNb].count = huffNode0[n1+1].count + huffNode0[n2+1].count
		huffNode0[n1+1].parent, huffNode0[n2+1].parent = uint16(nodeNb), uint16(nodeNb)
		nodeNb++
	***REMOVED***

	// distribute weights (unlimited tree height)
	huffNode[nodeRoot].nbBits = 0
	for n := nodeRoot - 1; n >= startNode; n-- ***REMOVED***
		huffNode[n].nbBits = huffNode[huffNode[n].parent].nbBits + 1
	***REMOVED***
	for n := uint16(0); n <= nonNullRank; n++ ***REMOVED***
		huffNode[n].nbBits = huffNode[huffNode[n].parent].nbBits + 1
	***REMOVED***
	s.actualTableLog = s.setMaxHeight(int(nonNullRank))
	maxNbBits := s.actualTableLog

	// fill result into tree (val, nbBits)
	if maxNbBits > tableLogMax ***REMOVED***
		return fmt.Errorf("internal error: maxNbBits (%d) > tableLogMax (%d)", maxNbBits, tableLogMax)
	***REMOVED***
	var nbPerRank [tableLogMax + 1]uint16
	var valPerRank [tableLogMax + 1]uint16
	for _, v := range huffNode[:nonNullRank+1] ***REMOVED***
		nbPerRank[v.nbBits]++
	***REMOVED***
	// determine stating value per rank
	***REMOVED***
		min := uint16(0)
		for n := maxNbBits; n > 0; n-- ***REMOVED***
			// get starting value within each rank
			valPerRank[n] = min
			min += nbPerRank[n]
			min >>= 1
		***REMOVED***
	***REMOVED***

	// push nbBits per symbol, symbol order
	// TODO: changed `s.symbolLen` -> `nonNullRank+1` (micro-opt)
	for _, v := range huffNode[:nonNullRank+1] ***REMOVED***
		s.cTable[v.symbol].nBits = v.nbBits
	***REMOVED***

	// assign value within rank, symbol order
	for n, val := range s.cTable[:s.symbolLen] ***REMOVED***
		v := valPerRank[val.nBits]
		s.cTable[n].val = v
		valPerRank[val.nBits] = v + 1
	***REMOVED***

	return nil
***REMOVED***

// huffSort will sort symbols, decreasing order.
func (s *Scratch) huffSort() ***REMOVED***
	type rankPos struct ***REMOVED***
		base    uint32
		current uint32
	***REMOVED***

	// Clear nodes
	nodes := s.nodes[:huffNodesLen+1]
	s.nodes = nodes
	nodes = nodes[1 : huffNodesLen+1]

	// Sort into buckets based on length of symbol count.
	var rank [32]rankPos
	for _, v := range s.count[:s.symbolLen] ***REMOVED***
		r := highBit32(v+1) & 31
		rank[r].base++
	***REMOVED***
	for n := 30; n > 0; n-- ***REMOVED***
		rank[n-1].base += rank[n].base
	***REMOVED***
	for n := range rank[:] ***REMOVED***
		rank[n].current = rank[n].base
	***REMOVED***
	for n, c := range s.count[:s.symbolLen] ***REMOVED***
		r := (highBit32(c+1) + 1) & 31
		pos := rank[r].current
		rank[r].current++
		prev := nodes[(pos-1)&huffNodesMask]
		for pos > rank[r].base && c > prev.count ***REMOVED***
			nodes[pos&huffNodesMask] = prev
			pos--
			prev = nodes[(pos-1)&huffNodesMask]
		***REMOVED***
		nodes[pos&huffNodesMask] = nodeElt***REMOVED***count: c, symbol: byte(n)***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (s *Scratch) setMaxHeight(lastNonNull int) uint8 ***REMOVED***
	maxNbBits := s.TableLog
	huffNode := s.nodes[1 : huffNodesLen+1]
	//huffNode = huffNode[: huffNodesLen]

	largestBits := huffNode[lastNonNull].nbBits

	// early exit : no elt > maxNbBits
	if largestBits <= maxNbBits ***REMOVED***
		return largestBits
	***REMOVED***
	totalCost := int(0)
	baseCost := int(1) << (largestBits - maxNbBits)
	n := uint32(lastNonNull)

	for huffNode[n].nbBits > maxNbBits ***REMOVED***
		totalCost += baseCost - (1 << (largestBits - huffNode[n].nbBits))
		huffNode[n].nbBits = maxNbBits
		n--
	***REMOVED***
	// n stops at huffNode[n].nbBits <= maxNbBits

	for huffNode[n].nbBits == maxNbBits ***REMOVED***
		n--
	***REMOVED***
	// n end at index of smallest symbol using < maxNbBits

	// renorm totalCost
	totalCost >>= largestBits - maxNbBits /* note : totalCost is necessarily a multiple of baseCost */

	// repay normalized cost
	***REMOVED***
		const noSymbol = 0xF0F0F0F0
		var rankLast [tableLogMax + 2]uint32

		for i := range rankLast[:] ***REMOVED***
			rankLast[i] = noSymbol
		***REMOVED***

		// Get pos of last (smallest) symbol per rank
		***REMOVED***
			currentNbBits := uint8(maxNbBits)
			for pos := int(n); pos >= 0; pos-- ***REMOVED***
				if huffNode[pos].nbBits >= currentNbBits ***REMOVED***
					continue
				***REMOVED***
				currentNbBits = huffNode[pos].nbBits // < maxNbBits
				rankLast[maxNbBits-currentNbBits] = uint32(pos)
			***REMOVED***
		***REMOVED***

		for totalCost > 0 ***REMOVED***
			nBitsToDecrease := uint8(highBit32(uint32(totalCost))) + 1

			for ; nBitsToDecrease > 1; nBitsToDecrease-- ***REMOVED***
				highPos := rankLast[nBitsToDecrease]
				lowPos := rankLast[nBitsToDecrease-1]
				if highPos == noSymbol ***REMOVED***
					continue
				***REMOVED***
				if lowPos == noSymbol ***REMOVED***
					break
				***REMOVED***
				highTotal := huffNode[highPos].count
				lowTotal := 2 * huffNode[lowPos].count
				if highTotal <= lowTotal ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			// only triggered when no more rank 1 symbol left => find closest one (note : there is necessarily at least one !)
			// HUF_MAX_TABLELOG test just to please gcc 5+; but it should not be necessary
			// FIXME: try to remove
			for (nBitsToDecrease <= tableLogMax) && (rankLast[nBitsToDecrease] == noSymbol) ***REMOVED***
				nBitsToDecrease++
			***REMOVED***
			totalCost -= 1 << (nBitsToDecrease - 1)
			if rankLast[nBitsToDecrease-1] == noSymbol ***REMOVED***
				// this rank is no longer empty
				rankLast[nBitsToDecrease-1] = rankLast[nBitsToDecrease]
			***REMOVED***
			huffNode[rankLast[nBitsToDecrease]].nbBits++
			if rankLast[nBitsToDecrease] == 0 ***REMOVED***
				/* special case, reached largest symbol */
				rankLast[nBitsToDecrease] = noSymbol
			***REMOVED*** else ***REMOVED***
				rankLast[nBitsToDecrease]--
				if huffNode[rankLast[nBitsToDecrease]].nbBits != maxNbBits-nBitsToDecrease ***REMOVED***
					rankLast[nBitsToDecrease] = noSymbol /* this rank is now empty */
				***REMOVED***
			***REMOVED***
		***REMOVED***

		for totalCost < 0 ***REMOVED*** /* Sometimes, cost correction overshoot */
			if rankLast[1] == noSymbol ***REMOVED*** /* special case : no rank 1 symbol (using maxNbBits-1); let's create one from largest rank 0 (using maxNbBits) */
				for huffNode[n].nbBits == maxNbBits ***REMOVED***
					n--
				***REMOVED***
				huffNode[n+1].nbBits--
				rankLast[1] = n + 1
				totalCost++
				continue
			***REMOVED***
			huffNode[rankLast[1]+1].nbBits--
			rankLast[1]++
			totalCost++
		***REMOVED***
	***REMOVED***
	return maxNbBits
***REMOVED***

type nodeElt struct ***REMOVED***
	count  uint32
	parent uint16
	symbol byte
	nbBits uint8
***REMOVED***
