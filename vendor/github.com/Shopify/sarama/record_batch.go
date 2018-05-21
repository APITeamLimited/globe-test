package sarama

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/eapache/go-xerial-snappy"
	"github.com/pierrec/lz4"
)

const recordBatchOverhead = 49

type recordsArray []*Record

func (e recordsArray) encode(pe packetEncoder) error ***REMOVED***
	for _, r := range e ***REMOVED***
		if err := r.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (e recordsArray) decode(pd packetDecoder) error ***REMOVED***
	for i := range e ***REMOVED***
		rec := &Record***REMOVED******REMOVED***
		if err := rec.decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
		e[i] = rec
	***REMOVED***
	return nil
***REMOVED***

type RecordBatch struct ***REMOVED***
	FirstOffset           int64
	PartitionLeaderEpoch  int32
	Version               int8
	Codec                 CompressionCodec
	Control               bool
	LastOffsetDelta       int32
	FirstTimestamp        time.Time
	MaxTimestamp          time.Time
	ProducerID            int64
	ProducerEpoch         int16
	FirstSequence         int32
	Records               []*Record
	PartialTrailingRecord bool

	compressedRecords []byte
	recordsLen        int // uncompressed records size
***REMOVED***

func (b *RecordBatch) encode(pe packetEncoder) error ***REMOVED***
	if b.Version != 2 ***REMOVED***
		return PacketEncodingError***REMOVED***fmt.Sprintf("unsupported compression codec (%d)", b.Codec)***REMOVED***
	***REMOVED***
	pe.putInt64(b.FirstOffset)
	pe.push(&lengthField***REMOVED******REMOVED***)
	pe.putInt32(b.PartitionLeaderEpoch)
	pe.putInt8(b.Version)
	pe.push(newCRC32Field(crcCastagnoli))
	pe.putInt16(b.computeAttributes())
	pe.putInt32(b.LastOffsetDelta)

	if err := (Timestamp***REMOVED***&b.FirstTimestamp***REMOVED***).encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := (Timestamp***REMOVED***&b.MaxTimestamp***REMOVED***).encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt64(b.ProducerID)
	pe.putInt16(b.ProducerEpoch)
	pe.putInt32(b.FirstSequence)

	if err := pe.putArrayLength(len(b.Records)); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.compressedRecords == nil ***REMOVED***
		if err := b.encodeRecords(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := pe.putRawBytes(b.compressedRecords); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.pop(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return pe.pop()
***REMOVED***

func (b *RecordBatch) decode(pd packetDecoder) (err error) ***REMOVED***
	if b.FirstOffset, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***

	batchLen, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.PartitionLeaderEpoch, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.Version, err = pd.getInt8(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = pd.push(&crc32Field***REMOVED***polynomial: crcCastagnoli***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	attributes, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.Codec = CompressionCodec(int8(attributes) & compressionCodecMask)
	b.Control = attributes&controlMask == controlMask

	if b.LastOffsetDelta, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = (Timestamp***REMOVED***&b.FirstTimestamp***REMOVED***).decode(pd); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = (Timestamp***REMOVED***&b.MaxTimestamp***REMOVED***).decode(pd); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.ProducerEpoch, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.FirstSequence, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***

	numRecs, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if numRecs >= 0 ***REMOVED***
		b.Records = make([]*Record, numRecs)
	***REMOVED***

	bufSize := int(batchLen) - recordBatchOverhead
	recBuffer, err := pd.getRawBytes(bufSize)
	if err != nil ***REMOVED***
		if err == ErrInsufficientData ***REMOVED***
			b.PartialTrailingRecord = true
			b.Records = nil
			return nil
		***REMOVED***
		return err
	***REMOVED***

	if err = pd.pop(); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch b.Codec ***REMOVED***
	case CompressionNone:
	case CompressionGZIP:
		reader, err := gzip.NewReader(bytes.NewReader(recBuffer))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if recBuffer, err = ioutil.ReadAll(reader); err != nil ***REMOVED***
			return err
		***REMOVED***
	case CompressionSnappy:
		if recBuffer, err = snappy.Decode(recBuffer); err != nil ***REMOVED***
			return err
		***REMOVED***
	case CompressionLZ4:
		reader := lz4.NewReader(bytes.NewReader(recBuffer))
		if recBuffer, err = ioutil.ReadAll(reader); err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		return PacketDecodingError***REMOVED***fmt.Sprintf("invalid compression specified (%d)", b.Codec)***REMOVED***
	***REMOVED***

	b.recordsLen = len(recBuffer)
	err = decode(recBuffer, recordsArray(b.Records))
	if err == ErrInsufficientData ***REMOVED***
		b.PartialTrailingRecord = true
		b.Records = nil
		return nil
	***REMOVED***
	return err
***REMOVED***

func (b *RecordBatch) encodeRecords(pe packetEncoder) error ***REMOVED***
	var raw []byte
	var err error
	if raw, err = encode(recordsArray(b.Records), pe.metricRegistry()); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.recordsLen = len(raw)

	switch b.Codec ***REMOVED***
	case CompressionNone:
		b.compressedRecords = raw
	case CompressionGZIP:
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write(raw); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := writer.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***
		b.compressedRecords = buf.Bytes()
	case CompressionSnappy:
		b.compressedRecords = snappy.Encode(raw)
	case CompressionLZ4:
		var buf bytes.Buffer
		writer := lz4.NewWriter(&buf)
		if _, err := writer.Write(raw); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := writer.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***
		b.compressedRecords = buf.Bytes()
	default:
		return PacketEncodingError***REMOVED***fmt.Sprintf("unsupported compression codec (%d)", b.Codec)***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (b *RecordBatch) computeAttributes() int16 ***REMOVED***
	attr := int16(b.Codec) & int16(compressionCodecMask)
	if b.Control ***REMOVED***
		attr |= controlMask
	***REMOVED***
	return attr
***REMOVED***

func (b *RecordBatch) addRecord(r *Record) ***REMOVED***
	b.Records = append(b.Records, r)
***REMOVED***
