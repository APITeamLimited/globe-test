// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// ErrWrongClient is returned when a user attempts to pass in a session created by a different client than
// the method call is using.
var ErrWrongClient = errors.New("session was not created by this client")

var withTransactionTimeout = 120 * time.Second

// SessionContext combines the context.Context and mongo.Session interfaces. It should be used as the Context arguments
// to operations that should be executed in a session.
//
// Implementations of SessionContext are not safe for concurrent use by multiple goroutines.
//
// There are two ways to create a SessionContext and use it in a session/transaction. The first is to use one of the
// callback-based functions such as WithSession and UseSession. These functions create a SessionContext and pass it to
// the provided callback. The other is to use NewSessionContext to explicitly create a SessionContext.
type SessionContext interface ***REMOVED***
	context.Context
	Session
***REMOVED***

type sessionContext struct ***REMOVED***
	context.Context
	Session
***REMOVED***

type sessionKey struct ***REMOVED***
***REMOVED***

// NewSessionContext creates a new SessionContext associated with the given Context and Session parameters.
func NewSessionContext(ctx context.Context, sess Session) SessionContext ***REMOVED***
	return &sessionContext***REMOVED***
		Context: context.WithValue(ctx, sessionKey***REMOVED******REMOVED***, sess),
		Session: sess,
	***REMOVED***
***REMOVED***

// SessionFromContext extracts the mongo.Session object stored in a Context. This can be used on a SessionContext that
// was created implicitly through one of the callback-based session APIs or explicitly by calling NewSessionContext. If
// there is no Session stored in the provided Context, nil is returned.
func SessionFromContext(ctx context.Context) Session ***REMOVED***
	val := ctx.Value(sessionKey***REMOVED******REMOVED***)
	if val == nil ***REMOVED***
		return nil
	***REMOVED***

	sess, ok := val.(Session)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	return sess
***REMOVED***

// Session is an interface that represents a MongoDB logical session. Sessions can be used to enable causal consistency
// for a group of operations or to execute operations in an ACID transaction. A new Session can be created from a Client
// instance. A Session created from a Client must only be used to execute operations using that Client or a Database or
// Collection created from that Client. Custom implementations of this interface should not be used in production. For
// more information about sessions, and their use cases, see
// https://www.mongodb.com/docs/manual/reference/server-sessions/,
// https://www.mongodb.com/docs/manual/core/read-isolation-consistency-recency/#causal-consistency, and
// https://www.mongodb.com/docs/manual/core/transactions/.
//
// Implementations of Session are not safe for concurrent use by multiple goroutines.
//
// StartTransaction starts a new transaction, configured with the given options, on this session. This method will
// return an error if there is already a transaction in-progress for this session.
//
// CommitTransaction commits the active transaction for this session. This method will return an error if there is no
// active transaction for this session or the transaction has been aborted.
//
// AbortTransaction aborts the active transaction for this session. This method will return an error if there is no
// active transaction for this session or the transaction has been committed or aborted.
//
// WithTransaction starts a transaction on this session and runs the fn callback. Errors with the
// TransientTransactionError and UnknownTransactionCommitResult labels are retried for up to 120 seconds. Inside the
// callback, sessCtx must be used as the Context parameter for any operations that should be part of the transaction. If
// the ctx parameter already has a Session attached to it, it will be replaced by this session. The fn callback may be
// run multiple times during WithTransaction due to retry attempts, so it must be idempotent. Non-retryable operation
// errors or any operation errors that occur after the timeout expires will be returned without retrying. If the
// callback fails, the driver will call AbortTransaction. Because this method must succeed to ensure that server-side
// resources are properly cleaned up, context deadlines and cancellations will not be respected during this call. For a
// usage example, see the Client.StartSession method documentation.
//
// ClusterTime, OperationTime, Client, and ID return the session's current cluster time, the session's current operation
// time, the Client associated with the session, and the ID document associated with the session, respectively. The ID
// document for a session is in the form ***REMOVED***"id": <BSON binary value>***REMOVED***.
//
// EndSession method should abort any existing transactions and close the session.
//
// AdvanceClusterTime advances the cluster time for a session. This method will return an error if the session has ended.
//
// AdvanceOperationTime advances the operation time for a session. This method will return an error if the session has
// ended.
type Session interface ***REMOVED***
	// Functions to modify session state.
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
	WithTransaction(ctx context.Context, fn func(sessCtx SessionContext) (interface***REMOVED******REMOVED***, error),
		opts ...*options.TransactionOptions) (interface***REMOVED******REMOVED***, error)
	EndSession(context.Context)

	// Functions to retrieve session properties.
	ClusterTime() bson.Raw
	OperationTime() *primitive.Timestamp
	Client() *Client
	ID() bson.Raw

	// Functions to modify mutable session properties.
	AdvanceClusterTime(bson.Raw) error
	AdvanceOperationTime(*primitive.Timestamp) error

	session()
***REMOVED***

// XSession is an unstable interface for internal use only.
//
// Deprecated: This interface is unstable because it provides access to a session.Client object, which exists in the
// "x" package. It should not be used by applications and may be changed or removed in any release.
type XSession interface ***REMOVED***
	ClientSession() *session.Client
***REMOVED***

// sessionImpl represents a set of sequential operations executed by an application that are related in some way.
type sessionImpl struct ***REMOVED***
	clientSession       *session.Client
	client              *Client
	deployment          driver.Deployment
	didCommitAfterStart bool // true if commit was called after start with no other operations
***REMOVED***

var _ Session = &sessionImpl***REMOVED******REMOVED***
var _ XSession = &sessionImpl***REMOVED******REMOVED***

// ClientSession implements the XSession interface.
func (s *sessionImpl) ClientSession() *session.Client ***REMOVED***
	return s.clientSession
***REMOVED***

// ID implements the Session interface.
func (s *sessionImpl) ID() bson.Raw ***REMOVED***
	return bson.Raw(s.clientSession.SessionID)
***REMOVED***

// EndSession implements the Session interface.
func (s *sessionImpl) EndSession(ctx context.Context) ***REMOVED***
	if s.clientSession.TransactionInProgress() ***REMOVED***
		// ignore all errors aborting during an end session
		_ = s.AbortTransaction(ctx)
	***REMOVED***
	s.clientSession.EndSession()
***REMOVED***

// WithTransaction implements the Session interface.
func (s *sessionImpl) WithTransaction(ctx context.Context, fn func(sessCtx SessionContext) (interface***REMOVED******REMOVED***, error),
	opts ...*options.TransactionOptions) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	timeout := time.NewTimer(withTransactionTimeout)
	defer timeout.Stop()
	var err error
	for ***REMOVED***
		err = s.StartTransaction(opts...)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		res, err := fn(NewSessionContext(ctx, s))
		if err != nil ***REMOVED***
			if s.clientSession.TransactionRunning() ***REMOVED***
				// Wrap the user-provided Context in a new one that behaves like context.Background() for deadlines and
				// cancellations, but forwards Value requests to the original one.
				_ = s.AbortTransaction(internal.NewBackgroundContext(ctx))
			***REMOVED***

			select ***REMOVED***
			case <-timeout.C:
				return nil, err
			default:
			***REMOVED***

			if errorHasLabel(err, driver.TransientTransactionError) ***REMOVED***
				continue
			***REMOVED***
			return res, err
		***REMOVED***

		err = s.clientSession.CheckAbortTransaction()
		if err != nil ***REMOVED***
			return res, nil
		***REMOVED***

	CommitLoop:
		for ***REMOVED***
			err = s.CommitTransaction(ctx)
			// End when error is nil, as transaction has been committed.
			if err == nil ***REMOVED***
				return res, nil
			***REMOVED***

			select ***REMOVED***
			case <-timeout.C:
				return res, err
			default:
			***REMOVED***

			if cerr, ok := err.(CommandError); ok ***REMOVED***
				if cerr.HasErrorLabel(driver.UnknownTransactionCommitResult) && !cerr.IsMaxTimeMSExpiredError() ***REMOVED***
					continue
				***REMOVED***
				if cerr.HasErrorLabel(driver.TransientTransactionError) ***REMOVED***
					break CommitLoop
				***REMOVED***
			***REMOVED***
			return res, err
		***REMOVED***
	***REMOVED***
***REMOVED***

// StartTransaction implements the Session interface.
func (s *sessionImpl) StartTransaction(opts ...*options.TransactionOptions) error ***REMOVED***
	err := s.clientSession.CheckStartTransaction()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	s.didCommitAfterStart = false

	topts := options.MergeTransactionOptions(opts...)
	coreOpts := &session.TransactionOptions***REMOVED***
		ReadConcern:    topts.ReadConcern,
		ReadPreference: topts.ReadPreference,
		WriteConcern:   topts.WriteConcern,
		MaxCommitTime:  topts.MaxCommitTime,
	***REMOVED***

	return s.clientSession.StartTransaction(coreOpts)
***REMOVED***

// AbortTransaction implements the Session interface.
func (s *sessionImpl) AbortTransaction(ctx context.Context) error ***REMOVED***
	err := s.clientSession.CheckAbortTransaction()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Do not run the abort command if the transaction is in starting state
	if s.clientSession.TransactionStarting() || s.didCommitAfterStart ***REMOVED***
		return s.clientSession.AbortTransaction()
	***REMOVED***

	selector := makePinnedSelector(s.clientSession, description.WriteSelector())

	s.clientSession.Aborting = true
	_ = operation.NewAbortTransaction().Session(s.clientSession).ClusterClock(s.client.clock).Database("admin").
		Deployment(s.deployment).WriteConcern(s.clientSession.CurrentWc).ServerSelector(selector).
		Retry(driver.RetryOncePerCommand).CommandMonitor(s.client.monitor).
		RecoveryToken(bsoncore.Document(s.clientSession.RecoveryToken)).ServerAPI(s.client.serverAPI).Execute(ctx)

	s.clientSession.Aborting = false
	_ = s.clientSession.AbortTransaction()

	return nil
***REMOVED***

// CommitTransaction implements the Session interface.
func (s *sessionImpl) CommitTransaction(ctx context.Context) error ***REMOVED***
	err := s.clientSession.CheckCommitTransaction()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Do not run the commit command if the transaction is in started state
	if s.clientSession.TransactionStarting() || s.didCommitAfterStart ***REMOVED***
		s.didCommitAfterStart = true
		return s.clientSession.CommitTransaction()
	***REMOVED***

	if s.clientSession.TransactionCommitted() ***REMOVED***
		s.clientSession.RetryingCommit = true
	***REMOVED***

	selector := makePinnedSelector(s.clientSession, description.WriteSelector())

	s.clientSession.Committing = true
	op := operation.NewCommitTransaction().
		Session(s.clientSession).ClusterClock(s.client.clock).Database("admin").Deployment(s.deployment).
		WriteConcern(s.clientSession.CurrentWc).ServerSelector(selector).Retry(driver.RetryOncePerCommand).
		CommandMonitor(s.client.monitor).RecoveryToken(bsoncore.Document(s.clientSession.RecoveryToken)).
		ServerAPI(s.client.serverAPI)
	if s.clientSession.CurrentMct != nil ***REMOVED***
		op.MaxTimeMS(int64(*s.clientSession.CurrentMct / time.Millisecond))
	***REMOVED***

	err = op.Execute(ctx)
	// Return error without updating transaction state if it is a timeout, as the transaction has not
	// actually been committed.
	if IsTimeout(err) ***REMOVED***
		return replaceErrors(err)
	***REMOVED***
	s.clientSession.Committing = false
	commitErr := s.clientSession.CommitTransaction()

	// We set the write concern to majority for subsequent calls to CommitTransaction.
	s.clientSession.UpdateCommitTransactionWriteConcern()

	if err != nil ***REMOVED***
		return replaceErrors(err)
	***REMOVED***
	return commitErr
***REMOVED***

// ClusterTime implements the Session interface.
func (s *sessionImpl) ClusterTime() bson.Raw ***REMOVED***
	return s.clientSession.ClusterTime
***REMOVED***

// AdvanceClusterTime implements the Session interface.
func (s *sessionImpl) AdvanceClusterTime(d bson.Raw) error ***REMOVED***
	return s.clientSession.AdvanceClusterTime(d)
***REMOVED***

// OperationTime implements the Session interface.
func (s *sessionImpl) OperationTime() *primitive.Timestamp ***REMOVED***
	return s.clientSession.OperationTime
***REMOVED***

// AdvanceOperationTime implements the Session interface.
func (s *sessionImpl) AdvanceOperationTime(ts *primitive.Timestamp) error ***REMOVED***
	return s.clientSession.AdvanceOperationTime(ts)
***REMOVED***

// Client implements the Session interface.
func (s *sessionImpl) Client() *Client ***REMOVED***
	return s.client
***REMOVED***

// session implements the Session interface.
func (*sessionImpl) session() ***REMOVED***
***REMOVED***

// sessionFromContext checks for a sessionImpl in the argued context and returns the session if it
// exists
func sessionFromContext(ctx context.Context) *session.Client ***REMOVED***
	s := ctx.Value(sessionKey***REMOVED******REMOVED***)
	if ses, ok := s.(*sessionImpl); ses != nil && ok ***REMOVED***
		return ses.clientSession
	***REMOVED***

	return nil
***REMOVED***
