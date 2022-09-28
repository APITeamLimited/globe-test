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

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// DropCollection performs a drop operation.
type DropCollection struct ***REMOVED***
	session      *session.Client
	clock        *session.ClusterClock
	collection   string
	monitor      *event.CommandMonitor
	crypt        driver.Crypt
	database     string
	deployment   driver.Deployment
	selector     description.ServerSelector
	writeConcern *writeconcern.WriteConcern
	result       DropCollectionResult
	serverAPI    *driver.ServerAPIOptions
***REMOVED***

// DropCollectionResult represents a dropCollection result returned by the server.
type DropCollectionResult struct ***REMOVED***
	// The number of indexes in the dropped collection.
	NIndexesWas int32
	// The namespace of the dropped collection.
	Ns string
***REMOVED***

func buildDropCollectionResult(response bsoncore.Document) (DropCollectionResult, error) ***REMOVED***
	elements, err := response.Elements()
	if err != nil ***REMOVED***
		return DropCollectionResult***REMOVED******REMOVED***, err
	***REMOVED***
	dcr := DropCollectionResult***REMOVED******REMOVED***
	for _, element := range elements ***REMOVED***
		switch element.Key() ***REMOVED***
		case "nIndexesWas":
			var ok bool
			dcr.NIndexesWas, ok = element.Value().AsInt32OK()
			if !ok ***REMOVED***
				return dcr, fmt.Errorf("response field 'nIndexesWas' is type int32, but received BSON type %s", element.Value().Type)
			***REMOVED***
		case "ns":
			var ok bool
			dcr.Ns, ok = element.Value().StringValueOK()
			if !ok ***REMOVED***
				return dcr, fmt.Errorf("response field 'ns' is type string, but received BSON type %s", element.Value().Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dcr, nil
***REMOVED***

// NewDropCollection constructs and returns a new DropCollection.
func NewDropCollection() *DropCollection ***REMOVED***
	return &DropCollection***REMOVED******REMOVED***
***REMOVED***

// Result returns the result of executing this operation.
func (dc *DropCollection) Result() DropCollectionResult ***REMOVED*** return dc.result ***REMOVED***

func (dc *DropCollection) processResponse(info driver.ResponseInfo) error ***REMOVED***
	var err error
	dc.result, err = buildDropCollectionResult(info.ServerResponse)
	return err
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (dc *DropCollection) Execute(ctx context.Context) error ***REMOVED***
	if dc.deployment == nil ***REMOVED***
		return errors.New("the DropCollection operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         dc.command,
		ProcessResponseFn: dc.processResponse,
		Client:            dc.session,
		Clock:             dc.clock,
		CommandMonitor:    dc.monitor,
		Crypt:             dc.crypt,
		Database:          dc.database,
		Deployment:        dc.deployment,
		Selector:          dc.selector,
		WriteConcern:      dc.writeConcern,
		ServerAPI:         dc.serverAPI,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (dc *DropCollection) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst = bsoncore.AppendStringElement(dst, "drop", dc.collection)
	return dst, nil
***REMOVED***

// Session sets the session for this operation.
func (dc *DropCollection) Session(session *session.Client) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.session = session
	return dc
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (dc *DropCollection) ClusterClock(clock *session.ClusterClock) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.clock = clock
	return dc
***REMOVED***

// Collection sets the collection that this command will run against.
func (dc *DropCollection) Collection(collection string) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.collection = collection
	return dc
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (dc *DropCollection) CommandMonitor(monitor *event.CommandMonitor) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.monitor = monitor
	return dc
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (dc *DropCollection) Crypt(crypt driver.Crypt) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.crypt = crypt
	return dc
***REMOVED***

// Database sets the database to run this operation against.
func (dc *DropCollection) Database(database string) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.database = database
	return dc
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (dc *DropCollection) Deployment(deployment driver.Deployment) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.deployment = deployment
	return dc
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (dc *DropCollection) ServerSelector(selector description.ServerSelector) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.selector = selector
	return dc
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (dc *DropCollection) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.writeConcern = writeConcern
	return dc
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (dc *DropCollection) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropCollection ***REMOVED***
	if dc == nil ***REMOVED***
		dc = new(DropCollection)
	***REMOVED***

	dc.serverAPI = serverAPI
	return dc
***REMOVED***
