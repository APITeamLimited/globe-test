// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var defaultRegistry = bson.NewRegistryBuilder().Build()

type serverConfig struct ***REMOVED***
	clock              *session.ClusterClock
	compressionOpts    []string
	connectionOpts     []ConnectionOption
	appname            string
	heartbeatInterval  time.Duration
	heartbeatTimeout   time.Duration
	serverMonitor      *event.ServerMonitor
	registry           *bsoncodec.Registry
	monitoringDisabled bool
	serverAPI          *driver.ServerAPIOptions
	loadBalanced       bool

	// Connection pool options.
	maxConns             uint64
	minConns             uint64
	maxConnecting        uint64
	poolMonitor          *event.PoolMonitor
	poolMaxIdleTime      time.Duration
	poolMaintainInterval time.Duration
***REMOVED***

func newServerConfig(opts ...ServerOption) *serverConfig ***REMOVED***
	cfg := &serverConfig***REMOVED***
		heartbeatInterval: 10 * time.Second,
		heartbeatTimeout:  10 * time.Second,
		registry:          defaultRegistry,
	***REMOVED***

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		opt(cfg)
	***REMOVED***

	return cfg
***REMOVED***

// ServerOption configures a server.
type ServerOption func(*serverConfig)

func withMonitoringDisabled(fn func(bool) bool) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.monitoringDisabled = fn(cfg.monitoringDisabled)
	***REMOVED***
***REMOVED***

// WithConnectionOptions configures the server's connections.
func WithConnectionOptions(fn func(...ConnectionOption) []ConnectionOption) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.connectionOpts = fn(cfg.connectionOpts...)
	***REMOVED***
***REMOVED***

// WithCompressionOptions configures the server's compressors.
func WithCompressionOptions(fn func(...string) []string) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.compressionOpts = fn(cfg.compressionOpts...)
	***REMOVED***
***REMOVED***

// WithServerAppName configures the server's application name.
func WithServerAppName(fn func(string) string) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.appname = fn(cfg.appname)
	***REMOVED***
***REMOVED***

// WithHeartbeatInterval configures a server's heartbeat interval.
func WithHeartbeatInterval(fn func(time.Duration) time.Duration) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.heartbeatInterval = fn(cfg.heartbeatInterval)
	***REMOVED***
***REMOVED***

// WithHeartbeatTimeout configures how long to wait for a heartbeat socket to
// connection.
func WithHeartbeatTimeout(fn func(time.Duration) time.Duration) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.heartbeatTimeout = fn(cfg.heartbeatTimeout)
	***REMOVED***
***REMOVED***

// WithMaxConnections configures the maximum number of connections to allow for
// a given server. If max is 0, then maximum connection pool size is not limited.
func WithMaxConnections(fn func(uint64) uint64) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.maxConns = fn(cfg.maxConns)
	***REMOVED***
***REMOVED***

// WithMinConnections configures the minimum number of connections to allow for
// a given server. If min is 0, then there is no lower limit to the number of
// connections.
func WithMinConnections(fn func(uint64) uint64) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.minConns = fn(cfg.minConns)
	***REMOVED***
***REMOVED***

// WithMaxConnecting configures the maximum number of connections a connection
// pool may establish simultaneously. If maxConnecting is 0, the default value
// of 2 is used.
func WithMaxConnecting(fn func(uint64) uint64) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.maxConnecting = fn(cfg.maxConnecting)
	***REMOVED***
***REMOVED***

// WithConnectionPoolMaxIdleTime configures the maximum time that a connection can remain idle in the connection pool
// before being removed. If connectionPoolMaxIdleTime is 0, then no idle time is set and connections will not be removed
// because of their age
func WithConnectionPoolMaxIdleTime(fn func(time.Duration) time.Duration) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.poolMaxIdleTime = fn(cfg.poolMaxIdleTime)
	***REMOVED***
***REMOVED***

// WithConnectionPoolMaintainInterval configures the interval that the background connection pool
// maintenance goroutine runs.
func WithConnectionPoolMaintainInterval(fn func(time.Duration) time.Duration) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.poolMaintainInterval = fn(cfg.poolMaintainInterval)
	***REMOVED***
***REMOVED***

// WithConnectionPoolMonitor configures the monitor for all connection pool actions
func WithConnectionPoolMonitor(fn func(*event.PoolMonitor) *event.PoolMonitor) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.poolMonitor = fn(cfg.poolMonitor)
	***REMOVED***
***REMOVED***

// WithServerMonitor configures the monitor for all SDAM events for a server
func WithServerMonitor(fn func(*event.ServerMonitor) *event.ServerMonitor) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.serverMonitor = fn(cfg.serverMonitor)
	***REMOVED***
***REMOVED***

// WithClock configures the ClusterClock for the server to use.
func WithClock(fn func(clock *session.ClusterClock) *session.ClusterClock) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.clock = fn(cfg.clock)
	***REMOVED***
***REMOVED***

// WithRegistry configures the registry for the server to use when creating
// cursors.
func WithRegistry(fn func(*bsoncodec.Registry) *bsoncodec.Registry) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.registry = fn(cfg.registry)
	***REMOVED***
***REMOVED***

// WithServerAPI configures the server API options for the server to use.
func WithServerAPI(fn func(serverAPI *driver.ServerAPIOptions) *driver.ServerAPIOptions) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.serverAPI = fn(cfg.serverAPI)
	***REMOVED***
***REMOVED***

// WithServerLoadBalanced specifies whether or not the server is behind a load balancer.
func WithServerLoadBalanced(fn func(bool) bool) ServerOption ***REMOVED***
	return func(cfg *serverConfig) ***REMOVED***
		cfg.loadBalanced = fn(cfg.loadBalanced)
	***REMOVED***
***REMOVED***
