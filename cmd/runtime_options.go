package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/pflag"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
)

// TODO: move this whole file out of the cmd package? maybe when fixing
// https://github.com/k6io/k6/issues/883, since this code is fairly
// self-contained and easily testable now, without any global dependencies...

var userEnvVarName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func runtimeOptionFlagSet(includeSysEnv bool) *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.Bool("include-system-env-vars", includeSysEnv, "pass the real system environment variables to the runtime")
	flags.String("compatibility-mode", "extended",
		`JavaScript compiler compatibility mode, "extended" or "base"
base: pure goja - Golang JS VM supporting ES5.1+
extended: base + Babel with parts of ES2015 preset
		  slower to compile in case the script uses syntax unsupported by base
`)
	flags.StringP("type", "t", "", "override test type, \"js\" or \"archive\"")
	flags.StringArrayP("env", "e", nil, "add/override environment variable with `VAR=value`")
	flags.Bool("no-thresholds", false, "don't run thresholds")
	flags.Bool("no-summary", false, "don't show the summary at the end of the test")
	flags.String(
		"summary-export",
		"",
		"output the end-of-test summary report to JSON file",
	)
	return flags
***REMOVED***

func saveBoolFromEnv(env map[string]string, varName string, placeholder *null.Bool) error ***REMOVED***
	strValue, ok := env[varName]
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	val, err := strconv.ParseBool(strValue)
	if err != nil ***REMOVED***
		return fmt.Errorf("env var '%s' is not a valid boolean value: %w", varName, err)
	***REMOVED***
	// Only override if not explicitly set via the CLI flag
	if !placeholder.Valid ***REMOVED***
		*placeholder = null.BoolFrom(val)
	***REMOVED***
	return nil
***REMOVED***

func getRuntimeOptions(flags *pflag.FlagSet, environment map[string]string) (lib.RuntimeOptions, error) ***REMOVED***
	// TODO: refactor with composable helpers as a part of #883, to reduce copy-paste
	// TODO: get these options out of the JSON config file as well?
	opts := lib.RuntimeOptions***REMOVED***
		TestType:             getNullString(flags, "type"),
		IncludeSystemEnvVars: getNullBool(flags, "include-system-env-vars"),
		CompatibilityMode:    getNullString(flags, "compatibility-mode"),
		NoThresholds:         getNullBool(flags, "no-thresholds"),
		NoSummary:            getNullBool(flags, "no-summary"),
		SummaryExport:        getNullString(flags, "summary-export"),
		Env:                  make(map[string]string),
	***REMOVED***

	if envVar, ok := environment["K6_TYPE"]; ok && !opts.TestType.Valid ***REMOVED***
		// Only override if not explicitly set via the CLI flag
		opts.TestType = null.StringFrom(envVar)
	***REMOVED***
	if envVar, ok := environment["K6_COMPATIBILITY_MODE"]; ok && !opts.CompatibilityMode.Valid ***REMOVED***
		// Only override if not explicitly set via the CLI flag
		opts.CompatibilityMode = null.StringFrom(envVar)
	***REMOVED***
	if _, err := lib.ValidateCompatibilityMode(opts.CompatibilityMode.String); err != nil ***REMOVED***
		// some early validation
		return opts, err
	***REMOVED***

	if err := saveBoolFromEnv(environment, "K6_INCLUDE_SYSTEM_ENV_VARS", &opts.IncludeSystemEnvVars); err != nil ***REMOVED***
		return opts, err
	***REMOVED***
	if err := saveBoolFromEnv(environment, "K6_NO_THRESHOLDS", &opts.NoThresholds); err != nil ***REMOVED***
		return opts, err
	***REMOVED***
	if err := saveBoolFromEnv(environment, "K6_NO_SUMMARY", &opts.NoSummary); err != nil ***REMOVED***
		return opts, err
	***REMOVED***

	if envVar, ok := environment["K6_SUMMARY_EXPORT"]; ok ***REMOVED***
		if !opts.SummaryExport.Valid ***REMOVED***
			opts.SummaryExport = null.StringFrom(envVar)
		***REMOVED***
	***REMOVED***

	if envVar, ok := environment["SSLKEYLOGFILE"]; ok ***REMOVED***
		if !opts.KeyWriter.Valid ***REMOVED***
			opts.KeyWriter = null.StringFrom(envVar)
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
			return opts, fmt.Errorf("invalid environment variable name '%s'", k)
		***REMOVED***
		opts.Env[k] = v
	***REMOVED***

	return opts, nil
***REMOVED***
