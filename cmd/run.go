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
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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

	thresholdHaveFailedErrorCode = 99
	setupTimeoutErrorCode        = 100
	teardownTimeoutErrorCode     = 101
	genericTimeoutErrorCode      = 102
	genericEngineErrorCode       = 103
	invalidConfigErrorCode       = 104
	externalAbortErrorCode       = 105
	cannotStartRESTAPIErrorCode  = 106
)

// TODO: fix this, global variables are not very testable...
//nolint:gochecknoglobals
var runType = os.Getenv("K6_TYPE")

//nolint:funlen,gocognit,gocyclo
func getRunCmd(ctx context.Context, logger *logrus.Logger) *cobra.Command ***REMOVED***
	// runCmd represents the run command.
	runCmd := &cobra.Command***REMOVED***
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
			// TODO: disable in quiet mode?
			_, _ = BannerColor.Fprintf(stdout, "\n%s\n\n", consts.Banner())

			logger.Debug("Initializing the runner...")

			// Create the Runner.
			pwd, err := os.Getwd()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			filename := args[0]
			filesystems := loader.CreateFilesystems()
			src, err := loader.ReadSource(logger, filename, pwd, filesystems, os.Stdin)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			runtimeOptions, err := getRuntimeOptions(cmd.Flags(), buildEnvMap(os.Environ()))
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			r, err := newRunner(logger, src, runType, filesystems, runtimeOptions)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			logger.Debug("Getting the script options...")

			cliConf, err := getConfig(cmd.Flags())
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			conf, err := getConsolidatedConfig(afero.NewOsFs(), cliConf, r)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			conf, cerr := deriveAndValidateConfig(conf, r.IsExecutable)
			if cerr != nil ***REMOVED***
				return ExitCode***REMOVED***error: cerr, Code: invalidConfigErrorCode***REMOVED***
			***REMOVED***

			// Write options back to the runner too.
			if err = r.SetOptions(conf.Options); err != nil ***REMOVED***
				return err
			***REMOVED***

			// We prepare a bunch of contexts:
			//  - The runCtx is cancelled as soon as the Engine's run() lambda finishes,
			//    and can trigger things like the usage report and end of test summary.
			//    Crucially, metrics processing by the Engine will still work after this
			//    context is cancelled!
			//  - The lingerCtx is cancelled by Ctrl+C, and is used to wait for that
			//    event when k6 was ran with the --linger option.
			//  - The globalCtx is cancelled only after we're completely done with the
			//    test execution and any --linger has been cleared, so that the Engine
			//    can start winding down its metrics processing.
			globalCtx, globalCancel := context.WithCancel(ctx)
			defer globalCancel()
			lingerCtx, lingerCancel := context.WithCancel(globalCtx)
			defer lingerCancel()
			runCtx, runCancel := context.WithCancel(lingerCtx)
			defer runCancel()

			// Create a local execution scheduler wrapping the runner.
			logger.Debug("Initializing the execution scheduler...")
			execScheduler, err := local.NewExecutionScheduler(r, logger)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			executionState := execScheduler.GetState()

			// This is manually triggered after the Engine's Run() has completed,
			// and things like a single Ctrl+C don't affect it. We use it to make
			// sure that the progressbars finish updating with the latest execution
			// state one last time, after the test run has finished.
			progressCtx, progressCancel := context.WithCancel(globalCtx)
			defer progressCancel()
			initBar := execScheduler.GetInitProgressBar()
			progressBarWG := &sync.WaitGroup***REMOVED******REMOVED***
			progressBarWG.Add(1)
			go func() ***REMOVED***
				pbs := []*pb.ProgressBar***REMOVED***execScheduler.GetInitProgressBar()***REMOVED***
				for _, s := range execScheduler.GetExecutors() ***REMOVED***
					pbs = append(pbs, s.GetProgress())
				***REMOVED***
				showProgress(progressCtx, conf, pbs, logger)
				progressBarWG.Done()
			***REMOVED***()

			// Create an engine.
			initBar.Modify(pb.WithConstProgress(0, "Init engine"))
			engine, err := core.NewEngine(execScheduler, conf.Options, runtimeOptions, logger)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			executionPlan := execScheduler.GetExecutionPlan()
			// Create a collector and assign it to the engine if requested.
			initBar.Modify(pb.WithConstProgress(0, "Init metric outputs"))
			for _, out := range conf.Out ***REMOVED***
				t, arg := parseCollector(out)
				collector, cerr := newCollector(logger, t, arg, src, conf, executionPlan)
				if cerr != nil ***REMOVED***
					return cerr
				***REMOVED***
				if cerr = collector.Init(); cerr != nil ***REMOVED***
					return cerr
				***REMOVED***
				engine.Collectors = append(engine.Collectors, collector)
			***REMOVED***

			// Spin up the REST API server, if not disabled.
			if address != "" ***REMOVED***
				initBar.Modify(pb.WithConstProgress(0, "Init API server"))
				go func() ***REMOVED***
					logger.Debugf("Starting the REST API server on %s", address)
					if aerr := api.ListenAndServe(address, engine, logger); aerr != nil ***REMOVED***
						// Only exit k6 if the user has explicitly set the REST API address
						if cmd.Flags().Lookup("address").Changed ***REMOVED***
							logger.WithError(aerr).Error("Error from API server")
							os.Exit(cannotStartRESTAPIErrorCode)
						***REMOVED*** else ***REMOVED***
							logger.WithError(aerr).Warn("Error from API server")
						***REMOVED***
					***REMOVED***
				***REMOVED***()
			***REMOVED***

			printExecutionDescription(
				"local", filename, "", conf, execScheduler.GetState().ExecutionTuple,
				executionPlan, engine.Collectors)

			// Trap Interrupts, SIGINTs and SIGTERMs.
			sigC := make(chan os.Signal, 1)
			signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer signal.Stop(sigC)
			go func() ***REMOVED***
				sig := <-sigC
				logger.WithField("sig", sig).Debug("Stopping k6 in response to signal...")
				lingerCancel() // stop the test run, metric processing is cancelled below

				// If we get a second signal, we immediately exit, so something like
				// https://github.com/loadimpact/k6/issues/971 never happens again
				sig = <-sigC
				logger.WithField("sig", sig).Error("Aborting k6 in response to signal")
				globalCancel() // not that it matters, given the following command...
				os.Exit(externalAbortErrorCode)
			***REMOVED***()

			// Initialize the engine
			initBar.Modify(pb.WithConstProgress(0, "Init VUs..."))
			engineRun, engineWait, err := engine.Init(globalCtx, runCtx)
			if err != nil ***REMOVED***
				return getExitCodeFromEngine(err)
			***REMOVED***

			// Init has passed successfully, so unless disabled, make sure we send a
			// usage report after the context is done.
			if !conf.NoUsageReport.Bool ***REMOVED***
				reportDone := make(chan struct***REMOVED******REMOVED***)
				go func() ***REMOVED***
					<-runCtx.Done()
					_ = reportUsage(execScheduler)
					close(reportDone)
				***REMOVED***()
				defer func() ***REMOVED***
					select ***REMOVED***
					case <-reportDone:
					case <-time.After(3 * time.Second):
					***REMOVED***
				***REMOVED***()
			***REMOVED***

			// Start the test run
			initBar.Modify(pb.WithConstProgress(0, "Starting test..."))
			if err := engineRun(); err != nil ***REMOVED***
				return getExitCodeFromEngine(err)
			***REMOVED***
			runCancel()
			logger.Debug("Engine run terminated cleanly")

			progressCancel()
			progressBarWG.Wait()

			// Warn if no iterations could be completed.
			if executionState.GetFullIterationCount() == 0 ***REMOVED***
				logger.Warn("No script iterations finished, consider making the test duration longer")
			***REMOVED***

			data := ui.SummaryData***REMOVED***
				Metrics:   engine.Metrics,
				RootGroup: engine.ExecutionScheduler.GetRunner().GetDefaultGroup(),
				Time:      executionState.GetCurrentTestRunDuration(),
				TimeUnit:  conf.Options.SummaryTimeUnit.String,
			***REMOVED***
			// Print the end-of-test summary.
			if !runtimeOptions.NoSummary.Bool ***REMOVED***
				fprintf(stdout, "\n")

				s := ui.NewSummary(conf.SummaryTrendStats)
				s.SummarizeMetrics(stdout, "", data)

				fprintf(stdout, "\n")
			***REMOVED***

			if runtimeOptions.SummaryExport.ValueOrZero() != "" ***REMOVED*** //nolint:nestif
				f, err := os.Create(runtimeOptions.SummaryExport.String)
				if err != nil ***REMOVED***
					logger.WithError(err).Error("failed to create summary export file")
				***REMOVED*** else ***REMOVED***
					defer func() ***REMOVED***
						if err := f.Close(); err != nil ***REMOVED***
							logger.WithError(err).Error("failed to close summary export file")
						***REMOVED***
					***REMOVED***()
					s := ui.NewSummary(conf.SummaryTrendStats)
					if err := s.SummarizeMetricsJSON(f, data); err != nil ***REMOVED***
						logger.WithError(err).Error("failed to make summary export file")
					***REMOVED***
				***REMOVED***
			***REMOVED***

			if conf.Linger.Bool ***REMOVED***
				select ***REMOVED***
				case <-lingerCtx.Done():
					// do nothing, we were interrupted by Ctrl+C already
				default:
					logger.Debug("Linger set; waiting for Ctrl+C...")
					fprintf(stdout, "Linger set; waiting for Ctrl+C...")
					<-lingerCtx.Done()
					logger.Debug("Ctrl+C received, exiting...")
				***REMOVED***
			***REMOVED***
			globalCancel() // signal the Engine that it should wind down
			logger.Debug("Waiting for engine processes to finish...")
			engineWait()
			logger.Debug("Everything has finished, exiting k6!")
			if engine.IsTainted() ***REMOVED***
				return ExitCode***REMOVED***error: errors.New("some thresholds have failed"), Code: thresholdHaveFailedErrorCode***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***,
	***REMOVED***

	runCmd.Flags().SortFlags = false
	runCmd.Flags().AddFlagSet(runCmdFlagSet())

	return runCmd
***REMOVED***

func getExitCodeFromEngine(err error) ExitCode ***REMOVED***
	switch e := errors.Cause(err).(type) ***REMOVED***
	case lib.TimeoutError:
		switch e.Place() ***REMOVED***
		case consts.SetupFn:
			return ExitCode***REMOVED***error: err, Code: setupTimeoutErrorCode, Hint: e.Hint()***REMOVED***
		case consts.TeardownFn:
			return ExitCode***REMOVED***error: err, Code: teardownTimeoutErrorCode, Hint: e.Hint()***REMOVED***
		default:
			return ExitCode***REMOVED***error: err, Code: genericTimeoutErrorCode***REMOVED***
		***REMOVED***
	default:
		//nolint:golint
		return ExitCode***REMOVED***error: errors.New("Engine error"), Code: genericEngineErrorCode, Hint: err.Error()***REMOVED***
	***REMOVED***
***REMOVED***

func reportUsage(execScheduler *local.ExecutionScheduler) error ***REMOVED***
	execState := execScheduler.GetState()
	executorConfigs := execScheduler.GetExecutorConfigs()

	executors := make(map[string]int)
	for _, ec := range executorConfigs ***REMOVED***
		executors[ec.GetType()]++
	***REMOVED***

	body, err := json.Marshal(map[string]interface***REMOVED******REMOVED******REMOVED***
		"k6_version": consts.Version,
		"executors":  executors,
		"vus_max":    execState.GetInitializedVUsCount(),
		"iterations": execState.GetFullIterationCount(),
		"duration":   execState.GetCurrentTestRunDuration().String(),
		"goos":       runtime.GOOS,
		"goarch":     runtime.GOARCH,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	res, err := http.Post("https://reports.k6.io/", "application/json", bytes.NewBuffer(body))
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			_ = res.Body.Close()
		***REMOVED***
	***REMOVED***()

	return err
***REMOVED***

func runCmdFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(true))
	flags.AddFlagSet(configFlagSet())

	// TODO: Figure out a better way to handle the CLI flags:
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

// Creates a new runner.
func newRunner(
	logger *logrus.Logger, src *loader.SourceData, typ string, filesystems map[string]afero.Fs, rtOpts lib.RuntimeOptions,
) (lib.Runner, error) ***REMOVED***
	switch typ ***REMOVED***
	case "":
		return newRunner(logger, src, detectType(src.Data), filesystems, rtOpts)
	case typeJS:
		return js.New(logger, src, filesystems, rtOpts)
	case typeArchive:
		arc, err := lib.ReadArchive(bytes.NewReader(src.Data))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch arc.Type ***REMOVED***
		case typeJS:
			return js.NewFromArchive(logger, arc, rtOpts)
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
