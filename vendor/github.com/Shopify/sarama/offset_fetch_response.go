package sarama

type OffsetFetchResponseBlock struct ***REMOVED***
	Offset   int64
	Metadata string
	Err      KError
***REMOVED***

func (b *OffsetFetchResponseBlock) decode(pd packetDecoder) (err error) ***REMOVED***
	b.Offset, err = pd.getInt64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b.Metadata, err = pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.Err = KError(tmp)

	return nil
***REMOVED***

func (b *OffsetFetchResponseBlock) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt64(b.Offset)

	err = pe.putString(b.Metadata)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt16(int16(b.Err))

	return nil
***REMOVED***

type OffsetFetchResponse struct ***REMOVED***
	Blocks map[string]map[int32]*OffsetFetchResponseBlock
***REMOVED***

func (r *OffsetFetchResponse) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(r.Blocks)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range r.Blocks ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putArrayLength(len(partitions)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for partition, block := range partitions ***REMOVED***
			pe.putInt32(partition)
			if err := block.encode(pe); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetFetchResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	numTopics, err := pd.getArrayLength()
	if err != nil || numTopics == 0 ***REMOVED***
		return err
	***REMOVED***

	r.Blocks = make(map[string]map[int32]*OffsetFetchResponseBlock, numTopics)
	for i := 0; i < numTopics; i++ ***REMOVED***
		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		numBlocks, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if numBlocks == 0 ***REMOVED***
			r.Blocks[name] = nil
			continue
		***REMOVED***
		r.Blocks[name] = make(map[int32]*OffsetFetchResponseBlock, numBlocks)

		for j := 0; j < numBlocks; j++ ***REMOVED***
			id, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			block := new(OffsetFetchResponseBlock)
			err = block.decode(pd)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			r.Blocks[name][id] = block
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *OffsetFetchResponse) key() int16 ***REMOVED***
	return 9
***REMOVED***

func (r *OffsetFetchResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *OffsetFetchResponse) requiredVersion() KafkaVersion ***REMOVED***
	return minVersion
***REMOVED***

func (r *OffsetFetchResponse) GetBlock(topic string, partition int32) *OffsetFetchResponseBlock ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		return nil
	***REMOVED***

	if r.Blocks[topic] == nil ***REMOVED***
		return nil
	***REMOVED***

	return r.Blocks[topic][partition]
***REMOVED***

func (r *OffsetFetchResponse) AddBlock(topic string, partition int32, block *OffsetFetchResponseBlock) ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		r.Blocks = make(map[string]map[int32]*OffsetFetchResponseBlock)
	***REMOVED***
	partitions := r.Blocks[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]*OffsetFetchResponseBlock)
		r.Blocks[topic] = partitions
	***REMOVED***
	partitions[partition] = block
***REMOVED***
