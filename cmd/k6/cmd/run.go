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
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/loadimpact/k6/ui"

	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/loader"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	null "gopkg.in/guregu/null.v3"
)

const (
	typeJS      = "js"
	typeArchive = "archive"
)

var (
	runType       string
	linger        bool
	noUsageReport bool
)

// runCmd represents the run command.
var runCmd = &cobra.Command***REMOVED***
	Use:   "run",
	Short: "Start a load test",
	Long: `Start a load test.

This also exposes a REST API to interact with it. Various k6 subcommands offer
a commandline interface for interacting with it.`,
	Example: `
  # Run a single VU, once.
  k6 run script.js

  # Run a single VU, 10 times.
  k6 run -i 10 script.js

  # Run 5 VUs, splitting 10 iterations between them.
  k6 run -u 5 -i 10 script.js

  # Run 5 VUs for 10s.
  k6 run -u 5 -d 10s script.js

  # Ramp VUs from 0 to 100 over 10s, stay there for 60s, then 10s down to 0.
  k6 run -u 0 -s 10s:100 -s 60s -s 10s:0`[1:],
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		// Create the Runner.
		pwd, err := os.Getwd()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		src, err := readSource(args[0], pwd, afero.NewOsFs(), os.Stdin)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r, err := newRunner(src, runType, afero.NewOsFs())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Assemble options; start with the CLI-provided options to get shadowed (non-Valid)
		// defaults in there, override with Runner-provided ones, then merge the CLI opts in
		// on top to give them priority.
		cliOpts, err := getOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		opts := cliOpts.Apply(r.GetOptions()).Apply(cliOpts)

		// If -m/--max isn't specified, figure out the max that should be needed.
		if !opts.VUsMax.Valid ***REMOVED***
			opts.VUsMax = null.IntFrom(opts.VUs.Int64)
			for _, stage := range opts.Stages ***REMOVED***
				if stage.Target.Valid && stage.Target.Int64 > opts.VUsMax.Int64 ***REMOVED***
					opts.VUsMax = stage.Target
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// If -d/--duration, -i/--iterations and -s/--stage are all unset, run to one iteration.
		if !opts.Duration.Valid && !opts.Iterations.Valid && opts.Stages == nil ***REMOVED***
			opts.Iterations = null.IntFrom(1)
		***REMOVED***

		// Write options back to the runner too.
		r.ApplyOptions(opts)

		// Create an engine with a local executor, wrapping the Runner.
		engine, err := core.NewEngine(local.New(r), opts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Run the engine with a cancellable context.
		ctx, cancel := context.WithCancel(context.Background())
		errC := make(chan error)
		go func() ***REMOVED*** errC <- engine.Run(ctx) ***REMOVED***()

		// Trap Interrupts, SIGINTs and SIGTERMs.
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigC)

		// Prepare a progress bar.
		progress := ui.ProgressBar***REMOVED***
			Width: 60,
			Left: func() string ***REMOVED***
				if engine.Executor.IsPaused() ***REMOVED***
					return "  paused"
				***REMOVED*** else if engine.Executor.IsRunning() ***REMOVED***
					return " running"
				***REMOVED*** else ***REMOVED***
					return "    done"
				***REMOVED***
			***REMOVED***,
			Right: func() string ***REMOVED***
				if endIt := engine.Executor.GetEndIterations(); endIt.Valid ***REMOVED***
					return fmt.Sprintf("%d / %d", engine.Executor.GetIterations(), endIt.Int64)
				***REMOVED***
				precision := 100 * time.Millisecond
				atT := engine.Executor.GetTime()
				if endT := engine.Executor.GetEndTime(); endT.Valid ***REMOVED***
					return fmt.Sprintf("%s / %s",
						(atT/precision)*precision,
						(time.Duration(endT.Duration)/precision)*precision,
					)
				***REMOVED***
				return ((atT / precision) * precision).String()
			***REMOVED***,
		***REMOVED***

		// Ticker for progress bar updates. Less frequent updates for non-TTYs, none if quiet.
		updateFreq := 50 * time.Millisecond
		if !stdoutTTY ***REMOVED***
			updateFreq = 1 * time.Second
		***REMOVED***
		ticker := time.NewTicker(updateFreq)
		if quiet ***REMOVED***
			ticker.Stop()
		***REMOVED***
	mainLoop:
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				if quiet || !stdoutTTY ***REMOVED***
					l := log.WithFields(log.Fields***REMOVED***
						"t": engine.Executor.GetTime(),
						"i": engine.Executor.GetIterations(),
					***REMOVED***)
					fn := l.Info
					if quiet ***REMOVED***
						fn = l.Debug
					***REMOVED***
					if engine.Executor.IsPaused() ***REMOVED***
						fn("Paused")
					***REMOVED*** else ***REMOVED***
						fn("Running")
					***REMOVED***
					break
				***REMOVED***

				var prog float64
				if endIt := engine.Executor.GetEndIterations(); endIt.Valid ***REMOVED***
					prog = float64(engine.Executor.GetIterations()) / float64(endIt.Int64)
				***REMOVED*** else if endT := engine.Executor.GetEndTime(); endT.Valid ***REMOVED***
					prog = float64(engine.Executor.GetTime()) / float64(endT.Duration)
				***REMOVED***
				progress.Progress = prog
				fmt.Fprintf(stdout, "%s\x1b[0K\r", progress.String())
			case err := <-errC:
				if err != nil ***REMOVED***
					log.WithError(err).Error("Engine error")
				***REMOVED*** else ***REMOVED***
					log.Debug("Engine terminated cleanly")
				***REMOVED***
				cancel()
				break mainLoop
			case sig := <-sigC:
				log.WithField("sig", sig).Debug("Exiting in response to signal")
				cancel()
			***REMOVED***
		***REMOVED***
		if quiet || !stdoutTTY ***REMOVED***
			e := log.WithFields(log.Fields***REMOVED***
				"t": engine.Executor.GetTime(),
				"i": engine.Executor.GetIterations(),
			***REMOVED***)
			fn := e.Info
			if quiet ***REMOVED***
				fn = e.Debug
			***REMOVED***
			fn("Test finished")
		***REMOVED*** else ***REMOVED***
			progress.Progress = 1
			fmt.Fprintf(stdout, "%s\x1b[0K\n", progress.String())
		***REMOVED***

		// Print the end-of-test summary.
		if !quiet ***REMOVED***
			fmt.Fprintf(stdout, "\n")
			ui.Summarize(stdout, stdoutTTY, "", ui.SummaryData***REMOVED***
				Opts:    opts,
				Root:    engine.Executor.GetRunner().GetDefaultGroup(),
				Metrics: engine.Metrics,
			***REMOVED***)
		***REMOVED***

		return nil
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().SortFlags = false
	registerOptions(runCmd.Flags())

	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.StringVarP(&runType, "type", "t", "", "override file `type`, \"js\" or \"archive\"")
	flags.BoolVarP(&linger, "linger", "l", false, "keep the API server alive past test end")
	flags.BoolVar(&noUsageReport, "no-usage-report", false, "don't send analytics to the maintainers")
***REMOVED***

// Reads a source file from any supported destination.
func readSource(src, pwd string, fs afero.Fs, stdin io.Reader) (*lib.SourceData, error) ***REMOVED***
	if src == "-" ***REMOVED***
		data, err := ioutil.ReadAll(stdin)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return &lib.SourceData***REMOVED***Filename: "-", Data: data***REMOVED***, nil
	***REMOVED***
	abspath := filepath.Join(pwd, src)
	if ok, _ := afero.Exists(fs, abspath); ok ***REMOVED***
		src = abspath
	***REMOVED***
	return loader.Load(fs, pwd, src)
***REMOVED***

// Creates a new runner.
func newRunner(src *lib.SourceData, typ string, fs afero.Fs) (lib.Runner, error) ***REMOVED***
	switch typ ***REMOVED***
	case "":
		if _, err := tar.NewReader(bytes.NewReader(src.Data)).Next(); err == nil ***REMOVED***
			return newRunner(src, typeArchive, fs)
		***REMOVED***
		return newRunner(src, typeJS, fs)
	case typeJS:
		return js.New(src, fs)
	case typeArchive:
		arc, err := lib.ReadArchive(bytes.NewReader(src.Data))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch arc.Type ***REMOVED***
		case typeJS:
			return js.NewFromArchive(arc)
		default:
			return nil, errors.Errorf("archive requests unsupported runner: %s", arc.Type)
		***REMOVED***
	default:
		return nil, errors.Errorf("unknown -t/--type: %s", typ)
	***REMOVED***
***REMOVED***
