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
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"golang.org/x/text/unicode/norm"
)

const (
	GroupPrefix   = "█"
	DetailsPrefix = "↳"

	SuccMark = "✓"
	FailMark = "✗"
)

var (
	StdColor      = color.New()                          // Default color.
	SuccColor     = color.New(color.FgGreen)             // Successful stuff.
	FailColor     = color.New(color.FgRed)               // Failed stuff.
	GrayColor     = color.New(color.Faint)               // Padding and disabled stuff.
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

// Returns the actual width of the string.
func StrWidth(s string) (n int) ***REMOVED***
	var it norm.Iter
	it.InitString(norm.NFKD, s)

	inEscSeq := false
	inLongEscSeq := false
	for !it.Done() ***REMOVED***
		data := it.Next()

		// Skip over ANSI escape codes.
		if data[0] == '\x1b' ***REMOVED***
			inEscSeq = true
			continue
		***REMOVED***
		if inEscSeq && data[0] == '[' ***REMOVED***
			inLongEscSeq = true
			continue
		***REMOVED***
		if inEscSeq && inLongEscSeq && data[0] >= 0x40 && data[0] <= 0x7E ***REMOVED***
			inEscSeq = false
			inLongEscSeq = false
			continue
		***REMOVED***
		if inEscSeq && !inLongEscSeq && data[0] >= 0x40 && data[0] <= 0x5F ***REMOVED***
			inEscSeq = false
			continue
		***REMOVED***

		n++
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

func NonTrendMetricValueForSum(t time.Duration, m *stats.Metric) (data string, extra []string) ***REMOVED***
	switch sink := m.Sink.(type) ***REMOVED***
	case *stats.CounterSink:
		value := sink.Value
		rate := value / (float64(t) / float64(time.Second))
		return m.HumanizeValue(value), []string***REMOVED***m.HumanizeValue(rate) + "/s"***REMOVED***
	case *stats.GaugeSink:
		value := sink.Value
		min := sink.Min
		max := sink.Max
		return m.HumanizeValue(value), []string***REMOVED***
			"min=" + m.HumanizeValue(min),
			"max=" + m.HumanizeValue(max),
		***REMOVED***
	case *stats.RateSink:
		value := float64(sink.Trues) / float64(sink.Total)
		passes := sink.Trues
		fails := sink.Total - sink.Trues
		return m.HumanizeValue(value), []string***REMOVED***
			"✓ " + strconv.FormatInt(passes, 10),
			"✗ " + strconv.FormatInt(fails, 10),
		***REMOVED***
	default:
		return "[no data]", nil
	***REMOVED***
***REMOVED***

func DisplayNameForMetric(m *stats.Metric) string ***REMOVED***
	if m.Sub.Parent != "" ***REMOVED***
		return "***REMOVED*** " + m.Sub.Suffix + " ***REMOVED***"
	***REMOVED***
	return m.Name
***REMOVED***

func IndentForMetric(m *stats.Metric) string ***REMOVED***
	if m.Sub.Parent != "" ***REMOVED***
		return "  "
	***REMOVED***
	return ""
***REMOVED***

func SummarizeMetrics(w io.Writer, indent string, t time.Duration, metrics map[string]*stats.Metric) ***REMOVED***
	names := []string***REMOVED******REMOVED***
	nameLenMax := 0

	values := make(map[string]string)
	valueMaxLen := 0
	extras := make(map[string][]string)
	extraMaxLens := make([]int, 2)

	trendCols := make(map[string][]string)
	trendColMaxLens := make([]int, len(TrendColumns))

	for name, m := range metrics ***REMOVED***
		names = append(names, name)

		// When calculating widths for metrics, account for the indentation on submetrics.
		displayName := DisplayNameForMetric(m) + IndentForMetric(m)
		if l := StrWidth(displayName); l > nameLenMax ***REMOVED***
			nameLenMax = l
		***REMOVED***

		m.Sink.Calc()
		if sink, ok := m.Sink.(*stats.TrendSink); ok ***REMOVED***
			cols := make([]string, len(TrendColumns))
			for i, col := range TrendColumns ***REMOVED***
				value := m.HumanizeValue(col.Get(sink))
				if l := StrWidth(value); l > trendColMaxLens[i] ***REMOVED***
					trendColMaxLens[i] = l
				***REMOVED***
				cols[i] = value
			***REMOVED***
			trendCols[name] = cols
			continue
		***REMOVED***

		value, extra := NonTrendMetricValueForSum(t, m)
		values[name] = value
		if l := StrWidth(value); l > valueMaxLen ***REMOVED***
			valueMaxLen = l
		***REMOVED***
		extras[name] = extra
		if len(extra) > 1 ***REMOVED***
			for i, ex := range extra ***REMOVED***
				if l := StrWidth(ex); l > extraMaxLens[i] ***REMOVED***
					extraMaxLens[i] = l
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	sort.Strings(names)
	tmpCols := make([]string, len(TrendColumns))
	for _, name := range names ***REMOVED***
		m := metrics[name]

		mark := " "
		markColor := StdColor
		if m.Tainted.Valid ***REMOVED***
			if m.Tainted.Bool ***REMOVED***
				mark = FailMark
				markColor = FailColor
			***REMOVED*** else ***REMOVED***
				mark = SuccMark
				markColor = SuccColor
			***REMOVED***
		***REMOVED***

		fmtName := DisplayNameForMetric(m)
		fmtIndent := IndentForMetric(m)
		fmtName += GrayColor.Sprint(strings.Repeat(".", nameLenMax-StrWidth(fmtName)-StrWidth(fmtIndent)+3) + ":")

		fmtData := ""
		if cols := trendCols[name]; cols != nil ***REMOVED***
			for i, val := range cols ***REMOVED***
				tmpCols[i] = TrendColumns[i].Key + "=" + ValueColor.Sprint(val) + strings.Repeat(" ", trendColMaxLens[i]-StrWidth(val))
			***REMOVED***
			fmtData = strings.Join(tmpCols, " ")
		***REMOVED*** else ***REMOVED***
			value := values[name]
			fmtData = ValueColor.Sprint(value) + strings.Repeat(" ", valueMaxLen-StrWidth(value))

			extra := extras[name]
			switch len(extra) ***REMOVED***
			case 0:
			case 1:
				fmtData = fmtData + " " + ExtraColor.Sprint(extra[0])
			default:
				parts := make([]string, len(extra))
				for i, ex := range extra ***REMOVED***
					parts[i] = ExtraColor.Sprint(ex) + strings.Repeat(" ", extraMaxLens[i]-StrWidth(ex))
				***REMOVED***
				fmtData = fmtData + " " + ExtraColor.Sprint(strings.Join(parts, " "))
			***REMOVED***
		***REMOVED***
		fmt.Fprint(w, indent+fmtIndent+markColor.Sprint(mark)+" "+fmtName+" "+fmtData+"\n")
	***REMOVED***
***REMOVED***

// Summarizes a dataset and returns whether the test run was considered a success.
func Summarize(w io.Writer, indent string, data SummaryData) ***REMOVED***
	SummarizeGroup(w, indent+"    ", data.Root)
	SummarizeMetrics(w, indent+"  ", data.Time, data.Metrics)
***REMOVED***
