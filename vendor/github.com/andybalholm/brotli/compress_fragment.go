package brotli

import "encoding/binary"

/* Copyright 2015 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Function for fast encoding of an input fragment, independently from the input
   history. This function uses one-pass processing: when we find a backward
   match, we immediately emit the corresponding command and literal codes to
   the bit stream.

   Adapted from the CompressFragment() function in
   https://github.com/google/snappy/blob/master/snappy.cc */

const maxDistance_compress_fragment = 262128

func hash5(p []byte, shift uint) uint32 ***REMOVED***
	var h uint64 = (binary.LittleEndian.Uint64(p) << 24) * uint64(kHashMul32)
	return uint32(h >> shift)
***REMOVED***

func hashBytesAtOffset5(v uint64, offset int, shift uint) uint32 ***REMOVED***
	assert(offset >= 0)
	assert(offset <= 3)
	***REMOVED***
		var h uint64 = ((v >> uint(8*offset)) << 24) * uint64(kHashMul32)
		return uint32(h >> shift)
	***REMOVED***
***REMOVED***

func isMatch5(p1 []byte, p2 []byte) bool ***REMOVED***
	return binary.LittleEndian.Uint32(p1) == binary.LittleEndian.Uint32(p2) &&
		p1[4] == p2[4]
***REMOVED***

/* Builds a literal prefix code into "depths" and "bits" based on the statistics
   of the "input" string and stores it into the bit stream.
   Note that the prefix code here is built from the pre-LZ77 input, therefore
   we can only approximate the statistics of the actual literal stream.
   Moreover, for long inputs we build a histogram from a sample of the input
   and thus have to assign a non-zero depth for each literal.
   Returns estimated compression ratio millibytes/char for encoding given input
   with generated code. */
func buildAndStoreLiteralPrefixCode(input []byte, input_size uint, depths []byte, bits []uint16, bw *bitWriter) uint ***REMOVED***
	var histogram = [256]uint32***REMOVED***0***REMOVED***
	var histogram_total uint
	var i uint
	if input_size < 1<<15 ***REMOVED***
		for i = 0; i < input_size; i++ ***REMOVED***
			histogram[input[i]]++
		***REMOVED***

		histogram_total = input_size
		for i = 0; i < 256; i++ ***REMOVED***
			/* We weigh the first 11 samples with weight 3 to account for the
			   balancing effect of the LZ77 phase on the histogram. */
			var adjust uint32 = 2 * brotli_min_uint32_t(histogram[i], 11)
			histogram[i] += adjust
			histogram_total += uint(adjust)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		const kSampleRate uint = 29
		for i = 0; i < input_size; i += kSampleRate ***REMOVED***
			histogram[input[i]]++
		***REMOVED***

		histogram_total = (input_size + kSampleRate - 1) / kSampleRate
		for i = 0; i < 256; i++ ***REMOVED***
			/* We add 1 to each population count to avoid 0 bit depths (since this is
			   only a sample and we don't know if the symbol appears or not), and we
			   weigh the first 11 samples with weight 3 to account for the balancing
			   effect of the LZ77 phase on the histogram (more frequent symbols are
			   more likely to be in backward references instead as literals). */
			var adjust uint32 = 1 + 2*brotli_min_uint32_t(histogram[i], 11)
			histogram[i] += adjust
			histogram_total += uint(adjust)
		***REMOVED***
	***REMOVED***

	buildAndStoreHuffmanTreeFast(histogram[:], histogram_total, /* max_bits = */
		8, depths, bits, bw)
	***REMOVED***
		var literal_ratio uint = 0
		for i = 0; i < 256; i++ ***REMOVED***
			if histogram[i] != 0 ***REMOVED***
				literal_ratio += uint(histogram[i] * uint32(depths[i]))
			***REMOVED***
		***REMOVED***

		/* Estimated encoding ratio, millibytes per symbol. */
		return (literal_ratio * 125) / histogram_total
	***REMOVED***
***REMOVED***

/* Builds a command and distance prefix code (each 64 symbols) into "depth" and
   "bits" based on "histogram" and stores it into the bit stream. */
func buildAndStoreCommandPrefixCode1(histogram []uint32, depth []byte, bits []uint16, bw *bitWriter) ***REMOVED***
	var tree [129]huffmanTree
	var cmd_depth = [numCommandSymbols]byte***REMOVED***0***REMOVED***
	/* Tree size for building a tree over 64 symbols is 2 * 64 + 1. */

	var cmd_bits [64]uint16

	createHuffmanTree(histogram, 64, 15, tree[:], depth)
	createHuffmanTree(histogram[64:], 64, 14, tree[:], depth[64:])

	/* We have to jump through a few hoops here in order to compute
	   the command bits because the symbols are in a different order than in
	   the full alphabet. This looks complicated, but having the symbols
	   in this order in the command bits saves a few branches in the Emit*
	   functions. */
	copy(cmd_depth[:], depth[:24])

	copy(cmd_depth[24:][:], depth[40:][:8])
	copy(cmd_depth[32:][:], depth[24:][:8])
	copy(cmd_depth[40:][:], depth[48:][:8])
	copy(cmd_depth[48:][:], depth[32:][:8])
	copy(cmd_depth[56:][:], depth[56:][:8])
	convertBitDepthsToSymbols(cmd_depth[:], 64, cmd_bits[:])
	copy(bits, cmd_bits[:24])
	copy(bits[24:], cmd_bits[32:][:8])
	copy(bits[32:], cmd_bits[48:][:8])
	copy(bits[40:], cmd_bits[24:][:8])
	copy(bits[48:], cmd_bits[40:][:8])
	copy(bits[56:], cmd_bits[56:][:8])
	convertBitDepthsToSymbols(depth[64:], 64, bits[64:])
	***REMOVED***
		/* Create the bit length array for the full command alphabet. */
		var i uint
		for i := 0; i < int(64); i++ ***REMOVED***
			cmd_depth[i] = 0
		***REMOVED*** /* only 64 first values were used */
		copy(cmd_depth[:], depth[:8])
		copy(cmd_depth[64:][:], depth[8:][:8])
		copy(cmd_depth[128:][:], depth[16:][:8])
		copy(cmd_depth[192:][:], depth[24:][:8])
		copy(cmd_depth[384:][:], depth[32:][:8])
		for i = 0; i < 8; i++ ***REMOVED***
			cmd_depth[128+8*i] = depth[40+i]
			cmd_depth[256+8*i] = depth[48+i]
			cmd_depth[448+8*i] = depth[56+i]
		***REMOVED***

		storeHuffmanTree(cmd_depth[:], numCommandSymbols, tree[:], bw)
	***REMOVED***

	storeHuffmanTree(depth[64:], 64, tree[:], bw)
***REMOVED***

/* REQUIRES: insertlen < 6210 */
func emitInsertLen1(insertlen uint, depth []byte, bits []uint16, histo []uint32, bw *bitWriter) ***REMOVED***
	if insertlen < 6 ***REMOVED***
		var code uint = insertlen + 40
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		histo[code]++
	***REMOVED*** else if insertlen < 130 ***REMOVED***
		var tail uint = insertlen - 2
		var nbits uint32 = log2FloorNonZero(tail) - 1
		var prefix uint = tail >> nbits
		var inscode uint = uint((nbits << 1) + uint32(prefix) + 42)
		bw.writeBits(uint(depth[inscode]), uint64(bits[inscode]))
		bw.writeBits(uint(nbits), uint64(tail)-(uint64(prefix)<<nbits))
		histo[inscode]++
	***REMOVED*** else if insertlen < 2114 ***REMOVED***
		var tail uint = insertlen - 66
		var nbits uint32 = log2FloorNonZero(tail)
		var code uint = uint(nbits + 50)
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		bw.writeBits(uint(nbits), uint64(tail)-(uint64(uint(1))<<nbits))
		histo[code]++
	***REMOVED*** else ***REMOVED***
		bw.writeBits(uint(depth[61]), uint64(bits[61]))
		bw.writeBits(12, uint64(insertlen)-2114)
		histo[61]++
	***REMOVED***
***REMOVED***

func emitLongInsertLen(insertlen uint, depth []byte, bits []uint16, histo []uint32, bw *bitWriter) ***REMOVED***
	if insertlen < 22594 ***REMOVED***
		bw.writeBits(uint(depth[62]), uint64(bits[62]))
		bw.writeBits(14, uint64(insertlen)-6210)
		histo[62]++
	***REMOVED*** else ***REMOVED***
		bw.writeBits(uint(depth[63]), uint64(bits[63]))
		bw.writeBits(24, uint64(insertlen)-22594)
		histo[63]++
	***REMOVED***
***REMOVED***

func emitCopyLen1(copylen uint, depth []byte, bits []uint16, histo []uint32, bw *bitWriter) ***REMOVED***
	if copylen < 10 ***REMOVED***
		bw.writeBits(uint(depth[copylen+14]), uint64(bits[copylen+14]))
		histo[copylen+14]++
	***REMOVED*** else if copylen < 134 ***REMOVED***
		var tail uint = copylen - 6
		var nbits uint32 = log2FloorNonZero(tail) - 1
		var prefix uint = tail >> nbits
		var code uint = uint((nbits << 1) + uint32(prefix) + 20)
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		bw.writeBits(uint(nbits), uint64(tail)-(uint64(prefix)<<nbits))
		histo[code]++
	***REMOVED*** else if copylen < 2118 ***REMOVED***
		var tail uint = copylen - 70
		var nbits uint32 = log2FloorNonZero(tail)
		var code uint = uint(nbits + 28)
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		bw.writeBits(uint(nbits), uint64(tail)-(uint64(uint(1))<<nbits))
		histo[code]++
	***REMOVED*** else ***REMOVED***
		bw.writeBits(uint(depth[39]), uint64(bits[39]))
		bw.writeBits(24, uint64(copylen)-2118)
		histo[39]++
	***REMOVED***
***REMOVED***

func emitCopyLenLastDistance1(copylen uint, depth []byte, bits []uint16, histo []uint32, bw *bitWriter) ***REMOVED***
	if copylen < 12 ***REMOVED***
		bw.writeBits(uint(depth[copylen-4]), uint64(bits[copylen-4]))
		histo[copylen-4]++
	***REMOVED*** else if copylen < 72 ***REMOVED***
		var tail uint = copylen - 8
		var nbits uint32 = log2FloorNonZero(tail) - 1
		var prefix uint = tail >> nbits
		var code uint = uint((nbits << 1) + uint32(prefix) + 4)
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		bw.writeBits(uint(nbits), uint64(tail)-(uint64(prefix)<<nbits))
		histo[code]++
	***REMOVED*** else if copylen < 136 ***REMOVED***
		var tail uint = copylen - 8
		var code uint = (tail >> 5) + 30
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		bw.writeBits(5, uint64(tail)&31)
		bw.writeBits(uint(depth[64]), uint64(bits[64]))
		histo[code]++
		histo[64]++
	***REMOVED*** else if copylen < 2120 ***REMOVED***
		var tail uint = copylen - 72
		var nbits uint32 = log2FloorNonZero(tail)
		var code uint = uint(nbits + 28)
		bw.writeBits(uint(depth[code]), uint64(bits[code]))
		bw.writeBits(uint(nbits), uint64(tail)-(uint64(uint(1))<<nbits))
		bw.writeBits(uint(depth[64]), uint64(bits[64]))
		histo[code]++
		histo[64]++
	***REMOVED*** else ***REMOVED***
		bw.writeBits(uint(depth[39]), uint64(bits[39]))
		bw.writeBits(24, uint64(copylen)-2120)
		bw.writeBits(uint(depth[64]), uint64(bits[64]))
		histo[39]++
		histo[64]++
	***REMOVED***
***REMOVED***

func emitDistance1(distance uint, depth []byte, bits []uint16, histo []uint32, bw *bitWriter) ***REMOVED***
	var d uint = distance + 3
	var nbits uint32 = log2FloorNonZero(d) - 1
	var prefix uint = (d >> nbits) & 1
	var offset uint = (2 + prefix) << nbits
	var distcode uint = uint(2*(nbits-1) + uint32(prefix) + 80)
	bw.writeBits(uint(depth[distcode]), uint64(bits[distcode]))
	bw.writeBits(uint(nbits), uint64(d)-uint64(offset))
	histo[distcode]++
***REMOVED***

func emitLiterals(input []byte, len uint, depth []byte, bits []uint16, bw *bitWriter) ***REMOVED***
	var j uint
	for j = 0; j < len; j++ ***REMOVED***
		var lit byte = input[j]
		bw.writeBits(uint(depth[lit]), uint64(bits[lit]))
	***REMOVED***
***REMOVED***

/* REQUIRES: len <= 1 << 24. */
func storeMetaBlockHeader1(len uint, is_uncompressed bool, bw *bitWriter) ***REMOVED***
	var nibbles uint = 6

	/* ISLAST */
	bw.writeBits(1, 0)

	if len <= 1<<16 ***REMOVED***
		nibbles = 4
	***REMOVED*** else if len <= 1<<20 ***REMOVED***
		nibbles = 5
	***REMOVED***

	bw.writeBits(2, uint64(nibbles)-4)
	bw.writeBits(nibbles*4, uint64(len)-1)

	/* ISUNCOMPRESSED */
	bw.writeSingleBit(is_uncompressed)
***REMOVED***

var shouldMergeBlock_kSampleRate uint = 43

func shouldMergeBlock(data []byte, len uint, depths []byte) bool ***REMOVED***
	var histo = [256]uint***REMOVED***0***REMOVED***
	var i uint
	for i = 0; i < len; i += shouldMergeBlock_kSampleRate ***REMOVED***
		histo[data[i]]++
	***REMOVED***
	***REMOVED***
		var total uint = (len + shouldMergeBlock_kSampleRate - 1) / shouldMergeBlock_kSampleRate
		var r float64 = (fastLog2(total)+0.5)*float64(total) + 200
		for i = 0; i < 256; i++ ***REMOVED***
			r -= float64(histo[i]) * (float64(depths[i]) + fastLog2(histo[i]))
		***REMOVED***

		return r >= 0.0
	***REMOVED***
***REMOVED***

func shouldUseUncompressedMode(metablock_start []byte, next_emit []byte, insertlen uint, literal_ratio uint) bool ***REMOVED***
	var compressed uint = uint(-cap(next_emit) + cap(metablock_start))
	if compressed*50 > insertlen ***REMOVED***
		return false
	***REMOVED*** else ***REMOVED***
		return literal_ratio > 980
	***REMOVED***
***REMOVED***

func emitUncompressedMetaBlock1(data []byte, storage_ix_start uint, bw *bitWriter) ***REMOVED***
	bw.rewind(storage_ix_start)
	storeMetaBlockHeader1(uint(len(data)), true, bw)
	bw.jumpToByteBoundary()
	bw.writeBytes(data)
***REMOVED***

var kCmdHistoSeed = [128]uint32***REMOVED***
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 0, 0, 0, 0,
***REMOVED***

var compressFragmentFastImpl_kFirstBlockSize uint = 3 << 15
var compressFragmentFastImpl_kMergeBlockSize uint = 1 << 16

func compressFragmentFastImpl(in []byte, input_size uint, is_last bool, table []int, table_bits uint, cmd_depth []byte, cmd_bits []uint16, cmd_code_numbits *uint, cmd_code []byte, bw *bitWriter) ***REMOVED***
	var cmd_histo [128]uint32
	var ip_end int
	var next_emit int = 0
	var base_ip int = 0
	var input int = 0
	const kInputMarginBytes uint = windowGap
	const kMinMatchLen uint = 5
	var metablock_start int = input
	var block_size uint = brotli_min_size_t(input_size, compressFragmentFastImpl_kFirstBlockSize)
	var total_block_size uint = block_size
	var mlen_storage_ix uint = bw.getPos() + 3
	var lit_depth [256]byte
	var lit_bits [256]uint16
	var literal_ratio uint
	var ip int
	var last_distance int
	var shift uint = 64 - table_bits

	/* "next_emit" is a pointer to the first byte that is not covered by a
	   previous copy. Bytes between "next_emit" and the start of the next copy or
	   the end of the input will be emitted as literal bytes. */

	/* Save the start of the first block for position and distance computations.
	 */

	/* Save the bit position of the MLEN field of the meta-block header, so that
	   we can update it later if we decide to extend this meta-block. */
	storeMetaBlockHeader1(block_size, false, bw)

	/* No block splits, no contexts. */
	bw.writeBits(13, 0)

	literal_ratio = buildAndStoreLiteralPrefixCode(in[input:], block_size, lit_depth[:], lit_bits[:], bw)
	***REMOVED***
		/* Store the pre-compressed command and distance prefix codes. */
		var i uint
		for i = 0; i+7 < *cmd_code_numbits; i += 8 ***REMOVED***
			bw.writeBits(8, uint64(cmd_code[i>>3]))
		***REMOVED***
	***REMOVED***

	bw.writeBits(*cmd_code_numbits&7, uint64(cmd_code[*cmd_code_numbits>>3]))

	/* Initialize the command and distance histograms. We will gather
	   statistics of command and distance codes during the processing
	   of this block and use it to update the command and distance
	   prefix codes for the next block. */
emit_commands:
	copy(cmd_histo[:], kCmdHistoSeed[:])

	/* "ip" is the input pointer. */
	ip = input

	last_distance = -1
	ip_end = int(uint(input) + block_size)

	if block_size >= kInputMarginBytes ***REMOVED***
		var len_limit uint = brotli_min_size_t(block_size-kMinMatchLen, input_size-kInputMarginBytes)
		var ip_limit int = int(uint(input) + len_limit)
		/* For the last block, we need to keep a 16 bytes margin so that we can be
		   sure that all distances are at most window size - 16.
		   For all other blocks, we only need to keep a margin of 5 bytes so that
		   we don't go over the block size with a copy. */

		var next_hash uint32
		ip++
		for next_hash = hash5(in[ip:], shift); ; ***REMOVED***
			var skip uint32 = 32
			var next_ip int = ip
			/* Step 1: Scan forward in the input looking for a 5-byte-long match.
			   If we get close to exhausting the input then goto emit_remainder.

			   Heuristic match skipping: If 32 bytes are scanned with no matches
			   found, start looking only at every other byte. If 32 more bytes are
			   scanned, look at every third byte, etc.. When a match is found,
			   immediately go back to looking at every byte. This is a small loss
			   (~5% performance, ~0.1% density) for compressible data due to more
			   bookkeeping, but for non-compressible data (such as JPEG) it's a huge
			   win since the compressor quickly "realizes" the data is incompressible
			   and doesn't bother looking for matches everywhere.

			   The "skip" variable keeps track of how many bytes there are since the
			   last match; dividing it by 32 (i.e. right-shifting by five) gives the
			   number of bytes to move ahead for each iteration. */

			var candidate int
			assert(next_emit < ip)

		trawl:
			for ***REMOVED***
				var hash uint32 = next_hash
				var bytes_between_hash_lookups uint32 = skip >> 5
				skip++
				assert(hash == hash5(in[next_ip:], shift))
				ip = next_ip
				next_ip = int(uint32(ip) + bytes_between_hash_lookups)
				if next_ip > ip_limit ***REMOVED***
					goto emit_remainder
				***REMOVED***

				next_hash = hash5(in[next_ip:], shift)
				candidate = ip - last_distance
				if isMatch5(in[ip:], in[candidate:]) ***REMOVED***
					if candidate < ip ***REMOVED***
						table[hash] = int(ip - base_ip)
						break
					***REMOVED***
				***REMOVED***

				candidate = base_ip + table[hash]
				assert(candidate >= base_ip)
				assert(candidate < ip)

				table[hash] = int(ip - base_ip)
				if !(!isMatch5(in[ip:], in[candidate:])) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			/* Check copy distance. If candidate is not feasible, continue search.
			   Checking is done outside of hot loop to reduce overhead. */
			if ip-candidate > maxDistance_compress_fragment ***REMOVED***
				goto trawl
			***REMOVED***

			/* Step 2: Emit the found match together with the literal bytes from
			   "next_emit" to the bit stream, and then see if we can find a next match
			   immediately afterwards. Repeat until we find no match for the input
			   without emitting some literal bytes. */
			***REMOVED***
				var base int = ip
				/* > 0 */
				var matched uint = 5 + findMatchLengthWithLimit(in[candidate+5:], in[ip+5:], uint(ip_end-ip)-5)
				var distance int = int(base - candidate)
				/* We have a 5-byte match at ip, and we need to emit bytes in
				   [next_emit, ip). */

				var insert uint = uint(base - next_emit)
				ip += int(matched)
				if insert < 6210 ***REMOVED***
					emitInsertLen1(insert, cmd_depth, cmd_bits, cmd_histo[:], bw)
				***REMOVED*** else if shouldUseUncompressedMode(in[metablock_start:], in[next_emit:], insert, literal_ratio) ***REMOVED***
					emitUncompressedMetaBlock1(in[metablock_start:base], mlen_storage_ix-3, bw)
					input_size -= uint(base - input)
					input = base
					next_emit = input
					goto next_block
				***REMOVED*** else ***REMOVED***
					emitLongInsertLen(insert, cmd_depth, cmd_bits, cmd_histo[:], bw)
				***REMOVED***

				emitLiterals(in[next_emit:], insert, lit_depth[:], lit_bits[:], bw)
				if distance == last_distance ***REMOVED***
					bw.writeBits(uint(cmd_depth[64]), uint64(cmd_bits[64]))
					cmd_histo[64]++
				***REMOVED*** else ***REMOVED***
					emitDistance1(uint(distance), cmd_depth, cmd_bits, cmd_histo[:], bw)
					last_distance = distance
				***REMOVED***

				emitCopyLenLastDistance1(matched, cmd_depth, cmd_bits, cmd_histo[:], bw)

				next_emit = ip
				if ip >= ip_limit ***REMOVED***
					goto emit_remainder
				***REMOVED***

				/* We could immediately start working at ip now, but to improve
				   compression we first update "table" with the hashes of some positions
				   within the last copy. */
				***REMOVED***
					var input_bytes uint64 = binary.LittleEndian.Uint64(in[ip-3:])
					var prev_hash uint32 = hashBytesAtOffset5(input_bytes, 0, shift)
					var cur_hash uint32 = hashBytesAtOffset5(input_bytes, 3, shift)
					table[prev_hash] = int(ip - base_ip - 3)
					prev_hash = hashBytesAtOffset5(input_bytes, 1, shift)
					table[prev_hash] = int(ip - base_ip - 2)
					prev_hash = hashBytesAtOffset5(input_bytes, 2, shift)
					table[prev_hash] = int(ip - base_ip - 1)

					candidate = base_ip + table[cur_hash]
					table[cur_hash] = int(ip - base_ip)
				***REMOVED***
			***REMOVED***

			for isMatch5(in[ip:], in[candidate:]) ***REMOVED***
				var base int = ip
				/* We have a 5-byte match at ip, and no need to emit any literal bytes
				   prior to ip. */

				var matched uint = 5 + findMatchLengthWithLimit(in[candidate+5:], in[ip+5:], uint(ip_end-ip)-5)
				if ip-candidate > maxDistance_compress_fragment ***REMOVED***
					break
				***REMOVED***
				ip += int(matched)
				last_distance = int(base - candidate) /* > 0 */
				emitCopyLen1(matched, cmd_depth, cmd_bits, cmd_histo[:], bw)
				emitDistance1(uint(last_distance), cmd_depth, cmd_bits, cmd_histo[:], bw)

				next_emit = ip
				if ip >= ip_limit ***REMOVED***
					goto emit_remainder
				***REMOVED***

				/* We could immediately start working at ip now, but to improve
				   compression we first update "table" with the hashes of some positions
				   within the last copy. */
				***REMOVED***
					var input_bytes uint64 = binary.LittleEndian.Uint64(in[ip-3:])
					var prev_hash uint32 = hashBytesAtOffset5(input_bytes, 0, shift)
					var cur_hash uint32 = hashBytesAtOffset5(input_bytes, 3, shift)
					table[prev_hash] = int(ip - base_ip - 3)
					prev_hash = hashBytesAtOffset5(input_bytes, 1, shift)
					table[prev_hash] = int(ip - base_ip - 2)
					prev_hash = hashBytesAtOffset5(input_bytes, 2, shift)
					table[prev_hash] = int(ip - base_ip - 1)

					candidate = base_ip + table[cur_hash]
					table[cur_hash] = int(ip - base_ip)
				***REMOVED***
			***REMOVED***

			ip++
			next_hash = hash5(in[ip:], shift)
		***REMOVED***
	***REMOVED***

emit_remainder:
	assert(next_emit <= ip_end)
	input += int(block_size)
	input_size -= block_size
	block_size = brotli_min_size_t(input_size, compressFragmentFastImpl_kMergeBlockSize)

	/* Decide if we want to continue this meta-block instead of emitting the
	   last insert-only command. */
	if input_size > 0 && total_block_size+block_size <= 1<<20 && shouldMergeBlock(in[input:], block_size, lit_depth[:]) ***REMOVED***
		assert(total_block_size > 1<<16)

		/* Update the size of the current meta-block and continue emitting commands.
		   We can do this because the current size and the new size both have 5
		   nibbles. */
		total_block_size += block_size

		bw.updateBits(20, uint32(total_block_size-1), mlen_storage_ix)
		goto emit_commands
	***REMOVED***

	/* Emit the remaining bytes as literals. */
	if next_emit < ip_end ***REMOVED***
		var insert uint = uint(ip_end - next_emit)
		if insert < 6210 ***REMOVED***
			emitInsertLen1(insert, cmd_depth, cmd_bits, cmd_histo[:], bw)
			emitLiterals(in[next_emit:], insert, lit_depth[:], lit_bits[:], bw)
		***REMOVED*** else if shouldUseUncompressedMode(in[metablock_start:], in[next_emit:], insert, literal_ratio) ***REMOVED***
			emitUncompressedMetaBlock1(in[metablock_start:ip_end], mlen_storage_ix-3, bw)
		***REMOVED*** else ***REMOVED***
			emitLongInsertLen(insert, cmd_depth, cmd_bits, cmd_histo[:], bw)
			emitLiterals(in[next_emit:], insert, lit_depth[:], lit_bits[:], bw)
		***REMOVED***
	***REMOVED***

	next_emit = ip_end

	/* If we have more data, write a new meta-block header and prefix codes and
	   then continue emitting commands. */
next_block:
	if input_size > 0 ***REMOVED***
		metablock_start = input
		block_size = brotli_min_size_t(input_size, compressFragmentFastImpl_kFirstBlockSize)
		total_block_size = block_size

		/* Save the bit position of the MLEN field of the meta-block header, so that
		   we can update it later if we decide to extend this meta-block. */
		mlen_storage_ix = bw.getPos() + 3

		storeMetaBlockHeader1(block_size, false, bw)

		/* No block splits, no contexts. */
		bw.writeBits(13, 0)

		literal_ratio = buildAndStoreLiteralPrefixCode(in[input:], block_size, lit_depth[:], lit_bits[:], bw)
		buildAndStoreCommandPrefixCode1(cmd_histo[:], cmd_depth, cmd_bits, bw)
		goto emit_commands
	***REMOVED***

	if !is_last ***REMOVED***
		/* If this is not the last block, update the command and distance prefix
		   codes for the next block and store the compressed forms. */
		var bw bitWriter
		bw.dst = cmd_code
		buildAndStoreCommandPrefixCode1(cmd_histo[:], cmd_depth, cmd_bits, &bw)
		*cmd_code_numbits = bw.getPos()
	***REMOVED***
***REMOVED***

/* Compresses "input" string to bw as one or more complete meta-blocks.

   If "is_last" is 1, emits an additional empty last meta-block.

   "cmd_depth" and "cmd_bits" contain the command and distance prefix codes
   (see comment in encode.h) used for the encoding of this input fragment.
   If "is_last" is 0, they are updated to reflect the statistics
   of this input fragment, to be used for the encoding of the next fragment.

   "*cmd_code_numbits" is the number of bits of the compressed representation
   of the command and distance prefix codes, and "cmd_code" is an array of
   at least "(*cmd_code_numbits + 7) >> 3" size that contains the compressed
   command and distance prefix codes. If "is_last" is 0, these are also
   updated to represent the updated "cmd_depth" and "cmd_bits".

   REQUIRES: "input_size" is greater than zero, or "is_last" is 1.
   REQUIRES: "input_size" is less or equal to maximal metablock size (1 << 24).
   REQUIRES: All elements in "table[0..table_size-1]" are initialized to zero.
   REQUIRES: "table_size" is an odd (9, 11, 13, 15) power of two
   OUTPUT: maximal copy distance <= |input_size|
   OUTPUT: maximal copy distance <= BROTLI_MAX_BACKWARD_LIMIT(18) */
func compressFragmentFast(input []byte, input_size uint, is_last bool, table []int, table_size uint, cmd_depth []byte, cmd_bits []uint16, cmd_code_numbits *uint, cmd_code []byte, bw *bitWriter) ***REMOVED***
	var initial_storage_ix uint = bw.getPos()
	var table_bits uint = uint(log2FloorNonZero(table_size))

	if input_size == 0 ***REMOVED***
		assert(is_last)
		bw.writeBits(1, 1) /* islast */
		bw.writeBits(1, 1) /* isempty */
		bw.jumpToByteBoundary()
		return
	***REMOVED***

	compressFragmentFastImpl(input, input_size, is_last, table, table_bits, cmd_depth, cmd_bits, cmd_code_numbits, cmd_code, bw)

	/* If output is larger than single uncompressed block, rewrite it. */
	if bw.getPos()-initial_storage_ix > 31+(input_size<<3) ***REMOVED***
		emitUncompressedMetaBlock1(input[:input_size], initial_storage_ix, bw)
	***REMOVED***

	if is_last ***REMOVED***
		bw.writeBits(1, 1) /* islast */
		bw.writeBits(1, 1) /* isempty */
		bw.jumpToByteBoundary()
	***REMOVED***
***REMOVED***
