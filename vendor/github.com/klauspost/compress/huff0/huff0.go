// Package huff0 provides fast huffman encoding as used in zstd.
//
// See README.md at https://github.com/klauspost/compress/tree/master/huff0 for details.
package huff0

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sync"

	"github.com/klauspost/compress/fse"
)

const (
	maxSymbolValue = 255

	// zstandard limits tablelog to 11, see:
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#huffman-tree-description
	tableLogMax     = 11
	tableLogDefault = 11
	minTablelog     = 5
	huffNodesLen    = 512

	// BlockSizeMax is maximum input size for a single block uncompressed.
	BlockSizeMax = 1<<18 - 1
)

var (
	// ErrIncompressible is returned when input is judged to be too hard to compress.
	ErrIncompressible = errors.New("input is not compressible")

	// ErrUseRLE is returned from the compressor when the input is a single byte value repeated.
	ErrUseRLE = errors.New("input is single value repeated")

	// ErrTooBig is return if input is too large for a single block.
	ErrTooBig = errors.New("input too big")

	// ErrMaxDecodedSizeExceeded is return if input is too large for a single block.
	ErrMaxDecodedSizeExceeded = errors.New("maximum output size exceeded")
)

type ReusePolicy uint8

const (
	// ReusePolicyAllow will allow reuse if it produces smaller output.
	ReusePolicyAllow ReusePolicy = iota

	// ReusePolicyPrefer will re-use aggressively if possible.
	// This will not check if a new table will produce smaller output,
	// except if the current table is impossible to use or
	// compressed output is bigger than input.
	ReusePolicyPrefer

	// ReusePolicyNone will disable re-use of tables.
	// This is slightly faster than ReusePolicyAllow but may produce larger output.
	ReusePolicyNone

	// ReusePolicyMust must allow reuse and produce smaller output.
	ReusePolicyMust
)

type Scratch struct ***REMOVED***
	count [maxSymbolValue + 1]uint32

	// Per block parameters.
	// These can be used to override compression parameters of the block.
	// Do not touch, unless you know what you are doing.

	// Out is output buffer.
	// If the scratch is re-used before the caller is done processing the output,
	// set this field to nil.
	// Otherwise the output buffer will be re-used for next Compression/Decompression step
	// and allocation will be avoided.
	Out []byte

	// OutTable will contain the table data only, if a new table has been generated.
	// Slice of the returned data.
	OutTable []byte

	// OutData will contain the compressed data.
	// Slice of the returned data.
	OutData []byte

	// MaxDecodedSize will set the maximum allowed output size.
	// This value will automatically be set to BlockSizeMax if not set.
	// Decoders will return ErrMaxDecodedSizeExceeded is this limit is exceeded.
	MaxDecodedSize int

	br byteReader

	// MaxSymbolValue will override the maximum symbol value of the next block.
	MaxSymbolValue uint8

	// TableLog will attempt to override the tablelog for the next block.
	// Must be <= 11 and >= 5.
	TableLog uint8

	// Reuse will specify the reuse policy
	Reuse ReusePolicy

	// WantLogLess allows to specify a log 2 reduction that should at least be achieved,
	// otherwise the block will be returned as incompressible.
	// The reduction should then at least be (input size >> WantLogLess)
	// If WantLogLess == 0 any improvement will do.
	WantLogLess uint8

	symbolLen      uint16 // Length of active part of the symbol table.
	maxCount       int    // count of the most probable symbol
	clearCount     bool   // clear count
	actualTableLog uint8  // Selected tablelog.
	prevTableLog   uint8  // Tablelog for previous table
	prevTable      cTable // Table used for previous compression.
	cTable         cTable // compression table
	dt             dTable // decompression table
	nodes          []nodeElt
	tmpOut         [4][]byte
	fse            *fse.Scratch
	decPool        sync.Pool // *[4][256]byte buffers.
	huffWeight     [maxSymbolValue + 1]byte
***REMOVED***

// TransferCTable will transfer the previously used compression table.
func (s *Scratch) TransferCTable(src *Scratch) ***REMOVED***
	if cap(s.prevTable) < len(src.prevTable) ***REMOVED***
		s.prevTable = make(cTable, 0, maxSymbolValue+1)
	***REMOVED***
	s.prevTable = s.prevTable[:len(src.prevTable)]
	copy(s.prevTable, src.prevTable)
	s.prevTableLog = src.prevTableLog
***REMOVED***

func (s *Scratch) prepare(in []byte) (*Scratch, error) ***REMOVED***
	if len(in) > BlockSizeMax ***REMOVED***
		return nil, ErrTooBig
	***REMOVED***
	if s == nil ***REMOVED***
		s = &Scratch***REMOVED******REMOVED***
	***REMOVED***
	if s.MaxSymbolValue == 0 ***REMOVED***
		s.MaxSymbolValue = maxSymbolValue
	***REMOVED***
	if s.TableLog == 0 ***REMOVED***
		s.TableLog = tableLogDefault
	***REMOVED***
	if s.TableLog > tableLogMax || s.TableLog < minTablelog ***REMOVED***
		return nil, fmt.Errorf(" invalid tableLog %d (%d -> %d)", s.TableLog, minTablelog, tableLogMax)
	***REMOVED***
	if s.MaxDecodedSize <= 0 || s.MaxDecodedSize > BlockSizeMax ***REMOVED***
		s.MaxDecodedSize = BlockSizeMax
	***REMOVED***
	if s.clearCount && s.maxCount == 0 ***REMOVED***
		for i := range s.count ***REMOVED***
			s.count[i] = 0
		***REMOVED***
		s.clearCount = false
	***REMOVED***
	if cap(s.Out) == 0 ***REMOVED***
		s.Out = make([]byte, 0, len(in))
	***REMOVED***
	s.Out = s.Out[:0]

	s.OutTable = nil
	s.OutData = nil
	if cap(s.nodes) < huffNodesLen+1 ***REMOVED***
		s.nodes = make([]nodeElt, 0, huffNodesLen+1)
	***REMOVED***
	s.nodes = s.nodes[:0]
	if s.fse == nil ***REMOVED***
		s.fse = &fse.Scratch***REMOVED******REMOVED***
	***REMOVED***
	s.br.init(in)

	return s, nil
***REMOVED***

type cTable []cTableEntry

func (c cTable) write(s *Scratch) error ***REMOVED***
	var (
		// precomputed conversion table
		bitsToWeight [tableLogMax + 1]byte
		huffLog      = s.actualTableLog
		// last weight is not saved.
		maxSymbolValue = uint8(s.symbolLen - 1)
		huffWeight     = s.huffWeight[:256]
	)
	const (
		maxFSETableLog = 6
	)
	// convert to weight
	bitsToWeight[0] = 0
	for n := uint8(1); n < huffLog+1; n++ ***REMOVED***
		bitsToWeight[n] = huffLog + 1 - n
	***REMOVED***

	// Acquire histogram for FSE.
	hist := s.fse.Histogram()
	hist = hist[:256]
	for i := range hist[:16] ***REMOVED***
		hist[i] = 0
	***REMOVED***
	for n := uint8(0); n < maxSymbolValue; n++ ***REMOVED***
		v := bitsToWeight[c[n].nBits] & 15
		huffWeight[n] = v
		hist[v]++
	***REMOVED***

	// FSE compress if feasible.
	if maxSymbolValue >= 2 ***REMOVED***
		huffMaxCnt := uint32(0)
		huffMax := uint8(0)
		for i, v := range hist[:16] ***REMOVED***
			if v == 0 ***REMOVED***
				continue
			***REMOVED***
			huffMax = byte(i)
			if v > huffMaxCnt ***REMOVED***
				huffMaxCnt = v
			***REMOVED***
		***REMOVED***
		s.fse.HistogramFinished(huffMax, int(huffMaxCnt))
		s.fse.TableLog = maxFSETableLog
		b, err := fse.Compress(huffWeight[:maxSymbolValue], s.fse)
		if err == nil && len(b) < int(s.symbolLen>>1) ***REMOVED***
			s.Out = append(s.Out, uint8(len(b)))
			s.Out = append(s.Out, b...)
			return nil
		***REMOVED***
		// Unable to compress (RLE/uncompressible)
	***REMOVED***
	// write raw values as 4-bits (max : 15)
	if maxSymbolValue > (256 - 128) ***REMOVED***
		// should not happen : likely means source cannot be compressed
		return ErrIncompressible
	***REMOVED***
	op := s.Out
	// special case, pack weights 4 bits/weight.
	op = append(op, 128|(maxSymbolValue-1))
	// be sure it doesn't cause msan issue in final combination
	huffWeight[maxSymbolValue] = 0
	for n := uint16(0); n < uint16(maxSymbolValue); n += 2 ***REMOVED***
		op = append(op, (huffWeight[n]<<4)|huffWeight[n+1])
	***REMOVED***
	s.Out = op
	return nil
***REMOVED***

func (c cTable) estTableSize(s *Scratch) (sz int, err error) ***REMOVED***
	var (
		// precomputed conversion table
		bitsToWeight [tableLogMax + 1]byte
		huffLog      = s.actualTableLog
		// last weight is not saved.
		maxSymbolValue = uint8(s.symbolLen - 1)
		huffWeight     = s.huffWeight[:256]
	)
	const (
		maxFSETableLog = 6
	)
	// convert to weight
	bitsToWeight[0] = 0
	for n := uint8(1); n < huffLog+1; n++ ***REMOVED***
		bitsToWeight[n] = huffLog + 1 - n
	***REMOVED***

	// Acquire histogram for FSE.
	hist := s.fse.Histogram()
	hist = hist[:256]
	for i := range hist[:16] ***REMOVED***
		hist[i] = 0
	***REMOVED***
	for n := uint8(0); n < maxSymbolValue; n++ ***REMOVED***
		v := bitsToWeight[c[n].nBits] & 15
		huffWeight[n] = v
		hist[v]++
	***REMOVED***

	// FSE compress if feasible.
	if maxSymbolValue >= 2 ***REMOVED***
		huffMaxCnt := uint32(0)
		huffMax := uint8(0)
		for i, v := range hist[:16] ***REMOVED***
			if v == 0 ***REMOVED***
				continue
			***REMOVED***
			huffMax = byte(i)
			if v > huffMaxCnt ***REMOVED***
				huffMaxCnt = v
			***REMOVED***
		***REMOVED***
		s.fse.HistogramFinished(huffMax, int(huffMaxCnt))
		s.fse.TableLog = maxFSETableLog
		b, err := fse.Compress(huffWeight[:maxSymbolValue], s.fse)
		if err == nil && len(b) < int(s.symbolLen>>1) ***REMOVED***
			sz += 1 + len(b)
			return sz, nil
		***REMOVED***
		// Unable to compress (RLE/uncompressible)
	***REMOVED***
	// write raw values as 4-bits (max : 15)
	if maxSymbolValue > (256 - 128) ***REMOVED***
		// should not happen : likely means source cannot be compressed
		return 0, ErrIncompressible
	***REMOVED***
	// special case, pack weights 4 bits/weight.
	sz += 1 + int(maxSymbolValue/2)
	return sz, nil
***REMOVED***

// estimateSize returns the estimated size in bytes of the input represented in the
// histogram supplied.
func (c cTable) estimateSize(hist []uint32) int ***REMOVED***
	nbBits := uint32(7)
	for i, v := range c[:len(hist)] ***REMOVED***
		nbBits += uint32(v.nBits) * hist[i]
	***REMOVED***
	return int(nbBits >> 3)
***REMOVED***

// minSize returns the minimum possible size considering the shannon limit.
func (s *Scratch) minSize(total int) int ***REMOVED***
	nbBits := float64(7)
	fTotal := float64(total)
	for _, v := range s.count[:s.symbolLen] ***REMOVED***
		n := float64(v)
		if n > 0 ***REMOVED***
			nbBits += math.Log2(fTotal/n) * n
		***REMOVED***
	***REMOVED***
	return int(nbBits) >> 3
***REMOVED***

func highBit32(val uint32) (n uint32) ***REMOVED***
	return uint32(bits.Len32(val) - 1)
***REMOVED***
