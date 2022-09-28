// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// MongoDBX509 is the mechanism name for MongoDBX509.
const MongoDBX509 = "MONGODB-X509"

func newMongoDBX509Authenticator(cred *Cred) (Authenticator, error) ***REMOVED***
	return &MongoDBX509Authenticator***REMOVED***User: cred.Username***REMOVED***, nil
***REMOVED***

// MongoDBX509Authenticator uses X.509 certificates over TLS to authenticate a connection.
type MongoDBX509Authenticator struct ***REMOVED***
	User string
***REMOVED***

var _ SpeculativeAuthenticator = (*MongoDBX509Authenticator)(nil)

// x509 represents a X509 authentication conversation. This type implements the SpeculativeConversation interface so the
// conversation can be executed in multi-step speculative fashion.
type x509Conversation struct***REMOVED******REMOVED***

var _ SpeculativeConversation = (*x509Conversation)(nil)

// FirstMessage returns the first message to be sent to the server.
func (c *x509Conversation) FirstMessage() (bsoncore.Document, error) ***REMOVED***
	return createFirstX509Message(description.Server***REMOVED******REMOVED***, ""), nil
***REMOVED***

// createFirstX509Message creates the first message for the X509 conversation.
func createFirstX509Message(desc description.Server, user string) bsoncore.Document ***REMOVED***
	elements := [][]byte***REMOVED***
		bsoncore.AppendInt32Element(nil, "authenticate", 1),
		bsoncore.AppendStringElement(nil, "mechanism", MongoDBX509),
	***REMOVED***

	// Server versions < 3.4 require the username to be included in the message. Versions >= 3.4 will extract the
	// username from the certificate.
	if desc.WireVersion != nil && desc.WireVersion.Max < 5 ***REMOVED***
		elements = append(elements, bsoncore.AppendStringElement(nil, "user", user))
	***REMOVED***

	return bsoncore.BuildDocument(nil, elements...)
***REMOVED***

// Finish implements the SpeculativeConversation interface and is a no-op because an X509 conversation only has one
// step.
func (c *x509Conversation) Finish(context.Context, *Config, bsoncore.Document) error ***REMOVED***
	return nil
***REMOVED***

// CreateSpeculativeConversation creates a speculative conversation for X509 authentication.
func (a *MongoDBX509Authenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) ***REMOVED***
	return &x509Conversation***REMOVED******REMOVED***, nil
***REMOVED***

// Auth authenticates the provided connection by conducting an X509 authentication conversation.
func (a *MongoDBX509Authenticator) Auth(ctx context.Context, cfg *Config) error ***REMOVED***
	requestDoc := createFirstX509Message(cfg.Description, a.User)
	authCmd := operation.
		NewCommand(requestDoc).
		Database("$external").
		Deployment(driver.SingleConnectionDeployment***REMOVED***cfg.Connection***REMOVED***).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err := authCmd.Execute(ctx)
	if err != nil ***REMOVED***
		return newAuthError("round trip error", err)
	***REMOVED***

	return nil
***REMOVED***
