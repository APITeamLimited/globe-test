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

package lib

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// CookieJar implements a simplified version of net/http/cookiejar, that most notably can be
// cleared without reinstancing the whole thing.
type CookieJar struct ***REMOVED***
	cookies map[string][]*http.Cookie
***REMOVED***

func NewCookieJar() *CookieJar ***REMOVED***
	jar := &CookieJar***REMOVED******REMOVED***
	jar.Clear()
	return jar
***REMOVED***

func (j *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) ***REMOVED***
	if u.Scheme != "http" && u.Scheme != "https" ***REMOVED***
		return
	***REMOVED***
	j.cookies[cookieHostKey(u.Host)] = cookies
***REMOVED***

func (j *CookieJar) Cookies(u *url.URL) []*http.Cookie ***REMOVED***
	return j.cookies[cookieHostKey(u.Host)]
***REMOVED***

func (j *CookieJar) Clear() ***REMOVED***
	j.cookies = make(map[string][]*http.Cookie)
***REMOVED***

func cookieHostKey(host string) string ***REMOVED***
	if strings.HasSuffix(host, ".") ***REMOVED***
		host = host[:len(host)-1]
	***REMOVED***
	host, _, err := net.SplitHostPort(host)
	if err != nil ***REMOVED***
		return host
	***REMOVED***
	return host
***REMOVED***
