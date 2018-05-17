package sarama

type Resource struct ***REMOVED***
	ResourceType AclResourceType
	ResourceName string
***REMOVED***

func (r *Resource) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt8(int8(r.ResourceType))

	if err := pe.putString(r.ResourceName); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *Resource) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	resourceType, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.ResourceType = AclResourceType(resourceType)

	if r.ResourceName, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

type Acl struct ***REMOVED***
	Principal      string
	Host           string
	Operation      AclOperation
	PermissionType AclPermissionType
***REMOVED***

func (a *Acl) encode(pe packetEncoder) error ***REMOVED***
	if err := pe.putString(a.Principal); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putString(a.Host); err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt8(int8(a.Operation))
	pe.putInt8(int8(a.PermissionType))

	return nil
***REMOVED***

func (a *Acl) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	if a.Principal, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if a.Host, err = pd.getString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	operation, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.Operation = AclOperation(operation)

	permissionType, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.PermissionType = AclPermissionType(permissionType)

	return nil
***REMOVED***

type ResourceAcls struct ***REMOVED***
	Resource
	Acls []*Acl
***REMOVED***

func (r *ResourceAcls) encode(pe packetEncoder) error ***REMOVED***
	if err := r.Resource.encode(pe); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := pe.putArrayLength(len(r.Acls)); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, acl := range r.Acls ***REMOVED***
		if err := acl.encode(pe); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *ResourceAcls) decode(pd packetDecoder, version int16) error ***REMOVED***
	if err := r.Resource.decode(pd, version); err != nil ***REMOVED***
		return err
	***REMOVED***

	n, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Acls = make([]*Acl, n)
	for i := 0; i < n; i++ ***REMOVED***
		r.Acls[i] = new(Acl)
		if err := r.Acls[i].decode(pd, version); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
