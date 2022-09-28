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

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// CreateIndexes performs a createIndexes operation.
type CreateIndexes struct ***REMOVED***
	commitQuorum bsoncore.Value
	indexes      bsoncore.Document
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
	result       CreateIndexesResult
	serverAPI    *driver.ServerAPIOptions
	timeout      *time.Duration
***REMOVED***

// CreateIndexesResult represents a createIndexes result returned by the server.
type CreateIndexesResult struct ***REMOVED***
	// If the collection was created automatically.
	CreatedCollectionAutomatically bool
	// The number of indexes existing after this command.
	IndexesAfter int32
	// The number of indexes existing before this command.
	IndexesBefore int32
***REMOVED***

func buildCreateIndexesResult(response bsoncore.Document) (CreateIndexesResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return CreateIndexesResult***REMOVED******REMOVED***, err
	***REMOVED***
	cir := CreateIndexesResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "createdCollectionAutomatically":
			var ok bool
			cir.CreatedCollectionAutomatically, ok = element.Value().BooleanOK()
			if !ok ***REMOVED***
				return cir, fmt.Errorf("response field 'createdCollectionAutomatically' is type bool, but received BSON type %s", element.Value().Type)
			***REMOVED***
		case "indexesAfter":
			var ok bool
			cir.IndexesAfter, ok = element.Value().AsInt32OK()
			if !ok ***REMOVED***
				return cir, fmt.Errorf("response field 'indexesAfter' is type int32, but received BSON type %s", element.Value().Type)
			***REMOVED***
		case "indexesBefore":
			var ok bool
			cir.IndexesBefore, ok = element.Value().AsInt32OK()
			if !ok ***REMOVED***
				return cir, fmt.Errorf("response field 'indexesBefore' is type int32, but received BSON type %s", element.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return cir, nil
***REMOVED***

// NewCreateIndexes constructs and returns a new CreateIndexes.
func NewCreateIndexes(indexes bsoncore.Document) *CreateIndexes ***REMOVED***
	return &CreateIndexes***REMOVED***
		indexes: indexes,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (ci *CreateIndexes) Result() CreateIndexesResult ***REMOVED*** return ci.result ***REMOVED***

func (ci *CreateIndexes) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error
	ci.result, err = buildCreateIndexesResult(info.ServerResponse)
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (ci *CreateIndexes) Execute(ctx context.Context) error ***REMOVED***
	if ci.deployment == nil ***REMOVED***
		return errors.New("the CreateIndexes operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         ci.command,
		ProcessResponseFn: ci.processResponse,
		Client:            ci.session,
		Clock:             ci.clock,
		CommandMonitor:    ci.monitor,
		Crypt:             ci.crypt,
		Database:          ci.database,
		Deployment:        ci.deployment,
		Selector:          ci.selector,
		WriteConcern:      ci.writeConcern,
		ServerAPI:         ci.serverAPI,
		Timeout:           ci.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (ci *CreateIndexes) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "createIndexes", ci.collection)
	if ci.commitQuorum.Type != bsontype.Type(0) ***REMOVED***
		if desc.WireVersion == nil || !desc.WireVersion.Includes(9) ***REMOVED***
			return nil, errors.New("the 'commitQuorum' command parameter requires a minimum server wire version of 9")
		***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "commitQuorum", ci.commitQuorum)
	***REMOVED***
	if ci.indexes != nil ***REMOVED***
		dst = bsoncore.AppendArrayElement(dst, "indexes", ci.indexes)
	***REMOVED***
	// Only append specified maxTimeMS if timeout is not also specified.
	if ci.maxTimeMS != nil && ci.timeout == nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *ci.maxTimeMS)
	***REMOVED***
	return dst, nil
***REMOVED***

// CommitQuorum specifies the number of data-bearing members of a replica set, including the primary, that must
// complete the index builds successfully before the primary marks the indexes as ready. This should either be a
// string or int32 value.
func (ci *CreateIndexes) CommitQuorum(commitQuorum bsoncore.Value) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.commitQuorum = commitQuorum
	return ci
***REMOVED***

// Indexes specifies an array containing index specification documents for the indexes being created.
func (ci *CreateIndexes) Indexes(indexes bsoncore.Document) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.indexes = indexes
	return ci
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (ci *CreateIndexes) MaxTimeMS(maxTimeMS int64) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.maxTimeMS = &maxTimeMS
	return ci
***REMOVED***

// Session sets the session for this operation.
func (ci *CreateIndexes) Session(session *session.Client) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.session = session
	return ci
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (ci *CreateIndexes) ClusterClock(clock *session.ClusterClock) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.clock = clock
	return ci
***REMOVED***

// Collection sets the collection that this command will run against.
func (ci *CreateIndexes) Collection(collection string) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.collection = collection
	return ci
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (ci *CreateIndexes) CommandMonitor(monitor *event.CommandMonitor) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.monitor = monitor
	return ci
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (ci *CreateIndexes) Crypt(crypt driver.Crypt) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.crypt = crypt
	return ci
***REMOVED***

// Database sets the database to run this operation against.
func (ci *CreateIndexes) Database(database string) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.database = database
	return ci
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (ci *CreateIndexes) Deployment(deployment driver.Deployment) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.deployment = deployment
	return ci
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (ci *CreateIndexes) ServerSelector(selector description.ServerSelector) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.selector = selector
	return ci
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (ci *CreateIndexes) WriteConcern(writeConcern *writeconcern.WriteConcern) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.writeConcern = writeConcern
	return ci
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (ci *CreateIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.serverAPI = serverAPI
	return ci
***REMOVED***

// Timeout sets the timeout for this operation.
func (ci *CreateIndexes) Timeout(timeout *time.Duration) *CreateIndexes ***REMOVED***
	if ci == nil ***REMOVED***
		ci = new(CreateIndexes)
	***REMOVED***

	ci.timeout = timeout
	return ci
***REMOVED***
