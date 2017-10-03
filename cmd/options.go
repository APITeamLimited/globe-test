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
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var optionFlagSet = pflag.NewFlagSet("", 0)

func init() ***REMOVED***
	optionFlagSet.SortFlags = false
	optionFlagSet.Int64P("vus", "u", 1, "number of virtual users")
	optionFlagSet.Int64P("max", "m", 0, "max available virtual users")
	optionFlagSet.DurationP("duration", "d", 0, "test duration limit")
	optionFlagSet.Int64P("iterations", "i", 0, "script iteration limit")
	optionFlagSet.StringSliceP("stage", "s", nil, "add a `stage`, as `[duration]:[target]`")
	optionFlagSet.BoolP("paused", "p", false, "start the test in a paused state")
	optionFlagSet.Int64("max-redirects", 10, "follow at most n redirects")
	optionFlagSet.String("user-agent", "", "user agent for http requests")
	optionFlagSet.Bool("insecure-skip-tls-verify", false, "skip verification of TLS certificates")
	optionFlagSet.Bool("no-connection-reuse", false, "don't reuse connections between iterations")
	optionFlagSet.BoolP("throw", "w", false, "throw warnings (like failed http requests) as errors")
***REMOVED***

func getOptions(flags *pflag.FlagSet) (lib.Options, error) ***REMOVED***
	opts := lib.Options***REMOVED***
		VUs:                   getNullInt64(flags, "vus"),
		VUsMax:                getNullInt64(flags, "max"),
		Duration:              getNullDuration(flags, "duration"),
		Iterations:            getNullInt64(flags, "iterations"),
		Paused:                getNullBool(flags, "paused"),
		MaxRedirects:          getNullInt64(flags, "max-redirects"),
		UserAgent:             getNullString(flags, "user-agent"),
		InsecureSkipTLSVerify: getNullBool(flags, "insecure-skip-tls-verify"),
		NoConnectionReuse:     getNullBool(flags, "no-connection-reuse"),
		Throw:                 getNullBool(flags, "throw"),
	***REMOVED***

	stageStrings, err := flags.GetStringSlice("stage")
	if err != nil ***REMOVED***
		return opts, err
	***REMOVED***
	if len(stageStrings) > 0 ***REMOVED***
		opts.Stages = make([]lib.Stage, len(stageStrings))
		for i, s := range stageStrings ***REMOVED***
			var stage lib.Stage
			if err := stage.UnmarshalText([]byte(s)); err != nil ***REMOVED***
				return opts, errors.Wrapf(err, "stage %d", i)
			***REMOVED***
			opts.Stages[i] = stage
		***REMOVED***
	***REMOVED***
	return opts, nil
***REMOVED***
