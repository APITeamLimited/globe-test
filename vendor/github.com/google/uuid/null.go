// Copyright 2021 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var jsonNull = []byte("null")

// NullUUID represents a UUID that may be null.
// NullUUID implements the SQL driver.Scanner interface so
// it can be used as a scan destination:
//
//  var u uuid.NullUUID
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&u)
//  ...
//  if u.Valid ***REMOVED***
//     // use u.UUID
//  ***REMOVED*** else ***REMOVED***
//     // NULL value
//  ***REMOVED***
//
type NullUUID struct ***REMOVED***
	UUID  UUID
	Valid bool // Valid is true if UUID is not NULL
***REMOVED***

// Scan implements the SQL driver.Scanner interface.
func (nu *NullUUID) Scan(value interface***REMOVED******REMOVED***) error ***REMOVED***
	if value == nil ***REMOVED***
		nu.UUID, nu.Valid = Nil, false
		return nil
	***REMOVED***

	err := nu.UUID.Scan(value)
	if err != nil ***REMOVED***
		nu.Valid = false
		return err
	***REMOVED***

	nu.Valid = true
	return nil
***REMOVED***

// Value implements the driver Valuer interface.
func (nu NullUUID) Value() (driver.Value, error) ***REMOVED***
	if !nu.Valid ***REMOVED***
		return nil, nil
	***REMOVED***
	// Delegate to UUID Value function
	return nu.UUID.Value()
***REMOVED***

// MarshalBinary implements encoding.BinaryMarshaler.
func (nu NullUUID) MarshalBinary() ([]byte, error) ***REMOVED***
	if nu.Valid ***REMOVED***
		return nu.UUID[:], nil
	***REMOVED***

	return []byte(nil), nil
***REMOVED***

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (nu *NullUUID) UnmarshalBinary(data []byte) error ***REMOVED***
	if len(data) != 16 ***REMOVED***
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	***REMOVED***
	copy(nu.UUID[:], data)
	nu.Valid = true
	return nil
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
func (nu NullUUID) MarshalText() ([]byte, error) ***REMOVED***
	if nu.Valid ***REMOVED***
		return nu.UUID.MarshalText()
	***REMOVED***

	return jsonNull, nil
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
func (nu *NullUUID) UnmarshalText(data []byte) error ***REMOVED***
	id, err := ParseBytes(data)
	if err != nil ***REMOVED***
		nu.Valid = false
		return err
	***REMOVED***
	nu.UUID = id
	nu.Valid = true
	return nil
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (nu NullUUID) MarshalJSON() ([]byte, error) ***REMOVED***
	if nu.Valid ***REMOVED***
		return json.Marshal(nu.UUID)
	***REMOVED***

	return jsonNull, nil
***REMOVED***

// UnmarshalJSON implements json.Unmarshaler.
func (nu *NullUUID) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, jsonNull) ***REMOVED***
		*nu = NullUUID***REMOVED******REMOVED***
		return nil // valid null UUID
	***REMOVED***
	err := json.Unmarshal(data, &nu.UUID)
	nu.Valid = err == nil
	return err
***REMOVED***
