package null

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Bool is a nullable bool.
// It does not consider false values to be null.
// It will decode to null, not false, if null.
type Bool struct ***REMOVED***
	sql.NullBool
***REMOVED***

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool ***REMOVED***
	return Bool***REMOVED***
		NullBool: sql.NullBool***REMOVED***
			Bool:  b,
			Valid: valid,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// BoolFrom creates a new Bool that will always be valid.
func BoolFrom(b bool) Bool ***REMOVED***
	return NewBool(b, true)
***REMOVED***

// BoolFromPtr creates a new Bool that will be null if f is nil.
func BoolFromPtr(b *bool) Bool ***REMOVED***
	if b == nil ***REMOVED***
		return NewBool(false, false)
	***REMOVED***
	return NewBool(*b, true)
***REMOVED***

// ValueOrZero returns the inner value if valid, otherwise false.
func (b Bool) ValueOrZero() bool ***REMOVED***
	return b.Valid && b.Bool
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Bool.
// It also supports unmarshalling a sql.NullBool.
func (b *Bool) UnmarshalJSON(data []byte) error ***REMOVED***
	var err error
	var v interface***REMOVED******REMOVED***
	if err = json.Unmarshal(data, &v); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch x := v.(type) ***REMOVED***
	case bool:
		b.Bool = x
	case map[string]interface***REMOVED******REMOVED***:
		err = json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Bool", reflect.TypeOf(v).Name())
	***REMOVED***
	b.Valid = err == nil
	return err
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (b *Bool) UnmarshalText(text []byte) error ***REMOVED***
	str := string(text)
	switch str ***REMOVED***
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	***REMOVED***
	b.Valid = true
	return nil
***REMOVED***

// MarshalJSON implements json.Marshaler.
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) ***REMOVED***
	if !b.Valid ***REMOVED***
		return []byte("null"), nil
	***REMOVED***
	if !b.Bool ***REMOVED***
		return []byte("false"), nil
	***REMOVED***
	return []byte("true"), nil
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) ***REMOVED***
	if !b.Valid ***REMOVED***
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***
	if !b.Bool ***REMOVED***
		return []byte("false"), nil
	***REMOVED***
	return []byte("true"), nil
***REMOVED***

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) ***REMOVED***
	b.Bool = v
	b.Valid = true
***REMOVED***

// Ptr returns a pointer to this Bool's value, or a nil pointer if this Bool is null.
func (b Bool) Ptr() *bool ***REMOVED***
	if !b.Valid ***REMOVED***
		return nil
	***REMOVED***
	return &b.Bool
***REMOVED***

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
// A non-null Bool with a 0 value will not be considered zero.
func (b Bool) IsZero() bool ***REMOVED***
	return !b.Valid
***REMOVED***
