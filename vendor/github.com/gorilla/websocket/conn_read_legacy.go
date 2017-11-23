// Copyright 2016 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.5

package websocket

import "io"

func (c *Conn) read(n int) ([]byte, error) ***REMOVED***
	p, err := c.br.Peek(n)
	if err == io.EOF ***REMOVED***
		err = errUnexpectedEOF
	***REMOVED***
	if len(p) > 0 ***REMOVED***
		// advance over the bytes just read
		io.ReadFull(c.br, p)
	***REMOVED***
	return p, err
***REMOVED***
