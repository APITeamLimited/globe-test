// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/ocsp"
)

// Dialer is used to make network connections.
type Dialer interface ***REMOVED***
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
***REMOVED***

// DialerFunc is a type implemented by functions that can be used as a Dialer.
type DialerFunc func(ctx context.Context, network, address string) (net.Conn, error)

// DialContext implements the Dialer interface.
func (df DialerFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) ***REMOVED***
	return df(ctx, network, address)
***REMOVED***

// DefaultDialer is the Dialer implementation that is used by this package. Changing this
// will also change the Dialer used for this package. This should only be changed why all
// of the connections being made need to use a different Dialer. Most of the time, using a
// WithDialer option is more appropriate than changing this variable.
var DefaultDialer Dialer = &net.Dialer***REMOVED******REMOVED***

// Handshaker is the interface implemented by types that can perform a MongoDB
// handshake over a provided driver.Connection. This is used during connection
// initialization. Implementations must be goroutine safe.
type Handshaker = driver.Handshaker

// generationNumberFn is a callback type used by a connection to fetch its generation number given its service ID.
type generationNumberFn func(serviceID *primitive.ObjectID) uint64

type connectionConfig struct ***REMOVED***
	connectTimeout           time.Duration
	dialer                   Dialer
	handshaker               Handshaker
	idleTimeout              time.Duration
	cmdMonitor               *event.CommandMonitor
	readTimeout              time.Duration
	writeTimeout             time.Duration
	tlsConfig                *tls.Config
	compressors              []string
	zlibLevel                *int
	zstdLevel                *int
	ocspCache                ocsp.Cache
	disableOCSPEndpointCheck bool
	tlsConnectionSource      tlsConnectionSource
	loadBalanced             bool
	getGenerationFn          generationNumberFn
***REMOVED***

func newConnectionConfig(opts ...ConnectionOption) *connectionConfig ***REMOVED***
	cfg := &connectionConfig***REMOVED***
		connectTimeout:      30 * time.Second,
		dialer:              nil,
		tlsConnectionSource: defaultTLSConnectionSource,
	***REMOVED***

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		opt(cfg)
	***REMOVED***

	if cfg.dialer == nil ***REMOVED***
		cfg.dialer = &net.Dialer***REMOVED******REMOVED***
	***REMOVED***

	return cfg
***REMOVED***

// ConnectionOption is used to configure a connection.
type ConnectionOption func(*connectionConfig)

func withTLSConnectionSource(fn func(tlsConnectionSource) tlsConnectionSource) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.tlsConnectionSource = fn(c.tlsConnectionSource)
	***REMOVED***
***REMOVED***

// WithCompressors sets the compressors that can be used for communication.
func WithCompressors(fn func([]string) []string) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.compressors = fn(c.compressors)
	***REMOVED***
***REMOVED***

// WithConnectTimeout configures the maximum amount of time a dial will wait for a
// Connect to complete. The default is 30 seconds.
func WithConnectTimeout(fn func(time.Duration) time.Duration) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.connectTimeout = fn(c.connectTimeout)
	***REMOVED***
***REMOVED***

// WithDialer configures the Dialer to use when making a new connection to MongoDB.
func WithDialer(fn func(Dialer) Dialer) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.dialer = fn(c.dialer)
	***REMOVED***
***REMOVED***

// WithHandshaker configures the Handshaker that wll be used to initialize newly
// dialed connections.
func WithHandshaker(fn func(Handshaker) Handshaker) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.handshaker = fn(c.handshaker)
	***REMOVED***
***REMOVED***

// WithIdleTimeout configures the maximum idle time to allow for a connection.
func WithIdleTimeout(fn func(time.Duration) time.Duration) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.idleTimeout = fn(c.idleTimeout)
	***REMOVED***
***REMOVED***

// WithReadTimeout configures the maximum read time for a connection.
func WithReadTimeout(fn func(time.Duration) time.Duration) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.readTimeout = fn(c.readTimeout)
	***REMOVED***
***REMOVED***

// WithWriteTimeout configures the maximum write time for a connection.
func WithWriteTimeout(fn func(time.Duration) time.Duration) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.writeTimeout = fn(c.writeTimeout)
	***REMOVED***
***REMOVED***

// WithTLSConfig configures the TLS options for a connection.
func WithTLSConfig(fn func(*tls.Config) *tls.Config) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.tlsConfig = fn(c.tlsConfig)
	***REMOVED***
***REMOVED***

// WithMonitor configures a event for command monitoring.
func WithMonitor(fn func(*event.CommandMonitor) *event.CommandMonitor) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.cmdMonitor = fn(c.cmdMonitor)
	***REMOVED***
***REMOVED***

// WithZlibLevel sets the zLib compression level.
func WithZlibLevel(fn func(*int) *int) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.zlibLevel = fn(c.zlibLevel)
	***REMOVED***
***REMOVED***

// WithZstdLevel sets the zstd compression level.
func WithZstdLevel(fn func(*int) *int) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.zstdLevel = fn(c.zstdLevel)
	***REMOVED***
***REMOVED***

// WithOCSPCache specifies a cache to use for OCSP verification.
func WithOCSPCache(fn func(ocsp.Cache) ocsp.Cache) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.ocspCache = fn(c.ocspCache)
	***REMOVED***
***REMOVED***

// WithDisableOCSPEndpointCheck specifies whether or the driver should perform non-stapled OCSP verification. If set
// to true, the driver will only check stapled responses and will continue the connection without reaching out to
// OCSP responders.
func WithDisableOCSPEndpointCheck(fn func(bool) bool) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.disableOCSPEndpointCheck = fn(c.disableOCSPEndpointCheck)
	***REMOVED***
***REMOVED***

// WithConnectionLoadBalanced specifies whether or not the connection is to a server behind a load balancer.
func WithConnectionLoadBalanced(fn func(bool) bool) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.loadBalanced = fn(c.loadBalanced)
	***REMOVED***
***REMOVED***

func withGenerationNumberFn(fn func(generationNumberFn) generationNumberFn) ConnectionOption ***REMOVED***
	return func(c *connectionConfig) ***REMOVED***
		c.getGenerationFn = fn(c.getGenerationFn)
	***REMOVED***
***REMOVED***
