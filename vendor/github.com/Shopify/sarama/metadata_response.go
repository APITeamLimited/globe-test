package sarama

type PartitionMetadata struct ***REMOVED***
	Err      KError
	ID       int32
	Leader   int32
	Replicas []int32
	Isr      []int32
***REMOVED***

func (pm *PartitionMetadata) decode(pd packetDecoder) (err error) ***REMOVED***
	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pm.Err = KError(tmp)

	pm.ID, err = pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pm.Leader, err = pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pm.Replicas, err = pd.getInt32Array()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pm.Isr, err = pd.getInt32Array()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (pm *PartitionMetadata) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt16(int16(pm.Err))
	pe.putInt32(pm.ID)
	pe.putInt32(pm.Leader)

	err = pe.putInt32Array(pm.Replicas)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = pe.putInt32Array(pm.Isr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

type TopicMetadata struct ***REMOVED***
	Err        KError
	Name       string
	Partitions []*PartitionMetadata
***REMOVED***

func (tm *TopicMetadata) decode(pd packetDecoder) (err error) ***REMOVED***
	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	tm.Err = KError(tmp)

	tm.Name, err = pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	tm.Partitions = make([]*PartitionMetadata, n)
	for i := 0; i < n; i++ ***REMOVED***
		tm.Partitions[i] = new(PartitionMetadata)
		err = tm.Partitions[i].decode(pd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (tm *TopicMetadata) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt16(int16(tm.Err))

	err = pe.putString(tm.Name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = pe.putArrayLength(len(tm.Partitions))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, pm := range tm.Partitions ***REMOVED***
		err = pm.encode(pe)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type MetadataResponse struct ***REMOVED***
	Brokers []*Broker
	Topics  []*TopicMetadata
***REMOVED***

func (r *MetadataResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Brokers = make([]*Broker, n)
	for i := 0; i < n; i++ ***REMOVED***
		r.Brokers[i] = new(Broker)
		err = r.Brokers[i].decode(pd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	n, err = pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Topics = make([]*TopicMetadata, n)
	for i := 0; i < n; i++ ***REMOVED***
		r.Topics[i] = new(TopicMetadata)
		err = r.Topics[i].decode(pd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *MetadataResponse) encode(pe packetEncoder) error ***REMOVED***
	err := pe.putArrayLength(len(r.Brokers))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, broker := range r.Brokers ***REMOVED***
		err = broker.encode(pe)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	err = pe.putArrayLength(len(r.Topics))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, tm := range r.Topics ***REMOVED***
		err = tm.encode(pe)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *MetadataResponse) key() int16 ***REMOVED***
	return 3
***REMOVED***

func (r *MetadataResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *MetadataResponse) requiredVersion() KafkaVersion ***REMOVED***
	return minVersion
***REMOVED***

// testing API

func (r *MetadataResponse) AddBroker(addr string, id int32) ***REMOVED***
	r.Brokers = append(r.Brokers, &Broker***REMOVED***id: id, addr: addr***REMOVED***)
***REMOVED***

func (r *MetadataResponse) AddTopic(topic string, err KError) *TopicMetadata ***REMOVED***
	var tmatch *TopicMetadata

	for _, tm := range r.Topics ***REMOVED***
		if tm.Name == topic ***REMOVED***
			tmatch = tm
			goto foundTopic
		***REMOVED***
	***REMOVED***

	tmatch = new(TopicMetadata)
	tmatch.Name = topic
	r.Topics = append(r.Topics, tmatch)

foundTopic:

	tmatch.Err = err
	return tmatch
***REMOVED***

func (r *MetadataResponse) AddTopicPartition(topic string, partition, brokerID int32, replicas, isr []int32, err KError) ***REMOVED***
	tmatch := r.AddTopic(topic, ErrNoError)
	var pmatch *PartitionMetadata

	for _, pm := range tmatch.Partitions ***REMOVED***
		if pm.ID == partition ***REMOVED***
			pmatch = pm
			goto foundPartition
		***REMOVED***
	***REMOVED***

	pmatch = new(PartitionMetadata)
	pmatch.ID = partition
	tmatch.Partitions = append(tmatch.Partitions, pmatch)

foundPartition:

	pmatch.Leader = brokerID
	pmatch.Replicas = replicas
	pmatch.Isr = isr
	pmatch.Err = err

***REMOVED***
