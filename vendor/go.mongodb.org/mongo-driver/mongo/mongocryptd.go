// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	defaultServerSelectionTimeout = 10 * time.Second
	defaultURI                    = "mongodb://localhost:27020"
	defaultPath                   = "mongocryptd"
	serverSelectionTimeoutStr     = "server selection error"
)

var defaultTimeoutArgs = []string***REMOVED***"--idleShutdownTimeoutSecs=60"***REMOVED***
var databaseOpts = options.Database().SetReadConcern(readconcern.New()).SetReadPreference(readpref.Primary())

type mongocryptdClient struct ***REMOVED***
	bypassSpawn bool
	client      *Client
	path        string
	spawnArgs   []string
***REMOVED***

func newMongocryptdClient(cryptSharedLibAvailable bool, opts *options.AutoEncryptionOptions) (*mongocryptdClient, error) ***REMOVED***
	// create mcryptClient instance and spawn process if necessary
	var bypassSpawn bool
	var bypassAutoEncryption bool

	if bypass, ok := opts.ExtraOptions["mongocryptdBypassSpawn"]; ok ***REMOVED***
		bypassSpawn = bypass.(bool)
	***REMOVED***
	if opts.BypassAutoEncryption != nil ***REMOVED***
		bypassAutoEncryption = *opts.BypassAutoEncryption
	***REMOVED***

	bypassQueryAnalysis := opts.BypassQueryAnalysis != nil && *opts.BypassQueryAnalysis

	mc := &mongocryptdClient***REMOVED***
		// mongocryptd should not be spawned if any of these conditions are true:
		// - mongocryptdBypassSpawn is passed
		// - bypassAutoEncryption is true because mongocryptd is not used during decryption
		// - bypassQueryAnalysis is true because mongocryptd is not used during decryption
		// - the crypt_shared library is available because it replaces all mongocryptd functionality.
		bypassSpawn: bypassSpawn || bypassAutoEncryption || bypassQueryAnalysis || cryptSharedLibAvailable,
	***REMOVED***

	if !mc.bypassSpawn ***REMOVED***
		mc.path, mc.spawnArgs = createSpawnArgs(opts.ExtraOptions)
		if err := mc.spawnProcess(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// get connection string
	uri := defaultURI
	if u, ok := opts.ExtraOptions["mongocryptdURI"]; ok ***REMOVED***
		uri = u.(string)
	***REMOVED***

	// create client
	client, err := NewClient(options.Client().ApplyURI(uri).SetServerSelectionTimeout(defaultServerSelectionTimeout))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mc.client = client

	return mc, nil
***REMOVED***

// markCommand executes the given command on mongocryptd.
func (mc *mongocryptdClient) markCommand(ctx context.Context, dbName string, cmd bsoncore.Document) (bsoncore.Document, error) ***REMOVED***
	// Remove the explicit session from the context if one is set.
	// The explicit session will be from a different client.
	// If an explicit session is set, it is applied after automatic encryption.
	ctx = NewSessionContext(ctx, nil)
	db := mc.client.Database(dbName, databaseOpts)

	res, err := db.RunCommand(ctx, cmd).DecodeBytes()
	// propagate original result
	if err == nil ***REMOVED***
		return bsoncore.Document(res), nil
	***REMOVED***
	// wrap original error
	if mc.bypassSpawn || !strings.Contains(err.Error(), serverSelectionTimeoutStr) ***REMOVED***
		return nil, MongocryptdError***REMOVED***Wrapped: err***REMOVED***
	***REMOVED***

	// re-spawn and retry
	if err = mc.spawnProcess(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	res, err = db.RunCommand(ctx, cmd).DecodeBytes()
	if err != nil ***REMOVED***
		return nil, MongocryptdError***REMOVED***Wrapped: err***REMOVED***
	***REMOVED***
	return bsoncore.Document(res), nil
***REMOVED***

// connect connects the underlying Client instance. This must be called before performing any mark operations.
func (mc *mongocryptdClient) connect(ctx context.Context) error ***REMOVED***
	return mc.client.Connect(ctx)
***REMOVED***

// disconnect disconnects the underlying Client instance. This should be called after all operations have completed.
func (mc *mongocryptdClient) disconnect(ctx context.Context) error ***REMOVED***
	return mc.client.Disconnect(ctx)
***REMOVED***

func (mc *mongocryptdClient) spawnProcess() error ***REMOVED***
	// Ignore gosec warning about subprocess launched with externally-provided path variable.
	/* #nosec G204 */
	cmd := exec.Command(mc.path, mc.spawnArgs...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
***REMOVED***

// createSpawnArgs creates arguments to spawn mcryptClient. It returns the path and a slice of arguments.
func createSpawnArgs(opts map[string]interface***REMOVED******REMOVED***) (string, []string) ***REMOVED***
	var spawnArgs []string

	// get command path
	path := defaultPath
	if p, ok := opts["mongocryptdPath"]; ok ***REMOVED***
		path = p.(string)
	***REMOVED***

	// add specified options
	if sa, ok := opts["mongocryptdSpawnArgs"]; ok ***REMOVED***
		spawnArgs = append(spawnArgs, sa.([]string)...)
	***REMOVED***

	// add timeout options if necessary
	var foundTimeout bool
	for _, arg := range spawnArgs ***REMOVED***
		// need to use HasPrefix instead of doing an exact equality check because both
		// mongocryptd supports both [--idleShutdownTimeoutSecs, 0] and [--idleShutdownTimeoutSecs=0]
		if strings.HasPrefix(arg, "--idleShutdownTimeoutSecs") ***REMOVED***
			foundTimeout = true
			break
		***REMOVED***
	***REMOVED***
	if !foundTimeout ***REMOVED***
		spawnArgs = append(spawnArgs, defaultTimeoutArgs...)
	***REMOVED***

	return path, spawnArgs
***REMOVED***
