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
