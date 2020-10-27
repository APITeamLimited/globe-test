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
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/loader"
)

func getInspectCmd() *cobra.Command ***REMOVED***
	// inspectCmd represents the resume command
	inspectCmd := &cobra.Command***REMOVED***
		Use:   "inspect [file]",
		Short: "Inspect a script or archive",
		Long:  `Inspect a script or archive.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			// TODO: don't use the Global logger
			logger := logrus.StandardLogger()
			pwd, err := os.Getwd()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			filesystems := loader.CreateFilesystems()
			src, err := loader.ReadSource(logger, args[0], pwd, filesystems, os.Stdin)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			typ := runType
			if typ == "" ***REMOVED***
				typ = detectType(src.Data)
			***REMOVED***

			runtimeOptions, err := getRuntimeOptions(cmd.Flags(), buildEnvMap(os.Environ()))
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			var (
				opts lib.Options
				b    *js.Bundle
			)
			switch typ ***REMOVED***
			case typeArchive:
				var arc *lib.Archive
				arc, err = lib.ReadArchive(bytes.NewBuffer(src.Data))
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				b, err = js.NewBundleFromArchive(logger, arc, runtimeOptions)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				opts = b.Options
			case typeJS:
				b, err = js.NewBundle(logger, src, filesystems, runtimeOptions)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				opts = b.Options
			***REMOVED***

			data, err := json.MarshalIndent(opts, "", "  ")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			fmt.Println(string(data))
			return nil
		***REMOVED***,
	***REMOVED***
	return inspectCmd
***REMOVED***

func init() ***REMOVED***
	inspectCmd := getInspectCmd()
	RootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().SortFlags = false
	inspectCmd.Flags().AddFlagSet(runtimeOptionFlagSet(false))
	inspectCmd.Flags().StringVarP(&runType, "type", "t", runType, "override file `type`, \"js\" or \"archive\"")
***REMOVED***
