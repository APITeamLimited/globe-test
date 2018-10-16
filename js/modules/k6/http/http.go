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

package http

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"fmt"
	"net/http/httputil"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/netext"
	log "github.com/sirupsen/logrus"
)

const (
	HTTP_METHOD_GET     = "GET"
	HTTP_METHOD_POST    = "POST"
	HTTP_METHOD_PUT     = "PUT"
	HTTP_METHOD_DELETE  = "DELETE"
	HTTP_METHOD_HEAD    = "HEAD"
	HTTP_METHOD_PATCH   = "PATCH"
	HTTP_METHOD_OPTIONS = "OPTIONS"
)

// ErrJarForbiddenInInitContext is used when a cookie jar was made in the init context
var ErrJarForbiddenInInitContext = common.NewInitContextError("Making cookie jars in the init context is not supported")

type HTTPCookie struct ***REMOVED***
	Name, Value, Domain, Path string
	HttpOnly, Secure          bool
	MaxAge                    int
	Expires                   int64
***REMOVED***

type HTTPRequestCookie struct ***REMOVED***
	Name, Value string
	Replace     bool
***REMOVED***

type HTTP struct ***REMOVED***
	SSL_3_0                            string `js:"SSL_3_0"`
	TLS_1_0                            string `js:"TLS_1_0"`
	TLS_1_1                            string `js:"TLS_1_1"`
	TLS_1_2                            string `js:"TLS_1_2"`
	OCSP_STATUS_GOOD                   string `js:"OCSP_STATUS_GOOD"`
	OCSP_STATUS_REVOKED                string `js:"OCSP_STATUS_REVOKED"`
	OCSP_STATUS_SERVER_FAILED          string `js:"OCSP_STATUS_SERVER_FAILED"`
	OCSP_STATUS_UNKNOWN                string `js:"OCSP_STATUS_UNKNOWN"`
	OCSP_REASON_UNSPECIFIED            string `js:"OCSP_REASON_UNSPECIFIED"`
	OCSP_REASON_KEY_COMPROMISE         string `js:"OCSP_REASON_KEY_COMPROMISE"`
	OCSP_REASON_CA_COMPROMISE          string `js:"OCSP_REASON_CA_COMPROMISE"`
	OCSP_REASON_AFFILIATION_CHANGED    string `js:"OCSP_REASON_AFFILIATION_CHANGED"`
	OCSP_REASON_SUPERSEDED             string `js:"OCSP_REASON_SUPERSEDED"`
	OCSP_REASON_CESSATION_OF_OPERATION string `js:"OCSP_REASON_CESSATION_OF_OPERATION"`
	OCSP_REASON_CERTIFICATE_HOLD       string `js:"OCSP_REASON_CERTIFICATE_HOLD"`
	OCSP_REASON_REMOVE_FROM_CRL        string `js:"OCSP_REASON_REMOVE_FROM_CRL"`
	OCSP_REASON_PRIVILEGE_WITHDRAWN    string `js:"OCSP_REASON_PRIVILEGE_WITHDRAWN"`
	OCSP_REASON_AA_COMPROMISE          string `js:"OCSP_REASON_AA_COMPROMISE"`
***REMOVED***

func New() *HTTP ***REMOVED***
	//TODO: move this as an anonymous struct somewhere...
	return &HTTP***REMOVED***
		SSL_3_0:                            netext.SSL_3_0,
		TLS_1_0:                            netext.TLS_1_0,
		TLS_1_1:                            netext.TLS_1_1,
		TLS_1_2:                            netext.TLS_1_2,
		OCSP_STATUS_GOOD:                   netext.OCSP_STATUS_GOOD,
		OCSP_STATUS_REVOKED:                netext.OCSP_STATUS_REVOKED,
		OCSP_STATUS_SERVER_FAILED:          netext.OCSP_STATUS_SERVER_FAILED,
		OCSP_STATUS_UNKNOWN:                netext.OCSP_STATUS_UNKNOWN,
		OCSP_REASON_UNSPECIFIED:            netext.OCSP_REASON_UNSPECIFIED,
		OCSP_REASON_KEY_COMPROMISE:         netext.OCSP_REASON_KEY_COMPROMISE,
		OCSP_REASON_CA_COMPROMISE:          netext.OCSP_REASON_CA_COMPROMISE,
		OCSP_REASON_AFFILIATION_CHANGED:    netext.OCSP_REASON_AFFILIATION_CHANGED,
		OCSP_REASON_SUPERSEDED:             netext.OCSP_REASON_SUPERSEDED,
		OCSP_REASON_CESSATION_OF_OPERATION: netext.OCSP_REASON_CESSATION_OF_OPERATION,
		OCSP_REASON_CERTIFICATE_HOLD:       netext.OCSP_REASON_CERTIFICATE_HOLD,
		OCSP_REASON_REMOVE_FROM_CRL:        netext.OCSP_REASON_REMOVE_FROM_CRL,
		OCSP_REASON_PRIVILEGE_WITHDRAWN:    netext.OCSP_REASON_PRIVILEGE_WITHDRAWN,
		OCSP_REASON_AA_COMPROMISE:          netext.OCSP_REASON_AA_COMPROMISE,
	***REMOVED***
***REMOVED***

func (*HTTP) XCookieJar(ctx *context.Context) *HTTPCookieJar ***REMOVED***
	return newCookieJar(ctx)
***REMOVED***

func (*HTTP) CookieJar(ctx context.Context) (*HTTPCookieJar, error) ***REMOVED***
	state := common.GetState(ctx)
	if state == nil ***REMOVED***
		return nil, ErrJarForbiddenInInitContext
	***REMOVED***
	return &HTTPCookieJar***REMOVED***state.CookieJar, &ctx***REMOVED***, nil
***REMOVED***

func (*HTTP) mergeCookies(req *http.Request, jar *cookiejar.Jar, reqCookies map[string]*HTTPRequestCookie) map[string][]*HTTPRequestCookie ***REMOVED***
	allCookies := make(map[string][]*HTTPRequestCookie)
	for _, c := range jar.Cookies(req.URL) ***REMOVED***
		allCookies[c.Name] = append(allCookies[c.Name], &HTTPRequestCookie***REMOVED***Name: c.Name, Value: c.Value***REMOVED***)
	***REMOVED***
	for key, reqCookie := range reqCookies ***REMOVED***
		if jc := allCookies[key]; jc != nil && reqCookie.Replace ***REMOVED***
			allCookies[key] = []*HTTPRequestCookie***REMOVED******REMOVED***Name: key, Value: reqCookie.Value***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			allCookies[key] = append(allCookies[key], &HTTPRequestCookie***REMOVED***Name: key, Value: reqCookie.Value***REMOVED***)
		***REMOVED***
	***REMOVED***
	return allCookies
***REMOVED***

func (*HTTP) setRequestCookies(req *http.Request, reqCookies map[string][]*HTTPRequestCookie) ***REMOVED***
	for _, cookies := range reqCookies ***REMOVED***
		for _, c := range cookies ***REMOVED***
			req.AddCookie(&http.Cookie***REMOVED***Name: c.Name, Value: c.Value***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (*HTTP) debugRequest(state *common.State, req *http.Request, description string) ***REMOVED***
	if state.Options.HttpDebug.String != "" ***REMOVED***
		dump, err := httputil.DumpRequestOut(req, state.Options.HttpDebug.String == "full")
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		logDump(description, dump)
	***REMOVED***
***REMOVED***

func (*HTTP) debugResponse(state *common.State, res *http.Response, description string) ***REMOVED***
	if state.Options.HttpDebug.String != "" && res != nil ***REMOVED***
		dump, err := httputil.DumpResponse(res, state.Options.HttpDebug.String == "full")
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		logDump(description, dump)
	***REMOVED***
***REMOVED***

func logDump(description string, dump []byte) ***REMOVED***
	fmt.Printf("%s:\n%s\n", description, dump)
***REMOVED***
