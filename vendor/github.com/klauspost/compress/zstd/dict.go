package zstd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/klauspost/compress/huff0"
)

type dict struct ***REMOVED***
	id uint32

	litEnc              *huff0.Scratch
	llDec, ofDec, mlDec sequenceDec
	//llEnc, ofEnc, mlEnc []*fseEncoder
	offsets [3]int
	content []byte
***REMOVED***

var dictMagic = [4]byte***REMOVED***0x37, 0xa4, 0x30, 0xec***REMOVED***

// ID returns the dictionary id or 0 if d is nil.
func (d *dict) ID() uint32 ***REMOVED***
	if d == nil ***REMOVED***
		return 0
	***REMOVED***
	return d.id
***REMOVED***

// DictContentSize returns the dictionary content size or 0 if d is nil.
func (d *dict) DictContentSize() int ***REMOVED***
	if d == nil ***REMOVED***
		return 0
	***REMOVED***
	return len(d.content)
***REMOVED***

// Load a dictionary as described in
// https://github.com/facebook/zstd/blob/master/doc/zstd_compression_format.md#dictionary-format
func loadDict(b []byte) (*dict, error) ***REMOVED***
	// Check static field size.
	if len(b) <= 8+(3*4) ***REMOVED***
		return nil, io.ErrUnexpectedEOF
	***REMOVED***
	d := dict***REMOVED***
		llDec: sequenceDec***REMOVED***fse: &fseDecoder***REMOVED******REMOVED******REMOVED***,
		ofDec: sequenceDec***REMOVED***fse: &fseDecoder***REMOVED******REMOVED******REMOVED***,
		mlDec: sequenceDec***REMOVED***fse: &fseDecoder***REMOVED******REMOVED******REMOVED***,
	***REMOVED***
	if !bytes.Equal(b[:4], dictMagic[:]) ***REMOVED***
		return nil, ErrMagicMismatch
	***REMOVED***
	d.id = binary.LittleEndian.Uint32(b[4:8])
	if d.id == 0 ***REMOVED***
		return nil, errors.New("dictionaries cannot have ID 0")
	***REMOVED***

	// Read literal table
	var err error
	d.litEnc, b, err = huff0.ReadTable(b[8:], nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	d.litEnc.Reuse = huff0.ReusePolicyMust

	br := byteReader***REMOVED***
		b:   b,
		off: 0,
	***REMOVED***
	readDec := func(i tableIndex, dec *fseDecoder) error ***REMOVED***
		if err := dec.readNCount(&br, uint16(maxTableSymbol[i])); err != nil ***REMOVED***
			return err
		***REMOVED***
		if br.overread() ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		err = dec.transform(symbolTableX[i])
		if err != nil ***REMOVED***
			println("Transform table error:", err)
			return err
		***REMOVED***
		if debugDecoder || debugEncoder ***REMOVED***
			println("Read table ok", "symbolLen:", dec.symbolLen)
		***REMOVED***
		// Set decoders as predefined so they aren't reused.
		dec.preDefined = true
		return nil
	***REMOVED***

	if err := readDec(tableOffsets, d.ofDec.fse); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := readDec(tableMatchLengths, d.mlDec.fse); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := readDec(tableLiteralLengths, d.llDec.fse); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if br.remain() < 12 ***REMOVED***
		return nil, io.ErrUnexpectedEOF
	***REMOVED***

	d.offsets[0] = int(br.Uint32())
	br.advance(4)
	d.offsets[1] = int(br.Uint32())
	br.advance(4)
	d.offsets[2] = int(br.Uint32())
	br.advance(4)
	if d.offsets[0] <= 0 || d.offsets[1] <= 0 || d.offsets[2] <= 0 ***REMOVED***
		return nil, errors.New("invalid offset in dictionary")
	***REMOVED***
	d.content = make([]byte, br.remain())
	copy(d.content, br.unread())
	if d.offsets[0] > len(d.content) || d.offsets[1] > len(d.content) || d.offsets[2] > len(d.content) ***REMOVED***
		return nil, fmt.Errorf("initial offset bigger than dictionary content size %d, offsets: %v", len(d.content), d.offsets)
	***REMOVED***

	return &d, nil
***REMOVED***
