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
	"compress/gzip"
	stdlibjson "encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/output"
	"go.k6.io/k6/stats"
)

// TODO: add option for emitting proper JSON files (https://github.com/loadimpact/k6/issues/737)
const flushPeriod = 200 * time.Millisecond // TODO: make this configurable

// Output funnels all passed metrics to an (optionally gzipped) JSON file.
type Output struct ***REMOVED***
	output.SampleBuffer

	params          output.Params
	periodicFlusher *output.PeriodicFlusher

	logger      logrus.FieldLogger
	filename    string
	encoder     *stdlibjson.Encoder
	closeFn     func() error
	seenMetrics map[string]struct***REMOVED******REMOVED***
	thresholds  map[string][]*stats.Threshold
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
		o.encoder = stdlibjson.NewEncoder(o.params.StdOut)
		o.closeFn = func() error ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		logfile, err := o.params.FS.Create(o.filename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if strings.HasSuffix(o.filename, ".gz") ***REMOVED***
			outfile := gzip.NewWriter(logfile)

			o.closeFn = func() error ***REMOVED***
				_ = outfile.Close()
				return logfile.Close()
			***REMOVED***
			o.encoder = stdlibjson.NewEncoder(outfile)
		***REMOVED*** else ***REMOVED***
			o.closeFn = logfile.Close
			o.encoder = stdlibjson.NewEncoder(logfile)
		***REMOVED***
	***REMOVED***

	o.encoder.SetEscapeHTML(false)

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
func (o *Output) SetThresholds(thresholds map[string]stats.Thresholds) ***REMOVED***
	ths := make(map[string][]*stats.Threshold)
	for name, t := range thresholds ***REMOVED***
		ths[name] = append(ths[name], t.Thresholds...)
	***REMOVED***
	o.thresholds = ths
***REMOVED***

func (o *Output) flushMetrics() ***REMOVED***
	samples := o.GetBufferedSamples()
	start := time.Now()
	var count int
	for _, sc := range samples ***REMOVED***
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples ***REMOVED***
			sample := sample
			sample.Metric.Thresholds.Thresholds = o.thresholds[sample.Metric.Name]
			o.handleMetric(sample.Metric)
			err := o.encoder.Encode(WrapSample(sample))
			if err != nil ***REMOVED***
				// Skip metric if it can't be made into JSON or envelope is null.
				o.logger.WithError(err).Error("Sample couldn't be marshalled to JSON")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if count > 0 ***REMOVED***
		o.logger.WithField("t", time.Since(start)).WithField("count", count).Debug("Wrote metrics to JSON")
	***REMOVED***
***REMOVED***

func (o *Output) handleMetric(m *stats.Metric) ***REMOVED***
	if _, ok := o.seenMetrics[m.Name]; ok ***REMOVED***
		return
	***REMOVED***
	o.seenMetrics[m.Name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	err := o.encoder.Encode(wrapMetric(m))
	if err != nil ***REMOVED***
		o.logger.WithError(err).Error("Metric couldn't be marshalled to JSON")
	***REMOVED***
***REMOVED***
