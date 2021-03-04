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

package csv

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

// Collector saving output to csv implements the lib.Collector interface
type Collector struct ***REMOVED***
	closeFn      func() error
	fname        string
	resTags      []string
	ignoredTags  []string
	csvWriter    *csv.Writer
	csvLock      sync.Mutex
	buffer       []stats.Sample
	bufferLock   sync.Mutex
	row          []string
	saveInterval time.Duration
	logger       logrus.FieldLogger
***REMOVED***

// Verify that Collector implements lib.Collector
var _ lib.Collector = &Collector***REMOVED******REMOVED***

// New Creates new instance of CSV collector
func New(logger logrus.FieldLogger, fs afero.Fs, tags stats.TagSet, config Config) (*Collector, error) ***REMOVED***
	resTags := []string***REMOVED******REMOVED***
	ignoredTags := []string***REMOVED******REMOVED***
	for tag, flag := range tags ***REMOVED***
		if flag ***REMOVED***
			resTags = append(resTags, tag)
		***REMOVED*** else ***REMOVED***
			ignoredTags = append(ignoredTags, tag)
		***REMOVED***
	***REMOVED***
	sort.Strings(resTags)
	sort.Strings(ignoredTags)

	saveInterval := time.Duration(config.SaveInterval.Duration)
	fname := config.FileName.String

	if fname == "" || fname == "-" ***REMOVED***
		stdoutWriter := csv.NewWriter(os.Stdout)
		return &Collector***REMOVED***
			fname:        "-",
			resTags:      resTags,
			ignoredTags:  ignoredTags,
			csvWriter:    stdoutWriter,
			row:          make([]string, 3+len(resTags)+1),
			saveInterval: saveInterval,
			closeFn:      func() error ***REMOVED*** return nil ***REMOVED***,
			logger:       logger,
		***REMOVED***, nil
	***REMOVED***

	logFile, err := fs.Create(fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c := Collector***REMOVED***
		fname:        fname,
		resTags:      resTags,
		ignoredTags:  ignoredTags,
		row:          make([]string, 3+len(resTags)+1),
		saveInterval: saveInterval,
		logger:       logger,
	***REMOVED***

	if strings.HasSuffix(fname, ".gz") ***REMOVED***
		outfile := gzip.NewWriter(logFile)
		csvWriter := csv.NewWriter(outfile)
		c.csvWriter = csvWriter
		c.closeFn = func() error ***REMOVED***
			_ = outfile.Close()
			return logFile.Close()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		csvWriter := csv.NewWriter(logFile)
		c.csvWriter = csvWriter
		c.closeFn = logFile.Close
	***REMOVED***

	return &c, nil
***REMOVED***

// Init writes column names to csv file
func (c *Collector) Init() error ***REMOVED***
	header := MakeHeader(c.resTags)
	err := c.csvWriter.Write(header)
	if err != nil ***REMOVED***
		c.logger.WithField("filename", c.fname).Error("CSV: Error writing column names to file")
	***REMOVED***
	c.csvWriter.Flush()
	return nil
***REMOVED***

// Run just blocks until the context is done
func (c *Collector) Run(ctx context.Context) ***REMOVED***
	ticker := time.NewTicker(c.saveInterval)
	defer func() ***REMOVED***
		err := c.closeFn()
		if err != nil ***REMOVED***
			c.logger.WithField("filename", c.fname).Errorf("CSV: Error closing the file: %v", err)
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			c.writeToFile()
		case <-ctx.Done():
			c.writeToFile()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Collect Saves samples to buffer
func (c *Collector) Collect(scs []stats.SampleContainer) ***REMOVED***
	c.bufferLock.Lock()
	defer c.bufferLock.Unlock()
	for _, sc := range scs ***REMOVED***
		c.buffer = append(c.buffer, sc.GetSamples()...)
	***REMOVED***
***REMOVED***

// writeToFile Writes samples to the csv file
func (c *Collector) writeToFile() ***REMOVED***
	c.bufferLock.Lock()
	samples := c.buffer
	c.buffer = nil
	c.bufferLock.Unlock()

	if len(samples) > 0 ***REMOVED***
		c.csvLock.Lock()
		defer c.csvLock.Unlock()
		for _, sc := range samples ***REMOVED***
			for _, sample := range sc.GetSamples() ***REMOVED***
				sample := sample
				row := SampleToRow(&sample, c.resTags, c.ignoredTags, c.row)
				err := c.csvWriter.Write(row)
				if err != nil ***REMOVED***
					c.logger.WithField("filename", c.fname).Error("CSV: Error writing to file")
				***REMOVED***
			***REMOVED***
		***REMOVED***
		c.csvWriter.Flush()
	***REMOVED***
***REMOVED***

// Link returns a dummy string, it's only included to satisfy the lib.Collector interface
func (c *Collector) Link() string ***REMOVED***
	return c.fname
***REMOVED***

// MakeHeader creates list of column names for csv file
func MakeHeader(tags []string) []string ***REMOVED***
	tags = append(tags, "extra_tags")
	return append([]string***REMOVED***"metric_name", "timestamp", "metric_value"***REMOVED***, tags...)
***REMOVED***

// SampleToRow converts sample into array of strings
func SampleToRow(sample *stats.Sample, resTags []string, ignoredTags []string, row []string) []string ***REMOVED***
	row[0] = sample.Metric.Name
	row[1] = fmt.Sprintf("%d", sample.Time.Unix())
	row[2] = fmt.Sprintf("%f", sample.Value)
	sampleTags := sample.Tags.CloneTags()

	for ind, tag := range resTags ***REMOVED***
		row[ind+3] = sampleTags[tag]
	***REMOVED***

	extraTags := bytes.Buffer***REMOVED******REMOVED***
	prev := false
	for tag, val := range sampleTags ***REMOVED***
		if !IsStringInSlice(resTags, tag) && !IsStringInSlice(ignoredTags, tag) ***REMOVED***
			if prev ***REMOVED***
				if _, err := extraTags.WriteString("&"); err != nil ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			if _, err := extraTags.WriteString(tag); err != nil ***REMOVED***
				break
			***REMOVED***

			if _, err := extraTags.WriteString("="); err != nil ***REMOVED***
				break
			***REMOVED***

			if _, err := extraTags.WriteString(val); err != nil ***REMOVED***
				break
			***REMOVED***
			prev = true
		***REMOVED***
	***REMOVED***
	row[len(row)-1] = extraTags.String()

	return row
***REMOVED***

// IsStringInSlice returns whether the string is contained within a string slice
func IsStringInSlice(slice []string, str string) bool ***REMOVED***
	if index := sort.SearchStrings(slice, str); index == len(slice) || slice[index] != str ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***
