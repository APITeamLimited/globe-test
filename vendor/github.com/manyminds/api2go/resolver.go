package api2go

import "net/http"

type callbackResolver struct ***REMOVED***
	callback func(r http.Request) string
	r        http.Request
***REMOVED***

// NewCallbackResolver handles each resolve via
// your provided callback func
func NewCallbackResolver(callback func(http.Request) string) URLResolver ***REMOVED***
	return &callbackResolver***REMOVED***callback: callback***REMOVED***
***REMOVED***

// GetBaseURL calls the callback given in the constructor method
// to implement `URLResolver`
func (c callbackResolver) GetBaseURL() string ***REMOVED***
	return c.callback(c.r)
***REMOVED***

// SetRequest to implement `RequestAwareURLResolver`
func (c *callbackResolver) SetRequest(r http.Request) ***REMOVED***
	c.r = r
***REMOVED***

// staticResolver is only used
// for backwards compatible reasons
// and might be removed in the future
type staticResolver struct ***REMOVED***
	baseURL string
***REMOVED***

func (s staticResolver) GetBaseURL() string ***REMOVED***
	return s.baseURL
***REMOVED***

// NewStaticResolver returns a simple resolver that
// will always answer with the same url
func NewStaticResolver(baseURL string) URLResolver ***REMOVED***
	return &staticResolver***REMOVED***baseURL: baseURL***REMOVED***
***REMOVED***
