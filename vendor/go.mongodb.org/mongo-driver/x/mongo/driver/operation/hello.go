// Copyright (C) MongoDB, Inc. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"runtime"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/version"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Hello is used to run the handshake operation.
type Hello struct ***REMOVED***
	appname            string
	compressors        []string
	saslSupportedMechs string
	d                  driver.Deployment
	clock              *session.ClusterClock
	speculativeAuth    bsoncore.Document
	topologyVersion    *description.TopologyVersion
	maxAwaitTimeMS     *int64
	serverAPI          *driver.ServerAPIOptions
	loadBalanced       bool

	res bsoncore.Document
***REMOVED***

var _ driver.Handshaker = (*Hello)(nil)

// NewHello constructs a Hello.
func NewHello() *Hello ***REMOVED*** return &Hello***REMOVED******REMOVED*** ***REMOVED***

// AppName sets the application name in the client metadata sent in this operation.
func (h *Hello) AppName(appname string) *Hello ***REMOVED***
	h.appname = appname
	return h
***REMOVED***

// ClusterClock sets the cluster clock for this operation.
func (h *Hello) ClusterClock(clock *session.ClusterClock) *Hello ***REMOVED***
	if h == nil ***REMOVED***
		h = new(Hello)
	***REMOVED***

	h.clock = clock
	return h
***REMOVED***

// Compressors sets the compressors that can be used.
func (h *Hello) Compressors(compressors []string) *Hello ***REMOVED***
	h.compressors = compressors
	return h
***REMOVED***

// SASLSupportedMechs retrieves the supported SASL mechanism for the given user when this operation
// is run.
func (h *Hello) SASLSupportedMechs(username string) *Hello ***REMOVED***
	h.saslSupportedMechs = username
	return h
***REMOVED***

// Deployment sets the Deployment for this operation.
func (h *Hello) Deployment(d driver.Deployment) *Hello ***REMOVED***
	h.d = d
	return h
***REMOVED***

// SpeculativeAuthenticate sets the document to be used for speculative authentication.
func (h *Hello) SpeculativeAuthenticate(doc bsoncore.Document) *Hello ***REMOVED***
	h.speculativeAuth = doc
	return h
***REMOVED***

// TopologyVersion sets the TopologyVersion to be used for heartbeats.
func (h *Hello) TopologyVersion(tv *description.TopologyVersion) *Hello ***REMOVED***
	h.topologyVersion = tv
	return h
***REMOVED***

// MaxAwaitTimeMS sets the maximum time for the server to wait for topology changes during a heartbeat.
func (h *Hello) MaxAwaitTimeMS(awaitTime int64) *Hello ***REMOVED***
	h.maxAwaitTimeMS = &awaitTime
	return h
***REMOVED***

// ServerAPI sets the server API version for this operation.
func (h *Hello) ServerAPI(serverAPI *driver.ServerAPIOptions) *Hello ***REMOVED***
	h.serverAPI = serverAPI
	return h
***REMOVED***

// LoadBalanced specifies whether or not this operation is being sent over a connection to a load balanced cluster.
func (h *Hello) LoadBalanced(lb bool) *Hello ***REMOVED***
	h.loadBalanced = lb
	return h
***REMOVED***

// Result returns the result of executing this operation.
func (h *Hello) Result(addr address.Address) description.Server ***REMOVED***
	return description.NewServer(addr, bson.Raw(h.res))
***REMOVED***

// handshakeCommand appends all necessary command fields as well as client metadata, SASL supported mechs, and compression.
func (h *Hello) handshakeCommand(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	dst, err := h.command(dst, desc)
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***

	if h.saslSupportedMechs != "" ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, "saslSupportedMechs", h.saslSupportedMechs)
	***REMOVED***
	if h.speculativeAuth != nil ***REMOVED***
		dst = bsoncore.AppendDocumentElement(dst, "speculativeAuthenticate", h.speculativeAuth)
	***REMOVED***
	var idx int32
	idx, dst = bsoncore.AppendArrayElementStart(dst, "compression")
	for i, compressor := range h.compressors ***REMOVED***
		dst = bsoncore.AppendStringElement(dst, strconv.Itoa(i), compressor)
	***REMOVED***
	dst, _ = bsoncore.AppendArrayEnd(dst, idx)

	// append client metadata
	idx, dst = bsoncore.AppendDocumentElementStart(dst, "client")

	didx, dst := bsoncore.AppendDocumentElementStart(dst, "driver")
	dst = bsoncore.AppendStringElement(dst, "name", "mongo-go-driver")
	dst = bsoncore.AppendStringElement(dst, "version", version.Driver)
	dst, _ = bsoncore.AppendDocumentEnd(dst, didx)

	didx, dst = bsoncore.AppendDocumentElementStart(dst, "os")
	dst = bsoncore.AppendStringElement(dst, "type", runtime.GOOS)
	dst = bsoncore.AppendStringElement(dst, "architecture", runtime.GOARCH)
	dst, _ = bsoncore.AppendDocumentEnd(dst, didx)

	dst = bsoncore.AppendStringElement(dst, "platform", runtime.Version())
	if h.appname != "" ***REMOVED***
		didx, dst = bsoncore.AppendDocumentElementStart(dst, "application")
		dst = bsoncore.AppendStringElement(dst, "name", h.appname)
		dst, _ = bsoncore.AppendDocumentEnd(dst, didx)
	***REMOVED***
	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)

	return dst, nil
***REMOVED***

// command appends all necessary command fields.
func (h *Hello) command(dst []byte, desc description.SelectedServer) ([]byte, error) ***REMOVED***
	// Use "hello" if topology is LoadBalanced, API version is declared or server
	// has responded with "helloOk". Otherwise, use legacy hello.
	if desc.Kind == description.LoadBalanced || h.serverAPI != nil || desc.Server.HelloOK ***REMOVED***
		dst = bsoncore.AppendInt32Element(dst, "hello", 1)
	***REMOVED*** else ***REMOVED***
		dst = bsoncore.AppendInt32Element(dst, internal.LegacyHello, 1)
	***REMOVED***
	dst = bsoncore.AppendBooleanElement(dst, "helloOk", true)

	if tv := h.topologyVersion; tv != nil ***REMOVED***
		var tvIdx int32

		tvIdx, dst = bsoncore.AppendDocumentElementStart(dst, "topologyVersion")
		dst = bsoncore.AppendObjectIDElement(dst, "processId", tv.ProcessID)
		dst = bsoncore.AppendInt64Element(dst, "counter", tv.Counter)
		dst, _ = bsoncore.AppendDocumentEnd(dst, tvIdx)
	***REMOVED***
	if h.maxAwaitTimeMS != nil ***REMOVED***
		dst = bsoncore.AppendInt64Element(dst, "maxAwaitTimeMS", *h.maxAwaitTimeMS)
	***REMOVED***
	if h.loadBalanced ***REMOVED***
		// The loadBalanced parameter should only be added if it's true. We should never explicitly send
		// loadBalanced=false per the load balancing spec.
		dst = bsoncore.AppendBooleanElement(dst, "loadBalanced", true)
	***REMOVED***

	return dst, nil
***REMOVED***

// Execute runs this operation.
func (h *Hello) Execute(ctx context.Context) error ***REMOVED***
	if h.d == nil ***REMOVED***
		return errors.New("a Hello must have a Deployment set before Execute can be called")
	***REMOVED***

	return h.createOperation().Execute(ctx, nil)
***REMOVED***

// StreamResponse gets the next streaming Hello response from the server.
func (h *Hello) StreamResponse(ctx context.Context, conn driver.StreamerConnection) error ***REMOVED***
	return h.createOperation().ExecuteExhaust(ctx, conn, nil)
***REMOVED***

func (h *Hello) createOperation() driver.Operation ***REMOVED***
	return driver.Operation***REMOVED***
		Clock:      h.clock,
		CommandFn:  h.command,
		Database:   "admin",
		Deployment: h.d,
		ProcessResponseFn: func(info driver.ResponseInfo) error ***REMOVED***
			h.res = info.ServerResponse
			return nil
		***REMOVED***,
		ServerAPI: h.serverAPI,
	***REMOVED***
***REMOVED***

// GetHandshakeInformation performs the MongoDB handshake for the provided connection and returns the relevant
// information about the server. This function implements the driver.Handshaker interface.
func (h *Hello) GetHandshakeInformation(ctx context.Context, _ address.Address, c driver.Connection) (driver.HandshakeInformation, error) ***REMOVED***
	err := driver.Operation***REMOVED***
		Clock:      h.clock,
		CommandFn:  h.handshakeCommand,
		Deployment: driver.SingleConnectionDeployment***REMOVED***c***REMOVED***,
		Database:   "admin",
		ProcessResponseFn: func(info driver.ResponseInfo) error ***REMOVED***
			h.res = info.ServerResponse
			return nil
		***REMOVED***,
		ServerAPI: h.serverAPI,
	***REMOVED***.Execute(ctx, nil)
	if err != nil ***REMOVED***
		return driver.HandshakeInformation***REMOVED******REMOVED***, err
	***REMOVED***

	info := driver.HandshakeInformation***REMOVED***
		Description: h.Result(c.Address()),
	***REMOVED***
	if speculativeAuthenticate, ok := h.res.Lookup("speculativeAuthenticate").DocumentOK(); ok ***REMOVED***
		info.SpeculativeAuthenticate = speculativeAuthenticate
	***REMOVED***
	if serverConnectionID, ok := h.res.Lookup("connectionId").Int32OK(); ok ***REMOVED***
		info.ServerConnectionID = &serverConnectionID
	***REMOVED***
	// Cast to bson.Raw to lookup saslSupportedMechs to avoid converting from bsoncore.Value to bson.RawValue for the
	// StringSliceFromRawValue call.
	if saslSupportedMechs, lookupErr := bson.Raw(h.res).LookupErr("saslSupportedMechs"); lookupErr == nil ***REMOVED***
		info.SaslSupportedMechs, err = internal.StringSliceFromRawValue("saslSupportedMechs", saslSupportedMechs)
	***REMOVED***
	return info, err
***REMOVED***

// FinishHandshake implements the Handshaker interface. This is a no-op function because a non-authenticated connection
// does not do anything besides the initial Hello for a handshake.
func (h *Hello) FinishHandshake(context.Context, driver.Connection) error ***REMOVED***
	return nil
***REMOVED***
