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
	"github.com/spf13/pflag"

	"github.com/loadimpact/k6/loader"
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
		filesystems := loader.CreateFilesystems()
		src, err := loader.ReadSource(filename, pwd, filesystems, os.Stdin)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		runtimeOptions, err := getRuntimeOptions(cmd.Flags(), buildEnvMap(os.Environ()))
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r, err := newRunner(src, runType, filesystems, runtimeOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		cliOpts, err := getOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		conf, err := getConsolidatedConfig(afero.NewOsFs(), Config***REMOVED***Options: cliOpts***REMOVED***, r)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if _, cerr := deriveAndValidateConfig(conf, r.GetExports()); cerr != nil ***REMOVED***
			return ExitCode***REMOVED***error: cerr, Code: invalidConfigErrorCode***REMOVED***
		***REMOVED***

		err = r.SetOptions(conf.Options)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Archive.
		arc := r.MakeArchive()
		f, err := os.Create(archiveOut)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return arc.Write(f)
	***REMOVED***,
***REMOVED***

func archiveCmdFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(false))
	//TODO: figure out a better way to handle the CLI flags - global variables are not very testable... :/
	flags.StringVarP(&archiveOut, "archive-out", "O", archiveOut, "archive output filename")
	return flags
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().SortFlags = false
	archiveCmd.Flags().AddFlagSet(archiveCmdFlagSet())
***REMOVED***
