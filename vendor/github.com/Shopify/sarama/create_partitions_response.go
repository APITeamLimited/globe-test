package sarama

import "time"

type CreatePartitionsResponse struct ***REMOVED***
	ThrottleTime         time.Duration
	TopicPartitionErrors map[string]*TopicPartitionError
***REMOVED***

func (c *CreatePartitionsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(c.ThrottleTime / time.Millisecond))
	if err := pe.putArrayLength(len(c.TopicPartitionErrors)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, partitionError := range c.TopicPartitionErrors ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := partitionError.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *CreatePartitionsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.TopicPartitionErrors = make(map[string]*TopicPartitionError, n)
	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.TopicPartitionErrors[topic] = new(TopicPartitionError)
		if err := c.TopicPartitionErrors[topic].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *CreatePartitionsResponse) key() int16 ***REMOVED***
	return 37
***REMOVED***

func (r *CreatePartitionsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *CreatePartitionsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V1_0_0_0
***REMOVED***

type TopicPartitionError struct ***REMOVED***
	Err    KError
	ErrMsg *string
***REMOVED***

func (t *TopicPartitionError) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(t.Err))

	if err := pe.putNullableString(t.ErrMsg); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (t *TopicPartitionError) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.Err = KError(kerr)

	if t.ErrMsg, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
