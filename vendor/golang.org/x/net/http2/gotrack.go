// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Defensive debug-only utility to track that functions run on the
// goroutine that they're supposed to.

package http2

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var DebugGoroutines = os.Getenv("DEBUG_HTTP2_GOROUTINES") == "1"

type goroutineLock uint64

func newGoroutineLock() goroutineLock ***REMOVED***
	if !DebugGoroutines ***REMOVED***
		return 0
	***REMOVED***
	return goroutineLock(curGoroutineID())
***REMOVED***

func (g goroutineLock) check() ***REMOVED***
	if !DebugGoroutines ***REMOVED***
		return
	***REMOVED***
	if curGoroutineID() != uint64(g) ***REMOVED***
		panic("running on the wrong goroutine")
	***REMOVED***
***REMOVED***

func (g goroutineLock) checkNotOn() ***REMOVED***
	if !DebugGoroutines ***REMOVED***
		return
	***REMOVED***
	if curGoroutineID() == uint64(g) ***REMOVED***
		panic("running on the wrong goroutine")
	***REMOVED***
***REMOVED***

var goroutineSpace = []byte("goroutine ")

func curGoroutineID() uint64 ***REMOVED***
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 ***REMOVED***
		panic(fmt.Sprintf("No space found in %q", b))
	***REMOVED***
	b = b[:i]
	n, err := parseUintBytes(b, 10, 64)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	***REMOVED***
	return n
***REMOVED***

var littleBuf = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		buf := make([]byte, 64)
		return &buf
	***REMOVED***,
***REMOVED***

// parseUintBytes is like strconv.ParseUint, but using a []byte.
func parseUintBytes(s []byte, base int, bitSize int) (n uint64, err error) ***REMOVED***
	var cutoff, maxVal uint64

	if bitSize == 0 ***REMOVED***
		bitSize = int(strconv.IntSize)
	***REMOVED***

	s0 := s
	switch ***REMOVED***
	case len(s) < 1:
		err = strconv.ErrSyntax
		goto Error

	case 2 <= base && base <= 36:
		// valid base; nothing to do

	case base == 0:
		// Look for octal, hex prefix.
		switch ***REMOVED***
		case s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X'):
			base = 16
			s = s[2:]
			if len(s) < 1 ***REMOVED***
				err = strconv.ErrSyntax
				goto Error
			***REMOVED***
		case s[0] == '0':
			base = 8
		default:
			base = 10
		***REMOVED***

	default:
		err = errors.New("invalid base " + strconv.Itoa(base))
		goto Error
	***REMOVED***

	n = 0
	cutoff = cutoff64(base)
	maxVal = 1<<uint(bitSize) - 1

	for i := 0; i < len(s); i++ ***REMOVED***
		var v byte
		d := s[i]
		switch ***REMOVED***
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			n = 0
			err = strconv.ErrSyntax
			goto Error
		***REMOVED***
		if int(v) >= base ***REMOVED***
			n = 0
			err = strconv.ErrSyntax
			goto Error
		***REMOVED***

		if n >= cutoff ***REMOVED***
			// n*base overflows
			n = 1<<64 - 1
			err = strconv.ErrRange
			goto Error
		***REMOVED***
		n *= uint64(base)

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal ***REMOVED***
			// n+v overflows
			n = 1<<64 - 1
			err = strconv.ErrRange
			goto Error
		***REMOVED***
		n = n1
	***REMOVED***

	return n, nil

Error:
	return n, &strconv.NumError***REMOVED***Func: "ParseUint", Num: string(s0), Err: err***REMOVED***
***REMOVED***

// Return the first number n such that n*base >= 1<<64.
func cutoff64(base int) uint64 ***REMOVED***
	if base < 2 ***REMOVED***
		return 0
	***REMOVED***
	return (1<<64-1)/uint64(base) + 1
***REMOVED***
