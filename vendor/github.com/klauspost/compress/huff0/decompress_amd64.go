//go:build amd64 && !appengine && !noasm && gc
// +build amd64,!appengine,!noasm,gc

// This file contains the specialisation of Decoder.Decompress4X
// that uses an asm implementation of its main loop.
package huff0

import (
	"errors"
	"fmt"
)

// decompress4x_main_loop_x86 is an x86 assembler implementation
// of Decompress4X when tablelog > 8.
// go:noescape
func decompress4x_main_loop_x86(pbr0, pbr1, pbr2, pbr3 *bitReaderShifted,
	peekBits uint8, buf *byte, tbl *dEntrySingle) uint8

// decompress4x_8b_loop_x86 is an x86 assembler implementation
// of Decompress4X when tablelog <= 8 which decodes 4 entries
// per loop.
// go:noescape
func decompress4x_8b_loop_x86(pbr0, pbr1, pbr2, pbr3 *bitReaderShifted,
	peekBits uint8, buf *byte, tbl *dEntrySingle) uint8

// fallback8BitSize is the size where using Go version is faster.
const fallback8BitSize = 800

// Decompress4X will decompress a 4X encoded stream.
// The length of the supplied input must match the end of a block exactly.
// The *capacity* of the dst slice must match the destination size of
// the uncompressed data exactly.
func (d *Decoder) Decompress4X(dst, src []byte) ([]byte, error) ***REMOVED***
	if len(d.dt.single) == 0 ***REMOVED***
		return nil, errors.New("no table loaded")
	***REMOVED***
	if len(src) < 6+(4*1) ***REMOVED***
		return nil, errors.New("input too small")
	***REMOVED***

	use8BitTables := d.actualTableLog <= 8
	if cap(dst) < fallback8BitSize && use8BitTables ***REMOVED***
		return d.decompress4X8bit(dst, src)
	***REMOVED***
	var br [4]bitReaderShifted
	// Decode "jump table"
	start := 6
	for i := 0; i < 3; i++ ***REMOVED***
		length := int(src[i*2]) | (int(src[i*2+1]) << 8)
		if start+length >= len(src) ***REMOVED***
			return nil, errors.New("truncated input (or invalid offset)")
		***REMOVED***
		err := br[i].init(src[start : start+length])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		start += length
	***REMOVED***
	err := br[3].init(src[start:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// destination, offset to match first output
	dstSize := cap(dst)
	dst = dst[:dstSize]
	out := dst
	dstEvery := (dstSize + 3) / 4

	const tlSize = 1 << tableLogMax
	const tlMask = tlSize - 1
	single := d.dt.single[:tlSize]

	// Use temp table to avoid bound checks/append penalty.
	buf := d.buffer()
	var off uint8
	var decoded int

	const debug = false

	// see: bitReaderShifted.peekBitsFast()
	peekBits := uint8((64 - d.actualTableLog) & 63)

	// Decode 2 values from each decoder/loop.
	const bufoff = 256
	for ***REMOVED***
		if br[0].off < 4 || br[1].off < 4 || br[2].off < 4 || br[3].off < 4 ***REMOVED***
			break
		***REMOVED***

		if use8BitTables ***REMOVED***
			off = decompress4x_8b_loop_x86(&br[0], &br[1], &br[2], &br[3], peekBits, &buf[0][0], &single[0])
		***REMOVED*** else ***REMOVED***
			off = decompress4x_main_loop_x86(&br[0], &br[1], &br[2], &br[3], peekBits, &buf[0][0], &single[0])
		***REMOVED***
		if debug ***REMOVED***
			fmt.Print("DEBUG: ")
			fmt.Printf("off=%d,", off)
			for i := 0; i < 4; i++ ***REMOVED***
				fmt.Printf(" br[%d]=***REMOVED***bitsRead=%d, value=%x, off=%d***REMOVED***",
					i, br[i].bitsRead, br[i].value, br[i].off)
			***REMOVED***
			fmt.Println("")
		***REMOVED***

		if off != 0 ***REMOVED***
			break
		***REMOVED***

		if bufoff > dstEvery ***REMOVED***
			d.bufs.Put(buf)
			return nil, errors.New("corruption detected: stream overrun 1")
		***REMOVED***
		copy(out, buf[0][:])
		copy(out[dstEvery:], buf[1][:])
		copy(out[dstEvery*2:], buf[2][:])
		copy(out[dstEvery*3:], buf[3][:])
		out = out[bufoff:]
		decoded += bufoff * 4
		// There must at least be 3 buffers left.
		if len(out) < dstEvery*3 ***REMOVED***
			d.bufs.Put(buf)
			return nil, errors.New("corruption detected: stream overrun 2")
		***REMOVED***
	***REMOVED***
	if off > 0 ***REMOVED***
		ioff := int(off)
		if len(out) < dstEvery*3+ioff ***REMOVED***
			d.bufs.Put(buf)
			return nil, errors.New("corruption detected: stream overrun 3")
		***REMOVED***
		copy(out, buf[0][:off])
		copy(out[dstEvery:], buf[1][:off])
		copy(out[dstEvery*2:], buf[2][:off])
		copy(out[dstEvery*3:], buf[3][:off])
		decoded += int(off) * 4
		out = out[off:]
	***REMOVED***

	// Decode remaining.
	remainBytes := dstEvery - (decoded / 4)
	for i := range br ***REMOVED***
		offset := dstEvery * i
		endsAt := offset + remainBytes
		if endsAt > len(out) ***REMOVED***
			endsAt = len(out)
		***REMOVED***
		br := &br[i]
		bitsLeft := br.remaining()
		for bitsLeft > 0 ***REMOVED***
			br.fill()
			if offset >= endsAt ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 4")
			***REMOVED***

			// Read value and increment offset.
			val := br.peekBitsFast(d.actualTableLog)
			v := single[val&tlMask].entry
			nBits := uint8(v)
			br.advance(nBits)
			bitsLeft -= uint(nBits)
			out[offset] = uint8(v >> 8)
			offset++
		***REMOVED***
		if offset != endsAt ***REMOVED***
			d.bufs.Put(buf)
			return nil, fmt.Errorf("corruption detected: short output block %d, end %d != %d", i, offset, endsAt)
		***REMOVED***
		decoded += offset - dstEvery*i
		err = br.close()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	d.bufs.Put(buf)
	if dstSize != decoded ***REMOVED***
		return nil, errors.New("corruption detected: short output block")
	***REMOVED***
	return dst, nil
***REMOVED***
