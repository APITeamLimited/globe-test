package brotli

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

const (
	decoderResultError           = 0
	decoderResultSuccess         = 1
	decoderResultNeedsMoreInput  = 2
	decoderResultNeedsMoreOutput = 3
)

/**
 * Error code for detailed logging / production debugging.
 *
 * See ::BrotliDecoderGetErrorCode and ::BROTLI_LAST_ERROR_CODE.
 */
const (
	decoderNoError                          = 0
	decoderSuccess                          = 1
	decoderNeedsMoreInput                   = 2
	decoderNeedsMoreOutput                  = 3
	decoderErrorFormatExuberantNibble       = -1
	decoderErrorFormatReserved              = -2
	decoderErrorFormatExuberantMetaNibble   = -3
	decoderErrorFormatSimpleHuffmanAlphabet = -4
	decoderErrorFormatSimpleHuffmanSame     = -5
	decoderErrorFormatClSpace               = -6
	decoderErrorFormatHuffmanSpace          = -7
	decoderErrorFormatContextMapRepeat      = -8
	decoderErrorFormatBlockLength1          = -9
	decoderErrorFormatBlockLength2          = -10
	decoderErrorFormatTransform             = -11
	decoderErrorFormatDictionary            = -12
	decoderErrorFormatWindowBits            = -13
	decoderErrorFormatPadding1              = -14
	decoderErrorFormatPadding2              = -15
	decoderErrorFormatDistance              = -16
	decoderErrorDictionaryNotSet            = -19
	decoderErrorInvalidArguments            = -20
	decoderErrorAllocContextModes           = -21
	decoderErrorAllocTreeGroups             = -22
	decoderErrorAllocContextMap             = -25
	decoderErrorAllocRingBuffer1            = -26
	decoderErrorAllocRingBuffer2            = -27
	decoderErrorAllocBlockTypeTrees         = -30
	decoderErrorUnreachable                 = -31
)

const huffmanTableBits = 8

const huffmanTableMask = 0xFF

/* We need the slack region for the following reasons:
   - doing up to two 16-byte copies for fast backward copying
   - inserting transformed dictionary word (5 prefix + 24 base + 8 suffix) */
const kRingBufferWriteAheadSlack uint32 = 42

var kCodeLengthCodeOrder = [codeLengthCodes]byte***REMOVED***1, 2, 3, 4, 0, 5, 17, 6, 16, 7, 8, 9, 10, 11, 12, 13, 14, 15***REMOVED***

/* Static prefix code for the complex code length code lengths. */
var kCodeLengthPrefixLength = [16]byte***REMOVED***2, 2, 2, 3, 2, 2, 2, 4, 2, 2, 2, 3, 2, 2, 2, 4***REMOVED***

var kCodeLengthPrefixValue = [16]byte***REMOVED***0, 4, 3, 2, 0, 4, 3, 1, 0, 4, 3, 2, 0, 4, 3, 5***REMOVED***

/* Saves error code and converts it to BrotliDecoderResult. */
func saveErrorCode(s *Reader, e int) int ***REMOVED***
	s.error_code = int(e)
	switch e ***REMOVED***
	case decoderSuccess:
		return decoderResultSuccess

	case decoderNeedsMoreInput:
		return decoderResultNeedsMoreInput

	case decoderNeedsMoreOutput:
		return decoderResultNeedsMoreOutput

	default:
		return decoderResultError
	***REMOVED***
***REMOVED***

/* Decodes WBITS by reading 1 - 7 bits, or 0x11 for "Large Window Brotli".
   Precondition: bit-reader accumulator has at least 8 bits. */
func decodeWindowBits(s *Reader, br *bitReader) int ***REMOVED***
	var n uint32
	var large_window bool = s.large_window
	s.large_window = false
	takeBits(br, 1, &n)
	if n == 0 ***REMOVED***
		s.window_bits = 16
		return decoderSuccess
	***REMOVED***

	takeBits(br, 3, &n)
	if n != 0 ***REMOVED***
		s.window_bits = 17 + n
		return decoderSuccess
	***REMOVED***

	takeBits(br, 3, &n)
	if n == 1 ***REMOVED***
		if large_window ***REMOVED***
			takeBits(br, 1, &n)
			if n == 1 ***REMOVED***
				return decoderErrorFormatWindowBits
			***REMOVED***

			s.large_window = true
			return decoderSuccess
		***REMOVED*** else ***REMOVED***
			return decoderErrorFormatWindowBits
		***REMOVED***
	***REMOVED***

	if n != 0 ***REMOVED***
		s.window_bits = 8 + n
		return decoderSuccess
	***REMOVED***

	s.window_bits = 17
	return decoderSuccess
***REMOVED***

/* Decodes a number in the range [0..255], by reading 1 - 11 bits. */
func decodeVarLenUint8(s *Reader, br *bitReader, value *uint32) int ***REMOVED***
	var bits uint32
	switch s.substate_decode_uint8 ***REMOVED***
	case stateDecodeUint8None:
		if !safeReadBits(br, 1, &bits) ***REMOVED***
			return decoderNeedsMoreInput
		***REMOVED***

		if bits == 0 ***REMOVED***
			*value = 0
			return decoderSuccess
		***REMOVED***
		fallthrough

		/* Fall through. */
	case stateDecodeUint8Short:
		if !safeReadBits(br, 3, &bits) ***REMOVED***
			s.substate_decode_uint8 = stateDecodeUint8Short
			return decoderNeedsMoreInput
		***REMOVED***

		if bits == 0 ***REMOVED***
			*value = 1
			s.substate_decode_uint8 = stateDecodeUint8None
			return decoderSuccess
		***REMOVED***

		/* Use output value as a temporary storage. It MUST be persisted. */
		*value = bits
		fallthrough

		/* Fall through. */
	case stateDecodeUint8Long:
		if !safeReadBits(br, *value, &bits) ***REMOVED***
			s.substate_decode_uint8 = stateDecodeUint8Long
			return decoderNeedsMoreInput
		***REMOVED***

		*value = (1 << *value) + bits
		s.substate_decode_uint8 = stateDecodeUint8None
		return decoderSuccess

	default:
		return decoderErrorUnreachable
	***REMOVED***
***REMOVED***

/* Decodes a metablock length and flags by reading 2 - 31 bits. */
func decodeMetaBlockLength(s *Reader, br *bitReader) int ***REMOVED***
	var bits uint32
	var i int
	for ***REMOVED***
		switch s.substate_metablock_header ***REMOVED***
		case stateMetablockHeaderNone:
			if !safeReadBits(br, 1, &bits) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			if bits != 0 ***REMOVED***
				s.is_last_metablock = 1
			***REMOVED*** else ***REMOVED***
				s.is_last_metablock = 0
			***REMOVED***
			s.meta_block_remaining_len = 0
			s.is_uncompressed = 0
			s.is_metadata = 0
			if s.is_last_metablock == 0 ***REMOVED***
				s.substate_metablock_header = stateMetablockHeaderNibbles
				break
			***REMOVED***

			s.substate_metablock_header = stateMetablockHeaderEmpty
			fallthrough

			/* Fall through. */
		case stateMetablockHeaderEmpty:
			if !safeReadBits(br, 1, &bits) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			if bits != 0 ***REMOVED***
				s.substate_metablock_header = stateMetablockHeaderNone
				return decoderSuccess
			***REMOVED***

			s.substate_metablock_header = stateMetablockHeaderNibbles
			fallthrough

			/* Fall through. */
		case stateMetablockHeaderNibbles:
			if !safeReadBits(br, 2, &bits) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			s.size_nibbles = uint(byte(bits + 4))
			s.loop_counter = 0
			if bits == 3 ***REMOVED***
				s.is_metadata = 1
				s.substate_metablock_header = stateMetablockHeaderReserved
				break
			***REMOVED***

			s.substate_metablock_header = stateMetablockHeaderSize
			fallthrough

			/* Fall through. */
		case stateMetablockHeaderSize:
			i = s.loop_counter

			for ; i < int(s.size_nibbles); i++ ***REMOVED***
				if !safeReadBits(br, 4, &bits) ***REMOVED***
					s.loop_counter = i
					return decoderNeedsMoreInput
				***REMOVED***

				if uint(i+1) == s.size_nibbles && s.size_nibbles > 4 && bits == 0 ***REMOVED***
					return decoderErrorFormatExuberantNibble
				***REMOVED***

				s.meta_block_remaining_len |= int(bits << uint(i*4))
			***REMOVED***

			s.substate_metablock_header = stateMetablockHeaderUncompressed
			fallthrough

			/* Fall through. */
		case stateMetablockHeaderUncompressed:
			if s.is_last_metablock == 0 ***REMOVED***
				if !safeReadBits(br, 1, &bits) ***REMOVED***
					return decoderNeedsMoreInput
				***REMOVED***

				if bits != 0 ***REMOVED***
					s.is_uncompressed = 1
				***REMOVED*** else ***REMOVED***
					s.is_uncompressed = 0
				***REMOVED***
			***REMOVED***

			s.meta_block_remaining_len++
			s.substate_metablock_header = stateMetablockHeaderNone
			return decoderSuccess

		case stateMetablockHeaderReserved:
			if !safeReadBits(br, 1, &bits) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			if bits != 0 ***REMOVED***
				return decoderErrorFormatReserved
			***REMOVED***

			s.substate_metablock_header = stateMetablockHeaderBytes
			fallthrough

			/* Fall through. */
		case stateMetablockHeaderBytes:
			if !safeReadBits(br, 2, &bits) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			if bits == 0 ***REMOVED***
				s.substate_metablock_header = stateMetablockHeaderNone
				return decoderSuccess
			***REMOVED***

			s.size_nibbles = uint(byte(bits))
			s.substate_metablock_header = stateMetablockHeaderMetadata
			fallthrough

			/* Fall through. */
		case stateMetablockHeaderMetadata:
			i = s.loop_counter

			for ; i < int(s.size_nibbles); i++ ***REMOVED***
				if !safeReadBits(br, 8, &bits) ***REMOVED***
					s.loop_counter = i
					return decoderNeedsMoreInput
				***REMOVED***

				if uint(i+1) == s.size_nibbles && s.size_nibbles > 1 && bits == 0 ***REMOVED***
					return decoderErrorFormatExuberantMetaNibble
				***REMOVED***

				s.meta_block_remaining_len |= int(bits << uint(i*8))
			***REMOVED***

			s.meta_block_remaining_len++
			s.substate_metablock_header = stateMetablockHeaderNone
			return decoderSuccess

		default:
			return decoderErrorUnreachable
		***REMOVED***
	***REMOVED***
***REMOVED***

/* Decodes the Huffman code.
   This method doesn't read data from the bit reader, BUT drops the amount of
   bits that correspond to the decoded symbol.
   bits MUST contain at least 15 (BROTLI_HUFFMAN_MAX_CODE_LENGTH) valid bits. */
func decodeSymbol(bits uint32, table []huffmanCode, br *bitReader) uint32 ***REMOVED***
	table = table[bits&huffmanTableMask:]
	if table[0].bits > huffmanTableBits ***REMOVED***
		var nbits uint32 = uint32(table[0].bits) - huffmanTableBits
		dropBits(br, huffmanTableBits)
		table = table[uint32(table[0].value)+((bits>>huffmanTableBits)&bitMask(nbits)):]
	***REMOVED***

	dropBits(br, uint32(table[0].bits))
	return uint32(table[0].value)
***REMOVED***

/* Reads and decodes the next Huffman code from bit-stream.
   This method peeks 16 bits of input and drops 0 - 15 of them. */
func readSymbol(table []huffmanCode, br *bitReader) uint32 ***REMOVED***
	return decodeSymbol(get16BitsUnmasked(br), table, br)
***REMOVED***

/* Same as DecodeSymbol, but it is known that there is less than 15 bits of
   input are currently available. */
func safeDecodeSymbol(table []huffmanCode, br *bitReader, result *uint32) bool ***REMOVED***
	var val uint32
	var available_bits uint32 = getAvailableBits(br)
	if available_bits == 0 ***REMOVED***
		if table[0].bits == 0 ***REMOVED***
			*result = uint32(table[0].value)
			return true
		***REMOVED***

		return false /* No valid bits at all. */
	***REMOVED***

	val = uint32(getBitsUnmasked(br))
	table = table[val&huffmanTableMask:]
	if table[0].bits <= huffmanTableBits ***REMOVED***
		if uint32(table[0].bits) <= available_bits ***REMOVED***
			dropBits(br, uint32(table[0].bits))
			*result = uint32(table[0].value)
			return true
		***REMOVED*** else ***REMOVED***
			return false /* Not enough bits for the first level. */
		***REMOVED***
	***REMOVED***

	if available_bits <= huffmanTableBits ***REMOVED***
		return false /* Not enough bits to move to the second level. */
	***REMOVED***

	/* Speculatively drop HUFFMAN_TABLE_BITS. */
	val = (val & bitMask(uint32(table[0].bits))) >> huffmanTableBits

	available_bits -= huffmanTableBits
	table = table[uint32(table[0].value)+val:]
	if available_bits < uint32(table[0].bits) ***REMOVED***
		return false /* Not enough bits for the second level. */
	***REMOVED***

	dropBits(br, huffmanTableBits+uint32(table[0].bits))
	*result = uint32(table[0].value)
	return true
***REMOVED***

func safeReadSymbol(table []huffmanCode, br *bitReader, result *uint32) bool ***REMOVED***
	var val uint32
	if safeGetBits(br, 15, &val) ***REMOVED***
		*result = decodeSymbol(val, table, br)
		return true
	***REMOVED***

	return safeDecodeSymbol(table, br, result)
***REMOVED***

/* Makes a look-up in first level Huffman table. Peeks 8 bits. */
func preloadSymbol(safe int, table []huffmanCode, br *bitReader, bits *uint32, value *uint32) ***REMOVED***
	if safe != 0 ***REMOVED***
		return
	***REMOVED***

	table = table[getBits(br, huffmanTableBits):]
	*bits = uint32(table[0].bits)
	*value = uint32(table[0].value)
***REMOVED***

/* Decodes the next Huffman code using data prepared by PreloadSymbol.
   Reads 0 - 15 bits. Also peeks 8 following bits. */
func readPreloadedSymbol(table []huffmanCode, br *bitReader, bits *uint32, value *uint32) uint32 ***REMOVED***
	var result uint32 = *value
	var ext []huffmanCode
	if *bits > huffmanTableBits ***REMOVED***
		var val uint32 = get16BitsUnmasked(br)
		ext = table[val&huffmanTableMask:][*value:]
		var mask uint32 = bitMask((*bits - huffmanTableBits))
		dropBits(br, huffmanTableBits)
		ext = ext[(val>>huffmanTableBits)&mask:]
		dropBits(br, uint32(ext[0].bits))
		result = uint32(ext[0].value)
	***REMOVED*** else ***REMOVED***
		dropBits(br, *bits)
	***REMOVED***

	preloadSymbol(0, table, br, bits, value)
	return result
***REMOVED***

func log2Floor(x uint32) uint32 ***REMOVED***
	var result uint32 = 0
	for x != 0 ***REMOVED***
		x >>= 1
		result++
	***REMOVED***

	return result
***REMOVED***

/* Reads (s->symbol + 1) symbols.
   Totally 1..4 symbols are read, 1..11 bits each.
   The list of symbols MUST NOT contain duplicates. */
func readSimpleHuffmanSymbols(alphabet_size uint32, max_symbol uint32, s *Reader) int ***REMOVED***
	var br *bitReader = &s.br
	var max_bits uint32 = log2Floor(alphabet_size - 1)
	var i uint32 = s.sub_loop_counter
	/* max_bits == 1..11; symbol == 0..3; 1..44 bits will be read. */

	var num_symbols uint32 = s.symbol
	for i <= num_symbols ***REMOVED***
		var v uint32
		if !safeReadBits(br, max_bits, &v) ***REMOVED***
			s.sub_loop_counter = i
			s.substate_huffman = stateHuffmanSimpleRead
			return decoderNeedsMoreInput
		***REMOVED***

		if v >= max_symbol ***REMOVED***
			return decoderErrorFormatSimpleHuffmanAlphabet
		***REMOVED***

		s.symbols_lists_array[i] = uint16(v)
		i++
	***REMOVED***

	for i = 0; i < num_symbols; i++ ***REMOVED***
		var k uint32 = i + 1
		for ; k <= num_symbols; k++ ***REMOVED***
			if s.symbols_lists_array[i] == s.symbols_lists_array[k] ***REMOVED***
				return decoderErrorFormatSimpleHuffmanSame
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return decoderSuccess
***REMOVED***

/* Process single decoded symbol code length:
   A) reset the repeat variable
   B) remember code length (if it is not 0)
   C) extend corresponding index-chain
   D) reduce the Huffman space
   E) update the histogram */
func processSingleCodeLength(code_len uint32, symbol *uint32, repeat *uint32, space *uint32, prev_code_len *uint32, symbol_lists symbolList, code_length_histo []uint16, next_symbol []int) ***REMOVED***
	*repeat = 0
	if code_len != 0 ***REMOVED*** /* code_len == 1..15 */
		symbolListPut(symbol_lists, next_symbol[code_len], uint16(*symbol))
		next_symbol[code_len] = int(*symbol)
		*prev_code_len = code_len
		*space -= 32768 >> code_len
		code_length_histo[code_len]++
	***REMOVED***

	(*symbol)++
***REMOVED***

/* Process repeated symbol code length.
    A) Check if it is the extension of previous repeat sequence; if the decoded
       value is not BROTLI_REPEAT_PREVIOUS_CODE_LENGTH, then it is a new
       symbol-skip
    B) Update repeat variable
    C) Check if operation is feasible (fits alphabet)
    D) For each symbol do the same operations as in ProcessSingleCodeLength

   PRECONDITION: code_len == BROTLI_REPEAT_PREVIOUS_CODE_LENGTH or
                 code_len == BROTLI_REPEAT_ZERO_CODE_LENGTH */
func processRepeatedCodeLength(code_len uint32, repeat_delta uint32, alphabet_size uint32, symbol *uint32, repeat *uint32, space *uint32, prev_code_len *uint32, repeat_code_len *uint32, symbol_lists symbolList, code_length_histo []uint16, next_symbol []int) ***REMOVED***
	var old_repeat uint32 /* for BROTLI_REPEAT_ZERO_CODE_LENGTH */ /* for BROTLI_REPEAT_ZERO_CODE_LENGTH */
	var extra_bits uint32 = 3
	var new_len uint32 = 0
	if code_len == repeatPreviousCodeLength ***REMOVED***
		new_len = *prev_code_len
		extra_bits = 2
	***REMOVED***

	if *repeat_code_len != new_len ***REMOVED***
		*repeat = 0
		*repeat_code_len = new_len
	***REMOVED***

	old_repeat = *repeat
	if *repeat > 0 ***REMOVED***
		*repeat -= 2
		*repeat <<= extra_bits
	***REMOVED***

	*repeat += repeat_delta + 3
	repeat_delta = *repeat - old_repeat
	if *symbol+repeat_delta > alphabet_size ***REMOVED***
		*symbol = alphabet_size
		*space = 0xFFFFF
		return
	***REMOVED***

	if *repeat_code_len != 0 ***REMOVED***
		var last uint = uint(*symbol + repeat_delta)
		var next int = next_symbol[*repeat_code_len]
		for ***REMOVED***
			symbolListPut(symbol_lists, next, uint16(*symbol))
			next = int(*symbol)
			(*symbol)++
			if (*symbol) == uint32(last) ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		next_symbol[*repeat_code_len] = next
		*space -= repeat_delta << (15 - *repeat_code_len)
		code_length_histo[*repeat_code_len] = uint16(uint32(code_length_histo[*repeat_code_len]) + repeat_delta)
	***REMOVED*** else ***REMOVED***
		*symbol += repeat_delta
	***REMOVED***
***REMOVED***

/* Reads and decodes symbol codelengths. */
func readSymbolCodeLengths(alphabet_size uint32, s *Reader) int ***REMOVED***
	var br *bitReader = &s.br
	var symbol uint32 = s.symbol
	var repeat uint32 = s.repeat
	var space uint32 = s.space
	var prev_code_len uint32 = s.prev_code_len
	var repeat_code_len uint32 = s.repeat_code_len
	var symbol_lists symbolList = s.symbol_lists
	var code_length_histo []uint16 = s.code_length_histo[:]
	var next_symbol []int = s.next_symbol[:]
	if !warmupBitReader(br) ***REMOVED***
		return decoderNeedsMoreInput
	***REMOVED***
	var p []huffmanCode
	for symbol < alphabet_size && space > 0 ***REMOVED***
		p = s.table[:]
		var code_len uint32
		if !checkInputAmount(br, shortFillBitWindowRead) ***REMOVED***
			s.symbol = symbol
			s.repeat = repeat
			s.prev_code_len = prev_code_len
			s.repeat_code_len = repeat_code_len
			s.space = space
			return decoderNeedsMoreInput
		***REMOVED***

		fillBitWindow16(br)
		p = p[getBitsUnmasked(br)&uint64(bitMask(huffmanMaxCodeLengthCodeLength)):]
		dropBits(br, uint32(p[0].bits)) /* Use 1..5 bits. */
		code_len = uint32(p[0].value)   /* code_len == 0..17 */
		if code_len < repeatPreviousCodeLength ***REMOVED***
			processSingleCodeLength(code_len, &symbol, &repeat, &space, &prev_code_len, symbol_lists, code_length_histo, next_symbol) /* code_len == 16..17, extra_bits == 2..3 */
		***REMOVED*** else ***REMOVED***
			var extra_bits uint32
			if code_len == repeatPreviousCodeLength ***REMOVED***
				extra_bits = 2
			***REMOVED*** else ***REMOVED***
				extra_bits = 3
			***REMOVED***
			var repeat_delta uint32 = uint32(getBitsUnmasked(br)) & bitMask(extra_bits)
			dropBits(br, extra_bits)
			processRepeatedCodeLength(code_len, repeat_delta, alphabet_size, &symbol, &repeat, &space, &prev_code_len, &repeat_code_len, symbol_lists, code_length_histo, next_symbol)
		***REMOVED***
	***REMOVED***

	s.space = space
	return decoderSuccess
***REMOVED***

func safeReadSymbolCodeLengths(alphabet_size uint32, s *Reader) int ***REMOVED***
	var br *bitReader = &s.br
	var get_byte bool = false
	var p []huffmanCode
	for s.symbol < alphabet_size && s.space > 0 ***REMOVED***
		p = s.table[:]
		var code_len uint32
		var available_bits uint32
		var bits uint32 = 0
		if get_byte && !pullByte(br) ***REMOVED***
			return decoderNeedsMoreInput
		***REMOVED***
		get_byte = false
		available_bits = getAvailableBits(br)
		if available_bits != 0 ***REMOVED***
			bits = uint32(getBitsUnmasked(br))
		***REMOVED***

		p = p[bits&bitMask(huffmanMaxCodeLengthCodeLength):]
		if uint32(p[0].bits) > available_bits ***REMOVED***
			get_byte = true
			continue
		***REMOVED***

		code_len = uint32(p[0].value) /* code_len == 0..17 */
		if code_len < repeatPreviousCodeLength ***REMOVED***
			dropBits(br, uint32(p[0].bits))
			processSingleCodeLength(code_len, &s.symbol, &s.repeat, &s.space, &s.prev_code_len, s.symbol_lists, s.code_length_histo[:], s.next_symbol[:]) /* code_len == 16..17, extra_bits == 2..3 */
		***REMOVED*** else ***REMOVED***
			var extra_bits uint32 = code_len - 14
			var repeat_delta uint32 = (bits >> p[0].bits) & bitMask(extra_bits)
			if available_bits < uint32(p[0].bits)+extra_bits ***REMOVED***
				get_byte = true
				continue
			***REMOVED***

			dropBits(br, uint32(p[0].bits)+extra_bits)
			processRepeatedCodeLength(code_len, repeat_delta, alphabet_size, &s.symbol, &s.repeat, &s.space, &s.prev_code_len, &s.repeat_code_len, s.symbol_lists, s.code_length_histo[:], s.next_symbol[:])
		***REMOVED***
	***REMOVED***

	return decoderSuccess
***REMOVED***

/* Reads and decodes 15..18 codes using static prefix code.
   Each code is 2..4 bits long. In total 30..72 bits are used. */
func readCodeLengthCodeLengths(s *Reader) int ***REMOVED***
	var br *bitReader = &s.br
	var num_codes uint32 = s.repeat
	var space uint32 = s.space
	var i uint32 = s.sub_loop_counter
	for ; i < codeLengthCodes; i++ ***REMOVED***
		var code_len_idx byte = kCodeLengthCodeOrder[i]
		var ix uint32
		var v uint32
		if !safeGetBits(br, 4, &ix) ***REMOVED***
			var available_bits uint32 = getAvailableBits(br)
			if available_bits != 0 ***REMOVED***
				ix = uint32(getBitsUnmasked(br) & 0xF)
			***REMOVED*** else ***REMOVED***
				ix = 0
			***REMOVED***

			if uint32(kCodeLengthPrefixLength[ix]) > available_bits ***REMOVED***
				s.sub_loop_counter = i
				s.repeat = num_codes
				s.space = space
				s.substate_huffman = stateHuffmanComplex
				return decoderNeedsMoreInput
			***REMOVED***
		***REMOVED***

		v = uint32(kCodeLengthPrefixValue[ix])
		dropBits(br, uint32(kCodeLengthPrefixLength[ix]))
		s.code_length_code_lengths[code_len_idx] = byte(v)
		if v != 0 ***REMOVED***
			space = space - (32 >> v)
			num_codes++
			s.code_length_histo[v]++
			if space-1 >= 32 ***REMOVED***
				/* space is 0 or wrapped around. */
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if num_codes != 1 && space != 0 ***REMOVED***
		return decoderErrorFormatClSpace
	***REMOVED***

	return decoderSuccess
***REMOVED***

/* Decodes the Huffman tables.
   There are 2 scenarios:
    A) Huffman code contains only few symbols (1..4). Those symbols are read
       directly; their code lengths are defined by the number of symbols.
       For this scenario 4 - 49 bits will be read.

    B) 2-phase decoding:
    B.1) Small Huffman table is decoded; it is specified with code lengths
         encoded with predefined entropy code. 32 - 74 bits are used.
    B.2) Decoded table is used to decode code lengths of symbols in resulting
         Huffman table. In worst case 3520 bits are read. */
func readHuffmanCode(alphabet_size uint32, max_symbol uint32, table []huffmanCode, opt_table_size *uint32, s *Reader) int ***REMOVED***
	var br *bitReader = &s.br

	/* Unnecessary masking, but might be good for safety. */
	alphabet_size &= 0x7FF

	/* State machine. */
	for ***REMOVED***
		switch s.substate_huffman ***REMOVED***
		case stateHuffmanNone:
			if !safeReadBits(br, 2, &s.sub_loop_counter) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			/* The value is used as follows:
			   1 for simple code;
			   0 for no skipping, 2 skips 2 code lengths, 3 skips 3 code lengths */
			if s.sub_loop_counter != 1 ***REMOVED***
				s.space = 32
				s.repeat = 0 /* num_codes */
				var i int
				for i = 0; i <= huffmanMaxCodeLengthCodeLength; i++ ***REMOVED***
					s.code_length_histo[i] = 0
				***REMOVED***

				for i = 0; i < codeLengthCodes; i++ ***REMOVED***
					s.code_length_code_lengths[i] = 0
				***REMOVED***

				s.substate_huffman = stateHuffmanComplex
				continue
			***REMOVED***
			fallthrough

			/* Read symbols, codes & code lengths directly. */
		case stateHuffmanSimpleSize:
			if !safeReadBits(br, 2, &s.symbol) ***REMOVED*** /* num_symbols */
				s.substate_huffman = stateHuffmanSimpleSize
				return decoderNeedsMoreInput
			***REMOVED***

			s.sub_loop_counter = 0
			fallthrough

		case stateHuffmanSimpleRead:
			***REMOVED***
				var result int = readSimpleHuffmanSymbols(alphabet_size, max_symbol, s)
				if result != decoderSuccess ***REMOVED***
					return result
				***REMOVED***
			***REMOVED***
			fallthrough

		case stateHuffmanSimpleBuild:
			var table_size uint32
			if s.symbol == 3 ***REMOVED***
				var bits uint32
				if !safeReadBits(br, 1, &bits) ***REMOVED***
					s.substate_huffman = stateHuffmanSimpleBuild
					return decoderNeedsMoreInput
				***REMOVED***

				s.symbol += bits
			***REMOVED***

			table_size = buildSimpleHuffmanTable(table, huffmanTableBits, s.symbols_lists_array[:], s.symbol)
			if opt_table_size != nil ***REMOVED***
				*opt_table_size = table_size
			***REMOVED***

			s.substate_huffman = stateHuffmanNone
			return decoderSuccess

			/* Decode Huffman-coded code lengths. */
		case stateHuffmanComplex:
			***REMOVED***
				var i uint32
				var result int = readCodeLengthCodeLengths(s)
				if result != decoderSuccess ***REMOVED***
					return result
				***REMOVED***

				buildCodeLengthsHuffmanTable(s.table[:], s.code_length_code_lengths[:], s.code_length_histo[:])
				for i = 0; i < 16; i++ ***REMOVED***
					s.code_length_histo[i] = 0
				***REMOVED***

				for i = 0; i <= huffmanMaxCodeLength; i++ ***REMOVED***
					s.next_symbol[i] = int(i) - (huffmanMaxCodeLength + 1)
					symbolListPut(s.symbol_lists, s.next_symbol[i], 0xFFFF)
				***REMOVED***

				s.symbol = 0
				s.prev_code_len = initialRepeatedCodeLength
				s.repeat = 0
				s.repeat_code_len = 0
				s.space = 32768
				s.substate_huffman = stateHuffmanLengthSymbols
			***REMOVED***
			fallthrough

		case stateHuffmanLengthSymbols:
			var table_size uint32
			var result int = readSymbolCodeLengths(max_symbol, s)
			if result == decoderNeedsMoreInput ***REMOVED***
				result = safeReadSymbolCodeLengths(max_symbol, s)
			***REMOVED***

			if result != decoderSuccess ***REMOVED***
				return result
			***REMOVED***

			if s.space != 0 ***REMOVED***
				return decoderErrorFormatHuffmanSpace
			***REMOVED***

			table_size = buildHuffmanTable(table, huffmanTableBits, s.symbol_lists, s.code_length_histo[:])
			if opt_table_size != nil ***REMOVED***
				*opt_table_size = table_size
			***REMOVED***

			s.substate_huffman = stateHuffmanNone
			return decoderSuccess

		default:
			return decoderErrorUnreachable
		***REMOVED***
	***REMOVED***
***REMOVED***

/* Decodes a block length by reading 3..39 bits. */
func readBlockLength(table []huffmanCode, br *bitReader) uint32 ***REMOVED***
	var code uint32
	var nbits uint32
	code = readSymbol(table, br)
	nbits = kBlockLengthPrefixCode[code].nbits /* nbits == 2..24 */
	return kBlockLengthPrefixCode[code].offset + readBits(br, nbits)
***REMOVED***

/* WARNING: if state is not BROTLI_STATE_READ_BLOCK_LENGTH_NONE, then
   reading can't be continued with ReadBlockLength. */
func safeReadBlockLength(s *Reader, result *uint32, table []huffmanCode, br *bitReader) bool ***REMOVED***
	var index uint32
	if s.substate_read_block_length == stateReadBlockLengthNone ***REMOVED***
		if !safeReadSymbol(table, br, &index) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		index = s.block_length_index
	***REMOVED***
	***REMOVED***
		var bits uint32 /* nbits == 2..24 */
		var nbits uint32 = kBlockLengthPrefixCode[index].nbits
		if !safeReadBits(br, nbits, &bits) ***REMOVED***
			s.block_length_index = index
			s.substate_read_block_length = stateReadBlockLengthSuffix
			return false
		***REMOVED***

		*result = kBlockLengthPrefixCode[index].offset + bits
		s.substate_read_block_length = stateReadBlockLengthNone
		return true
	***REMOVED***
***REMOVED***

/* Transform:
    1) initialize list L with values 0, 1,... 255
    2) For each input element X:
    2.1) let Y = L[X]
    2.2) remove X-th element from L
    2.3) prepend Y to L
    2.4) append Y to output

   In most cases max(Y) <= 7, so most of L remains intact.
   To reduce the cost of initialization, we reuse L, remember the upper bound
   of Y values, and reinitialize only first elements in L.

   Most of input values are 0 and 1. To reduce number of branches, we replace
   inner for loop with do-while. */
func inverseMoveToFrontTransform(v []byte, v_len uint32, state *Reader) ***REMOVED***
	var mtf [256]byte
	var i int
	for i = 1; i < 256; i++ ***REMOVED***
		mtf[i] = byte(i)
	***REMOVED***
	var mtf_1 byte

	/* Transform the input. */
	for i = 0; uint32(i) < v_len; i++ ***REMOVED***
		var index int = int(v[i])
		var value byte = mtf[index]
		v[i] = value
		mtf_1 = value
		for index >= 1 ***REMOVED***
			index--
			mtf[index+1] = mtf[index]
		***REMOVED***

		mtf[0] = mtf_1
	***REMOVED***
***REMOVED***

/* Decodes a series of Huffman table using ReadHuffmanCode function. */
func huffmanTreeGroupDecode(group *huffmanTreeGroup, s *Reader) int ***REMOVED***
	if s.substate_tree_group != stateTreeGroupLoop ***REMOVED***
		s.next = group.codes
		s.htree_index = 0
		s.substate_tree_group = stateTreeGroupLoop
	***REMOVED***

	for s.htree_index < int(group.num_htrees) ***REMOVED***
		var table_size uint32
		var result int = readHuffmanCode(uint32(group.alphabet_size), uint32(group.max_symbol), s.next, &table_size, s)
		if result != decoderSuccess ***REMOVED***
			return result
		***REMOVED***
		group.htrees[s.htree_index] = s.next
		s.next = s.next[table_size:]
		s.htree_index++
	***REMOVED***

	s.substate_tree_group = stateTreeGroupNone
	return decoderSuccess
***REMOVED***

/* Decodes a context map.
   Decoding is done in 4 phases:
    1) Read auxiliary information (6..16 bits) and allocate memory.
       In case of trivial context map, decoding is finished at this phase.
    2) Decode Huffman table using ReadHuffmanCode function.
       This table will be used for reading context map items.
    3) Read context map items; "0" values could be run-length encoded.
    4) Optionally, apply InverseMoveToFront transform to the resulting map. */
func decodeContextMap(context_map_size uint32, num_htrees *uint32, context_map_arg *[]byte, s *Reader) int ***REMOVED***
	var br *bitReader = &s.br
	var result int = decoderSuccess

	switch int(s.substate_context_map) ***REMOVED***
	case stateContextMapNone:
		result = decodeVarLenUint8(s, br, num_htrees)
		if result != decoderSuccess ***REMOVED***
			return result
		***REMOVED***

		(*num_htrees)++
		s.context_index = 0
		*context_map_arg = make([]byte, uint(context_map_size))
		if *context_map_arg == nil ***REMOVED***
			return decoderErrorAllocContextMap
		***REMOVED***

		if *num_htrees <= 1 ***REMOVED***
			for i := 0; i < int(context_map_size); i++ ***REMOVED***
				(*context_map_arg)[i] = 0
			***REMOVED***
			return decoderSuccess
		***REMOVED***

		s.substate_context_map = stateContextMapReadPrefix
		fallthrough
	/* Fall through. */
	case stateContextMapReadPrefix:
		***REMOVED***
			var bits uint32

			/* In next stage ReadHuffmanCode uses at least 4 bits, so it is safe
			   to peek 4 bits ahead. */
			if !safeGetBits(br, 5, &bits) ***REMOVED***
				return decoderNeedsMoreInput
			***REMOVED***

			if bits&1 != 0 ***REMOVED*** /* Use RLE for zeros. */
				s.max_run_length_prefix = (bits >> 1) + 1
				dropBits(br, 5)
			***REMOVED*** else ***REMOVED***
				s.max_run_length_prefix = 0
				dropBits(br, 1)
			***REMOVED***

			s.substate_context_map = stateContextMapHuffman
		***REMOVED***
		fallthrough

		/* Fall through. */
	case stateContextMapHuffman:
		***REMOVED***
			var alphabet_size uint32 = *num_htrees + s.max_run_length_prefix
			result = readHuffmanCode(alphabet_size, alphabet_size, s.context_map_table[:], nil, s)
			if result != decoderSuccess ***REMOVED***
				return result
			***REMOVED***
			s.code = 0xFFFF
			s.substate_context_map = stateContextMapDecode
		***REMOVED***
		fallthrough

		/* Fall through. */
	case stateContextMapDecode:
		***REMOVED***
			var context_index uint32 = s.context_index
			var max_run_length_prefix uint32 = s.max_run_length_prefix
			var context_map []byte = *context_map_arg
			var code uint32 = s.code
			var skip_preamble bool = (code != 0xFFFF)
			for context_index < context_map_size || skip_preamble ***REMOVED***
				if !skip_preamble ***REMOVED***
					if !safeReadSymbol(s.context_map_table[:], br, &code) ***REMOVED***
						s.code = 0xFFFF
						s.context_index = context_index
						return decoderNeedsMoreInput
					***REMOVED***

					if code == 0 ***REMOVED***
						context_map[context_index] = 0
						context_index++
						continue
					***REMOVED***

					if code > max_run_length_prefix ***REMOVED***
						context_map[context_index] = byte(code - max_run_length_prefix)
						context_index++
						continue
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					skip_preamble = false
				***REMOVED***

				/* RLE sub-stage. */
				***REMOVED***
					var reps uint32
					if !safeReadBits(br, code, &reps) ***REMOVED***
						s.code = code
						s.context_index = context_index
						return decoderNeedsMoreInput
					***REMOVED***

					reps += 1 << code
					if context_index+reps > context_map_size ***REMOVED***
						return decoderErrorFormatContextMapRepeat
					***REMOVED***

					for ***REMOVED***
						context_map[context_index] = 0
						context_index++
						reps--
						if reps == 0 ***REMOVED***
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		fallthrough

	case stateContextMapTransform:
		var bits uint32
		if !safeReadBits(br, 1, &bits) ***REMOVED***
			s.substate_context_map = stateContextMapTransform
			return decoderNeedsMoreInput
		***REMOVED***

		if bits != 0 ***REMOVED***
			inverseMoveToFrontTransform(*context_map_arg, context_map_size, s)
		***REMOVED***

		s.substate_context_map = stateContextMapNone
		return decoderSuccess

	default:
		return decoderErrorUnreachable
	***REMOVED***
***REMOVED***

/* Decodes a command or literal and updates block type ring-buffer.
   Reads 3..54 bits. */
func decodeBlockTypeAndLength(safe int, s *Reader, tree_type int) bool ***REMOVED***
	var max_block_type uint32 = s.num_block_types[tree_type]
	type_tree := s.block_type_trees[tree_type*huffmanMaxSize258:]
	len_tree := s.block_len_trees[tree_type*huffmanMaxSize26:]
	var br *bitReader = &s.br
	var ringbuffer []uint32 = s.block_type_rb[tree_type*2:]
	var block_type uint32
	if max_block_type <= 1 ***REMOVED***
		return false
	***REMOVED***

	/* Read 0..15 + 3..39 bits. */
	if safe == 0 ***REMOVED***
		block_type = readSymbol(type_tree, br)
		s.block_length[tree_type] = readBlockLength(len_tree, br)
	***REMOVED*** else ***REMOVED***
		var memento bitReaderState
		bitReaderSaveState(br, &memento)
		if !safeReadSymbol(type_tree, br, &block_type) ***REMOVED***
			return false
		***REMOVED***
		if !safeReadBlockLength(s, &s.block_length[tree_type], len_tree, br) ***REMOVED***
			s.substate_read_block_length = stateReadBlockLengthNone
			bitReaderRestoreState(br, &memento)
			return false
		***REMOVED***
	***REMOVED***

	if block_type == 1 ***REMOVED***
		block_type = ringbuffer[1] + 1
	***REMOVED*** else if block_type == 0 ***REMOVED***
		block_type = ringbuffer[0]
	***REMOVED*** else ***REMOVED***
		block_type -= 2
	***REMOVED***

	if block_type >= max_block_type ***REMOVED***
		block_type -= max_block_type
	***REMOVED***

	ringbuffer[0] = ringbuffer[1]
	ringbuffer[1] = block_type
	return true
***REMOVED***

func detectTrivialLiteralBlockTypes(s *Reader) ***REMOVED***
	var i uint
	for i = 0; i < 8; i++ ***REMOVED***
		s.trivial_literal_contexts[i] = 0
	***REMOVED***
	for i = 0; uint32(i) < s.num_block_types[0]; i++ ***REMOVED***
		var offset uint = i << literalContextBits
		var error uint = 0
		var sample uint = uint(s.context_map[offset])
		var j uint
		for j = 0; j < 1<<literalContextBits; ***REMOVED***
			var k int
			for k = 0; k < 4; k++ ***REMOVED***
				error |= uint(s.context_map[offset+j]) ^ sample
				j++
			***REMOVED***
		***REMOVED***

		if error == 0 ***REMOVED***
			s.trivial_literal_contexts[i>>5] |= 1 << (i & 31)
		***REMOVED***
	***REMOVED***
***REMOVED***

func prepareLiteralDecoding(s *Reader) ***REMOVED***
	var context_mode byte
	var trivial uint
	var block_type uint32 = s.block_type_rb[1]
	var context_offset uint32 = block_type << literalContextBits
	s.context_map_slice = s.context_map[context_offset:]
	trivial = uint(s.trivial_literal_contexts[block_type>>5])
	s.trivial_literal_context = int((trivial >> (block_type & 31)) & 1)
	s.literal_htree = []huffmanCode(s.literal_hgroup.htrees[s.context_map_slice[0]])
	context_mode = s.context_modes[block_type] & 3
	s.context_lookup = getContextLUT(int(context_mode))
***REMOVED***

/* Decodes the block type and updates the state for literal context.
   Reads 3..54 bits. */
func decodeLiteralBlockSwitchInternal(safe int, s *Reader) bool ***REMOVED***
	if !decodeBlockTypeAndLength(safe, s, 0) ***REMOVED***
		return false
	***REMOVED***

	prepareLiteralDecoding(s)
	return true
***REMOVED***

func decodeLiteralBlockSwitch(s *Reader) ***REMOVED***
	decodeLiteralBlockSwitchInternal(0, s)
***REMOVED***

func safeDecodeLiteralBlockSwitch(s *Reader) bool ***REMOVED***
	return decodeLiteralBlockSwitchInternal(1, s)
***REMOVED***

/* Block switch for insert/copy length.
   Reads 3..54 bits. */
func decodeCommandBlockSwitchInternal(safe int, s *Reader) bool ***REMOVED***
	if !decodeBlockTypeAndLength(safe, s, 1) ***REMOVED***
		return false
	***REMOVED***

	s.htree_command = []huffmanCode(s.insert_copy_hgroup.htrees[s.block_type_rb[3]])
	return true
***REMOVED***

func decodeCommandBlockSwitch(s *Reader) ***REMOVED***
	decodeCommandBlockSwitchInternal(0, s)
***REMOVED***

func safeDecodeCommandBlockSwitch(s *Reader) bool ***REMOVED***
	return decodeCommandBlockSwitchInternal(1, s)
***REMOVED***

/* Block switch for distance codes.
   Reads 3..54 bits. */
func decodeDistanceBlockSwitchInternal(safe int, s *Reader) bool ***REMOVED***
	if !decodeBlockTypeAndLength(safe, s, 2) ***REMOVED***
		return false
	***REMOVED***

	s.dist_context_map_slice = s.dist_context_map[s.block_type_rb[5]<<distanceContextBits:]
	s.dist_htree_index = s.dist_context_map_slice[s.distance_context]
	return true
***REMOVED***

func decodeDistanceBlockSwitch(s *Reader) ***REMOVED***
	decodeDistanceBlockSwitchInternal(0, s)
***REMOVED***

func safeDecodeDistanceBlockSwitch(s *Reader) bool ***REMOVED***
	return decodeDistanceBlockSwitchInternal(1, s)
***REMOVED***

func unwrittenBytes(s *Reader, wrap bool) uint ***REMOVED***
	var pos uint
	if wrap && s.pos > s.ringbuffer_size ***REMOVED***
		pos = uint(s.ringbuffer_size)
	***REMOVED*** else ***REMOVED***
		pos = uint(s.pos)
	***REMOVED***
	var partial_pos_rb uint = (s.rb_roundtrips * uint(s.ringbuffer_size)) + pos
	return partial_pos_rb - s.partial_pos_out
***REMOVED***

/* Dumps output.
   Returns BROTLI_DECODER_NEEDS_MORE_OUTPUT only if there is more output to push
   and either ring-buffer is as big as window size, or |force| is true. */
func writeRingBuffer(s *Reader, available_out *uint, next_out *[]byte, total_out *uint, force bool) int ***REMOVED***
	start := s.ringbuffer[s.partial_pos_out&uint(s.ringbuffer_mask):]
	var to_write uint = unwrittenBytes(s, true)
	var num_written uint = *available_out
	if num_written > to_write ***REMOVED***
		num_written = to_write
	***REMOVED***

	if s.meta_block_remaining_len < 0 ***REMOVED***
		return decoderErrorFormatBlockLength1
	***REMOVED***

	if next_out != nil && *next_out == nil ***REMOVED***
		*next_out = start
	***REMOVED*** else ***REMOVED***
		if next_out != nil ***REMOVED***
			copy(*next_out, start[:num_written])
			*next_out = (*next_out)[num_written:]
		***REMOVED***
	***REMOVED***

	*available_out -= num_written
	s.partial_pos_out += num_written
	if total_out != nil ***REMOVED***
		*total_out = s.partial_pos_out
	***REMOVED***

	if num_written < to_write ***REMOVED***
		if s.ringbuffer_size == 1<<s.window_bits || force ***REMOVED***
			return decoderNeedsMoreOutput
		***REMOVED*** else ***REMOVED***
			return decoderSuccess
		***REMOVED***
	***REMOVED***

	/* Wrap ring buffer only if it has reached its maximal size. */
	if s.ringbuffer_size == 1<<s.window_bits && s.pos >= s.ringbuffer_size ***REMOVED***
		s.pos -= s.ringbuffer_size
		s.rb_roundtrips++
		if uint(s.pos) != 0 ***REMOVED***
			s.should_wrap_ringbuffer = 1
		***REMOVED*** else ***REMOVED***
			s.should_wrap_ringbuffer = 0
		***REMOVED***
	***REMOVED***

	return decoderSuccess
***REMOVED***

func wrapRingBuffer(s *Reader) ***REMOVED***
	if s.should_wrap_ringbuffer != 0 ***REMOVED***
		copy(s.ringbuffer, s.ringbuffer_end[:uint(s.pos)])
		s.should_wrap_ringbuffer = 0
	***REMOVED***
***REMOVED***

/* Allocates ring-buffer.

   s->ringbuffer_size MUST be updated by BrotliCalculateRingBufferSize before
   this function is called.

   Last two bytes of ring-buffer are initialized to 0, so context calculation
   could be done uniformly for the first two and all other positions. */
func ensureRingBuffer(s *Reader) bool ***REMOVED***
	var old_ringbuffer []byte = s.ringbuffer
	if s.ringbuffer_size == s.new_ringbuffer_size ***REMOVED***
		return true
	***REMOVED***

	s.ringbuffer = make([]byte, uint(s.new_ringbuffer_size)+uint(kRingBufferWriteAheadSlack))
	if s.ringbuffer == nil ***REMOVED***
		/* Restore previous value. */
		s.ringbuffer = old_ringbuffer

		return false
	***REMOVED***

	s.ringbuffer[s.new_ringbuffer_size-2] = 0
	s.ringbuffer[s.new_ringbuffer_size-1] = 0

	if !(old_ringbuffer == nil) ***REMOVED***
		copy(s.ringbuffer, old_ringbuffer[:uint(s.pos)])

		old_ringbuffer = nil
	***REMOVED***

	s.ringbuffer_size = s.new_ringbuffer_size
	s.ringbuffer_mask = s.new_ringbuffer_size - 1
	s.ringbuffer_end = s.ringbuffer[s.ringbuffer_size:]

	return true
***REMOVED***

func copyUncompressedBlockToOutput(available_out *uint, next_out *[]byte, total_out *uint, s *Reader) int ***REMOVED***
	/* TODO: avoid allocation for single uncompressed block. */
	if !ensureRingBuffer(s) ***REMOVED***
		return decoderErrorAllocRingBuffer1
	***REMOVED***

	/* State machine */
	for ***REMOVED***
		switch s.substate_uncompressed ***REMOVED***
		case stateUncompressedNone:
			***REMOVED***
				var nbytes int = int(getRemainingBytes(&s.br))
				if nbytes > s.meta_block_remaining_len ***REMOVED***
					nbytes = s.meta_block_remaining_len
				***REMOVED***

				if s.pos+nbytes > s.ringbuffer_size ***REMOVED***
					nbytes = s.ringbuffer_size - s.pos
				***REMOVED***

				/* Copy remaining bytes from s->br.buf_ to ring-buffer. */
				copyBytes(s.ringbuffer[s.pos:], &s.br, uint(nbytes))

				s.pos += nbytes
				s.meta_block_remaining_len -= nbytes
				if s.pos < 1<<s.window_bits ***REMOVED***
					if s.meta_block_remaining_len == 0 ***REMOVED***
						return decoderSuccess
					***REMOVED***

					return decoderNeedsMoreInput
				***REMOVED***

				s.substate_uncompressed = stateUncompressedWrite
			***REMOVED***
			fallthrough

		case stateUncompressedWrite:
			***REMOVED***
				result := writeRingBuffer(s, available_out, next_out, total_out, false)
				if result != decoderSuccess ***REMOVED***
					return result
				***REMOVED***

				if s.ringbuffer_size == 1<<s.window_bits ***REMOVED***
					s.max_distance = s.max_backward_distance
				***REMOVED***

				s.substate_uncompressed = stateUncompressedNone
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

/* Calculates the smallest feasible ring buffer.

   If we know the data size is small, do not allocate more ring buffer
   size than needed to reduce memory usage.

   When this method is called, metablock size and flags MUST be decoded. */
func calculateRingBufferSize(s *Reader) ***REMOVED***
	var window_size int = 1 << s.window_bits
	var new_ringbuffer_size int = window_size
	var min_size int
	/* We need at least 2 bytes of ring buffer size to get the last two
	   bytes for context from there */
	if s.ringbuffer_size != 0 ***REMOVED***
		min_size = s.ringbuffer_size
	***REMOVED*** else ***REMOVED***
		min_size = 1024
	***REMOVED***
	var output_size int

	/* If maximum is already reached, no further extension is retired. */
	if s.ringbuffer_size == window_size ***REMOVED***
		return
	***REMOVED***

	/* Metadata blocks does not touch ring buffer. */
	if s.is_metadata != 0 ***REMOVED***
		return
	***REMOVED***

	if s.ringbuffer == nil ***REMOVED***
		output_size = 0
	***REMOVED*** else ***REMOVED***
		output_size = s.pos
	***REMOVED***

	output_size += s.meta_block_remaining_len
	if min_size < output_size ***REMOVED***
		min_size = output_size
	***REMOVED***

	if !(s.canny_ringbuffer_allocation == 0) ***REMOVED***
		/* Reduce ring buffer size to save memory when server is unscrupulous.
		   In worst case memory usage might be 1.5x bigger for a short period of
		   ring buffer reallocation. */
		for new_ringbuffer_size>>1 >= min_size ***REMOVED***
			new_ringbuffer_size >>= 1
		***REMOVED***
	***REMOVED***

	s.new_ringbuffer_size = new_ringbuffer_size
***REMOVED***

/* Reads 1..256 2-bit context modes. */
func readContextModes(s *Reader) int ***REMOVED***
	var br *bitReader = &s.br
	var i int = s.loop_counter

	for i < int(s.num_block_types[0]) ***REMOVED***
		var bits uint32
		if !safeReadBits(br, 2, &bits) ***REMOVED***
			s.loop_counter = i
			return decoderNeedsMoreInput
		***REMOVED***

		s.context_modes[i] = byte(bits)
		i++
	***REMOVED***

	return decoderSuccess
***REMOVED***

func takeDistanceFromRingBuffer(s *Reader) ***REMOVED***
	if s.distance_code == 0 ***REMOVED***
		s.dist_rb_idx--
		s.distance_code = s.dist_rb[s.dist_rb_idx&3]

		/* Compensate double distance-ring-buffer roll for dictionary items. */
		s.distance_context = 1
	***REMOVED*** else ***REMOVED***
		var distance_code int = s.distance_code << 1
		const kDistanceShortCodeIndexOffset uint32 = 0xAAAFFF1B
		const kDistanceShortCodeValueOffset uint32 = 0xFA5FA500
		var v int = (s.dist_rb_idx + int(kDistanceShortCodeIndexOffset>>uint(distance_code))) & 0x3
		/* kDistanceShortCodeIndexOffset has 2-bit values from LSB:
		   3, 2, 1, 0, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2 */

		/* kDistanceShortCodeValueOffset has 2-bit values from LSB:
		   -0, 0,-0, 0,-1, 1,-2, 2,-3, 3,-1, 1,-2, 2,-3, 3 */
		s.distance_code = s.dist_rb[v]

		v = int(kDistanceShortCodeValueOffset>>uint(distance_code)) & 0x3
		if distance_code&0x3 != 0 ***REMOVED***
			s.distance_code += v
		***REMOVED*** else ***REMOVED***
			s.distance_code -= v
			if s.distance_code <= 0 ***REMOVED***
				/* A huge distance will cause a () soon.
				   This is a little faster than failing here. */
				s.distance_code = 0x7FFFFFFF
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func safeReadBitsMaybeZero(br *bitReader, n_bits uint32, val *uint32) bool ***REMOVED***
	if n_bits != 0 ***REMOVED***
		return safeReadBits(br, n_bits, val)
	***REMOVED*** else ***REMOVED***
		*val = 0
		return true
	***REMOVED***
***REMOVED***

/* Precondition: s->distance_code < 0. */
func readDistanceInternal(safe int, s *Reader, br *bitReader) bool ***REMOVED***
	var distval int
	var memento bitReaderState
	var distance_tree []huffmanCode = []huffmanCode(s.distance_hgroup.htrees[s.dist_htree_index])
	if safe == 0 ***REMOVED***
		s.distance_code = int(readSymbol(distance_tree, br))
	***REMOVED*** else ***REMOVED***
		var code uint32
		bitReaderSaveState(br, &memento)
		if !safeReadSymbol(distance_tree, br, &code) ***REMOVED***
			return false
		***REMOVED***

		s.distance_code = int(code)
	***REMOVED***

	/* Convert the distance code to the actual distance by possibly
	   looking up past distances from the s->ringbuffer. */
	s.distance_context = 0

	if s.distance_code&^0xF == 0 ***REMOVED***
		takeDistanceFromRingBuffer(s)
		s.block_length[2]--
		return true
	***REMOVED***

	distval = s.distance_code - int(s.num_direct_distance_codes)
	if distval >= 0 ***REMOVED***
		var nbits uint32
		var postfix int
		var offset int
		if safe == 0 && (s.distance_postfix_bits == 0) ***REMOVED***
			nbits = (uint32(distval) >> 1) + 1
			offset = ((2 + (distval & 1)) << nbits) - 4
			s.distance_code = int(s.num_direct_distance_codes) + offset + int(readBits(br, nbits))
		***REMOVED*** else ***REMOVED***
			/* This branch also works well when s->distance_postfix_bits == 0. */
			var bits uint32
			postfix = distval & s.distance_postfix_mask
			distval >>= s.distance_postfix_bits
			nbits = (uint32(distval) >> 1) + 1
			if safe != 0 ***REMOVED***
				if !safeReadBitsMaybeZero(br, nbits, &bits) ***REMOVED***
					s.distance_code = -1 /* Restore precondition. */
					bitReaderRestoreState(br, &memento)
					return false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				bits = readBits(br, nbits)
			***REMOVED***

			offset = ((2 + (distval & 1)) << nbits) - 4
			s.distance_code = int(s.num_direct_distance_codes) + ((offset + int(bits)) << s.distance_postfix_bits) + postfix
		***REMOVED***
	***REMOVED***

	s.distance_code = s.distance_code - numDistanceShortCodes + 1
	s.block_length[2]--
	return true
***REMOVED***

func readDistance(s *Reader, br *bitReader) ***REMOVED***
	readDistanceInternal(0, s, br)
***REMOVED***

func safeReadDistance(s *Reader, br *bitReader) bool ***REMOVED***
	return readDistanceInternal(1, s, br)
***REMOVED***

func readCommandInternal(safe int, s *Reader, br *bitReader, insert_length *int) bool ***REMOVED***
	var cmd_code uint32
	var insert_len_extra uint32 = 0
	var copy_length uint32
	var v cmdLutElement
	var memento bitReaderState
	if safe == 0 ***REMOVED***
		cmd_code = readSymbol(s.htree_command, br)
	***REMOVED*** else ***REMOVED***
		bitReaderSaveState(br, &memento)
		if !safeReadSymbol(s.htree_command, br, &cmd_code) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	v = kCmdLut[cmd_code]
	s.distance_code = int(v.distance_code)
	s.distance_context = int(v.context)
	s.dist_htree_index = s.dist_context_map_slice[s.distance_context]
	*insert_length = int(v.insert_len_offset)
	if safe == 0 ***REMOVED***
		if v.insert_len_extra_bits != 0 ***REMOVED***
			insert_len_extra = readBits(br, uint32(v.insert_len_extra_bits))
		***REMOVED***

		copy_length = readBits(br, uint32(v.copy_len_extra_bits))
	***REMOVED*** else ***REMOVED***
		if !safeReadBitsMaybeZero(br, uint32(v.insert_len_extra_bits), &insert_len_extra) || !safeReadBitsMaybeZero(br, uint32(v.copy_len_extra_bits), &copy_length) ***REMOVED***
			bitReaderRestoreState(br, &memento)
			return false
		***REMOVED***
	***REMOVED***

	s.copy_length = int(copy_length) + int(v.copy_len_offset)
	s.block_length[1]--
	*insert_length += int(insert_len_extra)
	return true
***REMOVED***

func readCommand(s *Reader, br *bitReader, insert_length *int) ***REMOVED***
	readCommandInternal(0, s, br, insert_length)
***REMOVED***

func safeReadCommand(s *Reader, br *bitReader, insert_length *int) bool ***REMOVED***
	return readCommandInternal(1, s, br, insert_length)
***REMOVED***

func checkInputAmountMaybeSafe(safe int, br *bitReader, num uint) bool ***REMOVED***
	if safe != 0 ***REMOVED***
		return true
	***REMOVED***

	return checkInputAmount(br, num)
***REMOVED***

func processCommandsInternal(safe int, s *Reader) int ***REMOVED***
	var pos int = s.pos
	var i int = s.loop_counter
	var result int = decoderSuccess
	var br *bitReader = &s.br
	var hc []huffmanCode

	if !checkInputAmountMaybeSafe(safe, br, 28) ***REMOVED***
		result = decoderNeedsMoreInput
		goto saveStateAndReturn
	***REMOVED***

	if safe == 0 ***REMOVED***
		warmupBitReader(br)
	***REMOVED***

	/* Jump into state machine. */
	if s.state == stateCommandBegin ***REMOVED***
		goto CommandBegin
	***REMOVED*** else if s.state == stateCommandInner ***REMOVED***
		goto CommandInner
	***REMOVED*** else if s.state == stateCommandPostDecodeLiterals ***REMOVED***
		goto CommandPostDecodeLiterals
	***REMOVED*** else if s.state == stateCommandPostWrapCopy ***REMOVED***
		goto CommandPostWrapCopy
	***REMOVED*** else ***REMOVED***
		return decoderErrorUnreachable
	***REMOVED***

CommandBegin:
	if safe != 0 ***REMOVED***
		s.state = stateCommandBegin
	***REMOVED***

	if !checkInputAmountMaybeSafe(safe, br, 28) ***REMOVED*** /* 156 bits + 7 bytes */
		s.state = stateCommandBegin
		result = decoderNeedsMoreInput
		goto saveStateAndReturn
	***REMOVED***

	if s.block_length[1] == 0 ***REMOVED***
		if safe != 0 ***REMOVED***
			if !safeDecodeCommandBlockSwitch(s) ***REMOVED***
				result = decoderNeedsMoreInput
				goto saveStateAndReturn
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			decodeCommandBlockSwitch(s)
		***REMOVED***

		goto CommandBegin
	***REMOVED***

	/* Read the insert/copy length in the command. */
	if safe != 0 ***REMOVED***
		if !safeReadCommand(s, br, &i) ***REMOVED***
			result = decoderNeedsMoreInput
			goto saveStateAndReturn
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		readCommand(s, br, &i)
	***REMOVED***

	if i == 0 ***REMOVED***
		goto CommandPostDecodeLiterals
	***REMOVED***

	s.meta_block_remaining_len -= i

CommandInner:
	if safe != 0 ***REMOVED***
		s.state = stateCommandInner
	***REMOVED***

	/* Read the literals in the command. */
	if s.trivial_literal_context != 0 ***REMOVED***
		var bits uint32
		var value uint32
		preloadSymbol(safe, s.literal_htree, br, &bits, &value)
		for ***REMOVED***
			if !checkInputAmountMaybeSafe(safe, br, 28) ***REMOVED*** /* 162 bits + 7 bytes */
				s.state = stateCommandInner
				result = decoderNeedsMoreInput
				goto saveStateAndReturn
			***REMOVED***

			if s.block_length[0] == 0 ***REMOVED***
				if safe != 0 ***REMOVED***
					if !safeDecodeLiteralBlockSwitch(s) ***REMOVED***
						result = decoderNeedsMoreInput
						goto saveStateAndReturn
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					decodeLiteralBlockSwitch(s)
				***REMOVED***

				preloadSymbol(safe, s.literal_htree, br, &bits, &value)
				if s.trivial_literal_context == 0 ***REMOVED***
					goto CommandInner
				***REMOVED***
			***REMOVED***

			if safe == 0 ***REMOVED***
				s.ringbuffer[pos] = byte(readPreloadedSymbol(s.literal_htree, br, &bits, &value))
			***REMOVED*** else ***REMOVED***
				var literal uint32
				if !safeReadSymbol(s.literal_htree, br, &literal) ***REMOVED***
					result = decoderNeedsMoreInput
					goto saveStateAndReturn
				***REMOVED***

				s.ringbuffer[pos] = byte(literal)
			***REMOVED***

			s.block_length[0]--
			pos++
			if pos == s.ringbuffer_size ***REMOVED***
				s.state = stateCommandInnerWrite
				i--
				goto saveStateAndReturn
			***REMOVED***
			i--
			if i == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var p1 byte = s.ringbuffer[(pos-1)&s.ringbuffer_mask]
		var p2 byte = s.ringbuffer[(pos-2)&s.ringbuffer_mask]
		for ***REMOVED***
			var context byte
			if !checkInputAmountMaybeSafe(safe, br, 28) ***REMOVED*** /* 162 bits + 7 bytes */
				s.state = stateCommandInner
				result = decoderNeedsMoreInput
				goto saveStateAndReturn
			***REMOVED***

			if s.block_length[0] == 0 ***REMOVED***
				if safe != 0 ***REMOVED***
					if !safeDecodeLiteralBlockSwitch(s) ***REMOVED***
						result = decoderNeedsMoreInput
						goto saveStateAndReturn
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					decodeLiteralBlockSwitch(s)
				***REMOVED***

				if s.trivial_literal_context != 0 ***REMOVED***
					goto CommandInner
				***REMOVED***
			***REMOVED***

			context = getContext(p1, p2, s.context_lookup)
			hc = []huffmanCode(s.literal_hgroup.htrees[s.context_map_slice[context]])
			p2 = p1
			if safe == 0 ***REMOVED***
				p1 = byte(readSymbol(hc, br))
			***REMOVED*** else ***REMOVED***
				var literal uint32
				if !safeReadSymbol(hc, br, &literal) ***REMOVED***
					result = decoderNeedsMoreInput
					goto saveStateAndReturn
				***REMOVED***

				p1 = byte(literal)
			***REMOVED***

			s.ringbuffer[pos] = p1
			s.block_length[0]--
			pos++
			if pos == s.ringbuffer_size ***REMOVED***
				s.state = stateCommandInnerWrite
				i--
				goto saveStateAndReturn
			***REMOVED***
			i--
			if i == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if s.meta_block_remaining_len <= 0 ***REMOVED***
		s.state = stateMetablockDone
		goto saveStateAndReturn
	***REMOVED***

CommandPostDecodeLiterals:
	if safe != 0 ***REMOVED***
		s.state = stateCommandPostDecodeLiterals
	***REMOVED***

	if s.distance_code >= 0 ***REMOVED***
		/* Implicit distance case. */
		if s.distance_code != 0 ***REMOVED***
			s.distance_context = 0
		***REMOVED*** else ***REMOVED***
			s.distance_context = 1
		***REMOVED***

		s.dist_rb_idx--
		s.distance_code = s.dist_rb[s.dist_rb_idx&3]
	***REMOVED*** else ***REMOVED***
		/* Read distance code in the command, unless it was implicitly zero. */
		if s.block_length[2] == 0 ***REMOVED***
			if safe != 0 ***REMOVED***
				if !safeDecodeDistanceBlockSwitch(s) ***REMOVED***
					result = decoderNeedsMoreInput
					goto saveStateAndReturn
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				decodeDistanceBlockSwitch(s)
			***REMOVED***
		***REMOVED***

		if safe != 0 ***REMOVED***
			if !safeReadDistance(s, br) ***REMOVED***
				result = decoderNeedsMoreInput
				goto saveStateAndReturn
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			readDistance(s, br)
		***REMOVED***
	***REMOVED***

	if s.max_distance != s.max_backward_distance ***REMOVED***
		if pos < s.max_backward_distance ***REMOVED***
			s.max_distance = pos
		***REMOVED*** else ***REMOVED***
			s.max_distance = s.max_backward_distance
		***REMOVED***
	***REMOVED***

	i = s.copy_length

	/* Apply copy of LZ77 back-reference, or static dictionary reference if
	   the distance is larger than the max LZ77 distance */
	if s.distance_code > s.max_distance ***REMOVED***
		/* The maximum allowed distance is BROTLI_MAX_ALLOWED_DISTANCE = 0x7FFFFFFC.
		   With this choice, no signed overflow can occur after decoding
		   a special distance code (e.g., after adding 3 to the last distance). */
		if s.distance_code > maxAllowedDistance ***REMOVED***
			return decoderErrorFormatDistance
		***REMOVED***

		if i >= minDictionaryWordLength && i <= maxDictionaryWordLength ***REMOVED***
			var address int = s.distance_code - s.max_distance - 1
			var words *dictionary = s.dictionary
			var trans *transforms = s.transforms
			var offset int = int(s.dictionary.offsets_by_length[i])
			var shift uint32 = uint32(s.dictionary.size_bits_by_length[i])
			var mask int = int(bitMask(shift))
			var word_idx int = address & mask
			var transform_idx int = address >> shift

			/* Compensate double distance-ring-buffer roll. */
			s.dist_rb_idx += s.distance_context

			offset += word_idx * i
			if words.data == nil ***REMOVED***
				return decoderErrorDictionaryNotSet
			***REMOVED***

			if transform_idx < int(trans.num_transforms) ***REMOVED***
				word := words.data[offset:]
				var len int = i
				if transform_idx == int(trans.cutOffTransforms[0]) ***REMOVED***
					copy(s.ringbuffer[pos:], word[:uint(len)])
				***REMOVED*** else ***REMOVED***
					len = transformDictionaryWord(s.ringbuffer[pos:], word, int(len), trans, transform_idx)
				***REMOVED***

				pos += int(len)
				s.meta_block_remaining_len -= int(len)
				if pos >= s.ringbuffer_size ***REMOVED***
					s.state = stateCommandPostWrite1
					goto saveStateAndReturn
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return decoderErrorFormatTransform
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return decoderErrorFormatDictionary
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var src_start int = (pos - s.distance_code) & s.ringbuffer_mask
		copy_dst := s.ringbuffer[pos:]
		copy_src := s.ringbuffer[src_start:]
		var dst_end int = pos + i
		var src_end int = src_start + i

		/* Update the recent distances cache. */
		s.dist_rb[s.dist_rb_idx&3] = s.distance_code

		s.dist_rb_idx++
		s.meta_block_remaining_len -= i

		/* There are 32+ bytes of slack in the ring-buffer allocation.
		   Also, we have 16 short codes, that make these 16 bytes irrelevant
		   in the ring-buffer. Let's copy over them as a first guess. */
		copy(copy_dst, copy_src[:16])

		if src_end > pos && dst_end > src_start ***REMOVED***
			/* Regions intersect. */
			goto CommandPostWrapCopy
		***REMOVED***

		if dst_end >= s.ringbuffer_size || src_end >= s.ringbuffer_size ***REMOVED***
			/* At least one region wraps. */
			goto CommandPostWrapCopy
		***REMOVED***

		pos += i
		if i > 16 ***REMOVED***
			if i > 32 ***REMOVED***
				copy(copy_dst[16:], copy_src[16:][:uint(i-16)])
			***REMOVED*** else ***REMOVED***
				/* This branch covers about 45% cases.
				   Fixed size short copy allows more compiler optimizations. */
				copy(copy_dst[16:], copy_src[16:][:16])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if s.meta_block_remaining_len <= 0 ***REMOVED***
		/* Next metablock, if any. */
		s.state = stateMetablockDone

		goto saveStateAndReturn
	***REMOVED*** else ***REMOVED***
		goto CommandBegin
	***REMOVED***
CommandPostWrapCopy:
	***REMOVED***
		var wrap_guard int = s.ringbuffer_size - pos
		for ***REMOVED***
			i--
			if i < 0 ***REMOVED***
				break
			***REMOVED***
			s.ringbuffer[pos] = s.ringbuffer[(pos-s.distance_code)&s.ringbuffer_mask]
			pos++
			wrap_guard--
			if wrap_guard == 0 ***REMOVED***
				s.state = stateCommandPostWrite2
				goto saveStateAndReturn
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if s.meta_block_remaining_len <= 0 ***REMOVED***
		/* Next metablock, if any. */
		s.state = stateMetablockDone

		goto saveStateAndReturn
	***REMOVED*** else ***REMOVED***
		goto CommandBegin
	***REMOVED***

saveStateAndReturn:
	s.pos = pos
	s.loop_counter = i
	return result
***REMOVED***

func processCommands(s *Reader) int ***REMOVED***
	return processCommandsInternal(0, s)
***REMOVED***

func safeProcessCommands(s *Reader) int ***REMOVED***
	return processCommandsInternal(1, s)
***REMOVED***

/* Returns the maximum number of distance symbols which can only represent
   distances not exceeding BROTLI_MAX_ALLOWED_DISTANCE. */

var maxDistanceSymbol_bound = [maxNpostfix + 1]uint32***REMOVED***0, 4, 12, 28***REMOVED***
var maxDistanceSymbol_diff = [maxNpostfix + 1]uint32***REMOVED***73, 126, 228, 424***REMOVED***

func maxDistanceSymbol(ndirect uint32, npostfix uint32) uint32 ***REMOVED***
	var postfix uint32 = 1 << npostfix
	if ndirect < maxDistanceSymbol_bound[npostfix] ***REMOVED***
		return ndirect + maxDistanceSymbol_diff[npostfix] + postfix
	***REMOVED*** else if ndirect > maxDistanceSymbol_bound[npostfix]+postfix ***REMOVED***
		return ndirect + maxDistanceSymbol_diff[npostfix]
	***REMOVED*** else ***REMOVED***
		return maxDistanceSymbol_bound[npostfix] + maxDistanceSymbol_diff[npostfix] + postfix
	***REMOVED***
***REMOVED***

/* Invariant: input stream is never overconsumed:
   - invalid input implies that the whole stream is invalid -> any amount of
     input could be read and discarded
   - when result is "needs more input", then at least one more byte is REQUIRED
     to complete decoding; all input data MUST be consumed by decoder, so
     client could swap the input buffer
   - when result is "needs more output" decoder MUST ensure that it doesn't
     hold more than 7 bits in bit reader; this saves client from swapping input
     buffer ahead of time
   - when result is "success" decoder MUST return all unused data back to input
     buffer; this is possible because the invariant is held on enter */
func decoderDecompressStream(s *Reader, available_in *uint, next_in *[]byte, available_out *uint, next_out *[]byte) int ***REMOVED***
	var result int = decoderSuccess
	var br *bitReader = &s.br

	/* Do not try to process further in a case of unrecoverable error. */
	if int(s.error_code) < 0 ***REMOVED***
		return decoderResultError
	***REMOVED***

	if *available_out != 0 && (next_out == nil || *next_out == nil) ***REMOVED***
		return saveErrorCode(s, decoderErrorInvalidArguments)
	***REMOVED***

	if *available_out == 0 ***REMOVED***
		next_out = nil
	***REMOVED***
	if s.buffer_length == 0 ***REMOVED*** /* Just connect bit reader to input stream. */
		br.input_len = *available_in
		br.input = *next_in
		br.byte_pos = 0
	***REMOVED*** else ***REMOVED***
		/* At least one byte of input is required. More than one byte of input may
		   be required to complete the transaction -> reading more data must be
		   done in a loop -> do it in a main loop. */
		result = decoderNeedsMoreInput

		br.input = s.buffer.u8[:]
		br.byte_pos = 0
	***REMOVED***

	/* State machine */
	for ***REMOVED***
		if result != decoderSuccess ***REMOVED***
			/* Error, needs more input/output. */
			if result == decoderNeedsMoreInput ***REMOVED***
				if s.ringbuffer != nil ***REMOVED*** /* Pro-actively push output. */
					var intermediate_result int = writeRingBuffer(s, available_out, next_out, nil, true)

					/* WriteRingBuffer checks s->meta_block_remaining_len validity. */
					if int(intermediate_result) < 0 ***REMOVED***
						result = intermediate_result
						break
					***REMOVED***
				***REMOVED***

				if s.buffer_length != 0 ***REMOVED*** /* Used with internal buffer. */
					if br.byte_pos == br.input_len ***REMOVED***
						/* Successfully finished read transaction.
						   Accumulator contains less than 8 bits, because internal buffer
						   is expanded byte-by-byte until it is enough to complete read. */
						s.buffer_length = 0

						/* Switch to input stream and restart. */
						result = decoderSuccess

						br.input_len = *available_in
						br.input = *next_in
						br.byte_pos = 0
						continue
					***REMOVED*** else if *available_in != 0 ***REMOVED***
						/* Not enough data in buffer, but can take one more byte from
						   input stream. */
						result = decoderSuccess

						s.buffer.u8[s.buffer_length] = (*next_in)[0]
						s.buffer_length++
						br.input_len = uint(s.buffer_length)
						*next_in = (*next_in)[1:]
						(*available_in)--

						/* Retry with more data in buffer. */
						continue
					***REMOVED***

					/* Can't finish reading and no more input. */
					break
					/* Input stream doesn't contain enough input. */
				***REMOVED*** else ***REMOVED***
					/* Copy tail to internal buffer and return. */
					*next_in = br.input[br.byte_pos:]

					*available_in = br.input_len - br.byte_pos
					for *available_in != 0 ***REMOVED***
						s.buffer.u8[s.buffer_length] = (*next_in)[0]
						s.buffer_length++
						*next_in = (*next_in)[1:]
						(*available_in)--
					***REMOVED***

					break
				***REMOVED***
			***REMOVED***

			/* Unreachable. */

			/* Fail or needs more output. */
			if s.buffer_length != 0 ***REMOVED***
				/* Just consumed the buffered input and produced some output. Otherwise
				   it would result in "needs more input". Reset internal buffer. */
				s.buffer_length = 0
			***REMOVED*** else ***REMOVED***
				/* Using input stream in last iteration. When decoder switches to input
				   stream it has less than 8 bits in accumulator, so it is safe to
				   return unused accumulator bits there. */
				bitReaderUnload(br)

				*available_in = br.input_len - br.byte_pos
				*next_in = br.input[br.byte_pos:]
			***REMOVED***

			break
		***REMOVED***

		switch s.state ***REMOVED***
		/* Prepare to the first read. */
		case stateUninited:
			if !warmupBitReader(br) ***REMOVED***
				result = decoderNeedsMoreInput
				break
			***REMOVED***

			/* Decode window size. */
			result = decodeWindowBits(s, br) /* Reads 1..8 bits. */
			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***

			if s.large_window ***REMOVED***
				s.state = stateLargeWindowBits
				break
			***REMOVED***

			s.state = stateInitialize

		case stateLargeWindowBits:
			if !safeReadBits(br, 6, &s.window_bits) ***REMOVED***
				result = decoderNeedsMoreInput
				break
			***REMOVED***

			if s.window_bits < largeMinWbits || s.window_bits > largeMaxWbits ***REMOVED***
				result = decoderErrorFormatWindowBits
				break
			***REMOVED***

			s.state = stateInitialize
			fallthrough

			/* Maximum distance, see section 9.1. of the spec. */
		/* Fall through. */
		case stateInitialize:
			s.max_backward_distance = (1 << s.window_bits) - windowGap

			/* Allocate memory for both block_type_trees and block_len_trees. */
			s.block_type_trees = make([]huffmanCode, (3 * (huffmanMaxSize258 + huffmanMaxSize26)))

			if s.block_type_trees == nil ***REMOVED***
				result = decoderErrorAllocBlockTypeTrees
				break
			***REMOVED***

			s.block_len_trees = s.block_type_trees[3*huffmanMaxSize258:]

			s.state = stateMetablockBegin
			fallthrough

			/* Fall through. */
		case stateMetablockBegin:
			decoderStateMetablockBegin(s)

			s.state = stateMetablockHeader
			fallthrough

			/* Fall through. */
		case stateMetablockHeader:
			result = decodeMetaBlockLength(s, br)
			/* Reads 2 - 31 bits. */
			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***

			if s.is_metadata != 0 || s.is_uncompressed != 0 ***REMOVED***
				if !bitReaderJumpToByteBoundary(br) ***REMOVED***
					result = decoderErrorFormatPadding1
					break
				***REMOVED***
			***REMOVED***

			if s.is_metadata != 0 ***REMOVED***
				s.state = stateMetadata
				break
			***REMOVED***

			if s.meta_block_remaining_len == 0 ***REMOVED***
				s.state = stateMetablockDone
				break
			***REMOVED***

			calculateRingBufferSize(s)
			if s.is_uncompressed != 0 ***REMOVED***
				s.state = stateUncompressed
				break
			***REMOVED***

			s.loop_counter = 0
			s.state = stateHuffmanCode0

		case stateUncompressed:
			result = copyUncompressedBlockToOutput(available_out, next_out, nil, s)
			if result == decoderSuccess ***REMOVED***
				s.state = stateMetablockDone
			***REMOVED***

		case stateMetadata:
			for ; s.meta_block_remaining_len > 0; s.meta_block_remaining_len-- ***REMOVED***
				var bits uint32

				/* Read one byte and ignore it. */
				if !safeReadBits(br, 8, &bits) ***REMOVED***
					result = decoderNeedsMoreInput
					break
				***REMOVED***
			***REMOVED***

			if result == decoderSuccess ***REMOVED***
				s.state = stateMetablockDone
			***REMOVED***

		case stateHuffmanCode0:
			if s.loop_counter >= 3 ***REMOVED***
				s.state = stateMetablockHeader2
				break
			***REMOVED***

			/* Reads 1..11 bits. */
			result = decodeVarLenUint8(s, br, &s.num_block_types[s.loop_counter])

			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***

			s.num_block_types[s.loop_counter]++
			if s.num_block_types[s.loop_counter] < 2 ***REMOVED***
				s.loop_counter++
				break
			***REMOVED***

			s.state = stateHuffmanCode1
			fallthrough

		case stateHuffmanCode1:
			***REMOVED***
				var alphabet_size uint32 = s.num_block_types[s.loop_counter] + 2
				var tree_offset int = s.loop_counter * huffmanMaxSize258
				result = readHuffmanCode(alphabet_size, alphabet_size, s.block_type_trees[tree_offset:], nil, s)
				if result != decoderSuccess ***REMOVED***
					break
				***REMOVED***
				s.state = stateHuffmanCode2
			***REMOVED***
			fallthrough

		case stateHuffmanCode2:
			***REMOVED***
				var alphabet_size uint32 = numBlockLenSymbols
				var tree_offset int = s.loop_counter * huffmanMaxSize26
				result = readHuffmanCode(alphabet_size, alphabet_size, s.block_len_trees[tree_offset:], nil, s)
				if result != decoderSuccess ***REMOVED***
					break
				***REMOVED***
				s.state = stateHuffmanCode3
			***REMOVED***
			fallthrough

		case stateHuffmanCode3:
			var tree_offset int = s.loop_counter * huffmanMaxSize26
			if !safeReadBlockLength(s, &s.block_length[s.loop_counter], s.block_len_trees[tree_offset:], br) ***REMOVED***
				result = decoderNeedsMoreInput
				break
			***REMOVED***

			s.loop_counter++
			s.state = stateHuffmanCode0

		case stateMetablockHeader2:
			***REMOVED***
				var bits uint32
				if !safeReadBits(br, 6, &bits) ***REMOVED***
					result = decoderNeedsMoreInput
					break
				***REMOVED***

				s.distance_postfix_bits = bits & bitMask(2)
				bits >>= 2
				s.num_direct_distance_codes = numDistanceShortCodes + (bits << s.distance_postfix_bits)
				s.distance_postfix_mask = int(bitMask(s.distance_postfix_bits))
				s.context_modes = make([]byte, uint(s.num_block_types[0]))
				if s.context_modes == nil ***REMOVED***
					result = decoderErrorAllocContextModes
					break
				***REMOVED***

				s.loop_counter = 0
				s.state = stateContextModes
			***REMOVED***
			fallthrough

		case stateContextModes:
			result = readContextModes(s)

			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***

			s.state = stateContextMap1
			fallthrough

		case stateContextMap1:
			result = decodeContextMap(s.num_block_types[0]<<literalContextBits, &s.num_literal_htrees, &s.context_map, s)

			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***

			detectTrivialLiteralBlockTypes(s)
			s.state = stateContextMap2
			fallthrough

		case stateContextMap2:
			***REMOVED***
				var num_direct_codes uint32 = s.num_direct_distance_codes - numDistanceShortCodes
				var num_distance_codes uint32
				var max_distance_symbol uint32
				if s.large_window ***REMOVED***
					num_distance_codes = uint32(distanceAlphabetSize(uint(s.distance_postfix_bits), uint(num_direct_codes), largeMaxDistanceBits))
					max_distance_symbol = maxDistanceSymbol(num_direct_codes, s.distance_postfix_bits)
				***REMOVED*** else ***REMOVED***
					num_distance_codes = uint32(distanceAlphabetSize(uint(s.distance_postfix_bits), uint(num_direct_codes), maxDistanceBits))
					max_distance_symbol = num_distance_codes
				***REMOVED***
				var allocation_success bool = true
				result = decodeContextMap(s.num_block_types[2]<<distanceContextBits, &s.num_dist_htrees, &s.dist_context_map, s)
				if result != decoderSuccess ***REMOVED***
					break
				***REMOVED***

				if !decoderHuffmanTreeGroupInit(s, &s.literal_hgroup, numLiteralSymbols, numLiteralSymbols, s.num_literal_htrees) ***REMOVED***
					allocation_success = false
				***REMOVED***

				if !decoderHuffmanTreeGroupInit(s, &s.insert_copy_hgroup, numCommandSymbols, numCommandSymbols, s.num_block_types[1]) ***REMOVED***
					allocation_success = false
				***REMOVED***

				if !decoderHuffmanTreeGroupInit(s, &s.distance_hgroup, num_distance_codes, max_distance_symbol, s.num_dist_htrees) ***REMOVED***
					allocation_success = false
				***REMOVED***

				if !allocation_success ***REMOVED***
					return saveErrorCode(s, decoderErrorAllocTreeGroups)
				***REMOVED***

				s.loop_counter = 0
				s.state = stateTreeGroup
			***REMOVED***
			fallthrough

		case stateTreeGroup:
			var hgroup *huffmanTreeGroup = nil
			switch s.loop_counter ***REMOVED***
			case 0:
				hgroup = &s.literal_hgroup
			case 1:
				hgroup = &s.insert_copy_hgroup
			case 2:
				hgroup = &s.distance_hgroup
			default:
				return saveErrorCode(s, decoderErrorUnreachable)
			***REMOVED***

			result = huffmanTreeGroupDecode(hgroup, s)
			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***
			s.loop_counter++
			if s.loop_counter >= 3 ***REMOVED***
				prepareLiteralDecoding(s)
				s.dist_context_map_slice = s.dist_context_map
				s.htree_command = []huffmanCode(s.insert_copy_hgroup.htrees[0])
				if !ensureRingBuffer(s) ***REMOVED***
					result = decoderErrorAllocRingBuffer2
					break
				***REMOVED***

				s.state = stateCommandBegin
			***REMOVED***

		case stateCommandBegin, stateCommandInner, stateCommandPostDecodeLiterals, stateCommandPostWrapCopy:
			result = processCommands(s)

			if result == decoderNeedsMoreInput ***REMOVED***
				result = safeProcessCommands(s)
			***REMOVED***

		case stateCommandInnerWrite, stateCommandPostWrite1, stateCommandPostWrite2:
			result = writeRingBuffer(s, available_out, next_out, nil, false)

			if result != decoderSuccess ***REMOVED***
				break
			***REMOVED***

			wrapRingBuffer(s)
			if s.ringbuffer_size == 1<<s.window_bits ***REMOVED***
				s.max_distance = s.max_backward_distance
			***REMOVED***

			if s.state == stateCommandPostWrite1 ***REMOVED***
				if s.meta_block_remaining_len == 0 ***REMOVED***
					/* Next metablock, if any. */
					s.state = stateMetablockDone
				***REMOVED*** else ***REMOVED***
					s.state = stateCommandBegin
				***REMOVED***
			***REMOVED*** else if s.state == stateCommandPostWrite2 ***REMOVED***
				s.state = stateCommandPostWrapCopy /* BROTLI_STATE_COMMAND_INNER_WRITE */
			***REMOVED*** else ***REMOVED***
				if s.loop_counter == 0 ***REMOVED***
					if s.meta_block_remaining_len == 0 ***REMOVED***
						s.state = stateMetablockDone
					***REMOVED*** else ***REMOVED***
						s.state = stateCommandPostDecodeLiterals
					***REMOVED***

					break
				***REMOVED***

				s.state = stateCommandInner
			***REMOVED***

		case stateMetablockDone:
			if s.meta_block_remaining_len < 0 ***REMOVED***
				result = decoderErrorFormatBlockLength2
				break
			***REMOVED***

			decoderStateCleanupAfterMetablock(s)
			if s.is_last_metablock == 0 ***REMOVED***
				s.state = stateMetablockBegin
				break
			***REMOVED***

			if !bitReaderJumpToByteBoundary(br) ***REMOVED***
				result = decoderErrorFormatPadding2
				break
			***REMOVED***

			if s.buffer_length == 0 ***REMOVED***
				bitReaderUnload(br)
				*available_in = br.input_len - br.byte_pos
				*next_in = br.input[br.byte_pos:]
			***REMOVED***

			s.state = stateDone
			fallthrough

		case stateDone:
			if s.ringbuffer != nil ***REMOVED***
				result = writeRingBuffer(s, available_out, next_out, nil, true)
				if result != decoderSuccess ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			return saveErrorCode(s, result)
		***REMOVED***
	***REMOVED***

	return saveErrorCode(s, result)
***REMOVED***

func decoderHasMoreOutput(s *Reader) bool ***REMOVED***
	/* After unrecoverable error remaining output is considered nonsensical. */
	if int(s.error_code) < 0 ***REMOVED***
		return false
	***REMOVED***

	return s.ringbuffer != nil && unwrittenBytes(s, false) != 0
***REMOVED***

func decoderGetErrorCode(s *Reader) int ***REMOVED***
	return int(s.error_code)
***REMOVED***

func decoderErrorString(c int) string ***REMOVED***
	switch c ***REMOVED***
	case decoderNoError:
		return "NO_ERROR"
	case decoderSuccess:
		return "SUCCESS"
	case decoderNeedsMoreInput:
		return "NEEDS_MORE_INPUT"
	case decoderNeedsMoreOutput:
		return "NEEDS_MORE_OUTPUT"
	case decoderErrorFormatExuberantNibble:
		return "EXUBERANT_NIBBLE"
	case decoderErrorFormatReserved:
		return "RESERVED"
	case decoderErrorFormatExuberantMetaNibble:
		return "EXUBERANT_META_NIBBLE"
	case decoderErrorFormatSimpleHuffmanAlphabet:
		return "SIMPLE_HUFFMAN_ALPHABET"
	case decoderErrorFormatSimpleHuffmanSame:
		return "SIMPLE_HUFFMAN_SAME"
	case decoderErrorFormatClSpace:
		return "CL_SPACE"
	case decoderErrorFormatHuffmanSpace:
		return "HUFFMAN_SPACE"
	case decoderErrorFormatContextMapRepeat:
		return "CONTEXT_MAP_REPEAT"
	case decoderErrorFormatBlockLength1:
		return "BLOCK_LENGTH_1"
	case decoderErrorFormatBlockLength2:
		return "BLOCK_LENGTH_2"
	case decoderErrorFormatTransform:
		return "TRANSFORM"
	case decoderErrorFormatDictionary:
		return "DICTIONARY"
	case decoderErrorFormatWindowBits:
		return "WINDOW_BITS"
	case decoderErrorFormatPadding1:
		return "PADDING_1"
	case decoderErrorFormatPadding2:
		return "PADDING_2"
	case decoderErrorFormatDistance:
		return "DISTANCE"
	case decoderErrorDictionaryNotSet:
		return "DICTIONARY_NOT_SET"
	case decoderErrorInvalidArguments:
		return "INVALID_ARGUMENTS"
	case decoderErrorAllocContextModes:
		return "CONTEXT_MODES"
	case decoderErrorAllocTreeGroups:
		return "TREE_GROUPS"
	case decoderErrorAllocContextMap:
		return "CONTEXT_MAP"
	case decoderErrorAllocRingBuffer1:
		return "RING_BUFFER_1"
	case decoderErrorAllocRingBuffer2:
		return "RING_BUFFER_2"
	case decoderErrorAllocBlockTypeTrees:
		return "BLOCK_TYPE_TREES"
	case decoderErrorUnreachable:
		return "UNREACHABLE"
	default:
		return "INVALID"
	***REMOVED***
***REMOVED***
