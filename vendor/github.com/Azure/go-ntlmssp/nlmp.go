// Package ntlmssp provides NTLM/Negotiate authentication over HTTP
//
// Protocol details from https://msdn.microsoft.com/en-us/library/cc236621.aspx,
// implementation hints from http://davenport.sourceforge.net/ntlm.html .
// This package only implements authentication, no key exchange or encryption. It
// only supports Unicode (UTF16LE) encoding of protocol strings, no OEM encoding.
// This package implements NTLMv2.
package ntlmssp

import (
	"crypto/hmac"
	"crypto/md5"
	"golang.org/x/crypto/md4"
	"strings"
)

func getNtlmV2Hash(password, username, target string) []byte ***REMOVED***
	return hmacMd5(getNtlmHash(password), toUnicode(strings.ToUpper(username)+target))
***REMOVED***

func getNtlmHash(password string) []byte ***REMOVED***
	hash := md4.New()
	hash.Write(toUnicode(password))
	return hash.Sum(nil)
***REMOVED***

func computeNtlmV2Response(ntlmV2Hash, serverChallenge, clientChallenge,
	timestamp, targetInfo []byte) []byte ***REMOVED***

	temp := []byte***REMOVED***1, 1, 0, 0, 0, 0, 0, 0***REMOVED***
	temp = append(temp, timestamp...)
	temp = append(temp, clientChallenge...)
	temp = append(temp, 0, 0, 0, 0)
	temp = append(temp, targetInfo...)
	temp = append(temp, 0, 0, 0, 0)

	NTProofStr := hmacMd5(ntlmV2Hash, serverChallenge, temp)
	return append(NTProofStr, temp...)
***REMOVED***

func computeLmV2Response(ntlmV2Hash, serverChallenge, clientChallenge []byte) []byte ***REMOVED***
	return append(hmacMd5(ntlmV2Hash, serverChallenge, clientChallenge), clientChallenge...)
***REMOVED***

func hmacMd5(key []byte, data ...[]byte) []byte ***REMOVED***
	mac := hmac.New(md5.New, key)
	for _, d := range data ***REMOVED***
		mac.Write(d)
	***REMOVED***
	return mac.Sum(nil)
***REMOVED***
