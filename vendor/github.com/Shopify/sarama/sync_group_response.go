package sarama

type SyncGroupResponse struct ***REMOVED***
	Err              KError
	MemberAssignment []byte
***REMOVED***

func (r *SyncGroupResponse) GetMemberAssignment() (*ConsumerGroupMemberAssignment, error) ***REMOVED***
	assignment := new(ConsumerGroupMemberAssignment)
	err := decode(r.MemberAssignment, assignment)
	return assignment, err
***REMOVED***

func (r *SyncGroupResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	return pe.putBytes(r.MemberAssignment)
***REMOVED***

func (r *SyncGroupResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Err = KError(kerr)

	r.MemberAssignment, err = pd.getBytes()
	return
***REMOVED***

func (r *SyncGroupResponse) key() int16 ***REMOVED***
	return 14
***REMOVED***

func (r *SyncGroupResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *SyncGroupResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
