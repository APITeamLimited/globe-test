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

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/executor"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/stats/csv"
	"github.com/loadimpact/k6/stats/datadog"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/loadimpact/k6/stats/kafka"
	"github.com/loadimpact/k6/stats/statsd"
	"github.com/loadimpact/k6/ui"
)

// configFlagSet returns a FlagSet with the default run configuration flags.
func configFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.StringArrayP("out", "o", []string***REMOVED******REMOVED***, "`uri` for an external metrics database")
	flags.BoolP("linger", "l", false, "keep the API server alive past test end")
	flags.Bool("no-usage-report", false, "don't send anonymous stats to the developers")
	return flags
***REMOVED***

type Config struct ***REMOVED***
	lib.Options

	Out           []string  `json:"out" envconfig:"K6_OUT"`
	Linger        null.Bool `json:"linger" envconfig:"K6_LINGER"`
	NoUsageReport null.Bool `json:"noUsageReport" envconfig:"K6_NO_USAGE_REPORT"`

	Collectors struct ***REMOVED***
		InfluxDB influxdb.Config `json:"influxdb"`
		Kafka    kafka.Config    `json:"kafka"`
		Cloud    cloud.Config    `json:"cloud"`
		StatsD   statsd.Config   `json:"statsd"`
		Datadog  datadog.Config  `json:"datadog"`
		CSV      csv.Config      `json:"csv"`
	***REMOVED*** `json:"collectors"`
***REMOVED***

// Validate checks if all of the specified options make sense
func (c Config) Validate() []error ***REMOVED***
	errors := c.Options.Validate()
	//TODO: validate all of the other options... that we should have already been validating...
	//TODO: maybe integrate an external validation lib: https://github.com/avelino/awesome-go#validation

	return errors
***REMOVED***

func (c Config) Apply(cfg Config) Config ***REMOVED***
	c.Options = c.Options.Apply(cfg.Options)
	if len(cfg.Out) > 0 ***REMOVED***
		c.Out = cfg.Out
	***REMOVED***
	if cfg.Linger.Valid ***REMOVED***
		c.Linger = cfg.Linger
	***REMOVED***
	if cfg.NoUsageReport.Valid ***REMOVED***
		c.NoUsageReport = cfg.NoUsageReport
	***REMOVED***
	c.Collectors.InfluxDB = c.Collectors.InfluxDB.Apply(cfg.Collectors.InfluxDB)
	c.Collectors.Cloud = c.Collectors.Cloud.Apply(cfg.Collectors.Cloud)
	c.Collectors.Kafka = c.Collectors.Kafka.Apply(cfg.Collectors.Kafka)
	c.Collectors.StatsD = c.Collectors.StatsD.Apply(cfg.Collectors.StatsD)
	c.Collectors.Datadog = c.Collectors.Datadog.Apply(cfg.Collectors.Datadog)
	c.Collectors.CSV = c.Collectors.CSV.Apply(cfg.Collectors.CSV)
	return c
***REMOVED***

// Gets configuration from CLI flags.
func getConfig(flags *pflag.FlagSet) (Config, error) ***REMOVED***
	opts, err := getOptions(flags)
	if err != nil ***REMOVED***
		return Config***REMOVED******REMOVED***, err
	***REMOVED***
	out, err := flags.GetStringArray("out")
	if err != nil ***REMOVED***
		return Config***REMOVED******REMOVED***, err
	***REMOVED***
	return Config***REMOVED***
		Options:       opts,
		Out:           out,
		Linger:        getNullBool(flags, "linger"),
		NoUsageReport: getNullBool(flags, "no-usage-report"),
	***REMOVED***, nil
***REMOVED***

// Reads the configuration file from the supplied filesystem and returns it and its path.
// It will first try to see if the user explicitly specified a custom config file and will
// try to read that. If there's a custom config specified and it couldn't be read or parsed,
// an error will be returned.
// If there's no custom config specified and no file exists in the default config path, it will
// return an empty config struct, the default config location and *no* error.
func readDiskConfig(fs afero.Fs) (Config, string, error) ***REMOVED***
	realConfigFilePath := configFilePath
	if realConfigFilePath == "" ***REMOVED***
		// The user didn't specify K6_CONFIG or --config, use the default path
		realConfigFilePath = defaultConfigFilePath
	***REMOVED***

	// Try to see if the file exists in the supplied filesystem
	if _, err := fs.Stat(realConfigFilePath); err != nil ***REMOVED***
		if os.IsNotExist(err) && configFilePath == "" ***REMOVED***
			// If the file doesn't exist, but it was the default config file (i.e. the user
			// didn't specify anything), silence the error
			err = nil
		***REMOVED***
		return Config***REMOVED******REMOVED***, realConfigFilePath, err
	***REMOVED***

	data, err := afero.ReadFile(fs, realConfigFilePath)
	if err != nil ***REMOVED***
		return Config***REMOVED******REMOVED***, realConfigFilePath, err
	***REMOVED***
	var conf Config
	err = json.Unmarshal(data, &conf)
	return conf, realConfigFilePath, err
***REMOVED***

// Serializes the configuration to a JSON file and writes it in the supplied
// location on the supplied filesystem
func writeDiskConfig(fs afero.Fs, configPath string, conf Config) error ***REMOVED***
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := fs.MkdirAll(filepath.Dir(configPath), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	return afero.WriteFile(fs, configPath, data, 0644)
***REMOVED***

// Reads configuration variables from the environment.
func readEnvConfig() (conf Config, err error) ***REMOVED***
	// TODO: replace envconfig and refactor the whole configuration from the groun up :/
	for _, err := range []error***REMOVED***
		envconfig.Process("", &conf),
		envconfig.Process("", &conf.Collectors.Cloud),
		envconfig.Process("", &conf.Collectors.InfluxDB),
		envconfig.Process("", &conf.Collectors.Kafka),
		envconfig.Process("k6_statsd", &conf.Collectors.StatsD),
		envconfig.Process("k6_datadog", &conf.Collectors.Datadog),
	***REMOVED*** ***REMOVED***
		return conf, err
	***REMOVED***
	return conf, nil
***REMOVED***

// Assemble the final consolidated configuration from all of the different sources:
// - start with the CLI-provided options to get shadowed (non-Valid) defaults in there
// - add the global file config options
// - if supplied, add the Runner-provided options
// - add the environment variables
// - merge the user-supplied CLI flags back in on top, to give them the greatest priority
// - set some defaults if they weren't previously specified
// TODO: add better validation, more explicit default values and improve consistency between formats
// TODO: accumulate all errors and differentiate between the layers?
func getConsolidatedConfig(fs afero.Fs, cliConf Config, runner lib.Runner) (conf Config, err error) ***REMOVED***
	cliConf.Collectors.InfluxDB = influxdb.NewConfig().Apply(cliConf.Collectors.InfluxDB)
	cliConf.Collectors.Cloud = cloud.NewConfig().Apply(cliConf.Collectors.Cloud)
	cliConf.Collectors.Kafka = kafka.NewConfig().Apply(cliConf.Collectors.Kafka)
	cliConf.Collectors.StatsD = statsd.NewConfig().Apply(cliConf.Collectors.StatsD)
	cliConf.Collectors.Datadog = datadog.NewConfig().Apply(cliConf.Collectors.Datadog)

	fileConf, _, err := readDiskConfig(fs)
	if err != nil ***REMOVED***
		return conf, err
	***REMOVED***
	envConf, err := readEnvConfig()
	if err != nil ***REMOVED***
		return conf, err
	***REMOVED***

	conf = cliConf.Apply(fileConf)
	if runner != nil ***REMOVED***
		conf = conf.Apply(Config***REMOVED***Options: runner.GetOptions()***REMOVED***)
	***REMOVED***
	conf = conf.Apply(envConf).Apply(cliConf)
	conf = applyDefault(conf)

	// TODO(imiric): Move this validation where it makes sense in the configuration
	// refactor of #883. This repeats the trend stats validation already done
	// for CLI flags in cmd.getOptions, in case other configuration sources
	// (e.g. env vars) overrode our default value. This is not done in
	// lib.Options.Validate to avoid circular imports.
	if err = ui.ValidateSummary(conf.SummaryTrendStats); err != nil ***REMOVED***
		return conf, err
	***REMOVED***

	return conf, nil
***REMOVED***

// applyDefault applies the default options value if it is not specified.
// This happens with types which are not supported by "gopkg.in/guregu/null.v3".
//
// Note that if you add option default value here, also add it in command line argument help text.
func applyDefault(conf Config) Config ***REMOVED***
	if conf.Options.SystemTags == nil ***REMOVED***
		conf.Options.SystemTags = &stats.DefaultSystemTagSet
	***REMOVED***
	if conf.Options.SummaryTrendStats == nil ***REMOVED***
		conf.Options.SummaryTrendStats = lib.DefaultSummaryTrendStats
	***REMOVED***
	defDNS := types.DefaultDNSConfig()
	if !conf.DNS.TTL.Valid ***REMOVED***
		conf.DNS.TTL = defDNS.TTL
	***REMOVED***
	if !conf.DNS.Select.Valid ***REMOVED***
		conf.DNS.Select = defDNS.Select
	***REMOVED***
	if !conf.DNS.Policy.Valid ***REMOVED***
		conf.DNS.Policy = defDNS.Policy
	***REMOVED***

	return conf
***REMOVED***

func deriveAndValidateConfig(conf Config, isExecutable func(string) bool) (result Config, err error) ***REMOVED***
	result = conf
	result.Options, err = executor.DeriveScenariosFromShortcuts(conf.Options)
	if err != nil ***REMOVED***
		return result, err
	***REMOVED***
	return result, validateConfig(result, isExecutable)
***REMOVED***

func validateConfig(conf Config, isExecutable func(string) bool) error ***REMOVED***
	errList := conf.Validate()

	for _, ec := range conf.Scenarios ***REMOVED***
		if err := validateScenarioConfig(ec, isExecutable); err != nil ***REMOVED***
			errList = append(errList, err)
		***REMOVED***
	***REMOVED***

	if len(errList) == 0 ***REMOVED***
		return nil
	***REMOVED***

	errMsgParts := []string***REMOVED***"There were problems with the specified script configuration:"***REMOVED***
	for _, err := range errList ***REMOVED***
		errMsgParts = append(errMsgParts, fmt.Sprintf("\t- %s", err.Error()))
	***REMOVED***

	return errors.New(strings.Join(errMsgParts, "\n"))
***REMOVED***

func validateScenarioConfig(conf lib.ExecutorConfig, isExecutable func(string) bool) error ***REMOVED***
	execFn := conf.GetExec()
	if !isExecutable(execFn) ***REMOVED***
		return fmt.Errorf("executor %s: function '%s' not found in exports", conf.GetName(), execFn)
	***REMOVED***
	return nil
***REMOVED***
