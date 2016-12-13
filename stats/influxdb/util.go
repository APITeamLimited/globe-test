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

package influxdb

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrNoDatabase = errors.New("influxdb output: no database specified")
)

func parseURL(u *url.URL) (client.Client, client.BatchPointsConfig, error) ***REMOVED***
	batchConf, err := makeBatchConfigFromURL(u)
	if err != nil ***REMOVED***
		return nil, client.BatchPointsConfig***REMOVED******REMOVED***, err
	***REMOVED***

	if u.Scheme == "udp" ***REMOVED***
		conf, err := makeUDPConfigFromURL(u)
		if err != nil ***REMOVED***
			return nil, batchConf, err
		***REMOVED***
		c, err := client.NewUDPClient(conf)
		if err != nil ***REMOVED***
			return nil, batchConf, err
		***REMOVED***
		return c, batchConf, nil
	***REMOVED***

	conf, err := makeHTTPConfigFromURL(u)
	if err != nil ***REMOVED***
		return nil, batchConf, err
	***REMOVED***
	c, err := client.NewHTTPClient(conf)
	if err != nil ***REMOVED***
		return nil, batchConf, err
	***REMOVED***
	return c, batchConf, nil
***REMOVED***

func makeUDPConfigFromURL(u *url.URL) (client.UDPConfig, error) ***REMOVED***
	payloadSize := 0
	payloadSizeS := u.Query().Get("payload_size")
	if payloadSizeS != "" ***REMOVED***
		s, err := strconv.ParseInt(payloadSizeS, 10, 32)
		if err != nil ***REMOVED***
			return client.UDPConfig***REMOVED******REMOVED***, err
		***REMOVED***
		payloadSize = int(s)
	***REMOVED***

	return client.UDPConfig***REMOVED***
		Addr:        u.Host,
		PayloadSize: payloadSize,
	***REMOVED***, nil
***REMOVED***

func makeHTTPConfigFromURL(u *url.URL) (client.HTTPConfig, error) ***REMOVED***
	q := u.Query()

	username := ""
	password := ""
	if u.User != nil ***REMOVED***
		username = u.User.Username()
		password, _ = u.User.Password()
	***REMOVED***

	timeout := 0 * time.Second
	if ts := q.Get("timeout"); ts != "" ***REMOVED***
		t, err := time.ParseDuration(ts)
		if err != nil ***REMOVED***
			return client.HTTPConfig***REMOVED******REMOVED***, err
		***REMOVED***
		timeout = t
	***REMOVED***
	insecureSkipVerify := q.Get("insecure_skip_verify") != ""

	return client.HTTPConfig***REMOVED***
		Addr:               u.Scheme + "://" + u.Host,
		Username:           username,
		Password:           password,
		Timeout:            timeout,
		InsecureSkipVerify: insecureSkipVerify,
	***REMOVED***, nil
***REMOVED***

func makeBatchConfigFromURL(u *url.URL) (client.BatchPointsConfig, error) ***REMOVED***
	if u.Path == "" || u.Path == "/" ***REMOVED***
		return client.BatchPointsConfig***REMOVED******REMOVED***, ErrNoDatabase
	***REMOVED***

	q := u.Query()
	return client.BatchPointsConfig***REMOVED***
		Database:         u.Path[1:], // strip leading "/"
		Precision:        q.Get("precision"),
		RetentionPolicy:  q.Get("retention_policy"),
		WriteConsistency: q.Get("write_consistency"),
	***REMOVED***, nil
***REMOVED***
