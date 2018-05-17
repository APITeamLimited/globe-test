package sarama

import (
	"encoding/binary"
	"time"
)

const (
	controlMask           = 0x20
	maximumRecordOverhead = 5*binary.MaxVarintLen32 + binary.MaxVarintLen64 + 1
)

type RecordHeader struct ***REMOVED***
	Key   []byte
	Value []byte
***REMOVED***

func (h *RecordHeader) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putVarintBytes(h.Key); err != nil ***REMOVED***
		return err
	***REMOVED***
	return pe.putVarintBytes(h.Value)
***REMOVED***

func (h *RecordHeader) decode(pd packetDecoder) (err error) ***REMOVED***
	if h.Key, err = pd.getVarintBytes(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if h.Value, err = pd.getVarintBytes(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

type Record struct ***REMOVED***
	Attributes     int8
	TimestampDelta time.Duration
	OffsetDelta    int64
	Key            []byte
	Value          []byte
	Headers        []*RecordHeader

	length varintLengthField
***REMOVED***

func (r *Record) encode(pe packetEncoder) error ***REMOVED***
	pe.push(&r.length)
	pe.putInt8(r.Attributes)
	pe.putVarint(int64(r.TimestampDelta / time.Millisecond))
	pe.putVarint(r.OffsetDelta)
	if err := pe.putVarintBytes(r.Key); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putVarintBytes(r.Value); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putVarint(int64(len(r.Headers)))

	for _, h := range r.Headers ***REMOVED***
		if err := h.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return pe.pop()
***REMOVED***

func (r *Record) decode(pd packetDecoder) (err error) ***REMOVED***
	if err = pd.push(&r.length); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.Attributes, err = pd.getInt8(); err != nil ***REMOVED***
		return err
	***REMOVED***

	timestamp, err := pd.getVarint()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.TimestampDelta = time.Duration(timestamp) * time.Millisecond

	if r.OffsetDelta, err = pd.getVarint(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.Key, err = pd.getVarintBytes(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.Value, err = pd.getVarintBytes(); err != nil ***REMOVED***
		return err
	***REMOVED***

	numHeaders, err := pd.getVarint()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if numHeaders >= 0 ***REMOVED***
		r.Headers = make([]*RecordHeader, numHeaders)
	***REMOVED***
	for i := int64(0); i < numHeaders; i++ ***REMOVED***
		hdr := new(RecordHeader)
		if err := hdr.decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Headers[i] = hdr
	***REMOVED***

	return pd.pop()
***REMOVED***
