package sarama

type ApiVersionsResponseBlock struct ***REMOVED***
	ApiKey     int16
	MinVersion int16
	MaxVersion int16
***REMOVED***

func (b *ApiVersionsResponseBlock) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(b.ApiKey)
	pe.putInt16(b.MinVersion)
	pe.putInt16(b.MaxVersion)
	return nil
***REMOVED***

func (b *ApiVersionsResponseBlock) decode(pd packetDecoder) error ***REMOVED***
	var err error

	if b.ApiKey, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.MinVersion, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.MaxVersion, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

type ApiVersionsResponse struct ***REMOVED***
	Err         KError
	ApiVersions []*ApiVersionsResponseBlock
***REMOVED***

func (r *ApiVersionsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	if err := pe.putArrayLength(len(r.ApiVersions)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, apiVersion := range r.ApiVersions ***REMOVED***
		if err := apiVersion.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *ApiVersionsResponse) decode(pd packetDecoder, version int16) error ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Err = KError(kerr)

	numBlocks, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.ApiVersions = make([]*ApiVersionsResponseBlock, numBlocks)
	for i := 0; i < numBlocks; i++ ***REMOVED***
		block := new(ApiVersionsResponseBlock)
		if err := block.decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
		r.ApiVersions[i] = block
	***REMOVED***

	return nil
***REMOVED***

func (r *ApiVersionsResponse) key() int16 ***REMOVED***
	return 18
***REMOVED***

func (r *ApiVersionsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ApiVersionsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_10_0_0
***REMOVED***
