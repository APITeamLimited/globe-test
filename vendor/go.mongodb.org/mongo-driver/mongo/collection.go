// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Collection is a handle to a MongoDB collection. It is safe for concurrent use by multiple goroutines.
type Collection struct ***REMOVED***
	client         *Client
	db             *Database
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	registry       *bsoncodec.Registry
***REMOVED***

// aggregateParams is used to store information to configure an Aggregate operation.
type aggregateParams struct ***REMOVED***
	ctx            context.Context
	pipeline       interface***REMOVED******REMOVED***
	client         *Client
	registry       *bsoncodec.Registry
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	retryRead      bool
	db             string
	col            string
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	readPreference *readpref.ReadPref
	opts           []*options.AggregateOptions
***REMOVED***

func closeImplicitSession(sess *session.Client) ***REMOVED***
	if sess != nil && sess.SessionType == session.Implicit ***REMOVED***
		sess.EndSession()
	***REMOVED***
***REMOVED***

func newCollection(db *Database, name string, opts ...*options.CollectionOptions) *Collection ***REMOVED***
	collOpt := options.MergeCollectionOptions(opts...)

	rc := db.readConcern
	if collOpt.ReadConcern != nil ***REMOVED***
		rc = collOpt.ReadConcern
	***REMOVED***

	wc := db.writeConcern
	if collOpt.WriteConcern != nil ***REMOVED***
		wc = collOpt.WriteConcern
	***REMOVED***

	rp := db.readPreference
	if collOpt.ReadPreference != nil ***REMOVED***
		rp = collOpt.ReadPreference
	***REMOVED***

	reg := db.registry
	if collOpt.Registry != nil ***REMOVED***
		reg = collOpt.Registry
	***REMOVED***

	readSelector := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(rp),
		description.LatencySelector(db.client.localThreshold),
	***REMOVED***)

	writeSelector := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.WriteSelector(),
		description.LatencySelector(db.client.localThreshold),
	***REMOVED***)

	coll := &Collection***REMOVED***
		client:         db.client,
		db:             db,
		name:           name,
		readPreference: rp,
		readConcern:    rc,
		writeConcern:   wc,
		readSelector:   readSelector,
		writeSelector:  writeSelector,
		registry:       reg,
	***REMOVED***

	return coll
***REMOVED***

func (coll *Collection) copy() *Collection ***REMOVED***
	return &Collection***REMOVED***
		client:         coll.client,
		db:             coll.db,
		name:           coll.name,
		readConcern:    coll.readConcern,
		writeConcern:   coll.writeConcern,
		readPreference: coll.readPreference,
		readSelector:   coll.readSelector,
		writeSelector:  coll.writeSelector,
		registry:       coll.registry,
	***REMOVED***
***REMOVED***

// Clone creates a copy of the Collection configured with the given CollectionOptions.
// The specified options are merged with the existing options on the collection, with the specified options taking
// precedence.
func (coll *Collection) Clone(opts ...*options.CollectionOptions) (*Collection, error) ***REMOVED***
	copyColl := coll.copy()
	optsColl := options.MergeCollectionOptions(opts...)

	if optsColl.ReadConcern != nil ***REMOVED***
		copyColl.readConcern = optsColl.ReadConcern
	***REMOVED***

	if optsColl.WriteConcern != nil ***REMOVED***
		copyColl.writeConcern = optsColl.WriteConcern
	***REMOVED***

	if optsColl.ReadPreference != nil ***REMOVED***
		copyColl.readPreference = optsColl.ReadPreference
	***REMOVED***

	if optsColl.Registry != nil ***REMOVED***
		copyColl.registry = optsColl.Registry
	***REMOVED***

	copyColl.readSelector = description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(copyColl.readPreference),
		description.LatencySelector(copyColl.client.localThreshold),
	***REMOVED***)

	return copyColl, nil
***REMOVED***

// Name returns the name of the collection.
func (coll *Collection) Name() string ***REMOVED***
	return coll.name
***REMOVED***

// Database returns the Database that was used to create the Collection.
func (coll *Collection) Database() *Database ***REMOVED***
	return coll.db
***REMOVED***

// BulkWrite performs a bulk write operation (https://www.mongodb.com/docs/manual/core/bulk-write-operations/).
//
// The models parameter must be a slice of operations to be executed in this bulk write. It cannot be nil or empty.
// All of the models must be non-nil. See the mongo.WriteModel documentation for a list of valid model types and
// examples of how they should be used.
//
// The opts parameter can be used to specify options for the operation (see the options.BulkWriteOptions documentation.)
func (coll *Collection) BulkWrite(ctx context.Context, models []WriteModel,
	opts ...*options.BulkWriteOptions) (*BulkWriteResult, error) ***REMOVED***

	if len(models) == 0 ***REMOVED***
		return nil, ErrEmptySlice
	***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err := coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wc := coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, coll.writeSelector)

	for _, model := range models ***REMOVED***
		if model == nil ***REMOVED***
			return nil, ErrNilDocument
		***REMOVED***
	***REMOVED***

	bwo := options.MergeBulkWriteOptions(opts...)

	op := bulkWrite***REMOVED***
		comment:                  bwo.Comment,
		ordered:                  bwo.Ordered,
		bypassDocumentValidation: bwo.BypassDocumentValidation,
		models:                   models,
		session:                  sess,
		collection:               coll,
		selector:                 selector,
		writeConcern:             wc,
		let:                      bwo.Let,
	***REMOVED***

	err = op.execute(ctx)

	return &op.result, replaceErrors(err)
***REMOVED***

func (coll *Collection) insert(ctx context.Context, documents []interface***REMOVED******REMOVED***,
	opts ...*options.InsertManyOptions) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	result := make([]interface***REMOVED******REMOVED***, len(documents))
	docs := make([]bsoncore.Document, len(documents))

	for i, doc := range documents ***REMOVED***
		var err error
		docs[i], result[i], err = transformAndEnsureID(coll.registry, doc)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err := coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wc := coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewInsert(docs...).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).Ordered(true).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout)
	imo := options.MergeInsertManyOptions(opts...)
	if imo.BypassDocumentValidation != nil && *imo.BypassDocumentValidation ***REMOVED***
		op = op.BypassDocumentValidation(*imo.BypassDocumentValidation)
	***REMOVED***
	if imo.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, imo.Comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	if imo.Ordered != nil ***REMOVED***
		op = op.Ordered(*imo.Ordered)
	***REMOVED***
	retry := driver.RetryNone
	if coll.client.retryWrites ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err = op.Execute(ctx)
	wce, ok := err.(driver.WriteCommandError)
	if !ok ***REMOVED***
		return result, err
	***REMOVED***

	// remove the ids that had writeErrors from result
	for i, we := range wce.WriteErrors ***REMOVED***
		// i indexes have been removed before the current error, so the index is we.Index-i
		idIndex := int(we.Index) - i
		// if the insert is ordered, nothing after the error was inserted
		if imo.Ordered == nil || *imo.Ordered ***REMOVED***
			result = result[:idIndex]
			break
		***REMOVED***
		result = append(result[:idIndex], result[idIndex+1:]...)
	***REMOVED***

	return result, err
***REMOVED***

// InsertOne executes an insert command to insert a single document into the collection.
//
// The document parameter must be the document to be inserted. It cannot be nil. If the document does not have an _id
// field when transformed into BSON, one will be added automatically to the marshalled document. The original document
// will not be modified. The _id can be retrieved from the InsertedID field of the returned InsertOneResult.
//
// The opts parameter can be used to specify options for the operation (see the options.InsertOneOptions documentation.)
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/insert/.
func (coll *Collection) InsertOne(ctx context.Context, document interface***REMOVED******REMOVED***,
	opts ...*options.InsertOneOptions) (*InsertOneResult, error) ***REMOVED***

	ioOpts := options.MergeInsertOneOptions(opts...)
	imOpts := options.InsertMany()

	if ioOpts.BypassDocumentValidation != nil && *ioOpts.BypassDocumentValidation ***REMOVED***
		imOpts.SetBypassDocumentValidation(*ioOpts.BypassDocumentValidation)
	***REMOVED***
	if ioOpts.Comment != nil ***REMOVED***
		imOpts.SetComment(ioOpts.Comment)
	***REMOVED***
	res, err := coll.insert(ctx, []interface***REMOVED******REMOVED******REMOVED***document***REMOVED***, imOpts)

	rr, err := processWriteError(err)
	if rr&rrOne == 0 ***REMOVED***
		return nil, err
	***REMOVED***
	return &InsertOneResult***REMOVED***InsertedID: res[0]***REMOVED***, err
***REMOVED***

// InsertMany executes an insert command to insert multiple documents into the collection. If write errors occur
// during the operation (e.g. duplicate key error), this method returns a BulkWriteException error.
//
// The documents parameter must be a slice of documents to insert. The slice cannot be nil or empty. The elements must
// all be non-nil. For any document that does not have an _id field when transformed into BSON, one will be added
// automatically to the marshalled document. The original document will not be modified. The _id values for the inserted
// documents can be retrieved from the InsertedIDs field of the returned InsertManyResult.
//
// The opts parameter can be used to specify options for the operation (see the options.InsertManyOptions documentation.)
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/insert/.
func (coll *Collection) InsertMany(ctx context.Context, documents []interface***REMOVED******REMOVED***,
	opts ...*options.InsertManyOptions) (*InsertManyResult, error) ***REMOVED***

	if len(documents) == 0 ***REMOVED***
		return nil, ErrEmptySlice
	***REMOVED***

	result, err := coll.insert(ctx, documents, opts...)
	rr, err := processWriteError(err)
	if rr&rrMany == 0 ***REMOVED***
		return nil, err
	***REMOVED***

	imResult := &InsertManyResult***REMOVED***InsertedIDs: result***REMOVED***
	writeException, ok := err.(WriteException)
	if !ok ***REMOVED***
		return imResult, err
	***REMOVED***

	// create and return a BulkWriteException
	bwErrors := make([]BulkWriteError, 0, len(writeException.WriteErrors))
	for _, we := range writeException.WriteErrors ***REMOVED***
		bwErrors = append(bwErrors, BulkWriteError***REMOVED***
			WriteError: we,
			Request:    nil,
		***REMOVED***)
	***REMOVED***

	return imResult, BulkWriteException***REMOVED***
		WriteErrors:       bwErrors,
		WriteConcernError: writeException.WriteConcernError,
		Labels:            writeException.Labels,
	***REMOVED***
***REMOVED***

func (coll *Collection) delete(ctx context.Context, filter interface***REMOVED******REMOVED***, deleteOne bool, expectedRr returnResult,
	opts ...*options.DeleteOptions) (*DeleteResult, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wc := coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, coll.writeSelector)

	var limit int32
	if deleteOne ***REMOVED***
		limit = 1
	***REMOVED***
	do := options.MergeDeleteOptions(opts...)
	didx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendDocumentElement(doc, "q", f)
	doc = bsoncore.AppendInt32Element(doc, "limit", limit)
	if do.Collation != nil ***REMOVED***
		doc = bsoncore.AppendDocumentElement(doc, "collation", do.Collation.ToDocument())
	***REMOVED***
	if do.Hint != nil ***REMOVED***
		hint, err := transformValue(coll.registry, do.Hint, false, "hint")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		doc = bsoncore.AppendValueElement(doc, "hint", hint)
	***REMOVED***
	doc, _ = bsoncore.AppendDocumentEnd(doc, didx)

	op := operation.NewDelete(doc).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).Ordered(true).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout)
	if do.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, do.Comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	if do.Hint != nil ***REMOVED***
		op = op.Hint(true)
	***REMOVED***
	if do.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(coll.registry, do.Let, true, "let")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op = op.Let(let)
	***REMOVED***

	// deleteMany cannot be retried
	retryMode := driver.RetryNone
	if deleteOne && coll.client.retryWrites ***REMOVED***
		retryMode = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retryMode)
	rr, err := processWriteError(op.Execute(ctx))
	if rr&expectedRr == 0 ***REMOVED***
		return nil, err
	***REMOVED***
	return &DeleteResult***REMOVED***DeletedCount: op.Result().N***REMOVED***, err
***REMOVED***

// DeleteOne executes a delete command to delete at most one document from the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// deleted. It cannot be nil. If the filter does not match any documents, the operation will succeed and a DeleteResult
// with a DeletedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
// matched set.
//
// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/delete/.
func (coll *Collection) DeleteOne(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.DeleteOptions) (*DeleteResult, error) ***REMOVED***

	return coll.delete(ctx, filter, true, rrOne, opts...)
***REMOVED***

// DeleteMany executes a delete command to delete documents from the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the documents to
// be deleted. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to delete all documents in the
// collection. If the filter does not match any documents, the operation will succeed and a DeleteResult with a
// DeletedCount of 0 will be returned.
//
// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/delete/.
func (coll *Collection) DeleteMany(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.DeleteOptions) (*DeleteResult, error) ***REMOVED***

	return coll.delete(ctx, filter, false, rrMany, opts...)
***REMOVED***

func (coll *Collection) updateOrReplace(ctx context.Context, filter bsoncore.Document, update interface***REMOVED******REMOVED***, multi bool,
	expectedRr returnResult, checkDollarKey bool, opts ...*options.UpdateOptions) (*UpdateResult, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	uo := options.MergeUpdateOptions(opts...)

	// collation, arrayFilters, upsert, and hint are included on the individual update documents rather than as part of the
	// command
	updateDoc, err := createUpdateDoc(filter, update, uo.Hint, uo.ArrayFilters, uo.Collation, uo.Upsert, multi,
		checkDollarKey, coll.registry)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wc := coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewUpdate(updateDoc).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).Hint(uo.Hint != nil).
		ArrayFilters(uo.ArrayFilters != nil).Ordered(true).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout)
	if uo.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(coll.registry, uo.Let, true, "let")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op = op.Let(let)
	***REMOVED***

	if uo.BypassDocumentValidation != nil && *uo.BypassDocumentValidation ***REMOVED***
		op = op.BypassDocumentValidation(*uo.BypassDocumentValidation)
	***REMOVED***
	if uo.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, uo.Comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	retry := driver.RetryNone
	// retryable writes are only enabled updateOne/replaceOne operations
	if !multi && coll.client.retryWrites ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)
	err = op.Execute(ctx)

	rr, err := processWriteError(err)
	if rr&expectedRr == 0 ***REMOVED***
		return nil, err
	***REMOVED***

	opRes := op.Result()
	res := &UpdateResult***REMOVED***
		MatchedCount:  opRes.N,
		ModifiedCount: opRes.NModified,
		UpsertedCount: int64(len(opRes.Upserted)),
	***REMOVED***
	if len(opRes.Upserted) > 0 ***REMOVED***
		res.UpsertedID = opRes.Upserted[0].ID
		res.MatchedCount--
	***REMOVED***

	return res, err
***REMOVED***

// UpdateByID executes an update command to update the document whose _id value matches the provided ID in the collection.
// This is equivalent to running UpdateOne(ctx, bson.D***REMOVED******REMOVED***"_id", id***REMOVED******REMOVED***, update, opts...).
//
// The id parameter is the _id of the document to be updated. It cannot be nil. If the ID does not match any documents,
// the operation will succeed and an UpdateResult with a MatchedCount of 0 will be returned.
//
// The update parameter must be a document containing update operators
// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be
// made to the selected document. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
func (coll *Collection) UpdateByID(ctx context.Context, id interface***REMOVED******REMOVED***, update interface***REMOVED******REMOVED***,
	opts ...*options.UpdateOptions) (*UpdateResult, error) ***REMOVED***
	if id == nil ***REMOVED***
		return nil, ErrNilValue
	***REMOVED***
	return coll.UpdateOne(ctx, bson.D***REMOVED******REMOVED***"_id", id***REMOVED******REMOVED***, update, opts...)
***REMOVED***

// UpdateOne executes an update command to update at most one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
// with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
// matched set and MatchedCount will equal 1.
//
// The update parameter must be a document containing update operators
// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be
// made to the selected document. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
func (coll *Collection) UpdateOne(ctx context.Context, filter interface***REMOVED******REMOVED***, update interface***REMOVED******REMOVED***,
	opts ...*options.UpdateOptions) (*UpdateResult, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return coll.updateOrReplace(ctx, f, update, false, rrOne, true, opts...)
***REMOVED***

// UpdateMany executes an update command to update documents in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the documents to be
// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
// with a MatchedCount of 0 will be returned.
//
// The update parameter must be a document containing update operators
// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be made
// to the selected documents. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
func (coll *Collection) UpdateMany(ctx context.Context, filter interface***REMOVED******REMOVED***, update interface***REMOVED******REMOVED***,
	opts ...*options.UpdateOptions) (*UpdateResult, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return coll.updateOrReplace(ctx, f, update, true, rrMany, true, opts...)
***REMOVED***

// ReplaceOne executes an update command to replace at most one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// replaced. It cannot be nil. If the filter does not match any documents, the operation will succeed and an
// UpdateResult with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be
// selected from the matched set and MatchedCount will equal 1.
//
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
//
// The opts parameter can be used to specify options for the operation (see the options.ReplaceOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
func (coll *Collection) ReplaceOne(ctx context.Context, filter interface***REMOVED******REMOVED***,
	replacement interface***REMOVED******REMOVED***, opts ...*options.ReplaceOptions) (*UpdateResult, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	r, err := transformBsoncoreDocument(coll.registry, replacement, true, "replacement")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := ensureNoDollarKey(r); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	updateOptions := make([]*options.UpdateOptions, 0, len(opts))
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		uOpts := options.Update()
		uOpts.BypassDocumentValidation = opt.BypassDocumentValidation
		uOpts.Collation = opt.Collation
		uOpts.Upsert = opt.Upsert
		uOpts.Hint = opt.Hint
		uOpts.Let = opt.Let
		uOpts.Comment = opt.Comment
		updateOptions = append(updateOptions, uOpts)
	***REMOVED***

	return coll.updateOrReplace(ctx, f, r, false, rrOne, false, updateOptions...)
***REMOVED***

// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
//
// The pipeline parameter must be an array of documents, each representing an aggregation stage. The pipeline cannot
// be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
// mongo.Pipeline type can be used. See
// https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/#db-collection-aggregate-stages for a list of
// valid stages in aggregations.
//
// The opts parameter can be used to specify options for the operation (see the options.AggregateOptions documentation.)
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/aggregate/.
func (coll *Collection) Aggregate(ctx context.Context, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.AggregateOptions) (*Cursor, error) ***REMOVED***
	a := aggregateParams***REMOVED***
		ctx:            ctx,
		pipeline:       pipeline,
		client:         coll.client,
		registry:       coll.registry,
		readConcern:    coll.readConcern,
		writeConcern:   coll.writeConcern,
		retryRead:      coll.client.retryReads,
		db:             coll.db.name,
		col:            coll.name,
		readSelector:   coll.readSelector,
		writeSelector:  coll.writeSelector,
		readPreference: coll.readPreference,
		opts:           opts,
	***REMOVED***
	return aggregate(a)
***REMOVED***

// aggregate is the helper method for Aggregate
func aggregate(a aggregateParams) (cur *Cursor, err error) ***REMOVED***
	if a.ctx == nil ***REMOVED***
		a.ctx = context.Background()
	***REMOVED***

	pipelineArr, hasOutputStage, err := transformAggregatePipeline(a.registry, a.pipeline)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(a.ctx)
	// Always close any created implicit sessions if aggregate returns an error.
	defer func() ***REMOVED***
		if err != nil && sess != nil ***REMOVED***
			closeImplicitSession(sess)
		***REMOVED***
	***REMOVED***()
	if sess == nil && a.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(a.client.sessionPool, a.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if err = a.client.validSession(sess); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var wc *writeconcern.WriteConcern
	if hasOutputStage ***REMOVED***
		wc = a.writeConcern
	***REMOVED***
	rc := a.readConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
		rc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		closeImplicitSession(sess)
		sess = nil
	***REMOVED***

	selector := makeReadPrefSelector(sess, a.readSelector, a.client.localThreshold)
	if hasOutputStage ***REMOVED***
		selector = makeOutputAggregateSelector(sess, a.readPreference, a.client.localThreshold)
	***REMOVED***

	ao := options.MergeAggregateOptions(a.opts...)
	cursorOpts := a.client.createBaseCursorOptions()

	op := operation.NewAggregate(pipelineArr).
		Session(sess).
		WriteConcern(wc).
		ReadConcern(rc).
		ReadPreference(a.readPreference).
		CommandMonitor(a.client.monitor).
		ServerSelector(selector).
		ClusterClock(a.client.clock).
		Database(a.db).
		Collection(a.col).
		Deployment(a.client.deployment).
		Crypt(a.client.cryptFLE).
		ServerAPI(a.client.serverAPI).
		HasOutputStage(hasOutputStage).
		Timeout(a.client.timeout)

	if ao.AllowDiskUse != nil ***REMOVED***
		op.AllowDiskUse(*ao.AllowDiskUse)
	***REMOVED***
	// ignore batchSize of 0 with $out
	if ao.BatchSize != nil && !(*ao.BatchSize == 0 && hasOutputStage) ***REMOVED***
		op.BatchSize(*ao.BatchSize)
		cursorOpts.BatchSize = *ao.BatchSize
	***REMOVED***
	if ao.BypassDocumentValidation != nil && *ao.BypassDocumentValidation ***REMOVED***
		op.BypassDocumentValidation(*ao.BypassDocumentValidation)
	***REMOVED***
	if ao.Collation != nil ***REMOVED***
		op.Collation(bsoncore.Document(ao.Collation.ToDocument()))
	***REMOVED***
	if ao.MaxTime != nil ***REMOVED***
		op.MaxTimeMS(int64(*ao.MaxTime / time.Millisecond))
	***REMOVED***
	if ao.MaxAwaitTime != nil ***REMOVED***
		cursorOpts.MaxTimeMS = int64(*ao.MaxAwaitTime / time.Millisecond)
	***REMOVED***
	if ao.Comment != nil ***REMOVED***
		op.Comment(*ao.Comment)

		commentVal, err := transformValue(a.registry, ao.Comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cursorOpts.Comment = commentVal
	***REMOVED***
	if ao.Hint != nil ***REMOVED***
		hintVal, err := transformValue(a.registry, ao.Hint, false, "hint")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Hint(hintVal)
	***REMOVED***
	if ao.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(a.registry, ao.Let, true, "let")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Let(let)
	***REMOVED***
	if ao.Custom != nil ***REMOVED***
		// Marshal all custom options before passing to the aggregate operation. Return
		// any errors from Marshaling.
		customOptions := make(map[string]bsoncore.Value)
		for optionName, optionValue := range ao.Custom ***REMOVED***
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(a.registry, optionValue)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			optionValueBSON := bsoncore.Value***REMOVED***Type: bsonType, Data: bsonData***REMOVED***
			customOptions[optionName] = optionValueBSON
		***REMOVED***
		op.CustomOptions(customOptions)
	***REMOVED***

	retry := driver.RetryNone
	if a.retryRead && !hasOutputStage ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err = op.Execute(a.ctx)
	if err != nil ***REMOVED***
		if wce, ok := err.(driver.WriteCommandError); ok && wce.WriteConcernError != nil ***REMOVED***
			return nil, *convertDriverWriteConcernError(wce.WriteConcernError)
		***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***

	bc, err := op.Result(cursorOpts)
	if err != nil ***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***
	cursor, err := newCursorWithSession(bc, a.registry, sess)
	return cursor, replaceErrors(err)
***REMOVED***

// CountDocuments returns the number of documents in the collection. For a fast count of the documents in the
// collection, see the EstimatedDocumentCount method.
//
// The filter parameter must be a document and can be used to select which documents contribute to the count. It
// cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to count all documents in the collection. This will
// result in a full collection scan.
//
// The opts parameter can be used to specify options for the operation (see the options.CountOptions documentation).
func (coll *Collection) CountDocuments(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.CountOptions) (int64, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	countOpts := options.MergeCountOptions(opts...)

	pipelineArr, err := countDocumentsAggregatePipeline(coll.registry, filter, countOpts)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***
	if err = coll.client.validSession(sess); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	rc := coll.readConcern
	if sess.TransactionRunning() ***REMOVED***
		rc = nil
	***REMOVED***

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewAggregate(pipelineArr).Session(sess).ReadConcern(rc).ReadPreference(coll.readPreference).
		CommandMonitor(coll.client.monitor).ServerSelector(selector).ClusterClock(coll.client.clock).Database(coll.db.name).
		Collection(coll.name).Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout)
	if countOpts.Collation != nil ***REMOVED***
		op.Collation(bsoncore.Document(countOpts.Collation.ToDocument()))
	***REMOVED***
	if countOpts.Comment != nil ***REMOVED***
		op.Comment(*countOpts.Comment)
	***REMOVED***
	if countOpts.MaxTime != nil ***REMOVED***
		op.MaxTimeMS(int64(*countOpts.MaxTime / time.Millisecond))
	***REMOVED***
	if countOpts.Hint != nil ***REMOVED***
		hintVal, err := transformValue(coll.registry, countOpts.Hint, false, "hint")
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		op.Hint(hintVal)
	***REMOVED***
	retry := driver.RetryNone
	if coll.client.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		return 0, replaceErrors(err)
	***REMOVED***

	batch := op.ResultCursorResponse().FirstBatch
	if batch == nil ***REMOVED***
		return 0, errors.New("invalid response from server, no 'firstBatch' field")
	***REMOVED***

	docs, err := batch.Documents()
	if err != nil || len(docs) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***

	val, ok := docs[0].Lookup("n").AsInt64OK()
	if !ok ***REMOVED***
		return 0, errors.New("invalid response from server, no 'n' field")
	***REMOVED***

	return val, nil
***REMOVED***

// EstimatedDocumentCount executes a count command and returns an estimate of the number of documents in the collection
// using collection metadata.
//
// The opts parameter can be used to specify options for the operation (see the options.EstimatedDocumentCountOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/count/.
func (coll *Collection) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)

	var err error
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	rc := coll.readConcern
	if sess.TransactionRunning() ***REMOVED***
		rc = nil
	***REMOVED***

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewCount().Session(sess).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).CommandMonitor(coll.client.monitor).
		Deployment(coll.client.deployment).ReadConcern(rc).ReadPreference(coll.readPreference).
		ServerSelector(selector).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout)

	co := options.MergeEstimatedDocumentCountOptions(opts...)
	if co.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, co.Comment, false, "comment")
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	if co.MaxTime != nil ***REMOVED***
		op = op.MaxTimeMS(int64(*co.MaxTime / time.Millisecond))
	***REMOVED***
	retry := driver.RetryNone
	if coll.client.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op.Retry(retry)

	err = op.Execute(ctx)

	return op.Result().N, replaceErrors(err)
***REMOVED***

// Distinct executes a distinct command to find the unique values for a specified field in the collection.
//
// The fieldName parameter specifies the field name for which distinct values should be returned.
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// considered. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to select all documents.
//
// The opts parameter can be used to specify options for the operation (see the options.DistinctOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/distinct/.
func (coll *Collection) Distinct(ctx context.Context, fieldName string, filter interface***REMOVED******REMOVED***,
	opts ...*options.DistinctOptions) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(ctx)

	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rc := coll.readConcern
	if sess.TransactionRunning() ***REMOVED***
		rc = nil
	***REMOVED***

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	option := options.MergeDistinctOptions(opts...)

	op := operation.NewDistinct(fieldName, f).
		Session(sess).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).CommandMonitor(coll.client.monitor).
		Deployment(coll.client.deployment).ReadConcern(rc).ReadPreference(coll.readPreference).
		ServerSelector(selector).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout)

	if option.Collation != nil ***REMOVED***
		op.Collation(bsoncore.Document(option.Collation.ToDocument()))
	***REMOVED***
	if option.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, option.Comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Comment(comment)
	***REMOVED***
	if option.MaxTime != nil ***REMOVED***
		op.MaxTimeMS(int64(*option.MaxTime / time.Millisecond))
	***REMOVED***
	retry := driver.RetryNone
	if coll.client.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***

	arr, ok := op.Result().Values.ArrayOK()
	if !ok ***REMOVED***
		return nil, fmt.Errorf("response field 'values' is type array, but received BSON type %s", op.Result().Values.Type)
	***REMOVED***

	values, err := arr.Values()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	retArray := make([]interface***REMOVED******REMOVED***, len(values))

	for i, val := range values ***REMOVED***
		raw := bson.RawValue***REMOVED***Type: val.Type, Value: val.Data***REMOVED***
		err = raw.Unmarshal(&retArray[i])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return retArray, replaceErrors(err)
***REMOVED***

// Find executes a find command and returns a Cursor over the matching documents in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// included in the result. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to include all documents.
//
// The opts parameter can be used to specify options for the operation (see the options.FindOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/find/.
func (coll *Collection) Find(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.FindOptions) (cur *Cursor, err error) ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(ctx)
	// Always close any created implicit sessions if Find returns an error.
	defer func() ***REMOVED***
		if err != nil && sess != nil ***REMOVED***
			closeImplicitSession(sess)
		***REMOVED***
	***REMOVED***()
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	err = coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rc := coll.readConcern
	if sess.TransactionRunning() ***REMOVED***
		rc = nil
	***REMOVED***

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewFind(f).
		Session(sess).ReadConcern(rc).ReadPreference(coll.readPreference).
		CommandMonitor(coll.client.monitor).ServerSelector(selector).
		ClusterClock(coll.client.clock).Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout)

	fo := options.MergeFindOptions(opts...)
	cursorOpts := coll.client.createBaseCursorOptions()

	if fo.AllowDiskUse != nil ***REMOVED***
		op.AllowDiskUse(*fo.AllowDiskUse)
	***REMOVED***
	if fo.AllowPartialResults != nil ***REMOVED***
		op.AllowPartialResults(*fo.AllowPartialResults)
	***REMOVED***
	if fo.BatchSize != nil ***REMOVED***
		cursorOpts.BatchSize = *fo.BatchSize
		op.BatchSize(*fo.BatchSize)
	***REMOVED***
	if fo.Collation != nil ***REMOVED***
		op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	***REMOVED***
	if fo.Comment != nil ***REMOVED***
		op.Comment(*fo.Comment)

		commentVal, err := transformValue(coll.registry, fo.Comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cursorOpts.Comment = commentVal
	***REMOVED***
	if fo.CursorType != nil ***REMOVED***
		switch *fo.CursorType ***REMOVED***
		case options.Tailable:
			op.Tailable(true)
		case options.TailableAwait:
			op.Tailable(true)
			op.AwaitData(true)
		***REMOVED***
	***REMOVED***
	if fo.Hint != nil ***REMOVED***
		hint, err := transformValue(coll.registry, fo.Hint, false, "hint")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Hint(hint)
	***REMOVED***
	if fo.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(coll.registry, fo.Let, true, "let")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Let(let)
	***REMOVED***
	if fo.Limit != nil ***REMOVED***
		limit := *fo.Limit
		if limit < 0 ***REMOVED***
			limit = -1 * limit
			op.SingleBatch(true)
		***REMOVED***
		cursorOpts.Limit = int32(limit)
		op.Limit(limit)
	***REMOVED***
	if fo.Max != nil ***REMOVED***
		max, err := transformBsoncoreDocument(coll.registry, fo.Max, true, "max")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Max(max)
	***REMOVED***
	if fo.MaxAwaitTime != nil ***REMOVED***
		cursorOpts.MaxTimeMS = int64(*fo.MaxAwaitTime / time.Millisecond)
	***REMOVED***
	if fo.MaxTime != nil ***REMOVED***
		op.MaxTimeMS(int64(*fo.MaxTime / time.Millisecond))
	***REMOVED***
	if fo.Min != nil ***REMOVED***
		min, err := transformBsoncoreDocument(coll.registry, fo.Min, true, "min")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Min(min)
	***REMOVED***
	if fo.NoCursorTimeout != nil ***REMOVED***
		op.NoCursorTimeout(*fo.NoCursorTimeout)
	***REMOVED***
	if fo.OplogReplay != nil ***REMOVED***
		op.OplogReplay(*fo.OplogReplay)
	***REMOVED***
	if fo.Projection != nil ***REMOVED***
		proj, err := transformBsoncoreDocument(coll.registry, fo.Projection, true, "projection")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Projection(proj)
	***REMOVED***
	if fo.ReturnKey != nil ***REMOVED***
		op.ReturnKey(*fo.ReturnKey)
	***REMOVED***
	if fo.ShowRecordID != nil ***REMOVED***
		op.ShowRecordID(*fo.ShowRecordID)
	***REMOVED***
	if fo.Skip != nil ***REMOVED***
		op.Skip(*fo.Skip)
	***REMOVED***
	if fo.Snapshot != nil ***REMOVED***
		op.Snapshot(*fo.Snapshot)
	***REMOVED***
	if fo.Sort != nil ***REMOVED***
		sort, err := transformBsoncoreDocument(coll.registry, fo.Sort, false, "sort")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Sort(sort)
	***REMOVED***
	retry := driver.RetryNone
	if coll.client.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	if err = op.Execute(ctx); err != nil ***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***

	bc, err := op.Result(cursorOpts)
	if err != nil ***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***
	return newCursorWithSession(bc, coll.registry, sess)
***REMOVED***

// FindOne executes a find command and returns a SingleResult for one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// returned. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments will be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The opts parameter can be used to specify options for this operation (see the options.FindOneOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/find/.
func (coll *Collection) FindOne(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.FindOneOptions) *SingleResult ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	findOpts := make([]*options.FindOptions, 0, len(opts))
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		findOpts = append(findOpts, &options.FindOptions***REMOVED***
			AllowPartialResults: opt.AllowPartialResults,
			BatchSize:           opt.BatchSize,
			Collation:           opt.Collation,
			Comment:             opt.Comment,
			CursorType:          opt.CursorType,
			Hint:                opt.Hint,
			Max:                 opt.Max,
			MaxAwaitTime:        opt.MaxAwaitTime,
			MaxTime:             opt.MaxTime,
			Min:                 opt.Min,
			NoCursorTimeout:     opt.NoCursorTimeout,
			OplogReplay:         opt.OplogReplay,
			Projection:          opt.Projection,
			ReturnKey:           opt.ReturnKey,
			ShowRecordID:        opt.ShowRecordID,
			Skip:                opt.Skip,
			Snapshot:            opt.Snapshot,
			Sort:                opt.Sort,
		***REMOVED***)
	***REMOVED***
	// Unconditionally send a limit to make sure only one document is returned and the cursor is not kept open
	// by the server.
	findOpts = append(findOpts, options.Find().SetLimit(-1))

	cursor, err := coll.Find(ctx, filter, findOpts...)
	return &SingleResult***REMOVED***cur: cursor, reg: coll.registry, err: replaceErrors(err)***REMOVED***
***REMOVED***

func (coll *Collection) findAndModify(ctx context.Context, op *operation.FindAndModify) *SingleResult ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)
	var err error
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***

	wc := coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, coll.writeSelector)

	retry := driver.RetryNone
	if coll.client.retryWrites ***REMOVED***
		retry = driver.RetryOnce
	***REMOVED***

	op = op.Session(sess).
		WriteConcern(wc).
		CommandMonitor(coll.client.monitor).
		ServerSelector(selector).
		ClusterClock(coll.client.clock).
		Database(coll.db.name).
		Collection(coll.name).
		Deployment(coll.client.deployment).
		Retry(retry).
		Crypt(coll.client.cryptFLE)

	_, err = processWriteError(op.Execute(ctx))
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***

	return &SingleResult***REMOVED***rdr: bson.Raw(op.Result().Value), reg: coll.registry***REMOVED***
***REMOVED***

// FindOneAndDelete executes a findAndModify command to delete at most one document in the collection. and returns the
// document as it appeared before deletion.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// deleted. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The opts parameter can be used to specify options for the operation (see the options.FindOneAndDeleteOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/findAndModify/.
func (coll *Collection) FindOneAndDelete(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.FindOneAndDeleteOptions) *SingleResult ***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***
	fod := options.MergeFindOneAndDeleteOptions(opts...)
	op := operation.NewFindAndModify(f).Remove(true).ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout)
	if fod.Collation != nil ***REMOVED***
		op = op.Collation(bsoncore.Document(fod.Collation.ToDocument()))
	***REMOVED***
	if fod.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, fod.Comment, true, "comment")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	if fod.MaxTime != nil ***REMOVED***
		op = op.MaxTimeMS(int64(*fod.MaxTime / time.Millisecond))
	***REMOVED***
	if fod.Projection != nil ***REMOVED***
		proj, err := transformBsoncoreDocument(coll.registry, fod.Projection, true, "projection")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Fields(proj)
	***REMOVED***
	if fod.Sort != nil ***REMOVED***
		sort, err := transformBsoncoreDocument(coll.registry, fod.Sort, false, "sort")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Sort(sort)
	***REMOVED***
	if fod.Hint != nil ***REMOVED***
		hint, err := transformValue(coll.registry, fod.Hint, false, "hint")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Hint(hint)
	***REMOVED***
	if fod.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(coll.registry, fod.Let, true, "let")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Let(let)
	***REMOVED***

	return coll.findAndModify(ctx, op)
***REMOVED***

// FindOneAndReplace executes a findAndModify command to replace at most one document in the collection
// and returns the document as it appeared before replacement.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// replaced. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
//
// The opts parameter can be used to specify options for the operation (see the options.FindOneAndReplaceOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/findAndModify/.
func (coll *Collection) FindOneAndReplace(ctx context.Context, filter interface***REMOVED******REMOVED***,
	replacement interface***REMOVED******REMOVED***, opts ...*options.FindOneAndReplaceOptions) *SingleResult ***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***
	r, err := transformBsoncoreDocument(coll.registry, replacement, true, "replacement")
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***
	if firstElem, err := r.IndexErr(0); err == nil && strings.HasPrefix(firstElem.Key(), "$") ***REMOVED***
		return &SingleResult***REMOVED***err: errors.New("replacement document cannot contain keys beginning with '$'")***REMOVED***
	***REMOVED***

	fo := options.MergeFindOneAndReplaceOptions(opts...)
	op := operation.NewFindAndModify(f).Update(bsoncore.Value***REMOVED***Type: bsontype.EmbeddedDocument, Data: r***REMOVED***).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout)
	if fo.BypassDocumentValidation != nil && *fo.BypassDocumentValidation ***REMOVED***
		op = op.BypassDocumentValidation(*fo.BypassDocumentValidation)
	***REMOVED***
	if fo.Collation != nil ***REMOVED***
		op = op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	***REMOVED***
	if fo.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, fo.Comment, true, "comment")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	if fo.MaxTime != nil ***REMOVED***
		op = op.MaxTimeMS(int64(*fo.MaxTime / time.Millisecond))
	***REMOVED***
	if fo.Projection != nil ***REMOVED***
		proj, err := transformBsoncoreDocument(coll.registry, fo.Projection, true, "projection")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Fields(proj)
	***REMOVED***
	if fo.ReturnDocument != nil ***REMOVED***
		op = op.NewDocument(*fo.ReturnDocument == options.After)
	***REMOVED***
	if fo.Sort != nil ***REMOVED***
		sort, err := transformBsoncoreDocument(coll.registry, fo.Sort, false, "sort")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Sort(sort)
	***REMOVED***
	if fo.Upsert != nil ***REMOVED***
		op = op.Upsert(*fo.Upsert)
	***REMOVED***
	if fo.Hint != nil ***REMOVED***
		hint, err := transformValue(coll.registry, fo.Hint, false, "hint")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Hint(hint)
	***REMOVED***
	if fo.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(coll.registry, fo.Let, true, "let")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Let(let)
	***REMOVED***

	return coll.findAndModify(ctx, op)
***REMOVED***

// FindOneAndUpdate executes a findAndModify command to update at most one document in the collection and returns the
// document as it appeared before updating.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// updated. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The update parameter must be a document containing update operators
// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be made
// to the selected document. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.FindOneAndUpdateOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/findAndModify/.
func (coll *Collection) FindOneAndUpdate(ctx context.Context, filter interface***REMOVED******REMOVED***,
	update interface***REMOVED******REMOVED***, opts ...*options.FindOneAndUpdateOptions) *SingleResult ***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	f, err := transformBsoncoreDocument(coll.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***

	fo := options.MergeFindOneAndUpdateOptions(opts...)
	op := operation.NewFindAndModify(f).ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout)

	u, err := transformUpdateValue(coll.registry, update, true)
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***
	op = op.Update(u)

	if fo.ArrayFilters != nil ***REMOVED***
		filtersDoc, err := fo.ArrayFilters.ToArrayDocument()
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.ArrayFilters(bsoncore.Document(filtersDoc))
	***REMOVED***
	if fo.BypassDocumentValidation != nil && *fo.BypassDocumentValidation ***REMOVED***
		op = op.BypassDocumentValidation(*fo.BypassDocumentValidation)
	***REMOVED***
	if fo.Collation != nil ***REMOVED***
		op = op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	***REMOVED***
	if fo.Comment != nil ***REMOVED***
		comment, err := transformValue(coll.registry, fo.Comment, true, "comment")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Comment(comment)
	***REMOVED***
	if fo.MaxTime != nil ***REMOVED***
		op = op.MaxTimeMS(int64(*fo.MaxTime / time.Millisecond))
	***REMOVED***
	if fo.Projection != nil ***REMOVED***
		proj, err := transformBsoncoreDocument(coll.registry, fo.Projection, true, "projection")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Fields(proj)
	***REMOVED***
	if fo.ReturnDocument != nil ***REMOVED***
		op = op.NewDocument(*fo.ReturnDocument == options.After)
	***REMOVED***
	if fo.Sort != nil ***REMOVED***
		sort, err := transformBsoncoreDocument(coll.registry, fo.Sort, false, "sort")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Sort(sort)
	***REMOVED***
	if fo.Upsert != nil ***REMOVED***
		op = op.Upsert(*fo.Upsert)
	***REMOVED***
	if fo.Hint != nil ***REMOVED***
		hint, err := transformValue(coll.registry, fo.Hint, false, "hint")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Hint(hint)
	***REMOVED***
	if fo.Let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(coll.registry, fo.Let, true, "let")
		if err != nil ***REMOVED***
			return &SingleResult***REMOVED***err: err***REMOVED***
		***REMOVED***
		op = op.Let(let)
	***REMOVED***

	return coll.findAndModify(ctx, op)
***REMOVED***

// Watch returns a change stream for all changes on the corresponding collection. See
// https://www.mongodb.com/docs/manual/changeStreams/ for more information about change streams.
//
// The Collection must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
// nil but can be empty. The stage documents must all be non-nil. See https://www.mongodb.com/docs/manual/changeStreams/ for
// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
// mongo.Pipeline***REMOVED******REMOVED*** type can be used.
//
// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
// documentation).
func (coll *Collection) Watch(ctx context.Context, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) ***REMOVED***

	csConfig := changeStreamConfig***REMOVED***
		readConcern:    coll.readConcern,
		readPreference: coll.readPreference,
		client:         coll.client,
		registry:       coll.registry,
		streamType:     CollectionStream,
		collectionName: coll.Name(),
		databaseName:   coll.db.Name(),
		crypt:          coll.client.cryptFLE,
	***REMOVED***
	return newChangeStream(ctx, csConfig, pipeline, opts...)
***REMOVED***

// Indexes returns an IndexView instance that can be used to perform operations on the indexes for the collection.
func (coll *Collection) Indexes() IndexView ***REMOVED***
	return IndexView***REMOVED***coll: coll***REMOVED***
***REMOVED***

// Drop drops the collection on the server. This method ignores "namespace not found" errors so it is safe to drop
// a collection that does not exist on the server.
func (coll *Collection) Drop(ctx context.Context) error ***REMOVED***
	// Follow Client-Side Encryption specification to check for encryptedFields.
	// Drop does not have an encryptedFields option. See: GODRIVER-2413.
	// Check for encryptedFields from the client EncryptedFieldsMap.
	// Check for encryptedFields from the server if EncryptedFieldsMap is set.
	ef := coll.db.getEncryptedFieldsFromMap(coll.name)
	if ef == nil && coll.db.client.encryptedFieldsMap != nil ***REMOVED***
		var err error
		if ef, err = coll.db.getEncryptedFieldsFromServer(ctx, coll.name); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if ef != nil ***REMOVED***
		return coll.dropEncryptedCollection(ctx, ef)
	***REMOVED***

	return coll.drop(ctx)
***REMOVED***

// dropEncryptedCollection drops a collection with EncryptedFields.
func (coll *Collection) dropEncryptedCollection(ctx context.Context, ef interface***REMOVED******REMOVED***) error ***REMOVED***
	efBSON, err := transformBsoncoreDocument(coll.registry, ef, true /* mapAllowed */, "encryptedFields")
	if err != nil ***REMOVED***
		return fmt.Errorf("error transforming document: %v", err)
	***REMOVED***

	// Drop the three encryption-related, associated collections: `escCollection`, `eccCollection` and `ecocCollection`.
	// Drop ESCCollection.
	escCollection, err := internal.GetEncryptedStateCollectionName(efBSON, coll.name, internal.EncryptedStateCollection)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := coll.db.Collection(escCollection).drop(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Drop ECCCollection.
	eccCollection, err := internal.GetEncryptedStateCollectionName(efBSON, coll.name, internal.EncryptedCacheCollection)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := coll.db.Collection(eccCollection).drop(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Drop ECOCCollection.
	ecocCollection, err := internal.GetEncryptedStateCollectionName(efBSON, coll.name, internal.EncryptedCompactionCollection)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := coll.db.Collection(ecocCollection).drop(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Drop the data collection.
	if err := coll.drop(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// drop drops a collection without EncryptedFields.
func (coll *Collection) drop(ctx context.Context) error ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err := coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	wc := coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewDropCollection().
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).
		ServerAPI(coll.client.serverAPI)
	err = op.Execute(ctx)

	// ignore namespace not found erorrs
	driverErr, ok := err.(driver.Error)
	if !ok || (ok && !driverErr.NamespaceNotFound()) ***REMOVED***
		return replaceErrors(err)
	***REMOVED***
	return nil
***REMOVED***

// makePinnedSelector makes a selector for a pinned session with a pinned server. Will attempt to do server selection on
// the pinned server but if that fails it will go through a list of default selectors
func makePinnedSelector(sess *session.Client, defaultSelector description.ServerSelector) description.ServerSelectorFunc ***REMOVED***
	return func(t description.Topology, svrs []description.Server) ([]description.Server, error) ***REMOVED***
		if sess != nil && sess.PinnedServer != nil ***REMOVED***
			// If there is a pinned server, try to find it in the list of candidates.
			for _, candidate := range svrs ***REMOVED***
				if candidate.Addr == sess.PinnedServer.Addr ***REMOVED***
					return []description.Server***REMOVED***candidate***REMOVED***, nil
				***REMOVED***
			***REMOVED***

			return nil, nil
		***REMOVED***

		return defaultSelector.SelectServer(t, svrs)
	***REMOVED***
***REMOVED***

func makeReadPrefSelector(sess *session.Client, selector description.ServerSelector, localThreshold time.Duration) description.ServerSelectorFunc ***REMOVED***
	if sess != nil && sess.TransactionRunning() ***REMOVED***
		selector = description.CompositeSelector([]description.ServerSelector***REMOVED***
			description.ReadPrefSelector(sess.CurrentRp),
			description.LatencySelector(localThreshold),
		***REMOVED***)
	***REMOVED***

	return makePinnedSelector(sess, selector)
***REMOVED***

func makeOutputAggregateSelector(sess *session.Client, rp *readpref.ReadPref, localThreshold time.Duration) description.ServerSelectorFunc ***REMOVED***
	if sess != nil && sess.TransactionRunning() ***REMOVED***
		// Use current transaction's read preference if available
		rp = sess.CurrentRp
	***REMOVED***

	selector := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.OutputAggregateSelector(rp),
		description.LatencySelector(localThreshold),
	***REMOVED***)
	return makePinnedSelector(sess, selector)
***REMOVED***
