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

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Aggregate represents an aggregate operation.
type Aggregate struct ***REMOVED***
	allowDiskUse             *bool
	batchSize                *int32
	bypassDocumentValidation *bool
	collation                bsoncore.Document
	comment                  *string
	hint                     bsoncore.Value
	maxTimeMS                *int64
	pipeline                 bsoncore.Document
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	readConcern              *readconcern.ReadConcern
	readPreference           *readpref.ReadPref
	retry                    *driver.RetryMode
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	crypt                    driver.Crypt
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	hasOutputStage           bool
	customOptions            map[string]bsoncore.Value
	timeout                  *time.Duration

	result driver.CursorResponse
***REMOVED***

// NewAggregate constructs and returns a new Aggregate.
func NewAggregate(pipeline bsoncore.Document) *Aggregate ***REMOVED***
	return &Aggregate***REMOVED***
		pipeline: pipeline,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (a *Aggregate) Result(opts driver.CursorOptions) (*driver.BatchCursor, error) ***REMOVED***

	clientSession := a.session

	clock := a.clock
	opts.ServerAPI = a.serverAPI
	return driver.NewBatchCursor(a.result, clientSession, clock, opts)
***REMOVED***

// ResultCursorResponse returns the underlying CursorResponse result of executing this
// operation.
func (a *Aggregate) ResultCursorResponse() driver.CursorResponse ***REMOVED***
	return a.result
***REMOVED***

func (a *Aggregate) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error

	a.result, err = driver.NewCursorResponse(info)
	return err

***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (a *Aggregate) Execute(ctx context.Context) error ***REMOVED***
	if a.deployment == nil ***REMOVED***
		return errors.New("the Aggregate operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         a.command,
		ProcessResponseFn: a.processResponse,

		Client:                         a.session,
		Clock:                          a.clock,
		CommandMonitor:                 a.monitor,
		Database:                       a.database,
		Deployment:                     a.deployment,
		ReadConcern:                    a.readConcern,
		ReadPreference:                 a.readPreference,
		Type:                           driver.Read,
		RetryMode:                      a.retry,
		Selector:                       a.selector,
		WriteConcern:                   a.writeConcern,
		Crypt:                          a.crypt,
		MinimumWriteConcernWireVersion: 5,
		ServerAPI:                      a.serverAPI,
		IsOutputAggregate:              a.hasOutputStage,
		Timeout:                        a.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (a *Aggregate) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	header := bsoncore.Value***REMOVED***Type: bsontype.String, Data: bsoncore.AppendString(nil, a.collection)***REMOVED***
	if a.collection == "" ***REMOVED***
		header = bsoncore.Value***REMOVED***Type: bsontype.Int32, Data: []byte***REMOVED***0x01, 0x00, 0x00, 0x00***REMOVED******REMOVED***
	***REMOVED***
	dst = bsoncore.AppendValueElement(dst, "aggregate", header)

	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)
	if a.allowDiskUse != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "allowDiskUse", *a.allowDiskUse)
	***REMOVED***
	if a.batchSize != nil ***REMOVED***
		cursorDoc = bsoncore.AppendInt32Element(cursorDoc, "batchSize", *a.batchSize)
	***REMOVED***
	if a.bypassDocumentValidation != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *a.bypassDocumentValidation)
	***REMOVED***
	if a.collation != nil ***REMOVED***

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) ***REMOVED***
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "collation", a.collation)
	***REMOVED***
	if a.comment != nil ***REMOVED***

		dst = bsoncore.AppendStringElement(dst, "comment", *a.comment)
	***REMOVED***
	if a.hint.Type != bsontype.Type(0) ***REMOVED***

		dst = bsoncore.AppendValueElement(dst, "hint", a.hint)
	***REMOVED***

	// Only append specified maxTimeMS if timeout is not also specified.
	if a.maxTimeMS != nil && a.timeout == nil ***REMOVED***

		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *a.maxTimeMS)
	***REMOVED***
	if a.pipeline != nil ***REMOVED***

		dst = bsoncore.AppendArrayElement(dst, "pipeline", a.pipeline)
	***REMOVED***
	if a.let != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "let", a.let)
	***REMOVED***
	for optionName, optionValue := range a.customOptions ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, optionName, optionValue)
	***REMOVED***
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)
	dst = bsoncore.AppendDocumentElement(dst, "cursor", cursorDoc)

	return dst, nil
***REMOVED***

// AllowDiskUse enables writing to temporary files. When true, aggregation stages can write to the dbPath/_tmp directory.
func (a *Aggregate) AllowDiskUse(allowDiskUse bool) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.allowDiskUse = &allowDiskUse
	return a
***REMOVED***

// BatchSize specifies the number of documents to return in every batch.
func (a *Aggregate) BatchSize(batchSize int32) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.batchSize = &batchSize
	return a
***REMOVED***

// BypassDocumentValidation allows the write to opt-out of document level validation. This only applies when the $out stage is specified.
func (a *Aggregate) BypassDocumentValidation(bypassDocumentValidation bool) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.bypassDocumentValidation = &bypassDocumentValidation
	return a
***REMOVED***

// Collation specifies a collation. This option is only valid for server versions 3.4 and above.
func (a *Aggregate) Collation(collation bsoncore.Document) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.collation = collation
	return a
***REMOVED***

// Comment specifies an arbitrary string to help trace the operation through the database profiler, currentOp, and logs.
func (a *Aggregate) Comment(comment string) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.comment = &comment
	return a
***REMOVED***

// Hint specifies the index to use.
func (a *Aggregate) Hint(hint bsoncore.Value) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.hint = hint
	return a
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (a *Aggregate) MaxTimeMS(maxTimeMS int64) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.maxTimeMS = &maxTimeMS
	return a
***REMOVED***

// Pipeline determines how data is transformed for an aggregation.
func (a *Aggregate) Pipeline(pipeline bsoncore.Document) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.pipeline = pipeline
	return a
***REMOVED***

// Session sets the session for this operation.
func (a *Aggregate) Session(session *session.Client) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.session = session
	return a
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (a *Aggregate) ClusterClock(clock *session.ClusterClock) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.clock = clock
	return a
***REMOVED***

// Collection sets the collection that this command will run against.
func (a *Aggregate) Collection(collection string) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.collection = collection
	return a
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (a *Aggregate) CommandMonitor(monitor *event.CommandMonitor) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.monitor = monitor
	return a
***REMOVED***

// Database sets the database to run this operation against.
func (a *Aggregate) Database(database string) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.database = database
	return a
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (a *Aggregate) Deployment(deployment driver.Deployment) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.deployment = deployment
	return a
***REMOVED***

// ReadConcern specifies the read concern for this operation.
func (a *Aggregate) ReadConcern(readConcern *readconcern.ReadConcern) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.readConcern = readConcern
	return a
***REMOVED***

// ReadPreference set the read preference used with this operation.
func (a *Aggregate) ReadPreference(readPreference *readpref.ReadPref) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.readPreference = readPreference
	return a
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (a *Aggregate) ServerSelector(selector description.ServerSelector) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.selector = selector
	return a
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (a *Aggregate) WriteConcern(writeConcern *writeconcern.WriteConcern) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.writeConcern = writeConcern
	return a
***REMOVED***

// Retry enables retryable writes for this operation. Retries are not handled automatically,
// instead a boolean is returned from Execute and SelectAndExecute that indicates if the
// operation can be retried. Retrying is handled by calling RetryExecute.
func (a *Aggregate) Retry(retry driver.RetryMode) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.retry = &retry
	return a
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (a *Aggregate) Crypt(crypt driver.Crypt) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.crypt = crypt
	return a
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (a *Aggregate) ServerAPI(serverAPI *driver.ServerAPIOptions) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.serverAPI = serverAPI
	return a
***REMOVED***

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (a *Aggregate) Let(let bsoncore.Document) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.let = let
	return a
***REMOVED***

// HasOutputStage specifies whether the aggregate contains an output stage. Used in determining when to
// append read preference at the operation level.
func (a *Aggregate) HasOutputStage(hos bool) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.hasOutputStage = hos
	return a
***REMOVED***

// CustomOptions specifies extra options to use in the aggregate command.
func (a *Aggregate) CustomOptions(co map[string]bsoncore.Value) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.customOptions = co
	return a
***REMOVED***

// Timeout sets the timeout for this operation.
func (a *Aggregate) Timeout(timeout *time.Duration) *Aggregate ***REMOVED***
	if a == nil ***REMOVED***
		a = new(Aggregate)
	***REMOVED***

	a.timeout = timeout
	return a
***REMOVED***
