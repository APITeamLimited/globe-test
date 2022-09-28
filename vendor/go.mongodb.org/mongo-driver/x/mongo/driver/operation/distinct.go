// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
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

// Distinct performs a distinct operation.
type Distinct struct ***REMOVED***
	collation      bsoncore.Document
	key            *string
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
	result         DistinctResult
	serverAPI      *driver.ServerAPIOptions
	timeout        *time.Duration
***REMOVED***

// DistinctResult represents a distinct result returned by the server.
type DistinctResult struct ***REMOVED***
	// The distinct values for the field.
	Values bsoncore.Value
***REMOVED***

func buildDistinctResult(response bsoncore.Document) (DistinctResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return DistinctResult***REMOVED******REMOVED***, err
	***REMOVED***
	dr := DistinctResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "values":
			dr.Values = element.Value()
		***REMOVED***
	***REMOVED***
	return dr, nil
***REMOVED***

// NewDistinct constructs and returns a new Distinct.
func NewDistinct(key string, query bsoncore.Document) *Distinct ***REMOVED***
	return &Distinct***REMOVED***
		key:   &key,
		query: query,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (d *Distinct) Result() DistinctResult ***REMOVED*** return d.result ***REMOVED***

func (d *Distinct) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error
	d.result, err = buildDistinctResult(info.ServerResponse)
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (d *Distinct) Execute(ctx context.Context) error ***REMOVED***
	if d.deployment == nil ***REMOVED***
		return errors.New("the Distinct operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         d.command,
		ProcessResponseFn: d.processResponse,
		RetryMode:         d.retry,
		Type:              driver.Read,
		Client:            d.session,
		Clock:             d.clock,
		CommandMonitor:    d.monitor,
		Crypt:             d.crypt,
		Database:          d.database,
		Deployment:        d.deployment,
		ReadConcern:       d.readConcern,
		ReadPreference:    d.readPreference,
		Selector:          d.selector,
		ServerAPI:         d.serverAPI,
		Timeout:           d.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (d *Distinct) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "distinct", d.collection)
	if d.collation != nil ***REMOVED***
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) ***REMOVED***
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "collation", d.collation)
	***REMOVED***
	if d.comment.Type != bsontype.Type(0) ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "comment", d.comment)
	***REMOVED***
	if d.key != nil ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "key", *d.key)
	***REMOVED***
	if d.maxTimeMS != nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", *d.maxTimeMS)
	***REMOVED***
	if d.query != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "query", d.query)
	***REMOVED***
	return dst, nil
***REMOVED***

// Collation specifies a collation to be used.
func (d *Distinct) Collation(collation bsoncore.Document) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.collation = collation
	return d
***REMOVED***

// Key specifies which field to return distinct values for.
func (d *Distinct) Key(key string) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.key = &key
	return d
***REMOVED***

// MaxTimeMS specifies the maximum amount of time to allow the query to run.
func (d *Distinct) MaxTimeMS(maxTimeMS int64) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.maxTimeMS = &maxTimeMS
	return d
***REMOVED***

// Query specifies which documents to return distinct values from.
func (d *Distinct) Query(query bsoncore.Document) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.query = query
	return d
***REMOVED***

// Session sets the session for this operation.
func (d *Distinct) Session(session *session.Client) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.session = session
	return d
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (d *Distinct) ClusterClock(clock *session.ClusterClock) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.clock = clock
	return d
***REMOVED***

// Collection sets the collection that this command will run against.
func (d *Distinct) Collection(collection string) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.collection = collection
	return d
***REMOVED***

// Comment sets a value to help trace an operation.
func (d *Distinct) Comment(comment bsoncore.Value) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.comment = comment
	return d
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (d *Distinct) CommandMonitor(monitor *event.CommandMonitor) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.monitor = monitor
	return d
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (d *Distinct) Crypt(crypt driver.Crypt) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.crypt = crypt
	return d
***REMOVED***

// Database sets the database to run this operation against.
func (d *Distinct) Database(database string) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.database = database
	return d
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (d *Distinct) Deployment(deployment driver.Deployment) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.deployment = deployment
	return d
***REMOVED***

// ReadConcern specifies the read concern for this operation.
func (d *Distinct) ReadConcern(readConcern *readconcern.ReadConcern) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.readConcern = readConcern
	return d
***REMOVED***

// ReadPreference set the read preference used with this operation.
func (d *Distinct) ReadPreference(readPreference *readpref.ReadPref) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.readPreference = readPreference
	return d
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (d *Distinct) ServerSelector(selector description.ServerSelector) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.selector = selector
	return d
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (d *Distinct) Retry(retry driver.RetryMode) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.retry = &retry
	return d
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (d *Distinct) ServerAPI(serverAPI *driver.ServerAPIOptions) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.serverAPI = serverAPI
	return d
***REMOVED***

// Timeout sets the timeout for this operation.
func (d *Distinct) Timeout(timeout *time.Duration) *Distinct ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Distinct)
	***REMOVED***

	d.timeout = timeout
	return d
***REMOVED***
