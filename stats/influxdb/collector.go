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
	"context"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
)

const (
	pushInterval = 1 * time.Second
)

var _ lib.Collector = &Collector***REMOVED******REMOVED***

type Collector struct ***REMOVED***
	config    Config
	client    client.Client
	batchConf client.BatchPointsConfig

	buffer     []stats.Sample
	bufferLock sync.Mutex
***REMOVED***

func New(conf Config) (*Collector, error) ***REMOVED***
	cl, err := MakeClient(conf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	batchConf := MakeBatchConfig(conf)
	return &Collector***REMOVED***
		config:    conf,
		client:    cl,
		batchConf: batchConf,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) Init() error ***REMOVED***
	// Try to create the database if it doesn't exist. Failure to do so is USUALLY harmless; it
	// usually means we're either a non-admin user to an existing DB or connecting over UDP.
	_, err := c.client.Query(client.NewQuery("CREATE DATABASE "+c.batchConf.Database, "", ""))
	if err != nil ***REMOVED***
		log.WithError(err).Debug("InfluxDB: Couldn't create database; most likely harmless")
	***REMOVED***

	return nil
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	log.Debug("InfluxDB: Running!")
	ticker := time.NewTicker(pushInterval)
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

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	c.bufferLock.Lock()
	c.buffer = append(c.buffer, samples...)
	c.bufferLock.Unlock()
***REMOVED***

func (c *Collector) commit() ***REMOVED***
	c.bufferLock.Lock()
	samples := c.buffer
	c.buffer = nil
	c.bufferLock.Unlock()

	log.Debug("InfluxDB: Committing...")
	batch, err := client.NewBatchPoints(c.batchConf)
	if err != nil ***REMOVED***
		log.WithError(err).Error("InfluxDB: Couldn't make a batch")
		return
	***REMOVED***

	for _, sample := range samples ***REMOVED***
		p, err := client.NewPoint(
			sample.Metric.Name,
			sample.Tags,
			map[string]interface***REMOVED******REMOVED******REMOVED***"value": sample.Value***REMOVED***,
			sample.Time,
		)
		if err != nil ***REMOVED***
			log.WithError(err).Error("InfluxDB: Couldn't make point from sample!")
			return
		***REMOVED***
		batch.AddPoint(p)
	***REMOVED***

	log.WithField("points", len(batch.Points())).Debug("InfluxDB: Writing...")
	startTime := time.Now()
	if err := c.client.Write(batch); err != nil ***REMOVED***
		log.WithError(err).Error("InfluxDB: Couldn't write stats")
	***REMOVED***
	t := time.Since(startTime)
	log.WithField("t", t).Debug("InfluxDB: Batch written!")
***REMOVED***
