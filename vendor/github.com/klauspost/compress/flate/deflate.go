// Copyright 2009 The Go Authors. All rights reserved.
// Copyright (c) 2015 Klaus Post
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	NoCompression      = 0
	BestSpeed          = 1
	BestCompression    = 9
	DefaultCompression = -1

	// HuffmanOnly disables Lempel-Ziv match searching and only performs Huffman
	// entropy encoding. This mode is useful in compressing data that has
	// already been compressed with an LZ style algorithm (e.g. Snappy or LZ4)
	// that lacks an entropy encoder. Compression gains are achieved when
	// certain bytes in the input stream occur more frequently than others.
	//
	// Note that HuffmanOnly produces a compressed output that is
	// RFC 1951 compliant. That is, any valid DEFLATE decompressor will
	// continue to be able to decompress this output.
	HuffmanOnly         = -2
	ConstantCompression = HuffmanOnly // compatibility alias.

	logWindowSize    = 15
	windowSize       = 1 << logWindowSize
	windowMask       = windowSize - 1
	logMaxOffsetSize = 15  // Standard DEFLATE
	minMatchLength   = 4   // The smallest match that the compressor looks for
	maxMatchLength   = 258 // The longest match for the compressor
	minOffsetSize    = 1   // The shortest offset that makes any sense

	// The maximum number of tokens we will encode at the time.
	// Smaller sizes usually creates less optimal blocks.
	// Bigger can make context switching slow.
	// We use this for levels 7-9, so we make it big.
	maxFlateBlockTokens = 1 << 15
	maxStoreBlockSize   = 65535
	hashBits            = 17 // After 17 performance degrades
	hashSize            = 1 << hashBits
	hashMask            = (1 << hashBits) - 1
	hashShift           = (hashBits + minMatchLength - 1) / minMatchLength
	maxHashOffset       = 1 << 28

	skipNever = math.MaxInt32

	debugDeflate = false
)

type compressionLevel struct ***REMOVED***
	good, lazy, nice, chain, fastSkipHashing, level int
***REMOVED***

// Compression levels have been rebalanced from zlib deflate defaults
// to give a bigger spread in speed and compression.
// See https://blog.klauspost.com/rebalancing-deflate-compression-levels/
var levels = []compressionLevel***REMOVED***
	***REMOVED******REMOVED***, // 0
	// Level 1-6 uses specialized algorithm - values not used
	***REMOVED***0, 0, 0, 0, 0, 1***REMOVED***,
	***REMOVED***0, 0, 0, 0, 0, 2***REMOVED***,
	***REMOVED***0, 0, 0, 0, 0, 3***REMOVED***,
	***REMOVED***0, 0, 0, 0, 0, 4***REMOVED***,
	***REMOVED***0, 0, 0, 0, 0, 5***REMOVED***,
	***REMOVED***0, 0, 0, 0, 0, 6***REMOVED***,
	// Levels 7-9 use increasingly more lazy matching
	// and increasingly stringent conditions for "good enough".
	***REMOVED***8, 12, 16, 24, skipNever, 7***REMOVED***,
	***REMOVED***16, 30, 40, 64, skipNever, 8***REMOVED***,
	***REMOVED***32, 258, 258, 1024, skipNever, 9***REMOVED***,
***REMOVED***

// advancedState contains state for the advanced levels, with bigger hash tables, etc.
type advancedState struct ***REMOVED***
	// deflate state
	length         int
	offset         int
	maxInsertIndex int
	chainHead      int
	hashOffset     int

	ii uint16 // position of last match, intended to overflow to reset.

	// input window: unprocessed data is window[index:windowEnd]
	index          int
	estBitsPerByte int
	hashMatch      [maxMatchLength + minMatchLength]uint32

	// Input hash chains
	// hashHead[hashValue] contains the largest inputIndex with the specified hash value
	// If hashHead[hashValue] is within the current window, then
	// hashPrev[hashHead[hashValue] & windowMask] contains the previous index
	// with the same hash value.
	hashHead [hashSize]uint32
	hashPrev [windowSize]uint32
***REMOVED***

type compressor struct ***REMOVED***
	compressionLevel

	h *huffmanEncoder
	w *huffmanBitWriter

	// compression algorithm
	fill func(*compressor, []byte) int // copy data to window
	step func(*compressor)             // process window

	window     []byte
	windowEnd  int
	blockStart int // window index where current tokens start
	err        error

	// queued output tokens
	tokens tokens
	fast   fastEnc
	state  *advancedState

	sync          bool // requesting flush
	byteAvailable bool // if true, still need to process window[index-1].
***REMOVED***

func (d *compressor) fillDeflate(b []byte) int ***REMOVED***
	s := d.state
	if s.index >= 2*windowSize-(minMatchLength+maxMatchLength) ***REMOVED***
		// shift the window by windowSize
		copy(d.window[:], d.window[windowSize:2*windowSize])
		s.index -= windowSize
		d.windowEnd -= windowSize
		if d.blockStart >= windowSize ***REMOVED***
			d.blockStart -= windowSize
		***REMOVED*** else ***REMOVED***
			d.blockStart = math.MaxInt32
		***REMOVED***
		s.hashOffset += windowSize
		if s.hashOffset > maxHashOffset ***REMOVED***
			delta := s.hashOffset - 1
			s.hashOffset -= delta
			s.chainHead -= delta
			// Iterate over slices instead of arrays to avoid copying
			// the entire table onto the stack (Issue #18625).
			for i, v := range s.hashPrev[:] ***REMOVED***
				if int(v) > delta ***REMOVED***
					s.hashPrev[i] = uint32(int(v) - delta)
				***REMOVED*** else ***REMOVED***
					s.hashPrev[i] = 0
				***REMOVED***
			***REMOVED***
			for i, v := range s.hashHead[:] ***REMOVED***
				if int(v) > delta ***REMOVED***
					s.hashHead[i] = uint32(int(v) - delta)
				***REMOVED*** else ***REMOVED***
					s.hashHead[i] = 0
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	n := copy(d.window[d.windowEnd:], b)
	d.windowEnd += n
	return n
***REMOVED***

func (d *compressor) writeBlock(tok *tokens, index int, eof bool) error ***REMOVED***
	if index > 0 || eof ***REMOVED***
		var window []byte
		if d.blockStart <= index ***REMOVED***
			window = d.window[d.blockStart:index]
		***REMOVED***
		d.blockStart = index
		//d.w.writeBlock(tok, eof, window)
		d.w.writeBlockDynamic(tok, eof, window, d.sync)
		return d.w.err
	***REMOVED***
	return nil
***REMOVED***

// writeBlockSkip writes the current block and uses the number of tokens
// to determine if the block should be stored on no matches, or
// only huffman encoded.
func (d *compressor) writeBlockSkip(tok *tokens, index int, eof bool) error ***REMOVED***
	if index > 0 || eof ***REMOVED***
		if d.blockStart <= index ***REMOVED***
			window := d.window[d.blockStart:index]
			// If we removed less than a 64th of all literals
			// we huffman compress the block.
			if int(tok.n) > len(window)-int(tok.n>>6) ***REMOVED***
				d.w.writeBlockHuff(eof, window, d.sync)
			***REMOVED*** else ***REMOVED***
				// Write a dynamic huffman block.
				d.w.writeBlockDynamic(tok, eof, window, d.sync)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			d.w.writeBlock(tok, eof, nil)
		***REMOVED***
		d.blockStart = index
		return d.w.err
	***REMOVED***
	return nil
***REMOVED***

// fillWindow will fill the current window with the supplied
// dictionary and calculate all hashes.
// This is much faster than doing a full encode.
// Should only be used after a start/reset.
func (d *compressor) fillWindow(b []byte) ***REMOVED***
	// Do not fill window if we are in store-only or huffman mode.
	if d.level <= 0 ***REMOVED***
		return
	***REMOVED***
	if d.fast != nil ***REMOVED***
		// encode the last data, but discard the result
		if len(b) > maxMatchOffset ***REMOVED***
			b = b[len(b)-maxMatchOffset:]
		***REMOVED***
		d.fast.Encode(&d.tokens, b)
		d.tokens.Reset()
		return
	***REMOVED***
	s := d.state
	// If we are given too much, cut it.
	if len(b) > windowSize ***REMOVED***
		b = b[len(b)-windowSize:]
	***REMOVED***
	// Add all to window.
	n := copy(d.window[d.windowEnd:], b)

	// Calculate 256 hashes at the time (more L1 cache hits)
	loops := (n + 256 - minMatchLength) / 256
	for j := 0; j < loops; j++ ***REMOVED***
		startindex := j * 256
		end := startindex + 256 + minMatchLength - 1
		if end > n ***REMOVED***
			end = n
		***REMOVED***
		tocheck := d.window[startindex:end]
		dstSize := len(tocheck) - minMatchLength + 1

		if dstSize <= 0 ***REMOVED***
			continue
		***REMOVED***

		dst := s.hashMatch[:dstSize]
		bulkHash4(tocheck, dst)
		var newH uint32
		for i, val := range dst ***REMOVED***
			di := i + startindex
			newH = val & hashMask
			// Get previous value with the same hash.
			// Our chain should point to the previous value.
			s.hashPrev[di&windowMask] = s.hashHead[newH]
			// Set the head of the hash chain to us.
			s.hashHead[newH] = uint32(di + s.hashOffset)
		***REMOVED***
	***REMOVED***
	// Update window information.
	d.windowEnd += n
	s.index = n
***REMOVED***

// Try to find a match starting at index whose length is greater than prevSize.
// We only look at chainCount possibilities before giving up.
// pos = s.index, prevHead = s.chainHead-s.hashOffset, prevLength=minMatchLength-1, lookahead
func (d *compressor) findMatch(pos int, prevHead int, lookahead int) (length, offset int, ok bool) ***REMOVED***
	minMatchLook := maxMatchLength
	if lookahead < minMatchLook ***REMOVED***
		minMatchLook = lookahead
	***REMOVED***

	win := d.window[0 : pos+minMatchLook]

	// We quit when we get a match that's at least nice long
	nice := len(win) - pos
	if d.nice < nice ***REMOVED***
		nice = d.nice
	***REMOVED***

	// If we've got a match that's good enough, only look in 1/4 the chain.
	tries := d.chain
	length = minMatchLength - 1

	wEnd := win[pos+length]
	wPos := win[pos:]
	minIndex := pos - windowSize
	if minIndex < 0 ***REMOVED***
		minIndex = 0
	***REMOVED***
	offset = 0

	cGain := 0
	if d.chain < 100 ***REMOVED***
		for i := prevHead; tries > 0; tries-- ***REMOVED***
			if wEnd == win[i+length] ***REMOVED***
				n := matchLen(win[i:i+minMatchLook], wPos)
				if n > length ***REMOVED***
					length = n
					offset = pos - i
					ok = true
					if n >= nice ***REMOVED***
						// The match is good enough that we don't try to find a better one.
						break
					***REMOVED***
					wEnd = win[pos+n]
				***REMOVED***
			***REMOVED***
			if i <= minIndex ***REMOVED***
				// hashPrev[i & windowMask] has already been overwritten, so stop now.
				break
			***REMOVED***
			i = int(d.state.hashPrev[i&windowMask]) - d.state.hashOffset
			if i < minIndex ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	// Some like it higher (CSV), some like it lower (JSON)
	const baseCost = 6
	// Base is 4 bytes at with an additional cost.
	// Matches must be better than this.
	for i := prevHead; tries > 0; tries-- ***REMOVED***
		if wEnd == win[i+length] ***REMOVED***
			n := matchLen(win[i:i+minMatchLook], wPos)
			if n > length ***REMOVED***
				// Calculate gain. Estimate
				newGain := d.h.bitLengthRaw(wPos[:n]) - int(offsetExtraBits[offsetCode(uint32(pos-i))]) - baseCost - int(lengthExtraBits[lengthCodes[(n-3)&255]])

				//fmt.Println(n, "gain:", newGain, "prev:", cGain, "raw:", d.h.bitLengthRaw(wPos[:n]))
				if newGain > cGain ***REMOVED***
					length = n
					offset = pos - i
					cGain = newGain
					ok = true
					if n >= nice ***REMOVED***
						// The match is good enough that we don't try to find a better one.
						break
					***REMOVED***
					wEnd = win[pos+n]
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if i <= minIndex ***REMOVED***
			// hashPrev[i & windowMask] has already been overwritten, so stop now.
			break
		***REMOVED***
		i = int(d.state.hashPrev[i&windowMask]) - d.state.hashOffset
		if i < minIndex ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (d *compressor) writeStoredBlock(buf []byte) error ***REMOVED***
	if d.w.writeStoredHeader(len(buf), false); d.w.err != nil ***REMOVED***
		return d.w.err
	***REMOVED***
	d.w.writeBytes(buf)
	return d.w.err
***REMOVED***

// hash4 returns a hash representation of the first 4 bytes
// of the supplied slice.
// The caller must ensure that len(b) >= 4.
func hash4(b []byte) uint32 ***REMOVED***
	return hash4u(binary.LittleEndian.Uint32(b), hashBits)
***REMOVED***

// bulkHash4 will compute hashes using the same
// algorithm as hash4
func bulkHash4(b []byte, dst []uint32) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return
	***REMOVED***
	hb := binary.LittleEndian.Uint32(b)

	dst[0] = hash4u(hb, hashBits)
	end := len(b) - 4 + 1
	for i := 1; i < end; i++ ***REMOVED***
		hb = (hb >> 8) | uint32(b[i+3])<<24
		dst[i] = hash4u(hb, hashBits)
	***REMOVED***
***REMOVED***

func (d *compressor) initDeflate() ***REMOVED***
	d.window = make([]byte, 2*windowSize)
	d.byteAvailable = false
	d.err = nil
	if d.state == nil ***REMOVED***
		return
	***REMOVED***
	s := d.state
	s.index = 0
	s.hashOffset = 1
	s.length = minMatchLength - 1
	s.offset = 0
	s.chainHead = -1
***REMOVED***

// deflateLazy is the same as deflate, but with d.fastSkipHashing == skipNever,
// meaning it always has lazy matching on.
func (d *compressor) deflateLazy() ***REMOVED***
	s := d.state
	// Sanity enables additional runtime tests.
	// It's intended to be used during development
	// to supplement the currently ad-hoc unit tests.
	const sanity = debugDeflate

	if d.windowEnd-s.index < minMatchLength+maxMatchLength && !d.sync ***REMOVED***
		return
	***REMOVED***
	if d.windowEnd != s.index && d.chain > 100 ***REMOVED***
		// Get literal huffman coder.
		if d.h == nil ***REMOVED***
			d.h = newHuffmanEncoder(maxFlateBlockTokens)
		***REMOVED***
		var tmp [256]uint16
		for _, v := range d.window[s.index:d.windowEnd] ***REMOVED***
			tmp[v]++
		***REMOVED***
		d.h.generate(tmp[:], 15)
	***REMOVED***

	s.maxInsertIndex = d.windowEnd - (minMatchLength - 1)

	for ***REMOVED***
		if sanity && s.index > d.windowEnd ***REMOVED***
			panic("index > windowEnd")
		***REMOVED***
		lookahead := d.windowEnd - s.index
		if lookahead < minMatchLength+maxMatchLength ***REMOVED***
			if !d.sync ***REMOVED***
				return
			***REMOVED***
			if sanity && s.index > d.windowEnd ***REMOVED***
				panic("index > windowEnd")
			***REMOVED***
			if lookahead == 0 ***REMOVED***
				// Flush current output block if any.
				if d.byteAvailable ***REMOVED***
					// There is still one pending token that needs to be flushed
					d.tokens.AddLiteral(d.window[s.index-1])
					d.byteAvailable = false
				***REMOVED***
				if d.tokens.n > 0 ***REMOVED***
					if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil ***REMOVED***
						return
					***REMOVED***
					d.tokens.Reset()
				***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		if s.index < s.maxInsertIndex ***REMOVED***
			// Update the hash
			hash := hash4(d.window[s.index:])
			ch := s.hashHead[hash]
			s.chainHead = int(ch)
			s.hashPrev[s.index&windowMask] = ch
			s.hashHead[hash] = uint32(s.index + s.hashOffset)
		***REMOVED***
		prevLength := s.length
		prevOffset := s.offset
		s.length = minMatchLength - 1
		s.offset = 0
		minIndex := s.index - windowSize
		if minIndex < 0 ***REMOVED***
			minIndex = 0
		***REMOVED***

		if s.chainHead-s.hashOffset >= minIndex && lookahead > prevLength && prevLength < d.lazy ***REMOVED***
			if newLength, newOffset, ok := d.findMatch(s.index, s.chainHead-s.hashOffset, lookahead); ok ***REMOVED***
				s.length = newLength
				s.offset = newOffset
			***REMOVED***
		***REMOVED***

		if prevLength >= minMatchLength && s.length <= prevLength ***REMOVED***
			// Check for better match at end...
			//
			// checkOff must be >=2 since we otherwise risk checking s.index
			// Offset of 2 seems to yield best results.
			const checkOff = 2
			prevIndex := s.index - 1
			if prevIndex+prevLength+checkOff < s.maxInsertIndex ***REMOVED***
				end := lookahead
				if lookahead > maxMatchLength ***REMOVED***
					end = maxMatchLength
				***REMOVED***
				end += prevIndex
				idx := prevIndex + prevLength - (4 - checkOff)
				h := hash4(d.window[idx:])
				ch2 := int(s.hashHead[h]) - s.hashOffset - prevLength + (4 - checkOff)
				if ch2 > minIndex ***REMOVED***
					length := matchLen(d.window[prevIndex:end], d.window[ch2:])
					// It seems like a pure length metric is best.
					if length > prevLength ***REMOVED***
						prevLength = length
						prevOffset = prevIndex - ch2
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// There was a match at the previous step, and the current match is
			// not better. Output the previous match.
			d.tokens.AddMatch(uint32(prevLength-3), uint32(prevOffset-minOffsetSize))

			// Insert in the hash table all strings up to the end of the match.
			// index and index-1 are already inserted. If there is not enough
			// lookahead, the last two strings are not inserted into the hash
			// table.
			newIndex := s.index + prevLength - 1
			// Calculate missing hashes
			end := newIndex
			if end > s.maxInsertIndex ***REMOVED***
				end = s.maxInsertIndex
			***REMOVED***
			end += minMatchLength - 1
			startindex := s.index + 1
			if startindex > s.maxInsertIndex ***REMOVED***
				startindex = s.maxInsertIndex
			***REMOVED***
			tocheck := d.window[startindex:end]
			dstSize := len(tocheck) - minMatchLength + 1
			if dstSize > 0 ***REMOVED***
				dst := s.hashMatch[:dstSize]
				bulkHash4(tocheck, dst)
				var newH uint32
				for i, val := range dst ***REMOVED***
					di := i + startindex
					newH = val & hashMask
					// Get previous value with the same hash.
					// Our chain should point to the previous value.
					s.hashPrev[di&windowMask] = s.hashHead[newH]
					// Set the head of the hash chain to us.
					s.hashHead[newH] = uint32(di + s.hashOffset)
				***REMOVED***
			***REMOVED***

			s.index = newIndex
			d.byteAvailable = false
			s.length = minMatchLength - 1
			if d.tokens.n == maxFlateBlockTokens ***REMOVED***
				// The block includes the current character
				if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil ***REMOVED***
					return
				***REMOVED***
				d.tokens.Reset()
			***REMOVED***
			s.ii = 0
		***REMOVED*** else ***REMOVED***
			// Reset, if we got a match this run.
			if s.length >= minMatchLength ***REMOVED***
				s.ii = 0
			***REMOVED***
			// We have a byte waiting. Emit it.
			if d.byteAvailable ***REMOVED***
				s.ii++
				d.tokens.AddLiteral(d.window[s.index-1])
				if d.tokens.n == maxFlateBlockTokens ***REMOVED***
					if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil ***REMOVED***
						return
					***REMOVED***
					d.tokens.Reset()
				***REMOVED***
				s.index++

				// If we have a long run of no matches, skip additional bytes
				// Resets when s.ii overflows after 64KB.
				if n := int(s.ii) - d.chain; n > 0 ***REMOVED***
					n = 1 + int(n>>6)
					for j := 0; j < n; j++ ***REMOVED***
						if s.index >= d.windowEnd-1 ***REMOVED***
							break
						***REMOVED***
						d.tokens.AddLiteral(d.window[s.index-1])
						if d.tokens.n == maxFlateBlockTokens ***REMOVED***
							if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil ***REMOVED***
								return
							***REMOVED***
							d.tokens.Reset()
						***REMOVED***
						// Index...
						if s.index < s.maxInsertIndex ***REMOVED***
							h := hash4(d.window[s.index:])
							ch := s.hashHead[h]
							s.chainHead = int(ch)
							s.hashPrev[s.index&windowMask] = ch
							s.hashHead[h] = uint32(s.index + s.hashOffset)
						***REMOVED***
						s.index++
					***REMOVED***
					// Flush last byte
					d.tokens.AddLiteral(d.window[s.index-1])
					d.byteAvailable = false
					// s.length = minMatchLength - 1 // not needed, since s.ii is reset above, so it should never be > minMatchLength
					if d.tokens.n == maxFlateBlockTokens ***REMOVED***
						if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil ***REMOVED***
							return
						***REMOVED***
						d.tokens.Reset()
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				s.index++
				d.byteAvailable = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *compressor) store() ***REMOVED***
	if d.windowEnd > 0 && (d.windowEnd == maxStoreBlockSize || d.sync) ***REMOVED***
		d.err = d.writeStoredBlock(d.window[:d.windowEnd])
		d.windowEnd = 0
	***REMOVED***
***REMOVED***

// fillWindow will fill the buffer with data for huffman-only compression.
// The number of bytes copied is returned.
func (d *compressor) fillBlock(b []byte) int ***REMOVED***
	n := copy(d.window[d.windowEnd:], b)
	d.windowEnd += n
	return n
***REMOVED***

// storeHuff will compress and store the currently added data,
// if enough has been accumulated or we at the end of the stream.
// Any error that occurred will be in d.err
func (d *compressor) storeHuff() ***REMOVED***
	if d.windowEnd < len(d.window) && !d.sync || d.windowEnd == 0 ***REMOVED***
		return
	***REMOVED***
	d.w.writeBlockHuff(false, d.window[:d.windowEnd], d.sync)
	d.err = d.w.err
	d.windowEnd = 0
***REMOVED***

// storeFast will compress and store the currently added data,
// if enough has been accumulated or we at the end of the stream.
// Any error that occurred will be in d.err
func (d *compressor) storeFast() ***REMOVED***
	// We only compress if we have maxStoreBlockSize.
	if d.windowEnd < len(d.window) ***REMOVED***
		if !d.sync ***REMOVED***
			return
		***REMOVED***
		// Handle extremely small sizes.
		if d.windowEnd < 128 ***REMOVED***
			if d.windowEnd == 0 ***REMOVED***
				return
			***REMOVED***
			if d.windowEnd <= 32 ***REMOVED***
				d.err = d.writeStoredBlock(d.window[:d.windowEnd])
			***REMOVED*** else ***REMOVED***
				d.w.writeBlockHuff(false, d.window[:d.windowEnd], true)
				d.err = d.w.err
			***REMOVED***
			d.tokens.Reset()
			d.windowEnd = 0
			d.fast.Reset()
			return
		***REMOVED***
	***REMOVED***

	d.fast.Encode(&d.tokens, d.window[:d.windowEnd])
	// If we made zero matches, store the block as is.
	if d.tokens.n == 0 ***REMOVED***
		d.err = d.writeStoredBlock(d.window[:d.windowEnd])
		// If we removed less than 1/16th, huffman compress the block.
	***REMOVED*** else if int(d.tokens.n) > d.windowEnd-(d.windowEnd>>4) ***REMOVED***
		d.w.writeBlockHuff(false, d.window[:d.windowEnd], d.sync)
		d.err = d.w.err
	***REMOVED*** else ***REMOVED***
		d.w.writeBlockDynamic(&d.tokens, false, d.window[:d.windowEnd], d.sync)
		d.err = d.w.err
	***REMOVED***
	d.tokens.Reset()
	d.windowEnd = 0
***REMOVED***

// write will add input byte to the stream.
// Unless an error occurs all bytes will be consumed.
func (d *compressor) write(b []byte) (n int, err error) ***REMOVED***
	if d.err != nil ***REMOVED***
		return 0, d.err
	***REMOVED***
	n = len(b)
	for len(b) > 0 ***REMOVED***
		if d.windowEnd == len(d.window) || d.sync ***REMOVED***
			d.step(d)
		***REMOVED***
		b = b[d.fill(d, b):]
		if d.err != nil ***REMOVED***
			return 0, d.err
		***REMOVED***
	***REMOVED***
	return n, d.err
***REMOVED***

func (d *compressor) syncFlush() error ***REMOVED***
	d.sync = true
	if d.err != nil ***REMOVED***
		return d.err
	***REMOVED***
	d.step(d)
	if d.err == nil ***REMOVED***
		d.w.writeStoredHeader(0, false)
		d.w.flush()
		d.err = d.w.err
	***REMOVED***
	d.sync = false
	return d.err
***REMOVED***

func (d *compressor) init(w io.Writer, level int) (err error) ***REMOVED***
	d.w = newHuffmanBitWriter(w)

	switch ***REMOVED***
	case level == NoCompression:
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).store
	case level == ConstantCompression:
		d.w.logNewTablePenalty = 10
		d.window = make([]byte, 32<<10)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).storeHuff
	case level == DefaultCompression:
		level = 5
		fallthrough
	case level >= 1 && level <= 6:
		d.w.logNewTablePenalty = 7
		d.fast = newFastEnc(level)
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).storeFast
	case 7 <= level && level <= 9:
		d.w.logNewTablePenalty = 8
		d.state = &advancedState***REMOVED******REMOVED***
		d.compressionLevel = levels[level]
		d.initDeflate()
		d.fill = (*compressor).fillDeflate
		d.step = (*compressor).deflateLazy
	default:
		return fmt.Errorf("flate: invalid compression level %d: want value in range [-2, 9]", level)
	***REMOVED***
	d.level = level
	return nil
***REMOVED***

// reset the state of the compressor.
func (d *compressor) reset(w io.Writer) ***REMOVED***
	d.w.reset(w)
	d.sync = false
	d.err = nil
	// We only need to reset a few things for Snappy.
	if d.fast != nil ***REMOVED***
		d.fast.Reset()
		d.windowEnd = 0
		d.tokens.Reset()
		return
	***REMOVED***
	switch d.compressionLevel.chain ***REMOVED***
	case 0:
		// level was NoCompression or ConstantCompresssion.
		d.windowEnd = 0
	default:
		s := d.state
		s.chainHead = -1
		for i := range s.hashHead ***REMOVED***
			s.hashHead[i] = 0
		***REMOVED***
		for i := range s.hashPrev ***REMOVED***
			s.hashPrev[i] = 0
		***REMOVED***
		s.hashOffset = 1
		s.index, d.windowEnd = 0, 0
		d.blockStart, d.byteAvailable = 0, false
		d.tokens.Reset()
		s.length = minMatchLength - 1
		s.offset = 0
		s.ii = 0
		s.maxInsertIndex = 0
	***REMOVED***
***REMOVED***

func (d *compressor) close() error ***REMOVED***
	if d.err != nil ***REMOVED***
		return d.err
	***REMOVED***
	d.sync = true
	d.step(d)
	if d.err != nil ***REMOVED***
		return d.err
	***REMOVED***
	if d.w.writeStoredHeader(0, true); d.w.err != nil ***REMOVED***
		return d.w.err
	***REMOVED***
	d.w.flush()
	d.w.reset(nil)
	return d.w.err
***REMOVED***

// NewWriter returns a new Writer compressing data at the given level.
// Following zlib, levels range from 1 (BestSpeed) to 9 (BestCompression);
// higher levels typically run slower but compress more.
// Level 0 (NoCompression) does not attempt any compression; it only adds the
// necessary DEFLATE framing.
// Level -1 (DefaultCompression) uses the default compression level.
// Level -2 (ConstantCompression) will use Huffman compression only, giving
// a very fast compression for all types of input, but sacrificing considerable
// compression efficiency.
//
// If level is in the range [-2, 9] then the error returned will be nil.
// Otherwise the error returned will be non-nil.
func NewWriter(w io.Writer, level int) (*Writer, error) ***REMOVED***
	var dw Writer
	if err := dw.d.init(w, level); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &dw, nil
***REMOVED***

// NewWriterDict is like NewWriter but initializes the new
// Writer with a preset dictionary.  The returned Writer behaves
// as if the dictionary had been written to it without producing
// any compressed output.  The compressed data written to w
// can only be decompressed by a Reader initialized with the
// same dictionary.
func NewWriterDict(w io.Writer, level int, dict []byte) (*Writer, error) ***REMOVED***
	zw, err := NewWriter(w, level)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	zw.d.fillWindow(dict)
	zw.dict = append(zw.dict, dict...) // duplicate dictionary for Reset method.
	return zw, err
***REMOVED***

// A Writer takes data written to it and writes the compressed
// form of that data to an underlying writer (see NewWriter).
type Writer struct ***REMOVED***
	d    compressor
	dict []byte
***REMOVED***

// Write writes data to w, which will eventually write the
// compressed form of data to its underlying writer.
func (w *Writer) Write(data []byte) (n int, err error) ***REMOVED***
	return w.d.write(data)
***REMOVED***

// Flush flushes any pending data to the underlying writer.
// It is useful mainly in compressed network protocols, to ensure that
// a remote reader has enough data to reconstruct a packet.
// Flush does not return until the data has been written.
// Calling Flush when there is no pending data still causes the Writer
// to emit a sync marker of at least 4 bytes.
// If the underlying writer returns an error, Flush returns that error.
//
// In the terminology of the zlib library, Flush is equivalent to Z_SYNC_FLUSH.
func (w *Writer) Flush() error ***REMOVED***
	// For more about flushing:
	// http://www.bolet.org/~pornin/deflate-flush.html
	return w.d.syncFlush()
***REMOVED***

// Close flushes and closes the writer.
func (w *Writer) Close() error ***REMOVED***
	return w.d.close()
***REMOVED***

// Reset discards the writer's state and makes it equivalent to
// the result of NewWriter or NewWriterDict called with dst
// and w's level and dictionary.
func (w *Writer) Reset(dst io.Writer) ***REMOVED***
	if len(w.dict) > 0 ***REMOVED***
		// w was created with NewWriterDict
		w.d.reset(dst)
		if dst != nil ***REMOVED***
			w.d.fillWindow(w.dict)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// w was created with NewWriter
		w.d.reset(dst)
	***REMOVED***
***REMOVED***

// ResetDict discards the writer's state and makes it equivalent to
// the result of NewWriter or NewWriterDict called with dst
// and w's level, but sets a specific dictionary.
func (w *Writer) ResetDict(dst io.Writer, dict []byte) ***REMOVED***
	w.dict = dict
	w.d.reset(dst)
	w.d.fillWindow(w.dict)
***REMOVED***
