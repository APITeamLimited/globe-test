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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Update performs an update operation.
type Update struct ***REMOVED***
	bypassDocumentValidation *bool
	comment                  bsoncore.Value
	ordered                  *bool
	updates                  []bsoncore.Document
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	hint                     *bool
	arrayFilters             *bool
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	result                   UpdateResult
	crypt                    driver.Crypt
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	timeout                  *time.Duration
***REMOVED***

// Upsert contains the information for an upsert in an Update operation.
type Upsert struct ***REMOVED***
	Index int64
	ID    interface***REMOVED******REMOVED*** `bson:"_id"`
***REMOVED***

// UpdateResult contains information for the result of an Update operation.
type UpdateResult struct ***REMOVED***
	// Number of documents matched.
	N int64
	// Number of documents modified.
	NModified int64
	// Information about upserted documents.
	Upserted []Upsert
***REMOVED***

func buildUpdateResult(response bsoncore.Document) (UpdateResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return UpdateResult***REMOVED******REMOVED***, err
	***REMOVED***
	ur := UpdateResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "nModified":
			var ok bool
			ur.NModified, ok = element.Value().AsInt64OK()
			if !ok ***REMOVED***
				return ur, fmt.Errorf("response field 'nModified' is type int32 or int64, but received BSON type %s", element.Value().Type)
			***REMOVED***
		case "n":
			var ok bool
			ur.N, ok = element.Value().AsInt64OK()
			if !ok ***REMOVED***
				return ur, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			***REMOVED***
		case "upserted":
			arr, ok := element.Value().ArrayOK()
			if !ok ***REMOVED***
				return ur, fmt.Errorf("response field 'upserted' is type array, but received BSON type %s", element.Value().Type)
			***REMOVED***

			var values []bsoncore.Value
			values, err = arr.Values()
			if err != nil ***REMOVED***
				break
			***REMOVED***

			for _, val := range values ***REMOVED***
				valDoc, ok := val.DocumentOK()
				if !ok ***REMOVED***
					return ur, fmt.Errorf("upserted value is type document, but received BSON type %s", val.Type)
				***REMOVED***
				var upsert Upsert
				if err = bson.Unmarshal(valDoc, &upsert); err != nil ***REMOVED***
					return ur, err
				***REMOVED***
				ur.Upserted = append(ur.Upserted, upsert)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ur, nil
***REMOVED***

// NewUpdate constructs and returns a new Update.
func NewUpdate(updates ...bsoncore.Document) *Update ***REMOVED***
	return &Update***REMOVED***
		updates: updates,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (u *Update) Result() UpdateResult ***REMOVED*** return u.result ***REMOVED***

func (u *Update) processResponse(info driver.ResponseInfo) error ***REMOVED***
	ur, err := buildUpdateResult(info.ServerResponse)

	u.result.N += ur.N
	u.result.NModified += ur.NModified
	if info.CurrentIndex > 0 ***REMOVED***
		for ind := range ur.Upserted ***REMOVED***
			ur.Upserted[ind].Index += int64(info.CurrentIndex)
		***REMOVED***
	***REMOVED***
	u.result.Upserted = append(u.result.Upserted, ur.Upserted...)
	return err

***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (u *Update) Execute(ctx context.Context) error ***REMOVED***
	if u.deployment == nil ***REMOVED***
		return errors.New("the Update operation must have a Deployment set before Execute can be called")
	***REMOVED***
	batches := &driver.Batches***REMOVED***
		Identifier: "updates",
		Documents:  u.updates,
		Ordered:    u.ordered,
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         u.command,
		ProcessResponseFn: u.processResponse,
		Batches:           batches,
		RetryMode:         u.retry,
		Type:              driver.Write,
		Client:            u.session,
		Clock:             u.clock,
		CommandMonitor:    u.monitor,
		Database:          u.database,
		Deployment:        u.deployment,
		Selector:          u.selector,
		WriteConcern:      u.writeConcern,
		Crypt:             u.crypt,
		ServerAPI:         u.serverAPI,
		Timeout:           u.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (u *Update) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "update", u.collection)
	if u.bypassDocumentValidation != nil &&
		(desc.WireVersion != nil && desc.WireVersion.Includes(4)) ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *u.bypassDocumentValidation)
	***REMOVED***
	if u.comment.Type != bsontype.Type(0) ***REMOVED***
		dst = bsoncore.AppendValueElement(dst, "comment", u.comment)
	***REMOVED***
	if u.ordered != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "ordered", *u.ordered)
	***REMOVED***
	if u.hint != nil && *u.hint ***REMOVED***

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) ***REMOVED***
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 5")
		***REMOVED***
		if !u.writeConcern.Acknowledged() ***REMOVED***
			return nil, errUnacknowledgedHint
		***REMOVED***
	***REMOVED***
	if u.arrayFilters != nil && *u.arrayFilters ***REMOVED***
		if desc.WireVersion == nil || !desc.WireVersion.Includes(6) ***REMOVED***
			return nil, errors.New("the 'arrayFilters' command parameter requires a minimum server wire version of 6")
		***REMOVED***
	***REMOVED***
	if u.let != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "let", u.let)
	***REMOVED***

	return dst, nil
***REMOVED***

// BypassDocumentValidation allows the operation to opt-out of document level validation. Valid
// for server versions >= 3.2. For servers < 3.2, this setting is ignored.
func (u *Update) BypassDocumentValidation(bypassDocumentValidation bool) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.bypassDocumentValidation = &bypassDocumentValidation
	return u
***REMOVED***

// Hint is a flag to indicate that the update document contains a hint. Hint is only supported by
// servers >= 4.2. Older servers >= 3.4 will report an error for using the hint option. For servers <
// 3.4, the driver will return an error if the hint option is used.
func (u *Update) Hint(hint bool) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.hint = &hint
	return u
***REMOVED***

// ArrayFilters is a flag to indicate that the update document contains an arrayFilters field. This option is only
// supported on server versions 3.6 and higher. For servers < 3.6, the driver will return an error.
func (u *Update) ArrayFilters(arrayFilters bool) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.arrayFilters = &arrayFilters
	return u
***REMOVED***

// Ordered sets ordered. If true, when a write fails, the operation will return the error, when
// false write failures do not stop execution of the operation.
func (u *Update) Ordered(ordered bool) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.ordered = &ordered
	return u
***REMOVED***

// Updates specifies an array of update statements to perform when this operation is executed.
// Each update document must have the following structure:
// ***REMOVED***q: <query>, u: <update>, multi: <boolean>, collation: Optional<Document>, arrayFitlers: Optional<Array>, hint: Optional<string/Document>***REMOVED***.
func (u *Update) Updates(updates ...bsoncore.Document) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.updates = updates
	return u
***REMOVED***

// Session sets the session for this operation.
func (u *Update) Session(session *session.Client) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.session = session
	return u
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (u *Update) ClusterClock(clock *session.ClusterClock) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.clock = clock
	return u
***REMOVED***

// Collection sets the collection that this command will run against.
func (u *Update) Collection(collection string) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.collection = collection
	return u
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (u *Update) CommandMonitor(monitor *event.CommandMonitor) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.monitor = monitor
	return u
***REMOVED***

// Comment sets a value to help trace an operation.
func (u *Update) Comment(comment bsoncore.Value) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.comment = comment
	return u
***REMOVED***

// Database sets the database to run this operation against.
func (u *Update) Database(database string) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.database = database
	return u
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (u *Update) Deployment(deployment driver.Deployment) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.deployment = deployment
	return u
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (u *Update) ServerSelector(selector description.ServerSelector) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.selector = selector
	return u
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (u *Update) WriteConcern(writeConcern *writeconcern.WriteConcern) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.writeConcern = writeConcern
	return u
***REMOVED***

// Retry enables retryable writes for this operation. Retries are not handled automatically,
// instead a boolean is returned from Execute and SelectAndExecute that indicates if the
// operation can be retried. Retrying is handled by calling RetryExecute.
func (u *Update) Retry(retry driver.RetryMode) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.retry = &retry
	return u
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (u *Update) Crypt(crypt driver.Crypt) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.crypt = crypt
	return u
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (u *Update) ServerAPI(serverAPI *driver.ServerAPIOptions) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.serverAPI = serverAPI
	return u
***REMOVED***

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (u *Update) Let(let bsoncore.Document) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.let = let
	return u
***REMOVED***

// Timeout sets the timeout for this operation.
func (u *Update) Timeout(timeout *time.Duration) *Update ***REMOVED***
	if u == nil ***REMOVED***
		u = new(Update)
	***REMOVED***

	u.timeout = timeout
	return u
***REMOVED***
