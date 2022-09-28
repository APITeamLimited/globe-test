// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// NonceGeneratorFcn defines a function that returns a string of high-quality
// random printable ASCII characters EXCLUDING the comma (',') character.  The
// default nonce generator provides Base64 encoding of 24 bytes from
// crypto/rand.
type NonceGeneratorFcn func() string

// derivedKeys collects the three cryptographically derived values
// into one struct for caching.
type derivedKeys struct ***REMOVED***
	ClientKey []byte
	StoredKey []byte
	ServerKey []byte
***REMOVED***

// KeyFactors represent the two server-provided factors needed to compute
// client credentials for authentication.  Salt is decoded bytes (i.e. not
// base64), but in string form so that KeyFactors can be used as a map key for
// cached credentials.
type KeyFactors struct ***REMOVED***
	Salt  string
	Iters int
***REMOVED***

// StoredCredentials are the values that a server must store for a given
// username to allow authentication.  They include the salt and iteration
// count, plus the derived values to authenticate a client and for the server
// to authenticate itself back to the client.
//
// NOTE: these are specific to a given hash function.  To allow a user to
// authenticate with either SCRAM-SHA-1 or SCRAM-SHA-256, two sets of
// StoredCredentials must be created and stored, one for each hash function.
type StoredCredentials struct ***REMOVED***
	KeyFactors
	StoredKey []byte
	ServerKey []byte
***REMOVED***

// CredentialLookup is a callback to provide StoredCredentials for a given
// username.  This is used to configure Server objects.
//
// NOTE: these are specific to a given hash function.  The callback provided
// to a Server with a given hash function must provide the corresponding
// StoredCredentials.
type CredentialLookup func(string) (StoredCredentials, error)

func defaultNonceGenerator() string ***REMOVED***
	raw := make([]byte, 24)
	nonce := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	rand.Read(raw)
	base64.StdEncoding.Encode(nonce, raw)
	return string(nonce)
***REMOVED***

func encodeName(s string) string ***REMOVED***
	return strings.Replace(strings.Replace(s, "=", "=3D", -1), ",", "=2C", -1)
***REMOVED***

func decodeName(s string) (string, error) ***REMOVED***
	// TODO Check for = not followed by 2C or 3D
	return strings.Replace(strings.Replace(s, "=2C", ",", -1), "=3D", "=", -1), nil
***REMOVED***

func computeHash(hg HashGeneratorFcn, b []byte) []byte ***REMOVED***
	h := hg()
	h.Write(b)
	return h.Sum(nil)
***REMOVED***

func computeHMAC(hg HashGeneratorFcn, key, data []byte) []byte ***REMOVED***
	mac := hmac.New(hg, key)
	mac.Write(data)
	return mac.Sum(nil)
***REMOVED***

func xorBytes(a, b []byte) []byte ***REMOVED***
	// TODO check a & b are same length, or just xor to smallest
	xor := make([]byte, len(a))
	for i := range a ***REMOVED***
		xor[i] = a[i] ^ b[i]
	***REMOVED***
	return xor
***REMOVED***
