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
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// DropDatabase performs a dropDatabase operation
type DropDatabase struct ***REMOVED***
	session      *session.Client
	clock        *session.ClusterClock
	monitor      *event.CommandMonitor
	crypt        driver.Crypt
	database     string
	deployment   driver.Deployment
	selector     description.ServerSelector
	writeConcern *writeconcern.WriteConcern
	serverAPI    *driver.ServerAPIOptions
***REMOVED***

// NewDropDatabase constructs and returns a new DropDatabase.
func NewDropDatabase() *DropDatabase ***REMOVED***
	return &DropDatabase***REMOVED******REMOVED***
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (dd *DropDatabase) Execute(ctx context.Context) error ***REMOVED***
	if dd.deployment == nil ***REMOVED***
		return errors.New("the DropDatabase operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:      dd.command,
		Client:         dd.session,
		Clock:          dd.clock,
		CommandMonitor: dd.monitor,
		Crypt:          dd.crypt,
		Database:       dd.database,
		Deployment:     dd.deployment,
		Selector:       dd.selector,
		WriteConcern:   dd.writeConcern,
		ServerAPI:      dd.serverAPI,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (dd *DropDatabase) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***

	dst = bsoncore.AppendInt32Element(dst, "dropDatabase", 1)
	return dst, nil
***REMOVED***

// Session sets the session for this operation.
func (dd *DropDatabase) Session(session *session.Client) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.session = session
	return dd
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (dd *DropDatabase) ClusterClock(clock *session.ClusterClock) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.clock = clock
	return dd
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (dd *DropDatabase) CommandMonitor(monitor *event.CommandMonitor) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.monitor = monitor
	return dd
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (dd *DropDatabase) Crypt(crypt driver.Crypt) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.crypt = crypt
	return dd
***REMOVED***

// Database sets the database to run this operation against.
func (dd *DropDatabase) Database(database string) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.database = database
	return dd
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (dd *DropDatabase) Deployment(deployment driver.Deployment) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.deployment = deployment
	return dd
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (dd *DropDatabase) ServerSelector(selector description.ServerSelector) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.selector = selector
	return dd
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (dd *DropDatabase) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.writeConcern = writeConcern
	return dd
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (dd *DropDatabase) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropDatabase ***REMOVED***
	if dd == nil ***REMOVED***
		dd = new(DropDatabase)
	***REMOVED***

	dd.serverAPI = serverAPI
	return dd
***REMOVED***
