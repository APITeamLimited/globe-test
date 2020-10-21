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

package grpcutil

import (
	"strconv"
	"time"
)

const maxTimeoutValue int64 = 100000000 - 1

// div does integer division and round-up the result. Note that this is
// equivalent to (d+r-1)/r but has less chance to overflow.
func div(d, r time.Duration) int64 ***REMOVED***
	if d%r > 0 ***REMOVED***
		return int64(d/r + 1)
	***REMOVED***
	return int64(d / r)
***REMOVED***

// EncodeDuration encodes the duration to the format grpc-timeout header
// accepts.
//
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
func EncodeDuration(t time.Duration) string ***REMOVED***
	// TODO: This is simplistic and not bandwidth efficient. Improve it.
	if t <= 0 ***REMOVED***
		return "0n"
	***REMOVED***
	if d := div(t, time.Nanosecond); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "n"
	***REMOVED***
	if d := div(t, time.Microsecond); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "u"
	***REMOVED***
	if d := div(t, time.Millisecond); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "m"
	***REMOVED***
	if d := div(t, time.Second); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "S"
	***REMOVED***
	if d := div(t, time.Minute); d <= maxTimeoutValue ***REMOVED***
		return strconv.FormatInt(d, 10) + "M"
	***REMOVED***
	// Note that maxTimeoutValue * time.Hour > MaxInt64.
	return strconv.FormatInt(div(t, time.Hour), 10) + "H"
***REMOVED***
