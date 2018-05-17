package sarama

type SyncGroupRequest struct ***REMOVED***
	GroupId          string
	GenerationId     int32
	MemberId         string
	GroupAssignments map[string][]byte
***REMOVED***

func (r *SyncGroupRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(r.GroupId); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt32(r.GenerationId)

	if err := pe.putString(r.MemberId); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(r.GroupAssignments)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for memberId, memberAssignment := range r.GroupAssignments ***REMOVED***
		if err := pe.putString(memberId); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putBytes(memberAssignment); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *SyncGroupRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if r.GroupId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if r.GenerationId, err = pd.getInt32(); err != nil ***REMOVED***
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

	r.GroupAssignments = make(map[string][]byte)
	for i := 0; i < n; i++ ***REMOVED***
		memberId, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		memberAssignment, err := pd.getBytes()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.GroupAssignments[memberId] = memberAssignment
	***REMOVED***

	return nil
***REMOVED***

func (r *SyncGroupRequest) key() int16 ***REMOVED***
	return 14
***REMOVED***

func (r *SyncGroupRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *SyncGroupRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***

func (r *SyncGroupRequest) AddGroupAssignment(memberId string, memberAssignment []byte) ***REMOVED***
	if r.GroupAssignments == nil ***REMOVED***
		r.GroupAssignments = make(map[string][]byte)
	***REMOVED***

	r.GroupAssignments[memberId] = memberAssignment
***REMOVED***

func (r *SyncGroupRequest) AddGroupAssignmentMember(memberId string, memberAssignment *ConsumerGroupMemberAssignment) error ***REMOVED***
	bin, err := encode(memberAssignment, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.AddGroupAssignment(memberId, bin)
	return nil
***REMOVED***
