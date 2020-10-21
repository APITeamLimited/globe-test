/*
 *
 * Copyright 2020 gRPC authors.
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

package channelz

import (
	"fmt"

	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.Component("channelz")

// Info logs and adds a trace event if channelz is on.
func Info(l grpclog.DepthLoggerV2, id int64, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if IsOn() ***REMOVED***
		AddTraceEvent(l, id, 1, &TraceEventDesc***REMOVED***
			Desc:     fmt.Sprint(args...),
			Severity: CtINFO,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		l.InfoDepth(1, args...)
	***REMOVED***
***REMOVED***

// Infof logs and adds a trace event if channelz is on.
func Infof(l grpclog.DepthLoggerV2, id int64, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	msg := fmt.Sprintf(format, args...)
	if IsOn() ***REMOVED***
		AddTraceEvent(l, id, 1, &TraceEventDesc***REMOVED***
			Desc:     msg,
			Severity: CtINFO,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		l.InfoDepth(1, msg)
	***REMOVED***
***REMOVED***

// Warning logs and adds a trace event if channelz is on.
func Warning(l grpclog.DepthLoggerV2, id int64, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if IsOn() ***REMOVED***
		AddTraceEvent(l, id, 1, &TraceEventDesc***REMOVED***
			Desc:     fmt.Sprint(args...),
			Severity: CtWarning,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		l.WarningDepth(1, args...)
	***REMOVED***
***REMOVED***

// Warningf logs and adds a trace event if channelz is on.
func Warningf(l grpclog.DepthLoggerV2, id int64, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	msg := fmt.Sprintf(format, args...)
	if IsOn() ***REMOVED***
		AddTraceEvent(l, id, 1, &TraceEventDesc***REMOVED***
			Desc:     msg,
			Severity: CtWarning,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		l.WarningDepth(1, msg)
	***REMOVED***
***REMOVED***

// Error logs and adds a trace event if channelz is on.
func Error(l grpclog.DepthLoggerV2, id int64, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if IsOn() ***REMOVED***
		AddTraceEvent(l, id, 1, &TraceEventDesc***REMOVED***
			Desc:     fmt.Sprint(args...),
			Severity: CtError,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		l.ErrorDepth(1, args...)
	***REMOVED***
***REMOVED***

// Errorf logs and adds a trace event if channelz is on.
func Errorf(l grpclog.DepthLoggerV2, id int64, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	msg := fmt.Sprintf(format, args...)
	if IsOn() ***REMOVED***
		AddTraceEvent(l, id, 1, &TraceEventDesc***REMOVED***
			Desc:     msg,
			Severity: CtError,
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		l.ErrorDepth(1, msg)
	***REMOVED***
***REMOVED***
