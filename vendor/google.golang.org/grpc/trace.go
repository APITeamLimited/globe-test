/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/trace"
)

// EnableTracing controls whether to trace RPCs using the golang.org/x/net/trace package.
// This should only be set before any RPCs are sent or received by this program.
var EnableTracing bool

// methodFamily returns the trace family for the given method.
// It turns "/pkg.Service/GetFoo" into "pkg.Service".
func methodFamily(m string) string ***REMOVED***
	m = strings.TrimPrefix(m, "/") // remove leading slash
	if i := strings.Index(m, "/"); i >= 0 ***REMOVED***
		m = m[:i] // remove everything from second slash
	***REMOVED***
	return m
***REMOVED***

// traceInfo contains tracing information for an RPC.
type traceInfo struct ***REMOVED***
	tr        trace.Trace
	firstLine firstLine
***REMOVED***

// firstLine is the first line of an RPC trace.
// It may be mutated after construction; remoteAddr specifically may change
// during client-side use.
type firstLine struct ***REMOVED***
	mu         sync.Mutex
	client     bool // whether this is a client (outgoing) RPC
	remoteAddr net.Addr
	deadline   time.Duration // may be zero
***REMOVED***

func (f *firstLine) SetRemoteAddr(addr net.Addr) ***REMOVED***
	f.mu.Lock()
	f.remoteAddr = addr
	f.mu.Unlock()
***REMOVED***

func (f *firstLine) String() string ***REMOVED***
	f.mu.Lock()
	defer f.mu.Unlock()

	var line bytes.Buffer
	io.WriteString(&line, "RPC: ")
	if f.client ***REMOVED***
		io.WriteString(&line, "to")
	***REMOVED*** else ***REMOVED***
		io.WriteString(&line, "from")
	***REMOVED***
	fmt.Fprintf(&line, " %v deadline:", f.remoteAddr)
	if f.deadline != 0 ***REMOVED***
		fmt.Fprint(&line, f.deadline)
	***REMOVED*** else ***REMOVED***
		io.WriteString(&line, "none")
	***REMOVED***
	return line.String()
***REMOVED***

const truncateSize = 100

func truncate(x string, l int) string ***REMOVED***
	if l > len(x) ***REMOVED***
		return x
	***REMOVED***
	return x[:l]
***REMOVED***

// payload represents an RPC request or response payload.
type payload struct ***REMOVED***
	sent bool        // whether this is an outgoing payload
	msg  interface***REMOVED******REMOVED*** // e.g. a proto.Message
	// TODO(dsymonds): add stringifying info to codec, and limit how much we hold here?
***REMOVED***

func (p payload) String() string ***REMOVED***
	if p.sent ***REMOVED***
		return truncate(fmt.Sprintf("sent: %v", p.msg), truncateSize)
	***REMOVED***
	return truncate(fmt.Sprintf("recv: %v", p.msg), truncateSize)
***REMOVED***

type fmtStringer struct ***REMOVED***
	format string
	a      []interface***REMOVED******REMOVED***
***REMOVED***

func (f *fmtStringer) String() string ***REMOVED***
	return fmt.Sprintf(f.format, f.a...)
***REMOVED***

type stringer string

func (s stringer) String() string ***REMOVED*** return string(s) ***REMOVED***
