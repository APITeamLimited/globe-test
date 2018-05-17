package sarama

type offsetRequestBlock struct ***REMOVED***
	time       int64
	maxOffsets int32 // Only used in version 0
***REMOVED***

func (b *offsetRequestBlock) encode(pe packetEncoder, version int16) error ***REMOVED***
	pe.putInt64(int64(b.time))
	if version == 0 ***REMOVED***
		pe.putInt32(b.maxOffsets)
	***REMOVED***

	return nil
***REMOVED***

func (b *offsetRequestBlock) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if b.time, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if version == 0 ***REMOVED***
		if b.maxOffsets, err = pd.getInt32(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type OffsetRequest struct ***REMOVED***
	Version int16
	blocks  map[string]map[int32]*offsetRequestBlock
***REMOVED***

func (r *OffsetRequest) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(-1) // replica ID is always -1 for clients
	err := pe.putArrayLength(len(r.blocks))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range r.blocks ***REMOVED***
		err = pe.putString(topic)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = pe.putArrayLength(len(partitions))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for partition, block := range partitions ***REMOVED***
			pe.putInt32(partition)
			if err = block.encode(pe, r.Version); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetRequest) decode(pd packetDecoder, version int16) error ***REMOVED***
	r.Version = version

	// Ignore replica ID
	if _, err := pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	blockCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if blockCount == 0 ***REMOVED***
		return nil
	***REMOVED***
	r.blocks = make(map[string]map[int32]*offsetRequestBlock)
	for i := 0; i < blockCount; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		partitionCount, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.blocks[topic] = make(map[int32]*offsetRequestBlock)
		for j := 0; j < partitionCount; j++ ***REMOVED***
			partition, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			block := &offsetRequestBlock***REMOVED******REMOVED***
			if err := block.decode(pd, version); err != nil ***REMOVED***
				return err
			***REMOVED***
			r.blocks[topic][partition] = block
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetRequest) key() int16 ***REMOVED***
	return 2
***REMOVED***

func (r *OffsetRequest) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *OffsetRequest) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_10_1_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *OffsetRequest) AddBlock(topic string, partitionID int32, time int64, maxOffsets int32) ***REMOVED***
	if r.blocks == nil ***REMOVED***
		r.blocks = make(map[string]map[int32]*offsetRequestBlock)
	***REMOVED***

	if r.blocks[topic] == nil ***REMOVED***
		r.blocks[topic] = make(map[int32]*offsetRequestBlock)
	***REMOVED***

	tmp := new(offsetRequestBlock)
	tmp.time = time
	if r.Version == 0 ***REMOVED***
		tmp.maxOffsets = maxOffsets
	***REMOVED***

	r.blocks[topic][partitionID] = tmp
***REMOVED***
