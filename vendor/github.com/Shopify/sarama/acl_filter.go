package sarama

type AclFilter struct ***REMOVED***
	ResourceType   AclResourceType
	ResourceName   *string
	Principal      *string
	Host           *string
	Operation      AclOperation
	PermissionType AclPermissionType
***REMOVED***

func (a *AclFilter) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt8(int8(a.ResourceType))
	if err := pe.putNullableString(a.ResourceName); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putNullableString(a.Principal); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := pe.putNullableString(a.Host); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt8(int8(a.Operation))
	pe.putInt8(int8(a.PermissionType))

	return nil
***REMOVED***

func (a *AclFilter) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	resourceType, err := pd.getInt8()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	a.ResourceType = AclResourceType(resourceType)

	if a.ResourceName, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if a.Principal, err = pd.getNullableString(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if a.Host, err = pd.getNullableString(); err != nil ***REMOVED***
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
