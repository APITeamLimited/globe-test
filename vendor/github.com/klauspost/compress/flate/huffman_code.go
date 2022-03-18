// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

import (
	"math"
	"math/bits"
)

const (
	maxBitsLimit = 16
	// number of valid literals
	literalCount = 286
)

// hcode is a huffman code with a bit code and bit length.
type hcode struct ***REMOVED***
	code uint16
	len  uint8
***REMOVED***

type huffmanEncoder struct ***REMOVED***
	codes    []hcode
	bitCount [17]int32

	// Allocate a reusable buffer with the longest possible frequency table.
	// Possible lengths are codegenCodeCount, offsetCodeCount and literalCount.
	// The largest of these is literalCount, so we allocate for that case.
	freqcache [literalCount + 1]literalNode
***REMOVED***

type literalNode struct ***REMOVED***
	literal uint16
	freq    uint16
***REMOVED***

// A levelInfo describes the state of the constructed tree for a given depth.
type levelInfo struct ***REMOVED***
	// Our level.  for better printing
	level int32

	// The frequency of the last node at this level
	lastFreq int32

	// The frequency of the next character to add to this level
	nextCharFreq int32

	// The frequency of the next pair (from level below) to add to this level.
	// Only valid if the "needed" value of the next lower level is 0.
	nextPairFreq int32

	// The number of chains remaining to generate for this level before moving
	// up to the next level
	needed int32
***REMOVED***

// set sets the code and length of an hcode.
func (h *hcode) set(code uint16, length uint8) ***REMOVED***
	h.len = length
	h.code = code
***REMOVED***

func reverseBits(number uint16, bitLength byte) uint16 ***REMOVED***
	return bits.Reverse16(number << ((16 - bitLength) & 15))
***REMOVED***

func maxNode() literalNode ***REMOVED*** return literalNode***REMOVED***math.MaxUint16, math.MaxUint16***REMOVED*** ***REMOVED***

func newHuffmanEncoder(size int) *huffmanEncoder ***REMOVED***
	// Make capacity to next power of two.
	c := uint(bits.Len32(uint32(size - 1)))
	return &huffmanEncoder***REMOVED***codes: make([]hcode, size, 1<<c)***REMOVED***
***REMOVED***

// Generates a HuffmanCode corresponding to the fixed literal table
func generateFixedLiteralEncoding() *huffmanEncoder ***REMOVED***
	h := newHuffmanEncoder(literalCount)
	codes := h.codes
	var ch uint16
	for ch = 0; ch < literalCount; ch++ ***REMOVED***
		var bits uint16
		var size uint8
		switch ***REMOVED***
		case ch < 144:
			// size 8, 000110000  .. 10111111
			bits = ch + 48
			size = 8
		case ch < 256:
			// size 9, 110010000 .. 111111111
			bits = ch + 400 - 144
			size = 9
		case ch < 280:
			// size 7, 0000000 .. 0010111
			bits = ch - 256
			size = 7
		default:
			// size 8, 11000000 .. 11000111
			bits = ch + 192 - 280
			size = 8
		***REMOVED***
		codes[ch] = hcode***REMOVED***code: reverseBits(bits, size), len: size***REMOVED***
	***REMOVED***
	return h
***REMOVED***

func generateFixedOffsetEncoding() *huffmanEncoder ***REMOVED***
	h := newHuffmanEncoder(30)
	codes := h.codes
	for ch := range codes ***REMOVED***
		codes[ch] = hcode***REMOVED***code: reverseBits(uint16(ch), 5), len: 5***REMOVED***
	***REMOVED***
	return h
***REMOVED***

var fixedLiteralEncoding = generateFixedLiteralEncoding()
var fixedOffsetEncoding = generateFixedOffsetEncoding()

func (h *huffmanEncoder) bitLength(freq []uint16) int ***REMOVED***
	var total int
	for i, f := range freq ***REMOVED***
		if f != 0 ***REMOVED***
			total += int(f) * int(h.codes[i].len)
		***REMOVED***
	***REMOVED***
	return total
***REMOVED***

func (h *huffmanEncoder) bitLengthRaw(b []byte) int ***REMOVED***
	var total int
	for _, f := range b ***REMOVED***
		total += int(h.codes[f].len)
	***REMOVED***
	return total
***REMOVED***

// canReuseBits returns the number of bits or math.MaxInt32 if the encoder cannot be reused.
func (h *huffmanEncoder) canReuseBits(freq []uint16) int ***REMOVED***
	var total int
	for i, f := range freq ***REMOVED***
		if f != 0 ***REMOVED***
			code := h.codes[i]
			if code.len == 0 ***REMOVED***
				return math.MaxInt32
			***REMOVED***
			total += int(f) * int(code.len)
		***REMOVED***
	***REMOVED***
	return total
***REMOVED***

// Return the number of literals assigned to each bit size in the Huffman encoding
//
// This method is only called when list.length >= 3
// The cases of 0, 1, and 2 literals are handled by special case code.
//
// list  An array of the literals with non-zero frequencies
//             and their associated frequencies. The array is in order of increasing
//             frequency, and has as its last element a special element with frequency
//             MaxInt32
// maxBits     The maximum number of bits that should be used to encode any literal.
//             Must be less than 16.
// return      An integer array in which array[i] indicates the number of literals
//             that should be encoded in i bits.
func (h *huffmanEncoder) bitCounts(list []literalNode, maxBits int32) []int32 ***REMOVED***
	if maxBits >= maxBitsLimit ***REMOVED***
		panic("flate: maxBits too large")
	***REMOVED***
	n := int32(len(list))
	list = list[0 : n+1]
	list[n] = maxNode()

	// The tree can't have greater depth than n - 1, no matter what. This
	// saves a little bit of work in some small cases
	if maxBits > n-1 ***REMOVED***
		maxBits = n - 1
	***REMOVED***

	// Create information about each of the levels.
	// A bogus "Level 0" whose sole purpose is so that
	// level1.prev.needed==0.  This makes level1.nextPairFreq
	// be a legitimate value that never gets chosen.
	var levels [maxBitsLimit]levelInfo
	// leafCounts[i] counts the number of literals at the left
	// of ancestors of the rightmost node at level i.
	// leafCounts[i][j] is the number of literals at the left
	// of the level j ancestor.
	var leafCounts [maxBitsLimit][maxBitsLimit]int32

	// Descending to only have 1 bounds check.
	l2f := int32(list[2].freq)
	l1f := int32(list[1].freq)
	l0f := int32(list[0].freq) + int32(list[1].freq)

	for level := int32(1); level <= maxBits; level++ ***REMOVED***
		// For every level, the first two items are the first two characters.
		// We initialize the levels as if we had already figured this out.
		levels[level] = levelInfo***REMOVED***
			level:        level,
			lastFreq:     l1f,
			nextCharFreq: l2f,
			nextPairFreq: l0f,
		***REMOVED***
		leafCounts[level][level] = 2
		if level == 1 ***REMOVED***
			levels[level].nextPairFreq = math.MaxInt32
		***REMOVED***
	***REMOVED***

	// We need a total of 2*n - 2 items at top level and have already generated 2.
	levels[maxBits].needed = 2*n - 4

	level := uint32(maxBits)
	for level < 16 ***REMOVED***
		l := &levels[level]
		if l.nextPairFreq == math.MaxInt32 && l.nextCharFreq == math.MaxInt32 ***REMOVED***
			// We've run out of both leafs and pairs.
			// End all calculations for this level.
			// To make sure we never come back to this level or any lower level,
			// set nextPairFreq impossibly large.
			l.needed = 0
			levels[level+1].nextPairFreq = math.MaxInt32
			level++
			continue
		***REMOVED***

		prevFreq := l.lastFreq
		if l.nextCharFreq < l.nextPairFreq ***REMOVED***
			// The next item on this row is a leaf node.
			n := leafCounts[level][level] + 1
			l.lastFreq = l.nextCharFreq
			// Lower leafCounts are the same of the previous node.
			leafCounts[level][level] = n
			e := list[n]
			if e.literal < math.MaxUint16 ***REMOVED***
				l.nextCharFreq = int32(e.freq)
			***REMOVED*** else ***REMOVED***
				l.nextCharFreq = math.MaxInt32
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// The next item on this row is a pair from the previous row.
			// nextPairFreq isn't valid until we generate two
			// more values in the level below
			l.lastFreq = l.nextPairFreq
			// Take leaf counts from the lower level, except counts[level] remains the same.
			if true ***REMOVED***
				save := leafCounts[level][level]
				leafCounts[level] = leafCounts[level-1]
				leafCounts[level][level] = save
			***REMOVED*** else ***REMOVED***
				copy(leafCounts[level][:level], leafCounts[level-1][:level])
			***REMOVED***
			levels[l.level-1].needed = 2
		***REMOVED***

		if l.needed--; l.needed == 0 ***REMOVED***
			// We've done everything we need to do for this level.
			// Continue calculating one level up. Fill in nextPairFreq
			// of that level with the sum of the two nodes we've just calculated on
			// this level.
			if l.level == maxBits ***REMOVED***
				// All done!
				break
			***REMOVED***
			levels[l.level+1].nextPairFreq = prevFreq + l.lastFreq
			level++
		***REMOVED*** else ***REMOVED***
			// If we stole from below, move down temporarily to replenish it.
			for levels[level-1].needed > 0 ***REMOVED***
				level--
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Somethings is wrong if at the end, the top level is null or hasn't used
	// all of the leaves.
	if leafCounts[maxBits][maxBits] != n ***REMOVED***
		panic("leafCounts[maxBits][maxBits] != n")
	***REMOVED***

	bitCount := h.bitCount[:maxBits+1]
	bits := 1
	counts := &leafCounts[maxBits]
	for level := maxBits; level > 0; level-- ***REMOVED***
		// chain.leafCount gives the number of literals requiring at least "bits"
		// bits to encode.
		bitCount[bits] = counts[level] - counts[level-1]
		bits++
	***REMOVED***
	return bitCount
***REMOVED***

// Look at the leaves and assign them a bit count and an encoding as specified
// in RFC 1951 3.2.2
func (h *huffmanEncoder) assignEncodingAndSize(bitCount []int32, list []literalNode) ***REMOVED***
	code := uint16(0)
	for n, bits := range bitCount ***REMOVED***
		code <<= 1
		if n == 0 || bits == 0 ***REMOVED***
			continue
		***REMOVED***
		// The literals list[len(list)-bits] .. list[len(list)-bits]
		// are encoded using "bits" bits, and get the values
		// code, code + 1, ....  The code values are
		// assigned in literal order (not frequency order).
		chunk := list[len(list)-int(bits):]

		sortByLiteral(chunk)
		for _, node := range chunk ***REMOVED***
			h.codes[node.literal] = hcode***REMOVED***code: reverseBits(code, uint8(n)), len: uint8(n)***REMOVED***
			code++
		***REMOVED***
		list = list[0 : len(list)-int(bits)]
	***REMOVED***
***REMOVED***

// Update this Huffman Code object to be the minimum code for the specified frequency count.
//
// freq  An array of frequencies, in which frequency[i] gives the frequency of literal i.
// maxBits  The maximum number of bits to use for any literal.
func (h *huffmanEncoder) generate(freq []uint16, maxBits int32) ***REMOVED***
	list := h.freqcache[:len(freq)+1]
	codes := h.codes[:len(freq)]
	// Number of non-zero literals
	count := 0
	// Set list to be the set of all non-zero literals and their frequencies
	for i, f := range freq ***REMOVED***
		if f != 0 ***REMOVED***
			list[count] = literalNode***REMOVED***uint16(i), f***REMOVED***
			count++
		***REMOVED*** else ***REMOVED***
			codes[i].len = 0
		***REMOVED***
	***REMOVED***
	list[count] = literalNode***REMOVED******REMOVED***

	list = list[:count]
	if count <= 2 ***REMOVED***
		// Handle the small cases here, because they are awkward for the general case code. With
		// two or fewer literals, everything has bit length 1.
		for i, node := range list ***REMOVED***
			// "list" is in order of increasing literal value.
			h.codes[node.literal].set(uint16(i), 1)
		***REMOVED***
		return
	***REMOVED***
	sortByFreq(list)

	// Get the number of literals for each bit count
	bitCount := h.bitCounts(list, maxBits)
	// And do the assignment
	h.assignEncodingAndSize(bitCount, list)
***REMOVED***

// atLeastOne clamps the result between 1 and 15.
func atLeastOne(v float32) float32 ***REMOVED***
	if v < 1 ***REMOVED***
		return 1
	***REMOVED***
	if v > 15 ***REMOVED***
		return 15
	***REMOVED***
	return v
***REMOVED***

// Unassigned values are assigned '1' in the histogram.
func fillHist(b []uint16) ***REMOVED***
	for i, v := range b ***REMOVED***
		if v == 0 ***REMOVED***
			b[i] = 1
		***REMOVED***
	***REMOVED***
***REMOVED***

func histogram(b []byte, h []uint16, fill bool) ***REMOVED***
	h = h[:256]
	for _, t := range b ***REMOVED***
		h[t]++
	***REMOVED***
	if fill ***REMOVED***
		fillHist(h)
	***REMOVED***
***REMOVED***
