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
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/stats"
	"io"
	"net/url"
	"os"
)

type Collector struct ***REMOVED***
	outfile io.WriteCloser
	fname   string
	types   []string
***REMOVED***

func (c *Collector) InTypeList(str string) bool ***REMOVED***
	for _, n := range c.types ***REMOVED***
		if n == str ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func New(u *url.URL) (*Collector, error) ***REMOVED***
	fname := u.Path
	if u.Path == "" ***REMOVED***
		fname = u.String()
	***REMOVED***

	logfile, err := os.Create(fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	t := make([]string, 16)
	return &Collector***REMOVED***
		outfile: logfile,
		fname:   fname,
		types:   t,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) String() string ***REMOVED***
	return "JSON"
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	log.WithField("filename", c.fname).Debug("JSON: Writing JSON metrics")
	<-ctx.Done()
	err := c.outfile.Close()
	if err == nil ***REMOVED***
		return
	***REMOVED***
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	for _, sample := range samples ***REMOVED***
		if !c.InTypeList(sample.Metric.Name) ***REMOVED***
			c.types = append(c.types, sample.Metric.Name)
			if env := Wrap(sample.Metric); env != nil ***REMOVED***
				row, err := json.Marshal(env)
				if err == nil ***REMOVED***
					row = append(row, '\n')
					_, err := c.outfile.Write(row)
					if err != nil ***REMOVED***
						log.WithField("filename", c.fname).Error("JSON: Error writing to file")
					***REMOVED***

				***REMOVED***
			***REMOVED***

		***REMOVED***

		env := Wrap(sample)
		row, err := json.Marshal(env)
		if err != nil || env == nil ***REMOVED***
			// Skip metric if it can't be made into JSON or envelope is null.
			continue
		***REMOVED***
		row = append(row, '\n')
		_, err = c.outfile.Write(row)
		if err != nil ***REMOVED***
			log.WithField("filename", c.fname).Error("JSON: Error writing to file")
		***REMOVED***
	***REMOVED***
***REMOVED***
