package zstd

import (
	"fmt"
	"math/bits"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

const (
	dictShardBits = 6
)

type fastBase struct ***REMOVED***
	// cur is the offset at the start of hist
	cur int32
	// maximum offset. Should be at least 2x block size.
	maxMatchOff int32
	hist        []byte
	crc         *xxhash.Digest
	tmp         [8]byte
	blk         *blockEnc
	lastDictID  uint32
	lowMem      bool
***REMOVED***

// CRC returns the underlying CRC writer.
func (e *fastBase) CRC() *xxhash.Digest ***REMOVED***
	return e.crc
***REMOVED***

// AppendCRC will append the CRC to the destination slice and return it.
func (e *fastBase) AppendCRC(dst []byte) []byte ***REMOVED***
	crc := e.crc.Sum(e.tmp[:0])
	dst = append(dst, crc[7], crc[6], crc[5], crc[4])
	return dst
***REMOVED***

// WindowSize returns the window size of the encoder,
// or a window size small enough to contain the input size, if > 0.
func (e *fastBase) WindowSize(size int64) int32 ***REMOVED***
	if size > 0 && size < int64(e.maxMatchOff) ***REMOVED***
		b := int32(1) << uint(bits.Len(uint(size)))
		// Keep minimum window.
		if b < 1024 ***REMOVED***
			b = 1024
		***REMOVED***
		return b
	***REMOVED***
	return e.maxMatchOff
***REMOVED***

// Block returns the current block.
func (e *fastBase) Block() *blockEnc ***REMOVED***
	return e.blk
***REMOVED***

func (e *fastBase) addBlock(src []byte) int32 ***REMOVED***
	if debugAsserts && e.cur > bufferReset ***REMOVED***
		panic(fmt.Sprintf("ecur (%d) > buffer reset (%d)", e.cur, bufferReset))
	***REMOVED***
	// check if we have space already
	if len(e.hist)+len(src) > cap(e.hist) ***REMOVED***
		if cap(e.hist) == 0 ***REMOVED***
			e.ensureHist(len(src))
		***REMOVED*** else ***REMOVED***
			if cap(e.hist) < int(e.maxMatchOff+maxCompressedBlockSize) ***REMOVED***
				panic(fmt.Errorf("unexpected buffer cap %d, want at least %d with window %d", cap(e.hist), e.maxMatchOff+maxCompressedBlockSize, e.maxMatchOff))
			***REMOVED***
			// Move down
			offset := int32(len(e.hist)) - e.maxMatchOff
			copy(e.hist[0:e.maxMatchOff], e.hist[offset:])
			e.cur += offset
			e.hist = e.hist[:e.maxMatchOff]
		***REMOVED***
	***REMOVED***
	s := int32(len(e.hist))
	e.hist = append(e.hist, src...)
	return s
***REMOVED***

// ensureHist will ensure that history can keep at least this many bytes.
func (e *fastBase) ensureHist(n int) ***REMOVED***
	if cap(e.hist) >= n ***REMOVED***
		return
	***REMOVED***
	l := e.maxMatchOff
	if (e.lowMem && e.maxMatchOff > maxCompressedBlockSize) || e.maxMatchOff <= maxCompressedBlockSize ***REMOVED***
		l += maxCompressedBlockSize
	***REMOVED*** else ***REMOVED***
		l += e.maxMatchOff
	***REMOVED***
	// Make it at least 1MB.
	if l < 1<<20 && !e.lowMem ***REMOVED***
		l = 1 << 20
	***REMOVED***
	// Make it at least the requested size.
	if l < int32(n) ***REMOVED***
		l = int32(n)
	***REMOVED***
	e.hist = make([]byte, 0, l)
***REMOVED***

// useBlock will replace the block with the provided one,
// but transfer recent offsets from the previous.
func (e *fastBase) UseBlock(enc *blockEnc) ***REMOVED***
	enc.reset(e.blk)
	e.blk = enc
***REMOVED***

func (e *fastBase) matchlenNoHist(s, t int32, src []byte) int32 ***REMOVED***
	// Extend the match to be as long as possible.
	return int32(matchLen(src[s:], src[t:]))
***REMOVED***

func (e *fastBase) matchlen(s, t int32, src []byte) int32 ***REMOVED***
	if debugAsserts ***REMOVED***
		if s < 0 ***REMOVED***
			err := fmt.Sprintf("s (%d) < 0", s)
			panic(err)
		***REMOVED***
		if t < 0 ***REMOVED***
			err := fmt.Sprintf("s (%d) < 0", s)
			panic(err)
		***REMOVED***
		if s-t > e.maxMatchOff ***REMOVED***
			err := fmt.Sprintf("s (%d) - t (%d) > maxMatchOff (%d)", s, t, e.maxMatchOff)
			panic(err)
		***REMOVED***
		if len(src)-int(s) > maxCompressedBlockSize ***REMOVED***
			panic(fmt.Sprintf("len(src)-s (%d) > maxCompressedBlockSize (%d)", len(src)-int(s), maxCompressedBlockSize))
		***REMOVED***
	***REMOVED***

	// Extend the match to be as long as possible.
	return int32(matchLen(src[s:], src[t:]))
***REMOVED***

// Reset the encoding table.
func (e *fastBase) resetBase(d *dict, singleBlock bool) ***REMOVED***
	if e.blk == nil ***REMOVED***
		e.blk = &blockEnc***REMOVED***lowMem: e.lowMem***REMOVED***
		e.blk.init()
	***REMOVED*** else ***REMOVED***
		e.blk.reset(nil)
	***REMOVED***
	e.blk.initNewEncode()
	if e.crc == nil ***REMOVED***
		e.crc = xxhash.New()
	***REMOVED*** else ***REMOVED***
		e.crc.Reset()
	***REMOVED***
	if d != nil ***REMOVED***
		low := e.lowMem
		if singleBlock ***REMOVED***
			e.lowMem = true
		***REMOVED***
		e.ensureHist(d.DictContentSize() + maxCompressedBlockSize)
		e.lowMem = low
	***REMOVED***

	// We offset current position so everything will be out of reach.
	// If above reset line, history will be purged.
	if e.cur < bufferReset ***REMOVED***
		e.cur += e.maxMatchOff + int32(len(e.hist))
	***REMOVED***
	e.hist = e.hist[:0]
	if d != nil ***REMOVED***
		// Set offsets (currently not used)
		for i, off := range d.offsets ***REMOVED***
			e.blk.recentOffsets[i] = uint32(off)
			e.blk.prevRecentOffsets[i] = e.blk.recentOffsets[i]
		***REMOVED***
		// Transfer litenc.
		e.blk.dictLitEnc = d.litEnc
		e.hist = append(e.hist, d.content...)
	***REMOVED***
***REMOVED***
