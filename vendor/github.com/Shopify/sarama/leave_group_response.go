package sarama

type LeaveGroupResponse struct ***REMOVED***
	Err KError
***REMOVED***

func (r *LeaveGroupResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	return nil
***REMOVED***

func (r *LeaveGroupResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Err = KError(kerr)

	return nil
***REMOVED***

func (r *LeaveGroupResponse) key() int16 ***REMOVED***
	return 13
***REMOVED***

func (r *LeaveGroupResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *LeaveGroupResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
