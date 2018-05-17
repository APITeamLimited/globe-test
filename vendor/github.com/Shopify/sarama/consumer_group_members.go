package sarama

type ConsumerGroupMemberMetadata struct ***REMOVED***
	Version  int16
	Topics   []string
	UserData []byte
***REMOVED***

func (m *ConsumerGroupMemberMetadata) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(m.Version)

	if err := pe.putStringArray(m.Topics); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putBytes(m.UserData); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (m *ConsumerGroupMemberMetadata) decode(pd packetDecoder) (err error) ***REMOVED***
	if m.Version, err = pd.getInt16(); err != nil ***REMOVED***
		return
	***REMOVED***

	if m.Topics, err = pd.getStringArray(); err != nil ***REMOVED***
		return
	***REMOVED***

	if m.UserData, err = pd.getBytes(); err != nil ***REMOVED***
		return
	***REMOVED***

	return nil
***REMOVED***

type ConsumerGroupMemberAssignment struct ***REMOVED***
	Version  int16
	Topics   map[string][]int32
	UserData []byte
***REMOVED***

func (m *ConsumerGroupMemberAssignment) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(m.Version)

	if err := pe.putArrayLength(len(m.Topics)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, partitions := range m.Topics ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putInt32Array(partitions); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := pe.putBytes(m.UserData); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (m *ConsumerGroupMemberAssignment) decode(pd packetDecoder) (err error) ***REMOVED***
	if m.Version, err = pd.getInt16(); err != nil ***REMOVED***
		return
	***REMOVED***

	var topicLen int
	if topicLen, err = pd.getArrayLength(); err != nil ***REMOVED***
		return
	***REMOVED***

	m.Topics = make(map[string][]int32, topicLen)
	for i := 0; i < topicLen; i++ ***REMOVED***
		var topic string
		if topic, err = pd.getString(); err != nil ***REMOVED***
			return
		***REMOVED***
		if m.Topics[topic], err = pd.getInt32Array(); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if m.UserData, err = pd.getBytes(); err != nil ***REMOVED***
		return
	***REMOVED***

	return nil
***REMOVED***
