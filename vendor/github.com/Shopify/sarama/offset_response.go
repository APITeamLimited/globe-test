package sarama

type OffsetResponseBlock struct ***REMOVED***
	Err       KError
	Offsets   []int64 // Version 0
	Offset    int64   // Version 1
	Timestamp int64   // Version 1
***REMOVED***

func (b *OffsetResponseBlock) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.Err = KError(tmp)

	if version == 0 ***REMOVED***
		b.Offsets, err = pd.getInt64Array()

		return err
	***REMOVED***

	b.Timestamp, err = pd.getInt64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b.Offset, err = pd.getInt64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// For backwards compatibility put the offset in the offsets array too
	b.Offsets = []int64***REMOVED***b.Offset***REMOVED***

	return nil
***REMOVED***

func (b *OffsetResponseBlock) encode(pe packetEncoder, version int16) (err error) ***REMOVED***
	pe.putInt16(int16(b.Err))

	if version == 0 ***REMOVED***
		return pe.putInt64Array(b.Offsets)
	***REMOVED***

	pe.putInt64(b.Timestamp)
	pe.putInt64(b.Offset)

	return nil
***REMOVED***

type OffsetResponse struct ***REMOVED***
	Version int16
	Blocks  map[string]map[int32]*OffsetResponseBlock
***REMOVED***

func (r *OffsetResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	numTopics, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Blocks = make(map[string]map[int32]*OffsetResponseBlock, numTopics)
	for i := 0; i < numTopics; i++ ***REMOVED***
		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		numBlocks, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.Blocks[name] = make(map[int32]*OffsetResponseBlock, numBlocks)

		for j := 0; j < numBlocks; j++ ***REMOVED***
			id, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			block := new(OffsetResponseBlock)
			err = block.decode(pd, version)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			r.Blocks[name][id] = block
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *OffsetResponse) GetBlock(topic string, partition int32) *OffsetResponseBlock ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		return nil
	***REMOVED***

	if r.Blocks[topic] == nil ***REMOVED***
		return nil
	***REMOVED***

	return r.Blocks[topic][partition]
***REMOVED***

/*
// [0 0 0 1 ntopics
0 8 109 121 95 116 111 112 105 99 topic
0 0 0 1 npartitions
0 0 0 0 id
0 0

0 0 0 1 0 0 0 0
0 1 1 1 0 0 0 1
0 8 109 121 95 116 111 112
105 99 0 0 0 1 0 0
0 0 0 0 0 0 0 1
0 0 0 0 0 1 1 1] <nil>

*/
func (r *OffsetResponse) encode(pe packetEncoder) (err error) ***REMOVED***
	if err = pe.putArrayLength(len(r.Blocks)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, partitions := range r.Blocks ***REMOVED***
		if err = pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err = pe.putArrayLength(len(partitions)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for partition, block := range partitions ***REMOVED***
			pe.putInt32(partition)
			if err = block.encode(pe, r.version()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *OffsetResponse) key() int16 ***REMOVED***
	return 2
***REMOVED***

func (r *OffsetResponse) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *OffsetResponse) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_10_1_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

// testing API

func (r *OffsetResponse) AddTopicPartition(topic string, partition int32, offset int64) ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		r.Blocks = make(map[string]map[int32]*OffsetResponseBlock)
	***REMOVED***
	byTopic, ok := r.Blocks[topic]
	if !ok ***REMOVED***
		byTopic = make(map[int32]*OffsetResponseBlock)
		r.Blocks[topic] = byTopic
	***REMOVED***
	byTopic[partition] = &OffsetResponseBlock***REMOVED***Offsets: []int64***REMOVED***offset***REMOVED***, Offset: offset***REMOVED***
***REMOVED***
