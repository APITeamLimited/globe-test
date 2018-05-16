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
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

type testCmdData struct ***REMOVED***
	Name  string
	Tests []testCmdTest
***REMOVED***

type testCmdTest struct ***REMOVED***
	Args     []string
	Expected []string
	Name     string
***REMOVED***

func TestConfigCmd(t *testing.T) ***REMOVED***

	testdata := []testCmdData***REMOVED***
		***REMOVED***
			Name: "Out",

			Tests: []testCmdTest***REMOVED***
				***REMOVED***
					Name:     "NoArgs",
					Args:     []string***REMOVED***""***REMOVED***,
					Expected: []string***REMOVED******REMOVED***,
				***REMOVED***,
				***REMOVED***
					Name:     "SingleArg",
					Args:     []string***REMOVED***"--out", "influxdb=http://localhost:8086/k6"***REMOVED***,
					Expected: []string***REMOVED***"influxdb=http://localhost:8086/k6"***REMOVED***,
				***REMOVED***,
				***REMOVED***
					Name:     "MultiArg",
					Args:     []string***REMOVED***"--out", "influxdb=http://localhost:8086/k6", "--out", "json=test.json"***REMOVED***,
					Expected: []string***REMOVED***"influxdb=http://localhost:8086/k6", "json=test.json"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, data := range testdata ***REMOVED***
		t.Run(data.Name, func(t *testing.T) ***REMOVED***
			for _, test := range data.Tests ***REMOVED***
				t.Run(`"`+test.Name+`"`, func(t *testing.T) ***REMOVED***
					fs := configFlagSet()
					fs.AddFlagSet(optionFlagSet())
					fs.Parse(test.Args)

					config, err := getConfig(fs)
					assert.NoError(t, err)
					assert.Equal(t, test.Expected, config.Out)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConfigEnv(t *testing.T) ***REMOVED***
	testdata := map[struct***REMOVED*** Name, Key string ***REMOVED***]map[string]func(Config)***REMOVED***
		***REMOVED***"Linger", "K6_LINGER"***REMOVED***: ***REMOVED***
			"":      func(c Config) ***REMOVED*** assert.Equal(t, null.Bool***REMOVED******REMOVED***, c.Linger) ***REMOVED***,
			"true":  func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(true), c.Linger) ***REMOVED***,
			"false": func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(false), c.Linger) ***REMOVED***,
		***REMOVED***,
		***REMOVED***"NoUsageReport", "K6_NO_USAGE_REPORT"***REMOVED***: ***REMOVED***
			"":      func(c Config) ***REMOVED*** assert.Equal(t, null.Bool***REMOVED******REMOVED***, c.NoUsageReport) ***REMOVED***,
			"true":  func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(true), c.NoUsageReport) ***REMOVED***,
			"false": func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(false), c.NoUsageReport) ***REMOVED***,
		***REMOVED***,
		***REMOVED***"Out", "K6_OUT"***REMOVED***: ***REMOVED***
			"":         func(c Config) ***REMOVED*** assert.Equal(t, []string***REMOVED***""***REMOVED***, c.Out) ***REMOVED***,
			"influxdb": func(c Config) ***REMOVED*** assert.Equal(t, []string***REMOVED***"influxdb"***REMOVED***, c.Out) ***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for field, data := range testdata ***REMOVED***
		os.Clearenv()
		t.Run(field.Name, func(t *testing.T) ***REMOVED***
			for value, fn := range data ***REMOVED***
				t.Run(`"`+value+`"`, func(t *testing.T) ***REMOVED***
					assert.NoError(t, os.Setenv(field.Key, value))
					var config Config
					assert.NoError(t, envconfig.Process("k6", &config))
					fn(config)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConfigApply(t *testing.T) ***REMOVED***
	t.Run("Linger", func(t *testing.T) ***REMOVED***
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Linger: null.BoolFrom(true)***REMOVED***)
		assert.Equal(t, null.BoolFrom(true), conf.Linger)
	***REMOVED***)
	t.Run("NoUsageReport", func(t *testing.T) ***REMOVED***
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***NoUsageReport: null.BoolFrom(true)***REMOVED***)
		assert.Equal(t, null.BoolFrom(true), conf.NoUsageReport)
	***REMOVED***)
	t.Run("Out", func(t *testing.T) ***REMOVED***
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Out: []string***REMOVED***"influxdb"***REMOVED******REMOVED***)
		assert.Equal(t, []string***REMOVED***"influxdb"***REMOVED***, conf.Out)

		conf = Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Out: []string***REMOVED***"influxdb", "json"***REMOVED******REMOVED***)
		assert.Equal(t, []string***REMOVED***"influxdb", "json"***REMOVED***, conf.Out)
	***REMOVED***)
***REMOVED***
