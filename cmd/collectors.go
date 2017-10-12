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
	"encoding"
	"strings"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats/influxdb"
	jsonc "github.com/loadimpact/k6/stats/json"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func parseCollector(s string) (t, arg string) ***REMOVED***
	parts := strings.SplitN(s, "=", 2)
	switch len(parts) ***REMOVED***
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	***REMOVED***
***REMOVED***

func newCollector(t, arg string, src *lib.SourceData, conf Config) (lib.Collector, error) ***REMOVED***
	loadConfig := func(out encoding.TextUnmarshaler) error ***REMOVED***
		if err := conf.ConfigureCollector(t, out); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := out.UnmarshalText([]byte(arg)); err != nil ***REMOVED***
			return err
		***REMOVED***
		return nil
	***REMOVED***

	switch t ***REMOVED***
	case collectorJSON:
		return jsonc.New(afero.NewOsFs(), arg)
	case collectorInfluxDB:
		var config influxdb.Config
		if err := loadConfig(&config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return influxdb.New(config)
	// case collectorCloud:
	// 	return cloud.New(arg, src, conf.Options, Version)
	default:
		return nil, errors.Errorf("unknown output type: %s", t)
	***REMOVED***
***REMOVED***
