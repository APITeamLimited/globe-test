package sarama

type AddOffsetsToTxnRequest struct ***REMOVED***
	TransactionalID string
	ProducerID      int64
	ProducerEpoch   int16
	GroupID         string
***REMOVED***

func (a *AddOffsetsToTxnRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(a.TransactionalID); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt64(a.ProducerID)

	pe.putInt16(a.ProducerEpoch)

	if err := pe.putString(a.GroupID); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (a *AddOffsetsToTxnRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if a.TransactionalID, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.ProducerEpoch, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if a.GroupID, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (a *AddOffsetsToTxnRequest) key() int16 ***REMOVED***
	return 25
***REMOVED***

func (a *AddOffsetsToTxnRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (a *AddOffsetsToTxnRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
