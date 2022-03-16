// Copyright 2016 The Gorilla WebSocket Authors. All rights reserved.  Use of
// this source code is governed by a BSD-style license that can be found in the
// LICENSE file.

//go:build !appengine
// +build !appengine

package websocket

import "unsafe"

const wordSize = int(unsafe.Sizeof(uintptr(0)))

func maskBytes(key [4]byte, pos int, b []byte) int ***REMOVED***
	// Mask one byte at a time for small buffers.
	if len(b) < 2*wordSize ***REMOVED***
		for i := range b ***REMOVED***
			b[i] ^= key[pos&3]
			pos++
		***REMOVED***
		return pos & 3
	***REMOVED***

	// Mask one byte at a time to word boundary.
	if n := int(uintptr(unsafe.Pointer(&b[0]))) % wordSize; n != 0 ***REMOVED***
		n = wordSize - n
		for i := range b[:n] ***REMOVED***
			b[i] ^= key[pos&3]
			pos++
		***REMOVED***
		b = b[n:]
	***REMOVED***

	// Create aligned word size key.
	var k [wordSize]byte
	for i := range k ***REMOVED***
		k[i] = key[(pos+i)&3]
	***REMOVED***
	kw := *(*uintptr)(unsafe.Pointer(&k))

	// Mask one word at a time.
	n := (len(b) / wordSize) * wordSize
	for i := 0; i < n; i += wordSize ***REMOVED***
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(i))) ^= kw
	***REMOVED***

	// Mask one byte at a time for remaining bytes.
	b = b[n:]
	for i := range b ***REMOVED***
		b[i] ^= key[pos&3]
		pos++
	***REMOVED***

	return pos & 3
***REMOVED***
