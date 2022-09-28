// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// BatchCursor is a batch implementation of a cursor. It returns documents in entire batches instead
// of one at a time. An individual document cursor can be built on top of this batch cursor.
type BatchCursor struct ***REMOVED***
	clientSession        *session.Client
	clock                *session.ClusterClock
	comment              bsoncore.Value
	database             string
	collection           string
	id                   int64
	err                  error
	server               Server
	serverDescription    description.Server
	errorProcessor       ErrorProcessor // This will only be set when pinning to a connection.
	connection           PinnedConnection
	batchSize            int32
	maxTimeMS            int64
	currentBatch         *bsoncore.DocumentSequence
	firstBatch           bool
	cmdMonitor           *event.CommandMonitor
	postBatchResumeToken bsoncore.Document
	crypt                Crypt
	serverAPI            *ServerAPIOptions

	// legacy server (< 3.2) fields
	limit       int32
	numReturned int32 // number of docs returned by server
***REMOVED***

// CursorResponse represents the response from a command the results in a cursor. A BatchCursor can
// be constructed from a CursorResponse.
type CursorResponse struct ***REMOVED***
	Server               Server
	ErrorProcessor       ErrorProcessor // This will only be set when pinning to a connection.
	Connection           PinnedConnection
	Desc                 description.Server
	FirstBatch           *bsoncore.DocumentSequence
	Database             string
	Collection           string
	ID                   int64
	postBatchResumeToken bsoncore.Document
***REMOVED***

// NewCursorResponse constructs a cursor response from the given response and server. This method
// can be used within the ProcessResponse method for an operation.
func NewCursorResponse(info ResponseInfo) (CursorResponse, error) ***REMOVED***
	response := info.ServerResponse
	cur, ok := response.Lookup("cursor").DocumentOK()
	if !ok ***REMOVED***
		return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("cursor should be an embedded document but is of BSON type %s", response.Lookup("cursor").Type)
	***REMOVED***
	elems, err := cur.Elements()
	if err != nil ***REMOVED***
		return CursorResponse***REMOVED******REMOVED***, err
	***REMOVED***
	curresp := CursorResponse***REMOVED***Server: info.Server, Desc: info.ConnectionDescription***REMOVED***

	for _, elem := range elems ***REMOVED***
		switch elem.Key() ***REMOVED***
		case "firstBatch":
			arr, ok := elem.Value().ArrayOK()
			if !ok ***REMOVED***
				return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("firstBatch should be an array but is a BSON %s", elem.Value().Type)
			***REMOVED***
			curresp.FirstBatch = &bsoncore.DocumentSequence***REMOVED***Style: bsoncore.ArrayStyle, Data: arr***REMOVED***
		case "ns":
			ns, ok := elem.Value().StringValueOK()
			if !ok ***REMOVED***
				return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("ns should be a string but is a BSON %s", elem.Value().Type)
			***REMOVED***
			index := strings.Index(ns, ".")
			if index == -1 ***REMOVED***
				return CursorResponse***REMOVED******REMOVED***, errors.New("ns field must contain a valid namespace, but is missing '.'")
			***REMOVED***
			curresp.Database = ns[:index]
			curresp.Collection = ns[index+1:]
		case "id":
			curresp.ID, ok = elem.Value().Int64OK()
			if !ok ***REMOVED***
				return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("id should be an int64 but it is a BSON %s", elem.Value().Type)
			***REMOVED***
		case "postBatchResumeToken":
			curresp.postBatchResumeToken, ok = elem.Value().DocumentOK()
			if !ok ***REMOVED***
				return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("post batch resume token should be a document but it is a BSON %s", elem.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If the deployment is behind a load balancer and the cursor has a non-zero ID, pin the cursor to a connection and
	// use the same connection to execute getMore and killCursors commands.
	if curresp.Desc.LoadBalanced() && curresp.ID != 0 ***REMOVED***
		// Cache the server as an ErrorProcessor to use when constructing deployments for cursor commands.
		ep, ok := curresp.Server.(ErrorProcessor)
		if !ok ***REMOVED***
			return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("expected Server used to establish a cursor to implement ErrorProcessor, but got %T", curresp.Server)
		***REMOVED***
		curresp.ErrorProcessor = ep

		refConn, ok := info.Connection.(PinnedConnection)
		if !ok ***REMOVED***
			return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("expected Connection used to establish a cursor to implement PinnedConnection, but got %T", info.Connection)
		***REMOVED***
		if err := refConn.PinToCursor(); err != nil ***REMOVED***
			return CursorResponse***REMOVED******REMOVED***, fmt.Errorf("error incrementing connection reference count when creating a cursor: %v", err)
		***REMOVED***
		curresp.Connection = refConn
	***REMOVED***

	return curresp, nil
***REMOVED***

// CursorOptions are extra options that are required to construct a BatchCursor.
type CursorOptions struct ***REMOVED***
	BatchSize      int32
	Comment        bsoncore.Value
	MaxTimeMS      int64
	Limit          int32
	CommandMonitor *event.CommandMonitor
	Crypt          Crypt
	ServerAPI      *ServerAPIOptions
***REMOVED***

// NewBatchCursor creates a new BatchCursor from the provided parameters.
func NewBatchCursor(cr CursorResponse, clientSession *session.Client, clock *session.ClusterClock, opts CursorOptions) (*BatchCursor, error) ***REMOVED***
	ds := cr.FirstBatch
	bc := &BatchCursor***REMOVED***
		clientSession:        clientSession,
		clock:                clock,
		comment:              opts.Comment,
		database:             cr.Database,
		collection:           cr.Collection,
		id:                   cr.ID,
		server:               cr.Server,
		connection:           cr.Connection,
		errorProcessor:       cr.ErrorProcessor,
		batchSize:            opts.BatchSize,
		maxTimeMS:            opts.MaxTimeMS,
		cmdMonitor:           opts.CommandMonitor,
		firstBatch:           true,
		postBatchResumeToken: cr.postBatchResumeToken,
		crypt:                opts.Crypt,
		serverAPI:            opts.ServerAPI,
		serverDescription:    cr.Desc,
	***REMOVED***

	if ds != nil ***REMOVED***
		bc.numReturned = int32(ds.DocumentCount())
	***REMOVED***
	if cr.Desc.WireVersion == nil || cr.Desc.WireVersion.Max < 4 ***REMOVED***
		bc.limit = opts.Limit

		// Take as many documents from the batch as needed.
		if bc.limit != 0 && bc.limit < bc.numReturned ***REMOVED***
			for i := int32(0); i < bc.limit; i++ ***REMOVED***
				_, err := ds.Next()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***
			ds.Data = ds.Data[:ds.Pos]
			ds.ResetIterator()
		***REMOVED***
	***REMOVED***

	bc.currentBatch = ds
	return bc, nil
***REMOVED***

// NewEmptyBatchCursor returns a batch cursor that is empty.
func NewEmptyBatchCursor() *BatchCursor ***REMOVED***
	return &BatchCursor***REMOVED***currentBatch: new(bsoncore.DocumentSequence)***REMOVED***
***REMOVED***

// NewBatchCursorFromDocuments returns a batch cursor with current batch set to a sequence-style
// DocumentSequence containing the provided documents.
func NewBatchCursorFromDocuments(documents []byte) *BatchCursor ***REMOVED***
	return &BatchCursor***REMOVED***
		currentBatch: &bsoncore.DocumentSequence***REMOVED***
			Data:  documents,
			Style: bsoncore.SequenceStyle,
		***REMOVED***,
		// BatchCursors created with this function have no associated ID nor server, so no getMore
		// calls will be made.
		id:     0,
		server: nil,
	***REMOVED***
***REMOVED***

// ID returns the cursor ID for this batch cursor.
func (bc *BatchCursor) ID() int64 ***REMOVED***
	return bc.id
***REMOVED***

// Next indicates if there is another batch available. Returning false does not necessarily indicate
// that the cursor is closed. This method will return false when an empty batch is returned.
//
// If Next returns true, there is a valid batch of documents available. If Next returns false, there
// is not a valid batch of documents available.
func (bc *BatchCursor) Next(ctx context.Context) bool ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	if bc.firstBatch ***REMOVED***
		bc.firstBatch = false
		return !bc.currentBatch.Empty()
	***REMOVED***

	if bc.id == 0 || bc.server == nil ***REMOVED***
		return false
	***REMOVED***

	bc.getMore(ctx)

	return !bc.currentBatch.Empty()
***REMOVED***

// Batch will return a DocumentSequence for the current batch of documents. The returned
// DocumentSequence is only valid until the next call to Next or Close.
func (bc *BatchCursor) Batch() *bsoncore.DocumentSequence ***REMOVED*** return bc.currentBatch ***REMOVED***

// Err returns the latest error encountered.
func (bc *BatchCursor) Err() error ***REMOVED*** return bc.err ***REMOVED***

// Close closes this batch cursor.
func (bc *BatchCursor) Close(ctx context.Context) error ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	err := bc.KillCursor(ctx)
	bc.id = 0
	bc.currentBatch.Data = nil
	bc.currentBatch.Style = 0
	bc.currentBatch.ResetIterator()

	connErr := bc.unpinConnection()
	if err == nil ***REMOVED***
		err = connErr
	***REMOVED***
	return err
***REMOVED***

func (bc *BatchCursor) unpinConnection() error ***REMOVED***
	if bc.connection == nil ***REMOVED***
		return nil
	***REMOVED***

	err := bc.connection.UnpinFromCursor()
	closeErr := bc.connection.Close()
	if err == nil && closeErr != nil ***REMOVED***
		err = closeErr
	***REMOVED***
	bc.connection = nil
	return err
***REMOVED***

// Server returns the server for this cursor.
func (bc *BatchCursor) Server() Server ***REMOVED***
	return bc.server
***REMOVED***

func (bc *BatchCursor) clearBatch() ***REMOVED***
	bc.currentBatch.Data = bc.currentBatch.Data[:0]
***REMOVED***

// KillCursor kills cursor on server without closing batch cursor
func (bc *BatchCursor) KillCursor(ctx context.Context) error ***REMOVED***
	if bc.server == nil || bc.id == 0 ***REMOVED***
		return nil
	***REMOVED***

	return Operation***REMOVED***
		CommandFn: func(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
			dst = bsoncore.AppendStringElement(dst, "killCursors", bc.collection)
			dst = bsoncore.BuildArrayElement(dst, "cursors", bsoncore.Value***REMOVED***Type: bsontype.Int64, Data: bsoncore.AppendInt64(nil, bc.id)***REMOVED***)
			return dst, nil
		***REMOVED***,
		Database:       bc.database,
		Deployment:     bc.getOperationDeployment(),
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyKillCursors,
		CommandMonitor: bc.cmdMonitor,
		ServerAPI:      bc.serverAPI,
	***REMOVED***.Execute(ctx, nil)
***REMOVED***

func (bc *BatchCursor) getMore(ctx context.Context) ***REMOVED***
	bc.clearBatch()
	if bc.id == 0 ***REMOVED***
		return
	***REMOVED***

	// Required for legacy operations which don't support limit.
	numToReturn := bc.batchSize
	if bc.limit != 0 && bc.numReturned+bc.batchSize >= bc.limit ***REMOVED***
		numToReturn = bc.limit - bc.numReturned
		if numToReturn <= 0 ***REMOVED***
			err := bc.Close(ctx)
			if err != nil ***REMOVED***
				bc.err = err
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	bc.err = Operation***REMOVED***
		CommandFn: func(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
			dst = bsoncore.AppendInt64Element(dst, "getMore", bc.id)
			dst = bsoncore.AppendStringElement(dst, "collection", bc.collection)
			if numToReturn > 0 ***REMOVED***
				dst = bsoncore.AppendInt32Element(dst, "batchSize", numToReturn)
			***REMOVED***
			if bc.maxTimeMS > 0 ***REMOVED***
				dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", bc.maxTimeMS)
			***REMOVED***
			// The getMore command does not support commenting pre-4.4.
			if bc.comment.Type != bsontype.Type(0) && bc.serverDescription.WireVersion.Max >= 9 ***REMOVED***
				dst = bsoncore.AppendValueElement(dst, "comment", bc.comment)
			***REMOVED***
			return dst, nil
		***REMOVED***,
		Database:   bc.database,
		Deployment: bc.getOperationDeployment(),
		ProcessResponseFn: func(info ResponseInfo) error ***REMOVED***
			response := info.ServerResponse
			id, ok := response.Lookup("cursor", "id").Int64OK()
			if !ok ***REMOVED***
				return fmt.Errorf("cursor.id should be an int64 but is a BSON %s", response.Lookup("cursor", "id").Type)
			***REMOVED***
			bc.id = id

			batch, ok := response.Lookup("cursor", "nextBatch").ArrayOK()
			if !ok ***REMOVED***
				return fmt.Errorf("cursor.nextBatch should be an array but is a BSON %s", response.Lookup("cursor", "nextBatch").Type)
			***REMOVED***
			bc.currentBatch.Style = bsoncore.ArrayStyle
			bc.currentBatch.Data = batch
			bc.currentBatch.ResetIterator()
			bc.numReturned += int32(bc.currentBatch.DocumentCount()) // Required for legacy operations which don't support limit.

			pbrt, err := response.LookupErr("cursor", "postBatchResumeToken")
			if err != nil ***REMOVED***
				// I don't really understand why we don't set bc.err here
				return nil
			***REMOVED***

			pbrtDoc, ok := pbrt.DocumentOK()
			if !ok ***REMOVED***
				bc.err = fmt.Errorf("expected BSON type for post batch resume token to be EmbeddedDocument but got %s", pbrt.Type)
				return nil
			***REMOVED***

			bc.postBatchResumeToken = pbrtDoc

			return nil
		***REMOVED***,
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyGetMore,
		CommandMonitor: bc.cmdMonitor,
		Crypt:          bc.crypt,
		ServerAPI:      bc.serverAPI,
	***REMOVED***.Execute(ctx, nil)

	// Once the cursor has been drained, we can unpin the connection if one is currently pinned.
	if bc.id == 0 ***REMOVED***
		err := bc.unpinConnection()
		if err != nil && bc.err == nil ***REMOVED***
			bc.err = err
		***REMOVED***
	***REMOVED***

	// If we're in load balanced mode and the pinned connection encounters a network error, we should not use it for
	// future commands. Per the spec, the connection will not be unpinned until the cursor is actually closed, but
	// we set the cursor ID to 0 to ensure the Close() call will not execute a killCursors command.
	if driverErr, ok := bc.err.(Error); ok && driverErr.NetworkError() && bc.connection != nil ***REMOVED***
		bc.id = 0
	***REMOVED***

	// Required for legacy operations which don't support limit.
	if bc.limit != 0 && bc.numReturned >= bc.limit ***REMOVED***
		// call KillCursor instead of Close because Close will clear out the data for the current batch.
		err := bc.KillCursor(ctx)
		if err != nil && bc.err == nil ***REMOVED***
			bc.err = err
		***REMOVED***
	***REMOVED***
***REMOVED***

// PostBatchResumeToken returns the latest seen post batch resume token.
func (bc *BatchCursor) PostBatchResumeToken() bsoncore.Document ***REMOVED***
	return bc.postBatchResumeToken
***REMOVED***

// SetBatchSize sets the batchSize for future getMores.
func (bc *BatchCursor) SetBatchSize(size int32) ***REMOVED***
	bc.batchSize = size
***REMOVED***

func (bc *BatchCursor) getOperationDeployment() Deployment ***REMOVED***
	if bc.connection != nil ***REMOVED***
		return &loadBalancedCursorDeployment***REMOVED***
			errorProcessor: bc.errorProcessor,
			conn:           bc.connection,
		***REMOVED***
	***REMOVED***
	return SingleServerDeployment***REMOVED***bc.server***REMOVED***
***REMOVED***

// loadBalancedCursorDeployment is used as a Deployment for getMore and killCursors commands when pinning to a
// connection in load balanced mode. This type also functions as an ErrorProcessor to ensure that SDAM errors are
// handled for these commands in this mode.
type loadBalancedCursorDeployment struct ***REMOVED***
	errorProcessor ErrorProcessor
	conn           PinnedConnection
***REMOVED***

var _ Deployment = (*loadBalancedCursorDeployment)(nil)
var _ Server = (*loadBalancedCursorDeployment)(nil)
var _ ErrorProcessor = (*loadBalancedCursorDeployment)(nil)

func (lbcd *loadBalancedCursorDeployment) SelectServer(_ context.Context, _ description.ServerSelector) (Server, error) ***REMOVED***
	return lbcd, nil
***REMOVED***

func (lbcd *loadBalancedCursorDeployment) Kind() description.TopologyKind ***REMOVED***
	return description.LoadBalanced
***REMOVED***

func (lbcd *loadBalancedCursorDeployment) Connection(_ context.Context) (Connection, error) ***REMOVED***
	return lbcd.conn, nil
***REMOVED***

// MinRTT always returns 0. It implements the driver.Server interface.
func (lbcd *loadBalancedCursorDeployment) MinRTT() time.Duration ***REMOVED***
	return 0
***REMOVED***

// RTT90 always returns 0. It implements the driver.Server interface.
func (lbcd *loadBalancedCursorDeployment) RTT90() time.Duration ***REMOVED***
	return 0
***REMOVED***

func (lbcd *loadBalancedCursorDeployment) ProcessError(err error, conn Connection) ProcessErrorResult ***REMOVED***
	return lbcd.errorProcessor.ProcessError(err, conn)
***REMOVED***
