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
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/types"
)

type Config struct ***REMOVED***
	// Connection.
	Addr             null.String        `json:"addr" envconfig:"K6_INFLUXDB_ADDR"`
	Username         null.String        `json:"username,omitempty" envconfig:"K6_INFLUXDB_USERNAME"`
	Password         null.String        `json:"password,omitempty" envconfig:"K6_INFLUXDB_PASSWORD"`
	Insecure         null.Bool          `json:"insecure,omitempty" envconfig:"K6_INFLUXDB_INSECURE"`
	PayloadSize      null.Int           `json:"payloadSize,omitempty" envconfig:"K6_INFLUXDB_PAYLOAD_SIZE"`
	PushInterval     types.NullDuration `json:"pushInterval,omitempty" envconfig:"K6_INFLUXDB_PUSH_INTERVAL"`
	ConcurrentWrites null.Int           `json:"concurrentWrites,omitempty" envconfig:"K6_INFLUXDB_CONCURRENT_WRITES"`

	// Samples.
	DB           null.String `json:"db" envconfig:"K6_INFLUXDB_DB"`
	Precision    null.String `json:"precision,omitempty" envconfig:"K6_INFLUXDB_PRECISION"`
	Retention    null.String `json:"retention,omitempty" envconfig:"K6_INFLUXDB_RETENTION"`
	Consistency  null.String `json:"consistency,omitempty" envconfig:"K6_INFLUXDB_CONSISTENCY"`
	TagsAsFields []string    `json:"tagsAsFields,omitempty" envconfig:"K6_INFLUXDB_TAGS_AS_FIELDS"`
***REMOVED***

// NewConfig creates a new InfluxDB output config with some default values.
func NewConfig() Config ***REMOVED***
	c := Config***REMOVED***
		Addr:         null.NewString("http://localhost:8086", false),
		DB:           null.NewString("k6", false),
		TagsAsFields: []string***REMOVED***"vu", "iter", "url"***REMOVED***,
		PushInterval: types.NewNullDuration(time.Second, false),

		// The minimum value of pow(2, N) for handling a stressful situation
		// with the default push interval set to 1s.
		// Concurrency is not expected for the normal use-case,
		// the response time should be lower than the push interval set value.
		// In case of spikes, the response time could go around 2s,
		// higher values will highlight a not sustainable situation
		// and the user should adjust the executed script
		// or the configuration based on the environment and rate expected.
		ConcurrentWrites: null.NewInt(4, false),
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
	if cfg.PushInterval.Valid ***REMOVED***
		c.PushInterval = cfg.PushInterval
	***REMOVED***

	if cfg.ConcurrentWrites.Valid ***REMOVED***
		c.ConcurrentWrites = cfg.ConcurrentWrites
	***REMOVED***
	return c
***REMOVED***

// ParseJSON parses the supplied JSON into a Config.
func ParseJSON(data json.RawMessage) (Config, error) ***REMOVED***
	conf := Config***REMOVED******REMOVED***
	err := json.Unmarshal(data, &conf)
	return conf, err
***REMOVED***

// ParseURL parses the supplied URL into a Config.
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
				return c, fmt.Errorf("insecure must be true or false, not %s", vs[0])
			***REMOVED***
		case "payload_size":
			var size int
			size, err = strconv.Atoi(vs[0])
			if err != nil ***REMOVED***
				return c, err
			***REMOVED***
			c.PayloadSize = null.IntFrom(int64(size))
		case "precision":
			c.Precision = null.StringFrom(vs[0])
		case "retention":
			c.Retention = null.StringFrom(vs[0])
		case "consistency":
			c.Consistency = null.StringFrom(vs[0])

		case "pushInterval":
			err = c.PushInterval.UnmarshalText([]byte(vs[0]))
			if err != nil ***REMOVED***
				return c, err
			***REMOVED***
		case "concurrentWrites":
			var writes int
			writes, err = strconv.Atoi(vs[0])
			if err != nil ***REMOVED***
				return c, err
			***REMOVED***
			c.ConcurrentWrites = null.IntFrom(int64(writes))
		case "tagsAsFields":
			c.TagsAsFields = vs
		default:
			return c, fmt.Errorf("unknown query parameter: %s", k)
		***REMOVED***
	***REMOVED***
	return c, err
***REMOVED***

// GetConsolidatedConfig combines ***REMOVED***default config values + JSON config +
// environment vars + URL config values***REMOVED***, and returns the final result.
func GetConsolidatedConfig(jsonRawConf json.RawMessage, env map[string]string, url string) (Config, error) ***REMOVED***
	result := NewConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf, err := ParseJSON(jsonRawConf)
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(jsonConf)
	***REMOVED***

	envConfig := Config***REMOVED******REMOVED***
	if err := envconfig.Process("", &envConfig); err != nil ***REMOVED***
		// TODO: get rid of envconfig and actually use the env parameter...
		return result, err
	***REMOVED***
	result = result.Apply(envConfig)

	if url != "" ***REMOVED***
		urlConf, err := ParseURL(url)
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(urlConf)
	***REMOVED***

	return result, nil
***REMOVED***
