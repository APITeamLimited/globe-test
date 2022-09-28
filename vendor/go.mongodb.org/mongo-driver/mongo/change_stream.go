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
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var (
	// ErrMissingResumeToken indicates that a change stream notification from the server did not contain a resume token.
	ErrMissingResumeToken = errors.New("cannot provide resume functionality when the resume token is missing")
	// ErrNilCursor indicates that the underlying cursor for the change stream is nil.
	ErrNilCursor = errors.New("cursor is nil")

	minResumableLabelWireVersion int32 = 9 // Wire version at which the server includes the resumable error label
	networkErrorLabel                  = "NetworkError"
	resumableErrorLabel                = "ResumableChangeStreamError"
	errorCursorNotFound          int32 = 43 // CursorNotFound error code

	// Allowlist of error codes that are considered resumable.
	resumableChangeStreamErrors = map[int32]struct***REMOVED******REMOVED******REMOVED***
		6:     ***REMOVED******REMOVED***, // HostUnreachable
		7:     ***REMOVED******REMOVED***, // HostNotFound
		89:    ***REMOVED******REMOVED***, // NetworkTimeout
		91:    ***REMOVED******REMOVED***, // ShutdownInProgress
		189:   ***REMOVED******REMOVED***, // PrimarySteppedDown
		262:   ***REMOVED******REMOVED***, // ExceededTimeLimit
		9001:  ***REMOVED******REMOVED***, // SocketException
		10107: ***REMOVED******REMOVED***, // NotPrimary
		11600: ***REMOVED******REMOVED***, // InterruptedAtShutdown
		11602: ***REMOVED******REMOVED***, // InterruptedDueToReplStateChange
		13435: ***REMOVED******REMOVED***, // NotPrimaryNoSecondaryOK
		13436: ***REMOVED******REMOVED***, // NotPrimaryOrSecondary
		63:    ***REMOVED******REMOVED***, // StaleShardVersion
		150:   ***REMOVED******REMOVED***, // StaleEpoch
		13388: ***REMOVED******REMOVED***, // StaleConfig
		234:   ***REMOVED******REMOVED***, // RetryChangeStream
		133:   ***REMOVED******REMOVED***, // FailedToSatisfyReadPreference
	***REMOVED***
)

// ChangeStream is used to iterate over a stream of events. Each event can be decoded into a Go type via the Decode
// method or accessed as raw BSON via the Current field. This type is not goroutine safe and must not be used
// concurrently by multiple goroutines. For more information about change streams, see
// https://www.mongodb.com/docs/manual/changeStreams/.
type ChangeStream struct ***REMOVED***
	// Current is the BSON bytes of the current event. This property is only valid until the next call to Next or
	// TryNext. If continued access is required, a copy must be made.
	Current bson.Raw

	aggregate       *operation.Aggregate
	pipelineSlice   []bsoncore.Document
	pipelineOptions map[string]bsoncore.Value
	cursor          changeStreamCursor
	cursorOptions   driver.CursorOptions
	batch           []bsoncore.Document
	resumeToken     bson.Raw
	err             error
	sess            *session.Client
	client          *Client
	registry        *bsoncodec.Registry
	streamType      StreamType
	options         *options.ChangeStreamOptions
	selector        description.ServerSelector
	operationTime   *primitive.Timestamp
	wireVersion     *description.VersionRange
***REMOVED***

type changeStreamConfig struct ***REMOVED***
	readConcern    *readconcern.ReadConcern
	readPreference *readpref.ReadPref
	client         *Client
	registry       *bsoncodec.Registry
	streamType     StreamType
	collectionName string
	databaseName   string
	crypt          driver.Crypt
***REMOVED***

func newChangeStream(ctx context.Context, config changeStreamConfig, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	cs := &ChangeStream***REMOVED***
		client:     config.client,
		registry:   config.registry,
		streamType: config.streamType,
		options:    options.MergeChangeStreamOptions(opts...),
		selector: description.CompositeSelector([]description.ServerSelector***REMOVED***
			description.ReadPrefSelector(config.readPreference),
			description.LatencySelector(config.client.localThreshold),
		***REMOVED***),
		cursorOptions: config.client.createBaseCursorOptions(),
	***REMOVED***

	cs.sess = sessionFromContext(ctx)
	if cs.sess == nil && cs.client.sessionPool != nil ***REMOVED***
		cs.sess, cs.err = session.NewClientSession(cs.client.sessionPool, cs.client.id, session.Implicit)
		if cs.err != nil ***REMOVED***
			return nil, cs.Err()
		***REMOVED***
	***REMOVED***
	if cs.err = cs.client.validSession(cs.sess); cs.err != nil ***REMOVED***
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	***REMOVED***

	cs.aggregate = operation.NewAggregate(nil).
		ReadPreference(config.readPreference).ReadConcern(config.readConcern).
		Deployment(cs.client.deployment).ClusterClock(cs.client.clock).
		CommandMonitor(cs.client.monitor).Session(cs.sess).ServerSelector(cs.selector).Retry(driver.RetryNone).
		ServerAPI(cs.client.serverAPI).Crypt(config.crypt).Timeout(cs.client.timeout)

	if cs.options.Collation != nil ***REMOVED***
		cs.aggregate.Collation(bsoncore.Document(cs.options.Collation.ToDocument()))
	***REMOVED***
	if comment := cs.options.Comment; comment != nil ***REMOVED***
		cs.aggregate.Comment(*comment)

		commentVal, err := transformValue(cs.registry, comment, true, "comment")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cs.cursorOptions.Comment = commentVal
	***REMOVED***
	if cs.options.BatchSize != nil ***REMOVED***
		cs.aggregate.BatchSize(*cs.options.BatchSize)
		cs.cursorOptions.BatchSize = *cs.options.BatchSize
	***REMOVED***
	if cs.options.MaxAwaitTime != nil ***REMOVED***
		cs.cursorOptions.MaxTimeMS = int64(*cs.options.MaxAwaitTime / time.Millisecond)
	***REMOVED***
	if cs.options.Custom != nil ***REMOVED***
		// Marshal all custom options before passing to the initial aggregate. Return
		// any errors from Marshaling.
		customOptions := make(map[string]bsoncore.Value)
		for optionName, optionValue := range cs.options.Custom ***REMOVED***
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(cs.registry, optionValue)
			if err != nil ***REMOVED***
				cs.err = err
				closeImplicitSession(cs.sess)
				return nil, cs.Err()
			***REMOVED***
			optionValueBSON := bsoncore.Value***REMOVED***Type: bsonType, Data: bsonData***REMOVED***
			customOptions[optionName] = optionValueBSON
		***REMOVED***
		cs.aggregate.CustomOptions(customOptions)
	***REMOVED***
	if cs.options.CustomPipeline != nil ***REMOVED***
		// Marshal all custom pipeline options before building pipeline slice. Return
		// any errors from Marshaling.
		cs.pipelineOptions = make(map[string]bsoncore.Value)
		for optionName, optionValue := range cs.options.CustomPipeline ***REMOVED***
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(cs.registry, optionValue)
			if err != nil ***REMOVED***
				cs.err = err
				closeImplicitSession(cs.sess)
				return nil, cs.Err()
			***REMOVED***
			optionValueBSON := bsoncore.Value***REMOVED***Type: bsonType, Data: bsonData***REMOVED***
			cs.pipelineOptions[optionName] = optionValueBSON
		***REMOVED***
	***REMOVED***

	switch cs.streamType ***REMOVED***
	case ClientStream:
		cs.aggregate.Database("admin")
	case DatabaseStream:
		cs.aggregate.Database(config.databaseName)
	case CollectionStream:
		cs.aggregate.Collection(config.collectionName).Database(config.databaseName)
	default:
		closeImplicitSession(cs.sess)
		return nil, fmt.Errorf("must supply a valid StreamType in config, instead of %v", cs.streamType)
	***REMOVED***

	// When starting a change stream, cache startAfter as the first resume token if it is set. If not, cache
	// resumeAfter. If neither is set, do not cache a resume token.
	resumeToken := cs.options.StartAfter
	if resumeToken == nil ***REMOVED***
		resumeToken = cs.options.ResumeAfter
	***REMOVED***
	var marshaledToken bson.Raw
	if resumeToken != nil ***REMOVED***
		if marshaledToken, cs.err = bson.Marshal(resumeToken); cs.err != nil ***REMOVED***
			closeImplicitSession(cs.sess)
			return nil, cs.Err()
		***REMOVED***
	***REMOVED***
	cs.resumeToken = marshaledToken

	if cs.err = cs.buildPipelineSlice(pipeline); cs.err != nil ***REMOVED***
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	***REMOVED***
	var pipelineArr bsoncore.Document
	pipelineArr, cs.err = cs.pipelineToBSON()
	cs.aggregate.Pipeline(pipelineArr)

	if cs.err = cs.executeOperation(ctx, false); cs.err != nil ***REMOVED***
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	***REMOVED***

	return cs, cs.Err()
***REMOVED***

func (cs *ChangeStream) createOperationDeployment(server driver.Server, connection driver.Connection) driver.Deployment ***REMOVED***
	return &changeStreamDeployment***REMOVED***
		topologyKind: cs.client.deployment.Kind(),
		server:       server,
		conn:         connection,
	***REMOVED***
***REMOVED***

func (cs *ChangeStream) executeOperation(ctx context.Context, resuming bool) error ***REMOVED***
	var server driver.Server
	var conn driver.Connection
	var err error

	if server, cs.err = cs.client.deployment.SelectServer(ctx, cs.selector); cs.err != nil ***REMOVED***
		return cs.Err()
	***REMOVED***
	if conn, cs.err = server.Connection(ctx); cs.err != nil ***REMOVED***
		return cs.Err()
	***REMOVED***
	defer conn.Close()
	cs.wireVersion = conn.Description().WireVersion

	cs.aggregate.Deployment(cs.createOperationDeployment(server, conn))

	if resuming ***REMOVED***
		cs.replaceOptions(cs.wireVersion)

		csOptDoc, err := cs.createPipelineOptionsDoc()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		pipIdx, pipDoc := bsoncore.AppendDocumentStart(nil)
		pipDoc = bsoncore.AppendDocumentElement(pipDoc, "$changeStream", csOptDoc)
		if pipDoc, cs.err = bsoncore.AppendDocumentEnd(pipDoc, pipIdx); cs.err != nil ***REMOVED***
			return cs.Err()
		***REMOVED***
		cs.pipelineSlice[0] = pipDoc

		var plArr bsoncore.Document
		if plArr, cs.err = cs.pipelineToBSON(); cs.err != nil ***REMOVED***
			return cs.Err()
		***REMOVED***
		cs.aggregate.Pipeline(plArr)
	***REMOVED***

	// If no deadline is set on the passed-in context, cs.client.timeout is set, and context is not already
	// a Timeout context, honor cs.client.timeout in new Timeout context for change stream operation execution
	// and potential retry.
	if _, deadlineSet := ctx.Deadline(); !deadlineSet && cs.client.timeout != nil && !internal.IsTimeoutContext(ctx) ***REMOVED***
		newCtx, cancelFunc := internal.MakeTimeoutContext(ctx, *cs.client.timeout)
		// Redefine ctx to be the new timeout-derived context.
		ctx = newCtx
		// Cancel the timeout-derived context at the end of executeOperation to avoid a context leak.
		defer cancelFunc()
	***REMOVED***
	if original := cs.aggregate.Execute(ctx); original != nil ***REMOVED***
		retryableRead := cs.client.retryReads && cs.wireVersion != nil && cs.wireVersion.Max >= 6
		if !retryableRead ***REMOVED***
			cs.err = replaceErrors(original)
			return cs.err
		***REMOVED***

		cs.err = original
		switch tt := original.(type) ***REMOVED***
		case driver.Error:
			if !tt.RetryableRead() ***REMOVED***
				break
			***REMOVED***

			server, err = cs.client.deployment.SelectServer(ctx, cs.selector)
			if err != nil ***REMOVED***
				break
			***REMOVED***

			conn.Close()
			conn, err = server.Connection(ctx)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			defer conn.Close()
			cs.wireVersion = conn.Description().WireVersion

			if cs.wireVersion == nil || cs.wireVersion.Max < 6 ***REMOVED***
				break
			***REMOVED***

			cs.aggregate.Deployment(cs.createOperationDeployment(server, conn))
			cs.err = cs.aggregate.Execute(ctx)
		***REMOVED***

		if cs.err != nil ***REMOVED***
			cs.err = replaceErrors(cs.err)
			return cs.Err()
		***REMOVED***

	***REMOVED***
	cs.err = nil

	cr := cs.aggregate.ResultCursorResponse()
	cr.Server = server

	cs.cursor, cs.err = driver.NewBatchCursor(cr, cs.sess, cs.client.clock, cs.cursorOptions)
	if cs.err = replaceErrors(cs.err); cs.err != nil ***REMOVED***
		return cs.Err()
	***REMOVED***

	cs.updatePbrtFromCommand()
	if cs.options.StartAtOperationTime == nil && cs.options.ResumeAfter == nil &&
		cs.options.StartAfter == nil && cs.wireVersion.Max >= 7 &&
		cs.emptyBatch() && cs.resumeToken == nil ***REMOVED***
		cs.operationTime = cs.sess.OperationTime
	***REMOVED***

	return cs.Err()
***REMOVED***

// Updates the post batch resume token after a successful aggregate or getMore operation.
func (cs *ChangeStream) updatePbrtFromCommand() ***REMOVED***
	// Only cache the pbrt if an empty batch was returned and a pbrt was included
	if pbrt := cs.cursor.PostBatchResumeToken(); cs.emptyBatch() && pbrt != nil ***REMOVED***
		cs.resumeToken = bson.Raw(pbrt)
	***REMOVED***
***REMOVED***

func (cs *ChangeStream) storeResumeToken() error ***REMOVED***
	// If cs.Current is the last document in the batch and a pbrt is included, cache the pbrt
	// Otherwise, cache the _id of the document
	var tokenDoc bson.Raw
	if len(cs.batch) == 0 ***REMOVED***
		if pbrt := cs.cursor.PostBatchResumeToken(); pbrt != nil ***REMOVED***
			tokenDoc = bson.Raw(pbrt)
		***REMOVED***
	***REMOVED***

	if tokenDoc == nil ***REMOVED***
		var ok bool
		tokenDoc, ok = cs.Current.Lookup("_id").DocumentOK()
		if !ok ***REMOVED***
			_ = cs.Close(context.Background())
			return ErrMissingResumeToken
		***REMOVED***
	***REMOVED***

	cs.resumeToken = tokenDoc
	return nil
***REMOVED***

func (cs *ChangeStream) buildPipelineSlice(pipeline interface***REMOVED******REMOVED***) error ***REMOVED***
	val := reflect.ValueOf(pipeline)
	if !val.IsValid() || !(val.Kind() == reflect.Slice) ***REMOVED***
		cs.err = errors.New("can only transform slices and arrays into aggregation pipelines, but got invalid")
		return cs.err
	***REMOVED***

	cs.pipelineSlice = make([]bsoncore.Document, 0, val.Len()+1)

	csIdx, csDoc := bsoncore.AppendDocumentStart(nil)

	csDocTemp, err := cs.createPipelineOptionsDoc()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	csDoc = bsoncore.AppendDocumentElement(csDoc, "$changeStream", csDocTemp)
	csDoc, cs.err = bsoncore.AppendDocumentEnd(csDoc, csIdx)
	if cs.err != nil ***REMOVED***
		return cs.err
	***REMOVED***
	cs.pipelineSlice = append(cs.pipelineSlice, csDoc)

	for i := 0; i < val.Len(); i++ ***REMOVED***
		var elem []byte
		elem, cs.err = transformBsoncoreDocument(cs.registry, val.Index(i).Interface(), true, fmt.Sprintf("pipeline stage :%v", i))
		if cs.err != nil ***REMOVED***
			return cs.err
		***REMOVED***

		cs.pipelineSlice = append(cs.pipelineSlice, elem)
	***REMOVED***

	return cs.err
***REMOVED***

func (cs *ChangeStream) createPipelineOptionsDoc() (bsoncore.Document, error) ***REMOVED***
	plDocIdx, plDoc := bsoncore.AppendDocumentStart(nil)

	if cs.streamType == ClientStream ***REMOVED***
		plDoc = bsoncore.AppendBooleanElement(plDoc, "allChangesForCluster", true)
	***REMOVED***

	if cs.options.FullDocument != nil ***REMOVED***
		// Only append a default "fullDocument" field if wire version is less than 6 (3.6). Otherwise,
		// the server will assume users want the default behavior, and "fullDocument" does not need to be
		// specified.
		if *cs.options.FullDocument != options.Default || (cs.wireVersion != nil && cs.wireVersion.Max < 6) ***REMOVED***
			plDoc = bsoncore.AppendStringElement(plDoc, "fullDocument", string(*cs.options.FullDocument))
		***REMOVED***
	***REMOVED***

	if cs.options.FullDocumentBeforeChange != nil ***REMOVED***
		plDoc = bsoncore.AppendStringElement(plDoc, "fullDocumentBeforeChange", string(*cs.options.FullDocumentBeforeChange))
	***REMOVED***

	if cs.options.ResumeAfter != nil ***REMOVED***
		var raDoc bsoncore.Document
		raDoc, cs.err = transformBsoncoreDocument(cs.registry, cs.options.ResumeAfter, true, "resumeAfter")
		if cs.err != nil ***REMOVED***
			return nil, cs.err
		***REMOVED***

		plDoc = bsoncore.AppendDocumentElement(plDoc, "resumeAfter", raDoc)
	***REMOVED***

	if cs.options.ShowExpandedEvents != nil ***REMOVED***
		plDoc = bsoncore.AppendBooleanElement(plDoc, "showExpandedEvents", *cs.options.ShowExpandedEvents)
	***REMOVED***

	if cs.options.StartAfter != nil ***REMOVED***
		var saDoc bsoncore.Document
		saDoc, cs.err = transformBsoncoreDocument(cs.registry, cs.options.StartAfter, true, "startAfter")
		if cs.err != nil ***REMOVED***
			return nil, cs.err
		***REMOVED***

		plDoc = bsoncore.AppendDocumentElement(plDoc, "startAfter", saDoc)
	***REMOVED***

	if cs.options.StartAtOperationTime != nil ***REMOVED***
		plDoc = bsoncore.AppendTimestampElement(plDoc, "startAtOperationTime", cs.options.StartAtOperationTime.T, cs.options.StartAtOperationTime.I)
	***REMOVED***

	// Append custom pipeline options.
	for optionName, optionValue := range cs.pipelineOptions ***REMOVED***
		plDoc = bsoncore.AppendValueElement(plDoc, optionName, optionValue)
	***REMOVED***

	if plDoc, cs.err = bsoncore.AppendDocumentEnd(plDoc, plDocIdx); cs.err != nil ***REMOVED***
		return nil, cs.err
	***REMOVED***

	return plDoc, nil
***REMOVED***

func (cs *ChangeStream) pipelineToBSON() (bsoncore.Document, error) ***REMOVED***
	pipelineDocIdx, pipelineArr := bsoncore.AppendArrayStart(nil)
	for i, doc := range cs.pipelineSlice ***REMOVED***
		pipelineArr = bsoncore.AppendDocumentElement(pipelineArr, strconv.Itoa(i), doc)
	***REMOVED***
	if pipelineArr, cs.err = bsoncore.AppendArrayEnd(pipelineArr, pipelineDocIdx); cs.err != nil ***REMOVED***
		return nil, cs.err
	***REMOVED***
	return pipelineArr, cs.err
***REMOVED***

func (cs *ChangeStream) replaceOptions(wireVersion *description.VersionRange) ***REMOVED***
	// Cached resume token: use the resume token as the resumeAfter option and set no other resume options
	if cs.resumeToken != nil ***REMOVED***
		cs.options.SetResumeAfter(cs.resumeToken)
		cs.options.SetStartAfter(nil)
		cs.options.SetStartAtOperationTime(nil)
		return
	***REMOVED***

	// No cached resume token but cached operation time: use the operation time as the startAtOperationTime option and
	// set no other resume options
	if (cs.sess.OperationTime != nil || cs.options.StartAtOperationTime != nil) && wireVersion.Max >= 7 ***REMOVED***
		opTime := cs.options.StartAtOperationTime
		if cs.operationTime != nil ***REMOVED***
			opTime = cs.sess.OperationTime
		***REMOVED***

		cs.options.SetStartAtOperationTime(opTime)
		cs.options.SetResumeAfter(nil)
		cs.options.SetStartAfter(nil)
		return
	***REMOVED***

	// No cached resume token or operation time: set none of the resume options
	cs.options.SetResumeAfter(nil)
	cs.options.SetStartAfter(nil)
	cs.options.SetStartAtOperationTime(nil)
***REMOVED***

// ID returns the ID for this change stream, or 0 if the cursor has been closed or exhausted.
func (cs *ChangeStream) ID() int64 ***REMOVED***
	if cs.cursor == nil ***REMOVED***
		return 0
	***REMOVED***
	return cs.cursor.ID()
***REMOVED***

// Decode will unmarshal the current event document into val and return any errors from the unmarshalling process
// without any modification. If val is nil or is a typed nil, an error will be returned.
func (cs *ChangeStream) Decode(val interface***REMOVED******REMOVED***) error ***REMOVED***
	if cs.cursor == nil ***REMOVED***
		return ErrNilCursor
	***REMOVED***

	return bson.UnmarshalWithRegistry(cs.registry, cs.Current, val)
***REMOVED***

// Err returns the last error seen by the change stream, or nil if no errors has occurred.
func (cs *ChangeStream) Err() error ***REMOVED***
	if cs.err != nil ***REMOVED***
		return replaceErrors(cs.err)
	***REMOVED***
	if cs.cursor == nil ***REMOVED***
		return nil
	***REMOVED***

	return replaceErrors(cs.cursor.Err())
***REMOVED***

// Close closes this change stream and the underlying cursor. Next and TryNext must not be called after Close has been
// called. Close is idempotent. After the first call, any subsequent calls will not change the state.
func (cs *ChangeStream) Close(ctx context.Context) error ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	defer closeImplicitSession(cs.sess)

	if cs.cursor == nil ***REMOVED***
		return nil // cursor is already closed
	***REMOVED***

	cs.err = replaceErrors(cs.cursor.Close(ctx))
	cs.cursor = nil
	return cs.Err()
***REMOVED***

// ResumeToken returns the last cached resume token for this change stream, or nil if a resume token has not been
// stored.
func (cs *ChangeStream) ResumeToken() bson.Raw ***REMOVED***
	return cs.resumeToken
***REMOVED***

// Next gets the next event for this change stream. It returns true if there were no errors and the next event document
// is available.
//
// Next blocks until an event is available, an error occurs, or ctx expires. If ctx expires, the error
// will be set to ctx.Err(). In an error case, Next will return false.
//
// If Next returns false, subsequent calls will also return false.
func (cs *ChangeStream) Next(ctx context.Context) bool ***REMOVED***
	return cs.next(ctx, false)
***REMOVED***

// TryNext attempts to get the next event for this change stream. It returns true if there were no errors and the next
// event document is available.
//
// TryNext returns false if the change stream is closed by the server, an error occurs when getting changes from the
// server, the next change is not yet available, or ctx expires. If ctx expires, the error will be set to ctx.Err().
//
// If TryNext returns false and an error occurred or the change stream was closed
// (i.e. cs.Err() != nil || cs.ID() == 0), subsequent attempts will also return false. Otherwise, it is safe to call
// TryNext again until a change is available.
//
// This method requires driver version >= 1.2.0.
func (cs *ChangeStream) TryNext(ctx context.Context) bool ***REMOVED***
	return cs.next(ctx, true)
***REMOVED***

func (cs *ChangeStream) next(ctx context.Context, nonBlocking bool) bool ***REMOVED***
	// return false right away if the change stream has already errored or if cursor is closed.
	if cs.err != nil ***REMOVED***
		return false
	***REMOVED***

	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	if len(cs.batch) == 0 ***REMOVED***
		cs.loopNext(ctx, nonBlocking)
		if cs.err != nil ***REMOVED***
			cs.err = replaceErrors(cs.err)
			return false
		***REMOVED***
		if len(cs.batch) == 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	// successfully got non-empty batch
	cs.Current = bson.Raw(cs.batch[0])
	cs.batch = cs.batch[1:]
	if cs.err = cs.storeResumeToken(); cs.err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func (cs *ChangeStream) loopNext(ctx context.Context, nonBlocking bool) ***REMOVED***
	for ***REMOVED***
		if cs.cursor == nil ***REMOVED***
			return
		***REMOVED***

		if cs.cursor.Next(ctx) ***REMOVED***
			// non-empty batch returned
			cs.batch, cs.err = cs.cursor.Batch().Documents()
			return
		***REMOVED***

		cs.err = replaceErrors(cs.cursor.Err())
		if cs.err == nil ***REMOVED***
			// Check if cursor is alive
			if cs.ID() == 0 ***REMOVED***
				return
			***REMOVED***

			// If a getMore was done but the batch was empty, the batch cursor will return false with no error.
			// Update the tracked resume token to catch the post batch resume token from the server response.
			cs.updatePbrtFromCommand()
			if nonBlocking ***REMOVED***
				// stop after a successful getMore, even though the batch was empty
				return
			***REMOVED***
			continue // loop getMore until a non-empty batch is returned or an error occurs
		***REMOVED***

		if !cs.isResumableError() ***REMOVED***
			return
		***REMOVED***

		// ignore error from cursor close because if the cursor is deleted or errors we tried to close it and will remake and try to get next batch
		_ = cs.cursor.Close(ctx)
		if cs.err = cs.executeOperation(ctx, true); cs.err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (cs *ChangeStream) isResumableError() bool ***REMOVED***
	commandErr, ok := cs.err.(CommandError)
	if !ok || commandErr.HasErrorLabel(networkErrorLabel) ***REMOVED***
		// All non-server errors or network errors are resumable.
		return true
	***REMOVED***

	if commandErr.Code == errorCursorNotFound ***REMOVED***
		return true
	***REMOVED***

	// For wire versions 9 and above, a server error is resumable if it has the ResumableChangeStreamError label.
	if cs.wireVersion != nil && cs.wireVersion.Includes(minResumableLabelWireVersion) ***REMOVED***
		return commandErr.HasErrorLabel(resumableErrorLabel)
	***REMOVED***

	// For wire versions below 9, a server error is resumable if its code is on the allowlist.
	_, resumable := resumableChangeStreamErrors[commandErr.Code]
	return resumable
***REMOVED***

// Returns true if the underlying cursor's batch is empty
func (cs *ChangeStream) emptyBatch() bool ***REMOVED***
	return cs.cursor.Batch().Empty()
***REMOVED***

// StreamType represents the cluster type against which a ChangeStream was created.
type StreamType uint8

// These constants represent valid change stream types. A change stream can be initialized over a collection, all
// collections in a database, or over a cluster.
const (
	CollectionStream StreamType = iota
	DatabaseStream
	ClientStream
)
