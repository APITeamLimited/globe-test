package snappy

import (
	"bytes"
	"encoding/binary"

	master "github.com/golang/snappy"
)

var xerialHeader = []byte***REMOVED***130, 83, 78, 65, 80, 80, 89, 0***REMOVED***

// Encode encodes data as snappy with no framing header.
func Encode(src []byte) []byte ***REMOVED***
	return master.Encode(nil, src)
***REMOVED***

// Decode decodes snappy data whether it is traditional unframed
// or includes the xerial framing format.
func Decode(src []byte) ([]byte, error) ***REMOVED***
	if !bytes.Equal(src[:8], xerialHeader) ***REMOVED***
		return master.Decode(nil, src)
	***REMOVED***

	var (
		pos   = uint32(16)
		max   = uint32(len(src))
		dst   = make([]byte, 0, len(src))
		chunk []byte
		err   error
	)
	for pos < max ***REMOVED***
		size := binary.BigEndian.Uint32(src[pos : pos+4])
		pos += 4

		chunk, err = master.Decode(chunk, src[pos:pos+size])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		pos += size
		dst = append(dst, chunk...)
	***REMOVED***
	return dst, nil
***REMOVED***
