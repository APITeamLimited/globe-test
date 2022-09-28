// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Copyright (C) MongoDB, Inc. 2018-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"

	"github.com/xdg-go/scram"
	"github.com/xdg-go/stringprep"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	// SCRAMSHA1 holds the mechanism name "SCRAM-SHA-1"
	SCRAMSHA1 = "SCRAM-SHA-1"

	// SCRAMSHA256 holds the mechanism name "SCRAM-SHA-256"
	SCRAMSHA256 = "SCRAM-SHA-256"
)

var (
	// Additional options for the saslStart command to enable a shorter SCRAM conversation
	scramStartOptions bsoncore.Document = bsoncore.BuildDocumentFromElements(nil,
		bsoncore.AppendBooleanElement(nil, "skipEmptyExchange", true),
	)
)

func newScramSHA1Authenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	passdigest := mongoPasswordDigest(cred.Username, cred.Password)
	client, err := scram.SHA1.NewClientUnprepped(cred.Username, passdigest, "")
	if err != nil ***REMOVED***
		return nil, newAuthError("error initializing SCRAM-SHA-1 client", err)
	***REMOVED***
	client.WithMinIterations(4096)
	return &ScramAuthenticator***REMOVED***
		mechanism: SCRAMSHA1,
		source:    cred.Source,
		client:    client,
	***REMOVED***, nil
***REMOVED***

func newScramSHA256Authenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	passprep, err := stringprep.SASLprep.Prepare(cred.Password)
	if err != nil ***REMOVED***
		return nil, newAuthError(fmt.Sprintf("error SASLprepping password '%s'", cred.Password), err)
	***REMOVED***
	client, err := scram.SHA256.NewClientUnprepped(cred.Username, passprep, "")
	if err != nil ***REMOVED***
		return nil, newAuthError("error initializing SCRAM-SHA-256 client", err)
	***REMOVED***
	client.WithMinIterations(4096)
	return &ScramAuthenticator***REMOVED***
		mechanism: SCRAMSHA256,
		source:    cred.Source,
		client:    client,
	***REMOVED***, nil
***REMOVED***

// ScramAuthenticator uses the SCRAM algorithm over SASL to authenticate a connection.
type ScramAuthenticator struct ***REMOVED***
	mechanism string
	source    string
	client    *scram.Client
***REMOVED***

var _ SpeculativeAuthenticator = (*ScramAuthenticator)(nil)

// Auth authenticates the provided connection by conducting a full SASL conversation.
func (a *ScramAuthenticator) Auth(ctx context.Context, cfg *Config) error ***REMOVED***
	err := ConductSaslConversation(ctx, cfg, a.source, a.createSaslClient())
	if err != nil ***REMOVED***
		return newAuthError("sasl conversation error", err)
	***REMOVED***
	return nil
***REMOVED***

// CreateSpeculativeConversation creates a speculative conversation for SCRAM authentication.
func (a *ScramAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) ***REMOVED***
	return newSaslConversation(a.createSaslClient(), a.source, true), nil
***REMOVED***

func (a *ScramAuthenticator) createSaslClient() SaslClient ***REMOVED***
	return &scramSaslAdapter***REMOVED***
		conversation: a.client.NewConversation(),
		mechanism:    a.mechanism,
	***REMOVED***
***REMOVED***

type scramSaslAdapter struct ***REMOVED***
	mechanism    string
	conversation *scram.ClientConversation
***REMOVED***

var _ SaslClient = (*scramSaslAdapter)(nil)
var _ ExtraOptionsSaslClient = (*scramSaslAdapter)(nil)

func (a *scramSaslAdapter) Start() (string, []byte, error) ***REMOVED***
	step, err := a.conversation.Step("")
	if err != nil ***REMOVED***
		return a.mechanism, nil, err
	***REMOVED***
	return a.mechanism, []byte(step), nil
***REMOVED***

func (a *scramSaslAdapter) Next(challenge []byte) ([]byte, error) ***REMOVED***
	step, err := a.conversation.Step(string(challenge))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return []byte(step), nil
***REMOVED***

func (a *scramSaslAdapter) Completed() bool ***REMOVED***
	return a.conversation.Done()
***REMOVED***

func (*scramSaslAdapter) StartCommandOptions() bsoncore.Document ***REMOVED***
	return scramStartOptions
***REMOVED***
