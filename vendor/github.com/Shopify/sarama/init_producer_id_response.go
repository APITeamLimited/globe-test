package sarama

import "time"

type InitProducerIDResponse struct ***REMOVED***
	ThrottleTime  time.Duration
	Err           KError
	ProducerID    int64
	ProducerEpoch int16
***REMOVED***

func (i *InitProducerIDResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(i.ThrottleTime / time.Millisecond))
	pe.putInt16(int16(i.Err))
	pe.putInt64(i.ProducerID)
	pe.putInt16(i.ProducerEpoch)

	return nil
***REMOVED***

func (i *InitProducerIDResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i.Err = KError(kerr)

	if i.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if i.ProducerEpoch, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (i *InitProducerIDResponse) key() int16 ***REMOVED***
	return 22
***REMOVED***

func (i *InitProducerIDResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (i *InitProducerIDResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
