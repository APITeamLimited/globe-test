package httpext

import (
	"context"
	"sync/atomic"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
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
// The processResponse callback can be used to modify the response, e.g.
// to replace the body.
func MakeBatchRequests(
	ctx context.Context, state *libWorker.State,
	requests []BatchParsedHTTPRequest,
	reqCount, globalLimit, perHostLimit int,
	processResponse func(*Response, ResponseType),
) <-chan error ***REMOVED***
	workers := globalLimit
	if reqCount < workers ***REMOVED***
		workers = reqCount
	***REMOVED***
	result := make(chan error, reqCount)
	perHostLimiter := libWorker.NewMultiSlotLimiter(perHostLimit)

	makeRequest := func(req BatchParsedHTTPRequest) ***REMOVED***
		if hl := perHostLimiter.Slot(req.URL.GetURL().Host); hl != nil ***REMOVED***
			hl.Begin()
			defer hl.End()
		***REMOVED***

		resp, err := MakeRequest(ctx, state, req.ParsedHTTPRequest)
		if resp != nil ***REMOVED***
			processResponse(resp, req.ParsedHTTPRequest.ResponseType)
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
