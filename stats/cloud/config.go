/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package cloud

import (
	"encoding/json"
)

type ConfigFields struct ***REMOVED***
	Token     string `json:"token" mapstructure:"token" envconfig:"CLOUD_TOKEN"`
	Name      string `json:"name" mapstructure:"name" envconfig:"CLOUD_NAME"`
	Host      string `json:"host" mapstructure:"host" envconfig:"CLOUD_HOST"`
	ProjectID int    `json:"project_id" mapstructure:"project_id" envconfig:"CLOUD_PROJECT_ID"`
***REMOVED***

type Config ConfigFields

func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.Token != "" ***REMOVED***
		c.Token = cfg.Token
	***REMOVED***
	if cfg.Name != "" ***REMOVED***
		c.Name = cfg.Name
	***REMOVED***
	if cfg.Host != "" ***REMOVED***
		c.Host = cfg.Host
	***REMOVED***
	if cfg.ProjectID != 0 ***REMOVED***
		c.ProjectID = cfg.ProjectID
	***REMOVED***
	return c
***REMOVED***

func (c *Config) UnmarshalText(data []byte) error ***REMOVED***
	if s := string(data); s != "" ***REMOVED***
		c.Name = s
	***REMOVED***
	return nil
***REMOVED***

func (c *Config) UnmarshalJSON(data []byte) error ***REMOVED***
	fields := ConfigFields(*c)
	if err := json.Unmarshal(data, &fields); err != nil ***REMOVED***
		return err
	***REMOVED***
	*c = Config(fields)
	return nil
***REMOVED***

func (c Config) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(ConfigFields(c))
***REMOVED***
