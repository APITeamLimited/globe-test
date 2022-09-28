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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// FindAndModify performs a findAndModify operation.
type FindAndModify struct ***REMOVED***
	arrayFilters             bsoncore.Document
	bypassDocumentValidation *bool
	collation                bsoncore.Document
	comment                  bsoncore.Value
	fields                   bsoncore.Document
	maxTimeMS                *int64
	newDocument              *bool
	query                    bsoncore.Document
	remove                   *bool
	sort                     bsoncore.Document
	update                   bsoncore.Value
	upsert                   *bool
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	crypt                    driver.Crypt
	hint                     bsoncore.Value
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	timeout                  *time.Duration

	result FindAndModifyResult
***REMOVED***

// LastErrorObject represents information about updates and upserts returned by the server.
type LastErrorObject struct ***REMOVED***
	// True if an update modified an existing document
	UpdatedExisting bool
	// Object ID of the upserted document.
	Upserted interface***REMOVED******REMOVED***
***REMOVED***

// FindAndModifyResult represents a findAndModify result returned by the server.
type FindAndModifyResult struct ***REMOVED***
	// Either the old or modified document, depending on the value of the new parameter.
	Value bsoncore.Document
	// Contains information about updates and upserts.
	LastErrorObject LastErrorObject
***REMOVED***

func buildFindAndModifyResult(response bsoncore.Document) (FindAndModifyResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return FindAndModifyResult***REMOVED******REMOVED***, err
	***REMOVED***
	famr := FindAndModifyResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "value":
			var ok bool
			famr.Value, ok = element.Value().DocumentOK()

			// The 'value' field returned by a FindAndModify can be null in the case that no document was found.
			if element.Value().Type != bsontype.Null && !ok ***REMOVED***
				return famr, fmt.Errorf("response field 'value' is type document or null, but received BSON type %s", element.Value().Type)
			***REMOVED***
		case "lastErrorObject":
			valDoc, ok := element.Value().DocumentOK()
			if !ok ***REMOVED***
				return famr, fmt.Errorf("response field 'lastErrorObject' is type document, but received BSON type %s", element.Value().Type)
			***REMOVED***

			var leo LastErrorObject
			if err = bson.Unmarshal(valDoc, &leo); err != nil ***REMOVED***
				return famr, err
			***REMOVED***
			famr.LastErrorObject = leo
		***REMOVED***
	***REMOVED***
	return famr, nil
***REMOVED***

// NewFindAndModify constructs and returns a new FindAndModify.
func NewFindAndModify(query bsoncore.Document) *FindAndModify ***REMOVED***
	return &FindAndModify***REMOVED***
		query: query,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (fam *FindAndModify) Result() FindAndModifyResult ***REMOVED*** return fam.result ***REMOVED***

func (fam *FindAndModify) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error

	fam.result, err = buildFindAndModifyResult(info.ServerResponse)
	return err

***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (fam *FindAndModify) Execute(ctx context.Context) error ***REMOVED***
	if fam.deployment == nil ***REMOVED***
		return errors.New("the FindAndModify operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         fam.command,
		ProcessResponseFn: fam.processResponse,

		RetryMode:      fam.retry,
		Type:           driver.Write,
		Client:         fam.session,
		Clock:          fam.clock,
		CommandMonitor: fam.monitor,
		Database:       fam.database,
		Deployment:     fam.deployment,
		Selector:       fam.selector,
		WriteConcern:   fam.writeConcern,
		Crypt:          fam.crypt,
		ServerAPI:      fam.serverAPI,
		Timeout:        fam.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (fam *FindAndModify) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "findAndModify", fam.collection)
	if fam.arrayFilters != nil ***REMOVED***

		if desc.WireVersion == nil || !desc.WireVersion.Includes(6) ***REMOVED***
			return nil, errors.New("the 'arrayFilters' command parameter requires a minimum server wire version of 6")
		***REMOVED***
		dst = bsoncore.AppendArrayElement(dst, "arrayFilters", fam.arrayFilters)
	***REMOVED***
	if fam.bypassDocumentValidation != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *fam.bypassDocumentValidation)
	***REMOVED***
	if fam.collation != nil ***REMOVED***

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) ***REMOVED***
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "collation", fam.collation)
	***REMOVED***
	if fam.comment.Type != bsontype.Type(0) ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "comment", fam.comment)
	***REMOVED***
	if fam.fields != nil ***REMOVED***

		dst = bsoncore.AppendDocumentElement(dst, "fields", fam.fields)
	***REMOVED***

	// Only append specified maxTimeMS if timeout is not also specified.
	if fam.maxTimeMS != nil && fam.timeout == nil ***REMOVED***

		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *fam.maxTimeMS)
	***REMOVED***
	if fam.newDocument != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "new", *fam.newDocument)
	***REMOVED***
	if fam.query != nil ***REMOVED***

		dst = bsoncore.AppendDocumentElement(dst, "query", fam.query)
	***REMOVED***
	if fam.remove != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "remove", *fam.remove)
	***REMOVED***
	if fam.sort != nil ***REMOVED***

		dst = bsoncore.AppendDocumentElement(dst, "sort", fam.sort)
	***REMOVED***
	if fam.update.Data != nil ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "update", fam.update)
	***REMOVED***
	if fam.upsert != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "upsert", *fam.upsert)
	***REMOVED***
	if fam.hint.Type != bsontype.Type(0) ***REMOVED***

		if desc.WireVersion == nil || !desc.WireVersion.Includes(8) ***REMOVED***
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 8")
		***REMOVED***
		if !fam.writeConcern.Acknowledged() ***REMOVED***
			return nil, errUnacknowledgedHint
		***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "hint", fam.hint)
	***REMOVED***
	if fam.let != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "let", fam.let)
	***REMOVED***

	return dst, nil
***REMOVED***

// ArrayFilters specifies an array of filter documents that determines which array elements to modify for an update operation on an array field.
func (fam *FindAndModify) ArrayFilters(arrayFilters bsoncore.Document) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.arrayFilters = arrayFilters
	return fam
***REMOVED***

// BypassDocumentValidation specifies if document validation can be skipped when executing the operation.
func (fam *FindAndModify) BypassDocumentValidation(bypassDocumentValidation bool) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.bypassDocumentValidation = &bypassDocumentValidation
	return fam
***REMOVED***

// Collation specifies a collation to be used.
func (fam *FindAndModify) Collation(collation bsoncore.Document) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.collation = collation
	return fam
***REMOVED***

// Comment sets a value to help trace an operation.
func (fam *FindAndModify) Comment(comment bsoncore.Value) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.comment = comment
	return fam
***REMOVED***

// Fields specifies a subset of fields to return.
func (fam *FindAndModify) Fields(fields bsoncore.Document) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.fields = fields
	return fam
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the operation to run.
func (fam *FindAndModify) MaxTimeMS(maxTimeMS int64) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.maxTimeMS = &maxTimeMS
	return fam
***REMOVED***

// NewDocument specifies whether to return the modified document or the original. Defaults to false (return original).
func (fam *FindAndModify) NewDocument(newDocument bool) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.newDocument = &newDocument
	return fam
***REMOVED***

// Query specifies the selection criteria for the modification.
func (fam *FindAndModify) Query(query bsoncore.Document) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.query = query
	return fam
***REMOVED***

// Remove specifies that the matched document should be removed. Defaults to false.
func (fam *FindAndModify) Remove(remove bool) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.remove = &remove
	return fam
***REMOVED***

// Sort determines which document the operation modifies if the query matches multiple documents.The first document matched by the sort order will be modified.
//
func (fam *FindAndModify) Sort(sort bsoncore.Document) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.sort = sort
	return fam
***REMOVED***

// Update specifies the update document to perform on the matched document.
func (fam *FindAndModify) Update(update bsoncore.Value) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.update = update
	return fam
***REMOVED***

// Upsert specifies whether or not to create a new document if no documents match the query when doing an update. Defaults to false.
func (fam *FindAndModify) Upsert(upsert bool) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.upsert = &upsert
	return fam
***REMOVED***

// Session sets the session for this operation.
func (fam *FindAndModify) Session(session *session.Client) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.session = session
	return fam
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (fam *FindAndModify) ClusterClock(clock *session.ClusterClock) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.clock = clock
	return fam
***REMOVED***

// Collection sets the collection that this command will run against.
func (fam *FindAndModify) Collection(collection string) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.collection = collection
	return fam
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (fam *FindAndModify) CommandMonitor(monitor *event.CommandMonitor) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.monitor = monitor
	return fam
***REMOVED***

// Database sets the database to run this operation against.
func (fam *FindAndModify) Database(database string) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.database = database
	return fam
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (fam *FindAndModify) Deployment(deployment driver.Deployment) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.deployment = deployment
	return fam
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (fam *FindAndModify) ServerSelector(selector description.ServerSelector) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.selector = selector
	return fam
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (fam *FindAndModify) WriteConcern(writeConcern *writeconcern.WriteConcern) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.writeConcern = writeConcern
	return fam
***REMOVED***

// Retry enables retryable writes for this operation. Retries are not handled automatically,
// instead a boolean is returned from Execute and SelectAndExecute that indicates if the
// operation can be retried. Retrying is handled by calling RetryExecute.
func (fam *FindAndModify) Retry(retry driver.RetryMode) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.retry = &retry
	return fam
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (fam *FindAndModify) Crypt(crypt driver.Crypt) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.crypt = crypt
	return fam
***REMOVED***

// Hint specifies the index to use.
func (fam *FindAndModify) Hint(hint bsoncore.Value) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.hint = hint
	return fam
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (fam *FindAndModify) ServerAPI(serverAPI *driver.ServerAPIOptions) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.serverAPI = serverAPI
	return fam
***REMOVED***

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (fam *FindAndModify) Let(let bsoncore.Document) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.let = let
	return fam
***REMOVED***

// Timeout sets the timeout for this operation.
func (fam *FindAndModify) Timeout(timeout *time.Duration) *FindAndModify ***REMOVED***
	if fam == nil ***REMOVED***
		fam = new(FindAndModify)
	***REMOVED***

	fam.timeout = timeout
	return fam
***REMOVED***
