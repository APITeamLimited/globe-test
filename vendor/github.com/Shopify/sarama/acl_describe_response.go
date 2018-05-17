package sarama

import "time"

type DescribeAclsResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Err          KError
	ErrMsg       *string
	ResourceAcls []*ResourceAcls
***REMOVED***

func (d *DescribeAclsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(d.ThrottleTime / time.Millisecond))
	pe.putInt16(int16(d.Err))

	if err := pe.putNullableString(d.ErrMsg); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(d.ResourceAcls)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, resourceAcl := range d.ResourceAcls ***REMOVED***
		if err := resourceAcl.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *DescribeAclsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Err = KError(kerr)

	errmsg, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if errmsg != "" ***REMOVED***
		d.ErrMsg = &errmsg
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.ResourceAcls = make([]*ResourceAcls, n)

	for i := 0; i < n; i++ ***REMOVED***
		d.ResourceAcls[i] = new(ResourceAcls)
		if err := d.ResourceAcls[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *DescribeAclsResponse) key() int16 ***REMOVED***
	return 29
***REMOVED***

func (d *DescribeAclsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *DescribeAclsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
