// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options // import "go.mongodb.org/mongo-driver/mongo/options"

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/youmark/pkcs8"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/tag"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// ContextDialer is an interface that can be implemented by types that can create connections. It should be used to
// provide a custom dialer when configuring a Client.
//
// DialContext should return a connection to the provided address on the given network.
type ContextDialer interface ***REMOVED***
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
***REMOVED***

// Credential can be used to provide authentication options when configuring a Client.
//
// AuthMechanism: the mechanism to use for authentication. Supported values include "SCRAM-SHA-256", "SCRAM-SHA-1",
// "MONGODB-CR", "PLAIN", "GSSAPI", "MONGODB-X509", and "MONGODB-AWS". This can also be set through the "authMechanism"
// URI option. (e.g. "authMechanism=PLAIN"). For more information, see
// https://www.mongodb.com/docs/manual/core/authentication-mechanisms/.
//
// AuthMechanismProperties can be used to specify additional configuration options for certain mechanisms. They can also
// be set through the "authMechanismProperites" URI option
// (e.g. "authMechanismProperties=SERVICE_NAME:service,CANONICALIZE_HOST_NAME:true"). Supported properties are:
//
// 1. SERVICE_NAME: The service name to use for GSSAPI authentication. The default is "mongodb".
//
// 2. CANONICALIZE_HOST_NAME: If "true", the driver will canonicalize the host name for GSSAPI authentication. The default
// is "false".
//
// 3. SERVICE_REALM: The service realm for GSSAPI authentication.
//
// 4. SERVICE_HOST: The host name to use for GSSAPI authentication. This should be specified if the host name to use for
// authentication is different than the one given for Client construction.
//
// 4. AWS_SESSION_TOKEN: The AWS token for MONGODB-AWS authentication. This is optional and used for authentication with
// temporary credentials.
//
// The SERVICE_HOST and CANONICALIZE_HOST_NAME properties must not be used at the same time on Linux and Darwin
// systems.
//
// AuthSource: the name of the database to use for authentication. This defaults to "$external" for MONGODB-X509,
// GSSAPI, and PLAIN and "admin" for all other mechanisms. This can also be set through the "authSource" URI option
// (e.g. "authSource=otherDb").
//
// Username: the username for authentication. This can also be set through the URI as a username:password pair before
// the first @ character. For example, a URI for user "user", password "pwd", and host "localhost:27017" would be
// "mongodb://user:pwd@localhost:27017". This is optional for X509 authentication and will be extracted from the
// client certificate if not specified.
//
// Password: the password for authentication. This must not be specified for X509 and is optional for GSSAPI
// authentication.
//
// PasswordSet: For GSSAPI, this must be true if a password is specified, even if the password is the empty string, and
// false if no password is specified, indicating that the password should be taken from the context of the running
// process. For other mechanisms, this field is ignored.
type Credential struct ***REMOVED***
	AuthMechanism           string
	AuthMechanismProperties map[string]string
	AuthSource              string
	Username                string
	Password                string
	PasswordSet             bool
***REMOVED***

// ClientOptions contains options to configure a Client instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ClientOptions struct ***REMOVED***
	AppName                  *string
	Auth                     *Credential
	AutoEncryptionOptions    *AutoEncryptionOptions
	ConnectTimeout           *time.Duration
	Compressors              []string
	Dialer                   ContextDialer
	Direct                   *bool
	DisableOCSPEndpointCheck *bool
	HeartbeatInterval        *time.Duration
	Hosts                    []string
	LoadBalanced             *bool
	LocalThreshold           *time.Duration
	MaxConnIdleTime          *time.Duration
	MaxPoolSize              *uint64
	MinPoolSize              *uint64
	MaxConnecting            *uint64
	PoolMonitor              *event.PoolMonitor
	Monitor                  *event.CommandMonitor
	ServerMonitor            *event.ServerMonitor
	ReadConcern              *readconcern.ReadConcern
	ReadPreference           *readpref.ReadPref
	Registry                 *bsoncodec.Registry
	ReplicaSet               *string
	RetryReads               *bool
	RetryWrites              *bool
	ServerAPIOptions         *ServerAPIOptions
	ServerSelectionTimeout   *time.Duration
	SRVMaxHosts              *int
	SRVServiceName           *string
	Timeout                  *time.Duration
	TLSConfig                *tls.Config
	WriteConcern             *writeconcern.WriteConcern
	ZlibLevel                *int
	ZstdLevel                *int

	err error
	uri string
	cs  *connstring.ConnString

	// AuthenticateToAnything skips server type checks when deciding if authentication is possible.
	//
	// Deprecated: This option is for internal use only and should not be set. It may be changed or removed in any
	// release.
	AuthenticateToAnything *bool

	// Crypt specifies a custom driver.Crypt to be used to encrypt and decrypt documents. The default is no
	// encryption.
	//
	// Deprecated: This option is for internal use only and should not be set (see GODRIVER-2149). It may be
	// changed or removed in any release.
	Crypt driver.Crypt

	// Deployment specifies a custom deployment to use for the new Client.
	//
	// Deprecated: This option is for internal use only and should not be set. It may be changed or removed in any
	// release.
	Deployment driver.Deployment

	// SocketTimeout specifies the timeout to be used for the Client's socket reads and writes.
	//
	// NOTE(benjirewis): SocketTimeout will be deprecated in a future release. The more general Timeout option
	// may be used in its place to control the amount of time that a single operation can run before returning
	// an error. Setting SocketTimeout and Timeout on a single client will result in undefined behavior.
	SocketTimeout *time.Duration
***REMOVED***

// Client creates a new ClientOptions instance.
func Client() *ClientOptions ***REMOVED***
	return new(ClientOptions)
***REMOVED***

// Validate validates the client options. This method will return the first error found.
func (c *ClientOptions) Validate() error ***REMOVED***
	if c.err != nil ***REMOVED***
		return c.err
	***REMOVED***
	c.err = c.validate()
	return c.err
***REMOVED***

func (c *ClientOptions) validate() error ***REMOVED***
	// Direct connections cannot be made if multiple hosts are specified or an SRV URI is used.
	if c.Direct != nil && *c.Direct ***REMOVED***
		if len(c.Hosts) > 1 ***REMOVED***
			return errors.New("a direct connection cannot be made if multiple hosts are specified")
		***REMOVED***
		if c.cs != nil && c.cs.Scheme == connstring.SchemeMongoDBSRV ***REMOVED***
			return errors.New("a direct connection cannot be made if an SRV URI is used")
		***REMOVED***
	***REMOVED***

	if c.MaxPoolSize != nil && c.MinPoolSize != nil && *c.MaxPoolSize != 0 && *c.MinPoolSize > *c.MaxPoolSize ***REMOVED***
		return fmt.Errorf("minPoolSize must be less than or equal to maxPoolSize, got minPoolSize=%d maxPoolSize=%d", *c.MinPoolSize, *c.MaxPoolSize)
	***REMOVED***

	// verify server API version if ServerAPIOptions are passed in.
	if c.ServerAPIOptions != nil ***REMOVED***
		if err := c.ServerAPIOptions.ServerAPIVersion.Validate(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Validation for load-balanced mode.
	if c.LoadBalanced != nil && *c.LoadBalanced ***REMOVED***
		if len(c.Hosts) > 1 ***REMOVED***
			return internal.ErrLoadBalancedWithMultipleHosts
		***REMOVED***
		if c.ReplicaSet != nil ***REMOVED***
			return internal.ErrLoadBalancedWithReplicaSet
		***REMOVED***
		if c.Direct != nil ***REMOVED***
			return internal.ErrLoadBalancedWithDirectConnection
		***REMOVED***
	***REMOVED***

	// Validation for srvMaxHosts.
	if c.SRVMaxHosts != nil && *c.SRVMaxHosts > 0 ***REMOVED***
		if c.ReplicaSet != nil ***REMOVED***
			return internal.ErrSRVMaxHostsWithReplicaSet
		***REMOVED***
		if c.LoadBalanced != nil && *c.LoadBalanced ***REMOVED***
			return internal.ErrSRVMaxHostsWithLoadBalanced
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetURI returns the original URI used to configure the ClientOptions instance. If ApplyURI was not called during
// construction, this returns "".
func (c *ClientOptions) GetURI() string ***REMOVED***
	return c.uri
***REMOVED***

// ApplyURI parses the given URI and sets options accordingly. The URI can contain host names, IPv4/IPv6 literals, or
// an SRV record that will be resolved when the Client is created. When using an SRV record, TLS support is
// implictly enabled. Specify the "tls=false" URI option to override this.
//
// If the connection string contains any options that have previously been set, it will overwrite them. Options that
// correspond to multiple URI parameters, such as WriteConcern, will be completely overwritten if any of the query
// parameters are specified. If an option is set on ClientOptions after this method is called, that option will override
// any option applied via the connection string.
//
// If the URI format is incorrect or there are conflicting options specified in the URI an error will be recorded and
// can be retrieved by calling Validate.
//
// For more information about the URI format, see https://www.mongodb.com/docs/manual/reference/connection-string/. See
// mongo.Connect documentation for examples of using URIs for different Client configurations.
func (c *ClientOptions) ApplyURI(uri string) *ClientOptions ***REMOVED***
	if c.err != nil ***REMOVED***
		return c
	***REMOVED***

	c.uri = uri
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil ***REMOVED***
		c.err = err
		return c
	***REMOVED***
	c.cs = &cs

	if cs.AppName != "" ***REMOVED***
		c.AppName = &cs.AppName
	***REMOVED***

	// Only create a Credential if there is a request for authentication via non-empty credentials in the URI.
	if cs.HasAuthParameters() ***REMOVED***
		c.Auth = &Credential***REMOVED***
			AuthMechanism:           cs.AuthMechanism,
			AuthMechanismProperties: cs.AuthMechanismProperties,
			AuthSource:              cs.AuthSource,
			Username:                cs.Username,
			Password:                cs.Password,
			PasswordSet:             cs.PasswordSet,
		***REMOVED***
	***REMOVED***

	if cs.ConnectSet ***REMOVED***
		direct := cs.Connect == connstring.SingleConnect
		c.Direct = &direct
	***REMOVED***

	if cs.DirectConnectionSet ***REMOVED***
		c.Direct = &cs.DirectConnection
	***REMOVED***

	if cs.ConnectTimeoutSet ***REMOVED***
		c.ConnectTimeout = &cs.ConnectTimeout
	***REMOVED***

	if len(cs.Compressors) > 0 ***REMOVED***
		c.Compressors = cs.Compressors
		if stringSliceContains(c.Compressors, "zlib") ***REMOVED***
			defaultLevel := wiremessage.DefaultZlibLevel
			c.ZlibLevel = &defaultLevel
		***REMOVED***
		if stringSliceContains(c.Compressors, "zstd") ***REMOVED***
			defaultLevel := wiremessage.DefaultZstdLevel
			c.ZstdLevel = &defaultLevel
		***REMOVED***
	***REMOVED***

	if cs.HeartbeatIntervalSet ***REMOVED***
		c.HeartbeatInterval = &cs.HeartbeatInterval
	***REMOVED***

	c.Hosts = cs.Hosts

	if cs.LoadBalancedSet ***REMOVED***
		c.LoadBalanced = &cs.LoadBalanced
	***REMOVED***

	if cs.LocalThresholdSet ***REMOVED***
		c.LocalThreshold = &cs.LocalThreshold
	***REMOVED***

	if cs.MaxConnIdleTimeSet ***REMOVED***
		c.MaxConnIdleTime = &cs.MaxConnIdleTime
	***REMOVED***

	if cs.MaxPoolSizeSet ***REMOVED***
		c.MaxPoolSize = &cs.MaxPoolSize
	***REMOVED***

	if cs.MinPoolSizeSet ***REMOVED***
		c.MinPoolSize = &cs.MinPoolSize
	***REMOVED***

	if cs.MaxConnectingSet ***REMOVED***
		c.MaxConnecting = &cs.MaxConnecting
	***REMOVED***

	if cs.ReadConcernLevel != "" ***REMOVED***
		c.ReadConcern = readconcern.New(readconcern.Level(cs.ReadConcernLevel))
	***REMOVED***

	if cs.ReadPreference != "" || len(cs.ReadPreferenceTagSets) > 0 || cs.MaxStalenessSet ***REMOVED***
		opts := make([]readpref.Option, 0, 1)

		tagSets := tag.NewTagSetsFromMaps(cs.ReadPreferenceTagSets)
		if len(tagSets) > 0 ***REMOVED***
			opts = append(opts, readpref.WithTagSets(tagSets...))
		***REMOVED***

		if cs.MaxStaleness != 0 ***REMOVED***
			opts = append(opts, readpref.WithMaxStaleness(cs.MaxStaleness))
		***REMOVED***

		mode, err := readpref.ModeFromString(cs.ReadPreference)
		if err != nil ***REMOVED***
			c.err = err
			return c
		***REMOVED***

		c.ReadPreference, c.err = readpref.New(mode, opts...)
		if c.err != nil ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***

	if cs.RetryWritesSet ***REMOVED***
		c.RetryWrites = &cs.RetryWrites
	***REMOVED***

	if cs.RetryReadsSet ***REMOVED***
		c.RetryReads = &cs.RetryReads
	***REMOVED***

	if cs.ReplicaSet != "" ***REMOVED***
		c.ReplicaSet = &cs.ReplicaSet
	***REMOVED***

	if cs.ServerSelectionTimeoutSet ***REMOVED***
		c.ServerSelectionTimeout = &cs.ServerSelectionTimeout
	***REMOVED***

	if cs.SocketTimeoutSet ***REMOVED***
		c.SocketTimeout = &cs.SocketTimeout
	***REMOVED***

	if cs.SRVMaxHosts != 0 ***REMOVED***
		c.SRVMaxHosts = &cs.SRVMaxHosts
	***REMOVED***

	if cs.SRVServiceName != "" ***REMOVED***
		c.SRVServiceName = &cs.SRVServiceName
	***REMOVED***

	if cs.SSL ***REMOVED***
		tlsConfig := new(tls.Config)

		if cs.SSLCaFileSet ***REMOVED***
			c.err = addCACertFromFile(tlsConfig, cs.SSLCaFile)
			if c.err != nil ***REMOVED***
				return c
			***REMOVED***
		***REMOVED***

		if cs.SSLInsecure ***REMOVED***
			tlsConfig.InsecureSkipVerify = true
		***REMOVED***

		var x509Subject string
		var keyPasswd string
		if cs.SSLClientCertificateKeyPasswordSet && cs.SSLClientCertificateKeyPassword != nil ***REMOVED***
			keyPasswd = cs.SSLClientCertificateKeyPassword()
		***REMOVED***
		if cs.SSLClientCertificateKeyFileSet ***REMOVED***
			x509Subject, err = addClientCertFromConcatenatedFile(tlsConfig, cs.SSLClientCertificateKeyFile, keyPasswd)
		***REMOVED*** else if cs.SSLCertificateFileSet || cs.SSLPrivateKeyFileSet ***REMOVED***
			x509Subject, err = addClientCertFromSeparateFiles(tlsConfig, cs.SSLCertificateFile,
				cs.SSLPrivateKeyFile, keyPasswd)
		***REMOVED***
		if err != nil ***REMOVED***
			c.err = err
			return c
		***REMOVED***

		// If a username wasn't specified fork x509, add one from the certificate.
		if c.Auth != nil && strings.ToLower(c.Auth.AuthMechanism) == "mongodb-x509" &&
			c.Auth.Username == "" ***REMOVED***

			// The Go x509 package gives the subject with the pairs in reverse order that we want.
			c.Auth.Username = extractX509UsernameFromSubject(x509Subject)
		***REMOVED***

		c.TLSConfig = tlsConfig
	***REMOVED***

	if cs.JSet || cs.WString != "" || cs.WNumberSet || cs.WTimeoutSet ***REMOVED***
		opts := make([]writeconcern.Option, 0, 1)

		if len(cs.WString) > 0 ***REMOVED***
			opts = append(opts, writeconcern.WTagSet(cs.WString))
		***REMOVED*** else if cs.WNumberSet ***REMOVED***
			opts = append(opts, writeconcern.W(cs.WNumber))
		***REMOVED***

		if cs.JSet ***REMOVED***
			opts = append(opts, writeconcern.J(cs.J))
		***REMOVED***

		if cs.WTimeoutSet ***REMOVED***
			opts = append(opts, writeconcern.WTimeout(cs.WTimeout))
		***REMOVED***

		c.WriteConcern = writeconcern.New(opts...)
	***REMOVED***

	if cs.ZlibLevelSet ***REMOVED***
		c.ZlibLevel = &cs.ZlibLevel
	***REMOVED***
	if cs.ZstdLevelSet ***REMOVED***
		c.ZstdLevel = &cs.ZstdLevel
	***REMOVED***

	if cs.SSLDisableOCSPEndpointCheckSet ***REMOVED***
		c.DisableOCSPEndpointCheck = &cs.SSLDisableOCSPEndpointCheck
	***REMOVED***

	if cs.TimeoutSet ***REMOVED***
		c.Timeout = &cs.Timeout
	***REMOVED***

	return c
***REMOVED***

// SetAppName specifies an application name that is sent to the server when creating new connections. It is used by the
// server to log connection and profiling information (e.g. slow query logs). This can also be set through the "appName"
// URI option (e.g "appName=example_application"). The default is empty, meaning no app name will be sent.
func (c *ClientOptions) SetAppName(s string) *ClientOptions ***REMOVED***
	c.AppName = &s
	return c
***REMOVED***

// SetAuth specifies a Credential containing options for configuring authentication. See the options.Credential
// documentation for more information about Credential fields. The default is an empty Credential, meaning no
// authentication will be configured.
func (c *ClientOptions) SetAuth(auth Credential) *ClientOptions ***REMOVED***
	c.Auth = &auth
	return c
***REMOVED***

// SetCompressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
// 3. "zstd" - requires server version >= 4.2, and driver version >= 1.2.0 with cgo support enabled or driver
// version >= 1.3.0 without cgo.
//
// If this option is specified, the driver will perform a negotiation with the server to determine a common list of of
// compressors and will use the first one in that list when performing operations. See
// https://www.mongodb.com/docs/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about configuring compression on the server and the server-side defaults.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func (c *ClientOptions) SetCompressors(comps []string) *ClientOptions ***REMOVED***
	c.Compressors = comps

	return c
***REMOVED***

// SetConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
// specified through SetDialer, this option must not be used. This can be set through ApplyURI with the
// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
// seconds.
func (c *ClientOptions) SetConnectTimeout(d time.Duration) *ClientOptions ***REMOVED***
	c.ConnectTimeout = &d
	return c
***REMOVED***

// SetDialer specifies a custom ContextDialer to be used to create new connections to the server. The default is a
// net.Dialer with the Timeout field set to ConnectTimeout. See https://golang.org/pkg/net/#Dialer for more information
// about the net.Dialer type.
func (c *ClientOptions) SetDialer(d ContextDialer) *ClientOptions ***REMOVED***
	c.Dialer = d
	return c
***REMOVED***

// SetDirect specifies whether or not a direct connect should be made. If set to true, the driver will only connect to
// the host provided in the URI and will not discover other hosts in the cluster. This can also be set through the
// "directConnection" URI option. This option cannot be set to true if multiple hosts are specified, either through
// ApplyURI or SetHosts, or an SRV URI is used.
//
// As of driver version 1.4, the "connect" URI option has been deprecated and replaced with "directConnection". The
// "connect" URI option has two values:
//
// 1. "connect=direct" for direct connections. This corresponds to "directConnection=true".
//
// 2. "connect=automatic" for automatic discovery. This corresponds to "directConnection=false"
//
// If the "connect" and "directConnection" URI options are both specified in the connection string, their values must
// not conflict. Direct connections are not valid if multiple hosts are specified or an SRV URI is used. The default
// value for this option is false.
func (c *ClientOptions) SetDirect(b bool) *ClientOptions ***REMOVED***
	c.Direct = &b
	return c
***REMOVED***

// SetHeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
func (c *ClientOptions) SetHeartbeatInterval(d time.Duration) *ClientOptions ***REMOVED***
	c.HeartbeatInterval = &d
	return c
***REMOVED***

// SetHosts specifies a list of host names or IP addresses for servers in a cluster. Both IPv4 and IPv6 addresses are
// supported. IPv6 literals must be enclosed in '[]' following RFC-2732 syntax.
//
// Hosts can also be specified as a comma-separated list in a URI. For example, to include "localhost:27017" and
// "localhost:27018", a URI could be "mongodb://localhost:27017,localhost:27018". The default is ["localhost:27017"]
func (c *ClientOptions) SetHosts(s []string) *ClientOptions ***REMOVED***
	c.Hosts = s
	return c
***REMOVED***

// SetLoadBalanced specifies whether or not the MongoDB deployment is hosted behind a load balancer. This can also be
// set through the "loadBalanced" URI option. The driver will error during Client configuration if this option is set
// to true and one of the following conditions are met:
//
// 1. Multiple hosts are specified, either via the ApplyURI or SetHosts methods. This includes the case where an SRV
// URI is used and the SRV record resolves to multiple hostnames.
// 2. A replica set name is specified, either via the URI or the SetReplicaSet method.
// 3. The options specify whether or not a direct connection should be made, either via the URI or the SetDirect method.
//
// The default value is false.
func (c *ClientOptions) SetLoadBalanced(lb bool) *ClientOptions ***REMOVED***
	c.LoadBalanced = &lb
	return c
***REMOVED***

// SetLocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
// "localThresholdMS=15000"). The default is 15 milliseconds.
func (c *ClientOptions) SetLocalThreshold(d time.Duration) *ClientOptions ***REMOVED***
	c.LocalThreshold = &d
	return c
***REMOVED***

// SetMaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
func (c *ClientOptions) SetMaxConnIdleTime(d time.Duration) *ClientOptions ***REMOVED***
	c.MaxConnIdleTime = &d
	return c
***REMOVED***

// SetMaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
// (e.g. "maxPoolSize=100"). If this is 0, maximum connection pool size is not limited. The default is 100.
func (c *ClientOptions) SetMaxPoolSize(u uint64) *ClientOptions ***REMOVED***
	c.MaxPoolSize = &u
	return c
***REMOVED***

// SetMinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
func (c *ClientOptions) SetMinPoolSize(u uint64) *ClientOptions ***REMOVED***
	c.MinPoolSize = &u
	return c
***REMOVED***

// SetMaxConnecting specifies the maximum number of connections a connection pool may establish simultaneously. This can
// also be set through the "maxConnecting" URI option (e.g. "maxConnecting=2"). If this is 0, the default is used. The
// default is 2. Values greater than 100 are not recommended.
func (c *ClientOptions) SetMaxConnecting(u uint64) *ClientOptions ***REMOVED***
	c.MaxConnecting = &u
	return c
***REMOVED***

// SetPoolMonitor specifies a PoolMonitor to receive connection pool events. See the event.PoolMonitor documentation
// for more information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetPoolMonitor(m *event.PoolMonitor) *ClientOptions ***REMOVED***
	c.PoolMonitor = m
	return c
***REMOVED***

// SetMonitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetMonitor(m *event.CommandMonitor) *ClientOptions ***REMOVED***
	c.Monitor = m
	return c
***REMOVED***

// SetServerMonitor specifies an SDAM monitor used to monitor SDAM events.
func (c *ClientOptions) SetServerMonitor(m *event.ServerMonitor) *ClientOptions ***REMOVED***
	c.ServerMonitor = m
	return c
***REMOVED***

// SetReadConcern specifies the read concern to use for read operations. A read concern level can also be set through
// the "readConcernLevel" URI option (e.g. "readConcernLevel=majority"). The default is nil, meaning the server will use
// its configured default.
func (c *ClientOptions) SetReadConcern(rc *readconcern.ReadConcern) *ClientOptions ***REMOVED***
	c.ReadConcern = rc

	return c
***REMOVED***

// SetReadPreference specifies the read preference to use for read operations. This can also be set through the
// following URI options:
//
// 1. "readPreference" - Specify the read preference mode (e.g. "readPreference=primary").
//
// 2. "readPreferenceTags": Specify one or more read preference tags
// (e.g. "readPreferenceTags=region:south,datacenter:A").
//
// 3. "maxStalenessSeconds" (or "maxStaleness"): Specify a maximum replication lag for reads from secondaries in a
// replica set (e.g. "maxStalenessSeconds=10").
//
// The default is readpref.Primary(). See https://www.mongodb.com/docs/manual/core/read-preference/#read-preference for
// more information about read preferences.
func (c *ClientOptions) SetReadPreference(rp *readpref.ReadPref) *ClientOptions ***REMOVED***
	c.ReadPreference = rp

	return c
***REMOVED***

// SetRegistry specifies the BSON registry to use for BSON marshalling/unmarshalling operations. The default is
// bson.DefaultRegistry.
func (c *ClientOptions) SetRegistry(registry *bsoncodec.Registry) *ClientOptions ***REMOVED***
	c.Registry = registry
	return c
***REMOVED***

// SetReplicaSet specifies the replica set name for the cluster. If specified, the cluster will be treated as a replica
// set and the driver will automatically discover all servers in the set, starting with the nodes specified through
// ApplyURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be
// considered as part of the set by the Client. This can also be set through the "replicaSet" URI option (e.g.
// "replicaSet=replset"). The default is empty.
func (c *ClientOptions) SetReplicaSet(s string) *ClientOptions ***REMOVED***
	c.ReplicaSet = &s
	return c
***REMOVED***

// SetRetryWrites specifies whether supported write operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are InsertOne, UpdateOne, ReplaceOne, DeleteOne, FindOneAndDelete, FindOneAndReplace,
// FindOneAndDelete, InsertMany, and BulkWrite. Note that BulkWrite requests must not include UpdateManyModel or
// DeleteManyModel instances to be considered retryable. Unacknowledged writes will not be retried, even if this option
// is set to true.
//
// This option requires server version >= 3.6 and a replica set or sharded cluster and will be ignored for any other
// cluster type. This can also be set through the "retryWrites" URI option (e.g. "retryWrites=true"). The default is
// true.
func (c *ClientOptions) SetRetryWrites(b bool) *ClientOptions ***REMOVED***
	c.RetryWrites = &b

	return c
***REMOVED***

// SetRetryReads specifies whether supported read operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are Find, FindOne, Aggregate without a $out stage, Distinct, CountDocuments,
// EstimatedDocumentCount, Watch (for Client, Database, and Collection), ListCollections, and ListDatabases. Note that
// operations run through RunCommand are not retried.
//
// This option requires server version >= 3.6 and driver version >= 1.1.0. The default is true.
func (c *ClientOptions) SetRetryReads(b bool) *ClientOptions ***REMOVED***
	c.RetryReads = &b
	return c
***REMOVED***

// SetServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
func (c *ClientOptions) SetServerSelectionTimeout(d time.Duration) *ClientOptions ***REMOVED***
	c.ServerSelectionTimeout = &d
	return c
***REMOVED***

// SetSocketTimeout specifies how long the driver will wait for a socket read or write to return before returning a
// network error. This can also be set through the "socketTimeoutMS" URI option (e.g. "socketTimeoutMS=1000"). The
// default value is 0, meaning no timeout is used and socket operations can block indefinitely.
//
// NOTE(benjirewis): SocketTimeout will be deprecated in a future release. The more general Timeout option may be used
// in its place to control the amount of time that a single operation can run before returning an error. Setting
// SocketTimeout and Timeout on a single client will result in undefined behavior.
func (c *ClientOptions) SetSocketTimeout(d time.Duration) *ClientOptions ***REMOVED***
	c.SocketTimeout = &d
	return c
***REMOVED***

// SetTimeout specifies the amount of time that a single operation run on this Client can execute before returning an error.
// The deadline of any operation run through the Client will be honored above any Timeout set on the Client; Timeout will only
// be honored if there is no deadline on the operation Context. Timeout can also be set through the "timeoutMS" URI option
// (e.g. "timeoutMS=1000"). The default value is nil, meaning operations do not inherit a timeout from the Client.
//
// If any Timeout is set (even 0) on the Client, the values of MaxTime on operations, TransactionOptions.MaxCommitTime and
// SessionOptions.DefaultMaxCommitTime will be ignored. Setting Timeout and ClientOptions.SocketTimeout or WriteConcern.wTimeout
// will result in undefined behavior.
//
// NOTE(benjirewis): SetTimeout represents unstable, provisional API. The behavior of the driver when a Timeout is specified is
// subject to change.
func (c *ClientOptions) SetTimeout(d time.Duration) *ClientOptions ***REMOVED***
	c.Timeout = &d
	return c
***REMOVED***

// SetTLSConfig specifies a tls.Config instance to use use to configure TLS on all connections created to the cluster.
// This can also be set through the following URI options:
//
// 1. "tls" (or "ssl"): Specify if TLS should be used (e.g. "tls=true").
//
// 2. Either "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile") or a combination of "tlsCertificateFile" and
// "tlsPrivateKeyFile". The "tlsCertificateKeyFile" option specifies a path to the client certificate and private key,
// which must be concatenated into one file. The "tlsCertificateFile" and "tlsPrivateKey" combination specifies separate
// paths to the client certificate and private key, respectively. Note that if "tlsCertificateKeyFile" is used, the
// other two options must not be specified.
//
// 3. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 4. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// 5. "tlsInsecure" (or "sslInsecure"): Specifies whether or not certificates and hostnames received from the server
// should be validated. If true (e.g. "tlsInsecure=true"), the TLS library will accept any certificate presented by the
// server and any host name in that certificate. Note that setting this to true makes TLS susceptible to
// man-in-the-middle attacks and should only be done for testing.
//
// The default is nil, meaning no TLS will be enabled.
func (c *ClientOptions) SetTLSConfig(cfg *tls.Config) *ClientOptions ***REMOVED***
	c.TLSConfig = cfg
	return c
***REMOVED***

// SetWriteConcern specifies the write concern to use to for write operations. This can also be set through the following
// URI options:
//
// 1. "w": Specify the number of nodes in the cluster that must acknowledge write operations before the operation
// returns or "majority" to specify that a majority of the nodes must acknowledge writes. This can either be an integer
// (e.g. "w=10") or the string "majority" (e.g. "w=majority").
//
// 2. "wTimeoutMS": Specify how long write operations should wait for the correct number of nodes to acknowledge the
// operation (e.g. "wTimeoutMS=1000").
//
// 3. "journal": Specifies whether or not write operations should be written to an on-disk journal on the server before
// returning (e.g. "journal=true").
//
// The default is nil, meaning the server will use its configured default.
func (c *ClientOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *ClientOptions ***REMOVED***
	c.WriteConcern = wc

	return c
***REMOVED***

// SetZlibLevel specifies the level for the zlib compressor. This option is ignored if zlib is not specified as a
// compressor through ApplyURI or SetCompressors. Supported values are -1 through 9, inclusive. -1 tells the zlib
// library to use its default, 0 means no compression, 1 means best speed, and 9 means best compression.
// This can also be set through the "zlibCompressionLevel" URI option (e.g. "zlibCompressionLevel=-1"). Defaults to -1.
func (c *ClientOptions) SetZlibLevel(level int) *ClientOptions ***REMOVED***
	c.ZlibLevel = &level

	return c
***REMOVED***

// SetZstdLevel sets the level for the zstd compressor. This option is ignored if zstd is not specified as a compressor
// through ApplyURI or SetCompressors. Supported values are 1 through 20, inclusive. 1 means best speed and 20 means
// best compression. This can also be set through the "zstdCompressionLevel" URI option. Defaults to 6.
func (c *ClientOptions) SetZstdLevel(level int) *ClientOptions ***REMOVED***
	c.ZstdLevel = &level
	return c
***REMOVED***

// SetAutoEncryptionOptions specifies an AutoEncryptionOptions instance to automatically encrypt and decrypt commands
// and their results. See the options.AutoEncryptionOptions documentation for more information about the supported
// options.
func (c *ClientOptions) SetAutoEncryptionOptions(opts *AutoEncryptionOptions) *ClientOptions ***REMOVED***
	c.AutoEncryptionOptions = opts
	return c
***REMOVED***

// SetDisableOCSPEndpointCheck specifies whether or not the driver should reach out to OCSP responders to verify the
// certificate status for certificates presented by the server that contain a list of OCSP responders.
//
// If set to true, the driver will verify the status of the certificate using a response stapled by the server, if there
// is one, but will not send an HTTP request to any responders if there is no staple. In this case, the driver will
// continue the connection even though the certificate status is not known.
//
// This can also be set through the tlsDisableOCSPEndpointCheck URI option. Both this URI option and tlsInsecure must
// not be set at the same time and will error if they are. The default value is false.
func (c *ClientOptions) SetDisableOCSPEndpointCheck(disableCheck bool) *ClientOptions ***REMOVED***
	c.DisableOCSPEndpointCheck = &disableCheck
	return c
***REMOVED***

// SetServerAPIOptions specifies a ServerAPIOptions instance used to configure the API version sent to the server
// when running commands. See the options.ServerAPIOptions documentation for more information about the supported
// options.
func (c *ClientOptions) SetServerAPIOptions(opts *ServerAPIOptions) *ClientOptions ***REMOVED***
	c.ServerAPIOptions = opts
	return c
***REMOVED***

// SetSRVMaxHosts specifies the maximum number of SRV results to randomly select during polling. To limit the number
// of hosts selected in SRV discovery, this function must be called before ApplyURI. This can also be set through
// the "srvMaxHosts" URI option.
func (c *ClientOptions) SetSRVMaxHosts(srvMaxHosts int) *ClientOptions ***REMOVED***
	c.SRVMaxHosts = &srvMaxHosts
	return c
***REMOVED***

// SetSRVServiceName specifies a custom SRV service name to use in SRV polling. To use a custom SRV service name
// in SRV discovery, this function must be called before ApplyURI. This can also be set through the "srvServiceName"
// URI option.
func (c *ClientOptions) SetSRVServiceName(srvName string) *ClientOptions ***REMOVED***
	c.SRVServiceName = &srvName
	return c
***REMOVED***

// MergeClientOptions combines the given *ClientOptions into a single *ClientOptions in a last one wins fashion.
// The specified options are merged with the existing options on the client, with the specified options taking
// precedence.
func MergeClientOptions(opts ...*ClientOptions) *ClientOptions ***REMOVED***
	c := Client()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***

		if opt.Dialer != nil ***REMOVED***
			c.Dialer = opt.Dialer
		***REMOVED***
		if opt.AppName != nil ***REMOVED***
			c.AppName = opt.AppName
		***REMOVED***
		if opt.Auth != nil ***REMOVED***
			c.Auth = opt.Auth
		***REMOVED***
		if opt.AuthenticateToAnything != nil ***REMOVED***
			c.AuthenticateToAnything = opt.AuthenticateToAnything
		***REMOVED***
		if opt.Compressors != nil ***REMOVED***
			c.Compressors = opt.Compressors
		***REMOVED***
		if opt.ConnectTimeout != nil ***REMOVED***
			c.ConnectTimeout = opt.ConnectTimeout
		***REMOVED***
		if opt.Crypt != nil ***REMOVED***
			c.Crypt = opt.Crypt
		***REMOVED***
		if opt.HeartbeatInterval != nil ***REMOVED***
			c.HeartbeatInterval = opt.HeartbeatInterval
		***REMOVED***
		if len(opt.Hosts) > 0 ***REMOVED***
			c.Hosts = opt.Hosts
		***REMOVED***
		if opt.LoadBalanced != nil ***REMOVED***
			c.LoadBalanced = opt.LoadBalanced
		***REMOVED***
		if opt.LocalThreshold != nil ***REMOVED***
			c.LocalThreshold = opt.LocalThreshold
		***REMOVED***
		if opt.MaxConnIdleTime != nil ***REMOVED***
			c.MaxConnIdleTime = opt.MaxConnIdleTime
		***REMOVED***
		if opt.MaxPoolSize != nil ***REMOVED***
			c.MaxPoolSize = opt.MaxPoolSize
		***REMOVED***
		if opt.MinPoolSize != nil ***REMOVED***
			c.MinPoolSize = opt.MinPoolSize
		***REMOVED***
		if opt.MaxConnecting != nil ***REMOVED***
			c.MaxConnecting = opt.MaxConnecting
		***REMOVED***
		if opt.PoolMonitor != nil ***REMOVED***
			c.PoolMonitor = opt.PoolMonitor
		***REMOVED***
		if opt.Monitor != nil ***REMOVED***
			c.Monitor = opt.Monitor
		***REMOVED***
		if opt.ServerAPIOptions != nil ***REMOVED***
			c.ServerAPIOptions = opt.ServerAPIOptions
		***REMOVED***
		if opt.ServerMonitor != nil ***REMOVED***
			c.ServerMonitor = opt.ServerMonitor
		***REMOVED***
		if opt.ReadConcern != nil ***REMOVED***
			c.ReadConcern = opt.ReadConcern
		***REMOVED***
		if opt.ReadPreference != nil ***REMOVED***
			c.ReadPreference = opt.ReadPreference
		***REMOVED***
		if opt.Registry != nil ***REMOVED***
			c.Registry = opt.Registry
		***REMOVED***
		if opt.ReplicaSet != nil ***REMOVED***
			c.ReplicaSet = opt.ReplicaSet
		***REMOVED***
		if opt.RetryWrites != nil ***REMOVED***
			c.RetryWrites = opt.RetryWrites
		***REMOVED***
		if opt.RetryReads != nil ***REMOVED***
			c.RetryReads = opt.RetryReads
		***REMOVED***
		if opt.ServerSelectionTimeout != nil ***REMOVED***
			c.ServerSelectionTimeout = opt.ServerSelectionTimeout
		***REMOVED***
		if opt.Direct != nil ***REMOVED***
			c.Direct = opt.Direct
		***REMOVED***
		if opt.SocketTimeout != nil ***REMOVED***
			c.SocketTimeout = opt.SocketTimeout
		***REMOVED***
		if opt.SRVMaxHosts != nil ***REMOVED***
			c.SRVMaxHosts = opt.SRVMaxHosts
		***REMOVED***
		if opt.SRVServiceName != nil ***REMOVED***
			c.SRVServiceName = opt.SRVServiceName
		***REMOVED***
		if opt.Timeout != nil ***REMOVED***
			c.Timeout = opt.Timeout
		***REMOVED***
		if opt.TLSConfig != nil ***REMOVED***
			c.TLSConfig = opt.TLSConfig
		***REMOVED***
		if opt.WriteConcern != nil ***REMOVED***
			c.WriteConcern = opt.WriteConcern
		***REMOVED***
		if opt.ZlibLevel != nil ***REMOVED***
			c.ZlibLevel = opt.ZlibLevel
		***REMOVED***
		if opt.ZstdLevel != nil ***REMOVED***
			c.ZstdLevel = opt.ZstdLevel
		***REMOVED***
		if opt.AutoEncryptionOptions != nil ***REMOVED***
			c.AutoEncryptionOptions = opt.AutoEncryptionOptions
		***REMOVED***
		if opt.Deployment != nil ***REMOVED***
			c.Deployment = opt.Deployment
		***REMOVED***
		if opt.DisableOCSPEndpointCheck != nil ***REMOVED***
			c.DisableOCSPEndpointCheck = opt.DisableOCSPEndpointCheck
		***REMOVED***
		if opt.err != nil ***REMOVED***
			c.err = opt.err
		***REMOVED***
		if opt.uri != "" ***REMOVED***
			c.uri = opt.uri
		***REMOVED***
		if opt.cs != nil ***REMOVED***
			c.cs = opt.cs
		***REMOVED***
	***REMOVED***

	return c
***REMOVED***

// addCACertFromFile adds a root CA certificate to the configuration given a path
// to the containing file.
func addCACertFromFile(cfg *tls.Config, file string) error ***REMOVED***
	data, err := ioutil.ReadFile(file)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if cfg.RootCAs == nil ***REMOVED***
		cfg.RootCAs = x509.NewCertPool()
	***REMOVED***
	if !cfg.RootCAs.AppendCertsFromPEM(data) ***REMOVED***
		return errors.New("the specified CA file does not contain any valid certificates")
	***REMOVED***

	return nil
***REMOVED***

func addClientCertFromSeparateFiles(cfg *tls.Config, keyFile, certFile, keyPassword string) (string, error) ***REMOVED***
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	certData, err := ioutil.ReadFile(certFile)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	data := make([]byte, 0, len(keyData)+len(certData)+1)
	data = append(data, keyData...)
	data = append(data, '\n')
	data = append(data, certData...)
	return addClientCertFromBytes(cfg, data, keyPassword)
***REMOVED***

func addClientCertFromConcatenatedFile(cfg *tls.Config, certKeyFile, keyPassword string) (string, error) ***REMOVED***
	data, err := ioutil.ReadFile(certKeyFile)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return addClientCertFromBytes(cfg, data, keyPassword)
***REMOVED***

// addClientCertFromBytes adds a client certificate to the configuration given a path to the
// containing file and returns the certificate's subject name.
func addClientCertFromBytes(cfg *tls.Config, data []byte, keyPasswd string) (string, error) ***REMOVED***
	var currentBlock *pem.Block
	var certDecodedBlock []byte
	var certBlocks, keyBlocks [][]byte

	remaining := data
	start := 0
	for ***REMOVED***
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil ***REMOVED***
			break
		***REMOVED***

		if currentBlock.Type == "CERTIFICATE" ***REMOVED***
			certBlock := data[start : len(data)-len(remaining)]
			certBlocks = append(certBlocks, certBlock)
			certDecodedBlock = currentBlock.Bytes
			start += len(certBlock)
		***REMOVED*** else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") ***REMOVED***
			isEncrypted := x509.IsEncryptedPEMBlock(currentBlock) || strings.Contains(currentBlock.Type, "ENCRYPTED PRIVATE KEY")
			if isEncrypted ***REMOVED***
				if keyPasswd == "" ***REMOVED***
					return "", fmt.Errorf("no password provided to decrypt private key")
				***REMOVED***

				var keyBytes []byte
				var err error
				// Process the X.509-encrypted or PKCS-encrypted PEM block.
				if x509.IsEncryptedPEMBlock(currentBlock) ***REMOVED***
					// Only covers encrypted PEM data with a DEK-Info header.
					keyBytes, err = x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
				***REMOVED*** else if strings.Contains(currentBlock.Type, "ENCRYPTED") ***REMOVED***
					// The pkcs8 package only handles the PKCS #5 v2.0 scheme.
					decrypted, err := pkcs8.ParsePKCS8PrivateKey(currentBlock.Bytes, []byte(keyPasswd))
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
					keyBytes, err = x509.MarshalPKCS8PrivateKey(decrypted)
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
				***REMOVED***
				var encoded bytes.Buffer
				pem.Encode(&encoded, &pem.Block***REMOVED***Type: currentBlock.Type, Bytes: keyBytes***REMOVED***)
				keyBlock := encoded.Bytes()
				keyBlocks = append(keyBlocks, keyBlock)
				start = len(data) - len(remaining)
			***REMOVED*** else ***REMOVED***
				keyBlock := data[start : len(data)-len(remaining)]
				keyBlocks = append(keyBlocks, keyBlock)
				start += len(keyBlock)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(certBlocks) == 0 ***REMOVED***
		return "", fmt.Errorf("failed to find CERTIFICATE")
	***REMOVED***
	if len(keyBlocks) == 0 ***REMOVED***
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	***REMOVED***

	cert, err := tls.X509KeyPair(bytes.Join(certBlocks, []byte("\n")), bytes.Join(keyBlocks, []byte("\n")))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	cfg.Certificates = append(cfg.Certificates, cert)

	// The documentation for the tls.X509KeyPair indicates that the Leaf certificate is not
	// retained.
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return crt.Subject.String(), nil
***REMOVED***

func stringSliceContains(source []string, target string) bool ***REMOVED***
	for _, str := range source ***REMOVED***
		if str == target ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// create a username for x509 authentication from an x509 certificate subject.
func extractX509UsernameFromSubject(subject string) string ***REMOVED***
	// the Go x509 package gives the subject with the pairs in the reverse order from what we want.
	pairs := strings.Split(subject, ",")
	for left, right := 0, len(pairs)-1; left < right; left, right = left+1, right-1 ***REMOVED***
		pairs[left], pairs[right] = pairs[right], pairs[left]
	***REMOVED***

	return strings.Join(pairs, ",")
***REMOVED***
