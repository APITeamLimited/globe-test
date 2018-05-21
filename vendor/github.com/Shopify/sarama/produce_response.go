package sarama

import (
	"fmt"
	"time"
)

type ProduceResponseBlock struct ***REMOVED***
	Err    KError
	Offset int64
	// only provided if Version >= 2 and the broker is configured with `LogAppendTime`
	Timestamp time.Time
***REMOVED***

func (b *ProduceResponseBlock) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.Err = KError(tmp)

	b.Offset, err = pd.getInt64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if version >= 2 ***REMOVED***
		if millis, err := pd.getInt64(); err != nil ***REMOVED***
			return err
		***REMOVED*** else if millis != -1 ***REMOVED***
			b.Timestamp = time.Unix(millis/1000, (millis%1000)*int64(time.Millisecond))
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (b *ProduceResponseBlock) encode(pe packetEncoder, version int16) (err error) ***REMOVED***
	pe.putInt16(int16(b.Err))
	pe.putInt64(b.Offset)

	if version >= 2 ***REMOVED***
		timestamp := int64(-1)
		if !b.Timestamp.Before(time.Unix(0, 0)) ***REMOVED***
			timestamp = b.Timestamp.UnixNano() / int64(time.Millisecond)
		***REMOVED*** else if !b.Timestamp.IsZero() ***REMOVED***
			return PacketEncodingError***REMOVED***fmt.Sprintf("invalid timestamp (%v)", b.Timestamp)***REMOVED***
		***REMOVED***
		pe.putInt64(timestamp)
	***REMOVED***

	return nil
***REMOVED***

type ProduceResponse struct ***REMOVED***
	Blocks       map[string]map[int32]*ProduceResponseBlock
	Version      int16
	ThrottleTime time.Duration // only provided if Version >= 1
***REMOVED***

func (r *ProduceResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.Version = version

	numTopics, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Blocks = make(map[string]map[int32]*ProduceResponseBlock, numTopics)
	for i := 0; i < numTopics; i++ ***REMOVED***
		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		numBlocks, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.Blocks[name] = make(map[int32]*ProduceResponseBlock, numBlocks)

		for j := 0; j < numBlocks; j++ ***REMOVED***
			id, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			block := new(ProduceResponseBlock)
			err = block.decode(pd, version)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			r.Blocks[name][id] = block
		***REMOVED***
	***REMOVED***

	if r.Version >= 1 ***REMOVED***
		millis, err := pd.getInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.ThrottleTime = time.Duration(millis) * time.Millisecond
	***REMOVED***

	return nil
***REMOVED***

func (r *ProduceResponse) encode(pe packetEncoder) error ***REMOVED***
	err := pe.putArrayLength(len(r.Blocks))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range r.Blocks ***REMOVED***
		err = pe.putString(topic)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = pe.putArrayLength(len(partitions))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for id, prb := range partitions ***REMOVED***
			pe.putInt32(id)
			err = prb.encode(pe, r.Version)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if r.Version >= 1 ***REMOVED***
		pe.putInt32(int32(r.ThrottleTime / time.Millisecond))
	***REMOVED***
	return nil
***REMOVED***

func (r *ProduceResponse) key() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ProduceResponse) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *ProduceResponse) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_9_0_0
	case 2:
		return V0_10_0_0
	case 3:
		return V0_11_0_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *ProduceResponse) GetBlock(topic string, partition int32) *ProduceResponseBlock ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		return nil
	***REMOVED***

	if r.Blocks[topic] == nil ***REMOVED***
		return nil
	***REMOVED***

	return r.Blocks[topic][partition]
***REMOVED***

// Testing API

func (r *ProduceResponse) AddTopicPartition(topic string, partition int32, err KError) ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		r.Blocks = make(map[string]map[int32]*ProduceResponseBlock)
	***REMOVED***
	byTopic, ok := r.Blocks[topic]
	if !ok ***REMOVED***
		byTopic = make(map[int32]*ProduceResponseBlock)
		r.Blocks[topic] = byTopic
	***REMOVED***
	byTopic[partition] = &ProduceResponseBlock***REMOVED***Err: err***REMOVED***
***REMOVED***
