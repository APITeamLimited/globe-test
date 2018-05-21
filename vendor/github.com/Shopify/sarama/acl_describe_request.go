package sarama

type DescribeAclsRequest struct ***REMOVED***
	AclFilter
***REMOVED***

func (d *DescribeAclsRequest) encode(pe packetEncoder) error ***REMOVED***
	return d.AclFilter.encode(pe)
***REMOVED***

func (d *DescribeAclsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	return d.AclFilter.decode(pd, version)
***REMOVED***

func (d *DescribeAclsRequest) key() int16 ***REMOVED***
	return 29
***REMOVED***

func (d *DescribeAclsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *DescribeAclsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
