/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ThomsonReutersEikon/go-ntlm/ntlm"
	"github.com/pkg/errors"
)

// NTLMNegotiator is a http.Roundtripper that automatically converts basic authentication to
// NTLM authentication.
type NTLMNegotiator struct***REMOVED*** http.RoundTripper ***REMOVED***

func (n NTLMNegotiator) RoundTrip(req *http.Request) (res *http.Response, err error) ***REMOVED***
	// Use default round tripper if not provided
	rt := n.RoundTripper
	if rt == nil ***REMOVED***
		rt = http.DefaultTransport
	***REMOVED***

	username, password, err := getCredentialsFromHeader(req.Header.Get("Authorization"))
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "get basic credentials from header failed")
	***REMOVED***

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

	// Sending first request to get challenge data
	req.Header.Set("Authorization", "NTLM TlRMTVNTUAABAAAABoIIAAAAAAAAAAAAAAAAAAAAAAA=")
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
	if ntlmChallenge == "" ***REMOVED***
		return nil, errors.New("Empty WWW-Authenticate header")
	***REMOVED***

	challengeBytes, err := base64.StdEncoding.DecodeString(strings.Replace(ntlmChallenge, "NTLM ", "", -1))
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
	req.Header.Set("Authorization", "NTLM "+base64.StdEncoding.EncodeToString(authenticate.Bytes()))

	return rt.RoundTrip(req)
***REMOVED***

func getCredentialsFromHeader(header string) (string, string, error) ***REMOVED***
	credBytes, err := base64.StdEncoding.DecodeString(strings.Replace(header, "Basic ", "", -1))
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	parts := strings.SplitN(string(credBytes), ":", 2)
	return parts[0], parts[1], nil
***REMOVED***

func ntlmHandler(username, password string) func(w http.ResponseWriter, r *http.Request) ***REMOVED***
	challenges := make(map[string]*ntlm.ChallengeMessage)
	return func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		// Make sure there is some kind of authentication
		if r.Header.Get("Authorization") == "" ***REMOVED***
			w.Header().Set("WWW-Authenticate", "NTLM")
			w.WriteHeader(401)
			return
		***REMOVED***

		// Parse the proxy authorization header
		auth := r.Header.Get("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		authType := parts[0]
		authPayload := parts[1]

		// Filter out unsupported authentication methods
		if authType != "NTLM" ***REMOVED***
			w.Header().Set("WWW-Authenticate", "NTLM")
			w.WriteHeader(401)
			return
		***REMOVED***

		// Decode base64 auth data and get NTLM message type
		rawAuthPayload, _ := base64.StdEncoding.DecodeString(authPayload)
		ntlmMessageType := binary.LittleEndian.Uint32(rawAuthPayload[8:12])

		// Handle NTLM negotiate message
		if ntlmMessageType == 1 ***REMOVED***
			session, err := ntlm.CreateServerSession(ntlm.Version2, ntlm.ConnectionOrientedMode)
			if err != nil ***REMOVED***
				return
			***REMOVED***

			session.SetUserInfo(username, password, "")

			challenge, err := session.GenerateChallengeMessage()
			if err != nil ***REMOVED***
				return
			***REMOVED***

			challenges[r.RemoteAddr] = challenge

			authPayload := base64.StdEncoding.EncodeToString(challenge.Bytes())

			w.Header().Set("WWW-Authenticate", "NTLM "+authPayload)
			w.WriteHeader(401)

			return
		***REMOVED***

		if ntlmMessageType == 3 ***REMOVED***
			challenge := challenges[r.RemoteAddr]
			if challenge == nil ***REMOVED***
				w.Header().Set("WWW-Authenticate", "NTLM")
				w.WriteHeader(401)
				return
			***REMOVED***

			msg, err := ntlm.ParseAuthenticateMessage(rawAuthPayload, 2)
			if err != nil ***REMOVED***
				msg2, err := ntlm.ParseAuthenticateMessage(rawAuthPayload, 1)

				if err != nil ***REMOVED***
					return
				***REMOVED***

				session, err := ntlm.CreateServerSession(ntlm.Version1, ntlm.ConnectionOrientedMode)
				if err != nil ***REMOVED***
					return
				***REMOVED***

				session.SetServerChallenge(challenge.ServerChallenge)
				session.SetUserInfo(username, password, "")

				err = session.ProcessAuthenticateMessage(msg2)
				if err != nil ***REMOVED***
					w.Header().Set("WWW-Authenticate", "NTLM")
					w.WriteHeader(401)
					return
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				session, err := ntlm.CreateServerSession(ntlm.Version2, ntlm.ConnectionOrientedMode)
				if err != nil ***REMOVED***
					return
				***REMOVED***

				session.SetServerChallenge(challenge.ServerChallenge)
				session.SetUserInfo(username, password, "")

				err = session.ProcessAuthenticateMessage(msg)
				if err != nil ***REMOVED***
					w.Header().Set("WWW-Authenticate", "NTLM")
					w.WriteHeader(401)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***

		data := "authenticated"
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		fmt.Fprint(w, string(data))
	***REMOVED***
***REMOVED***
