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
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/api"
	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/ui"
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

var (
	//TODO: fix this, global variables are not very testable...
	runType       = os.Getenv("K6_TYPE")
	runNoSetup    = os.Getenv("K6_NO_SETUP") != ""
	runNoTeardown = os.Getenv("K6_NO_TEARDOWN") != ""
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
  k6 run -u 0 -s 10s:100 -s 60s -s 10s:0

  # Send metrics to an influxdb server
  k6 run -o influxdb=http://1.2.3.4:8086/k6`[1:],
	Args: exactArgsWithMsg(1, "arg should either be \"-\", if reading script from stdin, or a path to a script file"),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		//TODO: disable in quiet mode?
		_, _ = BannerColor.Fprintf(stdout, "\n%s\n\n", consts.Banner)

		initBar := ui.ProgressBar***REMOVED***
			Width: 60,
			Left:  func() string ***REMOVED*** return "    init" ***REMOVED***,
		***REMOVED***

		// Create the Runner.
		fprintf(stdout, "%s runner\r", initBar.String())
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

		// If -m/--max isn't specified, figure out the max that should be needed.
		if !conf.VUsMax.Valid ***REMOVED***
			conf.VUsMax = null.NewInt(conf.VUs.Int64, conf.VUs.Valid)
			for _, stage := range conf.Stages ***REMOVED***
				if stage.Target.Valid && stage.Target.Int64 > conf.VUsMax.Int64 ***REMOVED***
					conf.VUsMax = stage.Target
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// If -d/--duration, -i/--iterations and -s/--stage are all unset, run to one iteration.
		if !conf.Duration.Valid && !conf.Iterations.Valid && len(conf.Stages) == 0 ***REMOVED***
			conf.Iterations = null.IntFrom(1)
		***REMOVED***

		if conf.Iterations.Valid && conf.Iterations.Int64 < conf.VUsMax.Int64 ***REMOVED***
			logrus.Warnf(
				"All iterations (%d in this test run) are shared between all VUs, so some of the %d VUs will not execute even a single iteration!",
				conf.Iterations.Int64, conf.VUsMax.Int64,
			)
		***REMOVED***

		//TODO: move a bunch of the logic above to a config "constructor" and to the Validate() method

		// If duration is explicitly set to 0, it means run forever.
		//TODO: just... handle this differently, e.g. as a part of the manual executor
		if conf.Duration.Valid && conf.Duration.Duration == 0 ***REMOVED***
			conf.Duration = types.NullDuration***REMOVED******REMOVED***
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

		// Create a local executor wrapping the runner.
		fprintf(stdout, "%s executor\r", initBar.String())
		ex := local.New(r)
		if runNoSetup ***REMOVED***
			ex.SetRunSetup(false)
		***REMOVED***
		if runNoTeardown ***REMOVED***
			ex.SetRunTeardown(false)
		***REMOVED***

		// Create an engine.
		fprintf(stdout, "%s   engine\r", initBar.String())
		engine, err := core.NewEngine(ex, conf.Options)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Configure the engine.
		if conf.NoThresholds.Valid ***REMOVED***
			engine.NoThresholds = conf.NoThresholds.Bool
		***REMOVED***
		if conf.NoSummary.Valid ***REMOVED***
			engine.NoSummary = conf.NoSummary.Bool
		***REMOVED***

		// Create a collector and assign it to the engine if requested.
		fprintf(stdout, "%s   collector\r", initBar.String())
		for _, out := range conf.Out ***REMOVED***
			t, arg := parseCollector(out)
			collector, err := newCollector(t, arg, src, conf)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := collector.Init(); err != nil ***REMOVED***
				return err
			***REMOVED***
			engine.Collectors = append(engine.Collectors, collector)
		***REMOVED***

		// Create an API server.
		fprintf(stdout, "%s   server\r", initBar.String())
		go func() ***REMOVED***
			if err := api.ListenAndServe(address, engine); err != nil ***REMOVED***
				logrus.WithError(err).Warn("Error from API server")
			***REMOVED***
		***REMOVED***()

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

			fprintf(stdout, "  execution: %s\n", ui.ValueColor.Sprint("local"))
			fprintf(stdout, "     output: %s%s\n", ui.ValueColor.Sprint(out), ui.ExtraColor.Sprint(link))
			fprintf(stdout, "     script: %s\n", ui.ValueColor.Sprint(filename))
			fprintf(stdout, "\n")

			duration := ui.GrayColor.Sprint("-")
			iterations := ui.GrayColor.Sprint("-")
			if conf.Duration.Valid ***REMOVED***
				duration = ui.ValueColor.Sprint(conf.Duration.Duration)
			***REMOVED***
			if conf.Iterations.Valid ***REMOVED***
				iterations = ui.ValueColor.Sprint(conf.Iterations.Int64)
			***REMOVED***
			vus := ui.ValueColor.Sprint(conf.VUs.Int64)
			max := ui.ValueColor.Sprint(conf.VUsMax.Int64)

			leftWidth := ui.StrWidth(duration)
			if l := ui.StrWidth(vus); l > leftWidth ***REMOVED***
				leftWidth = l
			***REMOVED***
			durationPad := strings.Repeat(" ", leftWidth-ui.StrWidth(duration))
			vusPad := strings.Repeat(" ", leftWidth-ui.StrWidth(vus))

			fprintf(stdout, "    duration: %s,%s iterations: %s\n", duration, durationPad, iterations)
			fprintf(stdout, "         vus: %s,%s max: %s\n", vus, vusPad, max)
			fprintf(stdout, "\n")
		***REMOVED***

		// Run the engine with a cancellable context.
		fprintf(stdout, "%s starting\r", initBar.String())
		ctx, cancel := context.WithCancel(context.Background())
		errC := make(chan error)
		go func() ***REMOVED*** errC <- engine.Run(ctx) ***REMOVED***()

		// Trap Interrupts, SIGINTs and SIGTERMs.
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigC)

		// If the user hasn't opted out: report usage.
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
					"k6_version":  consts.Version,
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
				stagesEndT := lib.SumStages(engine.Executor.GetStages())
				endT := engine.Executor.GetEndTime()
				if !endT.Valid || (stagesEndT.Valid && endT.Duration > stagesEndT.Duration) ***REMOVED***
					endT = stagesEndT
				***REMOVED***
				if endT.Valid ***REMOVED***
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
		if quiet || conf.HTTPDebug.Valid && conf.HTTPDebug.String != "" ***REMOVED***
			ticker.Stop()
		***REMOVED***
	mainLoop:
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				if quiet || !stdoutTTY ***REMOVED***
					l := logrus.WithFields(logrus.Fields***REMOVED***
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
				***REMOVED*** else ***REMOVED***
					stagesEndT := lib.SumStages(engine.Executor.GetStages())
					endT := engine.Executor.GetEndTime()
					if !endT.Valid || (stagesEndT.Valid && endT.Duration > stagesEndT.Duration) ***REMOVED***
						endT = stagesEndT
					***REMOVED***
					if endT.Valid ***REMOVED***
						prog = float64(engine.Executor.GetTime()) / float64(endT.Duration)
					***REMOVED***
				***REMOVED***
				progress.Progress = prog
				fprintf(stdout, "%s\x1b[0K\r", progress.String())
			case err := <-errC:
				cancel()
				if err == nil ***REMOVED***
					logrus.Debug("Engine terminated cleanly")
					break mainLoop
				***REMOVED***

				switch e := errors.Cause(err).(type) ***REMOVED***
				case lib.TimeoutError:
					switch e.Place() ***REMOVED***
					case "setup":
						logrus.WithField("hint", e.Hint()).Error(err)
						return ExitCode***REMOVED***errors.New("Setup timeout"), setupTimeoutErrorCode***REMOVED***
					case "teardown":
						logrus.WithField("hint", e.Hint()).Error(err)
						return ExitCode***REMOVED***errors.New("Teardown timeout"), teardownTimeoutErrorCode***REMOVED***
					default:
						logrus.WithError(err).Error("Engine timeout")
						return ExitCode***REMOVED***errors.New("Engine timeout"), genericTimeoutErrorCode***REMOVED***
					***REMOVED***
				default:
					logrus.WithError(err).Error("Engine error")
					return ExitCode***REMOVED***errors.New("Engine Error"), genericEngineErrorCode***REMOVED***
				***REMOVED***
			case sig := <-sigC:
				logrus.WithField("sig", sig).Debug("Exiting in response to signal")
				cancel()
			***REMOVED***
		***REMOVED***
		if quiet || !stdoutTTY ***REMOVED***
			e := logrus.WithFields(logrus.Fields***REMOVED***
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
			fprintf(stdout, "%s\x1b[0K\n", progress.String())
		***REMOVED***

		// Warn if no iterations could be completed.
		if engine.Executor.GetIterations() == 0 ***REMOVED***
			logrus.Warn("No data generated, because no script iterations finished, consider making the test duration longer")
		***REMOVED***

		// Print the end-of-test summary.
		if !conf.NoSummary.Bool ***REMOVED***
			fprintf(stdout, "\n")
			ui.Summarize(stdout, "", ui.SummaryData***REMOVED***
				Opts:    conf.Options,
				Root:    engine.Executor.GetRunner().GetDefaultGroup(),
				Metrics: engine.Metrics,
				Time:    engine.Executor.GetTime(),
			***REMOVED***)
			fprintf(stdout, "\n")
		***REMOVED***

		if conf.Linger.Bool ***REMOVED***
			logrus.Info("Linger set; waiting for Ctrl+C...")
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
	flags.BoolVar(&runNoSetup, "no-setup", runNoSetup, "don't run setup()")
	falseStr := "false" // avoiding goconst warnings...
	flags.Lookup("no-setup").DefValue = falseStr
	flags.BoolVar(&runNoTeardown, "no-teardown", runNoTeardown, "don't run teardown()")
	flags.Lookup("no-teardown").DefValue = falseStr
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
