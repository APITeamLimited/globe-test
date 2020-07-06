// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"encoding/json"
	"io"
)

// WriteJSON writes the JSON encoding of v as a message.
//
// Deprecated: Use c.WriteJSON instead.
func WriteJSON(c *Conn, v interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.WriteJSON(v)
***REMOVED***

// WriteJSON writes the JSON encoding of v as a message.
//
// See the documentation for encoding/json Marshal for details about the
// conversion of Go values to JSON.
func (c *Conn) WriteJSON(v interface***REMOVED******REMOVED***) error ***REMOVED***
	w, err := c.NextWriter(TextMessage)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err1 := json.NewEncoder(w).Encode(v)
	err2 := w.Close()
	if err1 != nil ***REMOVED***
		return err1
	***REMOVED***
	return err2
***REMOVED***

// ReadJSON reads the next JSON-encoded message from the connection and stores
// it in the value pointed to by v.
//
// Deprecated: Use c.ReadJSON instead.
func ReadJSON(c *Conn, v interface***REMOVED******REMOVED***) error ***REMOVED***
	return c.ReadJSON(v)
***REMOVED***

// ReadJSON reads the next JSON-encoded message from the connection and stores
// it in the value pointed to by v.
//
// See the documentation for the encoding/json Unmarshal function for details
// about the conversion of JSON to a Go value.
func (c *Conn) ReadJSON(v interface***REMOVED******REMOVED***) error ***REMOVED***
	_, r, err := c.NextReader()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = json.NewDecoder(r).Decode(v)
	if err == io.EOF ***REMOVED***
		// One value is expected in the message.
		err = io.ErrUnexpectedEOF
	***REMOVED***
	return err
***REMOVED***
