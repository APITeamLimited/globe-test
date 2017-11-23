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
	"net/http"
	"strconv"
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
// It obtains and refreshes certificates automatically,
// as well as providing them to a TLS server via tls.Config.
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

	// tokenCert is keyed by token domain name, which matches server name
	// of ClientHello. Keys always have ".acme.invalid" suffix.
	tokenCertMu sync.RWMutex
	tokenCert   map[string]*tls.Certificate

	// renewal tracks the set of domains currently running renewal timers.
	// It is keyed by domain name.
	renewalMu sync.Mutex
	renewal   map[string]*domainRenewal
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// check whether this is a token cert requested for TLS-SNI challenge
	if strings.HasSuffix(name, ".acme.invalid") ***REMOVED***
		m.tokenCertMu.RLock()
		defer m.tokenCertMu.RUnlock()
		if cert := m.tokenCert[name]; cert != nil ***REMOVED***
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
	if err := m.verify(ctx, domain); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	client, err := m.acmeClient(ctx)
	if err != nil ***REMOVED***
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

// verify starts a new identifier (domain) authorization flow.
// It prepares a challenge response and then blocks until the authorization
// is marked as "completed" by the CA (either succeeded or failed).
//
// verify returns nil iff the verification was successful.
func (m *Manager) verify(ctx context.Context, domain string) error ***REMOVED***
	client, err := m.acmeClient(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// start domain authorization and get the challenge
	authz, err := client.Authorize(ctx, domain)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// maybe don't need to at all
	if authz.Status == acme.StatusValid ***REMOVED***
		return nil
	***REMOVED***

	// pick a challenge: prefer tls-sni-02 over tls-sni-01
	// TODO: consider authz.Combinations
	var chal *acme.Challenge
	for _, c := range authz.Challenges ***REMOVED***
		if c.Type == "tls-sni-02" ***REMOVED***
			chal = c
			break
		***REMOVED***
		if c.Type == "tls-sni-01" ***REMOVED***
			chal = c
		***REMOVED***
	***REMOVED***
	if chal == nil ***REMOVED***
		return errors.New("acme/autocert: no supported challenge type found")
	***REMOVED***

	// create a token cert for the challenge response
	var (
		cert tls.Certificate
		name string
	)
	switch chal.Type ***REMOVED***
	case "tls-sni-01":
		cert, name, err = client.TLSSNI01ChallengeCert(chal.Token)
	case "tls-sni-02":
		cert, name, err = client.TLSSNI02ChallengeCert(chal.Token)
	default:
		err = fmt.Errorf("acme/autocert: unknown challenge type %q", chal.Type)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m.putTokenCert(ctx, name, &cert)
	defer func() ***REMOVED***
		// verification has ended at this point
		// don't need token cert anymore
		go m.deleteTokenCert(name)
	***REMOVED***()

	// ready to fulfill the challenge
	if _, err := client.Accept(ctx, chal); err != nil ***REMOVED***
		return err
	***REMOVED***
	// wait for the CA to validate
	_, err = client.WaitAuthorization(ctx, authz.URI)
	return err
***REMOVED***

// putTokenCert stores the cert under the named key in both m.tokenCert map
// and m.Cache.
func (m *Manager) putTokenCert(ctx context.Context, name string, cert *tls.Certificate) ***REMOVED***
	m.tokenCertMu.Lock()
	defer m.tokenCertMu.Unlock()
	if m.tokenCert == nil ***REMOVED***
		m.tokenCert = make(map[string]*tls.Certificate)
	***REMOVED***
	m.tokenCert[name] = cert
	m.cachePut(ctx, name, cert)
***REMOVED***

// deleteTokenCert removes the token certificate for the specified domain name
// from both m.tokenCert map and m.Cache.
func (m *Manager) deleteTokenCert(name string) ***REMOVED***
	m.tokenCertMu.Lock()
	defer m.tokenCertMu.Unlock()
	delete(m.tokenCert, name)
	if m.Cache != nil ***REMOVED***
		m.Cache.Delete(context.Background(), name)
	***REMOVED***
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

func retryAfter(v string) time.Duration ***REMOVED***
	if i, err := strconv.Atoi(v); err == nil ***REMOVED***
		return time.Duration(i) * time.Second
	***REMOVED***
	if t, err := http.ParseTime(v); err == nil ***REMOVED***
		return t.Sub(timeNow())
	***REMOVED***
	return time.Second
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
