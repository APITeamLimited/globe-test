// Code generated by "enumer -type=CompatibilityMode -transform=snake -trimprefix CompatibilityMode -output compatibility_mode_gen.go"; DO NOT EDIT.

//
package compiler

import (
	"fmt"
)

const _CompatibilityModeName = "extendedbase"

var _CompatibilityModeIndex = [...]uint8***REMOVED***0, 8, 12***REMOVED***

func (i CompatibilityMode) String() string ***REMOVED***
	i -= 1
	if i >= CompatibilityMode(len(_CompatibilityModeIndex)-1) ***REMOVED***
		return fmt.Sprintf("CompatibilityMode(%d)", i+1)
	***REMOVED***
	return _CompatibilityModeName[_CompatibilityModeIndex[i]:_CompatibilityModeIndex[i+1]]
***REMOVED***

var _CompatibilityModeValues = []CompatibilityMode***REMOVED***1, 2***REMOVED***

var _CompatibilityModeNameToValueMap = map[string]CompatibilityMode***REMOVED***
	_CompatibilityModeName[0:8]:  1,
	_CompatibilityModeName[8:12]: 2,
***REMOVED***

// CompatibilityModeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CompatibilityModeString(s string) (CompatibilityMode, error) ***REMOVED***
	if val, ok := _CompatibilityModeNameToValueMap[s]; ok ***REMOVED***
		return val, nil
	***REMOVED***
	return 0, fmt.Errorf("%s does not belong to CompatibilityMode values", s)
***REMOVED***

// CompatibilityModeValues returns all values of the enum
func CompatibilityModeValues() []CompatibilityMode ***REMOVED***
	return _CompatibilityModeValues
***REMOVED***

// IsACompatibilityMode returns "true" if the value is listed in the enum definition. "false" otherwise
func (i CompatibilityMode) IsACompatibilityMode() bool ***REMOVED***
	for _, v := range _CompatibilityModeValues ***REMOVED***
		if i == v ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
