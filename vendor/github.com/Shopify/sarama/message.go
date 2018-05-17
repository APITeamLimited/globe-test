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

// CompressionCodec represents the various compression codecs recognized by Kafka in messages.
type CompressionCodec int8

// only the last two bits are really used
const compressionCodecMask int8 = 0x03

const (
	CompressionNone   CompressionCodec = 0
	CompressionGZIP   CompressionCodec = 1
	CompressionSnappy CompressionCodec = 2
	CompressionLZ4    CompressionCodec = 3
)

type Message struct ***REMOVED***
	Codec     CompressionCodec // codec used to compress the message contents
	Key       []byte           // the message key, may be nil
	Value     []byte           // the message contents
	Set       *MessageSet      // the message set a message might wrap
	Version   int8             // v1 requires Kafka 0.10
	Timestamp time.Time        // the timestamp of the message (version 1+ only)

	compressedCache []byte
	compressedSize  int // used for computing the compression ratio metrics
***REMOVED***

func (m *Message) encode(pe packetEncoder) error ***REMOVED***
	pe.push(newCRC32Field(crcIEEE))

	pe.putInt8(m.Version)

	attributes := int8(m.Codec) & compressionCodecMask
	pe.putInt8(attributes)

	if m.Version >= 1 ***REMOVED***
		if err := (Timestamp***REMOVED***&m.Timestamp***REMOVED***).encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	err := pe.putBytes(m.Key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var payload []byte

	if m.compressedCache != nil ***REMOVED***
		payload = m.compressedCache
		m.compressedCache = nil
	***REMOVED*** else if m.Value != nil ***REMOVED***
		switch m.Codec ***REMOVED***
		case CompressionNone:
			payload = m.Value
		case CompressionGZIP:
			var buf bytes.Buffer
			writer := gzip.NewWriter(&buf)
			if _, err = writer.Write(m.Value); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err = writer.Close(); err != nil ***REMOVED***
				return err
			***REMOVED***
			m.compressedCache = buf.Bytes()
			payload = m.compressedCache
		case CompressionSnappy:
			tmp := snappy.Encode(m.Value)
			m.compressedCache = tmp
			payload = m.compressedCache
		case CompressionLZ4:
			var buf bytes.Buffer
			writer := lz4.NewWriter(&buf)
			if _, err = writer.Write(m.Value); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err = writer.Close(); err != nil ***REMOVED***
				return err
			***REMOVED***
			m.compressedCache = buf.Bytes()
			payload = m.compressedCache

		default:
			return PacketEncodingError***REMOVED***fmt.Sprintf("unsupported compression codec (%d)", m.Codec)***REMOVED***
		***REMOVED***
		// Keep in mind the compressed payload size for metric gathering
		m.compressedSize = len(payload)
	***REMOVED***

	if err = pe.putBytes(payload); err != nil ***REMOVED***
		return err
	***REMOVED***

	return pe.pop()
***REMOVED***

func (m *Message) decode(pd packetDecoder) (err error) ***REMOVED***
	err = pd.push(newCRC32Field(crcIEEE))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	m.Version, err = pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if m.Version > 1 ***REMOVED***
		return PacketDecodingError***REMOVED***fmt.Sprintf("unknown magic byte (%v)", m.Version)***REMOVED***
	***REMOVED***

	attribute, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m.Codec = CompressionCodec(attribute & compressionCodecMask)

	if m.Version == 1 ***REMOVED***
		if err := (Timestamp***REMOVED***&m.Timestamp***REMOVED***).decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	m.Key, err = pd.getBytes()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	m.Value, err = pd.getBytes()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Required for deep equal assertion during tests but might be useful
	// for future metrics about the compression ratio in fetch requests
	m.compressedSize = len(m.Value)

	switch m.Codec ***REMOVED***
	case CompressionNone:
		// nothing to do
	case CompressionGZIP:
		if m.Value == nil ***REMOVED***
			break
		***REMOVED***
		reader, err := gzip.NewReader(bytes.NewReader(m.Value))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if m.Value, err = ioutil.ReadAll(reader); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := m.decodeSet(); err != nil ***REMOVED***
			return err
		***REMOVED***
	case CompressionSnappy:
		if m.Value == nil ***REMOVED***
			break
		***REMOVED***
		if m.Value, err = snappy.Decode(m.Value); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := m.decodeSet(); err != nil ***REMOVED***
			return err
		***REMOVED***
	case CompressionLZ4:
		if m.Value == nil ***REMOVED***
			break
		***REMOVED***
		reader := lz4.NewReader(bytes.NewReader(m.Value))
		if m.Value, err = ioutil.ReadAll(reader); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := m.decodeSet(); err != nil ***REMOVED***
			return err
		***REMOVED***

	default:
		return PacketDecodingError***REMOVED***fmt.Sprintf("invalid compression specified (%d)", m.Codec)***REMOVED***
	***REMOVED***

	return pd.pop()
***REMOVED***

// decodes a message set from a previousy encoded bulk-message
func (m *Message) decodeSet() (err error) ***REMOVED***
	pd := realDecoder***REMOVED***raw: m.Value***REMOVED***
	m.Set = &MessageSet***REMOVED******REMOVED***
	return m.Set.decode(&pd)
***REMOVED***
