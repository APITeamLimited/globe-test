package api2go

import (
	"context"
	"time"
)

// APIContextAllocatorFunc to allow custom context implementations
type APIContextAllocatorFunc func(*API) APIContexter

// APIContexter embedding context.Context and requesting two helper functions
type APIContexter interface ***REMOVED***
	context.Context
	Set(key string, value interface***REMOVED******REMOVED***)
	Get(key string) (interface***REMOVED******REMOVED***, bool)
	Reset()
***REMOVED***

// APIContext api2go context for handlers, nil implementations related to Deadline and Done.
type APIContext struct ***REMOVED***
	keys map[string]interface***REMOVED******REMOVED***
***REMOVED***

// Set a string key value in the context
func (c *APIContext) Set(key string, value interface***REMOVED******REMOVED***) ***REMOVED***
	if c.keys == nil ***REMOVED***
		c.keys = make(map[string]interface***REMOVED******REMOVED***)
	***REMOVED***
	c.keys[key] = value
***REMOVED***

// Get a key value from the context
func (c *APIContext) Get(key string) (value interface***REMOVED******REMOVED***, exists bool) ***REMOVED***
	if c.keys != nil ***REMOVED***
		value, exists = c.keys[key]
	***REMOVED***
	return
***REMOVED***

// Reset resets all values on Context, making it safe to reuse
func (c *APIContext) Reset() ***REMOVED***
	c.keys = nil
***REMOVED***

// Deadline implements net/context
func (c *APIContext) Deadline() (deadline time.Time, ok bool) ***REMOVED***
	return
***REMOVED***

// Done implements net/context
func (c *APIContext) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

// Err implements net/context
func (c *APIContext) Err() error ***REMOVED***
	return nil
***REMOVED***

// Value implements net/context
func (c *APIContext) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if keyAsString, ok := key.(string); ok ***REMOVED***
		val, _ := c.Get(keyAsString)
		return val
	***REMOVED***
	return nil
***REMOVED***

// Compile time check
var _ APIContexter = &APIContext***REMOVED******REMOVED***

// ContextQueryParams fetches the QueryParams if Set
func ContextQueryParams(c *APIContext) map[string][]string ***REMOVED***
	qp, ok := c.Get("QueryParams")
	if ok == false ***REMOVED***
		qp = make(map[string][]string)
		c.Set("QueryParams", qp)
	***REMOVED***
	return qp.(map[string][]string)
***REMOVED***
