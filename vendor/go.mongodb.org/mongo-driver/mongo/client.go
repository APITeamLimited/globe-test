// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	mcopts "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

const (
	defaultLocalThreshold        = 15 * time.Millisecond
	defaultMaxPoolSize    uint64 = 100
)

var (
	// keyVaultCollOpts specifies options used to communicate with the key vault collection
	keyVaultCollOpts = options.Collection().SetReadConcern(readconcern.Majority()).
				SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	endSessionsBatchSize = 10000
)

// Client is a handle representing a pool of connections to a MongoDB deployment. It is safe for concurrent use by
// multiple goroutines.
//
// The Client type opens and closes connections automatically and maintains a pool of idle connections. For
// connection pool configuration options, see documentation for the ClientOptions type in the mongo/options package.
type Client struct ***REMOVED***
	id              uuid.UUID
	topologyOptions []topology.Option
	deployment      driver.Deployment
	localThreshold  time.Duration
	retryWrites     bool
	retryReads      bool
	clock           *session.ClusterClock
	readPreference  *readpref.ReadPref
	readConcern     *readconcern.ReadConcern
	writeConcern    *writeconcern.WriteConcern
	registry        *bsoncodec.Registry
	monitor         *event.CommandMonitor
	serverAPI       *driver.ServerAPIOptions
	serverMonitor   *event.ServerMonitor
	sessionPool     *session.Pool
	timeout         *time.Duration

	// client-side encryption fields
	keyVaultClientFLE  *Client
	keyVaultCollFLE    *Collection
	mongocryptdFLE     *mongocryptdClient
	cryptFLE           driver.Crypt
	metadataClientFLE  *Client
	internalClientFLE  *Client
	encryptedFieldsMap map[string]interface***REMOVED******REMOVED***
***REMOVED***

// Connect creates a new Client and then initializes it using the Connect method. This is equivalent to calling
// NewClient followed by Client.Connect.
//
// When creating an options.ClientOptions, the order the methods are called matters. Later Set*
// methods will overwrite the values from previous Set* method invocations. This includes the
// ApplyURI method. This allows callers to determine the order of precedence for option
// application. For instance, if ApplyURI is called before SetAuth, the Credential from
// SetAuth will overwrite the values from the connection string. If ApplyURI is called
// after SetAuth, then its values will overwrite those from SetAuth.
//
// The opts parameter is processed using options.MergeClientOptions, which will overwrite entire
// option fields of previous options, there is no partial overwriting. For example, if Username is
// set in the Auth field for the first option, and Password is set for the second but with no
// Username, after the merge the Username field will be empty.
//
// The NewClient function does not do any I/O and returns an error if the given options are invalid.
// The Client.Connect method starts background goroutines to monitor the state of the deployment and does not do
// any I/O in the main goroutine to prevent the main goroutine from blocking. Therefore, it will not error if the
// deployment is down.
//
// The Client.Ping method can be used to verify that the deployment is successfully connected and the
// Client was correctly configured.
func Connect(ctx context.Context, opts ...*options.ClientOptions) (*Client, error) ***REMOVED***
	c, err := NewClient(opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = c.Connect(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c, nil
***REMOVED***

// NewClient creates a new client to connect to a deployment specified by the uri.
//
// When creating an options.ClientOptions, the order the methods are called matters. Later Set*
// methods will overwrite the values from previous Set* method invocations. This includes the
// ApplyURI method. This allows callers to determine the order of precedence for option
// application. For instance, if ApplyURI is called before SetAuth, the Credential from
// SetAuth will overwrite the values from the connection string. If ApplyURI is called
// after SetAuth, then its values will overwrite those from SetAuth.
//
// The opts parameter is processed using options.MergeClientOptions, which will overwrite entire
// option fields of previous options, there is no partial overwriting. For example, if Username is
// set in the Auth field for the first option, and Password is set for the second but with no
// Username, after the merge the Username field will be empty.
func NewClient(opts ...*options.ClientOptions) (*Client, error) ***REMOVED***
	clientOpt := options.MergeClientOptions(opts...)

	id, err := uuid.New()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	client := &Client***REMOVED***id: id***REMOVED***

	err = client.configure(clientOpt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if client.deployment == nil ***REMOVED***
		client.deployment, err = topology.New(client.topologyOptions...)
		if err != nil ***REMOVED***
			return nil, replaceErrors(err)
		***REMOVED***
	***REMOVED***
	return client, nil
***REMOVED***

// Connect initializes the Client by starting background monitoring goroutines.
// If the Client was created using the NewClient function, this method must be called before a Client can be used.
//
// Connect starts background goroutines to monitor the state of the deployment and does not do any I/O in the main
// goroutine. The Client.Ping method can be used to verify that the connection was created successfully.
func (c *Client) Connect(ctx context.Context) error ***REMOVED***
	if connector, ok := c.deployment.(driver.Connector); ok ***REMOVED***
		err := connector.Connect()
		if err != nil ***REMOVED***
			return replaceErrors(err)
		***REMOVED***
	***REMOVED***

	if c.mongocryptdFLE != nil ***REMOVED***
		if err := c.mongocryptdFLE.connect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if c.internalClientFLE != nil ***REMOVED***
		if err := c.internalClientFLE.Connect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if c.keyVaultClientFLE != nil && c.keyVaultClientFLE != c.internalClientFLE && c.keyVaultClientFLE != c ***REMOVED***
		if err := c.keyVaultClientFLE.Connect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if c.metadataClientFLE != nil && c.metadataClientFLE != c.internalClientFLE && c.metadataClientFLE != c ***REMOVED***
		if err := c.metadataClientFLE.Connect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var updateChan <-chan description.Topology
	if subscriber, ok := c.deployment.(driver.Subscriber); ok ***REMOVED***
		sub, err := subscriber.Subscribe()
		if err != nil ***REMOVED***
			return replaceErrors(err)
		***REMOVED***
		updateChan = sub.Updates
	***REMOVED***
	c.sessionPool = session.NewPool(updateChan)
	return nil
***REMOVED***

// Disconnect closes sockets to the topology referenced by this Client. It will
// shut down any monitoring goroutines, close the idle connection pool, and will
// wait until all the in use connections have been returned to the connection
// pool and closed before returning. If the context expires via cancellation,
// deadline, or timeout before the in use connections have returned, the in use
// connections will be closed, resulting in the failure of any in flight read
// or write operations. If this method returns with no errors, all connections
// associated with this Client have been closed.
func (c *Client) Disconnect(ctx context.Context) error ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	c.endSessions(ctx)
	if c.mongocryptdFLE != nil ***REMOVED***
		if err := c.mongocryptdFLE.disconnect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if c.internalClientFLE != nil ***REMOVED***
		if err := c.internalClientFLE.Disconnect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if c.keyVaultClientFLE != nil && c.keyVaultClientFLE != c.internalClientFLE && c.keyVaultClientFLE != c ***REMOVED***
		if err := c.keyVaultClientFLE.Disconnect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if c.metadataClientFLE != nil && c.metadataClientFLE != c.internalClientFLE && c.metadataClientFLE != c ***REMOVED***
		if err := c.metadataClientFLE.Disconnect(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if c.cryptFLE != nil ***REMOVED***
		c.cryptFLE.Close()
	***REMOVED***

	if disconnector, ok := c.deployment.(driver.Disconnector); ok ***REMOVED***
		return replaceErrors(disconnector.Disconnect(ctx))
	***REMOVED***
	return nil
***REMOVED***

// Ping sends a ping command to verify that the client can connect to the deployment.
//
// The rp parameter is used to determine which server is selected for the operation.
// If it is nil, the client's read preference is used.
//
// If the server is down, Ping will try to select a server until the client's server selection timeout expires.
// This can be configured through the ClientOptions.SetServerSelectionTimeout option when creating a new Client.
// After the timeout expires, a server selection error is returned.
//
// Using Ping reduces application resilience because applications starting up will error if the server is temporarily
// unavailable or is failing over (e.g. during autoscaling due to a load spike).
func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	if rp == nil ***REMOVED***
		rp = c.readPreference
	***REMOVED***

	db := c.Database("admin")
	res := db.RunCommand(ctx, bson.D***REMOVED***
		***REMOVED***"ping", 1***REMOVED***,
	***REMOVED***, options.RunCmd().SetReadPreference(rp))

	return replaceErrors(res.Err())
***REMOVED***

// StartSession starts a new session configured with the given options.
//
// StartSession does not actually communicate with the server and will not error if the client is
// disconnected.
//
// StartSession is safe to call from multiple goroutines concurrently. However, Sessions returned by StartSession are
// not safe for concurrent use by multiple goroutines.
//
// If the DefaultReadConcern, DefaultWriteConcern, or DefaultReadPreference options are not set, the client's read
// concern, write concern, or read preference will be used, respectively.
func (c *Client) StartSession(opts ...*options.SessionOptions) (Session, error) ***REMOVED***
	if c.sessionPool == nil ***REMOVED***
		return nil, ErrClientDisconnected
	***REMOVED***

	sopts := options.MergeSessionOptions(opts...)
	coreOpts := &session.ClientOptions***REMOVED***
		DefaultReadConcern:    c.readConcern,
		DefaultReadPreference: c.readPreference,
		DefaultWriteConcern:   c.writeConcern,
	***REMOVED***
	if sopts.CausalConsistency != nil ***REMOVED***
		coreOpts.CausalConsistency = sopts.CausalConsistency
	***REMOVED***
	if sopts.DefaultReadConcern != nil ***REMOVED***
		coreOpts.DefaultReadConcern = sopts.DefaultReadConcern
	***REMOVED***
	if sopts.DefaultWriteConcern != nil ***REMOVED***
		coreOpts.DefaultWriteConcern = sopts.DefaultWriteConcern
	***REMOVED***
	if sopts.DefaultReadPreference != nil ***REMOVED***
		coreOpts.DefaultReadPreference = sopts.DefaultReadPreference
	***REMOVED***
	if sopts.DefaultMaxCommitTime != nil ***REMOVED***
		coreOpts.DefaultMaxCommitTime = sopts.DefaultMaxCommitTime
	***REMOVED***
	if sopts.Snapshot != nil ***REMOVED***
		coreOpts.Snapshot = sopts.Snapshot
	***REMOVED***

	sess, err := session.NewClientSession(c.sessionPool, c.id, session.Explicit, coreOpts)
	if err != nil ***REMOVED***
		return nil, replaceErrors(err)
	***REMOVED***

	// Writes are not retryable on standalones, so let operation determine whether to retry
	sess.RetryWrite = false
	sess.RetryRead = c.retryReads

	return &sessionImpl***REMOVED***
		clientSession: sess,
		client:        c,
		deployment:    c.deployment,
	***REMOVED***, nil
***REMOVED***

func (c *Client) endSessions(ctx context.Context) ***REMOVED***
	if c.sessionPool == nil ***REMOVED***
		return
	***REMOVED***

	sessionIDs := c.sessionPool.IDSlice()
	op := operation.NewEndSessions(nil).ClusterClock(c.clock).Deployment(c.deployment).
		ServerSelector(description.ReadPrefSelector(readpref.PrimaryPreferred())).CommandMonitor(c.monitor).
		Database("admin").Crypt(c.cryptFLE).ServerAPI(c.serverAPI)

	totalNumIDs := len(sessionIDs)
	var currentBatch []bsoncore.Document
	for i := 0; i < totalNumIDs; i++ ***REMOVED***
		currentBatch = append(currentBatch, sessionIDs[i])

		// If we are at the end of a batch or the end of the overall IDs array, execute the operation.
		if ((i+1)%endSessionsBatchSize) == 0 || i == totalNumIDs-1 ***REMOVED***
			// Ignore all errors when ending sessions.
			_, marshalVal, err := bson.MarshalValue(currentBatch)
			if err == nil ***REMOVED***
				_ = op.SessionIDs(marshalVal).Execute(ctx)
			***REMOVED***

			currentBatch = currentBatch[:0]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Client) configure(opts *options.ClientOptions) error ***REMOVED***
	var defaultOptions int
	// Set default options
	if opts.MaxPoolSize == nil ***REMOVED***
		defaultOptions++
		opts.SetMaxPoolSize(defaultMaxPoolSize)
	***REMOVED***
	if err := opts.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***

	var connOpts []topology.ConnectionOption
	var serverOpts []topology.ServerOption
	var topologyOpts []topology.Option

	// TODO(GODRIVER-814): Add tests for topology, server, and connection related options.

	// ServerAPIOptions need to be handled early as other client and server options below reference
	// c.serverAPI and serverOpts.serverAPI.
	if opts.ServerAPIOptions != nil ***REMOVED***
		// convert passed in options to driver form for client.
		c.serverAPI = convertToDriverAPIOptions(opts.ServerAPIOptions)

		serverOpts = append(serverOpts, topology.WithServerAPI(func(*driver.ServerAPIOptions) *driver.ServerAPIOptions ***REMOVED***
			return c.serverAPI
		***REMOVED***))
	***REMOVED***

	// ClusterClock
	c.clock = new(session.ClusterClock)

	// Pass down URI, SRV service name, and SRV max hosts so topology can poll SRV records correctly.
	topologyOpts = append(topologyOpts,
		topology.WithURI(func(uri string) string ***REMOVED*** return opts.GetURI() ***REMOVED***),
		topology.WithSRVServiceName(func(srvName string) string ***REMOVED***
			if opts.SRVServiceName != nil ***REMOVED***
				return *opts.SRVServiceName
			***REMOVED***
			return ""
		***REMOVED***),
		topology.WithSRVMaxHosts(func(srvMaxHosts int) int ***REMOVED***
			if opts.SRVMaxHosts != nil ***REMOVED***
				return *opts.SRVMaxHosts
			***REMOVED***
			return 0
		***REMOVED***),
	)

	// AppName
	var appName string
	if opts.AppName != nil ***REMOVED***
		appName = *opts.AppName

		serverOpts = append(serverOpts, topology.WithServerAppName(func(string) string ***REMOVED***
			return appName
		***REMOVED***))
	***REMOVED***
	// Compressors & ZlibLevel
	var comps []string
	if len(opts.Compressors) > 0 ***REMOVED***
		comps = opts.Compressors

		connOpts = append(connOpts, topology.WithCompressors(
			func(compressors []string) []string ***REMOVED***
				return append(compressors, comps...)
			***REMOVED***,
		))

		for _, comp := range comps ***REMOVED***
			switch comp ***REMOVED***
			case "zlib":
				connOpts = append(connOpts, topology.WithZlibLevel(func(level *int) *int ***REMOVED***
					return opts.ZlibLevel
				***REMOVED***))
			case "zstd":
				connOpts = append(connOpts, topology.WithZstdLevel(func(level *int) *int ***REMOVED***
					return opts.ZstdLevel
				***REMOVED***))
			***REMOVED***
		***REMOVED***

		serverOpts = append(serverOpts, topology.WithCompressionOptions(
			func(opts ...string) []string ***REMOVED*** return append(opts, comps...) ***REMOVED***,
		))
	***REMOVED***

	var loadBalanced bool
	if opts.LoadBalanced != nil ***REMOVED***
		loadBalanced = *opts.LoadBalanced
	***REMOVED***

	// Handshaker
	var handshaker = func(driver.Handshaker) driver.Handshaker ***REMOVED***
		return operation.NewHello().AppName(appName).Compressors(comps).ClusterClock(c.clock).
			ServerAPI(c.serverAPI).LoadBalanced(loadBalanced)
	***REMOVED***
	// Auth & Database & Password & Username
	if opts.Auth != nil ***REMOVED***
		cred := &auth.Cred***REMOVED***
			Username:    opts.Auth.Username,
			Password:    opts.Auth.Password,
			PasswordSet: opts.Auth.PasswordSet,
			Props:       opts.Auth.AuthMechanismProperties,
			Source:      opts.Auth.AuthSource,
		***REMOVED***
		mechanism := opts.Auth.AuthMechanism

		if len(cred.Source) == 0 ***REMOVED***
			switch strings.ToUpper(mechanism) ***REMOVED***
			case auth.MongoDBX509, auth.GSSAPI, auth.PLAIN:
				cred.Source = "$external"
			default:
				cred.Source = "admin"
			***REMOVED***
		***REMOVED***

		authenticator, err := auth.CreateAuthenticator(mechanism, cred)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		handshakeOpts := &auth.HandshakeOptions***REMOVED***
			AppName:       appName,
			Authenticator: authenticator,
			Compressors:   comps,
			ClusterClock:  c.clock,
			ServerAPI:     c.serverAPI,
			LoadBalanced:  loadBalanced,
		***REMOVED***
		if mechanism == "" ***REMOVED***
			// Required for SASL mechanism negotiation during handshake
			handshakeOpts.DBUser = cred.Source + "." + cred.Username
		***REMOVED***
		if opts.AuthenticateToAnything != nil && *opts.AuthenticateToAnything ***REMOVED***
			// Authenticate arbiters
			handshakeOpts.PerformAuthentication = func(serv description.Server) bool ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		handshaker = func(driver.Handshaker) driver.Handshaker ***REMOVED***
			return auth.Handshaker(nil, handshakeOpts)
		***REMOVED***
	***REMOVED***
	connOpts = append(connOpts, topology.WithHandshaker(handshaker))
	// ConnectTimeout
	if opts.ConnectTimeout != nil ***REMOVED***
		serverOpts = append(serverOpts, topology.WithHeartbeatTimeout(
			func(time.Duration) time.Duration ***REMOVED*** return *opts.ConnectTimeout ***REMOVED***,
		))
		connOpts = append(connOpts, topology.WithConnectTimeout(
			func(time.Duration) time.Duration ***REMOVED*** return *opts.ConnectTimeout ***REMOVED***,
		))
	***REMOVED***
	// Dialer
	if opts.Dialer != nil ***REMOVED***
		connOpts = append(connOpts, topology.WithDialer(
			func(topology.Dialer) topology.Dialer ***REMOVED*** return opts.Dialer ***REMOVED***,
		))
	***REMOVED***
	// Direct
	if opts.Direct != nil && *opts.Direct ***REMOVED***
		topologyOpts = append(topologyOpts, topology.WithMode(
			func(topology.MonitorMode) topology.MonitorMode ***REMOVED*** return topology.SingleMode ***REMOVED***,
		))
	***REMOVED***
	// HeartbeatInterval
	if opts.HeartbeatInterval != nil ***REMOVED***
		serverOpts = append(serverOpts, topology.WithHeartbeatInterval(
			func(time.Duration) time.Duration ***REMOVED*** return *opts.HeartbeatInterval ***REMOVED***,
		))
	***REMOVED***
	// Hosts
	hosts := []string***REMOVED***"localhost:27017"***REMOVED*** // default host
	if len(opts.Hosts) > 0 ***REMOVED***
		hosts = opts.Hosts
	***REMOVED***
	topologyOpts = append(topologyOpts, topology.WithSeedList(
		func(...string) []string ***REMOVED*** return hosts ***REMOVED***,
	))
	// LocalThreshold
	c.localThreshold = defaultLocalThreshold
	if opts.LocalThreshold != nil ***REMOVED***
		c.localThreshold = *opts.LocalThreshold
	***REMOVED***
	// MaxConIdleTime
	if opts.MaxConnIdleTime != nil ***REMOVED***
		connOpts = append(connOpts, topology.WithIdleTimeout(
			func(time.Duration) time.Duration ***REMOVED*** return *opts.MaxConnIdleTime ***REMOVED***,
		))
	***REMOVED***
	// MaxPoolSize
	if opts.MaxPoolSize != nil ***REMOVED***
		serverOpts = append(
			serverOpts,
			topology.WithMaxConnections(func(uint64) uint64 ***REMOVED*** return *opts.MaxPoolSize ***REMOVED***),
		)
	***REMOVED***
	// MinPoolSize
	if opts.MinPoolSize != nil ***REMOVED***
		serverOpts = append(
			serverOpts,
			topology.WithMinConnections(func(uint64) uint64 ***REMOVED*** return *opts.MinPoolSize ***REMOVED***),
		)
	***REMOVED***
	// MaxConnecting
	if opts.MaxConnecting != nil ***REMOVED***
		serverOpts = append(
			serverOpts,
			topology.WithMaxConnecting(func(uint64) uint64 ***REMOVED*** return *opts.MaxConnecting ***REMOVED***),
		)
	***REMOVED***
	// PoolMonitor
	if opts.PoolMonitor != nil ***REMOVED***
		serverOpts = append(
			serverOpts,
			topology.WithConnectionPoolMonitor(func(*event.PoolMonitor) *event.PoolMonitor ***REMOVED*** return opts.PoolMonitor ***REMOVED***),
		)
	***REMOVED***
	// Monitor
	if opts.Monitor != nil ***REMOVED***
		c.monitor = opts.Monitor
		connOpts = append(connOpts, topology.WithMonitor(
			func(*event.CommandMonitor) *event.CommandMonitor ***REMOVED*** return opts.Monitor ***REMOVED***,
		))
	***REMOVED***
	// ServerMonitor
	if opts.ServerMonitor != nil ***REMOVED***
		c.serverMonitor = opts.ServerMonitor
		serverOpts = append(
			serverOpts,
			topology.WithServerMonitor(func(*event.ServerMonitor) *event.ServerMonitor ***REMOVED*** return opts.ServerMonitor ***REMOVED***),
		)

		topologyOpts = append(
			topologyOpts,
			topology.WithTopologyServerMonitor(func(*event.ServerMonitor) *event.ServerMonitor ***REMOVED*** return opts.ServerMonitor ***REMOVED***),
		)
	***REMOVED***
	// ReadConcern
	c.readConcern = readconcern.New()
	if opts.ReadConcern != nil ***REMOVED***
		c.readConcern = opts.ReadConcern
	***REMOVED***
	// ReadPreference
	c.readPreference = readpref.Primary()
	if opts.ReadPreference != nil ***REMOVED***
		c.readPreference = opts.ReadPreference
	***REMOVED***
	// Registry
	c.registry = bson.DefaultRegistry
	if opts.Registry != nil ***REMOVED***
		c.registry = opts.Registry
	***REMOVED***
	// ReplicaSet
	if opts.ReplicaSet != nil ***REMOVED***
		topologyOpts = append(topologyOpts, topology.WithReplicaSetName(
			func(string) string ***REMOVED*** return *opts.ReplicaSet ***REMOVED***,
		))
	***REMOVED***
	// RetryWrites
	c.retryWrites = true // retry writes on by default
	if opts.RetryWrites != nil ***REMOVED***
		c.retryWrites = *opts.RetryWrites
	***REMOVED***
	c.retryReads = true
	if opts.RetryReads != nil ***REMOVED***
		c.retryReads = *opts.RetryReads
	***REMOVED***
	// ServerSelectionTimeout
	if opts.ServerSelectionTimeout != nil ***REMOVED***
		topologyOpts = append(topologyOpts, topology.WithServerSelectionTimeout(
			func(time.Duration) time.Duration ***REMOVED*** return *opts.ServerSelectionTimeout ***REMOVED***,
		))
	***REMOVED***
	// SocketTimeout
	if opts.SocketTimeout != nil ***REMOVED***
		connOpts = append(
			connOpts,
			topology.WithReadTimeout(func(time.Duration) time.Duration ***REMOVED*** return *opts.SocketTimeout ***REMOVED***),
			topology.WithWriteTimeout(func(time.Duration) time.Duration ***REMOVED*** return *opts.SocketTimeout ***REMOVED***),
		)
	***REMOVED***
	// Timeout
	c.timeout = opts.Timeout
	// TLSConfig
	if opts.TLSConfig != nil ***REMOVED***
		connOpts = append(connOpts, topology.WithTLSConfig(
			func(*tls.Config) *tls.Config ***REMOVED***
				return opts.TLSConfig
			***REMOVED***,
		))
	***REMOVED***
	// WriteConcern
	if opts.WriteConcern != nil ***REMOVED***
		c.writeConcern = opts.WriteConcern
	***REMOVED***
	// AutoEncryptionOptions
	if opts.AutoEncryptionOptions != nil ***REMOVED***
		if err := c.configureAutoEncryption(opts); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.cryptFLE = opts.Crypt
	***REMOVED***

	// OCSP cache
	ocspCache := ocsp.NewCache()
	connOpts = append(
		connOpts,
		topology.WithOCSPCache(func(ocsp.Cache) ocsp.Cache ***REMOVED*** return ocspCache ***REMOVED***),
	)

	// Disable communication with external OCSP responders.
	if opts.DisableOCSPEndpointCheck != nil ***REMOVED***
		connOpts = append(
			connOpts,
			topology.WithDisableOCSPEndpointCheck(func(bool) bool ***REMOVED*** return *opts.DisableOCSPEndpointCheck ***REMOVED***),
		)
	***REMOVED***

	// LoadBalanced
	if opts.LoadBalanced != nil ***REMOVED***
		topologyOpts = append(
			topologyOpts,
			topology.WithLoadBalanced(func(bool) bool ***REMOVED*** return *opts.LoadBalanced ***REMOVED***),
		)
		serverOpts = append(
			serverOpts,
			topology.WithServerLoadBalanced(func(bool) bool ***REMOVED*** return *opts.LoadBalanced ***REMOVED***),
		)
		connOpts = append(
			connOpts,
			topology.WithConnectionLoadBalanced(func(bool) bool ***REMOVED*** return *opts.LoadBalanced ***REMOVED***),
		)
	***REMOVED***

	serverOpts = append(
		serverOpts,
		topology.WithClock(func(*session.ClusterClock) *session.ClusterClock ***REMOVED*** return c.clock ***REMOVED***),
		topology.WithConnectionOptions(func(...topology.ConnectionOption) []topology.ConnectionOption ***REMOVED*** return connOpts ***REMOVED***),
	)
	topologyOpts = append(topologyOpts, topology.WithServerOptions(
		func(...topology.ServerOption) []topology.ServerOption ***REMOVED*** return serverOpts ***REMOVED***,
	))
	c.topologyOptions = topologyOpts

	// Deployment
	if opts.Deployment != nil ***REMOVED***
		// topology options: WithSeedlist, WithURI, WithSRVServiceName, WithSRVMaxHosts, and WithServerOptions
		// server options: WithClock and WithConnectionOptions + default maxPoolSize
		if len(serverOpts) > 2+defaultOptions || len(topologyOpts) > 5 ***REMOVED***
			return errors.New("cannot specify topology or server options with a deployment")
		***REMOVED***
		c.deployment = opts.Deployment
	***REMOVED***

	return nil
***REMOVED***

func (c *Client) configureAutoEncryption(clientOpts *options.ClientOptions) error ***REMOVED***
	c.encryptedFieldsMap = clientOpts.AutoEncryptionOptions.EncryptedFieldsMap
	if err := c.configureKeyVaultClientFLE(clientOpts); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := c.configureMetadataClientFLE(clientOpts); err != nil ***REMOVED***
		return err
	***REMOVED***

	mc, err := c.newMongoCrypt(clientOpts.AutoEncryptionOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the crypt_shared library was loaded successfully, signal to the mongocryptd client creator
	// that it can bypass spawning mongocryptd.
	cryptSharedLibAvailable := mc.CryptSharedLibVersionString() != ""
	mongocryptdFLE, err := newMongocryptdClient(cryptSharedLibAvailable, clientOpts.AutoEncryptionOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.mongocryptdFLE = mongocryptdFLE

	c.configureCryptFLE(mc, clientOpts.AutoEncryptionOptions)
	return nil
***REMOVED***

func (c *Client) getOrCreateInternalClient(clientOpts *options.ClientOptions) (*Client, error) ***REMOVED***
	if c.internalClientFLE != nil ***REMOVED***
		return c.internalClientFLE, nil
	***REMOVED***

	internalClientOpts := options.MergeClientOptions(clientOpts)
	internalClientOpts.AutoEncryptionOptions = nil
	internalClientOpts.SetMinPoolSize(0)
	var err error
	c.internalClientFLE, err = NewClient(internalClientOpts)
	return c.internalClientFLE, err
***REMOVED***

func (c *Client) configureKeyVaultClientFLE(clientOpts *options.ClientOptions) error ***REMOVED***
	// parse key vault options and create new key vault client
	var err error
	aeOpts := clientOpts.AutoEncryptionOptions
	switch ***REMOVED***
	case aeOpts.KeyVaultClientOptions != nil:
		c.keyVaultClientFLE, err = NewClient(aeOpts.KeyVaultClientOptions)
	case clientOpts.MaxPoolSize != nil && *clientOpts.MaxPoolSize == 0:
		c.keyVaultClientFLE = c
	default:
		c.keyVaultClientFLE, err = c.getOrCreateInternalClient(clientOpts)
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	dbName, collName := splitNamespace(aeOpts.KeyVaultNamespace)
	c.keyVaultCollFLE = c.keyVaultClientFLE.Database(dbName).Collection(collName, keyVaultCollOpts)
	return nil
***REMOVED***

func (c *Client) configureMetadataClientFLE(clientOpts *options.ClientOptions) error ***REMOVED***
	// parse key vault options and create new key vault client
	aeOpts := clientOpts.AutoEncryptionOptions
	if aeOpts.BypassAutoEncryption != nil && *aeOpts.BypassAutoEncryption ***REMOVED***
		// no need for a metadata client.
		return nil
	***REMOVED***
	if clientOpts.MaxPoolSize != nil && *clientOpts.MaxPoolSize == 0 ***REMOVED***
		c.metadataClientFLE = c
		return nil
	***REMOVED***

	var err error
	c.metadataClientFLE, err = c.getOrCreateInternalClient(clientOpts)
	return err
***REMOVED***

func (c *Client) newMongoCrypt(opts *options.AutoEncryptionOptions) (*mongocrypt.MongoCrypt, error) ***REMOVED***
	// convert schemas in SchemaMap to bsoncore documents
	cryptSchemaMap := make(map[string]bsoncore.Document)
	for k, v := range opts.SchemaMap ***REMOVED***
		schema, err := transformBsoncoreDocument(c.registry, v, true, "schemaMap")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cryptSchemaMap[k] = schema
	***REMOVED***

	// convert schemas in EncryptedFieldsMap to bsoncore documents
	cryptEncryptedFieldsMap := make(map[string]bsoncore.Document)
	for k, v := range opts.EncryptedFieldsMap ***REMOVED***
		encryptedFields, err := transformBsoncoreDocument(c.registry, v, true, "encryptedFieldsMap")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cryptEncryptedFieldsMap[k] = encryptedFields
	***REMOVED***

	kmsProviders, err := transformBsoncoreDocument(c.registry, opts.KmsProviders, true, "kmsProviders")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error creating KMS providers document: %v", err)
	***REMOVED***

	// Set the crypt_shared library override path from the "cryptSharedLibPath" extra option if one
	// was set.
	cryptSharedLibPath := ""
	if val, ok := opts.ExtraOptions["cryptSharedLibPath"]; ok ***REMOVED***
		str, ok := val.(string)
		if !ok ***REMOVED***
			return nil, fmt.Errorf(
				`expected AutoEncryption extra option "cryptSharedLibPath" to be a string, but is a %T`, val)
		***REMOVED***
		cryptSharedLibPath = str
	***REMOVED***

	// Explicitly disable loading the crypt_shared library if requested. Note that this is ONLY
	// intended for use from tests; there is no supported public API for explicitly disabling
	// loading the crypt_shared library.
	cryptSharedLibDisabled := false
	if v, ok := opts.ExtraOptions["__cryptSharedLibDisabledForTestOnly"]; ok ***REMOVED***
		cryptSharedLibDisabled = v.(bool)
	***REMOVED***

	bypassAutoEncryption := opts.BypassAutoEncryption != nil && *opts.BypassAutoEncryption
	bypassQueryAnalysis := opts.BypassQueryAnalysis != nil && *opts.BypassQueryAnalysis

	mc, err := mongocrypt.NewMongoCrypt(mcopts.MongoCrypt().
		SetKmsProviders(kmsProviders).
		SetLocalSchemaMap(cryptSchemaMap).
		SetBypassQueryAnalysis(bypassQueryAnalysis).
		SetEncryptedFieldsMap(cryptEncryptedFieldsMap).
		SetCryptSharedLibDisabled(cryptSharedLibDisabled || bypassAutoEncryption).
		SetCryptSharedLibOverridePath(cryptSharedLibPath))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var cryptSharedLibRequired bool
	if val, ok := opts.ExtraOptions["cryptSharedLibRequired"]; ok ***REMOVED***
		b, ok := val.(bool)
		if !ok ***REMOVED***
			return nil, fmt.Errorf(
				`expected AutoEncryption extra option "cryptSharedLibRequired" to be a bool, but is a %T`, val)
		***REMOVED***
		cryptSharedLibRequired = b
	***REMOVED***

	// If the "cryptSharedLibRequired" extra option is set to true, check the MongoCrypt version
	// string to confirm that the library was successfully loaded. If the version string is empty,
	// return an error indicating that we couldn't load the crypt_shared library.
	if cryptSharedLibRequired && mc.CryptSharedLibVersionString() == "" ***REMOVED***
		return nil, errors.New(
			`AutoEncryption extra option "cryptSharedLibRequired" is true, but we failed to load the crypt_shared library`)
	***REMOVED***

	return mc, nil
***REMOVED***

//nolint:unused // the unused linter thinks that this function is unreachable because "c.newMongoCrypt" always panics without the "cse" build tag set.
func (c *Client) configureCryptFLE(mc *mongocrypt.MongoCrypt, opts *options.AutoEncryptionOptions) ***REMOVED***
	bypass := opts.BypassAutoEncryption != nil && *opts.BypassAutoEncryption
	kr := keyRetriever***REMOVED***coll: c.keyVaultCollFLE***REMOVED***
	var cir collInfoRetriever
	// If bypass is true, c.metadataClientFLE is nil and the collInfoRetriever
	// will not be used. If bypass is false, to the parent client or the
	// internal client.
	if !bypass ***REMOVED***
		cir = collInfoRetriever***REMOVED***client: c.metadataClientFLE***REMOVED***
	***REMOVED***

	c.cryptFLE = driver.NewCrypt(&driver.CryptOptions***REMOVED***
		MongoCrypt:           mc,
		CollInfoFn:           cir.cryptCollInfo,
		KeyFn:                kr.cryptKeys,
		MarkFn:               c.mongocryptdFLE.markCommand,
		TLSConfig:            opts.TLSConfig,
		BypassAutoEncryption: bypass,
	***REMOVED***)
***REMOVED***

// validSession returns an error if the session doesn't belong to the client
func (c *Client) validSession(sess *session.Client) error ***REMOVED***
	if sess != nil && sess.ClientID != c.id ***REMOVED***
		return ErrWrongClient
	***REMOVED***
	return nil
***REMOVED***

// convertToDriverAPIOptions converts a options.ServerAPIOptions instance to a driver.ServerAPIOptions.
func convertToDriverAPIOptions(s *options.ServerAPIOptions) *driver.ServerAPIOptions ***REMOVED***
	driverOpts := driver.NewServerAPIOptions(string(s.ServerAPIVersion))
	if s.Strict != nil ***REMOVED***
		driverOpts.SetStrict(*s.Strict)
	***REMOVED***
	if s.DeprecationErrors != nil ***REMOVED***
		driverOpts.SetDeprecationErrors(*s.DeprecationErrors)
	***REMOVED***
	return driverOpts
***REMOVED***

// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
func (c *Client) Database(name string, opts ...*options.DatabaseOptions) *Database ***REMOVED***
	return newDatabase(c, name, opts...)
***REMOVED***

// ListDatabases executes a listDatabases command and returns the result.
//
// The filter parameter must be a document containing query operators and can be used to select which
// databases are included in the result. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to include
// all databases.
//
// The opts parameter can be used to specify options for this operation (see the options.ListDatabasesOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listDatabases/.
func (c *Client) ListDatabases(ctx context.Context, filter interface***REMOVED******REMOVED***, opts ...*options.ListDatabasesOptions) (ListDatabasesResult, error) ***REMOVED***
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	sess := sessionFromContext(ctx)

	err := c.validSession(sess)
	if err != nil ***REMOVED***
		return ListDatabasesResult***REMOVED******REMOVED***, err
	***REMOVED***
	if sess == nil && c.sessionPool != nil ***REMOVED***
		sess, err = session.NewClientSession(c.sessionPool, c.id, session.Implicit)
		if err != nil ***REMOVED***
			return ListDatabasesResult***REMOVED******REMOVED***, err
		***REMOVED***
		defer sess.EndSession()
	***REMOVED***

	err = c.validSession(sess)
	if err != nil ***REMOVED***
		return ListDatabasesResult***REMOVED******REMOVED***, err
	***REMOVED***

	filterDoc, err := transformBsoncoreDocument(c.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return ListDatabasesResult***REMOVED******REMOVED***, err
	***REMOVED***

	selector := description.CompositeSelector([]description.ServerSelector***REMOVED***
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(c.localThreshold),
	***REMOVED***)
	selector = makeReadPrefSelector(sess, selector, c.localThreshold)

	ldo := options.MergeListDatabasesOptions(opts...)
	op := operation.NewListDatabases(filterDoc).
		Session(sess).ReadPreference(c.readPreference).CommandMonitor(c.monitor).
		ServerSelector(selector).ClusterClock(c.clock).Database("admin").Deployment(c.deployment).Crypt(c.cryptFLE).
		ServerAPI(c.serverAPI).Timeout(c.timeout)

	if ldo.NameOnly != nil ***REMOVED***
		op = op.NameOnly(*ldo.NameOnly)
	***REMOVED***
	if ldo.AuthorizedDatabases != nil ***REMOVED***
		op = op.AuthorizedDatabases(*ldo.AuthorizedDatabases)
	***REMOVED***

	retry := driver.RetryNone
	if c.retryReads ***REMOVED***
		retry = driver.RetryOncePerCommand
	***REMOVED***
	op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil ***REMOVED***
		return ListDatabasesResult***REMOVED******REMOVED***, replaceErrors(err)
	***REMOVED***

	return newListDatabasesResultFromOperation(op.Result()), nil
***REMOVED***

// ListDatabaseNames executes a listDatabases command and returns a slice containing the names of all of the databases
// on the server.
//
// The filter parameter must be a document containing query operators and can be used to select which databases
// are included in the result. It cannot be nil. An empty document (e.g. bson.D***REMOVED******REMOVED***) should be used to include all
// databases.
//
// The opts parameter can be used to specify options for this operation (see the options.ListDatabasesOptions
// documentation.)
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listDatabases/.
func (c *Client) ListDatabaseNames(ctx context.Context, filter interface***REMOVED******REMOVED***, opts ...*options.ListDatabasesOptions) ([]string, error) ***REMOVED***
	opts = append(opts, options.ListDatabases().SetNameOnly(true))

	res, err := c.ListDatabases(ctx, filter, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	names := make([]string, 0)
	for _, spec := range res.Databases ***REMOVED***
		names = append(names, spec.Name)
	***REMOVED***

	return names, nil
***REMOVED***

// WithSession creates a new SessionContext from the ctx and sess parameters and uses it to call the fn callback. The
// SessionContext must be used as the Context parameter for any operations in the fn callback that should be executed
// under the session.
//
// WithSession is safe to call from multiple goroutines concurrently. However, the SessionContext passed to the
// WithSession callback function is not safe for concurrent use by multiple goroutines.
//
// If the ctx parameter already contains a Session, that Session will be replaced with the one provided.
//
// Any error returned by the fn callback will be returned without any modifications.
func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error ***REMOVED***
	return fn(NewSessionContext(ctx, sess))
***REMOVED***

// UseSession creates a new Session and uses it to create a new SessionContext, which is used to call the fn callback.
// The SessionContext parameter must be used as the Context parameter for any operations in the fn callback that should
// be executed under a session. After the callback returns, the created Session is ended, meaning that any in-progress
// transactions started by fn will be aborted even if fn returns an error.
//
// UseSession is safe to call from multiple goroutines concurrently. However, the SessionContext passed to the
// UseSession callback function is not safe for concurrent use by multiple goroutines.
//
// If the ctx parameter already contains a Session, that Session will be replaced with the newly created one.
//
// Any error returned by the fn callback will be returned without any modifications.
func (c *Client) UseSession(ctx context.Context, fn func(SessionContext) error) error ***REMOVED***
	return c.UseSessionWithOptions(ctx, options.Session(), fn)
***REMOVED***

// UseSessionWithOptions operates like UseSession but uses the given SessionOptions to create the Session.
//
// UseSessionWithOptions is safe to call from multiple goroutines concurrently. However, the SessionContext passed to
// the UseSessionWithOptions callback function is not safe for concurrent use by multiple goroutines.
func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(SessionContext) error) error ***REMOVED***
	defaultSess, err := c.StartSession(opts)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer defaultSess.EndSession(ctx)
	return fn(NewSessionContext(ctx, defaultSess))
***REMOVED***

// Watch returns a change stream for all changes on the deployment. See
// https://www.mongodb.com/docs/manual/changeStreams/ for more information about change streams.
//
// The client must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
// nil or empty. The stage documents must all be non-nil. See https://www.mongodb.com/docs/manual/changeStreams/ for a list
// of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the mongo.Pipeline***REMOVED******REMOVED***
// type can be used.
//
// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
// documentation).
func (c *Client) Watch(ctx context.Context, pipeline interface***REMOVED******REMOVED***,
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) ***REMOVED***
	if c.sessionPool == nil ***REMOVED***
		return nil, ErrClientDisconnected
	***REMOVED***

	csConfig := changeStreamConfig***REMOVED***
		readConcern:    c.readConcern,
		readPreference: c.readPreference,
		client:         c,
		registry:       c.registry,
		streamType:     ClientStream,
		crypt:          c.cryptFLE,
	***REMOVED***

	return newChangeStream(ctx, csConfig, pipeline, opts...)
***REMOVED***

// NumberSessionsInProgress returns the number of sessions that have been started for this client but have not been
// closed (i.e. EndSession has not been called).
func (c *Client) NumberSessionsInProgress() int ***REMOVED***
	return c.sessionPool.CheckedOut()
***REMOVED***

func (c *Client) createBaseCursorOptions() driver.CursorOptions ***REMOVED***
	return driver.CursorOptions***REMOVED***
		CommandMonitor: c.monitor,
		Crypt:          c.cryptFLE,
		ServerAPI:      c.serverAPI,
	***REMOVED***
***REMOVED***
