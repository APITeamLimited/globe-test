package base64vlq

import "io"

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

const (
	vlqBaseShift       = 5
	vlqBase            = 1 << vlqBaseShift
	vlqBaseMask        = vlqBase - 1
	vlqSignBit         = 1
	vlqContinuationBit = vlqBase
)

var decodeMap [256]byte

func init() ***REMOVED***
	for i := 0; i < len(encodeStd); i++ ***REMOVED***
		decodeMap[encodeStd[i]] = byte(i)
	***REMOVED***
***REMOVED***

func toVLQSigned(n int32) int32 ***REMOVED***
	if n < 0 ***REMOVED***
		return -n<<1 + 1
	***REMOVED***
	return n << 1
***REMOVED***

func fromVLQSigned(n int32) int32 ***REMOVED***
	isNeg := n&vlqSignBit != 0
	n >>= 1
	if isNeg ***REMOVED***
		return -n
	***REMOVED***
	return n
***REMOVED***

type Encoder struct ***REMOVED***
	w io.ByteWriter
***REMOVED***

func NewEncoder(w io.ByteWriter) *Encoder ***REMOVED***
	return &Encoder***REMOVED***
		w: w,
	***REMOVED***
***REMOVED***

func (enc Encoder) Encode(n int32) error ***REMOVED***
	n = toVLQSigned(n)
	for digit := int32(vlqContinuationBit); digit&vlqContinuationBit != 0; ***REMOVED***
		digit = n & vlqBaseMask
		n >>= vlqBaseShift
		if n > 0 ***REMOVED***
			digit |= vlqContinuationBit
		***REMOVED***

		err := enc.w.WriteByte(encodeStd[digit])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type Decoder struct ***REMOVED***
	r io.ByteReader
***REMOVED***

func NewDecoder(r io.ByteReader) Decoder ***REMOVED***
	return Decoder***REMOVED***
		r: r,
	***REMOVED***
***REMOVED***

func (dec Decoder) Decode() (n int32, err error) ***REMOVED***
	shift := uint(0)
	for continuation := true; continuation; ***REMOVED***
		c, err := dec.r.ReadByte()
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***

		c = decodeMap[c]
		continuation = c&vlqContinuationBit != 0
		n += int32(c&vlqBaseMask) << shift
		shift += vlqBaseShift
	***REMOVED***
	return fromVLQSigned(n), nil
***REMOVED***
