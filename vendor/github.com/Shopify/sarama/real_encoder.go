package sarama

import (
	"encoding/binary"

	"github.com/rcrowley/go-metrics"
)

type realEncoder struct ***REMOVED***
	raw      []byte
	off      int
	stack    []pushEncoder
	registry metrics.Registry
***REMOVED***

// primitives

func (re *realEncoder) putInt8(in int8) ***REMOVED***
	re.raw[re.off] = byte(in)
	re.off++
***REMOVED***

func (re *realEncoder) putInt16(in int16) ***REMOVED***
	binary.BigEndian.PutUint16(re.raw[re.off:], uint16(in))
	re.off += 2
***REMOVED***

func (re *realEncoder) putInt32(in int32) ***REMOVED***
	binary.BigEndian.PutUint32(re.raw[re.off:], uint32(in))
	re.off += 4
***REMOVED***

func (re *realEncoder) putInt64(in int64) ***REMOVED***
	binary.BigEndian.PutUint64(re.raw[re.off:], uint64(in))
	re.off += 8
***REMOVED***

func (re *realEncoder) putVarint(in int64) ***REMOVED***
	re.off += binary.PutVarint(re.raw[re.off:], in)
***REMOVED***

func (re *realEncoder) putArrayLength(in int) error ***REMOVED***
	re.putInt32(int32(in))
	return nil
***REMOVED***

func (re *realEncoder) putBool(in bool) ***REMOVED***
	if in ***REMOVED***
		re.putInt8(1)
		return
	***REMOVED***
	re.putInt8(0)
***REMOVED***

// collection

func (re *realEncoder) putRawBytes(in []byte) error ***REMOVED***
	copy(re.raw[re.off:], in)
	re.off += len(in)
	return nil
***REMOVED***

func (re *realEncoder) putBytes(in []byte) error ***REMOVED***
	if in == nil ***REMOVED***
		re.putInt32(-1)
		return nil
	***REMOVED***
	re.putInt32(int32(len(in)))
	return re.putRawBytes(in)
***REMOVED***

func (re *realEncoder) putVarintBytes(in []byte) error ***REMOVED***
	if in == nil ***REMOVED***
		re.putVarint(-1)
		return nil
	***REMOVED***
	re.putVarint(int64(len(in)))
	return re.putRawBytes(in)
***REMOVED***

func (re *realEncoder) putString(in string) error ***REMOVED***
	re.putInt16(int16(len(in)))
	copy(re.raw[re.off:], in)
	re.off += len(in)
	return nil
***REMOVED***

func (re *realEncoder) putNullableString(in *string) error ***REMOVED***
	if in == nil ***REMOVED***
		re.putInt16(-1)
		return nil
	***REMOVED***
	return re.putString(*in)
***REMOVED***

func (re *realEncoder) putStringArray(in []string) error ***REMOVED***
	err := re.putArrayLength(len(in))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range in ***REMOVED***
		if err := re.putString(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (re *realEncoder) putInt32Array(in []int32) error ***REMOVED***
	err := re.putArrayLength(len(in))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, val := range in ***REMOVED***
		re.putInt32(val)
	***REMOVED***
	return nil
***REMOVED***

func (re *realEncoder) putInt64Array(in []int64) error ***REMOVED***
	err := re.putArrayLength(len(in))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, val := range in ***REMOVED***
		re.putInt64(val)
	***REMOVED***
	return nil
***REMOVED***

func (re *realEncoder) offset() int ***REMOVED***
	return re.off
***REMOVED***

// stacks

func (re *realEncoder) push(in pushEncoder) ***REMOVED***
	in.saveOffset(re.off)
	re.off += in.reserveLength()
	re.stack = append(re.stack, in)
***REMOVED***

func (re *realEncoder) pop() error ***REMOVED***
	// this is go's ugly pop pattern (the inverse of append)
	in := re.stack[len(re.stack)-1]
	re.stack = re.stack[:len(re.stack)-1]

	return in.run(re.off, re.raw)
***REMOVED***

// we do record metrics during the real encoder pass
func (re *realEncoder) metricRegistry() metrics.Registry ***REMOVED***
	return re.registry
***REMOVED***
