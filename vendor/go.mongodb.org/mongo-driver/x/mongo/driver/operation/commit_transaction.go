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

// CommitTransaction attempts to commit a transaction.
type CommitTransaction struct ***REMOVED***
	maxTimeMS     *int64
	recoveryToken bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
	serverAPI     *driver.ServerAPIOptions
***REMOVED***

// NewCommitTransaction constructs and returns a new CommitTransaction.
func NewCommitTransaction() *CommitTransaction ***REMOVED***
	return &CommitTransaction***REMOVED******REMOVED***
***REMOVED***

func (ct *CommitTransaction) processResponse(driver.ResponseInfo) error ***REMOVED***
	var err error
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (ct *CommitTransaction) Execute(ctx context.Context) error ***REMOVED***
	if ct.deployment == nil ***REMOVED***
		return errors.New("the CommitTransaction operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         ct.command,
		ProcessResponseFn: ct.processResponse,
		RetryMode:         ct.retry,
		Type:              driver.Write,
		Client:            ct.session,
		Clock:             ct.clock,
		CommandMonitor:    ct.monitor,
		Crypt:             ct.crypt,
		Database:          ct.database,
		Deployment:        ct.deployment,
		Selector:          ct.selector,
		WriteConcern:      ct.writeConcern,
		ServerAPI:         ct.serverAPI,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (ct *CommitTransaction) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***

	dst = bsoncore.AppendInt32Element(dst, "commitTransaction", 1)
	if ct.maxTimeMS != nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *ct.maxTimeMS)
	***REMOVED***
	if ct.recoveryToken != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "recoveryToken", ct.recoveryToken)
	***REMOVED***
	return dst, nil
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (ct *CommitTransaction) MaxTimeMS(maxTimeMS int64) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.maxTimeMS = &maxTimeMS
	return ct
***REMOVED***

// RecoveryToken sets the recovery token to use when committing or aborting a sharded transaction.
func (ct *CommitTransaction) RecoveryToken(recoveryToken bsoncore.Document) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.recoveryToken = recoveryToken
	return ct
***REMOVED***

// Session sets the session for this operation.
func (ct *CommitTransaction) Session(session *session.Client) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.session = session
	return ct
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (ct *CommitTransaction) ClusterClock(clock *session.ClusterClock) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.clock = clock
	return ct
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (ct *CommitTransaction) CommandMonitor(monitor *event.CommandMonitor) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.monitor = monitor
	return ct
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (ct *CommitTransaction) Crypt(crypt driver.Crypt) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.crypt = crypt
	return ct
***REMOVED***

// Database sets the database to run this operation against.
func (ct *CommitTransaction) Database(database string) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.database = database
	return ct
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (ct *CommitTransaction) Deployment(deployment driver.Deployment) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.deployment = deployment
	return ct
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (ct *CommitTransaction) ServerSelector(selector description.ServerSelector) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.selector = selector
	return ct
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (ct *CommitTransaction) WriteConcern(writeConcern *writeconcern.WriteConcern) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.writeConcern = writeConcern
	return ct
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (ct *CommitTransaction) Retry(retry driver.RetryMode) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.retry = &retry
	return ct
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (ct *CommitTransaction) ServerAPI(serverAPI *driver.ServerAPIOptions) *CommitTransaction ***REMOVED***
	if ct == nil ***REMOVED***
		ct = new(CommitTransaction)
	***REMOVED***

	ct.serverAPI = serverAPI
	return ct
***REMOVED***
