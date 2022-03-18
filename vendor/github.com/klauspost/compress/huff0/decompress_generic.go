//go:build !amd64 || appengine || !gc || noasm
// +build !amd64 appengine !gc noasm

// This file contains a generic implementation of Decoder.Decompress4X.
package huff0

import (
	"errors"
	"fmt"
)

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
	if use8BitTables && d.actualTableLog <= 8 ***REMOVED***
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

	// Decode 2 values from each decoder/loop.
	const bufoff = 256
	for ***REMOVED***
		if br[0].off < 4 || br[1].off < 4 || br[2].off < 4 || br[3].off < 4 ***REMOVED***
			break
		***REMOVED***

		***REMOVED***
			const stream = 0
			const stream2 = 1
			br[stream].fillFast()
			br[stream2].fillFast()

			val := br[stream].peekBitsFast(d.actualTableLog)
			val2 := br[stream2].peekBitsFast(d.actualTableLog)
			v := single[val&tlMask]
			v2 := single[val2&tlMask]
			br[stream].advance(uint8(v.entry))
			br[stream2].advance(uint8(v2.entry))
			buf[stream][off] = uint8(v.entry >> 8)
			buf[stream2][off] = uint8(v2.entry >> 8)

			val = br[stream].peekBitsFast(d.actualTableLog)
			val2 = br[stream2].peekBitsFast(d.actualTableLog)
			v = single[val&tlMask]
			v2 = single[val2&tlMask]
			br[stream].advance(uint8(v.entry))
			br[stream2].advance(uint8(v2.entry))
			buf[stream][off+1] = uint8(v.entry >> 8)
			buf[stream2][off+1] = uint8(v2.entry >> 8)
		***REMOVED***

		***REMOVED***
			const stream = 2
			const stream2 = 3
			br[stream].fillFast()
			br[stream2].fillFast()

			val := br[stream].peekBitsFast(d.actualTableLog)
			val2 := br[stream2].peekBitsFast(d.actualTableLog)
			v := single[val&tlMask]
			v2 := single[val2&tlMask]
			br[stream].advance(uint8(v.entry))
			br[stream2].advance(uint8(v2.entry))
			buf[stream][off] = uint8(v.entry >> 8)
			buf[stream2][off] = uint8(v2.entry >> 8)

			val = br[stream].peekBitsFast(d.actualTableLog)
			val2 = br[stream2].peekBitsFast(d.actualTableLog)
			v = single[val&tlMask]
			v2 = single[val2&tlMask]
			br[stream].advance(uint8(v.entry))
			br[stream2].advance(uint8(v2.entry))
			buf[stream][off+1] = uint8(v.entry >> 8)
			buf[stream2][off+1] = uint8(v2.entry >> 8)
		***REMOVED***

		off += 2

		if off == 0 ***REMOVED***
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
