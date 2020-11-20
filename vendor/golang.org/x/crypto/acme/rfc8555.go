// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acme

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// DeactivateReg permanently disables an existing account associated with c.Key.
// A deactivated account can no longer request certificate issuance or access
// resources related to the account, such as orders or authorizations.
//
// It only works with CAs implementing RFC 8555.
func (c *Client) DeactivateReg(ctx context.Context) error ***REMOVED***
	url := string(c.accountKID(ctx))
	if url == "" ***REMOVED***
		return ErrNoAccount
	***REMOVED***
	req := json.RawMessage(`***REMOVED***"status": "deactivated"***REMOVED***`)
	res, err := c.post(ctx, nil, url, req, wantStatus(http.StatusOK))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	res.Body.Close()
	return nil
***REMOVED***

// registerRFC is quivalent to c.Register but for CAs implementing RFC 8555.
// It expects c.Discover to have already been called.
// TODO: Implement externalAccountBinding.
func (c *Client) registerRFC(ctx context.Context, acct *Account, prompt func(tosURL string) bool) (*Account, error) ***REMOVED***
	c.cacheMu.Lock() // guard c.kid access
	defer c.cacheMu.Unlock()

	req := struct ***REMOVED***
		TermsAgreed bool     `json:"termsOfServiceAgreed,omitempty"`
		Contact     []string `json:"contact,omitempty"`
	***REMOVED******REMOVED***
		Contact: acct.Contact,
	***REMOVED***
	if c.dir.Terms != "" ***REMOVED***
		req.TermsAgreed = prompt(c.dir.Terms)
	***REMOVED***
	res, err := c.post(ctx, c.Key, c.dir.RegURL, req, wantStatus(
		http.StatusOK,      // account with this key already registered
		http.StatusCreated, // new account created
	))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer res.Body.Close()
	a, err := responseAccount(res)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Cache Account URL even if we return an error to the caller.
	// It is by all means a valid and usable "kid" value for future requests.
	c.kid = keyID(a.URI)
	if res.StatusCode == http.StatusOK ***REMOVED***
		return nil, ErrAccountAlreadyExists
	***REMOVED***
	return a, nil
***REMOVED***

// updateGegRFC is equivalent to c.UpdateReg but for CAs implementing RFC 8555.
// It expects c.Discover to have already been called.
func (c *Client) updateRegRFC(ctx context.Context, a *Account) (*Account, error) ***REMOVED***
	url := string(c.accountKID(ctx))
	if url == "" ***REMOVED***
		return nil, ErrNoAccount
	***REMOVED***
	req := struct ***REMOVED***
		Contact []string `json:"contact,omitempty"`
	***REMOVED******REMOVED***
		Contact: a.Contact,
	***REMOVED***
	res, err := c.post(ctx, nil, url, req, wantStatus(http.StatusOK))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	return responseAccount(res)
***REMOVED***

// getGegRFC is equivalent to c.GetReg but for CAs implementing RFC 8555.
// It expects c.Discover to have already been called.
func (c *Client) getRegRFC(ctx context.Context) (*Account, error) ***REMOVED***
	req := json.RawMessage(`***REMOVED***"onlyReturnExisting": true***REMOVED***`)
	res, err := c.post(ctx, c.Key, c.dir.RegURL, req, wantStatus(http.StatusOK))
	if e, ok := err.(*Error); ok && e.ProblemType == "urn:ietf:params:acme:error:accountDoesNotExist" ***REMOVED***
		return nil, ErrNoAccount
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer res.Body.Close()
	return responseAccount(res)
***REMOVED***

func responseAccount(res *http.Response) (*Account, error) ***REMOVED***
	var v struct ***REMOVED***
		Status  string
		Contact []string
		Orders  string
	***REMOVED***
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: invalid account response: %v", err)
	***REMOVED***
	return &Account***REMOVED***
		URI:       res.Header.Get("Location"),
		Status:    v.Status,
		Contact:   v.Contact,
		OrdersURL: v.Orders,
	***REMOVED***, nil
***REMOVED***

// AuthorizeOrder initiates the order-based application for certificate issuance,
// as opposed to pre-authorization in Authorize.
// It is only supported by CAs implementing RFC 8555.
//
// The caller then needs to fetch each authorization with GetAuthorization,
// identify those with StatusPending status and fulfill a challenge using Accept.
// Once all authorizations are satisfied, the caller will typically want to poll
// order status using WaitOrder until it's in StatusReady state.
// To finalize the order and obtain a certificate, the caller submits a CSR with CreateOrderCert.
func (c *Client) AuthorizeOrder(ctx context.Context, id []AuthzID, opt ...OrderOption) (*Order, error) ***REMOVED***
	dir, err := c.Discover(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req := struct ***REMOVED***
		Identifiers []wireAuthzID `json:"identifiers"`
		NotBefore   string        `json:"notBefore,omitempty"`
		NotAfter    string        `json:"notAfter,omitempty"`
	***REMOVED******REMOVED******REMOVED***
	for _, v := range id ***REMOVED***
		req.Identifiers = append(req.Identifiers, wireAuthzID***REMOVED***
			Type:  v.Type,
			Value: v.Value,
		***REMOVED***)
	***REMOVED***
	for _, o := range opt ***REMOVED***
		switch o := o.(type) ***REMOVED***
		case orderNotBeforeOpt:
			req.NotBefore = time.Time(o).Format(time.RFC3339)
		case orderNotAfterOpt:
			req.NotAfter = time.Time(o).Format(time.RFC3339)
		default:
			// Package's fault if we let this happen.
			panic(fmt.Sprintf("unsupported order option type %T", o))
		***REMOVED***
	***REMOVED***

	res, err := c.post(ctx, nil, dir.OrderURL, req, wantStatus(http.StatusCreated))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	return responseOrder(res)
***REMOVED***

// GetOrder retrives an order identified by the given URL.
// For orders created with AuthorizeOrder, the url value is Order.URI.
//
// If a caller needs to poll an order until its status is final,
// see the WaitOrder method.
func (c *Client) GetOrder(ctx context.Context, url string) (*Order, error) ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	res, err := c.postAsGet(ctx, url, wantStatus(http.StatusOK))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()
	return responseOrder(res)
***REMOVED***

// WaitOrder polls an order from the given URL until it is in one of the final states,
// StatusReady, StatusValid or StatusInvalid, the CA responded with a non-retryable error
// or the context is done.
//
// It returns a non-nil Order only if its Status is StatusReady or StatusValid.
// In all other cases WaitOrder returns an error.
// If the Status is StatusInvalid, the returned error is of type *OrderError.
func (c *Client) WaitOrder(ctx context.Context, url string) (*Order, error) ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for ***REMOVED***
		res, err := c.postAsGet(ctx, url, wantStatus(http.StatusOK))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		o, err := responseOrder(res)
		res.Body.Close()
		switch ***REMOVED***
		case err != nil:
			// Skip and retry.
		case o.Status == StatusInvalid:
			return nil, &OrderError***REMOVED***OrderURL: o.URI, Status: o.Status***REMOVED***
		case o.Status == StatusReady || o.Status == StatusValid:
			return o, nil
		***REMOVED***

		d := retryAfter(res.Header.Get("Retry-After"))
		if d == 0 ***REMOVED***
			// Default retry-after.
			// Same reasoning as in WaitAuthorization.
			d = time.Second
		***REMOVED***
		t := time.NewTimer(d)
		select ***REMOVED***
		case <-ctx.Done():
			t.Stop()
			return nil, ctx.Err()
		case <-t.C:
			// Retry.
		***REMOVED***
	***REMOVED***
***REMOVED***

func responseOrder(res *http.Response) (*Order, error) ***REMOVED***
	var v struct ***REMOVED***
		Status         string
		Expires        time.Time
		Identifiers    []wireAuthzID
		NotBefore      time.Time
		NotAfter       time.Time
		Error          *wireError
		Authorizations []string
		Finalize       string
		Certificate    string
	***REMOVED***
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: error reading order: %v", err)
	***REMOVED***
	o := &Order***REMOVED***
		URI:         res.Header.Get("Location"),
		Status:      v.Status,
		Expires:     v.Expires,
		NotBefore:   v.NotBefore,
		NotAfter:    v.NotAfter,
		AuthzURLs:   v.Authorizations,
		FinalizeURL: v.Finalize,
		CertURL:     v.Certificate,
	***REMOVED***
	for _, id := range v.Identifiers ***REMOVED***
		o.Identifiers = append(o.Identifiers, AuthzID***REMOVED***Type: id.Type, Value: id.Value***REMOVED***)
	***REMOVED***
	if v.Error != nil ***REMOVED***
		o.Error = v.Error.error(nil /* headers */)
	***REMOVED***
	return o, nil
***REMOVED***

// CreateOrderCert submits the CSR (Certificate Signing Request) to a CA at the specified URL.
// The URL is the FinalizeURL field of an Order created with AuthorizeOrder.
//
// If the bundle argument is true, the returned value also contain the CA (issuer)
// certificate chain. Otherwise, only a leaf certificate is returned.
// The returned URL can be used to re-fetch the certificate using FetchCert.
//
// This method is only supported by CAs implementing RFC 8555. See CreateCert for pre-RFC CAs.
//
// CreateOrderCert returns an error if the CA's response is unreasonably large.
// Callers are encouraged to parse the returned value to ensure the certificate is valid and has the expected features.
func (c *Client) CreateOrderCert(ctx context.Context, url string, csr []byte, bundle bool) (der [][]byte, certURL string, err error) ***REMOVED***
	if _, err := c.Discover(ctx); err != nil ***REMOVED*** // required by c.accountKID
		return nil, "", err
	***REMOVED***

	// RFC describes this as "finalize order" request.
	req := struct ***REMOVED***
		CSR string `json:"csr"`
	***REMOVED******REMOVED***
		CSR: base64.RawURLEncoding.EncodeToString(csr),
	***REMOVED***
	res, err := c.post(ctx, nil, url, req, wantStatus(http.StatusOK))
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	defer res.Body.Close()
	o, err := responseOrder(res)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***

	// Wait for CA to issue the cert if they haven't.
	if o.Status != StatusValid ***REMOVED***
		o, err = c.WaitOrder(ctx, o.URI)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	// The only acceptable status post finalize and WaitOrder is "valid".
	if o.Status != StatusValid ***REMOVED***
		return nil, "", &OrderError***REMOVED***OrderURL: o.URI, Status: o.Status***REMOVED***
	***REMOVED***
	crt, err := c.fetchCertRFC(ctx, o.CertURL, bundle)
	return crt, o.CertURL, err
***REMOVED***

// fetchCertRFC downloads issued certificate from the given URL.
// It expects the CA to respond with PEM-encoded certificate chain.
//
// The URL argument is the CertURL field of Order.
func (c *Client) fetchCertRFC(ctx context.Context, url string, bundle bool) ([][]byte, error) ***REMOVED***
	res, err := c.postAsGet(ctx, url, wantStatus(http.StatusOK))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()

	// Get all the bytes up to a sane maximum.
	// Account very roughly for base64 overhead.
	const max = maxCertChainSize + maxCertChainSize/33
	b, err := ioutil.ReadAll(io.LimitReader(res.Body, max+1))
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("acme: fetch cert response stream: %v", err)
	***REMOVED***
	if len(b) > max ***REMOVED***
		return nil, errors.New("acme: certificate chain is too big")
	***REMOVED***

	// Decode PEM chain.
	var chain [][]byte
	for ***REMOVED***
		var p *pem.Block
		p, b = pem.Decode(b)
		if p == nil ***REMOVED***
			break
		***REMOVED***
		if p.Type != "CERTIFICATE" ***REMOVED***
			return nil, fmt.Errorf("acme: invalid PEM cert type %q", p.Type)
		***REMOVED***

		chain = append(chain, p.Bytes)
		if !bundle ***REMOVED***
			return chain, nil
		***REMOVED***
		if len(chain) > maxChainLen ***REMOVED***
			return nil, errors.New("acme: certificate chain is too long")
		***REMOVED***
	***REMOVED***
	if len(chain) == 0 ***REMOVED***
		return nil, errors.New("acme: certificate chain is empty")
	***REMOVED***
	return chain, nil
***REMOVED***

// sends a cert revocation request in either JWK form when key is non-nil or KID form otherwise.
func (c *Client) revokeCertRFC(ctx context.Context, key crypto.Signer, cert []byte, reason CRLReasonCode) error ***REMOVED***
	req := &struct ***REMOVED***
		Cert   string `json:"certificate"`
		Reason int    `json:"reason"`
	***REMOVED******REMOVED***
		Cert:   base64.RawURLEncoding.EncodeToString(cert),
		Reason: int(reason),
	***REMOVED***
	res, err := c.post(ctx, key, c.dir.RevokeURL, req, wantStatus(http.StatusOK))
	if err != nil ***REMOVED***
		if isAlreadyRevoked(err) ***REMOVED***
			// Assume it is not an error to revoke an already revoked cert.
			return nil
		***REMOVED***
		return err
	***REMOVED***
	defer res.Body.Close()
	return nil
***REMOVED***

func isAlreadyRevoked(err error) bool ***REMOVED***
	e, ok := err.(*Error)
	return ok && e.ProblemType == "urn:ietf:params:acme:error:alreadyRevoked"
***REMOVED***
