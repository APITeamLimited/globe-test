// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// EndSessions performs an endSessions operation.
type EndSessions struct ***REMOVED***
	sessionIDs bsoncore.Document
	session    *session.Client
	clock      *session.ClusterClock
	monitor    *event.CommandMonitor
	crypt      driver.Crypt
	database   string
	deployment driver.Deployment
	selector   description.ServerSelector
	serverAPI  *driver.ServerAPIOptions
***REMOVED***

// NewEndSessions constructs and returns a new EndSessions.
func NewEndSessions(sessionIDs bsoncore.Document) *EndSessions ***REMOVED***
	return &EndSessions***REMOVED***
		sessionIDs: sessionIDs,
	***REMOVED***
***REMOVED***

func (es *EndSessions) processResponse(driver.ResponseInfo) error ***REMOVED***
	var err error
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (es *EndSessions) Execute(ctx context.Context) error ***REMOVED***
	if es.deployment == nil ***REMOVED***
		return errors.New("the EndSessions operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         es.command,
		ProcessResponseFn: es.processResponse,
		Client:            es.session,
		Clock:             es.clock,
		CommandMonitor:    es.monitor,
		Crypt:             es.crypt,
		Database:          es.database,
		Deployment:        es.deployment,
		Selector:          es.selector,
		ServerAPI:         es.serverAPI,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (es *EndSessions) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	if es.sessionIDs != nil ***REMOVED***
		dst = bsoncore.AppendArrayElement(dst, "endSessions", es.sessionIDs)
	***REMOVED***
	return dst, nil
***REMOVED***

// SessionIDs specifies the sessions to be expired.
func (es *EndSessions) SessionIDs(sessionIDs bsoncore.Document) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.sessionIDs = sessionIDs
	return es
***REMOVED***

// Session sets the session for this operation.
func (es *EndSessions) Session(session *session.Client) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.session = session
	return es
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (es *EndSessions) ClusterClock(clock *session.ClusterClock) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.clock = clock
	return es
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (es *EndSessions) CommandMonitor(monitor *event.CommandMonitor) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.monitor = monitor
	return es
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (es *EndSessions) Crypt(crypt driver.Crypt) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.crypt = crypt
	return es
***REMOVED***

// Database sets the database to run this operation against.
func (es *EndSessions) Database(database string) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.database = database
	return es
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (es *EndSessions) Deployment(deployment driver.Deployment) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.deployment = deployment
	return es
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (es *EndSessions) ServerSelector(selector description.ServerSelector) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.selector = selector
	return es
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (es *EndSessions) ServerAPI(serverAPI *driver.ServerAPIOptions) *EndSessions ***REMOVED***
	if es == nil ***REMOVED***
		es = new(EndSessions)
	***REMOVED***

	es.serverAPI = serverAPI
	return es
***REMOVED***
