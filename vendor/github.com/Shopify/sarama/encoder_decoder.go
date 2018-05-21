package sarama

import (
	"fmt"

	"github.com/rcrowley/go-metrics"
)

// Encoder is the interface that wraps the basic Encode method.
// Anything implementing Encoder can be turned into bytes using Kafka's encoding rules.
type encoder interface ***REMOVED***
	encode(pe packetEncoder) error
***REMOVED***

// Encode takes an Encoder and turns it into bytes while potentially recording metrics.
func encode(e encoder, metricRegistry metrics.Registry) ([]byte, error) ***REMOVED***
	if e == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	var prepEnc prepEncoder
	var realEnc realEncoder

	err := e.encode(&prepEnc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if prepEnc.length < 0 || prepEnc.length > int(MaxRequestSize) ***REMOVED***
		return nil, PacketEncodingError***REMOVED***fmt.Sprintf("invalid request size (%d)", prepEnc.length)***REMOVED***
	***REMOVED***

	realEnc.raw = make([]byte, prepEnc.length)
	realEnc.registry = metricRegistry
	err = e.encode(&realEnc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return realEnc.raw, nil
***REMOVED***

// Decoder is the interface that wraps the basic Decode method.
// Anything implementing Decoder can be extracted from bytes using Kafka's encoding rules.
type decoder interface ***REMOVED***
	decode(pd packetDecoder) error
***REMOVED***

type versionedDecoder interface ***REMOVED***
	decode(pd packetDecoder, version int16) error
***REMOVED***

// Decode takes bytes and a Decoder and fills the fields of the decoder from the bytes,
// interpreted using Kafka's encoding rules.
func decode(buf []byte, in decoder) error ***REMOVED***
	if buf == nil ***REMOVED***
		return nil
	***REMOVED***

	helper := realDecoder***REMOVED***raw: buf***REMOVED***
	err := in.decode(&helper)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if helper.off != len(buf) ***REMOVED***
		return PacketDecodingError***REMOVED***"invalid length"***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func versionedDecode(buf []byte, in versionedDecoder, version int16) error ***REMOVED***
	if buf == nil ***REMOVED***
		return nil
	***REMOVED***

	helper := realDecoder***REMOVED***raw: buf***REMOVED***
	err := in.decode(&helper, version)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if helper.off != len(buf) ***REMOVED***
		return PacketDecodingError***REMOVED***"invalid length"***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
