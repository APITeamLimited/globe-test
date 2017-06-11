/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package netext

import (
	"context"
	"net"
	"sync/atomic"

	"github.com/viki-org/dnscache"
)

type Dialer struct ***REMOVED***
	net.Dialer

	Resolver *dnscache.Resolver
***REMOVED***

func NewDialer(dialer net.Dialer) *Dialer ***REMOVED***
	return &Dialer***REMOVED***
		Dialer:   dialer,
		Resolver: dnscache.New(0),
	***REMOVED***
***REMOVED***

func (d Dialer) DialContext(ctx context.Context, proto, addr string) (net.Conn, error) ***REMOVED***
	host, port, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ip, err := d.Resolver.FetchOne(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	conn, err := d.Dialer.DialContext(ctx, proto, ip.String()+":"+port)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if v := ctx.Value(ctxKeyTracer); v != nil ***REMOVED***
		tracer := v.(*Tracer)
		return TrackConn(conn, &tracer.bytesRead, &tracer.bytesWritten), nil
	***REMOVED***
	return conn, err
***REMOVED***

type Conn struct ***REMOVED***
	net.Conn

	BytesRead, BytesWritten *int64
***REMOVED***

func (c *Conn) Read(b []byte) (int, error) ***REMOVED***
	n, err := c.Conn.Read(b)
	atomic.AddInt64(c.BytesRead, int64(n))
	return n, err
***REMOVED***

func (c *Conn) Write(b []byte) (int, error) ***REMOVED***
	n, err := c.Conn.Write(b)
	atomic.AddInt64(c.BytesWritten, int64(n))
	return n, err
***REMOVED***

func TrackConn(conn net.Conn, bytesRead, bytesWritten *int64) *Conn ***REMOVED***
	return &Conn***REMOVED***conn, bytesRead, bytesWritten***REMOVED***
***REMOVED***
