// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go by Amazon.com, Inc. with code from:
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/request/request.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/signer/v4/v4.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/signer/v4/uri_path.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/types.go
// See THIRD-PARTY-NOTICES for original license terms

package awsv4

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	authorizationHeader     = "Authorization"
	authHeaderSignatureElem = "Signature="

	authHeaderPrefix = "AWS4-HMAC-SHA256"
	timeFormat       = "20060102T150405Z"
	shortTimeFormat  = "20060102"
	awsV4Request     = "aws4_request"

	// emptyStringSHA256 is a SHA256 of an empty string
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
)

var ignoredHeaders = rules***REMOVED***
	denylist***REMOVED***
		mapRule***REMOVED***
			authorizationHeader: struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"User-Agent":        struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			"X-Amzn-Trace-Id":   struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

// Signer applies AWS v4 signing to given request. Use this to sign requests
// that need to be signed with AWS V4 Signatures.
type Signer struct ***REMOVED***
	Credentials *StaticProvider
***REMOVED***

// NewSigner returns a Signer pointer configured with the credentials and optional
// option values provided. If not options are provided the Signer will use its
// default configuration.
func NewSigner(credentials *StaticProvider) *Signer ***REMOVED***
	v4 := &Signer***REMOVED***
		Credentials: credentials,
	***REMOVED***

	return v4
***REMOVED***

type signingCtx struct ***REMOVED***
	ServiceName      string
	Region           string
	Request          *http.Request
	Body             io.ReadSeeker
	Query            url.Values
	Time             time.Time
	SignedHeaderVals http.Header

	credValues Value

	bodyDigest       string
	signedHeaders    string
	canonicalHeaders string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
	authorization    string
***REMOVED***

// Sign signs AWS v4 requests with the provided body, service name, region the
// request is made to, and time the request is signed at. The signTime allows
// you to specify that a request is signed for the future, and cannot be
// used until then.
//
// Returns a list of HTTP headers that were included in the signature or an
// error if signing the request failed. Generally for signed requests this value
// is not needed as the full request context will be captured by the http.Request
// value. It is included for reference though.
//
// Sign will set the request's Body to be the `body` parameter passed in. If
// the body is not already an io.ReadCloser, it will be wrapped within one. If
// a `nil` body parameter passed to Sign, the request's Body field will be
// also set to nil. Its important to note that this functionality will not
// change the request's ContentLength of the request.
//
// Sign differs from Presign in that it will sign the request using HTTP
// header values. This type of signing is intended for http.Request values that
// will not be shared, or are shared in a way the header values on the request
// will not be lost.
//
// The requests body is an io.ReadSeeker so the SHA256 of the body can be
// generated. To bypass the signer computing the hash you can set the
// "X-Amz-Content-Sha256" header with a precomputed value. The signer will
// only compute the hash if the request header value is empty.
func (v4 Signer) Sign(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) ***REMOVED***
	return v4.signWithBody(r, body, service, region, signTime)
***REMOVED***

func (v4 Signer) signWithBody(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) ***REMOVED***
	ctx := &signingCtx***REMOVED***
		Request:     r,
		Body:        body,
		Query:       r.URL.Query(),
		Time:        signTime,
		ServiceName: service,
		Region:      region,
	***REMOVED***

	for key := range ctx.Query ***REMOVED***
		sort.Strings(ctx.Query[key])
	***REMOVED***

	if ctx.isRequestSigned() ***REMOVED***
		ctx.Time = time.Now()
	***REMOVED***

	var err error
	ctx.credValues, err = v4.Credentials.Retrieve()
	if err != nil ***REMOVED***
		return http.Header***REMOVED******REMOVED***, err
	***REMOVED***

	ctx.sanitizeHostForHeader()
	ctx.assignAmzQueryValues()
	if err := ctx.build(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var reader io.ReadCloser
	if body != nil ***REMOVED***
		var ok bool
		if reader, ok = body.(io.ReadCloser); !ok ***REMOVED***
			reader = ioutil.NopCloser(body)
		***REMOVED***
	***REMOVED***
	r.Body = reader

	return ctx.SignedHeaderVals, nil
***REMOVED***

// sanitizeHostForHeader removes default port from host and updates request.Host
func (ctx *signingCtx) sanitizeHostForHeader() ***REMOVED***
	r := ctx.Request
	host := getHost(r)
	port := portOnly(host)
	if port != "" && isDefaultPort(r.URL.Scheme, port) ***REMOVED***
		r.Host = stripPort(host)
	***REMOVED***
***REMOVED***

func (ctx *signingCtx) assignAmzQueryValues() ***REMOVED***
	if ctx.credValues.SessionToken != "" ***REMOVED***
		ctx.Request.Header.Set("X-Amz-Security-Token", ctx.credValues.SessionToken)
	***REMOVED***
***REMOVED***

func (ctx *signingCtx) build() error ***REMOVED***
	ctx.buildTime()             // no depends
	ctx.buildCredentialString() // no depends

	if err := ctx.buildBodyDigest(); err != nil ***REMOVED***
		return err
	***REMOVED***

	unsignedHeaders := ctx.Request.Header

	ctx.buildCanonicalHeaders(ignoredHeaders, unsignedHeaders)
	ctx.buildCanonicalString() // depends on canon headers / signed headers
	ctx.buildStringToSign()    // depends on canon string
	ctx.buildSignature()       // depends on string to sign

	parts := []string***REMOVED***
		authHeaderPrefix + " Credential=" + ctx.credValues.AccessKeyID + "/" + ctx.credentialString,
		"SignedHeaders=" + ctx.signedHeaders,
		authHeaderSignatureElem + ctx.signature,
	***REMOVED***
	ctx.Request.Header.Set(authorizationHeader, strings.Join(parts, ", "))

	return nil
***REMOVED***

// GetSignedRequestSignature attempts to extract the signature of the request.
// Returning an error if the request is unsigned, or unable to extract the
// signature.
func GetSignedRequestSignature(r *http.Request) ([]byte, error) ***REMOVED***

	if auth := r.Header.Get(authorizationHeader); len(auth) != 0 ***REMOVED***
		ps := strings.Split(auth, ", ")
		for _, p := range ps ***REMOVED***
			if idx := strings.Index(p, authHeaderSignatureElem); idx >= 0 ***REMOVED***
				sig := p[len(authHeaderSignatureElem):]
				if len(sig) == 0 ***REMOVED***
					return nil, fmt.Errorf("invalid request signature authorization header")
				***REMOVED***
				return hex.DecodeString(sig)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if sig := r.URL.Query().Get("X-Amz-Signature"); len(sig) != 0 ***REMOVED***
		return hex.DecodeString(sig)
	***REMOVED***

	return nil, fmt.Errorf("request not signed")
***REMOVED***

func (ctx *signingCtx) buildTime() ***REMOVED***
	ctx.Request.Header.Set("X-Amz-Date", formatTime(ctx.Time))
***REMOVED***

func (ctx *signingCtx) buildCredentialString() ***REMOVED***
	ctx.credentialString = buildSigningScope(ctx.Region, ctx.ServiceName, ctx.Time)
***REMOVED***

func (ctx *signingCtx) buildCanonicalHeaders(r rule, header http.Header) ***REMOVED***
	headers := make([]string, 0, len(header))
	headers = append(headers, "host")
	for k, v := range header ***REMOVED***
		if !r.IsValid(k) ***REMOVED***
			continue // ignored header
		***REMOVED***
		if ctx.SignedHeaderVals == nil ***REMOVED***
			ctx.SignedHeaderVals = make(http.Header)
		***REMOVED***

		lowerCaseKey := strings.ToLower(k)
		if _, ok := ctx.SignedHeaderVals[lowerCaseKey]; ok ***REMOVED***
			// include additional values
			ctx.SignedHeaderVals[lowerCaseKey] = append(ctx.SignedHeaderVals[lowerCaseKey], v...)
			continue
		***REMOVED***

		headers = append(headers, lowerCaseKey)
		ctx.SignedHeaderVals[lowerCaseKey] = v
	***REMOVED***
	sort.Strings(headers)

	ctx.signedHeaders = strings.Join(headers, ";")

	headerValues := make([]string, len(headers))
	for i, k := range headers ***REMOVED***
		if k == "host" ***REMOVED***
			if ctx.Request.Host != "" ***REMOVED***
				headerValues[i] = "host:" + ctx.Request.Host
			***REMOVED*** else ***REMOVED***
				headerValues[i] = "host:" + ctx.Request.URL.Host
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			headerValues[i] = k + ":" +
				strings.Join(ctx.SignedHeaderVals[k], ",")
		***REMOVED***
	***REMOVED***
	stripExcessSpaces(headerValues)
	ctx.canonicalHeaders = strings.Join(headerValues, "\n")
***REMOVED***

func getURIPath(u *url.URL) string ***REMOVED***
	var uri string

	if len(u.Opaque) > 0 ***REMOVED***
		uri = "/" + strings.Join(strings.Split(u.Opaque, "/")[3:], "/")
	***REMOVED*** else ***REMOVED***
		uri = u.EscapedPath()
	***REMOVED***

	if len(uri) == 0 ***REMOVED***
		uri = "/"
	***REMOVED***

	return uri
***REMOVED***

func (ctx *signingCtx) buildCanonicalString() ***REMOVED***
	ctx.Request.URL.RawQuery = strings.Replace(ctx.Query.Encode(), "+", "%20", -1)

	uri := getURIPath(ctx.Request.URL)

	uri = EscapePath(uri, false)

	ctx.canonicalString = strings.Join([]string***REMOVED***
		ctx.Request.Method,
		uri,
		ctx.Request.URL.RawQuery,
		ctx.canonicalHeaders + "\n",
		ctx.signedHeaders,
		ctx.bodyDigest,
	***REMOVED***, "\n")
***REMOVED***

func (ctx *signingCtx) buildStringToSign() ***REMOVED***
	ctx.stringToSign = strings.Join([]string***REMOVED***
		authHeaderPrefix,
		formatTime(ctx.Time),
		ctx.credentialString,
		hex.EncodeToString(hashSHA256([]byte(ctx.canonicalString))),
	***REMOVED***, "\n")
***REMOVED***

func (ctx *signingCtx) buildSignature() ***REMOVED***
	creds := deriveSigningKey(ctx.Region, ctx.ServiceName, ctx.credValues.SecretAccessKey, ctx.Time)
	signature := hmacSHA256(creds, []byte(ctx.stringToSign))
	ctx.signature = hex.EncodeToString(signature)
***REMOVED***

func (ctx *signingCtx) buildBodyDigest() error ***REMOVED***
	hash := ctx.Request.Header.Get("X-Amz-Content-Sha256")
	if hash == "" ***REMOVED***
		if ctx.Body == nil ***REMOVED***
			hash = emptyStringSHA256
		***REMOVED*** else ***REMOVED***
			hashBytes, err := makeSha256Reader(ctx.Body)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			hash = hex.EncodeToString(hashBytes)
		***REMOVED***
	***REMOVED***
	ctx.bodyDigest = hash

	return nil
***REMOVED***

// isRequestSigned returns if the request is currently signed or presigned
func (ctx *signingCtx) isRequestSigned() bool ***REMOVED***
	return ctx.Request.Header.Get("Authorization") != ""
***REMOVED***

func hmacSHA256(key []byte, data []byte) []byte ***REMOVED***
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
***REMOVED***

func hashSHA256(data []byte) []byte ***REMOVED***
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
***REMOVED***

// seekerLen attempts to get the number of bytes remaining at the seeker's
// current position.  Returns the number of bytes remaining or error.
func seekerLen(s io.Seeker) (int64, error) ***REMOVED***
	curOffset, err := s.Seek(0, io.SeekCurrent)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	endOffset, err := s.Seek(0, io.SeekEnd)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	_, err = s.Seek(curOffset, io.SeekStart)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return endOffset - curOffset, nil
***REMOVED***

func makeSha256Reader(reader io.ReadSeeker) (hashBytes []byte, err error) ***REMOVED***
	hash := sha256.New()
	start, err := reader.Seek(0, io.SeekCurrent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		// ensure error is return if unable to seek back to start of payload.
		_, err = reader.Seek(start, io.SeekStart)
	***REMOVED***()

	// Use CopyN to avoid allocating the 32KB buffer in io.Copy for bodies
	// smaller than 32KB. Fall back to io.Copy if we fail to determine the size.
	size, err := seekerLen(reader)
	if err != nil ***REMOVED***
		_, _ = io.Copy(hash, reader)
	***REMOVED*** else ***REMOVED***
		_, _ = io.CopyN(hash, reader, size)
	***REMOVED***

	return hash.Sum(nil), nil
***REMOVED***

const doubleSpace = "  "

// stripExcessSpaces will rewrite the passed in slice's string values to not
// contain multiple side-by-side spaces.
func stripExcessSpaces(vals []string) ***REMOVED***
	var j, k, l, m, spaces int
	for i, str := range vals ***REMOVED***
		// Trim trailing spaces
		for j = len(str) - 1; j >= 0 && str[j] == ' '; j-- ***REMOVED***
		***REMOVED***

		// Trim leading spaces
		for k = 0; k < j && str[k] == ' '; k++ ***REMOVED***
		***REMOVED***
		str = str[k : j+1]

		// Strip multiple spaces.
		j = strings.Index(str, doubleSpace)
		if j < 0 ***REMOVED***
			vals[i] = str
			continue
		***REMOVED***

		buf := []byte(str)
		for k, m, l = j, j, len(buf); k < l; k++ ***REMOVED***
			if buf[k] == ' ' ***REMOVED***
				if spaces == 0 ***REMOVED***
					// First space.
					buf[m] = buf[k]
					m++
				***REMOVED***
				spaces++
			***REMOVED*** else ***REMOVED***
				// End of multiple spaces.
				spaces = 0
				buf[m] = buf[k]
				m++
			***REMOVED***
		***REMOVED***

		vals[i] = string(buf[:m])
	***REMOVED***
***REMOVED***

func buildSigningScope(region, service string, dt time.Time) string ***REMOVED***
	return strings.Join([]string***REMOVED***
		formatShortTime(dt),
		region,
		service,
		awsV4Request,
	***REMOVED***, "/")
***REMOVED***

func deriveSigningKey(region, service, secretKey string, dt time.Time) []byte ***REMOVED***
	keyDate := hmacSHA256([]byte("AWS4"+secretKey), []byte(formatShortTime(dt)))
	keyRegion := hmacSHA256(keyDate, []byte(region))
	keyService := hmacSHA256(keyRegion, []byte(service))
	signingKey := hmacSHA256(keyService, []byte(awsV4Request))
	return signingKey
***REMOVED***

func formatShortTime(dt time.Time) string ***REMOVED***
	return dt.UTC().Format(shortTimeFormat)
***REMOVED***

func formatTime(dt time.Time) string ***REMOVED***
	return dt.UTC().Format(timeFormat)
***REMOVED***
