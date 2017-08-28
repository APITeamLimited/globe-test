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
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"golang.org/x/text/unicode/norm"
)

const (
	GroupPrefix   = "█"
	DetailsPrefix = "↪"

	SuccMark = "✓"
	FailMark = "✗"
)

var (
	SuccColor     = color.New(color.FgGreen)             // Successful stuff.
	FailColor     = color.New(color.FgRed)               // Failed stuff.
	NamePadColor  = color.New(color.Faint)               // Padding for metric names.
	ValueColor    = color.New(color.FgCyan)              // Values of all kinds.
	ExtraColor    = color.New(color.FgCyan, color.Faint) // Extra annotations for values.
	ExtraKeyColor = color.New(color.Faint)               // Keys inside extra annotations.

	TrendColumns = []TrendColumn***REMOVED***
		***REMOVED***"avg", func(s *stats.TrendSink) float64 ***REMOVED*** return s.Avg ***REMOVED******REMOVED***,
		***REMOVED***"min", func(s *stats.TrendSink) float64 ***REMOVED*** return s.Min ***REMOVED******REMOVED***,
		***REMOVED***"med", func(s *stats.TrendSink) float64 ***REMOVED*** return s.Med ***REMOVED******REMOVED***,
		***REMOVED***"max", func(s *stats.TrendSink) float64 ***REMOVED*** return s.Max ***REMOVED******REMOVED***,
		***REMOVED***"p(90)", func(s *stats.TrendSink) float64 ***REMOVED*** return s.P(0.90) ***REMOVED******REMOVED***,
		***REMOVED***"p(95)", func(s *stats.TrendSink) float64 ***REMOVED*** return s.P(0.95) ***REMOVED******REMOVED***,
	***REMOVED***
)

type TrendColumn struct ***REMOVED***
	Key string
	Get func(s *stats.TrendSink) float64
***REMOVED***

// Returns the number of unicode glyphs in a string.
func NumGlyph(s string) (n int) ***REMOVED***
	var it norm.Iter
	it.InitString(norm.NFKD, s)
	for !it.Done() ***REMOVED***
		n++
		it.Next()
	***REMOVED***
	return
***REMOVED***

// SummaryData represents data passed to Summarize.
type SummaryData struct ***REMOVED***
	Opts    lib.Options
	Root    *lib.Group
	Metrics map[string]*stats.Metric
	Time    time.Duration
***REMOVED***

func SummarizeCheck(w io.Writer, indent string, check *lib.Check) ***REMOVED***
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

func SummarizeGroup(w io.Writer, indent string, group *lib.Group) ***REMOVED***
	if group.Name != "" ***REMOVED***
		_, _ = fmt.Fprintf(w, "%s%s %s\n\n", indent, GroupPrefix, group.Name)
		indent = indent + "  "
	***REMOVED***

	var checkNames []string
	for _, check := range group.Checks ***REMOVED***
		checkNames = append(checkNames, check.Name)
	***REMOVED***
	sort.Strings(checkNames)
	for _, name := range checkNames ***REMOVED***
		SummarizeCheck(w, indent, group.Checks[name])
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
		SummarizeGroup(w, indent, group.Groups[name])
	***REMOVED***
***REMOVED***

func NonTrendMetricValueForSum(t time.Duration, m *stats.Metric) (data, extra string) ***REMOVED***
	m.Sink.Calc()
	switch sink := m.Sink.(type) ***REMOVED***
	case *stats.CounterSink:
		value := m.HumanizeValue(sink.Value)
		rate := m.HumanizeValue(sink.Value / float64(t/time.Second))
		return value, rate + "/s"
	case *stats.GaugeSink:
		value := m.HumanizeValue(sink.Value)
		min := m.HumanizeValue(sink.Min)
		max := m.HumanizeValue(sink.Max)
		return value, "min=" + min + " max=" + max
	case *stats.RateSink:
		value := m.HumanizeValue(float64(sink.Trues) / float64(sink.Total))
		passes := sink.Trues
		fails := sink.Total - sink.Trues
		return value, fmt.Sprintf("✓ %d ✗ %d", passes, fails)
	default:
		return "[no data]", ""
	***REMOVED***
***REMOVED***

func SummarizeMetrics(w io.Writer, indent string, t time.Duration, metrics map[string]*stats.Metric) ***REMOVED***
	names := []string***REMOVED******REMOVED***
	nameLenMax := 0

	values := make(map[string]string)
	extras := make(map[string]string)
	valueMaxLen := 0

	trendCols := make(map[string][]string)
	trendColMaxLens := make([]int, len(TrendColumns))

	for name, m := range metrics ***REMOVED***
		names = append(names, name)
		if l := NumGlyph(name); l > nameLenMax ***REMOVED***
			nameLenMax = l
		***REMOVED***

		if sink, ok := m.Sink.(*stats.TrendSink); ok ***REMOVED***
			cols := make([]string, len(TrendColumns))
			for i, col := range TrendColumns ***REMOVED***
				value := m.HumanizeValue(col.Get(sink))
				if l := NumGlyph(value); l > trendColMaxLens[i] ***REMOVED***
					trendColMaxLens[i] = l
				***REMOVED***
				cols[i] = value
			***REMOVED***
			trendCols[name] = cols
			continue
		***REMOVED***

		value, extra := NonTrendMetricValueForSum(t, m)
		values[name] = value
		extras[name] = extra
		if l := NumGlyph(value); l > valueMaxLen ***REMOVED***
			valueMaxLen = l
		***REMOVED***
	***REMOVED***

	sort.Strings(names)
	tmpCols := make([]string, len(TrendColumns))
	for _, name := range names ***REMOVED***
		m := metrics[name]

		mark := " "
		if m.Tainted.Valid ***REMOVED***
			if m.Tainted.Bool ***REMOVED***
				mark = FailColor.Sprint(FailMark)
			***REMOVED*** else ***REMOVED***
				mark = SuccColor.Sprint(SuccMark)
			***REMOVED***
		***REMOVED***

		fmtName := name + NamePadColor.Sprint(strings.Repeat(".", nameLenMax-NumGlyph(name)+3)+":")
		fmtData := ""
		if cols := trendCols[name]; cols != nil ***REMOVED***
			for i, val := range cols ***REMOVED***
				tmpCols[i] = TrendColumns[i].Key + "=" + ValueColor.Sprint(val) + strings.Repeat(" ", trendColMaxLens[i]-NumGlyph(val))
			***REMOVED***
			fmtData = strings.Join(tmpCols, " ")
		***REMOVED*** else ***REMOVED***
			value := values[name]
			fmtData = ValueColor.Sprint(value) + strings.Repeat(" ", valueMaxLen-NumGlyph(value))
			if extra := extras[name]; extra != "" ***REMOVED***
				fmtData = fmtData + " " + ExtraColor.Sprint(extra)
			***REMOVED***
		***REMOVED***
		fmt.Fprint(w, indent+mark+" "+fmtName+" "+fmtData+"\n")
	***REMOVED***
***REMOVED***

// Summarizes a dataset and returns whether the test run was considered a success.
func Summarize(w io.Writer, indent string, data SummaryData) ***REMOVED***
	SummarizeGroup(w, indent+"    ", data.Root)
	SummarizeMetrics(w, indent+"  ", data.Time, data.Metrics)
***REMOVED***
