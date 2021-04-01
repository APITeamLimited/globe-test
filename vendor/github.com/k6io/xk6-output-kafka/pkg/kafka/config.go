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

package kafka

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kubernetes/helm/pkg/strvals"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/types"
)

// Config is the config for the kafka collector
type Config struct ***REMOVED***
	// Connection.
	Brokers []string `json:"brokers" envconfig:"K6_KAFKA_BROKERS"`

	// Samples.
	Topic        null.String        `json:"topic" envconfig:"K6_KAFKA_TOPIC"`
	Format       null.String        `json:"format" envconfig:"K6_KAFKA_FORMAT"`
	PushInterval types.NullDuration `json:"push_interval" envconfig:"K6_KAFKA_PUSH_INTERVAL"`

	InfluxDBConfig influxdbConfig `json:"influxdb"`
***REMOVED***

// config is a duplicate of ConfigFields as we can not mapstructure.Decode into
// null types so we duplicate the struct with primitive types to Decode into
type config struct ***REMOVED***
	Brokers      []string `json:"brokers" mapstructure:"brokers" envconfig:"K6_KAFKA_BROKERS"`
	Topic        string   `json:"topic" mapstructure:"topic" envconfig:"K6_KAFKA_TOPIC"`
	Format       string   `json:"format" mapstructure:"format" envconfig:"K6_KAFKA_FORMAT"`
	PushInterval string   `json:"push_interval" mapstructure:"push_interval" envconfig:"K6_KAFKA_PUSH_INTERVAL"`

	InfluxDBConfig influxdbConfig `json:"influxdb" mapstructure:"influxdb"`
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		Format:         null.StringFrom("json"),
		PushInterval:   types.NullDurationFrom(1 * time.Second),
		InfluxDBConfig: newInfluxdbConfig(),
	***REMOVED***
***REMOVED***

func (c Config) Apply(cfg Config) Config ***REMOVED***
	if len(cfg.Brokers) > 0 ***REMOVED***
		c.Brokers = cfg.Brokers
	***REMOVED***
	if cfg.Format.Valid ***REMOVED***
		c.Format = cfg.Format
	***REMOVED***
	if cfg.Topic.Valid ***REMOVED***
		c.Topic = cfg.Topic
	***REMOVED***
	if cfg.PushInterval.Valid ***REMOVED***
		c.PushInterval = cfg.PushInterval
	***REMOVED***
	c.InfluxDBConfig = c.InfluxDBConfig.Apply(cfg.InfluxDBConfig)
	return c
***REMOVED***

// ParseArg takes an arg string and converts it to a config
func ParseArg(arg string) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***
	params, err := strvals.Parse(arg)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	if v, ok := params["brokers"].(string); ok ***REMOVED***
		params["brokers"] = []string***REMOVED***v***REMOVED***
	***REMOVED***

	if v, ok := params["influxdb"].(map[string]interface***REMOVED******REMOVED***); ok ***REMOVED***
		influxConfig, err := influxdbParseMap(v)
		if err != nil ***REMOVED***
			return c, err
		***REMOVED***
		c.InfluxDBConfig = influxConfig
	***REMOVED***
	delete(params, "influxdb")

	if v, ok := params["push_interval"].(string); ok ***REMOVED***
		err := c.PushInterval.UnmarshalText([]byte(v))
		if err != nil ***REMOVED***
			return c, err
		***REMOVED***
	***REMOVED***

	var cfg config
	err = mapstructure.Decode(params, &cfg)
	if err != nil ***REMOVED***
		return c, err
	***REMOVED***

	c.Brokers = cfg.Brokers
	c.Topic = null.StringFrom(cfg.Topic)
	c.Format = null.StringFrom(cfg.Format)

	return c, nil
***REMOVED***

// GetConsolidatedConfig combines ***REMOVED***default config values + JSON config +
// environment vars + arg config values***REMOVED***, and returns the final result.
func GetConsolidatedConfig(jsonRawConf json.RawMessage, env map[string]string, arg string) (Config, error) ***REMOVED***
	result := NewConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf := Config***REMOVED******REMOVED***
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil ***REMOVED***
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

	if arg != "" ***REMOVED***
		urlConf, err := ParseArg(arg)
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(urlConf)
	***REMOVED***

	return result, nil
***REMOVED***
