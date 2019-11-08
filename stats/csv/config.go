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

package csv

import (
	"fmt"
	"strings"
	"time"

	"github.com/loadimpact/k6/lib/types"
	"gopkg.in/guregu/null.v3"
)

// Config is the config for the csv collector
type Config struct ***REMOVED***
	// Samples.
	FileName     null.String        `json:"file_name" envconfig:"K6_CSV_FILENAME"`
	SaveInterval types.NullDuration `json:"save_interval" envconfig:"K6_CSV_SAVE_INTERVAL"`
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		FileName:     null.StringFrom("file.csv"),
		SaveInterval: types.NullDurationFrom(1 * time.Second),
	***REMOVED***
***REMOVED***

// Apply merges two configs by overwriting properties in the old config
func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.FileName.Valid ***REMOVED***
		c.FileName = cfg.FileName
	***REMOVED***
	if cfg.SaveInterval.Valid ***REMOVED***
		c.SaveInterval = cfg.SaveInterval
	***REMOVED***
	return c
***REMOVED***

// ParseArg takes an arg string and converts it to a config
func ParseArg(arg string) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***

	if !strings.Contains(arg, "=") ***REMOVED***
		c.FileName = null.StringFrom(arg)
		c.SaveInterval = types.NullDurationFrom(1 * time.Second)
		return c, nil
	***REMOVED***

	pairs := strings.Split(arg, ",")
	for _, pair := range pairs ***REMOVED***
		r := strings.SplitN(pair, "=", 2)
		if len(r) != 2 ***REMOVED***
			return c, fmt.Errorf("couldn't parse %q as argument for csv output", arg)
		***REMOVED***
		switch r[0] ***REMOVED***
		case "save_interval":
			err := c.SaveInterval.UnmarshalText([]byte(r[1]))
			if err != nil ***REMOVED***
				return c, err
			***REMOVED***
		case "file_name":
			c.FileName = null.StringFrom(r[1])
		default:
			return c, fmt.Errorf("unknown key %q as argument for csv output", r[0])
		***REMOVED***
	***REMOVED***

	return c, nil
***REMOVED***
