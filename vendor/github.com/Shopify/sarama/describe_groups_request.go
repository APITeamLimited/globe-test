package sarama

type DescribeGroupsRequest struct ***REMOVED***
	Groups []string
***REMOVED***

func (r *DescribeGroupsRequest) encode(pe packetEncoder) error ***REMOVED***
	return pe.putStringArray(r.Groups)
***REMOVED***

func (r *DescribeGroupsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.Groups, err = pd.getStringArray()
	return
***REMOVED***

func (r *DescribeGroupsRequest) key() int16 ***REMOVED***
	return 15
***REMOVED***

func (r *DescribeGroupsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *DescribeGroupsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***

func (r *DescribeGroupsRequest) AddGroup(group string) ***REMOVED***
	r.Groups = append(r.Groups, group)
***REMOVED***
