package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/ui"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	token string
	show  bool
)

func printToken(conf cloud.Config) ***REMOVED***
	label := "Token"
	displayLabel := " " + color.New(color.Faint, color.FgCyan).Sprint("["+label+"]")
	fmt.Fprintf(stdout, " "+displayLabel+": "+conf.Token+"\n")
***REMOVED***

// loginCloudCommand represents the 'login cloud' command
var loginCloudCommand = &cobra.Command***REMOVED***
	Use:   "cloud",
	Short: "Authenticate with Load Impact",
	Long: `Authenticate with Load Impact.

This will set the default Token used when just "k6 run -o cloud" is passed.`,

	Example: `
  # Show the stored token.
  k6 login cloud -s

  # Set up the token.
  k6 login cloud -t YOUR_TOKEN

  # Ask for your Load Impact user email and password to automatically set up the token.
  k6 login cloud`[1:],

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		fs := afero.NewOsFs()
		config, cdir, err := readDiskConfig(fs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		conf := config.Collectors.Cloud

		if show ***REMOVED***
			printToken(conf)
			return nil
		***REMOVED***

		if token != "" ***REMOVED***
			conf.Token = token
		***REMOVED*** else ***REMOVED***
			printToken(conf)

			form := ui.Form***REMOVED***
				Fields: []ui.Field***REMOVED***
					ui.StringField***REMOVED***
						Key:   "Email",
						Label: "Email",
					***REMOVED***,
					ui.StringField***REMOVED***
						Key:   "Password",
						Label: "Password",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***

			vals, err := form.Run(os.Stdin, stdout)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			email := vals["Email"].(string)
			password := vals["Password"].(string)
			client := cloud.NewClient("", conf.Host, Version)
			response, err := client.Login(email, password)
			if err != nil ***REMOVED***
				return errors.New("Failed to login.")
			***REMOVED***

			if response.APIToken == "" ***REMOVED***
				// TODO: instead of `Login`, we must create an endpoint `GetorCreateAPIToken`:
				//       Given an email and password, it will return your Token or create a new one.
				return errors.New("You have to create an API Token with your Load Impact account.")
			***REMOVED***

			conf.Token = response.APIToken
		***REMOVED***

		config.Collectors.Cloud = conf
		if err := writeDiskConfig(cdir, config); err != nil ***REMOVED***
			return err
		***REMOVED***

		printToken(conf)
		return nil

	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	loginCmd.AddCommand(loginCloudCommand)
	loginCloudCommand.Flags().StringVarP(&token, "token", "t", token, "setup the Load Impact Token")
	loginCloudCommand.Flags().BoolVarP(&show, "show", "s", false, "show the saved Load Impact Token")
***REMOVED***
