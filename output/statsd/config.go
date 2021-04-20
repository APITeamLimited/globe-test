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
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
)

// config defines the StatsD configuration.
type config struct ***REMOVED***
	Addr         null.String        `json:"addr,omitempty" envconfig:"K6_STATSD_ADDR"`
	BufferSize   null.Int           `json:"bufferSize,omitempty" envconfig:"K6_STATSD_BUFFER_SIZE"`
	Namespace    null.String        `json:"namespace,omitempty" envconfig:"K6_STATSD_NAMESPACE"`
	PushInterval types.NullDuration `json:"pushInterval,omitempty" envconfig:"K6_STATSD_PUSH_INTERVAL"`
	TagBlocklist stats.TagSet       `json:"tagBlocklist,omitempty" envconfig:"K6_STATSD_TAG_BLOCKLIST"`
	EnableTags   null.Bool          `json:"enableTags,omitempty" envconfig:"K6_STATSD_ENABLE_TAGS"`
***REMOVED***

func processTags(t stats.TagSet, tags map[string]string) []string ***REMOVED***
	var res []string
	for key, value := range tags ***REMOVED***
		if value != "" && !t[key] ***REMOVED***
			res = append(res, key+":"+value)
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c config) Apply(cfg config) config ***REMOVED***
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
	if cfg.TagBlocklist != nil ***REMOVED***
		c.TagBlocklist = cfg.TagBlocklist
	***REMOVED***
	if cfg.EnableTags.Valid ***REMOVED***
		c.EnableTags = cfg.EnableTags
	***REMOVED***

	return c
***REMOVED***

// newConfig creates a new Config instance with default values for some fields.
func newConfig() config ***REMOVED***
	return config***REMOVED***
		Addr:         null.NewString("localhost:8125", false),
		BufferSize:   null.NewInt(20, false),
		Namespace:    null.NewString("k6.", false),
		PushInterval: types.NewNullDuration(1*time.Second, false),
		TagBlocklist: stats.TagSet***REMOVED******REMOVED***,
		EnableTags:   null.NewBool(false, false),
	***REMOVED***
***REMOVED***

// getConsolidatedConfig combines ***REMOVED***default config values + JSON config +
// environment vars***REMOVED***, and returns the final result.
func getConsolidatedConfig(jsonRawConf json.RawMessage, env map[string]string, _ string) (config, error) ***REMOVED***
	result := newConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf := config***REMOVED******REMOVED***
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(jsonConf)
	***REMOVED***

	envConfig := config***REMOVED******REMOVED***
	_ = env // TODO: get rid of envconfig and actually use the env parameter...
	if err := envconfig.Process("", &envConfig); err != nil ***REMOVED***
		return result, err
	***REMOVED***
	result = result.Apply(envConfig)

	return result, nil
***REMOVED***
