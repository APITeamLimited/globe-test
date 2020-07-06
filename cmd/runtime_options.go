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
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
)

// TODO: move this whole file out of the cmd package? maybe when fixing
// https://github.com/loadimpact/k6/issues/883, since this code is fairly
// self-contained and easily testable now, without any global dependencies...

var userEnvVarName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func parseEnvKeyValue(kv string) (string, string) ***REMOVED***
	if idx := strings.IndexRune(kv, '='); idx != -1 ***REMOVED***
		return kv[:idx], kv[idx+1:]
	***REMOVED***
	return kv, ""
***REMOVED***

func buildEnvMap(environ []string) map[string]string ***REMOVED***
	env := make(map[string]string, len(environ))
	for _, kv := range environ ***REMOVED***
		k, v := parseEnvKeyValue(kv)
		env[k] = v
	***REMOVED***
	return env
***REMOVED***

func runtimeOptionFlagSet(includeSysEnv bool) *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.Bool("include-system-env-vars", includeSysEnv, "pass the real system environment variables to the runtime")
	flags.String("compatibility-mode", "extended",
		`JavaScript compiler compatibility mode, "extended" or "base"
base: pure Golang JS VM supporting ES5.1+
extended: base + Babel with ES2015 preset + core.js v2,
          slower and memory consuming but with greater JS support
`)
	flags.StringArrayP("env", "e", nil, "add/override environment variable with `VAR=value`")
	return flags
***REMOVED***

func getRuntimeOptions(flags *pflag.FlagSet, environment map[string]string) (lib.RuntimeOptions, error) ***REMOVED***
	opts := lib.RuntimeOptions***REMOVED***
		IncludeSystemEnvVars: getNullBool(flags, "include-system-env-vars"),
		CompatibilityMode:    getNullString(flags, "compatibility-mode"),
		Env:                  make(map[string]string),
	***REMOVED***

	if !opts.CompatibilityMode.Valid ***REMOVED*** // If not explicitly set via CLI flags, look for an environment variable
		if envVar, ok := environment["K6_COMPATIBILITY_MODE"]; ok ***REMOVED***
			opts.CompatibilityMode = null.StringFrom(envVar)
		***REMOVED***
	***REMOVED***
	if _, err := lib.ValidateCompatibilityMode(opts.CompatibilityMode.String); err != nil ***REMOVED***
		// some early validation
		return opts, err
	***REMOVED***

	if !opts.IncludeSystemEnvVars.Valid ***REMOVED*** // If not explicitly set via CLI flags, look for an environment variable
		if envVar, ok := environment["K6_INCLUDE_SYSTEM_ENV_VARS"]; ok ***REMOVED***
			val, err := strconv.ParseBool(envVar)
			if err != nil ***REMOVED***
				return opts, err
			***REMOVED***
			opts.IncludeSystemEnvVars = null.BoolFrom(val)
		***REMOVED***
	***REMOVED***

	if opts.IncludeSystemEnvVars.Bool ***REMOVED*** // If enabled, gather the actual system environment variables
		opts.Env = environment
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
