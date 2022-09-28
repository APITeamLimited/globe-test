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
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// ListDatabases performs a listDatabases operation.
type ListDatabases struct ***REMOVED***
	filter              bsoncore.Document
	authorizedDatabases *bool
	nameOnly            *bool
	session             *session.Client
	clock               *session.ClusterClock
	monitor             *event.CommandMonitor
	database            string
	deployment          driver.Deployment
	readPreference      *readpref.ReadPref
	retry               *driver.RetryMode
	selector            description.ServerSelector
	crypt               driver.Crypt
	serverAPI           *driver.ServerAPIOptions
	timeout             *time.Duration

	result ListDatabasesResult
***REMOVED***

// ListDatabasesResult represents a listDatabases result returned by the server.
type ListDatabasesResult struct ***REMOVED***
	// An array of documents, one document for each database
	Databases []databaseRecord
	// The sum of the size of all the database files on disk in bytes.
	TotalSize int64
***REMOVED***

type databaseRecord struct ***REMOVED***
	Name       string
	SizeOnDisk int64 `bson:"sizeOnDisk"`
	Empty      bool
***REMOVED***

func buildListDatabasesResult(response bsoncore.Document) (ListDatabasesResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return ListDatabasesResult***REMOVED******REMOVED***, err
	***REMOVED***
	ir := ListDatabasesResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "totalSize":
			var ok bool
			ir.TotalSize, ok = element.Value().AsInt64OK()
			if !ok ***REMOVED***
				return ir, fmt.Errorf("response field 'totalSize' is type int64, but received BSON type %s: %s", element.Value().Type, element.Value())
			***REMOVED***
		case "databases":
			arr, ok := element.Value().ArrayOK()
			if !ok ***REMOVED***
				return ir, fmt.Errorf("response field 'databases' is type array, but received BSON type %s", element.Value().Type)
			***REMOVED***

			var tmp bsoncore.Document
			err := bson.Unmarshal(arr, &tmp)
			if err != nil ***REMOVED***
				return ir, err
			***REMOVED***

			records, err := tmp.Elements()
			if err != nil ***REMOVED***
				return ir, err
			***REMOVED***

			ir.Databases = make([]databaseRecord, len(records))
			for i, val := range records ***REMOVED***
				valueDoc, ok := val.Value().DocumentOK()
				if !ok ***REMOVED***
					return ir, fmt.Errorf("'databases' element is type document, but received BSON type %s", val.Value().Type)
				***REMOVED***

				elems, err := valueDoc.Elements()
				if err != nil ***REMOVED***
					return ir, err
				***REMOVED***

				for _, elem := range elems ***REMOVED***
					switch elem.Key() ***REMOVED***
					case "name":
						ir.Databases[i].Name, ok = elem.Value().StringValueOK()
						if !ok ***REMOVED***
							return ir, fmt.Errorf("response field 'name' is type string, but received BSON type %s", elem.Value().Type)
						***REMOVED***
					case "sizeOnDisk":
						ir.Databases[i].SizeOnDisk, ok = elem.Value().AsInt64OK()
						if !ok ***REMOVED***
							return ir, fmt.Errorf("response field 'sizeOnDisk' is type int64, but received BSON type %s", elem.Value().Type)
						***REMOVED***
					case "empty":
						ir.Databases[i].Empty, ok = elem.Value().BooleanOK()
						if !ok ***REMOVED***
							return ir, fmt.Errorf("response field 'empty' is type bool, but received BSON type %s", elem.Value().Type)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ir, nil
***REMOVED***

// NewListDatabases constructs and returns a new ListDatabases.
func NewListDatabases(filter bsoncore.Document) *ListDatabases ***REMOVED***
	return &ListDatabases***REMOVED***
		filter: filter,
	***REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (ld *ListDatabases) Result() ListDatabasesResult ***REMOVED*** return ld.result ***REMOVED***

func (ld *ListDatabases) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error

	ld.result, err = buildListDatabasesResult(info.ServerResponse)
	return err

***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (ld *ListDatabases) Execute(ctx context.Context) error ***REMOVED***
	if ld.deployment == nil ***REMOVED***
		return errors.New("the ListDatabases operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         ld.command,
		ProcessResponseFn: ld.processResponse,

		Client:         ld.session,
		Clock:          ld.clock,
		CommandMonitor: ld.monitor,
		Database:       ld.database,
		Deployment:     ld.deployment,
		ReadPreference: ld.readPreference,
		RetryMode:      ld.retry,
		Type:           driver.Read,
		Selector:       ld.selector,
		Crypt:          ld.crypt,
		ServerAPI:      ld.serverAPI,
		Timeout:        ld.timeout,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (ld *ListDatabases) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendInt32Element(dst, "listDatabases", 1)
	if ld.filter != nil ***REMOVED***

		dst = bsoncore.AppendDocumentElement(dst, "filter", ld.filter)
	***REMOVED***
	if ld.nameOnly != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "nameOnly", *ld.nameOnly)
	***REMOVED***
	if ld.authorizedDatabases != nil ***REMOVED***

		dst = bsoncore.AppendBooleanElement(dst, "authorizedDatabases", *ld.authorizedDatabases)
	***REMOVED***

	return dst, nil
***REMOVED***

// Filter determines what results are returned from listDatabases.
func (ld *ListDatabases) Filter(filter bsoncore.Document) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.filter = filter
	return ld
***REMOVED***

// NameOnly specifies whether to only return database names.
func (ld *ListDatabases) NameOnly(nameOnly bool) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.nameOnly = &nameOnly
	return ld
***REMOVED***

// AuthorizedDatabases specifies whether to only return databases which the user is authorized to use."
func (ld *ListDatabases) AuthorizedDatabases(authorizedDatabases bool) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.authorizedDatabases = &authorizedDatabases
	return ld
***REMOVED***

// Session sets the session for this operation.
func (ld *ListDatabases) Session(session *session.Client) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.session = session
	return ld
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (ld *ListDatabases) ClusterClock(clock *session.ClusterClock) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.clock = clock
	return ld
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (ld *ListDatabases) CommandMonitor(monitor *event.CommandMonitor) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.monitor = monitor
	return ld
***REMOVED***

// Database sets the database to run this operation against.
func (ld *ListDatabases) Database(database string) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.database = database
	return ld
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (ld *ListDatabases) Deployment(deployment driver.Deployment) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.deployment = deployment
	return ld
***REMOVED***

// ReadPreference set the read preference used with this operation.
func (ld *ListDatabases) ReadPreference(readPreference *readpref.ReadPref) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.readPreference = readPreference
	return ld
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (ld *ListDatabases) ServerSelector(selector description.ServerSelector) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.selector = selector
	return ld
***REMOVED***

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (ld *ListDatabases) Retry(retry driver.RetryMode) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.retry = &retry
	return ld
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (ld *ListDatabases) Crypt(crypt driver.Crypt) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.crypt = crypt
	return ld
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (ld *ListDatabases) ServerAPI(serverAPI *driver.ServerAPIOptions) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.serverAPI = serverAPI
	return ld
***REMOVED***

// Timeout sets the timeout for this operation.
func (ld *ListDatabases) Timeout(timeout *time.Duration) *ListDatabases ***REMOVED***
	if ld == nil ***REMOVED***
		ld = new(ListDatabases)
	***REMOVED***

	ld.timeout = timeout
	return ld
***REMOVED***
