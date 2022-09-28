// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
)

// PLAIN is the mechanism name for PLAIN.
const PLAIN = "PLAIN"

func newPlainAuthenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	return &PlainAuthenticator***REMOVED***
		Username: cred.Username,
		Password: cred.Password,
	***REMOVED***, nil
***REMOVED***

// PlainAuthenticator uses the PLAIN algorithm over SASL to authenticate a connection.
type PlainAuthenticator struct ***REMOVED***
	Username string
	Password string
***REMOVED***

// Auth authenticates the connection.
func (a *PlainAuthenticator) Auth(ctx context.Context, cfg *Config) error ***REMOVED***
	return ConductSaslConversation(ctx, cfg, "$external", &plainSaslClient***REMOVED***
		username: a.Username,
		password: a.Password,
	***REMOVED***)
***REMOVED***

type plainSaslClient struct ***REMOVED***
	username string
	password string
***REMOVED***

var _ SaslClient = (*plainSaslClient)(nil)

func (c *plainSaslClient) Start() (string, []byte, error) ***REMOVED***
	b := []byte("\x00" + c.username + "\x00" + c.password)
	return PLAIN, b, nil
***REMOVED***

func (c *plainSaslClient) Next(challenge []byte) ([]byte, error) ***REMOVED***
	return nil, newAuthError("unexpected server challenge", nil)
***REMOVED***

func (c *plainSaslClient) Completed() bool ***REMOVED***
	return true
***REMOVED***
