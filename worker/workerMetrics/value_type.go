package workerMetrics

import "errors"

// Possible values for ValueType.
const (
	Default = ValueType(iota) // Values are presented as-is
	Time                      // Values are timestamps (nanoseconds)
	Data                      // Values are data amounts (bytes)
)

// ErrInvalidValueType indicates the serialized value type is invalid.
var ErrInvalidValueType = errors.New("invalid value type")

// ValueType holds the type of values a metric contains.
type ValueType int

// MarshalJSON serializes a ValueType to a JSON string.
func (t ValueType) MarshalJSON() ([]byte, error) ***REMOVED***
	txt, err := t.MarshalText()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return []byte(`"` + string(txt) + `"`), nil
***REMOVED***

// MarshalText serializes a ValueType as a human readable string.
func (t ValueType) MarshalText() ([]byte, error) ***REMOVED***
	switch t ***REMOVED***
	case Default:
		return []byte(defaultString), nil
	case Time:
		return []byte(timeString), nil
	case Data:
		return []byte(dataString), nil
	default:
		return nil, ErrInvalidValueType
	***REMOVED***
***REMOVED***

// UnmarshalText deserializes a ValueType from a string representation.
func (t *ValueType) UnmarshalText(data []byte) error ***REMOVED***
	switch string(data) ***REMOVED***
	case defaultString:
		*t = Default
	case timeString:
		*t = Time
	case dataString:
		*t = Data
	default:
		return ErrInvalidValueType
	***REMOVED***

	return nil
***REMOVED***

func (t ValueType) String() string ***REMOVED***
	switch t ***REMOVED***
	case Default:
		return defaultString
	case Time:
		return timeString
	case Data:
		return dataString
	default:
		return "[INVALID]"
	***REMOVED***
***REMOVED***
