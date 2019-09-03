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

	client "github.com/influxdata/influxdb1-client/v2"
	null "gopkg.in/guregu/null.v3"
)

func MakeClient(conf Config) (client.Client, error) ***REMOVED***
	if strings.HasPrefix(conf.Addr.String, "udp://") ***REMOVED***
		return client.NewUDPClient(client.UDPConfig***REMOVED***
			Addr:        strings.TrimPrefix(conf.Addr.String, "udp://"),
			PayloadSize: int(conf.PayloadSize.Int64),
		***REMOVED***)
	***REMOVED***
	if conf.Addr.String == "" ***REMOVED***
		conf.Addr = null.StringFrom("http://localhost:8086")
	***REMOVED***
	return client.NewHTTPClient(client.HTTPConfig***REMOVED***
		Addr:               conf.Addr.String,
		Username:           conf.Username.String,
		Password:           conf.Password.String,
		UserAgent:          "k6",
		InsecureSkipVerify: conf.Insecure.Bool,
	***REMOVED***)
***REMOVED***

func MakeBatchConfig(conf Config) client.BatchPointsConfig ***REMOVED***
	if !conf.DB.Valid || conf.DB.String == "" ***REMOVED***
		conf.DB = null.StringFrom("k6")
	***REMOVED***
	return client.BatchPointsConfig***REMOVED***
		Precision:        conf.Precision.String,
		Database:         conf.DB.String,
		RetentionPolicy:  conf.Retention.String,
		WriteConsistency: conf.Consistency.String,
	***REMOVED***
***REMOVED***
