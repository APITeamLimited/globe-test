/*
Package echo implements high performance, minimalist Go web framework.

Example:

  package main

  import (
    "net/http"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
  )

  // Handler
  func hello(c echo.Context) error ***REMOVED***
    return c.String(http.StatusOK, "Hello, World!")
  ***REMOVED***

  func main() ***REMOVED***
    // Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    e.GET("/", hello)

    // Start server
    e.Logger.Fatal(e.Start(":1323"))
  ***REMOVED***

Learn more at https://echo.labstack.com
*/
package echo

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	stdLog "log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
)

type (
	// Echo is the top-level framework instance.
	Echo struct ***REMOVED***
		stdLogger        *stdLog.Logger
		colorer          *color.Color
		premiddleware    []MiddlewareFunc
		middleware       []MiddlewareFunc
		maxParam         *int
		router           *Router
		notFoundHandler  HandlerFunc
		pool             sync.Pool
		Server           *http.Server
		TLSServer        *http.Server
		Listener         net.Listener
		TLSListener      net.Listener
		AutoTLSManager   autocert.Manager
		DisableHTTP2     bool
		Debug            bool
		HideBanner       bool
		HTTPErrorHandler HTTPErrorHandler
		Binder           Binder
		Validator        Validator
		Renderer         Renderer
		// Mutex            sync.RWMutex
		Logger Logger
	***REMOVED***

	// Route contains a handler and information for matching against requests.
	Route struct ***REMOVED***
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
	***REMOVED***

	// HTTPError represents an error that occurred while handling a request.
	HTTPError struct ***REMOVED***
		Code    int
		Message interface***REMOVED******REMOVED***
		Inner   error // Stores the error returned by an external dependency
	***REMOVED***

	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// HandlerFunc defines a function to server HTTP requests.
	HandlerFunc func(Context) error

	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error, Context)

	// Validator is the interface that wraps the Validate function.
	Validator interface ***REMOVED***
		Validate(i interface***REMOVED******REMOVED***) error
	***REMOVED***

	// Renderer is the interface that wraps the Render function.
	Renderer interface ***REMOVED***
		Render(io.Writer, string, interface***REMOVED******REMOVED***, Context) error
	***REMOVED***

	// Map defines a generic map of type `map[string]interface***REMOVED******REMOVED***`.
	Map map[string]interface***REMOVED******REMOVED***

	// i is the interface for Echo and Group.
	i interface ***REMOVED***
		GET(string, HandlerFunc, ...MiddlewareFunc) *Route
	***REMOVED***
)

// HTTP methods
const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=UTF-8"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

const (
	version = "3.2.6"
	website = "https://echo.labstack.com"
	// http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Echo
	banner = `
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ %s
High performance, minimalist Go web framework
%s
____________________________________O/_______
                                    O\
`
)

var (
	methods = [...]string***REMOVED***
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	***REMOVED***
)

// Errors
var (
	ErrUnsupportedMediaType        = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   = NewHTTPError(http.StatusForbidden)
	ErrMethodNotAllowed            = NewHTTPError(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrValidatorNotRegistered      = errors.New("Validator not registered")
	ErrRendererNotRegistered       = errors.New("Renderer not registered")
	ErrInvalidRedirectCode         = errors.New("Invalid redirect status code")
	ErrCookieNotFound              = errors.New("Cookie not found")
)

// Error handlers
var (
	NotFoundHandler = func(c Context) error ***REMOVED***
		return ErrNotFound
	***REMOVED***

	MethodNotAllowedHandler = func(c Context) error ***REMOVED***
		return ErrMethodNotAllowed
	***REMOVED***
)

// New creates an instance of Echo.
func New() (e *Echo) ***REMOVED***
	e = &Echo***REMOVED***
		Server:    new(http.Server),
		TLSServer: new(http.Server),
		AutoTLSManager: autocert.Manager***REMOVED***
			Prompt: autocert.AcceptTOS,
		***REMOVED***,
		Logger:   log.New("echo"),
		colorer:  color.New(),
		maxParam: new(int),
	***REMOVED***
	e.Server.Handler = e
	e.TLSServer.Handler = e
	e.HTTPErrorHandler = e.DefaultHTTPErrorHandler
	e.Binder = &DefaultBinder***REMOVED******REMOVED***
	e.Logger.SetLevel(log.ERROR)
	e.stdLogger = stdLog.New(e.Logger.Output(), e.Logger.Prefix()+": ", 0)
	e.pool.New = func() interface***REMOVED******REMOVED*** ***REMOVED***
		return e.NewContext(nil, nil)
	***REMOVED***
	e.router = NewRouter(e)
	return
***REMOVED***

// NewContext returns a Context instance.
func (e *Echo) NewContext(r *http.Request, w http.ResponseWriter) Context ***REMOVED***
	return &context***REMOVED***
		request:  r,
		response: NewResponse(w, e),
		store:    make(Map),
		echo:     e,
		pvalues:  make([]string, *e.maxParam),
		handler:  NotFoundHandler,
	***REMOVED***
***REMOVED***

// Router returns router.
func (e *Echo) Router() *Router ***REMOVED***
	return e.router
***REMOVED***

// DefaultHTTPErrorHandler is the default HTTP error handler. It sends a JSON response
// with status code.
func (e *Echo) DefaultHTTPErrorHandler(err error, c Context) ***REMOVED***
	var (
		code = http.StatusInternalServerError
		msg  interface***REMOVED******REMOVED***
	)

	if he, ok := err.(*HTTPError); ok ***REMOVED***
		code = he.Code
		msg = he.Message
		if he.Inner != nil ***REMOVED***
			msg = fmt.Sprintf("%v, %v", err, he.Inner)
		***REMOVED***
	***REMOVED*** else if e.Debug ***REMOVED***
		msg = err.Error()
	***REMOVED*** else ***REMOVED***
		msg = http.StatusText(code)
	***REMOVED***
	if _, ok := msg.(string); ok ***REMOVED***
		msg = Map***REMOVED***"message": msg***REMOVED***
	***REMOVED***

	e.Logger.Error(err)

	// Send response
	if !c.Response().Committed ***REMOVED***
		if c.Request().Method == HEAD ***REMOVED*** // Issue #608
			err = c.NoContent(code)
		***REMOVED*** else ***REMOVED***
			err = c.JSON(code, msg)
		***REMOVED***
		if err != nil ***REMOVED***
			e.Logger.Error(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Pre adds middleware to the chain which is run before router.
func (e *Echo) Pre(middleware ...MiddlewareFunc) ***REMOVED***
	e.premiddleware = append(e.premiddleware, middleware...)
***REMOVED***

// Use adds middleware to the chain which is run after router.
func (e *Echo) Use(middleware ...MiddlewareFunc) ***REMOVED***
	e.middleware = append(e.middleware, middleware...)
***REMOVED***

// CONNECT registers a new CONNECT route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(CONNECT, path, h, m...)
***REMOVED***

// DELETE registers a new DELETE route for a path with matching handler in the router
// with optional route-level middleware.
func (e *Echo) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(DELETE, path, h, m...)
***REMOVED***

// GET registers a new GET route for a path with matching handler in the router
// with optional route-level middleware.
func (e *Echo) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(GET, path, h, m...)
***REMOVED***

// HEAD registers a new HEAD route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(HEAD, path, h, m...)
***REMOVED***

// OPTIONS registers a new OPTIONS route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(OPTIONS, path, h, m...)
***REMOVED***

// PATCH registers a new PATCH route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(PATCH, path, h, m...)
***REMOVED***

// POST registers a new POST route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(POST, path, h, m...)
***REMOVED***

// PUT registers a new PUT route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(PUT, path, h, m...)
***REMOVED***

// TRACE registers a new TRACE route for a path with matching handler in the
// router with optional route-level middleware.
func (e *Echo) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route ***REMOVED***
	return e.Add(TRACE, path, h, m...)
***REMOVED***

// Any registers a new route for all HTTP methods and path with matching handler
// in the router with optional route-level middleware.
func (e *Echo) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route ***REMOVED***
	routes := make([]*Route, len(methods))
	for i, m := range methods ***REMOVED***
		routes[i] = e.Add(m, path, handler, middleware...)
	***REMOVED***
	return routes
***REMOVED***

// Match registers a new route for multiple HTTP methods and path with matching
// handler in the router with optional route-level middleware.
func (e *Echo) Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route ***REMOVED***
	routes := make([]*Route, len(methods))
	for i, m := range methods ***REMOVED***
		routes[i] = e.Add(m, path, handler, middleware...)
	***REMOVED***
	return routes
***REMOVED***

// Static registers a new route with path prefix to serve static files from the
// provided root directory.
func (e *Echo) Static(prefix, root string) *Route ***REMOVED***
	if root == "" ***REMOVED***
		root = "." // For security we want to restrict to CWD.
	***REMOVED***
	return static(e, prefix, root)
***REMOVED***

func static(i i, prefix, root string) *Route ***REMOVED***
	h := func(c Context) error ***REMOVED***
		p, err := PathUnescape(c.Param("*"))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		name := filepath.Join(root, path.Clean("/"+p)) // "/"+ for security
		return c.File(name)
	***REMOVED***
	i.GET(prefix, h)
	if prefix == "/" ***REMOVED***
		return i.GET(prefix+"*", h)
	***REMOVED***

	return i.GET(prefix+"/*", h)
***REMOVED***

// File registers a new route with path to serve a static file.
func (e *Echo) File(path, file string) *Route ***REMOVED***
	return e.GET(path, func(c Context) error ***REMOVED***
		return c.File(file)
	***REMOVED***)
***REMOVED***

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (e *Echo) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route ***REMOVED***
	name := handlerName(handler)
	e.router.Add(method, path, func(c Context) error ***REMOVED***
		h := handler
		// Chain middleware
		for i := len(middleware) - 1; i >= 0; i-- ***REMOVED***
			h = middleware[i](h)
		***REMOVED***
		return h(c)
	***REMOVED***)
	r := &Route***REMOVED***
		Method: method,
		Path:   path,
		Name:   name,
	***REMOVED***
	e.router.routes[method+path] = r
	return r
***REMOVED***

// Group creates a new router group with prefix and optional group-level middleware.
func (e *Echo) Group(prefix string, m ...MiddlewareFunc) (g *Group) ***REMOVED***
	g = &Group***REMOVED***prefix: prefix, echo: e***REMOVED***
	g.Use(m...)
	return
***REMOVED***

// URI generates a URI from handler.
func (e *Echo) URI(handler HandlerFunc, params ...interface***REMOVED******REMOVED***) string ***REMOVED***
	name := handlerName(handler)
	return e.Reverse(name, params...)
***REMOVED***

// URL is an alias for `URI` function.
func (e *Echo) URL(h HandlerFunc, params ...interface***REMOVED******REMOVED***) string ***REMOVED***
	return e.URI(h, params...)
***REMOVED***

// Reverse generates an URL from route name and provided parameters.
func (e *Echo) Reverse(name string, params ...interface***REMOVED******REMOVED***) string ***REMOVED***
	uri := new(bytes.Buffer)
	ln := len(params)
	n := 0
	for _, r := range e.router.routes ***REMOVED***
		if r.Name == name ***REMOVED***
			for i, l := 0, len(r.Path); i < l; i++ ***REMOVED***
				if r.Path[i] == ':' && n < ln ***REMOVED***
					for ; i < l && r.Path[i] != '/'; i++ ***REMOVED***
					***REMOVED***
					uri.WriteString(fmt.Sprintf("%v", params[n]))
					n++
				***REMOVED***
				if i < l ***REMOVED***
					uri.WriteByte(r.Path[i])
				***REMOVED***
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return uri.String()
***REMOVED***

// Routes returns the registered routes.
func (e *Echo) Routes() []*Route ***REMOVED***
	routes := []*Route***REMOVED******REMOVED***
	for _, v := range e.router.routes ***REMOVED***
		routes = append(routes, v)
	***REMOVED***
	return routes
***REMOVED***

// AcquireContext returns an empty `Context` instance from the pool.
// You must return the context by calling `ReleaseContext()`.
func (e *Echo) AcquireContext() Context ***REMOVED***
	return e.pool.Get().(Context)
***REMOVED***

// ReleaseContext returns the `Context` instance back to the pool.
// You must call it after `AcquireContext()`.
func (e *Echo) ReleaseContext(c Context) ***REMOVED***
	e.pool.Put(c)
***REMOVED***

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	// Acquire lock
	// e.Mutex.RLock()
	// defer e.Mutex.RUnlock()

	// Acquire context
	c := e.pool.Get().(*context)
	defer e.pool.Put(c)
	c.Reset(r, w)

	// Middleware
	h := func(c Context) error ***REMOVED***
		method := r.Method
		path := r.URL.RawPath
		if path == "" ***REMOVED***
			path = r.URL.Path
		***REMOVED***
		e.router.Find(method, path, c)
		h := c.Handler()
		for i := len(e.middleware) - 1; i >= 0; i-- ***REMOVED***
			h = e.middleware[i](h)
		***REMOVED***
		return h(c)
	***REMOVED***

	// Premiddleware
	for i := len(e.premiddleware) - 1; i >= 0; i-- ***REMOVED***
		h = e.premiddleware[i](h)
	***REMOVED***

	// Execute chain
	if err := h(c); err != nil ***REMOVED***
		e.HTTPErrorHandler(err, c)
	***REMOVED***
***REMOVED***

// Start starts an HTTP server.
func (e *Echo) Start(address string) error ***REMOVED***
	e.Server.Addr = address
	return e.StartServer(e.Server)
***REMOVED***

// StartTLS starts an HTTPS server.
func (e *Echo) StartTLS(address string, certFile, keyFile string) (err error) ***REMOVED***
	if certFile == "" || keyFile == "" ***REMOVED***
		return errors.New("invalid tls configuration")
	***REMOVED***
	s := e.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)
	s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return e.startTLS(address)
***REMOVED***

// StartAutoTLS starts an HTTPS server using certificates automatically installed from https://letsencrypt.org.
func (e *Echo) StartAutoTLS(address string) error ***REMOVED***
	s := e.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.GetCertificate = e.AutoTLSManager.GetCertificate
	return e.startTLS(address)
***REMOVED***

func (e *Echo) startTLS(address string) error ***REMOVED***
	s := e.TLSServer
	s.Addr = address
	if !e.DisableHTTP2 ***REMOVED***
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	***REMOVED***
	return e.StartServer(e.TLSServer)
***REMOVED***

// StartServer starts a custom http server.
func (e *Echo) StartServer(s *http.Server) (err error) ***REMOVED***
	// Setup
	e.colorer.SetOutput(e.Logger.Output())
	s.ErrorLog = e.stdLogger
	s.Handler = e
	if e.Debug ***REMOVED***
		e.Logger.SetLevel(log.DEBUG)
	***REMOVED***

	if !e.HideBanner ***REMOVED***
		e.colorer.Printf(banner, e.colorer.Red("v"+version), e.colorer.Blue(website))
	***REMOVED***

	if s.TLSConfig == nil ***REMOVED***
		if e.Listener == nil ***REMOVED***
			e.Listener, err = newListener(s.Addr)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if !e.HideBanner ***REMOVED***
			e.colorer.Printf("⇨ http server started on %s\n", e.colorer.Green(e.Listener.Addr()))
		***REMOVED***
		return s.Serve(e.Listener)
	***REMOVED***
	if e.TLSListener == nil ***REMOVED***
		l, err := newListener(s.Addr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		e.TLSListener = tls.NewListener(l, s.TLSConfig)
	***REMOVED***
	if !e.HideBanner ***REMOVED***
		e.colorer.Printf("⇨ https server started on %s\n", e.colorer.Green(e.TLSListener.Addr()))
	***REMOVED***
	return s.Serve(e.TLSListener)
***REMOVED***

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message ...interface***REMOVED******REMOVED***) *HTTPError ***REMOVED***
	he := &HTTPError***REMOVED***Code: code, Message: http.StatusText(code)***REMOVED***
	if len(message) > 0 ***REMOVED***
		he.Message = message[0]
	***REMOVED***
	return he
***REMOVED***

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string ***REMOVED***
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
***REMOVED***

// WrapHandler wraps `http.Handler` into `echo.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc ***REMOVED***
	return func(c Context) error ***REMOVED***
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	***REMOVED***
***REMOVED***

// WrapMiddleware wraps `func(http.Handler) http.Handler` into `echo.MiddlewareFunc`
func WrapMiddleware(m func(http.Handler) http.Handler) MiddlewareFunc ***REMOVED***
	return func(next HandlerFunc) HandlerFunc ***REMOVED***
		return func(c Context) (err error) ***REMOVED***
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
				c.SetRequest(r)
				err = next(c)
			***REMOVED***)).ServeHTTP(c.Response(), c.Request())
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func handlerName(h HandlerFunc) string ***REMOVED***
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func ***REMOVED***
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	***REMOVED***
	return t.String()
***REMOVED***

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct ***REMOVED***
	*net.TCPListener
***REMOVED***

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) ***REMOVED***
	tc, err := ln.AcceptTCP()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
***REMOVED***

func newListener(address string) (*tcpKeepAliveListener, error) ***REMOVED***
	l, err := net.Listen("tcp", address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &tcpKeepAliveListener***REMOVED***l.(*net.TCPListener)***REMOVED***, nil
***REMOVED***
