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

package ui

import (
	"fmt"
	"io"
	"sort"

	"github.com/fatih/color"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

const (
	GroupPrefix   = "█"
	DetailsPrefix = "↪"

	SuccMark = "✓"
	FailMark = "✗"
)

var (
	SuccColor = color.New(color.FgGreen)
	FailColor = color.New(color.FgRed)
)

// SummaryData represents data passed to Summarize.
type SummaryData struct ***REMOVED***
	Opts    lib.Options
	Root    *lib.Group
	Metrics map[string]*stats.Metric
***REMOVED***

func SummarizeCheck(w io.Writer, tty bool, indent string, check *lib.Check) ***REMOVED***
	mark := SuccMark
	color := SuccColor
	if check.Fails > 0 ***REMOVED***
		mark = FailMark
		color = FailColor
	***REMOVED***
	_, _ = color.Fprintf(w, "%s%s %s\n", indent, mark, check.Name)
	if check.Fails > 0 ***REMOVED***
		_, _ = color.Fprintf(w, "%s %s  %d%% — %s %d / %s %d\n",
			indent, DetailsPrefix,
			int(100*(float64(check.Passes)/float64(check.Fails+check.Passes))),
			SuccMark, check.Passes, FailMark, check.Fails,
		)
	***REMOVED***
***REMOVED***

func SummarizeGroup(w io.Writer, tty bool, indent string, group *lib.Group) ***REMOVED***
	if group.Name != "" ***REMOVED***
		_, _ = fmt.Fprintf(w, "%s%s %s\n\n", indent, GroupPrefix, group.Name)
	***REMOVED***

	var checkNames []string
	for _, check := range group.Checks ***REMOVED***
		checkNames = append(checkNames, check.Name)
	***REMOVED***
	sort.Strings(checkNames)
	for _, name := range checkNames ***REMOVED***
		SummarizeCheck(w, tty, indent+"  ", group.Checks[name])
	***REMOVED***
	if len(checkNames) > 0 ***REMOVED***
		fmt.Fprintf(w, "\n")
	***REMOVED***

	var groupNames []string
	for _, grp := range group.Groups ***REMOVED***
		groupNames = append(groupNames, grp.Name)
	***REMOVED***
	sort.Strings(groupNames)
	for _, name := range groupNames ***REMOVED***
		SummarizeGroup(w, tty, indent+"  ", group.Groups[name])
	***REMOVED***
***REMOVED***

// Summarizes a dataset and returns whether the test run was considered a success.
func Summarize(w io.Writer, tty bool, indent string, data SummaryData) ***REMOVED***
	SummarizeGroup(w, tty, indent, data.Root)
***REMOVED***
