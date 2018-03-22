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
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ThomsonReutersEikon/go-ntlm/ntlm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var challenges map[string]*ntlm.ChallengeMessage

func TestGetCredentialsFromHeader(t *testing.T) ***REMOVED***
	user, pass, err := getCredentialsFromHeader("Basic Ym9iOnBhc3M=")
	require.NoError(t, err)

	assert.Equal(t, "bob", user)
	assert.Equal(t, "pass", pass)
***REMOVED***

func TestNTLMServer(t *testing.T) ***REMOVED***
	s := httptest.NewServer(http.HandlerFunc(ntlmHandler("bob", "pass")))
	defer s.Close()

	client := &http.Client***REMOVED***
		Transport: NTLMNegotiator***REMOVED***
			RoundTripper: &http.Transport***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***

	url := strings.Replace(s.URL, "http://", "http://bob:pass@", -1)

	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	require.NoError(t, err)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "authenticated", string(body))
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
		fmt.Fprint(w, data)
	***REMOVED***
***REMOVED***
