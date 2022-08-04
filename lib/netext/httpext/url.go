package httpext

import (
	"fmt"
	"net/url"
)

// A URL wraps net.URL, and preserves the template (if any) the URL was constructed from.
type URL struct ***REMOVED***
	u        *url.URL
	Name     string // http://example.com/thing/$***REMOVED******REMOVED***/
	URL      string // http://example.com/thing/1234/
	CleanURL string // URL with masked user credentials, used for output
***REMOVED***

// NewURL returns a new URL for the provided url and name. The error is returned if the url provided
// can't be parsed
func NewURL(urlString, name string) (URL, error) ***REMOVED***
	u, err := url.Parse(urlString)
	if err != nil ***REMOVED***
		return URL***REMOVED******REMOVED***, NewK6Error(invalidURLErrorCode,
			fmt.Sprintf("%s: %s", invalidURLErrorCodeMsg, err), err)
	***REMOVED***
	newURL := URL***REMOVED***u: u, Name: name, URL: urlString***REMOVED***
	newURL.CleanURL = newURL.Clean()
	if urlString == name ***REMOVED***
		newURL.Name = newURL.CleanURL
	***REMOVED***
	return newURL, nil
***REMOVED***

// Clean returns an output-safe representation of URL
func (u URL) Clean() string ***REMOVED***
	if u.CleanURL != "" ***REMOVED***
		return u.CleanURL
	***REMOVED***

	if u.u == nil || u.u.User == nil ***REMOVED***
		return u.URL
	***REMOVED***

	if password, passwordOk := u.u.User.Password(); passwordOk ***REMOVED***
		// here 3 is for the '://' and 4 is because of '://' and ':' between the credentials
		return u.URL[:len(u.u.Scheme)+3] + "****:****" + u.URL[len(u.u.Scheme)+4+len(u.u.User.Username())+len(password):]
	***REMOVED***

	// here 3 in both places is for the '://'
	return u.URL[:len(u.u.Scheme)+3] + "****" + u.URL[len(u.u.Scheme)+3+len(u.u.User.Username()):]
***REMOVED***

// GetURL returns the internal url.URL
func (u URL) GetURL() *url.URL ***REMOVED***
	return u.u
***REMOVED***

// ToURL tries to convert anything passed to it to a k6 URL struct
func ToURL(u interface***REMOVED******REMOVED***) (URL, error) ***REMOVED***
	switch tu := u.(type) ***REMOVED***
	case URL:
		// Handling of http.url`http://example.com/***REMOVED***$id***REMOVED***`
		return tu, nil
	case string:
		// Handling of "http://example.com/"
		return NewURL(tu, tu)
	default:
		return URL***REMOVED******REMOVED***, NewK6Error(invalidURLErrorCode,
			fmt.Sprintf("%s: '#%v'", invalidURLErrorCodeMsg, u), nil)
	***REMOVED***
***REMOVED***
