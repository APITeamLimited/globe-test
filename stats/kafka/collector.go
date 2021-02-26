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

package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"

	"github.com/loadimpact/k6/lib"
	jsono "github.com/loadimpact/k6/output/json"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/influxdb"
)

// Collector implements the lib.Collector interface and should be used only for testing
type Collector struct ***REMOVED***
	Producer sarama.SyncProducer
	Config   Config

	Samples []stats.Sample
	logger  logrus.FieldLogger
	lock    sync.Mutex
***REMOVED***

// New creates an instance of the collector
func New(logger logrus.FieldLogger, conf Config) (*Collector, error) ***REMOVED***
	producer, err := sarama.NewSyncProducer(conf.Brokers, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Collector***REMOVED***
		Producer: producer,
		Config:   conf,
		logger:   logger,
	***REMOVED***, nil
***REMOVED***

// Init does nothing, it's only included to satisfy the lib.Collector interface
func (c *Collector) Init() error ***REMOVED*** return nil ***REMOVED***

// Run just blocks until the context is done
func (c *Collector) Run(ctx context.Context) ***REMOVED***
	c.logger.Debug("Kafka: Running!")
	ticker := time.NewTicker(time.Duration(c.Config.PushInterval.Duration))
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			c.pushMetrics()
		case <-ctx.Done():
			c.pushMetrics()

			err := c.Producer.Close()
			if err != nil ***REMOVED***
				c.logger.WithError(err).Error("Kafka: Failed to close producer.")
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Collect just appends all of the samples passed to it to the internal sample slice.
// According to the the lib.Collector interface, it should never be called concurrently,
// so there's no locking on purpose - that way Go's race condition detector can actually
// detect incorrect usage.
// Also, theoretically the collector doesn't have to actually Run() before samples start
// being collected, it only has to be initialized.
func (c *Collector) Collect(scs []stats.SampleContainer) ***REMOVED***
	c.lock.Lock()
	for _, sc := range scs ***REMOVED***
		c.Samples = append(c.Samples, sc.GetSamples()...)
	***REMOVED***
	c.lock.Unlock()
***REMOVED***

// Link returns a dummy string, it's only included to satisfy the lib.Collector interface
func (c *Collector) Link() string ***REMOVED***
	return ""
***REMOVED***

// GetRequiredSystemTags returns which sample tags are needed by this collector
func (c *Collector) GetRequiredSystemTags() stats.SystemTagSet ***REMOVED***
	return stats.SystemTagSet(0) // There are no required tags for this collector
***REMOVED***

// SetRunStatus does nothing in the Kafka collector
func (c *Collector) SetRunStatus(status lib.RunStatus) ***REMOVED******REMOVED***

func (c *Collector) formatSamples(samples stats.Samples) ([]string, error) ***REMOVED***
	var metrics []string

	switch c.Config.Format.String ***REMOVED***
	case "influxdb":
		i, err := influxdb.New(c.logger, c.Config.InfluxDBConfig)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		metrics, err = i.Format(samples)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	default:
		for _, sample := range samples ***REMOVED***
			env := jsono.WrapSample(sample)
			metric, err := json.Marshal(env)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			metrics = append(metrics, string(metric))
		***REMOVED***
	***REMOVED***

	return metrics, nil
***REMOVED***

func (c *Collector) pushMetrics() ***REMOVED***
	startTime := time.Now()

	c.lock.Lock()
	samples := c.Samples
	c.Samples = nil
	c.lock.Unlock()

	// Format the samples
	formattedSamples, err := c.formatSamples(samples)
	if err != nil ***REMOVED***
		c.logger.WithError(err).Error("Kafka: Couldn't format the samples")
		return
	***REMOVED***

	// Send the samples
	c.logger.Debug("Kafka: Delivering...")

	for _, sample := range formattedSamples ***REMOVED***
		msg := &sarama.ProducerMessage***REMOVED***Topic: c.Config.Topic.String, Value: sarama.StringEncoder(sample)***REMOVED***
		partition, offset, err := c.Producer.SendMessage(msg)
		if err != nil ***REMOVED***
			c.logger.WithError(err).Error("Kafka: failed to send message.")
		***REMOVED*** else ***REMOVED***
			c.logger.WithFields(logrus.Fields***REMOVED***
				"partition": partition,
				"offset":    offset,
			***REMOVED***).Debug("Kafka: message sent.")
		***REMOVED***
	***REMOVED***

	t := time.Since(startTime)
	c.logger.WithField("t", t).Debug("Kafka: Delivered!")
***REMOVED***
