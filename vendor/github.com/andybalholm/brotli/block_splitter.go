package brotli

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Block split point selection utilities. */

type blockSplit struct ***REMOVED***
	num_types          uint
	num_blocks         uint
	types              []byte
	lengths            []uint32
	types_alloc_size   uint
	lengths_alloc_size uint
***REMOVED***

const (
	kMaxLiteralHistograms        uint    = 100
	kMaxCommandHistograms        uint    = 50
	kLiteralBlockSwitchCost      float64 = 28.1
	kCommandBlockSwitchCost      float64 = 13.5
	kDistanceBlockSwitchCost     float64 = 14.6
	kLiteralStrideLength         uint    = 70
	kCommandStrideLength         uint    = 40
	kSymbolsPerLiteralHistogram  uint    = 544
	kSymbolsPerCommandHistogram  uint    = 530
	kSymbolsPerDistanceHistogram uint    = 544
	kMinLengthForBlockSplitting  uint    = 128
	kIterMulForRefining          uint    = 2
	kMinItersForRefining         uint    = 100
)

func countLiterals(cmds []command, num_commands uint) uint ***REMOVED***
	var total_length uint = 0
	/* Count how many we have. */

	var i uint
	for i = 0; i < num_commands; i++ ***REMOVED***
		total_length += uint(cmds[i].insert_len_)
	***REMOVED***

	return total_length
***REMOVED***

func copyLiteralsToByteArray(cmds []command, num_commands uint, data []byte, offset uint, mask uint, literals []byte) ***REMOVED***
	var pos uint = 0
	var from_pos uint = offset & mask
	var i uint
	for i = 0; i < num_commands; i++ ***REMOVED***
		var insert_len uint = uint(cmds[i].insert_len_)
		if from_pos+insert_len > mask ***REMOVED***
			var head_size uint = mask + 1 - from_pos
			copy(literals[pos:], data[from_pos:][:head_size])
			from_pos = 0
			pos += head_size
			insert_len -= head_size
		***REMOVED***

		if insert_len > 0 ***REMOVED***
			copy(literals[pos:], data[from_pos:][:insert_len])
			pos += insert_len
		***REMOVED***

		from_pos = uint((uint32(from_pos+insert_len) + commandCopyLen(&cmds[i])) & uint32(mask))
	***REMOVED***
***REMOVED***

func myRand(seed *uint32) uint32 ***REMOVED***
	/* Initial seed should be 7. In this case, loop length is (1 << 29). */
	*seed *= 16807

	return *seed
***REMOVED***

func bitCost(count uint) float64 ***REMOVED***
	if count == 0 ***REMOVED***
		return -2.0
	***REMOVED*** else ***REMOVED***
		return fastLog2(count)
	***REMOVED***
***REMOVED***

const histogramsPerBatch = 64

const clustersPerBatch = 16

func initBlockSplit(self *blockSplit) ***REMOVED***
	self.num_types = 0
	self.num_blocks = 0
	self.types = nil
	self.lengths = nil
	self.types_alloc_size = 0
	self.lengths_alloc_size = 0
***REMOVED***

func destroyBlockSplit(self *blockSplit) ***REMOVED***
	self.types = nil
	self.lengths = nil
***REMOVED***

func splitBlock(cmds []command, num_commands uint, data []byte, pos uint, mask uint, params *encoderParams, literal_split *blockSplit, insert_and_copy_split *blockSplit, dist_split *blockSplit) ***REMOVED***
	***REMOVED***
		var literals_count uint = countLiterals(cmds, num_commands)
		var literals []byte = make([]byte, literals_count)

		/* Create a continuous array of literals. */
		copyLiteralsToByteArray(cmds, num_commands, data, pos, mask, literals)

		/* Create the block split on the array of literals.
		   Literal histograms have alphabet size 256. */
		splitByteVectorLiteral(literals, literals_count, kSymbolsPerLiteralHistogram, kMaxLiteralHistograms, kLiteralStrideLength, kLiteralBlockSwitchCost, params, literal_split)

		literals = nil
	***REMOVED***
	***REMOVED***
		var insert_and_copy_codes []uint16 = make([]uint16, num_commands)
		/* Compute prefix codes for commands. */

		var i uint
		for i = 0; i < num_commands; i++ ***REMOVED***
			insert_and_copy_codes[i] = cmds[i].cmd_prefix_
		***REMOVED***

		/* Create the block split on the array of command prefixes. */
		splitByteVectorCommand(insert_and_copy_codes, num_commands, kSymbolsPerCommandHistogram, kMaxCommandHistograms, kCommandStrideLength, kCommandBlockSwitchCost, params, insert_and_copy_split)

		/* TODO: reuse for distances? */

		insert_and_copy_codes = nil
	***REMOVED***
	***REMOVED***
		var distance_prefixes []uint16 = make([]uint16, num_commands)
		var j uint = 0
		/* Create a continuous array of distance prefixes. */

		var i uint
		for i = 0; i < num_commands; i++ ***REMOVED***
			var cmd *command = &cmds[i]
			if commandCopyLen(cmd) != 0 && cmd.cmd_prefix_ >= 128 ***REMOVED***
				distance_prefixes[j] = cmd.dist_prefix_ & 0x3FF
				j++
			***REMOVED***
		***REMOVED***

		/* Create the block split on the array of distance prefixes. */
		splitByteVectorDistance(distance_prefixes, j, kSymbolsPerDistanceHistogram, kMaxCommandHistograms, kCommandStrideLength, kDistanceBlockSwitchCost, params, dist_split)

		distance_prefixes = nil
	***REMOVED***
***REMOVED***
