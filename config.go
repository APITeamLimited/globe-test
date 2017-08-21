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
	"encoding/json"
	"os"

	"github.com/ghodss/yaml"
	"github.com/shibukawa/configdir"
	log "github.com/Sirupsen/logrus"
)

const configFilename = "config.yml"

type ConfigCollectors map[string]interface***REMOVED******REMOVED***

func (env *ConfigCollectors) UnmarshalJSON(data []byte) error ***REMOVED***
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil ***REMOVED***
		return err
	***REMOVED***

	res := make(map[string]interface***REMOVED******REMOVED***, len(raw))
	for t, data := range raw ***REMOVED***
		c := collectorOfType(t)
		if c == nil ***REMOVED***
			log.Debugf("Config for unknown collector '%s'; skipping...", t)
			continue
		***REMOVED***

		cconf := c.MakeConfig()
		if err := json.Unmarshal(data, cconf); err != nil ***REMOVED***
			return err
		***REMOVED***
		res[t] = cconf
	***REMOVED***

	*env = res
	return nil
***REMOVED***

func (env ConfigCollectors) Get(t string) interface***REMOVED******REMOVED*** ***REMOVED***
	conf, ok := env[t]
	if !ok ***REMOVED***
		conf = collectorOfType(t).MakeConfig()
	***REMOVED***
	return conf
***REMOVED***

// Global application configuration.
type Config struct ***REMOVED***
	// Collector-specific data placeholders.
	Collectors ConfigCollectors `json:"collectors,omitempty"`
***REMOVED***

func LoadConfig() (conf Config, err error) ***REMOVED***
	conf.Collectors = make(ConfigCollectors)

	cdir := configdir.New("loadimpact", "k6")
	folders := cdir.QueryFolders(configdir.Global)
	data, err := folders[0].ReadFile(configFilename)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return conf, nil
		***REMOVED***
		return conf, err
	***REMOVED***

	if err := yaml.Unmarshal(data, &conf); err != nil ***REMOVED***
		return conf, err
	***REMOVED***
	return conf, nil
***REMOVED***

func (c Config) Store() error ***REMOVED***
	data, err := yaml.Marshal(c)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cdir := configdir.New("loadimpact", "k6")
	folders := cdir.QueryFolders(configdir.Global)
	return folders[0].WriteFile(configFilename, data)
***REMOVED***
