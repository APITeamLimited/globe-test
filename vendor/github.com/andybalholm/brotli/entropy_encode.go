package brotli

import "math"

/* Copyright 2010 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Entropy encoding (Huffman) utilities. */

/* A node of a Huffman tree. */
type huffmanTree struct ***REMOVED***
	total_count_          uint32
	index_left_           int16
	index_right_or_value_ int16
***REMOVED***

func initHuffmanTree(self *huffmanTree, count uint32, left int16, right int16) ***REMOVED***
	self.total_count_ = count
	self.index_left_ = left
	self.index_right_or_value_ = right
***REMOVED***

/* Input size optimized Shell sort. */
type huffmanTreeComparator func(*huffmanTree, *huffmanTree) bool

var sortHuffmanTreeItems_gaps = []uint***REMOVED***132, 57, 23, 10, 4, 1***REMOVED***

func sortHuffmanTreeItems(items []huffmanTree, n uint, comparator huffmanTreeComparator) ***REMOVED***
	if n < 13 ***REMOVED***
		/* Insertion sort. */
		var i uint
		for i = 1; i < n; i++ ***REMOVED***
			var tmp huffmanTree = items[i]
			var k uint = i
			var j uint = i - 1
			for comparator(&tmp, &items[j]) ***REMOVED***
				items[k] = items[j]
				k = j
				tmp10 := j
				j--
				if tmp10 == 0 ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			items[k] = tmp
		***REMOVED***

		return
	***REMOVED*** else ***REMOVED***
		var g int
		if n < 57 ***REMOVED***
			g = 2
		***REMOVED*** else ***REMOVED***
			g = 0
		***REMOVED***
		for ; g < 6; g++ ***REMOVED***
			var gap uint = sortHuffmanTreeItems_gaps[g]
			var i uint
			for i = gap; i < n; i++ ***REMOVED***
				var j uint = i
				var tmp huffmanTree = items[i]
				for ; j >= gap && comparator(&tmp, &items[j-gap]); j -= gap ***REMOVED***
					items[j] = items[j-gap]
				***REMOVED***

				items[j] = tmp
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

/* Returns 1 if assignment of depths succeeded, otherwise 0. */
func setDepth(p0 int, pool []huffmanTree, depth []byte, max_depth int) bool ***REMOVED***
	var stack [16]int
	var level int = 0
	var p int = p0
	assert(max_depth <= 15)
	stack[0] = -1
	for ***REMOVED***
		if pool[p].index_left_ >= 0 ***REMOVED***
			level++
			if level > max_depth ***REMOVED***
				return false
			***REMOVED***
			stack[level] = int(pool[p].index_right_or_value_)
			p = int(pool[p].index_left_)
			continue
		***REMOVED*** else ***REMOVED***
			depth[pool[p].index_right_or_value_] = byte(level)
		***REMOVED***

		for level >= 0 && stack[level] == -1 ***REMOVED***
			level--
		***REMOVED***
		if level < 0 ***REMOVED***
			return true
		***REMOVED***
		p = stack[level]
		stack[level] = -1
	***REMOVED***
***REMOVED***

/* Sort the root nodes, least popular first. */
func sortHuffmanTree(v0 *huffmanTree, v1 *huffmanTree) bool ***REMOVED***
	if v0.total_count_ != v1.total_count_ ***REMOVED***
		return v0.total_count_ < v1.total_count_
	***REMOVED***

	return v0.index_right_or_value_ > v1.index_right_or_value_
***REMOVED***

/* This function will create a Huffman tree.

   The catch here is that the tree cannot be arbitrarily deep.
   Brotli specifies a maximum depth of 15 bits for "code trees"
   and 7 bits for "code length code trees."

   count_limit is the value that is to be faked as the minimum value
   and this minimum value is raised until the tree matches the
   maximum length requirement.

   This algorithm is not of excellent performance for very long data blocks,
   especially when population counts are longer than 2**tree_limit, but
   we are not planning to use this with extremely long blocks.

   See http://en.wikipedia.org/wiki/Huffman_coding */
func createHuffmanTree(data []uint32, length uint, tree_limit int, tree []huffmanTree, depth []byte) ***REMOVED***
	var count_limit uint32
	var sentinel huffmanTree
	initHuffmanTree(&sentinel, math.MaxUint32, -1, -1)

	/* For block sizes below 64 kB, we never need to do a second iteration
	   of this loop. Probably all of our block sizes will be smaller than
	   that, so this loop is mostly of academic interest. If we actually
	   would need this, we would be better off with the Katajainen algorithm. */
	for count_limit = 1; ; count_limit *= 2 ***REMOVED***
		var n uint = 0
		var i uint
		var j uint
		var k uint
		for i = length; i != 0; ***REMOVED***
			i--
			if data[i] != 0 ***REMOVED***
				var count uint32 = brotli_max_uint32_t(data[i], count_limit)
				initHuffmanTree(&tree[n], count, -1, int16(i))
				n++
			***REMOVED***
		***REMOVED***

		if n == 1 ***REMOVED***
			depth[tree[0].index_right_or_value_] = 1 /* Only one element. */
			break
		***REMOVED***

		sortHuffmanTreeItems(tree, n, huffmanTreeComparator(sortHuffmanTree))

		/* The nodes are:
		   [0, n): the sorted leaf nodes that we start with.
		   [n]: we add a sentinel here.
		   [n + 1, 2n): new parent nodes are added here, starting from
		                (n+1). These are naturally in ascending order.
		   [2n]: we add a sentinel at the end as well.
		   There will be (2n+1) elements at the end. */
		tree[n] = sentinel

		tree[n+1] = sentinel

		i = 0     /* Points to the next leaf node. */
		j = n + 1 /* Points to the next non-leaf node. */
		for k = n - 1; k != 0; k-- ***REMOVED***
			var left uint
			var right uint
			if tree[i].total_count_ <= tree[j].total_count_ ***REMOVED***
				left = i
				i++
			***REMOVED*** else ***REMOVED***
				left = j
				j++
			***REMOVED***

			if tree[i].total_count_ <= tree[j].total_count_ ***REMOVED***
				right = i
				i++
			***REMOVED*** else ***REMOVED***
				right = j
				j++
			***REMOVED***
			***REMOVED***
				/* The sentinel node becomes the parent node. */
				var j_end uint = 2*n - k
				tree[j_end].total_count_ = tree[left].total_count_ + tree[right].total_count_
				tree[j_end].index_left_ = int16(left)
				tree[j_end].index_right_or_value_ = int16(right)

				/* Add back the last sentinel node. */
				tree[j_end+1] = sentinel
			***REMOVED***
		***REMOVED***

		if setDepth(int(2*n-1), tree[0:], depth, tree_limit) ***REMOVED***
			/* We need to pack the Huffman tree in tree_limit bits. If this was not
			   successful, add fake entities to the lowest values and retry. */
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func reverse(v []byte, start uint, end uint) ***REMOVED***
	end--
	for start < end ***REMOVED***
		var tmp byte = v[start]
		v[start] = v[end]
		v[end] = tmp
		start++
		end--
	***REMOVED***
***REMOVED***

func writeHuffmanTreeRepetitions(previous_value byte, value byte, repetitions uint, tree_size *uint, tree []byte, extra_bits_data []byte) ***REMOVED***
	assert(repetitions > 0)
	if previous_value != value ***REMOVED***
		tree[*tree_size] = value
		extra_bits_data[*tree_size] = 0
		(*tree_size)++
		repetitions--
	***REMOVED***

	if repetitions == 7 ***REMOVED***
		tree[*tree_size] = value
		extra_bits_data[*tree_size] = 0
		(*tree_size)++
		repetitions--
	***REMOVED***

	if repetitions < 3 ***REMOVED***
		var i uint
		for i = 0; i < repetitions; i++ ***REMOVED***
			tree[*tree_size] = value
			extra_bits_data[*tree_size] = 0
			(*tree_size)++
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var start uint = *tree_size
		repetitions -= 3
		for ***REMOVED***
			tree[*tree_size] = repeatPreviousCodeLength
			extra_bits_data[*tree_size] = byte(repetitions & 0x3)
			(*tree_size)++
			repetitions >>= 2
			if repetitions == 0 ***REMOVED***
				break
			***REMOVED***

			repetitions--
		***REMOVED***

		reverse(tree, start, *tree_size)
		reverse(extra_bits_data, start, *tree_size)
	***REMOVED***
***REMOVED***

func writeHuffmanTreeRepetitionsZeros(repetitions uint, tree_size *uint, tree []byte, extra_bits_data []byte) ***REMOVED***
	if repetitions == 11 ***REMOVED***
		tree[*tree_size] = 0
		extra_bits_data[*tree_size] = 0
		(*tree_size)++
		repetitions--
	***REMOVED***

	if repetitions < 3 ***REMOVED***
		var i uint
		for i = 0; i < repetitions; i++ ***REMOVED***
			tree[*tree_size] = 0
			extra_bits_data[*tree_size] = 0
			(*tree_size)++
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var start uint = *tree_size
		repetitions -= 3
		for ***REMOVED***
			tree[*tree_size] = repeatZeroCodeLength
			extra_bits_data[*tree_size] = byte(repetitions & 0x7)
			(*tree_size)++
			repetitions >>= 3
			if repetitions == 0 ***REMOVED***
				break
			***REMOVED***

			repetitions--
		***REMOVED***

		reverse(tree, start, *tree_size)
		reverse(extra_bits_data, start, *tree_size)
	***REMOVED***
***REMOVED***

/* Change the population counts in a way that the consequent
   Huffman tree compression, especially its RLE-part will be more
   likely to compress this data more efficiently.

   length contains the size of the histogram.
   counts contains the population counts.
   good_for_rle is a buffer of at least length size */
func optimizeHuffmanCountsForRLE(length uint, counts []uint32, good_for_rle []byte) ***REMOVED***
	var nonzero_count uint = 0
	var stride uint
	var limit uint
	var sum uint
	var streak_limit uint = 1240
	var i uint
	/* Let's make the Huffman code more compatible with RLE encoding. */
	for i = 0; i < length; i++ ***REMOVED***
		if counts[i] != 0 ***REMOVED***
			nonzero_count++
		***REMOVED***
	***REMOVED***

	if nonzero_count < 16 ***REMOVED***
		return
	***REMOVED***

	for length != 0 && counts[length-1] == 0 ***REMOVED***
		length--
	***REMOVED***

	if length == 0 ***REMOVED***
		return /* All zeros. */
	***REMOVED***

	/* Now counts[0..length - 1] does not have trailing zeros. */
	***REMOVED***
		var nonzeros uint = 0
		var smallest_nonzero uint32 = 1 << 30
		for i = 0; i < length; i++ ***REMOVED***
			if counts[i] != 0 ***REMOVED***
				nonzeros++
				if smallest_nonzero > counts[i] ***REMOVED***
					smallest_nonzero = counts[i]
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if nonzeros < 5 ***REMOVED***
			/* Small histogram will model it well. */
			return
		***REMOVED***

		if smallest_nonzero < 4 ***REMOVED***
			var zeros uint = length - nonzeros
			if zeros < 6 ***REMOVED***
				for i = 1; i < length-1; i++ ***REMOVED***
					if counts[i-1] != 0 && counts[i] == 0 && counts[i+1] != 0 ***REMOVED***
						counts[i] = 1
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if nonzeros < 28 ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	/* 2) Let's mark all population counts that already can be encoded
	   with an RLE code. */
	for i := 0; i < int(length); i++ ***REMOVED***
		good_for_rle[i] = 0
	***REMOVED***
	***REMOVED***
		var symbol uint32 = counts[0]
		/* Let's not spoil any of the existing good RLE codes.
		   Mark any seq of 0's that is longer as 5 as a good_for_rle.
		   Mark any seq of non-0's that is longer as 7 as a good_for_rle. */

		var step uint = 0
		for i = 0; i <= length; i++ ***REMOVED***
			if i == length || counts[i] != symbol ***REMOVED***
				if (symbol == 0 && step >= 5) || (symbol != 0 && step >= 7) ***REMOVED***
					var k uint
					for k = 0; k < step; k++ ***REMOVED***
						good_for_rle[i-k-1] = 1
					***REMOVED***
				***REMOVED***

				step = 1
				if i != length ***REMOVED***
					symbol = counts[i]
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				step++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	/* 3) Let's replace those population counts that lead to more RLE codes.
	   Math here is in 24.8 fixed point representation. */
	stride = 0

	limit = uint(256*(counts[0]+counts[1]+counts[2])/3 + 420)
	sum = 0
	for i = 0; i <= length; i++ ***REMOVED***
		if i == length || good_for_rle[i] != 0 || (i != 0 && good_for_rle[i-1] != 0) || (256*counts[i]-uint32(limit)+uint32(streak_limit)) >= uint32(2*streak_limit) ***REMOVED***
			if stride >= 4 || (stride >= 3 && sum == 0) ***REMOVED***
				var k uint
				var count uint = (sum + stride/2) / stride
				/* The stride must end, collapse what we have, if we have enough (4). */
				if count == 0 ***REMOVED***
					count = 1
				***REMOVED***

				if sum == 0 ***REMOVED***
					/* Don't make an all zeros stride to be upgraded to ones. */
					count = 0
				***REMOVED***

				for k = 0; k < stride; k++ ***REMOVED***
					/* We don't want to change value at counts[i],
					   that is already belonging to the next stride. Thus - 1. */
					counts[i-k-1] = uint32(count)
				***REMOVED***
			***REMOVED***

			stride = 0
			sum = 0
			if i < length-2 ***REMOVED***
				/* All interesting strides have a count of at least 4, */
				/* at least when non-zeros. */
				limit = uint(256*(counts[i]+counts[i+1]+counts[i+2])/3 + 420)
			***REMOVED*** else if i < length ***REMOVED***
				limit = uint(256 * counts[i])
			***REMOVED*** else ***REMOVED***
				limit = 0
			***REMOVED***
		***REMOVED***

		stride++
		if i != length ***REMOVED***
			sum += uint(counts[i])
			if stride >= 4 ***REMOVED***
				limit = (256*sum + stride/2) / stride
			***REMOVED***

			if stride == 4 ***REMOVED***
				limit += 120
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func decideOverRLEUse(depth []byte, length uint, use_rle_for_non_zero *bool, use_rle_for_zero *bool) ***REMOVED***
	var total_reps_zero uint = 0
	var total_reps_non_zero uint = 0
	var count_reps_zero uint = 1
	var count_reps_non_zero uint = 1
	var i uint
	for i = 0; i < length; ***REMOVED***
		var value byte = depth[i]
		var reps uint = 1
		var k uint
		for k = i + 1; k < length && depth[k] == value; k++ ***REMOVED***
			reps++
		***REMOVED***

		if reps >= 3 && value == 0 ***REMOVED***
			total_reps_zero += reps
			count_reps_zero++
		***REMOVED***

		if reps >= 4 && value != 0 ***REMOVED***
			total_reps_non_zero += reps
			count_reps_non_zero++
		***REMOVED***

		i += reps
	***REMOVED***

	*use_rle_for_non_zero = total_reps_non_zero > count_reps_non_zero*2
	*use_rle_for_zero = total_reps_zero > count_reps_zero*2
***REMOVED***

/* Write a Huffman tree from bit depths into the bit-stream representation
   of a Huffman tree. The generated Huffman tree is to be compressed once
   more using a Huffman tree */
func writeHuffmanTree(depth []byte, length uint, tree_size *uint, tree []byte, extra_bits_data []byte) ***REMOVED***
	var previous_value byte = initialRepeatedCodeLength
	var i uint
	var use_rle_for_non_zero bool = false
	var use_rle_for_zero bool = false
	var new_length uint = length
	/* Throw away trailing zeros. */
	for i = 0; i < length; i++ ***REMOVED***
		if depth[length-i-1] == 0 ***REMOVED***
			new_length--
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	/* First gather statistics on if it is a good idea to do RLE. */
	if length > 50 ***REMOVED***
		/* Find RLE coding for longer codes.
		   Shorter codes seem not to benefit from RLE. */
		decideOverRLEUse(depth, new_length, &use_rle_for_non_zero, &use_rle_for_zero)
	***REMOVED***

	/* Actual RLE coding. */
	for i = 0; i < new_length; ***REMOVED***
		var value byte = depth[i]
		var reps uint = 1
		if (value != 0 && use_rle_for_non_zero) || (value == 0 && use_rle_for_zero) ***REMOVED***
			var k uint
			for k = i + 1; k < new_length && depth[k] == value; k++ ***REMOVED***
				reps++
			***REMOVED***
		***REMOVED***

		if value == 0 ***REMOVED***
			writeHuffmanTreeRepetitionsZeros(reps, tree_size, tree, extra_bits_data)
		***REMOVED*** else ***REMOVED***
			writeHuffmanTreeRepetitions(previous_value, value, reps, tree_size, tree, extra_bits_data)
			previous_value = value
		***REMOVED***

		i += reps
	***REMOVED***
***REMOVED***

var reverseBits_kLut = [16]uint***REMOVED***
	0x00,
	0x08,
	0x04,
	0x0C,
	0x02,
	0x0A,
	0x06,
	0x0E,
	0x01,
	0x09,
	0x05,
	0x0D,
	0x03,
	0x0B,
	0x07,
	0x0F,
***REMOVED***

func reverseBits(num_bits uint, bits uint16) uint16 ***REMOVED***
	var retval uint = reverseBits_kLut[bits&0x0F]
	var i uint
	for i = 4; i < num_bits; i += 4 ***REMOVED***
		retval <<= 4
		bits = uint16(bits >> 4)
		retval |= reverseBits_kLut[bits&0x0F]
	***REMOVED***

	retval >>= ((0 - num_bits) & 0x03)
	return uint16(retval)
***REMOVED***

/* 0..15 are values for bits */
const maxHuffmanBits = 16

/* Get the actual bit values for a tree of bit depths. */
func convertBitDepthsToSymbols(depth []byte, len uint, bits []uint16) ***REMOVED***
	var bl_count = [maxHuffmanBits]uint16***REMOVED***0***REMOVED***
	var next_code [maxHuffmanBits]uint16
	var i uint
	/* In Brotli, all bit depths are [1..15]
	   0 bit depth means that the symbol does not exist. */

	var code int = 0
	for i = 0; i < len; i++ ***REMOVED***
		bl_count[depth[i]]++
	***REMOVED***

	bl_count[0] = 0
	next_code[0] = 0
	for i = 1; i < maxHuffmanBits; i++ ***REMOVED***
		code = (code + int(bl_count[i-1])) << 1
		next_code[i] = uint16(code)
	***REMOVED***

	for i = 0; i < len; i++ ***REMOVED***
		if depth[i] != 0 ***REMOVED***
			bits[i] = reverseBits(uint(depth[i]), next_code[depth[i]])
			next_code[depth[i]]++
		***REMOVED***
	***REMOVED***
***REMOVED***
