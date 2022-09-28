// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// SaslClient is the client piece of a sasl conversation.
type SaslClient interface ***REMOVED***
	Start() (string, []byte, error)
	Next(challenge []byte) ([]byte, error)
	Completed() bool
***REMOVED***

// SaslClientCloser is a SaslClient that has resources to clean up.
type SaslClientCloser interface ***REMOVED***
	SaslClient
	Close()
***REMOVED***

// ExtraOptionsSaslClient is a SaslClient that appends options to the saslStart command.
type ExtraOptionsSaslClient interface ***REMOVED***
	StartCommandOptions() bsoncore.Document
***REMOVED***

// saslConversation represents a SASL conversation. This type implements the SpeculativeConversation interface so the
// conversation can be executed in multi-step speculative fashion.
type saslConversation struct ***REMOVED***
	client      SaslClient
	source      string
	mechanism   string
	speculative bool
***REMOVED***

var _ SpeculativeConversation = (*saslConversation)(nil)

func newSaslConversation(client SaslClient, source string, speculative bool) *saslConversation ***REMOVED***
	authSource := source
	if authSource == "" ***REMOVED***
		authSource = defaultAuthDB
	***REMOVED***
	return &saslConversation***REMOVED***
		client:      client,
		source:      authSource,
		speculative: speculative,
	***REMOVED***
***REMOVED***

// FirstMessage returns the first message to be sent to the server. This message contains a "db" field so it can be used
// for speculative authentication.
func (sc *saslConversation) FirstMessage() (bsoncore.Document, error) ***REMOVED***
	var payload []byte
	var err error
	sc.mechanism, payload, err = sc.client.Start()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	saslCmdElements := [][]byte***REMOVED***
		bsoncore.AppendInt32Element(nil, "saslStart", 1),
		bsoncore.AppendStringElement(nil, "mechanism", sc.mechanism),
		bsoncore.AppendBinaryElement(nil, "payload", 0x00, payload),
	***REMOVED***
	if sc.speculative ***REMOVED***
		// The "db" field is only appended for speculative auth because the hello command is executed against admin
		// so this is needed to tell the server the user's auth source. For a non-speculative attempt, the SASL commands
		// will be executed against the auth source.
		saslCmdElements = append(saslCmdElements, bsoncore.AppendStringElement(nil, "db", sc.source))
	***REMOVED***
	if extraOptionsClient, ok := sc.client.(ExtraOptionsSaslClient); ok ***REMOVED***
		optionsDoc := extraOptionsClient.StartCommandOptions()
		saslCmdElements = append(saslCmdElements, bsoncore.AppendDocumentElement(nil, "options", optionsDoc))
	***REMOVED***

	return bsoncore.BuildDocumentFromElements(nil, saslCmdElements...), nil
***REMOVED***

type saslResponse struct ***REMOVED***
	ConversationID int    `bson:"conversationId"`
	Code           int    `bson:"code"`
	Done           bool   `bson:"done"`
	Payload        []byte `bson:"payload"`
***REMOVED***

// Finish completes the conversation based on the first server response to authenticate the given connection.
func (sc *saslConversation) Finish(ctx context.Context, cfg *Config, firstResponse bsoncore.Document) error ***REMOVED***
	if closer, ok := sc.client.(SaslClientCloser); ok ***REMOVED***
		defer closer.Close()
	***REMOVED***

	var saslResp saslResponse
	err := bson.Unmarshal(firstResponse, &saslResp)
	if err != nil ***REMOVED***
		fullErr := fmt.Errorf("unmarshal error: %v", err)
		return newError(fullErr, sc.mechanism)
	***REMOVED***

	cid := saslResp.ConversationID
	var payload []byte
	var rdr bsoncore.Document
	for ***REMOVED***
		if saslResp.Code != 0 ***REMOVED***
			return newError(err, sc.mechanism)
		***REMOVED***

		if saslResp.Done && sc.client.Completed() ***REMOVED***
			return nil
		***REMOVED***

		payload, err = sc.client.Next(saslResp.Payload)
		if err != nil ***REMOVED***
			return newError(err, sc.mechanism)
		***REMOVED***

		if saslResp.Done && sc.client.Completed() ***REMOVED***
			return nil
		***REMOVED***

		doc := bsoncore.BuildDocumentFromElements(nil,
			bsoncore.AppendInt32Element(nil, "saslContinue", 1),
			bsoncore.AppendInt32Element(nil, "conversationId", int32(cid)),
			bsoncore.AppendBinaryElement(nil, "payload", 0x00, payload),
		)
		saslContinueCmd := operation.NewCommand(doc).
			Database(sc.source).
			Deployment(driver.SingleConnectionDeployment***REMOVED***cfg.Connection***REMOVED***).
			ClusterClock(cfg.ClusterClock).
			ServerAPI(cfg.ServerAPI)

		err = saslContinueCmd.Execute(ctx)
		if err != nil ***REMOVED***
			return newError(err, sc.mechanism)
		***REMOVED***
		rdr = saslContinueCmd.Result()

		err = bson.Unmarshal(rdr, &saslResp)
		if err != nil ***REMOVED***
			fullErr := fmt.Errorf("unmarshal error: %v", err)
			return newError(fullErr, sc.mechanism)
		***REMOVED***
	***REMOVED***
***REMOVED***

// ConductSaslConversation runs a full SASL conversation to authenticate the given connection.
func ConductSaslConversation(ctx context.Context, cfg *Config, authSource string, client SaslClient) error ***REMOVED***
	// Create a non-speculative SASL conversation.
	conversation := newSaslConversation(client, authSource, false)

	saslStartDoc, err := conversation.FirstMessage()
	if err != nil ***REMOVED***
		return newError(err, conversation.mechanism)
	***REMOVED***
	saslStartCmd := operation.NewCommand(saslStartDoc).
		Database(authSource).
		Deployment(driver.SingleConnectionDeployment***REMOVED***cfg.Connection***REMOVED***).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	if err := saslStartCmd.Execute(ctx); err != nil ***REMOVED***
		return newError(err, conversation.mechanism)
	***REMOVED***

	return conversation.Finish(ctx, cfg, saslStartCmd.Result())
***REMOVED***
