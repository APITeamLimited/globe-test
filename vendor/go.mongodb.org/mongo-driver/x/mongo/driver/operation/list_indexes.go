// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// ListIndexes performs a listIndexes operation.
type ListIndexes struct ***REMOVED***
	batchSize  *int32
	maxTimeMS  *int64
	session    *session.Client
	clock      *session.ClusterClock
	collection string
	monitor    *event.CommandMonitor
	database   string
	deployment driver.Deployment
	selector   description.ServerSelector
	retry      *driver.RetryMode
	crypt      driver.Crypt
	serverAPI  *driver.ServerAPIOptions
	timeout    *time.Duration

	result driver.CursorResponse
***REMOVED***

// NewListIndexes constructs and returns a new ListIndexes.
func NewListIndexes() *ListIndexes ***REMOVED***
	return &ListIndexes***REMOVED******REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (li *ListIndexes) Result(opts driver.CursorOptions) (*driver.BatchCursor, error) ***REMOVED***

	clientSession := li.session

	clock := li.clock
	opts.ServerAPI = li.serverAPI
	return driver.NewBatchCursor(li.result, clientSession, clock, opts)
***REMOVED***

func (li *ListIndexes) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error

	li.result, err = driver.NewCursorResponse(info)
	return err

***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (li *ListIndexes) Execute(ctx context.Context) error ***REMOVED***
	if li.deployment == nil ***REMOVED***
		return errors.New("the ListIndexes operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         li.command,
		ProcessResponseFn: li.processResponse,

		Client:         li.session,
		Clock:          li.clock,
		CommandMonitor: li.monitor,
		Database:       li.database,
		Deployment:     li.deployment,
		Selector:       li.selector,
		Crypt:          li.crypt,
		Legacy:         driver.LegacyListIndexes,
		RetryMode:      li.retry,
		Type:           driver.Read,
		ServerAPI:      li.serverAPI,
		Timeout:        li.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (li *ListIndexes) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "listIndexes", li.collection)
	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)

	if li.batchSize != nil ***REMOVED***

		cursorDoc = bsoncore.AppendInt32Element(cursorDoc, "batchSize", *li.batchSize)
	***REMOVED***

	// Only append specified maxTimeMS if timeout is not also specified.
	if li.maxTimeMS != nil && li.timeout == nil ***REMOVED***

		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *li.maxTimeMS)
	***REMOVED***
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)
	dst = bsoncore.AppendDocumentElement(dst, "cursor", cursorDoc)

	return dst, nil
***REMOVED***

// BatchSize specifies the number of documents to return in every batch.
func (li *ListIndexes) BatchSize(batchSize int32) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.batchSize = &batchSize
	return li
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (li *ListIndexes) MaxTimeMS(maxTimeMS int64) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.maxTimeMS = &maxTimeMS
	return li
***REMOVED***

// Session sets the session for this operation.
func (li *ListIndexes) Session(session *session.Client) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.session = session
	return li
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (li *ListIndexes) ClusterClock(clock *session.ClusterClock) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.clock = clock
	return li
***REMOVED***

// Collection sets the collection that this command will run against.
func (li *ListIndexes) Collection(collection string) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.collection = collection
	return li
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (li *ListIndexes) CommandMonitor(monitor *event.CommandMonitor) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.monitor = monitor
	return li
***REMOVED***

// Database sets the database to run this operation against.
func (li *ListIndexes) Database(database string) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.database = database
	return li
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (li *ListIndexes) Deployment(deployment driver.Deployment) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.deployment = deployment
	return li
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (li *ListIndexes) ServerSelector(selector description.ServerSelector) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.selector = selector
	return li
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (li *ListIndexes) Retry(retry driver.RetryMode) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.retry = &retry
	return li
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (li *ListIndexes) Crypt(crypt driver.Crypt) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.crypt = crypt
	return li
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (li *ListIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.serverAPI = serverAPI
	return li
***REMOVED***

// Timeout sets the timeout for this operation.
func (li *ListIndexes) Timeout(timeout *time.Duration) *ListIndexes ***REMOVED***
	if li == nil ***REMOVED***
		li = new(ListIndexes)
	***REMOVED***

	li.timeout = timeout
	return li
***REMOVED***
