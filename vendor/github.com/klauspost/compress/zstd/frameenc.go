// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/bits"
)

type frameHeader struct ***REMOVED***
	ContentSize   uint64
	WindowSize    uint32
	SingleSegment bool
	Checksum      bool
	DictID        uint32 // Not stored.
***REMOVED***

const maxHeaderSize = 14

func (f frameHeader) appendTo(dst []byte) ([]byte, error) ***REMOVED***
	dst = append(dst, frameMagic...)
	var fhd uint8
	if f.Checksum ***REMOVED***
		fhd |= 1 << 2
	***REMOVED***
	if f.SingleSegment ***REMOVED***
		fhd |= 1 << 5
	***REMOVED***
	var fcs uint8
	if f.ContentSize >= 256 ***REMOVED***
		fcs++
	***REMOVED***
	if f.ContentSize >= 65536+256 ***REMOVED***
		fcs++
	***REMOVED***
	if f.ContentSize >= 0xffffffff ***REMOVED***
		fcs++
	***REMOVED***
	fhd |= fcs << 6

	dst = append(dst, fhd)
	if !f.SingleSegment ***REMOVED***
		const winLogMin = 10
		windowLog := (bits.Len32(f.WindowSize-1) - winLogMin) << 3
		dst = append(dst, uint8(windowLog))
	***REMOVED***
	if f.SingleSegment && f.ContentSize == 0 ***REMOVED***
		return nil, errors.New("single segment, but no size set")
	***REMOVED***
	switch fcs ***REMOVED***
	case 0:
		if f.SingleSegment ***REMOVED***
			dst = append(dst, uint8(f.ContentSize))
		***REMOVED***
		// Unless SingleSegment is set, framessizes < 256 are nto stored.
	case 1:
		f.ContentSize -= 256
		dst = append(dst, uint8(f.ContentSize), uint8(f.ContentSize>>8))
	case 2:
		dst = append(dst, uint8(f.ContentSize), uint8(f.ContentSize>>8), uint8(f.ContentSize>>16), uint8(f.ContentSize>>24))
	case 3:
		dst = append(dst, uint8(f.ContentSize), uint8(f.ContentSize>>8), uint8(f.ContentSize>>16), uint8(f.ContentSize>>24),
			uint8(f.ContentSize>>32), uint8(f.ContentSize>>40), uint8(f.ContentSize>>48), uint8(f.ContentSize>>56))
	default:
		panic("invalid fcs")
	***REMOVED***
	return dst, nil
***REMOVED***

const skippableFrameHeader = 4 + 4

// calcSkippableFrame will return a total size to be added for written
// to be divisible by multiple.
// The value will always be > skippableFrameHeader.
// The function will panic if written < 0 or wantMultiple <= 0.
func calcSkippableFrame(written, wantMultiple int64) int ***REMOVED***
	if wantMultiple <= 0 ***REMOVED***
		panic("wantMultiple <= 0")
	***REMOVED***
	if written < 0 ***REMOVED***
		panic("written < 0")
	***REMOVED***
	leftOver := written % wantMultiple
	if leftOver == 0 ***REMOVED***
		return 0
	***REMOVED***
	toAdd := wantMultiple - leftOver
	for toAdd < skippableFrameHeader ***REMOVED***
		toAdd += wantMultiple
	***REMOVED***
	return int(toAdd)
***REMOVED***

// skippableFrame will add a skippable frame with a total size of bytes.
// total should be >= skippableFrameHeader and < math.MaxUint32.
func skippableFrame(dst []byte, total int, r io.Reader) ([]byte, error) ***REMOVED***
	if total == 0 ***REMOVED***
		return dst, nil
	***REMOVED***
	if total < skippableFrameHeader ***REMOVED***
		return dst, fmt.Errorf("requested skippable frame (%d) < 8", total)
	***REMOVED***
	if int64(total) > math.MaxUint32 ***REMOVED***
		return dst, fmt.Errorf("requested skippable frame (%d) > max uint32", total)
	***REMOVED***
	dst = append(dst, 0x50, 0x2a, 0x4d, 0x18)
	f := uint32(total - skippableFrameHeader)
	dst = append(dst, uint8(f), uint8(f>>8), uint8(f>>16), uint8(f>>24))
	start := len(dst)
	dst = append(dst, make([]byte, f)...)
	_, err := io.ReadFull(r, dst[start:])
	return dst, err
***REMOVED***
