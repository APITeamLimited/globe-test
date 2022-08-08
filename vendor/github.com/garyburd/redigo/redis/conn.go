// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	_ ConnWithTimeout = (*conn)(nil)
)

// conn is the low-level implementation of Conn
type conn struct ***REMOVED***
	// Shared
	mu      sync.Mutex
	pending int
	err     error
	conn    net.Conn

	// Read
	readTimeout time.Duration
	br          *bufio.Reader

	// Write
	writeTimeout time.Duration
	bw           *bufio.Writer

	// Scratch space for formatting argument length.
	// '*' or '$', length, "\r\n"
	lenScratch [32]byte

	// Scratch space for formatting integers and floats.
	numScratch [40]byte
***REMOVED***

// DialTimeout acts like Dial but takes timeouts for establishing the
// connection to the server, writing a command and reading a reply.
//
// Deprecated: Use Dial with options instead.
func DialTimeout(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (Conn, error) ***REMOVED***
	return Dial(network, address,
		DialConnectTimeout(connectTimeout),
		DialReadTimeout(readTimeout),
		DialWriteTimeout(writeTimeout))
***REMOVED***

// DialOption specifies an option for dialing a Redis server.
type DialOption struct ***REMOVED***
	f func(*dialOptions)
***REMOVED***

type dialOptions struct ***REMOVED***
	readTimeout  time.Duration
	writeTimeout time.Duration
	dialer       *net.Dialer
	dial         func(network, addr string) (net.Conn, error)
	db           int
	password     string
	useTLS       bool
	skipVerify   bool
	tlsConfig    *tls.Config
***REMOVED***

// DialReadTimeout specifies the timeout for reading a single command reply.
func DialReadTimeout(d time.Duration) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.readTimeout = d
	***REMOVED******REMOVED***
***REMOVED***

// DialWriteTimeout specifies the timeout for writing a single command.
func DialWriteTimeout(d time.Duration) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.writeTimeout = d
	***REMOVED******REMOVED***
***REMOVED***

// DialConnectTimeout specifies the timeout for connecting to the Redis server when
// no DialNetDial option is specified.
func DialConnectTimeout(d time.Duration) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.dialer.Timeout = d
	***REMOVED******REMOVED***
***REMOVED***

// DialKeepAlive specifies the keep-alive period for TCP connections to the Redis server
// when no DialNetDial option is specified.
// If zero, keep-alives are not enabled. If no DialKeepAlive option is specified then
// the default of 5 minutes is used to ensure that half-closed TCP sessions are detected.
func DialKeepAlive(d time.Duration) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.dialer.KeepAlive = d
	***REMOVED******REMOVED***
***REMOVED***

// DialNetDial specifies a custom dial function for creating TCP
// connections, otherwise a net.Dialer customized via the other options is used.
// DialNetDial overrides DialConnectTimeout and DialKeepAlive.
func DialNetDial(dial func(network, addr string) (net.Conn, error)) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.dial = dial
	***REMOVED******REMOVED***
***REMOVED***

// DialDatabase specifies the database to select when dialing a connection.
func DialDatabase(db int) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.db = db
	***REMOVED******REMOVED***
***REMOVED***

// DialPassword specifies the password to use when connecting to
// the Redis server.
func DialPassword(password string) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.password = password
	***REMOVED******REMOVED***
***REMOVED***

// DialTLSConfig specifies the config to use when a TLS connection is dialed.
// Has no effect when not dialing a TLS connection.
func DialTLSConfig(c *tls.Config) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.tlsConfig = c
	***REMOVED******REMOVED***
***REMOVED***

// DialTLSSkipVerify disables server name verification when connecting over
// TLS. Has no effect when not dialing a TLS connection.
func DialTLSSkipVerify(skip bool) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.skipVerify = skip
	***REMOVED******REMOVED***
***REMOVED***

// DialUseTLS specifies whether TLS should be used when connecting to the
// server. This option is ignore by DialURL.
func DialUseTLS(useTLS bool) DialOption ***REMOVED***
	return DialOption***REMOVED***func(do *dialOptions) ***REMOVED***
		do.useTLS = useTLS
	***REMOVED******REMOVED***
***REMOVED***

// Dial connects to the Redis server at the given network and
// address using the specified options.
func Dial(network, address string, options ...DialOption) (Conn, error) ***REMOVED***
	do := dialOptions***REMOVED***
		dialer: &net.Dialer***REMOVED***
			KeepAlive: time.Minute * 5,
		***REMOVED***,
	***REMOVED***
	for _, option := range options ***REMOVED***
		option.f(&do)
	***REMOVED***
	if do.dial == nil ***REMOVED***
		do.dial = do.dialer.Dial
	***REMOVED***

	netConn, err := do.dial(network, address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if do.useTLS ***REMOVED***
		var tlsConfig *tls.Config
		if do.tlsConfig == nil ***REMOVED***
			tlsConfig = &tls.Config***REMOVED***InsecureSkipVerify: do.skipVerify***REMOVED***
		***REMOVED*** else ***REMOVED***
			tlsConfig = cloneTLSConfig(do.tlsConfig)
		***REMOVED***
		if tlsConfig.ServerName == "" ***REMOVED***
			host, _, err := net.SplitHostPort(address)
			if err != nil ***REMOVED***
				netConn.Close()
				return nil, err
			***REMOVED***
			tlsConfig.ServerName = host
		***REMOVED***

		tlsConn := tls.Client(netConn, tlsConfig)
		if err := tlsConn.Handshake(); err != nil ***REMOVED***
			netConn.Close()
			return nil, err
		***REMOVED***
		netConn = tlsConn
	***REMOVED***

	c := &conn***REMOVED***
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  do.readTimeout,
		writeTimeout: do.writeTimeout,
	***REMOVED***

	if do.password != "" ***REMOVED***
		if _, err := c.Do("AUTH", do.password); err != nil ***REMOVED***
			netConn.Close()
			return nil, err
		***REMOVED***
	***REMOVED***

	if do.db != 0 ***REMOVED***
		if _, err := c.Do("SELECT", do.db); err != nil ***REMOVED***
			netConn.Close()
			return nil, err
		***REMOVED***
	***REMOVED***

	return c, nil
***REMOVED***

var pathDBRegexp = regexp.MustCompile(`/(\d*)\z`)

// DialURL connects to a Redis server at the given URL using the Redis
// URI scheme. URLs should follow the draft IANA specification for the
// scheme (https://www.iana.org/assignments/uri-schemes/prov/redis).
func DialURL(rawurl string, options ...DialOption) (Conn, error) ***REMOVED***
	u, err := url.Parse(rawurl)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if u.Scheme != "redis" && u.Scheme != "rediss" ***REMOVED***
		return nil, fmt.Errorf("invalid redis URL scheme: %s", u.Scheme)
	***REMOVED***

	// As per the IANA draft spec, the host defaults to localhost and
	// the port defaults to 6379.
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil ***REMOVED***
		// assume port is missing
		host = u.Host
		port = "6379"
	***REMOVED***
	if host == "" ***REMOVED***
		host = "localhost"
	***REMOVED***
	address := net.JoinHostPort(host, port)

	if u.User != nil ***REMOVED***
		password, isSet := u.User.Password()
		if isSet ***REMOVED***
			options = append(options, DialPassword(password))
		***REMOVED***
	***REMOVED***

	match := pathDBRegexp.FindStringSubmatch(u.Path)
	if len(match) == 2 ***REMOVED***
		db := 0
		if len(match[1]) > 0 ***REMOVED***
			db, err = strconv.Atoi(match[1])
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("invalid database: %s", u.Path[1:])
			***REMOVED***
		***REMOVED***
		if db != 0 ***REMOVED***
			options = append(options, DialDatabase(db))
		***REMOVED***
	***REMOVED*** else if u.Path != "" ***REMOVED***
		return nil, fmt.Errorf("invalid database: %s", u.Path[1:])
	***REMOVED***

	options = append(options, DialUseTLS(u.Scheme == "rediss"))

	return Dial("tcp", address, options...)
***REMOVED***

// NewConn returns a new Redigo connection for the given net connection.
func NewConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) Conn ***REMOVED***
	return &conn***REMOVED***
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	***REMOVED***
***REMOVED***

func (c *conn) Close() error ***REMOVED***
	c.mu.Lock()
	err := c.err
	if c.err == nil ***REMOVED***
		c.err = errors.New("redigo: closed")
		err = c.conn.Close()
	***REMOVED***
	c.mu.Unlock()
	return err
***REMOVED***

func (c *conn) fatal(err error) error ***REMOVED***
	c.mu.Lock()
	if c.err == nil ***REMOVED***
		c.err = err
		// Close connection to force errors on subsequent calls and to unblock
		// other reader or writer.
		c.conn.Close()
	***REMOVED***
	c.mu.Unlock()
	return err
***REMOVED***

func (c *conn) Err() error ***REMOVED***
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
***REMOVED***

func (c *conn) writeLen(prefix byte, n int) error ***REMOVED***
	c.lenScratch[len(c.lenScratch)-1] = '\n'
	c.lenScratch[len(c.lenScratch)-2] = '\r'
	i := len(c.lenScratch) - 3
	for ***REMOVED***
		c.lenScratch[i] = byte('0' + n%10)
		i -= 1
		n = n / 10
		if n == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	c.lenScratch[i] = prefix
	_, err := c.bw.Write(c.lenScratch[i:])
	return err
***REMOVED***

func (c *conn) writeString(s string) error ***REMOVED***
	c.writeLen('$', len(s))
	c.bw.WriteString(s)
	_, err := c.bw.WriteString("\r\n")
	return err
***REMOVED***

func (c *conn) writeBytes(p []byte) error ***REMOVED***
	c.writeLen('$', len(p))
	c.bw.Write(p)
	_, err := c.bw.WriteString("\r\n")
	return err
***REMOVED***

func (c *conn) writeInt64(n int64) error ***REMOVED***
	return c.writeBytes(strconv.AppendInt(c.numScratch[:0], n, 10))
***REMOVED***

func (c *conn) writeFloat64(n float64) error ***REMOVED***
	return c.writeBytes(strconv.AppendFloat(c.numScratch[:0], n, 'g', -1, 64))
***REMOVED***

func (c *conn) writeCommand(cmd string, args []interface***REMOVED******REMOVED***) error ***REMOVED***
	c.writeLen('*', 1+len(args))
	if err := c.writeString(cmd); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, arg := range args ***REMOVED***
		if err := c.writeArg(arg, true); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *conn) writeArg(arg interface***REMOVED******REMOVED***, argumentTypeOK bool) (err error) ***REMOVED***
	switch arg := arg.(type) ***REMOVED***
	case string:
		return c.writeString(arg)
	case []byte:
		return c.writeBytes(arg)
	case int:
		return c.writeInt64(int64(arg))
	case int64:
		return c.writeInt64(arg)
	case float64:
		return c.writeFloat64(arg)
	case bool:
		if arg ***REMOVED***
			return c.writeString("1")
		***REMOVED*** else ***REMOVED***
			return c.writeString("0")
		***REMOVED***
	case nil:
		return c.writeString("")
	case Argument:
		if argumentTypeOK ***REMOVED***
			return c.writeArg(arg.RedisArg(), false)
		***REMOVED***
		// See comment in default clause below.
		var buf bytes.Buffer
		fmt.Fprint(&buf, arg)
		return c.writeBytes(buf.Bytes())
	default:
		// This default clause is intended to handle builtin numeric types.
		// The function should return an error for other types, but this is not
		// done for compatibility with previous versions of the package.
		var buf bytes.Buffer
		fmt.Fprint(&buf, arg)
		return c.writeBytes(buf.Bytes())
	***REMOVED***
***REMOVED***

type protocolError string

func (pe protocolError) Error() string ***REMOVED***
	return fmt.Sprintf("redigo: %s (possible server error or unsupported concurrent read by application)", string(pe))
***REMOVED***

func (c *conn) readLine() ([]byte, error) ***REMOVED***
	p, err := c.br.ReadSlice('\n')
	if err == bufio.ErrBufferFull ***REMOVED***
		return nil, protocolError("long response line")
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	i := len(p) - 2
	if i < 0 || p[i] != '\r' ***REMOVED***
		return nil, protocolError("bad response line terminator")
	***REMOVED***
	return p[:i], nil
***REMOVED***

// parseLen parses bulk string and array lengths.
func parseLen(p []byte) (int, error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return -1, protocolError("malformed length")
	***REMOVED***

	if p[0] == '-' && len(p) == 2 && p[1] == '1' ***REMOVED***
		// handle $-1 and $-1 null replies.
		return -1, nil
	***REMOVED***

	var n int
	for _, b := range p ***REMOVED***
		n *= 10
		if b < '0' || b > '9' ***REMOVED***
			return -1, protocolError("illegal bytes in length")
		***REMOVED***
		n += int(b - '0')
	***REMOVED***

	return n, nil
***REMOVED***

// parseInt parses an integer reply.
func parseInt(p []byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return 0, protocolError("malformed integer")
	***REMOVED***

	var negate bool
	if p[0] == '-' ***REMOVED***
		negate = true
		p = p[1:]
		if len(p) == 0 ***REMOVED***
			return 0, protocolError("malformed integer")
		***REMOVED***
	***REMOVED***

	var n int64
	for _, b := range p ***REMOVED***
		n *= 10
		if b < '0' || b > '9' ***REMOVED***
			return 0, protocolError("illegal bytes in length")
		***REMOVED***
		n += int64(b - '0')
	***REMOVED***

	if negate ***REMOVED***
		n = -n
	***REMOVED***
	return n, nil
***REMOVED***

var (
	okReply   interface***REMOVED******REMOVED*** = "OK"
	pongReply interface***REMOVED******REMOVED*** = "PONG"
)

func (c *conn) readReply() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	line, err := c.readLine()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(line) == 0 ***REMOVED***
		return nil, protocolError("short response line")
	***REMOVED***
	switch line[0] ***REMOVED***
	case '+':
		switch ***REMOVED***
		case len(line) == 3 && line[1] == 'O' && line[2] == 'K':
			// Avoid allocation for frequent "+OK" response.
			return okReply, nil
		case len(line) == 5 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G':
			// Avoid allocation in PING command benchmarks :)
			return pongReply, nil
		default:
			return string(line[1:]), nil
		***REMOVED***
	case '-':
		return Error(string(line[1:])), nil
	case ':':
		return parseInt(line[1:])
	case '$':
		n, err := parseLen(line[1:])
		if n < 0 || err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		p := make([]byte, n)
		_, err = io.ReadFull(c.br, p)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if line, err := c.readLine(); err != nil ***REMOVED***
			return nil, err
		***REMOVED*** else if len(line) != 0 ***REMOVED***
			return nil, protocolError("bad bulk string format")
		***REMOVED***
		return p, nil
	case '*':
		n, err := parseLen(line[1:])
		if n < 0 || err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		r := make([]interface***REMOVED******REMOVED***, n)
		for i := range r ***REMOVED***
			r[i], err = c.readReply()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return r, nil
	***REMOVED***
	return nil, protocolError("unexpected response line")
***REMOVED***

func (c *conn) Send(cmd string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	c.mu.Lock()
	c.pending += 1
	c.mu.Unlock()
	if c.writeTimeout != 0 ***REMOVED***
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	***REMOVED***
	if err := c.writeCommand(cmd, args); err != nil ***REMOVED***
		return c.fatal(err)
	***REMOVED***
	return nil
***REMOVED***

func (c *conn) Flush() error ***REMOVED***
	if c.writeTimeout != 0 ***REMOVED***
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	***REMOVED***
	if err := c.bw.Flush(); err != nil ***REMOVED***
		return c.fatal(err)
	***REMOVED***
	return nil
***REMOVED***

func (c *conn) Receive() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return c.ReceiveWithTimeout(c.readTimeout)
***REMOVED***

func (c *conn) ReceiveWithTimeout(timeout time.Duration) (reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	var deadline time.Time
	if timeout != 0 ***REMOVED***
		deadline = time.Now().Add(timeout)
	***REMOVED***
	c.conn.SetReadDeadline(deadline)

	if reply, err = c.readReply(); err != nil ***REMOVED***
		return nil, c.fatal(err)
	***REMOVED***
	// When using pub/sub, the number of receives can be greater than the
	// number of sends. To enable normal use of the connection after
	// unsubscribing from all channels, we do not decrement pending to a
	// negative value.
	//
	// The pending field is decremented after the reply is read to handle the
	// case where Receive is called before Send.
	c.mu.Lock()
	if c.pending > 0 ***REMOVED***
		c.pending -= 1
	***REMOVED***
	c.mu.Unlock()
	if err, ok := reply.(Error); ok ***REMOVED***
		return nil, err
	***REMOVED***
	return
***REMOVED***

func (c *conn) Do(cmd string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return c.DoWithTimeout(c.readTimeout, cmd, args...)
***REMOVED***

func (c *conn) DoWithTimeout(readTimeout time.Duration, cmd string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	c.mu.Lock()
	pending := c.pending
	c.pending = 0
	c.mu.Unlock()

	if cmd == "" && pending == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	if c.writeTimeout != 0 ***REMOVED***
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	***REMOVED***

	if cmd != "" ***REMOVED***
		if err := c.writeCommand(cmd, args); err != nil ***REMOVED***
			return nil, c.fatal(err)
		***REMOVED***
	***REMOVED***

	if err := c.bw.Flush(); err != nil ***REMOVED***
		return nil, c.fatal(err)
	***REMOVED***

	var deadline time.Time
	if readTimeout != 0 ***REMOVED***
		deadline = time.Now().Add(readTimeout)
	***REMOVED***
	c.conn.SetReadDeadline(deadline)

	if cmd == "" ***REMOVED***
		reply := make([]interface***REMOVED******REMOVED***, pending)
		for i := range reply ***REMOVED***
			r, e := c.readReply()
			if e != nil ***REMOVED***
				return nil, c.fatal(e)
			***REMOVED***
			reply[i] = r
		***REMOVED***
		return reply, nil
	***REMOVED***

	var err error
	var reply interface***REMOVED******REMOVED***
	for i := 0; i <= pending; i++ ***REMOVED***
		var e error
		if reply, e = c.readReply(); e != nil ***REMOVED***
			return nil, c.fatal(e)
		***REMOVED***
		if e, ok := reply.(Error); ok && err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	return reply, err
***REMOVED***
