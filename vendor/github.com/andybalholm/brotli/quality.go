package brotli

const fastOnePassCompressionQuality = 0

const fastTwoPassCompressionQuality = 1

const zopflificationQuality = 10

const hqZopflificationQuality = 11

const maxQualityForStaticEntropyCodes = 2

const minQualityForBlockSplit = 4

const minQualityForNonzeroDistanceParams = 4

const minQualityForOptimizeHistograms = 4

const minQualityForExtensiveReferenceSearch = 5

const minQualityForContextModeling = 5

const minQualityForHqContextModeling = 7

const minQualityForHqBlockSplitting = 10

/* For quality below MIN_QUALITY_FOR_BLOCK_SPLIT there is no block splitting,
   so we buffer at most this much literals and commands. */
const maxNumDelayedSymbols = 0x2FFF

/* Returns hash-table size for quality levels 0 and 1. */
func maxHashTableSize(quality int) uint ***REMOVED***
	if quality == fastOnePassCompressionQuality ***REMOVED***
		return 1 << 15
	***REMOVED*** else ***REMOVED***
		return 1 << 17
	***REMOVED***
***REMOVED***

/* The maximum length for which the zopflification uses distinct distances. */
const maxZopfliLenQuality10 = 150

const maxZopfliLenQuality11 = 325

/* Do not thoroughly search when a long copy is found. */
const longCopyQuickStep = 16384

func maxZopfliLen(params *encoderParams) uint ***REMOVED***
	if params.quality <= 10 ***REMOVED***
		return maxZopfliLenQuality10
	***REMOVED*** else ***REMOVED***
		return maxZopfliLenQuality11
	***REMOVED***
***REMOVED***

/* Number of best candidates to evaluate to expand Zopfli chain. */
func maxZopfliCandidates(params *encoderParams) uint ***REMOVED***
	if params.quality <= 10 ***REMOVED***
		return 1
	***REMOVED*** else ***REMOVED***
		return 5
	***REMOVED***
***REMOVED***

func sanitizeParams(params *encoderParams) ***REMOVED***
	params.quality = brotli_min_int(maxQuality, brotli_max_int(minQuality, params.quality))
	if params.quality <= maxQualityForStaticEntropyCodes ***REMOVED***
		params.large_window = false
	***REMOVED***

	if params.lgwin < minWindowBits ***REMOVED***
		params.lgwin = minWindowBits
	***REMOVED*** else ***REMOVED***
		var max_lgwin int
		if params.large_window ***REMOVED***
			max_lgwin = largeMaxWindowBits
		***REMOVED*** else ***REMOVED***
			max_lgwin = maxWindowBits
		***REMOVED***
		if params.lgwin > uint(max_lgwin) ***REMOVED***
			params.lgwin = uint(max_lgwin)
		***REMOVED***
	***REMOVED***
***REMOVED***

/* Returns optimized lg_block value. */
func computeLgBlock(params *encoderParams) int ***REMOVED***
	var lgblock int = params.lgblock
	if params.quality == fastOnePassCompressionQuality || params.quality == fastTwoPassCompressionQuality ***REMOVED***
		lgblock = int(params.lgwin)
	***REMOVED*** else if params.quality < minQualityForBlockSplit ***REMOVED***
		lgblock = 14
	***REMOVED*** else if lgblock == 0 ***REMOVED***
		lgblock = 16
		if params.quality >= 9 && params.lgwin > uint(lgblock) ***REMOVED***
			lgblock = brotli_min_int(18, int(params.lgwin))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		lgblock = brotli_min_int(maxInputBlockBits, brotli_max_int(minInputBlockBits, lgblock))
	***REMOVED***

	return lgblock
***REMOVED***

/* Returns log2 of the size of main ring buffer area.
   Allocate at least lgwin + 1 bits for the ring buffer so that the newly
   added block fits there completely and we still get lgwin bits and at least
   read_block_size_bits + 1 bits because the copy tail length needs to be
   smaller than ring-buffer size. */
func computeRbBits(params *encoderParams) int ***REMOVED***
	return 1 + brotli_max_int(int(params.lgwin), params.lgblock)
***REMOVED***

func maxMetablockSize(params *encoderParams) uint ***REMOVED***
	var bits int = brotli_min_int(computeRbBits(params), maxInputBlockBits)
	return uint(1) << uint(bits)
***REMOVED***

/* When searching for backward references and have not seen matches for a long
   time, we can skip some match lookups. Unsuccessful match lookups are very
   expensive and this kind of a heuristic speeds up compression quite a lot.
   At first 8 byte strides are taken and every second byte is put to hasher.
   After 4x more literals stride by 16 bytes, every put 4-th byte to hasher.
   Applied only to qualities 2 to 9. */
func literalSpreeLengthForSparseSearch(params *encoderParams) uint ***REMOVED***
	if params.quality < 9 ***REMOVED***
		return 64
	***REMOVED*** else ***REMOVED***
		return 512
	***REMOVED***
***REMOVED***

func chooseHasher(params *encoderParams, hparams *hasherParams) ***REMOVED***
	if params.quality > 9 ***REMOVED***
		hparams.type_ = 10
	***REMOVED*** else if params.quality == 4 && params.size_hint >= 1<<20 ***REMOVED***
		hparams.type_ = 54
	***REMOVED*** else if params.quality < 5 ***REMOVED***
		hparams.type_ = params.quality
	***REMOVED*** else if params.lgwin <= 16 ***REMOVED***
		if params.quality < 7 ***REMOVED***
			hparams.type_ = 40
		***REMOVED*** else if params.quality < 9 ***REMOVED***
			hparams.type_ = 41
		***REMOVED*** else ***REMOVED***
			hparams.type_ = 42
		***REMOVED***
	***REMOVED*** else if params.size_hint >= 1<<20 && params.lgwin >= 19 ***REMOVED***
		hparams.type_ = 6
		hparams.block_bits = params.quality - 1
		hparams.bucket_bits = 15
		hparams.hash_len = 5
		if params.quality < 7 ***REMOVED***
			hparams.num_last_distances_to_check = 4
		***REMOVED*** else if params.quality < 9 ***REMOVED***
			hparams.num_last_distances_to_check = 10
		***REMOVED*** else ***REMOVED***
			hparams.num_last_distances_to_check = 16
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		hparams.type_ = 5
		hparams.block_bits = params.quality - 1
		if params.quality < 7 ***REMOVED***
			hparams.bucket_bits = 14
		***REMOVED*** else ***REMOVED***
			hparams.bucket_bits = 15
		***REMOVED***
		if params.quality < 7 ***REMOVED***
			hparams.num_last_distances_to_check = 4
		***REMOVED*** else if params.quality < 9 ***REMOVED***
			hparams.num_last_distances_to_check = 10
		***REMOVED*** else ***REMOVED***
			hparams.num_last_distances_to_check = 16
		***REMOVED***
	***REMOVED***

	if params.lgwin > 24 ***REMOVED***
		/* Different hashers for large window brotli: not for qualities <= 2,
		   these are too fast for large window. Not for qualities >= 10: their
		   hasher already works well with large window. So the changes are:
		   H3 --> H35: for quality 3.
		   H54 --> H55: for quality 4 with size hint > 1MB
		   H6 --> H65: for qualities 5, 6, 7, 8, 9. */
		if hparams.type_ == 3 ***REMOVED***
			hparams.type_ = 35
		***REMOVED***

		if hparams.type_ == 54 ***REMOVED***
			hparams.type_ = 55
		***REMOVED***

		if hparams.type_ == 6 ***REMOVED***
			hparams.type_ = 65
		***REMOVED***
	***REMOVED***
***REMOVED***
