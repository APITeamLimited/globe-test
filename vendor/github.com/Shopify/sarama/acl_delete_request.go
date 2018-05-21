package sarama

type DeleteAclsRequest struct ***REMOVED***
	Filters []*AclFilter
***REMOVED***

func (d *DeleteAclsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(d.Filters)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, filter := range d.Filters ***REMOVED***
		if err := filter.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *DeleteAclsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	d.Filters = make([]*AclFilter, n)
	for i := 0; i < n; i++ ***REMOVED***
		d.Filters[i] = new(AclFilter)
		if err := d.Filters[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *DeleteAclsRequest) key() int16 ***REMOVED***
	return 31
***REMOVED***

func (d *DeleteAclsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *DeleteAclsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
