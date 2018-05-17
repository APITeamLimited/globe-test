package sarama

type AddPartitionsToTxnRequest struct ***REMOVED***
	TransactionalID string
	ProducerID      int64
	ProducerEpoch   int16
	TopicPartitions map[string][]int32
***REMOVED***

func (a *AddPartitionsToTxnRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(a.TransactionalID); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt64(a.ProducerID)
	pe.putInt16(a.ProducerEpoch)

	if err := pe.putArrayLength(len(a.TopicPartitions)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range a.TopicPartitions ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putInt32Array(partitions); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (a *AddPartitionsToTxnRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if a.TransactionalID, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.ProducerEpoch, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	a.TopicPartitions = make(map[string][]int32)
	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		partitions, err := pd.getInt32Array()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		a.TopicPartitions[topic] = partitions
	***REMOVED***

	return nil
***REMOVED***

func (a *AddPartitionsToTxnRequest) key() int16 ***REMOVED***
	return 24
***REMOVED***

func (a *AddPartitionsToTxnRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *AddPartitionsToTxnRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
