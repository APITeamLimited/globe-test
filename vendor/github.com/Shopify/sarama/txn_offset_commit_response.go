package sarama

import (
	"time"
)

type TxnOffsetCommitResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Topics       map[string][]*PartitionError
***REMOVED***

func (t *TxnOffsetCommitResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(t.ThrottleTime / time.Millisecond))
	if err := pe.putArrayLength(len(t.Topics)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, e := range t.Topics ***REMOVED***
		if err := pe.putString(topic); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putArrayLength(len(e)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, partitionError := range e ***REMOVED***
			if err := partitionError.encode(pe); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (t *TxnOffsetCommitResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	t.Topics = make(map[string][]*PartitionError)

	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		m, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		t.Topics[topic] = make([]*PartitionError, m)

		for j := 0; j < m; j++ ***REMOVED***
			t.Topics[topic][j] = new(PartitionError)
			if err := t.Topics[topic][j].decode(pd, version); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (a *TxnOffsetCommitResponse) key() int16 ***REMOVED***
	return 28
***REMOVED***

func (a *TxnOffsetCommitResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *TxnOffsetCommitResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
