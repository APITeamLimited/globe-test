// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session // import "go.mongodb.org/mongo-driver/x/mongo/driver/session"

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrSessionEnded is returned when a client session is used after a call to endSession().
var ErrSessionEnded = errors.New("ended session was used")

// ErrNoTransactStarted is returned if a transaction operation is called when no transaction has started.
var ErrNoTransactStarted = errors.New("no transaction started")

// ErrTransactInProgress is returned if startTransaction() is called when a transaction is in progress.
var ErrTransactInProgress = errors.New("transaction already in progress")

// ErrAbortAfterCommit is returned when abort is called after a commit.
var ErrAbortAfterCommit = errors.New("cannot call abortTransaction after calling commitTransaction")

// ErrAbortTwice is returned if abort is called after transaction is already aborted.
var ErrAbortTwice = errors.New("cannot call abortTransaction twice")

// ErrCommitAfterAbort is returned if commit is called after an abort.
var ErrCommitAfterAbort = errors.New("cannot call commitTransaction after calling abortTransaction")

// ErrUnackWCUnsupported is returned if an unacknowledged write concern is supported for a transaction.
var ErrUnackWCUnsupported = errors.New("transactions do not support unacknowledged write concerns")

// ErrSnapshotTransaction is returned if an transaction is started on a snapshot session.
var ErrSnapshotTransaction = errors.New("transactions are not supported in snapshot sessions")

// Type describes the type of the session
type Type uint8

// These constants are the valid types for a client session.
const (
	Explicit Type = iota
	Implicit
)

// TransactionState indicates the state of the transactions FSM.
type TransactionState uint8

// Client Session states
const (
	None TransactionState = iota
	Starting
	InProgress
	Committed
	Aborted
)

// String implements the fmt.Stringer interface.
func (s TransactionState) String() string ***REMOVED***
	switch s ***REMOVED***
	case None:
		return "none"
	case Starting:
		return "starting"
	case InProgress:
		return "in progress"
	case Committed:
		return "committed"
	case Aborted:
		return "aborted"
	default:
		return "unknown"
	***REMOVED***
***REMOVED***

// LoadBalancedTransactionConnection represents a connection that's pinned by a ClientSession because it's being used
// to execute a transaction when running against a load balancer. This interface is a copy of driver.PinnedConnection
// and exists to be able to pin transactions to a connection without causing an import cycle.
type LoadBalancedTransactionConnection interface ***REMOVED***
	// Functions copied over from driver.Connection.
	WriteWireMessage(context.Context, []byte) error
	ReadWireMessage(ctx context.Context, dst []byte) ([]byte, error)
	Description() description.Server
	Close() error
	ID() string
	ServerConnectionID() *int32
	Address() address.Address
	Stale() bool

	// Functions copied over from driver.PinnedConnection that are not part of Connection or Expirable.
	PinToCursor() error
	PinToTransaction() error
	UnpinFromCursor() error
	UnpinFromTransaction() error
***REMOVED***

// Client is a session for clients to run commands.
type Client struct ***REMOVED***
	*Server
	ClientID       uuid.UUID
	ClusterTime    bson.Raw
	Consistent     bool // causal consistency
	OperationTime  *primitive.Timestamp
	SessionType    Type
	Terminated     bool
	RetryingCommit bool
	Committing     bool
	Aborting       bool
	RetryWrite     bool
	RetryRead      bool
	Snapshot       bool

	// options for the current transaction
	// most recently set by transactionopt
	CurrentRc  *readconcern.ReadConcern
	CurrentRp  *readpref.ReadPref
	CurrentWc  *writeconcern.WriteConcern
	CurrentMct *time.Duration

	// default transaction options
	transactionRc            *readconcern.ReadConcern
	transactionRp            *readpref.ReadPref
	transactionWc            *writeconcern.WriteConcern
	transactionMaxCommitTime *time.Duration

	pool             *Pool
	TransactionState TransactionState
	PinnedServer     *description.Server
	RecoveryToken    bson.Raw
	PinnedConnection LoadBalancedTransactionConnection
	SnapshotTime     *primitive.Timestamp
***REMOVED***

func getClusterTime(clusterTime bson.Raw) (uint32, uint32) ***REMOVED***
	if clusterTime == nil ***REMOVED***
		return 0, 0
	***REMOVED***

	clusterTimeVal, err := clusterTime.LookupErr("$clusterTime")
	if err != nil ***REMOVED***
		return 0, 0
	***REMOVED***

	timestampVal, err := bson.Raw(clusterTimeVal.Value).LookupErr("clusterTime")
	if err != nil ***REMOVED***
		return 0, 0
	***REMOVED***

	return timestampVal.Timestamp()
***REMOVED***

// MaxClusterTime compares 2 clusterTime documents and returns the document representing the highest cluster time.
func MaxClusterTime(ct1, ct2 bson.Raw) bson.Raw ***REMOVED***
	epoch1, ord1 := getClusterTime(ct1)
	epoch2, ord2 := getClusterTime(ct2)

	if epoch1 > epoch2 ***REMOVED***
		return ct1
	***REMOVED*** else if epoch1 < epoch2 ***REMOVED***
		return ct2
	***REMOVED*** else if ord1 > ord2 ***REMOVED***
		return ct1
	***REMOVED*** else if ord1 < ord2 ***REMOVED***
		return ct2
	***REMOVED***

	return ct1
***REMOVED***

// NewClientSession creates a Client.
func NewClientSession(pool *Pool, clientID uuid.UUID, sessionType Type, opts ...*ClientOptions) (*Client, error) ***REMOVED***
	mergedOpts := mergeClientOptions(opts...)

	c := &Client***REMOVED***
		ClientID:    clientID,
		SessionType: sessionType,
		pool:        pool,
	***REMOVED***
	if mergedOpts.DefaultReadPreference != nil ***REMOVED***
		c.transactionRp = mergedOpts.DefaultReadPreference
	***REMOVED***
	if mergedOpts.DefaultReadConcern != nil ***REMOVED***
		c.transactionRc = mergedOpts.DefaultReadConcern
	***REMOVED***
	if mergedOpts.DefaultWriteConcern != nil ***REMOVED***
		c.transactionWc = mergedOpts.DefaultWriteConcern
	***REMOVED***
	if mergedOpts.DefaultMaxCommitTime != nil ***REMOVED***
		c.transactionMaxCommitTime = mergedOpts.DefaultMaxCommitTime
	***REMOVED***
	if mergedOpts.Snapshot != nil ***REMOVED***
		c.Snapshot = *mergedOpts.Snapshot
	***REMOVED***

	// The default for causalConsistency is true, unless Snapshot is enabled, then it's false. Set
	// the default and then allow any explicit causalConsistency setting to override it.
	c.Consistent = !c.Snapshot
	if mergedOpts.CausalConsistency != nil ***REMOVED***
		c.Consistent = *mergedOpts.CausalConsistency
	***REMOVED***

	if c.Consistent && c.Snapshot ***REMOVED***
		return nil, errors.New("causal consistency and snapshot cannot both be set for a session")
	***REMOVED***

	servSess, err := pool.GetSession()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.Server = servSess

	return c, nil
***REMOVED***

// AdvanceClusterTime updates the session's cluster time.
func (c *Client) AdvanceClusterTime(clusterTime bson.Raw) error ***REMOVED***
	if c.Terminated ***REMOVED***
		return ErrSessionEnded
	***REMOVED***
	c.ClusterTime = MaxClusterTime(c.ClusterTime, clusterTime)
	return nil
***REMOVED***

// AdvanceOperationTime updates the session's operation time.
func (c *Client) AdvanceOperationTime(opTime *primitive.Timestamp) error ***REMOVED***
	if c.Terminated ***REMOVED***
		return ErrSessionEnded
	***REMOVED***

	if c.OperationTime == nil ***REMOVED***
		c.OperationTime = opTime
		return nil
	***REMOVED***

	if opTime.T > c.OperationTime.T ***REMOVED***
		c.OperationTime = opTime
	***REMOVED*** else if (opTime.T == c.OperationTime.T) && (opTime.I > c.OperationTime.I) ***REMOVED***
		c.OperationTime = opTime
	***REMOVED***

	return nil
***REMOVED***

// UpdateUseTime sets the session's last used time to the current time. This must be called whenever the session is
// used to send a command to the server to ensure that the session is not prematurely marked expired in the driver's
// session pool. If the session has already been ended, this method will return ErrSessionEnded.
func (c *Client) UpdateUseTime() error ***REMOVED***
	if c.Terminated ***REMOVED***
		return ErrSessionEnded
	***REMOVED***
	c.updateUseTime()
	return nil
***REMOVED***

// UpdateRecoveryToken updates the session's recovery token from the server response.
func (c *Client) UpdateRecoveryToken(response bson.Raw) ***REMOVED***
	if c == nil ***REMOVED***
		return
	***REMOVED***

	token, err := response.LookupErr("recoveryToken")
	if err != nil ***REMOVED***
		return
	***REMOVED***

	c.RecoveryToken = token.Document()
***REMOVED***

// UpdateSnapshotTime updates the session's value for the atClusterTime field of ReadConcern.
func (c *Client) UpdateSnapshotTime(response bsoncore.Document) ***REMOVED***
	if c == nil ***REMOVED***
		return
	***REMOVED***

	subDoc := response
	if cur, ok := response.Lookup("cursor").DocumentOK(); ok ***REMOVED***
		subDoc = cur
	***REMOVED***

	ssTimeElem, err := subDoc.LookupErr("atClusterTime")
	if err != nil ***REMOVED***
		// atClusterTime not included by the server
		return
	***REMOVED***

	t, i := ssTimeElem.Timestamp()
	c.SnapshotTime = &primitive.Timestamp***REMOVED***
		T: t,
		I: i,
	***REMOVED***
***REMOVED***

// ClearPinnedResources clears the pinned server and/or connection associated with the session.
func (c *Client) ClearPinnedResources() error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***

	c.PinnedServer = nil
	if c.PinnedConnection != nil ***REMOVED***
		if err := c.PinnedConnection.UnpinFromTransaction(); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := c.PinnedConnection.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	c.PinnedConnection = nil
	return nil
***REMOVED***

// UnpinConnection gracefully unpins the connection associated with the session if there is one. This is done via
// the pinned connection's UnpinFromTransaction function.
func (c *Client) UnpinConnection() error ***REMOVED***
	if c == nil || c.PinnedConnection == nil ***REMOVED***
		return nil
	***REMOVED***

	err := c.PinnedConnection.UnpinFromTransaction()
	closeErr := c.PinnedConnection.Close()
	if err == nil && closeErr != nil ***REMOVED***
		err = closeErr
	***REMOVED***
	c.PinnedConnection = nil
	return err
***REMOVED***

// EndSession ends the session.
func (c *Client) EndSession() ***REMOVED***
	if c.Terminated ***REMOVED***
		return
	***REMOVED***

	c.Terminated = true
	c.pool.ReturnSession(c.Server)
***REMOVED***

// TransactionInProgress returns true if the client session is in an active transaction.
func (c *Client) TransactionInProgress() bool ***REMOVED***
	return c.TransactionState == InProgress
***REMOVED***

// TransactionStarting returns true if the client session is starting a transaction.
func (c *Client) TransactionStarting() bool ***REMOVED***
	return c.TransactionState == Starting
***REMOVED***

// TransactionRunning returns true if the client session has started the transaction
// and it hasn't been committed or aborted
func (c *Client) TransactionRunning() bool ***REMOVED***
	return c != nil && (c.TransactionState == Starting || c.TransactionState == InProgress)
***REMOVED***

// TransactionCommitted returns true of the client session just committed a transaction.
func (c *Client) TransactionCommitted() bool ***REMOVED***
	return c.TransactionState == Committed
***REMOVED***

// CheckStartTransaction checks to see if allowed to start transaction and returns
// an error if not allowed
func (c *Client) CheckStartTransaction() error ***REMOVED***
	if c.TransactionState == InProgress || c.TransactionState == Starting ***REMOVED***
		return ErrTransactInProgress
	***REMOVED***
	if c.Snapshot ***REMOVED***
		return ErrSnapshotTransaction
	***REMOVED***
	return nil
***REMOVED***

// StartTransaction initializes the transaction options and advances the state machine.
// It does not contact the server to start the transaction.
func (c *Client) StartTransaction(opts *TransactionOptions) error ***REMOVED***
	err := c.CheckStartTransaction()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.IncrementTxnNumber()
	c.RetryingCommit = false

	if opts != nil ***REMOVED***
		c.CurrentRc = opts.ReadConcern
		c.CurrentRp = opts.ReadPreference
		c.CurrentWc = opts.WriteConcern
		c.CurrentMct = opts.MaxCommitTime
	***REMOVED***

	if c.CurrentRc == nil ***REMOVED***
		c.CurrentRc = c.transactionRc
	***REMOVED***

	if c.CurrentRp == nil ***REMOVED***
		c.CurrentRp = c.transactionRp
	***REMOVED***

	if c.CurrentWc == nil ***REMOVED***
		c.CurrentWc = c.transactionWc
	***REMOVED***

	if c.CurrentMct == nil ***REMOVED***
		c.CurrentMct = c.transactionMaxCommitTime
	***REMOVED***

	if !writeconcern.AckWrite(c.CurrentWc) ***REMOVED***
		_ = c.clearTransactionOpts()
		return ErrUnackWCUnsupported
	***REMOVED***

	c.TransactionState = Starting
	return c.ClearPinnedResources()
***REMOVED***

// CheckCommitTransaction checks to see if allowed to commit transaction and returns
// an error if not allowed.
func (c *Client) CheckCommitTransaction() error ***REMOVED***
	if c.TransactionState == None ***REMOVED***
		return ErrNoTransactStarted
	***REMOVED*** else if c.TransactionState == Aborted ***REMOVED***
		return ErrCommitAfterAbort
	***REMOVED***
	return nil
***REMOVED***

// CommitTransaction updates the state for a successfully committed transaction and returns
// an error if not permissible.  It does not actually perform the commit.
func (c *Client) CommitTransaction() error ***REMOVED***
	err := c.CheckCommitTransaction()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.TransactionState = Committed
	return nil
***REMOVED***

// UpdateCommitTransactionWriteConcern will set the write concern to majority and potentially set  a
// w timeout of 10 seconds. This should be called after a commit transaction operation fails with a
// retryable error or after a successful commit transaction operation.
func (c *Client) UpdateCommitTransactionWriteConcern() ***REMOVED***
	wc := c.CurrentWc
	timeout := 10 * time.Second
	if wc != nil && wc.GetWTimeout() != 0 ***REMOVED***
		timeout = wc.GetWTimeout()
	***REMOVED***
	c.CurrentWc = wc.WithOptions(writeconcern.WMajority(), writeconcern.WTimeout(timeout))
***REMOVED***

// CheckAbortTransaction checks to see if allowed to abort transaction and returns
// an error if not allowed.
func (c *Client) CheckAbortTransaction() error ***REMOVED***
	if c.TransactionState == None ***REMOVED***
		return ErrNoTransactStarted
	***REMOVED*** else if c.TransactionState == Committed ***REMOVED***
		return ErrAbortAfterCommit
	***REMOVED*** else if c.TransactionState == Aborted ***REMOVED***
		return ErrAbortTwice
	***REMOVED***
	return nil
***REMOVED***

// AbortTransaction updates the state for a successfully aborted transaction and returns
// an error if not permissible.  It does not actually perform the abort.
func (c *Client) AbortTransaction() error ***REMOVED***
	err := c.CheckAbortTransaction()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.TransactionState = Aborted
	return c.clearTransactionOpts()
***REMOVED***

// StartCommand updates the session's internal state at the beginning of an operation. This must be called before
// server selection is done for the operation as the session's state can impact the result of that process.
func (c *Client) StartCommand() error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***

	// If we're executing the first operation using this session after a transaction, we must ensure that the session
	// is not pinned to any resources.
	if !c.TransactionRunning() && !c.Committing && !c.Aborting ***REMOVED***
		return c.ClearPinnedResources()
	***REMOVED***
	return nil
***REMOVED***

// ApplyCommand advances the state machine upon command execution. This must be called after server selection is
// complete.
func (c *Client) ApplyCommand(desc description.Server) error ***REMOVED***
	if c.Committing ***REMOVED***
		// Do not change state if committing after already committed
		return nil
	***REMOVED***
	if c.TransactionState == Starting ***REMOVED***
		c.TransactionState = InProgress
		// If this is in a transaction and the server is a mongos, pin it
		if desc.Kind == description.Mongos ***REMOVED***
			c.PinnedServer = &desc
		***REMOVED***
	***REMOVED*** else if c.TransactionState == Committed || c.TransactionState == Aborted ***REMOVED***
		c.TransactionState = None
		return c.clearTransactionOpts()
	***REMOVED***

	return nil
***REMOVED***

func (c *Client) clearTransactionOpts() error ***REMOVED***
	c.RetryingCommit = false
	c.Aborting = false
	c.Committing = false
	c.CurrentWc = nil
	c.CurrentRp = nil
	c.CurrentRc = nil
	c.RecoveryToken = nil

	return c.ClearPinnedResources()
***REMOVED***
