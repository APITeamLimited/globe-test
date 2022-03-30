/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package json

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/klauspost/compress/gzip"

	"github.com/mailru/easyjson/jwriter"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

// TODO: add option for emitting proper JSON files (https://github.com/k6io/k6/issues/737)
const flushPeriod = 200 * time.Millisecond // TODO: make this configurable

// Output funnels all passed metrics to an (optionally gzipped) JSON file.
type Output struct ***REMOVED***
	output.SampleBuffer

	params          output.Params
	periodicFlusher *output.PeriodicFlusher

	logger      logrus.FieldLogger
	filename    string
	out         io.Writer
	closeFn     func() error
	seenMetrics map[string]struct***REMOVED******REMOVED***
	thresholds  map[string][]*metrics.Threshold
***REMOVED***

// New returns a new JSON output.
func New(params output.Params) (output.Output, error) ***REMOVED***
	return &Output***REMOVED***
		params:   params,
		filename: params.ConfigArgument,
		logger: params.Logger.WithFields(logrus.Fields***REMOVED***
			"output":   "json",
			"filename": params.ConfigArgument,
		***REMOVED***),
		seenMetrics: make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***, nil
***REMOVED***

// Description returns a human-readable description of the output.
func (o *Output) Description() string ***REMOVED***
	if o.filename == "" || o.filename == "-" ***REMOVED***
		return "json(stdout)"
	***REMOVED***
	return fmt.Sprintf("json (%s)", o.filename)
***REMOVED***

// Start tries to open the specified JSON file and starts the goroutine for
// metric flushing. If gzip encoding is specified, it also handles that.
func (o *Output) Start() error ***REMOVED***
	o.logger.Debug("Starting...")

	if o.filename == "" || o.filename == "-" ***REMOVED***
		w := bufio.NewWriter(o.params.StdOut)
		o.closeFn = func() error ***REMOVED***
			return w.Flush()
		***REMOVED***
		o.out = w
	***REMOVED*** else ***REMOVED***
		logfile, err := o.params.FS.Create(o.filename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		w := bufio.NewWriter(logfile)

		if strings.HasSuffix(o.filename, ".gz") ***REMOVED***
			outfile := gzip.NewWriter(w)

			o.closeFn = func() error ***REMOVED***
				_ = outfile.Close()
				_ = w.Flush()
				return logfile.Close()
			***REMOVED***
			o.out = outfile
		***REMOVED*** else ***REMOVED***
			o.closeFn = func() error ***REMOVED***
				_ = w.Flush()
				return logfile.Close()
			***REMOVED***
			o.out = logfile
		***REMOVED***
	***REMOVED***

	pf, err := output.NewPeriodicFlusher(flushPeriod, o.flushMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.logger.Debug("Started!")
	o.periodicFlusher = pf

	return nil
***REMOVED***

// Stop flushes any remaining metrics and stops the goroutine.
func (o *Output) Stop() error ***REMOVED***
	o.logger.Debug("Stopping...")
	defer o.logger.Debug("Stopped!")
	o.periodicFlusher.Stop()
	return o.closeFn()
***REMOVED***

// SetThresholds receives the thresholds before the output is Start()-ed.
func (o *Output) SetThresholds(thresholds map[string]metrics.Thresholds) ***REMOVED***
	ths := make(map[string][]*metrics.Threshold)
	for name, t := range thresholds ***REMOVED***
		ths[name] = append(ths[name], t.Thresholds...)
	***REMOVED***
	o.thresholds = ths
***REMOVED***

func (o *Output) flushMetrics() ***REMOVED***
	samples := o.GetBufferedSamples()
	start := time.Now()
	var count int
	jw := new(jwriter.Writer)
	for _, sc := range samples ***REMOVED***
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples ***REMOVED***
			sample := sample
			sample.Metric.Thresholds.Thresholds = o.thresholds[sample.Metric.Name]
			o.handleMetric(sample.Metric, jw)
			wrapSample(sample).MarshalEasyJSON(jw)
			jw.RawByte('\n')
		***REMOVED***
	***REMOVED***

	if _, err := jw.DumpTo(o.out); err != nil ***REMOVED***
		// Skip metric if it can't be made into JSON or envelope is null.
		o.logger.WithError(err).Error("Sample couldn't be marshalled to JSON")
	***REMOVED***
	if count > 0 ***REMOVED***
		o.logger.WithField("t", time.Since(start)).WithField("count", count).Debug("Wrote metrics to JSON")
	***REMOVED***
***REMOVED***

func (o *Output) handleMetric(m *metrics.Metric, jw *jwriter.Writer) ***REMOVED***
	if _, ok := o.seenMetrics[m.Name]; ok ***REMOVED***
		return
	***REMOVED***
	o.seenMetrics[m.Name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	wrapMetric(m).MarshalEasyJSON(jw)
	jw.RawByte('\n')
***REMOVED***
