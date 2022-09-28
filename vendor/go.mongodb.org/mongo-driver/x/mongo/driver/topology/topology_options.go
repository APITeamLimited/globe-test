// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// Option is a configuration option for a topology.
type Option func(*config) error

type config struct ***REMOVED***
	mode                   MonitorMode
	replicaSetName         string
	seedList               []string
	serverOpts             []ServerOption
	cs                     connstring.ConnString // This must not be used for any logic in topology.Topology.
	uri                    string
	serverSelectionTimeout time.Duration
	serverMonitor          *event.ServerMonitor
	srvMaxHosts            int
	srvServiceName         string
	loadBalanced           bool
***REMOVED***

func newConfig(opts ...Option) (*config, error) ***REMOVED***
	cfg := &config***REMOVED***
		seedList:               []string***REMOVED***"localhost:27017"***REMOVED***,
		serverSelectionTimeout: 30 * time.Second,
	***REMOVED***

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		err := opt(cfg)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return cfg, nil
***REMOVED***

// WithConnString configures the topology using the connection string.
func WithConnString(fn func(connstring.ConnString) connstring.ConnString) Option ***REMOVED***
	return func(c *config) error ***REMOVED***
		cs := fn(c.cs)
		c.cs = cs

		if cs.ServerSelectionTimeoutSet ***REMOVED***
			c.serverSelectionTimeout = cs.ServerSelectionTimeout
		***REMOVED***

		var connOpts []ConnectionOption

		if cs.AppName != "" ***REMOVED***
			c.serverOpts = append(c.serverOpts, WithServerAppName(func(string) string ***REMOVED*** return cs.AppName ***REMOVED***))
		***REMOVED***

		if cs.Connect == connstring.SingleConnect || (cs.DirectConnectionSet && cs.DirectConnection) ***REMOVED***
			c.mode = SingleMode
		***REMOVED***

		c.seedList = cs.Hosts

		if cs.ConnectTimeout > 0 ***REMOVED***
			c.serverOpts = append(c.serverOpts, WithHeartbeatTimeout(func(time.Duration) time.Duration ***REMOVED*** return cs.ConnectTimeout ***REMOVED***))
			connOpts = append(connOpts, WithConnectTimeout(func(time.Duration) time.Duration ***REMOVED*** return cs.ConnectTimeout ***REMOVED***))
		***REMOVED***

		if cs.SocketTimeoutSet ***REMOVED***
			connOpts = append(
				connOpts,
				WithReadTimeout(func(time.Duration) time.Duration ***REMOVED*** return cs.SocketTimeout ***REMOVED***),
				WithWriteTimeout(func(time.Duration) time.Duration ***REMOVED*** return cs.SocketTimeout ***REMOVED***),
			)
		***REMOVED***

		if cs.HeartbeatInterval > 0 ***REMOVED***
			c.serverOpts = append(c.serverOpts, WithHeartbeatInterval(func(time.Duration) time.Duration ***REMOVED*** return cs.HeartbeatInterval ***REMOVED***))
		***REMOVED***

		if cs.MaxConnIdleTime > 0 ***REMOVED***
			connOpts = append(connOpts, WithIdleTimeout(func(time.Duration) time.Duration ***REMOVED*** return cs.MaxConnIdleTime ***REMOVED***))
		***REMOVED***

		if cs.MaxPoolSizeSet ***REMOVED***
			c.serverOpts = append(c.serverOpts, WithMaxConnections(func(uint64) uint64 ***REMOVED*** return cs.MaxPoolSize ***REMOVED***))
		***REMOVED***

		if cs.MinPoolSizeSet ***REMOVED***
			c.serverOpts = append(c.serverOpts, WithMinConnections(func(u uint64) uint64 ***REMOVED*** return cs.MinPoolSize ***REMOVED***))
		***REMOVED***

		if cs.ReplicaSet != "" ***REMOVED***
			c.replicaSetName = cs.ReplicaSet
		***REMOVED***

		var x509Username string
		if cs.SSL ***REMOVED***
			tlsConfig := new(tls.Config)

			if cs.SSLCaFileSet ***REMOVED***
				err := addCACertFromFile(tlsConfig, cs.SSLCaFile)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***

			if cs.SSLInsecure ***REMOVED***
				tlsConfig.InsecureSkipVerify = true
			***REMOVED***

			if cs.SSLClientCertificateKeyFileSet ***REMOVED***
				var keyPasswd string
				if cs.SSLClientCertificateKeyPasswordSet && cs.SSLClientCertificateKeyPassword != nil ***REMOVED***
					keyPasswd = cs.SSLClientCertificateKeyPassword()
				***REMOVED***
				s, err := addClientCertFromFile(tlsConfig, cs.SSLClientCertificateKeyFile, keyPasswd)
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				// The Go x509 package gives the subject with the pairs in reverse order that we want.
				pairs := strings.Split(s, ",")
				b := bytes.NewBufferString("")

				for i := len(pairs) - 1; i >= 0; i-- ***REMOVED***
					b.WriteString(pairs[i])

					if i > 0 ***REMOVED***
						b.WriteString(",")
					***REMOVED***
				***REMOVED***

				x509Username = b.String()
			***REMOVED***

			connOpts = append(connOpts, WithTLSConfig(func(*tls.Config) *tls.Config ***REMOVED*** return tlsConfig ***REMOVED***))
		***REMOVED***

		if cs.Username != "" || cs.AuthMechanism == auth.MongoDBX509 || cs.AuthMechanism == auth.GSSAPI ***REMOVED***
			cred := &auth.Cred***REMOVED***
				Source:      "admin",
				Username:    cs.Username,
				Password:    cs.Password,
				PasswordSet: cs.PasswordSet,
				Props:       cs.AuthMechanismProperties,
			***REMOVED***

			if cs.AuthSource != "" ***REMOVED***
				cred.Source = cs.AuthSource
			***REMOVED*** else ***REMOVED***
				switch cs.AuthMechanism ***REMOVED***
				case auth.MongoDBX509:
					if cred.Username == "" ***REMOVED***
						cred.Username = x509Username
					***REMOVED***
					fallthrough
				case auth.GSSAPI, auth.PLAIN:
					cred.Source = "$external"
				default:
					cred.Source = cs.Database
				***REMOVED***
			***REMOVED***

			authenticator, err := auth.CreateAuthenticator(cs.AuthMechanism, cred)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			connOpts = append(connOpts, WithHandshaker(func(h Handshaker) Handshaker ***REMOVED***
				options := &auth.HandshakeOptions***REMOVED***
					AppName:       cs.AppName,
					Authenticator: authenticator,
					Compressors:   cs.Compressors,
					LoadBalanced:  cs.LoadBalancedSet && cs.LoadBalanced,
				***REMOVED***
				if cs.AuthMechanism == "" ***REMOVED***
					// Required for SASL mechanism negotiation during handshake
					options.DBUser = cred.Source + "." + cred.Username
				***REMOVED***
				return auth.Handshaker(h, options)
			***REMOVED***))
		***REMOVED*** else ***REMOVED***
			// We need to add a non-auth Handshaker to the connection options
			connOpts = append(connOpts, WithHandshaker(func(h driver.Handshaker) driver.Handshaker ***REMOVED***
				return operation.NewHello().
					AppName(cs.AppName).
					Compressors(cs.Compressors).
					LoadBalanced(cs.LoadBalancedSet && cs.LoadBalanced)
			***REMOVED***))
		***REMOVED***

		if len(cs.Compressors) > 0 ***REMOVED***
			connOpts = append(connOpts, WithCompressors(func(compressors []string) []string ***REMOVED***
				return append(compressors, cs.Compressors...)
			***REMOVED***))

			for _, comp := range cs.Compressors ***REMOVED***
				switch comp ***REMOVED***
				case "zlib":
					connOpts = append(connOpts, WithZlibLevel(func(level *int) *int ***REMOVED***
						return &cs.ZlibLevel
					***REMOVED***))
				case "zstd":
					connOpts = append(connOpts, WithZstdLevel(func(level *int) *int ***REMOVED***
						return &cs.ZstdLevel
					***REMOVED***))
				***REMOVED***
			***REMOVED***

			c.serverOpts = append(c.serverOpts, WithCompressionOptions(func(opts ...string) []string ***REMOVED***
				return append(opts, cs.Compressors...)
			***REMOVED***))
		***REMOVED***

		// LoadBalanced
		if cs.LoadBalancedSet ***REMOVED***
			c.loadBalanced = cs.LoadBalanced
			c.serverOpts = append(c.serverOpts, WithServerLoadBalanced(func(bool) bool ***REMOVED***
				return cs.LoadBalanced
			***REMOVED***))
			connOpts = append(connOpts, WithConnectionLoadBalanced(func(bool) bool ***REMOVED***
				return cs.LoadBalanced
			***REMOVED***))
		***REMOVED***

		if len(connOpts) > 0 ***REMOVED***
			c.serverOpts = append(c.serverOpts, WithConnectionOptions(func(opts ...ConnectionOption) []ConnectionOption ***REMOVED***
				return append(opts, connOpts...)
			***REMOVED***))
		***REMOVED***

		return nil
	***REMOVED***
***REMOVED***

// WithMode configures the topology's monitor mode.
func WithMode(fn func(MonitorMode) MonitorMode) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.mode = fn(cfg.mode)
		return nil
	***REMOVED***
***REMOVED***

// WithReplicaSetName configures the topology's default replica set name.
func WithReplicaSetName(fn func(string) string) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.replicaSetName = fn(cfg.replicaSetName)
		return nil
	***REMOVED***
***REMOVED***

// WithSeedList configures a topology's seed list.
func WithSeedList(fn func(...string) []string) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.seedList = fn(cfg.seedList...)
		return nil
	***REMOVED***
***REMOVED***

// WithServerOptions configures a topology's server options for when a new server
// needs to be created.
func WithServerOptions(fn func(...ServerOption) []ServerOption) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.serverOpts = fn(cfg.serverOpts...)
		return nil
	***REMOVED***
***REMOVED***

// WithServerSelectionTimeout configures a topology's server selection timeout.
// A server selection timeout of 0 means there is no timeout for server selection.
func WithServerSelectionTimeout(fn func(time.Duration) time.Duration) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.serverSelectionTimeout = fn(cfg.serverSelectionTimeout)
		return nil
	***REMOVED***
***REMOVED***

// WithTopologyServerMonitor configures the monitor for all SDAM events
func WithTopologyServerMonitor(fn func(*event.ServerMonitor) *event.ServerMonitor) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.serverMonitor = fn(cfg.serverMonitor)
		return nil
	***REMOVED***
***REMOVED***

// WithURI specifies the URI that was used to create the topology.
func WithURI(fn func(string) string) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.uri = fn(cfg.uri)
		return nil
	***REMOVED***
***REMOVED***

// WithLoadBalanced specifies whether or not the cluster is behind a load balancer.
func WithLoadBalanced(fn func(bool) bool) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.loadBalanced = fn(cfg.loadBalanced)
		return nil
	***REMOVED***
***REMOVED***

// WithSRVMaxHosts specifies the SRV host limit that was used to create the topology.
func WithSRVMaxHosts(fn func(int) int) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.srvMaxHosts = fn(cfg.srvMaxHosts)
		return nil
	***REMOVED***
***REMOVED***

// WithSRVServiceName specifies the SRV service name that was used to create the topology.
func WithSRVServiceName(fn func(string) string) Option ***REMOVED***
	return func(cfg *config) error ***REMOVED***
		cfg.srvServiceName = fn(cfg.srvServiceName)
		return nil
	***REMOVED***
***REMOVED***

// addCACertFromFile adds a root CA certificate to the configuration given a path
// to the containing file.
func addCACertFromFile(cfg *tls.Config, file string) error ***REMOVED***
	data, err := ioutil.ReadFile(file)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	certBytes, err := loadCert(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if cfg.RootCAs == nil ***REMOVED***
		cfg.RootCAs = x509.NewCertPool()
	***REMOVED***

	cfg.RootCAs.AddCert(cert)

	return nil
***REMOVED***

func loadCert(data []byte) ([]byte, error) ***REMOVED***
	var certBlock *pem.Block

	for certBlock == nil ***REMOVED***
		if len(data) == 0 ***REMOVED***
			return nil, errors.New(".pem file must have both a CERTIFICATE and an RSA PRIVATE KEY section")
		***REMOVED***

		block, rest := pem.Decode(data)
		if block == nil ***REMOVED***
			return nil, errors.New("invalid .pem file")
		***REMOVED***

		switch block.Type ***REMOVED***
		case "CERTIFICATE":
			if certBlock != nil ***REMOVED***
				return nil, errors.New("multiple CERTIFICATE sections in .pem file")
			***REMOVED***

			certBlock = block
		***REMOVED***

		data = rest
	***REMOVED***

	return certBlock.Bytes, nil
***REMOVED***

// addClientCertFromFile adds a client certificate to the configuration given a path to the
// containing file and returns the certificate's subject name.
func addClientCertFromFile(cfg *tls.Config, clientFile, keyPasswd string) (string, error) ***REMOVED***
	data, err := ioutil.ReadFile(clientFile)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var currentBlock *pem.Block
	var certBlock, certDecodedBlock, keyBlock []byte

	remaining := data
	start := 0
	for ***REMOVED***
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil ***REMOVED***
			break
		***REMOVED***

		if currentBlock.Type == "CERTIFICATE" ***REMOVED***
			certBlock = data[start : len(data)-len(remaining)]
			certDecodedBlock = currentBlock.Bytes
			start += len(certBlock)
		***REMOVED*** else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") ***REMOVED***
			if keyPasswd != "" && x509.IsEncryptedPEMBlock(currentBlock) ***REMOVED***
				var encoded bytes.Buffer
				buf, err := x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***

				pem.Encode(&encoded, &pem.Block***REMOVED***Type: currentBlock.Type, Bytes: buf***REMOVED***)
				keyBlock = encoded.Bytes()
				start = len(data) - len(remaining)
			***REMOVED*** else ***REMOVED***
				keyBlock = data[start : len(data)-len(remaining)]
				start += len(keyBlock)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(certBlock) == 0 ***REMOVED***
		return "", fmt.Errorf("failed to find CERTIFICATE")
	***REMOVED***
	if len(keyBlock) == 0 ***REMOVED***
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	***REMOVED***

	cert, err := tls.X509KeyPair(certBlock, keyBlock)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	cfg.Certificates = append(cfg.Certificates, cert)

	// The documentation for the tls.X509KeyPair indicates that the Leaf certificate is not
	// retained.
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return crt.Subject.String(), nil
***REMOVED***
