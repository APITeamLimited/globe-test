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

	"github.com/kubernetes/helm/pkg/strvals"
	"github.com/loadimpact/k6/lib/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	null "gopkg.in/guregu/null.v3"
)

type Config struct ***REMOVED***
	// Connection.
	Addr        null.String `json:"addr" envconfig:"INFLUXDB_ADDR"`
	Username    null.String `json:"username,omitempty" envconfig:"INFLUXDB_USERNAME"`
	Password    null.String `json:"password,omitempty" envconfig:"INFLUXDB_PASSWORD"`
	Insecure    null.Bool   `json:"insecure,omitempty" envconfig:"INFLUXDB_INSECURE"`
	PayloadSize null.Int    `json:"payloadSize,omitempty" envconfig:"INFLUXDB_PAYLOAD_SIZE"`

	// Samples.
	DB           null.String `json:"db" envconfig:"INFLUXDB_DB"`
	Precision    null.String `json:"precision,omitempty" envconfig:"INFLUXDB_PRECISION"`
	Retention    null.String `json:"retention,omitempty" envconfig:"INFLUXDB_RETENTION"`
	Consistency  null.String `json:"consistency,omitempty" envconfig:"INFLUXDB_CONSISTENCY"`
	TagsAsFields []string    `json:"tagsAsFields,omitempty" envconfig:"INFLUXDB_TAGS_AS_FIELDS"`
***REMOVED***

func NewConfig() *Config ***REMOVED***
	c := &Config***REMOVED***
		Addr:         null.NewString("http://localhost:8086", false),
		DB:           null.NewString("k6", false),
		TagsAsFields: []string***REMOVED***"vu", "iter", "url"***REMOVED***,
	***REMOVED***
	return c
***REMOVED***

func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.Addr.Valid ***REMOVED***
		c.Addr = cfg.Addr
	***REMOVED***
	if cfg.Username.Valid ***REMOVED***
		c.Username = cfg.Username
	***REMOVED***
	if cfg.Password.Valid ***REMOVED***
		c.Password = cfg.Password
	***REMOVED***
	if cfg.Insecure.Valid ***REMOVED***
		c.Insecure = cfg.Insecure
	***REMOVED***
	if cfg.PayloadSize.Valid && cfg.PayloadSize.Int64 > 0 ***REMOVED***
		c.PayloadSize = cfg.PayloadSize
	***REMOVED***
	if cfg.DB.Valid ***REMOVED***
		c.DB = cfg.DB
	***REMOVED***
	if cfg.Precision.Valid ***REMOVED***
		c.Precision = cfg.Precision
	***REMOVED***
	if cfg.Retention.Valid ***REMOVED***
		c.Retention = cfg.Retention
	***REMOVED***
	if cfg.Consistency.Valid ***REMOVED***
		c.Consistency = cfg.Consistency
	***REMOVED***
	if len(cfg.TagsAsFields) > 0 ***REMOVED***
		c.TagsAsFields = cfg.TagsAsFields
	***REMOVED***
	return c
***REMOVED***

// ParseArg parses an argument string into a Config
func ParseArg(arg string) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***
	params, err := strvals.Parse(arg)

	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	c, err = ParseMap(params)
	return c, err
***REMOVED***

// ParseMap parses a map[string]interface***REMOVED******REMOVED*** into a Config
func ParseMap(m map[string]interface***REMOVED******REMOVED***) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***
	if v, ok := m["tagsAsFields"].(string); ok ***REMOVED***
		m["tagsAsFields"] = []string***REMOVED***v***REMOVED***
	***REMOVED***
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig***REMOVED***
		DecodeHook: types.NullDecoder,
		Result:     &c,
	***REMOVED***)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	err = dec.Decode(m)
	return c, err
***REMOVED***

func ParseURL(text string) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***
	u, err := url.Parse(text)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***
	if u.Host != "" ***REMOVED***
		c.Addr = null.StringFrom(u.Scheme + "://" + u.Host)
	***REMOVED***
	if db := strings.TrimPrefix(u.Path, "/"); db != "" ***REMOVED***
		c.DB = null.StringFrom(db)
	***REMOVED***
	if u.User != nil ***REMOVED***
		c.Username = null.StringFrom(u.User.Username())
		pass, _ := u.User.Password()
		c.Password = null.StringFrom(pass)
	***REMOVED***
	for k, vs := range u.Query() ***REMOVED***
		switch k ***REMOVED***
		case "insecure":
			switch vs[0] ***REMOVED***
			case "":
			case "false":
				c.Insecure = null.BoolFrom(false)
			case "true":
				c.Insecure = null.BoolFrom(true)
			default:
				return c, errors.Errorf("insecure must be true or false, not %s", vs[0])
			***REMOVED***
		case "payload_size":
			var size int
			size, err = strconv.Atoi(vs[0])
			c.PayloadSize = null.IntFrom(int64(size))
		case "precision":
			c.Precision = null.StringFrom(vs[0])
		case "retention":
			c.Retention = null.StringFrom(vs[0])
		case "consistency":
			c.Consistency = null.StringFrom(vs[0])
		case "tagsAsFields":
			c.TagsAsFields = vs
		default:
			return c, errors.Errorf("unknown query parameter: %s", k)
		***REMOVED***
	***REMOVED***
	return c, err
***REMOVED***
