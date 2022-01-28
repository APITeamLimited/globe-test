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
	"context"

	"github.com/spf13/cobra"
	"gopkg.in/guregu/null.v3"

	v1 "go.k6.io/k6/api/v1"
	"go.k6.io/k6/api/v1/client"
)

func getPauseCmd(ctx context.Context, globalFlags *commandFlags) *cobra.Command ***REMOVED***
	// pauseCmd represents the pause command
	pauseCmd := &cobra.Command***REMOVED***
		Use:   "pause",
		Short: "Pause a running test",
		Long: `Pause a running test.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			c, err := client.New(globalFlags.address)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			status, err := c.SetStatus(ctx, v1.Status***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			return yamlPrint(globalFlags.stdout, status)
		***REMOVED***,
	***REMOVED***
	return pauseCmd
***REMOVED***
