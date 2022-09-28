// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build gssapi && (windows || linux || darwin)
// +build gssapi
// +build windows linux darwin

package auth

import (
	"context"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/internal/gssapi"
)

// GSSAPI is the mechanism name for GSSAPI.
const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	if cred.Source != "" && cred.Source != "$external" ***REMOVED***
		return nil, newAuthError("GSSAPI source must be empty or $external", nil)
	***REMOVED***

	return &GSSAPIAuthenticator***REMOVED***
		Username:    cred.Username,
		Password:    cred.Password,
		PasswordSet: cred.PasswordSet,
		Props:       cred.Props,
	***REMOVED***, nil
***REMOVED***

// GSSAPIAuthenticator uses the GSSAPI algorithm over SASL to authenticate a connection.
type GSSAPIAuthenticator struct ***REMOVED***
	Username    string
	Password    string
	PasswordSet bool
	Props       map[string]string
***REMOVED***

// Auth authenticates the connection.
func (a *GSSAPIAuthenticator) Auth(ctx context.Context, cfg *Config) error ***REMOVED***
	target := cfg.Description.Addr.String()
	hostname, _, err := net.SplitHostPort(target)
	if err != nil ***REMOVED***
		return newAuthError(fmt.Sprintf("invalid endpoint (%s) specified: %s", target, err), nil)
	***REMOVED***

	client, err := gssapi.New(hostname, a.Username, a.Password, a.PasswordSet, a.Props)

	if err != nil ***REMOVED***
		return newAuthError("error creating gssapi", err)
	***REMOVED***
	return ConductSaslConversation(ctx, cfg, "$external", client)
***REMOVED***
