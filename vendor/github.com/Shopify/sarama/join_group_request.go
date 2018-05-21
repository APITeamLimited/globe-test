package sarama

type GroupProtocol struct ***REMOVED***
	Name     string
	Metadata []byte
***REMOVED***

func (p *GroupProtocol) decode(pd packetDecoder) (err error) ***REMOVED***
	p.Name, err = pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.Metadata, err = pd.getBytes()
	return err
***REMOVED***

func (p *GroupProtocol) encode(pe packetEncoder) (err error) ***REMOVED***
	if err := pe.putString(p.Name); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putBytes(p.Metadata); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

type JoinGroupRequest struct ***REMOVED***
	GroupId               string
	SessionTimeout        int32
	MemberId              string
	ProtocolType          string
	GroupProtocols        map[string][]byte // deprecated; use OrderedGroupProtocols
	OrderedGroupProtocols []*GroupProtocol
***REMOVED***

func (r *JoinGroupRequest) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(r.GroupId); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt32(r.SessionTimeout)
	if err := pe.putString(r.MemberId); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(r.ProtocolType); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(r.GroupProtocols) > 0 ***REMOVED***
		if len(r.OrderedGroupProtocols) > 0 ***REMOVED***
			return PacketDecodingError***REMOVED***"cannot specify both GroupProtocols and OrderedGroupProtocols on JoinGroupRequest"***REMOVED***
		***REMOVED***

		if err := pe.putArrayLength(len(r.GroupProtocols)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for name, metadata := range r.GroupProtocols ***REMOVED***
			if err := pe.putString(name); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := pe.putBytes(metadata); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := pe.putArrayLength(len(r.OrderedGroupProtocols)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, protocol := range r.OrderedGroupProtocols ***REMOVED***
			if err := protocol.encode(pe); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *JoinGroupRequest) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if r.GroupId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.SessionTimeout, err = pd.getInt32(); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.MemberId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.ProtocolType, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil
	***REMOVED***

	r.GroupProtocols = make(map[string][]byte)
	for i := 0; i < n; i++ ***REMOVED***
		protocol := &GroupProtocol***REMOVED******REMOVED***
		if err := protocol.decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
		r.GroupProtocols[protocol.Name] = protocol.Metadata
		r.OrderedGroupProtocols = append(r.OrderedGroupProtocols, protocol)
	***REMOVED***

	return nil
***REMOVED***

func (r *JoinGroupRequest) key() int16 ***REMOVED***
	return 11
***REMOVED***

func (r *JoinGroupRequest) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *JoinGroupRequest) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***

func (r *JoinGroupRequest) AddGroupProtocol(name string, metadata []byte) ***REMOVED***
	r.OrderedGroupProtocols = append(r.OrderedGroupProtocols, &GroupProtocol***REMOVED***
		Name:     name,
		Metadata: metadata,
	***REMOVED***)
***REMOVED***

func (r *JoinGroupRequest) AddGroupProtocolMetadata(name string, metadata *ConsumerGroupMemberMetadata) error ***REMOVED***
	bin, err := encode(metadata, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.AddGroupProtocol(name, bin)
	return nil
***REMOVED***
