// +build echo,!gorillamux,!gingonic

package routing

import (
	"net/http"

	"github.com/labstack/echo"
)

type echoRouter struct ***REMOVED***
	echo *echo.Echo
***REMOVED***

func (e echoRouter) Handler() http.Handler ***REMOVED***
	return e.echo
***REMOVED***

func (e echoRouter) Handle(protocol, route string, handler HandlerFunc) ***REMOVED***
	echoHandlerFunc := func(c echo.Context) error ***REMOVED***
		params := map[string]string***REMOVED******REMOVED***

		for i, p := range c.ParamNames() ***REMOVED***
			params[p] = c.ParamValues()[i]
		***REMOVED***

		handler(c.Response(), c.Request(), params)

		return nil
	***REMOVED***
	e.echo.Add(protocol, route, echoHandlerFunc)
***REMOVED***

// Echo created a new api2go router to use with the echo framework
func Echo(e *echo.Echo) Routeable ***REMOVED***
	return &echoRouter***REMOVED***echo: e***REMOVED***
***REMOVED***
