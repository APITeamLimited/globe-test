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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"net/url"
	"time"
)

const pushInterval = 1 * time.Second

type Collector struct ***REMOVED***
	u         *url.URL
	client    client.Client
	batchConf client.BatchPointsConfig
	buffer    []stats.Sample
***REMOVED***

func New(s string, opts lib.Options) (*Collector, error) ***REMOVED***
	u, err := url.Parse(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cl, batchConf, err := parseURL(u)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Collector***REMOVED***
		u:         u,
		client:    cl,
		batchConf: batchConf,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) String() string ***REMOVED***
	return fmt.Sprintf("influxdb (%s)", c.u.Host)
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
	c.buffer = append(c.buffer, samples...)
***REMOVED***

func (c *Collector) commit() ***REMOVED***
	samples := c.buffer
	c.buffer = nil

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
