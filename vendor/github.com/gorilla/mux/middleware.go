package mux

import "net/http"

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc func(http.Handler) http.Handler

// middleware interface is anything which implements a MiddlewareFunc named Middleware.
type middleware interface ***REMOVED***
	Middleware(handler http.Handler) http.Handler
***REMOVED***

// MiddlewareFunc also implements the middleware interface.
func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler ***REMOVED***
	return mw(handler)
***REMOVED***

// Use appends a MiddlewareFunc to the chain. Middleware can be used to intercept or otherwise modify requests and/or responses, and are executed in the order that they are applied to the Router.
func (r *Router) Use(mwf MiddlewareFunc) ***REMOVED***
	r.middlewares = append(r.middlewares, mwf)
***REMOVED***

// useInterface appends a middleware to the chain. Middleware can be used to intercept or otherwise modify requests and/or responses, and are executed in the order that they are applied to the Router.
func (r *Router) useInterface(mw middleware) ***REMOVED***
	r.middlewares = append(r.middlewares, mw)
***REMOVED***
