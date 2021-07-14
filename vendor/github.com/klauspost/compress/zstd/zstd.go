// Package zstd provides decompression of zstandard files.
//
// For advanced usage and examples, go to the README: https://github.com/klauspost/compress/tree/master/zstd#zstd
package zstd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"math/bits"
)

// enable debug printing
const debug = false

// enable encoding debug printing
const debugEncoder = debug

// enable decoding debug printing
const debugDecoder = debug

// Enable extra assertions.
const debugAsserts = debug || false

// print sequence details
const debugSequences = false

// print detailed matching information
const debugMatches = false

// force encoder to use predefined tables.
const forcePreDef = false

// zstdMinMatch is the minimum zstd match length.
const zstdMinMatch = 3

// Reset the buffer offset when reaching this.
const bufferReset = math.MaxInt32 - MaxWindowSize

var (
	// ErrReservedBlockType is returned when a reserved block type is found.
	// Typically this indicates wrong or corrupted input.
	ErrReservedBlockType = errors.New("invalid input: reserved block type encountered")

	// ErrCompressedSizeTooBig is returned when a block is bigger than allowed.
	// Typically this indicates wrong or corrupted input.
	ErrCompressedSizeTooBig = errors.New("invalid input: compressed size too big")

	// ErrBlockTooSmall is returned when a block is too small to be decoded.
	// Typically returned on invalid input.
	ErrBlockTooSmall = errors.New("block too small")

	// ErrMagicMismatch is returned when a "magic" number isn't what is expected.
	// Typically this indicates wrong or corrupted input.
	ErrMagicMismatch = errors.New("invalid input: magic number mismatch")

	// ErrWindowSizeExceeded is returned when a reference exceeds the valid window size.
	// Typically this indicates wrong or corrupted input.
	ErrWindowSizeExceeded = errors.New("window size exceeded")

	// ErrWindowSizeTooSmall is returned when no window size is specified.
	// Typically this indicates wrong or corrupted input.
	ErrWindowSizeTooSmall = errors.New("invalid input: window size was too small")

	// ErrDecoderSizeExceeded is returned if decompressed size exceeds the configured limit.
	ErrDecoderSizeExceeded = errors.New("decompressed size exceeds configured limit")

	// ErrUnknownDictionary is returned if the dictionary ID is unknown.
	// For the time being dictionaries are not supported.
	ErrUnknownDictionary = errors.New("unknown dictionary")

	// ErrFrameSizeExceeded is returned if the stated frame size is exceeded.
	// This is only returned if SingleSegment is specified on the frame.
	ErrFrameSizeExceeded = errors.New("frame size exceeded")

	// ErrCRCMismatch is returned if CRC mismatches.
	ErrCRCMismatch = errors.New("CRC check failed")

	// ErrDecoderClosed will be returned if the Decoder was used after
	// Close has been called.
	ErrDecoderClosed = errors.New("decoder used after Close")

	// ErrDecoderNilInput is returned when a nil Reader was provided
	// and an operation other than Reset/DecodeAll/Close was attempted.
	ErrDecoderNilInput = errors.New("nil input provided as reader")
)

func println(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if debug || debugDecoder || debugEncoder ***REMOVED***
		log.Println(a...)
	***REMOVED***
***REMOVED***

func printf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if debug || debugDecoder || debugEncoder ***REMOVED***
		log.Printf(format, a...)
	***REMOVED***
***REMOVED***

// matchLenFast does matching, but will not match the last up to 7 bytes.
func matchLenFast(a, b []byte) int ***REMOVED***
	endI := len(a) & (math.MaxInt32 - 7)
	for i := 0; i < endI; i += 8 ***REMOVED***
		if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
			return i + bits.TrailingZeros64(diff)>>3
		***REMOVED***
	***REMOVED***
	return endI
***REMOVED***

// matchLen returns the maximum length.
// a must be the shortest of the two.
// The function also returns whether all bytes matched.
func matchLen(a, b []byte) int ***REMOVED***
	b = b[:len(a)]
	for i := 0; i < len(a)-7; i += 8 ***REMOVED***
		if diff := load64(a, i) ^ load64(b, i); diff != 0 ***REMOVED***
			return i + (bits.TrailingZeros64(diff) >> 3)
		***REMOVED***
	***REMOVED***

	checked := (len(a) >> 3) << 3
	a = a[checked:]
	b = b[checked:]
	for i := range a ***REMOVED***
		if a[i] != b[i] ***REMOVED***
			return i + checked
		***REMOVED***
	***REMOVED***
	return len(a) + checked
***REMOVED***

func load3232(b []byte, i int32) uint32 ***REMOVED***
	return binary.LittleEndian.Uint32(b[i:])
***REMOVED***

func load6432(b []byte, i int32) uint64 ***REMOVED***
	return binary.LittleEndian.Uint64(b[i:])
***REMOVED***

func load64(b []byte, i int) uint64 ***REMOVED***
	return binary.LittleEndian.Uint64(b[i:])
***REMOVED***

type byter interface ***REMOVED***
	Bytes() []byte
	Len() int
***REMOVED***

var _ byter = &bytes.Buffer***REMOVED******REMOVED***
