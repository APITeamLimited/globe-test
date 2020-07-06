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
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/ui"
	"github.com/loadimpact/k6/ui/pb"
)

const (
	cloudFailedToGetProgressErrorCode = 98
	cloudTestRunFailedErrorCode       = 99
)

var (
	exitOnRunning = os.Getenv("K6_EXIT_ON_RUNNING") != ""
)

var cloudCmd = &cobra.Command***REMOVED***
	Use:   "cloud",
	Short: "Run a test on the cloud",
	Long: `Run a test on the cloud.

This will execute the test on the Load Impact cloud service. Use "k6 login cloud" to authenticate.`,
	Example: `
        k6 cloud script.js`[1:],
	Args: exactArgsWithMsg(1, "arg should either be \"-\", if reading script from stdin, or a path to a script file"),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		//TODO: disable in quiet mode?
		_, _ = BannerColor.Fprintf(stdout, "\n%s\n\n", consts.Banner)

		progressBar := pb.New(
			pb.WithConstLeft(" Init"),
			pb.WithConstProgress(0, "Parsing script"),
		)
		printBar(progressBar)

		// Runner
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

		modifyAndPrintBar(progressBar, pb.WithConstProgress(0, "Getting script options"))
		r, err := newRunner(src, runType, filesystems, runtimeOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		modifyAndPrintBar(progressBar, pb.WithConstProgress(0, "Consolidating options"))
		cliOpts, err := getOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		conf, err := getConsolidatedConfig(afero.NewOsFs(), Config***REMOVED***Options: cliOpts***REMOVED***, r)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		derivedConf, cerr := deriveAndValidateConfig(conf, r.IsExecutable)
		if cerr != nil ***REMOVED***
			return ExitCode***REMOVED***error: cerr, Code: invalidConfigErrorCode***REMOVED***
		***REMOVED***

		//TODO: validate for usage of execution segment
		//TODO: validate for externally controlled executor (i.e. executors that aren't distributable)
		//TODO: move those validations to a separate function and reuse validateConfig()?

		err = r.SetOptions(conf.Options)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Cloud config
		cloudConfig := cloud.NewConfig().Apply(derivedConf.Collectors.Cloud)
		if err = envconfig.Process("", &cloudConfig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if !cloudConfig.Token.Valid ***REMOVED***
			return errors.New("Not logged in, please use `k6 login cloud`.")
		***REMOVED***

		modifyAndPrintBar(progressBar, pb.WithConstProgress(0, "Building the archive"))
		arc := r.MakeArchive()
		// TODO: Fix this
		// We reuse cloud.Config for parsing options.ext.loadimpact, but this probably shouldn't be
		// done as the idea of options.ext is that they are extensible without touching k6. But in
		// order for this to happen we shouldn't actually marshall cloud.Config on top of it because
		// it will be missing some fields that aren't actually mentioned in the struct.
		// So in order for use to copy the fields that we need for loadimpact's api we unmarshal in
		// map[string]interface***REMOVED******REMOVED*** and copy what we need if it isn't set already
		var tmpCloudConfig map[string]interface***REMOVED******REMOVED***
		if val, ok := arc.Options.External["loadimpact"]; ok ***REMOVED***
			var dec = json.NewDecoder(bytes.NewReader(val))
			dec.UseNumber() // otherwise float64 are used
			if err := dec.Decode(&tmpCloudConfig); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if err := cloud.MergeFromExternal(arc.Options.External, &cloudConfig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if tmpCloudConfig == nil ***REMOVED***
			tmpCloudConfig = make(map[string]interface***REMOVED******REMOVED***, 3)
		***REMOVED***

		if _, ok := tmpCloudConfig["token"]; !ok && cloudConfig.Token.Valid ***REMOVED***
			tmpCloudConfig["token"] = cloudConfig.Token
		***REMOVED***
		if _, ok := tmpCloudConfig["name"]; !ok && cloudConfig.Name.Valid ***REMOVED***
			tmpCloudConfig["name"] = cloudConfig.Name
		***REMOVED***
		if _, ok := tmpCloudConfig["projectID"]; !ok && cloudConfig.ProjectID.Valid ***REMOVED***
			tmpCloudConfig["projectID"] = cloudConfig.ProjectID
		***REMOVED***

		if arc.Options.External == nil ***REMOVED***
			arc.Options.External = make(map[string]json.RawMessage)
		***REMOVED***
		arc.Options.External["loadimpact"], err = json.Marshal(tmpCloudConfig)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		name := cloudConfig.Name.String
		if !cloudConfig.Name.Valid || cloudConfig.Name.String == "" ***REMOVED***
			name = filepath.Base(filename)
		***REMOVED***

		// Start cloud test run
		modifyAndPrintBar(progressBar, pb.WithConstProgress(0, "Validating script options"))
		client := cloud.NewClient(cloudConfig.Token.String, cloudConfig.Host.String, consts.Version)
		if err := client.ValidateOptions(arc.Options); err != nil ***REMOVED***
			return err
		***REMOVED***

		modifyAndPrintBar(progressBar, pb.WithConstProgress(0, "Uploading archive"))
		refID, err := client.StartCloudTestRun(name, cloudConfig.ProjectID.Int64, arc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		et, err := lib.NewExecutionTuple(derivedConf.ExecutionSegment, derivedConf.ExecutionSegmentSequence)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		testURL := cloud.URLForResults(refID, cloudConfig)
		printExecutionDescription("cloud", filename, testURL, derivedConf, et, nil)

		modifyAndPrintBar(
			progressBar,
			pb.WithConstLeft(" Run "),
			pb.WithConstProgress(0, "Initializing the cloud test"),
		)

		// The quiet option hides the progress bar and disallow aborting the test
		if quiet ***REMOVED***
			return nil
		***REMOVED***

		// Trap Interrupts, SIGINTs and SIGTERMs.
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigC)

		var (
			startTime   time.Time
			maxDuration time.Duration
		)
		executionPlan := derivedConf.Scenarios.GetFullExecutionRequirements(et)
		maxDuration, _ = lib.GetEndOffset(executionPlan)

		testProgress := &cloud.TestProgressResponse***REMOVED******REMOVED***
		progressBar.Modify(
			pb.WithProgress(func() (float64, []string) ***REMOVED***
				statusText := testProgress.RunStatusText

				if testProgress.RunStatus == lib.RunStatusRunning ***REMOVED***
					if startTime.IsZero() ***REMOVED***
						startTime = time.Now()
					***REMOVED***
					spent := time.Since(startTime)
					if spent > maxDuration ***REMOVED***
						statusText = maxDuration.String()
					***REMOVED*** else ***REMOVED***
						statusText = fmt.Sprintf("%s/%s", pb.GetFixedLengthDuration(spent, maxDuration), maxDuration)
					***REMOVED***
				***REMOVED***

				return testProgress.Progress, []string***REMOVED***statusText***REMOVED***
			***REMOVED***),
		)

		var progressErr error
		ticker := time.NewTicker(time.Millisecond * 2000)
		shouldExitLoop := false

	runningLoop:
		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				testProgress, progressErr = client.GetTestProgress(refID)
				if progressErr == nil ***REMOVED***
					if (testProgress.RunStatus > lib.RunStatusRunning) || (exitOnRunning && testProgress.RunStatus == lib.RunStatusRunning) ***REMOVED***
						shouldExitLoop = true
					***REMOVED***
					printBar(progressBar)
				***REMOVED*** else ***REMOVED***
					logrus.WithError(progressErr).Error("Test progress error")
				***REMOVED***
				if shouldExitLoop ***REMOVED***
					break runningLoop
				***REMOVED***
			case sig := <-sigC:
				logrus.WithField("sig", sig).Print("Exiting in response to signal...")
				err := client.StopCloudTestRun(refID)
				if err != nil ***REMOVED***
					logrus.WithError(err).Error("Stop cloud test error")
				***REMOVED***
				shouldExitLoop = true // Exit after the next GetTestProgress call
			***REMOVED***
		***REMOVED***

		if testProgress == nil ***REMOVED***
			//nolint:golint
			return ExitCode***REMOVED***error: errors.New("Test progress error"), Code: cloudFailedToGetProgressErrorCode***REMOVED***
		***REMOVED***

		fprintf(stdout, "     test status: %s\n", ui.ValueColor.Sprint(testProgress.RunStatusText))

		if testProgress.ResultStatus == cloud.ResultStatusFailed ***REMOVED***
			//nolint:golint
			return ExitCode***REMOVED***error: errors.New("The test has failed"), Code: cloudTestRunFailedErrorCode***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***,
***REMOVED***

func cloudCmdFlagSet() *pflag.FlagSet ***REMOVED***
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SortFlags = false
	flags.AddFlagSet(optionFlagSet())
	flags.AddFlagSet(runtimeOptionFlagSet(false))

	//TODO: Figure out a better way to handle the CLI flags:
	// - the default value is specified in this way so we don't overwrire whatever
	//   was specified via the environment variable
	// - global variables are not very testable... :/
	flags.BoolVar(&exitOnRunning, "exit-on-running", exitOnRunning, "exits when test reaches the running status")
	// We also need to explicitly set the default value for the usage message here, so setting
	// K6_EXIT_ON_RUNNING=true won't affect the usage message
	flags.Lookup("exit-on-running").DefValue = "false"

	return flags
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(cloudCmd)
	cloudCmd.Flags().SortFlags = false
	cloudCmd.Flags().AddFlagSet(cloudCmdFlagSet())
***REMOVED***
