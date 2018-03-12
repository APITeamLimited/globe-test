// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package autocert provides automatic access to certificates from Let's Encrypt
// and any other ACME-based CA.
//
// This package is a work in progress and makes no API stability promises.
package autocert

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/acme"
)

// createCertRetryAfter is how much time to wait before removing a failed state
// entry due to an unsuccessful createCert call.
// This is a variable instead of a const for testing.
// TODO: Consider making it configurable or an exp backoff?
var createCertRetryAfter = time.Minute

// pseudoRand is safe for concurrent use.
var pseudoRand *lockedMathRand

func init() ***REMOVED***
	src := mathrand.NewSource(timeNow().UnixNano())
	pseudoRand = &lockedMathRand***REMOVED***rnd: mathrand.New(src)***REMOVED***
***REMOVED***

// AcceptTOS is a Manager.Prompt function that always returns true to
// indicate acceptance of the CA's Terms of Service during account
// registration.
func AcceptTOS(tosURL string) bool ***REMOVED*** return true ***REMOVED***

// HostPolicy specifies which host names the Manager is allowed to respond to.
// It returns a non-nil error if the host should be rejected.
// The returned error is accessible via tls.Conn.Handshake and its callers.
// See Manager's HostPolicy field and GetCertificate method docs for more details.
type HostPolicy func(ctx context.Context, host string) error

// HostWhitelist returns a policy where only the specified host names are allowed.
// Only exact matches are currently supported. Subdomains, regexp or wildcard
// will not match.
func HostWhitelist(hosts ...string) HostPolicy ***REMOVED***
	whitelist := make(map[string]bool, len(hosts))
	for _, h := range hosts ***REMOVED***
		whitelist[h] = true
	***REMOVED***
	return func(_ context.Context, host string) error ***REMOVED***
		if !whitelist[host] ***REMOVED***
			return errors.New("acme/autocert: host not configured")
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// defaultHostPolicy is used when Manager.HostPolicy is not set.
func defaultHostPolicy(context.Context, string) error ***REMOVED***
	return nil
***REMOVED***

// Manager is a stateful certificate manager built on top of acme.Client.
// It obtains and refreshes certificates automatically using "tls-sni-01",
// "tls-sni-02" and "http-01" challenge types, as well as providing them
// to a TLS server via tls.Config.
//
// You must specify a cache implementation, such as DirCache,
// to reuse obtained certificates across program restarts.
// Otherwise your server is very likely to exceed the certificate
// issuer's request rate limits.
type Manager struct ***REMOVED***
	// Prompt specifies a callback function to conditionally accept a CA's Terms of Service (TOS).
	// The registration may require the caller to agree to the CA's TOS.
	// If so, Manager calls Prompt with a TOS URL provided by the CA. Prompt should report
	// whether the caller agrees to the terms.
	//
	// To always accept the terms, the callers can use AcceptTOS.
	Prompt func(tosURL string) bool

	// Cache optionally stores and retrieves previously-obtained certificates.
	// If nil, certs will only be cached for the lifetime of the Manager.
	//
	// Manager passes the Cache certificates data encoded in PEM, with private/public
	// parts combined in a single Cache.Put call, private key first.
	Cache Cache

	// HostPolicy controls which domains the Manager will attempt
	// to retrieve new certificates for. It does not affect cached certs.
	//
	// If non-nil, HostPolicy is called before requesting a new cert.
	// If nil, all hosts are currently allowed. This is not recommended,
	// as it opens a potential attack where clients connect to a server
	// by IP address and pretend to be asking for an incorrect host name.
	// Manager will attempt to obtain a certificate for that host, incorrectly,
	// eventually reaching the CA's rate limit for certificate requests
	// and making it impossible to obtain actual certificates.
	//
	// See GetCertificate for more details.
	HostPolicy HostPolicy

	// RenewBefore optionally specifies how early certificates should
	// be renewed before they expire.
	//
	// If zero, they're renewed 30 days before expiration.
	RenewBefore time.Duration

	// Client is used to perform low-level operations, such as account registration
	// and requesting new certificates.
	// If Client is nil, a zero-value acme.Client is used with acme.LetsEncryptURL
	// directory endpoint and a newly-generated ECDSA P-256 key.
	//
	// Mutating the field after the first call of GetCertificate method will have no effect.
	Client *acme.Client

	// Email optionally specifies a contact email address.
	// This is used by CAs, such as Let's Encrypt, to notify about problems
	// with issued certificates.
	//
	// If the Client's account key is already registered, Email is not used.
	Email string

	// ForceRSA makes the Manager generate certificates with 2048-bit RSA keys.
	//
	// If false, a default is used. Currently the default
	// is EC-based keys using the P-256 curve.
	ForceRSA bool

	clientMu sync.Mutex
	client   *acme.Client // initialized by acmeClient method

	stateMu sync.Mutex
	state   map[string]*certState // keyed by domain name

	// renewal tracks the set of domains currently running renewal timers.
	// It is keyed by domain name.
	renewalMu sync.Mutex
	renewal   map[string]*domainRenewal

	// tokensMu guards the rest of the fields: tryHTTP01, certTokens and httpTokens.
	tokensMu sync.RWMutex
	// tryHTTP01 indicates whether the Manager should try "http-01" challenge type
	// during the authorization flow.
	tryHTTP01 bool
	// httpTokens contains response body values for http-01 challenges
	// and is keyed by the URL path at which a challenge response is expected
	// to be provisioned.
	// The entries are stored for the duration of the authorization flow.
	httpTokens map[string][]byte
	// certTokens contains temporary certificates for tls-sni challenges
	// and is keyed by token domain name, which matches server name of ClientHello.
	// Keys always have ".acme.invalid" suffix.
	// The entries are stored for the duration of the authorization flow.
	certTokens map[string]*tls.Certificate
***REMOVED***

// GetCertificate implements the tls.Config.GetCertificate hook.
// It provides a TLS certificate for hello.ServerName host, including answering
// *.acme.invalid (TLS-SNI) challenges. All other fields of hello are ignored.
//
// If m.HostPolicy is non-nil, GetCertificate calls the policy before requesting
// a new cert. A non-nil error returned from m.HostPolicy halts TLS negotiation.
// The error is propagated back to the caller of GetCertificate and is user-visible.
// This does not affect cached certs. See HostPolicy field description for more details.
func (m *Manager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) ***REMOVED***
	if m.Prompt == nil ***REMOVED***
		return nil, errors.New("acme/autocert: Manager.Prompt not set")
	***REMOVED***

	name := hello.ServerName
	if name == "" ***REMOVED***
		return nil, errors.New("acme/autocert: missing server name")
	***REMOVED***
	if !strings.Contains(strings.Trim(name, "."), ".") ***REMOVED***
		return nil, errors.New("acme/autocert: server name component count invalid")
	***REMOVED***
	if strings.ContainsAny(name, `/\`) ***REMOVED***
		return nil, errors.New("acme/autocert: server name contains invalid character")
	***REMOVED***

	// In the worst-case scenario, the timeout needs to account for caching, host policy,
	// domain ownership verification and certificate issuance.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// check whether this is a token cert requested for TLS-SNI challenge
	if strings.HasSuffix(name, ".acme.invalid") ***REMOVED***
		m.tokensMu.RLock()
		defer m.tokensMu.RUnlock()
		if cert := m.certTokens[name]; cert != nil ***REMOVED***
			return cert, nil
		***REMOVED***
		if cert, err := m.cacheGet(ctx, name); err == nil ***REMOVED***
			return cert, nil
		***REMOVED***
		// TODO: cache error results?
		return nil, fmt.Errorf("acme/autocert: no token cert for %q", name)
	***REMOVED***

	// regular domain
	name = strings.TrimSuffix(name, ".") // golang.org/issue/18114
	cert, err := m.cert(ctx, name)
	if err == nil ***REMOVED***
		return cert, nil
	***REMOVED***
	if err != ErrCacheMiss ***REMOVED***
		return nil, err
	***REMOVED***

	// first-time
	if err := m.hostPolicy()(ctx, name); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cert, err = m.createCert(ctx, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.cachePut(ctx, name, cert)
	return cert, nil
***REMOVED***

// HTTPHandler configures the Manager to provision ACME "http-01" challenge responses.
// It returns an http.Handler that responds to the challenges and must be
// running on port 80. If it receives a request that is not an ACME challenge,
// it delegates the request to the optional fallback handler.
//
// If fallback is nil, the returned handler redirects all GET and HEAD requests
// to the default TLS port 443 with 302 Found status code, preserving the original
// request path and query. It responds with 400 Bad Request to all other HTTP methods.
// The fallback is not protected by the optional HostPolicy.
//
// Because the fallback handler is run with unencrypted port 80 requests,
// the fallback should not serve TLS-only requests.
//
// If HTTPHandler is never called, the Manager will only use TLS SNI
// challenges for domain verification.
func (m *Manager) HTTPHandler(fallback http.Handler) http.Handler ***REMOVED***
	m.tokensMu.Lock()
	defer m.tokensMu.Unlock()
	m.tryHTTP01 = true

	if fallback == nil ***REMOVED***
		fallback = http.HandlerFunc(handleHTTPRedirect)
	***REMOVED***
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if !strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge/") ***REMOVED***
			fallback.ServeHTTP(w, r)
			return
		***REMOVED***
		// A reasonable context timeout for cache and host policy only,
		// because we don't wait for a new certificate issuance here.
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()
		if err := m.hostPolicy()(ctx, r.Host); err != nil ***REMOVED***
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		***REMOVED***
		data, err := m.httpToken(ctx, r.URL.Path)
		if err != nil ***REMOVED***
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		***REMOVED***
		w.Write(data)
	***REMOVED***)
***REMOVED***

func handleHTTPRedirect(w http.ResponseWriter, r *http.Request) ***REMOVED***
	if r.Method != "GET" && r.Method != "HEAD" ***REMOVED***
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	***REMOVED***
	target := "https://" + stripPort(r.Host) + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusFound)
***REMOVED***

func stripPort(hostport string) string ***REMOVED***
	host, _, err := net.SplitHostPort(hostport)
	if err != nil ***REMOVED***
		return hostport
	***REMOVED***
	return net.JoinHostPort(host, "443")
***REMOVED***

// cert returns an existing certificate either from m.state or cache.
// If a certificate is found in cache but not in m.state, the latter will be filled
// with the cached value.
func (m *Manager) cert(ctx context.Context, name string) (*tls.Certificate, error) ***REMOVED***
	m.stateMu.Lock()
	if s, ok := m.state[name]; ok ***REMOVED***
		m.stateMu.Unlock()
		s.RLock()
		defer s.RUnlock()
		return s.tlscert()
	***REMOVED***
	defer m.stateMu.Unlock()
	cert, err := m.cacheGet(ctx, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	signer, ok := cert.PrivateKey.(crypto.Signer)
	if !ok ***REMOVED***
		return nil, errors.New("acme/autocert: private key cannot sign")
	***REMOVED***
	if m.state == nil ***REMOVED***
		m.state = make(map[string]*certState)
	***REMOVED***
	s := &certState***REMOVED***
		key:  signer,
		cert: cert.Certificate,
		leaf: cert.Leaf,
	***REMOVED***
	m.state[name] = s
	go m.renew(name, s.key, s.leaf.NotAfter)
	return cert, nil
***REMOVED***

// cacheGet always returns a valid certificate, or an error otherwise.
// If a cached certficate exists but is not valid, ErrCacheMiss is returned.
func (m *Manager) cacheGet(ctx context.Context, domain string) (*tls.Certificate, error) ***REMOVED***
	if m.Cache == nil ***REMOVED***
		return nil, ErrCacheMiss
	***REMOVED***
	data, err := m.Cache.Get(ctx, domain)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// private
	priv, pub := pem.Decode(data)
	if priv == nil || !strings.Contains(priv.Type, "PRIVATE") ***REMOVED***
		return nil, ErrCacheMiss
	***REMOVED***
	privKey, err := parsePrivateKey(priv.Bytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// public
	var pubDER [][]byte
	for len(pub) > 0 ***REMOVED***
		var b *pem.Block
		b, pub = pem.Decode(pub)
		if b == nil ***REMOVED***
			break
		***REMOVED***
		pubDER = append(pubDER, b.Bytes)
	***REMOVED***
	if len(pub) > 0 ***REMOVED***
		// Leftover content not consumed by pem.Decode. Corrupt. Ignore.
		return nil, ErrCacheMiss
	***REMOVED***

	// verify and create TLS cert
	leaf, err := validCert(domain, pubDER, privKey)
	if err != nil ***REMOVED***
		return nil, ErrCacheMiss
	***REMOVED***
	tlscert := &tls.Certificate***REMOVED***
		Certificate: pubDER,
		PrivateKey:  privKey,
		Leaf:        leaf,
	***REMOVED***
	return tlscert, nil
***REMOVED***

func (m *Manager) cachePut(ctx context.Context, domain string, tlscert *tls.Certificate) error ***REMOVED***
	if m.Cache == nil ***REMOVED***
		return nil
	***REMOVED***

	// contains PEM-encoded data
	var buf bytes.Buffer

	// private
	switch key := tlscert.PrivateKey.(type) ***REMOVED***
	case *ecdsa.PrivateKey:
		if err := encodeECDSAKey(&buf, key); err != nil ***REMOVED***
			return err
		***REMOVED***
	case *rsa.PrivateKey:
		b := x509.MarshalPKCS1PrivateKey(key)
		pb := &pem.Block***REMOVED***Type: "RSA PRIVATE KEY", Bytes: b***REMOVED***
		if err := pem.Encode(&buf, pb); err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		return errors.New("acme/autocert: unknown private key type")
	***REMOVED***

	// public
	for _, b := range tlscert.Certificate ***REMOVED***
		pb := &pem.Block***REMOVED***Type: "CERTIFICATE", Bytes: b***REMOVED***
		if err := pem.Encode(&buf, pb); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return m.Cache.Put(ctx, domain, buf.Bytes())
***REMOVED***

func encodeECDSAKey(w io.Writer, key *ecdsa.PrivateKey) error ***REMOVED***
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pb := &pem.Block***REMOVED***Type: "EC PRIVATE KEY", Bytes: b***REMOVED***
	return pem.Encode(w, pb)
***REMOVED***

// createCert starts the domain ownership verification and returns a certificate
// for that domain upon success.
//
// If the domain is already being verified, it waits for the existing verification to complete.
// Either way, createCert blocks for the duration of the whole process.
func (m *Manager) createCert(ctx context.Context, domain string) (*tls.Certificate, error) ***REMOVED***
	// TODO: maybe rewrite this whole piece using sync.Once
	state, err := m.certState(domain)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// state may exist if another goroutine is already working on it
	// in which case just wait for it to finish
	if !state.locked ***REMOVED***
		state.RLock()
		defer state.RUnlock()
		return state.tlscert()
	***REMOVED***

	// We are the first; state is locked.
	// Unblock the readers when domain ownership is verified
	// and we got the cert or the process failed.
	defer state.Unlock()
	state.locked = false

	der, leaf, err := m.authorizedCert(ctx, state.key, domain)
	if err != nil ***REMOVED***
		// Remove the failed state after some time,
		// making the manager call createCert again on the following TLS hello.
		time.AfterFunc(createCertRetryAfter, func() ***REMOVED***
			defer testDidRemoveState(domain)
			m.stateMu.Lock()
			defer m.stateMu.Unlock()
			// Verify the state hasn't changed and it's still invalid
			// before deleting.
			s, ok := m.state[domain]
			if !ok ***REMOVED***
				return
			***REMOVED***
			if _, err := validCert(domain, s.cert, s.key); err == nil ***REMOVED***
				return
			***REMOVED***
			delete(m.state, domain)
		***REMOVED***)
		return nil, err
	***REMOVED***
	state.cert = der
	state.leaf = leaf
	go m.renew(domain, state.key, state.leaf.NotAfter)
	return state.tlscert()
***REMOVED***

// certState returns a new or existing certState.
// If a new certState is returned, state.exist is false and the state is locked.
// The returned error is non-nil only in the case where a new state could not be created.
func (m *Manager) certState(domain string) (*certState, error) ***REMOVED***
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if m.state == nil ***REMOVED***
		m.state = make(map[string]*certState)
	***REMOVED***
	// existing state
	if state, ok := m.state[domain]; ok ***REMOVED***
		return state, nil
	***REMOVED***

	// new locked state
	var (
		err error
		key crypto.Signer
	)
	if m.ForceRSA ***REMOVED***
		key, err = rsa.GenerateKey(rand.Reader, 2048)
	***REMOVED*** else ***REMOVED***
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	state := &certState***REMOVED***
		key:    key,
		locked: true,
	***REMOVED***
	state.Lock() // will be unlocked by m.certState caller
	m.state[domain] = state
	return state, nil
***REMOVED***

// authorizedCert starts the domain ownership verification process and requests a new cert upon success.
// The key argument is the certificate private key.
func (m *Manager) authorizedCert(ctx context.Context, key crypto.Signer, domain string) (der [][]byte, leaf *x509.Certificate, err error) ***REMOVED***
	client, err := m.acmeClient(ctx)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if err := m.verify(ctx, client, domain); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	csr, err := certRequest(key, domain)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	der, _, err = client.CreateCert(ctx, csr, 0, true)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	leaf, err = validCert(domain, der, key)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return der, leaf, nil
***REMOVED***

// verify runs the identifier (domain) authorization flow
// using each applicable ACME challenge type.
func (m *Manager) verify(ctx context.Context, client *acme.Client, domain string) error ***REMOVED***
	// The list of challenge types we'll try to fulfill
	// in this specific order.
	challengeTypes := []string***REMOVED***"tls-sni-02", "tls-sni-01"***REMOVED***
	m.tokensMu.RLock()
	if m.tryHTTP01 ***REMOVED***
		challengeTypes = append(challengeTypes, "http-01")
	***REMOVED***
	m.tokensMu.RUnlock()

	var nextTyp int // challengeType index of the next challenge type to try
	for ***REMOVED***
		// Start domain authorization and get the challenge.
		authz, err := client.Authorize(ctx, domain)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		// No point in accepting challenges if the authorization status
		// is in a final state.
		switch authz.Status ***REMOVED***
		case acme.StatusValid:
			return nil // already authorized
		case acme.StatusInvalid:
			return fmt.Errorf("acme/autocert: invalid authorization %q", authz.URI)
		***REMOVED***

		// Pick the next preferred challenge.
		var chal *acme.Challenge
		for chal == nil && nextTyp < len(challengeTypes) ***REMOVED***
			chal = pickChallenge(challengeTypes[nextTyp], authz.Challenges)
			nextTyp++
		***REMOVED***
		if chal == nil ***REMOVED***
			return fmt.Errorf("acme/autocert: unable to authorize %q; tried %q", domain, challengeTypes)
		***REMOVED***
		cleanup, err := m.fulfill(ctx, client, chal)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		defer cleanup()
		if _, err := client.Accept(ctx, chal); err != nil ***REMOVED***
			continue
		***REMOVED***

		// A challenge is fulfilled and accepted: wait for the CA to validate.
		if _, err := client.WaitAuthorization(ctx, authz.URI); err == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// fulfill provisions a response to the challenge chal.
// The cleanup is non-nil only if provisioning succeeded.
func (m *Manager) fulfill(ctx context.Context, client *acme.Client, chal *acme.Challenge) (cleanup func(), err error) ***REMOVED***
	switch chal.Type ***REMOVED***
	case "tls-sni-01":
		cert, name, err := client.TLSSNI01ChallengeCert(chal.Token)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m.putCertToken(ctx, name, &cert)
		return func() ***REMOVED*** go m.deleteCertToken(name) ***REMOVED***, nil
	case "tls-sni-02":
		cert, name, err := client.TLSSNI02ChallengeCert(chal.Token)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m.putCertToken(ctx, name, &cert)
		return func() ***REMOVED*** go m.deleteCertToken(name) ***REMOVED***, nil
	case "http-01":
		resp, err := client.HTTP01ChallengeResponse(chal.Token)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		p := client.HTTP01ChallengePath(chal.Token)
		m.putHTTPToken(ctx, p, resp)
		return func() ***REMOVED*** go m.deleteHTTPToken(p) ***REMOVED***, nil
	***REMOVED***
	return nil, fmt.Errorf("acme/autocert: unknown challenge type %q", chal.Type)
***REMOVED***

func pickChallenge(typ string, chal []*acme.Challenge) *acme.Challenge ***REMOVED***
	for _, c := range chal ***REMOVED***
		if c.Type == typ ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// putCertToken stores the cert under the named key in both m.certTokens map
// and m.Cache.
func (m *Manager) putCertToken(ctx context.Context, name string, cert *tls.Certificate) ***REMOVED***
	m.tokensMu.Lock()
	defer m.tokensMu.Unlock()
	if m.certTokens == nil ***REMOVED***
		m.certTokens = make(map[string]*tls.Certificate)
	***REMOVED***
	m.certTokens[name] = cert
	m.cachePut(ctx, name, cert)
***REMOVED***

// deleteCertToken removes the token certificate for the specified domain name
// from both m.certTokens map and m.Cache.
func (m *Manager) deleteCertToken(name string) ***REMOVED***
	m.tokensMu.Lock()
	defer m.tokensMu.Unlock()
	delete(m.certTokens, name)
	if m.Cache != nil ***REMOVED***
		m.Cache.Delete(context.Background(), name)
	***REMOVED***
***REMOVED***

// httpToken retrieves an existing http-01 token value from an in-memory map
// or the optional cache.
func (m *Manager) httpToken(ctx context.Context, tokenPath string) ([]byte, error) ***REMOVED***
	m.tokensMu.RLock()
	defer m.tokensMu.RUnlock()
	if v, ok := m.httpTokens[tokenPath]; ok ***REMOVED***
		return v, nil
	***REMOVED***
	if m.Cache == nil ***REMOVED***
		return nil, fmt.Errorf("acme/autocert: no token at %q", tokenPath)
	***REMOVED***
	return m.Cache.Get(ctx, httpTokenCacheKey(tokenPath))
***REMOVED***

// putHTTPToken stores an http-01 token value using tokenPath as key
// in both in-memory map and the optional Cache.
//
// It ignores any error returned from Cache.Put.
func (m *Manager) putHTTPToken(ctx context.Context, tokenPath, val string) ***REMOVED***
	m.tokensMu.Lock()
	defer m.tokensMu.Unlock()
	if m.httpTokens == nil ***REMOVED***
		m.httpTokens = make(map[string][]byte)
	***REMOVED***
	b := []byte(val)
	m.httpTokens[tokenPath] = b
	if m.Cache != nil ***REMOVED***
		m.Cache.Put(ctx, httpTokenCacheKey(tokenPath), b)
	***REMOVED***
***REMOVED***

// deleteHTTPToken removes an http-01 token value from both in-memory map
// and the optional Cache, ignoring any error returned from the latter.
//
// If m.Cache is non-nil, it blocks until Cache.Delete returns without a timeout.
func (m *Manager) deleteHTTPToken(tokenPath string) ***REMOVED***
	m.tokensMu.Lock()
	defer m.tokensMu.Unlock()
	delete(m.httpTokens, tokenPath)
	if m.Cache != nil ***REMOVED***
		m.Cache.Delete(context.Background(), httpTokenCacheKey(tokenPath))
	***REMOVED***
***REMOVED***

// httpTokenCacheKey returns a key at which an http-01 token value may be stored
// in the Manager's optional Cache.
func httpTokenCacheKey(tokenPath string) string ***REMOVED***
	return "http-01-" + path.Base(tokenPath)
***REMOVED***

// renew starts a cert renewal timer loop, one per domain.
//
// The loop is scheduled in two cases:
// - a cert was fetched from cache for the first time (wasn't in m.state)
// - a new cert was created by m.createCert
//
// The key argument is a certificate private key.
// The exp argument is the cert expiration time (NotAfter).
func (m *Manager) renew(domain string, key crypto.Signer, exp time.Time) ***REMOVED***
	m.renewalMu.Lock()
	defer m.renewalMu.Unlock()
	if m.renewal[domain] != nil ***REMOVED***
		// another goroutine is already on it
		return
	***REMOVED***
	if m.renewal == nil ***REMOVED***
		m.renewal = make(map[string]*domainRenewal)
	***REMOVED***
	dr := &domainRenewal***REMOVED***m: m, domain: domain, key: key***REMOVED***
	m.renewal[domain] = dr
	dr.start(exp)
***REMOVED***

// stopRenew stops all currently running cert renewal timers.
// The timers are not restarted during the lifetime of the Manager.
func (m *Manager) stopRenew() ***REMOVED***
	m.renewalMu.Lock()
	defer m.renewalMu.Unlock()
	for name, dr := range m.renewal ***REMOVED***
		delete(m.renewal, name)
		dr.stop()
	***REMOVED***
***REMOVED***

func (m *Manager) accountKey(ctx context.Context) (crypto.Signer, error) ***REMOVED***
	const keyName = "acme_account.key"

	genKey := func() (*ecdsa.PrivateKey, error) ***REMOVED***
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	***REMOVED***

	if m.Cache == nil ***REMOVED***
		return genKey()
	***REMOVED***

	data, err := m.Cache.Get(ctx, keyName)
	if err == ErrCacheMiss ***REMOVED***
		key, err := genKey()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var buf bytes.Buffer
		if err := encodeECDSAKey(&buf, key); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if err := m.Cache.Put(ctx, keyName, buf.Bytes()); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return key, nil
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	priv, _ := pem.Decode(data)
	if priv == nil || !strings.Contains(priv.Type, "PRIVATE") ***REMOVED***
		return nil, errors.New("acme/autocert: invalid account key found in cache")
	***REMOVED***
	return parsePrivateKey(priv.Bytes)
***REMOVED***

func (m *Manager) acmeClient(ctx context.Context) (*acme.Client, error) ***REMOVED***
	m.clientMu.Lock()
	defer m.clientMu.Unlock()
	if m.client != nil ***REMOVED***
		return m.client, nil
	***REMOVED***

	client := m.Client
	if client == nil ***REMOVED***
		client = &acme.Client***REMOVED***DirectoryURL: acme.LetsEncryptURL***REMOVED***
	***REMOVED***
	if client.Key == nil ***REMOVED***
		var err error
		client.Key, err = m.accountKey(ctx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	var contact []string
	if m.Email != "" ***REMOVED***
		contact = []string***REMOVED***"mailto:" + m.Email***REMOVED***
	***REMOVED***
	a := &acme.Account***REMOVED***Contact: contact***REMOVED***
	_, err := client.Register(ctx, a, m.Prompt)
	if ae, ok := err.(*acme.Error); err == nil || ok && ae.StatusCode == http.StatusConflict ***REMOVED***
		// conflict indicates the key is already registered
		m.client = client
		err = nil
	***REMOVED***
	return m.client, err
***REMOVED***

func (m *Manager) hostPolicy() HostPolicy ***REMOVED***
	if m.HostPolicy != nil ***REMOVED***
		return m.HostPolicy
	***REMOVED***
	return defaultHostPolicy
***REMOVED***

func (m *Manager) renewBefore() time.Duration ***REMOVED***
	if m.RenewBefore > renewJitter ***REMOVED***
		return m.RenewBefore
	***REMOVED***
	return 720 * time.Hour // 30 days
***REMOVED***

// certState is ready when its mutex is unlocked for reading.
type certState struct ***REMOVED***
	sync.RWMutex
	locked bool              // locked for read/write
	key    crypto.Signer     // private key for cert
	cert   [][]byte          // DER encoding
	leaf   *x509.Certificate // parsed cert[0]; always non-nil if cert != nil
***REMOVED***

// tlscert creates a tls.Certificate from s.key and s.cert.
// Callers should wrap it in s.RLock() and s.RUnlock().
func (s *certState) tlscert() (*tls.Certificate, error) ***REMOVED***
	if s.key == nil ***REMOVED***
		return nil, errors.New("acme/autocert: missing signer")
	***REMOVED***
	if len(s.cert) == 0 ***REMOVED***
		return nil, errors.New("acme/autocert: missing certificate")
	***REMOVED***
	return &tls.Certificate***REMOVED***
		PrivateKey:  s.key,
		Certificate: s.cert,
		Leaf:        s.leaf,
	***REMOVED***, nil
***REMOVED***

// certRequest creates a certificate request for the given common name cn
// and optional SANs.
func certRequest(key crypto.Signer, cn string, san ...string) ([]byte, error) ***REMOVED***
	req := &x509.CertificateRequest***REMOVED***
		Subject:  pkix.Name***REMOVED***CommonName: cn***REMOVED***,
		DNSNames: san,
	***REMOVED***
	return x509.CreateCertificateRequest(rand.Reader, req, key)
***REMOVED***

// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
//
// Inspired by parsePrivateKey in crypto/tls/tls.go.
func parsePrivateKey(der []byte) (crypto.Signer, error) ***REMOVED***
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil ***REMOVED***
		return key, nil
	***REMOVED***
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil ***REMOVED***
		switch key := key.(type) ***REMOVED***
		case *rsa.PrivateKey:
			return key, nil
		case *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("acme/autocert: unknown private key type in PKCS#8 wrapping")
		***REMOVED***
	***REMOVED***
	if key, err := x509.ParseECPrivateKey(der); err == nil ***REMOVED***
		return key, nil
	***REMOVED***

	return nil, errors.New("acme/autocert: failed to parse private key")
***REMOVED***

// validCert parses a cert chain provided as der argument and verifies the leaf, der[0],
// corresponds to the private key, as well as the domain match and expiration dates.
// It doesn't do any revocation checking.
//
// The returned value is the verified leaf cert.
func validCert(domain string, der [][]byte, key crypto.Signer) (leaf *x509.Certificate, err error) ***REMOVED***
	// parse public part(s)
	var n int
	for _, b := range der ***REMOVED***
		n += len(b)
	***REMOVED***
	pub := make([]byte, n)
	n = 0
	for _, b := range der ***REMOVED***
		n += copy(pub[n:], b)
	***REMOVED***
	x509Cert, err := x509.ParseCertificates(pub)
	if len(x509Cert) == 0 ***REMOVED***
		return nil, errors.New("acme/autocert: no public key found")
	***REMOVED***
	// verify the leaf is not expired and matches the domain name
	leaf = x509Cert[0]
	now := timeNow()
	if now.Before(leaf.NotBefore) ***REMOVED***
		return nil, errors.New("acme/autocert: certificate is not valid yet")
	***REMOVED***
	if now.After(leaf.NotAfter) ***REMOVED***
		return nil, errors.New("acme/autocert: expired certificate")
	***REMOVED***
	if err := leaf.VerifyHostname(domain); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// ensure the leaf corresponds to the private key
	switch pub := leaf.PublicKey.(type) ***REMOVED***
	case *rsa.PublicKey:
		prv, ok := key.(*rsa.PrivateKey)
		if !ok ***REMOVED***
			return nil, errors.New("acme/autocert: private key type does not match public key type")
		***REMOVED***
		if pub.N.Cmp(prv.N) != 0 ***REMOVED***
			return nil, errors.New("acme/autocert: private key does not match public key")
		***REMOVED***
	case *ecdsa.PublicKey:
		prv, ok := key.(*ecdsa.PrivateKey)
		if !ok ***REMOVED***
			return nil, errors.New("acme/autocert: private key type does not match public key type")
		***REMOVED***
		if pub.X.Cmp(prv.X) != 0 || pub.Y.Cmp(prv.Y) != 0 ***REMOVED***
			return nil, errors.New("acme/autocert: private key does not match public key")
		***REMOVED***
	default:
		return nil, errors.New("acme/autocert: unknown public key algorithm")
	***REMOVED***
	return leaf, nil
***REMOVED***

type lockedMathRand struct ***REMOVED***
	sync.Mutex
	rnd *mathrand.Rand
***REMOVED***

func (r *lockedMathRand) int63n(max int64) int64 ***REMOVED***
	r.Lock()
	n := r.rnd.Int63n(max)
	r.Unlock()
	return n
***REMOVED***

// For easier testing.
var (
	timeNow = time.Now

	// Called when a state is removed.
	testDidRemoveState = func(domain string) ***REMOVED******REMOVED***
)
