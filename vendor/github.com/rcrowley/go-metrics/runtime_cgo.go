// +build cgo
// +build !appengine

package metrics

import "runtime"

func numCgoCall() int64 ***REMOVED***
	return runtime.NumCgoCall()
***REMOVED***
