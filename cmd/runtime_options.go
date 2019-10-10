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

func runtimeOptionFlagSet(includeSysEnv bool) *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.Bool("include-system-env-vars", includeSysEnv, "pass the real system environment variables to the runtime")
	flags.StringArrayP("env", "e", nil, "add/override environment variable with `VAR=value`")
	return flags
***REMOVED***

func getRuntimeOptions(flags *pflag.FlagSet) (lib.RuntimeOptions, error) ***REMOVED***
	opts := lib.RuntimeOptions***REMOVED***
		IncludeSystemEnvVars: getNullBool(flags, "include-system-env-vars"),
		Env:                  make(map[string]string),
	***REMOVED***

	// If enabled, gather the actual system environment variables
	if opts.IncludeSystemEnvVars.Bool ***REMOVED***
		opts.Env = collectEnv()
	***REMOVED***

	// Set/overwrite environment variables with custom user-supplied values
	envVars, err := flags.GetStringArray("env")
	if err != nil ***REMOVED***
		return opts, err
	***REMOVED***

	for _, kv := range envVars ***REMOVED***
		k, v := parseEnvKeyValue(kv)
		// Allow only alphanumeric ASCII variable names for now
		if !userEnvVarName.MatchString(k) ***REMOVED***
			return opts, errors.Errorf("Invalid environment variable name '%s'", k)
		***REMOVED***
		opts.Env[k] = v
	***REMOVED***

	return opts, nil
***REMOVED***
