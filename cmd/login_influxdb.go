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

package cmd

import (
	"os"
	"time"

	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats/influxdb"
	"github.com/loadimpact/k6/ui"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// loginInfluxDBCommand represents the 'login influxdb' command
var loginInfluxDBCommand = &cobra.Command***REMOVED***
	Use:   "influxdb [uri]",
	Short: "Authenticate with InfluxDB",
	Long: `Authenticate with InfluxDB.

This will set the default server used when just "-o influxdb" is passed.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		fs := afero.NewOsFs()
		config, cdir, err := readDiskConfig(fs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		conf := influxdb.NewConfig().Apply(config.Collectors.InfluxDB)
		if len(args) > 0 ***REMOVED***
			urlConf, err := influxdb.ParseURL(args[0])
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			conf = conf.Apply(urlConf)
		***REMOVED***

		form := ui.Form***REMOVED***
			Fields: []ui.Field***REMOVED***
				ui.StringField***REMOVED***
					Key:     "Addr",
					Label:   "Address",
					Default: conf.Addr.String,
				***REMOVED***,
				ui.StringField***REMOVED***
					Key:     "DB",
					Label:   "Database",
					Default: conf.DB.String,
				***REMOVED***,
				ui.StringField***REMOVED***
					Key:     "Username",
					Label:   "Username",
					Default: conf.Username.String,
				***REMOVED***,
				ui.PasswordField***REMOVED***
					Key:   "Password",
					Label: "Password",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		vals, err := form.Run(os.Stdin, stdout)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig***REMOVED***
			DecodeHook: types.NullDecoder,
			Result:     &conf,
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := dec.Decode(vals); err != nil ***REMOVED***
			return err
		***REMOVED***

		coll, err := influxdb.New(conf)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, _, err := coll.Client.Ping(10 * time.Second); err != nil ***REMOVED***
			return err
		***REMOVED***

		config.Collectors.InfluxDB = conf
		return writeDiskConfig(fs, cdir, config)
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	loginCmd.AddCommand(loginInfluxDBCommand)
***REMOVED***
