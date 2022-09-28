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

// Create represents a create operation.
type Create struct ***REMOVED***
	capped                       *bool
	collation                    bsoncore.Document
	changeStreamPreAndPostImages bsoncore.Document
	collectionName               *string
	indexOptionDefaults          bsoncore.Document
	max                          *int64
	pipeline                     bsoncore.Document
	size                         *int64
	storageEngine                bsoncore.Document
	validationAction             *string
	validationLevel              *string
	validator                    bsoncore.Document
	viewOn                       *string
	session                      *session.Client
	clock                        *session.ClusterClock
	monitor                      *event.CommandMonitor
	crypt                        driver.Crypt
	database                     string
	deployment                   driver.Deployment
	selector                     description.ServerSelector
	writeConcern                 *writeconcern.WriteConcern
	serverAPI                    *driver.ServerAPIOptions
	expireAfterSeconds           *int64
	timeSeries                   bsoncore.Document
	encryptedFields              bsoncore.Document
	clusteredIndex               bsoncore.Document
***REMOVED***

// NewCreate constructs and returns a new Create.
func NewCreate(collectionName string) *Create ***REMOVED***
	return &Create***REMOVED***
		collectionName: &collectionName,
	***REMOVED***
***REMOVED***

func (c *Create) processResponse(driver.ResponseInfo) error ***REMOVED***
	return nil
***REMOVED***

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (c *Create) Execute(ctx context.Context) error ***REMOVED***
	if c.deployment == nil ***REMOVED***
		return errors.New("the Create operation must have a Deployment set before Execute can be called")
	***REMOVED***

	return driver.Operation***REMOVED***
		CommandFn:         c.command,
		ProcessResponseFn: c.processResponse,
		Client:            c.session,
		Clock:             c.clock,
		CommandMonitor:    c.monitor,
		Crypt:             c.crypt,
		Database:          c.database,
		Deployment:        c.deployment,
		Selector:          c.selector,
		WriteConcern:      c.writeConcern,
		ServerAPI:         c.serverAPI,
	***REMOVED***.Execute(ctx, nil)

***REMOVED***

func (c *Create) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	if c.collectionName != nil ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "create", *c.collectionName)
	***REMOVED***
	if c.capped != nil ***REMOVED***
		dst = bsoncore.AppendBooleanElement(dst, "capped", *c.capped)
	***REMOVED***
	if c.changeStreamPreAndPostImages != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "changeStreamPreAndPostImages", c.changeStreamPreAndPostImages)
	***REMOVED***
	if c.collation != nil ***REMOVED***
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) ***REMOVED***
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "collation", c.collation)
	***REMOVED***
	if c.indexOptionDefaults != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "indexOptionDefaults", c.indexOptionDefaults)
	***REMOVED***
	if c.max != nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "max", *c.max)
	***REMOVED***
	if c.pipeline != nil ***REMOVED***
		dst = bsoncore.AppendArrayElement(dst, "pipeline", c.pipeline)
	***REMOVED***
	if c.size != nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "size", *c.size)
	***REMOVED***
	if c.storageEngine != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "storageEngine", c.storageEngine)
	***REMOVED***
	if c.validationAction != nil ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "validationAction", *c.validationAction)
	***REMOVED***
	if c.validationLevel != nil ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "validationLevel", *c.validationLevel)
	***REMOVED***
	if c.validator != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "validator", c.validator)
	***REMOVED***
	if c.viewOn != nil ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "viewOn", *c.viewOn)
	***REMOVED***
	if c.expireAfterSeconds != nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "expireAfterSeconds", *c.expireAfterSeconds)
	***REMOVED***
	if c.timeSeries != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "timeseries", c.timeSeries)
	***REMOVED***
	if c.encryptedFields != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "encryptedFields", c.encryptedFields)
	***REMOVED***
	if c.clusteredIndex != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "clusteredIndex", c.clusteredIndex)
	***REMOVED***
	return dst, nil
***REMOVED***

// Capped specifies if the collection is capped.
func (c *Create) Capped(capped bool) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.capped = &capped
	return c
***REMOVED***

// Collation specifies a collation. This option is only valid for server versions 3.4 and above.
func (c *Create) Collation(collation bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.collation = collation
	return c
***REMOVED***

// ChangeStreamPreAndPostImages specifies how change streams opened against the collection can return pre-
// and post-images of updated documents. This option is only valid for server versions 6.0 and above.
func (c *Create) ChangeStreamPreAndPostImages(csppi bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.changeStreamPreAndPostImages = csppi
	return c
***REMOVED***

// CollectionName specifies the name of the collection to create.
func (c *Create) CollectionName(collectionName string) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.collectionName = &collectionName
	return c
***REMOVED***

// IndexOptionDefaults specifies a default configuration for indexes on the collection.
func (c *Create) IndexOptionDefaults(indexOptionDefaults bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.indexOptionDefaults = indexOptionDefaults
	return c
***REMOVED***

// Max specifies the maximum number of documents allowed in a capped collection.
func (c *Create) Max(max int64) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.max = &max
	return c
***REMOVED***

// Pipeline specifies the agggregtion pipeline to be run against the source to create the view.
func (c *Create) Pipeline(pipeline bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.pipeline = pipeline
	return c
***REMOVED***

// Size specifies the maximum size in bytes for a capped collection.
func (c *Create) Size(size int64) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.size = &size
	return c
***REMOVED***

// StorageEngine specifies the storage engine to use for the index.
func (c *Create) StorageEngine(storageEngine bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.storageEngine = storageEngine
	return c
***REMOVED***

// ValidationAction specifies what should happen if a document being inserted does not pass validation.
func (c *Create) ValidationAction(validationAction string) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.validationAction = &validationAction
	return c
***REMOVED***

// ValidationLevel specifies how strictly the server applies validation rules to existing documents in the collection
// during update operations.
func (c *Create) ValidationLevel(validationLevel string) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.validationLevel = &validationLevel
	return c
***REMOVED***

// Validator specifies validation rules for the collection.
func (c *Create) Validator(validator bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.validator = validator
	return c
***REMOVED***

// ViewOn specifies the name of the source collection or view on which the view will be created.
func (c *Create) ViewOn(viewOn string) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.viewOn = &viewOn
	return c
***REMOVED***

// Session sets the session for this operation.
func (c *Create) Session(session *session.Client) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.session = session
	return c
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (c *Create) ClusterClock(clock *session.ClusterClock) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.clock = clock
	return c
***REMOVED***

// CommandMonitor sets the monitor to use for APM events.
func (c *Create) CommandMonitor(monitor *event.CommandMonitor) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.monitor = monitor
	return c
***REMOVED***

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (c *Create) Crypt(crypt driver.Crypt) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.crypt = crypt
	return c
***REMOVED***

// Database sets the database to run this operation against.
func (c *Create) Database(database string) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.database = database
	return c
***REMOVED***

// Deployment sets the deployment to use for this operation.
func (c *Create) Deployment(deployment driver.Deployment) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.deployment = deployment
	return c
***REMOVED***

// ServerSelector sets the selector used to retrieve a server.
func (c *Create) ServerSelector(selector description.ServerSelector) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.selector = selector
	return c
***REMOVED***

// WriteConcern sets the write concern for this operation.
func (c *Create) WriteConcern(writeConcern *writeconcern.WriteConcern) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.writeConcern = writeConcern
	return c
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (c *Create) ServerAPI(serverAPI *driver.ServerAPIOptions) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.serverAPI = serverAPI
	return c
***REMOVED***

// ExpireAfterSeconds sets the seconds to wait before deleting old time-series data.
func (c *Create) ExpireAfterSeconds(eas int64) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.expireAfterSeconds = &eas
	return c
***REMOVED***

// TimeSeries sets the time series options for this operation.
func (c *Create) TimeSeries(timeSeries bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.timeSeries = timeSeries
	return c
***REMOVED***

// EncryptedFields sets the EncryptedFields for this operation.
func (c *Create) EncryptedFields(ef bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.encryptedFields = ef
	return c
***REMOVED***

// ClusteredIndex sets the ClusteredIndex option for this operation.
func (c *Create) ClusteredIndex(ci bsoncore.Document) *Create ***REMOVED***
	if c == nil ***REMOVED***
		c = new(Create)
	***REMOVED***

	c.clusteredIndex = ci
	return c
***REMOVED***
