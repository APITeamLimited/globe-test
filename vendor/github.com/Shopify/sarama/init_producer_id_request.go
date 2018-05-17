package sarama

import "time"

type InitProducerIDRequest struct ***REMOVED***
	TransactionalID    *string
	TransactionTimeout time.Duration
***REMOVED***

func (i *InitProducerIDRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putNullableString(i.TransactionalID); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt32(int32(i.TransactionTimeout / time.Millisecond))

	return nil
***REMOVED***

func (i *InitProducerIDRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if i.TransactionalID, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	timeout, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i.TransactionTimeout = time.Duration(timeout) * time.Millisecond

	return nil
***REMOVED***

func (i *InitProducerIDRequest) key() int16 ***REMOVED***
	return 22
***REMOVED***

func (i *InitProducerIDRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (i *InitProducerIDRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
