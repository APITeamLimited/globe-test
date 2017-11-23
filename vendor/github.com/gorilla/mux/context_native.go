// +build go1.7

package mux

import (
	"context"
	"net/http"
)

func contextGet(r *http.Request, key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return r.Context().Value(key)
***REMOVED***

func contextSet(r *http.Request, key, val interface***REMOVED******REMOVED***) *http.Request ***REMOVED***
	if val == nil ***REMOVED***
		return r
	***REMOVED***

	return r.WithContext(context.WithValue(r.Context(), key, val))
***REMOVED***

func contextClear(r *http.Request) ***REMOVED***
	return
***REMOVED***
