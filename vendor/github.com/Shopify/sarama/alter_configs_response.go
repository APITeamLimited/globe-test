package sarama

import "time"

type AlterConfigsResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Resources    []*AlterConfigsResourceResponse
***REMOVED***

type AlterConfigsResourceResponse struct ***REMOVED***
	ErrorCode int16
	ErrorMsg  string
	Type      ConfigResourceType
	Name      string
***REMOVED***

func (ct *AlterConfigsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(ct.ThrottleTime / time.Millisecond))

	if err := pe.putArrayLength(len(ct.Resources)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for i := range ct.Resources ***REMOVED***
		pe.putInt16(ct.Resources[i].ErrorCode)
		err := pe.putString(ct.Resources[i].ErrorMsg)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		pe.putInt8(int8(ct.Resources[i].Type))
		err = pe.putString(ct.Resources[i].Name)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (acr *AlterConfigsResponse) decode(pd packetDecoder, version int16) error ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	acr.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	responseCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	acr.Resources = make([]*AlterConfigsResourceResponse, responseCount)

	for i := range acr.Resources ***REMOVED***
		acr.Resources[i] = new(AlterConfigsResourceResponse)

		errCode, err := pd.getInt16()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		acr.Resources[i].ErrorCode = errCode

		e, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		acr.Resources[i].ErrorMsg = e

		t, err := pd.getInt8()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		acr.Resources[i].Type = ConfigResourceType(t)

		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		acr.Resources[i].Name = name
	***REMOVED***

	return nil
***REMOVED***

func (r *AlterConfigsResponse) key() int16 ***REMOVED***
	return 32
***REMOVED***

func (r *AlterConfigsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *AlterConfigsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
