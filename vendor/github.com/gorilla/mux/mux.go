// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
)

var (
	ErrMethodMismatch = errors.New("method is not allowed")
	ErrNotFound       = errors.New("no matching route was found")
)

// NewRouter returns a new router instance.
func NewRouter() *Router ***REMOVED***
	return &Router***REMOVED***namedRoutes: make(map[string]*Route), KeepContext: false***REMOVED***
***REMOVED***

// Router registers routes to be matched and dispatches a handler.
//
// It implements the http.Handler interface, so it can be registered to serve
// requests:
//
//     var router = mux.NewRouter()
//
//     func main() ***REMOVED***
//         http.Handle("/", router)
//     ***REMOVED***
//
// Or, for Google App Engine, register it in a init() function:
//
//     func init() ***REMOVED***
//         http.Handle("/", router)
//     ***REMOVED***
//
// This will send all incoming requests to the router.
type Router struct ***REMOVED***
	// Configurable Handler to be used when no route matches.
	NotFoundHandler http.Handler

	// Configurable Handler to be used when the request method does not match the route.
	MethodNotAllowedHandler http.Handler

	// Parent route, if this is a subrouter.
	parent parentRoute
	// Routes to be matched, in order.
	routes []*Route
	// Routes by name for URL building.
	namedRoutes map[string]*Route
	// See Router.StrictSlash(). This defines the flag for new routes.
	strictSlash bool
	// See Router.SkipClean(). This defines the flag for new routes.
	skipClean bool
	// If true, do not clear the request context after handling the request.
	// This has no effect when go1.7+ is used, since the context is stored
	// on the request itself.
	KeepContext bool
	// see Router.UseEncodedPath(). This defines a flag for all routes.
	useEncodedPath bool
	// Slice of middlewares to be called after a match is found
	middlewares []middleware
***REMOVED***

// Match attempts to match the given request against the router's registered routes.
//
// If the request matches a route of this router or one of its subrouters the Route,
// Handler, and Vars fields of the the match argument are filled and this function
// returns true.
//
// If the request does not match any of this router's or its subrouters' routes
// then this function returns false. If available, a reason for the match failure
// will be filled in the match argument's MatchErr field. If the match failure type
// (eg: not found) has a registered handler, the handler is assigned to the Handler
// field of the match argument.
func (r *Router) Match(req *http.Request, match *RouteMatch) bool ***REMOVED***
	for _, route := range r.routes ***REMOVED***
		if route.Match(req, match) ***REMOVED***
			// Build middleware chain if no error was found
			if match.MatchErr == nil ***REMOVED***
				for i := len(r.middlewares) - 1; i >= 0; i-- ***REMOVED***
					match.Handler = r.middlewares[i].Middleware(match.Handler)
				***REMOVED***
			***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	if match.MatchErr == ErrMethodMismatch ***REMOVED***
		if r.MethodNotAllowedHandler != nil ***REMOVED***
			match.Handler = r.MethodNotAllowedHandler
			return true
		***REMOVED*** else ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// Closest match for a router (includes sub-routers)
	if r.NotFoundHandler != nil ***REMOVED***
		match.Handler = r.NotFoundHandler
		match.MatchErr = ErrNotFound
		return true
	***REMOVED***

	match.MatchErr = ErrNotFound
	return false
***REMOVED***

// ServeHTTP dispatches the handler registered in the matched route.
//
// When there is a match, the route variables can be retrieved calling
// mux.Vars(request).
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) ***REMOVED***
	if !r.skipClean ***REMOVED***
		path := req.URL.Path
		if r.useEncodedPath ***REMOVED***
			path = req.URL.EscapedPath()
		***REMOVED***
		// Clean path to canonical form and redirect.
		if p := cleanPath(path); p != path ***REMOVED***

			// Added 3 lines (Philip Schlump) - It was dropping the query string and #whatever from query.
			// This matches with fix in go 1.2 r.c. 4 for same problem.  Go Issue:
			// http://code.google.com/p/go/issues/detail?id=5252
			url := *req.URL
			url.Path = p
			p = url.String()

			w.Header().Set("Location", p)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		***REMOVED***
	***REMOVED***
	var match RouteMatch
	var handler http.Handler
	if r.Match(req, &match) ***REMOVED***
		handler = match.Handler
		req = setVars(req, match.Vars)
		req = setCurrentRoute(req, match.Route)
	***REMOVED***

	if handler == nil && match.MatchErr == ErrMethodMismatch ***REMOVED***
		handler = methodNotAllowedHandler()
	***REMOVED***

	if handler == nil ***REMOVED***
		handler = http.NotFoundHandler()
	***REMOVED***

	if !r.KeepContext ***REMOVED***
		defer contextClear(req)
	***REMOVED***

	handler.ServeHTTP(w, req)
***REMOVED***

// Get returns a route registered with the given name.
func (r *Router) Get(name string) *Route ***REMOVED***
	return r.getNamedRoutes()[name]
***REMOVED***

// GetRoute returns a route registered with the given name. This method
// was renamed to Get() and remains here for backwards compatibility.
func (r *Router) GetRoute(name string) *Route ***REMOVED***
	return r.getNamedRoutes()[name]
***REMOVED***

// StrictSlash defines the trailing slash behavior for new routes. The initial
// value is false.
//
// When true, if the route path is "/path/", accessing "/path" will perform a redirect
// to the former and vice versa. In other words, your application will always
// see the path as specified in the route.
//
// When false, if the route path is "/path", accessing "/path/" will not match
// this route and vice versa.
//
// The re-direct is a HTTP 301 (Moved Permanently). Note that when this is set for
// routes with a non-idempotent method (e.g. POST, PUT), the subsequent re-directed
// request will be made as a GET by most clients. Use middleware or client settings
// to modify this behaviour as needed.
//
// Special case: when a route sets a path prefix using the PathPrefix() method,
// strict slash is ignored for that route because the redirect behavior can't
// be determined from a prefix alone. However, any subrouters created from that
// route inherit the original StrictSlash setting.
func (r *Router) StrictSlash(value bool) *Router ***REMOVED***
	r.strictSlash = value
	return r
***REMOVED***

// SkipClean defines the path cleaning behaviour for new routes. The initial
// value is false. Users should be careful about which routes are not cleaned
//
// When true, if the route path is "/path//to", it will remain with the double
// slash. This is helpful if you have a route like: /fetch/http://xkcd.com/534/
//
// When false, the path will be cleaned, so /fetch/http://xkcd.com/534/ will
// become /fetch/http/xkcd.com/534
func (r *Router) SkipClean(value bool) *Router ***REMOVED***
	r.skipClean = value
	return r
***REMOVED***

// UseEncodedPath tells the router to match the encoded original path
// to the routes.
// For eg. "/path/foo%2Fbar/to" will match the path "/path/***REMOVED***var***REMOVED***/to".
//
// If not called, the router will match the unencoded path to the routes.
// For eg. "/path/foo%2Fbar/to" will match the path "/path/foo/bar/to"
func (r *Router) UseEncodedPath() *Router ***REMOVED***
	r.useEncodedPath = true
	return r
***REMOVED***

// ----------------------------------------------------------------------------
// parentRoute
// ----------------------------------------------------------------------------

func (r *Router) getBuildScheme() string ***REMOVED***
	if r.parent != nil ***REMOVED***
		return r.parent.getBuildScheme()
	***REMOVED***
	return ""
***REMOVED***

// getNamedRoutes returns the map where named routes are registered.
func (r *Router) getNamedRoutes() map[string]*Route ***REMOVED***
	if r.namedRoutes == nil ***REMOVED***
		if r.parent != nil ***REMOVED***
			r.namedRoutes = r.parent.getNamedRoutes()
		***REMOVED*** else ***REMOVED***
			r.namedRoutes = make(map[string]*Route)
		***REMOVED***
	***REMOVED***
	return r.namedRoutes
***REMOVED***

// getRegexpGroup returns regexp definitions from the parent route, if any.
func (r *Router) getRegexpGroup() *routeRegexpGroup ***REMOVED***
	if r.parent != nil ***REMOVED***
		return r.parent.getRegexpGroup()
	***REMOVED***
	return nil
***REMOVED***

func (r *Router) buildVars(m map[string]string) map[string]string ***REMOVED***
	if r.parent != nil ***REMOVED***
		m = r.parent.buildVars(m)
	***REMOVED***
	return m
***REMOVED***

// ----------------------------------------------------------------------------
// Route factories
// ----------------------------------------------------------------------------

// NewRoute registers an empty route.
func (r *Router) NewRoute() *Route ***REMOVED***
	route := &Route***REMOVED***parent: r, strictSlash: r.strictSlash, skipClean: r.skipClean, useEncodedPath: r.useEncodedPath***REMOVED***
	r.routes = append(r.routes, route)
	return route
***REMOVED***

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(path string, handler http.Handler) *Route ***REMOVED***
	return r.NewRoute().Path(path).Handler(handler)
***REMOVED***

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *Route ***REMOVED***
	return r.NewRoute().Path(path).HandlerFunc(f)
***REMOVED***

// Headers registers a new route with a matcher for request header values.
// See Route.Headers().
func (r *Router) Headers(pairs ...string) *Route ***REMOVED***
	return r.NewRoute().Headers(pairs...)
***REMOVED***

// Host registers a new route with a matcher for the URL host.
// See Route.Host().
func (r *Router) Host(tpl string) *Route ***REMOVED***
	return r.NewRoute().Host(tpl)
***REMOVED***

// MatcherFunc registers a new route with a custom matcher function.
// See Route.MatcherFunc().
func (r *Router) MatcherFunc(f MatcherFunc) *Route ***REMOVED***
	return r.NewRoute().MatcherFunc(f)
***REMOVED***

// Methods registers a new route with a matcher for HTTP methods.
// See Route.Methods().
func (r *Router) Methods(methods ...string) *Route ***REMOVED***
	return r.NewRoute().Methods(methods...)
***REMOVED***

// Path registers a new route with a matcher for the URL path.
// See Route.Path().
func (r *Router) Path(tpl string) *Route ***REMOVED***
	return r.NewRoute().Path(tpl)
***REMOVED***

// PathPrefix registers a new route with a matcher for the URL path prefix.
// See Route.PathPrefix().
func (r *Router) PathPrefix(tpl string) *Route ***REMOVED***
	return r.NewRoute().PathPrefix(tpl)
***REMOVED***

// Queries registers a new route with a matcher for URL query values.
// See Route.Queries().
func (r *Router) Queries(pairs ...string) *Route ***REMOVED***
	return r.NewRoute().Queries(pairs...)
***REMOVED***

// Schemes registers a new route with a matcher for URL schemes.
// See Route.Schemes().
func (r *Router) Schemes(schemes ...string) *Route ***REMOVED***
	return r.NewRoute().Schemes(schemes...)
***REMOVED***

// BuildVarsFunc registers a new route with a custom function for modifying
// route variables before building a URL.
func (r *Router) BuildVarsFunc(f BuildVarsFunc) *Route ***REMOVED***
	return r.NewRoute().BuildVarsFunc(f)
***REMOVED***

// Walk walks the router and all its sub-routers, calling walkFn for each route
// in the tree. The routes are walked in the order they were added. Sub-routers
// are explored depth-first.
func (r *Router) Walk(walkFn WalkFunc) error ***REMOVED***
	return r.walk(walkFn, []*Route***REMOVED******REMOVED***)
***REMOVED***

// SkipRouter is used as a return value from WalkFuncs to indicate that the
// router that walk is about to descend down to should be skipped.
var SkipRouter = errors.New("skip this router")

// WalkFunc is the type of the function called for each route visited by Walk.
// At every invocation, it is given the current route, and the current router,
// and a list of ancestor routes that lead to the current route.
type WalkFunc func(route *Route, router *Router, ancestors []*Route) error

func (r *Router) walk(walkFn WalkFunc, ancestors []*Route) error ***REMOVED***
	for _, t := range r.routes ***REMOVED***
		err := walkFn(t, r, ancestors)
		if err == SkipRouter ***REMOVED***
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, sr := range t.matchers ***REMOVED***
			if h, ok := sr.(*Router); ok ***REMOVED***
				ancestors = append(ancestors, t)
				err := h.walk(walkFn, ancestors)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				ancestors = ancestors[:len(ancestors)-1]
			***REMOVED***
		***REMOVED***
		if h, ok := t.handler.(*Router); ok ***REMOVED***
			ancestors = append(ancestors, t)
			err := h.walk(walkFn, ancestors)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			ancestors = ancestors[:len(ancestors)-1]
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// ----------------------------------------------------------------------------
// Context
// ----------------------------------------------------------------------------

// RouteMatch stores information about a matched route.
type RouteMatch struct ***REMOVED***
	Route   *Route
	Handler http.Handler
	Vars    map[string]string

	// MatchErr is set to appropriate matching error
	// It is set to ErrMethodMismatch if there is a mismatch in
	// the request method and route method
	MatchErr error
***REMOVED***

type contextKey int

const (
	varsKey contextKey = iota
	routeKey
)

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string ***REMOVED***
	if rv := contextGet(r, varsKey); rv != nil ***REMOVED***
		return rv.(map[string]string)
	***REMOVED***
	return nil
***REMOVED***

// CurrentRoute returns the matched route for the current request, if any.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns, unless the KeepContext option is set on the
// Router.
func CurrentRoute(r *http.Request) *Route ***REMOVED***
	if rv := contextGet(r, routeKey); rv != nil ***REMOVED***
		return rv.(*Route)
	***REMOVED***
	return nil
***REMOVED***

func setVars(r *http.Request, val interface***REMOVED******REMOVED***) *http.Request ***REMOVED***
	return contextSet(r, varsKey, val)
***REMOVED***

func setCurrentRoute(r *http.Request, val interface***REMOVED******REMOVED***) *http.Request ***REMOVED***
	return contextSet(r, routeKey, val)
***REMOVED***

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string ***REMOVED***
	if p == "" ***REMOVED***
		return "/"
	***REMOVED***
	if p[0] != '/' ***REMOVED***
		p = "/" + p
	***REMOVED***
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" ***REMOVED***
		np += "/"
	***REMOVED***

	return np
***REMOVED***

// uniqueVars returns an error if two slices contain duplicated strings.
func uniqueVars(s1, s2 []string) error ***REMOVED***
	for _, v1 := range s1 ***REMOVED***
		for _, v2 := range s2 ***REMOVED***
			if v1 == v2 ***REMOVED***
				return fmt.Errorf("mux: duplicated route variable %q", v2)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// checkPairs returns the count of strings passed in, and an error if
// the count is not an even number.
func checkPairs(pairs ...string) (int, error) ***REMOVED***
	length := len(pairs)
	if length%2 != 0 ***REMOVED***
		return length, fmt.Errorf(
			"mux: number of parameters must be multiple of 2, got %v", pairs)
	***REMOVED***
	return length, nil
***REMOVED***

// mapFromPairsToString converts variadic string parameters to a
// string to string map.
func mapFromPairsToString(pairs ...string) (map[string]string, error) ***REMOVED***
	length, err := checkPairs(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 ***REMOVED***
		m[pairs[i]] = pairs[i+1]
	***REMOVED***
	return m, nil
***REMOVED***

// mapFromPairsToRegex converts variadic string parameters to a
// string to regex map.
func mapFromPairsToRegex(pairs ...string) (map[string]*regexp.Regexp, error) ***REMOVED***
	length, err := checkPairs(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m := make(map[string]*regexp.Regexp, length/2)
	for i := 0; i < length; i += 2 ***REMOVED***
		regex, err := regexp.Compile(pairs[i+1])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m[pairs[i]] = regex
	***REMOVED***
	return m, nil
***REMOVED***

// matchInArray returns true if the given string value is in the array.
func matchInArray(arr []string, value string) bool ***REMOVED***
	for _, v := range arr ***REMOVED***
		if v == value ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// matchMapWithString returns true if the given key/value pairs exist in a given map.
func matchMapWithString(toCheck map[string]string, toMatch map[string][]string, canonicalKey bool) bool ***REMOVED***
	for k, v := range toCheck ***REMOVED***
		// Check if key exists.
		if canonicalKey ***REMOVED***
			k = http.CanonicalHeaderKey(k)
		***REMOVED***
		if values := toMatch[k]; values == nil ***REMOVED***
			return false
		***REMOVED*** else if v != "" ***REMOVED***
			// If value was defined as an empty string we only check that the
			// key exists. Otherwise we also check for equality.
			valueExists := false
			for _, value := range values ***REMOVED***
				if v == value ***REMOVED***
					valueExists = true
					break
				***REMOVED***
			***REMOVED***
			if !valueExists ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// matchMapWithRegex returns true if the given key/value pairs exist in a given map compiled against
// the given regex
func matchMapWithRegex(toCheck map[string]*regexp.Regexp, toMatch map[string][]string, canonicalKey bool) bool ***REMOVED***
	for k, v := range toCheck ***REMOVED***
		// Check if key exists.
		if canonicalKey ***REMOVED***
			k = http.CanonicalHeaderKey(k)
		***REMOVED***
		if values := toMatch[k]; values == nil ***REMOVED***
			return false
		***REMOVED*** else if v != nil ***REMOVED***
			// If value was defined as an empty string we only check that the
			// key exists. Otherwise we also check for equality.
			valueExists := false
			for _, value := range values ***REMOVED***
				if v.MatchString(value) ***REMOVED***
					valueExists = true
					break
				***REMOVED***
			***REMOVED***
			if !valueExists ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// methodNotAllowed replies to the request with an HTTP status code 405.
func methodNotAllowed(w http.ResponseWriter, r *http.Request) ***REMOVED***
	w.WriteHeader(http.StatusMethodNotAllowed)
***REMOVED***

// methodNotAllowedHandler returns a simple request handler
// that replies to each request with a status code 405.
func methodNotAllowedHandler() http.Handler ***REMOVED*** return http.HandlerFunc(methodNotAllowed) ***REMOVED***
