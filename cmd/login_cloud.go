package cmd

import (
	"fmt"
	"os"

	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/ui"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// loginCloudCommand represents the 'login cloud' command
var loginCloudCommand = &cobra.Command***REMOVED***
	Use:   "cloud",
	Short: "Authenticate with Load Impact",
	Long: `Authenticate with Load Impact.

This will set the default token used when just "k6 run -o cloud" is passed.`,
	Example: `
  # Show the stored token.
  k6 login cloud -s

  # Store a token.
  k6 login cloud -t YOUR_TOKEN

  # Log in with an email/password.
  k6 login cloud`[1:],
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		fs := afero.NewOsFs()
		config, cdir, err := readDiskConfig(fs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		show := getNullBool(cmd.Flags(), "show")
		token := getNullString(cmd.Flags(), "token")

		conf := config.Collectors.Cloud

		switch ***REMOVED***
		case show.Bool:
		case token.Valid:
			conf.Token = token.String
		default:
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
			res, err := client.Login(email, password)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if res.Token == "" ***REMOVED***
				return errors.New(`Your account has no API token, please generate one: "https://app.loadimpact.com/account/token".`)
			***REMOVED***

			conf.Token = res.Token
		***REMOVED***

		config.Collectors.Cloud = conf
		if err := writeDiskConfig(fs, cdir, config); err != nil ***REMOVED***
			return err
		***REMOVED***

		fmt.Fprintf(stdout, "  token: %s\n", ui.ValueColor.Sprint(conf.Token))
		return nil
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	loginCmd.AddCommand(loginCloudCommand)
	loginCloudCommand.Flags().StringP("token", "t", "", "specify `token` to use")
	loginCloudCommand.Flags().BoolP("show", "s", false, "display saved token and exit")
***REMOVED***
