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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"go.k6.io/k6/api"
	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/output"
	"go.k6.io/k6/ui/pb"
)

// cmdRun handles the `k6 run` sub-command
type cmdRun struct ***REMOVED***
	gs *globalState
***REMOVED***

// TODO: split apart some more
//nolint:funlen,gocognit,gocyclo,cyclop
func (c *cmdRun) run(cmd *cobra.Command, args []string) error ***REMOVED***
	printBanner(c.gs)

	test, err := loadAndConfigureTest(c.gs, cmd, args, getConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Write the full consolidated *and derived* options back to the Runner.
	conf := test.derivedConfig
	if err = test.initRunner.SetOptions(conf.Options); err != nil ***REMOVED***
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
	globalCtx, globalCancel := context.WithCancel(c.gs.ctx)
	defer globalCancel()
	lingerCtx, lingerCancel := context.WithCancel(globalCtx)
	defer lingerCancel()
	runCtx, runCancel := context.WithCancel(lingerCtx)
	defer runCancel()

	logger := c.gs.logger
	// Create a local execution scheduler wrapping the runner.
	logger.Debug("Initializing the execution scheduler...")
	execScheduler, err := local.NewExecutionScheduler(test.initRunner, test.builtInMetrics, logger)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

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
		showProgress(progressCtx, c.gs, pbs, logger)
		progressBarWG.Done()
	***REMOVED***()

	// Create all outputs.
	executionPlan := execScheduler.GetExecutionPlan()
	outputs, err := createOutputs(c.gs, test, executionPlan)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// TODO: remove
	// Create the engine.
	initBar.Modify(pb.WithConstProgress(0, "Init engine"))
	engine, err := core.NewEngine(
		execScheduler, conf.Options, test.runtimeOptions,
		outputs, logger, test.metricsRegistry,
	)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Spin up the REST API server, if not disabled.
	if c.gs.flags.address != "" ***REMOVED***
		initBar.Modify(pb.WithConstProgress(0, "Init API server"))
		go func() ***REMOVED***
			logger.Debugf("Starting the REST API server on %s", c.gs.flags.address)
			// TODO: send the ExecutionState and MetricsEngine instead of the Engine
			if aerr := api.ListenAndServe(c.gs.flags.address, engine, logger); aerr != nil ***REMOVED***
				// Only exit k6 if the user has explicitly set the REST API address
				if cmd.Flags().Lookup("address").Changed ***REMOVED***
					logger.WithError(aerr).Error("Error from API server")
					c.gs.osExit(int(exitcodes.CannotStartRESTAPI))
				***REMOVED*** else ***REMOVED***
					logger.WithError(aerr).Warn("Error from API server")
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// We do this here so we can get any output URLs below.
	initBar.Modify(pb.WithConstProgress(0, "Starting outputs"))
	outputManager := output.NewManager(outputs, logger, func(err error) ***REMOVED***
		if err != nil ***REMOVED***
			logger.WithError(err).Error("Received error to stop from output")
		***REMOVED***
		runCancel()
	***REMOVED***)
	err = outputManager.StartOutputs()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer outputManager.StopOutputs()

	printExecutionDescription(
		c.gs, "local", args[0], "", conf, execScheduler.GetState().ExecutionTuple, executionPlan, outputs,
	)

	// Trap Interrupts, SIGINTs and SIGTERMs.
	gracefulStop := func(sig os.Signal) ***REMOVED***
		logger.WithField("sig", sig).Debug("Stopping k6 in response to signal...")
		lingerCancel() // stop the test run, metric processing is cancelled below
	***REMOVED***
	hardStop := func(sig os.Signal) ***REMOVED***
		logger.WithField("sig", sig).Error("Aborting k6 in response to signal")
		globalCancel() // not that it matters, given the following command...
	***REMOVED***
	stopSignalHandling := handleTestAbortSignals(c.gs, gracefulStop, hardStop)
	defer stopSignalHandling()

	// Initialize the engine
	initBar.Modify(pb.WithConstProgress(0, "Init VUs..."))
	engineRun, engineWait, err := engine.Init(globalCtx, runCtx)
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		return errext.WithExitCodeIfNone(err, exitcodes.GenericEngine)
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
	var interrupt error
	err = engineRun()
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		if common.IsInterruptError(err) ***REMOVED***
			// Don't return here since we need to work with --linger,
			// show the end-of-test summary and exit cleanly.
			interrupt = err
		***REMOVED***
		if !conf.Linger.Bool && interrupt == nil ***REMOVED***
			return errext.WithExitCodeIfNone(err, exitcodes.GenericEngine)
		***REMOVED***
	***REMOVED***
	runCancel()
	logger.Debug("Engine run terminated cleanly")

	progressCancel()
	progressBarWG.Wait()

	executionState := execScheduler.GetState()
	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 ***REMOVED***
		logger.Warn("No script iterations finished, consider making the test duration longer")
	***REMOVED***

	// Handle the end-of-test summary.
	if !test.runtimeOptions.NoSummary.Bool ***REMOVED***
		summaryResult, err := test.initRunner.HandleSummary(globalCtx, &lib.Summary***REMOVED***
			Metrics:         engine.Metrics,
			RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
			TestRunDuration: executionState.GetCurrentTestRunDuration(),
			NoColor:         c.gs.flags.noColor,
			UIState: lib.UIState***REMOVED***
				IsStdOutTTY: c.gs.stdOut.isTTY,
				IsStdErrTTY: c.gs.stdErr.isTTY,
			***REMOVED***,
		***REMOVED***)
		if err == nil ***REMOVED***
			err = handleSummaryResult(c.gs.fs, c.gs.stdOut, c.gs.stdErr, summaryResult)
		***REMOVED***
		if err != nil ***REMOVED***
			logger.WithError(err).Error("failed to handle the end-of-test summary")
		***REMOVED***
	***REMOVED***

	if conf.Linger.Bool ***REMOVED***
		select ***REMOVED***
		case <-lingerCtx.Done():
			// do nothing, we were interrupted by Ctrl+C already
		default:
			logger.Debug("Linger set; waiting for Ctrl+C...")
			if !c.gs.flags.quiet ***REMOVED***
				printToStdout(c.gs, "Linger set; waiting for Ctrl+C...")
			***REMOVED***
			<-lingerCtx.Done()
			logger.Debug("Ctrl+C received, exiting...")
		***REMOVED***
	***REMOVED***
	globalCancel() // signal the Engine that it should wind down
	logger.Debug("Waiting for engine processes to finish...")
	engineWait()
	logger.Debug("Everything has finished, exiting k6!")
	if interrupt != nil ***REMOVED***
		return interrupt
	***REMOVED***
	if engine.IsTainted() ***REMOVED***
		return errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed)
	***REMOVED***
	return nil
***REMOVED***

func (c *cmdRun) flagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(true))
	flags.AddFlagSet(configFlagSet())
	return flags
***REMOVED***

func getCmdRun(gs *globalState) *cobra.Command ***REMOVED***
	c := &cmdRun***REMOVED***
		gs: gs,
	***REMOVED***

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
		RunE: c.run,
	***REMOVED***

	runCmd.Flags().SortFlags = false
	runCmd.Flags().AddFlagSet(c.flagSet())

	return runCmd
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
	res, err := http.Post("https://reports.k6.io/", "application/json", bytes.NewBuffer(body)) //nolint:noctx
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			_ = res.Body.Close()
		***REMOVED***
	***REMOVED***()

	return err
***REMOVED***

func handleSummaryResult(fs afero.Fs, stdOut, stdErr io.Writer, result map[string]io.Reader) error ***REMOVED***
	var errs []error

	getWriter := func(path string) (io.Writer, error) ***REMOVED***
		switch path ***REMOVED***
		case "stdout":
			return stdOut, nil
		case "stderr":
			return stdErr, nil
		default:
			return fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		***REMOVED***
	***REMOVED***

	for path, value := range result ***REMOVED***
		if writer, err := getWriter(path); err != nil ***REMOVED***
			errs = append(errs, fmt.Errorf("could not open '%s': %w", path, err))
		***REMOVED*** else if n, err := io.Copy(writer, value); err != nil ***REMOVED***
			errs = append(errs, fmt.Errorf("error saving summary to '%s' after %d bytes: %w", path, n, err))
		***REMOVED***
	***REMOVED***

	return consolidateErrorMessage(errs, "Could not save some summary information:")
***REMOVED***
