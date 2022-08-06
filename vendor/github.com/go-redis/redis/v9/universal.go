package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

// UniversalOptions information is required by UniversalClient to establish
// connections.
type UniversalOptions struct ***REMOVED***
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Addrs []string

	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int

	// Common options.

	Dialer    func(ctx context.Context, network, addr string) (net.Conn, error)
	OnConnect func(ctx context.Context, cn *Conn) error

	Username         string
	Password         string
	SentinelUsername string
	SentinelPassword string

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool

	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration

	TLSConfig *tls.Config

	// Only cluster clients.

	MaxRedirects   int
	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool

	// The sentinel master name.
	// Only failover clients.

	MasterName string
***REMOVED***

// Cluster returns cluster options created from the universal options.
func (o *UniversalOptions) Cluster() *ClusterOptions ***REMOVED***
	if len(o.Addrs) == 0 ***REMOVED***
		o.Addrs = []string***REMOVED***"127.0.0.1:6379"***REMOVED***
	***REMOVED***

	return &ClusterOptions***REMOVED***
		Addrs:     o.Addrs,
		Dialer:    o.Dialer,
		OnConnect: o.OnConnect,

		Username: o.Username,
		Password: o.Password,

		MaxRedirects:   o.MaxRedirects,
		ReadOnly:       o.ReadOnly,
		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolFIFO: o.PoolFIFO,

		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		MinIdleConns:    o.MinIdleConns,
		MaxIdleConns:    o.MaxIdleConns,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		ConnMaxLifetime: o.ConnMaxLifetime,

		TLSConfig: o.TLSConfig,
	***REMOVED***
***REMOVED***

// Failover returns failover options created from the universal options.
func (o *UniversalOptions) Failover() *FailoverOptions ***REMOVED***
	if len(o.Addrs) == 0 ***REMOVED***
		o.Addrs = []string***REMOVED***"127.0.0.1:26379"***REMOVED***
	***REMOVED***

	return &FailoverOptions***REMOVED***
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,

		Dialer:    o.Dialer,
		OnConnect: o.OnConnect,

		DB:               o.DB,
		Username:         o.Username,
		Password:         o.Password,
		SentinelUsername: o.SentinelUsername,
		SentinelPassword: o.SentinelPassword,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolFIFO:        o.PoolFIFO,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		MinIdleConns:    o.MinIdleConns,
		MaxIdleConns:    o.MaxIdleConns,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		ConnMaxLifetime: o.ConnMaxLifetime,

		TLSConfig: o.TLSConfig,
	***REMOVED***
***REMOVED***

// Simple returns basic options created from the universal options.
func (o *UniversalOptions) Simple() *Options ***REMOVED***
	addr := "127.0.0.1:6379"
	if len(o.Addrs) > 0 ***REMOVED***
		addr = o.Addrs[0]
	***REMOVED***

	return &Options***REMOVED***
		Addr:      addr,
		Dialer:    o.Dialer,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Username: o.Username,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolFIFO:        o.PoolFIFO,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		MinIdleConns:    o.MinIdleConns,
		MaxIdleConns:    o.MaxIdleConns,
		ConnMaxIdleTime: o.ConnMaxIdleTime,
		ConnMaxLifetime: o.ConnMaxLifetime,

		TLSConfig: o.TLSConfig,
	***REMOVED***
***REMOVED***

// --------------------------------------------------------------------

// UniversalClient is an abstract client which - based on the provided options -
// represents either a ClusterClient, a FailoverClient, or a single-node Client.
// This can be useful for testing cluster-specific applications locally or having different
// clients in different environments.
type UniversalClient interface ***REMOVED***
	Cmdable
	AddHook(Hook)
	Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error
	Do(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd
	Process(ctx context.Context, cmd Cmder) error
	Subscribe(ctx context.Context, channels ...string) *PubSub
	PSubscribe(ctx context.Context, channels ...string) *PubSub
	Close() error
	PoolStats() *PoolStats
***REMOVED***

var (
	_ UniversalClient = (*Client)(nil)
	_ UniversalClient = (*ClusterClient)(nil)
	_ UniversalClient = (*Ring)(nil)
)

// NewUniversalClient returns a new multi client. The type of the returned client depends
// on the following conditions:
//
// 1. If the MasterName option is specified, a sentinel-backed FailoverClient is returned.
// 2. if the number of Addrs is two or more, a ClusterClient is returned.
// 3. Otherwise, a single-node Client is returned.
func NewUniversalClient(opts *UniversalOptions) UniversalClient ***REMOVED***
	if opts.MasterName != "" ***REMOVED***
		return NewFailoverClient(opts.Failover())
	***REMOVED*** else if len(opts.Addrs) > 1 ***REMOVED***
		return NewClusterClient(opts.Cluster())
	***REMOVED***
	return NewClient(opts.Simple())
***REMOVED***
