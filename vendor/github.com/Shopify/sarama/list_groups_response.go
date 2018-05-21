package sarama

type ListGroupsResponse struct ***REMOVED***
	Err    KError
	Groups map[string]string
***REMOVED***

func (r *ListGroupsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))

	if err := pe.putArrayLength(len(r.Groups)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for groupId, protocolType := range r.Groups ***REMOVED***
		if err := pe.putString(groupId); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putString(protocolType); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *ListGroupsResponse) decode(pd packetDecoder, version int16) error ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Err = KError(kerr)

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil
	***REMOVED***

	r.Groups = make(map[string]string)
	for i := 0; i < n; i++ ***REMOVED***
		groupId, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		protocolType, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.Groups[groupId] = protocolType
	***REMOVED***

	return nil
***REMOVED***

func (r *ListGroupsResponse) key() int16 ***REMOVED***
	return 16
***REMOVED***

func (r *ListGroupsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ListGroupsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
