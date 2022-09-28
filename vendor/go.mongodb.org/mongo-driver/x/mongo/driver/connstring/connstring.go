// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package connstring // import "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/internal/randutil"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/dns"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// random is a package-global pseudo-random number generator.
var random = randutil.NewLockedRand()

// ParseAndValidate parses the provided URI into a ConnString object.
// It check that all values are valid.
func ParseAndValidate(s string) (ConnString, error) ***REMOVED***
	p := parser***REMOVED***dnsResolver: dns.DefaultResolver***REMOVED***
	err := p.parse(s)
	if err != nil ***REMOVED***
		return p.ConnString, internal.WrapErrorf(err, "error parsing uri")
	***REMOVED***
	err = p.ConnString.Validate()
	if err != nil ***REMOVED***
		return p.ConnString, internal.WrapErrorf(err, "error validating uri")
	***REMOVED***
	return p.ConnString, nil
***REMOVED***

// Parse parses the provided URI into a ConnString object
// but does not check that all values are valid. Use `ConnString.Validate()`
// to run the validation checks separately.
func Parse(s string) (ConnString, error) ***REMOVED***
	p := parser***REMOVED***dnsResolver: dns.DefaultResolver***REMOVED***
	err := p.parse(s)
	if err != nil ***REMOVED***
		err = internal.WrapErrorf(err, "error parsing uri")
	***REMOVED***
	return p.ConnString, err
***REMOVED***

// ConnString represents a connection string to mongodb.
type ConnString struct ***REMOVED***
	Original                           string
	AppName                            string
	AuthMechanism                      string
	AuthMechanismProperties            map[string]string
	AuthMechanismPropertiesSet         bool
	AuthSource                         string
	AuthSourceSet                      bool
	Compressors                        []string
	Connect                            ConnectMode
	ConnectSet                         bool
	DirectConnection                   bool
	DirectConnectionSet                bool
	ConnectTimeout                     time.Duration
	ConnectTimeoutSet                  bool
	Database                           string
	HeartbeatInterval                  time.Duration
	HeartbeatIntervalSet               bool
	Hosts                              []string
	J                                  bool
	JSet                               bool
	LoadBalanced                       bool
	LoadBalancedSet                    bool
	LocalThreshold                     time.Duration
	LocalThresholdSet                  bool
	MaxConnIdleTime                    time.Duration
	MaxConnIdleTimeSet                 bool
	MaxPoolSize                        uint64
	MaxPoolSizeSet                     bool
	MinPoolSize                        uint64
	MinPoolSizeSet                     bool
	MaxConnecting                      uint64
	MaxConnectingSet                   bool
	Password                           string
	PasswordSet                        bool
	ReadConcernLevel                   string
	ReadPreference                     string
	ReadPreferenceTagSets              []map[string]string
	RetryWrites                        bool
	RetryWritesSet                     bool
	RetryReads                         bool
	RetryReadsSet                      bool
	MaxStaleness                       time.Duration
	MaxStalenessSet                    bool
	ReplicaSet                         string
	Scheme                             string
	ServerSelectionTimeout             time.Duration
	ServerSelectionTimeoutSet          bool
	SocketTimeout                      time.Duration
	SocketTimeoutSet                   bool
	SRVMaxHosts                        int
	SRVServiceName                     string
	SSL                                bool
	SSLSet                             bool
	SSLClientCertificateKeyFile        string
	SSLClientCertificateKeyFileSet     bool
	SSLClientCertificateKeyPassword    func() string
	SSLClientCertificateKeyPasswordSet bool
	SSLCertificateFile                 string
	SSLCertificateFileSet              bool
	SSLPrivateKeyFile                  string
	SSLPrivateKeyFileSet               bool
	SSLInsecure                        bool
	SSLInsecureSet                     bool
	SSLCaFile                          string
	SSLCaFileSet                       bool
	SSLDisableOCSPEndpointCheck        bool
	SSLDisableOCSPEndpointCheckSet     bool
	Timeout                            time.Duration
	TimeoutSet                         bool
	WString                            string
	WNumber                            int
	WNumberSet                         bool
	Username                           string
	UsernameSet                        bool
	ZlibLevel                          int
	ZlibLevelSet                       bool
	ZstdLevel                          int
	ZstdLevelSet                       bool

	WTimeout              time.Duration
	WTimeoutSet           bool
	WTimeoutSetFromOption bool

	Options        map[string][]string
	UnknownOptions map[string][]string
***REMOVED***

func (u *ConnString) String() string ***REMOVED***
	return u.Original
***REMOVED***

// HasAuthParameters returns true if this ConnString has any authentication parameters set and therefore represents
// a request for authentication.
func (u *ConnString) HasAuthParameters() bool ***REMOVED***
	// Check all auth parameters except for AuthSource because an auth source without other credentials is semantically
	// valid and must not be interpreted as a request for authentication.
	return u.AuthMechanism != "" || u.AuthMechanismProperties != nil || u.UsernameSet || u.PasswordSet
***REMOVED***

// Validate checks that the Auth and SSL parameters are valid values.
func (u *ConnString) Validate() error ***REMOVED***
	p := parser***REMOVED***
		dnsResolver: dns.DefaultResolver,
		ConnString:  *u,
	***REMOVED***
	return p.validate()
***REMOVED***

// ConnectMode informs the driver on how to connect
// to the server.
type ConnectMode uint8

var _ fmt.Stringer = ConnectMode(0)

// ConnectMode constants.
const (
	AutoConnect ConnectMode = iota
	SingleConnect
)

// String implements the fmt.Stringer interface.
func (c ConnectMode) String() string ***REMOVED***
	switch c ***REMOVED***
	case AutoConnect:
		return "automatic"
	case SingleConnect:
		return "direct"
	default:
		return "unknown"
	***REMOVED***
***REMOVED***

// Scheme constants
const (
	SchemeMongoDB    = "mongodb"
	SchemeMongoDBSRV = "mongodb+srv"
)

type parser struct ***REMOVED***
	ConnString

	dnsResolver *dns.Resolver
	tlsssl      *bool // used to determine if tls and ssl options are both specified and set differently.
***REMOVED***

func (p *parser) parse(original string) error ***REMOVED***
	p.Original = original
	uri := original

	var err error
	if strings.HasPrefix(uri, SchemeMongoDBSRV+"://") ***REMOVED***
		p.Scheme = SchemeMongoDBSRV
		// remove the scheme
		uri = uri[len(SchemeMongoDBSRV)+3:]
	***REMOVED*** else if strings.HasPrefix(uri, SchemeMongoDB+"://") ***REMOVED***
		p.Scheme = SchemeMongoDB
		// remove the scheme
		uri = uri[len(SchemeMongoDB)+3:]
	***REMOVED*** else ***REMOVED***
		return fmt.Errorf("scheme must be \"mongodb\" or \"mongodb+srv\"")
	***REMOVED***

	if idx := strings.Index(uri, "@"); idx != -1 ***REMOVED***
		userInfo := uri[:idx]
		uri = uri[idx+1:]

		username := userInfo
		var password string

		if idx := strings.Index(userInfo, ":"); idx != -1 ***REMOVED***
			username = userInfo[:idx]
			password = userInfo[idx+1:]
			p.PasswordSet = true
		***REMOVED***

		// Validate and process the username.
		if strings.Contains(username, "/") ***REMOVED***
			return fmt.Errorf("unescaped slash in username")
		***REMOVED***
		p.Username, err = url.PathUnescape(username)
		if err != nil ***REMOVED***
			return internal.WrapErrorf(err, "invalid username")
		***REMOVED***
		p.UsernameSet = true

		// Validate and process the password.
		if strings.Contains(password, ":") ***REMOVED***
			return fmt.Errorf("unescaped colon in password")
		***REMOVED***
		if strings.Contains(password, "/") ***REMOVED***
			return fmt.Errorf("unescaped slash in password")
		***REMOVED***
		p.Password, err = url.PathUnescape(password)
		if err != nil ***REMOVED***
			return internal.WrapErrorf(err, "invalid password")
		***REMOVED***
	***REMOVED***

	// fetch the hosts field
	hosts := uri
	if idx := strings.IndexAny(uri, "/?@"); idx != -1 ***REMOVED***
		if uri[idx] == '@' ***REMOVED***
			return fmt.Errorf("unescaped @ sign in user info")
		***REMOVED***
		if uri[idx] == '?' ***REMOVED***
			return fmt.Errorf("must have a / before the query ?")
		***REMOVED***
		hosts = uri[:idx]
	***REMOVED***

	parsedHosts := strings.Split(hosts, ",")
	uri = uri[len(hosts):]
	extractedDatabase, err := extractDatabaseFromURI(uri)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	uri = extractedDatabase.uri
	p.Database = extractedDatabase.db

	// grab connection arguments from URI
	connectionArgsFromQueryString, err := extractQueryArgsFromURI(uri)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// grab connection arguments from TXT record and enable SSL if "mongodb+srv://"
	var connectionArgsFromTXT []string
	if p.Scheme == SchemeMongoDBSRV ***REMOVED***
		connectionArgsFromTXT, err = p.dnsResolver.GetConnectionArgsFromTXT(hosts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// SSL is enabled by default for SRV, but can be manually disabled with "ssl=false".
		p.SSL = true
		p.SSLSet = true
	***REMOVED***

	// add connection arguments from URI and TXT records to connstring
	connectionArgPairs := make([]string, 0, len(connectionArgsFromTXT)+len(connectionArgsFromQueryString))
	connectionArgPairs = append(connectionArgPairs, connectionArgsFromTXT...)
	connectionArgPairs = append(connectionArgPairs, connectionArgsFromQueryString...)

	for _, pair := range connectionArgPairs ***REMOVED***
		err := p.addOption(pair)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// do SRV lookup if "mongodb+srv://"
	if p.Scheme == SchemeMongoDBSRV ***REMOVED***
		parsedHosts, err = p.dnsResolver.ParseHosts(hosts, p.SRVServiceName, true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// If p.SRVMaxHosts is non-zero and is less than the number of hosts, randomly
		// select SRVMaxHosts hosts from parsedHosts.
		if p.SRVMaxHosts > 0 && p.SRVMaxHosts < len(parsedHosts) ***REMOVED***
			random.Shuffle(len(parsedHosts), func(i, j int) ***REMOVED***
				parsedHosts[i], parsedHosts[j] = parsedHosts[j], parsedHosts[i]
			***REMOVED***)
			parsedHosts = parsedHosts[:p.SRVMaxHosts]
		***REMOVED***
	***REMOVED***

	for _, host := range parsedHosts ***REMOVED***
		err = p.addHost(host)
		if err != nil ***REMOVED***
			return internal.WrapErrorf(err, "invalid host %q", host)
		***REMOVED***
	***REMOVED***
	if len(p.Hosts) == 0 ***REMOVED***
		return fmt.Errorf("must have at least 1 host")
	***REMOVED***

	err = p.setDefaultAuthParams(extractedDatabase.db)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If WTimeout was set from manual options passed in, set WTImeoutSet to true.
	if p.WTimeoutSetFromOption ***REMOVED***
		p.WTimeoutSet = true
	***REMOVED***

	return nil
***REMOVED***

func (p *parser) validate() error ***REMOVED***
	var err error

	err = p.validateAuth()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = p.validateSSL(); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Check for invalid write concern (i.e. w=0 and j=true)
	if p.WNumberSet && p.WNumber == 0 && p.JSet && p.J ***REMOVED***
		return writeconcern.ErrInconsistent
	***REMOVED***

	// Check for invalid use of direct connections.
	if (p.ConnectSet && p.Connect == SingleConnect) || (p.DirectConnectionSet && p.DirectConnection) ***REMOVED***
		if len(p.Hosts) > 1 ***REMOVED***
			return errors.New("a direct connection cannot be made if multiple hosts are specified")
		***REMOVED***
		if p.Scheme == SchemeMongoDBSRV ***REMOVED***
			return errors.New("a direct connection cannot be made if an SRV URI is used")
		***REMOVED***
		if p.LoadBalancedSet && p.LoadBalanced ***REMOVED***
			return internal.ErrLoadBalancedWithDirectConnection
		***REMOVED***
	***REMOVED***

	// Validation for load-balanced mode.
	if p.LoadBalancedSet && p.LoadBalanced ***REMOVED***
		if len(p.Hosts) > 1 ***REMOVED***
			return internal.ErrLoadBalancedWithMultipleHosts
		***REMOVED***
		if p.ReplicaSet != "" ***REMOVED***
			return internal.ErrLoadBalancedWithReplicaSet
		***REMOVED***
	***REMOVED***

	// Check for invalid use of SRVMaxHosts.
	if p.SRVMaxHosts > 0 ***REMOVED***
		if p.ReplicaSet != "" ***REMOVED***
			return internal.ErrSRVMaxHostsWithReplicaSet
		***REMOVED***
		if p.LoadBalanced ***REMOVED***
			return internal.ErrSRVMaxHostsWithLoadBalanced
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (p *parser) setDefaultAuthParams(dbName string) error ***REMOVED***
	// We do this check here rather than in validateAuth because this function is called as part of parsing and sets
	// the value of AuthSource if authentication is enabled.
	if p.AuthSourceSet && p.AuthSource == "" ***REMOVED***
		return errors.New("authSource must be non-empty when supplied in a URI")
	***REMOVED***

	switch strings.ToLower(p.AuthMechanism) ***REMOVED***
	case "plain":
		if p.AuthSource == "" ***REMOVED***
			p.AuthSource = dbName
			if p.AuthSource == "" ***REMOVED***
				p.AuthSource = "$external"
			***REMOVED***
		***REMOVED***
	case "gssapi":
		if p.AuthMechanismProperties == nil ***REMOVED***
			p.AuthMechanismProperties = map[string]string***REMOVED***
				"SERVICE_NAME": "mongodb",
			***REMOVED***
		***REMOVED*** else if v, ok := p.AuthMechanismProperties["SERVICE_NAME"]; !ok || v == "" ***REMOVED***
			p.AuthMechanismProperties["SERVICE_NAME"] = "mongodb"
		***REMOVED***
		fallthrough
	case "mongodb-aws", "mongodb-x509":
		if p.AuthSource == "" ***REMOVED***
			p.AuthSource = "$external"
		***REMOVED*** else if p.AuthSource != "$external" ***REMOVED***
			return fmt.Errorf("auth source must be $external")
		***REMOVED***
	case "mongodb-cr":
		fallthrough
	case "scram-sha-1":
		fallthrough
	case "scram-sha-256":
		if p.AuthSource == "" ***REMOVED***
			p.AuthSource = dbName
			if p.AuthSource == "" ***REMOVED***
				p.AuthSource = "admin"
			***REMOVED***
		***REMOVED***
	case "":
		// Only set auth source if there is a request for authentication via non-empty credentials.
		if p.AuthSource == "" && (p.AuthMechanismProperties != nil || p.Username != "" || p.PasswordSet) ***REMOVED***
			p.AuthSource = dbName
			if p.AuthSource == "" ***REMOVED***
				p.AuthSource = "admin"
			***REMOVED***
		***REMOVED***
	default:
		return fmt.Errorf("invalid auth mechanism")
	***REMOVED***
	return nil
***REMOVED***

func (p *parser) validateAuth() error ***REMOVED***
	switch strings.ToLower(p.AuthMechanism) ***REMOVED***
	case "mongodb-cr":
		if p.Username == "" ***REMOVED***
			return fmt.Errorf("username required for MONGO-CR")
		***REMOVED***
		if p.Password == "" ***REMOVED***
			return fmt.Errorf("password required for MONGO-CR")
		***REMOVED***
		if p.AuthMechanismProperties != nil ***REMOVED***
			return fmt.Errorf("MONGO-CR cannot have mechanism properties")
		***REMOVED***
	case "mongodb-x509":
		if p.Password != "" ***REMOVED***
			return fmt.Errorf("password cannot be specified for MONGO-X509")
		***REMOVED***
		if p.AuthMechanismProperties != nil ***REMOVED***
			return fmt.Errorf("MONGO-X509 cannot have mechanism properties")
		***REMOVED***
	case "mongodb-aws":
		if p.Username != "" && p.Password == "" ***REMOVED***
			return fmt.Errorf("username without password is invalid for MONGODB-AWS")
		***REMOVED***
		if p.Username == "" && p.Password != "" ***REMOVED***
			return fmt.Errorf("password without username is invalid for MONGODB-AWS")
		***REMOVED***
		var token bool
		for k := range p.AuthMechanismProperties ***REMOVED***
			if k != "AWS_SESSION_TOKEN" ***REMOVED***
				return fmt.Errorf("invalid auth property for MONGODB-AWS")
			***REMOVED***
			token = true
		***REMOVED***
		if token && p.Username == "" && p.Password == "" ***REMOVED***
			return fmt.Errorf("token without username and password is invalid for MONGODB-AWS")
		***REMOVED***
	case "gssapi":
		if p.Username == "" ***REMOVED***
			return fmt.Errorf("username required for GSSAPI")
		***REMOVED***
		for k := range p.AuthMechanismProperties ***REMOVED***
			if k != "SERVICE_NAME" && k != "CANONICALIZE_HOST_NAME" && k != "SERVICE_REALM" ***REMOVED***
				return fmt.Errorf("invalid auth property for GSSAPI")
			***REMOVED***
		***REMOVED***
	case "plain":
		if p.Username == "" ***REMOVED***
			return fmt.Errorf("username required for PLAIN")
		***REMOVED***
		if p.Password == "" ***REMOVED***
			return fmt.Errorf("password required for PLAIN")
		***REMOVED***
		if p.AuthMechanismProperties != nil ***REMOVED***
			return fmt.Errorf("PLAIN cannot have mechanism properties")
		***REMOVED***
	case "scram-sha-1":
		if p.Username == "" ***REMOVED***
			return fmt.Errorf("username required for SCRAM-SHA-1")
		***REMOVED***
		if p.Password == "" ***REMOVED***
			return fmt.Errorf("password required for SCRAM-SHA-1")
		***REMOVED***
		if p.AuthMechanismProperties != nil ***REMOVED***
			return fmt.Errorf("SCRAM-SHA-1 cannot have mechanism properties")
		***REMOVED***
	case "scram-sha-256":
		if p.Username == "" ***REMOVED***
			return fmt.Errorf("username required for SCRAM-SHA-256")
		***REMOVED***
		if p.Password == "" ***REMOVED***
			return fmt.Errorf("password required for SCRAM-SHA-256")
		***REMOVED***
		if p.AuthMechanismProperties != nil ***REMOVED***
			return fmt.Errorf("SCRAM-SHA-256 cannot have mechanism properties")
		***REMOVED***
	case "":
		if p.UsernameSet && p.Username == "" ***REMOVED***
			return fmt.Errorf("username required if URI contains user info")
		***REMOVED***
	default:
		return fmt.Errorf("invalid auth mechanism")
	***REMOVED***
	return nil
***REMOVED***

func (p *parser) validateSSL() error ***REMOVED***
	if !p.SSL ***REMOVED***
		return nil
	***REMOVED***

	if p.SSLClientCertificateKeyFileSet ***REMOVED***
		if p.SSLCertificateFileSet || p.SSLPrivateKeyFileSet ***REMOVED***
			return errors.New("the sslClientCertificateKeyFile/tlsCertificateKeyFile URI option cannot be provided " +
				"along with tlsCertificateFile or tlsPrivateKeyFile")
		***REMOVED***
		return nil
	***REMOVED***
	if p.SSLCertificateFileSet && !p.SSLPrivateKeyFileSet ***REMOVED***
		return errors.New("the tlsPrivateKeyFile URI option must be provided if the tlsCertificateFile option is specified")
	***REMOVED***
	if p.SSLPrivateKeyFileSet && !p.SSLCertificateFileSet ***REMOVED***
		return errors.New("the tlsCertificateFile URI option must be provided if the tlsPrivateKeyFile option is specified")
	***REMOVED***

	if p.SSLInsecureSet && p.SSLDisableOCSPEndpointCheckSet ***REMOVED***
		return errors.New("the sslInsecure/tlsInsecure URI option cannot be provided along with " +
			"tlsDisableOCSPEndpointCheck ")
	***REMOVED***
	return nil
***REMOVED***

func (p *parser) addHost(host string) error ***REMOVED***
	if host == "" ***REMOVED***
		return nil
	***REMOVED***
	host, err := url.QueryUnescape(host)
	if err != nil ***REMOVED***
		return internal.WrapErrorf(err, "invalid host %q", host)
	***REMOVED***

	_, port, err := net.SplitHostPort(host)
	// this is unfortunate that SplitHostPort actually requires
	// a port to exist.
	if err != nil ***REMOVED***
		if addrError, ok := err.(*net.AddrError); !ok || addrError.Err != "missing port in address" ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if port != "" ***REMOVED***
		d, err := strconv.Atoi(port)
		if err != nil ***REMOVED***
			return internal.WrapErrorf(err, "port must be an integer")
		***REMOVED***
		if d <= 0 || d >= 65536 ***REMOVED***
			return fmt.Errorf("port must be in the range [1, 65535]")
		***REMOVED***
	***REMOVED***
	p.Hosts = append(p.Hosts, host)
	return nil
***REMOVED***

func (p *parser) addOption(pair string) error ***REMOVED***
	kv := strings.SplitN(pair, "=", 2)
	if len(kv) != 2 || kv[0] == "" ***REMOVED***
		return fmt.Errorf("invalid option")
	***REMOVED***

	key, err := url.QueryUnescape(kv[0])
	if err != nil ***REMOVED***
		return internal.WrapErrorf(err, "invalid option key %q", kv[0])
	***REMOVED***

	value, err := url.QueryUnescape(kv[1])
	if err != nil ***REMOVED***
		return internal.WrapErrorf(err, "invalid option value %q", kv[1])
	***REMOVED***

	lowerKey := strings.ToLower(key)
	switch lowerKey ***REMOVED***
	case "appname":
		p.AppName = value
	case "authmechanism":
		p.AuthMechanism = value
	case "authmechanismproperties":
		p.AuthMechanismProperties = make(map[string]string)
		pairs := strings.Split(value, ",")
		for _, pair := range pairs ***REMOVED***
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) != 2 || kv[0] == "" ***REMOVED***
				return fmt.Errorf("invalid authMechanism property")
			***REMOVED***
			p.AuthMechanismProperties[kv[0]] = kv[1]
		***REMOVED***
		p.AuthMechanismPropertiesSet = true
	case "authsource":
		p.AuthSource = value
		p.AuthSourceSet = true
	case "compressors":
		compressors := strings.Split(value, ",")
		if len(compressors) < 1 ***REMOVED***
			return fmt.Errorf("must have at least 1 compressor")
		***REMOVED***
		p.Compressors = compressors
	case "connect":
		switch strings.ToLower(value) ***REMOVED***
		case "automatic":
		case "direct":
			p.Connect = SingleConnect
		default:
			return fmt.Errorf("invalid 'connect' value: %q", value)
		***REMOVED***
		if p.DirectConnectionSet ***REMOVED***
			expectedValue := p.Connect == SingleConnect // directConnection should be true if connect=direct
			if p.DirectConnection != expectedValue ***REMOVED***
				return fmt.Errorf("options connect=%q and directConnection=%v conflict", value, p.DirectConnection)
			***REMOVED***
		***REMOVED***

		p.ConnectSet = true
	case "directconnection":
		switch strings.ToLower(value) ***REMOVED***
		case "true":
			p.DirectConnection = true
		case "false":
		default:
			return fmt.Errorf("invalid 'directConnection' value: %q", value)
		***REMOVED***

		if p.ConnectSet ***REMOVED***
			expectedValue := AutoConnect
			if p.DirectConnection ***REMOVED***
				expectedValue = SingleConnect
			***REMOVED***

			if p.Connect != expectedValue ***REMOVED***
				return fmt.Errorf("options connect=%q and directConnection=%q conflict", p.Connect, value)
			***REMOVED***
		***REMOVED***
		p.DirectConnectionSet = true
	case "connecttimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.ConnectTimeout = time.Duration(n) * time.Millisecond
		p.ConnectTimeoutSet = true
	case "heartbeatintervalms", "heartbeatfrequencyms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.HeartbeatInterval = time.Duration(n) * time.Millisecond
		p.HeartbeatIntervalSet = true
	case "journal":
		switch value ***REMOVED***
		case "true":
			p.J = true
		case "false":
			p.J = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		p.JSet = true
	case "loadbalanced":
		switch value ***REMOVED***
		case "true":
			p.LoadBalanced = true
		case "false":
			p.LoadBalanced = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		p.LoadBalancedSet = true
	case "localthresholdms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.LocalThreshold = time.Duration(n) * time.Millisecond
		p.LocalThresholdSet = true
	case "maxidletimems":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.MaxConnIdleTime = time.Duration(n) * time.Millisecond
		p.MaxConnIdleTimeSet = true
	case "maxpoolsize":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.MaxPoolSize = uint64(n)
		p.MaxPoolSizeSet = true
	case "minpoolsize":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.MinPoolSize = uint64(n)
		p.MinPoolSizeSet = true
	case "maxconnecting":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.MaxConnecting = uint64(n)
		p.MaxConnectingSet = true
	case "readconcernlevel":
		p.ReadConcernLevel = value
	case "readpreference":
		p.ReadPreference = value
	case "readpreferencetags":
		if value == "" ***REMOVED***
			// If "readPreferenceTags=" is supplied, append an empty map to tag sets to
			// represent a wild-card.
			p.ReadPreferenceTagSets = append(p.ReadPreferenceTagSets, map[string]string***REMOVED******REMOVED***)
			break
		***REMOVED***

		tags := make(map[string]string)
		items := strings.Split(value, ",")
		for _, item := range items ***REMOVED***
			parts := strings.Split(item, ":")
			if len(parts) != 2 ***REMOVED***
				return fmt.Errorf("invalid value for %q: %q", key, value)
			***REMOVED***
			tags[parts[0]] = parts[1]
		***REMOVED***
		p.ReadPreferenceTagSets = append(p.ReadPreferenceTagSets, tags)
	case "maxstaleness", "maxstalenessseconds":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.MaxStaleness = time.Duration(n) * time.Second
		p.MaxStalenessSet = true
	case "replicaset":
		p.ReplicaSet = value
	case "retrywrites":
		switch value ***REMOVED***
		case "true":
			p.RetryWrites = true
		case "false":
			p.RetryWrites = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		p.RetryWritesSet = true
	case "retryreads":
		switch value ***REMOVED***
		case "true":
			p.RetryReads = true
		case "false":
			p.RetryReads = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		p.RetryReadsSet = true
	case "serverselectiontimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.ServerSelectionTimeout = time.Duration(n) * time.Millisecond
		p.ServerSelectionTimeoutSet = true
	case "sockettimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.SocketTimeout = time.Duration(n) * time.Millisecond
		p.SocketTimeoutSet = true
	case "srvmaxhosts":
		// srvMaxHosts can only be set on URIs with the "mongodb+srv" scheme
		if p.Scheme != SchemeMongoDBSRV ***REMOVED***
			return fmt.Errorf("cannot specify srvMaxHosts on non-SRV URI")
		***REMOVED***

		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.SRVMaxHosts = n
	case "srvservicename":
		// srvServiceName can only be set on URIs with the "mongodb+srv" scheme
		if p.Scheme != SchemeMongoDBSRV ***REMOVED***
			return fmt.Errorf("cannot specify srvServiceName on non-SRV URI")
		***REMOVED***

		// srvServiceName must be between 1 and 62 characters according to
		// our specification. Empty service names are not valid, and the service
		// name (including prepended underscore) should not exceed the 63 character
		// limit for DNS query subdomains.
		if len(value) < 1 || len(value) > 62 ***REMOVED***
			return fmt.Errorf("srvServiceName value must be between 1 and 62 characters")
		***REMOVED***
		p.SRVServiceName = value
	case "ssl", "tls":
		switch value ***REMOVED***
		case "true":
			p.SSL = true
		case "false":
			p.SSL = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		if p.tlsssl != nil && *p.tlsssl != p.SSL ***REMOVED***
			return errors.New("tls and ssl options, when both specified, must be equivalent")
		***REMOVED***

		p.tlsssl = new(bool)
		*p.tlsssl = p.SSL

		p.SSLSet = true
	case "sslclientcertificatekeyfile", "tlscertificatekeyfile":
		p.SSL = true
		p.SSLSet = true
		p.SSLClientCertificateKeyFile = value
		p.SSLClientCertificateKeyFileSet = true
	case "sslclientcertificatekeypassword", "tlscertificatekeyfilepassword":
		p.SSLClientCertificateKeyPassword = func() string ***REMOVED*** return value ***REMOVED***
		p.SSLClientCertificateKeyPasswordSet = true
	case "tlscertificatefile":
		p.SSL = true
		p.SSLSet = true
		p.SSLCertificateFile = value
		p.SSLCertificateFileSet = true
	case "tlsprivatekeyfile":
		p.SSL = true
		p.SSLSet = true
		p.SSLPrivateKeyFile = value
		p.SSLPrivateKeyFileSet = true
	case "sslinsecure", "tlsinsecure":
		switch value ***REMOVED***
		case "true":
			p.SSLInsecure = true
		case "false":
			p.SSLInsecure = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		p.SSLInsecureSet = true
	case "sslcertificateauthorityfile", "tlscafile":
		p.SSL = true
		p.SSLSet = true
		p.SSLCaFile = value
		p.SSLCaFileSet = true
	case "timeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.Timeout = time.Duration(n) * time.Millisecond
		p.TimeoutSet = true
	case "tlsdisableocspendpointcheck":
		p.SSL = true
		p.SSLSet = true

		switch value ***REMOVED***
		case "true":
			p.SSLDisableOCSPEndpointCheck = true
		case "false":
			p.SSLDisableOCSPEndpointCheck = false
		default:
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.SSLDisableOCSPEndpointCheckSet = true
	case "w":
		if w, err := strconv.Atoi(value); err == nil ***REMOVED***
			if w < 0 ***REMOVED***
				return fmt.Errorf("invalid value for %q: %q", key, value)
			***REMOVED***

			p.WNumber = w
			p.WNumberSet = true
			p.WString = ""
			break
		***REMOVED***

		p.WString = value
		p.WNumberSet = false

	case "wtimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.WTimeout = time.Duration(n) * time.Millisecond
		p.WTimeoutSet = true
	case "wtimeout":
		// Defer to wtimeoutms, but not to a manually-set option.
		if p.WTimeoutSet ***REMOVED***
			break
		***REMOVED***
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***
		p.WTimeout = time.Duration(n) * time.Millisecond
	case "zlibcompressionlevel":
		level, err := strconv.Atoi(value)
		if err != nil || (level < -1 || level > 9) ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		if level == -1 ***REMOVED***
			level = wiremessage.DefaultZlibLevel
		***REMOVED***
		p.ZlibLevel = level
		p.ZlibLevelSet = true
	case "zstdcompressionlevel":
		const maxZstdLevel = 22 // https://github.com/facebook/zstd/blob/a880ca239b447968493dd2fed3850e766d6305cc/contrib/linux-kernel/lib/zstd/compress.c#L3291
		level, err := strconv.Atoi(value)
		if err != nil || (level < -1 || level > maxZstdLevel) ***REMOVED***
			return fmt.Errorf("invalid value for %q: %q", key, value)
		***REMOVED***

		if level == -1 ***REMOVED***
			level = wiremessage.DefaultZstdLevel
		***REMOVED***
		p.ZstdLevel = level
		p.ZstdLevelSet = true
	default:
		if p.UnknownOptions == nil ***REMOVED***
			p.UnknownOptions = make(map[string][]string)
		***REMOVED***
		p.UnknownOptions[lowerKey] = append(p.UnknownOptions[lowerKey], value)
	***REMOVED***

	if p.Options == nil ***REMOVED***
		p.Options = make(map[string][]string)
	***REMOVED***
	p.Options[lowerKey] = append(p.Options[lowerKey], value)

	return nil
***REMOVED***

func extractQueryArgsFromURI(uri string) ([]string, error) ***REMOVED***
	if len(uri) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	if uri[0] != '?' ***REMOVED***
		return nil, errors.New("must have a ? separator between path and query")
	***REMOVED***

	uri = uri[1:]
	if len(uri) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***
	return strings.FieldsFunc(uri, func(r rune) bool ***REMOVED*** return r == ';' || r == '&' ***REMOVED***), nil

***REMOVED***

type extractedDatabase struct ***REMOVED***
	uri string
	db  string
***REMOVED***

// extractDatabaseFromURI is a helper function to retrieve information about
// the database from the passed in URI. It accepts as an argument the currently
// parsed URI and returns the remainder of the uri, the database it found,
// and any error it encounters while parsing.
func extractDatabaseFromURI(uri string) (extractedDatabase, error) ***REMOVED***
	if len(uri) == 0 ***REMOVED***
		return extractedDatabase***REMOVED******REMOVED***, nil
	***REMOVED***

	if uri[0] != '/' ***REMOVED***
		return extractedDatabase***REMOVED******REMOVED***, errors.New("must have a / separator between hosts and path")
	***REMOVED***

	uri = uri[1:]
	if len(uri) == 0 ***REMOVED***
		return extractedDatabase***REMOVED******REMOVED***, nil
	***REMOVED***

	database := uri
	if idx := strings.IndexRune(uri, '?'); idx != -1 ***REMOVED***
		database = uri[:idx]
	***REMOVED***

	escapedDatabase, err := url.QueryUnescape(database)
	if err != nil ***REMOVED***
		return extractedDatabase***REMOVED******REMOVED***, internal.WrapErrorf(err, "invalid database %q", database)
	***REMOVED***

	uri = uri[len(database):]

	return extractedDatabase***REMOVED***
		uri: uri,
		db:  escapedDatabase,
	***REMOVED***, nil
***REMOVED***
