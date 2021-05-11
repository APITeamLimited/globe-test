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

	"golang.org/x/text/unicode/norm"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/stats"
)

// TODO: delete everything here after we move it to a JS function

const (
	groupPrefix   = "█"
	detailsPrefix = "↳"

	succMark = "✓"
	failMark = "✗"
)

// Summary handles test summary output
type Summary struct ***REMOVED***
	trendColumns        []string
	trendValueResolvers map[string]func(s *stats.TrendSink) float64
***REMOVED***

// NewSummary returns a new Summary instance, used for writing a
// summary/report of the test metrics data.
func NewSummary(cols []string) *Summary ***REMOVED***
	s := Summary***REMOVED***trendColumns: cols***REMOVED***

	s.trendValueResolvers, _ = stats.GetResolversForTrendColumns(cols)
	return &s
***REMOVED***

// StrWidth returns the actual width of the string.
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

func summarizeCheck(w io.Writer, indent string, check *lib.Check) ***REMOVED***
	mark := succMark
	color := SuccColor
	if check.Fails > 0 ***REMOVED***
		mark = failMark
		color = FailColor
	***REMOVED***
	_, _ = color.Fprintf(w, "%s%s %s\n", indent, mark, check.Name)
	if check.Fails > 0 ***REMOVED***
		_, _ = color.Fprintf(w, "%s %s  %d%% — %s %d / %s %d\n",
			indent, detailsPrefix,
			int(100*(float64(check.Passes)/float64(check.Fails+check.Passes))),
			succMark, check.Passes, failMark, check.Fails,
		)
	***REMOVED***
***REMOVED***

func summarizeGroup(w io.Writer, indent string, group *lib.Group) ***REMOVED***
	if group.Name != "" ***REMOVED***
		_, _ = fmt.Fprintf(w, "%s%s %s\n\n", indent, groupPrefix, group.Name)
		indent = indent + "  "
	***REMOVED***

	var checkNames []string
	for _, check := range group.OrderedChecks ***REMOVED***
		checkNames = append(checkNames, check.Name)
	***REMOVED***
	for _, name := range checkNames ***REMOVED***
		summarizeCheck(w, indent, group.Checks[name])
	***REMOVED***
	if len(checkNames) > 0 ***REMOVED***
		_, _ = fmt.Fprintf(w, "\n")
	***REMOVED***

	var groupNames []string
	for _, grp := range group.OrderedGroups ***REMOVED***
		groupNames = append(groupNames, grp.Name)
	***REMOVED***
	for _, name := range groupNames ***REMOVED***
		summarizeGroup(w, indent, group.Groups[name])
	***REMOVED***
***REMOVED***

func nonTrendMetricValueForSum(t time.Duration, timeUnit string, m *stats.Metric) (data string, extra []string) ***REMOVED***
	switch sink := m.Sink.(type) ***REMOVED***
	case *stats.CounterSink:
		value := sink.Value
		rate := 0.0
		if t > 0 ***REMOVED***
			rate = value / (float64(t) / float64(time.Second))
		***REMOVED***
		return m.HumanizeValue(value, timeUnit), []string***REMOVED***m.HumanizeValue(rate, timeUnit) + "/s"***REMOVED***
	case *stats.GaugeSink:
		value := sink.Value
		min := sink.Min
		max := sink.Max
		return m.HumanizeValue(value, timeUnit), []string***REMOVED***
			"min=" + m.HumanizeValue(min, timeUnit),
			"max=" + m.HumanizeValue(max, timeUnit),
		***REMOVED***
	case *stats.RateSink:
		value := float64(sink.Trues) / float64(sink.Total)
		passes := sink.Trues
		fails := sink.Total - sink.Trues
		return m.HumanizeValue(value, timeUnit), []string***REMOVED***
			"✓ " + strconv.FormatInt(passes, 10),
			"✗ " + strconv.FormatInt(fails, 10),
		***REMOVED***
	default:
		return "[no data]", nil
	***REMOVED***
***REMOVED***

func displayNameForMetric(m *stats.Metric) string ***REMOVED***
	if m.Sub.Parent != "" ***REMOVED***
		return "***REMOVED*** " + m.Sub.Suffix + " ***REMOVED***"
	***REMOVED***
	return m.Name
***REMOVED***

func indentForMetric(m *stats.Metric) string ***REMOVED***
	if m.Sub.Parent != "" ***REMOVED***
		return "  "
	***REMOVED***
	return ""
***REMOVED***

// nolint:funlen
func (s *Summary) summarizeMetrics(w io.Writer, indent string, t time.Duration,
	timeUnit string, metrics map[string]*stats.Metric) ***REMOVED***
	names := []string***REMOVED******REMOVED***
	nameLenMax := 0

	values := make(map[string]string)
	valueMaxLen := 0
	extras := make(map[string][]string)
	extraMaxLens := make([]int, 2)

	trendCols := make(map[string][]string)
	trendColMaxLens := make([]int, len(s.trendColumns))

	for name, m := range metrics ***REMOVED***
		names = append(names, name)

		// When calculating widths for metrics, account for the indentation on submetrics.
		displayName := displayNameForMetric(m) + indentForMetric(m)
		if l := StrWidth(displayName); l > nameLenMax ***REMOVED***
			nameLenMax = l
		***REMOVED***

		m.Sink.Calc()
		if sink, ok := m.Sink.(*stats.TrendSink); ok ***REMOVED***
			cols := make([]string, len(s.trendColumns))

			for i, tc := range s.trendColumns ***REMOVED***
				var value string

				resolver := s.trendValueResolvers[tc]

				v := resolver(sink)
				if tc != "count" ***REMOVED*** // sigh
					value = m.HumanizeValue(v, timeUnit)
				***REMOVED*** else ***REMOVED***
					value = strconv.FormatInt(int64(v), 10)
				***REMOVED***
				if l := StrWidth(value); l > trendColMaxLens[i] ***REMOVED***
					trendColMaxLens[i] = l
				***REMOVED***
				cols[i] = value
			***REMOVED***
			trendCols[name] = cols
			continue
		***REMOVED***

		value, extra := nonTrendMetricValueForSum(t, timeUnit, m)
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

	tmpCols := make([]string, len(s.trendColumns))
	for _, name := range names ***REMOVED***
		m := metrics[name]

		mark := " "
		markColor := StdColor
		if m.Tainted.Valid ***REMOVED***
			if m.Tainted.Bool ***REMOVED***
				mark = failMark
				markColor = FailColor
			***REMOVED*** else ***REMOVED***
				mark = succMark
				markColor = SuccColor
			***REMOVED***
		***REMOVED***

		fmtName := displayNameForMetric(m)
		fmtIndent := indentForMetric(m)
		fmtName += GrayColor.Sprint(strings.Repeat(".", nameLenMax-StrWidth(fmtName)-StrWidth(fmtIndent)+3) + ":")

		var fmtData string
		if cols := trendCols[name]; cols != nil ***REMOVED***
			for i, val := range cols ***REMOVED***
				tmpCols[i] = s.trendColumns[i] + "=" + ValueColor.Sprint(val) +
					strings.Repeat(" ", trendColMaxLens[i]-StrWidth(val))
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
		_, _ = fmt.Fprint(w, indent+fmtIndent+markColor.Sprint(mark)+" "+fmtName+" "+fmtData+"\n")
	***REMOVED***
***REMOVED***

// SummaryData represents data passed to Summary.SummarizeMetrics
type SummaryData struct ***REMOVED***
	Metrics   map[string]*stats.Metric
	RootGroup *lib.Group
	Time      time.Duration
	TimeUnit  string
***REMOVED***

// SummarizeMetrics creates a summary of provided metrics and writes it to w.
func (s *Summary) SummarizeMetrics(w io.Writer, indent string, data SummaryData) ***REMOVED***
	if data.RootGroup != nil ***REMOVED***
		summarizeGroup(w, indent+"    ", data.RootGroup)
	***REMOVED***

	s.summarizeMetrics(w, indent+"  ", data.Time, data.TimeUnit, data.Metrics)
***REMOVED***
