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

package datadog

import (
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/statsd/common"
)

type tagHandler stats.SystemTagSet

func (t tagHandler) processTags(tags map[string]string) []string ***REMOVED***
	var res []string
	for key, value := range tags ***REMOVED***
		if v, err := stats.SystemTagSetString(key); err == nil ***REMOVED***
			if value != "" && t&tagHandler(v) != 0 ***REMOVED***
				res = append(res, key+":"+value)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

// Config defines the datadog configuration
type Config struct ***REMOVED***
	common.Config

	TagBlacklist stats.SystemTagMap `json:"tagBlacklist,omitempty" envconfig:"TAG_BLACKLIST"`
***REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c Config) Apply(cfg Config) Config ***REMOVED***
	c.Config = c.Config.Apply(cfg.Config)

	if cfg.TagBlacklist != nil ***REMOVED***
		c.TagBlacklist = cfg.TagBlacklist
	***REMOVED***

	return c
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		Config:       common.NewConfig(),
		TagBlacklist: stats.SystemTagMap***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// New creates a new statsd connector client
func New(conf Config) (*common.Collector, error) ***REMOVED***
	return &common.Collector***REMOVED***
		Config:      conf.Config,
		Type:        "datadog",
		ProcessTags: tagHandler(conf.TagBlacklist.ToTagSet()).processTags,
	***REMOVED***, nil
***REMOVED***
