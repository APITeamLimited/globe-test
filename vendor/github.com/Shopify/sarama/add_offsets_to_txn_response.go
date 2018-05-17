package sarama

import (
	"time"
)

type AddOffsetsToTxnResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Err          KError
***REMOVED***

func (a *AddOffsetsToTxnResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(a.ThrottleTime / time.Millisecond))
	pe.putInt16(int16(a.Err))
	return nil
***REMOVED***

func (a *AddOffsetsToTxnResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.Err = KError(kerr)

	return nil
***REMOVED***

func (a *AddOffsetsToTxnResponse) key() int16 ***REMOVED***
	return 25
***REMOVED***

func (a *AddOffsetsToTxnResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *AddOffsetsToTxnResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
