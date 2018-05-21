package sarama

import "time"

type DeleteTopicsRequest struct ***REMOVED***
	Topics  []string
	Timeout time.Duration
***REMOVED***

func (d *DeleteTopicsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putStringArray(d.Topics); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt32(int32(d.Timeout / time.Millisecond))

	return nil
***REMOVED***

func (d *DeleteTopicsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if d.Topics, err = pd.getStringArray(); err != nil ***REMOVED***
		return err
	***REMOVED***
	timeout, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Timeout = time.Duration(timeout) * time.Millisecond
	return nil
***REMOVED***

func (d *DeleteTopicsRequest) key() int16 ***REMOVED***
	return 20
***REMOVED***

func (d *DeleteTopicsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *DeleteTopicsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_10_1_0
***REMOVED***
