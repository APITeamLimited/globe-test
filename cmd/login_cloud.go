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
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/ui"
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
		// TODO: don't use a global... or maybe change the logger?
		logger := logrus.StandardLogger()
		fs := afero.NewOsFs()

		k6Conf, err := getConsolidatedConfig(fs, Config***REMOVED******REMOVED***, nil)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		currentDiskConf, configPath, err := readDiskConfig(fs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		show := getNullBool(cmd.Flags(), "show")
		reset := getNullBool(cmd.Flags(), "reset")
		token := getNullString(cmd.Flags(), "token")

		newCloudConf := cloud.NewConfig().Apply(currentDiskConf.Collectors.Cloud)

		switch ***REMOVED***
		case reset.Valid:
			newCloudConf.Token = null.StringFromPtr(nil)
			fprintf(stdout, "  token reset\n")
		case show.Bool:
		case token.Valid:
			newCloudConf.Token = token
		default:
			form := ui.Form***REMOVED***
				Fields: []ui.Field***REMOVED***
					ui.StringField***REMOVED***
						Key:   "Email",
						Label: "Email",
					***REMOVED***,
					ui.PasswordField***REMOVED***
						Key:   "Password",
						Label: "Password",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***
			if !terminal.IsTerminal(syscall.Stdin) ***REMOVED***
				logger.Warn("Stdin is not a terminal, falling back to plain text input")
			***REMOVED***
			vals, err := form.Run(os.Stdin, stdout)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			email := vals["Email"].(string)
			password := vals["Password"].(string)

			client := cloud.NewClient("", k6Conf.Collectors.Cloud.Host.String, consts.Version)
			res, err := client.Login(email, password)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if res.Token == "" ***REMOVED***
				return errors.New(`your account has no API token, please generate one at https://app.k6.io/account/api-token`)
			***REMOVED***

			newCloudConf.Token = null.StringFrom(res.Token)
		***REMOVED***

		currentDiskConf.Collectors.Cloud = newCloudConf
		if err := writeDiskConfig(fs, configPath, currentDiskConf); err != nil ***REMOVED***
			return err
		***REMOVED***

		if newCloudConf.Token.Valid ***REMOVED***
			fprintf(stdout, "  token: %s\n", ui.ValueColor.Sprint(newCloudConf.Token.String))
		***REMOVED***
		return nil
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	loginCmd.AddCommand(loginCloudCommand)
	loginCloudCommand.Flags().StringP("token", "t", "", "specify `token` to use")
	loginCloudCommand.Flags().BoolP("show", "s", false, "display saved token and exit")
	loginCloudCommand.Flags().BoolP("reset", "r", false, "reset token")
***REMOVED***
