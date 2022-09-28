// Copyright 2019 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"io"
	"strings"
)

// JoinMessages concatenates received messages to create a single io.Reader.
// The string term is appended to each message. The returned reader does not
// support concurrent calls to the Read method.
func JoinMessages(c *Conn, term string) io.Reader ***REMOVED***
	return &joinReader***REMOVED***c: c, term: term***REMOVED***
***REMOVED***

type joinReader struct ***REMOVED***
	c    *Conn
	term string
	r    io.Reader
***REMOVED***

func (r *joinReader) Read(p []byte) (int, error) ***REMOVED***
	if r.r == nil ***REMOVED***
		var err error
		_, r.r, err = r.c.NextReader()
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if r.term != "" ***REMOVED***
			r.r = io.MultiReader(r.r, strings.NewReader(r.term))
		***REMOVED***
	***REMOVED***
	n, err := r.r.Read(p)
	if err == io.EOF ***REMOVED***
		err = nil
		r.r = nil
	***REMOVED***
	return n, err
***REMOVED***
