package sarama

type fetchRequestBlock struct ***REMOVED***
	fetchOffset int64
	maxBytes    int32
***REMOVED***

func (b *fetchRequestBlock) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt64(b.fetchOffset)
	pe.putInt32(b.maxBytes)
	return nil
***REMOVED***

func (b *fetchRequestBlock) decode(pd packetDecoder) (err error) ***REMOVED***
	if b.fetchOffset, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if b.maxBytes, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// FetchRequest (API key 1) will fetch Kafka messages. Version 3 introduced the MaxBytes field. See
// https://issues.apache.org/jira/browse/KAFKA-2063 for a discussion of the issues leading up to that.  The KIP is at
// https://cwiki.apache.org/confluence/display/KAFKA/KIP-74%3A+Add+Fetch+Response+Size+Limit+in+Bytes
type FetchRequest struct ***REMOVED***
	MaxWaitTime int32
	MinBytes    int32
	MaxBytes    int32
	Version     int16
	Isolation   IsolationLevel
	blocks      map[string]map[int32]*fetchRequestBlock
***REMOVED***

type IsolationLevel int8

const (
	ReadUncommitted IsolationLevel = 0
	ReadCommitted   IsolationLevel = 1
)

func (r *FetchRequest) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt32(-1) // replica ID is always -1 for clients
	pe.putInt32(r.MaxWaitTime)
	pe.putInt32(r.MinBytes)
	if r.Version >= 3 ***REMOVED***
		pe.putInt32(r.MaxBytes)
	***REMOVED***
	if r.Version >= 4 ***REMOVED***
		pe.putInt8(int8(r.Isolation))
	***REMOVED***
	err = pe.putArrayLength(len(r.blocks))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, blocks := range r.blocks ***REMOVED***
		err = pe.putString(topic)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = pe.putArrayLength(len(blocks))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for partition, block := range blocks ***REMOVED***
			pe.putInt32(partition)
			err = block.encode(pe)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *FetchRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.Version = version
	if _, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.MaxWaitTime, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.MinBytes, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.Version >= 3 ***REMOVED***
		if r.MaxBytes, err = pd.getInt32(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if r.Version >= 4 ***REMOVED***
		isolation, err := pd.getInt8()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Isolation = IsolationLevel(isolation)
	***REMOVED***
	topicCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if topicCount == 0 ***REMOVED***
		return nil
	***REMOVED***
	r.blocks = make(map[string]map[int32]*fetchRequestBlock)
	for i := 0; i < topicCount; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		partitionCount, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.blocks[topic] = make(map[int32]*fetchRequestBlock)
		for j := 0; j < partitionCount; j++ ***REMOVED***
			partition, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			fetchBlock := &fetchRequestBlock***REMOVED******REMOVED***
			if err = fetchBlock.decode(pd); err != nil ***REMOVED***
				return err
			***REMOVED***
			r.blocks[topic][partition] = fetchBlock
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *FetchRequest) key() int16 ***REMOVED***
	return 1
***REMOVED***

func (r *FetchRequest) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *FetchRequest) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_9_0_0
	case 2:
		return V0_10_0_0
	case 3:
		return V0_10_1_0
	case 4:
		return V0_11_0_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *FetchRequest) AddBlock(topic string, partitionID int32, fetchOffset int64, maxBytes int32) ***REMOVED***
	if r.blocks == nil ***REMOVED***
		r.blocks = make(map[string]map[int32]*fetchRequestBlock)
	***REMOVED***

	if r.blocks[topic] == nil ***REMOVED***
		r.blocks[topic] = make(map[int32]*fetchRequestBlock)
	***REMOVED***

	tmp := new(fetchRequestBlock)
	tmp.maxBytes = maxBytes
	tmp.fetchOffset = fetchOffset

	r.blocks[topic][partitionID] = tmp
***REMOVED***
