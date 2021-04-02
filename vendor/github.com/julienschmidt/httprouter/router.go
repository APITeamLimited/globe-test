// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package httprouter is a trie based high performance HTTP request router.
//
// A trivial example is:
//
//  package main
//
//  import (
//      "fmt"
//      "github.com/julienschmidt/httprouter"
//      "net/http"
//      "log"
//  )
//
//  func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) ***REMOVED***
//      fmt.Fprint(w, "Welcome!\n")
//  ***REMOVED***
//
//  func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) ***REMOVED***
//      fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//  ***REMOVED***
//
//  func main() ***REMOVED***
//      router := httprouter.New()
//      router.GET("/", Index)
//      router.GET("/hello/:name", Hello)
//
//      log.Fatal(http.ListenAndServe(":8080", router))
//  ***REMOVED***
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH and DELETE shortcut functions exist to
// register handles, for all other methods router.Handle can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//  Path: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//  Path: /files/*filepath
//
//  Requests:
//   /files/                             match: filepath="/"
//   /files/LICENSE                      match: filepath="/LICENSE"
//   /files/templates/article.html       match: filepath="/templates/article.html"
//   /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//  // by the name of the parameter
//  user := ps.ByName("user") // defined by :user or *user
//
//  // by the index of the parameter. This way you can also get the name (key)
//  thirdKey   := ps[2].Key   // the name of the 3rd parameter
//  thirdValue := ps[2].Value // the value of the 3rd parameter
package httprouter

import (
	"net/http"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
type Handle func(http.ResponseWriter, *http.Request, Params)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct ***REMOVED***
	Key   string
	Value string
***REMOVED***

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string ***REMOVED***
	for i := range ps ***REMOVED***
		if ps[i].Key == name ***REMOVED***
			return ps[i].Value
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct ***REMOVED***
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface***REMOVED******REMOVED***)
***REMOVED***

// Make sure the Router conforms with the http.Handler interface
var _ http.Handler = New()

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router ***REMOVED***
	return &Router***REMOVED***
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	***REMOVED***
***REMOVED***

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle Handle) ***REMOVED***
	r.Handle("GET", path, handle)
***REMOVED***

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle Handle) ***REMOVED***
	r.Handle("HEAD", path, handle)
***REMOVED***

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle Handle) ***REMOVED***
	r.Handle("OPTIONS", path, handle)
***REMOVED***

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle Handle) ***REMOVED***
	r.Handle("POST", path, handle)
***REMOVED***

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle Handle) ***REMOVED***
	r.Handle("PUT", path, handle)
***REMOVED***

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle Handle) ***REMOVED***
	r.Handle("PATCH", path, handle)
***REMOVED***

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle Handle) ***REMOVED***
	r.Handle("DELETE", path, handle)
***REMOVED***

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handle Handle) ***REMOVED***
	if path[0] != '/' ***REMOVED***
		panic("path must begin with '/' in path '" + path + "'")
	***REMOVED***

	if r.trees == nil ***REMOVED***
		r.trees = make(map[string]*node)
	***REMOVED***

	root := r.trees[method]
	if root == nil ***REMOVED***
		root = new(node)
		r.trees[method] = root
	***REMOVED***

	root.addRoute(path, handle)
***REMOVED***

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) ***REMOVED***
	r.Handler(method, path, handler)
***REMOVED***

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path string, root http.FileSystem) ***REMOVED***
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" ***REMOVED***
		panic("path must end with /*filepath in path '" + path + "'")
	***REMOVED***

	fileServer := http.FileServer(root)

	r.GET(path, func(w http.ResponseWriter, req *http.Request, ps Params) ***REMOVED***
		req.URL.Path = ps.ByName("filepath")
		fileServer.ServeHTTP(w, req)
	***REMOVED***)
***REMOVED***

func (r *Router) recv(w http.ResponseWriter, req *http.Request) ***REMOVED***
	if rcv := recover(); rcv != nil ***REMOVED***
		r.PanicHandler(w, req, rcv)
	***REMOVED***
***REMOVED***

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (Handle, Params, bool) ***REMOVED***
	if root := r.trees[method]; root != nil ***REMOVED***
		return root.getValue(path)
	***REMOVED***
	return nil, nil, false
***REMOVED***

func (r *Router) allowed(path, reqMethod string) (allow string) ***REMOVED***
	if path == "*" ***REMOVED*** // server-wide
		for method := range r.trees ***REMOVED***
			if method == "OPTIONS" ***REMOVED***
				continue
			***REMOVED***

			// add request method to list of allowed methods
			if len(allow) == 0 ***REMOVED***
				allow = method
			***REMOVED*** else ***REMOVED***
				allow += ", " + method
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // specific path
		for method := range r.trees ***REMOVED***
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" ***REMOVED***
				continue
			***REMOVED***

			handle, _, _ := r.trees[method].getValue(path)
			if handle != nil ***REMOVED***
				// add request method to list of allowed methods
				if len(allow) == 0 ***REMOVED***
					allow = method
				***REMOVED*** else ***REMOVED***
					allow += ", " + method
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(allow) > 0 ***REMOVED***
		allow += ", OPTIONS"
	***REMOVED***
	return
***REMOVED***

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) ***REMOVED***
	if r.PanicHandler != nil ***REMOVED***
		defer r.recv(w, req)
	***REMOVED***

	path := req.URL.Path

	if root := r.trees[req.Method]; root != nil ***REMOVED***
		if handle, ps, tsr := root.getValue(path); handle != nil ***REMOVED***
			handle(w, req, ps)
			return
		***REMOVED*** else if req.Method != "CONNECT" && path != "/" ***REMOVED***
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" ***REMOVED***
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			***REMOVED***

			if tsr && r.RedirectTrailingSlash ***REMOVED***
				if len(path) > 1 && path[len(path)-1] == '/' ***REMOVED***
					req.URL.Path = path[:len(path)-1]
				***REMOVED*** else ***REMOVED***
					req.URL.Path = path + "/"
				***REMOVED***
				http.Redirect(w, req, req.URL.String(), code)
				return
			***REMOVED***

			// Try to fix the request path
			if r.RedirectFixedPath ***REMOVED***
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found ***REMOVED***
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if req.Method == "OPTIONS" && r.HandleOPTIONS ***REMOVED***
		// Handle OPTIONS requests
		if allow := r.allowed(path, req.Method); len(allow) > 0 ***REMOVED***
			w.Header().Set("Allow", allow)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Handle 405
		if r.HandleMethodNotAllowed ***REMOVED***
			if allow := r.allowed(path, req.Method); len(allow) > 0 ***REMOVED***
				w.Header().Set("Allow", allow)
				if r.MethodNotAllowed != nil ***REMOVED***
					r.MethodNotAllowed.ServeHTTP(w, req)
				***REMOVED*** else ***REMOVED***
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Handle 404
	if r.NotFound != nil ***REMOVED***
		r.NotFound.ServeHTTP(w, req)
	***REMOVED*** else ***REMOVED***
		http.NotFound(w, req)
	***REMOVED***
***REMOVED***
