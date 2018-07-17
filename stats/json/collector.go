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
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Collector struct ***REMOVED***
	outfile     io.WriteCloser
	fname       string
	seenMetrics []string
***REMOVED***

// Verify that Collector implements lib.Collector
var _ lib.Collector = &Collector***REMOVED******REMOVED***

// Similar to ioutil.NopCloser, but for writers
type nopCloser struct ***REMOVED***
	io.Writer
***REMOVED***

func (nopCloser) Close() error ***REMOVED*** return nil ***REMOVED***

func (c *Collector) HasSeenMetric(str string) bool ***REMOVED***
	for _, n := range c.seenMetrics ***REMOVED***
		if n == str ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func New(fs afero.Fs, fname string) (*Collector, error) ***REMOVED***
	if fname == "" || fname == "-" ***REMOVED***
		return &Collector***REMOVED***
			outfile: nopCloser***REMOVED***os.Stdout***REMOVED***,
			fname:   "-",
		***REMOVED***, nil
	***REMOVED***

	logfile, err := fs.Create(fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Collector***REMOVED***
		outfile: logfile,
		fname:   fname,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) Init() error ***REMOVED***
	return nil
***REMOVED***

func (c *Collector) SetRunStatus(status lib.RunStatus) ***REMOVED******REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	log.WithField("filename", c.fname).Debug("JSON: Writing JSON metrics")
	<-ctx.Done()
	_ = c.outfile.Close()
***REMOVED***

func (c *Collector) HandleMetric(m *stats.Metric) ***REMOVED***
	if c.HasSeenMetric(m.Name) ***REMOVED***
		return
	***REMOVED***

	c.seenMetrics = append(c.seenMetrics, m.Name)
	env := WrapMetric(m)
	row, err := json.Marshal(env)

	if env == nil || err != nil ***REMOVED***
		log.WithField("filename", c.fname).Warning(
			"JSON: Envelope is nil or Metric couldn't be marshalled to JSON")
		return
	***REMOVED***

	row = append(row, '\n')
	_, err = c.outfile.Write(row)
	if err != nil ***REMOVED***
		log.WithField("filename", c.fname).Error("JSON: Error writing to file")
	***REMOVED***
***REMOVED***

func (c *Collector) Collect(scs []stats.SampleContainer) ***REMOVED***
	for _, sc := range scs ***REMOVED***
		for _, sample := range sc.GetSamples() ***REMOVED***
			c.HandleMetric(sample.Metric)

			env := WrapSample(&sample)
			row, err := json.Marshal(env)

			if err != nil || env == nil ***REMOVED***
				// Skip metric if it can't be made into JSON or envelope is null.
				log.WithField("filename", c.fname).Warning(
					"JSON: Envelope is nil or Sample couldn't be marshalled to JSON")
				continue
			***REMOVED***
			row = append(row, '\n')
			_, err = c.outfile.Write(row)
			if err != nil ***REMOVED***
				log.WithField("filename", c.fname).Error("JSON: Error writing to file")
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) Link() string ***REMOVED***
	return ""
***REMOVED***

// GetRequiredSystemTags returns which sample tags are needed by this collector
func (c *Collector) GetRequiredSystemTags() lib.TagSet ***REMOVED***
	return lib.TagSet***REMOVED******REMOVED*** // There are no required tags for this collector
***REMOVED***
