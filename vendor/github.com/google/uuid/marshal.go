// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import "fmt"

// MarshalText implements encoding.TextMarshaler.
func (uuid UUID) MarshalText() ([]byte, error) ***REMOVED***
	var js [36]byte
	encodeHex(js[:], uuid)
	return js[:], nil
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
func (uuid *UUID) UnmarshalText(data []byte) error ***REMOVED***
	id, err := ParseBytes(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*uuid = id
	return nil
***REMOVED***

// MarshalBinary implements encoding.BinaryMarshaler.
func (uuid UUID) MarshalBinary() ([]byte, error) ***REMOVED***
	return uuid[:], nil
***REMOVED***

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (uuid *UUID) UnmarshalBinary(data []byte) error ***REMOVED***
	if len(data) != 16 ***REMOVED***
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	***REMOVED***
	copy(uuid[:], data)
	return nil
***REMOVED***
