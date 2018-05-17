package sarama

import "time"

type CreatePartitionsRequest struct ***REMOVED***
	TopicPartitions map[string]*TopicPartition
	Timeout         time.Duration
	ValidateOnly    bool
***REMOVED***

func (c *CreatePartitionsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(c.TopicPartitions)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, partition := range c.TopicPartitions ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := partition.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	pe.putInt32(int32(c.Timeout / time.Millisecond))

	pe.putBool(c.ValidateOnly)

	return nil
***REMOVED***

func (c *CreatePartitionsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.TopicPartitions = make(map[string]*TopicPartition, n)
	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.TopicPartitions[topic] = new(TopicPartition)
		if err := c.TopicPartitions[topic].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	timeout, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Timeout = time.Duration(timeout) * time.Millisecond

	if c.ValidateOnly, err = pd.getBool(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *CreatePartitionsRequest) key() int16 ***REMOVED***
	return 37
***REMOVED***

func (r *CreatePartitionsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *CreatePartitionsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V1_0_0_0
***REMOVED***

type TopicPartition struct ***REMOVED***
	Count      int32
	Assignment [][]int32
***REMOVED***

func (t *TopicPartition) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(t.Count)

	if len(t.Assignment) == 0 ***REMOVED***
		pe.putInt32(-1)
		return nil
	***REMOVED***

	if err := pe.putArrayLength(len(t.Assignment)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, assign := range t.Assignment ***REMOVED***
		if err := pe.putInt32Array(assign); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (t *TopicPartition) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if t.Count, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n <= 0 ***REMOVED***
		return nil
	***REMOVED***
	t.Assignment = make([][]int32, n)

	for i := 0; i < int(n); i++ ***REMOVED***
		if t.Assignment[i], err = pd.getInt32Array(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
