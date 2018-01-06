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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/ui"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var cloudCmd = &cobra.Command***REMOVED***
	Use:   "cloud",
	Short: "Run a test on the cloud",
	Long: `Run a test on the cloud.

This will execute the test on the Load Impact cloud service. Use "k6 login cloud" to authenticate.`,
	Example: `
        k6 cloud script.js`[1:],
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
		_, _ = BannerColor.Fprint(stdout, Banner+"\n\n")
		fmt.Fprint(stdout, "  Uploading script to the cloud..")

		// Runner
		pwd, err := os.Getwd()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		filename := args[0]
		src, err := readSource(filename, pwd, afero.NewOsFs(), os.Stdin)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r, err := newRunner(src, runType, afero.NewOsFs())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Options
		fs := afero.NewOsFs()
		fileConf, _, err := readDiskConfig(fs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		options, err := getOptions(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cliConf := Config***REMOVED***Options: options***REMOVED***

		envConf, err := readEnvConfig()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		conf := cliConf.Apply(fileConf).Apply(Config***REMOVED***Options: r.GetOptions()***REMOVED***).Apply(envConf).Apply(cliConf)
		r.SetOptions(conf.Options)

		// Cloud config
		cloudConfig := conf.Collectors.Cloud
		if err := envconfig.Process("k6", &cloudConfig); err != nil ***REMOVED***
			return err
		***REMOVED***
		if cloudConfig.Token == "" ***REMOVED***
			return errors.New("Not logged in, please use `k6 login cloud`.")
		***REMOVED***

		// Create a ticker to add a dot to the console every 0.5s
		ticker := time.NewTicker(time.Millisecond * 500)
		go func() ***REMOVED***
			for range ticker.C ***REMOVED***
				fmt.Fprint(stdout, ".")
			***REMOVED***
		***REMOVED***()

		// Start cloud test run
		client := cloud.NewClient(cloudConfig.Token, cloudConfig.Host, Version)

		arc := r.MakeArchive()
		if err := client.ValidateOptions(arc.Options); err != nil ***REMOVED***
			return err
		***REMOVED***

		if val, ok := arc.Options.External["loadimpact"]; ok ***REMOVED***
			if err := mapstructure.Decode(val, &cloudConfig); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		name := cloudConfig.Name
		if name == "" ***REMOVED***
			name = filepath.Base(filename)
		***REMOVED***

		refID, err := client.StartCloudTestRun(name, arc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ticker.Stop()

		testURL := cloud.URLForResults(refID, cloudConfig)
		fmt.Fprint(stdout, "\n\n")
		fmt.Fprintf(stdout, "     execution: %s\n", ui.ValueColor.Sprint("cloud"))
		fmt.Fprintf(stdout, "     script: %s\n", ui.ValueColor.Sprint(filename))
		fmt.Fprintf(stdout, "     output: %s\n", ui.ValueColor.Sprint(testURL))
		return nil
	***REMOVED***,
***REMOVED***

func init() ***REMOVED***
	RootCmd.AddCommand(cloudCmd)
	cloudCmd.Flags().AddFlagSet(optionFlagSet())
***REMOVED***
