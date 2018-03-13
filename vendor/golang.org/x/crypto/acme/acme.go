// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package acme provides an implementation of the
// Automatic Certificate Management Environment (ACME) spec.
// See https://tools.ietf.org/html/draft-ietf-acme-acme-02 for details.
//
// Most common scenarios will want to use autocert subdirectory instead,
// which provides automatic access to certificates from Let's Encrypt
// and any other ACME-based CA.
//
// This package is a work in progress and makes no API stability promises.
package acme

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LetsEncryptURL is the Directory endpoint of Let's Encrypt CA.
const LetsEncryptURL = "https://acme-v01.api.letsencrypt.org/directory"

const (
	maxChainLen = 5       // max depth and breadth of a certificate chain
	maxCertSize = 1 << 20 // max size of a certificate, in bytes

	// Max number of collected nonces kept in memory.
	// Expect usual peak of 1 or 2.
	maxNonces = 100
)

// Client is an ACME client.
// The only required field is Key. An example of creating a client with a new key
// is as follows:
//
// 	key, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil ***REMOVED***
// 		log.Fatal(err)
// 	***REMOVED***
// 	client := &Client***REMOVED***Key: key***REMOVED***
//
type Client struct ***REMOVED***
	// Key is the account key used to register with a CA and sign requests.
	// Key.Public() must return a *rsa.PublicKey or *ecdsa.PublicKey.
	Key crypto.Signer

	// HTTPClient optionally specifies an HTTP client to use
	// instead of http.DefaultClient.
	HTTPClient *http.Client

	// DirectoryURL points to the CA directory endpoint.
	// If empty, LetsEncryptURL is used.
	// Mutating this value after a successful call of Client's Discover method
	// will have no effect.
	DirectoryURL string

	dirMu sync.Mutex // guards writes to dir
	dir   *Directory // cached result of Client's Discover method

	noncesMu sync.Mutex
	nonces   map[string]struct***REMOVED******REMOVED*** // nonces collected from previous responses
***REMOVED***

// Discover performs ACME server discovery using c.DirectoryURL.
//
// It caches successful result. So, subsequent calls will not result in
// a network round-trip. This also means mutating c.DirectoryURL after successful call
// of this method will have no effect.
func (c *Client) Discover(ctx context.Context) (Directory, error) ***REMOVED***
	c.dirMu.Lock()
	defer c.dirMu.Unlock()
	if c.dir != nil ***REMOVED***
		return *c.dir, nil
	***REMOVED***

	dirURL := c.DirectoryURL
	if dirURL == "" ***REMOVED***
		dirURL = LetsEncryptURL
	***REMOVED***
	res, err := c.get(ctx, dirURL)
	if err != nil ***REMOVED***
		return Directory***REMOVED******REMOVED***, err
	***REMOVED***
	defer res.Body.Close()
	c.addNonce(res.Header)
	if res.StatusCode != http.StatusOK ***REMOVED***
		return Directory***REMOVED******REMOVED***, responseError(res)
	***REMOVED***

	var v struct ***REMOVED***
		Reg    string `json:"new-reg"`
		Authz  string `json:"new-authz"`
		Cert   string `json:"new-cert"`
		Revoke string `json:"revoke-cert"`
		Meta   struct ***REMOVED***
			Terms   string   `json:"terms-of-service"`
			Website string   `json:"website"`
			CAA     []string `json:"caa-identities"`
		***REMOVED***
	***REMOVED***
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return Directory***REMOVED******REMOVED***, err
	***REMOVED***
	c.dir = &Directory***REMOVED***
		RegURL:    v.Reg,
		AuthzURL:  v.Authz,
		CertURL:   v.Cert,
		RevokeURL: v.Revoke,
		Terms:     v.Meta.Terms,
		Website:   v.Meta.Website,
		CAA:       v.Meta.CAA,
	***REMOVED***
	return *c.dir, nil
***REMOVED***

// CreateCert requests a new certificate using the Certificate Signing Request csr encoded in DER format.
// The exp argument indicates the desired certificate validity duration. CA may issue a certificate
// with a different duration.
// If the bundle argument is true, the returned value will also contain the CA (issuer) certificate chain.
//
// In the case where CA server does not provide the issued certificate in the response,
// CreateCert will poll certURL using c.FetchCert, which will result in additional round-trips.
// In such a scenario, the caller can cancel the polling with ctx.
//
// CreateCert returns an error if the CA's response or chain was unreasonably large.
// Callers are encouraged to parse the returned value to ensure the certificate is valid and has the expected features.
func (c *Client) CreateCert(ctx context.Context, csr []byte, exp time.Duration, bundle bool) (der [][]byte, certURL string, err error) ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***

	req := struct ***REMOVED***
		Resource  string `json:"resource"`
		CSR       string `json:"csr"`
		NotBefore string `json:"notBefore,omitempty"`
		NotAfter  string `json:"notAfter,omitempty"`
	***REMOVED******REMOVED***
		Resource: "new-cert",
		CSR:      base64.RawURLEncoding.EncodeToString(csr),
	***REMOVED***
	now := timeNow()
	req.NotBefore = now.Format(time.RFC3339)
	if exp > 0 ***REMOVED***
		req.NotAfter = now.Add(exp).Format(time.RFC3339)
	***REMOVED***

	res, err := c.retryPostJWS(ctx, c.Key, c.dir.CertURL, req)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated ***REMOVED***
		return nil, "", responseError(res)
	***REMOVED***

	curl := res.Header.Get("Location") // cert permanent URL
	if res.ContentLength == 0 ***REMOVED***
		// no cert in the body; poll until we get it
		cert, err := c.FetchCert(ctx, curl, bundle)
		return cert, curl, err
	***REMOVED***
	// slurp issued cert and CA chain, if requested
	cert, err := c.responseCert(ctx, res, bundle)
	return cert, curl, err
***REMOVED***

// FetchCert retrieves already issued certificate from the given url, in DER format.
// It retries the request until the certificate is successfully retrieved,
// context is cancelled by the caller or an error response is received.
//
// The returned value will also contain the CA (issuer) certificate if the bundle argument is true.
//
// FetchCert returns an error if the CA's response or chain was unreasonably large.
// Callers are encouraged to parse the returned value to ensure the certificate is valid
// and has expected features.
func (c *Client) FetchCert(ctx context.Context, url string, bundle bool) ([][]byte, error) ***REMOVED***
	for ***REMOVED***
		res, err := c.get(ctx, url)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK ***REMOVED***
			return c.responseCert(ctx, res, bundle)
		***REMOVED***
		if res.StatusCode > 299 ***REMOVED***
			return nil, responseError(res)
		***REMOVED***
		d := retryAfter(res.Header.Get("Retry-After"), 3*time.Second)
		select ***REMOVED***
		case <-time.After(d):
			// retry
		case <-ctx.Done():
			return nil, ctx.Err()
		***REMOVED***
	***REMOVED***
***REMOVED***

// RevokeCert revokes a previously issued certificate cert, provided in DER format.
//
// The key argument, used to sign the request, must be authorized
// to revoke the certificate. It's up to the CA to decide which keys are authorized.
// For instance, the key pair of the certificate may be authorized.
// If the key is nil, c.Key is used instead.
func (c *Client) RevokeCert(ctx context.Context, key crypto.Signer, cert []byte, reason CRLReasonCode) error ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	body := &struct ***REMOVED***
		Resource string `json:"resource"`
		Cert     string `json:"certificate"`
		Reason   int    `json:"reason"`
	***REMOVED******REMOVED***
		Resource: "revoke-cert",
		Cert:     base64.RawURLEncoding.EncodeToString(cert),
		Reason:   int(reason),
	***REMOVED***
	if key == nil ***REMOVED***
		key = c.Key
	***REMOVED***
	res, err := c.retryPostJWS(ctx, key, c.dir.RevokeURL, body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK ***REMOVED***
		return responseError(res)
	***REMOVED***
	return nil
***REMOVED***

// AcceptTOS always returns true to indicate the acceptance of a CA's Terms of Service
// during account registration. See Register method of Client for more details.
func AcceptTOS(tosURL string) bool ***REMOVED*** return true ***REMOVED***

// Register creates a new account registration by following the "new-reg" flow.
// It returns the registered account. The account is not modified.
//
// The registration may require the caller to agree to the CA's Terms of Service (TOS).
// If so, and the account has not indicated the acceptance of the terms (see Account for details),
// Register calls prompt with a TOS URL provided by the CA. Prompt should report
// whether the caller agrees to the terms. To always accept the terms, the caller can use AcceptTOS.
func (c *Client) Register(ctx context.Context, a *Account, prompt func(tosURL string) bool) (*Account, error) ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var err error
	if a, err = c.doReg(ctx, c.dir.RegURL, "new-reg", a); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var accept bool
	if a.CurrentTerms != "" && a.CurrentTerms != a.AgreedTerms ***REMOVED***
		accept = prompt(a.CurrentTerms)
	***REMOVED***
	if accept ***REMOVED***
		a.AgreedTerms = a.CurrentTerms
		a, err = c.UpdateReg(ctx, a)
	***REMOVED***
	return a, err
***REMOVED***

// GetReg retrieves an existing registration.
// The url argument is an Account URI.
func (c *Client) GetReg(ctx context.Context, url string) (*Account, error) ***REMOVED***
	a, err := c.doReg(ctx, url, "reg", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	a.URI = url
	return a, nil
***REMOVED***

// UpdateReg updates an existing registration.
// It returns an updated account copy. The provided account is not modified.
func (c *Client) UpdateReg(ctx context.Context, a *Account) (*Account, error) ***REMOVED***
	uri := a.URI
	a, err := c.doReg(ctx, uri, "reg", a)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	a.URI = uri
	return a, nil
***REMOVED***

// Authorize performs the initial step in an authorization flow.
// The caller will then need to choose from and perform a set of returned
// challenges using c.Accept in order to successfully complete authorization.
//
// If an authorization has been previously granted, the CA may return
// a valid authorization (Authorization.Status is StatusValid). If so, the caller
// need not fulfill any challenge and can proceed to requesting a certificate.
func (c *Client) Authorize(ctx context.Context, domain string) (*Authorization, error) ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	type authzID struct ***REMOVED***
		Type  string `json:"type"`
		Value string `json:"value"`
	***REMOVED***
	req := struct ***REMOVED***
		Resource   string  `json:"resource"`
		Identifier authzID `json:"identifier"`
	***REMOVED******REMOVED***
		Resource:   "new-authz",
		Identifier: authzID***REMOVED***Type: "dns", Value: domain***REMOVED***,
	***REMOVED***
	res, err := c.retryPostJWS(ctx, c.Key, c.dir.AuthzURL, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated ***REMOVED***
		return nil, responseError(res)
	***REMOVED***

	var v wireAuthz
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: invalid response: %v", err)
	***REMOVED***
	if v.Status != StatusPending && v.Status != StatusValid ***REMOVED***
		return nil, fmt.Errorf("acme: unexpected status: %s", v.Status)
	***REMOVED***
	return v.authorization(res.Header.Get("Location")), nil
***REMOVED***

// GetAuthorization retrieves an authorization identified by the given URL.
//
// If a caller needs to poll an authorization until its status is final,
// see the WaitAuthorization method.
func (c *Client) GetAuthorization(ctx context.Context, url string) (*Authorization, error) ***REMOVED***
	res, err := c.get(ctx, url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted ***REMOVED***
		return nil, responseError(res)
	***REMOVED***
	var v wireAuthz
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: invalid response: %v", err)
	***REMOVED***
	return v.authorization(url), nil
***REMOVED***

// RevokeAuthorization relinquishes an existing authorization identified
// by the given URL.
// The url argument is an Authorization.URI value.
//
// If successful, the caller will be required to obtain a new authorization
// using the Authorize method before being able to request a new certificate
// for the domain associated with the authorization.
//
// It does not revoke existing certificates.
func (c *Client) RevokeAuthorization(ctx context.Context, url string) error ***REMOVED***
	req := struct ***REMOVED***
		Resource string `json:"resource"`
		Status   string `json:"status"`
		Delete   bool   `json:"delete"`
	***REMOVED******REMOVED***
		Resource: "authz",
		Status:   "deactivated",
		Delete:   true,
	***REMOVED***
	res, err := c.retryPostJWS(ctx, c.Key, url, req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK ***REMOVED***
		return responseError(res)
	***REMOVED***
	return nil
***REMOVED***

// WaitAuthorization polls an authorization at the given URL
// until it is in one of the final states, StatusValid or StatusInvalid,
// the ACME CA responded with a 4xx error code, or the context is done.
//
// It returns a non-nil Authorization only if its Status is StatusValid.
// In all other cases WaitAuthorization returns an error.
// If the Status is StatusInvalid, the returned error is of type *AuthorizationError.
func (c *Client) WaitAuthorization(ctx context.Context, url string) (*Authorization, error) ***REMOVED***
	sleep := sleeper(ctx)
	for ***REMOVED***
		res, err := c.get(ctx, url)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if res.StatusCode >= 400 && res.StatusCode <= 499 ***REMOVED***
			// Non-retriable error. For instance, Let's Encrypt may return 404 Not Found
			// when requesting an expired authorization.
			defer res.Body.Close()
			return nil, responseError(res)
		***REMOVED***

		retry := res.Header.Get("Retry-After")
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted ***REMOVED***
			res.Body.Close()
			if err := sleep(retry, 1); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			continue
		***REMOVED***
		var raw wireAuthz
		err = json.NewDecoder(res.Body).Decode(&raw)
		res.Body.Close()
		if err != nil ***REMOVED***
			if err := sleep(retry, 0); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			continue
		***REMOVED***
		if raw.Status == StatusValid ***REMOVED***
			return raw.authorization(url), nil
		***REMOVED***
		if raw.Status == StatusInvalid ***REMOVED***
			return nil, raw.error(url)
		***REMOVED***
		if err := sleep(retry, 0); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***

// GetChallenge retrieves the current status of an challenge.
//
// A client typically polls a challenge status using this method.
func (c *Client) GetChallenge(ctx context.Context, url string) (*Challenge, error) ***REMOVED***
	res, err := c.get(ctx, url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted ***REMOVED***
		return nil, responseError(res)
	***REMOVED***
	v := wireChallenge***REMOVED***URI: url***REMOVED***
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: invalid response: %v", err)
	***REMOVED***
	return v.challenge(), nil
***REMOVED***

// Accept informs the server that the client accepts one of its challenges
// previously obtained with c.Authorize.
//
// The server will then perform the validation asynchronously.
func (c *Client) Accept(ctx context.Context, chal *Challenge) (*Challenge, error) ***REMOVED***
	auth, err := keyAuth(c.Key.Public(), chal.Token)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req := struct ***REMOVED***
		Resource string `json:"resource"`
		Type     string `json:"type"`
		Auth     string `json:"keyAuthorization"`
	***REMOVED******REMOVED***
		Resource: "challenge",
		Type:     chal.Type,
		Auth:     auth,
	***REMOVED***
	res, err := c.retryPostJWS(ctx, c.Key, chal.URI, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	// Note: the protocol specifies 200 as the expected response code, but
	// letsencrypt seems to be returning 202.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted ***REMOVED***
		return nil, responseError(res)
	***REMOVED***

	var v wireChallenge
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: invalid response: %v", err)
	***REMOVED***
	return v.challenge(), nil
***REMOVED***

// DNS01ChallengeRecord returns a DNS record value for a dns-01 challenge response.
// A TXT record containing the returned value must be provisioned under
// "_acme-challenge" name of the domain being validated.
//
// The token argument is a Challenge.Token value.
func (c *Client) DNS01ChallengeRecord(token string) (string, error) ***REMOVED***
	ka, err := keyAuth(c.Key.Public(), token)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	b := sha256.Sum256([]byte(ka))
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
***REMOVED***

// HTTP01ChallengeResponse returns the response for an http-01 challenge.
// Servers should respond with the value to HTTP requests at the URL path
// provided by HTTP01ChallengePath to validate the challenge and prove control
// over a domain name.
//
// The token argument is a Challenge.Token value.
func (c *Client) HTTP01ChallengeResponse(token string) (string, error) ***REMOVED***
	return keyAuth(c.Key.Public(), token)
***REMOVED***

// HTTP01ChallengePath returns the URL path at which the response for an http-01 challenge
// should be provided by the servers.
// The response value can be obtained with HTTP01ChallengeResponse.
//
// The token argument is a Challenge.Token value.
func (c *Client) HTTP01ChallengePath(token string) string ***REMOVED***
	return "/.well-known/acme-challenge/" + token
***REMOVED***

// TLSSNI01ChallengeCert creates a certificate for TLS-SNI-01 challenge response.
// Servers can present the certificate to validate the challenge and prove control
// over a domain name.
//
// The implementation is incomplete in that the returned value is a single certificate,
// computed only for Z0 of the key authorization. ACME CAs are expected to update
// their implementations to use the newer version, TLS-SNI-02.
// For more details on TLS-SNI-01 see https://tools.ietf.org/html/draft-ietf-acme-acme-01#section-7.3.
//
// The token argument is a Challenge.Token value.
// If a WithKey option is provided, its private part signs the returned cert,
// and the public part is used to specify the signee.
// If no WithKey option is provided, a new ECDSA key is generated using P-256 curve.
//
// The returned certificate is valid for the next 24 hours and must be presented only when
// the server name of the client hello matches exactly the returned name value.
func (c *Client) TLSSNI01ChallengeCert(token string, opt ...CertOption) (cert tls.Certificate, name string, err error) ***REMOVED***
	ka, err := keyAuth(c.Key.Public(), token)
	if err != nil ***REMOVED***
		return tls.Certificate***REMOVED******REMOVED***, "", err
	***REMOVED***
	b := sha256.Sum256([]byte(ka))
	h := hex.EncodeToString(b[:])
	name = fmt.Sprintf("%s.%s.acme.invalid", h[:32], h[32:])
	cert, err = tlsChallengeCert([]string***REMOVED***name***REMOVED***, opt)
	if err != nil ***REMOVED***
		return tls.Certificate***REMOVED******REMOVED***, "", err
	***REMOVED***
	return cert, name, nil
***REMOVED***

// TLSSNI02ChallengeCert creates a certificate for TLS-SNI-02 challenge response.
// Servers can present the certificate to validate the challenge and prove control
// over a domain name. For more details on TLS-SNI-02 see
// https://tools.ietf.org/html/draft-ietf-acme-acme-03#section-7.3.
//
// The token argument is a Challenge.Token value.
// If a WithKey option is provided, its private part signs the returned cert,
// and the public part is used to specify the signee.
// If no WithKey option is provided, a new ECDSA key is generated using P-256 curve.
//
// The returned certificate is valid for the next 24 hours and must be presented only when
// the server name in the client hello matches exactly the returned name value.
func (c *Client) TLSSNI02ChallengeCert(token string, opt ...CertOption) (cert tls.Certificate, name string, err error) ***REMOVED***
	b := sha256.Sum256([]byte(token))
	h := hex.EncodeToString(b[:])
	sanA := fmt.Sprintf("%s.%s.token.acme.invalid", h[:32], h[32:])

	ka, err := keyAuth(c.Key.Public(), token)
	if err != nil ***REMOVED***
		return tls.Certificate***REMOVED******REMOVED***, "", err
	***REMOVED***
	b = sha256.Sum256([]byte(ka))
	h = hex.EncodeToString(b[:])
	sanB := fmt.Sprintf("%s.%s.ka.acme.invalid", h[:32], h[32:])

	cert, err = tlsChallengeCert([]string***REMOVED***sanA, sanB***REMOVED***, opt)
	if err != nil ***REMOVED***
		return tls.Certificate***REMOVED******REMOVED***, "", err
	***REMOVED***
	return cert, sanA, nil
***REMOVED***

// doReg sends all types of registration requests.
// The type of request is identified by typ argument, which is a "resource"
// in the ACME spec terms.
//
// A non-nil acct argument indicates whether the intention is to mutate data
// of the Account. Only Contact and Agreement of its fields are used
// in such cases.
func (c *Client) doReg(ctx context.Context, url string, typ string, acct *Account) (*Account, error) ***REMOVED***
	req := struct ***REMOVED***
		Resource  string   `json:"resource"`
		Contact   []string `json:"contact,omitempty"`
		Agreement string   `json:"agreement,omitempty"`
	***REMOVED******REMOVED***
		Resource: typ,
	***REMOVED***
	if acct != nil ***REMOVED***
		req.Contact = acct.Contact
		req.Agreement = acct.AgreedTerms
	***REMOVED***
	res, err := c.retryPostJWS(ctx, c.Key, url, req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 ***REMOVED***
		return nil, responseError(res)
	***REMOVED***

	var v struct ***REMOVED***
		Contact        []string
		Agreement      string
		Authorizations string
		Certificates   string
	***REMOVED***
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: invalid response: %v", err)
	***REMOVED***
	var tos string
	if v := linkHeader(res.Header, "terms-of-service"); len(v) > 0 ***REMOVED***
		tos = v[0]
	***REMOVED***
	var authz string
	if v := linkHeader(res.Header, "next"); len(v) > 0 ***REMOVED***
		authz = v[0]
	***REMOVED***
	return &Account***REMOVED***
		URI:            res.Header.Get("Location"),
		Contact:        v.Contact,
		AgreedTerms:    v.Agreement,
		CurrentTerms:   tos,
		Authz:          authz,
		Authorizations: v.Authorizations,
		Certificates:   v.Certificates,
	***REMOVED***, nil
***REMOVED***

// retryPostJWS will retry calls to postJWS if there is a badNonce error,
// clearing the stored nonces after each error.
// If the response was 4XX-5XX, then responseError is called on the body,
// the body is closed, and the error returned.
func (c *Client) retryPostJWS(ctx context.Context, key crypto.Signer, url string, body interface***REMOVED******REMOVED***) (*http.Response, error) ***REMOVED***
	sleep := sleeper(ctx)
	for ***REMOVED***
		res, err := c.postJWS(ctx, key, url, body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// handle errors 4XX-5XX with responseError
		if res.StatusCode >= 400 && res.StatusCode <= 599 ***REMOVED***
			err := responseError(res)
			res.Body.Close()
			// according to spec badNonce is urn:ietf:params:acme:error:badNonce
			// however, acme servers in the wild return their version of the error
			// https://tools.ietf.org/html/draft-ietf-acme-acme-02#section-5.4
			if ae, ok := err.(*Error); ok && strings.HasSuffix(strings.ToLower(ae.ProblemType), ":badnonce") ***REMOVED***
				// clear any nonces that we might've stored that might now be
				// considered bad
				c.clearNonces()
				retry := res.Header.Get("Retry-After")
				if err := sleep(retry, 1); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		return res, nil
	***REMOVED***
***REMOVED***

// postJWS signs the body with the given key and POSTs it to the provided url.
// The body argument must be JSON-serializable.
func (c *Client) postJWS(ctx context.Context, key crypto.Signer, url string, body interface***REMOVED******REMOVED***) (*http.Response, error) ***REMOVED***
	nonce, err := c.popNonce(ctx, url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	b, err := jwsEncodeJSON(body, key, nonce)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	res, err := c.post(ctx, url, "application/jose+json", bytes.NewReader(b))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.addNonce(res.Header)
	return res, nil
***REMOVED***

// popNonce returns a nonce value previously stored with c.addNonce
// or fetches a fresh one from the given URL.
func (c *Client) popNonce(ctx context.Context, url string) (string, error) ***REMOVED***
	c.noncesMu.Lock()
	defer c.noncesMu.Unlock()
	if len(c.nonces) == 0 ***REMOVED***
		return c.fetchNonce(ctx, url)
	***REMOVED***
	var nonce string
	for nonce = range c.nonces ***REMOVED***
		delete(c.nonces, nonce)
		break
	***REMOVED***
	return nonce, nil
***REMOVED***

// clearNonces clears any stored nonces
func (c *Client) clearNonces() ***REMOVED***
	c.noncesMu.Lock()
	defer c.noncesMu.Unlock()
	c.nonces = make(map[string]struct***REMOVED******REMOVED***)
***REMOVED***

// addNonce stores a nonce value found in h (if any) for future use.
func (c *Client) addNonce(h http.Header) ***REMOVED***
	v := nonceFromHeader(h)
	if v == "" ***REMOVED***
		return
	***REMOVED***
	c.noncesMu.Lock()
	defer c.noncesMu.Unlock()
	if len(c.nonces) >= maxNonces ***REMOVED***
		return
	***REMOVED***
	if c.nonces == nil ***REMOVED***
		c.nonces = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	c.nonces[v] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

func (c *Client) httpClient() *http.Client ***REMOVED***
	if c.HTTPClient != nil ***REMOVED***
		return c.HTTPClient
	***REMOVED***
	return http.DefaultClient
***REMOVED***

func (c *Client) get(ctx context.Context, urlStr string) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c.do(ctx, req)
***REMOVED***

func (c *Client) head(ctx context.Context, urlStr string) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("HEAD", urlStr, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c.do(ctx, req)
***REMOVED***

func (c *Client) post(ctx context.Context, urlStr, contentType string, body io.Reader) (*http.Response, error) ***REMOVED***
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Set("Content-Type", contentType)
	return c.do(ctx, req)
***REMOVED***

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) ***REMOVED***
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

func (c *Client) fetchNonce(ctx context.Context, url string) (string, error) ***REMOVED***
	resp, err := c.head(ctx, url)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer resp.Body.Close()
	nonce := nonceFromHeader(resp.Header)
	if nonce == "" ***REMOVED***
		if resp.StatusCode > 299 ***REMOVED***
			return "", responseError(resp)
		***REMOVED***
		return "", errors.New("acme: nonce not found")
	***REMOVED***
	return nonce, nil
***REMOVED***

func nonceFromHeader(h http.Header) string ***REMOVED***
	return h.Get("Replay-Nonce")
***REMOVED***

func (c *Client) responseCert(ctx context.Context, res *http.Response, bundle bool) ([][]byte, error) ***REMOVED***
	b, err := ioutil.ReadAll(io.LimitReader(res.Body, maxCertSize+1))
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: response stream: %v", err)
	***REMOVED***
	if len(b) > maxCertSize ***REMOVED***
		return nil, errors.New("acme: certificate is too big")
	***REMOVED***
	cert := [][]byte***REMOVED***b***REMOVED***
	if !bundle ***REMOVED***
		return cert, nil
	***REMOVED***

	// Append CA chain cert(s).
	// At least one is required according to the spec:
	// https://tools.ietf.org/html/draft-ietf-acme-acme-03#section-6.3.1
	up := linkHeader(res.Header, "up")
	if len(up) == 0 ***REMOVED***
		return nil, errors.New("acme: rel=up link not found")
	***REMOVED***
	if len(up) > maxChainLen ***REMOVED***
		return nil, errors.New("acme: rel=up link is too large")
	***REMOVED***
	for _, url := range up ***REMOVED***
		cc, err := c.chainCert(ctx, url, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cert = append(cert, cc...)
	***REMOVED***
	return cert, nil
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

// chainCert fetches CA certificate chain recursively by following "up" links.
// Each recursive call increments the depth by 1, resulting in an error
// if the recursion level reaches maxChainLen.
//
// First chainCert call starts with depth of 0.
func (c *Client) chainCert(ctx context.Context, url string, depth int) ([][]byte, error) ***REMOVED***
	if depth >= maxChainLen ***REMOVED***
		return nil, errors.New("acme: certificate chain is too deep")
	***REMOVED***

	res, err := c.get(ctx, url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK ***REMOVED***
		return nil, responseError(res)
	***REMOVED***
	b, err := ioutil.ReadAll(io.LimitReader(res.Body, maxCertSize+1))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(b) > maxCertSize ***REMOVED***
		return nil, errors.New("acme: certificate is too big")
	***REMOVED***
	chain := [][]byte***REMOVED***b***REMOVED***

	uplink := linkHeader(res.Header, "up")
	if len(uplink) > maxChainLen ***REMOVED***
		return nil, errors.New("acme: certificate chain is too large")
	***REMOVED***
	for _, up := range uplink ***REMOVED***
		cc, err := c.chainCert(ctx, up, depth+1)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		chain = append(chain, cc...)
	***REMOVED***

	return chain, nil
***REMOVED***

// linkHeader returns URI-Reference values of all Link headers
// with relation-type rel.
// See https://tools.ietf.org/html/rfc5988#section-5 for details.
func linkHeader(h http.Header, rel string) []string ***REMOVED***
	var links []string
	for _, v := range h["Link"] ***REMOVED***
		parts := strings.Split(v, ";")
		for _, p := range parts ***REMOVED***
			p = strings.TrimSpace(p)
			if !strings.HasPrefix(p, "rel=") ***REMOVED***
				continue
			***REMOVED***
			if v := strings.Trim(p[4:], `"`); v == rel ***REMOVED***
				links = append(links, strings.Trim(parts[0], "<>"))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return links
***REMOVED***

// sleeper returns a function that accepts the Retry-After HTTP header value
// and an increment that's used with backoff to increasingly sleep on
// consecutive calls until the context is done. If the Retry-After header
// cannot be parsed, then backoff is used with a maximum sleep time of 10
// seconds.
func sleeper(ctx context.Context) func(ra string, inc int) error ***REMOVED***
	var count int
	return func(ra string, inc int) error ***REMOVED***
		count += inc
		d := backoff(count, 10*time.Second)
		d = retryAfter(ra, d)
		wakeup := time.NewTimer(d)
		defer wakeup.Stop()
		select ***REMOVED***
		case <-ctx.Done():
			return ctx.Err()
		case <-wakeup.C:
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// retryAfter parses a Retry-After HTTP header value,
// trying to convert v into an int (seconds) or use http.ParseTime otherwise.
// It returns d if v cannot be parsed.
func retryAfter(v string, d time.Duration) time.Duration ***REMOVED***
	if i, err := strconv.Atoi(v); err == nil ***REMOVED***
		return time.Duration(i) * time.Second
	***REMOVED***
	t, err := http.ParseTime(v)
	if err != nil ***REMOVED***
		return d
	***REMOVED***
	return t.Sub(timeNow())
***REMOVED***

// backoff computes a duration after which an n+1 retry iteration should occur
// using truncated exponential backoff algorithm.
//
// The n argument is always bounded between 0 and 30.
// The max argument defines upper bound for the returned value.
func backoff(n int, max time.Duration) time.Duration ***REMOVED***
	if n < 0 ***REMOVED***
		n = 0
	***REMOVED***
	if n > 30 ***REMOVED***
		n = 30
	***REMOVED***
	var d time.Duration
	if x, err := rand.Int(rand.Reader, big.NewInt(1000)); err == nil ***REMOVED***
		d = time.Duration(x.Int64()) * time.Millisecond
	***REMOVED***
	d += time.Duration(1<<uint(n)) * time.Second
	if d > max ***REMOVED***
		return max
	***REMOVED***
	return d
***REMOVED***

// keyAuth generates a key authorization string for a given token.
func keyAuth(pub crypto.PublicKey, token string) (string, error) ***REMOVED***
	th, err := JWKThumbprint(pub)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return fmt.Sprintf("%s.%s", token, th), nil
***REMOVED***

// tlsChallengeCert creates a temporary certificate for TLS-SNI challenges
// with the given SANs and auto-generated public/private key pair.
// The Subject Common Name is set to the first SAN to aid debugging.
// To create a cert with a custom key pair, specify WithKey option.
func tlsChallengeCert(san []string, opt []CertOption) (tls.Certificate, error) ***REMOVED***
	var (
		key  crypto.Signer
		tmpl *x509.Certificate
	)
	for _, o := range opt ***REMOVED***
		switch o := o.(type) ***REMOVED***
		case *certOptKey:
			if key != nil ***REMOVED***
				return tls.Certificate***REMOVED******REMOVED***, errors.New("acme: duplicate key option")
			***REMOVED***
			key = o.key
		case *certOptTemplate:
			var t = *(*x509.Certificate)(o) // shallow copy is ok
			tmpl = &t
		default:
			// package's fault, if we let this happen:
			panic(fmt.Sprintf("unsupported option type %T", o))
		***REMOVED***
	***REMOVED***
	if key == nil ***REMOVED***
		var err error
		if key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil ***REMOVED***
			return tls.Certificate***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	if tmpl == nil ***REMOVED***
		tmpl = &x509.Certificate***REMOVED***
			SerialNumber:          big.NewInt(1),
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(24 * time.Hour),
			BasicConstraintsValid: true,
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage***REMOVED***x509.ExtKeyUsageServerAuth***REMOVED***,
		***REMOVED***
	***REMOVED***
	tmpl.DNSNames = san
	if len(san) > 0 ***REMOVED***
		tmpl.Subject.CommonName = san[0]
	***REMOVED***

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil ***REMOVED***
		return tls.Certificate***REMOVED******REMOVED***, err
	***REMOVED***
	return tls.Certificate***REMOVED***
		Certificate: [][]byte***REMOVED***der***REMOVED***,
		PrivateKey:  key,
	***REMOVED***, nil
***REMOVED***

// encodePEM returns b encoded as PEM with block of type typ.
func encodePEM(typ string, b []byte) []byte ***REMOVED***
	pb := &pem.Block***REMOVED***Type: typ, Bytes: b***REMOVED***
	return pem.EncodeToMemory(pb)
***REMOVED***

// timeNow is useful for testing for fixed current time.
var timeNow = time.Now
