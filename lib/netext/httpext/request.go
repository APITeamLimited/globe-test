/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package httpext

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	ntlmssp "github.com/Azure/go-ntlmssp"
	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

// HTTPRequestCookie is a representation of a cookie used for request objects
type HTTPRequestCookie struct ***REMOVED***
	Name, Value string
	Replace     bool
***REMOVED***

// A URL wraps net.URL, and preserves the template (if any) the URL was constructed from.
type URL struct ***REMOVED***
	u    *url.URL
	Name string // http://example.com/thing/$***REMOVED******REMOVED***/
	URL  string // http://example.com/thing/1234/
***REMOVED***

// NewURL returns a new URL for the provided url and name. The error is returned if the url provided
// can't be parsed
func NewURL(urlString, name string) (URL, error) ***REMOVED***
	u, err := url.Parse(urlString)
	return URL***REMOVED***u: u, Name: name, URL: urlString***REMOVED***, err
***REMOVED***

// GetURL returns the internal url.URL
func (u URL) GetURL() *url.URL ***REMOVED***
	return u.u
***REMOVED***

// CompressionType is used to specify what compression is to be used to compress the body of a
// request
// The conversion and validation methods are auto-generated with https://github.com/alvaroloes/enumer:
//nolint: lll
//go:generate enumer -type=CompressionType -transform=snake -trimprefix CompressionType -output compression_type_gen.go
type CompressionType uint

const (
	// CompressionTypeGzip compresses through gzip
	CompressionTypeGzip CompressionType = iota
	// CompressionTypeDeflate compresses through flate
	CompressionTypeDeflate
	// CompressionTypeZstd compresses through zstd
	CompressionTypeZstd
	// CompressionTypeBr compresses through brotli
	CompressionTypeBr
	// TODO: add compress(lzw), maybe bzip2 and others listed at
	// https://en.wikipedia.org/wiki/HTTP_compression#Content-Encoding_tokens
)

// Request represent an http request
type Request struct ***REMOVED***
	Method  string                          `json:"method"`
	URL     string                          `json:"url"`
	Headers map[string][]string             `json:"headers"`
	Body    string                          `json:"body"`
	Cookies map[string][]*HTTPRequestCookie `json:"cookies"`
***REMOVED***

// ParsedHTTPRequest a represantion of a request after it has been parsed from a user script
type ParsedHTTPRequest struct ***REMOVED***
	URL          *URL
	Body         *bytes.Buffer
	Req          *http.Request
	Timeout      time.Duration
	Auth         string
	Throw        bool
	ResponseType ResponseType
	Compressions []CompressionType
	Redirects    null.Int
	ActiveJar    *cookiejar.Jar
	Cookies      map[string]*HTTPRequestCookie
	Tags         map[string]string
***REMOVED***

// Matches non-compliant io.Closer implementations (e.g. zstd.Decoder)
type ncloser interface ***REMOVED***
	Close()
***REMOVED***

type readCloser struct ***REMOVED***
	io.Reader
***REMOVED***

// Close readers with differing Close() implementations
func (r readCloser) Close() error ***REMOVED***
	var err error
	switch v := r.Reader.(type) ***REMOVED***
	case io.Closer:
		err = v.Close()
	case ncloser:
		v.Close()
	***REMOVED***
	return err
***REMOVED***

func stdCookiesToHTTPRequestCookies(cookies []*http.Cookie) map[string][]*HTTPRequestCookie ***REMOVED***
	var result = make(map[string][]*HTTPRequestCookie, len(cookies))
	for _, cookie := range cookies ***REMOVED***
		result[cookie.Name] = append(result[cookie.Name],
			&HTTPRequestCookie***REMOVED***Name: cookie.Name, Value: cookie.Value***REMOVED***)
	***REMOVED***
	return result
***REMOVED***

func compressBody(algos []CompressionType, body io.ReadCloser) (*bytes.Buffer, string, error) ***REMOVED***
	var contentEncoding string
	var prevBuf io.Reader = body
	var buf *bytes.Buffer
	for _, compressionType := range algos ***REMOVED***
		if buf != nil ***REMOVED***
			prevBuf = buf
		***REMOVED***
		buf = new(bytes.Buffer)

		if contentEncoding != "" ***REMOVED***
			contentEncoding += ", "
		***REMOVED***
		contentEncoding += compressionType.String()
		var w io.WriteCloser
		switch compressionType ***REMOVED***
		case CompressionTypeGzip:
			w = gzip.NewWriter(buf)
		case CompressionTypeDeflate:
			w = zlib.NewWriter(buf)
		case CompressionTypeZstd:
			w, _ = zstd.NewWriter(buf)
		case CompressionTypeBr:
			w = brotli.NewWriter(buf)
		default:
			return nil, "", fmt.Errorf("unknown compressionType %s", compressionType)
		***REMOVED***
		// we don't close in defer because zlib will write it's checksum again if it closes twice :(
		var _, err = io.Copy(w, prevBuf)
		if err != nil ***REMOVED***
			_ = w.Close()
			return nil, "", err
		***REMOVED***

		if err = w.Close(); err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
	***REMOVED***

	return buf, contentEncoding, body.Close()
***REMOVED***

//nolint:gochecknoglobals
var decompressionErrors = [...]error***REMOVED***
	zlib.ErrChecksum, zlib.ErrDictionary, zlib.ErrHeader,
	gzip.ErrChecksum, gzip.ErrHeader,
	//TODO: handle brotli errors - currently unexported
	zstd.ErrReservedBlockType, zstd.ErrCompressedSizeTooBig, zstd.ErrBlockTooSmall, zstd.ErrMagicMismatch,
	zstd.ErrWindowSizeExceeded, zstd.ErrWindowSizeTooSmall, zstd.ErrDecoderSizeExceeded, zstd.ErrUnknownDictionary,
	zstd.ErrFrameSizeExceeded, zstd.ErrCRCMismatch, zstd.ErrDecoderClosed,
***REMOVED***

func newDecompressionError(originalErr error) K6Error ***REMOVED***
	return NewK6Error(
		responseDecompressionErrorCode,
		fmt.Sprintf("error decompressing response body (%s)", originalErr.Error()),
		originalErr,
	)
***REMOVED***

func wrapDecompressionError(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	// TODO: something more optimized? for example, we won't get zstd errors if
	// we don't use it... maybe the code that builds the decompression readers
	// could also add an appropriate error-wrapper layer?
	for _, decErr := range decompressionErrors ***REMOVED***
		if err == decErr ***REMOVED***
			return newDecompressionError(err)
		***REMOVED***
	***REMOVED***
	if strings.HasPrefix(err.Error(), "brotli: ") ***REMOVED*** //TODO: submit an upstream patch and fix...
		return newDecompressionError(err)
	***REMOVED***
	return err
***REMOVED***

func readResponseBody(
	state *lib.State, respType ResponseType, resp *http.Response, respErr error,
) (interface***REMOVED******REMOVED***, error) ***REMOVED***

	if resp == nil || respErr != nil ***REMOVED***
		return nil, respErr
	***REMOVED***

	if respType == ResponseTypeNone ***REMOVED***
		_, err := io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
		if err != nil ***REMOVED***
			respErr = err
		***REMOVED***
		return nil, respErr
	***REMOVED***

	rc := &readCloser***REMOVED***resp.Body***REMOVED***
	// Ensure that the entire response body is read and closed, e.g. in case of decoding errors
	defer func(respBody io.ReadCloser) ***REMOVED***
		_, _ = io.Copy(ioutil.Discard, respBody)
		_ = respBody.Close()
	***REMOVED***(resp.Body)

	// Transparently decompress the body if it's has a content-encoding we
	// support. If not, simply return it as it is.
	contentEncoding := strings.TrimSpace(resp.Header.Get("Content-Encoding"))
	//TODO: support stacked compressions, e.g. `deflate, gzip`
	if compression, err := CompressionTypeString(contentEncoding); err == nil ***REMOVED***
		var decoder io.Reader
		var err error
		switch compression ***REMOVED***
		case CompressionTypeDeflate:
			decoder, err = zlib.NewReader(resp.Body)
		case CompressionTypeGzip:
			decoder, err = gzip.NewReader(resp.Body)
		case CompressionTypeZstd:
			decoder, err = zstd.NewReader(resp.Body)
		case CompressionTypeBr:
			decoder = brotli.NewReader(resp.Body)
		default:
			// We have not implemented a compression ... :(
			err = fmt.Errorf(
				"unsupported compression type %s - this is a bug in k6, please report it",
				compression,
			)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, newDecompressionError(err)
		***REMOVED***
		rc = &readCloser***REMOVED***decoder***REMOVED***
	***REMOVED***

	buf := state.BPool.Get()
	defer state.BPool.Put(buf)
	buf.Reset()
	_, err := io.Copy(buf, rc.Reader)
	if err != nil ***REMOVED***
		respErr = wrapDecompressionError(err)
	***REMOVED***

	err = rc.Close()
	if err != nil && respErr == nil ***REMOVED*** // Don't overwrite previous errors
		respErr = wrapDecompressionError(err)
	***REMOVED***

	var result interface***REMOVED******REMOVED***
	// Binary or string
	switch respType ***REMOVED***
	case ResponseTypeText:
		result = buf.String()
	case ResponseTypeBinary:
		// Copy the data to a new slice before we return the buffer to the pool,
		// because buf.Bytes() points to the underlying buffer byte slice.
		binData := make([]byte, buf.Len())
		copy(binData, buf.Bytes())
		result = binData
	default:
		respErr = fmt.Errorf("unknown responseType %s", respType)
	***REMOVED***

	return result, respErr
***REMOVED***

//TODO: move as a response method? or constructor?
func updateK6Response(k6Response *Response, finishedReq *finishedRequest) ***REMOVED***
	k6Response.ErrorCode = int(finishedReq.errorCode)
	k6Response.Error = finishedReq.errorMsg
	trail := finishedReq.trail

	if trail.ConnRemoteAddr != nil ***REMOVED***
		remoteHost, remotePortStr, _ := net.SplitHostPort(trail.ConnRemoteAddr.String())
		remotePort, _ := strconv.Atoi(remotePortStr)
		k6Response.RemoteIP = remoteHost
		k6Response.RemotePort = remotePort
	***REMOVED***
	k6Response.Timings = ResponseTimings***REMOVED***
		Duration:       stats.D(trail.Duration),
		Blocked:        stats.D(trail.Blocked),
		Connecting:     stats.D(trail.Connecting),
		TLSHandshaking: stats.D(trail.TLSHandshaking),
		Sending:        stats.D(trail.Sending),
		Waiting:        stats.D(trail.Waiting),
		Receiving:      stats.D(trail.Receiving),
	***REMOVED***
***REMOVED***

// MakeRequest makes http request for tor the provided ParsedHTTPRequest
func MakeRequest(ctx context.Context, preq *ParsedHTTPRequest) (*Response, error) ***REMOVED***
	state := lib.GetState(ctx)

	respReq := &Request***REMOVED***
		Method:  preq.Req.Method,
		URL:     preq.Req.URL.String(),
		Cookies: stdCookiesToHTTPRequestCookies(preq.Req.Cookies()),
		Headers: preq.Req.Header,
	***REMOVED***

	if preq.Body != nil ***REMOVED***
		// TODO: maybe hide this behind of flag in order for this to not happen for big post/puts?
		// should we set this after the compression? what will be the point ?
		respReq.Body = preq.Body.String()

		if len(preq.Compressions) > 0 ***REMOVED***
			compressedBody, contentEncoding, err := compressBody(preq.Compressions, ioutil.NopCloser(preq.Body))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			preq.Body = compressedBody

			currentContentEncoding := preq.Req.Header.Get("Content-Encoding")
			if currentContentEncoding == "" ***REMOVED***
				preq.Req.Header.Set("Content-Encoding", contentEncoding)
			***REMOVED*** else if currentContentEncoding != contentEncoding ***REMOVED***
				state.Logger.Warningf(
					"There's a mismatch between the desired `compression` the manually set `Content-Encoding` header "+
						"in the %s request for '%s', the custom header has precedence and won't be overwritten. "+
						"This may result in invalid data being sent to the server.", preq.Req.Method, preq.Req.URL,
				)
			***REMOVED***
		***REMOVED***

		preq.Req.ContentLength = int64(preq.Body.Len()) // This will make Go set the content-length header
		preq.Req.GetBody = func() (io.ReadCloser, error) ***REMOVED***
			//  using `Bytes()` should reuse the same buffer and as such help with the memory usage. We
			//  should not be writing to it any way so there shouldn't be way to corrupt it (?)
			return ioutil.NopCloser(bytes.NewBuffer(preq.Body.Bytes())), nil
		***REMOVED***
		// as per the documentation using GetBody still requires setting the Body.
		preq.Req.Body, _ = preq.Req.GetBody()
	***REMOVED***

	if contentLengthHeader := preq.Req.Header.Get("Content-Length"); contentLengthHeader != "" ***REMOVED***
		// The content-length header was set by the user, delete it (since Go
		// will set it automatically) and warn if there were differences
		preq.Req.Header.Del("Content-Length")
		length, err := strconv.Atoi(contentLengthHeader)
		if err != nil || preq.Req.ContentLength != int64(length) ***REMOVED***
			state.Logger.Warnf(
				"The specified Content-Length header %q in the %s request for %s "+
					"doesn't match the actual request body length of %d, so it will be ignored!",
				contentLengthHeader, preq.Req.Method, preq.Req.URL, preq.Req.ContentLength,
			)
		***REMOVED***
	***REMOVED***

	tags := state.Options.RunTags.CloneTags()
	for k, v := range preq.Tags ***REMOVED***
		tags[k] = v
	***REMOVED***

	if state.Options.SystemTags["method"] ***REMOVED***
		tags["method"] = preq.Req.Method
	***REMOVED***
	if state.Options.SystemTags["url"] ***REMOVED***
		tags["url"] = preq.URL.URL
	***REMOVED***

	// Only set the name system tag if the user didn't explicitly set it beforehand
	if _, ok := tags["name"]; !ok && state.Options.SystemTags["name"] ***REMOVED***
		tags["name"] = preq.URL.Name
	***REMOVED***
	if state.Options.SystemTags["group"] ***REMOVED***
		tags["group"] = state.Group.Path
	***REMOVED***
	if state.Options.SystemTags["vu"] ***REMOVED***
		tags["vu"] = strconv.FormatInt(state.Vu, 10)
	***REMOVED***
	if state.Options.SystemTags["iter"] ***REMOVED***
		tags["iter"] = strconv.FormatInt(state.Iteration, 10)
	***REMOVED***

	// Check rate limit *after* we've prepared a request; no need to wait with that part.
	if rpsLimit := state.RPSLimit; rpsLimit != nil ***REMOVED***
		if err := rpsLimit.Wait(ctx); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	tracerTransport := newTransport(state, tags)
	var transport http.RoundTripper = tracerTransport

	if state.Options.HttpDebug.String != "" ***REMOVED***
		transport = httpDebugTransport***REMOVED***
			originalTransport: transport,
			httpDebugOption:   state.Options.HttpDebug.String,
		***REMOVED***
	***REMOVED***

	if preq.Auth == "digest" ***REMOVED***
		transport = digestTransport***REMOVED***originalTransport: transport***REMOVED***
	***REMOVED*** else if preq.Auth == "ntlm" ***REMOVED***
		transport = ntlmssp.Negotiator***REMOVED***RoundTripper: transport***REMOVED***
	***REMOVED***

	resp := &Response***REMOVED***ctx: ctx, URL: preq.URL.URL, Request: *respReq***REMOVED***
	client := http.Client***REMOVED***
		Transport: transport,
		Timeout:   preq.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error ***REMOVED***
			resp.URL = req.URL.String()

			// Update active jar with cookies found in "Set-Cookie" header(s) of redirect response
			if preq.ActiveJar != nil ***REMOVED***
				if respCookies := req.Response.Cookies(); len(respCookies) > 0 ***REMOVED***
					preq.ActiveJar.SetCookies(req.URL, respCookies)
				***REMOVED***
				req.Header.Del("Cookie")
				SetRequestCookies(req, preq.ActiveJar, preq.Cookies)
			***REMOVED***

			if l := len(via); int64(l) > preq.Redirects.Int64 ***REMOVED***
				if !preq.Redirects.Valid ***REMOVED***
					url := req.URL
					if l > 0 ***REMOVED***
						url = via[0].URL
					***REMOVED***
					state.Logger.WithFields(log.Fields***REMOVED***"url": url.String()***REMOVED***).Warnf(
						"Stopped after %d redirects and returned the redirection; pass ***REMOVED*** redirects: n ***REMOVED***"+
							" in request params or set global maxRedirects to silence this", l)
				***REMOVED***
				return http.ErrUseLastResponse
			***REMOVED***
			return nil
		***REMOVED***,
	***REMOVED***

	mreq := preq.Req.WithContext(ctx)
	res, resErr := client.Do(mreq)

	resp.Body, resErr = readResponseBody(state, preq.ResponseType, res, resErr)
	finishedReq := tracerTransport.processLastSavedRequest(wrapDecompressionError(resErr))
	if finishedReq != nil ***REMOVED***
		updateK6Response(resp, finishedReq)
	***REMOVED***

	if resErr == nil ***REMOVED***
		if preq.ActiveJar != nil ***REMOVED***
			if rc := res.Cookies(); len(rc) > 0 ***REMOVED***
				preq.ActiveJar.SetCookies(res.Request.URL, rc)
			***REMOVED***
		***REMOVED***

		resp.URL = res.Request.URL.String()
		resp.Status = res.StatusCode
		resp.Proto = res.Proto

		if res.TLS != nil ***REMOVED***
			resp.setTLSInfo(res.TLS)
		***REMOVED***

		resp.Headers = make(map[string]string, len(res.Header))
		for k, vs := range res.Header ***REMOVED***
			resp.Headers[k] = strings.Join(vs, ", ")
		***REMOVED***

		resCookies := res.Cookies()
		resp.Cookies = make(map[string][]*HTTPCookie, len(resCookies))
		for _, c := range resCookies ***REMOVED***
			resp.Cookies[c.Name] = append(resp.Cookies[c.Name], &HTTPCookie***REMOVED***
				Name:     c.Name,
				Value:    c.Value,
				Domain:   c.Domain,
				Path:     c.Path,
				HTTPOnly: c.HttpOnly,
				Secure:   c.Secure,
				MaxAge:   c.MaxAge,
				Expires:  c.Expires.UnixNano() / 1000000,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	if resErr != nil ***REMOVED***
		// Do *not* log errors about the contex being cancelled.
		select ***REMOVED***
		case <-ctx.Done():
		default:
			state.Logger.WithField("error", resErr).Warn("Request Failed")
		***REMOVED***

		if preq.Throw ***REMOVED***
			return nil, resErr
		***REMOVED***
	***REMOVED***

	return resp, nil
***REMOVED***

// SetRequestCookies sets the cookies of the requests getting those cookies both from the jar and
// from the reqCookies map. The Replace field of the HTTPRequestCookie will be taken into account
func SetRequestCookies(req *http.Request, jar *cookiejar.Jar, reqCookies map[string]*HTTPRequestCookie) ***REMOVED***
	var replacedCookies = make(map[string]struct***REMOVED******REMOVED***)
	for key, reqCookie := range reqCookies ***REMOVED***
		req.AddCookie(&http.Cookie***REMOVED***Name: key, Value: reqCookie.Value***REMOVED***)
		if reqCookie.Replace ***REMOVED***
			replacedCookies[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	for _, c := range jar.Cookies(req.URL) ***REMOVED***
		if _, ok := replacedCookies[c.Name]; !ok ***REMOVED***
			req.AddCookie(&http.Cookie***REMOVED***Name: c.Name, Value: c.Value***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***
