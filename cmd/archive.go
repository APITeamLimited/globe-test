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
	"github.com/spf13/pflag"
)

// cmdArchive handles the `k6 archive` sub-command
type cmdArchive struct ***REMOVED***
	gs *globalState

	archiveOut string
***REMOVED***

func (c *cmdArchive) run(cmd *cobra.Command, args []string) error ***REMOVED***
	test, err := loadAndConfigureTest(c.gs, cmd, args, getPartialConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// It's important to NOT set the derived options back to the runner
	// here, only the consolidated ones. Otherwise, if the script used
	// an execution shortcut option (e.g. `iterations` or `duration`),
	// we will have multiple conflicting execution options since the
	// derivation will set `scenarios` as well.
	testRunState, err := test.buildTestRunState(test.consolidatedConfig.Options)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Archive.
	arc := testRunState.Runner.MakeArchive()
	f, err := c.gs.fs.Create(c.archiveOut)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = arc.Write(f)
	if cerr := f.Close(); err == nil && cerr != nil ***REMOVED***
		err = cerr
	***REMOVED***
	return err
***REMOVED***

func (c *cmdArchive) flagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(false))
	flags.StringVarP(&c.archiveOut, "archive-out", "O", c.archiveOut, "archive output filename")
	return flags
***REMOVED***

func getCmdArchive(gs *globalState) *cobra.Command ***REMOVED***
	c := &cmdArchive***REMOVED***
		gs:         gs,
		archiveOut: "archive.tar",
	***REMOVED***

	archiveCmd := &cobra.Command***REMOVED***
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
		RunE: c.run,
	***REMOVED***

	archiveCmd.Flags().SortFlags = false
	archiveCmd.Flags().AddFlagSet(c.flagSet())

	return archiveCmd
***REMOVED***
