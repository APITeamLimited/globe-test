/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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
	"os"
	"regexp"
	"strings"

	"github.com/loadimpact/k6/lib"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var userEnvVarName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func parseEnvKeyValue(kv string) (string, string) ***REMOVED***
	if idx := strings.IndexRune(kv, '='); idx != -1 ***REMOVED***
		return kv[:idx], kv[idx+1:]
	***REMOVED***
	return kv, ""
***REMOVED***

func collectEnv() map[string]string ***REMOVED***
	env := make(map[string]string)
	for _, kv := range os.Environ() ***REMOVED***
		k, v := parseEnvKeyValue(kv)
		env[k] = v
	***REMOVED***
	return env
***REMOVED***

func runtimeOptionFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.Bool("no-system-env-vars", false, "don't pass actual system environment variables to the runtime")
	flags.StringSliceP("env", "e", nil, "add/override environment variable with `VAR=value`")
	return flags
***REMOVED***

func getRuntimeOptions(flags *pflag.FlagSet) (lib.RuntimeOptions, error) ***REMOVED***
	opts := lib.RuntimeOptions***REMOVED***
		NoSystemEnvVars: getNullBool(flags, "no-system-env-vars"),
	***REMOVED***

	// If not disabled, gather the actual system environment variables
	if !opts.NoSystemEnvVars.Bool ***REMOVED***
		opts.Env = collectEnv()
	***REMOVED***

	// Set/overwrite environment varialbes with custom user-supplied values
	envVars, err := flags.GetStringSlice("env")
	if err != nil ***REMOVED***
		return opts, err
	***REMOVED***
	if len(envVars) > 0 ***REMOVED***
		// Initialize opts.Env if NoSystemEnvVars was enabled
		if opts.Env == nil ***REMOVED***
			opts.Env = make(map[string]string)
		***REMOVED***

		for _, kv := range envVars ***REMOVED***
			k, v := parseEnvKeyValue(kv)
			// Allow only alphanumeric ASCII variable names for now
			if !userEnvVarName.MatchString(k) ***REMOVED***
				return opts, errors.Errorf("Invalid environment variable name '%s'", k)
			***REMOVED***
			opts.Env[k] = v
		***REMOVED***
	***REMOVED***

	return opts, nil
***REMOVED***
