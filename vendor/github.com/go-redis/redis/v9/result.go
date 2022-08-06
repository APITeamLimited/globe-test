package redis

import "time"

// NewCmdResult returns a Cmd initialised with val and err for testing.
func NewCmdResult(val interface***REMOVED******REMOVED***, err error) *Cmd ***REMOVED***
	var cmd Cmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewSliceResult returns a SliceCmd initialised with val and err for testing.
func NewSliceResult(val []interface***REMOVED******REMOVED***, err error) *SliceCmd ***REMOVED***
	var cmd SliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewStatusResult returns a StatusCmd initialised with val and err for testing.
func NewStatusResult(val string, err error) *StatusCmd ***REMOVED***
	var cmd StatusCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewIntResult returns an IntCmd initialised with val and err for testing.
func NewIntResult(val int64, err error) *IntCmd ***REMOVED***
	var cmd IntCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewDurationResult returns a DurationCmd initialised with val and err for testing.
func NewDurationResult(val time.Duration, err error) *DurationCmd ***REMOVED***
	var cmd DurationCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewBoolResult returns a BoolCmd initialised with val and err for testing.
func NewBoolResult(val bool, err error) *BoolCmd ***REMOVED***
	var cmd BoolCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewStringResult returns a StringCmd initialised with val and err for testing.
func NewStringResult(val string, err error) *StringCmd ***REMOVED***
	var cmd StringCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewFloatResult returns a FloatCmd initialised with val and err for testing.
func NewFloatResult(val float64, err error) *FloatCmd ***REMOVED***
	var cmd FloatCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewStringSliceResult returns a StringSliceCmd initialised with val and err for testing.
func NewStringSliceResult(val []string, err error) *StringSliceCmd ***REMOVED***
	var cmd StringSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewBoolSliceResult returns a BoolSliceCmd initialised with val and err for testing.
func NewBoolSliceResult(val []bool, err error) *BoolSliceCmd ***REMOVED***
	var cmd BoolSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewStringStringMapResult returns a StringStringMapCmd initialised with val and err for testing.
func NewStringStringMapResult(val map[string]string, err error) *MapStringStringCmd ***REMOVED***
	var cmd MapStringStringCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewStringIntMapCmdResult returns a StringIntMapCmd initialised with val and err for testing.
func NewStringIntMapCmdResult(val map[string]int64, err error) *StringIntMapCmd ***REMOVED***
	var cmd StringIntMapCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewTimeCmdResult returns a TimeCmd initialised with val and err for testing.
func NewTimeCmdResult(val time.Time, err error) *TimeCmd ***REMOVED***
	var cmd TimeCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewZSliceCmdResult returns a ZSliceCmd initialised with val and err for testing.
func NewZSliceCmdResult(val []Z, err error) *ZSliceCmd ***REMOVED***
	var cmd ZSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewZWithKeyCmdResult returns a NewZWithKeyCmd initialised with val and err for testing.
func NewZWithKeyCmdResult(val *ZWithKey, err error) *ZWithKeyCmd ***REMOVED***
	var cmd ZWithKeyCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewScanCmdResult returns a ScanCmd initialised with val and err for testing.
func NewScanCmdResult(keys []string, cursor uint64, err error) *ScanCmd ***REMOVED***
	var cmd ScanCmd
	cmd.page = keys
	cmd.cursor = cursor
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewClusterSlotsCmdResult returns a ClusterSlotsCmd initialised with val and err for testing.
func NewClusterSlotsCmdResult(val []ClusterSlot, err error) *ClusterSlotsCmd ***REMOVED***
	var cmd ClusterSlotsCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewGeoLocationCmdResult returns a GeoLocationCmd initialised with val and err for testing.
func NewGeoLocationCmdResult(val []GeoLocation, err error) *GeoLocationCmd ***REMOVED***
	var cmd GeoLocationCmd
	cmd.locations = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewGeoPosCmdResult returns a GeoPosCmd initialised with val and err for testing.
func NewGeoPosCmdResult(val []*GeoPos, err error) *GeoPosCmd ***REMOVED***
	var cmd GeoPosCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewCommandsInfoCmdResult returns a CommandsInfoCmd initialised with val and err for testing.
func NewCommandsInfoCmdResult(val map[string]*CommandInfo, err error) *CommandsInfoCmd ***REMOVED***
	var cmd CommandsInfoCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewXMessageSliceCmdResult returns a XMessageSliceCmd initialised with val and err for testing.
func NewXMessageSliceCmdResult(val []XMessage, err error) *XMessageSliceCmd ***REMOVED***
	var cmd XMessageSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewXStreamSliceCmdResult returns a XStreamSliceCmd initialised with val and err for testing.
func NewXStreamSliceCmdResult(val []XStream, err error) *XStreamSliceCmd ***REMOVED***
	var cmd XStreamSliceCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***

// NewXPendingResult returns a XPendingCmd initialised with val and err for testing.
func NewXPendingResult(val *XPending, err error) *XPendingCmd ***REMOVED***
	var cmd XPendingCmd
	cmd.val = val
	cmd.SetErr(err)
	return &cmd
***REMOVED***
