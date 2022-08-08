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
	"bytes"
	"fmt"
	"log"
	"time"
)

var (
	_ ConnWithTimeout = (*loggingConn)(nil)
)

// NewLoggingConn returns a logging wrapper around a connection.
func NewLoggingConn(conn Conn, logger *log.Logger, prefix string) Conn ***REMOVED***
	if prefix != "" ***REMOVED***
		prefix = prefix + "."
	***REMOVED***
	return &loggingConn***REMOVED***conn, logger, prefix***REMOVED***
***REMOVED***

type loggingConn struct ***REMOVED***
	Conn
	logger *log.Logger
	prefix string
***REMOVED***

func (c *loggingConn) Close() error ***REMOVED***
	err := c.Conn.Close()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%sClose() -> (%v)", c.prefix, err)
	c.logger.Output(2, buf.String())
	return err
***REMOVED***

func (c *loggingConn) printValue(buf *bytes.Buffer, v interface***REMOVED******REMOVED***) ***REMOVED***
	const chop = 32
	switch v := v.(type) ***REMOVED***
	case []byte:
		if len(v) > chop ***REMOVED***
			fmt.Fprintf(buf, "%q...", v[:chop])
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(buf, "%q", v)
		***REMOVED***
	case string:
		if len(v) > chop ***REMOVED***
			fmt.Fprintf(buf, "%q...", v[:chop])
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(buf, "%q", v)
		***REMOVED***
	case []interface***REMOVED******REMOVED***:
		if len(v) == 0 ***REMOVED***
			buf.WriteString("[]")
		***REMOVED*** else ***REMOVED***
			sep := "["
			fin := "]"
			if len(v) > chop ***REMOVED***
				v = v[:chop]
				fin = "...]"
			***REMOVED***
			for _, vv := range v ***REMOVED***
				buf.WriteString(sep)
				c.printValue(buf, vv)
				sep = ", "
			***REMOVED***
			buf.WriteString(fin)
		***REMOVED***
	default:
		fmt.Fprint(buf, v)
	***REMOVED***
***REMOVED***

func (c *loggingConn) print(method, commandName string, args []interface***REMOVED******REMOVED***, reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s%s(", c.prefix, method)
	if method != "Receive" ***REMOVED***
		buf.WriteString(commandName)
		for _, arg := range args ***REMOVED***
			buf.WriteString(", ")
			c.printValue(&buf, arg)
		***REMOVED***
	***REMOVED***
	buf.WriteString(") -> (")
	if method != "Send" ***REMOVED***
		c.printValue(&buf, reply)
		buf.WriteString(", ")
	***REMOVED***
	fmt.Fprintf(&buf, "%v)", err)
	c.logger.Output(3, buf.String())
***REMOVED***

func (c *loggingConn) Do(commandName string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	reply, err := c.Conn.Do(commandName, args...)
	c.print("Do", commandName, args, reply, err)
	return reply, err
***REMOVED***

func (c *loggingConn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	reply, err := DoWithTimeout(c.Conn, timeout, commandName, args...)
	c.print("DoWithTimeout", commandName, args, reply, err)
	return reply, err
***REMOVED***

func (c *loggingConn) Send(commandName string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	err := c.Conn.Send(commandName, args...)
	c.print("Send", commandName, args, nil, err)
	return err
***REMOVED***

func (c *loggingConn) Receive() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	reply, err := c.Conn.Receive()
	c.print("Receive", "", nil, reply, err)
	return reply, err
***REMOVED***

func (c *loggingConn) ReceiveWithTimeout(timeout time.Duration) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	reply, err := ReceiveWithTimeout(c.Conn, timeout)
	c.print("ReceiveWithTimeout", "", nil, reply, err)
	return reply, err
***REMOVED***
