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

// Subscription represents a subscribe or unsubscribe notification.
type Subscription struct ***REMOVED***
	// Kind is "subscribe", "unsubscribe", "psubscribe" or "punsubscribe"
	Kind string

	// The channel that was changed.
	Channel string

	// The current number of subscriptions for connection.
	Count int
***REMOVED***

// Message represents a message notification.
type Message struct ***REMOVED***
	// The originating channel.
	Channel string

	// The message data.
	Data []byte
***REMOVED***

// PMessage represents a pmessage notification.
type PMessage struct ***REMOVED***
	// The matched pattern.
	Pattern string

	// The originating channel.
	Channel string

	// The message data.
	Data []byte
***REMOVED***

// Pong represents a pubsub pong notification.
type Pong struct ***REMOVED***
	Data string
***REMOVED***

// PubSubConn wraps a Conn with convenience methods for subscribers.
type PubSubConn struct ***REMOVED***
	Conn Conn
***REMOVED***

// Close closes the connection.
func (c PubSubConn) Close() error ***REMOVED***
	return c.Conn.Close()
***REMOVED***

// Subscribe subscribes the connection to the specified channels.
func (c PubSubConn) Subscribe(channel ...interface***REMOVED******REMOVED***) error ***REMOVED***
	c.Conn.Send("SUBSCRIBE", channel...)
	return c.Conn.Flush()
***REMOVED***

// PSubscribe subscribes the connection to the given patterns.
func (c PubSubConn) PSubscribe(channel ...interface***REMOVED******REMOVED***) error ***REMOVED***
	c.Conn.Send("PSUBSCRIBE", channel...)
	return c.Conn.Flush()
***REMOVED***

// Unsubscribe unsubscribes the connection from the given channels, or from all
// of them if none is given.
func (c PubSubConn) Unsubscribe(channel ...interface***REMOVED******REMOVED***) error ***REMOVED***
	c.Conn.Send("UNSUBSCRIBE", channel...)
	return c.Conn.Flush()
***REMOVED***

// PUnsubscribe unsubscribes the connection from the given patterns, or from all
// of them if none is given.
func (c PubSubConn) PUnsubscribe(channel ...interface***REMOVED******REMOVED***) error ***REMOVED***
	c.Conn.Send("PUNSUBSCRIBE", channel...)
	return c.Conn.Flush()
***REMOVED***

// Ping sends a PING to the server with the specified data.
//
// The connection must be subscribed to at least one channel or pattern when
// calling this method.
func (c PubSubConn) Ping(data string) error ***REMOVED***
	c.Conn.Send("PING", data)
	return c.Conn.Flush()
***REMOVED***

// Receive returns a pushed message as a Subscription, Message, PMessage, Pong
// or error. The return value is intended to be used directly in a type switch
// as illustrated in the PubSubConn example.
func (c PubSubConn) Receive() interface***REMOVED******REMOVED*** ***REMOVED***
	return c.receiveInternal(c.Conn.Receive())
***REMOVED***

// ReceiveWithTimeout is like Receive, but it allows the application to
// override the connection's default timeout.
func (c PubSubConn) ReceiveWithTimeout(timeout time.Duration) interface***REMOVED******REMOVED*** ***REMOVED***
	return c.receiveInternal(ReceiveWithTimeout(c.Conn, timeout))
***REMOVED***

func (c PubSubConn) receiveInternal(replyArg interface***REMOVED******REMOVED***, errArg error) interface***REMOVED******REMOVED*** ***REMOVED***
	reply, err := Values(replyArg, errArg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var kind string
	reply, err = Scan(reply, &kind)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch kind ***REMOVED***
	case "message":
		var m Message
		if _, err := Scan(reply, &m.Channel, &m.Data); err != nil ***REMOVED***
			return err
		***REMOVED***
		return m
	case "pmessage":
		var pm PMessage
		if _, err := Scan(reply, &pm.Pattern, &pm.Channel, &pm.Data); err != nil ***REMOVED***
			return err
		***REMOVED***
		return pm
	case "subscribe", "psubscribe", "unsubscribe", "punsubscribe":
		s := Subscription***REMOVED***Kind: kind***REMOVED***
		if _, err := Scan(reply, &s.Channel, &s.Count); err != nil ***REMOVED***
			return err
		***REMOVED***
		return s
	case "pong":
		var p Pong
		if _, err := Scan(reply, &p.Data); err != nil ***REMOVED***
			return err
		***REMOVED***
		return p
	***REMOVED***
	return errors.New("redigo: unknown pubsub notification")
***REMOVED***
