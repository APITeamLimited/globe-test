// Copyright 2016 The Gorilla WebSocket Authors. All rights reserved.  Use of
// this source code is governed by a BSD-style license that can be found in the
// LICENSE file.

//go:build appengine
// +build appengine

package websocket

func maskBytes(key [4]byte, pos int, b []byte) int ***REMOVED***
	for i := range b ***REMOVED***
		b[i] ^= key[pos&3]
		pos++
	***REMOVED***
	return pos & 3
***REMOVED***
