package sarama

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/rcrowley/go-metrics"
)

type prepEncoder struct ***REMOVED***
	stack  []pushEncoder
	length int
***REMOVED***

// primitives

func (pe *prepEncoder) putInt8(in int8) ***REMOVED***
	pe.length++
***REMOVED***

func (pe *prepEncoder) putInt16(in int16) ***REMOVED***
	pe.length += 2
***REMOVED***

func (pe *prepEncoder) putInt32(in int32) ***REMOVED***
	pe.length += 4
***REMOVED***

func (pe *prepEncoder) putInt64(in int64) ***REMOVED***
	pe.length += 8
***REMOVED***

func (pe *prepEncoder) putVarint(in int64) ***REMOVED***
	var buf [binary.MaxVarintLen64]byte
	pe.length += binary.PutVarint(buf[:], in)
***REMOVED***

func (pe *prepEncoder) putArrayLength(in int) error ***REMOVED***
	if in > math.MaxInt32 ***REMOVED***
		return PacketEncodingError***REMOVED***fmt.Sprintf("array too long (%d)", in)***REMOVED***
	***REMOVED***
	pe.length += 4
	return nil
***REMOVED***

func (pe *prepEncoder) putBool(in bool) ***REMOVED***
	pe.length++
***REMOVED***

// arrays

func (pe *prepEncoder) putBytes(in []byte) error ***REMOVED***
	pe.length += 4
	if in == nil ***REMOVED***
		return nil
	***REMOVED***
	return pe.putRawBytes(in)
***REMOVED***

func (pe *prepEncoder) putVarintBytes(in []byte) error ***REMOVED***
	if in == nil ***REMOVED***
		pe.putVarint(-1)
		return nil
	***REMOVED***
	pe.putVarint(int64(len(in)))
	return pe.putRawBytes(in)
***REMOVED***

func (pe *prepEncoder) putRawBytes(in []byte) error ***REMOVED***
	if len(in) > math.MaxInt32 ***REMOVED***
		return PacketEncodingError***REMOVED***fmt.Sprintf("byteslice too long (%d)", len(in))***REMOVED***
	***REMOVED***
	pe.length += len(in)
	return nil
***REMOVED***

func (pe *prepEncoder) putNullableString(in *string) error ***REMOVED***
	if in == nil ***REMOVED***
		pe.length += 2
		return nil
	***REMOVED***
	return pe.putString(*in)
***REMOVED***

func (pe *prepEncoder) putString(in string) error ***REMOVED***
	pe.length += 2
	if len(in) > math.MaxInt16 ***REMOVED***
		return PacketEncodingError***REMOVED***fmt.Sprintf("string too long (%d)", len(in))***REMOVED***
	***REMOVED***
	pe.length += len(in)
	return nil
***REMOVED***

func (pe *prepEncoder) putStringArray(in []string) error ***REMOVED***
	err := pe.putArrayLength(len(in))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, str := range in ***REMOVED***
		if err := pe.putString(str); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (pe *prepEncoder) putInt32Array(in []int32) error ***REMOVED***
	err := pe.putArrayLength(len(in))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.length += 4 * len(in)
	return nil
***REMOVED***

func (pe *prepEncoder) putInt64Array(in []int64) error ***REMOVED***
	err := pe.putArrayLength(len(in))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.length += 8 * len(in)
	return nil
***REMOVED***

func (pe *prepEncoder) offset() int ***REMOVED***
	return pe.length
***REMOVED***

// stackable

func (pe *prepEncoder) push(in pushEncoder) ***REMOVED***
	in.saveOffset(pe.length)
	pe.length += in.reserveLength()
	pe.stack = append(pe.stack, in)
***REMOVED***

func (pe *prepEncoder) pop() error ***REMOVED***
	in := pe.stack[len(pe.stack)-1]
	pe.stack = pe.stack[:len(pe.stack)-1]
	if dpe, ok := in.(dynamicPushEncoder); ok ***REMOVED***
		pe.length += dpe.adjustLength(pe.length)
	***REMOVED***

	return nil
***REMOVED***

// we do not record metrics during the prep encoder pass
func (pe *prepEncoder) metricRegistry() metrics.Registry ***REMOVED***
	return nil
***REMOVED***
