// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/description"
)

func newDefaultAuthenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	scram, err := newScramSHA256Authenticator(cred)
	if err != nil ***REMOVED***
		return nil, newAuthError("failed to create internal authenticator", err)
	***REMOVED***
	speculative, ok := scram.(SpeculativeAuthenticator)
	if !ok ***REMOVED***
		typeErr := fmt.Errorf("expected SCRAM authenticator to be SpeculativeAuthenticator but got %T", scram)
		return nil, newAuthError("failed to create internal authenticator", typeErr)
	***REMOVED***

	return &DefaultAuthenticator***REMOVED***
		Cred:                     cred,
		speculativeAuthenticator: speculative,
	***REMOVED***, nil
***REMOVED***

// DefaultAuthenticator uses SCRAM-SHA-1 or MONGODB-CR depending
// on the server version.
type DefaultAuthenticator struct ***REMOVED***
	Cred *Cred

	// The authenticator to use for speculative authentication. Because the correct auth mechanism is unknown when doing
	// the initial hello, SCRAM-SHA-256 is used for the speculative attempt.
	speculativeAuthenticator SpeculativeAuthenticator
***REMOVED***

var _ SpeculativeAuthenticator = (*DefaultAuthenticator)(nil)

// CreateSpeculativeConversation creates a speculative conversation for SCRAM authentication.
func (a *DefaultAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) ***REMOVED***
	return a.speculativeAuthenticator.CreateSpeculativeConversation()
***REMOVED***

// Auth authenticates the connection.
func (a *DefaultAuthenticator) Auth(ctx context.Context, cfg *Config) error ***REMOVED***
	var actual Authenticator
	var err error

	switch chooseAuthMechanism(cfg) ***REMOVED***
	case SCRAMSHA256:
		actual, err = newScramSHA256Authenticator(a.Cred)
	case SCRAMSHA1:
		actual, err = newScramSHA1Authenticator(a.Cred)
	default:
		actual, err = newMongoDBCRAuthenticator(a.Cred)
	***REMOVED***

	if err != nil ***REMOVED***
		return newAuthError("error creating authenticator", err)
	***REMOVED***

	return actual.Auth(ctx, cfg)
***REMOVED***

// If a server provides a list of supported mechanisms, we choose
// SCRAM-SHA-256 if it exists or else MUST use SCRAM-SHA-1.
// Otherwise, we decide based on what is supported.
func chooseAuthMechanism(cfg *Config) string ***REMOVED***
	if saslSupportedMechs := cfg.HandshakeInfo.SaslSupportedMechs; saslSupportedMechs != nil ***REMOVED***
		for _, v := range saslSupportedMechs ***REMOVED***
			if v == SCRAMSHA256 ***REMOVED***
				return v
			***REMOVED***
		***REMOVED***
		return SCRAMSHA1
	***REMOVED***

	if err := scramSHA1Supported(cfg.HandshakeInfo.Description.WireVersion); err == nil ***REMOVED***
		return SCRAMSHA1
	***REMOVED***

	return MONGODBCR
***REMOVED***

// scramSHA1Supported returns an error if the given server version does not support scram-sha-1.
func scramSHA1Supported(wireVersion *description.VersionRange) error ***REMOVED***
	if wireVersion != nil && wireVersion.Max < 3 ***REMOVED***
		return fmt.Errorf("SCRAM-SHA-1 is only supported for servers 3.0 or newer")
	***REMOVED***

	return nil
***REMOVED***
