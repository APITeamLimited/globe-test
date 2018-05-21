package sarama

import (
	"time"
)

type CreateTopicsRequest struct ***REMOVED***
	Version int16

	TopicDetails map[string]*TopicDetail
	Timeout      time.Duration
	ValidateOnly bool
***REMOVED***

func (c *CreateTopicsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(c.TopicDetails)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, detail := range c.TopicDetails ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := detail.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	pe.putInt32(int32(c.Timeout / time.Millisecond))

	if c.Version >= 1 ***REMOVED***
		pe.putBool(c.ValidateOnly)
	***REMOVED***

	return nil
***REMOVED***

func (c *CreateTopicsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.TopicDetails = make(map[string]*TopicDetail, n)

	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.TopicDetails[topic] = new(TopicDetail)
		if err = c.TopicDetails[topic].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	timeout, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Timeout = time.Duration(timeout) * time.Millisecond

	if version >= 1 ***REMOVED***
		c.ValidateOnly, err = pd.getBool()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		c.Version = version
	***REMOVED***

	return nil
***REMOVED***

func (c *CreateTopicsRequest) key() int16 ***REMOVED***
	return 19
***REMOVED***

func (c *CreateTopicsRequest) version() int16 ***REMOVED***
	return c.Version
***REMOVED***

func (c *CreateTopicsRequest) requiredVersion() KafkaVersion ***REMOVED***
	switch c.Version ***REMOVED***
	case 2:
		return V1_0_0_0
	case 1:
		return V0_11_0_0
	default:
		return V0_10_1_0
	***REMOVED***
***REMOVED***

type TopicDetail struct ***REMOVED***
	NumPartitions     int32
	ReplicationFactor int16
	ReplicaAssignment map[int32][]int32
	ConfigEntries     map[string]*string
***REMOVED***

func (t *TopicDetail) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(t.NumPartitions)
	pe.putInt16(t.ReplicationFactor)

	if err := pe.putArrayLength(len(t.ReplicaAssignment)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for partition, assignment := range t.ReplicaAssignment ***REMOVED***
		pe.putInt32(partition)
		if err := pe.putInt32Array(assignment); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := pe.putArrayLength(len(t.ConfigEntries)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for configKey, configValue := range t.ConfigEntries ***REMOVED***
		if err := pe.putString(configKey); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putNullableString(configValue); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (t *TopicDetail) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if t.NumPartitions, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if t.ReplicationFactor, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if n > 0 ***REMOVED***
		t.ReplicaAssignment = make(map[int32][]int32, n)
		for i := 0; i < n; i++ ***REMOVED***
			replica, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if t.ReplicaAssignment[replica], err = pd.getInt32Array(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	n, err = pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if n > 0 ***REMOVED***
		t.ConfigEntries = make(map[string]*string, n)
		for i := 0; i < n; i++ ***REMOVED***
			configKey, err := pd.getString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if t.ConfigEntries[configKey], err = pd.getNullableString(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
