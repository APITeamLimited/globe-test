package sarama

type OffsetFetchRequest struct ***REMOVED***
	ConsumerGroup string
	Version       int16
	partitions    map[string][]int32
***REMOVED***

func (r *OffsetFetchRequest) encode(pe packetEncoder) (err error) ***REMOVED***
	if r.Version < 0 || r.Version > 1 ***REMOVED***
		return PacketEncodingError***REMOVED***"invalid or unsupported OffsetFetchRequest version field"***REMOVED***
	***REMOVED***

	if err = pe.putString(r.ConsumerGroup); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = pe.putArrayLength(len(r.partitions)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range r.partitions ***REMOVED***
		if err = pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err = pe.putInt32Array(partitions); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetFetchRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.Version = version
	if r.ConsumerGroup, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	partitionCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if partitionCount == 0 ***REMOVED***
		return nil
	***REMOVED***
	r.partitions = make(map[string][]int32)
	for i := 0; i < partitionCount; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		partitions, err := pd.getInt32Array()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.partitions[topic] = partitions
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetFetchRequest) key() int16 ***REMOVED***
	return 9
***REMOVED***

func (r *OffsetFetchRequest) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *OffsetFetchRequest) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_8_2_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *OffsetFetchRequest) AddPartition(topic string, partitionID int32) ***REMOVED***
	if r.partitions == nil ***REMOVED***
		r.partitions = make(map[string][]int32)
	***REMOVED***

	r.partitions[topic] = append(r.partitions[topic], partitionID)
***REMOVED***
