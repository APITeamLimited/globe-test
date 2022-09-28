package brotli

import "math"

/* The distance symbols effectively used by "Large Window Brotli" (32-bit). */
const numHistogramDistanceSymbols = 544

type histogramLiteral struct ***REMOVED***
	data_        [numLiteralSymbols]uint32
	total_count_ uint
	bit_cost_    float64
***REMOVED***

func histogramClearLiteral(self *histogramLiteral) ***REMOVED***
	self.data_ = [numLiteralSymbols]uint32***REMOVED******REMOVED***
	self.total_count_ = 0
	self.bit_cost_ = math.MaxFloat64
***REMOVED***

func clearHistogramsLiteral(array []histogramLiteral, length uint) ***REMOVED***
	var i uint
	for i = 0; i < length; i++ ***REMOVED***
		histogramClearLiteral(&array[i:][0])
	***REMOVED***
***REMOVED***

func histogramAddLiteral(self *histogramLiteral, val uint) ***REMOVED***
	self.data_[val]++
	self.total_count_++
***REMOVED***

func histogramAddVectorLiteral(self *histogramLiteral, p []byte, n uint) ***REMOVED***
	self.total_count_ += n
	n += 1
	for ***REMOVED***
		n--
		if n == 0 ***REMOVED***
			break
		***REMOVED***
		self.data_[p[0]]++
		p = p[1:]
	***REMOVED***
***REMOVED***

func histogramAddHistogramLiteral(self *histogramLiteral, v *histogramLiteral) ***REMOVED***
	var i uint
	self.total_count_ += v.total_count_
	for i = 0; i < numLiteralSymbols; i++ ***REMOVED***
		self.data_[i] += v.data_[i]
	***REMOVED***
***REMOVED***

func histogramDataSizeLiteral() uint ***REMOVED***
	return numLiteralSymbols
***REMOVED***

type histogramCommand struct ***REMOVED***
	data_        [numCommandSymbols]uint32
	total_count_ uint
	bit_cost_    float64
***REMOVED***

func histogramClearCommand(self *histogramCommand) ***REMOVED***
	self.data_ = [numCommandSymbols]uint32***REMOVED******REMOVED***
	self.total_count_ = 0
	self.bit_cost_ = math.MaxFloat64
***REMOVED***

func clearHistogramsCommand(array []histogramCommand, length uint) ***REMOVED***
	var i uint
	for i = 0; i < length; i++ ***REMOVED***
		histogramClearCommand(&array[i:][0])
	***REMOVED***
***REMOVED***

func histogramAddCommand(self *histogramCommand, val uint) ***REMOVED***
	self.data_[val]++
	self.total_count_++
***REMOVED***

func histogramAddVectorCommand(self *histogramCommand, p []uint16, n uint) ***REMOVED***
	self.total_count_ += n
	n += 1
	for ***REMOVED***
		n--
		if n == 0 ***REMOVED***
			break
		***REMOVED***
		self.data_[p[0]]++
		p = p[1:]
	***REMOVED***
***REMOVED***

func histogramAddHistogramCommand(self *histogramCommand, v *histogramCommand) ***REMOVED***
	var i uint
	self.total_count_ += v.total_count_
	for i = 0; i < numCommandSymbols; i++ ***REMOVED***
		self.data_[i] += v.data_[i]
	***REMOVED***
***REMOVED***

func histogramDataSizeCommand() uint ***REMOVED***
	return numCommandSymbols
***REMOVED***

type histogramDistance struct ***REMOVED***
	data_        [numDistanceSymbols]uint32
	total_count_ uint
	bit_cost_    float64
***REMOVED***

func histogramClearDistance(self *histogramDistance) ***REMOVED***
	self.data_ = [numDistanceSymbols]uint32***REMOVED******REMOVED***
	self.total_count_ = 0
	self.bit_cost_ = math.MaxFloat64
***REMOVED***

func clearHistogramsDistance(array []histogramDistance, length uint) ***REMOVED***
	var i uint
	for i = 0; i < length; i++ ***REMOVED***
		histogramClearDistance(&array[i:][0])
	***REMOVED***
***REMOVED***

func histogramAddDistance(self *histogramDistance, val uint) ***REMOVED***
	self.data_[val]++
	self.total_count_++
***REMOVED***

func histogramAddVectorDistance(self *histogramDistance, p []uint16, n uint) ***REMOVED***
	self.total_count_ += n
	n += 1
	for ***REMOVED***
		n--
		if n == 0 ***REMOVED***
			break
		***REMOVED***
		self.data_[p[0]]++
		p = p[1:]
	***REMOVED***
***REMOVED***

func histogramAddHistogramDistance(self *histogramDistance, v *histogramDistance) ***REMOVED***
	var i uint
	self.total_count_ += v.total_count_
	for i = 0; i < numDistanceSymbols; i++ ***REMOVED***
		self.data_[i] += v.data_[i]
	***REMOVED***
***REMOVED***

func histogramDataSizeDistance() uint ***REMOVED***
	return numDistanceSymbols
***REMOVED***

type blockSplitIterator struct ***REMOVED***
	split_  *blockSplit
	idx_    uint
	type_   uint
	length_ uint
***REMOVED***

func initBlockSplitIterator(self *blockSplitIterator, split *blockSplit) ***REMOVED***
	self.split_ = split
	self.idx_ = 0
	self.type_ = 0
	if len(split.lengths) > 0 ***REMOVED***
		self.length_ = uint(split.lengths[0])
	***REMOVED*** else ***REMOVED***
		self.length_ = 0
	***REMOVED***
***REMOVED***

func blockSplitIteratorNext(self *blockSplitIterator) ***REMOVED***
	if self.length_ == 0 ***REMOVED***
		self.idx_++
		self.type_ = uint(self.split_.types[self.idx_])
		self.length_ = uint(self.split_.lengths[self.idx_])
	***REMOVED***

	self.length_--
***REMOVED***

func buildHistogramsWithContext(cmds []command, literal_split *blockSplit, insert_and_copy_split *blockSplit, dist_split *blockSplit, ringbuffer []byte, start_pos uint, mask uint, prev_byte byte, prev_byte2 byte, context_modes []int, literal_histograms []histogramLiteral, insert_and_copy_histograms []histogramCommand, copy_dist_histograms []histogramDistance) ***REMOVED***
	var pos uint = start_pos
	var literal_it blockSplitIterator
	var insert_and_copy_it blockSplitIterator
	var dist_it blockSplitIterator

	initBlockSplitIterator(&literal_it, literal_split)
	initBlockSplitIterator(&insert_and_copy_it, insert_and_copy_split)
	initBlockSplitIterator(&dist_it, dist_split)
	for i := range cmds ***REMOVED***
		var cmd *command = &cmds[i]
		var j uint
		blockSplitIteratorNext(&insert_and_copy_it)
		histogramAddCommand(&insert_and_copy_histograms[insert_and_copy_it.type_], uint(cmd.cmd_prefix_))

		/* TODO: unwrap iterator blocks. */
		for j = uint(cmd.insert_len_); j != 0; j-- ***REMOVED***
			var context uint
			blockSplitIteratorNext(&literal_it)
			context = literal_it.type_
			if context_modes != nil ***REMOVED***
				var lut contextLUT = getContextLUT(context_modes[context])
				context = (context << literalContextBits) + uint(getContext(prev_byte, prev_byte2, lut))
			***REMOVED***

			histogramAddLiteral(&literal_histograms[context], uint(ringbuffer[pos&mask]))
			prev_byte2 = prev_byte
			prev_byte = ringbuffer[pos&mask]
			pos++
		***REMOVED***

		pos += uint(commandCopyLen(cmd))
		if commandCopyLen(cmd) != 0 ***REMOVED***
			prev_byte2 = ringbuffer[(pos-2)&mask]
			prev_byte = ringbuffer[(pos-1)&mask]
			if cmd.cmd_prefix_ >= 128 ***REMOVED***
				var context uint
				blockSplitIteratorNext(&dist_it)
				context = uint(uint32(dist_it.type_<<distanceContextBits) + commandDistanceContext(cmd))
				histogramAddDistance(&copy_dist_histograms[context], uint(cmd.dist_prefix_)&0x3FF)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
