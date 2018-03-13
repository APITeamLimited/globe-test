package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Time is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Time struct ***REMOVED***
	Time  time.Time
	Valid bool
***REMOVED***

// Scan implements the Scanner interface.
func (t *Time) Scan(value interface***REMOVED******REMOVED***) error ***REMOVED***
	var err error
	switch x := value.(type) ***REMOVED***
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	***REMOVED***
	t.Valid = err == nil
	return err
***REMOVED***

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) ***REMOVED***
	if !t.Valid ***REMOVED***
		return nil, nil
	***REMOVED***
	return t.Time, nil
***REMOVED***

// NewTime creates a new Time.
func NewTime(t time.Time, valid bool) Time ***REMOVED***
	return Time***REMOVED***
		Time:  t,
		Valid: valid,
	***REMOVED***
***REMOVED***

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time ***REMOVED***
	return NewTime(t, true)
***REMOVED***

// TimeFromPtr creates a new Time that will be null if t is nil.
func TimeFromPtr(t *time.Time) Time ***REMOVED***
	if t == nil ***REMOVED***
		return NewTime(time.Time***REMOVED******REMOVED***, false)
	***REMOVED***
	return NewTime(*t, true)
***REMOVED***

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Time) ValueOrZero() time.Time ***REMOVED***
	if !t.Valid ***REMOVED***
		return time.Time***REMOVED******REMOVED***
	***REMOVED***
	return t.Time
***REMOVED***

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) ***REMOVED***
	if !t.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	return t.Time.MarshalJSON()
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *Time) UnmarshalJSON(data []byte) error ***REMOVED***
	var err error
	var v interface***REMOVED******REMOVED***
	if err = json.Unmarshal(data, &v); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch x := v.(type) ***REMOVED***
	case string:
		err = t.Time.UnmarshalJSON(data)
	case map[string]interface***REMOVED******REMOVED***:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK ***REMOVED***
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		***REMOVED***
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	***REMOVED***
	t.Valid = err == nil
	return err
***REMOVED***

func (t Time) MarshalText() ([]byte, error) ***REMOVED***
	if !t.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	return t.Time.MarshalText()
***REMOVED***

func (t *Time) UnmarshalText(text []byte) error ***REMOVED***
	str := string(text)
	if str == "" || str == "null" ***REMOVED***
		t.Valid = false
		return nil
	***REMOVED***
	if err := t.Time.UnmarshalText(text); err != nil ***REMOVED***
		return err
	***REMOVED***
	t.Valid = true
	return nil
***REMOVED***

// SetValid changes this Time's value and sets it to be non-null.
func (t *Time) SetValid(v time.Time) ***REMOVED***
	t.Time = v
	t.Valid = true
***REMOVED***

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t Time) Ptr() *time.Time ***REMOVED***
	if !t.Valid ***REMOVED***
		return nil
	***REMOVED***
	return &t.Time
***REMOVED***
