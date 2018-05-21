package sarama

type SaslHandshakeRequest struct ***REMOVED***
	Mechanism string
***REMOVED***

func (r *SaslHandshakeRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(r.Mechanism); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *SaslHandshakeRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if r.Mechanism, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *SaslHandshakeRequest) key() int16 ***REMOVED***
	return 17
***REMOVED***

func (r *SaslHandshakeRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *SaslHandshakeRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_10_0_0
***REMOVED***
