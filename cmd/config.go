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
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/scheduler"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/loadimpact/k6/stats/kafka"
	"github.com/pkg/errors"
	"github.com/shibukawa/configdir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	null "gopkg.in/guregu/null.v3"
)

const configFilename = "config.json"

var configDirs = configdir.New("loadimpact", "k6")
var configFile = os.Getenv("K6_CONFIG") // overridden by `-c` flag!

// configFileFlagSet returns a FlagSet that contains flags needed for specifying a config file.
func configFileFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.StringVarP(&configFile, "config", "c", configFile, "specify config file to read")
	return flags
***REMOVED***

// configFlagSet returns a FlagSet with the default run configuration flags.
func configFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.StringArrayP("out", "o", []string***REMOVED******REMOVED***, "`uri` for an external metrics database")
	flags.BoolP("linger", "l", false, "keep the API server alive past test end")
	flags.Bool("no-usage-report", false, "don't send anonymous stats to the developers")
	flags.Bool("no-thresholds", false, "don't run thresholds")
	flags.Bool("no-summary", false, "don't show the summary at the end of the test")
	flags.AddFlagSet(configFileFlagSet())
	return flags
***REMOVED***

type Config struct ***REMOVED***
	lib.Options

	Out           []string  `json:"out" envconfig:"out"`
	Linger        null.Bool `json:"linger" envconfig:"linger"`
	NoUsageReport null.Bool `json:"noUsageReport" envconfig:"no_usage_report"`
	NoThresholds  null.Bool `json:"noThresholds" envconfig:"no_thresholds"`
	NoSummary     null.Bool `json:"noSummary" envconfig:"no_summary"`

	Collectors struct ***REMOVED***
		InfluxDB influxdb.Config `json:"influxdb"`
		Kafka    kafka.Config    `json:"kafka"`
		Cloud    cloud.Config    `json:"cloud"`
	***REMOVED*** `json:"collectors"`
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
	if cfg.NoThresholds.Valid ***REMOVED***
		c.NoThresholds = cfg.NoThresholds
	***REMOVED***
	if cfg.NoSummary.Valid ***REMOVED***
		c.NoSummary = cfg.NoSummary
	***REMOVED***
	c.Collectors.InfluxDB = c.Collectors.InfluxDB.Apply(cfg.Collectors.InfluxDB)
	c.Collectors.Cloud = c.Collectors.Cloud.Apply(cfg.Collectors.Cloud)
	c.Collectors.Kafka = c.Collectors.Kafka.Apply(cfg.Collectors.Kafka)
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
		NoThresholds:  getNullBool(flags, "no-thresholds"),
		NoSummary:     getNullBool(flags, "no-summary"),
	***REMOVED***, nil
***REMOVED***

// Reads a configuration file from disk.
func readDiskConfig(fs afero.Fs) (Config, *configdir.Config, error) ***REMOVED***
	if configFile != "" ***REMOVED***
		data, err := ioutil.ReadFile(configFile)
		if err != nil ***REMOVED***
			return Config***REMOVED******REMOVED***, nil, err
		***REMOVED***
		var conf Config
		err = json.Unmarshal(data, &conf)
		return conf, nil, err
	***REMOVED***

	cdir := configDirs.QueryFolderContainsFile(configFilename)
	if cdir == nil ***REMOVED***
		return Config***REMOVED******REMOVED***, configDirs.QueryFolders(configdir.Global)[0], nil
	***REMOVED***
	data, err := cdir.ReadFile(configFilename)
	if err != nil ***REMOVED***
		return Config***REMOVED******REMOVED***, cdir, err
	***REMOVED***
	var conf Config
	err = json.Unmarshal(data, &conf)
	return conf, cdir, err
***REMOVED***

// Writes configuration back to disk.
func writeDiskConfig(fs afero.Fs, cdir *configdir.Config, conf Config) error ***REMOVED***
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if configFile != "" ***REMOVED***
		return afero.WriteFile(fs, configFilename, data, 0644)
	***REMOVED***
	return cdir.WriteFile(configFilename, data)
***REMOVED***

// Reads configuration variables from the environment.
func readEnvConfig() (conf Config, err error) ***REMOVED***
	// TODO: replace envconfig and refactor the whole configuration from the groun up :/
	for _, err := range []error***REMOVED***
		envconfig.Process("k6", &conf),
		envconfig.Process("k6", &conf.Collectors.Cloud),
		envconfig.Process("k6", &conf.Collectors.InfluxDB),
		envconfig.Process("k6", &conf.Collectors.Kafka),
	***REMOVED*** ***REMOVED***
		return conf, err
	***REMOVED***
	return conf, nil
***REMOVED***

// This checks for conflicting options and turns any shortcut options (i.e. duration, iterations,
// stages) into the proper scheduler configuration
func buildExecutionConfig(conf Config) (Config, error) ***REMOVED***
	result := conf
	if conf.Duration.Valid ***REMOVED***
		if conf.Iterations.Valid ***REMOVED***
			//TODO: make this an error in the next version
			log.Warnf("Specifying both duration and iterations is deprecated and won't be supported in the future k6 versions")
		***REMOVED***

		if conf.Stages != nil ***REMOVED***
			//TODO: make this an error in the next version
			log.Warnf("Specifying both duration and stages is deprecated and won't be supported in the future k6 versions")
		***REMOVED***

		if conf.Execution != nil ***REMOVED***
			return result, errors.New("specifying both duration and execution is not supported")
		***REMOVED***

		ds := scheduler.NewConstantLoopingVUsConfig(lib.DefaultSchedulerName)
		ds.VUs = conf.VUs
		ds.Duration = conf.Duration
		result.Execution = scheduler.ConfigMap***REMOVED***lib.DefaultSchedulerName: ds***REMOVED***
	***REMOVED*** else if conf.Stages != nil ***REMOVED***
		if conf.Iterations.Valid ***REMOVED***
			//TODO: make this an error in the next version
			log.Warnf("Specifying both iterations and stages is deprecated and won't be supported in the future k6 versions")
		***REMOVED***

		if conf.Execution != nil ***REMOVED***
			return conf, errors.New("specifying both stages and execution is not supported")
		***REMOVED***

		ds := scheduler.NewVariableLoopingVUsConfig(lib.DefaultSchedulerName)
		ds.StartVUs = conf.VUs
		for _, s := range conf.Stages ***REMOVED***
			if s.Duration.Valid ***REMOVED***
				ds.Stages = append(ds.Stages, scheduler.Stage***REMOVED***Duration: s.Duration, Target: s.Target***REMOVED***)
			***REMOVED***
		***REMOVED***
		result.Execution = scheduler.ConfigMap***REMOVED***lib.DefaultSchedulerName: ds***REMOVED***
	***REMOVED*** else if conf.Iterations.Valid || conf.Execution == nil ***REMOVED***
		// Either shared iterations were explicitly specified via the shortcut option, or no execution
		// parameters were specified in any way, which will run the default 1 iteration in 1 VU
		if conf.Iterations.Valid && conf.Execution != nil ***REMOVED***
			return conf, errors.New("specifying both iterations and execution is not supported")
		***REMOVED***

		ds := scheduler.NewSharedIterationsConfig(lib.DefaultSchedulerName)
		ds.VUs = conf.VUs
		if conf.Iterations.Valid ***REMOVED*** // TODO: fix where the default iterations value is set... sigh...
			ds.Iterations = conf.Iterations
		***REMOVED***

		result.Execution = scheduler.ConfigMap***REMOVED***lib.DefaultSchedulerName: ds***REMOVED***
	***REMOVED***

	//TODO: validate the config; questions:
	// - separately validate the duration, iterations and stages for better error messages?
	// - or reuse the execution validation somehow, at the end? or something mixed?
	// - here or in getConsolidatedConfig() or somewhere else?

	return result, nil
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

	return buildExecutionConfig(conf)
***REMOVED***
