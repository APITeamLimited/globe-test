package sarama

type TxnOffsetCommitRequest struct ***REMOVED***
	TransactionalID string
	GroupID         string
	ProducerID      int64
	ProducerEpoch   int16
	Topics          map[string][]*PartitionOffsetMetadata
***REMOVED***

func (t *TxnOffsetCommitRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(t.TransactionalID); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(t.GroupID); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt64(t.ProducerID)
	pe.putInt16(t.ProducerEpoch)

	if err := pe.putArrayLength(len(t.Topics)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for topic, partitions := range t.Topics ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putArrayLength(len(partitions)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, partition := range partitions ***REMOVED***
			if err := partition.encode(pe); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (t *TxnOffsetCommitRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if t.TransactionalID, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if t.GroupID, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if t.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if t.ProducerEpoch, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	t.Topics = make(map[string][]*PartitionOffsetMetadata)
	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		m, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		t.Topics[topic] = make([]*PartitionOffsetMetadata, m)

		for j := 0; j < m; j++ ***REMOVED***
			partitionOffsetMetadata := new(PartitionOffsetMetadata)
			if err := partitionOffsetMetadata.decode(pd, version); err != nil ***REMOVED***
				return err
			***REMOVED***
			t.Topics[topic][j] = partitionOffsetMetadata
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (a *TxnOffsetCommitRequest) key() int16 ***REMOVED***
	return 28
***REMOVED***

func (a *TxnOffsetCommitRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *TxnOffsetCommitRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***

type PartitionOffsetMetadata struct ***REMOVED***
	Partition int32
	Offset    int64
	Metadata  *string
***REMOVED***

func (p *PartitionOffsetMetadata) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(p.Partition)
	pe.putInt64(p.Offset)
	if err := pe.putNullableString(p.Metadata); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (p *PartitionOffsetMetadata) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if p.Partition, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if p.Offset, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if p.Metadata, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
