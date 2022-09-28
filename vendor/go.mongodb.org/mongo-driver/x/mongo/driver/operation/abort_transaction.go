// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// AbortTransaction performs an abortTransaction operation.
type AbortTransaction struct ***REMOVED***
	recoveryToken bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
	serverAPI     *driver.ServerAPIOptions
***REMOVED***

// NewAbortTransaction constructs and returns a new AbortTransaction.
func NewAbortTransaction() *AbortTransaction ***REMOVED***
	return &AbortTransaction***REMOVED******REMOVED***
***REMOVED***

func (at *AbortTransaction) processResponse(driver.ResponseInfo) error ***REMOVED***
	var err error
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (at *AbortTransaction) Execute(ctx context.Context) error ***REMOVED***
	if at.deployment == nil ***REMOVED***
		return errors.New("the AbortTransaction operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         at.command,
		ProcessResponseFn: at.processResponse,
		RetryMode:         at.retry,
		Type:              driver.Write,
		Client:            at.session,
		Clock:             at.clock,
		CommandMonitor:    at.monitor,
		Crypt:             at.crypt,
		Database:          at.database,
		Deployment:        at.deployment,
		Selector:          at.selector,
		WriteConcern:      at.writeConcern,
		ServerAPI:         at.serverAPI,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (at *AbortTransaction) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***

	dst = bsoncore.AppendInt32Element(dst, "abortTransaction", 1)
	if at.recoveryToken != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "recoveryToken", at.recoveryToken)
	***REMOVED***
	return dst, nil
***REMOVED***

// RecoveryToken sets the recovery token to use when committing or aborting a sharded transaction.
func (at *AbortTransaction) RecoveryToken(recoveryToken bsoncore.Document) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.recoveryToken = recoveryToken
	return at
***REMOVED***

// Session sets the session for this operation.
func (at *AbortTransaction) Session(session *session.Client) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.session = session
	return at
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (at *AbortTransaction) ClusterClock(clock *session.ClusterClock) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.clock = clock
	return at
***REMOVED***

// Collection sets the collection that this command will run against.
func (at *AbortTransaction) Collection(collection string) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.collection = collection
	return at
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (at *AbortTransaction) CommandMonitor(monitor *event.CommandMonitor) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.monitor = monitor
	return at
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (at *AbortTransaction) Crypt(crypt driver.Crypt) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.crypt = crypt
	return at
***REMOVED***

// Database sets the database to run this operation against.
func (at *AbortTransaction) Database(database string) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.database = database
	return at
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (at *AbortTransaction) Deployment(deployment driver.Deployment) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.deployment = deployment
	return at
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (at *AbortTransaction) ServerSelector(selector description.ServerSelector) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.selector = selector
	return at
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (at *AbortTransaction) WriteConcern(writeConcern *writeconcern.WriteConcern) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.writeConcern = writeConcern
	return at
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (at *AbortTransaction) Retry(retry driver.RetryMode) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.retry = &retry
	return at
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (at *AbortTransaction) ServerAPI(serverAPI *driver.ServerAPIOptions) *AbortTransaction ***REMOVED***
	if at == nil ***REMOVED***
		at = new(AbortTransaction)
	***REMOVED***

	at.serverAPI = serverAPI
	return at
***REMOVED***
