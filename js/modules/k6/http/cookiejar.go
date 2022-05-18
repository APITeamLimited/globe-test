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
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strings"
	"time"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
)

// ErrJarForbiddenInInitContext is used when a cookie jar was made in the init context
// TODO: unexport this? there's no reason for this to be exported
var ErrJarForbiddenInInitContext = common.NewInitContextError("Making cookie jars in the init context is not supported")

// CookieJar is cookiejar.Jar wrapper to be used in js scripts
type CookieJar struct ***REMOVED***
	moduleInstance *ModuleInstance
	// js is to make it not be accessible from inside goja/js, the json is
	// for when it is returned from setup().
	Jar *cookiejar.Jar `js:"-" json:"-"`
***REMOVED***

// CookiesForURL return the cookies for a given url as a map of key and values
func (j CookieJar) CookiesForURL(url string) map[string][]string ***REMOVED***
	u, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	cookies := j.Jar.Cookies(u)
	objs := make(map[string][]string, len(cookies))
	for _, c := range cookies ***REMOVED***
		objs[c.Name] = append(objs[c.Name], c.Value)
	***REMOVED***
	return objs
***REMOVED***

// Set sets a cookie for a particular url with the given name value and additional opts
func (j CookieJar) Set(url, name, value string, opts goja.Value) (bool, error) ***REMOVED***
	rt := j.moduleInstance.vu.Runtime()

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
						return false, fmt.Errorf(`unable to parse "expires" date string "%s": %w`, expires, err)
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
	j.Jar.SetCookies(u, []*http.Cookie***REMOVED***&c***REMOVED***)
	return true, nil
***REMOVED***

// Clear all cookies for a particular URL
func (j CookieJar) Clear(url string) error ***REMOVED***
	u, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cookies := j.Jar.Cookies(u)
	for _, c := range cookies ***REMOVED***
		c.MaxAge = -1
	***REMOVED***
	j.Jar.SetCookies(u, cookies)

	return nil
***REMOVED***

// Delete cookies for a particular URL
func (j CookieJar) Delete(url, name string) error ***REMOVED***
	if name == "" ***REMOVED***
		return errors.New("cookie: is null")
	***REMOVED***

	u, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c := http.Cookie***REMOVED***Name: name, MaxAge: -1***REMOVED***
	j.Jar.SetCookies(u, []*http.Cookie***REMOVED***&c***REMOVED***)

	return nil
***REMOVED***
