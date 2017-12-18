package cmd

import (
	"os"

	"github.com/loadimpact/k6/ui"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// loginCloudCommand represents the 'login cloud' command
var loginCloudCommand = &cobra.Command***REMOVED***
	Use:   "cloud",
	Short: "Authenticate with Load Impact",
	Long: `Authenticate with Load Impact.

This will set the default server used when just "-o cloud" is passed.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		fs := afero.NewOsFs()
		config, cdir, err := readDiskConfig(fs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		conf := config.Collectors.Cloud
		form := ui.Form***REMOVED***
			Fields: []ui.Field***REMOVED***
				ui.StringField***REMOVED***
					Key:     "token",
					Label:   "API Token",
					Default: conf.Token,
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

		config.Collectors.Cloud = conf
		return writeDiskConfig(fs, cdir, config)
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	loginCmd.AddCommand(loginCloudCommand)
***REMOVED***
