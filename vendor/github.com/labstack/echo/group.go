package echo

import (
	"path"
)

type (
	// Group is a set of sub-routes for a specified route. It can be used for inner
	// routes that share a common middleware or functionality that should be separate
	// from the parent echo instance while still inheriting from it.
	Group struct ***REMOVED***
		prefix     string
		middleware []MiddlewareFunc
		echo       *Echo
	***REMOVED***
)

// Use implements `Echo#Use()` for sub-routes within the Group.
func (g *Group) Use(middleware ...MiddlewareFunc) ***REMOVED***
	g.middleware = append(g.middleware, middleware...)
	// Allow all requests to reach the group as they might get dropped if router
	// doesn't find a match, making none of the group middleware process.
	for _, p := range []string***REMOVED***"", "/*"***REMOVED*** ***REMOVED***
		g.echo.Any(path.Clean(g.prefix+p), func(c Context) error ***REMOVED***
			return NotFoundHandler(c)
		***REMOVED***, g.middleware...)
	***REMOVED***
***REMOVED***

// CONNECT implements `Echo#CONNECT()` for sub-routes within the Group.
func (g *Group) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(CONNECT, path, h, m...)
***REMOVED***

// DELETE implements `Echo#DELETE()` for sub-routes within the Group.
func (g *Group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(DELETE, path, h, m...)
***REMOVED***

// GET implements `Echo#GET()` for sub-routes within the Group.
func (g *Group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(GET, path, h, m...)
***REMOVED***

// HEAD implements `Echo#HEAD()` for sub-routes within the Group.
func (g *Group) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(HEAD, path, h, m...)
***REMOVED***

// OPTIONS implements `Echo#OPTIONS()` for sub-routes within the Group.
func (g *Group) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(OPTIONS, path, h, m...)
***REMOVED***

// PATCH implements `Echo#PATCH()` for sub-routes within the Group.
func (g *Group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(PATCH, path, h, m...)
***REMOVED***

// POST implements `Echo#POST()` for sub-routes within the Group.
func (g *Group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(POST, path, h, m...)
***REMOVED***

// PUT implements `Echo#PUT()` for sub-routes within the Group.
func (g *Group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(PUT, path, h, m...)
***REMOVED***

// TRACE implements `Echo#TRACE()` for sub-routes within the Group.
func (g *Group) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return g.Add(TRACE, path, h, m...)
***REMOVED***

// Any implements `Echo#Any()` for sub-routes within the Group.
func (g *Group) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route ***REMOVED***
	routes := make([]*Route, len(methods))
	for i, m := range methods ***REMOVED***
		routes[i] = g.Add(m, path, handler, middleware...)
	***REMOVED***
	return routes
***REMOVED***

// Match implements `Echo#Match()` for sub-routes within the Group.
func (g *Group) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route ***REMOVED***
	routes := make([]*Route, len(methods))
	for i, m := range methods ***REMOVED***
		routes[i] = g.Add(m, path, handler, middleware...)
	***REMOVED***
	return routes
***REMOVED***

// Group creates a new sub-group with prefix and optional sub-group-level middleware.
func (g *Group) Group(prefix string, middleware ...MiddlewareFunc) *Group ***REMOVED***
	m := []MiddlewareFunc***REMOVED******REMOVED***
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	return g.echo.Group(g.prefix+prefix, m...)
***REMOVED***

// Static implements `Echo#Static()` for sub-routes within the Group.
func (g *Group) Static(prefix, root string) ***REMOVED***
	static(g, prefix, root)
***REMOVED***

// File implements `Echo#File()` for sub-routes within the Group.
func (g *Group) File(path, file string) ***REMOVED***
	g.echo.File(g.prefix+path, file)
***REMOVED***

// Add implements `Echo#Add()` for sub-routes within the Group.
func (g *Group) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route ***REMOVED***
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := []MiddlewareFunc***REMOVED******REMOVED***
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	return g.echo.Add(method, g.prefix+path, handler, m...)
***REMOVED***
