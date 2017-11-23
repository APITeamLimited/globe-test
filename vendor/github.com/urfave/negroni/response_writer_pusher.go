//+build go1.8

package negroni

import (
	"fmt"
	"net/http"
)

func (rw *responseWriter) Push(target string, opts *http.PushOptions) error ***REMOVED***
	pusher, ok := rw.ResponseWriter.(http.Pusher)
	if ok ***REMOVED***
		return pusher.Push(target, opts)
	***REMOVED***
	return fmt.Errorf("the ResponseWriter doesn't support the Pusher interface")
***REMOVED***
