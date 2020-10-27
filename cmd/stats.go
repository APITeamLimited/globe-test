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

	"github.com/loadimpact/k6/api/v1/client"
	"github.com/loadimpact/k6/ui"
)

func getStatsCmd() *cobra.Command ***REMOVED***
	// statsCmd represents the stats command
	statsCmd := &cobra.Command***REMOVED***
		Use:   "stats",
		Short: "Show test metrics",
		Long: `Show test metrics.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			c, err := client.New(address)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			metrics, err := c.Metrics(context.Background())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			ui.Dump(stdout, metrics)
			return nil
		***REMOVED***,
	***REMOVED***
	return statsCmd
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(getStatsCmd())
***REMOVED***
