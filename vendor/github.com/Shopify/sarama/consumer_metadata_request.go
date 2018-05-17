package sarama

type ConsumerMetadataRequest struct ***REMOVED***
	ConsumerGroup string
***REMOVED***

func (r *ConsumerMetadataRequest) encode(pe packetEncoder) error ***REMOVED***
	return pe.putString(r.ConsumerGroup)
***REMOVED***

func (r *ConsumerMetadataRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.ConsumerGroup, err = pd.getString()
	return err
***REMOVED***

func (r *ConsumerMetadataRequest) key() int16 ***REMOVED***
	return 10
***REMOVED***

func (r *ConsumerMetadataRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ConsumerMetadataRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_8_2_0
***REMOVED***
