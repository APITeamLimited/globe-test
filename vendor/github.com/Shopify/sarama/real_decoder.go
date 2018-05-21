package sarama

import (
	"encoding/binary"
	"math"
)

var errInvalidArrayLength = PacketDecodingError***REMOVED***"invalid array length"***REMOVED***
var errInvalidByteSliceLength = PacketDecodingError***REMOVED***"invalid byteslice length"***REMOVED***
var errInvalidByteSliceLengthType = PacketDecodingError***REMOVED***"invalid byteslice length type"***REMOVED***
var errInvalidStringLength = PacketDecodingError***REMOVED***"invalid string length"***REMOVED***
var errInvalidSubsetSize = PacketDecodingError***REMOVED***"invalid subset size"***REMOVED***
var errVarintOverflow = PacketDecodingError***REMOVED***"varint overflow"***REMOVED***
var errInvalidBool = PacketDecodingError***REMOVED***"invalid bool"***REMOVED***

type realDecoder struct ***REMOVED***
	raw   []byte
	off   int
	stack []pushDecoder
***REMOVED***

// primitives

func (rd *realDecoder) getInt8() (int8, error) ***REMOVED***
	if rd.remaining() < 1 ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED***
	tmp := int8(rd.raw[rd.off])
	rd.off++
	return tmp, nil
***REMOVED***

func (rd *realDecoder) getInt16() (int16, error) ***REMOVED***
	if rd.remaining() < 2 ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED***
	tmp := int16(binary.BigEndian.Uint16(rd.raw[rd.off:]))
	rd.off += 2
	return tmp, nil
***REMOVED***

func (rd *realDecoder) getInt32() (int32, error) ***REMOVED***
	if rd.remaining() < 4 ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED***
	tmp := int32(binary.BigEndian.Uint32(rd.raw[rd.off:]))
	rd.off += 4
	return tmp, nil
***REMOVED***

func (rd *realDecoder) getInt64() (int64, error) ***REMOVED***
	if rd.remaining() < 8 ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED***
	tmp := int64(binary.BigEndian.Uint64(rd.raw[rd.off:]))
	rd.off += 8
	return tmp, nil
***REMOVED***

func (rd *realDecoder) getVarint() (int64, error) ***REMOVED***
	tmp, n := binary.Varint(rd.raw[rd.off:])
	if n == 0 ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED***
	if n < 0 ***REMOVED***
		rd.off -= n
		return -1, errVarintOverflow
	***REMOVED***
	rd.off += n
	return tmp, nil
***REMOVED***

func (rd *realDecoder) getArrayLength() (int, error) ***REMOVED***
	if rd.remaining() < 4 ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED***
	tmp := int(int32(binary.BigEndian.Uint32(rd.raw[rd.off:])))
	rd.off += 4
	if tmp > rd.remaining() ***REMOVED***
		rd.off = len(rd.raw)
		return -1, ErrInsufficientData
	***REMOVED*** else if tmp > 2*math.MaxUint16 ***REMOVED***
		return -1, errInvalidArrayLength
	***REMOVED***
	return tmp, nil
***REMOVED***

func (rd *realDecoder) getBool() (bool, error) ***REMOVED***
	b, err := rd.getInt8()
	if err != nil || b == 0 ***REMOVED***
		return false, err
	***REMOVED***
	if b != 1 ***REMOVED***
		return false, errInvalidBool
	***REMOVED***
	return true, nil
***REMOVED***

// collections

func (rd *realDecoder) getBytes() ([]byte, error) ***REMOVED***
	tmp, err := rd.getInt32()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if tmp == -1 ***REMOVED***
		return nil, nil
	***REMOVED***

	return rd.getRawBytes(int(tmp))
***REMOVED***

func (rd *realDecoder) getVarintBytes() ([]byte, error) ***REMOVED***
	tmp, err := rd.getVarint()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if tmp == -1 ***REMOVED***
		return nil, nil
	***REMOVED***

	return rd.getRawBytes(int(tmp))
***REMOVED***

func (rd *realDecoder) getStringLength() (int, error) ***REMOVED***
	length, err := rd.getInt16()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n := int(length)

	switch ***REMOVED***
	case n < -1:
		return 0, errInvalidStringLength
	case n > rd.remaining():
		rd.off = len(rd.raw)
		return 0, ErrInsufficientData
	***REMOVED***

	return n, nil
***REMOVED***

func (rd *realDecoder) getString() (string, error) ***REMOVED***
	n, err := rd.getStringLength()
	if err != nil || n == -1 ***REMOVED***
		return "", err
	***REMOVED***

	tmpStr := string(rd.raw[rd.off : rd.off+n])
	rd.off += n
	return tmpStr, nil
***REMOVED***

func (rd *realDecoder) getNullableString() (*string, error) ***REMOVED***
	n, err := rd.getStringLength()
	if err != nil || n == -1 ***REMOVED***
		return nil, err
	***REMOVED***

	tmpStr := string(rd.raw[rd.off : rd.off+n])
	rd.off += n
	return &tmpStr, err
***REMOVED***

func (rd *realDecoder) getInt32Array() ([]int32, error) ***REMOVED***
	if rd.remaining() < 4 ***REMOVED***
		rd.off = len(rd.raw)
		return nil, ErrInsufficientData
	***REMOVED***
	n := int(binary.BigEndian.Uint32(rd.raw[rd.off:]))
	rd.off += 4

	if rd.remaining() < 4*n ***REMOVED***
		rd.off = len(rd.raw)
		return nil, ErrInsufficientData
	***REMOVED***

	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	if n < 0 ***REMOVED***
		return nil, errInvalidArrayLength
	***REMOVED***

	ret := make([]int32, n)
	for i := range ret ***REMOVED***
		ret[i] = int32(binary.BigEndian.Uint32(rd.raw[rd.off:]))
		rd.off += 4
	***REMOVED***
	return ret, nil
***REMOVED***

func (rd *realDecoder) getInt64Array() ([]int64, error) ***REMOVED***
	if rd.remaining() < 4 ***REMOVED***
		rd.off = len(rd.raw)
		return nil, ErrInsufficientData
	***REMOVED***
	n := int(binary.BigEndian.Uint32(rd.raw[rd.off:]))
	rd.off += 4

	if rd.remaining() < 8*n ***REMOVED***
		rd.off = len(rd.raw)
		return nil, ErrInsufficientData
	***REMOVED***

	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	if n < 0 ***REMOVED***
		return nil, errInvalidArrayLength
	***REMOVED***

	ret := make([]int64, n)
	for i := range ret ***REMOVED***
		ret[i] = int64(binary.BigEndian.Uint64(rd.raw[rd.off:]))
		rd.off += 8
	***REMOVED***
	return ret, nil
***REMOVED***

func (rd *realDecoder) getStringArray() ([]string, error) ***REMOVED***
	if rd.remaining() < 4 ***REMOVED***
		rd.off = len(rd.raw)
		return nil, ErrInsufficientData
	***REMOVED***
	n := int(binary.BigEndian.Uint32(rd.raw[rd.off:]))
	rd.off += 4

	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	if n < 0 ***REMOVED***
		return nil, errInvalidArrayLength
	***REMOVED***

	ret := make([]string, n)
	for i := range ret ***REMOVED***
		str, err := rd.getString()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		ret[i] = str
	***REMOVED***
	return ret, nil
***REMOVED***

// subsets

func (rd *realDecoder) remaining() int ***REMOVED***
	return len(rd.raw) - rd.off
***REMOVED***

func (rd *realDecoder) getSubset(length int) (packetDecoder, error) ***REMOVED***
	buf, err := rd.getRawBytes(length)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &realDecoder***REMOVED***raw: buf***REMOVED***, nil
***REMOVED***

func (rd *realDecoder) getRawBytes(length int) ([]byte, error) ***REMOVED***
	if length < 0 ***REMOVED***
		return nil, errInvalidByteSliceLength
	***REMOVED*** else if length > rd.remaining() ***REMOVED***
		rd.off = len(rd.raw)
		return nil, ErrInsufficientData
	***REMOVED***

	start := rd.off
	rd.off += length
	return rd.raw[start:rd.off], nil
***REMOVED***

func (rd *realDecoder) peek(offset, length int) (packetDecoder, error) ***REMOVED***
	if rd.remaining() < offset+length ***REMOVED***
		return nil, ErrInsufficientData
	***REMOVED***
	off := rd.off + offset
	return &realDecoder***REMOVED***raw: rd.raw[off : off+length]***REMOVED***, nil
***REMOVED***

// stacks

func (rd *realDecoder) push(in pushDecoder) error ***REMOVED***
	in.saveOffset(rd.off)

	var reserve int
	if dpd, ok := in.(dynamicPushDecoder); ok ***REMOVED***
		if err := dpd.decode(rd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		reserve = in.reserveLength()
		if rd.remaining() < reserve ***REMOVED***
			rd.off = len(rd.raw)
			return ErrInsufficientData
		***REMOVED***
	***REMOVED***

	rd.stack = append(rd.stack, in)

	rd.off += reserve

	return nil
***REMOVED***

func (rd *realDecoder) pop() error ***REMOVED***
	// this is go's ugly pop pattern (the inverse of append)
	in := rd.stack[len(rd.stack)-1]
	rd.stack = rd.stack[:len(rd.stack)-1]

	return in.check(rd.off, rd.raw)
***REMOVED***
