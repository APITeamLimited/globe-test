package brotli

/* Dictionary data (words and transforms) for 1 possible context */
type encoderDictionary struct ***REMOVED***
	words                 *dictionary
	cutoffTransformsCount uint32
	cutoffTransforms      uint64
	hash_table            []uint16
	buckets               []uint16
	dict_words            []dictWord
***REMOVED***

func initEncoderDictionary(dict *encoderDictionary) ***REMOVED***
	dict.words = getDictionary()

	dict.hash_table = kStaticDictionaryHash[:]
	dict.buckets = kStaticDictionaryBuckets[:]
	dict.dict_words = kStaticDictionaryWords[:]

	dict.cutoffTransformsCount = kCutoffTransformsCount
	dict.cutoffTransforms = kCutoffTransforms
***REMOVED***
