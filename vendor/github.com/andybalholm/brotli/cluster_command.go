package brotli

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Computes the bit cost reduction by combining out[idx1] and out[idx2] and if
   it is below a threshold, stores the pair (idx1, idx2) in the *pairs queue. */
func compareAndPushToQueueCommand(out []histogramCommand, cluster_size []uint32, idx1 uint32, idx2 uint32, max_num_pairs uint, pairs []histogramPair, num_pairs *uint) ***REMOVED***
	var is_good_pair bool = false
	var p histogramPair
	p.idx2 = 0
	p.idx1 = p.idx2
	p.cost_combo = 0
	p.cost_diff = p.cost_combo
	if idx1 == idx2 ***REMOVED***
		return
	***REMOVED***

	if idx2 < idx1 ***REMOVED***
		var t uint32 = idx2
		idx2 = idx1
		idx1 = t
	***REMOVED***

	p.idx1 = idx1
	p.idx2 = idx2
	p.cost_diff = 0.5 * clusterCostDiff(uint(cluster_size[idx1]), uint(cluster_size[idx2]))
	p.cost_diff -= out[idx1].bit_cost_
	p.cost_diff -= out[idx2].bit_cost_

	if out[idx1].total_count_ == 0 ***REMOVED***
		p.cost_combo = out[idx2].bit_cost_
		is_good_pair = true
	***REMOVED*** else if out[idx2].total_count_ == 0 ***REMOVED***
		p.cost_combo = out[idx1].bit_cost_
		is_good_pair = true
	***REMOVED*** else ***REMOVED***
		var threshold float64
		if *num_pairs == 0 ***REMOVED***
			threshold = 1e99
		***REMOVED*** else ***REMOVED***
			threshold = brotli_max_double(0.0, pairs[0].cost_diff)
		***REMOVED***
		var combo histogramCommand = out[idx1]
		var cost_combo float64
		histogramAddHistogramCommand(&combo, &out[idx2])
		cost_combo = populationCostCommand(&combo)
		if cost_combo < threshold-p.cost_diff ***REMOVED***
			p.cost_combo = cost_combo
			is_good_pair = true
		***REMOVED***
	***REMOVED***

	if is_good_pair ***REMOVED***
		p.cost_diff += p.cost_combo
		if *num_pairs > 0 && histogramPairIsLess(&pairs[0], &p) ***REMOVED***
			/* Replace the top of the queue if needed. */
			if *num_pairs < max_num_pairs ***REMOVED***
				pairs[*num_pairs] = pairs[0]
				(*num_pairs)++
			***REMOVED***

			pairs[0] = p
		***REMOVED*** else if *num_pairs < max_num_pairs ***REMOVED***
			pairs[*num_pairs] = p
			(*num_pairs)++
		***REMOVED***
	***REMOVED***
***REMOVED***

func histogramCombineCommand(out []histogramCommand, cluster_size []uint32, symbols []uint32, clusters []uint32, pairs []histogramPair, num_clusters uint, symbols_size uint, max_clusters uint, max_num_pairs uint) uint ***REMOVED***
	var cost_diff_threshold float64 = 0.0
	var min_cluster_size uint = 1
	var num_pairs uint = 0
	***REMOVED***
		/* We maintain a vector of histogram pairs, with the property that the pair
		   with the maximum bit cost reduction is the first. */
		var idx1 uint
		for idx1 = 0; idx1 < num_clusters; idx1++ ***REMOVED***
			var idx2 uint
			for idx2 = idx1 + 1; idx2 < num_clusters; idx2++ ***REMOVED***
				compareAndPushToQueueCommand(out, cluster_size, clusters[idx1], clusters[idx2], max_num_pairs, pairs[0:], &num_pairs)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for num_clusters > min_cluster_size ***REMOVED***
		var best_idx1 uint32
		var best_idx2 uint32
		var i uint
		if pairs[0].cost_diff >= cost_diff_threshold ***REMOVED***
			cost_diff_threshold = 1e99
			min_cluster_size = max_clusters
			continue
		***REMOVED***

		/* Take the best pair from the top of heap. */
		best_idx1 = pairs[0].idx1

		best_idx2 = pairs[0].idx2
		histogramAddHistogramCommand(&out[best_idx1], &out[best_idx2])
		out[best_idx1].bit_cost_ = pairs[0].cost_combo
		cluster_size[best_idx1] += cluster_size[best_idx2]
		for i = 0; i < symbols_size; i++ ***REMOVED***
			if symbols[i] == best_idx2 ***REMOVED***
				symbols[i] = best_idx1
			***REMOVED***
		***REMOVED***

		for i = 0; i < num_clusters; i++ ***REMOVED***
			if clusters[i] == best_idx2 ***REMOVED***
				copy(clusters[i:], clusters[i+1:][:num_clusters-i-1])
				break
			***REMOVED***
		***REMOVED***

		num_clusters--
		***REMOVED***
			/* Remove pairs intersecting the just combined best pair. */
			var copy_to_idx uint = 0
			for i = 0; i < num_pairs; i++ ***REMOVED***
				var p *histogramPair = &pairs[i]
				if p.idx1 == best_idx1 || p.idx2 == best_idx1 || p.idx1 == best_idx2 || p.idx2 == best_idx2 ***REMOVED***
					/* Remove invalid pair from the queue. */
					continue
				***REMOVED***

				if histogramPairIsLess(&pairs[0], p) ***REMOVED***
					/* Replace the top of the queue if needed. */
					var front histogramPair = pairs[0]
					pairs[0] = *p
					pairs[copy_to_idx] = front
				***REMOVED*** else ***REMOVED***
					pairs[copy_to_idx] = *p
				***REMOVED***

				copy_to_idx++
			***REMOVED***

			num_pairs = copy_to_idx
		***REMOVED***

		/* Push new pairs formed with the combined histogram to the heap. */
		for i = 0; i < num_clusters; i++ ***REMOVED***
			compareAndPushToQueueCommand(out, cluster_size, best_idx1, clusters[i], max_num_pairs, pairs[0:], &num_pairs)
		***REMOVED***
	***REMOVED***

	return num_clusters
***REMOVED***

/* What is the bit cost of moving histogram from cur_symbol to candidate. */
func histogramBitCostDistanceCommand(histogram *histogramCommand, candidate *histogramCommand) float64 ***REMOVED***
	if histogram.total_count_ == 0 ***REMOVED***
		return 0.0
	***REMOVED*** else ***REMOVED***
		var tmp histogramCommand = *histogram
		histogramAddHistogramCommand(&tmp, candidate)
		return populationCostCommand(&tmp) - candidate.bit_cost_
	***REMOVED***
***REMOVED***
