// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
)

// MongoDBAWS is the mechanism name for MongoDBAWS.
const MongoDBAWS = "MONGODB-AWS"

func newMongoDBAWSAuthenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	if cred.Source != "" && cred.Source != "$external" ***REMOVED***
		return nil, newAuthError("MONGODB-AWS source must be empty or $external", nil)
	***REMOVED***
	return &MongoDBAWSAuthenticator***REMOVED***
		source:       cred.Source,
		username:     cred.Username,
		password:     cred.Password,
		sessionToken: cred.Props["AWS_SESSION_TOKEN"],
	***REMOVED***, nil
***REMOVED***

// MongoDBAWSAuthenticator uses AWS-IAM credentials over SASL to authenticate a connection.
type MongoDBAWSAuthenticator struct ***REMOVED***
	source       string
	username     string
	password     string
	sessionToken string
***REMOVED***

// Auth authenticates the connection.
func (a *MongoDBAWSAuthenticator) Auth(ctx context.Context, cfg *Config) error ***REMOVED***
	adapter := &awsSaslAdapter***REMOVED***
		conversation: &awsConversation***REMOVED***
			username: a.username,
			password: a.password,
			token:    a.sessionToken,
		***REMOVED***,
	***REMOVED***
	err := ConductSaslConversation(ctx, cfg, a.source, adapter)
	if err != nil ***REMOVED***
		return newAuthError("sasl conversation error", err)
	***REMOVED***
	return nil
***REMOVED***

type awsSaslAdapter struct ***REMOVED***
	conversation *awsConversation
***REMOVED***

var _ SaslClient = (*awsSaslAdapter)(nil)

func (a *awsSaslAdapter) Start() (string, []byte, error) ***REMOVED***
	step, err := a.conversation.Step(nil)
	if err != nil ***REMOVED***
		return MongoDBAWS, nil, err
	***REMOVED***
	return MongoDBAWS, step, nil
***REMOVED***

func (a *awsSaslAdapter) Next(challenge []byte) ([]byte, error) ***REMOVED***
	step, err := a.conversation.Step(challenge)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return step, nil
***REMOVED***

func (a *awsSaslAdapter) Completed() bool ***REMOVED***
	return a.conversation.Done()
***REMOVED***
