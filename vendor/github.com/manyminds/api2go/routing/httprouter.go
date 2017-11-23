package routing

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HTTPRouter default router implementation for api2go
type HTTPRouter struct ***REMOVED***
	router *httprouter.Router
***REMOVED***

// Handle each method like before and wrap them into julienschmidt handler func style
func (h HTTPRouter) Handle(protocol, route string, handler HandlerFunc) ***REMOVED***
	wrappedCallback := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) ***REMOVED***
		params := map[string]string***REMOVED******REMOVED***
		for _, p := range ps ***REMOVED***
			params[p.Key] = p.Value
		***REMOVED***

		handler(w, r, params)
	***REMOVED***

	h.router.Handle(protocol, route, wrappedCallback)
***REMOVED***

// Handler returns the router
func (h HTTPRouter) Handler() http.Handler ***REMOVED***
	return h.router
***REMOVED***

// SetRedirectTrailingSlash wraps this internal functionality of
// the julienschmidt router.
func (h *HTTPRouter) SetRedirectTrailingSlash(enabled bool) ***REMOVED***
	h.router.RedirectTrailingSlash = enabled
***REMOVED***

// GetRouteParameter implemention will extract the param the julienschmidt way
func (h HTTPRouter) GetRouteParameter(r http.Request, param string) string ***REMOVED***
	path := httprouter.CleanPath(r.URL.Path)
	_, params, _ := h.router.Lookup(r.Method, path)
	return params.ByName(param)
***REMOVED***

// NewHTTPRouter returns a new instance of julienschmidt/httprouter
// this is the default router when using api2go
func NewHTTPRouter(prefix string, notAllowedHandler http.Handler) Routeable ***REMOVED***
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.MethodNotAllowed = notAllowedHandler
	return &HTTPRouter***REMOVED***router: router***REMOVED***
***REMOVED***
