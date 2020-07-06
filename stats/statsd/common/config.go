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

package common

import (
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/types"
)

// Config defines the statsd configuration
type Config struct ***REMOVED***
	Addr         null.String        `json:"addr,omitempty" envconfig:"ADDR"`
	BufferSize   null.Int           `json:"bufferSize,omitempty" envconfig:"BUFFER_SIZE"`
	Namespace    null.String        `json:"namespace,omitempty" envconfig:"NAMESPACE"`
	PushInterval types.NullDuration `json:"pushInterval,omitempty" envconfig:"PUSH_INTERVAL"`
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		Addr:         null.NewString("localhost:8125", false),
		BufferSize:   null.NewInt(20, false),
		Namespace:    null.NewString("k6.", false),
		PushInterval: types.NewNullDuration(1*time.Second, false),
	***REMOVED***
***REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c Config) Apply(cfg Config) Config ***REMOVED***
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

	return c
***REMOVED***
