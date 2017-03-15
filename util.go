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

package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/loadimpact/k6/lib"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
)

func dumpYAML(v interface***REMOVED******REMOVED***) error ***REMOVED***
	bytes, err := yaml.Marshal(v)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := os.Stdout.Write(bytes); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// cliBool returns a CLI argument as a bool, which is invalid if not given.
func cliBool(cc *cli.Context, name string) null.Bool ***REMOVED***
	return null.NewBool(cc.Bool(name), cc.IsSet(name))
***REMOVED***

// cliInt64 returns a CLI argument as an int64, which is invalid if not given.
func cliInt64(cc *cli.Context, name string) null.Int ***REMOVED***
	return null.NewInt(cc.Int64(name), cc.IsSet(name))
***REMOVED***

// cliDuration returns a CLI argument as a duration string, which is invalid if not given.
func cliDuration(cc *cli.Context, name string) null.String ***REMOVED***
	return null.NewString(cc.Duration(name).String(), cc.IsSet(name))
***REMOVED***

func roundDuration(d, to time.Duration) time.Duration ***REMOVED***
	return d - (d % to)
***REMOVED***

func ParseStage(s string) (lib.Stage, error) ***REMOVED***
	parts := strings.SplitN(s, ":", 2)

	var stage lib.Stage
	if parts[0] != "" ***REMOVED***
		d, err := time.ParseDuration(parts[0])
		if err != nil ***REMOVED***
			return stage, err
		***REMOVED***
		stage.Duration = d
	***REMOVED***
	if len(parts) > 1 && parts[1] != "" ***REMOVED***
		vus, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil ***REMOVED***
			return stage, err
		***REMOVED***
		stage.Target = null.IntFrom(vus)
	***REMOVED***
	return stage, nil
***REMOVED***
