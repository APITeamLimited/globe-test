// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin/render"
)

// Version is Framework's version
const Version = "v1.1.4"

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")
var defaultAppEngine bool

type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() HandlerFunc ***REMOVED***
	length := len(c)
	if length > 0 ***REMOVED***
		return c[length-1]
	***REMOVED***
	return nil
***REMOVED***

type (
	RoutesInfo []RouteInfo
	RouteInfo  struct ***REMOVED***
		Method  string
		Path    string
		Handler string
	***REMOVED***

	// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
	// Create an instance of Engine, by using New() or Default()
	Engine struct ***REMOVED***
		RouterGroup
		HTMLRender  render.HTMLRender
		allNoRoute  HandlersChain
		allNoMethod HandlersChain
		noRoute     HandlersChain
		noMethod    HandlersChain
		pool        sync.Pool
		trees       methodTrees

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
		ForwardedByClientIP    bool

		// #726 #755 If enabled, it will thrust some headers starting with
		// 'X-AppEngine...' for better integration with that PaaS.
		AppEngine bool

		// If enabled, the url.RawPath will be used to find parameters.
		UseRawPath bool
		// If true, the path value will be unescaped.
		// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
		// as url.Path gonna be used, which is already unescaped.
		UnescapePathValues bool
	***REMOVED***
)

var _ IRouter = &Engine***REMOVED******REMOVED***

// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true
func New() *Engine ***REMOVED***
	debugPrintWARNINGNew()
	engine := &Engine***REMOVED***
		RouterGroup: RouterGroup***REMOVED***
			Handlers: nil,
			basePath: "/",
			root:     true,
		***REMOVED***,
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		AppEngine:              defaultAppEngine,
		UseRawPath:             false,
		UnescapePathValues:     true,
		trees:                  make(methodTrees, 0, 9),
	***REMOVED***
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface***REMOVED******REMOVED*** ***REMOVED***
		return engine.allocateContext()
	***REMOVED***
	return engine
***REMOVED***

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine ***REMOVED***
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
***REMOVED***

func (engine *Engine) allocateContext() *Context ***REMOVED***
	return &Context***REMOVED***engine: engine***REMOVED***
***REMOVED***

func (engine *Engine) LoadHTMLGlob(pattern string) ***REMOVED***
	if IsDebugging() ***REMOVED***
		debugPrintLoadTemplate(template.Must(template.ParseGlob(pattern)))
		engine.HTMLRender = render.HTMLDebug***REMOVED***Glob: pattern***REMOVED***
	***REMOVED*** else ***REMOVED***
		templ := template.Must(template.ParseGlob(pattern))
		engine.SetHTMLTemplate(templ)
	***REMOVED***
***REMOVED***

func (engine *Engine) LoadHTMLFiles(files ...string) ***REMOVED***
	if IsDebugging() ***REMOVED***
		engine.HTMLRender = render.HTMLDebug***REMOVED***Files: files***REMOVED***
	***REMOVED*** else ***REMOVED***
		templ := template.Must(template.ParseFiles(files...))
		engine.SetHTMLTemplate(templ)
	***REMOVED***
***REMOVED***

func (engine *Engine) SetHTMLTemplate(templ *template.Template) ***REMOVED***
	if len(engine.trees) > 0 ***REMOVED***
		debugPrintWARNINGSetHTMLTemplate()
	***REMOVED***
	engine.HTMLRender = render.HTMLProduction***REMOVED***Template: templ***REMOVED***
***REMOVED***

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) ***REMOVED***
	engine.noRoute = handlers
	engine.rebuild404Handlers()
***REMOVED***

// NoMethod sets the handlers called when... TODO
func (engine *Engine) NoMethod(handlers ...HandlerFunc) ***REMOVED***
	engine.noMethod = handlers
	engine.rebuild405Handlers()
***REMOVED***

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes ***REMOVED***
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
***REMOVED***

func (engine *Engine) rebuild404Handlers() ***REMOVED***
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
***REMOVED***

func (engine *Engine) rebuild405Handlers() ***REMOVED***
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
***REMOVED***

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) ***REMOVED***
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(len(method) > 0, "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	debugPrintRoute(method, path, handlers)
	root := engine.trees.get(method)
	if root == nil ***REMOVED***
		root = new(node)
		engine.trees = append(engine.trees, methodTree***REMOVED***method: method, root: root***REMOVED***)
	***REMOVED***
	root.addRoute(path, handlers)
***REMOVED***

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) ***REMOVED***
	for _, tree := range engine.trees ***REMOVED***
		routes = iterate("", tree.method, routes, tree.root)
	***REMOVED***
	return routes
***REMOVED***

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo ***REMOVED***
	path += root.path
	if len(root.handlers) > 0 ***REMOVED***
		routes = append(routes, RouteInfo***REMOVED***
			Method:  method,
			Path:    path,
			Handler: nameOfFunction(root.handlers.Last()),
		***REMOVED***)
	***REMOVED***
	for _, child := range root.children ***REMOVED***
		routes = iterate(path, method, routes, child)
	***REMOVED***
	return routes
***REMOVED***

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) ***REMOVED***
	defer func() ***REMOVED*** debugPrintError(err) ***REMOVED***()

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine)
	return
***REMOVED***

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) (err error) ***REMOVED***
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	defer func() ***REMOVED*** debugPrintError(err) ***REMOVED***()

	err = http.ListenAndServeTLS(addr, certFile, keyFile, engine)
	return
***REMOVED***

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) ***REMOVED***
	debugPrint("Listening and serving HTTP on unix:/%s", file)
	defer func() ***REMOVED*** debugPrintError(err) ***REMOVED***()

	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer listener.Close()
	err = http.Serve(listener, engine)
	return
***REMOVED***

// Conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) ***REMOVED***
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
***REMOVED***

// Re-enter a context that has been rewritten.
// This can be done by setting c.Request.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (engine *Engine) HandleContext(c *Context) ***REMOVED***
	c.reset()
	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
***REMOVED***

func (engine *Engine) handleHTTPRequest(context *Context) ***REMOVED***
	httpMethod := context.Request.Method
	var path string
	var unescape bool
	if engine.UseRawPath && len(context.Request.URL.RawPath) > 0 ***REMOVED***
		path = context.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	***REMOVED*** else ***REMOVED***
		path = context.Request.URL.Path
		unescape = false
	***REMOVED***

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ ***REMOVED***
		if t[i].method == httpMethod ***REMOVED***
			root := t[i].root
			// Find route in tree
			handlers, params, tsr := root.getValue(path, context.Params, unescape)
			if handlers != nil ***REMOVED***
				context.handlers = handlers
				context.Params = params
				context.Next()
				context.writermem.WriteHeaderNow()
				return

			***REMOVED*** else if httpMethod != "CONNECT" && path != "/" ***REMOVED***
				if tsr && engine.RedirectTrailingSlash ***REMOVED***
					redirectTrailingSlash(context)
					return
				***REMOVED***
				if engine.RedirectFixedPath && redirectFixedPath(context, root, engine.RedirectFixedPath) ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	// TODO: unit test
	if engine.HandleMethodNotAllowed ***REMOVED***
		for _, tree := range engine.trees ***REMOVED***
			if tree.method != httpMethod ***REMOVED***
				if handlers, _, _ := tree.root.getValue(path, nil, unescape); handlers != nil ***REMOVED***
					context.handlers = engine.allNoMethod
					serveError(context, 405, default405Body)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	context.handlers = engine.allNoRoute
	serveError(context, 404, default404Body)
***REMOVED***

var mimePlain = []string***REMOVED***MIMEPlain***REMOVED***

func serveError(c *Context, code int, defaultMessage []byte) ***REMOVED***
	c.writermem.status = code
	c.Next()
	if !c.writermem.Written() ***REMOVED***
		if c.writermem.Status() == code ***REMOVED***
			c.writermem.Header()["Content-Type"] = mimePlain
			c.Writer.Write(defaultMessage)
		***REMOVED*** else ***REMOVED***
			c.writermem.WriteHeaderNow()
		***REMOVED***
	***REMOVED***
***REMOVED***

func redirectTrailingSlash(c *Context) ***REMOVED***
	req := c.Request
	path := req.URL.Path
	code := 301 // Permanent redirect, request with GET method
	if req.Method != "GET" ***REMOVED***
		code = 307
	***REMOVED***

	if len(path) > 1 && path[len(path)-1] == '/' ***REMOVED***
		req.URL.Path = path[:len(path)-1]
	***REMOVED*** else ***REMOVED***
		req.URL.Path = path + "/"
	***REMOVED***
	debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
	http.Redirect(c.Writer, req, req.URL.String(), code)
	c.writermem.WriteHeaderNow()
***REMOVED***

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool ***REMOVED***
	req := c.Request
	path := req.URL.Path

	fixedPath, found := root.findCaseInsensitivePath(
		cleanPath(path),
		trailingSlash,
	)
	if found ***REMOVED***
		code := 301 // Permanent redirect, request with GET method
		if req.Method != "GET" ***REMOVED***
			code = 307
		***REMOVED***
		req.URL.Path = string(fixedPath)
		debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
		http.Redirect(c.Writer, req, req.URL.String(), code)
		c.writermem.WriteHeaderNow()
		return true
	***REMOVED***
	return false
***REMOVED***
