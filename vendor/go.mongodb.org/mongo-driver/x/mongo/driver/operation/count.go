// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Count represents a count operation.
type Count struct ***REMOVED***
	maxTimeMS      *int64
	query          bsoncore.Document
	session        *session.Client
	clock          *session.ClusterClock
	collection     string
	comment        bsoncore.Value
	monitor        *event.CommandMonitor
	crypt          driver.Crypt
	database       string
	deployment     driver.Deployment
	readConcern    *readconcern.ReadConcern
	readPreference *readpref.ReadPref
	selector       description.ServerSelector
	retry          *driver.RetryMode
	result         CountResult
	serverAPI      *driver.ServerAPIOptions
	timeout        *time.Duration
***REMOVED***

// CountResult represents a count result returned by the server.
type CountResult struct ***REMOVED***
	// The number of documents found
	N int64
***REMOVED***

func buildCountResult(response bsoncore.Document) (CountResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return CountResult***REMOVED******REMOVED***, err
	***REMOVED***
	cr := CountResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "n": // for count using original command
			var ok bool
			cr.N, ok = element.Value().AsInt64OK()
			if !ok ***REMOVED***
				return cr, fmt.Errorf("response field 'n' is type int64, but received BSON type %s",
					element.Value().Type)
			***REMOVED***
		case "cursor": // for count using aggregate with $collStats
			firstBatch, err := element.Value().Document().LookupErr("firstBatch")
			if err != nil ***REMOVED***
				return cr, err
			***REMOVED***

			// get count value from first batch
			val := firstBatch.Array().Index(0)
			count, err := val.Document().LookupErr("n")
			if err != nil ***REMOVED***
				return cr, err
			***REMOVED***

			// use count as Int64 for result
			var ok bool
			cr.N, ok = count.AsInt64OK()
			if !ok ***REMOVED***
				return cr, fmt.Errorf("response field 'n' is type int64, but received BSON type %s",
					element.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return cr, nil
***REMOVED***

// NewCount constructs and returns a new Count.
func NewCount() *Count ***REMOVED***
	return &Count***REMOVED******REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (c *Count) Result() CountResult ***REMOVED*** return c.result ***REMOVED***

func (c *Count) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error
	c.result, err = buildCountResult(info.ServerResponse)
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (c *Count) Execute(ctx context.Context) error ***REMOVED***
	if c.deployment == nil ***REMOVED***
		return errors.New("the Count operation must have a Deployment set before Execute can be called")
	***REMOVED***

	err := driver.Operation***REMOVED***
		CommandFn:         c.command,
		ProcessResponseFn: c.processResponse,
		RetryMode:         c.retry,
		Type:              driver.Read,
		Client:            c.session,
		Clock:             c.clock,
		CommandMonitor:    c.monitor,
		Crypt:             c.crypt,
		Database:          c.database,
		Deployment:        c.deployment,
		ReadConcern:       c.readConcern,
		ReadPreference:    c.readPreference,
		Selector:          c.selector,
		ServerAPI:         c.serverAPI,
		Timeout:           c.timeout,
	***REMOVED***.Execute(ctx, nil)

	// Swallow error if NamespaceNotFound(26) is returned from aggregate on non-existent namespace
	if err != nil ***REMOVED***
		dErr, ok := err.(driver.Error)
		if ok && dErr.Code == 26 ***REMOVED***
			err = nil
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (c *Count) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "count", c.collection)
	if c.query != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "query", c.query)
	***REMOVED***

	// Only append specified maxTimeMS if timeout is not also specified.
	if c.maxTimeMS != nil && c.timeout == nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *c.maxTimeMS)
	***REMOVED***
	if c.comment.Type != bsontype.Type(0) ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "comment", c.comment)
	***REMOVED***
	return dst, nil
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (c *Count) MaxTimeMS(maxTimeMS int64) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.maxTimeMS = &maxTimeMS
	return c
***REMOVED***

// Query determines what results are returned from find.
func (c *Count) Query(query bsoncore.Document) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.query = query
	return c
***REMOVED***

// Session sets the session for this operation.
func (c *Count) Session(session *session.Client) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.session = session
	return c
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (c *Count) ClusterClock(clock *session.ClusterClock) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.clock = clock
	return c
***REMOVED***

// Collection sets the collection that this command will run against.
func (c *Count) Collection(collection string) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.collection = collection
	return c
***REMOVED***

// Comment sets a value to help trace an operation.
func (c *Count) Comment(comment bsoncore.Value) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.comment = comment
	return c
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (c *Count) CommandMonitor(monitor *event.CommandMonitor) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.monitor = monitor
	return c
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (c *Count) Crypt(crypt driver.Crypt) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.crypt = crypt
	return c
***REMOVED***

// Database sets the database to run this operation against.
func (c *Count) Database(database string) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.database = database
	return c
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (c *Count) Deployment(deployment driver.Deployment) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.deployment = deployment
	return c
***REMOVED***

// ReadConcern specifies the read concern for this operation.
func (c *Count) ReadConcern(readConcern *readconcern.ReadConcern) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.readConcern = readConcern
	return c
***REMOVED***

// ReadPreference set the read preference used with this operation.
func (c *Count) ReadPreference(readPreference *readpref.ReadPref) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.readPreference = readPreference
	return c
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (c *Count) ServerSelector(selector description.ServerSelector) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.selector = selector
	return c
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (c *Count) Retry(retry driver.RetryMode) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.retry = &retry
	return c
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (c *Count) ServerAPI(serverAPI *driver.ServerAPIOptions) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.serverAPI = serverAPI
	return c
***REMOVED***

// Timeout sets the timeout for this operation.
func (c *Count) Timeout(timeout *time.Duration) *Count ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Count)
	***REMOVED***

	c.timeout = timeout
	return c
***REMOVED***
