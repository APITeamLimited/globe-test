/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package common

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
)

// Config is the common configuration interface for StatsD/Datadog.
type Config interface ***REMOVED***
	GetAddr() null.String
	GetBufferSize() null.Int
	GetNamespace() null.String
	GetPushInterval() types.NullDuration
***REMOVED***

var _ lib.Collector = &Collector***REMOVED******REMOVED***

// Collector sends result data to statsd daemons with the ability to send to datadog as well
type Collector struct ***REMOVED***
	Config Config
	Type   string
	// ProcessTags is called on a map of all tags for each metric and returns a slice representation
	// of those tags that should be sent. No tags are send in case of ProcessTags being null
	ProcessTags func(map[string]string) []string

	Logger     logrus.FieldLogger
	client     *statsd.Client
	startTime  time.Time
	buffer     []*Sample
	bufferLock sync.Mutex
***REMOVED***

// Init sets up the collector
func (c *Collector) Init() (err error) ***REMOVED***
	c.Logger = c.Logger.WithField("type", c.Type)
	if address := c.Config.GetAddr().String; address == "" ***REMOVED***
		err = fmt.Errorf(
			"connection string is invalid. Received: \"%+s\"",
			address,
		)
		c.Logger.Error(err)

		return err
	***REMOVED***

	c.client, err = statsd.NewBuffered(c.Config.GetAddr().String, int(c.Config.GetBufferSize().Int64))

	if err != nil ***REMOVED***
		c.Logger.Errorf("Couldn't make buffered client, %s", err)
		return err
	***REMOVED***

	if namespace := c.Config.GetNamespace().String; namespace != "" ***REMOVED***
		c.client.Namespace = namespace
	***REMOVED***

	return nil
***REMOVED***

// Link returns the address of the client
func (c *Collector) Link() string ***REMOVED***
	return c.Config.GetAddr().String
***REMOVED***

// Run the collector
func (c *Collector) Run(ctx context.Context) ***REMOVED***
	c.Logger.Debugf("%s: Running!", c.Type)
	ticker := time.NewTicker(time.Duration(c.Config.GetPushInterval().Duration))
	c.startTime = time.Now()

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			c.pushMetrics()
		case <-ctx.Done():
			c.pushMetrics()
			c.finish()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Collect metrics
func (c *Collector) Collect(containers []stats.SampleContainer) ***REMOVED***
	var pointSamples []*Sample

	for _, container := range containers ***REMOVED***
		for _, sample := range container.GetSamples() ***REMOVED***
			pointSamples = append(pointSamples, generateDataPoint(sample))
		***REMOVED***
	***REMOVED***

	if len(pointSamples) > 0 ***REMOVED***
		c.bufferLock.Lock()
		c.buffer = append(c.buffer, pointSamples...)
		c.bufferLock.Unlock()
	***REMOVED***
***REMOVED***

func (c *Collector) pushMetrics() ***REMOVED***
	c.bufferLock.Lock()
	if len(c.buffer) == 0 ***REMOVED***
		c.bufferLock.Unlock()
		return
	***REMOVED***
	buffer := c.buffer
	c.buffer = nil
	c.bufferLock.Unlock()

	c.Logger.
		WithField("samples", len(buffer)).
		Debug("Pushing metrics to server")

	if err := c.commit(buffer); err != nil ***REMOVED***
		c.Logger.
			WithError(err).
			Error("Couldn't commit a batch")
	***REMOVED***
***REMOVED***

func (c *Collector) finish() ***REMOVED***
	// Close when context is done
	if err := c.client.Close(); err != nil ***REMOVED***
		c.Logger.Warnf("Error closing the client, %+v", err)
	***REMOVED***
***REMOVED***

func (c *Collector) commit(data []*Sample) error ***REMOVED***
	var errorCount int
	for _, entry := range data ***REMOVED***
		if err := c.dispatch(entry); err != nil ***REMOVED***
			// No need to return error if just one metric didn't go through
			c.Logger.WithError(err).Debugf("Error while sending metric %s", entry.Metric)
			errorCount++
		***REMOVED***
	***REMOVED***
	if errorCount != 0 ***REMOVED***
		c.Logger.Warnf("Couldn't send %d out of %d metrics. Enable debug logging to see individual errors",
			errorCount, len(data))
	***REMOVED***
	return c.client.Flush()
***REMOVED***

func (c *Collector) dispatch(entry *Sample) error ***REMOVED***
	var tagList []string
	if c.ProcessTags != nil ***REMOVED***
		tagList = c.ProcessTags(entry.Tags)
	***REMOVED***

	switch entry.Type ***REMOVED***
	case stats.Counter:
		return c.client.Count(entry.Metric, int64(entry.Value), tagList, 1)
	case stats.Trend:
		return c.client.TimeInMilliseconds(entry.Metric, entry.Value, tagList, 1)
	case stats.Gauge:
		return c.client.Gauge(entry.Metric, entry.Value, tagList, 1)
	case stats.Rate:
		if check := entry.Tags["check"]; check != "" ***REMOVED***
			return c.client.Count(
				checkToString(check, entry.Value),
				1,
				tagList,
				1,
			)
		***REMOVED***
		return c.client.Count(entry.Metric, int64(entry.Value), tagList, 1)
	default:
		return fmt.Errorf("unsupported metric type %s", entry.Type)
	***REMOVED***
***REMOVED***

func checkToString(check string, value float64) string ***REMOVED***
	label := "pass"
	if value == 0 ***REMOVED***
		label = "fail"
	***REMOVED***
	return "check." + check + "." + label
***REMOVED***
