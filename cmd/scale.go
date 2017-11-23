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

	"github.com/loadimpact/k6/api/v1"
	"github.com/loadimpact/k6/api/v1/client"
	"github.com/loadimpact/k6/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command***REMOVED***
	Use:   "scale",
	Short: "Scale a running test",
	Long: `Scale a running test.

  Use the global --address flag to specify the URL to the API server.`,
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		vus := getNullInt64(cmd.Flags(), "vus")
		max := getNullInt64(cmd.Flags(), "max")
		if !vus.Valid && !max.Valid ***REMOVED***
			return errors.New("Specify either -u/--vus or -m/--max")
		***REMOVED***

		c, err := client.New(address)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		status, err := c.SetStatus(context.Background(), v1.Status***REMOVED***VUs: vus, VUsMax: max***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ui.Dump(stdout, status)
		return nil
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(scaleCmd)

	scaleCmd.Flags().Int64P("vus", "u", 1, "number of virtual users")
	scaleCmd.Flags().Int64P("max", "m", 0, "max available virtual users")
***REMOVED***