package csv

// TimeFormat custom enum type
type TimeFormat string

// valid defined values for TimeFormat
const (
	Unix    TimeFormat = "unix"
	RFC3399 TimeFormat = "rfc3399"
)

// IsValid validates TimeFormat
func (timeFormat TimeFormat) IsValid() bool ***REMOVED***
	switch timeFormat ***REMOVED***
	case Unix, RFC3399:
		return true
	***REMOVED***
	return false
***REMOVED***
