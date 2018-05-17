package sarama

type EndTxnRequest struct ***REMOVED***
	TransactionalID   string
	ProducerID        int64
	ProducerEpoch     int16
	TransactionResult bool
***REMOVED***

func (a *EndTxnRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(a.TransactionalID); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt64(a.ProducerID)

	pe.putInt16(a.ProducerEpoch)

	pe.putBool(a.TransactionResult)

	return nil
***REMOVED***

func (a *EndTxnRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if a.TransactionalID, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.ProducerEpoch, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.TransactionResult, err = pd.getBool(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (a *EndTxnRequest) key() int16 ***REMOVED***
	return 26
***REMOVED***

func (a *EndTxnRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *EndTxnRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
