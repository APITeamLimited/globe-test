package goHttpDigestClient

import (
	"io"
	"io/ioutil"
	"net/http"
)

// if option is set, get challenge at construct time
// if option not set, ever digest auth will send 2 request
type Client struct ***REMOVED***
	is_init bool
	option  ClientOption
	http.Client
***REMOVED***

type ClientOption struct ***REMOVED***
	username string
	password string
***REMOVED***

// create new Client instance
func NewClient(username, password string) *Client ***REMOVED***
	opt := &ClientOption***REMOVED***username: username, password: password***REMOVED***
	// here need more attention
	return &Client***REMOVED***option: *opt, is_init: false***REMOVED***
***REMOVED***

func GetChallengeFromHeader(h *http.Header) Challenge ***REMOVED***
	return NewChallenge(h.Get(KEY_WWW_Authenticate))
***REMOVED***

func (c *Client) Do(req *http.Request, opt *ClientOption) (*http.Response, error) ***REMOVED***
	res, err := c.Client.Do(req)
	if res.StatusCode == http.StatusUnauthorized ***REMOVED***
		challenge := GetChallengeFromHeader(&res.Header)
		challenge.ComputeResponse(req.Method, req.URL.RequestURI(), getStrFromIO(req.Body), opt.username, opt.password)
		authorization := challenge.ToAuthorizationStr()
		req.Header.Set(KEY_AUTHORIZATION, authorization)
		return c.Client.Do(req)
	***REMOVED*** else ***REMOVED***
		return res, err
	***REMOVED***
***REMOVED***

// From ReadCloser to string
func getStrFromIO(r io.ReadCloser) string ***REMOVED***
	if r == nil ***REMOVED***
		return ""
	***REMOVED***
	if b, err := ioutil.ReadAll(r); err == nil ***REMOVED***
		return string(b)
	***REMOVED*** else ***REMOVED***
		return ""
	***REMOVED***
***REMOVED***

// static Defualt Client
var DefaultClient = &Client***REMOVED***is_init: true***REMOVED***

// Default Client Do　Request
func Do(req *http.Request, opt *ClientOption) (*http.Response, error) ***REMOVED***
	return DefaultClient.Do(req, opt)
***REMOVED***
