package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Int is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct ***REMOVED***
	sql.NullInt64
***REMOVED***

// NewInt creates a new Int
func NewInt(i int64, valid bool) Int ***REMOVED***
	return Int***REMOVED***
		NullInt64: sql.NullInt64***REMOVED***
			Int64: i,
			Valid: valid,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int64) Int ***REMOVED***
	return NewInt(i, true)
***REMOVED***

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int ***REMOVED***
	if i == nil ***REMOVED***
		return NewInt(0, false)
	***REMOVED***
	return NewInt(*i, true)
***REMOVED***

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int) ValueOrZero() int64 ***REMOVED***
	if !i.Valid ***REMOVED***
		return 0
	***REMOVED***
	return i.Int64
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error ***REMOVED***
	var err error
	var v interface***REMOVED******REMOVED***
	if err = json.Unmarshal(data, &v); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch x := v.(type) ***REMOVED***
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case string:
		str := string(x)
		if len(str) == 0 ***REMOVED***
			i.Valid = false
			return nil
		***REMOVED***
		i.Int64, err = strconv.ParseInt(str, 10, 64)
	case map[string]interface***REMOVED******REMOVED***:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	***REMOVED***
	i.Valid = err == nil
	return err
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error ***REMOVED***
	str := string(text)
	if str == "" || str == "null" ***REMOVED***
		i.Valid = false
		return nil
	***REMOVED***
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
***REMOVED***

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) ***REMOVED***
	if !i.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) ***REMOVED***
	if !i.Valid ***REMOVED***
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
***REMOVED***

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) ***REMOVED***
	i.Int64 = n
	i.Valid = true
***REMOVED***

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 ***REMOVED***
	if !i.Valid ***REMOVED***
		return nil
	***REMOVED***
	return &i.Int64
***REMOVED***

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool ***REMOVED***
	return !i.Valid
***REMOVED***
