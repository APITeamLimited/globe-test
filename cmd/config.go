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

	"github.com/kelseyhightower/envconfig"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/shibukawa/configdir"
	"github.com/spf13/pflag"
	null "gopkg.in/guregu/null.v3"
)

const configFilename = "config.json"

var (
	configDirs    = configdir.New("loadimpact", "k6")
	configFlagSet = pflag.NewFlagSet("", 0)
)

func init() ***REMOVED***
	configFlagSet.SortFlags = false
	configFlagSet.StringP("out", "o", "", "`uri` for an external metrics database")
	configFlagSet.BoolP("linger", "l", false, "keep the API server alive past test end")
	configFlagSet.Bool("no-usage-report", false, "don't send anonymous stats to the developers")
***REMOVED***

type Config struct ***REMOVED***
	lib.Options

	Out           null.String `json:"out" envconfig:"out"`
	Linger        null.Bool   `json:"linger" envconfig:"linger"`
	NoUsageReport null.Bool   `json:"noUsageReport" envconfig:"no_usage_report"`

	Collectors struct ***REMOVED***
		InfluxDB influxdb.Config `json:"influxdb"`
		Cloud    cloud.Config    `json:"cloud"`
	***REMOVED*** `json:"collectors"`
***REMOVED***

func (c Config) Apply(cfg Config) Config ***REMOVED***
	c.Options = c.Options.Apply(cfg.Options)
	if cfg.Out.Valid ***REMOVED***
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
	return c
***REMOVED***

// Gets configuration from CLI flags.
func getConfig(flags *pflag.FlagSet) (Config, error) ***REMOVED***
	opts, err := getOptions(flags)
	if err != nil ***REMOVED***
		return Config***REMOVED******REMOVED***, err
	***REMOVED***
	return Config***REMOVED***
		Options:       opts,
		Out:           getNullString(flags, "out"),
		Linger:        getNullBool(flags, "linger"),
		NoUsageReport: getNullBool(flags, "no-usage-report"),
	***REMOVED***, nil
***REMOVED***

// Reads a configuration file from disk.
func readDiskConfig() (Config, *configdir.Config, error) ***REMOVED***
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
func writeDiskConfig(cdir *configdir.Config, conf Config) error ***REMOVED***
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return cdir.WriteFile(configFilename, data)
***REMOVED***

// Reads configuration variables from the environment.
func readEnvConfig() (conf Config, err error) ***REMOVED***
	err = envconfig.Process("k6", &conf)
	return conf, err
***REMOVED***
