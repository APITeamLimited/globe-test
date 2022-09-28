// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"crypto/tls"
	"fmt"
)

// ClientEncryptionOptions represents all possible options used to configure a ClientEncryption instance.
type ClientEncryptionOptions struct ***REMOVED***
	KeyVaultNamespace string
	KmsProviders      map[string]map[string]interface***REMOVED******REMOVED***
	TLSConfig         map[string]*tls.Config
***REMOVED***

// ClientEncryption creates a new ClientEncryptionOptions instance.
func ClientEncryption() *ClientEncryptionOptions ***REMOVED***
	return &ClientEncryptionOptions***REMOVED******REMOVED***
***REMOVED***

// SetKeyVaultNamespace specifies the namespace of the key vault collection. This is required.
func (c *ClientEncryptionOptions) SetKeyVaultNamespace(ns string) *ClientEncryptionOptions ***REMOVED***
	c.KeyVaultNamespace = ns
	return c
***REMOVED***

// SetKmsProviders specifies options for KMS providers. This is required.
func (c *ClientEncryptionOptions) SetKmsProviders(providers map[string]map[string]interface***REMOVED******REMOVED***) *ClientEncryptionOptions ***REMOVED***
	c.KmsProviders = providers
	return c
***REMOVED***

// SetTLSConfig specifies tls.Config instances for each KMS provider to use to configure TLS on all connections created
// to the KMS provider.
//
// This should only be used to set custom TLS configurations. By default, the connection will use an empty tls.Config***REMOVED******REMOVED*** with MinVersion set to tls.VersionTLS12.
func (c *ClientEncryptionOptions) SetTLSConfig(tlsOpts map[string]*tls.Config) *ClientEncryptionOptions ***REMOVED***
	tlsConfigs := make(map[string]*tls.Config)
	for provider, config := range tlsOpts ***REMOVED***
		// use TLS min version 1.2 to enforce more secure hash algorithms and advanced cipher suites
		if config.MinVersion == 0 ***REMOVED***
			config.MinVersion = tls.VersionTLS12
		***REMOVED***
		tlsConfigs[provider] = config
	***REMOVED***
	c.TLSConfig = tlsConfigs
	return c
***REMOVED***

// BuildTLSConfig specifies tls.Config options for each KMS provider to use to configure TLS on all connections created
// to the KMS provider. The input map should contain a mapping from each KMS provider to a document containing the necessary
// options, as follows:
//
// ***REMOVED***
//		"kmip": ***REMOVED***
//			"tlsCertificateKeyFile": "foo.pem",
// 			"tlsCAFile": "fooCA.pem"
//		***REMOVED***
// ***REMOVED***
//
// Currently, the following TLS options are supported:
//
// 1. "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile"): The "tlsCertificateKeyFile" option specifies a path to
// the client certificate and private key, which must be concatenated into one file.
//
// 2. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 3. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// This should only be used to set custom TLS options. By default, the connection will use an empty tls.Config***REMOVED******REMOVED*** with MinVersion set to tls.VersionTLS12.
func BuildTLSConfig(tlsOpts map[string]interface***REMOVED******REMOVED***) (*tls.Config, error) ***REMOVED***
	// use TLS min version 1.2 to enforce more secure hash algorithms and advanced cipher suites
	cfg := &tls.Config***REMOVED***MinVersion: tls.VersionTLS12***REMOVED***

	for name := range tlsOpts ***REMOVED***
		var err error
		switch name ***REMOVED***
		case "tlsCertificateKeyFile", "sslClientCertificateKeyFile":
			clientCertPath, ok := tlsOpts[name].(string)
			if !ok ***REMOVED***
				return nil, fmt.Errorf("expected %q value to be of type string, got %T", name, tlsOpts[name])
			***REMOVED***
			// apply custom key file password if found, otherwise use empty string
			if keyPwd, found := tlsOpts["tlsCertificateKeyFilePassword"].(string); found ***REMOVED***
				_, err = addClientCertFromConcatenatedFile(cfg, clientCertPath, keyPwd)
			***REMOVED*** else if keyPwd, found := tlsOpts["sslClientCertificateKeyPassword"].(string); found ***REMOVED***
				_, err = addClientCertFromConcatenatedFile(cfg, clientCertPath, keyPwd)
			***REMOVED*** else ***REMOVED***
				_, err = addClientCertFromConcatenatedFile(cfg, clientCertPath, "")
			***REMOVED***
		case "tlsCertificateKeyFilePassword", "sslClientCertificateKeyPassword":
			continue
		case "tlsCAFile", "sslCertificateAuthorityFile":
			caPath, ok := tlsOpts[name].(string)
			if !ok ***REMOVED***
				return nil, fmt.Errorf("expected %q value to be of type string, got %T", name, tlsOpts[name])
			***REMOVED***
			err = addCACertFromFile(cfg, caPath)
		default:
			return nil, fmt.Errorf("unrecognized TLS option %v", name)
		***REMOVED***

		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return cfg, nil
***REMOVED***

// MergeClientEncryptionOptions combines the argued ClientEncryptionOptions in a last-one wins fashion.
func MergeClientEncryptionOptions(opts ...*ClientEncryptionOptions) *ClientEncryptionOptions ***REMOVED***
	ceo := ClientEncryption()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***

		if opt.KeyVaultNamespace != "" ***REMOVED***
			ceo.KeyVaultNamespace = opt.KeyVaultNamespace
		***REMOVED***
		if opt.KmsProviders != nil ***REMOVED***
			ceo.KmsProviders = opt.KmsProviders
		***REMOVED***
		if opt.TLSConfig != nil ***REMOVED***
			ceo.TLSConfig = opt.TLSConfig
		***REMOVED***
	***REMOVED***

	return ceo
***REMOVED***
