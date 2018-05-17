package sarama

// ReceiveTime is a special value for the timestamp field of Offset Commit Requests which
// tells the broker to set the timestamp to the time at which the request was received.
// The timestamp is only used if message version 1 is used, which requires kafka 0.8.2.
const ReceiveTime int64 = -1

// GroupGenerationUndefined is a special value for the group generation field of
// Offset Commit Requests that should be used when a consumer group does not rely
// on Kafka for partition management.
const GroupGenerationUndefined = -1

type offsetCommitRequestBlock struct ***REMOVED***
	offset    int64
	timestamp int64
	metadata  string
***REMOVED***

func (b *offsetCommitRequestBlock) encode(pe packetEncoder, version int16) error ***REMOVED***
	pe.putInt64(b.offset)
	if version == 1 ***REMOVED***
		pe.putInt64(b.timestamp)
	***REMOVED*** else if b.timestamp != 0 ***REMOVED***
		Logger.Println("Non-zero timestamp specified for OffsetCommitRequest not v1, it will be ignored")
	***REMOVED***

	return pe.putString(b.metadata)
***REMOVED***

func (b *offsetCommitRequestBlock) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if b.offset, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if version == 1 ***REMOVED***
		if b.timestamp, err = pd.getInt64(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	b.metadata, err = pd.getString()
	return err
***REMOVED***

type OffsetCommitRequest struct ***REMOVED***
	ConsumerGroup           string
	ConsumerGroupGeneration int32  // v1 or later
	ConsumerID              string // v1 or later
	RetentionTime           int64  // v2 or later

	// Version can be:
	// - 0 (kafka 0.8.1 and later)
	// - 1 (kafka 0.8.2 and later)
	// - 2 (kafka 0.9.0 and later)
	Version int16
	blocks  map[string]map[int32]*offsetCommitRequestBlock
***REMOVED***

func (r *OffsetCommitRequest) encode(pe packetEncoder) error ***REMOVED***
	if r.Version < 0 || r.Version > 2 ***REMOVED***
		return PacketEncodingError***REMOVED***"invalid or unsupported OffsetCommitRequest version field"***REMOVED***
	***REMOVED***

	if err := pe.putString(r.ConsumerGroup); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.Version >= 1 ***REMOVED***
		pe.putInt32(r.ConsumerGroupGeneration)
		if err := pe.putString(r.ConsumerID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if r.ConsumerGroupGeneration != 0 ***REMOVED***
			Logger.Println("Non-zero ConsumerGroupGeneration specified for OffsetCommitRequest v0, it will be ignored")
		***REMOVED***
		if r.ConsumerID != "" ***REMOVED***
			Logger.Println("Non-empty ConsumerID specified for OffsetCommitRequest v0, it will be ignored")
		***REMOVED***
	***REMOVED***

	if r.Version >= 2 ***REMOVED***
		pe.putInt64(r.RetentionTime)
	***REMOVED*** else if r.RetentionTime != 0 ***REMOVED***
		Logger.Println("Non-zero RetentionTime specified for OffsetCommitRequest version <2, it will be ignored")
	***REMOVED***

	if err := pe.putArrayLength(len(r.blocks)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range r.blocks ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putArrayLength(len(partitions)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for partition, block := range partitions ***REMOVED***
			pe.putInt32(partition)
			if err := block.encode(pe, r.Version); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetCommitRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.Version = version

	if r.ConsumerGroup, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if r.Version >= 1 ***REMOVED***
		if r.ConsumerGroupGeneration, err = pd.getInt32(); err != nil ***REMOVED***
			return err
		***REMOVED***
		if r.ConsumerID, err = pd.getString(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if r.Version >= 2 ***REMOVED***
		if r.RetentionTime, err = pd.getInt64(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	topicCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if topicCount == 0 ***REMOVED***
		return nil
	***REMOVED***
	r.blocks = make(map[string]map[int32]*offsetCommitRequestBlock)
	for i := 0; i < topicCount; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		partitionCount, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.blocks[topic] = make(map[int32]*offsetCommitRequestBlock)
		for j := 0; j < partitionCount; j++ ***REMOVED***
			partition, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			block := &offsetCommitRequestBlock***REMOVED******REMOVED***
			if err := block.decode(pd, r.Version); err != nil ***REMOVED***
				return err
			***REMOVED***
			r.blocks[topic][partition] = block
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetCommitRequest) key() int16 ***REMOVED***
	return 8
***REMOVED***

func (r *OffsetCommitRequest) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *OffsetCommitRequest) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_8_2_0
	case 2:
		return V0_9_0_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *OffsetCommitRequest) AddBlock(topic string, partitionID int32, offset int64, timestamp int64, metadata string) ***REMOVED***
	if r.blocks == nil ***REMOVED***
		r.blocks = make(map[string]map[int32]*offsetCommitRequestBlock)
	***REMOVED***

	if r.blocks[topic] == nil ***REMOVED***
		r.blocks[topic] = make(map[int32]*offsetCommitRequestBlock)
	***REMOVED***

	r.blocks[topic][partitionID] = &offsetCommitRequestBlock***REMOVED***offset, timestamp, metadata***REMOVED***
***REMOVED***
