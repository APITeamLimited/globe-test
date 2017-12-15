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
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var archiveOut = "archive.tar"

// archiveCmd represents the pause command
var archiveCmd = &cobra.Command***REMOVED***
	Use:   "archive",
	Short: "Create an archive",
	Long: `Create an archive.
	 
An archive is a fully self-contained test run, and can be executed identically elsewhere.`,
	Example: `
  # Archive a test run.
  k6 archive -u 10 -d 10s -O myarchive.tar script.js
  
  # Run the resulting archive.
  k6 run myarchive.tar`[1:],
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		// Runner.
		pwd, err := os.Getwd()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		filename := args[0]
		src, err := readSource(filename, pwd, afero.NewOsFs(), os.Stdin)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r, err := newRunner(src, runType, afero.NewOsFs())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Options.
		cliOpts, err := getOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		fileConf, _, err := readDiskConfig()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		envConf, err := readEnvConfig()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		opts := cliOpts.Apply(fileConf.Options).Apply(r.GetOptions()).Apply(envConf.Options).Apply(cliOpts)
		r.SetOptions(opts)

		// Archive.
		arc := r.MakeArchive()
		f, err := os.Create(archiveOut)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return arc.Write(f)
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().AddFlagSet(optionFlagSet())
	archiveCmd.Flags().StringVarP(&archiveOut, "archive-out", "O", archiveOut, "archive output filename")
***REMOVED***