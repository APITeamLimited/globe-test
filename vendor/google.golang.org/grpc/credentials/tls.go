/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package credentials

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"

	"google.golang.org/grpc/credentials/internal"
	credinternal "google.golang.org/grpc/internal/credentials"
)

// TLSInfo contains the auth information for a TLS authenticated connection.
// It implements the AuthInfo interface.
type TLSInfo struct ***REMOVED***
	State tls.ConnectionState
	CommonAuthInfo
	// This API is experimental.
	SPIFFEID *url.URL
***REMOVED***

// AuthType returns the type of TLSInfo as a string.
func (t TLSInfo) AuthType() string ***REMOVED***
	return "tls"
***REMOVED***

// GetSecurityValue returns security info requested by channelz.
func (t TLSInfo) GetSecurityValue() ChannelzSecurityValue ***REMOVED***
	v := &TLSChannelzSecurityValue***REMOVED***
		StandardName: cipherSuiteLookup[t.State.CipherSuite],
	***REMOVED***
	// Currently there's no way to get LocalCertificate info from tls package.
	if len(t.State.PeerCertificates) > 0 ***REMOVED***
		v.RemoteCertificate = t.State.PeerCertificates[0].Raw
	***REMOVED***
	return v
***REMOVED***

// tlsCreds is the credentials required for authenticating a connection using TLS.
type tlsCreds struct ***REMOVED***
	// TLS configuration
	config *tls.Config
***REMOVED***

func (c tlsCreds) Info() ProtocolInfo ***REMOVED***
	return ProtocolInfo***REMOVED***
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
		ServerName:       c.config.ServerName,
	***REMOVED***
***REMOVED***

func (c *tlsCreds) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (_ net.Conn, _ AuthInfo, err error) ***REMOVED***
	// use local cfg to avoid clobbering ServerName if using multiple endpoints
	cfg := cloneTLSConfig(c.config)
	if cfg.ServerName == "" ***REMOVED***
		serverName, _, err := net.SplitHostPort(authority)
		if err != nil ***REMOVED***
			// If the authority had no host port or if the authority cannot be parsed, use it as-is.
			serverName = authority
		***REMOVED***
		cfg.ServerName = serverName
	***REMOVED***
	conn := tls.Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() ***REMOVED***
		errChannel <- conn.Handshake()
		close(errChannel)
	***REMOVED***()
	select ***REMOVED***
	case err := <-errChannel:
		if err != nil ***REMOVED***
			conn.Close()
			return nil, nil, err
		***REMOVED***
	case <-ctx.Done():
		conn.Close()
		return nil, nil, ctx.Err()
	***REMOVED***
	tlsInfo := TLSInfo***REMOVED***
		State: conn.ConnectionState(),
		CommonAuthInfo: CommonAuthInfo***REMOVED***
			SecurityLevel: PrivacyAndIntegrity,
		***REMOVED***,
	***REMOVED***
	id := credinternal.SPIFFEIDFromState(conn.ConnectionState())
	if id != nil ***REMOVED***
		tlsInfo.SPIFFEID = id
	***REMOVED***
	return internal.WrapSyscallConn(rawConn, conn), tlsInfo, nil
***REMOVED***

func (c *tlsCreds) ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error) ***REMOVED***
	conn := tls.Server(rawConn, c.config)
	if err := conn.Handshake(); err != nil ***REMOVED***
		conn.Close()
		return nil, nil, err
	***REMOVED***
	tlsInfo := TLSInfo***REMOVED***
		State: conn.ConnectionState(),
		CommonAuthInfo: CommonAuthInfo***REMOVED***
			SecurityLevel: PrivacyAndIntegrity,
		***REMOVED***,
	***REMOVED***
	id := credinternal.SPIFFEIDFromState(conn.ConnectionState())
	if id != nil ***REMOVED***
		tlsInfo.SPIFFEID = id
	***REMOVED***
	return internal.WrapSyscallConn(rawConn, conn), tlsInfo, nil
***REMOVED***

func (c *tlsCreds) Clone() TransportCredentials ***REMOVED***
	return NewTLS(c.config)
***REMOVED***

func (c *tlsCreds) OverrideServerName(serverNameOverride string) error ***REMOVED***
	c.config.ServerName = serverNameOverride
	return nil
***REMOVED***

const alpnProtoStrH2 = "h2"

func appendH2ToNextProtos(ps []string) []string ***REMOVED***
	for _, p := range ps ***REMOVED***
		if p == alpnProtoStrH2 ***REMOVED***
			return ps
		***REMOVED***
	***REMOVED***
	ret := make([]string, 0, len(ps)+1)
	ret = append(ret, ps...)
	return append(ret, alpnProtoStrH2)
***REMOVED***

// NewTLS uses c to construct a TransportCredentials based on TLS.
func NewTLS(c *tls.Config) TransportCredentials ***REMOVED***
	tc := &tlsCreds***REMOVED***cloneTLSConfig(c)***REMOVED***
	tc.config.NextProtos = appendH2ToNextProtos(tc.config.NextProtos)
	return tc
***REMOVED***

// NewClientTLSFromCert constructs TLS credentials from the provided root
// certificate authority certificate(s) to validate server connections. If
// certificates to establish the identity of the client need to be included in
// the credentials (eg: for mTLS), use NewTLS instead, where a complete
// tls.Config can be specified.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header
// field) in requests.
func NewClientTLSFromCert(cp *x509.CertPool, serverNameOverride string) TransportCredentials ***REMOVED***
	return NewTLS(&tls.Config***REMOVED***ServerName: serverNameOverride, RootCAs: cp***REMOVED***)
***REMOVED***

// NewClientTLSFromFile constructs TLS credentials from the provided root
// certificate authority certificate file(s) to validate server connections. If
// certificates to establish the identity of the client need to be included in
// the credentials (eg: for mTLS), use NewTLS instead, where a complete
// tls.Config can be specified.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header
// field) in requests.
func NewClientTLSFromFile(certFile, serverNameOverride string) (TransportCredentials, error) ***REMOVED***
	b, err := ioutil.ReadFile(certFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) ***REMOVED***
		return nil, fmt.Errorf("credentials: failed to append certificates")
	***REMOVED***
	return NewTLS(&tls.Config***REMOVED***ServerName: serverNameOverride, RootCAs: cp***REMOVED***), nil
***REMOVED***

// NewServerTLSFromCert constructs TLS credentials from the input certificate for server.
func NewServerTLSFromCert(cert *tls.Certificate) TransportCredentials ***REMOVED***
	return NewTLS(&tls.Config***REMOVED***Certificates: []tls.Certificate***REMOVED****cert***REMOVED******REMOVED***)
***REMOVED***

// NewServerTLSFromFile constructs TLS credentials from the input certificate file and key
// file for server.
func NewServerTLSFromFile(certFile, keyFile string) (TransportCredentials, error) ***REMOVED***
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewTLS(&tls.Config***REMOVED***Certificates: []tls.Certificate***REMOVED***cert***REMOVED******REMOVED***), nil
***REMOVED***

// TLSChannelzSecurityValue defines the struct that TLS protocol should return
// from GetSecurityValue(), containing security info like cipher and certificate used.
//
// This API is EXPERIMENTAL.
type TLSChannelzSecurityValue struct ***REMOVED***
	ChannelzSecurityValue
	StandardName      string
	LocalCertificate  []byte
	RemoteCertificate []byte
***REMOVED***

var cipherSuiteLookup = map[uint16]string***REMOVED***
	tls.TLS_RSA_WITH_RC4_128_SHA:                "TLS_RSA_WITH_RC4_128_SHA",
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	tls.TLS_RSA_WITH_AES_128_CBC_SHA:            "TLS_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_RSA_WITH_AES_256_CBC_SHA:            "TLS_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256:         "TLS_RSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384:         "TLS_RSA_WITH_AES_256_GCM_SHA384",
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:        "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:    "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:          "TLS_ECDHE_RSA_WITH_RC4_128_SHA",
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:     "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:      "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:   "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384: "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	tls.TLS_FALLBACK_SCSV:                       "TLS_FALLBACK_SCSV",
	tls.TLS_RSA_WITH_AES_128_CBC_SHA256:         "TLS_RSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:   "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:    "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:  "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
***REMOVED***

// cloneTLSConfig returns a shallow clone of the exported
// fields of cfg, ignoring the unexported sync.Once, which
// contains a mutex and must not be copied.
//
// If cfg is nil, a new zero tls.Config is returned.
//
// TODO: inline this function if possible.
func cloneTLSConfig(cfg *tls.Config) *tls.Config ***REMOVED***
	if cfg == nil ***REMOVED***
		return &tls.Config***REMOVED******REMOVED***
	***REMOVED***

	return cfg.Clone()
***REMOVED***
