package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Float is a nullable float64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float struct ***REMOVED***
	sql.NullFloat64
***REMOVED***

// NewFloat creates a new Float
func NewFloat(f float64, valid bool) Float ***REMOVED***
	return Float***REMOVED***
		NullFloat64: sql.NullFloat64***REMOVED***
			Float64: f,
			Valid:   valid,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// FloatFrom creates a new Float that will always be valid.
func FloatFrom(f float64) Float ***REMOVED***
	return NewFloat(f, true)
***REMOVED***

// FloatFromPtr creates a new Float that be null if f is nil.
func FloatFromPtr(f *float64) Float ***REMOVED***
	if f == nil ***REMOVED***
		return NewFloat(0, false)
	***REMOVED***
	return NewFloat(*f, true)
***REMOVED***

// ValueOrZero returns the inner value if valid, otherwise zero.
func (f Float) ValueOrZero() float64 ***REMOVED***
	if !f.Valid ***REMOVED***
		return 0
	***REMOVED***
	return f.Float64
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float.
// It also supports unmarshalling a sql.NullFloat64.
func (f *Float) UnmarshalJSON(data []byte) error ***REMOVED***
	var err error
	var v interface***REMOVED******REMOVED***
	if err = json.Unmarshal(data, &v); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch x := v.(type) ***REMOVED***
	case float64:
		f.Float64 = float64(x)
	case string:
		str := string(x)
		if len(str) == 0 ***REMOVED***
			f.Valid = false
			return nil
		***REMOVED***
		f.Float64, err = strconv.ParseFloat(str, 64)
	case map[string]interface***REMOVED******REMOVED***:
		err = json.Unmarshal(data, &f.NullFloat64)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float", reflect.TypeOf(v).Name())
	***REMOVED***
	f.Valid = err == nil
	return err
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float) UnmarshalText(text []byte) error ***REMOVED***
	str := string(text)
	if str == "" || str == "null" ***REMOVED***
		f.Valid = false
		return nil
	***REMOVED***
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil
	return err
***REMOVED***

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float is null.
func (f Float) MarshalJSON() ([]byte, error) ***REMOVED***
	if !f.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	if math.IsInf(f.Float64, 0) || math.IsNaN(f.Float64) ***REMOVED***
		return nil, &json.UnsupportedValueError***REMOVED***
			Value: reflect.ValueOf(f.Float64),
			Str:   strconv.FormatFloat(f.Float64, 'g', -1, 64),
		***REMOVED***
	***REMOVED***
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float is null.
func (f Float) MarshalText() ([]byte, error) ***REMOVED***
	if !f.Valid ***REMOVED***
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
***REMOVED***

// SetValid changes this Float's value and also sets it to be non-null.
func (f *Float) SetValid(n float64) ***REMOVED***
	f.Float64 = n
	f.Valid = true
***REMOVED***

// Ptr returns a pointer to this Float's value, or a nil pointer if this Float is null.
func (f Float) Ptr() *float64 ***REMOVED***
	if !f.Valid ***REMOVED***
		return nil
	***REMOVED***
	return &f.Float64
***REMOVED***

// IsZero returns true for invalid Floats, for future omitempty support (Go 1.4?)
// A non-null Float with a 0 value will not be considered zero.
func (f Float) IsZero() bool ***REMOVED***
	return !f.Valid
***REMOVED***
