package sarama

type AlterConfigsRequest struct ***REMOVED***
	Resources    []*AlterConfigsResource
	ValidateOnly bool
***REMOVED***

type AlterConfigsResource struct ***REMOVED***
	Type          ConfigResourceType
	Name          string
	ConfigEntries map[string]*string
***REMOVED***

func (acr *AlterConfigsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(acr.Resources)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, r := range acr.Resources ***REMOVED***
		if err := r.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	pe.putBool(acr.ValidateOnly)
	return nil
***REMOVED***

func (acr *AlterConfigsRequest) decode(pd packetDecoder, version int16) error ***REMOVED***
	resourceCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	acr.Resources = make([]*AlterConfigsResource, resourceCount)
	for i := range acr.Resources ***REMOVED***
		r := &AlterConfigsResource***REMOVED******REMOVED***
		err = r.decode(pd, version)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		acr.Resources[i] = r
	***REMOVED***

	validateOnly, err := pd.getBool()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	acr.ValidateOnly = validateOnly

	return nil
***REMOVED***

func (ac *AlterConfigsResource) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt8(int8(ac.Type))

	if err := pe.putString(ac.Name); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(ac.ConfigEntries)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for configKey, configValue := range ac.ConfigEntries ***REMOVED***
		if err := pe.putString(configKey); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := pe.putNullableString(configValue); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ac *AlterConfigsResource) decode(pd packetDecoder, version int16) error ***REMOVED***
	t, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	ac.Type = ConfigResourceType(t)

	name, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	ac.Name = name

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if n > 0 ***REMOVED***
		ac.ConfigEntries = make(map[string]*string, n)
		for i := 0; i < n; i++ ***REMOVED***
			configKey, err := pd.getString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if ac.ConfigEntries[configKey], err = pd.getNullableString(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (acr *AlterConfigsRequest) key() int16 ***REMOVED***
	return 33
***REMOVED***

func (acr *AlterConfigsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (acr *AlterConfigsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
