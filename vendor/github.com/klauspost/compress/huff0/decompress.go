package huff0

import (
	"errors"
	"fmt"
	"io"

	"github.com/klauspost/compress/fse"
)

type dTable struct ***REMOVED***
	single []dEntrySingle
	double []dEntryDouble
***REMOVED***

// single-symbols decoding
type dEntrySingle struct ***REMOVED***
	byte  uint8
	nBits uint8
***REMOVED***

// double-symbols decoding
type dEntryDouble struct ***REMOVED***
	seq   uint16
	nBits uint8
	len   uint8
***REMOVED***

// ReadTable will read a table from the input.
// The size of the input may be larger than the table definition.
// Any content remaining after the table definition will be returned.
// If no Scratch is provided a new one is allocated.
// The returned Scratch can be used for decoding input using this table.
func ReadTable(in []byte, s *Scratch) (s2 *Scratch, remain []byte, err error) ***REMOVED***
	s, err = s.prepare(in)
	if err != nil ***REMOVED***
		return s, nil, err
	***REMOVED***
	if len(in) <= 1 ***REMOVED***
		return s, nil, errors.New("input too small for table")
	***REMOVED***
	iSize := in[0]
	in = in[1:]
	if iSize >= 128 ***REMOVED***
		// Uncompressed
		oSize := iSize - 127
		iSize = (oSize + 1) / 2
		if int(iSize) > len(in) ***REMOVED***
			return s, nil, errors.New("input too small for table")
		***REMOVED***
		for n := uint8(0); n < oSize; n += 2 ***REMOVED***
			v := in[n/2]
			s.huffWeight[n] = v >> 4
			s.huffWeight[n+1] = v & 15
		***REMOVED***
		s.symbolLen = uint16(oSize)
		in = in[iSize:]
	***REMOVED*** else ***REMOVED***
		if len(in) <= int(iSize) ***REMOVED***
			return s, nil, errors.New("input too small for table")
		***REMOVED***
		// FSE compressed weights
		s.fse.DecompressLimit = 255
		hw := s.huffWeight[:]
		s.fse.Out = hw
		b, err := fse.Decompress(in[:iSize], s.fse)
		s.fse.Out = nil
		if err != nil ***REMOVED***
			return s, nil, err
		***REMOVED***
		if len(b) > 255 ***REMOVED***
			return s, nil, errors.New("corrupt input: output table too large")
		***REMOVED***
		s.symbolLen = uint16(len(b))
		in = in[iSize:]
	***REMOVED***

	// collect weight stats
	var rankStats [tableLogMax + 1]uint32
	weightTotal := uint32(0)
	for _, v := range s.huffWeight[:s.symbolLen] ***REMOVED***
		if v > tableLogMax ***REMOVED***
			return s, nil, errors.New("corrupt input: weight too large")
		***REMOVED***
		rankStats[v]++
		weightTotal += (1 << (v & 15)) >> 1
	***REMOVED***
	if weightTotal == 0 ***REMOVED***
		return s, nil, errors.New("corrupt input: weights zero")
	***REMOVED***

	// get last non-null symbol weight (implied, total must be 2^n)
	***REMOVED***
		tableLog := highBit32(weightTotal) + 1
		if tableLog > tableLogMax ***REMOVED***
			return s, nil, errors.New("corrupt input: tableLog too big")
		***REMOVED***
		s.actualTableLog = uint8(tableLog)
		// determine last weight
		***REMOVED***
			total := uint32(1) << tableLog
			rest := total - weightTotal
			verif := uint32(1) << highBit32(rest)
			lastWeight := highBit32(rest) + 1
			if verif != rest ***REMOVED***
				// last value must be a clean power of 2
				return s, nil, errors.New("corrupt input: last value not power of two")
			***REMOVED***
			s.huffWeight[s.symbolLen] = uint8(lastWeight)
			s.symbolLen++
			rankStats[lastWeight]++
		***REMOVED***
	***REMOVED***

	if (rankStats[1] < 2) || (rankStats[1]&1 != 0) ***REMOVED***
		// by construction : at least 2 elts of rank 1, must be even
		return s, nil, errors.New("corrupt input: min elt size, even check failed ")
	***REMOVED***

	// TODO: Choose between single/double symbol decoding

	// Calculate starting value for each rank
	***REMOVED***
		var nextRankStart uint32
		for n := uint8(1); n < s.actualTableLog+1; n++ ***REMOVED***
			current := nextRankStart
			nextRankStart += rankStats[n] << (n - 1)
			rankStats[n] = current
		***REMOVED***
	***REMOVED***

	// fill DTable (always full size)
	tSize := 1 << tableLogMax
	if len(s.dt.single) != tSize ***REMOVED***
		s.dt.single = make([]dEntrySingle, tSize)
	***REMOVED***

	for n, w := range s.huffWeight[:s.symbolLen] ***REMOVED***
		length := (uint32(1) << w) >> 1
		d := dEntrySingle***REMOVED***
			byte:  uint8(n),
			nBits: s.actualTableLog + 1 - w,
		***REMOVED***
		for u := rankStats[w]; u < rankStats[w]+length; u++ ***REMOVED***
			s.dt.single[u] = d
		***REMOVED***
		rankStats[w] += length
	***REMOVED***
	return s, in, nil
***REMOVED***

// Decompress1X will decompress a 1X encoded stream.
// The length of the supplied input must match the end of a block exactly.
// Before this is called, the table must be initialized with ReadTable unless
// the encoder re-used the table.
func (s *Scratch) Decompress1X(in []byte) (out []byte, err error) ***REMOVED***
	if len(s.dt.single) == 0 ***REMOVED***
		return nil, errors.New("no table loaded")
	***REMOVED***
	var br bitReader
	err = br.init(in)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s.Out = s.Out[:0]

	decode := func() byte ***REMOVED***
		val := br.peekBitsFast(s.actualTableLog) /* note : actualTableLog >= 1 */
		v := s.dt.single[val]
		br.bitsRead += v.nBits
		return v.byte
	***REMOVED***
	hasDec := func(v dEntrySingle) byte ***REMOVED***
		br.bitsRead += v.nBits
		return v.byte
	***REMOVED***

	// Avoid bounds check by always having full sized table.
	const tlSize = 1 << tableLogMax
	const tlMask = tlSize - 1
	dt := s.dt.single[:tlSize]

	// Use temp table to avoid bound checks/append penalty.
	var tmp = s.huffWeight[:256]
	var off uint8

	for br.off >= 8 ***REMOVED***
		br.fillFast()
		tmp[off+0] = hasDec(dt[br.peekBitsFast(s.actualTableLog)&tlMask])
		tmp[off+1] = hasDec(dt[br.peekBitsFast(s.actualTableLog)&tlMask])
		br.fillFast()
		tmp[off+2] = hasDec(dt[br.peekBitsFast(s.actualTableLog)&tlMask])
		tmp[off+3] = hasDec(dt[br.peekBitsFast(s.actualTableLog)&tlMask])
		off += 4
		if off == 0 ***REMOVED***
			s.Out = append(s.Out, tmp...)
		***REMOVED***
	***REMOVED***

	s.Out = append(s.Out, tmp[:off]...)

	for !br.finished() ***REMOVED***
		br.fill()
		s.Out = append(s.Out, decode())
	***REMOVED***
	return s.Out, br.close()
***REMOVED***

// Decompress4X will decompress a 4X encoded stream.
// Before this is called, the table must be initialized with ReadTable unless
// the encoder re-used the table.
// The length of the supplied input must match the end of a block exactly.
// The destination size of the uncompressed data must be known and provided.
func (s *Scratch) Decompress4X(in []byte, dstSize int) (out []byte, err error) ***REMOVED***
	if len(s.dt.single) == 0 ***REMOVED***
		return nil, errors.New("no table loaded")
	***REMOVED***
	if len(in) < 6+(4*1) ***REMOVED***
		return nil, errors.New("input too small")
	***REMOVED***
	// TODO: We do not detect when we overrun a buffer, except if the last one does.

	var br [4]bitReader
	start := 6
	for i := 0; i < 3; i++ ***REMOVED***
		length := int(in[i*2]) | (int(in[i*2+1]) << 8)
		if start+length >= len(in) ***REMOVED***
			return nil, errors.New("truncated input (or invalid offset)")
		***REMOVED***
		err = br[i].init(in[start : start+length])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		start += length
	***REMOVED***
	err = br[3].init(in[start:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Prepare output
	if cap(s.Out) < dstSize ***REMOVED***
		s.Out = make([]byte, 0, dstSize)
	***REMOVED***
	s.Out = s.Out[:dstSize]
	// destination, offset to match first output
	dstOut := s.Out
	dstEvery := (dstSize + 3) / 4

	decode := func(br *bitReader) byte ***REMOVED***
		val := br.peekBitsFast(s.actualTableLog) /* note : actualTableLog >= 1 */
		v := s.dt.single[val]
		br.bitsRead += v.nBits
		return v.byte
	***REMOVED***

	// Use temp table to avoid bound checks/append penalty.
	var tmp = s.huffWeight[:256]
	var off uint8

	// Decode 2 values from each decoder/loop.
	const bufoff = 256 / 4
bigloop:
	for ***REMOVED***
		for i := range br ***REMOVED***
			if br[i].off < 4 ***REMOVED***
				break bigloop
			***REMOVED***
			br[i].fillFast()
		***REMOVED***
		tmp[off] = decode(&br[0])
		tmp[off+bufoff] = decode(&br[1])
		tmp[off+bufoff*2] = decode(&br[2])
		tmp[off+bufoff*3] = decode(&br[3])
		tmp[off+1] = decode(&br[0])
		tmp[off+1+bufoff] = decode(&br[1])
		tmp[off+1+bufoff*2] = decode(&br[2])
		tmp[off+1+bufoff*3] = decode(&br[3])
		off += 2
		if off == bufoff ***REMOVED***
			if bufoff > dstEvery ***REMOVED***
				return nil, errors.New("corruption detected: stream overrun")
			***REMOVED***
			copy(dstOut, tmp[:bufoff])
			copy(dstOut[dstEvery:], tmp[bufoff:bufoff*2])
			copy(dstOut[dstEvery*2:], tmp[bufoff*2:bufoff*3])
			copy(dstOut[dstEvery*3:], tmp[bufoff*3:bufoff*4])
			off = 0
			dstOut = dstOut[bufoff:]
			// There must at least be 3 buffers left.
			if len(dstOut) < dstEvery*3+3 ***REMOVED***
				return nil, errors.New("corruption detected: stream overrun")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if off > 0 ***REMOVED***
		ioff := int(off)
		if len(dstOut) < dstEvery*3+ioff ***REMOVED***
			return nil, errors.New("corruption detected: stream overrun")
		***REMOVED***
		copy(dstOut, tmp[:off])
		copy(dstOut[dstEvery:dstEvery+ioff], tmp[bufoff:bufoff*2])
		copy(dstOut[dstEvery*2:dstEvery*2+ioff], tmp[bufoff*2:bufoff*3])
		copy(dstOut[dstEvery*3:dstEvery*3+ioff], tmp[bufoff*3:bufoff*4])
		dstOut = dstOut[off:]
	***REMOVED***

	for i := range br ***REMOVED***
		offset := dstEvery * i
		br := &br[i]
		for !br.finished() ***REMOVED***
			br.fill()
			if offset >= len(dstOut) ***REMOVED***
				return nil, errors.New("corruption detected: stream overrun")
			***REMOVED***
			dstOut[offset] = decode(br)
			offset++
		***REMOVED***
		err = br.close()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return s.Out, nil
***REMOVED***

// matches will compare a decoding table to a coding table.
// Errors are written to the writer.
// Nothing will be written if table is ok.
func (s *Scratch) matches(ct cTable, w io.Writer) ***REMOVED***
	if s == nil || len(s.dt.single) == 0 ***REMOVED***
		return
	***REMOVED***
	dt := s.dt.single[:1<<s.actualTableLog]
	tablelog := s.actualTableLog
	ok := 0
	broken := 0
	for sym, enc := range ct ***REMOVED***
		errs := 0
		broken++
		if enc.nBits == 0 ***REMOVED***
			for _, dec := range dt ***REMOVED***
				if dec.byte == byte(sym) ***REMOVED***
					fmt.Fprintf(w, "symbol %x has decoder, but no encoder\n", sym)
					errs++
					break
				***REMOVED***
			***REMOVED***
			if errs == 0 ***REMOVED***
				broken--
			***REMOVED***
			continue
		***REMOVED***
		// Unused bits in input
		ub := tablelog - enc.nBits
		top := enc.val << ub
		// decoder looks at top bits.
		dec := dt[top]
		if dec.nBits != enc.nBits ***REMOVED***
			fmt.Fprintf(w, "symbol 0x%x bit size mismatch (enc: %d, dec:%d).\n", sym, enc.nBits, dec.nBits)
			errs++
		***REMOVED***
		if dec.byte != uint8(sym) ***REMOVED***
			fmt.Fprintf(w, "symbol 0x%x decoder output mismatch (enc: %d, dec:%d).\n", sym, sym, dec.byte)
			errs++
		***REMOVED***
		if errs > 0 ***REMOVED***
			fmt.Fprintf(w, "%d errros in base, stopping\n", errs)
			continue
		***REMOVED***
		// Ensure that all combinations are covered.
		for i := uint16(0); i < (1 << ub); i++ ***REMOVED***
			vval := top | i
			dec := dt[vval]
			if dec.nBits != enc.nBits ***REMOVED***
				fmt.Fprintf(w, "symbol 0x%x bit size mismatch (enc: %d, dec:%d).\n", vval, enc.nBits, dec.nBits)
				errs++
			***REMOVED***
			if dec.byte != uint8(sym) ***REMOVED***
				fmt.Fprintf(w, "symbol 0x%x decoder output mismatch (enc: %d, dec:%d).\n", vval, sym, dec.byte)
				errs++
			***REMOVED***
			if errs > 20 ***REMOVED***
				fmt.Fprintf(w, "%d errros, stopping\n", errs)
				break
			***REMOVED***
		***REMOVED***
		if errs == 0 ***REMOVED***
			ok++
			broken--
		***REMOVED***
	***REMOVED***
	if broken > 0 ***REMOVED***
		fmt.Fprintf(w, "%d broken, %d ok\n", broken, ok)
	***REMOVED***
***REMOVED***
