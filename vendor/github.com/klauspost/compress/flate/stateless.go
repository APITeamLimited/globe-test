package flate

import (
	"io"
	"math"
	"sync"
)

const (
	maxStatelessBlock = math.MaxInt16
	// dictionary will be taken from maxStatelessBlock, so limit it.
	maxStatelessDict = 8 << 10

	slTableBits  = 13
	slTableSize  = 1 << slTableBits
	slTableShift = 32 - slTableBits
)

type statelessWriter struct ***REMOVED***
	dst    io.Writer
	closed bool
***REMOVED***

func (s *statelessWriter) Close() error ***REMOVED***
	if s.closed ***REMOVED***
		return nil
	***REMOVED***
	s.closed = true
	// Emit EOF block
	return StatelessDeflate(s.dst, nil, true, nil)
***REMOVED***

func (s *statelessWriter) Write(p []byte) (n int, err error) ***REMOVED***
	err = StatelessDeflate(s.dst, p, false, nil)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return len(p), nil
***REMOVED***

func (s *statelessWriter) Reset(w io.Writer) ***REMOVED***
	s.dst = w
	s.closed = false
***REMOVED***

// NewStatelessWriter will do compression but without maintaining any state
// between Write calls.
// There will be no memory kept between Write calls,
// but compression and speed will be suboptimal.
// Because of this, the size of actual Write calls will affect output size.
func NewStatelessWriter(dst io.Writer) io.WriteCloser ***REMOVED***
	return &statelessWriter***REMOVED***dst: dst***REMOVED***
***REMOVED***

// bitWriterPool contains bit writers that can be reused.
var bitWriterPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return newHuffmanBitWriter(nil)
	***REMOVED***,
***REMOVED***

// StatelessDeflate allows to compress directly to a Writer without retaining state.
// When returning everything will be flushed.
// Up to 8KB of an optional dictionary can be given which is presumed to presumed to precede the block.
// Longer dictionaries will be truncated and will still produce valid output.
// Sending nil dictionary is perfectly fine.
func StatelessDeflate(out io.Writer, in []byte, eof bool, dict []byte) error ***REMOVED***
	var dst tokens
	bw := bitWriterPool.Get().(*huffmanBitWriter)
	bw.reset(out)
	defer func() ***REMOVED***
		// don't keep a reference to our output
		bw.reset(nil)
		bitWriterPool.Put(bw)
	***REMOVED***()
	if eof && len(in) == 0 ***REMOVED***
		// Just write an EOF block.
		// Could be faster...
		bw.writeStoredHeader(0, true)
		bw.flush()
		return bw.err
	***REMOVED***

	// Truncate dict
	if len(dict) > maxStatelessDict ***REMOVED***
		dict = dict[len(dict)-maxStatelessDict:]
	***REMOVED***

	for len(in) > 0 ***REMOVED***
		todo := in
		if len(todo) > maxStatelessBlock-len(dict) ***REMOVED***
			todo = todo[:maxStatelessBlock-len(dict)]
		***REMOVED***
		in = in[len(todo):]
		uncompressed := todo
		if len(dict) > 0 ***REMOVED***
			// combine dict and source
			bufLen := len(todo) + len(dict)
			combined := make([]byte, bufLen)
			copy(combined, dict)
			copy(combined[len(dict):], todo)
			todo = combined
		***REMOVED***
		// Compress
		statelessEnc(&dst, todo, int16(len(dict)))
		isEof := eof && len(in) == 0

		if dst.n == 0 ***REMOVED***
			bw.writeStoredHeader(len(uncompressed), isEof)
			if bw.err != nil ***REMOVED***
				return bw.err
			***REMOVED***
			bw.writeBytes(uncompressed)
		***REMOVED*** else if int(dst.n) > len(uncompressed)-len(uncompressed)>>4 ***REMOVED***
			// If we removed less than 1/16th, huffman compress the block.
			bw.writeBlockHuff(isEof, uncompressed, len(in) == 0)
		***REMOVED*** else ***REMOVED***
			bw.writeBlockDynamic(&dst, isEof, uncompressed, len(in) == 0)
		***REMOVED***
		if len(in) > 0 ***REMOVED***
			// Retain a dict if we have more
			dict = todo[len(todo)-maxStatelessDict:]
			dst.Reset()
		***REMOVED***
		if bw.err != nil ***REMOVED***
			return bw.err
		***REMOVED***
	***REMOVED***
	if !eof ***REMOVED***
		// Align, only a stored block can do that.
		bw.writeStoredHeader(0, false)
	***REMOVED***
	bw.flush()
	return bw.err
***REMOVED***

func hashSL(u uint32) uint32 ***REMOVED***
	return (u * 0x1e35a7bd) >> slTableShift
***REMOVED***

func load3216(b []byte, i int16) uint32 ***REMOVED***
	// Help the compiler eliminate bounds checks on the read so it can be done in a single read.
	b = b[i:]
	b = b[:4]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
***REMOVED***

func load6416(b []byte, i int16) uint64 ***REMOVED***
	// Help the compiler eliminate bounds checks on the read so it can be done in a single read.
	b = b[i:]
	b = b[:8]
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
***REMOVED***

func statelessEnc(dst *tokens, src []byte, startAt int16) ***REMOVED***
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	type tableEntry struct ***REMOVED***
		offset int16
	***REMOVED***

	var table [slTableSize]tableEntry

	// This check isn't in the Snappy implementation, but there, the caller
	// instead of the callee handles this case.
	if len(src)-int(startAt) < minNonLiteralBlockSize ***REMOVED***
		// We do not fill the token table.
		// This will be picked up by caller.
		dst.n = 0
		return
	***REMOVED***
	// Index until startAt
	if startAt > 0 ***REMOVED***
		cv := load3232(src, 0)
		for i := int16(0); i < startAt; i++ ***REMOVED***
			table[hashSL(cv)] = tableEntry***REMOVED***offset: i***REMOVED***
			cv = (cv >> 8) | (uint32(src[i+4]) << 24)
		***REMOVED***
	***REMOVED***

	s := startAt + 1
	nextEmit := startAt
	// sLimit is when to stop looking for offset/length copies. The inputMargin
	// lets us use a fast path for emitLiteral in the main loop, while we are
	// looking for copies.
	sLimit := int16(len(src) - inputMargin)

	// nextEmit is where in src the next emitLiteral should start from.
	cv := load3216(src, s)

	for ***REMOVED***
		const skipLog = 5
		const doEvery = 2

		nextS := s
		var candidate tableEntry
		for ***REMOVED***
			nextHash := hashSL(cv)
			candidate = table[nextHash]
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit || nextS <= 0 ***REMOVED***
				goto emitRemainder
			***REMOVED***

			now := load6416(src, nextS)
			table[nextHash] = tableEntry***REMOVED***offset: s***REMOVED***
			nextHash = hashSL(uint32(now))

			if cv == load3216(src, candidate.offset) ***REMOVED***
				table[nextHash] = tableEntry***REMOVED***offset: nextS***REMOVED***
				break
			***REMOVED***

			// Do one right away...
			cv = uint32(now)
			s = nextS
			nextS++
			candidate = table[nextHash]
			now >>= 8
			table[nextHash] = tableEntry***REMOVED***offset: s***REMOVED***

			if cv == load3216(src, candidate.offset) ***REMOVED***
				table[nextHash] = tableEntry***REMOVED***offset: nextS***REMOVED***
				break
			***REMOVED***
			cv = uint32(now)
			s = nextS
		***REMOVED***

		// A 4-byte match has been found. We'll later see if more than 4 bytes
		// match. But, prior to the match, src[nextEmit:s] are unmatched. Emit
		// them as literal bytes.
		for ***REMOVED***
			// Invariant: we have a 4-byte match at s, and no need to emit any
			// literal bytes prior to s.

			// Extend the 4-byte match as long as possible.
			t := candidate.offset
			l := int16(matchLen(src[s+4:], src[t+4:]) + 4)

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

			// Save the match found
			dst.AddMatchLong(int32(l), uint32(s-t-baseMatchOffset))
			s += l
			nextEmit = s
			if nextS >= s ***REMOVED***
				s = nextS + 1
			***REMOVED***
			if s >= sLimit ***REMOVED***
				goto emitRemainder
			***REMOVED***

			// We could immediately start working at s now, but to improve
			// compression we first update the hash table at s-2 and at s. If
			// another emitCopy is not our next move, also calculate nextHash
			// at s+1. At least on GOARCH=amd64, these three hash calculations
			// are faster as one load64 call (with some shifts) instead of
			// three load32 calls.
			x := load6416(src, s-2)
			o := s - 2
			prevHash := hashSL(uint32(x))
			table[prevHash] = tableEntry***REMOVED***offset: o***REMOVED***
			x >>= 16
			currHash := hashSL(uint32(x))
			candidate = table[currHash]
			table[currHash] = tableEntry***REMOVED***offset: o + 2***REMOVED***

			if uint32(x) != load3216(src, candidate.offset) ***REMOVED***
				cv = uint32(x >> 8)
				s++
				break
			***REMOVED***
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
