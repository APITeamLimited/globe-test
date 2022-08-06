package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9/internal/pool"
)

// Limiter is the interface of a rate limiter or a circuit breaker.
type Limiter interface ***REMOVED***
	// Allow returns nil if operation is allowed or an error otherwise.
	// If operation is allowed client must ReportResult of the operation
	// whether it is a success or a failure.
	Allow() error
	// ReportResult reports the result of the previously allowed operation.
	// nil indicates a success, non-nil error usually indicates a failure.
	ReportResult(result error)
***REMOVED***

// Options keeps the settings to setup redis connection.
type Options struct ***REMOVED***
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	// Hook that is called when new connection is established.
	OnConnect func(ctx context.Context, cn *Conn) error

	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string
	// CredentialsProvider allows the username and password to be updated
	// before reconnecting. It should return the current username and password.
	CredentialsProvider func() (username string, password string)

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher overhead compared to lifo.
	PoolFIFO bool
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Maximum number of idle connections.
	MaxIdleConns int
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	ConnMaxLifetime time.Duration

	// Enables read only queries on slave nodes.
	readOnly bool

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	// Limiter interface used to implemented circuit breaker or rate limiter.
	Limiter Limiter
***REMOVED***

func (opt *Options) init() ***REMOVED***
	if opt.Addr == "" ***REMOVED***
		opt.Addr = "localhost:6379"
	***REMOVED***
	if opt.Network == "" ***REMOVED***
		if strings.HasPrefix(opt.Addr, "/") ***REMOVED***
			opt.Network = "unix"
		***REMOVED*** else ***REMOVED***
			opt.Network = "tcp"
		***REMOVED***
	***REMOVED***
	if opt.DialTimeout == 0 ***REMOVED***
		opt.DialTimeout = 5 * time.Second
	***REMOVED***
	if opt.Dialer == nil ***REMOVED***
		opt.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) ***REMOVED***
			netDialer := &net.Dialer***REMOVED***
				Timeout:   opt.DialTimeout,
				KeepAlive: 5 * time.Minute,
			***REMOVED***
			if opt.TLSConfig == nil ***REMOVED***
				return netDialer.DialContext(ctx, network, addr)
			***REMOVED***
			return tls.DialWithDialer(netDialer, network, addr, opt.TLSConfig)
		***REMOVED***
	***REMOVED***
	if opt.PoolSize == 0 ***REMOVED***
		opt.PoolSize = 10 * runtime.GOMAXPROCS(0)
	***REMOVED***
	switch opt.ReadTimeout ***REMOVED***
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 3 * time.Second
	***REMOVED***
	switch opt.WriteTimeout ***REMOVED***
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	***REMOVED***
	if opt.PoolTimeout == 0 ***REMOVED***
		opt.PoolTimeout = opt.ReadTimeout + time.Second
	***REMOVED***
	if opt.ConnMaxIdleTime == 0 ***REMOVED***
		opt.ConnMaxIdleTime = 30 * time.Minute
	***REMOVED***

	if opt.MaxRetries == -1 ***REMOVED***
		opt.MaxRetries = 0
	***REMOVED*** else if opt.MaxRetries == 0 ***REMOVED***
		opt.MaxRetries = 3
	***REMOVED***
	switch opt.MinRetryBackoff ***REMOVED***
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	***REMOVED***
	switch opt.MaxRetryBackoff ***REMOVED***
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	***REMOVED***
***REMOVED***

func (opt *Options) clone() *Options ***REMOVED***
	clone := *opt
	return &clone
***REMOVED***

// ParseURL parses an URL into Options that can be used to connect to Redis.
// Scheme is required.
// There are two connection types: by tcp socket and by unix socket.
// Tcp connection:
//		redis://<user>:<password>@<host>:<port>/<db_number>
// Unix connection:
//		unix://<user>:<password>@</path/to/redis.sock>?db=<db_number>
// Most Option fields can be set using query parameters, with the following restrictions:
//	- field names are mapped using snake-case conversion: to set MaxRetries, use max_retries
//	- only scalar type fields are supported (bool, int, time.Duration)
//	- for time.Duration fields, values must be a valid input for time.ParseDuration();
//	  additionally a plain integer as value (i.e. without unit) is intepreted as seconds
//	- to disable a duration field, use value less than or equal to 0; to use the default
//	  value, leave the value blank or remove the parameter
//	- only the last value is interpreted if a parameter is given multiple times
//	- fields "network", "addr", "username" and "password" can only be set using other
//	  URL attributes (scheme, host, userinfo, resp.), query paremeters using these
//	  names will be treated as unknown parameters
//	- unknown parameter names will result in an error
// Examples:
//		redis://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2
//		is equivalent to:
//		&Options***REMOVED***
//			Network:     "tcp",
//			Addr:        "localhost:6789",
//			DB:          1,               // path "/3" was overridden by "&db=1"
//			DialTimeout: 3 * time.Second, // no time unit = seconds
//			ReadTimeout: 6 * time.Second,
//			MaxRetries:  2,
//		***REMOVED***
func ParseURL(redisURL string) (*Options, error) ***REMOVED***
	u, err := url.Parse(redisURL)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch u.Scheme ***REMOVED***
	case "redis", "rediss":
		return setupTCPConn(u)
	case "unix":
		return setupUnixConn(u)
	default:
		return nil, fmt.Errorf("redis: invalid URL scheme: %s", u.Scheme)
	***REMOVED***
***REMOVED***

func setupTCPConn(u *url.URL) (*Options, error) ***REMOVED***
	o := &Options***REMOVED***Network: "tcp"***REMOVED***

	o.Username, o.Password = getUserPassword(u)

	h, p, err := net.SplitHostPort(u.Host)
	if err != nil ***REMOVED***
		h = u.Host
	***REMOVED***
	if h == "" ***REMOVED***
		h = "localhost"
	***REMOVED***
	if p == "" ***REMOVED***
		p = "6379"
	***REMOVED***
	o.Addr = net.JoinHostPort(h, p)

	f := strings.FieldsFunc(u.Path, func(r rune) bool ***REMOVED***
		return r == '/'
	***REMOVED***)
	switch len(f) ***REMOVED***
	case 0:
		o.DB = 0
	case 1:
		if o.DB, err = strconv.Atoi(f[0]); err != nil ***REMOVED***
			return nil, fmt.Errorf("redis: invalid database number: %q", f[0])
		***REMOVED***
	default:
		return nil, fmt.Errorf("redis: invalid URL path: %s", u.Path)
	***REMOVED***

	if u.Scheme == "rediss" ***REMOVED***
		o.TLSConfig = &tls.Config***REMOVED***
			ServerName: h,
			MinVersion: tls.VersionTLS12,
		***REMOVED***
	***REMOVED***

	return setupConnParams(u, o)
***REMOVED***

func setupUnixConn(u *url.URL) (*Options, error) ***REMOVED***
	o := &Options***REMOVED***
		Network: "unix",
	***REMOVED***

	if strings.TrimSpace(u.Path) == "" ***REMOVED*** // path is required with unix connection
		return nil, errors.New("redis: empty unix socket path")
	***REMOVED***
	o.Addr = u.Path
	o.Username, o.Password = getUserPassword(u)
	return setupConnParams(u, o)
***REMOVED***

type queryOptions struct ***REMOVED***
	q   url.Values
	err error
***REMOVED***

func (o *queryOptions) has(name string) bool ***REMOVED***
	return len(o.q[name]) > 0
***REMOVED***

func (o *queryOptions) string(name string) string ***REMOVED***
	vs := o.q[name]
	if len(vs) == 0 ***REMOVED***
		return ""
	***REMOVED***
	delete(o.q, name) // enable detection of unknown parameters
	return vs[len(vs)-1]
***REMOVED***

func (o *queryOptions) int(name string) int ***REMOVED***
	s := o.string(name)
	if s == "" ***REMOVED***
		return 0
	***REMOVED***
	i, err := strconv.Atoi(s)
	if err == nil ***REMOVED***
		return i
	***REMOVED***
	if o.err == nil ***REMOVED***
		o.err = fmt.Errorf("redis: invalid %s number: %s", name, err)
	***REMOVED***
	return 0
***REMOVED***

func (o *queryOptions) duration(name string) time.Duration ***REMOVED***
	s := o.string(name)
	if s == "" ***REMOVED***
		return 0
	***REMOVED***
	// try plain number first
	if i, err := strconv.Atoi(s); err == nil ***REMOVED***
		if i <= 0 ***REMOVED***
			// disable timeouts
			return -1
		***REMOVED***
		return time.Duration(i) * time.Second
	***REMOVED***
	dur, err := time.ParseDuration(s)
	if err == nil ***REMOVED***
		return dur
	***REMOVED***
	if o.err == nil ***REMOVED***
		o.err = fmt.Errorf("redis: invalid %s duration: %w", name, err)
	***REMOVED***
	return 0
***REMOVED***

func (o *queryOptions) bool(name string) bool ***REMOVED***
	switch s := o.string(name); s ***REMOVED***
	case "true", "1":
		return true
	case "false", "0", "":
		return false
	default:
		if o.err == nil ***REMOVED***
			o.err = fmt.Errorf("redis: invalid %s boolean: expected true/false/1/0 or an empty string, got %q", name, s)
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (o *queryOptions) remaining() []string ***REMOVED***
	if len(o.q) == 0 ***REMOVED***
		return nil
	***REMOVED***
	keys := make([]string, 0, len(o.q))
	for k := range o.q ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	sort.Strings(keys)
	return keys
***REMOVED***

// setupConnParams converts query parameters in u to option value in o.
func setupConnParams(u *url.URL, o *Options) (*Options, error) ***REMOVED***
	q := queryOptions***REMOVED***q: u.Query()***REMOVED***

	// compat: a future major release may use q.int("db")
	if tmp := q.string("db"); tmp != "" ***REMOVED***
		db, err := strconv.Atoi(tmp)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("redis: invalid database number: %w", err)
		***REMOVED***
		o.DB = db
	***REMOVED***

	o.MaxRetries = q.int("max_retries")
	o.MinRetryBackoff = q.duration("min_retry_backoff")
	o.MaxRetryBackoff = q.duration("max_retry_backoff")
	o.DialTimeout = q.duration("dial_timeout")
	o.ReadTimeout = q.duration("read_timeout")
	o.WriteTimeout = q.duration("write_timeout")
	o.PoolFIFO = q.bool("pool_fifo")
	o.PoolSize = q.int("pool_size")
	o.PoolTimeout = q.duration("pool_timeout")
	o.MinIdleConns = q.int("min_idle_conns")
	o.MaxIdleConns = q.int("max_idle_conns")
	if q.has("conn_max_idle_time") ***REMOVED***
		o.ConnMaxIdleTime = q.duration("conn_max_idle_time")
	***REMOVED*** else ***REMOVED***
		o.ConnMaxIdleTime = q.duration("idle_timeout")
	***REMOVED***
	if q.has("conn_max_lifetime") ***REMOVED***
		o.ConnMaxLifetime = q.duration("conn_max_lifetime")
	***REMOVED*** else ***REMOVED***
		o.ConnMaxLifetime = q.duration("max_conn_age")
	***REMOVED***
	if q.err != nil ***REMOVED***
		return nil, q.err
	***REMOVED***

	// any parameters left?
	if r := q.remaining(); len(r) > 0 ***REMOVED***
		return nil, fmt.Errorf("redis: unexpected option: %s", strings.Join(r, ", "))
	***REMOVED***

	return o, nil
***REMOVED***

func getUserPassword(u *url.URL) (string, string) ***REMOVED***
	var user, password string
	if u.User != nil ***REMOVED***
		user = u.User.Username()
		if p, ok := u.User.Password(); ok ***REMOVED***
			password = p
		***REMOVED***
	***REMOVED***
	return user, password
***REMOVED***

func newConnPool(opt *Options) *pool.ConnPool ***REMOVED***
	return pool.NewConnPool(&pool.Options***REMOVED***
		Dialer: func(ctx context.Context) (net.Conn, error) ***REMOVED***
			return opt.Dialer(ctx, opt.Network, opt.Addr)
		***REMOVED***,
		PoolFIFO:        opt.PoolFIFO,
		PoolSize:        opt.PoolSize,
		PoolTimeout:     opt.PoolTimeout,
		MinIdleConns:    opt.MinIdleConns,
		MaxIdleConns:    opt.MaxIdleConns,
		ConnMaxIdleTime: opt.ConnMaxIdleTime,
		ConnMaxLifetime: opt.ConnMaxLifetime,
	***REMOVED***)
***REMOVED***
