package cmd

import (
	"encoding/json"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/output/influxdb"
	"go.k6.io/k6/ui"
)

//nolint:funlen
func getCmdLoginInfluxDB(globalState *globalState) *cobra.Command ***REMOVED***
	// loginInfluxDBCommand represents the 'login influxdb' command
	loginInfluxDBCommand := &cobra.Command***REMOVED***
		Use:   "influxdb [uri]",
		Short: "Authenticate with InfluxDB",
		Long: `Authenticate with InfluxDB.

This will set the default server used when just "-o influxdb" is passed.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			config, err := readDiskConfig(globalState)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			conf := influxdb.NewConfig()
			jsonConf := config.Collectors["influxdb"]
			if jsonConf != nil ***REMOVED***
				jsonConfParsed, jsonerr := influxdb.ParseJSON(jsonConf)
				if jsonerr != nil ***REMOVED***
					return jsonerr
				***REMOVED***
				conf = conf.Apply(jsonConfParsed)
			***REMOVED***
			if len(args) > 0 ***REMOVED***
				urlConf, err := influxdb.ParseURL(args[0]) //nolint:govet
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
			if !term.IsTerminal(int(syscall.Stdin)) ***REMOVED*** //nolint:unconvert
				globalState.logger.Warn("Stdin is not a terminal, falling back to plain text input")
			***REMOVED***
			vals, err := form.Run(globalState.stdIn, globalState.stdOut)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			conf.Addr = null.StringFrom(vals["Addr"])
			conf.DB = null.StringFrom(vals["DB"])
			conf.Username = null.StringFrom(vals["Username"])
			conf.Password = null.StringFrom(vals["Password"])

			client, err := influxdb.MakeClient(conf)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if _, _, err = client.Ping(10 * time.Second); err != nil ***REMOVED***
				return err
			***REMOVED***

			if config.Collectors == nil ***REMOVED***
				config.Collectors = make(map[string]json.RawMessage)
			***REMOVED***
			config.Collectors["influxdb"], err = json.Marshal(conf)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			return writeDiskConfig(globalState, config)
		***REMOVED***,
	***REMOVED***
	return loginInfluxDBCommand
***REMOVED***
