// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// ErrInvalidIndexValue is returned if an index is created with a keys document that has a value that is not a number
// or string.
var ErrInvalidIndexValue = errors.New("invalid index value")

// ErrNonStringIndexName is returned if an index is created with a name that is not a string.
var ErrNonStringIndexName = errors.New("index name must be a string")

// ErrMultipleIndexDrop is returned if multiple indexes would be dropped from a call to IndexView.DropOne.
var ErrMultipleIndexDrop = errors.New("multiple indexes would be dropped")

// IndexView is a type that can be used to create, drop, and list indexes on a collection. An IndexView for a collection
// can be created by a call to Collection.Indexes().
type IndexView struct ***REMOVED***
	coll *Collection
***REMOVED***

// IndexModel represents a new index to be created.
type IndexModel struct ***REMOVED***
	// A document describing which keys should be used for the index. It cannot be nil. This must be an order-preserving
	// type such as bson.D. Map types such as bson.M are not valid. See https://www.mongodb.com/docs/manual/indexes/#indexes
	// for examples of valid documents.
	Keys interface***REMOVED******REMOVED***

	// The options to use to create the index.
	Options *options.IndexOptions
***REMOVED***

func isNamespaceNotFoundError(err error) bool ***REMOVED***
	if de, ok := err.(driver.Error); ok ***REMOVED***
		return de.Code == 26
	***REMOVED***
	return false
***REMOVED***

// List executes a listIndexes command and returns a cursor over the indexes in the collection.
//
// The opts parameter can be used to specify options for this operation (see the options.ListIndexesOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listIndexes/.
func (iv IndexView) List(ctx context.Context, opts ...*options.ListIndexesOptions) (*Cursor, error) ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && iv.coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(iv.coll.client.sessionPool, iv.coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	err := iv.coll.client.validSession(sess)
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, err
	***REMOVED***

	selector := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(iv.coll.client.localThreshold),
	***REMOVED***)
	selector = makeReadPrefSelector(sess, selector, iv.coll.client.localThreshold)
	op := operation.NewListIndexes().
		Session(sess).CommandMonitor(iv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(iv.coll.client.clock).
		Database(iv.coll.db.name).Collection(iv.coll.name).
		Deployment(iv.coll.client.deployment).ServerAPI(iv.coll.client.serverAPI).
		Timeout(iv.coll.client.timeout)

	cursorOpts := iv.coll.client.createBaseCursorOptions()
	lio := options.MergeListIndexesOptions(opts...)
	if lio.BatchSize != nil ***REMOVED***
		op = op.BatchSize(*lio.BatchSize)
		cursorOpts.BatchSize = *lio.BatchSize
	***REMOVED***
	if lio.MaxTime != nil ***REMOVED***
		op = op.MaxTimeMS(int64(*lio.MaxTime / time.Millisecond))
	***REMOVED***
	retry := driver.RetryNone
	if iv.coll.client.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		// for namespaceNotFound errors, return an empty cursor and do not throw an error
		closeImplicitSession(sess)
		if isNamespaceNotFoundError(err) ***REMOVED***
			return newEmptyCursor(), nil
		***REMOVED***

		return nil, replaceErrors(err)
	***REMOVED***

	bc, err := op.Result(cursorOpts)
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	***REMOVED***
	cursor, err := newCursorWithSession(bc, iv.coll.registry, sess)
	return cursor, replaceErrors(err)
***REMOVED***

// ListSpecifications executes a List command and returns a slice of returned IndexSpecifications
func (iv IndexView) ListSpecifications(ctx context.Context, opts ...*options.ListIndexesOptions) ([]*IndexSpecification, error) ***REMOVED***
	cursor, err := iv.List(ctx, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var results []*IndexSpecification
	err = cursor.All(ctx, &results)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ns := iv.coll.db.Name() + "." + iv.coll.Name()
	for _, res := range results ***REMOVED***
		// Pre-4.4 servers report a namespace in their responses, so we only set Namespace manually if it was not in
		// the response.
		res.Namespace = ns
	***REMOVED***

	return results, nil
***REMOVED***

// CreateOne executes a createIndexes command to create an index on the collection and returns the name of the new
// index. See the IndexView.CreateMany documentation for more information and an example.
func (iv IndexView) CreateOne(ctx context.Context, model IndexModel, opts ...*options.CreateIndexesOptions) (string, error) ***REMOVED***
	names, err := iv.CreateMany(ctx, []IndexModel***REMOVED***model***REMOVED***, opts...)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return names[0], nil
***REMOVED***

// CreateMany executes a createIndexes command to create multiple indexes on the collection and returns the names of
// the new indexes.
//
// For each IndexModel in the models parameter, the index name can be specified via the Options field. If a name is not
// given, it will be generated from the Keys document.
//
// The opts parameter can be used to specify options for this operation (see the options.CreateIndexesOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/createIndexes/.
func (iv IndexView) CreateMany(ctx context.Context, models []IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) ***REMOVED***
	names := make([]string, 0, len(models))

	var indexes bsoncore.Document
	aidx, indexes := bsoncore.AppendArrayStart(indexes)

	for i, model := range models ***REMOVED***
		if model.Keys == nil ***REMOVED***
			return nil, fmt.Errorf("index model keys cannot be nil")
		***REMOVED***

		keys, err := transformBsoncoreDocument(iv.coll.registry, model.Keys, false, "keys")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		name, err := getOrGenerateIndexName(keys, model)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		names = append(names, name)

		var iidx int32
		iidx, indexes = bsoncore.AppendDocumentElementStart(indexes, strconv.Itoa(i))
		indexes = bsoncore.AppendDocumentElement(indexes, "key", keys)

		if model.Options == nil ***REMOVED***
			model.Options = options.Index()
		***REMOVED***
		model.Options.SetName(name)

		optsDoc, err := iv.createOptionsDoc(model.Options)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		indexes = bsoncore.AppendDocument(indexes, optsDoc)

		indexes, err = bsoncore.AppendDocumentEnd(indexes, iidx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	indexes, err := bsoncore.AppendArrayEnd(indexes, aidx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(ctx)

	if sess == nil && iv.coll.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(iv.coll.client.sessionPool, iv.coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = iv.coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wc := iv.coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, iv.coll.writeSelector)

	option := options.MergeCreateIndexesOptions(opts...)

	op := operation.NewCreateIndexes(indexes).
		Session(sess).WriteConcern(wc).ClusterClock(iv.coll.client.clock).
		Database(iv.coll.db.name).Collection(iv.coll.name).CommandMonitor(iv.coll.client.monitor).
		Deployment(iv.coll.client.deployment).ServerSelector(selector).ServerAPI(iv.coll.client.serverAPI).
		Timeout(iv.coll.client.timeout)

	if option.MaxTime != nil ***REMOVED***
		op.MaxTimeMS(int64(*option.MaxTime / time.Millisecond))
	***REMOVED***
	if option.CommitQuorum != nil ***REMOVED***
		commitQuorum, err := transformValue(iv.coll.registry, option.CommitQuorum, true, "commitQuorum")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		op.CommitQuorum(commitQuorum)
	***REMOVED***

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		_, err = processWriteError(err)
		return nil, err
	***REMOVED***

	return names, nil
***REMOVED***

func (iv IndexView) createOptionsDoc(opts *options.IndexOptions) (bsoncore.Document, error) ***REMOVED***
	optsDoc := bsoncore.Document***REMOVED******REMOVED***
	if opts.Background != nil ***REMOVED***
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "background", *opts.Background)
	***REMOVED***
	if opts.ExpireAfterSeconds != nil ***REMOVED***
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "expireAfterSeconds", *opts.ExpireAfterSeconds)
	***REMOVED***
	if opts.Name != nil ***REMOVED***
		optsDoc = bsoncore.AppendStringElement(optsDoc, "name", *opts.Name)
	***REMOVED***
	if opts.Sparse != nil ***REMOVED***
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "sparse", *opts.Sparse)
	***REMOVED***
	if opts.StorageEngine != nil ***REMOVED***
		doc, err := transformBsoncoreDocument(iv.coll.registry, opts.StorageEngine, true, "storageEngine")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "storageEngine", doc)
	***REMOVED***
	if opts.Unique != nil ***REMOVED***
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "unique", *opts.Unique)
	***REMOVED***
	if opts.Version != nil ***REMOVED***
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "v", *opts.Version)
	***REMOVED***
	if opts.DefaultLanguage != nil ***REMOVED***
		optsDoc = bsoncore.AppendStringElement(optsDoc, "default_language", *opts.DefaultLanguage)
	***REMOVED***
	if opts.LanguageOverride != nil ***REMOVED***
		optsDoc = bsoncore.AppendStringElement(optsDoc, "language_override", *opts.LanguageOverride)
	***REMOVED***
	if opts.TextVersion != nil ***REMOVED***
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "textIndexVersion", *opts.TextVersion)
	***REMOVED***
	if opts.Weights != nil ***REMOVED***
		doc, err := transformBsoncoreDocument(iv.coll.registry, opts.Weights, true, "weights")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "weights", doc)
	***REMOVED***
	if opts.SphereVersion != nil ***REMOVED***
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "2dsphereIndexVersion", *opts.SphereVersion)
	***REMOVED***
	if opts.Bits != nil ***REMOVED***
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "bits", *opts.Bits)
	***REMOVED***
	if opts.Max != nil ***REMOVED***
		optsDoc = bsoncore.AppendDoubleElement(optsDoc, "max", *opts.Max)
	***REMOVED***
	if opts.Min != nil ***REMOVED***
		optsDoc = bsoncore.AppendDoubleElement(optsDoc, "min", *opts.Min)
	***REMOVED***
	if opts.BucketSize != nil ***REMOVED***
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "bucketSize", *opts.BucketSize)
	***REMOVED***
	if opts.PartialFilterExpression != nil ***REMOVED***
		doc, err := transformBsoncoreDocument(iv.coll.registry, opts.PartialFilterExpression, true, "partialFilterExpression")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "partialFilterExpression", doc)
	***REMOVED***
	if opts.Collation != nil ***REMOVED***
		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "collation", bsoncore.Document(opts.Collation.ToDocument()))
	***REMOVED***
	if opts.WildcardProjection != nil ***REMOVED***
		doc, err := transformBsoncoreDocument(iv.coll.registry, opts.WildcardProjection, true, "wildcardProjection")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "wildcardProjection", doc)
	***REMOVED***
	if opts.Hidden != nil ***REMOVED***
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "hidden", *opts.Hidden)
	***REMOVED***

	return optsDoc, nil
***REMOVED***

func (iv IndexView) drop(ctx context.Context, name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && iv.coll.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(iv.coll.client.sessionPool, iv.coll.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err := iv.coll.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wc := iv.coll.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, iv.coll.writeSelector)

	dio := options.MergeDropIndexesOptions(opts...)
	op := operation.NewDropIndexes(name).
		Session(sess).WriteConcern(wc).CommandMonitor(iv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(iv.coll.client.clock).
		Database(iv.coll.db.name).Collection(iv.coll.name).
		Deployment(iv.coll.client.deployment).ServerAPI(iv.coll.client.serverAPI).
		Timeout(iv.coll.client.timeout)
	if dio.MaxTime != nil ***REMOVED***
		op.MaxTimeMS(int64(*dio.MaxTime / time.Millisecond))
	***REMOVED***

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***

	// TODO: it's weird to return a bson.Raw here because we have to convert the result back to BSON
	ridx, res := bsoncore.AppendDocumentStart(nil)
	res = bsoncore.AppendInt32Element(res, "nIndexesWas", op.Result().NIndexesWas)
	res, _ = bsoncore.AppendDocumentEnd(res, ridx)
	return res, nil
***REMOVED***

// DropOne executes a dropIndexes operation to drop an index on the collection. If the operation succeeds, this returns
// a BSON document in the form ***REMOVED***nIndexesWas: <int32>***REMOVED***. The "nIndexesWas" field in the response contains the number of
// indexes that existed prior to the drop.
//
// The name parameter should be the name of the index to drop. If the name is "*", ErrMultipleIndexDrop will be returned
// without running the command because doing so would drop all indexes.
//
// The opts parameter can be used to specify options for this operation (see the options.DropIndexesOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/dropIndexes/.
func (iv IndexView) DropOne(ctx context.Context, name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) ***REMOVED***
	if name == "*" ***REMOVED***
		return nil, ErrMultipleIndexDrop
	***REMOVED***

	return iv.drop(ctx, name, opts...)
***REMOVED***

// DropAll executes a dropIndexes operation to drop all indexes on the collection. If the operation succeeds, this
// returns a BSON document in the form ***REMOVED***nIndexesWas: <int32>***REMOVED***. The "nIndexesWas" field in the response contains the
// number of indexes that existed prior to the drop.
//
// The opts parameter can be used to specify options for this operation (see the options.DropIndexesOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/dropIndexes/.
func (iv IndexView) DropAll(ctx context.Context, opts ...*options.DropIndexesOptions) (bson.Raw, error) ***REMOVED***
	return iv.drop(ctx, "*", opts...)
***REMOVED***

func getOrGenerateIndexName(keySpecDocument bsoncore.Document, model IndexModel) (string, error) ***REMOVED***
	if model.Options != nil && model.Options.Name != nil ***REMOVED***
		return *model.Options.Name, nil
	***REMOVED***

	name := bytes.NewBufferString("")
	first := true

	elems, err := keySpecDocument.Elements()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	for _, elem := range elems ***REMOVED***
		if !first ***REMOVED***
			_, err := name.WriteRune('_')
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***

		_, err := name.WriteString(elem.Key())
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		_, err = name.WriteRune('_')
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		var value string

		bsonValue := elem.Value()
		switch bsonValue.Type ***REMOVED***
		case bsontype.Int32:
			value = fmt.Sprintf("%d", bsonValue.Int32())
		case bsontype.Int64:
			value = fmt.Sprintf("%d", bsonValue.Int64())
		case bsontype.String:
			value = bsonValue.StringValue()
		default:
			return "", ErrInvalidIndexValue
		***REMOVED***

		_, err = name.WriteString(value)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		first = false
	***REMOVED***

	return name.String(), nil
***REMOVED***
