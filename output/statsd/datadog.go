/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package statsd

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/output"
	"go.k6.io/k6/stats"
	"gopkg.in/guregu/null.v3"
)

// TODO delete this whole file

type datadogConfig struct ***REMOVED***
	Addr         null.String        `json:"addr,omitempty" envconfig:"K6_DATADOG_ADDR"`
	BufferSize   null.Int           `json:"bufferSize,omitempty" envconfig:"K6_DATADOG_BUFFER_SIZE"`
	Namespace    null.String        `json:"namespace,omitempty" envconfig:"K6_DATADOG_NAMESPACE"`
	PushInterval types.NullDuration `json:"pushInterval,omitempty" envconfig:"K6_DATADOG_PUSH_INTERVAL"`
	TagBlacklist stats.TagSet       `json:"tagBlacklist,omitempty" envconfig:"K6_DATADOG_TAG_BLACKLIST"`
***REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c datadogConfig) Apply(cfg datadogConfig) datadogConfig ***REMOVED***
	if cfg.Addr.Valid ***REMOVED***
		c.Addr = cfg.Addr
	***REMOVED***
	if cfg.BufferSize.Valid ***REMOVED***
		c.BufferSize = cfg.BufferSize
	***REMOVED***
	if cfg.Namespace.Valid ***REMOVED***
		c.Namespace = cfg.Namespace
	***REMOVED***
	if cfg.PushInterval.Valid ***REMOVED***
		c.PushInterval = cfg.PushInterval
	***REMOVED***
	if cfg.TagBlacklist != nil ***REMOVED***
		c.TagBlacklist = cfg.TagBlacklist
	***REMOVED***

	return c
***REMOVED***

// NewdatadogConfig creates a new datadogConfig instance with default values for some fields.
func newdatadogConfig() datadogConfig ***REMOVED***
	return datadogConfig***REMOVED***
		Addr:         null.NewString("localhost:8125", false),
		BufferSize:   null.NewInt(20, false),
		Namespace:    null.NewString("k6.", false),
		PushInterval: types.NewNullDuration(1*time.Second, false),
		TagBlacklist: (stats.TagVU | stats.TagIter | stats.TagURL).Map(),
	***REMOVED***
***REMOVED***

// GetConsolidateddatadogConfig combines ***REMOVED***default config values + JSON config +
// environment vars***REMOVED***, and returns the final result.
func getConsolidatedDatadogConfig(jsonRawConf json.RawMessage) (datadogConfig, error) ***REMOVED***
	result := newdatadogConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf := datadogConfig***REMOVED******REMOVED***
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(jsonConf)
	***REMOVED***

	envdatadogConfig := datadogConfig***REMOVED******REMOVED***
	if err := envconfig.Process("", &envdatadogConfig); err != nil ***REMOVED***
		// TODO: get rid of envconfig and actually use the env parameter...
		return result, err
	***REMOVED***
	result = result.Apply(envdatadogConfig)

	return result, nil
***REMOVED***

// NewDatadog creates a new statsd connector client with tags enabled
// TODO delete this
func NewDatadog(params output.Params) (*Output, error) ***REMOVED***
	conf, err := getConsolidatedDatadogConfig(params.JSONConfig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	logger := params.Logger.WithFields(logrus.Fields***REMOVED***"output": "statsd"***REMOVED***)
	statsdConfig := config***REMOVED***
		Addr:         conf.Addr,
		BufferSize:   conf.BufferSize,
		Namespace:    conf.Namespace,
		PushInterval: conf.PushInterval,
		TagBlocklist: conf.TagBlacklist,
		EnableTags:   null.NewBool(true, false),
	***REMOVED***

	return &Output***REMOVED***
		config: statsdConfig,
		logger: logger,
	***REMOVED***, nil
***REMOVED***
