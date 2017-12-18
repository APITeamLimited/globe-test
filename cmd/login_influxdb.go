package cmd

import (
	"os"
	"time"

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

		conf := config.Collectors.InfluxDB
		if len(args) > 0 ***REMOVED***
			if err := conf.UnmarshalText([]byte(args[0])); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if conf.Addr == "" ***REMOVED***
			conf.Addr = "http://localhost:8086"
		***REMOVED***
		if conf.DB == "" ***REMOVED***
			conf.DB = "k6"
		***REMOVED***

		form := ui.Form***REMOVED***
			Fields: []ui.Field***REMOVED***
				ui.StringField***REMOVED***
					Key:     "Addr",
					Label:   "Address",
					Default: conf.Addr,
				***REMOVED***,
				ui.StringField***REMOVED***
					Key:     "DB",
					Label:   "Database",
					Default: conf.DB,
				***REMOVED***,
				ui.StringField***REMOVED***
					Key:     "Username",
					Label:   "Username",
					Default: conf.Username,
				***REMOVED***,
				ui.StringField***REMOVED***
					Key:     "Password",
					Label:   "Password",
					Default: conf.Password,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		vals, err := form.Run(os.Stdin, stdout)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := mapstructure.Decode(vals, &conf); err != nil ***REMOVED***
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
