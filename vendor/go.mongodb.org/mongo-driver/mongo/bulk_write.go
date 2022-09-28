// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

type bulkWriteBatch struct ***REMOVED***
	models   []WriteModel
	canRetry bool
	indexes  []int
***REMOVED***

// bulkWrite perfoms a bulkwrite operation
type bulkWrite struct ***REMOVED***
	comment                  interface***REMOVED******REMOVED***
	ordered                  *bool
	bypassDocumentValidation *bool
	models                   []WriteModel
	session                  *session.Client
	collection               *Collection
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	result                   BulkWriteResult
	let                      interface***REMOVED******REMOVED***
***REMOVED***

func (bw *bulkWrite) execute(ctx context.Context) error ***REMOVED***
	ordered := true
	if bw.ordered != nil ***REMOVED***
		ordered = *bw.ordered
	***REMOVED***

	batches := createBatches(bw.models, ordered)
	bw.result = BulkWriteResult***REMOVED***
		UpsertedIDs: make(map[int64]interface***REMOVED******REMOVED***),
	***REMOVED***

	bwErr := BulkWriteException***REMOVED***
		WriteErrors: make([]BulkWriteError, 0),
	***REMOVED***

	var lastErr error
	continueOnError := !ordered
	for _, batch := range batches ***REMOVED***
		if len(batch.models) == 0 ***REMOVED***
			continue
		***REMOVED***

		batchRes, batchErr, err := bw.runBatch(ctx, batch)

		bw.mergeResults(batchRes)

		bwErr.WriteConcernError = batchErr.WriteConcernError
		bwErr.Labels = append(bwErr.Labels, batchErr.Labels...)

		bwErr.WriteErrors = append(bwErr.WriteErrors, batchErr.WriteErrors...)

		commandErrorOccurred := err != nil && err != driver.ErrUnacknowledgedWrite
		writeErrorOccurred := len(batchErr.WriteErrors) > 0 || batchErr.WriteConcernError != nil
		if !continueOnError && (commandErrorOccurred || writeErrorOccurred) ***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			return bwErr
		***REMOVED***

		if err != nil ***REMOVED***
			lastErr = err
		***REMOVED***
	***REMOVED***

	bw.result.MatchedCount -= bw.result.UpsertedCount
	if lastErr != nil ***REMOVED***
		_, lastErr = processWriteError(lastErr)
		return lastErr
	***REMOVED***
	if len(bwErr.WriteErrors) > 0 || bwErr.WriteConcernError != nil ***REMOVED***
		return bwErr
	***REMOVED***
	return nil
***REMOVED***

func (bw *bulkWrite) runBatch(ctx context.Context, batch bulkWriteBatch) (BulkWriteResult, BulkWriteException, error) ***REMOVED***
	batchRes := BulkWriteResult***REMOVED***
		UpsertedIDs: make(map[int64]interface***REMOVED******REMOVED***),
	***REMOVED***
	batchErr := BulkWriteException***REMOVED******REMOVED***

	var writeErrors []driver.WriteError
	switch batch.models[0].(type) ***REMOVED***
	case *InsertOneModel:
		res, err := bw.runInsert(ctx, batch)
		if err != nil ***REMOVED***
			writeErr, ok := err.(driver.WriteCommandError)
			if !ok ***REMOVED***
				return BulkWriteResult***REMOVED******REMOVED***, batchErr, err
			***REMOVED***
			writeErrors = writeErr.WriteErrors
			batchErr.Labels = writeErr.Labels
			batchErr.WriteConcernError = convertDriverWriteConcernError(writeErr.WriteConcernError)
		***REMOVED***
		batchRes.InsertedCount = res.N
	case *DeleteOneModel, *DeleteManyModel:
		res, err := bw.runDelete(ctx, batch)
		if err != nil ***REMOVED***
			writeErr, ok := err.(driver.WriteCommandError)
			if !ok ***REMOVED***
				return BulkWriteResult***REMOVED******REMOVED***, batchErr, err
			***REMOVED***
			writeErrors = writeErr.WriteErrors
			batchErr.Labels = writeErr.Labels
			batchErr.WriteConcernError = convertDriverWriteConcernError(writeErr.WriteConcernError)
		***REMOVED***
		batchRes.DeletedCount = res.N
	case *ReplaceOneModel, *UpdateOneModel, *UpdateManyModel:
		res, err := bw.runUpdate(ctx, batch)
		if err != nil ***REMOVED***
			writeErr, ok := err.(driver.WriteCommandError)
			if !ok ***REMOVED***
				return BulkWriteResult***REMOVED******REMOVED***, batchErr, err
			***REMOVED***
			writeErrors = writeErr.WriteErrors
			batchErr.Labels = writeErr.Labels
			batchErr.WriteConcernError = convertDriverWriteConcernError(writeErr.WriteConcernError)
		***REMOVED***
		batchRes.MatchedCount = res.N
		batchRes.ModifiedCount = res.NModified
		batchRes.UpsertedCount = int64(len(res.Upserted))
		for _, upsert := range res.Upserted ***REMOVED***
			batchRes.UpsertedIDs[int64(batch.indexes[upsert.Index])] = upsert.ID
		***REMOVED***
	***REMOVED***

	batchErr.WriteErrors = make([]BulkWriteError, 0, len(writeErrors))
	convWriteErrors := writeErrorsFromDriverWriteErrors(writeErrors)
	for _, we := range convWriteErrors ***REMOVED***
		request := batch.models[we.Index]
		we.Index = batch.indexes[we.Index]
		batchErr.WriteErrors = append(batchErr.WriteErrors, BulkWriteError***REMOVED***
			WriteError: we,
			Request:    request,
		***REMOVED***)
	***REMOVED***
	return batchRes, batchErr, nil
***REMOVED***

func (bw *bulkWrite) runInsert(ctx context.Context, batch bulkWriteBatch) (operation.InsertResult, error) ***REMOVED***
	docs := make([]bsoncore.Document, len(batch.models))
	var i int
	for _, model := range batch.models ***REMOVED***
		converted := model.(*InsertOneModel)
		doc, _, err := transformAndEnsureID(bw.collection.registry, converted.Document)
		if err != nil ***REMOVED***
			return operation.InsertResult***REMOVED******REMOVED***, err
		***REMOVED***

		docs[i] = doc
		i++
	***REMOVED***

	op := operation.NewInsert(docs...).
		Session(bw.session).WriteConcern(bw.writeConcern).CommandMonitor(bw.collection.client.monitor).
		ServerSelector(bw.selector).ClusterClock(bw.collection.client.clock).
		Database(bw.collection.db.name).Collection(bw.collection.name).
		Deployment(bw.collection.client.deployment).Crypt(bw.collection.client.cryptFLE).
		ServerAPI(bw.collection.client.serverAPI).Timeout(bw.collection.client.timeout)
	if bw.comment != nil ***REMOVED***
		comment, err := transformValue(bw.collection.registry, bw.comment, true, "comment")
		if err != nil ***REMOVED***
			return op.Result(), err
		***REMOVED***
		op.Comment(comment)
	***REMOVED***
	if bw.bypassDocumentValidation != nil && *bw.bypassDocumentValidation ***REMOVED***
		op = op.BypassDocumentValidation(*bw.bypassDocumentValidation)
	***REMOVED***
	if bw.ordered != nil ***REMOVED***
		op = op.Ordered(*bw.ordered)
	***REMOVED***

	retry := driver.RetryNone
	if bw.collection.client.retryWrites && batch.canRetry ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err := op.Execute(ctx)

	return op.Result(), err
***REMOVED***

func (bw *bulkWrite) runDelete(ctx context.Context, batch bulkWriteBatch) (operation.DeleteResult, error) ***REMOVED***
	docs := make([]bsoncore.Document, len(batch.models))
	var i int
	var hasHint bool

	for _, model := range batch.models ***REMOVED***
		var doc bsoncore.Document
		var err error

		switch converted := model.(type) ***REMOVED***
		case *DeleteOneModel:
			doc, err = createDeleteDoc(converted.Filter, converted.Collation, converted.Hint, true, bw.collection.registry)
			hasHint = hasHint || (converted.Hint != nil)
		case *DeleteManyModel:
			doc, err = createDeleteDoc(converted.Filter, converted.Collation, converted.Hint, false, bw.collection.registry)
			hasHint = hasHint || (converted.Hint != nil)
		***REMOVED***

		if err != nil ***REMOVED***
			return operation.DeleteResult***REMOVED******REMOVED***, err
		***REMOVED***

		docs[i] = doc
		i++
	***REMOVED***

	op := operation.NewDelete(docs...).
		Session(bw.session).WriteConcern(bw.writeConcern).CommandMonitor(bw.collection.client.monitor).
		ServerSelector(bw.selector).ClusterClock(bw.collection.client.clock).
		Database(bw.collection.db.name).Collection(bw.collection.name).
		Deployment(bw.collection.client.deployment).Crypt(bw.collection.client.cryptFLE).Hint(hasHint).
		ServerAPI(bw.collection.client.serverAPI).Timeout(bw.collection.client.timeout)
	if bw.comment != nil ***REMOVED***
		comment, err := transformValue(bw.collection.registry, bw.comment, true, "comment")
		if err != nil ***REMOVED***
			return op.Result(), err
		***REMOVED***
		op.Comment(comment)
	***REMOVED***
	if bw.let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(bw.collection.registry, bw.let, true, "let")
		if err != nil ***REMOVED***
			return operation.DeleteResult***REMOVED******REMOVED***, err
		***REMOVED***
		op = op.Let(let)
	***REMOVED***
	if bw.ordered != nil ***REMOVED***
		op = op.Ordered(*bw.ordered)
	***REMOVED***
	retry := driver.RetryNone
	if bw.collection.client.retryWrites && batch.canRetry ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err := op.Execute(ctx)

	return op.Result(), err
***REMOVED***

func createDeleteDoc(filter interface***REMOVED******REMOVED***, collation *options.Collation, hint interface***REMOVED******REMOVED***, deleteOne bool,
	registry *bsoncodec.Registry) (bsoncore.Document, error) ***REMOVED***

	f, err := transformBsoncoreDocument(registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var limit int32
	if deleteOne ***REMOVED***
		limit = 1
	***REMOVED***
	didx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendDocumentElement(doc, "q", f)
	doc = bsoncore.AppendInt32Element(doc, "limit", limit)
	if collation != nil ***REMOVED***
		doc = bsoncore.AppendDocumentElement(doc, "collation", collation.ToDocument())
	***REMOVED***
	if hint != nil ***REMOVED***
		hintVal, err := transformValue(registry, hint, false, "hint")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		doc = bsoncore.AppendValueElement(doc, "hint", hintVal)
	***REMOVED***
	doc, _ = bsoncore.AppendDocumentEnd(doc, didx)

	return doc, nil
***REMOVED***

func (bw *bulkWrite) runUpdate(ctx context.Context, batch bulkWriteBatch) (operation.UpdateResult, error) ***REMOVED***
	docs := make([]bsoncore.Document, len(batch.models))
	var hasHint bool
	var hasArrayFilters bool
	for i, model := range batch.models ***REMOVED***
		var doc bsoncore.Document
		var err error

		switch converted := model.(type) ***REMOVED***
		case *ReplaceOneModel:
			doc, err = createUpdateDoc(converted.Filter, converted.Replacement, converted.Hint, nil, converted.Collation, converted.Upsert, false,
				false, bw.collection.registry)
			hasHint = hasHint || (converted.Hint != nil)
		case *UpdateOneModel:
			doc, err = createUpdateDoc(converted.Filter, converted.Update, converted.Hint, converted.ArrayFilters, converted.Collation, converted.Upsert, false,
				true, bw.collection.registry)
			hasHint = hasHint || (converted.Hint != nil)
			hasArrayFilters = hasArrayFilters || (converted.ArrayFilters != nil)
		case *UpdateManyModel:
			doc, err = createUpdateDoc(converted.Filter, converted.Update, converted.Hint, converted.ArrayFilters, converted.Collation, converted.Upsert, true,
				true, bw.collection.registry)
			hasHint = hasHint || (converted.Hint != nil)
			hasArrayFilters = hasArrayFilters || (converted.ArrayFilters != nil)
		***REMOVED***
		if err != nil ***REMOVED***
			return operation.UpdateResult***REMOVED******REMOVED***, err
		***REMOVED***

		docs[i] = doc
	***REMOVED***

	op := operation.NewUpdate(docs...).
		Session(bw.session).WriteConcern(bw.writeConcern).CommandMonitor(bw.collection.client.monitor).
		ServerSelector(bw.selector).ClusterClock(bw.collection.client.clock).
		Database(bw.collection.db.name).Collection(bw.collection.name).
		Deployment(bw.collection.client.deployment).Crypt(bw.collection.client.cryptFLE).Hint(hasHint).
		ArrayFilters(hasArrayFilters).ServerAPI(bw.collection.client.serverAPI).Timeout(bw.collection.client.timeout)
	if bw.comment != nil ***REMOVED***
		comment, err := transformValue(bw.collection.registry, bw.comment, true, "comment")
		if err != nil ***REMOVED***
			return op.Result(), err
		***REMOVED***
		op.Comment(comment)
	***REMOVED***
	if bw.let != nil ***REMOVED***
		let, err := transformBsoncoreDocument(bw.collection.registry, bw.let, true, "let")
		if err != nil ***REMOVED***
			return operation.UpdateResult***REMOVED******REMOVED***, err
		***REMOVED***
		op = op.Let(let)
	***REMOVED***
	if bw.ordered != nil ***REMOVED***
		op = op.Ordered(*bw.ordered)
	***REMOVED***
	if bw.bypassDocumentValidation != nil && *bw.bypassDocumentValidation ***REMOVED***
		op = op.BypassDocumentValidation(*bw.bypassDocumentValidation)
	***REMOVED***
	retry := driver.RetryNone
	if bw.collection.client.retryWrites && batch.canRetry ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err := op.Execute(ctx)

	return op.Result(), err
***REMOVED***
func createUpdateDoc(
	filter interface***REMOVED******REMOVED***,
	update interface***REMOVED******REMOVED***,
	hint interface***REMOVED******REMOVED***,
	arrayFilters *options.ArrayFilters,
	collation *options.Collation,
	upsert *bool,
	multi bool,
	checkDollarKey bool,
	registry *bsoncodec.Registry,
) (bsoncore.Document, error) ***REMOVED***
	f, err := transformBsoncoreDocument(registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	uidx, updateDoc := bsoncore.AppendDocumentStart(nil)
	updateDoc = bsoncore.AppendDocumentElement(updateDoc, "q", f)

	u, err := transformUpdateValue(registry, update, checkDollarKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	updateDoc = bsoncore.AppendValueElement(updateDoc, "u", u)

	if multi ***REMOVED***
		updateDoc = bsoncore.AppendBooleanElement(updateDoc, "multi", multi)
	***REMOVED***

	if arrayFilters != nil ***REMOVED***
		arr, err := arrayFilters.ToArrayDocument()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		updateDoc = bsoncore.AppendArrayElement(updateDoc, "arrayFilters", arr)
	***REMOVED***

	if collation != nil ***REMOVED***
		updateDoc = bsoncore.AppendDocumentElement(updateDoc, "collation", bsoncore.Document(collation.ToDocument()))
	***REMOVED***

	if upsert != nil ***REMOVED***
		updateDoc = bsoncore.AppendBooleanElement(updateDoc, "upsert", *upsert)
	***REMOVED***

	if hint != nil ***REMOVED***
		hintVal, err := transformValue(registry, hint, false, "hint")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		updateDoc = bsoncore.AppendValueElement(updateDoc, "hint", hintVal)
	***REMOVED***

	updateDoc, _ = bsoncore.AppendDocumentEnd(updateDoc, uidx)
	return updateDoc, nil
***REMOVED***

func createBatches(models []WriteModel, ordered bool) []bulkWriteBatch ***REMOVED***
	if ordered ***REMOVED***
		return createOrderedBatches(models)
	***REMOVED***

	batches := make([]bulkWriteBatch, 5)
	batches[insertCommand].canRetry = true
	batches[deleteOneCommand].canRetry = true
	batches[updateOneCommand].canRetry = true

	// TODO(GODRIVER-1157): fix batching once operation retryability is fixed
	for i, model := range models ***REMOVED***
		switch model.(type) ***REMOVED***
		case *InsertOneModel:
			batches[insertCommand].models = append(batches[insertCommand].models, model)
			batches[insertCommand].indexes = append(batches[insertCommand].indexes, i)
		case *DeleteOneModel:
			batches[deleteOneCommand].models = append(batches[deleteOneCommand].models, model)
			batches[deleteOneCommand].indexes = append(batches[deleteOneCommand].indexes, i)
		case *DeleteManyModel:
			batches[deleteManyCommand].models = append(batches[deleteManyCommand].models, model)
			batches[deleteManyCommand].indexes = append(batches[deleteManyCommand].indexes, i)
		case *ReplaceOneModel, *UpdateOneModel:
			batches[updateOneCommand].models = append(batches[updateOneCommand].models, model)
			batches[updateOneCommand].indexes = append(batches[updateOneCommand].indexes, i)
		case *UpdateManyModel:
			batches[updateManyCommand].models = append(batches[updateManyCommand].models, model)
			batches[updateManyCommand].indexes = append(batches[updateManyCommand].indexes, i)
		***REMOVED***
	***REMOVED***

	return batches
***REMOVED***

func createOrderedBatches(models []WriteModel) []bulkWriteBatch ***REMOVED***
	var batches []bulkWriteBatch
	var prevKind writeCommandKind = -1
	i := -1 // batch index

	for ind, model := range models ***REMOVED***
		var createNewBatch bool
		var canRetry bool
		var newKind writeCommandKind

		// TODO(GODRIVER-1157): fix batching once operation retryability is fixed
		switch model.(type) ***REMOVED***
		case *InsertOneModel:
			createNewBatch = prevKind != insertCommand
			canRetry = true
			newKind = insertCommand
		case *DeleteOneModel:
			createNewBatch = prevKind != deleteOneCommand
			canRetry = true
			newKind = deleteOneCommand
		case *DeleteManyModel:
			createNewBatch = prevKind != deleteManyCommand
			newKind = deleteManyCommand
		case *ReplaceOneModel, *UpdateOneModel:
			createNewBatch = prevKind != updateOneCommand
			canRetry = true
			newKind = updateOneCommand
		case *UpdateManyModel:
			createNewBatch = prevKind != updateManyCommand
			newKind = updateManyCommand
		***REMOVED***

		if createNewBatch ***REMOVED***
			batches = append(batches, bulkWriteBatch***REMOVED***
				models:   []WriteModel***REMOVED***model***REMOVED***,
				canRetry: canRetry,
				indexes:  []int***REMOVED***ind***REMOVED***,
			***REMOVED***)
			i++
		***REMOVED*** else ***REMOVED***
			batches[i].models = append(batches[i].models, model)
			if !canRetry ***REMOVED***
				batches[i].canRetry = false // don't make it true if it was already false
			***REMOVED***
			batches[i].indexes = append(batches[i].indexes, ind)
		***REMOVED***

		prevKind = newKind
	***REMOVED***

	return batches
***REMOVED***

func (bw *bulkWrite) mergeResults(newResult BulkWriteResult) ***REMOVED***
	bw.result.InsertedCount += newResult.InsertedCount
	bw.result.MatchedCount += newResult.MatchedCount
	bw.result.ModifiedCount += newResult.ModifiedCount
	bw.result.DeletedCount += newResult.DeletedCount
	bw.result.UpsertedCount += newResult.UpsertedCount

	for index, upsertID := range newResult.UpsertedIDs ***REMOVED***
		bw.result.UpsertedIDs[index] = upsertID
	***REMOVED***
***REMOVED***

// WriteCommandKind is the type of command represented by a Write
type writeCommandKind int8

// These constants represent the valid types of write commands.
const (
	insertCommand writeCommandKind = iota
	updateOneCommand
	updateManyCommand
	deleteOneCommand
	deleteManyCommand
)
