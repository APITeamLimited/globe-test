// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"strconv"
)

const AuthUserKey = "user"

type (
	Accounts map[string]string
	authPair struct ***REMOVED***
		Value string
		User  string
	***REMOVED***
	authPairs []authPair
)

func (a authPairs) searchCredential(authValue string) (string, bool) ***REMOVED***
	if len(authValue) == 0 ***REMOVED***
		return "", false
	***REMOVED***
	for _, pair := range a ***REMOVED***
		if pair.Value == authValue ***REMOVED***
			return pair.User, true
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***

// BasicAuthForRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the user name and the value is the password, as well as the name of the Realm.
// If the realm is empty, "Authorization Required" will be used by default.
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc ***REMOVED***
	if realm == "" ***REMOVED***
		realm = "Authorization Required"
	***REMOVED***
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) ***REMOVED***
		// Search user in the slice of allowed credentials
		user, found := pairs.searchCredential(c.Request.Header.Get("Authorization"))
		if !found ***REMOVED***
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(401)
		***REMOVED*** else ***REMOVED***
			// The user credentials was found, set user's id to key AuthUserKey in this context, the userId can be read later using
			// c.MustGet(gin.AuthUserKey)
			c.Set(AuthUserKey, user)
		***REMOVED***
	***REMOVED***
***REMOVED***

// BasicAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
func BasicAuth(accounts Accounts) HandlerFunc ***REMOVED***
	return BasicAuthForRealm(accounts, "")
***REMOVED***

func processAccounts(accounts Accounts) authPairs ***REMOVED***
	assert1(len(accounts) > 0, "Empty list of authorized credentials")
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts ***REMOVED***
		assert1(len(user) > 0, "User can not be empty")
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair***REMOVED***
			Value: value,
			User:  user,
		***REMOVED***)
	***REMOVED***
	return pairs
***REMOVED***

func authorizationHeader(user, password string) string ***REMOVED***
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
***REMOVED***

func secureCompare(given, actual string) bool ***REMOVED***
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 ***REMOVED***
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	***REMOVED***
	/* Securely compare actual to itself to keep constant time, but always return false */
	return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
***REMOVED***
