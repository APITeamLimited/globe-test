// +build gingonic,!gorillamux,!echo

package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginRouter struct ***REMOVED***
	router *gin.Engine
***REMOVED***

func (g ginRouter) Handler() http.Handler ***REMOVED***
	return g.router
***REMOVED***

func (g ginRouter) Handle(protocol, route string, handler HandlerFunc) ***REMOVED***
	wrappedCallback := func(c *gin.Context) ***REMOVED***
		params := map[string]string***REMOVED******REMOVED***
		for _, p := range c.Params ***REMOVED***
			params[p.Key] = p.Value
		***REMOVED***

		handler(c.Writer, c.Request, params)
	***REMOVED***

	g.router.Handle(protocol, route, wrappedCallback)
***REMOVED***

//Gin creates a new api2go router to use with the gin framework
func Gin(g *gin.Engine) Routeable ***REMOVED***
	return &ginRouter***REMOVED***router: g***REMOVED***
***REMOVED***
