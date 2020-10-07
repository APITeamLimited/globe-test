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

package json

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

type Collector struct ***REMOVED***
	closeFn     func() error
	fname       string
	seenMetrics []string
	logger      logrus.FieldLogger

	encoder *json.Encoder

	buffer     []stats.Sample
	bufferLock sync.Mutex
***REMOVED***

// Verify that Collector implements lib.Collector
var _ lib.Collector = &Collector***REMOVED******REMOVED***

func (c *Collector) HasSeenMetric(str string) bool ***REMOVED***
	for _, n := range c.seenMetrics ***REMOVED***
		if n == str ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// New return new JSON collector
func New(logger logrus.FieldLogger, fs afero.Fs, fname string) (*Collector, error) ***REMOVED***
	c := &Collector***REMOVED***
		fname:  fname,
		logger: logger,
	***REMOVED***
	if fname == "" || fname == "-" ***REMOVED***
		c.encoder = json.NewEncoder(os.Stdout)
		c.closeFn = func() error ***REMOVED***
			return nil
		***REMOVED***
		return c, nil
	***REMOVED***
	logfile, err := fs.Create(c.fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if strings.HasSuffix(c.fname, ".gz") ***REMOVED***
		outfile := gzip.NewWriter(logfile)

		c.closeFn = func() error ***REMOVED***
			_ = outfile.Close()
			return logfile.Close()
		***REMOVED***
		c.encoder = json.NewEncoder(outfile)
	***REMOVED*** else ***REMOVED***
		c.closeFn = logfile.Close
		c.encoder = json.NewEncoder(logfile)
	***REMOVED***

	return c, nil
***REMOVED***

func (c *Collector) Init() error ***REMOVED***
	return nil
***REMOVED***

func (c *Collector) SetRunStatus(status lib.RunStatus) ***REMOVED******REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	const timeout = 200
	c.logger.Debug("JSON output: Running!")
	ticker := time.NewTicker(time.Millisecond * timeout)
	defer func() ***REMOVED***
		_ = c.closeFn()
	***REMOVED***()
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			c.commit()
		case <-ctx.Done():
			c.commit()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) HandleMetric(m *stats.Metric) ***REMOVED***
	if c.HasSeenMetric(m.Name) ***REMOVED***
		return
	***REMOVED***

	c.seenMetrics = append(c.seenMetrics, m.Name)
	err := c.encoder.Encode(WrapMetric(m))
	if err != nil ***REMOVED***
		c.logger.WithField("filename", c.fname).WithError(err).Warning(
			"JSON: Envelope is nil or Metric couldn't be marshalled to JSON")
		return
	***REMOVED***
***REMOVED***

func (c *Collector) Collect(scs []stats.SampleContainer) ***REMOVED***
	c.bufferLock.Lock()
	defer c.bufferLock.Unlock()
	for _, sc := range scs ***REMOVED***
		c.buffer = append(c.buffer, sc.GetSamples()...)
	***REMOVED***
***REMOVED***

func (c *Collector) commit() ***REMOVED***
	c.bufferLock.Lock()
	samples := c.buffer
	c.buffer = nil
	c.bufferLock.Unlock()
	start := time.Now()
	var count int
	for _, sc := range samples ***REMOVED***
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range sc.GetSamples() ***REMOVED***
			sample := sample
			c.HandleMetric(sample.Metric)
			err := c.encoder.Encode(WrapSample(&sample))
			if err != nil ***REMOVED***
				// Skip metric if it can't be made into JSON or envelope is null.
				c.logger.WithField("filename", c.fname).WithError(err).Warning(
					"JSON: Sample couldn't be marshalled to JSON")
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if count > 0 ***REMOVED***
		c.logger.WithField("filename", c.fname).WithField("t", time.Since(start)).
			WithField("count", count).Debug("JSON: Wrote JSON metrics")
	***REMOVED***
***REMOVED***

func (c *Collector) Link() string ***REMOVED***
	return ""
***REMOVED***

// GetRequiredSystemTags returns which sample tags are needed by this collector
func (c *Collector) GetRequiredSystemTags() stats.SystemTagSet ***REMOVED***
	return stats.SystemTagSet(0) // There are no required tags for this collector
***REMOVED***
