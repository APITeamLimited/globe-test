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

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/types"
)

func getInspectCmd(logger *logrus.Logger, globalFlags *commandFlags) *cobra.Command ***REMOVED***
	var addExecReqs bool

	// inspectCmd represents the inspect command
	inspectCmd := &cobra.Command***REMOVED***
		Use:   "inspect [file]",
		Short: "Inspect a script or archive",
		Long:  `Inspect a script or archive.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error ***REMOVED***
			src, filesystems, err := readSource(args[0], logger)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			runtimeOptions, err := getRuntimeOptions(cmd.Flags(), buildEnvMap(os.Environ()))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			registry := metrics.NewRegistry()

			var b *js.Bundle
			typ := globalFlags.runType
			if typ == "" ***REMOVED***
				typ = detectType(src.Data)
			***REMOVED***
			switch typ ***REMOVED***
			// this is an exhaustive list
			case typeArchive:
				var arc *lib.Archive
				arc, err = lib.ReadArchive(bytes.NewBuffer(src.Data))
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				b, err = js.NewBundleFromArchive(logger, arc, runtimeOptions, registry)

			case typeJS:
				b, err = js.NewBundle(logger, src, filesystems, runtimeOptions, registry)
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			// ATM, output can take 2 forms: standard (equal to lib.Options struct) and extended, with additional fields.
			inspectOutput := interface***REMOVED******REMOVED***(b.Options)

			if addExecReqs ***REMOVED***
				inspectOutput, err = addExecRequirements(b, logger, globalFlags)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***

			data, err := json.MarshalIndent(inspectOutput, "", "  ")
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			fmt.Println(string(data)) //nolint:forbidigo // yes we want to just print it

			return nil
		***REMOVED***,
	***REMOVED***

	inspectCmd.Flags().SortFlags = false
	inspectCmd.Flags().AddFlagSet(runtimeOptionFlagSet(false))
	inspectCmd.Flags().StringVarP(&globalFlags.runType, "type", "t", globalFlags.runType, "override file `type`, \"js\" or \"archive\"") //nolint:lll
	inspectCmd.Flags().BoolVar(&addExecReqs,
		"execution-requirements",
		false,
		"include calculations of execution requirements for the test")

	return inspectCmd
***REMOVED***

func addExecRequirements(b *js.Bundle, logger *logrus.Logger, globalFlags *commandFlags) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	conf, err := getConsolidatedConfig(
		afero.NewOsFs(), Config***REMOVED******REMOVED***, b.Options, buildEnvMap(os.Environ()), globalFlags)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	conf, err = deriveAndValidateConfig(conf, b.IsExecutable, logger)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	et, err := lib.NewExecutionTuple(conf.ExecutionSegment, conf.ExecutionSegmentSequence)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	executionPlan := conf.Scenarios.GetFullExecutionRequirements(et)
	duration, _ := lib.GetEndOffset(executionPlan)

	return struct ***REMOVED***
		lib.Options
		TotalDuration types.NullDuration `json:"totalDuration"`
		MaxVUs        uint64             `json:"maxVUs"`
	***REMOVED******REMOVED***
		conf.Options,
		types.NewNullDuration(duration, true),
		lib.GetMaxPossibleVUs(executionPlan),
	***REMOVED***, nil
***REMOVED***
