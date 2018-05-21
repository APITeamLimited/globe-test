package sarama

import (
	"time"
)

type AddPartitionsToTxnResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Errors       map[string][]*PartitionError
***REMOVED***

func (a *AddPartitionsToTxnResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(a.ThrottleTime / time.Millisecond))
	if err := pe.putArrayLength(len(a.Errors)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, e := range a.Errors ***REMOVED***
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

func (a *AddPartitionsToTxnResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	a.Errors = make(map[string][]*PartitionError)

	for i := 0; i < n; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		m, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		a.Errors[topic] = make([]*PartitionError, m)

		for j := 0; j < m; j++ ***REMOVED***
			a.Errors[topic][j] = new(PartitionError)
			if err := a.Errors[topic][j].decode(pd, version); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (a *AddPartitionsToTxnResponse) key() int16 ***REMOVED***
	return 24
***REMOVED***

func (a *AddPartitionsToTxnResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *AddPartitionsToTxnResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***

type PartitionError struct ***REMOVED***
	Partition int32
	Err       KError
***REMOVED***

func (p *PartitionError) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(p.Partition)
	pe.putInt16(int16(p.Err))
	return nil
***REMOVED***

func (p *PartitionError) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if p.Partition, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***

	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.Err = KError(kerr)

	return nil
***REMOVED***
