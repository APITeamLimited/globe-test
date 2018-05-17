package sarama

import "time"

type DeleteTopicsResponse struct ***REMOVED***
	Version         int16
	ThrottleTime    time.Duration
	TopicErrorCodes map[string]KError
***REMOVED***

func (d *DeleteTopicsResponse) encode(pe packetEncoder) error ***REMOVED***
	if d.Version >= 1 ***REMOVED***
		pe.putInt32(int32(d.ThrottleTime / time.Millisecond))
	***REMOVED***

	if err := pe.putArrayLength(len(d.TopicErrorCodes)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, errorCode := range d.TopicErrorCodes ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		pe.putInt16(int16(errorCode))
	***REMOVED***

	return nil
***REMOVED***

func (d *DeleteTopicsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if version >= 1 ***REMOVED***
		throttleTime, err := pd.getInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		d.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

		d.Version = version
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	d.TopicErrorCodes = make(map[string]KError, n)

	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		errorCode, err := pd.getInt16()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		d.TopicErrorCodes[topic] = KError(errorCode)
	***REMOVED***

	return nil
***REMOVED***

func (d *DeleteTopicsResponse) key() int16 ***REMOVED***
	return 20
***REMOVED***

func (d *DeleteTopicsResponse) version() int16 ***REMOVED***
	return d.Version
***REMOVED***

func (d *DeleteTopicsResponse) requiredVersion() KafkaVersion ***REMOVED***
	switch d.Version ***REMOVED***
	case 1:
		return V0_11_0_0
	default:
		return V0_10_1_0
	***REMOVED***
***REMOVED***
