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
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// ListCollections performs a listCollections operation.
type ListCollections struct ***REMOVED***
	filter                bsoncore.Document
	nameOnly              *bool
	authorizedCollections *bool
	session               *session.Client
	clock                 *session.ClusterClock
	monitor               *event.CommandMonitor
	crypt                 driver.Crypt
	database              string
	deployment            driver.Deployment
	readPreference        *readpref.ReadPref
	selector              description.ServerSelector
	retry                 *driver.RetryMode
	result                driver.CursorResponse
	batchSize             *int32
	serverAPI             *driver.ServerAPIOptions
	timeout               *time.Duration
***REMOVED***

// NewListCollections constructs and returns a new ListCollections.
func NewListCollections(filter bsoncore.Document) *ListCollections ***REMOVED***
	return &ListCollections***REMOVED***
		filter: filter,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (lc *ListCollections) Result(opts driver.CursorOptions) (*driver.ListCollectionsBatchCursor, error) ***REMOVED***
	opts.ServerAPI = lc.serverAPI
	bc, err := driver.NewBatchCursor(lc.result, lc.session, lc.clock, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	desc := lc.result.Desc
	if desc.WireVersion == nil || desc.WireVersion.Max < 3 ***REMOVED***
		return driver.NewLegacyListCollectionsBatchCursor(bc)
	***REMOVED***
	return driver.NewListCollectionsBatchCursor(bc)
***REMOVED***

func (lc *ListCollections) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error
	lc.result, err = driver.NewCursorResponse(info)
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (lc *ListCollections) Execute(ctx context.Context) error ***REMOVED***
	if lc.deployment == nil ***REMOVED***
		return errors.New("the ListCollections operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         lc.command,
		ProcessResponseFn: lc.processResponse,
		RetryMode:         lc.retry,
		Type:              driver.Read,
		Client:            lc.session,
		Clock:             lc.clock,
		CommandMonitor:    lc.monitor,
		Crypt:             lc.crypt,
		Database:          lc.database,
		Deployment:        lc.deployment,
		ReadPreference:    lc.readPreference,
		Selector:          lc.selector,
		Legacy:            driver.LegacyListCollections,
		ServerAPI:         lc.serverAPI,
		Timeout:           lc.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (lc *ListCollections) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendInt32Element(dst, "listCollections", 1)
	if lc.filter != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "filter", lc.filter)
	***REMOVED***
	if lc.nameOnly != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "nameOnly", *lc.nameOnly)
	***REMOVED***
	if lc.authorizedCollections != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "authorizedCollections", *lc.authorizedCollections)
	***REMOVED***

	cursorDoc := bsoncore.NewDocumentBuilder()
	if lc.batchSize != nil ***REMOVED***
		cursorDoc.AppendInt32("batchSize", *lc.batchSize)
	***REMOVED***
	dst = bsoncore.AppendDocumentElement(dst, "cursor", cursorDoc.Build())

	return dst, nil
***REMOVED***

// Filter determines what results are returned from listCollections.
func (lc *ListCollections) Filter(filter bsoncore.Document) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.filter = filter
	return lc
***REMOVED***

// NameOnly specifies whether to only return collection names.
func (lc *ListCollections) NameOnly(nameOnly bool) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.nameOnly = &nameOnly
	return lc
***REMOVED***

// AuthorizedCollections specifies whether to only return collections the user
// is authorized to use.
func (lc *ListCollections) AuthorizedCollections(authorizedCollections bool) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.authorizedCollections = &authorizedCollections
	return lc
***REMOVED***

// Session sets the session for this operation.
func (lc *ListCollections) Session(session *session.Client) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.session = session
	return lc
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (lc *ListCollections) ClusterClock(clock *session.ClusterClock) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.clock = clock
	return lc
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (lc *ListCollections) CommandMonitor(monitor *event.CommandMonitor) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.monitor = monitor
	return lc
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (lc *ListCollections) Crypt(crypt driver.Crypt) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.crypt = crypt
	return lc
***REMOVED***

// Database sets the database to run this operation against.
func (lc *ListCollections) Database(database string) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.database = database
	return lc
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (lc *ListCollections) Deployment(deployment driver.Deployment) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.deployment = deployment
	return lc
***REMOVED***

// ReadPreference set the read preference used with this operation.
func (lc *ListCollections) ReadPreference(readPreference *readpref.ReadPref) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.readPreference = readPreference
	return lc
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (lc *ListCollections) ServerSelector(selector description.ServerSelector) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.selector = selector
	return lc
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (lc *ListCollections) Retry(retry driver.RetryMode) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.retry = &retry
	return lc
***REMOVED***

// BatchSize specifies the number of documents to return in every batch.
func (lc *ListCollections) BatchSize(batchSize int32) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.batchSize = &batchSize
	return lc
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (lc *ListCollections) ServerAPI(serverAPI *driver.ServerAPIOptions) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.serverAPI = serverAPI
	return lc
***REMOVED***

// Timeout sets the timeout for this operation.
func (lc *ListCollections) Timeout(timeout *time.Duration) *ListCollections ***REMOVED***
	if lc == nil ***REMOVED***
		lc = new(ListCollections)
	***REMOVED***

	lc.timeout = timeout
	return lc
***REMOVED***
