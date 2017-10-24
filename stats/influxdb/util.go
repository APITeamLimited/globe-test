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
	"strings"

	client "github.com/influxdata/influxdb/client/v2"
)

func MakeClient(conf Config) (client.Client, error) ***REMOVED***
	if strings.HasPrefix(conf.Addr, "udp://") ***REMOVED***
		return client.NewUDPClient(client.UDPConfig***REMOVED***
			Addr:        strings.TrimPrefix(conf.Addr, "udp://"),
			PayloadSize: conf.PayloadSize,
		***REMOVED***)
	***REMOVED***
	if conf.Addr == "" ***REMOVED***
		conf.Addr = "http://localhost:8086"
	***REMOVED***
	return client.NewHTTPClient(client.HTTPConfig***REMOVED***
		Addr:               conf.Addr,
		Username:           conf.Username,
		Password:           conf.Password,
		UserAgent:          "k6",
		InsecureSkipVerify: conf.Insecure,
	***REMOVED***)
***REMOVED***

func MakeBatchConfig(conf Config) client.BatchPointsConfig ***REMOVED***
	if conf.DB == "" ***REMOVED***
		conf.DB = "k6"
	***REMOVED***
	return client.BatchPointsConfig***REMOVED***
		Precision:        conf.Precision,
		Database:         conf.DB,
		RetentionPolicy:  conf.Retention,
		WriteConsistency: conf.Consistency,
	***REMOVED***
***REMOVED***
