// Code generated by "enumer -type=SystemTagSet -transform=snake -trimprefix=Tag -output system_tag_set_gen.go"; DO NOT EDIT.

//
package stats

import (
	"fmt"
)

const _SystemTagSetName = "protosubprotostatusmethodurlnamegroupcheckerrorerror_codetls_versionscenarioitervuocsp_statusip"

var _SystemTagSetMap = map[SystemTagSet]string***REMOVED***
	1:     _SystemTagSetName[0:5],
	2:     _SystemTagSetName[5:13],
	4:     _SystemTagSetName[13:19],
	8:     _SystemTagSetName[19:25],
	16:    _SystemTagSetName[25:28],
	32:    _SystemTagSetName[28:32],
	64:    _SystemTagSetName[32:37],
	128:   _SystemTagSetName[37:42],
	256:   _SystemTagSetName[42:47],
	512:   _SystemTagSetName[47:57],
	1024:  _SystemTagSetName[57:68],
	2048:  _SystemTagSetName[68:76],
	4096:  _SystemTagSetName[76:80],
	8192:  _SystemTagSetName[80:82],
	16384: _SystemTagSetName[82:93],
	32768: _SystemTagSetName[93:95],
***REMOVED***

func (i SystemTagSet) String() string ***REMOVED***
	if str, ok := _SystemTagSetMap[i]; ok ***REMOVED***
		return str
	***REMOVED***
	return fmt.Sprintf("SystemTagSet(%d)", i)
***REMOVED***

var _SystemTagSetValues = []SystemTagSet***REMOVED***1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768***REMOVED***

var _SystemTagSetNameToValueMap = map[string]SystemTagSet***REMOVED***
	_SystemTagSetName[0:5]:   1,
	_SystemTagSetName[5:13]:  2,
	_SystemTagSetName[13:19]: 4,
	_SystemTagSetName[19:25]: 8,
	_SystemTagSetName[25:28]: 16,
	_SystemTagSetName[28:32]: 32,
	_SystemTagSetName[32:37]: 64,
	_SystemTagSetName[37:42]: 128,
	_SystemTagSetName[42:47]: 256,
	_SystemTagSetName[47:57]: 512,
	_SystemTagSetName[57:68]: 1024,
	_SystemTagSetName[68:76]: 2048,
	_SystemTagSetName[76:80]: 4096,
	_SystemTagSetName[80:82]: 8192,
	_SystemTagSetName[82:93]: 16384,
	_SystemTagSetName[93:95]: 32768,
***REMOVED***

// SystemTagSetString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SystemTagSetString(s string) (SystemTagSet, error) ***REMOVED***
	if val, ok := _SystemTagSetNameToValueMap[s]; ok ***REMOVED***
		return val, nil
	***REMOVED***
	return 0, fmt.Errorf("%s does not belong to SystemTagSet values", s)
***REMOVED***

// SystemTagSetValues returns all values of the enum
func SystemTagSetValues() []SystemTagSet ***REMOVED***
	return _SystemTagSetValues
***REMOVED***

// IsASystemTagSet returns "true" if the value is listed in the enum definition. "false" otherwise
func (i SystemTagSet) IsASystemTagSet() bool ***REMOVED***
	_, ok := _SystemTagSetMap[i]
	return ok
***REMOVED***
