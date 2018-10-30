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
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/ui"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
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
		_, _ = BannerColor.Fprint(stdout, Banner+"\n\n")
		initBar := ui.ProgressBar***REMOVED***
			Width: 60,
			Left:  func() string ***REMOVED*** return "    uploading script" ***REMOVED***,
		***REMOVED***
		fprintf(stdout, "%s \r", initBar.String())

		// Runner
		pwd, err := os.Getwd()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		filename := args[0]
		fs := afero.NewOsFs()
		src, err := readSource(filename, pwd, fs, os.Stdin)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		runtimeOptions, err := getRuntimeOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r, err := newRunner(src, runType, fs, runtimeOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		cliOpts, err := getOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		conf, err := getConsolidatedConfig(fs, Config***REMOVED***Options: cliOpts***REMOVED***, r)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.SetOptions(conf.Options)

		// Cloud config
		cloudConfig := cloud.NewConfig().Apply(conf.Collectors.Cloud)
		if err := envconfig.Process("k6", &cloudConfig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if !cloudConfig.Token.Valid ***REMOVED***
			return errors.New("Not logged in, please use `k6 login cloud`.")
		***REMOVED***

		// Start cloud test run
		client := cloud.NewClient(cloudConfig.Token.String, cloudConfig.Host.String, Version)

		arc := r.MakeArchive()
		if val, ok := arc.Options.External["loadimpact"]; ok ***REMOVED***
			if err := json.Unmarshal(val, &cloudConfig); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if arc.Options.External == nil ***REMOVED***
			arc.Options.External = make(map[string]json.RawMessage)
		***REMOVED***
		arc.Options.External["loadimpact"], err = json.Marshal(cloudConfig)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		name := cloudConfig.Name.String
		if !cloudConfig.Name.Valid || cloudConfig.Name.String == "" ***REMOVED***
			name = filepath.Base(filename)
		***REMOVED***

		if err := client.ValidateOptions(arc.Options); err != nil ***REMOVED***
			return err
		***REMOVED***

		refID, err := client.StartCloudTestRun(name, cloudConfig.ProjectID.Int64, arc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		testURL := cloud.URLForResults(refID, cloudConfig)
		fprintf(stdout, "\n\n")
		fprintf(stdout, "     execution: %s\n", ui.ValueColor.Sprint("cloud"))
		fprintf(stdout, "     script: %s\n", ui.ValueColor.Sprint(filename))
		fprintf(stdout, "     output: %s\n", ui.ValueColor.Sprint(testURL))
		fprintf(stdout, "\n")

		// The quiet option hides the progress bar and disallow aborting the test
		if quiet ***REMOVED***
			return nil
		***REMOVED***

		// Trap Interrupts, SIGINTs and SIGTERMs.
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigC)

		var progressErr error
		testProgress := &cloud.TestProgressResponse***REMOVED******REMOVED***
		progress := ui.ProgressBar***REMOVED***
			Width: 60,
			Left: func() string ***REMOVED***
				return "  " + testProgress.RunStatusText
			***REMOVED***,
		***REMOVED***

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
					progress.Progress = testProgress.Progress
					fprintf(stdout, "%s\x1b[0K\r", progress.String())
				***REMOVED*** else ***REMOVED***
					log.WithError(progressErr).Error("Test progress error")
				***REMOVED***
				if shouldExitLoop ***REMOVED***
					break runningLoop
				***REMOVED***
			case sig := <-sigC:
				log.WithField("sig", sig).Print("Exiting in response to signal...")
				err := client.StopCloudTestRun(refID)
				if err != nil ***REMOVED***
					log.WithError(err).Error("Stop cloud test error")
				***REMOVED***
				shouldExitLoop = true // Exit after the next GetTestProgress call
			***REMOVED***
		***REMOVED***

		if testProgress == nil ***REMOVED***
			return ExitCode***REMOVED***errors.New("Test progress error"), 98***REMOVED***
		***REMOVED***

		fprintf(stdout, "     test status: %s\n", ui.ValueColor.Sprint(testProgress.RunStatusText))

		if testProgress.ResultStatus == cloud.ResultStatusFailed ***REMOVED***
			return ExitCode***REMOVED***errors.New("The test has failed"), 99***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(cloudCmd)
	cloudCmd.Flags().SortFlags = false
	cloudCmd.Flags().AddFlagSet(optionFlagSet())
	cloudCmd.Flags().AddFlagSet(runtimeOptionFlagSet(false))
	cloudCmd.Flags().BoolVar(&exitOnRunning, "exit-on-running", exitOnRunning, "exits when test reaches the running status")
***REMOVED***
