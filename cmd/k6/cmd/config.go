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

package cmd

import (
	"github.com/loadimpact/k6/lib"
	null "gopkg.in/guregu/null.v3"
)

type Config struct ***REMOVED***
	lib.Options

	Linger        null.Bool `json:"linger"`        // DEPRECATED; will be removed.
	NoUsageReport null.Bool `json:"noUsageReport"` // DEPRECATED; will be removed.
***REMOVED***

func ConfigFromEnv(env []string) ***REMOVED***
***REMOVED***

func (c Config) Apply(cfg Config) Config ***REMOVED***
	c.Options = c.Options.Apply(cfg.Options)
	if cfg.Linger.Valid ***REMOVED***
		c.Linger = cfg.Linger
	***REMOVED***
	if cfg.NoUsageReport.Valid ***REMOVED***
		c.NoUsageReport = cfg.NoUsageReport
	***REMOVED***
	return c
***REMOVED***
