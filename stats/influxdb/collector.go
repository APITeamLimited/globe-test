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
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	null "gopkg.in/guregu/null.v3"
)

const (
	pushInterval = 1 * time.Second

	defaultURL = "http://localhost:8086/k6"
)

var _ lib.AuthenticatedCollector = &Collector***REMOVED******REMOVED***

type Config struct ***REMOVED***
	DefaultURL null.String `json:"default_url,omitempty"`
***REMOVED***

type Collector struct ***REMOVED***
	u          *url.URL
	client     client.Client
	batchConf  client.BatchPointsConfig
	buffer     []stats.Sample
	bufferLock sync.Mutex
***REMOVED***

func New(s string, conf_ interface***REMOVED******REMOVED***, opts lib.Options) (*Collector, error) ***REMOVED***
	conf := conf_.(*Config)

	if s == "" ***REMOVED***
		s = conf.DefaultURL.String
	***REMOVED***
	if s == "" ***REMOVED***
		s = defaultURL
	***REMOVED***

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

func (c *Collector) Init() error ***REMOVED***
	// Try to create the database if it doesn't exist. Failure to do so is USUALLY harmless; it
	// usually means we're either a non-admin user to an existing DB or connecting over UDP.
	_, err := c.client.Query(client.NewQuery("CREATE DATABASE "+c.batchConf.Database, "", ""))
	if err != nil ***REMOVED***
		log.WithError(err).Debug("InfluxDB: Couldn't create database; most likely harmless")
	***REMOVED***

	return nil
***REMOVED***

func (c *Collector) MakeConfig() interface***REMOVED******REMOVED*** ***REMOVED***
	return &Config***REMOVED******REMOVED***
***REMOVED***

func (c *Collector) Login(conf_ interface***REMOVED******REMOVED***, in io.Reader, out io.Writer) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	conf := conf_.(*Config)

	buf := bufio.NewReader(in)
	fmt.Fprintf(out, "default connection url: ")
	ustr, err := buf.ReadString('\n')
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ustr = strings.TrimSpace(ustr)

	u, err := url.Parse(ustr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cl, _, err := parseURL(u)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, _, err := cl.Ping(5 * time.Second); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	conf.DefaultURL = null.StringFrom(u.String())

	return conf, nil
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
