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

// Insert performs an insert operation.
type Insert struct ***REMOVED***
	bypassDocumentValidation *bool
	comment                  bsoncore.Value
	documents                []bsoncore.Document
	ordered                  *bool
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	crypt                    driver.Crypt
	database                 string
	deployment               driver.Deployment
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	result                   InsertResult
	serverAPI                *driver.ServerAPIOptions
	timeout                  *time.Duration
***REMOVED***

// InsertResult represents an insert result returned by the server.
type InsertResult struct ***REMOVED***
	// Number of documents successfully inserted.
	N int64
***REMOVED***

func buildInsertResult(response bsoncore.Document) (InsertResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return InsertResult***REMOVED******REMOVED***, err
	***REMOVED***
	ir := InsertResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "n":
			var ok bool
			ir.N, ok = element.Value().AsInt64OK()
			if !ok ***REMOVED***
				return ir, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ir, nil
***REMOVED***

// NewInsert constructs and returns a new Insert.
func NewInsert(documents ...bsoncore.Document) *Insert ***REMOVED***
	return &Insert***REMOVED***
		documents: documents,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (i *Insert) Result() InsertResult ***REMOVED*** return i.result ***REMOVED***

func (i *Insert) processResponse(info driver.ResponseInfo) error ***REMOVED***
	ir, err := buildInsertResult(info.ServerResponse)
	i.result.N += ir.N
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (i *Insert) Execute(ctx context.Context) error ***REMOVED***
	if i.deployment == nil ***REMOVED***
		return errors.New("the Insert operation must have a Deployment set before Execute can be called")
	***REMOVED***
	batches := &driver.Batches***REMOVED***
		Identifier: "documents",
		Documents:  i.documents,
		Ordered:    i.ordered,
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         i.command,
		ProcessResponseFn: i.processResponse,
		Batches:           batches,
		RetryMode:         i.retry,
		Type:              driver.Write,
		Client:            i.session,
		Clock:             i.clock,
		CommandMonitor:    i.monitor,
		Crypt:             i.crypt,
		Database:          i.database,
		Deployment:        i.deployment,
		Selector:          i.selector,
		WriteConcern:      i.writeConcern,
		ServerAPI:         i.serverAPI,
		Timeout:           i.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (i *Insert) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "insert", i.collection)
	if i.bypassDocumentValidation != nil && (desc.WireVersion != nil && desc.WireVersion.Includes(4)) ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *i.bypassDocumentValidation)
	***REMOVED***
	if i.comment.Type != bsontype.Type(0) ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "comment", i.comment)
	***REMOVED***
	if i.ordered != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "ordered", *i.ordered)
	***REMOVED***
	return dst, nil
***REMOVED***

// BypassDocumentValidation allows the operation to opt-out of document level validation. Valid
// for server versions >= 3.2. For servers < 3.2, this setting is ignored.
func (i *Insert) BypassDocumentValidation(bypassDocumentValidation bool) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.bypassDocumentValidation = &bypassDocumentValidation
	return i
***REMOVED***

// Comment sets a value to help trace an operation.
func (i *Insert) Comment(comment bsoncore.Value) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.comment = comment
	return i
***REMOVED***

// Documents adds documents to this operation that will be inserted when this operation is
// executed.
func (i *Insert) Documents(documents ...bsoncore.Document) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.documents = documents
	return i
***REMOVED***

// Ordered sets ordered. If true, when a write fails, the operation will return the error, when
// false write failures do not stop execution of the operation.
func (i *Insert) Ordered(ordered bool) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.ordered = &ordered
	return i
***REMOVED***

// Session sets the session for this operation.
func (i *Insert) Session(session *session.Client) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.session = session
	return i
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (i *Insert) ClusterClock(clock *session.ClusterClock) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.clock = clock
	return i
***REMOVED***

// Collection sets the collection that this command will run against.
func (i *Insert) Collection(collection string) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.collection = collection
	return i
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (i *Insert) CommandMonitor(monitor *event.CommandMonitor) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.monitor = monitor
	return i
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (i *Insert) Crypt(crypt driver.Crypt) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.crypt = crypt
	return i
***REMOVED***

// Database sets the database to run this operation against.
func (i *Insert) Database(database string) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.database = database
	return i
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (i *Insert) Deployment(deployment driver.Deployment) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.deployment = deployment
	return i
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (i *Insert) ServerSelector(selector description.ServerSelector) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.selector = selector
	return i
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (i *Insert) WriteConcern(writeConcern *writeconcern.WriteConcern) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.writeConcern = writeConcern
	return i
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (i *Insert) Retry(retry driver.RetryMode) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.retry = &retry
	return i
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (i *Insert) ServerAPI(serverAPI *driver.ServerAPIOptions) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.serverAPI = serverAPI
	return i
***REMOVED***

// Timeout sets the timeout for this operation.
func (i *Insert) Timeout(timeout *time.Duration) *Insert ***REMOVED***
	if i == nil ***REMOVED***
		i = new(Insert)
	***REMOVED***

	i.timeout = timeout
	return i
***REMOVED***
