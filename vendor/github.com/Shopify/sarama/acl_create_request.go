package sarama

type CreateAclsRequest struct ***REMOVED***
	AclCreations []*AclCreation
***REMOVED***

func (c *CreateAclsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(c.AclCreations)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, aclCreation := range c.AclCreations ***REMOVED***
		if err := aclCreation.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *CreateAclsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.AclCreations = make([]*AclCreation, n)

	for i := 0; i < n; i++ ***REMOVED***
		c.AclCreations[i] = new(AclCreation)
		if err := c.AclCreations[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *CreateAclsRequest) key() int16 ***REMOVED***
	return 30
***REMOVED***

func (d *CreateAclsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *CreateAclsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***

type AclCreation struct ***REMOVED***
	Resource
	Acl
***REMOVED***

func (a *AclCreation) encode(pe packetEncoder) error ***REMOVED***
	if err := a.Resource.encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := a.Acl.encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (a *AclCreation) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if err := a.Resource.decode(pd, version); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := a.Acl.decode(pd, version); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
