// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const defaultLocalThreshold = 15 * time.Millisecond

var dollarCmd = [...]byte***REMOVED***'.', '$', 'c', 'm', 'd'***REMOVED***

var (
	// ErrNoDocCommandResponse occurs when the server indicated a response existed, but none was found.
	ErrNoDocCommandResponse = errors.New("command returned no documents")
	// ErrMultiDocCommandResponse occurs when the server sent multiple documents in response to a command.
	ErrMultiDocCommandResponse = errors.New("command returned multiple documents")
	// ErrReplyDocumentMismatch occurs when the number of documents returned in an OP_QUERY does not match the numberReturned field.
	ErrReplyDocumentMismatch = errors.New("number of documents returned does not match numberReturned field")
	// ErrNonPrimaryReadPref is returned when a read is attempted in a transaction with a non-primary read preference.
	ErrNonPrimaryReadPref = errors.New("read preference in a transaction must be primary")
)

const (
	// maximum BSON object size when client side encryption is enabled
	cryptMaxBsonObjectSize uint32 = 2097152
	// minimum wire version necessary to use automatic encryption
	cryptMinWireVersion int32 = 8
	// minimum wire version necessary to use read snapshots
	readSnapshotMinWireVersion int32 = 13
)

// RetryablePoolError is a connection pool error that can be retried while executing an operation.
type RetryablePoolError interface ***REMOVED***
	Retryable() bool
***REMOVED***

// InvalidOperationError is returned from Validate and indicates that a required field is missing
// from an instance of Operation.
type InvalidOperationError struct***REMOVED*** MissingField string ***REMOVED***

func (err InvalidOperationError) Error() string ***REMOVED***
	return "the " + err.MissingField + " field must be set on Operation"
***REMOVED***

// opReply stores information returned in an OP_REPLY response from the server.
// The err field stores any error that occurred when decoding or validating the OP_REPLY response.
type opReply struct ***REMOVED***
	responseFlags wiremessage.ReplyFlag
	cursorID      int64
	startingFrom  int32
	numReturned   int32
	documents     []bsoncore.Document
	err           error
***REMOVED***

// startedInformation keeps track of all of the information necessary for monitoring started events.
type startedInformation struct ***REMOVED***
	cmd                      bsoncore.Document
	requestID                int32
	cmdName                  string
	documentSequenceIncluded bool
	connID                   string
	serverConnID             *int32
	redacted                 bool
	serviceID                *primitive.ObjectID
***REMOVED***

// finishedInformation keeps track of all of the information necessary for monitoring success and failure events.
type finishedInformation struct ***REMOVED***
	cmdName      string
	requestID    int32
	response     bsoncore.Document
	cmdErr       error
	connID       string
	serverConnID *int32
	startTime    time.Time
	redacted     bool
	serviceID    *primitive.ObjectID
***REMOVED***

// ResponseInfo contains the context required to parse a server response.
type ResponseInfo struct ***REMOVED***
	ServerResponse        bsoncore.Document
	Server                Server
	Connection            Connection
	ConnectionDescription description.Server
	CurrentIndex          int
***REMOVED***

// Operation is used to execute an operation. It contains all of the common code required to
// select a server, transform an operation into a command, write the command to a connection from
// the selected server, read a response from that connection, process the response, and potentially
// retry.
//
// The required fields are Database, CommandFn, and Deployment. All other fields are optional.
//
// While an Operation can be constructed manually, drivergen should be used to generate an
// implementation of an operation instead. This will ensure that there are helpers for constructing
// the operation and that this type isn't configured incorrectly.
type Operation struct ***REMOVED***
	// CommandFn is used to create the command that will be wrapped in a wire message and sent to
	// the server. This function should only add the elements of the command and not start or end
	// the enclosing BSON document. Per the command API, the first element must be the name of the
	// command to run. This field is required.
	CommandFn func(dst []byte, desc description.SelectedServer) ([]byte, error)

	// Database is the database that the command will be run against. This field is required.
	Database string

	// Deployment is the MongoDB Deployment to use. While most of the time this will be multiple
	// servers, commands that need to run against a single, preselected server can use the
	// SingleServerDeployment type. Commands that need to run on a preselected connection can use
	// the SingleConnectionDeployment type.
	Deployment Deployment

	// ProcessResponseFn is called after a response to the command is returned. The server is
	// provided for types like Cursor that are required to run subsequent commands using the same
	// server.
	ProcessResponseFn func(ResponseInfo) error

	// Selector is the server selector that's used during both initial server selection and
	// subsequent selection for retries. Depending on the Deployment implementation, the
	// SelectServer method may not actually be called.
	Selector description.ServerSelector

	// ReadPreference is the read preference that will be attached to the command. If this field is
	// not specified a default read preference of primary will be used.
	ReadPreference *readpref.ReadPref

	// ReadConcern is the read concern used when running read commands. This field should not be set
	// for write operations. If this field is set, it will be encoded onto the commands sent to the
	// server.
	ReadConcern *readconcern.ReadConcern

	// MinimumReadConcernWireVersion specifies the minimum wire version to add the read concern to
	// the command being executed.
	MinimumReadConcernWireVersion int32

	// WriteConcern is the write concern used when running write commands. This field should not be
	// set for read operations. If this field is set, it will be encoded onto the commands sent to
	// the server.
	WriteConcern *writeconcern.WriteConcern

	// MinimumWriteConcernWireVersion specifies the minimum wire version to add the write concern to
	// the command being executed.
	MinimumWriteConcernWireVersion int32

	// Client is the session used with this operation. This can be either an implicit or explicit
	// session. If the server selected does not support sessions and Client is specified the
	// behavior depends on the session type. If the session is implicit, the session fields will not
	// be encoded onto the command. If the session is explicit, an error will be returned. The
	// caller is responsible for ensuring that this field is nil if the Deployment does not support
	// sessions.
	Client *session.Client

	// Clock is a cluster clock, different from the one contained within a session.Client. This
	// allows updating cluster times for a global cluster clock while allowing individual session's
	// cluster clocks to be only updated as far as the last command that's been run.
	Clock *session.ClusterClock

	// RetryMode specifies how to retry. There are three modes that enable retry: RetryOnce,
	// RetryOncePerCommand, and RetryContext. For more information about what these modes do, please
	// refer to their definitions. Both RetryMode and Type must be set for retryability to be enabled.
	RetryMode *RetryMode

	// Type specifies the kind of operation this is. There is only one mode that enables retry: Write.
	// For more information about what this mode does, please refer to it's definition. Both Type and
	// RetryMode must be set for retryability to be enabled.
	Type Type

	// Batches contains the documents that are split when executing a write command that potentially
	// has more documents than can fit in a single command. This should only be specified for
	// commands that are batch compatible. For more information, please refer to the definition of
	// Batches.
	Batches *Batches

	// Legacy sets the legacy type for this operation. There are only 3 types that require legacy
	// support: find, getMore, and killCursors. For more information about LegacyOperationKind,
	// please refer to it's definition.
	Legacy LegacyOperationKind

	// CommandMonitor specifies the monitor to use for APM events. If this field is not set,
	// no events will be reported.
	CommandMonitor *event.CommandMonitor

	// Crypt specifies a Crypt object to use for automatic client side encryption and decryption.
	Crypt Crypt

	// ServerAPI specifies options used to configure the API version sent to the server.
	ServerAPI *ServerAPIOptions

	// IsOutputAggregate specifies whether this operation is an aggregate with an output stage. If true,
	// read preference will not be added to the command on wire versions < 13.
	IsOutputAggregate bool

	// Timeout is the amount of time that this operation can execute before returning an error. The default value
	// nil, which means that the timeout of the operation's caller will be used.
	Timeout *time.Duration

	// cmdName is only set when serializing OP_MSG and is used internally in readWireMessage.
	cmdName string
***REMOVED***

// shouldEncrypt returns true if this operation should automatically be encrypted.
func (op Operation) shouldEncrypt() bool ***REMOVED***
	return op.Crypt != nil && !op.Crypt.BypassAutoEncryption()
***REMOVED***

// selectServer handles performing server selection for an operation.
func (op Operation) selectServer(ctx context.Context) (Server, error) ***REMOVED***
	if err := op.Validate(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	selector := op.Selector
	if selector == nil ***REMOVED***
		rp := op.ReadPreference
		if rp == nil ***REMOVED***
			rp = readpref.Primary()
		***REMOVED***
		selector = description.CompositeSelector([]description.ServerSelector***REMOVED***
			description.ReadPrefSelector(rp),
			description.LatencySelector(defaultLocalThreshold),
		***REMOVED***)
	***REMOVED***

	return op.Deployment.SelectServer(ctx, selector)
***REMOVED***

// getServerAndConnection should be used to retrieve a Server and Connection to execute an operation.
func (op Operation) getServerAndConnection(ctx context.Context) (Server, Connection, error) ***REMOVED***
	server, err := op.selectServer(ctx)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// If the provided client session has a pinned connection, it should be used for the operation because this
	// indicates that we're in a transaction and the target server is behind a load balancer.
	if op.Client != nil && op.Client.PinnedConnection != nil ***REMOVED***
		return server, op.Client.PinnedConnection, nil
	***REMOVED***

	// Otherwise, default to checking out a connection from the server's pool.
	conn, err := server.Connection(ctx)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// If we're in load balanced mode and this is the first operation in a transaction, pin the session to a connection.
	if conn.Description().LoadBalanced() && op.Client != nil && op.Client.TransactionStarting() ***REMOVED***
		pinnedConn, ok := conn.(PinnedConnection)
		if !ok ***REMOVED***
			// Close the original connection to avoid a leak.
			_ = conn.Close()
			return nil, nil, fmt.Errorf("expected Connection used to start a transaction to be a PinnedConnection, but got %T", conn)
		***REMOVED***
		if err := pinnedConn.PinToTransaction(); err != nil ***REMOVED***
			// Close the original connection to avoid a leak.
			_ = conn.Close()
			return nil, nil, fmt.Errorf("error incrementing connection reference count when starting a transaction: %v", err)
		***REMOVED***
		op.Client.PinnedConnection = pinnedConn
	***REMOVED***

	return server, conn, nil
***REMOVED***

// Validate validates this operation, ensuring the fields are set properly.
func (op Operation) Validate() error ***REMOVED***
	if op.CommandFn == nil ***REMOVED***
		return InvalidOperationError***REMOVED***MissingField: "CommandFn"***REMOVED***
	***REMOVED***
	if op.Deployment == nil ***REMOVED***
		return InvalidOperationError***REMOVED***MissingField: "Deployment"***REMOVED***
	***REMOVED***
	if op.Database == "" ***REMOVED***
		return InvalidOperationError***REMOVED***MissingField: "Database"***REMOVED***
	***REMOVED***
	if op.Client != nil && !writeconcern.AckWrite(op.WriteConcern) ***REMOVED***
		return errors.New("session provided for an unacknowledged write")
	***REMOVED***
	return nil
***REMOVED***

// Execute runs this operation. The scratch parameter will be used and overwritten (potentially many
// times), this should mainly be used to enable pooling of byte slices.
func (op Operation) Execute(ctx context.Context, scratch []byte) error ***REMOVED***
	err := op.Validate()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If no deadline is set on the passed-in context, op.Timeout is set, and context is not already
	// a Timeout context, honor op.Timeout in new Timeout context for operation execution.
	if _, deadlineSet := ctx.Deadline(); !deadlineSet && op.Timeout != nil && !internal.IsTimeoutContext(ctx) ***REMOVED***
		newCtx, cancelFunc := internal.MakeTimeoutContext(ctx, *op.Timeout)
		// Redefine ctx to be the new timeout-derived context.
		ctx = newCtx
		// Cancel the timeout-derived context at the end of Execute to avoid a context leak.
		defer cancelFunc()
	***REMOVED***

	if op.Client != nil ***REMOVED***
		if err := op.Client.StartCommand(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var retries int
	if op.RetryMode != nil ***REMOVED***
		switch op.Type ***REMOVED***
		case Write:
			if op.Client == nil ***REMOVED***
				break
			***REMOVED***
			switch *op.RetryMode ***REMOVED***
			case RetryOnce, RetryOncePerCommand:
				retries = 1
			case RetryContext:
				retries = -1
			***REMOVED***
		case Read:
			switch *op.RetryMode ***REMOVED***
			case RetryOnce, RetryOncePerCommand:
				retries = 1
			case RetryContext:
				retries = -1
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var srvr Server
	var conn Connection
	var res bsoncore.Document
	var operationErr WriteCommandError
	var prevErr error
	batching := op.Batches.Valid()
	retryEnabled := op.RetryMode != nil && op.RetryMode.Enabled()
	retrySupported := false
	first := true
	currIndex := 0

	// resetForRetry records the error that caused the retry, decrements retries, and resets the
	// retry loop variables to request a new server and a new connection for the next attempt.
	resetForRetry := func(err error) ***REMOVED***
		retries--
		prevErr = err
		// If we got a connection, close it immediately to release pool resources for
		// subsequent retries.
		if conn != nil ***REMOVED***
			conn.Close()
		***REMOVED***
		// Set the server and connection to nil to request a new server and connection.
		srvr = nil
		conn = nil
	***REMOVED***

	for ***REMOVED***
		// If the server or connection are nil, try to select a new server and get a new connection.
		if srvr == nil || conn == nil ***REMOVED***
			srvr, conn, err = op.getServerAndConnection(ctx)
			if err != nil ***REMOVED***
				// If the returned error is retryable and there are retries remaining (negative
				// retries means retry indefinitely), then retry the operation. Set the server
				// and connection to nil to request a new server and connection.
				if rerr, ok := err.(RetryablePoolError); ok && rerr.Retryable() && retries != 0 ***REMOVED***
					resetForRetry(err)
					continue
				***REMOVED***

				// If this is a retry and there's an error from a previous attempt, return the previous
				// error instead of the current connection error.
				if prevErr != nil ***REMOVED***
					return prevErr
				***REMOVED***
				return err
			***REMOVED***
			defer conn.Close()
		***REMOVED***

		// Run steps that must only be run on the first attempt, but not again for retries.
		if first ***REMOVED***
			// Determine if retries are supported for the current operation on the current server
			// description. Per the retryable writes specification, only determine this for the
			// first server selected:
			//
			//   If the server selected for the first attempt of a retryable write operation does
			//   not support retryable writes, drivers MUST execute the write as if retryable writes
			//   were not enabled.
			retrySupported = op.retryable(conn.Description())

			// If retries are supported for the current operation on the current server description,
			// client retries are enabled, the operation type is write, and we haven't incremented
			// the txn number yet, enable retry writes on the session and increment the txn number.
			// Calling IncrementTxnNumber() for server descriptions or topologies that do not
			// support retries (e.g. standalone topologies) will cause server errors. Only do this
			// check for the first attempt to keep retried writes in the same transaction.
			if retrySupported && op.RetryMode != nil && op.Type == Write && op.Client != nil ***REMOVED***
				op.Client.RetryWrite = false
				if op.RetryMode.Enabled() ***REMOVED***
					op.Client.RetryWrite = true
					if !op.Client.Committing && !op.Client.Aborting ***REMOVED***
						op.Client.IncrementTxnNumber()
					***REMOVED***
				***REMOVED***
			***REMOVED***

			first = false
		***REMOVED***

		desc := description.SelectedServer***REMOVED***Server: conn.Description(), Kind: op.Deployment.Kind()***REMOVED***
		scratch = scratch[:0]
		if desc.WireVersion == nil || desc.WireVersion.Max < 4 ***REMOVED***
			switch op.Legacy ***REMOVED***
			case LegacyFind:
				return op.legacyFind(ctx, scratch, srvr, conn, desc)
			case LegacyGetMore:
				return op.legacyGetMore(ctx, scratch, srvr, conn, desc)
			case LegacyKillCursors:
				return op.legacyKillCursors(ctx, scratch, srvr, conn, desc)
			***REMOVED***
		***REMOVED***
		if desc.WireVersion == nil || desc.WireVersion.Max < 3 ***REMOVED***
			switch op.Legacy ***REMOVED***
			case LegacyListCollections:
				return op.legacyListCollections(ctx, scratch, srvr, conn, desc)
			case LegacyListIndexes:
				return op.legacyListIndexes(ctx, scratch, srvr, conn, desc)
			***REMOVED***
		***REMOVED***

		if batching ***REMOVED***
			targetBatchSize := desc.MaxDocumentSize
			maxDocSize := desc.MaxDocumentSize
			if op.shouldEncrypt() ***REMOVED***
				// For client-side encryption, we want the batch to be split at 2 MiB instead of 16MiB.
				// If there's only one document in the batch, it can be up to 16MiB, so we set target batch size to
				// 2MiB but max document size to 16MiB. This will allow the AdvanceBatch call to create a batch
				// with a single large document.
				targetBatchSize = cryptMaxBsonObjectSize
			***REMOVED***

			err = op.Batches.AdvanceBatch(int(desc.MaxBatchCount), int(targetBatchSize), int(maxDocSize))
			if err != nil ***REMOVED***
				// TODO(GODRIVER-982): Should we also be returning operationErr?
				return err
			***REMOVED***
		***REMOVED***

		// Calculate value of 'maxTimeMS' field to potentially append to the wire message based on the current
		// context's deadline and the 90th percentile RTT if the ctx is a Timeout Context.
		var maxTimeMS uint64
		if internal.IsTimeoutContext(ctx) ***REMOVED***
			if deadline, ok := ctx.Deadline(); ok ***REMOVED***
				remainingTimeout := time.Until(deadline)

				maxTimeMSVal := int64(remainingTimeout/time.Millisecond) -
					int64(srvr.RTT90()/time.Millisecond)

				// A maxTimeMS value <= 0 indicates that we are already at or past the Context's deadline.
				if maxTimeMSVal <= 0 ***REMOVED***
					return internal.WrapErrorf(ErrDeadlineWouldBeExceeded,
						"Context deadline has already been surpassed by %v", remainingTimeout)
				***REMOVED***
				maxTimeMS = uint64(maxTimeMSVal)
			***REMOVED***
		***REMOVED***

		// convert to wire message
		if len(scratch) > 0 ***REMOVED***
			scratch = scratch[:0]
		***REMOVED***
		wm, startedInfo, err := op.createWireMessage(ctx, scratch, desc, maxTimeMS, conn)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// set extra data and send event if possible
		startedInfo.connID = conn.ID()
		startedInfo.cmdName = op.getCommandName(startedInfo.cmd)
		op.cmdName = startedInfo.cmdName
		startedInfo.redacted = op.redactCommand(startedInfo.cmdName, startedInfo.cmd)
		startedInfo.serviceID = conn.Description().ServiceID
		startedInfo.serverConnID = conn.ServerConnectionID()
		op.publishStartedEvent(ctx, startedInfo)

		// get the moreToCome flag information before we compress
		moreToCome := wiremessage.IsMsgMoreToCome(wm)

		// compress wiremessage if allowed
		if compressor, ok := conn.(Compressor); ok && op.canCompress(startedInfo.cmdName) ***REMOVED***
			wm, err = compressor.CompressWireMessage(wm, nil)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		finishedInfo := finishedInformation***REMOVED***
			cmdName:      startedInfo.cmdName,
			requestID:    startedInfo.requestID,
			startTime:    time.Now(),
			connID:       startedInfo.connID,
			serverConnID: startedInfo.serverConnID,
			redacted:     startedInfo.redacted,
			serviceID:    startedInfo.serviceID,
		***REMOVED***

		// Check for possible context error. If no context error, check if there's enough time to perform a
		// round trip before the Context deadline. If ctx is a Timeout Context, use the 90th percentile RTT
		// as a threshold. Otherwise, use the minimum observed RTT.
		if ctx.Err() != nil ***REMOVED***
			err = ctx.Err()
		***REMOVED*** else if deadline, ok := ctx.Deadline(); ok ***REMOVED***
			if internal.IsTimeoutContext(ctx) && time.Now().Add(srvr.RTT90()).After(deadline) ***REMOVED***
				err = internal.WrapErrorf(ErrDeadlineWouldBeExceeded,
					"Remaining timeout %v applied from Timeout is less than 90th percentile RTT", time.Until(deadline))
			***REMOVED*** else if time.Now().Add(srvr.MinRTT()).After(deadline) ***REMOVED***
				err = context.DeadlineExceeded
			***REMOVED***
		***REMOVED***

		if err == nil ***REMOVED***
			// roundtrip using either the full roundTripper or a special one for when the moreToCome
			// flag is set
			var roundTrip = op.roundTrip
			if moreToCome ***REMOVED***
				roundTrip = op.moreToComeRoundTrip
			***REMOVED***
			res, err = roundTrip(ctx, conn, wm)

			if ep, ok := srvr.(ErrorProcessor); ok ***REMOVED***
				_ = ep.ProcessError(err, conn)
			***REMOVED***
		***REMOVED***

		finishedInfo.response = res
		finishedInfo.cmdErr = err
		op.publishFinishedEvent(ctx, finishedInfo)

		var perr error
		switch tt := err.(type) ***REMOVED***
		case WriteCommandError:
			if e := err.(WriteCommandError); retrySupported && op.Type == Write && e.UnsupportedStorageEngine() ***REMOVED***
				return ErrUnsupportedStorageEngine
			***REMOVED***

			connDesc := conn.Description()
			retryableErr := tt.Retryable(connDesc.WireVersion)
			preRetryWriteLabelVersion := connDesc.WireVersion != nil && connDesc.WireVersion.Max < 9
			inTransaction := op.Client != nil &&
				!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning()
			// If retry is enabled and the operation isn't in a transaction, add a RetryableWriteError label for
			// retryable errors from pre-4.4 servers
			if retryableErr && preRetryWriteLabelVersion && retryEnabled && !inTransaction ***REMOVED***
				tt.Labels = append(tt.Labels, RetryableWriteError)
			***REMOVED***

			// If retries are supported for the current operation on the first server description,
			// the error is considered retryable, and there are retries remaining (negative retries
			// means retry indefinitely), then retry the operation.
			if retrySupported && retryableErr && retries != 0 ***REMOVED***
				if op.Client != nil && op.Client.Committing ***REMOVED***
					// Apply majority write concern for retries
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				***REMOVED***
				resetForRetry(tt)
				continue
			***REMOVED***

			// If the operation isn't being retried, process the response
			if op.ProcessResponseFn != nil ***REMOVED***
				info := ResponseInfo***REMOVED***
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				***REMOVED***
				_ = op.ProcessResponseFn(info)
			***REMOVED***

			if batching && len(tt.WriteErrors) > 0 && currIndex > 0 ***REMOVED***
				for i := range tt.WriteErrors ***REMOVED***
					tt.WriteErrors[i].Index += int64(currIndex)
				***REMOVED***
			***REMOVED***

			// If batching is enabled and either ordered is the default (which is true) or
			// explicitly set to true and we have write errors, return the errors.
			if batching && (op.Batches.Ordered == nil || *op.Batches.Ordered) && len(tt.WriteErrors) > 0 ***REMOVED***
				return tt
			***REMOVED***
			if op.Client != nil && op.Client.Committing && tt.WriteConcernError != nil ***REMOVED***
				// When running commitTransaction we return WriteConcernErrors as an Error.
				err := Error***REMOVED***
					Name:    tt.WriteConcernError.Name,
					Code:    int32(tt.WriteConcernError.Code),
					Message: tt.WriteConcernError.Message,
					Labels:  tt.Labels,
					Raw:     tt.Raw,
				***REMOVED***
				// The UnknownTransactionCommitResult label is added to all writeConcernErrors besides unknownReplWriteConcernCode
				// and unsatisfiableWriteConcernCode
				if err.Code != unknownReplWriteConcernCode && err.Code != unsatisfiableWriteConcernCode ***REMOVED***
					err.Labels = append(err.Labels, UnknownTransactionCommitResult)
				***REMOVED***
				if retryableErr && retryEnabled ***REMOVED***
					err.Labels = append(err.Labels, RetryableWriteError)
				***REMOVED***
				return err
			***REMOVED***
			operationErr.WriteConcernError = tt.WriteConcernError
			operationErr.WriteErrors = append(operationErr.WriteErrors, tt.WriteErrors...)
			operationErr.Labels = tt.Labels
			operationErr.Raw = tt.Raw
		case Error:
			if tt.HasErrorLabel(TransientTransactionError) || tt.HasErrorLabel(UnknownTransactionCommitResult) ***REMOVED***
				if err := op.Client.ClearPinnedResources(); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***

			if e := err.(Error); retrySupported && op.Type == Write && e.UnsupportedStorageEngine() ***REMOVED***
				return ErrUnsupportedStorageEngine
			***REMOVED***

			connDesc := conn.Description()
			var retryableErr bool
			if op.Type == Write ***REMOVED***
				retryableErr = tt.RetryableWrite(connDesc.WireVersion)
				preRetryWriteLabelVersion := connDesc.WireVersion != nil && connDesc.WireVersion.Max < 9
				inTransaction := op.Client != nil &&
					!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning()
				// If retryWrites is enabled and the operation isn't in a transaction, add a RetryableWriteError label
				// for network errors and retryable errors from pre-4.4 servers
				if retryEnabled && !inTransaction &&
					(tt.HasErrorLabel(NetworkError) || (retryableErr && preRetryWriteLabelVersion)) ***REMOVED***
					tt.Labels = append(tt.Labels, RetryableWriteError)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				retryableErr = tt.RetryableRead()
			***REMOVED***

			// If retries are supported for the current operation on the first server description,
			// the error is considered retryable, and there are retries remaining (negative retries
			// means retry indefinitely), then retry the operation.
			if retrySupported && retryableErr && retries != 0 ***REMOVED***
				if op.Client != nil && op.Client.Committing ***REMOVED***
					// Apply majority write concern for retries
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				***REMOVED***
				resetForRetry(tt)
				continue
			***REMOVED***

			// If the operation isn't being retried, process the response
			if op.ProcessResponseFn != nil ***REMOVED***
				info := ResponseInfo***REMOVED***
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				***REMOVED***
				_ = op.ProcessResponseFn(info)
			***REMOVED***

			if op.Client != nil && op.Client.Committing && (retryableErr || tt.Code == 50) ***REMOVED***
				// If we got a retryable error or MaxTimeMSExpired error, we add UnknownTransactionCommitResult.
				tt.Labels = append(tt.Labels, UnknownTransactionCommitResult)
			***REMOVED***
			return tt
		case nil:
			if moreToCome ***REMOVED***
				return ErrUnacknowledgedWrite
			***REMOVED***
			if op.ProcessResponseFn != nil ***REMOVED***
				info := ResponseInfo***REMOVED***
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				***REMOVED***
				perr = op.ProcessResponseFn(info)
			***REMOVED***
			if perr != nil ***REMOVED***
				return perr
			***REMOVED***
		default:
			if op.ProcessResponseFn != nil ***REMOVED***
				info := ResponseInfo***REMOVED***
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				***REMOVED***
				_ = op.ProcessResponseFn(info)
			***REMOVED***
			return err
		***REMOVED***

		// If we're batching and there are batches remaining, advance to the next batch. This isn't
		// a retry, so increment the transaction number, reset the retries number, and don't set
		// server or connection to nil to continue using the same connection.
		if batching && len(op.Batches.Documents) > 0 ***REMOVED***
			if retrySupported && op.Client != nil && op.RetryMode != nil ***REMOVED***
				if *op.RetryMode > RetryNone ***REMOVED***
					op.Client.IncrementTxnNumber()
				***REMOVED***
				if *op.RetryMode == RetryOncePerCommand ***REMOVED***
					retries = 1
				***REMOVED***
			***REMOVED***
			currIndex += len(op.Batches.Current)
			op.Batches.ClearBatch()
			continue
		***REMOVED***
		break
	***REMOVED***
	if len(operationErr.WriteErrors) > 0 || operationErr.WriteConcernError != nil ***REMOVED***
		return operationErr
	***REMOVED***
	return nil
***REMOVED***

// Retryable writes are supported if the server supports sessions, the operation is not
// within a transaction, and the write is acknowledged
func (op Operation) retryable(desc description.Server) bool ***REMOVED***
	switch op.Type ***REMOVED***
	case Write:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) ***REMOVED***
			return true
		***REMOVED***
		if retryWritesSupported(desc) &&
			desc.WireVersion != nil && desc.WireVersion.Max >= 6 &&
			op.Client != nil && !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) &&
			writeconcern.AckWrite(op.WriteConcern) ***REMOVED***
			return true
		***REMOVED***
	case Read:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) ***REMOVED***
			return true
		***REMOVED***
		if desc.WireVersion != nil && desc.WireVersion.Max >= 6 &&
			(op.Client == nil || !(op.Client.TransactionInProgress() || op.Client.TransactionStarting())) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// roundTrip writes a wiremessage to the connection and then reads a wiremessage. The wm parameter
// is reused when reading the wiremessage.
func (op Operation) roundTrip(ctx context.Context, conn Connection, wm []byte) ([]byte, error) ***REMOVED***
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil ***REMOVED***
		return nil, op.networkError(err)
	***REMOVED***

	return op.readWireMessage(ctx, conn, wm)
***REMOVED***

func (op Operation) readWireMessage(ctx context.Context, conn Connection, wm []byte) ([]byte, error) ***REMOVED***
	var err error

	wm, err = conn.ReadWireMessage(ctx, wm[:0])
	if err != nil ***REMOVED***
		return nil, op.networkError(err)
	***REMOVED***

	// If we're using a streamable connection, we set its streaming state based on the moreToCome flag in the server
	// response.
	if streamer, ok := conn.(StreamerConnection); ok ***REMOVED***
		streamer.SetStreaming(wiremessage.IsMsgMoreToCome(wm))
	***REMOVED***

	// decompress wiremessage
	wm, err = op.decompressWireMessage(wm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// decode
	res, err := op.decodeResult(wm)
	// Update cluster/operation time and recovery tokens before handling the error to ensure we're properly updating
	// everything.
	op.updateClusterTimes(res)
	op.updateOperationTime(res)
	op.Client.UpdateRecoveryToken(bson.Raw(res))

	// Update snapshot time if operation was a "find", "aggregate" or "distinct".
	if op.cmdName == "find" || op.cmdName == "aggregate" || op.cmdName == "distinct" ***REMOVED***
		op.Client.UpdateSnapshotTime(res)
	***REMOVED***

	if err != nil ***REMOVED***
		return res, err
	***REMOVED***

	// If there is no error, automatically attempt to decrypt all results if client side encryption is enabled.
	if op.Crypt != nil ***REMOVED***
		return op.Crypt.Decrypt(ctx, res)
	***REMOVED***
	return res, nil
***REMOVED***

// networkError wraps the provided error in an Error with label "NetworkError" and, if a transaction
// is running or committing, the appropriate transaction state labels. The returned error indicates
// the operation should be retried for reads and writes. If err is nil, networkError returns nil.
func (op Operation) networkError(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	labels := []string***REMOVED***NetworkError***REMOVED***
	if op.Client != nil ***REMOVED***
		op.Client.MarkDirty()
	***REMOVED***
	if op.Client != nil && op.Client.TransactionRunning() && !op.Client.Committing ***REMOVED***
		labels = append(labels, TransientTransactionError)
	***REMOVED***
	if op.Client != nil && op.Client.Committing ***REMOVED***
		labels = append(labels, UnknownTransactionCommitResult)
	***REMOVED***
	return Error***REMOVED***Message: err.Error(), Labels: labels, Wrapped: err***REMOVED***
***REMOVED***

// moreToComeRoundTrip writes a wiremessage to the provided connection. This is used when an OP_MSG is
// being sent with  the moreToCome bit set.
func (op *Operation) moreToComeRoundTrip(ctx context.Context, conn Connection, wm []byte) ([]byte, error) ***REMOVED***
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil ***REMOVED***
		if op.Client != nil ***REMOVED***
			op.Client.MarkDirty()
		***REMOVED***
		err = Error***REMOVED***Message: err.Error(), Labels: []string***REMOVED***TransientTransactionError, NetworkError***REMOVED***, Wrapped: err***REMOVED***
	***REMOVED***
	return bsoncore.BuildDocument(nil, bsoncore.AppendInt32Element(nil, "ok", 1)), err
***REMOVED***

// decompressWireMessage handles decompressing a wiremessage. If the wiremessage
// is not compressed, this method will return the wiremessage.
func (Operation) decompressWireMessage(wm []byte) ([]byte, error) ***REMOVED***
	// read the header and ensure this is a compressed wire message
	length, reqid, respto, opcode, rem, ok := wiremessage.ReadHeader(wm)
	if !ok || len(wm) < int(length) ***REMOVED***
		return nil, errors.New("malformed wire message: insufficient bytes")
	***REMOVED***
	if opcode != wiremessage.OpCompressed ***REMOVED***
		return wm, nil
	***REMOVED***
	// get the original opcode and uncompressed size
	opcode, rem, ok = wiremessage.ReadCompressedOriginalOpCode(rem)
	if !ok ***REMOVED***
		return nil, errors.New("malformed OP_COMPRESSED: missing original opcode")
	***REMOVED***
	uncompressedSize, rem, ok := wiremessage.ReadCompressedUncompressedSize(rem)
	if !ok ***REMOVED***
		return nil, errors.New("malformed OP_COMPRESSED: missing uncompressed size")
	***REMOVED***
	// get the compressor ID and decompress the message
	compressorID, rem, ok := wiremessage.ReadCompressedCompressorID(rem)
	if !ok ***REMOVED***
		return nil, errors.New("malformed OP_COMPRESSED: missing compressor ID")
	***REMOVED***
	compressedSize := length - 25 // header (16) + original opcode (4) + uncompressed size (4) + compressor ID (1)
	// return the original wiremessage
	msg, rem, ok := wiremessage.ReadCompressedCompressedMessage(rem, compressedSize)
	if !ok ***REMOVED***
		return nil, errors.New("malformed OP_COMPRESSED: insufficient bytes for compressed wiremessage")
	***REMOVED***

	header := make([]byte, 0, uncompressedSize+16)
	header = wiremessage.AppendHeader(header, uncompressedSize+16, reqid, respto, opcode)
	opts := CompressionOpts***REMOVED***
		Compressor:       compressorID,
		UncompressedSize: uncompressedSize,
	***REMOVED***
	uncompressed, err := DecompressPayload(msg, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return append(header, uncompressed...), nil
***REMOVED***

func (op Operation) createWireMessage(
	ctx context.Context,
	dst []byte,
	desc description.SelectedServer,
	maxTimeMS uint64,
	conn Connection) ([]byte, startedInformation, error) ***REMOVED***

	// If topology is not LoadBalanced, API version is not declared, and wire version is unknown
	// or less than 6, use OP_QUERY. Otherwise, use OP_MSG.
	if desc.Kind != description.LoadBalanced && op.ServerAPI == nil &&
		(desc.WireVersion == nil || desc.WireVersion.Max < wiremessage.OpmsgWireVersion) ***REMOVED***
		return op.createQueryWireMessage(maxTimeMS, dst, desc)
	***REMOVED***
	return op.createMsgWireMessage(ctx, maxTimeMS, dst, desc, conn)
***REMOVED***

func (op Operation) addBatchArray(dst []byte) []byte ***REMOVED***
	aidx, dst := bsoncore.AppendArrayElementStart(dst, op.Batches.Identifier)
	for i, doc := range op.Batches.Current ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, strconv.Itoa(i), doc)
	***REMOVED***
	dst, _ = bsoncore.AppendArrayEnd(dst, aidx)
	return dst
***REMOVED***

func (op Operation) createQueryWireMessage(maxTimeMS uint64, dst []byte, desc description.SelectedServer) ([]byte, startedInformation, error) ***REMOVED***
	var info startedInformation
	flags := op.secondaryOK(desc)
	var wmindex int32
	info.requestID = wiremessage.NextRequestID()
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, flags)
	// FullCollectionName
	dst = append(dst, op.Database...)
	dst = append(dst, dollarCmd[:]...)
	dst = append(dst, 0x00)
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, -1)

	wrapper := int32(-1)
	rp, err := op.createReadPref(desc, true)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***
	if len(rp) > 0 ***REMOVED***
		wrapper, dst = bsoncore.AppendDocumentStart(dst)
		dst = bsoncore.AppendHeader(dst, bsontype.EmbeddedDocument, "$query")
	***REMOVED***
	idx, dst := bsoncore.AppendDocumentStart(dst)
	dst, err = op.CommandFn(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***

	if op.Batches != nil && len(op.Batches.Current) > 0 ***REMOVED***
		dst = op.addBatchArray(dst)
	***REMOVED***

	dst, err = op.addReadConcern(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***

	dst, err = op.addWriteConcern(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***

	dst, err = op.addSession(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***

	dst = op.addClusterTime(dst, desc)
	dst = op.addServerAPI(dst)
	// If maxTimeMS is greater than 0 append it to wire message. A maxTimeMS value of 0 only explicitly
	// specifies the default behavior of no timeout server-side.
	if maxTimeMS > 0 ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", int64(maxTimeMS))
	***REMOVED***

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)
	// Command monitoring only reports the document inside $query
	info.cmd = dst[idx:]

	if len(rp) > 0 ***REMOVED***
		var err error
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
		dst, err = bsoncore.AppendDocumentEnd(dst, wrapper)
		if err != nil ***REMOVED***
			return dst, info, err
		***REMOVED***
	***REMOVED***

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), info, nil
***REMOVED***

func (op Operation) createMsgWireMessage(ctx context.Context, maxTimeMS uint64, dst []byte, desc description.SelectedServer,
	conn Connection) ([]byte, startedInformation, error) ***REMOVED***

	var info startedInformation
	var flags wiremessage.MsgFlag
	var wmindex int32
	// We set the MoreToCome bit if we have a write concern, it's unacknowledged, and we either
	// aren't batching or we are encoding the last batch.
	if op.WriteConcern != nil && !writeconcern.AckWrite(op.WriteConcern) && (op.Batches == nil || len(op.Batches.Documents) == 0) ***REMOVED***
		flags = wiremessage.MoreToCome
	***REMOVED***
	// Set the ExhaustAllowed flag if the connection supports streaming. This will tell the server that it can
	// respond with the MoreToCome flag and then stream responses over this connection.
	if streamer, ok := conn.(StreamerConnection); ok && streamer.SupportsStreaming() ***REMOVED***
		flags |= wiremessage.ExhaustAllowed
	***REMOVED***

	info.requestID = wiremessage.NextRequestID()
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpMsg)
	dst = wiremessage.AppendMsgFlags(dst, flags)
	// Body
	dst = wiremessage.AppendMsgSectionType(dst, wiremessage.SingleDocument)

	idx, dst := bsoncore.AppendDocumentStart(dst)

	dst, err := op.addCommandFields(ctx, dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***
	dst, err = op.addReadConcern(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***
	dst, err = op.addWriteConcern(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***
	dst, err = op.addSession(dst, desc)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***

	dst = op.addClusterTime(dst, desc)
	dst = op.addServerAPI(dst)
	// If maxTimeMS is greater than 0 append it to wire message. A maxTimeMS value of 0 only explicitly
	// specifies the default behavior of no timeout server-side.
	if maxTimeMS > 0 ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", int64(maxTimeMS))
	***REMOVED***

	dst = bsoncore.AppendStringElement(dst, "$db", op.Database)
	rp, err := op.createReadPref(desc, false)
	if err != nil ***REMOVED***
		return dst, info, err
	***REMOVED***
	if len(rp) > 0 ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
	***REMOVED***

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)
	// The command document for monitoring shouldn't include the type 1 payload as a document sequence
	info.cmd = dst[idx:]

	// add batch as a document sequence if auto encryption is not enabled
	// if auto encryption is enabled, the batch will already be an array in the command document
	if !op.shouldEncrypt() && op.Batches != nil && len(op.Batches.Current) > 0 ***REMOVED***
		info.documentSequenceIncluded = true
		dst = wiremessage.AppendMsgSectionType(dst, wiremessage.DocumentSequence)
		idx, dst = bsoncore.ReserveLength(dst)

		dst = append(dst, op.Batches.Identifier...)
		dst = append(dst, 0x00)

		for _, doc := range op.Batches.Current ***REMOVED***
			dst = append(dst, doc...)
		***REMOVED***

		dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	***REMOVED***

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), info, nil
***REMOVED***

// addCommandFields adds the fields for a command to the wire message in dst. This assumes that the start of the document
// has already been added and does not add the final 0 byte.
func (op Operation) addCommandFields(ctx context.Context, dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	if !op.shouldEncrypt() ***REMOVED***
		return op.CommandFn(dst, desc)
	***REMOVED***

	if desc.WireVersion.Max < cryptMinWireVersion ***REMOVED***
		return dst, errors.New("auto-encryption requires a MongoDB version of 4.2")
	***REMOVED***

	// create temporary command document
	cidx, cmdDst := bsoncore.AppendDocumentStart(nil)
	var err error
	cmdDst, err = op.CommandFn(cmdDst, desc)
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***
	// use a BSON array instead of a type 1 payload because mongocryptd will convert to arrays regardless
	if op.Batches != nil && len(op.Batches.Current) > 0 ***REMOVED***
		cmdDst = op.addBatchArray(cmdDst)
	***REMOVED***
	cmdDst, _ = bsoncore.AppendDocumentEnd(cmdDst, cidx)

	// encrypt the command
	encrypted, err := op.Crypt.Encrypt(ctx, op.Database, cmdDst)
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***
	// append encrypted command to original destination, removing the first 4 bytes (length) and final byte (terminator)
	dst = append(dst, encrypted[4:len(encrypted)-1]...)
	return dst, nil
***REMOVED***

// addServerAPI adds the relevant fields for server API specification to the wire message in dst.
func (op Operation) addServerAPI(dst []byte) []byte ***REMOVED***
	sa := op.ServerAPI
	if sa == nil ***REMOVED***
		return dst
	***REMOVED***

	dst = bsoncore.AppendStringElement(dst, "apiVersion", sa.ServerAPIVersion)
	if sa.Strict != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "apiStrict", *sa.Strict)
	***REMOVED***
	if sa.DeprecationErrors != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "apiDeprecationErrors", *sa.DeprecationErrors)
	***REMOVED***
	return dst
***REMOVED***

func (op Operation) addReadConcern(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	if op.MinimumReadConcernWireVersion > 0 && (desc.WireVersion == nil || !desc.WireVersion.Includes(op.MinimumReadConcernWireVersion)) ***REMOVED***
		return dst, nil
	***REMOVED***
	rc := op.ReadConcern
	client := op.Client
	// Starting transaction's read concern overrides all others
	if client != nil && client.TransactionStarting() && client.CurrentRc != nil ***REMOVED***
		rc = client.CurrentRc
	***REMOVED***

	// start transaction must append afterclustertime IF causally consistent and operation time exists
	if rc == nil && client != nil && client.TransactionStarting() && client.Consistent && client.OperationTime != nil ***REMOVED***
		rc = readconcern.New()
	***REMOVED***

	if client != nil && client.Snapshot ***REMOVED***
		if desc.WireVersion.Max < readSnapshotMinWireVersion ***REMOVED***
			return dst, errors.New("snapshot reads require MongoDB 5.0 or later")
		***REMOVED***
		rc = readconcern.Snapshot()
	***REMOVED***

	if rc == nil ***REMOVED***
		return dst, nil
	***REMOVED***

	_, data, err := rc.MarshalBSONValue() // always returns a document
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***

	if sessionsSupported(desc.WireVersion) && client != nil ***REMOVED***
		if client.Consistent && client.OperationTime != nil ***REMOVED***
			data = data[:len(data)-1] // remove the null byte
			data = bsoncore.AppendTimestampElement(data, "afterClusterTime", client.OperationTime.T, client.OperationTime.I)
			data, _ = bsoncore.AppendDocumentEnd(data, 0)
		***REMOVED***
		if client.Snapshot && client.SnapshotTime != nil ***REMOVED***
			data = data[:len(data)-1] // remove the null byte
			data = bsoncore.AppendTimestampElement(data, "atClusterTime", client.SnapshotTime.T, client.SnapshotTime.I)
			data, _ = bsoncore.AppendDocumentEnd(data, 0)
		***REMOVED***
	***REMOVED***

	if len(data) == bsoncore.EmptyDocumentLength ***REMOVED***
		return dst, nil
	***REMOVED***
	return bsoncore.AppendDocumentElement(dst, "readConcern", data), nil
***REMOVED***

func (op Operation) addWriteConcern(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	if op.MinimumWriteConcernWireVersion > 0 && (desc.WireVersion == nil || !desc.WireVersion.Includes(op.MinimumWriteConcernWireVersion)) ***REMOVED***
		return dst, nil
	***REMOVED***
	wc := op.WriteConcern
	if wc == nil ***REMOVED***
		return dst, nil
	***REMOVED***

	t, data, err := wc.MarshalBSONValue()
	if err == writeconcern.ErrEmptyWriteConcern ***REMOVED***
		return dst, nil
	***REMOVED***
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***

	return append(bsoncore.AppendHeader(dst, t, "writeConcern"), data...), nil
***REMOVED***

func (op Operation) addSession(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	client := op.Client
	if client == nil || !sessionsSupported(desc.WireVersion) || desc.SessionTimeoutMinutes == 0 ***REMOVED***
		return dst, nil
	***REMOVED***
	if err := client.UpdateUseTime(); err != nil ***REMOVED***
		return dst, err
	***REMOVED***
	dst = bsoncore.AppendDocumentElement(dst, "lsid", client.SessionID)

	var addedTxnNumber bool
	if op.Type == Write && client.RetryWrite ***REMOVED***
		addedTxnNumber = true
		dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
	***REMOVED***
	if client.TransactionRunning() || client.RetryingCommit ***REMOVED***
		if !addedTxnNumber ***REMOVED***
			dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
		***REMOVED***
		if client.TransactionStarting() ***REMOVED***
			dst = bsoncore.AppendBooleanElement(dst, "startTransaction", true)
		***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "autocommit", false)
	***REMOVED***

	return dst, client.ApplyCommand(desc.Server)
***REMOVED***

func (op Operation) addClusterTime(dst []byte, desc description.SelectedServer) []byte ***REMOVED***
	client, clock := op.Client, op.Clock
	if (clock == nil && client == nil) || !sessionsSupported(desc.WireVersion) ***REMOVED***
		return dst
	***REMOVED***
	clusterTime := clock.GetClusterTime()
	if client != nil ***REMOVED***
		clusterTime = session.MaxClusterTime(clusterTime, client.ClusterTime)
	***REMOVED***
	if clusterTime == nil ***REMOVED***
		return dst
	***REMOVED***
	val, err := clusterTime.LookupErr("$clusterTime")
	if err != nil ***REMOVED***
		return dst
	***REMOVED***
	return append(bsoncore.AppendHeader(dst, val.Type, "$clusterTime"), val.Value...)
	// return bsoncore.AppendDocumentElement(dst, "$clusterTime", clusterTime)
***REMOVED***

// updateClusterTimes updates the cluster times for the session and cluster clock attached to this
// operation. While the session's AdvanceClusterTime may return an error, this method does not
// because an error being returned from this method will not be returned further up.
func (op Operation) updateClusterTimes(response bsoncore.Document) ***REMOVED***
	// Extract cluster time.
	value, err := response.LookupErr("$clusterTime")
	if err != nil ***REMOVED***
		// $clusterTime not included by the server
		return
	***REMOVED***
	clusterTime := bsoncore.BuildDocumentFromElements(nil, bsoncore.AppendValueElement(nil, "$clusterTime", value))

	sess, clock := op.Client, op.Clock

	if sess != nil ***REMOVED***
		_ = sess.AdvanceClusterTime(bson.Raw(clusterTime))
	***REMOVED***

	if clock != nil ***REMOVED***
		clock.AdvanceClusterTime(bson.Raw(clusterTime))
	***REMOVED***
***REMOVED***

// updateOperationTime updates the operation time on the session attached to this operation. While
// the session's AdvanceOperationTime method may return an error, this method does not because an
// error being returned from this method will not be returned further up.
func (op Operation) updateOperationTime(response bsoncore.Document) ***REMOVED***
	sess := op.Client
	if sess == nil ***REMOVED***
		return
	***REMOVED***

	opTimeElem, err := response.LookupErr("operationTime")
	if err != nil ***REMOVED***
		// operationTime not included by the server
		return
	***REMOVED***

	t, i := opTimeElem.Timestamp()
	_ = sess.AdvanceOperationTime(&primitive.Timestamp***REMOVED***
		T: t,
		I: i,
	***REMOVED***)
***REMOVED***

func (op Operation) getReadPrefBasedOnTransaction() (*readpref.ReadPref, error) ***REMOVED***
	if op.Client != nil && op.Client.TransactionRunning() ***REMOVED***
		// Transaction's read preference always takes priority
		rp := op.Client.CurrentRp
		// Reads in a transaction must have read preference primary
		// This must not be checked in startTransaction
		if rp != nil && !op.Client.TransactionStarting() && rp.Mode() != readpref.PrimaryMode ***REMOVED***
			return nil, ErrNonPrimaryReadPref
		***REMOVED***
		return rp, nil
	***REMOVED***
	return op.ReadPreference, nil
***REMOVED***

func (op Operation) createReadPref(desc description.SelectedServer, isOpQuery bool) (bsoncore.Document, error) ***REMOVED***
	// TODO(GODRIVER-2231): Instead of checking if isOutputAggregate and desc.Server.WireVersion.Max < 13, somehow check
	// TODO if supplied readPreference was "overwritten" with primary in description.selectForReplicaSet.
	if desc.Server.Kind == description.Standalone || (isOpQuery && desc.Server.Kind != description.Mongos) ||
		op.Type == Write || (op.IsOutputAggregate && desc.Server.WireVersion.Max < 13) ***REMOVED***
		// Don't send read preference for:
		// 1. all standalones
		// 2. non-mongos when using OP_QUERY
		// 3. all writes
		// 4. when operation is an aggregate with an output stage, and selected server's wire
		//    version is < 13
		return nil, nil
	***REMOVED***

	idx, doc := bsoncore.AppendDocumentStart(nil)
	rp, err := op.getReadPrefBasedOnTransaction()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if rp == nil ***REMOVED***
		if desc.Kind == description.Single && desc.Server.Kind != description.Mongos ***REMOVED***
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		***REMOVED***
		return nil, nil
	***REMOVED***

	switch rp.Mode() ***REMOVED***
	case readpref.PrimaryMode:
		if desc.Server.Kind == description.Mongos ***REMOVED***
			return nil, nil
		***REMOVED***
		if desc.Kind == description.Single ***REMOVED***
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		***REMOVED***
		doc = bsoncore.AppendStringElement(doc, "mode", "primary")
	case readpref.PrimaryPreferredMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
	case readpref.SecondaryPreferredMode:
		_, ok := rp.MaxStaleness()
		if desc.Server.Kind == description.Mongos && isOpQuery && !ok && len(rp.TagSets()) == 0 && rp.HedgeEnabled() == nil ***REMOVED***
			return nil, nil
		***REMOVED***
		doc = bsoncore.AppendStringElement(doc, "mode", "secondaryPreferred")
	case readpref.SecondaryMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "secondary")
	case readpref.NearestMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "nearest")
	***REMOVED***

	sets := make([]bsoncore.Document, 0, len(rp.TagSets()))
	for _, ts := range rp.TagSets() ***REMOVED***
		i, set := bsoncore.AppendDocumentStart(nil)
		for _, t := range ts ***REMOVED***
			set = bsoncore.AppendStringElement(set, t.Name, t.Value)
		***REMOVED***
		set, _ = bsoncore.AppendDocumentEnd(set, i)
		sets = append(sets, set)
	***REMOVED***
	if len(sets) > 0 ***REMOVED***
		var aidx int32
		aidx, doc = bsoncore.AppendArrayElementStart(doc, "tags")
		for i, set := range sets ***REMOVED***
			doc = bsoncore.AppendDocumentElement(doc, strconv.Itoa(i), set)
		***REMOVED***
		doc, _ = bsoncore.AppendArrayEnd(doc, aidx)
	***REMOVED***

	if d, ok := rp.MaxStaleness(); ok ***REMOVED***
		doc = bsoncore.AppendInt32Element(doc, "maxStalenessSeconds", int32(d.Seconds()))
	***REMOVED***

	if hedgeEnabled := rp.HedgeEnabled(); hedgeEnabled != nil ***REMOVED***
		var hedgeIdx int32
		hedgeIdx, doc = bsoncore.AppendDocumentElementStart(doc, "hedge")
		doc = bsoncore.AppendBooleanElement(doc, "enabled", *hedgeEnabled)
		doc, err = bsoncore.AppendDocumentEnd(doc, hedgeIdx)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("error creating hedge document: %v", err)
		***REMOVED***
	***REMOVED***

	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc, nil
***REMOVED***

func (op Operation) secondaryOK(desc description.SelectedServer) wiremessage.QueryFlag ***REMOVED***
	if desc.Kind == description.Single && desc.Server.Kind != description.Mongos ***REMOVED***
		return wiremessage.SecondaryOK
	***REMOVED***

	if rp := op.ReadPreference; rp != nil && rp.Mode() != readpref.PrimaryMode ***REMOVED***
		return wiremessage.SecondaryOK
	***REMOVED***

	return 0
***REMOVED***

func (Operation) canCompress(cmd string) bool ***REMOVED***
	if cmd == internal.LegacyHello || cmd == "hello" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "authenticate" ||
		cmd == "createUser" || cmd == "updateUser" || cmd == "copydbSaslStart" || cmd == "copydbgetnonce" || cmd == "copydb" ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// decodeOpReply extracts the necessary information from an OP_REPLY wire message.
// includesHeader: specifies whether or not wm includes the message header
// Returns the decoded OP_REPLY. If the err field of the returned opReply is non-nil, an error occurred while decoding
// or validating the response and the other fields are undefined.
func (Operation) decodeOpReply(wm []byte, includesHeader bool) opReply ***REMOVED***
	var reply opReply
	var ok bool

	if includesHeader ***REMOVED***
		wmLength := len(wm)
		var length int32
		var opcode wiremessage.OpCode
		length, _, _, opcode, wm, ok = wiremessage.ReadHeader(wm)
		if !ok || int(length) > wmLength ***REMOVED***
			reply.err = errors.New("malformed wire message: insufficient bytes")
			return reply
		***REMOVED***
		if opcode != wiremessage.OpReply ***REMOVED***
			reply.err = errors.New("malformed wire message: incorrect opcode")
			return reply
		***REMOVED***
	***REMOVED***

	reply.responseFlags, wm, ok = wiremessage.ReadReplyFlags(wm)
	if !ok ***REMOVED***
		reply.err = errors.New("malformed OP_REPLY: missing flags")
		return reply
	***REMOVED***
	reply.cursorID, wm, ok = wiremessage.ReadReplyCursorID(wm)
	if !ok ***REMOVED***
		reply.err = errors.New("malformed OP_REPLY: missing cursorID")
		return reply
	***REMOVED***
	reply.startingFrom, wm, ok = wiremessage.ReadReplyStartingFrom(wm)
	if !ok ***REMOVED***
		reply.err = errors.New("malformed OP_REPLY: missing startingFrom")
		return reply
	***REMOVED***
	reply.numReturned, wm, ok = wiremessage.ReadReplyNumberReturned(wm)
	if !ok ***REMOVED***
		reply.err = errors.New("malformed OP_REPLY: missing numberReturned")
		return reply
	***REMOVED***
	reply.documents, wm, ok = wiremessage.ReadReplyDocuments(wm)
	if !ok ***REMOVED***
		reply.err = errors.New("malformed OP_REPLY: could not read documents from reply")
	***REMOVED***

	if reply.responseFlags&wiremessage.QueryFailure == wiremessage.QueryFailure ***REMOVED***
		reply.err = QueryFailureError***REMOVED***
			Message:  "command failure",
			Response: reply.documents[0],
		***REMOVED***
		return reply
	***REMOVED***
	if reply.responseFlags&wiremessage.CursorNotFound == wiremessage.CursorNotFound ***REMOVED***
		reply.err = ErrCursorNotFound
		return reply
	***REMOVED***
	if reply.numReturned != int32(len(reply.documents)) ***REMOVED***
		reply.err = ErrReplyDocumentMismatch
		return reply
	***REMOVED***

	return reply
***REMOVED***

func (op Operation) decodeResult(wm []byte) (bsoncore.Document, error) ***REMOVED***
	wmLength := len(wm)
	length, _, _, opcode, wm, ok := wiremessage.ReadHeader(wm)
	if !ok || int(length) > wmLength ***REMOVED***
		return nil, errors.New("malformed wire message: insufficient bytes")
	***REMOVED***

	wm = wm[:wmLength-16] // constrain to just this wiremessage, incase there are multiple in the slice

	switch opcode ***REMOVED***
	case wiremessage.OpReply:
		reply := op.decodeOpReply(wm, false)
		if reply.err != nil ***REMOVED***
			return nil, reply.err
		***REMOVED***
		if reply.numReturned == 0 ***REMOVED***
			return nil, ErrNoDocCommandResponse
		***REMOVED***
		if reply.numReturned > 1 ***REMOVED***
			return nil, ErrMultiDocCommandResponse
		***REMOVED***
		rdr := reply.documents[0]
		if err := rdr.Validate(); err != nil ***REMOVED***
			return nil, NewCommandResponseError("malformed OP_REPLY: invalid document", err)
		***REMOVED***

		return rdr, ExtractErrorFromServerResponse(rdr)
	case wiremessage.OpMsg:
		_, wm, ok = wiremessage.ReadMsgFlags(wm)
		if !ok ***REMOVED***
			return nil, errors.New("malformed wire message: missing OP_MSG flags")
		***REMOVED***

		var res bsoncore.Document
		for len(wm) > 0 ***REMOVED***
			var stype wiremessage.SectionType
			stype, wm, ok = wiremessage.ReadMsgSectionType(wm)
			if !ok ***REMOVED***
				return nil, errors.New("malformed wire message: insuffienct bytes to read section type")
			***REMOVED***

			switch stype ***REMOVED***
			case wiremessage.SingleDocument:
				res, wm, ok = wiremessage.ReadMsgSectionSingleDocument(wm)
				if !ok ***REMOVED***
					return nil, errors.New("malformed wire message: insufficient bytes to read single document")
				***REMOVED***
			case wiremessage.DocumentSequence:
				// TODO(GODRIVER-617): Implement document sequence returns.
				_, _, wm, ok = wiremessage.ReadMsgSectionDocumentSequence(wm)
				if !ok ***REMOVED***
					return nil, errors.New("malformed wire message: insufficient bytes to read document sequence")
				***REMOVED***
			default:
				return nil, fmt.Errorf("malformed wire message: uknown section type %v", stype)
			***REMOVED***
		***REMOVED***

		err := res.Validate()
		if err != nil ***REMOVED***
			return nil, NewCommandResponseError("malformed OP_MSG: invalid document", err)
		***REMOVED***

		return res, ExtractErrorFromServerResponse(res)
	default:
		return nil, fmt.Errorf("cannot decode result from %s", opcode)
	***REMOVED***
***REMOVED***

// getCommandName returns the name of the command from the given BSON document.
func (op Operation) getCommandName(doc []byte) string ***REMOVED***
	// skip 4 bytes for document length and 1 byte for element type
	idx := bytes.IndexByte(doc[5:], 0x00) // look for the 0 byte after the command name
	return string(doc[5 : idx+5])
***REMOVED***

func (op *Operation) redactCommand(cmd string, doc bsoncore.Document) bool ***REMOVED***
	if cmd == "authenticate" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "createUser" ||
		cmd == "updateUser" || cmd == "copydbgetnonce" || cmd == "copydbsaslstart" || cmd == "copydb" ***REMOVED***

		return true
	***REMOVED***
	if strings.ToLower(cmd) != internal.LegacyHelloLowercase && cmd != "hello" ***REMOVED***
		return false
	***REMOVED***

	// A hello without speculative authentication can be monitored.
	_, err := doc.LookupErr("speculativeAuthenticate")
	return err == nil
***REMOVED***

// publishStartedEvent publishes a CommandStartedEvent to the operation's command monitor if possible. If the command is
// an unacknowledged write, a CommandSucceededEvent will be published as well. If started events are not being monitored,
// no events are published.
func (op Operation) publishStartedEvent(ctx context.Context, info startedInformation) ***REMOVED***
	if op.CommandMonitor == nil || op.CommandMonitor.Started == nil ***REMOVED***
		return
	***REMOVED***

	// Make a copy of the command. Redact if the command is security sensitive and cannot be monitored.
	// If there was a type 1 payload for the current batch, convert it to a BSON array.
	cmdCopy := bson.Raw***REMOVED******REMOVED***
	if !info.redacted ***REMOVED***
		cmdCopy = make([]byte, len(info.cmd))
		copy(cmdCopy, info.cmd)
		if info.documentSequenceIncluded ***REMOVED***
			cmdCopy = cmdCopy[:len(info.cmd)-1] // remove 0 byte at end
			cmdCopy = op.addBatchArray(cmdCopy)
			cmdCopy, _ = bsoncore.AppendDocumentEnd(cmdCopy, 0) // add back 0 byte and update length
		***REMOVED***
	***REMOVED***

	started := &event.CommandStartedEvent***REMOVED***
		Command:            cmdCopy,
		DatabaseName:       op.Database,
		CommandName:        info.cmdName,
		RequestID:          int64(info.requestID),
		ConnectionID:       info.connID,
		ServerConnectionID: info.serverConnID,
		ServiceID:          info.serviceID,
	***REMOVED***
	op.CommandMonitor.Started(ctx, started)
***REMOVED***

// publishFinishedEvent publishes either a CommandSucceededEvent or a CommandFailedEvent to the operation's command
// monitor if possible. If success/failure events aren't being monitored, no events are published.
func (op Operation) publishFinishedEvent(ctx context.Context, info finishedInformation) ***REMOVED***
	success := info.cmdErr == nil
	if _, ok := info.cmdErr.(WriteCommandError); ok ***REMOVED***
		success = true
	***REMOVED***
	if op.CommandMonitor == nil || (success && op.CommandMonitor.Succeeded == nil) || (!success && op.CommandMonitor.Failed == nil) ***REMOVED***
		return
	***REMOVED***

	var durationNanos int64
	var emptyTime time.Time
	if info.startTime != emptyTime ***REMOVED***
		durationNanos = time.Since(info.startTime).Nanoseconds()
	***REMOVED***

	finished := event.CommandFinishedEvent***REMOVED***
		CommandName:        info.cmdName,
		RequestID:          int64(info.requestID),
		ConnectionID:       info.connID,
		DurationNanos:      durationNanos,
		ServerConnectionID: info.serverConnID,
		ServiceID:          info.serviceID,
	***REMOVED***

	if success ***REMOVED***
		res := bson.Raw***REMOVED******REMOVED***
		// Only copy the reply for commands that are not security sensitive
		if !info.redacted ***REMOVED***
			res = make([]byte, len(info.response))
			copy(res, info.response)
		***REMOVED***
		successEvent := &event.CommandSucceededEvent***REMOVED***
			Reply:                res,
			CommandFinishedEvent: finished,
		***REMOVED***
		op.CommandMonitor.Succeeded(ctx, successEvent)
		return
	***REMOVED***

	failedEvent := &event.CommandFailedEvent***REMOVED***
		Failure:              info.cmdErr.Error(),
		CommandFinishedEvent: finished,
	***REMOVED***
	op.CommandMonitor.Failed(ctx, failedEvent)
***REMOVED***

// sessionsSupported returns true of the given server version indicates that it supports sessions.
func sessionsSupported(wireVersion *description.VersionRange) bool ***REMOVED***
	return wireVersion != nil && wireVersion.Max >= 6
***REMOVED***

// retryWritesSupported returns true if this description represents a server that supports retryable writes.
func retryWritesSupported(s description.Server) bool ***REMOVED***
	return s.SessionTimeoutMinutes != 0 && s.Kind != description.Standalone
***REMOVED***
