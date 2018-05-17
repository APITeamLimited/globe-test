package sarama

type MetadataRequest struct ***REMOVED***
	Topics []string
***REMOVED***

func (r *MetadataRequest) encode(pe packetEncoder) error ***REMOVED***
	err := pe.putArrayLength(len(r.Topics))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for i := range r.Topics ***REMOVED***
		err = pe.putString(r.Topics[i])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *MetadataRequest) decode(pd packetDecoder, version int16) error ***REMOVED***
	topicCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if topicCount == 0 ***REMOVED***
		return nil
	***REMOVED***

	r.Topics = make([]string, topicCount)
	for i := range r.Topics ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Topics[i] = topic
	***REMOVED***
	return nil
***REMOVED***

func (r *MetadataRequest) key() int16 ***REMOVED***
	return 3
***REMOVED***

func (r *MetadataRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *MetadataRequest) requiredVersion() KafkaVersion ***REMOVED***
	return minVersion
***REMOVED***
