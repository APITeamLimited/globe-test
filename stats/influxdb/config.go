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
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ConfigFields struct ***REMOVED***
	// Connection.
	Addr        string `json:"addr"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Insecure    bool   `json:"insecure,omitempty"`
	PayloadSize int    `json:"payload_size,omitempty"`

	// Samples.
	Database    string `json:"database"`
	Precision   string `json:"precision,omitempty"`
	Retention   string `json:"retention,omitempty"`
	Consistency string `json:"consistency,omitempty"`
***REMOVED***

type Config ConfigFields

func (c *Config) UnmarshalText(text []byte) error ***REMOVED***
	u, err := url.Parse(string(text))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if u.Host != "" ***REMOVED***
		c.Addr = u.Scheme + "://" + u.Host
	***REMOVED***
	if db := strings.TrimPrefix(u.Path, "/"); db != "" ***REMOVED***
		c.Database = db
	***REMOVED***
	if u.User != nil ***REMOVED***
		c.Username = u.User.Username()
		c.Password, _ = u.User.Password()
	***REMOVED***
	for k, vs := range u.Query() ***REMOVED***
		switch k ***REMOVED***
		case "insecure":
			switch vs[0] ***REMOVED***
			case "":
			case "false":
				c.Insecure = false
			case "true":
				c.Insecure = true
			default:
				return errors.Errorf("insecure must be true or false, not %s", vs[0])
			***REMOVED***
		case "payload_size":
			c.PayloadSize, err = strconv.Atoi(vs[0])
		case "precision":
			c.Precision = vs[0]
		case "retention":
			c.Retention = vs[0]
		case "consistency":
			c.Consistency = vs[0]
		default:
			return errors.Errorf("unknown query parameter: %s", k)
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (c *Config) UnmarshalJSON(data []byte) error ***REMOVED***
	fields := ConfigFields(*c)
	if err := json.Unmarshal(data, &fields); err != nil ***REMOVED***
		return err
	***REMOVED***
	*c = Config(fields)
	return nil
***REMOVED***

func (c Config) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(ConfigFields(c))
***REMOVED***
