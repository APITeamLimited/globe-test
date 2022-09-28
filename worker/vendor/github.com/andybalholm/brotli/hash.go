package brotli

import (
	"encoding/binary"
	"fmt"
)

type hasherCommon struct ***REMOVED***
	params           hasherParams
	is_prepared_     bool
	dict_num_lookups uint
	dict_num_matches uint
***REMOVED***

func (h *hasherCommon) Common() *hasherCommon ***REMOVED***
	return h
***REMOVED***

type hasherHandle interface ***REMOVED***
	Common() *hasherCommon
	Initialize(params *encoderParams)
	Prepare(one_shot bool, input_size uint, data []byte)
	StitchToPreviousBlock(num_bytes uint, position uint, ringbuffer []byte, ringbuffer_mask uint)
	HashTypeLength() uint
	StoreLookahead() uint
	PrepareDistanceCache(distance_cache []int)
	FindLongestMatch(dictionary *encoderDictionary, data []byte, ring_buffer_mask uint, distance_cache []int, cur_ix uint, max_length uint, max_backward uint, gap uint, max_distance uint, out *hasherSearchResult)
	StoreRange(data []byte, mask uint, ix_start uint, ix_end uint)
	Store(data []byte, mask uint, ix uint)
***REMOVED***

const kCutoffTransformsCount uint32 = 10

/*   0,  12,   27,    23,    42,    63,    56,    48,    59,    64 */
/* 0+0, 4+8, 8+19, 12+11, 16+26, 20+43, 24+32, 28+20, 32+27, 36+28 */
const kCutoffTransforms uint64 = 0x071B520ADA2D3200

type hasherSearchResult struct ***REMOVED***
	len            uint
	distance       uint
	score          uint
	len_code_delta int
***REMOVED***

/* kHashMul32 multiplier has these properties:
   * The multiplier must be odd. Otherwise we may lose the highest bit.
   * No long streaks of ones or zeros.
   * There is no effort to ensure that it is a prime, the oddity is enough
     for this use.
   * The number has been tuned heuristically against compression benchmarks. */
const kHashMul32 uint32 = 0x1E35A7BD

const kHashMul64 uint64 = 0x1E35A7BD1E35A7BD

const kHashMul64Long uint64 = 0x1FE35A7BD3579BD3

func hash14(data []byte) uint32 ***REMOVED***
	var h uint32 = binary.LittleEndian.Uint32(data) * kHashMul32

	/* The higher bits contain more mixture from the multiplication,
	   so we take our results from there. */
	return h >> (32 - 14)
***REMOVED***

func prepareDistanceCache(distance_cache []int, num_distances int) ***REMOVED***
	if num_distances > 4 ***REMOVED***
		var last_distance int = distance_cache[0]
		distance_cache[4] = last_distance - 1
		distance_cache[5] = last_distance + 1
		distance_cache[6] = last_distance - 2
		distance_cache[7] = last_distance + 2
		distance_cache[8] = last_distance - 3
		distance_cache[9] = last_distance + 3
		if num_distances > 10 ***REMOVED***
			var next_last_distance int = distance_cache[1]
			distance_cache[10] = next_last_distance - 1
			distance_cache[11] = next_last_distance + 1
			distance_cache[12] = next_last_distance - 2
			distance_cache[13] = next_last_distance + 2
			distance_cache[14] = next_last_distance - 3
			distance_cache[15] = next_last_distance + 3
		***REMOVED***
	***REMOVED***
***REMOVED***

const literalByteScore = 135

const distanceBitPenalty = 30

/* Score must be positive after applying maximal penalty. */
const scoreBase = (distanceBitPenalty * 8 * 8)

/* Usually, we always choose the longest backward reference. This function
   allows for the exception of that rule.

   If we choose a backward reference that is further away, it will
   usually be coded with more bits. We approximate this by assuming
   log2(distance). If the distance can be expressed in terms of the
   last four distances, we use some heuristic constants to estimate
   the bits cost. For the first up to four literals we use the bit
   cost of the literals from the literal cost model, after that we
   use the average bit cost of the cost model.

   This function is used to sometimes discard a longer backward reference
   when it is not much longer and the bit cost for encoding it is more
   than the saved literals.

   backward_reference_offset MUST be positive. */
func backwardReferenceScore(copy_length uint, backward_reference_offset uint) uint ***REMOVED***
	return scoreBase + literalByteScore*uint(copy_length) - distanceBitPenalty*uint(log2FloorNonZero(backward_reference_offset))
***REMOVED***

func backwardReferenceScoreUsingLastDistance(copy_length uint) uint ***REMOVED***
	return literalByteScore*uint(copy_length) + scoreBase + 15
***REMOVED***

func backwardReferencePenaltyUsingLastDistance(distance_short_code uint) uint ***REMOVED***
	return uint(39) + ((0x1CA10 >> (distance_short_code & 0xE)) & 0xE)
***REMOVED***

func testStaticDictionaryItem(dictionary *encoderDictionary, item uint, data []byte, max_length uint, max_backward uint, max_distance uint, out *hasherSearchResult) bool ***REMOVED***
	var len uint
	var word_idx uint
	var offset uint
	var matchlen uint
	var backward uint
	var score uint
	len = item & 0x1F
	word_idx = item >> 5
	offset = uint(dictionary.words.offsets_by_length[len]) + len*word_idx
	if len > max_length ***REMOVED***
		return false
	***REMOVED***

	matchlen = findMatchLengthWithLimit(data, dictionary.words.data[offset:], uint(len))
	if matchlen+uint(dictionary.cutoffTransformsCount) <= len || matchlen == 0 ***REMOVED***
		return false
	***REMOVED***
	***REMOVED***
		var cut uint = len - matchlen
		var transform_id uint = (cut << 2) + uint((dictionary.cutoffTransforms>>(cut*6))&0x3F)
		backward = max_backward + 1 + word_idx + (transform_id << dictionary.words.size_bits_by_length[len])
	***REMOVED***

	if backward > max_distance ***REMOVED***
		return false
	***REMOVED***

	score = backwardReferenceScore(matchlen, backward)
	if score < out.score ***REMOVED***
		return false
	***REMOVED***

	out.len = matchlen
	out.len_code_delta = int(len) - int(matchlen)
	out.distance = backward
	out.score = score
	return true
***REMOVED***

func searchInStaticDictionary(dictionary *encoderDictionary, handle hasherHandle, data []byte, max_length uint, max_backward uint, max_distance uint, out *hasherSearchResult, shallow bool) ***REMOVED***
	var key uint
	var i uint
	var self *hasherCommon = handle.Common()
	if self.dict_num_matches < self.dict_num_lookups>>7 ***REMOVED***
		return
	***REMOVED***

	key = uint(hash14(data) << 1)
	for i = 0; ; (func() ***REMOVED*** i++; key++ ***REMOVED***)() ***REMOVED***
		var tmp uint
		if shallow ***REMOVED***
			tmp = 1
		***REMOVED*** else ***REMOVED***
			tmp = 2
		***REMOVED***
		if i >= tmp ***REMOVED***
			break
		***REMOVED***
		var item uint = uint(dictionary.hash_table[key])
		self.dict_num_lookups++
		if item != 0 ***REMOVED***
			var item_matches bool = testStaticDictionaryItem(dictionary, item, data, max_length, max_backward, max_distance, out)
			if item_matches ***REMOVED***
				self.dict_num_matches++
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type backwardMatch struct ***REMOVED***
	distance        uint32
	length_and_code uint32
***REMOVED***

func initBackwardMatch(self *backwardMatch, dist uint, len uint) ***REMOVED***
	self.distance = uint32(dist)
	self.length_and_code = uint32(len << 5)
***REMOVED***

func initDictionaryBackwardMatch(self *backwardMatch, dist uint, len uint, len_code uint) ***REMOVED***
	self.distance = uint32(dist)
	var tmp uint
	if len == len_code ***REMOVED***
		tmp = 0
	***REMOVED*** else ***REMOVED***
		tmp = len_code
	***REMOVED***
	self.length_and_code = uint32(len<<5 | tmp)
***REMOVED***

func backwardMatchLength(self *backwardMatch) uint ***REMOVED***
	return uint(self.length_and_code >> 5)
***REMOVED***

func backwardMatchLengthCode(self *backwardMatch) uint ***REMOVED***
	var code uint = uint(self.length_and_code) & 31
	if code != 0 ***REMOVED***
		return code
	***REMOVED*** else ***REMOVED***
		return backwardMatchLength(self)
	***REMOVED***
***REMOVED***

func hasherReset(handle hasherHandle) ***REMOVED***
	if handle == nil ***REMOVED***
		return
	***REMOVED***
	handle.Common().is_prepared_ = false
***REMOVED***

func newHasher(typ int) hasherHandle ***REMOVED***
	switch typ ***REMOVED***
	case 2:
		return &hashLongestMatchQuickly***REMOVED***
			bucketBits:    16,
			bucketSweep:   1,
			hashLen:       5,
			useDictionary: true,
		***REMOVED***
	case 3:
		return &hashLongestMatchQuickly***REMOVED***
			bucketBits:    16,
			bucketSweep:   2,
			hashLen:       5,
			useDictionary: false,
		***REMOVED***
	case 4:
		return &hashLongestMatchQuickly***REMOVED***
			bucketBits:    17,
			bucketSweep:   4,
			hashLen:       5,
			useDictionary: true,
		***REMOVED***
	case 5:
		return new(h5)
	case 6:
		return new(h6)
	case 10:
		return new(h10)
	case 35:
		return &hashComposite***REMOVED***
			ha: newHasher(3),
			hb: &hashRolling***REMOVED***jump: 4***REMOVED***,
		***REMOVED***
	case 40:
		return &hashForgetfulChain***REMOVED***
			bucketBits:              15,
			numBanks:                1,
			bankBits:                16,
			numLastDistancesToCheck: 4,
		***REMOVED***
	case 41:
		return &hashForgetfulChain***REMOVED***
			bucketBits:              15,
			numBanks:                1,
			bankBits:                16,
			numLastDistancesToCheck: 10,
		***REMOVED***
	case 42:
		return &hashForgetfulChain***REMOVED***
			bucketBits:              15,
			numBanks:                512,
			bankBits:                9,
			numLastDistancesToCheck: 16,
		***REMOVED***
	case 54:
		return &hashLongestMatchQuickly***REMOVED***
			bucketBits:    20,
			bucketSweep:   4,
			hashLen:       7,
			useDictionary: false,
		***REMOVED***
	case 55:
		return &hashComposite***REMOVED***
			ha: newHasher(54),
			hb: &hashRolling***REMOVED***jump: 4***REMOVED***,
		***REMOVED***
	case 65:
		return &hashComposite***REMOVED***
			ha: newHasher(6),
			hb: &hashRolling***REMOVED***jump: 1***REMOVED***,
		***REMOVED***
	***REMOVED***

	panic(fmt.Sprintf("unknown hasher type: %d", typ))
***REMOVED***

func hasherSetup(handle *hasherHandle, params *encoderParams, data []byte, position uint, input_size uint, is_last bool) ***REMOVED***
	var self hasherHandle = nil
	var common *hasherCommon = nil
	var one_shot bool = (position == 0 && is_last)
	if *handle == nil ***REMOVED***
		chooseHasher(params, &params.hasher)
		self = newHasher(params.hasher.type_)

		*handle = self
		common = self.Common()
		common.params = params.hasher
		self.Initialize(params)
	***REMOVED***

	self = *handle
	common = self.Common()
	if !common.is_prepared_ ***REMOVED***
		self.Prepare(one_shot, input_size, data)

		if position == 0 ***REMOVED***
			common.dict_num_lookups = 0
			common.dict_num_matches = 0
		***REMOVED***

		common.is_prepared_ = true
	***REMOVED***
***REMOVED***

func initOrStitchToPreviousBlock(handle *hasherHandle, data []byte, mask uint, params *encoderParams, position uint, input_size uint, is_last bool) ***REMOVED***
	var self hasherHandle
	hasherSetup(handle, params, data, position, input_size, is_last)
	self = *handle
	self.StitchToPreviousBlock(input_size, position, data, mask)
***REMOVED***
