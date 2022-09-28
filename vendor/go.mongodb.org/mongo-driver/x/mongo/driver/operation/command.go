// Copyright (C) MongoDB, Inc. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Command is used to run a generic operation.
type Command struct ***REMOVED***
	command        bsoncore.Document
	readConcern    *readconcern.ReadConcern
	database       string
	deployment     driver.Deployment
	selector       description.ServerSelector
	readPreference *readpref.ReadPref
	clock          *session.ClusterClock
	session        *session.Client
	monitor        *event.CommandMonitor
	resultResponse bsoncore.Document
	resultCursor   *driver.BatchCursor
	crypt          driver.Crypt
	serverAPI      *driver.ServerAPIOptions
	createCursor   bool
	cursorOpts     driver.CursorOptions
	timeout        *time.Duration
***REMOVED***

// NewCommand constructs and returns a new Command. Once the operation is executed, the result may only be accessed via
// the Result() function.
func NewCommand(command bsoncore.Document) *Command ***REMOVED***
	return &Command***REMOVED***
		command: command,
	***REMOVED***
***REMOVED***

// NewCursorCommand constructs a new Command. Once the operation is executed, the server response will be used to
// construct a cursor, which can be accessed via the ResultCursor() function.
func NewCursorCommand(command bsoncore.Document, cursorOpts driver.CursorOptions) *Command ***REMOVED***
	return &Command***REMOVED***
		command:      command,
		cursorOpts:   cursorOpts,
		createCursor: true,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (c *Command) Result() bsoncore.Document ***REMOVED*** return c.resultResponse ***REMOVED***

// ResultCursor returns the BatchCursor that was constructed using the command response. If the operation was not
// configured to create a cursor (i.e. it was created using NewCommand rather than NewCursorCommand), this function
// will return nil and an error.
func (c *Command) ResultCursor() (*driver.BatchCursor, error) ***REMOVED***
	if !c.createCursor ***REMOVED***
		return nil, errors.New("command operation was not configured to create a cursor, but a result cursor was requested")
	***REMOVED***
	return c.resultCursor, nil
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (c *Command) Execute(ctx context.Context) error ***REMOVED***
	if c.deployment == nil ***REMOVED***
		return errors.New("the Command operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn: func(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
			return append(dst, c.command[4:len(c.command)-1]...), nil
		***REMOVED***,
		ProcessResponseFn: func(info driver.ResponseInfo) error ***REMOVED***
			c.resultResponse = info.ServerResponse

			if c.createCursor ***REMOVED***
				cursorRes, err := driver.NewCursorResponse(info)
				if err != nil ***REMOVED***
					return err
				***REMOVED***

				c.resultCursor, err = driver.NewBatchCursor(cursorRes, c.session, c.clock, c.cursorOpts)
				return err
			***REMOVED***

			return nil
		***REMOVED***,
		Client:         c.session,
		Clock:          c.clock,
		CommandMonitor: c.monitor,
		Database:       c.database,
		Deployment:     c.deployment,
		ReadPreference: c.readPreference,
		Selector:       c.selector,
		Crypt:          c.crypt,
		ServerAPI:      c.serverAPI,
		Timeout:        c.timeout,
	***REMOVED***.Execute(ctx, nil)
***REMOVED***

// Session sets the session for this operation.
func (c *Command) Session(session *session.Client) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.session = session
	return c
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (c *Command) ClusterClock(clock *session.ClusterClock) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.clock = clock
	return c
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (c *Command) CommandMonitor(monitor *event.CommandMonitor) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.monitor = monitor
	return c
***REMOVED***

// Database sets the database to run this operation against.
func (c *Command) Database(database string) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.database = database
	return c
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (c *Command) Deployment(deployment driver.Deployment) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.deployment = deployment
	return c
***REMOVED***

// ReadConcern specifies the read concern for this operation.
func (c *Command) ReadConcern(readConcern *readconcern.ReadConcern) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.readConcern = readConcern
	return c
***REMOVED***

// ReadPreference set the read preference used with this operation.
func (c *Command) ReadPreference(readPreference *readpref.ReadPref) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.readPreference = readPreference
	return c
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (c *Command) ServerSelector(selector description.ServerSelector) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.selector = selector
	return c
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (c *Command) Crypt(crypt driver.Crypt) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.crypt = crypt
	return c
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (c *Command) ServerAPI(serverAPI *driver.ServerAPIOptions) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.serverAPI = serverAPI
	return c
***REMOVED***

// Timeout sets the timeout for this operation.
func (c *Command) Timeout(timeout *time.Duration) *Command ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Command)
	***REMOVED***

	c.timeout = timeout
	return c
***REMOVED***
