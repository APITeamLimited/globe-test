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

package influxdb

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"

	"go.k6.io/k6/output"
	"go.k6.io/k6/stats"
)

// FieldKind defines Enum for tag-to-field type conversion
type FieldKind int

const (
	// String field (default)
	String FieldKind = iota
	// Int field
	Int
	// Float field
	Float
	// Bool field
	Bool
)

// Output is the influxdb Output struct
type Output struct ***REMOVED***
	output.SampleBuffer

	params          output.Params
	periodicFlusher *output.PeriodicFlusher

	Client    client.Client
	Config    Config
	BatchConf client.BatchPointsConfig

	logger      logrus.FieldLogger
	semaphoreCh chan struct***REMOVED******REMOVED***
	fieldKinds  map[string]FieldKind
***REMOVED***

// New returns new influxdb output
func New(params output.Params) (output.Output, error) ***REMOVED***
	return newOutput(params)
***REMOVED***

func newOutput(params output.Params) (*Output, error) ***REMOVED***
	conf, err := GetConsolidatedConfig(params.JSONConfig, params.Environment, params.ConfigArgument)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cl, err := MakeClient(conf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	batchConf := MakeBatchConfig(conf)
	if conf.ConcurrentWrites.Int64 <= 0 ***REMOVED***
		return nil, errors.New("influxdb's ConcurrentWrites must be a positive number")
	***REMOVED***
	fldKinds, err := MakeFieldKinds(conf)
	return &Output***REMOVED***
		params: params,
		logger: params.Logger.WithFields(logrus.Fields***REMOVED***
			"output": "InfluxDBv1",
		***REMOVED***),
		Client:      cl,
		Config:      conf,
		BatchConf:   batchConf,
		semaphoreCh: make(chan struct***REMOVED******REMOVED***, conf.ConcurrentWrites.Int64),
		fieldKinds:  fldKinds,
	***REMOVED***, err
***REMOVED***

func (o *Output) extractTagsToValues(tags map[string]string, values map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	for tag, kind := range o.fieldKinds ***REMOVED***
		if val, ok := tags[tag]; ok ***REMOVED***
			var v interface***REMOVED******REMOVED***
			var err error
			switch kind ***REMOVED***
			case String:
				v = val
			case Bool:
				v, err = strconv.ParseBool(val)
			case Float:
				v, err = strconv.ParseFloat(val, 64)
			case Int:
				v, err = strconv.ParseInt(val, 10, 64)
			***REMOVED***
			if err == nil ***REMOVED***
				values[tag] = v
			***REMOVED*** else ***REMOVED***
				values[tag] = val
			***REMOVED***
			delete(tags, tag)
		***REMOVED***
	***REMOVED***
	return values
***REMOVED***

func (o *Output) batchFromSamples(containers []stats.SampleContainer) (client.BatchPoints, error) ***REMOVED***
	batch, err := client.NewBatchPoints(o.BatchConf)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("couldn't make a batch: %w", err)
	***REMOVED***

	type cacheItem struct ***REMOVED***
		tags   map[string]string
		values map[string]interface***REMOVED******REMOVED***
	***REMOVED***
	cache := map[*stats.SampleTags]cacheItem***REMOVED******REMOVED***
	for _, container := range containers ***REMOVED***
		samples := container.GetSamples()
		for _, sample := range samples ***REMOVED***
			var tags map[string]string
			values := make(map[string]interface***REMOVED******REMOVED***)
			if cached, ok := cache[sample.Tags]; ok ***REMOVED***
				tags = cached.tags
				for k, v := range cached.values ***REMOVED***
					values[k] = v
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				tags = sample.Tags.CloneTags()
				o.extractTagsToValues(tags, values)
				cache[sample.Tags] = cacheItem***REMOVED***tags, values***REMOVED***
			***REMOVED***
			values["value"] = sample.Value
			var p *client.Point
			p, err = client.NewPoint(
				sample.Metric.Name,
				tags,
				values,
				sample.Time,
			)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("couldn't make point from sample: %w", err)
			***REMOVED***
			batch.AddPoint(p)
		***REMOVED***
	***REMOVED***

	return batch, nil
***REMOVED***

// Description returns a human-readable description of the output.
func (o *Output) Description() string ***REMOVED***
	return fmt.Sprintf("InfluxDBv1 (%s)", o.Config.Addr.String)
***REMOVED***

// Start tries to open the specified JSON file and starts the goroutine for
// metric flushing. If gzip encoding is specified, it also handles that.
func (o *Output) Start() error ***REMOVED***
	o.logger.Debug("Starting...")
	// Try to create the database if it doesn't exist. Failure to do so is USUALLY harmless; it
	// usually means we're either a non-admin user to an existing DB or connecting over UDP.
	_, err := o.Client.Query(client.NewQuery("CREATE DATABASE "+o.BatchConf.Database, "", ""))
	if err != nil ***REMOVED***
		o.logger.WithError(err).Debug("InfluxDB: Couldn't create database; most likely harmless")
	***REMOVED***

	pf, err := output.NewPeriodicFlusher(time.Duration(o.Config.PushInterval.Duration), o.flushMetrics)
	if err != nil ***REMOVED***
		return err //nolint:wrapcheck
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
	return nil
***REMOVED***

func (o *Output) flushMetrics() ***REMOVED***
	samples := o.GetBufferedSamples()

	o.semaphoreCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	defer func() ***REMOVED***
		<-o.semaphoreCh
	***REMOVED***()
	o.logger.Debug("Committing...")
	o.logger.WithField("samples", len(samples)).Debug("Writing...")

	batch, err := o.batchFromSamples(samples)
	if err != nil ***REMOVED***
		o.logger.WithError(err).Error("Couldn't create batch from samples")
		return
	***REMOVED***

	o.logger.WithField("points", len(batch.Points())).Debug("Writing...")
	startTime := time.Now()
	if err := o.Client.Write(batch); err != nil ***REMOVED***
		o.logger.WithError(err).Error("Couldn't write stats")
	***REMOVED***
	t := time.Since(startTime)
	o.logger.WithField("t", t).Debug("Batch written!")
***REMOVED***
