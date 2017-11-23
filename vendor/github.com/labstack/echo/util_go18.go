// +build go1.8

package echo

import "net/url"

// PathUnescape is wraps `url.PathUnescape`
func PathUnescape(s string) (string, error) ***REMOVED***
	return url.PathUnescape(s)
***REMOVED***
