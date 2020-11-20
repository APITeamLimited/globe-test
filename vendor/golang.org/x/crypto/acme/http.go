// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acme

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// retryTimer encapsulates common logic for retrying unsuccessful requests.
// It is not safe for concurrent use.
type retryTimer struct ***REMOVED***
	// backoffFn provides backoff delay sequence for retries.
	// See Client.RetryBackoff doc comment.
	backoffFn func(n int, r *http.Request, res *http.Response) time.Duration
	// n is the current retry attempt.
	n int
***REMOVED***

func (t *retryTimer) inc() ***REMOVED***
	t.n++
***REMOVED***

// backoff pauses the current goroutine as described in Client.RetryBackoff.
func (t *retryTimer) backoff(ctx context.Context, r *http.Request, res *http.Response) error ***REMOVED***
	d := t.backoffFn(t.n, r, res)
	if d <= 0 ***REMOVED***
		return fmt.Errorf("acme: no more retries for %s; tried %d time(s)", r.URL, t.n)
	***REMOVED***
	wakeup := time.NewTimer(d)
	defer wakeup.Stop()
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	case <-wakeup.C:
		return nil
	***REMOVED***
***REMOVED***

func (c *Client) retryTimer() *retryTimer ***REMOVED***
	f := c.RetryBackoff
	if f == nil ***REMOVED***
		f = defaultBackoff
	***REMOVED***
	return &retryTimer***REMOVED***backoffFn: f***REMOVED***
***REMOVED***

// defaultBackoff provides default Client.RetryBackoff implementation
// using a truncated exponential backoff algorithm,
// as described in Client.RetryBackoff.
//
// The n argument is always bounded between 1 and 30.
// The returned value is always greater than 0.
func defaultBackoff(n int, r *http.Request, res *http.Response) time.Duration ***REMOVED***
	const max = 10 * time.Second
	var jitter time.Duration
	if x, err := rand.Int(rand.Reader, big.NewInt(1000)); err == nil ***REMOVED***
		// Set the minimum to 1ms to avoid a case where
		// an invalid Retry-After value is parsed into 0 below,
		// resulting in the 0 returned value which would unintentionally
		// stop the retries.
		jitter = (1 + time.Duration(x.Int64())) * time.Millisecond
	***REMOVED***
	if v, ok := res.Header["Retry-After"]; ok ***REMOVED***
		return retryAfter(v[0]) + jitter
	***REMOVED***

	if n < 1 ***REMOVED***
		n = 1
	***REMOVED***
	if n > 30 ***REMOVED***
		n = 30
	***REMOVED***
	d := time.Duration(1<<uint(n-1))*time.Second + jitter
	if d > max ***REMOVED***
		return max
	***REMOVED***
	return d
***REMOVED***

// retryAfter parses a Retry-After HTTP header value,
// trying to convert v into an int (seconds) or use http.ParseTime otherwise.
// It returns zero value if v cannot be parsed.
func retryAfter(v string) time.Duration ***REMOVED***
	if i, err := strconv.Atoi(v); err == nil ***REMOVED***
		return time.Duration(i) * time.Second
	***REMOVED***
	t, err := http.ParseTime(v)
	if err != nil ***REMOVED***
		return 0
	***REMOVED***
	return t.Sub(timeNow())
***REMOVED***

// resOkay is a function that reports whether the provided response is okay.
// It is expected to keep the response body unread.
type resOkay func(*http.Response) bool

// wantStatus returns a function which reports whether the code
// matches the status code of a response.
func wantStatus(codes ...int) resOkay ***REMOVED***
	return func(res *http.Response) bool ***REMOVED***
		for _, code := range codes ***REMOVED***
			if code == res.StatusCode ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

// get issues an unsigned GET request to the specified URL.
// It returns a non-error value only when ok reports true.
//
// get retries unsuccessful attempts according to c.RetryBackoff
// until the context is done or a non-retriable error is received.
func (c *Client) get(ctx context.Context, url string, ok resOkay) (*http.Response, error) ***REMOVED***
	retry := c.retryTimer()
	for ***REMOVED***
		req, err := http.NewRequest("GET", url, nil)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		res, err := c.doNoRetry(ctx, req)
		switch ***REMOVED***
		case err != nil:
			return nil, err
		case ok(res):
			return res, nil
		case isRetriable(res.StatusCode):
			retry.inc()
			resErr := responseError(res)
			res.Body.Close()
			// Ignore the error value from retry.backoff
			// and return the one from last retry, as received from the CA.
			if retry.backoff(ctx, req, res) != nil ***REMOVED***
				return nil, resErr
			***REMOVED***
		default:
			defer res.Body.Close()
			return nil, responseError(res)
		***REMOVED***
	***REMOVED***
***REMOVED***

// postAsGet is POST-as-GET, a replacement for GET in RFC8555
// as described in https://tools.ietf.org/html/rfc8555#section-6.3.
// It makes a POST request in KID form with zero JWS payload.
// See nopayload doc comments in jws.go.
func (c *Client) postAsGet(ctx context.Context, url string, ok resOkay) (*http.Response, error) ***REMOVED***
	return c.post(ctx, nil, url, noPayload, ok)
***REMOVED***

// post issues a signed POST request in JWS format using the provided key
// to the specified URL. If key is nil, c.Key is used instead.
// It returns a non-error value only when ok reports true.
//
// post retries unsuccessful attempts according to c.RetryBackoff
// until the context is done or a non-retriable error is received.
// It uses postNoRetry to make individual requests.
func (c *Client) post(ctx context.Context, key crypto.Signer, url string, body interface***REMOVED******REMOVED***, ok resOkay) (*http.Response, error) ***REMOVED***
	retry := c.retryTimer()
	for ***REMOVED***
		res, req, err := c.postNoRetry(ctx, key, url, body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if ok(res) ***REMOVED***
			return res, nil
		***REMOVED***
		resErr := responseError(res)
		res.Body.Close()
		switch ***REMOVED***
		// Check for bad nonce before isRetriable because it may have been returned
		// with an unretriable response code such as 400 Bad Request.
		case isBadNonce(resErr):
			// Consider any previously stored nonce values to be invalid.
			c.clearNonces()
		case !isRetriable(res.StatusCode):
			return nil, resErr
		***REMOVED***
		retry.inc()
		// Ignore the error value from retry.backoff
		// and return the one from last retry, as received from the CA.
		if err := retry.backoff(ctx, req, res); err != nil ***REMOVED***
			return nil, resErr
		***REMOVED***
	***REMOVED***
***REMOVED***

// postNoRetry signs the body with the given key and POSTs it to the provided url.
// It is used by c.post to retry unsuccessful attempts.
// The body argument must be JSON-serializable.
//
// If key argument is nil, c.Key is used to sign the request.
// If key argument is nil and c.accountKID returns a non-zero keyID,
// the request is sent in KID form. Otherwise, JWK form is used.
//
// In practice, when interfacing with RFC-compliant CAs most requests are sent in KID form
// and JWK is used only when KID is unavailable: new account endpoint and certificate
// revocation requests authenticated by a cert key.
// See jwsEncodeJSON for other details.
func (c *Client) postNoRetry(ctx context.Context, key crypto.Signer, url string, body interface***REMOVED******REMOVED***) (*http.Response, *http.Request, error) ***REMOVED***
	kid := noKeyID
	if key == nil ***REMOVED***
		key = c.Key
		kid = c.accountKID(ctx)
	***REMOVED***
	nonce, err := c.popNonce(ctx, url)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	b, err := jwsEncodeJSON(body, key, kid, nonce, url)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	req.Header.Set("Content-Type", "application/jose+json")
	res, err := c.doNoRetry(ctx, req)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	c.addNonce(res.Header)
	return res, req, nil
***REMOVED***

// doNoRetry issues a request req, replacing its context (if any) with ctx.
func (c *Client) doNoRetry(ctx context.Context, req *http.Request) (*http.Response, error) ***REMOVED***
	req.Header.Set("User-Agent", c.userAgent())
	res, err := c.httpClient().Do(req.WithContext(ctx))
	if err != nil ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			// Prefer the unadorned context error.
			// (The acme package had tests assuming this, previously from ctxhttp's
			// behavior, predating net/http supporting contexts natively)
			// TODO(bradfitz): reconsider this in the future. But for now this
			// requires no test updates.
			return nil, ctx.Err()
		default:
			return nil, err
		***REMOVED***
	***REMOVED***
	return res, nil
***REMOVED***

func (c *Client) httpClient() *http.Client ***REMOVED***
	if c.HTTPClient != nil ***REMOVED***
		return c.HTTPClient
	***REMOVED***
	return http.DefaultClient
***REMOVED***

// packageVersion is the version of the module that contains this package, for
// sending as part of the User-Agent header. It's set in version_go112.go.
var packageVersion string

// userAgent returns the User-Agent header value. It includes the package name,
// the module version (if available), and the c.UserAgent value (if set).
func (c *Client) userAgent() string ***REMOVED***
	ua := "golang.org/x/crypto/acme"
	if packageVersion != "" ***REMOVED***
		ua += "@" + packageVersion
	***REMOVED***
	if c.UserAgent != "" ***REMOVED***
		ua = c.UserAgent + " " + ua
	***REMOVED***
	return ua
***REMOVED***

// isBadNonce reports whether err is an ACME "badnonce" error.
func isBadNonce(err error) bool ***REMOVED***
	// According to the spec badNonce is urn:ietf:params:acme:error:badNonce.
	// However, ACME servers in the wild return their versions of the error.
	// See https://tools.ietf.org/html/draft-ietf-acme-acme-02#section-5.4
	// and https://github.com/letsencrypt/boulder/blob/0e07eacb/docs/acme-divergences.md#section-66.
	ae, ok := err.(*Error)
	return ok && strings.HasSuffix(strings.ToLower(ae.ProblemType), ":badnonce")
***REMOVED***

// isRetriable reports whether a request can be retried
// based on the response status code.
//
// Note that a "bad nonce" error is returned with a non-retriable 400 Bad Request code.
// Callers should parse the response and check with isBadNonce.
func isRetriable(code int) bool ***REMOVED***
	return code <= 399 || code >= 500 || code == http.StatusTooManyRequests
***REMOVED***

// responseError creates an error of Error type from resp.
func responseError(resp *http.Response) error ***REMOVED***
	// don't care if ReadAll returns an error:
	// json.Unmarshal will fail in that case anyway
	b, _ := ioutil.ReadAll(resp.Body)
	e := &wireError***REMOVED***Status: resp.StatusCode***REMOVED***
	if err := json.Unmarshal(b, e); err != nil ***REMOVED***
		// this is not a regular error response:
		// populate detail with anything we received,
		// e.Status will already contain HTTP response code value
		e.Detail = string(b)
		if e.Detail == "" ***REMOVED***
			e.Detail = resp.Status
		***REMOVED***
	***REMOVED***
	return e.error(resp.Header)
***REMOVED***
