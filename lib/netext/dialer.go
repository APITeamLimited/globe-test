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
	"strings"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/viki-org/dnscache"
)

type Dialer struct ***REMOVED***
	net.Dialer

	Resolver  *dnscache.Resolver
	Blacklist []*net.IPNet
	Hosts     map[string]net.IP

	BytesRead    int64
	BytesWritten int64
***REMOVED***

func NewDialer(dialer net.Dialer) *Dialer ***REMOVED***
	return &Dialer***REMOVED***
		Dialer:   dialer,
		Resolver: dnscache.New(0),
	***REMOVED***
***REMOVED***

func (d *Dialer) DialContext(ctx context.Context, proto, addr string) (net.Conn, error) ***REMOVED***
	delimiter := strings.LastIndex(addr, ":")
	host := addr[:delimiter]

	// lookup for domain defined in Hosts option before trying to resolve DNS.
	ip, ok := d.Hosts[host]
	if !ok ***REMOVED***
		var err error
		ip, err = d.Resolver.FetchOne(host)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	for _, net := range d.Blacklist ***REMOVED***
		if net.Contains(ip) ***REMOVED***
			return nil, errors.Errorf("IP (%s) is in a blacklisted range (%s)", ip, net)
		***REMOVED***
	***REMOVED***
	ipStr := ip.String()
	if strings.ContainsRune(ipStr, ':') ***REMOVED***
		ipStr = "[" + ipStr + "]"
	***REMOVED***
	conn, err := d.Dialer.DialContext(ctx, proto, ipStr+":"+addr[delimiter+1:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	conn = &Conn***REMOVED***conn, &d.BytesRead, &d.BytesWritten***REMOVED***
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
