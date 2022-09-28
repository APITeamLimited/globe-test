// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// DropIndexes performs an dropIndexes operation.
type DropIndexes struct ***REMOVED***
	index        *string
	maxTimeMS    *int64
	session      *session.Client
	clock        *session.ClusterClock
	collection   string
	monitor      *event.CommandMonitor
	crypt        driver.Crypt
	database     string
	deployment   driver.Deployment
	selector     description.ServerSelector
	writeConcern *writeconcern.WriteConcern
	result       DropIndexesResult
	serverAPI    *driver.ServerAPIOptions
	timeout      *time.Duration
***REMOVED***

// DropIndexesResult represents a dropIndexes result returned by the server.
type DropIndexesResult struct ***REMOVED***
	// Number of indexes that existed before the drop was executed.
	NIndexesWas int32
***REMOVED***

func buildDropIndexesResult(response bsoncore.Document) (DropIndexesResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return DropIndexesResult***REMOVED******REMOVED***, err
	***REMOVED***
	dir := DropIndexesResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "nIndexesWas":
			var ok bool
			dir.NIndexesWas, ok = element.Value().AsInt32OK()
			if !ok ***REMOVED***
				return dir, fmt.Errorf("response field 'nIndexesWas' is type int32, but received BSON type %s", element.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dir, nil
***REMOVED***

// NewDropIndexes constructs and returns a new DropIndexes.
func NewDropIndexes(index string) *DropIndexes ***REMOVED***
	return &DropIndexes***REMOVED***
		index: &index,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (di *DropIndexes) Result() DropIndexesResult ***REMOVED*** return di.result ***REMOVED***

func (di *DropIndexes) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error
	di.result, err = buildDropIndexesResult(info.ServerResponse)
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (di *DropIndexes) Execute(ctx context.Context) error ***REMOVED***
	if di.deployment == nil ***REMOVED***
		return errors.New("the DropIndexes operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         di.command,
		ProcessResponseFn: di.processResponse,
		Client:            di.session,
		Clock:             di.clock,
		CommandMonitor:    di.monitor,
		Crypt:             di.crypt,
		Database:          di.database,
		Deployment:        di.deployment,
		Selector:          di.selector,
		WriteConcern:      di.writeConcern,
		ServerAPI:         di.serverAPI,
		Timeout:           di.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (di *DropIndexes) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "dropIndexes", di.collection)
	if di.index != nil ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "index", *di.index)
	***REMOVED***
	// Only append specified maxTimeMS if timeout is not also specified.
	if di.maxTimeMS != nil && di.timeout == nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *di.maxTimeMS)
	***REMOVED***
	return dst, nil
***REMOVED***

// Index specifies the name of the index to drop. If '*' is specified, all indexes will be dropped.
//
func (di *DropIndexes) Index(index string) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.index = &index
	return di
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (di *DropIndexes) MaxTimeMS(maxTimeMS int64) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.maxTimeMS = &maxTimeMS
	return di
***REMOVED***

// Session sets the session for this operation.
func (di *DropIndexes) Session(session *session.Client) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.session = session
	return di
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (di *DropIndexes) ClusterClock(clock *session.ClusterClock) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.clock = clock
	return di
***REMOVED***

// Collection sets the collection that this command will run against.
func (di *DropIndexes) Collection(collection string) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.collection = collection
	return di
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (di *DropIndexes) CommandMonitor(monitor *event.CommandMonitor) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.monitor = monitor
	return di
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (di *DropIndexes) Crypt(crypt driver.Crypt) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.crypt = crypt
	return di
***REMOVED***

// Database sets the database to run this operation against.
func (di *DropIndexes) Database(database string) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.database = database
	return di
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (di *DropIndexes) Deployment(deployment driver.Deployment) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.deployment = deployment
	return di
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (di *DropIndexes) ServerSelector(selector description.ServerSelector) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.selector = selector
	return di
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (di *DropIndexes) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.writeConcern = writeConcern
	return di
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (di *DropIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.serverAPI = serverAPI
	return di
***REMOVED***

// Timeout sets the timeout for this operation.
func (di *DropIndexes) Timeout(timeout *time.Duration) *DropIndexes ***REMOVED***
	if di == nil ***REMOVED***
		di = new(DropIndexes)
	***REMOVED***

	di.timeout = timeout
	return di
***REMOVED***
