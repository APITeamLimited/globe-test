package sarama

type DescribeGroupsResponse struct ***REMOVED***
	Groups []*GroupDescription
***REMOVED***

func (r *DescribeGroupsResponse) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putArrayLength(len(r.Groups)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, groupDescription := range r.Groups ***REMOVED***
		if err := groupDescription.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *DescribeGroupsResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Groups = make([]*GroupDescription, n)
	for i := 0; i < n; i++ ***REMOVED***
		r.Groups[i] = new(GroupDescription)
		if err := r.Groups[i].decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *DescribeGroupsResponse) key() int16 ***REMOVED***
	return 15
***REMOVED***

func (r *DescribeGroupsResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *DescribeGroupsResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_9_0_0
***REMOVED***

type GroupDescription struct ***REMOVED***
	Err          KError
	GroupId      string
	State        string
	ProtocolType string
	Protocol     string
	Members      map[string]*GroupMemberDescription
***REMOVED***

func (gd *GroupDescription) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(gd.Err))

	if err := pe.putString(gd.GroupId); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(gd.State); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(gd.ProtocolType); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(gd.Protocol); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(gd.Members)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for memberId, groupMemberDescription := range gd.Members ***REMOVED***
		if err := pe.putString(memberId); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := groupMemberDescription.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (gd *GroupDescription) decode(pd packetDecoder) (err error) ***REMOVED***
	kerr, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	gd.Err = KError(kerr)

	if gd.GroupId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if gd.State, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if gd.ProtocolType, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if gd.Protocol, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil
	***REMOVED***

	gd.Members = make(map[string]*GroupMemberDescription)
	for i := 0; i < n; i++ ***REMOVED***
		memberId, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		gd.Members[memberId] = new(GroupMemberDescription)
		if err := gd.Members[memberId].decode(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type GroupMemberDescription struct ***REMOVED***
	ClientId         string
	ClientHost       string
	MemberMetadata   []byte
	MemberAssignment []byte
***REMOVED***

func (gmd *GroupMemberDescription) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(gmd.ClientId); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putString(gmd.ClientHost); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putBytes(gmd.MemberMetadata); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putBytes(gmd.MemberAssignment); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (gmd *GroupMemberDescription) decode(pd packetDecoder) (err error) ***REMOVED***
	if gmd.ClientId, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if gmd.ClientHost, err = pd.getString(); err != nil ***REMOVED***
		return
	***REMOVED***
	if gmd.MemberMetadata, err = pd.getBytes(); err != nil ***REMOVED***
		return
	***REMOVED***
	if gmd.MemberAssignment, err = pd.getBytes(); err != nil ***REMOVED***
		return
	***REMOVED***

	return nil
***REMOVED***

func (gmd *GroupMemberDescription) GetMemberAssignment() (*ConsumerGroupMemberAssignment, error) ***REMOVED***
	assignment := new(ConsumerGroupMemberAssignment)
	err := decode(gmd.MemberAssignment, assignment)
	return assignment, err
***REMOVED***

func (gmd *GroupMemberDescription) GetMemberMetadata() (*ConsumerGroupMemberMetadata, error) ***REMOVED***
	metadata := new(ConsumerGroupMemberMetadata)
	err := decode(gmd.MemberMetadata, metadata)
	return metadata, err
***REMOVED***
