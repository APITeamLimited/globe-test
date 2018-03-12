// Package null contains SQL types that consider zero input and null input as separate values,
// with convenient support for JSON and text marshaling.
// Types in this package will always encode to their null value if null.
// Use the zero subpackage if you want zero values and null to be treated the same.
package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

// String is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
type String struct ***REMOVED***
	sql.NullString
***REMOVED***

// StringFrom creates a new String that will never be blank.
func StringFrom(s string) String ***REMOVED***
	return NewString(s, true)
***REMOVED***

// StringFromPtr creates a new String that be null if s is nil.
func StringFromPtr(s *string) String ***REMOVED***
	if s == nil ***REMOVED***
		return NewString("", false)
	***REMOVED***
	return NewString(*s, true)
***REMOVED***

// ValueOrZero returns the inner value if valid, otherwise zero.
func (s String) ValueOrZero() string ***REMOVED***
	if !s.Valid ***REMOVED***
		return ""
	***REMOVED***
	return s.String
***REMOVED***

// NewString creates a new String
func NewString(s string, valid bool) String ***REMOVED***
	return String***REMOVED***
		NullString: sql.NullString***REMOVED***
			String: s,
			Valid:  valid,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.NullString.
func (s *String) UnmarshalJSON(data []byte) error ***REMOVED***
	var err error
	var v interface***REMOVED******REMOVED***
	if err = json.Unmarshal(data, &v); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch x := v.(type) ***REMOVED***
	case string:
		s.String = x
	case map[string]interface***REMOVED******REMOVED***:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.String", reflect.TypeOf(v).Name())
	***REMOVED***
	s.Valid = err == nil
	return err
***REMOVED***

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (s String) MarshalJSON() ([]byte, error) ***REMOVED***
	if !s.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	return json.Marshal(s.String)
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (s String) MarshalText() ([]byte, error) ***REMOVED***
	if !s.Valid ***REMOVED***
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***
	return []byte(s.String), nil
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error ***REMOVED***
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
***REMOVED***

// SetValid changes this String's value and also sets it to be non-null.
func (s *String) SetValid(v string) ***REMOVED***
	s.String = v
	s.Valid = true
***REMOVED***

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string ***REMOVED***
	if !s.Valid ***REMOVED***
		return nil
	***REMOVED***
	return &s.String
***REMOVED***

// IsZero returns true for null strings, for potential future omitempty support.
func (s String) IsZero() bool ***REMOVED***
	return !s.Valid
***REMOVED***
