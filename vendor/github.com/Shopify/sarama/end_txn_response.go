package sarama

import (
	"time"
)

type EndTxnResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Err          KError
***REMOVED***

func (e *EndTxnResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(e.ThrottleTime / time.Millisecond))
	pe.putInt16(int16(e.Err))
	return nil
***REMOVED***

func (e *EndTxnResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	e.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	e.Err = KError(kerr)

	return nil
***REMOVED***

func (e *EndTxnResponse) key() int16 ***REMOVED***
	return 25
***REMOVED***

func (e *EndTxnResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (e *EndTxnResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
