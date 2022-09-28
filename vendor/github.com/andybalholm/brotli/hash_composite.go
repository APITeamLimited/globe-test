package brotli

/* Copyright 2018 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

func (h *hashComposite) HashTypeLength() uint ***REMOVED***
	var a uint = h.ha.HashTypeLength()
	var b uint = h.hb.HashTypeLength()
	if a > b ***REMOVED***
		return a
	***REMOVED*** else ***REMOVED***
		return b
	***REMOVED***
***REMOVED***

func (h *hashComposite) StoreLookahead() uint ***REMOVED***
	var a uint = h.ha.StoreLookahead()
	var b uint = h.hb.StoreLookahead()
	if a > b ***REMOVED***
		return a
	***REMOVED*** else ***REMOVED***
		return b
	***REMOVED***
***REMOVED***

/* Composite hasher: This hasher allows to combine two other hashers, HASHER_A
   and HASHER_B. */
type hashComposite struct ***REMOVED***
	hasherCommon
	ha     hasherHandle
	hb     hasherHandle
	params *encoderParams
***REMOVED***

func (h *hashComposite) Initialize(params *encoderParams) ***REMOVED***
	h.params = params
***REMOVED***

/* TODO: Initialize of the hashers is defered to Prepare (and params
   remembered here) because we don't get the one_shot and input_size params
   here that are needed to know the memory size of them. Instead provide
   those params to all hashers InitializehashComposite */
func (h *hashComposite) Prepare(one_shot bool, input_size uint, data []byte) ***REMOVED***
	if h.ha == nil ***REMOVED***
		var common_a *hasherCommon
		var common_b *hasherCommon

		common_a = h.ha.Common()
		common_a.params = h.params.hasher
		common_a.is_prepared_ = false
		common_a.dict_num_lookups = 0
		common_a.dict_num_matches = 0
		h.ha.Initialize(h.params)

		common_b = h.hb.Common()
		common_b.params = h.params.hasher
		common_b.is_prepared_ = false
		common_b.dict_num_lookups = 0
		common_b.dict_num_matches = 0
		h.hb.Initialize(h.params)
	***REMOVED***

	h.ha.Prepare(one_shot, input_size, data)
	h.hb.Prepare(one_shot, input_size, data)
***REMOVED***

func (h *hashComposite) Store(data []byte, mask uint, ix uint) ***REMOVED***
	h.ha.Store(data, mask, ix)
	h.hb.Store(data, mask, ix)
***REMOVED***

func (h *hashComposite) StoreRange(data []byte, mask uint, ix_start uint, ix_end uint) ***REMOVED***
	h.ha.StoreRange(data, mask, ix_start, ix_end)
	h.hb.StoreRange(data, mask, ix_start, ix_end)
***REMOVED***

func (h *hashComposite) StitchToPreviousBlock(num_bytes uint, position uint, ringbuffer []byte, ring_buffer_mask uint) ***REMOVED***
	h.ha.StitchToPreviousBlock(num_bytes, position, ringbuffer, ring_buffer_mask)
	h.hb.StitchToPreviousBlock(num_bytes, position, ringbuffer, ring_buffer_mask)
***REMOVED***

func (h *hashComposite) PrepareDistanceCache(distance_cache []int) ***REMOVED***
	h.ha.PrepareDistanceCache(distance_cache)
	h.hb.PrepareDistanceCache(distance_cache)
***REMOVED***

func (h *hashComposite) FindLongestMatch(dictionary *encoderDictionary, data []byte, ring_buffer_mask uint, distance_cache []int, cur_ix uint, max_length uint, max_backward uint, gap uint, max_distance uint, out *hasherSearchResult) ***REMOVED***
	h.ha.FindLongestMatch(dictionary, data, ring_buffer_mask, distance_cache, cur_ix, max_length, max_backward, gap, max_distance, out)
	h.hb.FindLongestMatch(dictionary, data, ring_buffer_mask, distance_cache, cur_ix, max_length, max_backward, gap, max_distance, out)
***REMOVED***
