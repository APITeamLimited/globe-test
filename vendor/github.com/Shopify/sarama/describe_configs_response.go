package sarama

import "time"

type DescribeConfigsResponse struct ***REMOVED***
	ThrottleTime time.Duration
	Resources    []*ResourceResponse
***REMOVED***

type ResourceResponse struct ***REMOVED***
	ErrorCode int16
	ErrorMsg  string
	Type      ConfigResourceType
	Name      string
	Configs   []*ConfigEntry
***REMOVED***

type ConfigEntry struct ***REMOVED***
	Name      string
	Value     string
	ReadOnly  bool
	Default   bool
	Sensitive bool
***REMOVED***

func (r *DescribeConfigsResponse) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt32(int32(r.ThrottleTime / time.Millisecond))
	if err = pe.putArrayLength(len(r.Resources)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, c := range r.Resources ***REMOVED***
		if err = c.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *DescribeConfigsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Resources = make([]*ResourceResponse, n)
	for i := 0; i < n; i++ ***REMOVED***
		rr := &ResourceResponse***REMOVED******REMOVED***
		if err := rr.decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Resources[i] = rr
	***REMOVED***

	return nil
***REMOVED***

func (r *DescribeConfigsResponse) key() int16 ***REMOVED***
	return 32
***REMOVED***

func (r *DescribeConfigsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *DescribeConfigsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***

func (r *ResourceResponse) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt16(r.ErrorCode)

	if err = pe.putString(r.ErrorMsg); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt8(int8(r.Type))

	if err = pe.putString(r.Name); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = pe.putArrayLength(len(r.Configs)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, c := range r.Configs ***REMOVED***
		if err = c.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *ResourceResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	ec, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.ErrorCode = ec

	em, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.ErrorMsg = em

	t, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Type = ConfigResourceType(t)

	name, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Name = name

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Configs = make([]*ConfigEntry, n)
	for i := 0; i < n; i++ ***REMOVED***
		c := &ConfigEntry***REMOVED******REMOVED***
		if err := c.decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
		r.Configs[i] = c
	***REMOVED***
	return nil
***REMOVED***

func (r *ConfigEntry) encode(pe packetEncoder) (err error) ***REMOVED***
	if err = pe.putString(r.Name); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = pe.putString(r.Value); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putBool(r.ReadOnly)
	pe.putBool(r.Default)
	pe.putBool(r.Sensitive)
	return nil
***REMOVED***

func (r *ConfigEntry) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	name, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Name = name

	value, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Value = value

	read, err := pd.getBool()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.ReadOnly = read

	de, err := pd.getBool()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Default = de

	sensitive, err := pd.getBool()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Sensitive = sensitive
	return nil
***REMOVED***
