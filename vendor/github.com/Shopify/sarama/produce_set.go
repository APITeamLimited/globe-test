package sarama

import (
	"encoding/binary"
	"time"
)

type partitionSet struct ***REMOVED***
	msgs          []*ProducerMessage
	recordsToSend Records
	bufferBytes   int
***REMOVED***

type produceSet struct ***REMOVED***
	parent *asyncProducer
	msgs   map[string]map[int32]*partitionSet

	bufferBytes int
	bufferCount int
***REMOVED***

func newProduceSet(parent *asyncProducer) *produceSet ***REMOVED***
	return &produceSet***REMOVED***
		msgs:   make(map[string]map[int32]*partitionSet),
		parent: parent,
	***REMOVED***
***REMOVED***

func (ps *produceSet) add(msg *ProducerMessage) error ***REMOVED***
	var err error
	var key, val []byte

	if msg.Key != nil ***REMOVED***
		if key, err = msg.Key.Encode(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if msg.Value != nil ***REMOVED***
		if val, err = msg.Value.Encode(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	timestamp := msg.Timestamp
	if msg.Timestamp.IsZero() ***REMOVED***
		timestamp = time.Now()
	***REMOVED***

	partitions := ps.msgs[msg.Topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]*partitionSet)
		ps.msgs[msg.Topic] = partitions
	***REMOVED***

	var size int

	set := partitions[msg.Partition]
	if set == nil ***REMOVED***
		if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) ***REMOVED***
			batch := &RecordBatch***REMOVED***
				FirstTimestamp: timestamp,
				Version:        2,
				ProducerID:     -1, /* No producer id */
				Codec:          ps.parent.conf.Producer.Compression,
			***REMOVED***
			set = &partitionSet***REMOVED***recordsToSend: newDefaultRecords(batch)***REMOVED***
			size = recordBatchOverhead
		***REMOVED*** else ***REMOVED***
			set = &partitionSet***REMOVED***recordsToSend: newLegacyRecords(new(MessageSet))***REMOVED***
		***REMOVED***
		partitions[msg.Partition] = set
	***REMOVED***

	set.msgs = append(set.msgs, msg)
	if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) ***REMOVED***
		// We are being conservative here to avoid having to prep encode the record
		size += maximumRecordOverhead
		rec := &Record***REMOVED***
			Key:            key,
			Value:          val,
			TimestampDelta: timestamp.Sub(set.recordsToSend.recordBatch.FirstTimestamp),
		***REMOVED***
		size += len(key) + len(val)
		if len(msg.Headers) > 0 ***REMOVED***
			rec.Headers = make([]*RecordHeader, len(msg.Headers))
			for i := range msg.Headers ***REMOVED***
				rec.Headers[i] = &msg.Headers[i]
				size += len(rec.Headers[i].Key) + len(rec.Headers[i].Value) + 2*binary.MaxVarintLen32
			***REMOVED***
		***REMOVED***
		set.recordsToSend.recordBatch.addRecord(rec)
	***REMOVED*** else ***REMOVED***
		msgToSend := &Message***REMOVED***Codec: CompressionNone, Key: key, Value: val***REMOVED***
		if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) ***REMOVED***
			msgToSend.Timestamp = timestamp
			msgToSend.Version = 1
		***REMOVED***
		set.recordsToSend.msgSet.addMessage(msgToSend)
		size = producerMessageOverhead + len(key) + len(val)
	***REMOVED***

	set.bufferBytes += size
	ps.bufferBytes += size
	ps.bufferCount++

	return nil
***REMOVED***

func (ps *produceSet) buildRequest() *ProduceRequest ***REMOVED***
	req := &ProduceRequest***REMOVED***
		RequiredAcks: ps.parent.conf.Producer.RequiredAcks,
		Timeout:      int32(ps.parent.conf.Producer.Timeout / time.Millisecond),
	***REMOVED***
	if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) ***REMOVED***
		req.Version = 2
	***REMOVED***
	if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) ***REMOVED***
		req.Version = 3
	***REMOVED***

	for topic, partitionSet := range ps.msgs ***REMOVED***
		for partition, set := range partitionSet ***REMOVED***
			if req.Version >= 3 ***REMOVED***
				rb := set.recordsToSend.recordBatch
				if len(rb.Records) > 0 ***REMOVED***
					rb.LastOffsetDelta = int32(len(rb.Records) - 1)
					for i, record := range rb.Records ***REMOVED***
						record.OffsetDelta = int64(i)
					***REMOVED***
				***REMOVED***

				req.AddBatch(topic, partition, rb)
				continue
			***REMOVED***
			if ps.parent.conf.Producer.Compression == CompressionNone ***REMOVED***
				req.AddSet(topic, partition, set.recordsToSend.msgSet)
			***REMOVED*** else ***REMOVED***
				// When compression is enabled, the entire set for each partition is compressed
				// and sent as the payload of a single fake "message" with the appropriate codec
				// set and no key. When the server sees a message with a compression codec, it
				// decompresses the payload and treats the result as its message set.

				if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) ***REMOVED***
					// If our version is 0.10 or later, assign relative offsets
					// to the inner messages. This lets the broker avoid
					// recompressing the message set.
					// (See https://cwiki.apache.org/confluence/display/KAFKA/KIP-31+-+Move+to+relative+offsets+in+compressed+message+sets
					// for details on relative offsets.)
					for i, msg := range set.recordsToSend.msgSet.Messages ***REMOVED***
						msg.Offset = int64(i)
					***REMOVED***
				***REMOVED***
				payload, err := encode(set.recordsToSend.msgSet, ps.parent.conf.MetricRegistry)
				if err != nil ***REMOVED***
					Logger.Println(err) // if this happens, it's basically our fault.
					panic(err)
				***REMOVED***
				compMsg := &Message***REMOVED***
					Codec: ps.parent.conf.Producer.Compression,
					Key:   nil,
					Value: payload,
					Set:   set.recordsToSend.msgSet, // Provide the underlying message set for accurate metrics
				***REMOVED***
				if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) ***REMOVED***
					compMsg.Version = 1
					compMsg.Timestamp = set.recordsToSend.msgSet.Messages[0].Msg.Timestamp
				***REMOVED***
				req.AddMessage(topic, partition, compMsg)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return req
***REMOVED***

func (ps *produceSet) eachPartition(cb func(topic string, partition int32, msgs []*ProducerMessage)) ***REMOVED***
	for topic, partitionSet := range ps.msgs ***REMOVED***
		for partition, set := range partitionSet ***REMOVED***
			cb(topic, partition, set.msgs)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ps *produceSet) dropPartition(topic string, partition int32) []*ProducerMessage ***REMOVED***
	if ps.msgs[topic] == nil ***REMOVED***
		return nil
	***REMOVED***
	set := ps.msgs[topic][partition]
	if set == nil ***REMOVED***
		return nil
	***REMOVED***
	ps.bufferBytes -= set.bufferBytes
	ps.bufferCount -= len(set.msgs)
	delete(ps.msgs[topic], partition)
	return set.msgs
***REMOVED***

func (ps *produceSet) wouldOverflow(msg *ProducerMessage) bool ***REMOVED***
	version := 1
	if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) ***REMOVED***
		version = 2
	***REMOVED***

	switch ***REMOVED***
	// Would we overflow our maximum possible size-on-the-wire? 10KiB is arbitrary overhead for safety.
	case ps.bufferBytes+msg.byteSize(version) >= int(MaxRequestSize-(10*1024)):
		return true
	// Would we overflow the size-limit of a compressed message-batch for this partition?
	case ps.parent.conf.Producer.Compression != CompressionNone &&
		ps.msgs[msg.Topic] != nil && ps.msgs[msg.Topic][msg.Partition] != nil &&
		ps.msgs[msg.Topic][msg.Partition].bufferBytes+msg.byteSize(version) >= ps.parent.conf.Producer.MaxMessageBytes:
		return true
	// Would we overflow simply in number of messages?
	case ps.parent.conf.Producer.Flush.MaxMessages > 0 && ps.bufferCount >= ps.parent.conf.Producer.Flush.MaxMessages:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (ps *produceSet) readyToFlush() bool ***REMOVED***
	switch ***REMOVED***
	// If we don't have any messages, nothing else matters
	case ps.empty():
		return false
	// If all three config values are 0, we always flush as-fast-as-possible
	case ps.parent.conf.Producer.Flush.Frequency == 0 && ps.parent.conf.Producer.Flush.Bytes == 0 && ps.parent.conf.Producer.Flush.Messages == 0:
		return true
	// If we've passed the message trigger-point
	case ps.parent.conf.Producer.Flush.Messages > 0 && ps.bufferCount >= ps.parent.conf.Producer.Flush.Messages:
		return true
	// If we've passed the byte trigger-point
	case ps.parent.conf.Producer.Flush.Bytes > 0 && ps.bufferBytes >= ps.parent.conf.Producer.Flush.Bytes:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (ps *produceSet) empty() bool ***REMOVED***
	return ps.bufferCount == 0
***REMOVED***
