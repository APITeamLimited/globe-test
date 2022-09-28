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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
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

var (
	defaultRunCmdOpts = []*options.RunCmdOptions***REMOVED***options.RunCmd().SetReadPreference(readpref.Primary())***REMOVED***
)

// Database is a handle to a MongoDB database. It is safe for concurrent use by multiple goroutines.
type Database struct ***REMOVED***
	client         *Client
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	registry       *bsoncodec.Registry
***REMOVED***

func newDatabase(client *Client, name string, opts ...*options.DatabaseOptions) *Database ***REMOVED***
	dbOpt := options.MergeDatabaseOptions(opts...)

	rc := client.readConcern
	if dbOpt.ReadConcern != nil ***REMOVED***
		rc = dbOpt.ReadConcern
	***REMOVED***

	rp := client.readPreference
	if dbOpt.ReadPreference != nil ***REMOVED***
		rp = dbOpt.ReadPreference
	***REMOVED***

	wc := client.writeConcern
	if dbOpt.WriteConcern != nil ***REMOVED***
		wc = dbOpt.WriteConcern
	***REMOVED***

	reg := client.registry
	if dbOpt.Registry != nil ***REMOVED***
		reg = dbOpt.Registry
	***REMOVED***

	db := &Database***REMOVED***
		client:         client,
		name:           name,
		readPreference: rp,
		readConcern:    rc,
		writeConcern:   wc,
		registry:       reg,
	***REMOVED***

	db.readSelector = description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(db.readPreference),
		description.LatencySelector(db.client.localThreshold),
	***REMOVED***)

	db.writeSelector = description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.WriteSelector(),
		description.LatencySelector(db.client.localThreshold),
	***REMOVED***)

	return db
***REMOVED***

// Client returns the Client the Database was created from.
func (db *Database) Client() *Client ***REMOVED***
	return db.client
***REMOVED***

// Name returns the name of the database.
func (db *Database) Name() string ***REMOVED***
	return db.name
***REMOVED***

// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection ***REMOVED***
	return newCollection(db, name, opts...)
***REMOVED***

// Aggregate executes an aggregate command the database. This requires MongoDB version >= 3.6 and driver version >=
// 1.1.0.
//
// The pipeline parameter must be a slice of documents, each representing an aggregation stage. The pipeline
// cannot be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
// mongo.Pipeline type can be used. See
// https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/#db-aggregate-stages for a list of valid
// stages in database-level aggregations.
//
// The opts parameter can be used to specify options for this operation (see the options.AggregateOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/aggregate/.
func (db *Database) Aggregate(ctx context.Context, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.AggregateOptions) (*Cursor, error) ***REMOVED***
	a := aggregateParams***REMOVED***
		ctx:            ctx,
		pipeline:       pipeline,
		client:         db.client,
		registry:       db.registry,
		readConcern:    db.readConcern,
		writeConcern:   db.writeConcern,
		retryRead:      db.client.retryReads,
		db:             db.name,
		readSelector:   db.readSelector,
		writeSelector:  db.writeSelector,
		readPreference: db.readPreference,
		opts:           opts,
	***REMOVED***
	return aggregate(a)
***REMOVED***

func (db *Database) processRunCommand(ctx context.Context, cmd interface***REMOVED******REMOVED***,
	cursorCommand bool, opts ...*options.RunCmdOptions) (*operation.Command, *session.Client, error) ***REMOVED***
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, sess, err
		***REMOVED***
	***REMOVED***

	err := db.client.validSession(sess)
	if err != nil ***REMOVED***
		return nil, sess, err
	***REMOVED***

	ro := options.MergeRunCmdOptions(append(defaultRunCmdOpts, opts...)...)
	if sess != nil && sess.TransactionRunning() && ro.ReadPreference != nil && ro.ReadPreference.Mode() != readpref.PrimaryMode ***REMOVED***
		return nil, sess, errors.New("read preference in a transaction must be primary")
	***REMOVED***

	runCmdDoc, err := transformBsoncoreDocument(db.registry, cmd, false, "cmd")
	if err != nil ***REMOVED***
		return nil, sess, err
	***REMOVED***
	readSelect := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(ro.ReadPreference),
		description.LatencySelector(db.client.localThreshold),
	***REMOVED***)
	if sess != nil && sess.PinnedServer != nil ***REMOVED***
		readSelect = makePinnedSelector(sess, readSelect)
	***REMOVED***

	var op *operation.Command
	switch cursorCommand ***REMOVED***
	case true:
		cursorOpts := db.client.createBaseCursorOptions()
		op = operation.NewCursorCommand(runCmdDoc, cursorOpts)
	default:
		op = operation.NewCommand(runCmdDoc)
	***REMOVED***
	return op.Session(sess).CommandMonitor(db.client.monitor).
		ServerSelector(readSelect).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).ReadConcern(db.readConcern).
		Crypt(db.client.cryptFLE).ReadPreference(ro.ReadPreference).ServerAPI(db.client.serverAPI).
		Timeout(db.client.timeout), sess, nil
***REMOVED***

// RunCommand executes the given command against the database. This function does not obey the Database's read
// preference. To specify a read preference, the RunCmdOptions.ReadPreference option must be used.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
//
// The behavior of RunCommand is undefined if the command document contains any of the following:
// - A session ID or any transaction-specific fields
// - API versioning options when an API version is already declared on the Client
// - maxTimeMS when Timeout is set on the Client
func (db *Database) RunCommand(ctx context.Context, runCommand interface***REMOVED******REMOVED***, opts ...*options.RunCmdOptions) *SingleResult ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	op, sess, err := db.processRunCommand(ctx, runCommand, false, opts...)
	defer closeImplicitSession(sess)
	if err != nil ***REMOVED***
		return &SingleResult***REMOVED***err: err***REMOVED***
	***REMOVED***

	err = op.Execute(ctx)
	// RunCommand can be used to run a write, thus execute may return a write error
	_, convErr := processWriteError(err)
	return &SingleResult***REMOVED***
		err: convErr,
		rdr: bson.Raw(op.Result()),
		reg: db.registry,
	***REMOVED***
***REMOVED***

// RunCommandCursor executes the given command against the database and parses the response as a cursor. If the command
// being executed does not return a cursor (e.g. insert), the command will be executed on the server and an error will
// be returned because the server response cannot be parsed as a cursor. This function does not obey the Database's read
// preference. To specify a read preference, the RunCmdOptions.ReadPreference option must be used.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
//
// The behavior of RunCommandCursor is undefined if the command document contains any of the following:
// - A session ID or any transaction-specific fields
// - API versioning options when an API version is already declared on the Client
// - maxTimeMS when Timeout is set on the Client
func (db *Database) RunCommandCursor(ctx context.Context, runCommand interface***REMOVED******REMOVED***, opts ...*options.RunCmdOptions) (*Cursor, error) ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	op, sess, err := db.processRunCommand(ctx, runCommand, true, opts...)
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	***REMOVED***

	if err = op.Execute(ctx); err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	***REMOVED***

	bc, err := op.ResultCursor()
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	***REMOVED***
	cursor, err := newCursorWithSession(bc, db.registry, sess)
	return cursor, replaceErrors(err)
***REMOVED***

// Drop drops the database on the server. This method ignores "namespace not found" errors so it is safe to drop
// a database that does not exist on the server.
func (db *Database) Drop(ctx context.Context) error ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err := db.client.validSession(sess)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	wc := db.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, db.writeSelector)

	op := operation.NewDropDatabase().
		Session(sess).WriteConcern(wc).CommandMonitor(db.client.monitor).
		ServerSelector(selector).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.cryptFLE).
		ServerAPI(db.client.serverAPI)

	err = op.Execute(ctx)

	driverErr, ok := err.(driver.Error)
	if err != nil && (!ok || !driverErr.NamespaceNotFound()) ***REMOVED***
		return replaceErrors(err)
	***REMOVED***
	return nil
***REMOVED***

// ListCollectionSpecifications executes a listCollections command and returns a slice of CollectionSpecification
// instances representing the collections in the database.
//
// The filter parameter must be a document containing query operators and can be used to select which collections
// are included in the result. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to include all
// collections.
//
// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listCollections/.
//
// BUG(benjirewis): ListCollectionSpecifications prevents listing more than 100 collections per database when running
// against MongoDB version 2.6.
func (db *Database) ListCollectionSpecifications(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.ListCollectionsOptions) ([]*CollectionSpecification, error) ***REMOVED***

	cursor, err := db.ListCollections(ctx, filter, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var specs []*CollectionSpecification
	err = cursor.All(ctx, &specs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, spec := range specs ***REMOVED***
		// Pre-4.4 servers report a namespace in their responses, so we only set Namespace manually if it was not in
		// the response.
		if spec.IDIndex != nil && spec.IDIndex.Namespace == "" ***REMOVED***
			spec.IDIndex.Namespace = db.name + "." + spec.Name
		***REMOVED***
	***REMOVED***
	return specs, nil
***REMOVED***

// ListCollections executes a listCollections command and returns a cursor over the collections in the database.
//
// The filter parameter must be a document containing query operators and can be used to select which collections
// are included in the result. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to include all
// collections.
//
// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listCollections/.
//
// BUG(benjirewis): ListCollections prevents listing more than 100 collections per database when running against
// MongoDB version 2.6.
func (db *Database) ListCollections(ctx context.Context, filter interface***REMOVED******REMOVED***, opts ...*options.ListCollectionsOptions) (*Cursor, error) ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	filterDoc, err := transformBsoncoreDocument(db.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	err = db.client.validSession(sess)
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, err
	***REMOVED***

	selector := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(db.client.localThreshold),
	***REMOVED***)
	selector = makeReadPrefSelector(sess, selector, db.client.localThreshold)

	lco := options.MergeListCollectionsOptions(opts...)
	op := operation.NewListCollections(filterDoc).
		Session(sess).ReadPreference(db.readPreference).CommandMonitor(db.client.monitor).
		ServerSelector(selector).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.cryptFLE).
		ServerAPI(db.client.serverAPI).Timeout(db.client.timeout)

	cursorOpts := db.client.createBaseCursorOptions()
	if lco.NameOnly != nil ***REMOVED***
		op = op.NameOnly(*lco.NameOnly)
	***REMOVED***
	if lco.BatchSize != nil ***REMOVED***
		cursorOpts.BatchSize = *lco.BatchSize
		op = op.BatchSize(*lco.BatchSize)
	***REMOVED***
	if lco.AuthorizedCollections != nil ***REMOVED***
		op = op.AuthorizedCollections(*lco.AuthorizedCollections)
	***REMOVED***

	retry := driver.RetryNone
	if db.client.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	***REMOVED***

	bc, err := op.Result(cursorOpts)
	if err != nil ***REMOVED***
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	***REMOVED***
	cursor, err := newCursorWithSession(bc, db.registry, sess)
	return cursor, replaceErrors(err)
***REMOVED***

// ListCollectionNames executes a listCollections command and returns a slice containing the names of the collections
// in the database. This method requires driver version >= 1.1.0.
//
// The filter parameter must be a document containing query operators and can be used to select which collections
// are included in the result. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to include all
// collections.
//
// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listCollections/.
//
// BUG(benjirewis): ListCollectionNames prevents listing more than 100 collections per database when running against
// MongoDB version 2.6.
func (db *Database) ListCollectionNames(ctx context.Context, filter interface***REMOVED******REMOVED***, opts ...*options.ListCollectionsOptions) ([]string, error) ***REMOVED***
	opts = append(opts, options.ListCollections().SetNameOnly(true))

	res, err := db.ListCollections(ctx, filter, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer res.Close(ctx)

	names := make([]string, 0)
	for res.Next(ctx) ***REMOVED***
		elem, err := res.Current.LookupErr("name")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if elem.Type != bson.TypeString ***REMOVED***
			return nil, fmt.Errorf("incorrect type for 'name'. got %v. want %v", elem.Type, bson.TypeString)
		***REMOVED***

		elemName := elem.StringValue()
		names = append(names, elemName)
	***REMOVED***

	res.Close(ctx)
	return names, nil
***REMOVED***

// ReadConcern returns the read concern used to configure the Database object.
func (db *Database) ReadConcern() *readconcern.ReadConcern ***REMOVED***
	return db.readConcern
***REMOVED***

// ReadPreference returns the read preference used to configure the Database object.
func (db *Database) ReadPreference() *readpref.ReadPref ***REMOVED***
	return db.readPreference
***REMOVED***

// WriteConcern returns the write concern used to configure the Database object.
func (db *Database) WriteConcern() *writeconcern.WriteConcern ***REMOVED***
	return db.writeConcern
***REMOVED***

// Watch returns a change stream for all changes to the corresponding database. See
// https://www.mongodb.com/docs/manual/changeStreams/ for more information about change streams.
//
// The Database must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be a slice of documents, each representing a pipeline stage. The pipeline cannot be
// nil but can be empty. The stage documents must all be non-nil. See https://www.mongodb.com/docs/manual/changeStreams/ for
// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
// mongo.Pipeline***REMOVED******REMOVED*** type can be used.
//
// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
// documentation).
func (db *Database) Watch(ctx context.Context, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) ***REMOVED***

	csConfig := changeStreamConfig***REMOVED***
		readConcern:    db.readConcern,
		readPreference: db.readPreference,
		client:         db.client,
		registry:       db.registry,
		streamType:     DatabaseStream,
		databaseName:   db.Name(),
		crypt:          db.client.cryptFLE,
	***REMOVED***
	return newChangeStream(ctx, csConfig, pipeline, opts...)
***REMOVED***

// CreateCollection executes a create command to explicitly create a new collection with the specified name on the
// server. If the collection being created already exists, this method will return a mongo.CommandError. This method
// requires driver version 1.4.0 or higher.
//
// The opts parameter can be used to specify options for the operation (see the options.CreateCollectionOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/create/.
func (db *Database) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error ***REMOVED***
	cco := options.MergeCreateCollectionOptions(opts...)
	// Follow Client-Side Encryption specification to check for encryptedFields.
	// Check for encryptedFields from create options.
	ef := cco.EncryptedFields
	// Check for encryptedFields from the client EncryptedFieldsMap.
	if ef == nil ***REMOVED***
		ef = db.getEncryptedFieldsFromMap(name)
	***REMOVED***
	if ef != nil ***REMOVED***
		return db.createCollectionWithEncryptedFields(ctx, name, ef, opts...)
	***REMOVED***

	return db.createCollection(ctx, name, opts...)
***REMOVED***

// getEncryptedFieldsFromServer tries to get an "encryptedFields" document associated with collectionName by running the "listCollections" command.
// Returns nil and no error if the listCollections command succeeds, but "encryptedFields" is not present.
func (db *Database) getEncryptedFieldsFromServer(ctx context.Context, collectionName string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	// Check if collection has an EncryptedFields configured server-side.
	collSpecs, err := db.ListCollectionSpecifications(ctx, bson.D***REMOVED******REMOVED***"name", collectionName***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(collSpecs) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	if len(collSpecs) > 1 ***REMOVED***
		return nil, fmt.Errorf("expected 1 or 0 results from listCollections, got %v", len(collSpecs))
	***REMOVED***
	collSpec := collSpecs[0]
	rawValue, err := collSpec.Options.LookupErr("encryptedFields")
	if err == bsoncore.ErrElementNotFound ***REMOVED***
		return nil, nil
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	encryptedFields, ok := rawValue.DocumentOK()
	if !ok ***REMOVED***
		return nil, fmt.Errorf("expected encryptedFields of %v to be document, got %v", collectionName, rawValue.Type)
	***REMOVED***

	return encryptedFields, nil
***REMOVED***

// getEncryptedFieldsFromServer tries to get an "encryptedFields" document associated with collectionName by checking the client EncryptedFieldsMap.
// Returns nil and no error if an EncryptedFieldsMap is not configured, or does not contain an entry for collectionName.
func (db *Database) getEncryptedFieldsFromMap(collectionName string) interface***REMOVED******REMOVED*** ***REMOVED***
	// Check the EncryptedFieldsMap
	efMap := db.client.encryptedFieldsMap
	if efMap == nil ***REMOVED***
		return nil
	***REMOVED***

	namespace := db.name + "." + collectionName

	ef, ok := efMap[namespace]
	if ok ***REMOVED***
		return ef
	***REMOVED***
	return nil
***REMOVED***

// createCollectionWithEncryptedFields creates a collection with an EncryptedFields.
func (db *Database) createCollectionWithEncryptedFields(ctx context.Context, name string, ef interface***REMOVED******REMOVED***, opts ...*options.CreateCollectionOptions) error ***REMOVED***
	efBSON, err := transformBsoncoreDocument(db.registry, ef, true /* mapAllowed */, "encryptedFields")
	if err != nil ***REMOVED***
		return fmt.Errorf("error transforming document: %v", err)
	***REMOVED***

	// Create the three encryption-related, associated collections: `escCollection`, `eccCollection` and `ecocCollection`.

	stateCollectionOpts := options.CreateCollection().
		SetClusteredIndex(bson.D***REMOVED******REMOVED***"key", bson.D***REMOVED******REMOVED***"_id", 1***REMOVED******REMOVED******REMOVED***, ***REMOVED***"unique", true***REMOVED******REMOVED***)
	// Create ESCCollection.
	escCollection, err := internal.GetEncryptedStateCollectionName(efBSON, name, internal.EncryptedStateCollection)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := db.createCollection(ctx, escCollection, stateCollectionOpts); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create ECCCollection.
	eccCollection, err := internal.GetEncryptedStateCollectionName(efBSON, name, internal.EncryptedCacheCollection)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := db.createCollection(ctx, eccCollection, stateCollectionOpts); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create ECOCCollection.
	ecocCollection, err := internal.GetEncryptedStateCollectionName(efBSON, name, internal.EncryptedCompactionCollection)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := db.createCollection(ctx, ecocCollection, stateCollectionOpts); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create a data collection with the 'encryptedFields' option.
	op, err := db.createCollectionOperation(name, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	op.EncryptedFields(efBSON)
	if err := db.executeCreateOperation(ctx, op); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Create an index on the __safeContent__ field in the collection @collectionName.
	if _, err := db.Collection(name).Indexes().CreateOne(ctx, IndexModel***REMOVED***Keys: bson.D***REMOVED******REMOVED***"__safeContent__", 1***REMOVED******REMOVED******REMOVED***); err != nil ***REMOVED***
		return fmt.Errorf("error creating safeContent index: %v", err)
	***REMOVED***

	return nil
***REMOVED***

// createCollection creates a collection without EncryptedFields.
func (db *Database) createCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error ***REMOVED***
	op, err := db.createCollectionOperation(name, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return db.executeCreateOperation(ctx, op)
***REMOVED***

func (db *Database) createCollectionOperation(name string, opts ...*options.CreateCollectionOptions) (*operation.Create, error) ***REMOVED***
	cco := options.MergeCreateCollectionOptions(opts...)
	op := operation.NewCreate(name).ServerAPI(db.client.serverAPI)

	if cco.Capped != nil ***REMOVED***
		op.Capped(*cco.Capped)
	***REMOVED***
	if cco.Collation != nil ***REMOVED***
		op.Collation(bsoncore.Document(cco.Collation.ToDocument()))
	***REMOVED***
	if cco.ChangeStreamPreAndPostImages != nil ***REMOVED***
		csppi, err := transformBsoncoreDocument(db.registry, cco.ChangeStreamPreAndPostImages, true, "changeStreamPreAndPostImages")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.ChangeStreamPreAndPostImages(csppi)
	***REMOVED***
	if cco.DefaultIndexOptions != nil ***REMOVED***
		idx, doc := bsoncore.AppendDocumentStart(nil)
		if cco.DefaultIndexOptions.StorageEngine != nil ***REMOVED***
			storageEngine, err := transformBsoncoreDocument(db.registry, cco.DefaultIndexOptions.StorageEngine, true, "storageEngine")
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			doc = bsoncore.AppendDocumentElement(doc, "storageEngine", storageEngine)
		***REMOVED***
		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		op.IndexOptionDefaults(doc)
	***REMOVED***
	if cco.MaxDocuments != nil ***REMOVED***
		op.Max(*cco.MaxDocuments)
	***REMOVED***
	if cco.SizeInBytes != nil ***REMOVED***
		op.Size(*cco.SizeInBytes)
	***REMOVED***
	if cco.StorageEngine != nil ***REMOVED***
		storageEngine, err := transformBsoncoreDocument(db.registry, cco.StorageEngine, true, "storageEngine")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.StorageEngine(storageEngine)
	***REMOVED***
	if cco.ValidationAction != nil ***REMOVED***
		op.ValidationAction(*cco.ValidationAction)
	***REMOVED***
	if cco.ValidationLevel != nil ***REMOVED***
		op.ValidationLevel(*cco.ValidationLevel)
	***REMOVED***
	if cco.Validator != nil ***REMOVED***
		validator, err := transformBsoncoreDocument(db.registry, cco.Validator, true, "validator")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.Validator(validator)
	***REMOVED***
	if cco.ExpireAfterSeconds != nil ***REMOVED***
		op.ExpireAfterSeconds(*cco.ExpireAfterSeconds)
	***REMOVED***
	if cco.TimeSeriesOptions != nil ***REMOVED***
		idx, doc := bsoncore.AppendDocumentStart(nil)
		doc = bsoncore.AppendStringElement(doc, "timeField", cco.TimeSeriesOptions.TimeField)

		if cco.TimeSeriesOptions.MetaField != nil ***REMOVED***
			doc = bsoncore.AppendStringElement(doc, "metaField", *cco.TimeSeriesOptions.MetaField)
		***REMOVED***
		if cco.TimeSeriesOptions.Granularity != nil ***REMOVED***
			doc = bsoncore.AppendStringElement(doc, "granularity", *cco.TimeSeriesOptions.Granularity)
		***REMOVED***

		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		op.TimeSeries(doc)
	***REMOVED***
	if cco.ClusteredIndex != nil ***REMOVED***
		clusteredIndex, err := transformBsoncoreDocument(db.registry, cco.ClusteredIndex, true, "clusteredIndex")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		op.ClusteredIndex(clusteredIndex)
	***REMOVED***

	return op, nil
***REMOVED***

// CreateView executes a create command to explicitly create a view on the server. See
// https://www.mongodb.com/docs/manual/core/views/ for more information about views. This method requires driver version >=
// 1.4.0 and MongoDB version >= 3.4.
//
// The viewName parameter specifies the name of the view to create.
//
// The viewOn parameter specifies the name of the collection or view on which this view will be created
//
// The pipeline parameter specifies an aggregation pipeline that will be exececuted against the source collection or
// view to create this view.
//
// The opts parameter can be used to specify options for the operation (see the options.CreateViewOptions
// documentation).
func (db *Database) CreateView(ctx context.Context, viewName, viewOn string, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.CreateViewOptions) error ***REMOVED***

	pipelineArray, _, err := transformAggregatePipeline(db.registry, pipeline)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	op := operation.NewCreate(viewName).
		ViewOn(viewOn).
		Pipeline(pipelineArray).
		ServerAPI(db.client.serverAPI)
	cvo := options.MergeCreateViewOptions(opts...)
	if cvo.Collation != nil ***REMOVED***
		op.Collation(bsoncore.Document(cvo.Collation.ToDocument()))
	***REMOVED***

	return db.executeCreateOperation(ctx, op)
***REMOVED***

func (db *Database) executeCreateOperation(ctx context.Context, op *operation.Create) error ***REMOVED***
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil ***REMOVED***
		var err error
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err := db.client.validSession(sess)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	wc := db.writeConcern
	if sess.TransactionRunning() ***REMOVED***
		wc = nil
	***REMOVED***
	if !writeconcern.AckWrite(wc) ***REMOVED***
		sess = nil
	***REMOVED***

	selector := makePinnedSelector(sess, db.writeSelector)
	op = op.Session(sess).
		WriteConcern(wc).
		CommandMonitor(db.client.monitor).
		ServerSelector(selector).
		ClusterClock(db.client.clock).
		Database(db.name).
		Deployment(db.client.deployment).
		Crypt(db.client.cryptFLE)

	return replaceErrors(op.Execute(ctx))
***REMOVED***
