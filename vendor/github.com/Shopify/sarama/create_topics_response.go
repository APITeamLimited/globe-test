package sarama

import "time"

type CreateTopicsResponse struct ***REMOVED***
	Version      int16
	ThrottleTime time.Duration
	TopicErrors  map[string]*TopicError
***REMOVED***

func (c *CreateTopicsResponse) encode(pe packetEncoder) error ***REMOVED***
	if c.Version >= 2 ***REMOVED***
		pe.putInt32(int32(c.ThrottleTime / time.Millisecond))
	***REMOVED***

	if err := pe.putArrayLength(len(c.TopicErrors)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, topicError := range c.TopicErrors ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := topicError.encode(pe, c.Version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *CreateTopicsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	c.Version = version

	if version >= 2 ***REMOVED***
		throttleTime, err := pd.getInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.ThrottleTime = time.Duration(throttleTime) * time.Millisecond
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.TopicErrors = make(map[string]*TopicError, n)
	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.TopicErrors[topic] = new(TopicError)
		if err := c.TopicErrors[topic].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *CreateTopicsResponse) key() int16 ***REMOVED***
	return 19
***REMOVED***

func (c *CreateTopicsResponse) version() int16 ***REMOVED***
	return c.Version
***REMOVED***

func (c *CreateTopicsResponse) requiredVersion() KafkaVersion ***REMOVED***
	switch c.Version ***REMOVED***
	case 2:
		return V1_0_0_0
	case 1:
		return V0_11_0_0
	default:
		return V0_10_1_0
	***REMOVED***
***REMOVED***

type TopicError struct ***REMOVED***
	Err    KError
	ErrMsg *string
***REMOVED***

func (t *TopicError) encode(pe packetEncoder, version int16) error ***REMOVED***
	pe.putInt16(int16(t.Err))

	if version >= 1 ***REMOVED***
		if err := pe.putNullableString(t.ErrMsg); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (t *TopicError) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kErr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.Err = KError(kErr)

	if version >= 1 ***REMOVED***
		if t.ErrMsg, err = pd.getNullableString(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
