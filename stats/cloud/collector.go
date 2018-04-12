/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package cloud

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

const (
	TestName           = "k6 test"
	MetricPushInterval = 1 * time.Second
)

// Collector sends result data to the Load Impact cloud service.
type Collector struct ***REMOVED***
	config      Config
	referenceID string

	duration   int64
	thresholds map[string][]*stats.Threshold
	client     *Client

	anonymous bool

	sampleBuffer []*Sample
	sampleMu     sync.Mutex
***REMOVED***

// New creates a new cloud collector
func New(conf Config, src *lib.SourceData, opts lib.Options, version string) (*Collector, error) ***REMOVED***
	if val, ok := opts.External["loadimpact"]; ok ***REMOVED***
		if err := mapstructure.Decode(val, &conf); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if conf.Name == "" ***REMOVED***
		conf.Name = filepath.Base(src.Filename)
	***REMOVED***
	if conf.Name == "-" ***REMOVED***
		conf.Name = TestName
	***REMOVED***

	thresholds := make(map[string][]*stats.Threshold)
	for name, t := range opts.Thresholds ***REMOVED***
		thresholds[name] = append(thresholds[name], t.Thresholds...)
	***REMOVED***

	// Sum test duration from options. -1 for unknown duration.
	var duration int64 = -1
	if len(opts.Stages) > 0 ***REMOVED***
		duration = sumStages(opts.Stages)
	***REMOVED*** else if opts.Duration.Valid ***REMOVED***
		duration = int64(time.Duration(opts.Duration.Duration).Seconds())
	***REMOVED***

	if conf.Token == "" && conf.DeprecatedToken != "" ***REMOVED***
		log.Warn("K6CLOUD_TOKEN is deprecated and will be removed. Use K6_CLOUD_TOKEN instead.")
		conf.Token = conf.DeprecatedToken
	***REMOVED***

	return &Collector***REMOVED***
		config:     conf,
		thresholds: thresholds,
		client:     NewClient(conf.Token, conf.Host, version),
		anonymous:  conf.Token == "",
		duration:   duration,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) Init() error ***REMOVED***
	thresholds := make(map[string][]string)

	for name, t := range c.thresholds ***REMOVED***
		for _, threshold := range t ***REMOVED***
			thresholds[name] = append(thresholds[name], threshold.Source)
		***REMOVED***
	***REMOVED***

	testRun := &TestRun***REMOVED***
		Name:       c.config.Name,
		ProjectID:  c.config.ProjectID,
		Thresholds: thresholds,
		Duration:   c.duration,
	***REMOVED***

	response, err := c.client.CreateTestRun(testRun)

	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.referenceID = response.ReferenceID

	log.WithFields(log.Fields***REMOVED***
		"name":        c.config.Name,
		"projectId":   c.config.ProjectID,
		"duration":    c.duration,
		"referenceId": c.referenceID,
	***REMOVED***).Debug("Cloud: Initialized")
	return nil
***REMOVED***

func (c *Collector) Link() string ***REMOVED***
	return URLForResults(c.referenceID, c.config)
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	timer := time.NewTicker(MetricPushInterval)

	for ***REMOVED***
		select ***REMOVED***
		case <-timer.C:
			c.pushMetrics()
		case <-ctx.Done():
			c.pushMetrics()
			c.testFinished()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) IsReady() bool ***REMOVED***
	return true
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	if c.referenceID == "" ***REMOVED***
		return
	***REMOVED***

	var cloudSamples []*Sample
	var httpJSON *Sample
	var iterationJSON *Sample
	for _, samp := range samples ***REMOVED***

		name := samp.Metric.Name
		if name == "http_reqs" ***REMOVED***
			httpJSON = &Sample***REMOVED***
				Type:   "Points",
				Metric: "http_req_li_all",
				Data: SampleData***REMOVED***
					Type:   samp.Metric.Type,
					Time:   samp.Time,
					Tags:   samp.Tags,
					Values: make(map[string]float64),
				***REMOVED***,
			***REMOVED***
			httpJSON.Data.Values[name] = samp.Value
			cloudSamples = append(cloudSamples, httpJSON)
		***REMOVED*** else if name == "data_sent" ***REMOVED***
			iterationJSON = &Sample***REMOVED***
				Type:   "Points",
				Metric: "iter_li_all",
				Data: SampleData***REMOVED***
					Type:   samp.Metric.Type,
					Time:   samp.Time,
					Tags:   samp.Tags,
					Values: make(map[string]float64),
				***REMOVED***,
			***REMOVED***
			iterationJSON.Data.Values[name] = samp.Value
			cloudSamples = append(cloudSamples, iterationJSON)
		***REMOVED*** else if name == "data_received" || name == "iteration_duration" ***REMOVED***
			//TODO: make sure that tags match
			iterationJSON.Data.Values[name] = samp.Value
		***REMOVED*** else if strings.HasPrefix(name, "http_req_") ***REMOVED***
			//TODO: make sure that tags match
			httpJSON.Data.Values[name] = samp.Value
		***REMOVED*** else ***REMOVED***
			sampleJSON := &Sample***REMOVED***
				Type:   "Point",
				Metric: name,
				Data: SampleData***REMOVED***
					Type:  samp.Metric.Type,
					Time:  samp.Time,
					Value: samp.Value,
					Tags:  samp.Tags,
				***REMOVED***,
			***REMOVED***
			cloudSamples = append(cloudSamples, sampleJSON)
		***REMOVED***
	***REMOVED***

	if len(cloudSamples) > 0 ***REMOVED***
		c.sampleMu.Lock()
		c.sampleBuffer = append(c.sampleBuffer, cloudSamples...)
		c.sampleMu.Unlock()
	***REMOVED***
***REMOVED***

func (c *Collector) pushMetrics() ***REMOVED***
	c.sampleMu.Lock()
	if len(c.sampleBuffer) == 0 ***REMOVED***
		c.sampleMu.Unlock()
		return
	***REMOVED***
	buffer := c.sampleBuffer
	c.sampleBuffer = nil
	c.sampleMu.Unlock()

	log.WithFields(log.Fields***REMOVED***
		"samples": len(buffer),
	***REMOVED***).Debug("Pushing metrics to cloud")

	err := c.client.PushMetric(c.referenceID, c.config.NoCompress, buffer)
	if err != nil ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"error": err,
		***REMOVED***).Warn("Failed to send metrics to cloud")
	***REMOVED***
***REMOVED***

func (c *Collector) testFinished() ***REMOVED***
	if c.referenceID == "" ***REMOVED***
		return
	***REMOVED***

	testTainted := false
	thresholdResults := make(ThresholdResult)
	for name, thresholds := range c.thresholds ***REMOVED***
		thresholdResults[name] = make(map[string]bool)
		for _, t := range thresholds ***REMOVED***
			thresholdResults[name][t.Source] = t.Failed
			if t.Failed ***REMOVED***
				testTainted = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"ref":     c.referenceID,
		"tainted": testTainted,
	***REMOVED***).Debug("Sending test finished")

	err := c.client.TestFinished(c.referenceID, thresholdResults, testTainted)
	if err != nil ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"error": err,
		***REMOVED***).Warn("Failed to send test finished to cloud")
	***REMOVED***
***REMOVED***

func sumStages(stages []lib.Stage) int64 ***REMOVED***
	var total time.Duration
	for _, stage := range stages ***REMOVED***
		total += time.Duration(stage.Duration.Duration)
	***REMOVED***

	return int64(total.Seconds())
***REMOVED***

// GetRequiredSystemTags returns which sample tags are needed by this collector
func (c *Collector) GetRequiredSystemTags() lib.TagSet ***REMOVED***
	return lib.GetTagSet("name", "method", "status", "error", "check", "group")
***REMOVED***
