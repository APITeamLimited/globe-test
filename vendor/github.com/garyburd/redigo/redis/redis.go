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
	"errors"
	"time"
)

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string ***REMOVED*** return string(err) ***REMOVED***

// Conn represents a connection to a Redis server.
type Conn interface ***REMOVED***
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value when the connection is not usable.
	Err() error

	// Do sends a command to the server and returns the received reply.
	Do(commandName string, args ...interface***REMOVED******REMOVED***) (reply interface***REMOVED******REMOVED***, err error)

	// Send writes the command to the client's output buffer.
	Send(commandName string, args ...interface***REMOVED******REMOVED***) error

	// Flush flushes the output buffer to the Redis server.
	Flush() error

	// Receive receives a single reply from the Redis server
	Receive() (reply interface***REMOVED******REMOVED***, err error)
***REMOVED***

// Argument is the interface implemented by an object which wants to control how
// the object is converted to Redis bulk strings.
type Argument interface ***REMOVED***
	// RedisArg returns a value to be encoded as a bulk string per the
	// conversions listed in the section 'Executing Commands'.
	// Implementations should typically return a []byte or string.
	RedisArg() interface***REMOVED******REMOVED***
***REMOVED***

// Scanner is implemented by an object which wants to control its value is
// interpreted when read from Redis.
type Scanner interface ***REMOVED***
	// RedisScan assigns a value from a Redis value. The argument src is one of
	// the reply types listed in the section `Executing Commands`.
	//
	// An error should be returned if the value cannot be stored without
	// loss of information.
	RedisScan(src interface***REMOVED******REMOVED***) error
***REMOVED***

// ConnWithTimeout is an optional interface that allows the caller to override
// a connection's default read timeout. This interface is useful for executing
// the BLPOP, BRPOP, BRPOPLPUSH, XREAD and other commands that block at the
// server.
//
// A connection's default read timeout is set with the DialReadTimeout dial
// option. Applications should rely on the default timeout for commands that do
// not block at the server.
//
// All of the Conn implementations in this package satisfy the ConnWithTimeout
// interface.
//
// Use the DoWithTimeout and ReceiveWithTimeout helper functions to simplify
// use of this interface.
type ConnWithTimeout interface ***REMOVED***
	Conn

	// Do sends a command to the server and returns the received reply.
	// The timeout overrides the read timeout set when dialing the
	// connection.
	DoWithTimeout(timeout time.Duration, commandName string, args ...interface***REMOVED******REMOVED***) (reply interface***REMOVED******REMOVED***, err error)

	// Receive receives a single reply from the Redis server. The timeout
	// overrides the read timeout set when dialing the connection.
	ReceiveWithTimeout(timeout time.Duration) (reply interface***REMOVED******REMOVED***, err error)
***REMOVED***

var errTimeoutNotSupported = errors.New("redis: connection does not support ConnWithTimeout")

// DoWithTimeout executes a Redis command with the specified read timeout. If
// the connection does not satisfy the ConnWithTimeout interface, then an error
// is returned.
func DoWithTimeout(c Conn, timeout time.Duration, cmd string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	cwt, ok := c.(ConnWithTimeout)
	if !ok ***REMOVED***
		return nil, errTimeoutNotSupported
	***REMOVED***
	return cwt.DoWithTimeout(timeout, cmd, args...)
***REMOVED***

// ReceiveWithTimeout receives a reply with the specified read timeout. If the
// connection does not satisfy the ConnWithTimeout interface, then an error is
// returned.
func ReceiveWithTimeout(c Conn, timeout time.Duration) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	cwt, ok := c.(ConnWithTimeout)
	if !ok ***REMOVED***
		return nil, errTimeoutNotSupported
	***REMOVED***
	return cwt.ReceiveWithTimeout(timeout)
***REMOVED***
