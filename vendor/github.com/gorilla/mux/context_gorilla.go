// +build !go1.7

package mux

import (
	"net/http"

	"github.com/gorilla/context"
)

func contextGet(r *http.Request, key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return context.Get(r, key)
***REMOVED***

func contextSet(r *http.Request, key, val interface***REMOVED******REMOVED***) *http.Request ***REMOVED***
	if val == nil ***REMOVED***
		return r
	***REMOVED***

	context.Set(r, key, val)
	return r
***REMOVED***

func contextClear(r *http.Request) ***REMOVED***
	context.Clear(r)
***REMOVED***
