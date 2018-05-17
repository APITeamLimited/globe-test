package sarama

type HeartbeatResponse struct ***REMOVED***
	Err KError
***REMOVED***

func (r *HeartbeatResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	return nil
***REMOVED***

func (r *HeartbeatResponse) decode(pd packetDecoder, version int16) error ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Err = KError(kerr)

	return nil
***REMOVED***

func (r *HeartbeatResponse) key() int16 ***REMOVED***
	return 12
***REMOVED***

func (r *HeartbeatResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *HeartbeatResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
