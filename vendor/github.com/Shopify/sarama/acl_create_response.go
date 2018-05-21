package sarama

import "time"

type CreateAclsResponse struct ***REMOVED***
	ThrottleTime         time.Duration
	AclCreationResponses []*AclCreationResponse
***REMOVED***

func (c *CreateAclsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(c.ThrottleTime / time.Millisecond))

	if err := pe.putArrayLength(len(c.AclCreationResponses)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, aclCreationResponse := range c.AclCreationResponses ***REMOVED***
		if err := aclCreationResponse.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *CreateAclsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.AclCreationResponses = make([]*AclCreationResponse, n)
	for i := 0; i < n; i++ ***REMOVED***
		c.AclCreationResponses[i] = new(AclCreationResponse)
		if err := c.AclCreationResponses[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *CreateAclsResponse) key() int16 ***REMOVED***
	return 30
***REMOVED***

func (d *CreateAclsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *CreateAclsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***

type AclCreationResponse struct ***REMOVED***
	Err    KError
	ErrMsg *string
***REMOVED***

func (a *AclCreationResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(a.Err))

	if err := pe.putNullableString(a.ErrMsg); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (a *AclCreationResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.Err = KError(kerr)

	if a.ErrMsg, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
