// +build gorillamux,!gingonic,!echo

package routing

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type gorillamuxRouter struct ***REMOVED***
	router *mux.Router
***REMOVED***

func (gm gorillamuxRouter) Handler() http.Handler ***REMOVED***
	return gm.router
***REMOVED***

func (gm gorillamuxRouter) Handle(protocol, route string, handler HandlerFunc) ***REMOVED***
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		handler(w, r, mux.Vars(r))
	***REMOVED***

	// The request path will have parameterized segments indicated as :name.  Convert
	// that notation to the ***REMOVED***name***REMOVED*** notation used by Gorilla mux.
	orig := strings.Split(route, "/")
	var mod []string
	for _, s := range orig ***REMOVED***
		if len(s) > 0 && s[0] == ':' ***REMOVED***
			s = fmt.Sprintf("***REMOVED***%s***REMOVED***", s[1:])
		***REMOVED***
		mod = append(mod, s)
	***REMOVED***
	modroute := strings.Join(mod, "/")

	gm.router.HandleFunc(modroute, wrappedHandler).Methods(protocol)
***REMOVED***

//Gorilla creates a new api2go router to use with the Gorilla mux framework
func Gorilla(gm *mux.Router) Routeable ***REMOVED***
	return &gorillamuxRouter***REMOVED***router: gm***REMOVED***
***REMOVED***
