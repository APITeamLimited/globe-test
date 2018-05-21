package sarama

type ApiVersionsRequest struct ***REMOVED***
***REMOVED***

func (r *ApiVersionsRequest) encode(pe packetEncoder) error ***REMOVED***
	return nil
***REMOVED***

func (r *ApiVersionsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	return nil
***REMOVED***

func (r *ApiVersionsRequest) key() int16 ***REMOVED***
	return 18
***REMOVED***

func (r *ApiVersionsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ApiVersionsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_10_0_0
***REMOVED***
