package sarama

type ListGroupsRequest struct ***REMOVED***
***REMOVED***

func (r *ListGroupsRequest) encode(pe packetEncoder) error ***REMOVED***
	return nil
***REMOVED***

func (r *ListGroupsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	return nil
***REMOVED***

func (r *ListGroupsRequest) key() int16 ***REMOVED***
	return 16
***REMOVED***

func (r *ListGroupsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ListGroupsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
