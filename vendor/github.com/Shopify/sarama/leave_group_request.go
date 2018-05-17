package sarama

type LeaveGroupRequest struct ***REMOVED***
	GroupId  string
	MemberId string
***REMOVED***

func (r *LeaveGroupRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(r.GroupId); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(r.MemberId); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *LeaveGroupRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if r.GroupId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if r.MemberId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	return nil
***REMOVED***

func (r *LeaveGroupRequest) key() int16 ***REMOVED***
	return 13
***REMOVED***

func (r *LeaveGroupRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *LeaveGroupRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
