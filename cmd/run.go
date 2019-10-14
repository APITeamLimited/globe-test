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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/loadimpact/k6/api"
	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/ui"
	"github.com/loadimpact/k6/ui/pb"
)

const (
	typeJS      = "js"
	typeArchive = "archive"

	thresholdHaveFailedErroCode = 99
	setupTimeoutErrorCode       = 100
	teardownTimeoutErrorCode    = 101
	genericTimeoutErrorCode     = 102
	genericEngineErrorCode      = 103
	invalidConfigErrorCode      = 104
)

//TODO: fix this, global variables are not very testable...
//nolint:gochecknoglobals
var runType = os.Getenv("K6_TYPE")

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
  k6 run -u 0 -s 10s:100 -s 60s -s 10s:0

  # Send metrics to an influxdb server
  k6 run -o influxdb=http://1.2.3.4:8086/k6`[1:],
	Args: exactArgsWithMsg(1, "arg should either be \"-\", if reading script from stdin, or a path to a script file"),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		//TODO: disable in quiet mode?
		_, _ = BannerColor.Fprintf(stdout, "\n%s\n\n", consts.Banner)

		initBar := pb.New(pb.WithConstLeft("   init"))

		// Create the Runner.
		fprintf(stdout, "%s runner\r", initBar.String()) //TODO use printBar()
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

		runtimeOptions, err := getRuntimeOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r, err := newRunner(src, runType, filesystems, runtimeOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		fprintf(stdout, "%s options\r", initBar.String())

		cliConf, err := getConfig(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		conf, err := getConsolidatedConfig(afero.NewOsFs(), cliConf, r)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		conf, cerr := deriveAndValidateConfig(conf)
		if cerr != nil ***REMOVED***
			return ExitCode***REMOVED***cerr, invalidConfigErrorCode***REMOVED***
		***REMOVED***

		// If summary trend stats are defined, update the UI to reflect them
		if len(conf.SummaryTrendStats) > 0 ***REMOVED***
			ui.UpdateTrendColumns(conf.SummaryTrendStats)
		***REMOVED***

		// Write options back to the runner too.
		if err = r.SetOptions(conf.Options); err != nil ***REMOVED***
			return err
		***REMOVED***

		//TODO: don't use a global... or maybe change the logger?
		logger := logrus.StandardLogger()

		ctx, cancel := context.WithCancel(context.Background()) //TODO: move even earlier?
		defer cancel()

		// Create a local execution scheduler wrapping the runner.
		fprintf(stdout, "%s execution scheduler\r", initBar.String())
		execScheduler, err := local.NewExecutionScheduler(r, logger)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		executionState := execScheduler.GetState()
		initBar = execScheduler.GetInitProgressBar()
		progressBarWG := &sync.WaitGroup***REMOVED******REMOVED***
		progressBarWG.Add(1)
		go func() ***REMOVED***
			showProgress(ctx, conf, execScheduler)
			progressBarWG.Done()
		***REMOVED***()

		// Create an engine.
		initBar.Modify(pb.WithConstProgress(0, "Init engine"))
		engine, err := core.NewEngine(execScheduler, conf.Options, logger)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		//TODO: the engine should just probably have a copy of the config...
		// Configure the engine.
		if conf.NoThresholds.Valid ***REMOVED***
			engine.NoThresholds = conf.NoThresholds.Bool
		***REMOVED***
		if conf.NoSummary.Valid ***REMOVED***
			engine.NoSummary = conf.NoSummary.Bool
		***REMOVED***

		// Create a collector and assign it to the engine if requested.
		initBar.Modify(pb.WithConstProgress(0, "Init metric outputs"))
		for _, out := range conf.Out ***REMOVED***
			t, arg := parseCollector(out)
			collector, err := newCollector(t, arg, src, conf, execScheduler.GetExecutionPlan())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := collector.Init(); err != nil ***REMOVED***
				return err
			***REMOVED***
			engine.Collectors = append(engine.Collectors, collector)
		***REMOVED***

		// Create an API server.
		if address != "" ***REMOVED***
			initBar.Modify(pb.WithConstProgress(0, "Init API server"))
			go func() ***REMOVED***
				if err := api.ListenAndServe(address, engine); err != nil ***REMOVED***
					logger.WithError(err).Warn("Error from API server")
				***REMOVED***
			***REMOVED***()
		***REMOVED***

		// Write the big banner.
		***REMOVED***
			out := "-"
			link := ""

			for idx, collector := range engine.Collectors ***REMOVED***
				if out != "-" ***REMOVED***
					out = out + "; " + conf.Out[idx]
				***REMOVED*** else ***REMOVED***
					out = conf.Out[idx]
				***REMOVED***

				if l := collector.Link(); l != "" ***REMOVED***
					link = link + " (" + l + ")"
				***REMOVED***
			***REMOVED***

			fprintf(stdout, "   executor: %s\n", ui.ValueColor.Sprint("local"))
			fprintf(stdout, "     output: %s%s\n", ui.ValueColor.Sprint(out), ui.ExtraColor.Sprint(link))
			fprintf(stdout, "     script: %s\n", ui.ValueColor.Sprint(filename))
			fprintf(stdout, "\n")

			plan := execScheduler.GetExecutionPlan()
			executors := execScheduler.GetExecutors()
			maxDuration, _ := lib.GetEndOffset(plan)

			fprintf(stdout, "  execution: %s\n", ui.ValueColor.Sprintf(
				"(%.2f%%) %d executors, %d max VUs, %s max duration (incl. graceful stop):",
				conf.ExecutionSegment.FloatLength()*100, len(executors),
				lib.GetMaxPossibleVUs(plan), maxDuration),
			)
			for _, sched := range executors ***REMOVED***
				fprintf(stdout, "           * %s: %s\n",
					sched.GetConfig().GetName(), sched.GetConfig().GetDescription(conf.ExecutionSegment))
			***REMOVED***
			fprintf(stdout, "\n")
		***REMOVED***

		// Run the engine with a cancellable context.
		errC := make(chan error)
		go func() ***REMOVED***
			initBar.Modify(pb.WithConstProgress(0, "Init VUs"))
			if err := engine.Init(ctx); err != nil ***REMOVED***
				errC <- err
			***REMOVED*** else ***REMOVED***
				initBar.Modify(pb.WithConstProgress(0, "Start test"))
				errC <- engine.Run(ctx)
			***REMOVED***
		***REMOVED***()

		// Trap Interrupts, SIGINTs and SIGTERMs.
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigC)

		// If the user hasn't opted out: report usage.
		//TODO: fix
		//TODO: move to a separate function
		/*
			if !conf.NoUsageReport.Bool ***REMOVED***
				go func() ***REMOVED***
					u := "http://k6reports.loadimpact.com/"
					mime := "application/json"
					var endTSeconds float64
					if endT := engine.Executor.GetEndTime(); endT.Valid ***REMOVED***
						endTSeconds = time.Duration(endT.Duration).Seconds()
					***REMOVED***
					var stagesEndTSeconds float64
					if stagesEndT := lib.SumStages(engine.Executor.GetStages()); stagesEndT.Valid ***REMOVED***
						stagesEndTSeconds = time.Duration(stagesEndT.Duration).Seconds()
					***REMOVED***
					body, err := json.Marshal(map[string]interface***REMOVED******REMOVED******REMOVED***
						"k6_version":  Version,
						"vus_max":     engine.Executor.GetVUsMax(),
						"iterations":  engine.Executor.GetEndIterations(),
						"duration":    endTSeconds,
						"st_duration": stagesEndTSeconds,
						"goos":        runtime.GOOS,
						"goarch":      runtime.GOARCH,
					***REMOVED***)
					if err != nil ***REMOVED***
						panic(err) // This should never happen!!
					***REMOVED***
					_, _ = http.Post(u, mime, bytes.NewBuffer(body))
				***REMOVED***()
			***REMOVED***
		*/

		// Ticker for progress bar updates. Less frequent updates for non-TTYs, none if quiet.
		updateFreq := 50 * time.Millisecond
		if !stdoutTTY ***REMOVED***
			updateFreq = 1 * time.Second
		***REMOVED***
		ticker := time.NewTicker(updateFreq)
		if quiet || conf.HttpDebug.Valid && conf.HttpDebug.String != "" ***REMOVED***
			ticker.Stop()
		***REMOVED***
	mainLoop:
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				if quiet || !stdoutTTY ***REMOVED***
					l := logrus.WithFields(logrus.Fields***REMOVED***
						"t": executionState.GetCurrentTestRunDuration(),
						"i": executionState.GetFullIterationCount(),
					***REMOVED***)
					fn := l.Info
					if quiet ***REMOVED***
						fn = l.Debug
					***REMOVED***
					if executionState.IsPaused() ***REMOVED***
						fn("Paused")
					***REMOVED*** else ***REMOVED***
						fn("Running")
					***REMOVED***
				***REMOVED***
			case err := <-errC:
				cancel()
				if err == nil ***REMOVED***
					logger.Debug("Engine terminated cleanly")
					break mainLoop
				***REMOVED***

				switch e := errors.Cause(err).(type) ***REMOVED***
				case lib.TimeoutError:
					switch e.Place() ***REMOVED***
					case "setup":
						logger.WithField("hint", e.Hint()).Error(err)
						return ExitCode***REMOVED***errors.New("Setup timeout"), setupTimeoutErrorCode***REMOVED***
					case "teardown":
						logger.WithField("hint", e.Hint()).Error(err)
						return ExitCode***REMOVED***errors.New("Teardown timeout"), teardownTimeoutErrorCode***REMOVED***
					default:
						logger.WithError(err).Error("Engine timeout")
						return ExitCode***REMOVED***errors.New("Engine timeout"), genericTimeoutErrorCode***REMOVED***
					***REMOVED***
				default:
					logger.WithError(err).Error("Engine error")
					return ExitCode***REMOVED***errors.New("Engine Error"), genericEngineErrorCode***REMOVED***
				***REMOVED***
			case sig := <-sigC:
				logger.WithField("sig", sig).Debug("Exiting in response to signal")
				cancel()
				//TODO: Actually exit on a second Ctrl+C, even if some of the iterations are stuck.
				// This is currently problematic because of https://github.com/loadimpact/k6/issues/971,
				// but with uninterruptible iterations it will be even more problematic.
			***REMOVED***
		***REMOVED***
		if quiet || !stdoutTTY ***REMOVED***
			e := logger.WithFields(logrus.Fields***REMOVED***
				"t": executionState.GetCurrentTestRunDuration(),
				"i": executionState.GetFullIterationCount(),
			***REMOVED***)
			fn := e.Info
			if quiet ***REMOVED***
				fn = e.Debug
			***REMOVED***
			fn("Test finished")
		***REMOVED***

		progressBarWG.Wait()

		// Warn if no iterations could be completed.
		if executionState.GetFullIterationCount() == 0 ***REMOVED***
			logger.Warn("No data generated, because no script iterations finished, consider making the test duration longer")
		***REMOVED***

		// Print the end-of-test summary.
		if !conf.NoSummary.Bool ***REMOVED***
			fprintf(stdout, "\n")
			ui.Summarize(stdout, "", ui.SummaryData***REMOVED***
				Opts:    conf.Options,
				Root:    engine.ExecutionScheduler.GetRunner().GetDefaultGroup(),
				Metrics: engine.Metrics,
				Time:    executionState.GetCurrentTestRunDuration(),
			***REMOVED***)
			fprintf(stdout, "\n")
		***REMOVED***

		if conf.Linger.Bool ***REMOVED***
			logger.Info("Linger set; waiting for Ctrl+C...")
			<-sigC
		***REMOVED***

		if engine.IsTainted() ***REMOVED***
			return ExitCode***REMOVED***errors.New("some thresholds have failed"), thresholdHaveFailedErroCode***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***,
***REMOVED***

func runCmdFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(true))
	flags.AddFlagSet(configFlagSet())

	//TODO: Figure out a better way to handle the CLI flags:
	// - the default values are specified in this way so we don't overwrire whatever
	//   was specified via the environment variables
	// - but we need to manually specify the DefValue, since that's the default value
	//   that will be used in the help/usage message - if we don't set it, the environment
	//   variables will affect the usage message
	// - and finally, global variables are not very testable... :/
	flags.StringVarP(&runType, "type", "t", runType, "override file `type`, \"js\" or \"archive\"")
	flags.Lookup("type").DefValue = ""
	return flags
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().SortFlags = false
	runCmd.Flags().AddFlagSet(runCmdFlagSet())
***REMOVED***

// Creates a new runner.
func newRunner(
	src *loader.SourceData, typ string, filesystems map[string]afero.Fs, rtOpts lib.RuntimeOptions,
) (lib.Runner, error) ***REMOVED***
	switch typ ***REMOVED***
	case "":
		return newRunner(src, detectType(src.Data), filesystems, rtOpts)
	case typeJS:
		return js.New(src, filesystems, rtOpts)
	case typeArchive:
		arc, err := lib.ReadArchive(bytes.NewReader(src.Data))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch arc.Type ***REMOVED***
		case typeJS:
			return js.NewFromArchive(arc, rtOpts)
		default:
			return nil, errors.Errorf("archive requests unsupported runner: %s", arc.Type)
		***REMOVED***
	default:
		return nil, errors.Errorf("unknown -t/--type: %s", typ)
	***REMOVED***
***REMOVED***

func detectType(data []byte) string ***REMOVED***
	if _, err := tar.NewReader(bytes.NewReader(data)).Next(); err == nil ***REMOVED***
		return typeArchive
	***REMOVED***
	return typeJS
***REMOVED***
