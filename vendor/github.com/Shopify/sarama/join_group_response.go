package sarama

type JoinGroupResponse struct ***REMOVED***
	Err           KError
	GenerationId  int32
	GroupProtocol string
	LeaderId      string
	MemberId      string
	Members       map[string][]byte
***REMOVED***

func (r *JoinGroupResponse) GetMembers() (map[string]ConsumerGroupMemberMetadata, error) ***REMOVED***
	members := make(map[string]ConsumerGroupMemberMetadata, len(r.Members))
	for id, bin := range r.Members ***REMOVED***
		meta := new(ConsumerGroupMemberMetadata)
		if err := decode(bin, meta); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		members[id] = *meta
	***REMOVED***
	return members, nil
***REMOVED***

func (r *JoinGroupResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	pe.putInt32(r.GenerationId)

	if err := pe.putString(r.GroupProtocol); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(r.LeaderId); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(r.MemberId); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(r.Members)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for memberId, memberMetadata := range r.Members ***REMOVED***
		if err := pe.putString(memberId); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := pe.putBytes(memberMetadata); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *JoinGroupResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Err = KError(kerr)

	if r.GenerationId, err = pd.getInt32(); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.GroupProtocol, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.LeaderId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.MemberId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil
	***REMOVED***

	r.Members = make(map[string][]byte)
	for i := 0; i < n; i++ ***REMOVED***
		memberId, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		memberMetadata, err := pd.getBytes()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.Members[memberId] = memberMetadata
	***REMOVED***

	return nil
***REMOVED***

func (r *JoinGroupResponse) key() int16 ***REMOVED***
	return 11
***REMOVED***

func (r *JoinGroupResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *JoinGroupResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
