/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package api

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/api/common"
	"github.com/loadimpact/k6/api/v1"
	"github.com/loadimpact/k6/api/v2"
	"github.com/loadimpact/k6/lib"
	"github.com/urfave/negroni"
	"net/http"
)

const (
	staticRoot = "../web/dist"

	notFoundText = "UI unavailable. If you're using a custom build of k6, please remember to run `make`."
)

func NewHandler(root string) http.Handler ***REMOVED***
	mux := http.NewServeMux()
	mux.Handle("/v1/", v1.NewHandler())
	mux.Handle("/v2/", v2.NewHandler())
	mux.Handle("/ping", HandlePing())
	mux.Handle("/", HandleStatic(root))
	return mux
***REMOVED***

func ListenAndServe(addr string, engine *lib.Engine) error ***REMOVED***
	mux := NewHandler(staticRoot)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseFunc(WithEngine(engine))
	n.UseFunc(NewLogger(log.StandardLogger()))
	n.UseHandler(mux)

	return http.ListenAndServe(addr, n)
***REMOVED***

func NewLogger(l *log.Logger) negroni.HandlerFunc ***REMOVED***
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
		next(rw, r)

		res := rw.(negroni.ResponseWriter)
		l.WithField("status", res.Status()).Debugf("%s %s", r.Method, r.URL.Path)
	***REMOVED***
***REMOVED***

func WithEngine(engine *lib.Engine) negroni.HandlerFunc ***REMOVED***
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
		r = r.WithContext(common.WithEngine(r.Context(), engine))
		next(rw, r)
	***REMOVED***)
***REMOVED***

func HandlePing() http.Handler ***REMOVED***
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
		rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(rw, "ok")
	***REMOVED***)
***REMOVED***

func HandleStatic(root string) http.Handler ***REMOVED***
	box, err := rice.FindBox(root)
	if err != nil ***REMOVED***
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) ***REMOVED***
			rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rw, notFoundText)
		***REMOVED***)
	***REMOVED***
	return http.FileServer(box.HTTPBox())
***REMOVED***
