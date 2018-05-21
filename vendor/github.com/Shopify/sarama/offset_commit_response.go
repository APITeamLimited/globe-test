package sarama

type OffsetCommitResponse struct ***REMOVED***
	Errors map[string]map[int32]KError
***REMOVED***

func (r *OffsetCommitResponse) AddError(topic string, partition int32, kerror KError) ***REMOVED***
	if r.Errors == nil ***REMOVED***
		r.Errors = make(map[string]map[int32]KError)
	***REMOVED***
	partitions := r.Errors[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]KError)
		r.Errors[topic] = partitions
	***REMOVED***
	partitions[partition] = kerror
***REMOVED***

func (r *OffsetCommitResponse) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(r.Errors)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range r.Errors ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putArrayLength(len(partitions)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for partition, kerror := range partitions ***REMOVED***
			pe.putInt32(partition)
			pe.putInt16(int16(kerror))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *OffsetCommitResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	numTopics, err := pd.getArrayLength()
	if err != nil || numTopics == 0 ***REMOVED***
		return err
	***REMOVED***

	r.Errors = make(map[string]map[int32]KError, numTopics)
	for i := 0; i < numTopics; i++ ***REMOVED***
		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		numErrors, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.Errors[name] = make(map[int32]KError, numErrors)

		for j := 0; j < numErrors; j++ ***REMOVED***
			id, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			tmp, err := pd.getInt16()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			r.Errors[name][id] = KError(tmp)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *OffsetCommitResponse) key() int16 ***REMOVED***
	return 8
***REMOVED***

func (r *OffsetCommitResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *OffsetCommitResponse) requiredVersion() KafkaVersion ***REMOVED***
	return minVersion
***REMOVED***
