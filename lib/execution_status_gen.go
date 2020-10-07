// Code generated by "enumer -type=ExecutionStatus -trimprefix ExecutionStatus -output execution_status_gen.go"; DO NOT EDIT.

//
package lib

import (
	"fmt"
)

const _ExecutionStatusName = "CreatedInitVUsInitExecutorsInitDonePausedBeforeRunStartedSetupRunningTeardownEnded"

var _ExecutionStatusIndex = [...]uint8***REMOVED***0, 7, 14, 27, 35, 50, 57, 62, 69, 77, 82***REMOVED***

func (i ExecutionStatus) String() string ***REMOVED***
	if i >= ExecutionStatus(len(_ExecutionStatusIndex)-1) ***REMOVED***
		return fmt.Sprintf("ExecutionStatus(%d)", i)
	***REMOVED***
	return _ExecutionStatusName[_ExecutionStatusIndex[i]:_ExecutionStatusIndex[i+1]]
***REMOVED***

var _ExecutionStatusValues = []ExecutionStatus***REMOVED***0, 1, 2, 3, 4, 5, 6, 7, 8, 9***REMOVED***

var _ExecutionStatusNameToValueMap = map[string]ExecutionStatus***REMOVED***
	_ExecutionStatusName[0:7]:   0,
	_ExecutionStatusName[7:14]:  1,
	_ExecutionStatusName[14:27]: 2,
	_ExecutionStatusName[27:35]: 3,
	_ExecutionStatusName[35:50]: 4,
	_ExecutionStatusName[50:57]: 5,
	_ExecutionStatusName[57:62]: 6,
	_ExecutionStatusName[62:69]: 7,
	_ExecutionStatusName[69:77]: 8,
	_ExecutionStatusName[77:82]: 9,
***REMOVED***

// ExecutionStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ExecutionStatusString(s string) (ExecutionStatus, error) ***REMOVED***
	if val, ok := _ExecutionStatusNameToValueMap[s]; ok ***REMOVED***
		return val, nil
	***REMOVED***
	return 0, fmt.Errorf("%s does not belong to ExecutionStatus values", s)
***REMOVED***

// ExecutionStatusValues returns all values of the enum
func ExecutionStatusValues() []ExecutionStatus ***REMOVED***
	return _ExecutionStatusValues
***REMOVED***

// IsAExecutionStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ExecutionStatus) IsAExecutionStatus() bool ***REMOVED***
	for _, v := range _ExecutionStatusValues ***REMOVED***
		if i == v ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
