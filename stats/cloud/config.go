/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package cloud

import (
	"time"

	"github.com/loadimpact/k6/lib/types"
	"gopkg.in/guregu/null.v3"
)

// Config holds all the neccessary data and options for sending metrics to the Load Impact cloud.
type Config struct ***REMOVED***
	Token              null.String        `json:"token" envconfig:"CLOUD_TOKEN"`
	DeprecatedToken    null.String        `envconfig:"K6CLOUD_TOKEN"`
	Name               null.String        `json:"name" envconfig:"CLOUD_NAME"`
	Host               null.String        `json:"host" envconfig:"CLOUD_HOST"`
	NoCompress         null.Bool          `json:"noCompress" envconfig:"CLOUD_NO_COMPRESS"`
	ProjectID          null.Int           `json:"projectID" envconfig:"CLOUD_PROJECT_ID"`
	MetricPushInterval types.NullDuration `json:"metricPushInterval" envconfig:"CLOUD_METRIC_PUSH_INTERVAL"`

	// If specified and greater than 0, sample aggregation with that period is enabled.
	AggregationPeriod types.NullDuration `json:"aggregationPeriod" envconfig:"CLOUD_AGGREGATION_PERIOD"`

	// If aggregation is enabled, this specifies how long we'll wait for period samples to accomulate before pushing them to the cloud.
	AggregationPushDelay types.NullDuration `json:"aggregationPushDelay" envconfig:"CLOUD_AGGREGATION_PUSH_DELAY"`

	// If AggregationPeriod is positive, but the collected samples for a certain period are less than this number, they won't be aggregated.
	AggregationMinSamples null.Int `json:"aggregationMinSamples" envconfig:"CLOUD_AGGREGATION_MIN_SAMPLES"`
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		Host:                  null.StringFrom("https://ingest.loadimpact.com"),
		MetricPushInterval:    types.NullDurationFrom(1 * time.Second),
		AggregationPushDelay:  types.NullDurationFrom(3 * time.Second),
		AggregationMinSamples: null.IntFrom(100),
	***REMOVED***
***REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.Token.Valid ***REMOVED***
		c.Token = cfg.Token
	***REMOVED***
	if cfg.DeprecatedToken.Valid ***REMOVED***
		c.DeprecatedToken = cfg.DeprecatedToken
	***REMOVED***
	if cfg.Name.Valid ***REMOVED***
		c.Name = cfg.Name
	***REMOVED***
	if cfg.Host.Valid ***REMOVED***
		c.Host = cfg.Host
	***REMOVED***
	if cfg.NoCompress.Valid ***REMOVED***
		c.NoCompress = cfg.NoCompress
	***REMOVED***
	if cfg.ProjectID.Valid ***REMOVED***
		c.ProjectID = cfg.ProjectID
	***REMOVED***
	if cfg.MetricPushInterval.Valid ***REMOVED***
		c.MetricPushInterval = cfg.MetricPushInterval
	***REMOVED***
	if cfg.AggregationPeriod.Valid ***REMOVED***
		c.AggregationPeriod = cfg.AggregationPeriod
	***REMOVED***
	if cfg.AggregationPushDelay.Valid ***REMOVED***
		c.AggregationPushDelay = cfg.AggregationPushDelay
	***REMOVED***
	if cfg.AggregationMinSamples.Valid ***REMOVED***
		c.AggregationMinSamples = cfg.AggregationMinSamples
	***REMOVED***
	return c
***REMOVED***
