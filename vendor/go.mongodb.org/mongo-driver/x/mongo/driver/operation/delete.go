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
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Delete performs a delete operation
type Delete struct ***REMOVED***
	comment      bsoncore.Value
	deletes      []bsoncore.Document
	ordered      *bool
	session      *session.Client
	clock        *session.ClusterClock
	collection   string
	monitor      *event.CommandMonitor
	crypt        driver.Crypt
	database     string
	deployment   driver.Deployment
	selector     description.ServerSelector
	writeConcern *writeconcern.WriteConcern
	retry        *driver.RetryMode
	hint         *bool
	result       DeleteResult
	serverAPI    *driver.ServerAPIOptions
	let          bsoncore.Document
	timeout      *time.Duration
***REMOVED***

// DeleteResult represents a delete result returned by the server.
type DeleteResult struct ***REMOVED***
	// Number of documents successfully deleted.
	N int64
***REMOVED***

func buildDeleteResult(response bsoncore.Document) (DeleteResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return DeleteResult***REMOVED******REMOVED***, err
	***REMOVED***
	dr := DeleteResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "n":
			var ok bool
			dr.N, ok = element.Value().AsInt64OK()
			if !ok ***REMOVED***
				return dr, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dr, nil
***REMOVED***

// NewDelete constructs and returns a new Delete.
func NewDelete(deletes ...bsoncore.Document) *Delete ***REMOVED***
	return &Delete***REMOVED***
		deletes: deletes,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (d *Delete) Result() DeleteResult ***REMOVED*** return d.result ***REMOVED***

func (d *Delete) processResponse(info driver.ResponseInfo) error ***REMOVED***
	dr, err := buildDeleteResult(info.ServerResponse)
	d.result.N += dr.N
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (d *Delete) Execute(ctx context.Context) error ***REMOVED***
	if d.deployment == nil ***REMOVED***
		return errors.New("the Delete operation must have a Deployment set before Execute can be called")
	***REMOVED***
	batches := &driver.Batches***REMOVED***
		Identifier: "deletes",
		Documents:  d.deletes,
		Ordered:    d.ordered,
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         d.command,
		ProcessResponseFn: d.processResponse,
		Batches:           batches,
		RetryMode:         d.retry,
		Type:              driver.Write,
		Client:            d.session,
		Clock:             d.clock,
		CommandMonitor:    d.monitor,
		Crypt:             d.crypt,
		Database:          d.database,
		Deployment:        d.deployment,
		Selector:          d.selector,
		WriteConcern:      d.writeConcern,
		ServerAPI:         d.serverAPI,
		Timeout:           d.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (d *Delete) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "delete", d.collection)
	if d.comment.Type != bsontype.Type(0) ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "comment", d.comment)
	***REMOVED***
	if d.ordered != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "ordered", *d.ordered)
	***REMOVED***
	if d.hint != nil && *d.hint ***REMOVED***
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) ***REMOVED***
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 5")
		***REMOVED***
		if !d.writeConcern.Acknowledged() ***REMOVED***
			return nil, errUnacknowledgedHint
		***REMOVED***
	***REMOVED***
	if d.let != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "let", d.let)
	***REMOVED***
	return dst, nil
***REMOVED***

// Deletes adds documents to this operation that will be used to determine what documents to delete when this operation
// is executed. These documents should have the form ***REMOVED***q: <query>, limit: <integer limit>, collation: <document>***REMOVED***. The
// collation field is optional. If limit is 0, there will be no limit on the number of documents deleted.
func (d *Delete) Deletes(deletes ...bsoncore.Document) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.deletes = deletes
	return d
***REMOVED***

// Ordered sets ordered. If true, when a write fails, the operation will return the error, when
// false write failures do not stop execution of the operation.
func (d *Delete) Ordered(ordered bool) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.ordered = &ordered
	return d
***REMOVED***

// Session sets the session for this operation.
func (d *Delete) Session(session *session.Client) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.session = session
	return d
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (d *Delete) ClusterClock(clock *session.ClusterClock) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.clock = clock
	return d
***REMOVED***

// Collection sets the collection that this command will run against.
func (d *Delete) Collection(collection string) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.collection = collection
	return d
***REMOVED***

// Comment sets a value to help trace an operation.
func (d *Delete) Comment(comment bsoncore.Value) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.comment = comment
	return d
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (d *Delete) CommandMonitor(monitor *event.CommandMonitor) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.monitor = monitor
	return d
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (d *Delete) Crypt(crypt driver.Crypt) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.crypt = crypt
	return d
***REMOVED***

// Database sets the database to run this operation against.
func (d *Delete) Database(database string) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.database = database
	return d
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (d *Delete) Deployment(deployment driver.Deployment) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.deployment = deployment
	return d
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (d *Delete) ServerSelector(selector description.ServerSelector) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.selector = selector
	return d
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (d *Delete) WriteConcern(writeConcern *writeconcern.WriteConcern) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.writeConcern = writeConcern
	return d
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (d *Delete) Retry(retry driver.RetryMode) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.retry = &retry
	return d
***REMOVED***

// Hint is a flag to indicate that the update document contains a hint. Hint is only supported by
// servers >= 4.4. Older servers >= 3.4 will report an error for using the hint option. For servers <
// 3.4, the driver will return an error if the hint option is used.
func (d *Delete) Hint(hint bool) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.hint = &hint
	return d
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (d *Delete) ServerAPI(serverAPI *driver.ServerAPIOptions) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.serverAPI = serverAPI
	return d
***REMOVED***

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (d *Delete) Let(let bsoncore.Document) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.let = let
	return d
***REMOVED***

// Timeout sets the timeout for this operation.
func (d *Delete) Timeout(timeout *time.Duration) *Delete ***REMOVED***
	if d == nil ***REMOVED***
		d = new(Delete)
	***REMOVED***

	d.timeout = timeout
	return d
***REMOVED***
