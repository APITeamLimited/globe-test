package ntlmssp

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetDomain : parse domain name from based on slashes in the input
func GetDomain(user string) (string, string) ***REMOVED***
	domain := ""

	if strings.Contains(user, "\\") ***REMOVED***
		ucomponents := strings.SplitN(user, "\\", 2)
		domain = ucomponents[0]
		user = ucomponents[1]
	***REMOVED***
	return user, domain
***REMOVED***

//Negotiator is a http.Roundtripper decorator that automatically
//converts basic authentication to NTLM/Negotiate authentication when appropriate.
type Negotiator struct***REMOVED*** http.RoundTripper ***REMOVED***

//RoundTrip sends the request to the server, handling any authentication
//re-sends as needed.
func (l Negotiator) RoundTrip(req *http.Request) (res *http.Response, err error) ***REMOVED***
	// Use default round tripper if not provided
	rt := l.RoundTripper
	if rt == nil ***REMOVED***
		rt = http.DefaultTransport
	***REMOVED***
	// If it is not basic auth, just round trip the request as usual
	reqauth := authheader(req.Header.Values("Authorization"))
	if !reqauth.IsBasic() ***REMOVED***
		return rt.RoundTrip(req)
	***REMOVED***
	reqauthBasic := reqauth.Basic()
	// Save request body
	body := bytes.Buffer***REMOVED******REMOVED***
	if req.Body != nil ***REMOVED***
		_, err = body.ReadFrom(req.Body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
	***REMOVED***
	// first try anonymous, in case the server still finds us
	// authenticated from previous traffic
	req.Header.Del("Authorization")
	res, err = rt.RoundTrip(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if res.StatusCode != http.StatusUnauthorized ***REMOVED***
		return res, err
	***REMOVED***
	resauth := authheader(res.Header.Values("Www-Authenticate"))
	if !resauth.IsNegotiate() && !resauth.IsNTLM() ***REMOVED***
		// Unauthorized, Negotiate not requested, let's try with basic auth
		req.Header.Set("Authorization", string(reqauthBasic))
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

		res, err = rt.RoundTrip(req)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if res.StatusCode != http.StatusUnauthorized ***REMOVED***
			return res, err
		***REMOVED***
		resauth = authheader(res.Header.Values("Www-Authenticate"))
	***REMOVED***

	if resauth.IsNegotiate() || resauth.IsNTLM() ***REMOVED***
		// 401 with request:Basic and response:Negotiate
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()

		// recycle credentials
		u, p, err := reqauth.GetBasicCreds()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// get domain from username
		domain := ""
		u, domain = GetDomain(u)

		// send negotiate
		negotiateMessage, err := NewNegotiateMessage(domain, "")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if resauth.IsNTLM() ***REMOVED***
			req.Header.Set("Authorization", "NTLM "+base64.StdEncoding.EncodeToString(negotiateMessage))
		***REMOVED*** else ***REMOVED***
			req.Header.Set("Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(negotiateMessage))
		***REMOVED***

		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

		res, err = rt.RoundTrip(req)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// receive challenge?
		resauth = authheader(res.Header.Values("Www-Authenticate"))
		challengeMessage, err := resauth.GetData()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !(resauth.IsNegotiate() || resauth.IsNTLM()) || len(challengeMessage) == 0 ***REMOVED***
			// Negotiation failed, let client deal with response
			return res, nil
		***REMOVED***
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()

		// send authenticate
		authenticateMessage, err := ProcessChallenge(challengeMessage, u, p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if resauth.IsNTLM() ***REMOVED***
			req.Header.Set("Authorization", "NTLM "+base64.StdEncoding.EncodeToString(authenticateMessage))
		***REMOVED*** else ***REMOVED***
			req.Header.Set("Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(authenticateMessage))
		***REMOVED***

		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

		return rt.RoundTrip(req)
	***REMOVED***

	return res, err
***REMOVED***
