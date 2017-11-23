package negroni

import (
	"log"
	"net/http"
	"os"
)

const (
	// DefaultAddress is used if no other is specified.
	DefaultAddress = ":8080"
)

// Handler handler is an interface that objects can implement to be registered to serve as middleware
// in the Negroni middleware stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next http.HandlerFunc
// passed in.
//
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
type Handler interface ***REMOVED***
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
***REMOVED***

// HandlerFunc is an adapter to allow the use of ordinary functions as Negroni handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
	h(rw, r, next)
***REMOVED***

type middleware struct ***REMOVED***
	handler Handler
	next    *middleware
***REMOVED***

func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
***REMOVED***

// Wrap converts a http.Handler into a negroni.Handler so it can be used as a Negroni
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func Wrap(handler http.Handler) Handler ***REMOVED***
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
		handler.ServeHTTP(rw, r)
		next(rw, r)
	***REMOVED***)
***REMOVED***

// WrapFunc converts a http.HandlerFunc into a negroni.Handler so it can be used as a Negroni
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func WrapFunc(handlerFunc http.HandlerFunc) Handler ***REMOVED***
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
		handlerFunc(rw, r)
		next(rw, r)
	***REMOVED***)
***REMOVED***

// Negroni is a stack of Middleware Handlers that can be invoked as an http.Handler.
// Negroni middleware is evaluated in the order that they are added to the stack using
// the Use and UseHandler methods.
type Negroni struct ***REMOVED***
	middleware middleware
	handlers   []Handler
***REMOVED***

// New returns a new Negroni instance with no middleware preconfigured.
func New(handlers ...Handler) *Negroni ***REMOVED***
	return &Negroni***REMOVED***
		handlers:   handlers,
		middleware: build(handlers),
	***REMOVED***
***REMOVED***

// With returns a new Negroni instance that is a combination of the negroni
// receiver's handlers and the provided handlers.
func (n *Negroni) With(handlers ...Handler) *Negroni ***REMOVED***
	return New(
		append(n.handlers, handlers...)...,
	)
***REMOVED***

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Classic() *Negroni ***REMOVED***
	return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))
***REMOVED***

func (n *Negroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	n.middleware.ServeHTTP(NewResponseWriter(rw), r)
***REMOVED***

// Use adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) Use(handler Handler) ***REMOVED***
	if handler == nil ***REMOVED***
		panic("handler cannot be nil")
	***REMOVED***

	n.handlers = append(n.handlers, handler)
	n.middleware = build(n.handlers)
***REMOVED***

// UseFunc adds a Negroni-style handler function onto the middleware stack.
func (n *Negroni) UseFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) ***REMOVED***
	n.Use(HandlerFunc(handlerFunc))
***REMOVED***

// UseHandler adds a http.Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) UseHandler(handler http.Handler) ***REMOVED***
	n.Use(Wrap(handler))
***REMOVED***

// UseHandlerFunc adds a http.HandlerFunc-style handler function onto the middleware stack.
func (n *Negroni) UseHandlerFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request)) ***REMOVED***
	n.UseHandler(http.HandlerFunc(handlerFunc))
***REMOVED***

// Run is a convenience function that runs the negroni stack as an HTTP
// server. The addr string, if provided, takes the same format as http.ListenAndServe.
// If no address is provided but the PORT environment variable is set, the PORT value is used.
// If neither is provided, the address' value will equal the DefaultAddress constant.
func (n *Negroni) Run(addr ...string) ***REMOVED***
	l := log.New(os.Stdout, "[negroni] ", 0)
	finalAddr := detectAddress(addr...)
	l.Printf("listening on %s", finalAddr)
	l.Fatal(http.ListenAndServe(finalAddr, n))
***REMOVED***

func detectAddress(addr ...string) string ***REMOVED***
	if len(addr) > 0 ***REMOVED***
		return addr[0]
	***REMOVED***
	if port := os.Getenv("PORT"); port != "" ***REMOVED***
		return ":" + port
	***REMOVED***
	return DefaultAddress
***REMOVED***

// Returns a list of all the handlers in the current Negroni middleware chain.
func (n *Negroni) Handlers() []Handler ***REMOVED***
	return n.handlers
***REMOVED***

func build(handlers []Handler) middleware ***REMOVED***
	var next middleware

	if len(handlers) == 0 ***REMOVED***
		return voidMiddleware()
	***REMOVED*** else if len(handlers) > 1 ***REMOVED***
		next = build(handlers[1:])
	***REMOVED*** else ***REMOVED***
		next = voidMiddleware()
	***REMOVED***

	return middleware***REMOVED***handlers[0], &next***REMOVED***
***REMOVED***

func voidMiddleware() middleware ***REMOVED***
	return middleware***REMOVED***
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED******REMOVED***),
		&middleware***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***
