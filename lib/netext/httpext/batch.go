/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package httpext

import (
	"context"
	"sync/atomic"

	"github.com/loadimpact/k6/lib"
)

// BatchParsedHTTPRequest extends the normal parsed HTTP request with a pointer
// to a Response object, so that the batch goroutines can concurrently store the
// responses they receive, without any locking.
type BatchParsedHTTPRequest struct ***REMOVED***
	*ParsedHTTPRequest
	Response *Response // this is modified by MakeBatchRequests()
***REMOVED***

// MakeBatchRequests concurrently makes multiple requests. It spawns
// min(reqCount, globalLimit) goroutines that asynchronously process all
// requests coming from the requests channel. Responses are recorded in the
// pointers contained in each BatchParsedHTTPRequest object, so they need to be
// pre-initialized. In addition, each processed request would emit either a nil
// value, or an error, via the returned errors channel. The goroutines exit when
// the requests channel is closed.
func MakeBatchRequests(
	ctx context.Context,
	requests []BatchParsedHTTPRequest,
	reqCount, globalLimit, perHostLimit int,
) <-chan error ***REMOVED***
	workers := globalLimit
	if reqCount < workers ***REMOVED***
		workers = reqCount
	***REMOVED***
	result := make(chan error, reqCount)
	perHostLimiter := lib.NewMultiSlotLimiter(perHostLimit)

	makeRequest := func(req BatchParsedHTTPRequest) ***REMOVED***
		if hl := perHostLimiter.Slot(req.URL.GetURL().Host); hl != nil ***REMOVED***
			hl.Begin()
			defer hl.End()
		***REMOVED***

		resp, err := MakeRequest(ctx, req.ParsedHTTPRequest)
		if resp != nil ***REMOVED***
			*req.Response = *resp
		***REMOVED***
		result <- err
	***REMOVED***

	counter, i32reqCount := int32(-1), int32(reqCount)
	for i := 0; i < workers; i++ ***REMOVED***
		go func() ***REMOVED***
			for ***REMOVED***
				reqNum := atomic.AddInt32(&counter, 1)
				if reqNum >= i32reqCount ***REMOVED***
					return
				***REMOVED***
				makeRequest(requests[reqNum])
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	return result
***REMOVED***
