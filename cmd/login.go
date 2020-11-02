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
	"github.com/spf13/cobra"
)

func getLoginCmd() *cobra.Command ***REMOVED***
	// loginCmd represents the login command
	loginCmd := &cobra.Command***REMOVED***
		Use:   "login",
		Short: "Authenticate with a service",
		Long: `Authenticate with a service.

Logging into a service changes the default when just "-o [type]" is passed with
no parameters, you can always override the stored credentials by passing some
on the commandline.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			return cmd.Usage()
		***REMOVED***,
	***REMOVED***
	return loginCmd
***REMOVED***

func init() ***REMOVED***
	loginCmd := getLoginCmd()
	RootCmd.AddCommand(loginCmd)

	loginCloudCommand := getLoginCloudCommand()
	loginCmd.AddCommand(loginCloudCommand)
	loginCloudCommand.Flags().StringP("token", "t", "", "specify `token` to use")
	loginCloudCommand.Flags().BoolP("show", "s", false, "display saved token and exit")
	loginCloudCommand.Flags().BoolP("reset", "r", false, "reset token")

	loginInfluxDBCommand := getLoginInfluxDBCommand()
	loginCmd.AddCommand(loginInfluxDBCommand)
***REMOVED***
