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
	neturl "net/url"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/pkg/errors"
)

// HTTPCookieJar is cookiejar.Jar wrapper to be used in js scripts
type HTTPCookieJar struct ***REMOVED***
	jar *cookiejar.Jar
	ctx *context.Context
***REMOVED***

func newCookieJar(ctxPtr *context.Context) *HTTPCookieJar ***REMOVED***
	jar, err := cookiejar.New(nil)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(*ctxPtr), err)
	***REMOVED***
	return &HTTPCookieJar***REMOVED***jar, ctxPtr***REMOVED***
***REMOVED***

// CookiesForURL return the cookies for a given url as a map of key and values
func (j HTTPCookieJar) CookiesForURL(url string) map[string][]string ***REMOVED***
	u, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	cookies := j.jar.Cookies(u)
	objs := make(map[string][]string, len(cookies))
	for _, c := range cookies ***REMOVED***
		objs[c.Name] = append(objs[c.Name], c.Value)
	***REMOVED***
	return objs
***REMOVED***

// Set sets a cookie for a particular url with the given name value and additional opts
func (j HTTPCookieJar) Set(url, name, value string, opts goja.Value) (bool, error) ***REMOVED***
	rt := common.GetRuntime(*j.ctx)

	u, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	c := http.Cookie***REMOVED***Name: name, Value: value***REMOVED***
	paramsV := opts
	if paramsV != nil && !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
		params := paramsV.ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch strings.ToLower(k) ***REMOVED***
			case "path":
				c.Path = params.Get(k).String()
			case "domain":
				c.Domain = params.Get(k).String()
			case "expires":
				var t time.Time
				expires := params.Get(k).String()
				if expires != "" ***REMOVED***
					t, err = time.Parse(time.RFC1123, expires)
					if err != nil ***REMOVED***
						return false, errors.Errorf("unable to parse \"expires\" date string \"%s\" with: %s", expires, err.Error())
					***REMOVED***
				***REMOVED***
				c.Expires = t
			case "max_age":
				c.MaxAge = int(params.Get(k).ToInteger())
			case "secure":
				c.Secure = params.Get(k).ToBoolean()
			case "http_only":
				c.HttpOnly = params.Get(k).ToBoolean()
			***REMOVED***
		***REMOVED***
	***REMOVED***
	j.jar.SetCookies(u, []*http.Cookie***REMOVED***&c***REMOVED***)
	return true, nil
***REMOVED***
