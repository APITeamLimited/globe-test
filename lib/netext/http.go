package netext

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/ThomsonReutersEikon/go-ntlm/ntlm"
	"github.com/pkg/errors"
)

type HTTPTransport struct ***REMOVED***
	*http.Transport

	mu        sync.Mutex
	authCache map[string]bool
***REMOVED***

func NewHTTPTransport(transport *http.Transport) *HTTPTransport ***REMOVED***
	return &HTTPTransport***REMOVED***
		Transport: transport,
		authCache: make(map[string]bool),
	***REMOVED***
***REMOVED***

func (t *HTTPTransport) RoundTrip(req *http.Request) (res *http.Response, err error) ***REMOVED***
	if t.Transport == nil ***REMOVED***
		return nil, errors.New("no roundtrip defined")
	***REMOVED***

	// checking if the request needs ntlm authentication
	if GetAuth(req.Context()) == "ntlm" && req.URL.User != nil ***REMOVED***
		return t.roundtripWithNTLM(req)
	***REMOVED***

	return t.Transport.RoundTrip(req)
***REMOVED***

func (t *HTTPTransport) roundtripWithNTLM(req *http.Request) (res *http.Response, err error) ***REMOVED***
	rt := t.Transport

	username := req.URL.User.Username()
	password, _ := req.URL.User.Password()

	// Save request body
	body := bytes.Buffer***REMOVED******REMOVED***
	if req.Body != nil ***REMOVED***
		_, err = body.ReadFrom(req.Body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if err := req.Body.Close(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
	***REMOVED***

	// before making the request check if there is a cached authorization.
	if _, ok := t.getAuthCache(req.URL.String()); ok ***REMOVED***
		req.Header.Del("Authorization")
	***REMOVED*** else ***REMOVED***
		req.Header.Set("Authorization", "NTLM TlRMTVNTUAABAAAAB4IAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAAAAAAAMAA=")
	***REMOVED***

	res, err = rt.RoundTrip(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if res.StatusCode != http.StatusUnauthorized ***REMOVED***
		return res, err
	***REMOVED***

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := res.Body.Close(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

	// retrieve Www-Authenticate header from response
	ntlmChallenge := res.Header.Get("WWW-Authenticate")
	if !strings.HasPrefix(ntlmChallenge, "NTLM ") ***REMOVED***
		return nil, errors.New("Invalid WWW-Authenticate header")
	***REMOVED***

	challengeBytes, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ntlmChallenge, "NTLM "))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	session, err := ntlm.CreateClientSession(ntlm.Version2, ntlm.ConnectionlessMode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	session.SetUserInfo(username, password, "")

	// parse NTLM challenge
	challenge, err := ntlm.ParseChallengeMessage(challengeBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = session.ProcessChallengeMessage(challenge)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// authenticate user
	authenticate, err := session.GenerateAuthenticateMessage()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// set NTLM Authorization header
	header := "NTLM " + base64.StdEncoding.EncodeToString(authenticate.Bytes())
	req.Header.Set("Authorization", header)

	t.setAuthCache(req.URL.String(), true)

	return rt.RoundTrip(req)
***REMOVED***

func (t *HTTPTransport) setAuthCache(key string, value bool) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()

	t.authCache[key] = value
***REMOVED***

func (t *HTTPTransport) getAuthCache(key string) (bool, bool) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()

	value, ok := t.authCache[key]
	return value, ok
***REMOVED***
