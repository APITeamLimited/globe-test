// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements sql.Scanner so UUIDs can be read from databases transparently.
// Currently, database types that map to string and []byte are supported. Please
// consult database-specific driver documentation for matching types.
func (uuid *UUID) Scan(src interface***REMOVED******REMOVED***) error ***REMOVED***
	switch src := src.(type) ***REMOVED***
	case nil:
		return nil

	case string:
		// if an empty UUID comes from a table, we return a null UUID
		if src == "" ***REMOVED***
			return nil
		***REMOVED***

		// see Parse for required string format
		u, err := Parse(src)
		if err != nil ***REMOVED***
			return fmt.Errorf("Scan: %v", err)
		***REMOVED***

		*uuid = u

	case []byte:
		// if an empty UUID comes from a table, we return a null UUID
		if len(src) == 0 ***REMOVED***
			return nil
		***REMOVED***

		// assumes a simple slice of bytes if 16 bytes
		// otherwise attempts to parse
		if len(src) != 16 ***REMOVED***
			return uuid.Scan(string(src))
		***REMOVED***
		copy((*uuid)[:], src)

	default:
		return fmt.Errorf("Scan: unable to scan type %T into UUID", src)
	***REMOVED***

	return nil
***REMOVED***

// Value implements sql.Valuer so that UUIDs can be written to databases
// transparently. Currently, UUIDs map to strings. Please consult
// database-specific driver documentation for matching types.
func (uuid UUID) Value() (driver.Value, error) ***REMOVED***
	return uuid.String(), nil
***REMOVED***
