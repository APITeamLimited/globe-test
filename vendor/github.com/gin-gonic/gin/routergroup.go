// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

type (
	IRouter interface ***REMOVED***
		IRoutes
		Group(string, ...HandlerFunc) *RouterGroup
	***REMOVED***

	IRoutes interface ***REMOVED***
		Use(...HandlerFunc) IRoutes

		Handle(string, string, ...HandlerFunc) IRoutes
		Any(string, ...HandlerFunc) IRoutes
		GET(string, ...HandlerFunc) IRoutes
		POST(string, ...HandlerFunc) IRoutes
		DELETE(string, ...HandlerFunc) IRoutes
		PATCH(string, ...HandlerFunc) IRoutes
		PUT(string, ...HandlerFunc) IRoutes
		OPTIONS(string, ...HandlerFunc) IRoutes
		HEAD(string, ...HandlerFunc) IRoutes

		StaticFile(string, string) IRoutes
		Static(string, string) IRoutes
		StaticFS(string, http.FileSystem) IRoutes
	***REMOVED***

	// RouterGroup is used internally to configure router, a RouterGroup is associated with a prefix
	// and an array of handlers (middleware)
	RouterGroup struct ***REMOVED***
		Handlers HandlersChain
		basePath string
		engine   *Engine
		root     bool
	***REMOVED***
)

var _ IRouter = &RouterGroup***REMOVED******REMOVED***

// Use adds middleware to the group, see example code in github.
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes ***REMOVED***
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
***REMOVED***

// Group creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup ***REMOVED***
	return &RouterGroup***REMOVED***
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	***REMOVED***
***REMOVED***

func (group *RouterGroup) BasePath() string ***REMOVED***
	return group.basePath
***REMOVED***

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes ***REMOVED***
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
***REMOVED***

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in github.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil ***REMOVED***
		panic("http method " + httpMethod + " is not valid")
	***REMOVED***
	return group.handle(httpMethod, relativePath, handlers)
***REMOVED***

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("POST", relativePath, handlers)
***REMOVED***

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("GET", relativePath, handlers)
***REMOVED***

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("DELETE", relativePath, handlers)
***REMOVED***

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("PATCH", relativePath, handlers)
***REMOVED***

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("PUT", relativePath, handlers)
***REMOVED***

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("OPTIONS", relativePath, handlers)
***REMOVED***

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	return group.handle("HEAD", relativePath, handlers)
***REMOVED***

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes ***REMOVED***
	group.handle("GET", relativePath, handlers)
	group.handle("POST", relativePath, handlers)
	group.handle("PUT", relativePath, handlers)
	group.handle("PATCH", relativePath, handlers)
	group.handle("HEAD", relativePath, handlers)
	group.handle("OPTIONS", relativePath, handlers)
	group.handle("DELETE", relativePath, handlers)
	group.handle("CONNECT", relativePath, handlers)
	group.handle("TRACE", relativePath, handlers)
	return group.returnObj()
***REMOVED***

// StaticFile registers a single route in order to server a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (group *RouterGroup) StaticFile(relativePath, filepath string) IRoutes ***REMOVED***
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") ***REMOVED***
		panic("URL parameters can not be used when serving a static file")
	***REMOVED***
	handler := func(c *Context) ***REMOVED***
		c.File(filepath)
	***REMOVED***
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group.returnObj()
***REMOVED***

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (group *RouterGroup) Static(relativePath, root string) IRoutes ***REMOVED***
	return group.StaticFS(relativePath, Dir(root, false))
***REMOVED***

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes ***REMOVED***
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") ***REMOVED***
		panic("URL parameters can not be used when serving a static folder")
	***REMOVED***
	handler := group.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group.returnObj()
***REMOVED***

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc ***REMOVED***
	absolutePath := group.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	_, nolisting := fs.(*onlyfilesFS)
	return func(c *Context) ***REMOVED***
		if nolisting ***REMOVED***
			c.Writer.WriteHeader(404)
		***REMOVED***
		fileServer.ServeHTTP(c.Writer, c.Request)
	***REMOVED***
***REMOVED***

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain ***REMOVED***
	finalSize := len(group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) ***REMOVED***
		panic("too many handlers")
	***REMOVED***
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
***REMOVED***

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string ***REMOVED***
	return joinPaths(group.basePath, relativePath)
***REMOVED***

func (group *RouterGroup) returnObj() IRoutes ***REMOVED***
	if group.root ***REMOVED***
		return group.engine
	***REMOVED***
	return group
***REMOVED***
