// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"time"

	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Server is an open session with the server.
type Server struct ***REMOVED***
	SessionID bsoncore.Document
	TxnNumber int64
	LastUsed  time.Time
	Dirty     bool
***REMOVED***

// returns whether or not a session has expired given a timeout in minutes
// a session is considered expired if it has less than 1 minute left before becoming stale
func (ss *Server) expired(topoDesc topologyDescription) bool ***REMOVED***
	// There is no server monitoring in LB mode, so we do not track session timeout minutes from server hello responses
	// and never consider sessions to be expired.
	if topoDesc.kind == description.LoadBalanced ***REMOVED***
		return false
	***REMOVED***

	if topoDesc.timeoutMinutes <= 0 ***REMOVED***
		return true
	***REMOVED***
	timeUnused := time.Since(ss.LastUsed).Minutes()
	return timeUnused > float64(topoDesc.timeoutMinutes-1)
***REMOVED***

// update the last used time for this session.
// must be called whenever this server session is used to send a command to the server.
func (ss *Server) updateUseTime() ***REMOVED***
	ss.LastUsed = time.Now()
***REMOVED***

func newServerSession() (*Server, error) ***REMOVED***
	id, err := uuid.New()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	idx, idDoc := bsoncore.AppendDocumentStart(nil)
	idDoc = bsoncore.AppendBinaryElement(idDoc, "id", UUIDSubtype, id[:])
	idDoc, _ = bsoncore.AppendDocumentEnd(idDoc, idx)

	return &Server***REMOVED***
		SessionID: idDoc,
		LastUsed:  time.Now(),
	***REMOVED***, nil
***REMOVED***

// IncrementTxnNumber increments the transaction number.
func (ss *Server) IncrementTxnNumber() ***REMOVED***
	ss.TxnNumber++
***REMOVED***

// MarkDirty marks the session as dirty.
func (ss *Server) MarkDirty() ***REMOVED***
	ss.Dirty = true
***REMOVED***

// UUIDSubtype is the BSON binary subtype that a UUID should be encoded as
const UUIDSubtype byte = 4
