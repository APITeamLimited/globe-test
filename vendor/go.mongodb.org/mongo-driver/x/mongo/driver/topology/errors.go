// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/description"
)

// ConnectionError represents a connection error.
type ConnectionError struct ***REMOVED***
	ConnectionID string
	Wrapped      error

	// init will be set to true if this error occurred during connection initialization or
	// during a connection handshake.
	init    bool
	message string
***REMOVED***

// Error implements the error interface.
func (e ConnectionError) Error() string ***REMOVED***
	message := e.message
	if e.init ***REMOVED***
		fullMsg := "error occurred during connection handshake"
		if message != "" ***REMOVED***
			fullMsg = fmt.Sprintf("%s: %s", fullMsg, message)
		***REMOVED***
		message = fullMsg
	***REMOVED***
	if e.Wrapped != nil && message != "" ***REMOVED***
		return fmt.Sprintf("connection(%s) %s: %s", e.ConnectionID, message, e.Wrapped.Error())
	***REMOVED***
	if e.Wrapped != nil ***REMOVED***
		return fmt.Sprintf("connection(%s) %s", e.ConnectionID, e.Wrapped.Error())
	***REMOVED***
	return fmt.Sprintf("connection(%s) %s", e.ConnectionID, message)
***REMOVED***

// Unwrap returns the underlying error.
func (e ConnectionError) Unwrap() error ***REMOVED***
	return e.Wrapped
***REMOVED***

// ServerSelectionError represents a Server Selection error.
type ServerSelectionError struct ***REMOVED***
	Desc    description.Topology
	Wrapped error
***REMOVED***

// Error implements the error interface.
func (e ServerSelectionError) Error() string ***REMOVED***
	if e.Wrapped != nil ***REMOVED***
		return fmt.Sprintf("server selection error: %s, current topology: ***REMOVED*** %s ***REMOVED***", e.Wrapped.Error(), e.Desc.String())
	***REMOVED***
	return fmt.Sprintf("server selection error: current topology: ***REMOVED*** %s ***REMOVED***", e.Desc.String())
***REMOVED***

// Unwrap returns the underlying error.
func (e ServerSelectionError) Unwrap() error ***REMOVED***
	return e.Wrapped
***REMOVED***

// WaitQueueTimeoutError represents a timeout when requesting a connection from the pool
type WaitQueueTimeoutError struct ***REMOVED***
	Wrapped                      error
	PinnedCursorConnections      uint64
	PinnedTransactionConnections uint64
	maxPoolSize                  uint64
	totalConnectionCount         int
***REMOVED***

// Error implements the error interface.
func (w WaitQueueTimeoutError) Error() string ***REMOVED***
	errorMsg := "timed out while checking out a connection from connection pool"
	if w.Wrapped != nil ***REMOVED***
		errorMsg = fmt.Sprintf("%s: %s", errorMsg, w.Wrapped.Error())
	***REMOVED***

	return fmt.Sprintf(
		"%s; maxPoolSize: %d, connections in use by cursors: %d"+
			", connections in use by transactions: %d, connections in use by other operations: %d",
		errorMsg,
		w.maxPoolSize,
		w.PinnedCursorConnections,
		w.PinnedTransactionConnections,
		uint64(w.totalConnectionCount)-w.PinnedCursorConnections-w.PinnedTransactionConnections)
***REMOVED***

// Unwrap returns the underlying error.
func (w WaitQueueTimeoutError) Unwrap() error ***REMOVED***
	return w.Wrapped
***REMOVED***
