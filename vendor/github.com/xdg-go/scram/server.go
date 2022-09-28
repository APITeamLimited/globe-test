// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import "sync"

// Server implements the server side of SCRAM authentication.  It holds
// configuration values needed to initialize new server-side conversations.
// Generally, this can be persistent within an application.
type Server struct ***REMOVED***
	sync.RWMutex
	credentialCB CredentialLookup
	nonceGen     NonceGeneratorFcn
	hashGen      HashGeneratorFcn
***REMOVED***

func newServer(cl CredentialLookup, fcn HashGeneratorFcn) (*Server, error) ***REMOVED***
	return &Server***REMOVED***
		credentialCB: cl,
		nonceGen:     defaultNonceGenerator,
		hashGen:      fcn,
	***REMOVED***, nil
***REMOVED***

// WithNonceGenerator replaces the default nonce generator (base64 encoding of
// 24 bytes from crypto/rand) with a custom generator.  This is provided for
// testing or for users with custom nonce requirements.
func (s *Server) WithNonceGenerator(ng NonceGeneratorFcn) *Server ***REMOVED***
	s.Lock()
	defer s.Unlock()
	s.nonceGen = ng
	return s
***REMOVED***

// NewConversation constructs a server-side authentication conversation.
// Conversations cannot be reused, so this must be called for each new
// authentication attempt.
func (s *Server) NewConversation() *ServerConversation ***REMOVED***
	s.RLock()
	defer s.RUnlock()
	return &ServerConversation***REMOVED***
		nonceGen:     s.nonceGen,
		hashGen:      s.hashGen,
		credentialCB: s.credentialCB,
	***REMOVED***
***REMOVED***
