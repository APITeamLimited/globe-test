package sarama

type SaslHandshakeResponse struct ***REMOVED***
	Err               KError
	EnabledMechanisms []string
***REMOVED***

func (r *SaslHandshakeResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	return pe.putStringArray(r.EnabledMechanisms)
***REMOVED***

func (r *SaslHandshakeResponse) decode(pd packetDecoder, version int16) error ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Err = KError(kerr)

	if r.EnabledMechanisms, err = pd.getStringArray(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *SaslHandshakeResponse) key() int16 ***REMOVED***
	return 17
***REMOVED***

func (r *SaslHandshakeResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *SaslHandshakeResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_10_0_0
***REMOVED***
