// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.11

package http2

import (
	"net/http/httptrace"
	"net/textproto"
)

func traceHasWroteHeaderField(trace *httptrace.ClientTrace) bool ***REMOVED***
	return trace != nil && trace.WroteHeaderField != nil
***REMOVED***

func traceWroteHeaderField(trace *httptrace.ClientTrace, k, v string) ***REMOVED***
	if trace != nil && trace.WroteHeaderField != nil ***REMOVED***
		trace.WroteHeaderField(k, []string***REMOVED***v***REMOVED***)
	***REMOVED***
***REMOVED***

func traceGot1xxResponseFunc(trace *httptrace.ClientTrace) func(int, textproto.MIMEHeader) error ***REMOVED***
	if trace != nil ***REMOVED***
		return trace.Got1xxResponse
	***REMOVED***
	return nil
***REMOVED***
