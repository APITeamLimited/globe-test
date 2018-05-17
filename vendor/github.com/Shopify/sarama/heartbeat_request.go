package sarama

type HeartbeatRequest struct ***REMOVED***
	GroupId      string
	GenerationId int32
	MemberId     string
***REMOVED***

func (r *HeartbeatRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(r.GroupId); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt32(r.GenerationId)

	if err := pe.putString(r.MemberId); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *HeartbeatRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if r.GroupId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if r.GenerationId, err = pd.getInt32(); err != nil ***REMOVED***
		return
	***REMOVED***
	if r.MemberId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	return nil
***REMOVED***

func (r *HeartbeatRequest) key() int16 ***REMOVED***
	return 12
***REMOVED***

func (r *HeartbeatRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *HeartbeatRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***
