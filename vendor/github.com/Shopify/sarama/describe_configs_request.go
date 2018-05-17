package sarama

type ConfigResource struct ***REMOVED***
	Type        ConfigResourceType
	Name        string
	ConfigNames []string
***REMOVED***

type DescribeConfigsRequest struct ***REMOVED***
	Resources []*ConfigResource
***REMOVED***

func (r *DescribeConfigsRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(r.Resources)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, c := range r.Resources ***REMOVED***
		pe.putInt8(int8(c.Type))
		if err := pe.putString(c.Name); err != nil ***REMOVED***
			return err
		***REMOVED***

		if len(c.ConfigNames) == 0 ***REMOVED***
			pe.putInt32(-1)
			continue
		***REMOVED***
		if err := pe.putStringArray(c.ConfigNames); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *DescribeConfigsRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Resources = make([]*ConfigResource, n)

	for i := 0; i < n; i++ ***REMOVED***
		r.Resources[i] = &ConfigResource***REMOVED******REMOVED***
		t, err := pd.getInt8()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Resources[i].Type = ConfigResourceType(t)
		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Resources[i].Name = name

		confLength, err := pd.getArrayLength()

		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if confLength == -1 ***REMOVED***
			continue
		***REMOVED***

		cfnames := make([]string, confLength)
		for i := 0; i < confLength; i++ ***REMOVED***
			s, err := pd.getString()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			cfnames[i] = s
		***REMOVED***
		r.Resources[i].ConfigNames = cfnames
	***REMOVED***

	return nil
***REMOVED***

func (r *DescribeConfigsRequest) key() int16 ***REMOVED***
	return 32
***REMOVED***

func (r *DescribeConfigsRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *DescribeConfigsRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***
