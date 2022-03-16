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

	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/lib/metrics"
)

func getArchiveCmd(globalState *globalState) *cobra.Command ***REMOVED*** // nolint: funlen
	archiveOut := "archive.tar"
	// archiveCmd represents the archive command
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
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			src, filesystems, err := readSource(globalState, args[0])
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			runtimeOptions, err := getRuntimeOptions(cmd.Flags(), globalState.envVars)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			registry := metrics.NewRegistry()
			builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
			r, err := newRunner(
				globalState.logger, src, globalState.flags.runType,
				filesystems, runtimeOptions, builtinMetrics, registry,
			)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			cliOpts, err := getOptions(cmd.Flags())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			conf, err := getConsolidatedConfig(globalState, Config***REMOVED***Options: cliOpts***REMOVED***, r.GetOptions())
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			// Parse the thresholds, only if the --no-threshold flag is not set.
			// If parsing the threshold expressions failed, consider it as an
			// invalid configuration error.
			if !runtimeOptions.NoThresholds.Bool ***REMOVED***
				for _, thresholds := range conf.Options.Thresholds ***REMOVED***
					err = thresholds.Parse()
					if err != nil ***REMOVED***
						return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
					***REMOVED***
				***REMOVED***
			***REMOVED***

			_, err = deriveAndValidateConfig(conf, r.IsExecutable, globalState.logger)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			err = r.SetOptions(conf.Options)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			// Archive.
			arc := r.MakeArchive()
			f, err := globalState.fs.Create(archiveOut)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			err = arc.Write(f)
			if cerr := f.Close(); err == nil && cerr != nil ***REMOVED***
				err = cerr
			***REMOVED***
			return err
		***REMOVED***,
	***REMOVED***

	archiveCmd.Flags().SortFlags = false
	archiveCmd.Flags().AddFlagSet(archiveCmdFlagSet(&archiveOut))

	return archiveCmd
***REMOVED***

func archiveCmdFlagSet(archiveOut *string) *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(false))
	flags.StringVarP(archiveOut, "archive-out", "O", *archiveOut, "archive output filename")
	return flags
***REMOVED***
