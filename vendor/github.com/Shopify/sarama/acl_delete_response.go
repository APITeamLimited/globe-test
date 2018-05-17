package sarama

import "time"

type DeleteAclsResponse struct ***REMOVED***
	ThrottleTime    time.Duration
	FilterResponses []*FilterResponse
***REMOVED***

func (a *DeleteAclsResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt32(int32(a.ThrottleTime / time.Millisecond))

	if err := pe.putArrayLength(len(a.FilterResponses)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, filterResponse := range a.FilterResponses ***REMOVED***
		if err := filterResponse.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (a *DeleteAclsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	throttleTime, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.ThrottleTime = time.Duration(throttleTime) * time.Millisecond

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.FilterResponses = make([]*FilterResponse, n)

	for i := 0; i < n; i++ ***REMOVED***
		a.FilterResponses[i] = new(FilterResponse)
		if err := a.FilterResponses[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *DeleteAclsResponse) key() int16 ***REMOVED***
	return 31
***REMOVED***

func (d *DeleteAclsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (d *DeleteAclsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_11_0_0
***REMOVED***

type FilterResponse struct ***REMOVED***
	Err          KError
	ErrMsg       *string
	MatchingAcls []*MatchingAcl
***REMOVED***

func (f *FilterResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(f.Err))
	if err := pe.putNullableString(f.ErrMsg); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(f.MatchingAcls)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, matchingAcl := range f.MatchingAcls ***REMOVED***
		if err := matchingAcl.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (f *FilterResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.Err = KError(kerr)

	if f.ErrMsg, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	f.MatchingAcls = make([]*MatchingAcl, n)
	for i := 0; i < n; i++ ***REMOVED***
		f.MatchingAcls[i] = new(MatchingAcl)
		if err := f.MatchingAcls[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type MatchingAcl struct ***REMOVED***
	Err    KError
	ErrMsg *string
	Resource
	Acl
***REMOVED***

func (m *MatchingAcl) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(m.Err))
	if err := pe.putNullableString(m.ErrMsg); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := m.Resource.encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := m.Acl.encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (m *MatchingAcl) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m.Err = KError(kerr)

	if m.ErrMsg, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := m.Resource.decode(pd, version); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := m.Acl.decode(pd, version); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
