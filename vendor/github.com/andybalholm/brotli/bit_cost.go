package brotli

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Functions to estimate the bit cost of Huffman trees. */
func shannonEntropy(population []uint32, size uint, total *uint) float64 ***REMOVED***
	var sum uint = 0
	var retval float64 = 0
	var population_end []uint32 = population[size:]
	var p uint
	for -cap(population) < -cap(population_end) ***REMOVED***
		p = uint(population[0])
		population = population[1:]
		sum += p
		retval -= float64(p) * fastLog2(p)
	***REMOVED***

	if sum != 0 ***REMOVED***
		retval += float64(sum) * fastLog2(sum)
	***REMOVED***
	*total = sum
	return retval
***REMOVED***

func bitsEntropy(population []uint32, size uint) float64 ***REMOVED***
	var sum uint
	var retval float64 = shannonEntropy(population, size, &sum)
	if retval < float64(sum) ***REMOVED***
		/* At least one bit per literal is needed. */
		retval = float64(sum)
	***REMOVED***

	return retval
***REMOVED***

const kOneSymbolHistogramCost float64 = 12
const kTwoSymbolHistogramCost float64 = 20
const kThreeSymbolHistogramCost float64 = 28
const kFourSymbolHistogramCost float64 = 37

func populationCostLiteral(histogram *histogramLiteral) float64 ***REMOVED***
	var data_size uint = histogramDataSizeLiteral()
	var count int = 0
	var s [5]uint
	var bits float64 = 0.0
	var i uint
	if histogram.total_count_ == 0 ***REMOVED***
		return kOneSymbolHistogramCost
	***REMOVED***

	for i = 0; i < data_size; i++ ***REMOVED***
		if histogram.data_[i] > 0 ***REMOVED***
			s[count] = i
			count++
			if count > 4 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if count == 1 ***REMOVED***
		return kOneSymbolHistogramCost
	***REMOVED***

	if count == 2 ***REMOVED***
		return kTwoSymbolHistogramCost + float64(histogram.total_count_)
	***REMOVED***

	if count == 3 ***REMOVED***
		var histo0 uint32 = histogram.data_[s[0]]
		var histo1 uint32 = histogram.data_[s[1]]
		var histo2 uint32 = histogram.data_[s[2]]
		var histomax uint32 = brotli_max_uint32_t(histo0, brotli_max_uint32_t(histo1, histo2))
		return kThreeSymbolHistogramCost + 2*(float64(histo0)+float64(histo1)+float64(histo2)) - float64(histomax)
	***REMOVED***

	if count == 4 ***REMOVED***
		var histo [4]uint32
		var h23 uint32
		var histomax uint32
		for i = 0; i < 4; i++ ***REMOVED***
			histo[i] = histogram.data_[s[i]]
		***REMOVED***

		/* Sort */
		for i = 0; i < 4; i++ ***REMOVED***
			var j uint
			for j = i + 1; j < 4; j++ ***REMOVED***
				if histo[j] > histo[i] ***REMOVED***
					var tmp uint32 = histo[j]
					histo[j] = histo[i]
					histo[i] = tmp
				***REMOVED***
			***REMOVED***
		***REMOVED***

		h23 = histo[2] + histo[3]
		histomax = brotli_max_uint32_t(h23, histo[0])
		return kFourSymbolHistogramCost + 3*float64(h23) + 2*(float64(histo[0])+float64(histo[1])) - float64(histomax)
	***REMOVED***
	***REMOVED***
		var max_depth uint = 1
		var depth_histo = [codeLengthCodes]uint32***REMOVED***0***REMOVED***
		/* In this loop we compute the entropy of the histogram and simultaneously
		   build a simplified histogram of the code length codes where we use the
		   zero repeat code 17, but we don't use the non-zero repeat code 16. */

		var log2total float64 = fastLog2(histogram.total_count_)
		for i = 0; i < data_size; ***REMOVED***
			if histogram.data_[i] > 0 ***REMOVED***
				var log2p float64 = log2total - fastLog2(uint(histogram.data_[i]))
				/* Compute -log2(P(symbol)) = -log2(count(symbol)/total_count) =
				   = log2(total_count) - log2(count(symbol)) */

				var depth uint = uint(log2p + 0.5)
				/* Approximate the bit depth by round(-log2(P(symbol))) */
				bits += float64(histogram.data_[i]) * log2p

				if depth > 15 ***REMOVED***
					depth = 15
				***REMOVED***

				if depth > max_depth ***REMOVED***
					max_depth = depth
				***REMOVED***

				depth_histo[depth]++
				i++
			***REMOVED*** else ***REMOVED***
				var reps uint32 = 1
				/* Compute the run length of zeros and add the appropriate number of 0
				   and 17 code length codes to the code length code histogram. */

				var k uint
				for k = i + 1; k < data_size && histogram.data_[k] == 0; k++ ***REMOVED***
					reps++
				***REMOVED***

				i += uint(reps)
				if i == data_size ***REMOVED***
					/* Don't add any cost for the last zero run, since these are encoded
					   only implicitly. */
					break
				***REMOVED***

				if reps < 3 ***REMOVED***
					depth_histo[0] += reps
				***REMOVED*** else ***REMOVED***
					reps -= 2
					for reps > 0 ***REMOVED***
						depth_histo[repeatZeroCodeLength]++

						/* Add the 3 extra bits for the 17 code length code. */
						bits += 3

						reps >>= 3
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		/* Add the estimated encoding cost of the code length code histogram. */
		bits += float64(18 + 2*max_depth)

		/* Add the entropy of the code length code histogram. */
		bits += bitsEntropy(depth_histo[:], codeLengthCodes)
	***REMOVED***

	return bits
***REMOVED***

func populationCostCommand(histogram *histogramCommand) float64 ***REMOVED***
	var data_size uint = histogramDataSizeCommand()
	var count int = 0
	var s [5]uint
	var bits float64 = 0.0
	var i uint
	if histogram.total_count_ == 0 ***REMOVED***
		return kOneSymbolHistogramCost
	***REMOVED***

	for i = 0; i < data_size; i++ ***REMOVED***
		if histogram.data_[i] > 0 ***REMOVED***
			s[count] = i
			count++
			if count > 4 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if count == 1 ***REMOVED***
		return kOneSymbolHistogramCost
	***REMOVED***

	if count == 2 ***REMOVED***
		return kTwoSymbolHistogramCost + float64(histogram.total_count_)
	***REMOVED***

	if count == 3 ***REMOVED***
		var histo0 uint32 = histogram.data_[s[0]]
		var histo1 uint32 = histogram.data_[s[1]]
		var histo2 uint32 = histogram.data_[s[2]]
		var histomax uint32 = brotli_max_uint32_t(histo0, brotli_max_uint32_t(histo1, histo2))
		return kThreeSymbolHistogramCost + 2*(float64(histo0)+float64(histo1)+float64(histo2)) - float64(histomax)
	***REMOVED***

	if count == 4 ***REMOVED***
		var histo [4]uint32
		var h23 uint32
		var histomax uint32
		for i = 0; i < 4; i++ ***REMOVED***
			histo[i] = histogram.data_[s[i]]
		***REMOVED***

		/* Sort */
		for i = 0; i < 4; i++ ***REMOVED***
			var j uint
			for j = i + 1; j < 4; j++ ***REMOVED***
				if histo[j] > histo[i] ***REMOVED***
					var tmp uint32 = histo[j]
					histo[j] = histo[i]
					histo[i] = tmp
				***REMOVED***
			***REMOVED***
		***REMOVED***

		h23 = histo[2] + histo[3]
		histomax = brotli_max_uint32_t(h23, histo[0])
		return kFourSymbolHistogramCost + 3*float64(h23) + 2*(float64(histo[0])+float64(histo[1])) - float64(histomax)
	***REMOVED***
	***REMOVED***
		var max_depth uint = 1
		var depth_histo = [codeLengthCodes]uint32***REMOVED***0***REMOVED***
		/* In this loop we compute the entropy of the histogram and simultaneously
		   build a simplified histogram of the code length codes where we use the
		   zero repeat code 17, but we don't use the non-zero repeat code 16. */

		var log2total float64 = fastLog2(histogram.total_count_)
		for i = 0; i < data_size; ***REMOVED***
			if histogram.data_[i] > 0 ***REMOVED***
				var log2p float64 = log2total - fastLog2(uint(histogram.data_[i]))
				/* Compute -log2(P(symbol)) = -log2(count(symbol)/total_count) =
				   = log2(total_count) - log2(count(symbol)) */

				var depth uint = uint(log2p + 0.5)
				/* Approximate the bit depth by round(-log2(P(symbol))) */
				bits += float64(histogram.data_[i]) * log2p

				if depth > 15 ***REMOVED***
					depth = 15
				***REMOVED***

				if depth > max_depth ***REMOVED***
					max_depth = depth
				***REMOVED***

				depth_histo[depth]++
				i++
			***REMOVED*** else ***REMOVED***
				var reps uint32 = 1
				/* Compute the run length of zeros and add the appropriate number of 0
				   and 17 code length codes to the code length code histogram. */

				var k uint
				for k = i + 1; k < data_size && histogram.data_[k] == 0; k++ ***REMOVED***
					reps++
				***REMOVED***

				i += uint(reps)
				if i == data_size ***REMOVED***
					/* Don't add any cost for the last zero run, since these are encoded
					   only implicitly. */
					break
				***REMOVED***

				if reps < 3 ***REMOVED***
					depth_histo[0] += reps
				***REMOVED*** else ***REMOVED***
					reps -= 2
					for reps > 0 ***REMOVED***
						depth_histo[repeatZeroCodeLength]++

						/* Add the 3 extra bits for the 17 code length code. */
						bits += 3

						reps >>= 3
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		/* Add the estimated encoding cost of the code length code histogram. */
		bits += float64(18 + 2*max_depth)

		/* Add the entropy of the code length code histogram. */
		bits += bitsEntropy(depth_histo[:], codeLengthCodes)
	***REMOVED***

	return bits
***REMOVED***

func populationCostDistance(histogram *histogramDistance) float64 ***REMOVED***
	var data_size uint = histogramDataSizeDistance()
	var count int = 0
	var s [5]uint
	var bits float64 = 0.0
	var i uint
	if histogram.total_count_ == 0 ***REMOVED***
		return kOneSymbolHistogramCost
	***REMOVED***

	for i = 0; i < data_size; i++ ***REMOVED***
		if histogram.data_[i] > 0 ***REMOVED***
			s[count] = i
			count++
			if count > 4 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if count == 1 ***REMOVED***
		return kOneSymbolHistogramCost
	***REMOVED***

	if count == 2 ***REMOVED***
		return kTwoSymbolHistogramCost + float64(histogram.total_count_)
	***REMOVED***

	if count == 3 ***REMOVED***
		var histo0 uint32 = histogram.data_[s[0]]
		var histo1 uint32 = histogram.data_[s[1]]
		var histo2 uint32 = histogram.data_[s[2]]
		var histomax uint32 = brotli_max_uint32_t(histo0, brotli_max_uint32_t(histo1, histo2))
		return kThreeSymbolHistogramCost + 2*(float64(histo0)+float64(histo1)+float64(histo2)) - float64(histomax)
	***REMOVED***

	if count == 4 ***REMOVED***
		var histo [4]uint32
		var h23 uint32
		var histomax uint32
		for i = 0; i < 4; i++ ***REMOVED***
			histo[i] = histogram.data_[s[i]]
		***REMOVED***

		/* Sort */
		for i = 0; i < 4; i++ ***REMOVED***
			var j uint
			for j = i + 1; j < 4; j++ ***REMOVED***
				if histo[j] > histo[i] ***REMOVED***
					var tmp uint32 = histo[j]
					histo[j] = histo[i]
					histo[i] = tmp
				***REMOVED***
			***REMOVED***
		***REMOVED***

		h23 = histo[2] + histo[3]
		histomax = brotli_max_uint32_t(h23, histo[0])
		return kFourSymbolHistogramCost + 3*float64(h23) + 2*(float64(histo[0])+float64(histo[1])) - float64(histomax)
	***REMOVED***
	***REMOVED***
		var max_depth uint = 1
		var depth_histo = [codeLengthCodes]uint32***REMOVED***0***REMOVED***
		/* In this loop we compute the entropy of the histogram and simultaneously
		   build a simplified histogram of the code length codes where we use the
		   zero repeat code 17, but we don't use the non-zero repeat code 16. */

		var log2total float64 = fastLog2(histogram.total_count_)
		for i = 0; i < data_size; ***REMOVED***
			if histogram.data_[i] > 0 ***REMOVED***
				var log2p float64 = log2total - fastLog2(uint(histogram.data_[i]))
				/* Compute -log2(P(symbol)) = -log2(count(symbol)/total_count) =
				   = log2(total_count) - log2(count(symbol)) */

				var depth uint = uint(log2p + 0.5)
				/* Approximate the bit depth by round(-log2(P(symbol))) */
				bits += float64(histogram.data_[i]) * log2p

				if depth > 15 ***REMOVED***
					depth = 15
				***REMOVED***

				if depth > max_depth ***REMOVED***
					max_depth = depth
				***REMOVED***

				depth_histo[depth]++
				i++
			***REMOVED*** else ***REMOVED***
				var reps uint32 = 1
				/* Compute the run length of zeros and add the appropriate number of 0
				   and 17 code length codes to the code length code histogram. */

				var k uint
				for k = i + 1; k < data_size && histogram.data_[k] == 0; k++ ***REMOVED***
					reps++
				***REMOVED***

				i += uint(reps)
				if i == data_size ***REMOVED***
					/* Don't add any cost for the last zero run, since these are encoded
					   only implicitly. */
					break
				***REMOVED***

				if reps < 3 ***REMOVED***
					depth_histo[0] += reps
				***REMOVED*** else ***REMOVED***
					reps -= 2
					for reps > 0 ***REMOVED***
						depth_histo[repeatZeroCodeLength]++

						/* Add the 3 extra bits for the 17 code length code. */
						bits += 3

						reps >>= 3
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		/* Add the estimated encoding cost of the code length code histogram. */
		bits += float64(18 + 2*max_depth)

		/* Add the entropy of the code length code histogram. */
		bits += bitsEntropy(depth_histo[:], codeLengthCodes)
	***REMOVED***

	return bits
***REMOVED***
