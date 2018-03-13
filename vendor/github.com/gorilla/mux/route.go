// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Route stores information to match a request and build URLs.
type Route struct ***REMOVED***
	// Parent where the route was registered (a Router).
	parent parentRoute
	// Request handler for the route.
	handler http.Handler
	// List of matchers.
	matchers []matcher
	// Manager for the variables from host and path.
	regexp *routeRegexpGroup
	// If true, when the path pattern is "/path/", accessing "/path" will
	// redirect to the former and vice versa.
	strictSlash bool
	// If true, when the path pattern is "/path//to", accessing "/path//to"
	// will not redirect
	skipClean bool
	// If true, "/path/foo%2Fbar/to" will match the path "/path/***REMOVED***var***REMOVED***/to"
	useEncodedPath bool
	// The scheme used when building URLs.
	buildScheme string
	// If true, this route never matches: it is only used to build URLs.
	buildOnly bool
	// The name used to build URLs.
	name string
	// Error resulted from building a route.
	err error

	buildVarsFunc BuildVarsFunc
***REMOVED***

func (r *Route) SkipClean() bool ***REMOVED***
	return r.skipClean
***REMOVED***

// Match matches the route against the request.
func (r *Route) Match(req *http.Request, match *RouteMatch) bool ***REMOVED***
	if r.buildOnly || r.err != nil ***REMOVED***
		return false
	***REMOVED***

	var matchErr error

	// Match everything.
	for _, m := range r.matchers ***REMOVED***
		if matched := m.Match(req, match); !matched ***REMOVED***
			if _, ok := m.(methodMatcher); ok ***REMOVED***
				matchErr = ErrMethodMismatch
				continue
			***REMOVED***
			matchErr = nil
			return false
		***REMOVED***
	***REMOVED***

	if matchErr != nil ***REMOVED***
		match.MatchErr = matchErr
		return false
	***REMOVED***

	if match.MatchErr == ErrMethodMismatch ***REMOVED***
		// We found a route which matches request method, clear MatchErr
		match.MatchErr = nil
		// Then override the mis-matched handler
		match.Handler = r.handler
	***REMOVED***

	// Yay, we have a match. Let's collect some info about it.
	if match.Route == nil ***REMOVED***
		match.Route = r
	***REMOVED***
	if match.Handler == nil ***REMOVED***
		match.Handler = r.handler
	***REMOVED***
	if match.Vars == nil ***REMOVED***
		match.Vars = make(map[string]string)
	***REMOVED***

	// Set variables.
	if r.regexp != nil ***REMOVED***
		r.regexp.setMatch(req, match, r)
	***REMOVED***
	return true
***REMOVED***

// ----------------------------------------------------------------------------
// Route attributes
// ----------------------------------------------------------------------------

// GetError returns an error resulted from building the route, if any.
func (r *Route) GetError() error ***REMOVED***
	return r.err
***REMOVED***

// BuildOnly sets the route to never match: it is only used to build URLs.
func (r *Route) BuildOnly() *Route ***REMOVED***
	r.buildOnly = true
	return r
***REMOVED***

// Handler --------------------------------------------------------------------

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route ***REMOVED***
	if r.err == nil ***REMOVED***
		r.handler = handler
	***REMOVED***
	return r
***REMOVED***

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route ***REMOVED***
	return r.Handler(http.HandlerFunc(f))
***REMOVED***

// GetHandler returns the handler for the route, if any.
func (r *Route) GetHandler() http.Handler ***REMOVED***
	return r.handler
***REMOVED***

// Name -----------------------------------------------------------------------

// Name sets the name for the route, used to build URLs.
// If the name was registered already it will be overwritten.
func (r *Route) Name(name string) *Route ***REMOVED***
	if r.name != "" ***REMOVED***
		r.err = fmt.Errorf("mux: route already has name %q, can't set %q",
			r.name, name)
	***REMOVED***
	if r.err == nil ***REMOVED***
		r.name = name
		r.getNamedRoutes()[name] = r
	***REMOVED***
	return r
***REMOVED***

// GetName returns the name for the route, if any.
func (r *Route) GetName() string ***REMOVED***
	return r.name
***REMOVED***

// ----------------------------------------------------------------------------
// Matchers
// ----------------------------------------------------------------------------

// matcher types try to match a request.
type matcher interface ***REMOVED***
	Match(*http.Request, *RouteMatch) bool
***REMOVED***

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m matcher) *Route ***REMOVED***
	if r.err == nil ***REMOVED***
		r.matchers = append(r.matchers, m)
	***REMOVED***
	return r
***REMOVED***

// addRegexpMatcher adds a host or path matcher and builder to a route.
func (r *Route) addRegexpMatcher(tpl string, typ regexpType) error ***REMOVED***
	if r.err != nil ***REMOVED***
		return r.err
	***REMOVED***
	r.regexp = r.getRegexpGroup()
	if typ == regexpTypePath || typ == regexpTypePrefix ***REMOVED***
		if len(tpl) > 0 && tpl[0] != '/' ***REMOVED***
			return fmt.Errorf("mux: path must start with a slash, got %q", tpl)
		***REMOVED***
		if r.regexp.path != nil ***REMOVED***
			tpl = strings.TrimRight(r.regexp.path.template, "/") + tpl
		***REMOVED***
	***REMOVED***
	rr, err := newRouteRegexp(tpl, typ, routeRegexpOptions***REMOVED***
		strictSlash:    r.strictSlash,
		useEncodedPath: r.useEncodedPath,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, q := range r.regexp.queries ***REMOVED***
		if err = uniqueVars(rr.varsN, q.varsN); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if typ == regexpTypeHost ***REMOVED***
		if r.regexp.path != nil ***REMOVED***
			if err = uniqueVars(rr.varsN, r.regexp.path.varsN); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		r.regexp.host = rr
	***REMOVED*** else ***REMOVED***
		if r.regexp.host != nil ***REMOVED***
			if err = uniqueVars(rr.varsN, r.regexp.host.varsN); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if typ == regexpTypeQuery ***REMOVED***
			r.regexp.queries = append(r.regexp.queries, rr)
		***REMOVED*** else ***REMOVED***
			r.regexp.path = rr
		***REMOVED***
	***REMOVED***
	r.addMatcher(rr)
	return nil
***REMOVED***

// Headers --------------------------------------------------------------------

// headerMatcher matches the request against header values.
type headerMatcher map[string]string

func (m headerMatcher) Match(r *http.Request, match *RouteMatch) bool ***REMOVED***
	return matchMapWithString(m, r.Header, true)
***REMOVED***

// Headers adds a matcher for request header values.
// It accepts a sequence of key/value pairs to be matched. For example:
//
//     r := mux.NewRouter()
//     r.Headers("Content-Type", "application/json",
//               "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both request header values match.
// If the value is an empty string, it will match any value if the key is set.
func (r *Route) Headers(pairs ...string) *Route ***REMOVED***
	if r.err == nil ***REMOVED***
		var headers map[string]string
		headers, r.err = mapFromPairsToString(pairs...)
		return r.addMatcher(headerMatcher(headers))
	***REMOVED***
	return r
***REMOVED***

// headerRegexMatcher matches the request against the route given a regex for the header
type headerRegexMatcher map[string]*regexp.Regexp

func (m headerRegexMatcher) Match(r *http.Request, match *RouteMatch) bool ***REMOVED***
	return matchMapWithRegex(m, r.Header, true)
***REMOVED***

// HeadersRegexp accepts a sequence of key/value pairs, where the value has regex
// support. For example:
//
//     r := mux.NewRouter()
//     r.HeadersRegexp("Content-Type", "application/(text|json)",
//               "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both the request header matches both regular expressions.
// If the value is an empty string, it will match any value if the key is set.
// Use the start and end of string anchors (^ and $) to match an exact value.
func (r *Route) HeadersRegexp(pairs ...string) *Route ***REMOVED***
	if r.err == nil ***REMOVED***
		var headers map[string]*regexp.Regexp
		headers, r.err = mapFromPairsToRegex(pairs...)
		return r.addMatcher(headerRegexMatcher(headers))
	***REMOVED***
	return r
***REMOVED***

// Host -----------------------------------------------------------------------

// Host adds a matcher for the URL host.
// It accepts a template with zero or more URL variables enclosed by ***REMOVED******REMOVED***.
// Variables can define an optional regexp pattern to be matched:
//
// - ***REMOVED***name***REMOVED*** matches anything until the next dot.
//
// - ***REMOVED***name:pattern***REMOVED*** matches the given regexp pattern.
//
// For example:
//
//     r := mux.NewRouter()
//     r.Host("www.example.com")
//     r.Host("***REMOVED***subdomain***REMOVED***.domain.com")
//     r.Host("***REMOVED***subdomain:[a-z]+***REMOVED***.domain.com")
//
// Variable names must be unique in a given route. They can be retrieved
// calling mux.Vars(request).
func (r *Route) Host(tpl string) *Route ***REMOVED***
	r.err = r.addRegexpMatcher(tpl, regexpTypeHost)
	return r
***REMOVED***

// MatcherFunc ----------------------------------------------------------------

// MatcherFunc is the function signature used by custom matchers.
type MatcherFunc func(*http.Request, *RouteMatch) bool

// Match returns the match for a given request.
func (m MatcherFunc) Match(r *http.Request, match *RouteMatch) bool ***REMOVED***
	return m(r, match)
***REMOVED***

// MatcherFunc adds a custom function to be used as request matcher.
func (r *Route) MatcherFunc(f MatcherFunc) *Route ***REMOVED***
	return r.addMatcher(f)
***REMOVED***

// Methods --------------------------------------------------------------------

// methodMatcher matches the request against HTTP methods.
type methodMatcher []string

func (m methodMatcher) Match(r *http.Request, match *RouteMatch) bool ***REMOVED***
	return matchInArray(m, r.Method)
***REMOVED***

// Methods adds a matcher for HTTP methods.
// It accepts a sequence of one or more methods to be matched, e.g.:
// "GET", "POST", "PUT".
func (r *Route) Methods(methods ...string) *Route ***REMOVED***
	for k, v := range methods ***REMOVED***
		methods[k] = strings.ToUpper(v)
	***REMOVED***
	return r.addMatcher(methodMatcher(methods))
***REMOVED***

// Path -----------------------------------------------------------------------

// Path adds a matcher for the URL path.
// It accepts a template with zero or more URL variables enclosed by ***REMOVED******REMOVED***. The
// template must start with a "/".
// Variables can define an optional regexp pattern to be matched:
//
// - ***REMOVED***name***REMOVED*** matches anything until the next slash.
//
// - ***REMOVED***name:pattern***REMOVED*** matches the given regexp pattern.
//
// For example:
//
//     r := mux.NewRouter()
//     r.Path("/products/").Handler(ProductsHandler)
//     r.Path("/products/***REMOVED***key***REMOVED***").Handler(ProductsHandler)
//     r.Path("/articles/***REMOVED***category***REMOVED***/***REMOVED***id:[0-9]+***REMOVED***").
//       Handler(ArticleHandler)
//
// Variable names must be unique in a given route. They can be retrieved
// calling mux.Vars(request).
func (r *Route) Path(tpl string) *Route ***REMOVED***
	r.err = r.addRegexpMatcher(tpl, regexpTypePath)
	return r
***REMOVED***

// PathPrefix -----------------------------------------------------------------

// PathPrefix adds a matcher for the URL path prefix. This matches if the given
// template is a prefix of the full URL path. See Route.Path() for details on
// the tpl argument.
//
// Note that it does not treat slashes specially ("/foobar/" will be matched by
// the prefix "/foo") so you may want to use a trailing slash here.
//
// Also note that the setting of Router.StrictSlash() has no effect on routes
// with a PathPrefix matcher.
func (r *Route) PathPrefix(tpl string) *Route ***REMOVED***
	r.err = r.addRegexpMatcher(tpl, regexpTypePrefix)
	return r
***REMOVED***

// Query ----------------------------------------------------------------------

// Queries adds a matcher for URL query values.
// It accepts a sequence of key/value pairs. Values may define variables.
// For example:
//
//     r := mux.NewRouter()
//     r.Queries("foo", "bar", "id", "***REMOVED***id:[0-9]+***REMOVED***")
//
// The above route will only match if the URL contains the defined queries
// values, e.g.: ?foo=bar&id=42.
//
// It the value is an empty string, it will match any value if the key is set.
//
// Variables can define an optional regexp pattern to be matched:
//
// - ***REMOVED***name***REMOVED*** matches anything until the next slash.
//
// - ***REMOVED***name:pattern***REMOVED*** matches the given regexp pattern.
func (r *Route) Queries(pairs ...string) *Route ***REMOVED***
	length := len(pairs)
	if length%2 != 0 ***REMOVED***
		r.err = fmt.Errorf(
			"mux: number of parameters must be multiple of 2, got %v", pairs)
		return nil
	***REMOVED***
	for i := 0; i < length; i += 2 ***REMOVED***
		if r.err = r.addRegexpMatcher(pairs[i]+"="+pairs[i+1], regexpTypeQuery); r.err != nil ***REMOVED***
			return r
		***REMOVED***
	***REMOVED***

	return r
***REMOVED***

// Schemes --------------------------------------------------------------------

// schemeMatcher matches the request against URL schemes.
type schemeMatcher []string

func (m schemeMatcher) Match(r *http.Request, match *RouteMatch) bool ***REMOVED***
	return matchInArray(m, r.URL.Scheme)
***REMOVED***

// Schemes adds a matcher for URL schemes.
// It accepts a sequence of schemes to be matched, e.g.: "http", "https".
func (r *Route) Schemes(schemes ...string) *Route ***REMOVED***
	for k, v := range schemes ***REMOVED***
		schemes[k] = strings.ToLower(v)
	***REMOVED***
	if r.buildScheme == "" && len(schemes) > 0 ***REMOVED***
		r.buildScheme = schemes[0]
	***REMOVED***
	return r.addMatcher(schemeMatcher(schemes))
***REMOVED***

// BuildVarsFunc --------------------------------------------------------------

// BuildVarsFunc is the function signature used by custom build variable
// functions (which can modify route variables before a route's URL is built).
type BuildVarsFunc func(map[string]string) map[string]string

// BuildVarsFunc adds a custom function to be used to modify build variables
// before a route's URL is built.
func (r *Route) BuildVarsFunc(f BuildVarsFunc) *Route ***REMOVED***
	r.buildVarsFunc = f
	return r
***REMOVED***

// Subrouter ------------------------------------------------------------------

// Subrouter creates a subrouter for the route.
//
// It will test the inner routes only if the parent route matched. For example:
//
//     r := mux.NewRouter()
//     s := r.Host("www.example.com").Subrouter()
//     s.HandleFunc("/products/", ProductsHandler)
//     s.HandleFunc("/products/***REMOVED***key***REMOVED***", ProductHandler)
//     s.HandleFunc("/articles/***REMOVED***category***REMOVED***/***REMOVED***id:[0-9]+***REMOVED***"), ArticleHandler)
//
// Here, the routes registered in the subrouter won't be tested if the host
// doesn't match.
func (r *Route) Subrouter() *Router ***REMOVED***
	router := &Router***REMOVED***parent: r, strictSlash: r.strictSlash***REMOVED***
	r.addMatcher(router)
	return router
***REMOVED***

// ----------------------------------------------------------------------------
// URL building
// ----------------------------------------------------------------------------

// URL builds a URL for the route.
//
// It accepts a sequence of key/value pairs for the route variables. For
// example, given this route:
//
//     r := mux.NewRouter()
//     r.HandleFunc("/articles/***REMOVED***category***REMOVED***/***REMOVED***id:[0-9]+***REMOVED***", ArticleHandler).
//       Name("article")
//
// ...a URL for it can be built using:
//
//     url, err := r.Get("article").URL("category", "technology", "id", "42")
//
// ...which will return an url.URL with the following path:
//
//     "/articles/technology/42"
//
// This also works for host variables:
//
//     r := mux.NewRouter()
//     r.Host("***REMOVED***subdomain***REMOVED***.domain.com").
//       HandleFunc("/articles/***REMOVED***category***REMOVED***/***REMOVED***id:[0-9]+***REMOVED***", ArticleHandler).
//       Name("article")
//
//     // url.String() will be "http://news.domain.com/articles/technology/42"
//     url, err := r.Get("article").URL("subdomain", "news",
//                                      "category", "technology",
//                                      "id", "42")
//
// All variables defined in the route are required, and their values must
// conform to the corresponding patterns.
func (r *Route) URL(pairs ...string) (*url.URL, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return nil, r.err
	***REMOVED***
	if r.regexp == nil ***REMOVED***
		return nil, errors.New("mux: route doesn't have a host or path")
	***REMOVED***
	values, err := r.prepareVars(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var scheme, host, path string
	queries := make([]string, 0, len(r.regexp.queries))
	if r.regexp.host != nil ***REMOVED***
		if host, err = r.regexp.host.url(values); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		scheme = "http"
		if s := r.getBuildScheme(); s != "" ***REMOVED***
			scheme = s
		***REMOVED***
	***REMOVED***
	if r.regexp.path != nil ***REMOVED***
		if path, err = r.regexp.path.url(values); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	for _, q := range r.regexp.queries ***REMOVED***
		var query string
		if query, err = q.url(values); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		queries = append(queries, query)
	***REMOVED***
	return &url.URL***REMOVED***
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: strings.Join(queries, "&"),
	***REMOVED***, nil
***REMOVED***

// URLHost builds the host part of the URL for a route. See Route.URL().
//
// The route must have a host defined.
func (r *Route) URLHost(pairs ...string) (*url.URL, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return nil, r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.host == nil ***REMOVED***
		return nil, errors.New("mux: route doesn't have a host")
	***REMOVED***
	values, err := r.prepareVars(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	host, err := r.regexp.host.url(values)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	u := &url.URL***REMOVED***
		Scheme: "http",
		Host:   host,
	***REMOVED***
	if s := r.getBuildScheme(); s != "" ***REMOVED***
		u.Scheme = s
	***REMOVED***
	return u, nil
***REMOVED***

// URLPath builds the path part of the URL for a route. See Route.URL().
//
// The route must have a path defined.
func (r *Route) URLPath(pairs ...string) (*url.URL, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return nil, r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.path == nil ***REMOVED***
		return nil, errors.New("mux: route doesn't have a path")
	***REMOVED***
	values, err := r.prepareVars(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	path, err := r.regexp.path.url(values)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &url.URL***REMOVED***
		Path: path,
	***REMOVED***, nil
***REMOVED***

// GetPathTemplate returns the template used to build the
// route match.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An error will be returned if the route does not define a path.
func (r *Route) GetPathTemplate() (string, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return "", r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.path == nil ***REMOVED***
		return "", errors.New("mux: route doesn't have a path")
	***REMOVED***
	return r.regexp.path.template, nil
***REMOVED***

// GetPathRegexp returns the expanded regular expression used to match route path.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An error will be returned if the route does not define a path.
func (r *Route) GetPathRegexp() (string, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return "", r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.path == nil ***REMOVED***
		return "", errors.New("mux: route does not have a path")
	***REMOVED***
	return r.regexp.path.regexp.String(), nil
***REMOVED***

// GetQueriesRegexp returns the expanded regular expressions used to match the
// route queries.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An empty list will be returned if the route does not have queries.
func (r *Route) GetQueriesRegexp() ([]string, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return nil, r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.queries == nil ***REMOVED***
		return nil, errors.New("mux: route doesn't have queries")
	***REMOVED***
	var queries []string
	for _, query := range r.regexp.queries ***REMOVED***
		queries = append(queries, query.regexp.String())
	***REMOVED***
	return queries, nil
***REMOVED***

// GetQueriesTemplates returns the templates used to build the
// query matching.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An empty list will be returned if the route does not define queries.
func (r *Route) GetQueriesTemplates() ([]string, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return nil, r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.queries == nil ***REMOVED***
		return nil, errors.New("mux: route doesn't have queries")
	***REMOVED***
	var queries []string
	for _, query := range r.regexp.queries ***REMOVED***
		queries = append(queries, query.template)
	***REMOVED***
	return queries, nil
***REMOVED***

// GetMethods returns the methods the route matches against
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An empty list will be returned if route does not have methods.
func (r *Route) GetMethods() ([]string, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return nil, r.err
	***REMOVED***
	for _, m := range r.matchers ***REMOVED***
		if methods, ok := m.(methodMatcher); ok ***REMOVED***
			return []string(methods), nil
		***REMOVED***
	***REMOVED***
	return nil, nil
***REMOVED***

// GetHostTemplate returns the template used to build the
// route match.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An error will be returned if the route does not define a host.
func (r *Route) GetHostTemplate() (string, error) ***REMOVED***
	if r.err != nil ***REMOVED***
		return "", r.err
	***REMOVED***
	if r.regexp == nil || r.regexp.host == nil ***REMOVED***
		return "", errors.New("mux: route doesn't have a host")
	***REMOVED***
	return r.regexp.host.template, nil
***REMOVED***

// prepareVars converts the route variable pairs into a map. If the route has a
// BuildVarsFunc, it is invoked.
func (r *Route) prepareVars(pairs ...string) (map[string]string, error) ***REMOVED***
	m, err := mapFromPairsToString(pairs...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.buildVars(m), nil
***REMOVED***

func (r *Route) buildVars(m map[string]string) map[string]string ***REMOVED***
	if r.parent != nil ***REMOVED***
		m = r.parent.buildVars(m)
	***REMOVED***
	if r.buildVarsFunc != nil ***REMOVED***
		m = r.buildVarsFunc(m)
	***REMOVED***
	return m
***REMOVED***

// ----------------------------------------------------------------------------
// parentRoute
// ----------------------------------------------------------------------------

// parentRoute allows routes to know about parent host and path definitions.
type parentRoute interface ***REMOVED***
	getBuildScheme() string
	getNamedRoutes() map[string]*Route
	getRegexpGroup() *routeRegexpGroup
	buildVars(map[string]string) map[string]string
***REMOVED***

func (r *Route) getBuildScheme() string ***REMOVED***
	if r.buildScheme != "" ***REMOVED***
		return r.buildScheme
	***REMOVED***
	if r.parent != nil ***REMOVED***
		return r.parent.getBuildScheme()
	***REMOVED***
	return ""
***REMOVED***

// getNamedRoutes returns the map where named routes are registered.
func (r *Route) getNamedRoutes() map[string]*Route ***REMOVED***
	if r.parent == nil ***REMOVED***
		// During tests router is not always set.
		r.parent = NewRouter()
	***REMOVED***
	return r.parent.getNamedRoutes()
***REMOVED***

// getRegexpGroup returns regexp definitions from this route.
func (r *Route) getRegexpGroup() *routeRegexpGroup ***REMOVED***
	if r.regexp == nil ***REMOVED***
		if r.parent == nil ***REMOVED***
			// During tests router is not always set.
			r.parent = NewRouter()
		***REMOVED***
		regexp := r.parent.getRegexpGroup()
		if regexp == nil ***REMOVED***
			r.regexp = new(routeRegexpGroup)
		***REMOVED*** else ***REMOVED***
			// Copy.
			r.regexp = &routeRegexpGroup***REMOVED***
				host:    regexp.host,
				path:    regexp.path,
				queries: regexp.queries,
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return r.regexp
***REMOVED***
