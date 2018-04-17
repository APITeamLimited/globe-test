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
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Config struct ***REMOVED***
	// Connection.
	Addr        string `json:"addr" envconfig:"INFLUXDB_ADDR"`
	Username    string `json:"username,omitempty" envconfig:"INFLUXDB_USERNAME"`
	Password    string `json:"password,omitempty" envconfig:"INFLUXDB_PASSWORD"`
	Insecure    bool   `json:"insecure,omitempty" envconfig:"INFLUXDB_INSECURE"`
	PayloadSize int    `json:"payloadSize,omitempty" envconfig:"INFLUXDB_PAYLOAD_SIZE"`

	// Samples.
	DB           string   `json:"db" envconfig:"INFLUXDB_DB"`
	Precision    string   `json:"precision,omitempty" envconfig:"INFLUXDB_PRECISION"`
	Retention    string   `json:"retention,omitempty" envconfig:"INFLUXDB_RETENTION"`
	Consistency  string   `json:"consistency,omitempty" envconfig:"INFLUXDB_CONSISTENCY"`
	TagsAsFields []string `json:"tagsAsFields,omitempty" envconfig:"INFLUXDB_TAGS_AS_FIELDS"`
***REMOVED***

func NewConfig() *Config ***REMOVED***
	c := &Config***REMOVED***TagsAsFields: []string***REMOVED***"vu", "iter", "url"***REMOVED******REMOVED***
	return c
***REMOVED***

func (c Config) Apply(cfg Config) Config ***REMOVED***
	//TODO: fix this, use nullable values like all other configs...
	if cfg.Addr != "" ***REMOVED***
		c.Addr = cfg.Addr
	***REMOVED***
	if cfg.Username != "" ***REMOVED***
		c.Username = cfg.Username
	***REMOVED***
	if cfg.Password != "" ***REMOVED***
		c.Password = cfg.Password
	***REMOVED***
	if cfg.Insecure ***REMOVED***
		c.Insecure = cfg.Insecure
	***REMOVED***
	if cfg.PayloadSize > 0 ***REMOVED***
		c.PayloadSize = cfg.PayloadSize
	***REMOVED***
	if cfg.DB != "" ***REMOVED***
		c.DB = cfg.DB
	***REMOVED***
	if cfg.Precision != "" ***REMOVED***
		c.Precision = cfg.Precision
	***REMOVED***
	if cfg.Retention != "" ***REMOVED***
		c.Retention = cfg.Retention
	***REMOVED***
	if cfg.Consistency != "" ***REMOVED***
		c.Consistency = cfg.Consistency
	***REMOVED***
	if len(cfg.TagsAsFields) > 0 ***REMOVED***
		c.TagsAsFields = cfg.TagsAsFields
	***REMOVED***
	return c
***REMOVED***

func ParseURL(text string) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***
	u, err := url.Parse(text)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***
	if u.Host != "" ***REMOVED***
		c.Addr = u.Scheme + "://" + u.Host
	***REMOVED***
	if db := strings.TrimPrefix(u.Path, "/"); db != "" ***REMOVED***
		c.DB = db
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
				return c, errors.Errorf("insecure must be true or false, not %s", vs[0])
			***REMOVED***
		case "payload_size":
			c.PayloadSize, err = strconv.Atoi(vs[0])
		case "precision":
			c.Precision = vs[0]
		case "retention":
			c.Retention = vs[0]
		case "consistency":
			c.Consistency = vs[0]
		case "tagsAsFields":
			c.TagsAsFields = vs
		default:
			return c, errors.Errorf("unknown query parameter: %s", k)
		***REMOVED***
	***REMOVED***
	return c, err
***REMOVED***
