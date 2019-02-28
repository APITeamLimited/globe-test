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
	"encoding/json"
	"path/filepath"
	"sync"
	"time"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/pkg/errors"

	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
)

// TestName is the default Load Impact Cloud test name
const TestName = "k6 test"

// Collector sends result data to the Load Impact cloud service.
type Collector struct ***REMOVED***
	config      Config
	referenceID string

	duration   int64
	thresholds map[string][]*stats.Threshold
	client     *Client

	anonymous bool
	runStatus lib.RunStatus

	bufferMutex      sync.Mutex
	bufferHTTPTrails []*httpext.Trail
	bufferSamples    []*Sample

	opts lib.Options

	// TODO: optimize this
	//
	// Since the real-time metrics refactoring (https://github.com/loadimpact/k6/pull/678),
	// we should no longer have to handle metrics that have times long in the past. So instead of a
	// map, we can probably use a simple slice (or even an array!) as a ring buffer to store the
	// aggregation buckets. This should save us a some time, since it would make the lookups and WaitPeriod
	// checks basically O(1). And even if for some reason there are occasional metrics with past times that
	// don't fit in the chosen ring buffer size, we could just send them along to the buffer unaggregated
	aggrBuckets map[int64]aggregationBucket
***REMOVED***

// Verify that Collector implements lib.Collector
var _ lib.Collector = &Collector***REMOVED******REMOVED***

// MergeFromExternal merges three fields from json in a loadimact key of the provided external map
func MergeFromExternal(external map[string]json.RawMessage, conf *Config) error ***REMOVED***
	if val, ok := external["loadimpact"]; ok ***REMOVED***
		// TODO: Important! Separate configs and fix the whole 2 configs mess!
		tmpConfig := Config***REMOVED******REMOVED***
		if err := json.Unmarshal(val, &tmpConfig); err != nil ***REMOVED***
			return err
		***REMOVED***
		// Only take out the ProjectID, Name and Token from the options.ext.loadimpact map:
		if tmpConfig.ProjectID.Valid ***REMOVED***
			conf.ProjectID = tmpConfig.ProjectID
		***REMOVED***
		if tmpConfig.Name.Valid ***REMOVED***
			conf.Name = tmpConfig.Name
		***REMOVED***
		if tmpConfig.Token.Valid ***REMOVED***
			conf.Token = tmpConfig.Token
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// New creates a new cloud collector
func New(conf Config, src *lib.SourceData, opts lib.Options, version string) (*Collector, error) ***REMOVED***
	if err := MergeFromExternal(opts.External, &conf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if conf.AggregationPeriod.Duration > 0 && (opts.SystemTags["vu"] || opts.SystemTags["iter"]) ***REMOVED***
		return nil, errors.New("Aggregation cannot be enabled if the 'vu' or 'iter' system tag is also enabled")
	***REMOVED***

	if !conf.Name.Valid || conf.Name.String == "" ***REMOVED***
		conf.Name = null.StringFrom(filepath.Base(src.Filename))
	***REMOVED***
	if conf.Name.String == "-" ***REMOVED***
		conf.Name = null.StringFrom(TestName)
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

	if duration == -1 ***REMOVED***
		return nil, errors.New("Tests with unspecified duration are not allowed when using Load Impact Insights")
	***REMOVED***

	if !conf.Token.Valid && conf.DeprecatedToken.Valid ***REMOVED***
		log.Warn("K6CLOUD_TOKEN is deprecated and will be removed. Use K6_CLOUD_TOKEN instead.")
		conf.Token = conf.DeprecatedToken
	***REMOVED***

	return &Collector***REMOVED***
		config:      conf,
		thresholds:  thresholds,
		client:      NewClient(conf.Token.String, conf.Host.String, version),
		anonymous:   !conf.Token.Valid,
		duration:    duration,
		opts:        opts,
		aggrBuckets: map[int64]aggregationBucket***REMOVED******REMOVED***,
	***REMOVED***, nil
***REMOVED***

// Init is called between the collector's creation and the call to Run().
// You should do any lengthy setup here rather than in New.
func (c *Collector) Init() error ***REMOVED***
	thresholds := make(map[string][]string)

	for name, t := range c.thresholds ***REMOVED***
		for _, threshold := range t ***REMOVED***
			thresholds[name] = append(thresholds[name], threshold.Source)
		***REMOVED***
	***REMOVED***

	testRun := &TestRun***REMOVED***
		Name:       c.config.Name.String,
		ProjectID:  c.config.ProjectID.Int64,
		VUsMax:     c.opts.VUsMax.Int64,
		Thresholds: thresholds,
		Duration:   c.duration,
	***REMOVED***

	response, err := c.client.CreateTestRun(testRun)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.referenceID = response.ReferenceID

	if response.ConfigOverride != nil ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"override": response.ConfigOverride,
		***REMOVED***).Debug("Cloud: overriding config options")
		c.config = c.config.Apply(*response.ConfigOverride)
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"name":        c.config.Name,
		"projectId":   c.config.ProjectID,
		"duration":    c.duration,
		"referenceId": c.referenceID,
	***REMOVED***).Debug("Cloud: Initialized")
	return nil
***REMOVED***

// Link return a link that is shown to the user.
func (c *Collector) Link() string ***REMOVED***
	return URLForResults(c.referenceID, c.config)
***REMOVED***

// Run is called in a goroutine and starts the collector. Should commit samples to the backend
// at regular intervals and when the context is terminated.
func (c *Collector) Run(ctx context.Context) ***REMOVED***
	wg := sync.WaitGroup***REMOVED******REMOVED***

	// If enabled, start periodically aggregating the collected HTTP trails
	if c.config.AggregationPeriod.Duration > 0 ***REMOVED***
		wg.Add(1)
		aggregationTicker := time.NewTicker(time.Duration(c.config.AggregationCalcInterval.Duration))

		go func() ***REMOVED***
			for ***REMOVED***
				select ***REMOVED***
				case <-aggregationTicker.C:
					c.aggregateHTTPTrails(time.Duration(c.config.AggregationWaitPeriod.Duration))
				case <-ctx.Done():
					c.aggregateHTTPTrails(0)
					c.flushHTTPTrails()
					c.pushMetrics()
					wg.Done()
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	defer func() ***REMOVED***
		wg.Wait()
		c.testFinished()
	***REMOVED***()

	pushTicker := time.NewTicker(time.Duration(c.config.MetricPushInterval.Duration))
	for ***REMOVED***
		select ***REMOVED***
		case <-pushTicker.C:
			c.pushMetrics()
		case <-ctx.Done():
			c.pushMetrics()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Collect receives a set of samples. This method is never called concurrently, and only while
// the context for Run() is valid, but should defer as much work as possible to Run().
func (c *Collector) Collect(sampleContainers []stats.SampleContainer) ***REMOVED***
	if c.referenceID == "" ***REMOVED***
		return
	***REMOVED***

	newSamples := []*Sample***REMOVED******REMOVED***
	newHTTPTrails := []*httpext.Trail***REMOVED******REMOVED***

	for _, sampleContainer := range sampleContainers ***REMOVED***
		switch sc := sampleContainer.(type) ***REMOVED***
		case *httpext.Trail:
			// Check if aggregation is enabled,
			if c.config.AggregationPeriod.Duration > 0 ***REMOVED***
				newHTTPTrails = append(newHTTPTrails, sc)
			***REMOVED*** else ***REMOVED***
				newSamples = append(newSamples, NewSampleFromTrail(sc))
			***REMOVED***
		case *netext.NetTrail:
			//TODO: aggregate?
			values := map[string]float64***REMOVED***
				metrics.DataSent.Name:     float64(sc.BytesWritten),
				metrics.DataReceived.Name: float64(sc.BytesRead),
			***REMOVED***

			if sc.FullIteration ***REMOVED***
				values[metrics.IterationDuration.Name] = stats.D(sc.EndTime.Sub(sc.StartTime))
			***REMOVED***

			newSamples = append(newSamples, &Sample***REMOVED***
				Type:   DataTypeMap,
				Metric: "iter_li_all",
				Data: &SampleDataMap***REMOVED***
					Time:   Timestamp(sc.GetTime()),
					Tags:   sc.GetTags(),
					Values: values,
				***REMOVED******REMOVED***)
		default:
			for _, sample := range sampleContainer.GetSamples() ***REMOVED***
				newSamples = append(newSamples, &Sample***REMOVED***
					Type:   DataTypeSingle,
					Metric: sample.Metric.Name,
					Data: &SampleDataSingle***REMOVED***
						Type:  sample.Metric.Type,
						Time:  Timestamp(sample.Time),
						Tags:  sample.Tags,
						Value: sample.Value,
					***REMOVED***,
				***REMOVED***)
			***REMOVED***

		***REMOVED***
	***REMOVED***

	if len(newSamples) > 0 || len(newHTTPTrails) > 0 ***REMOVED***
		c.bufferMutex.Lock()
		c.bufferSamples = append(c.bufferSamples, newSamples...)
		c.bufferHTTPTrails = append(c.bufferHTTPTrails, newHTTPTrails...)
		c.bufferMutex.Unlock()
	***REMOVED***
***REMOVED***

func (c *Collector) aggregateHTTPTrails(waitPeriod time.Duration) ***REMOVED***
	c.bufferMutex.Lock()
	newHTTPTrails := c.bufferHTTPTrails
	c.bufferHTTPTrails = nil
	c.bufferMutex.Unlock()

	aggrPeriod := int64(c.config.AggregationPeriod.Duration)

	// Distribute all newly buffered HTTP trails into buckets and sub-buckets
	for _, trail := range newHTTPTrails ***REMOVED***
		trailTags := trail.GetTags()
		bucketID := trail.GetTime().UnixNano() / aggrPeriod

		// Get or create a time bucket for that trail period
		bucket, ok := c.aggrBuckets[bucketID]
		if !ok ***REMOVED***
			bucket = aggregationBucket***REMOVED******REMOVED***
			c.aggrBuckets[bucketID] = bucket
		***REMOVED***

		// Either use an existing subbucket key or use the trail tags as a new one
		subBucketKey := trailTags
		subBucket, ok := bucket[subBucketKey]
		if !ok ***REMOVED***
			for sbTags, sb := range bucket ***REMOVED***
				if trailTags.IsEqual(sbTags) ***REMOVED***
					subBucketKey = sbTags
					subBucket = sb
				***REMOVED***
			***REMOVED***
		***REMOVED***
		bucket[subBucketKey] = append(subBucket, trail)
	***REMOVED***

	// Which buckets are still new and we'll wait for trails to accumulate before aggregating
	bucketCutoffID := time.Now().Add(-waitPeriod).UnixNano() / aggrPeriod
	iqrRadius := c.config.AggregationOutlierIqrRadius.Float64
	iqrLowerCoef := c.config.AggregationOutlierIqrCoefLower.Float64
	iqrUpperCoef := c.config.AggregationOutlierIqrCoefUpper.Float64
	newSamples := []*Sample***REMOVED******REMOVED***

	// Handle all aggregation buckets older than bucketCutoffID
	for bucketID, subBuckets := range c.aggrBuckets ***REMOVED***
		if bucketID > bucketCutoffID ***REMOVED***
			continue
		***REMOVED***

		for tags, httpTrails := range subBuckets ***REMOVED***
			trailCount := int64(len(httpTrails))
			if trailCount < c.config.AggregationMinSamples.Int64 ***REMOVED***
				for _, trail := range httpTrails ***REMOVED***
					newSamples = append(newSamples, NewSampleFromTrail(trail))
				***REMOVED***
				continue
			***REMOVED***

			aggrData := &SampleDataAggregatedHTTPReqs***REMOVED***
				Time: Timestamp(time.Unix(0, bucketID*aggrPeriod+aggrPeriod/2)),
				Type: "aggregated_trend",
				Tags: tags,
			***REMOVED***

			if c.config.AggregationSkipOutlierDetection.Bool ***REMOVED***
				// Simply add up all HTTP trails, no outlier detection
				for _, trail := range httpTrails ***REMOVED***
					aggrData.Add(trail)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				connDurations := make(durations, trailCount)
				reqDurations := make(durations, trailCount)
				for i, trail := range httpTrails ***REMOVED***
					connDurations[i] = trail.ConnDuration
					reqDurations[i] = trail.Duration
				***REMOVED***

				var minConnDur, maxConnDur, minReqDur, maxReqDur time.Duration
				if trailCount < c.config.AggregationOutlierAlgoThreshold.Int64 ***REMOVED***
					// Since there are fewer samples, we'll use the interpolation-enabled and
					// more precise sorting-based algorithm
					minConnDur, maxConnDur = connDurations.SortGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef, true)
					minReqDur, maxReqDur = reqDurations.SortGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef, true)
				***REMOVED*** else ***REMOVED***
					minConnDur, maxConnDur = connDurations.SelectGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef)
					minReqDur, maxReqDur = reqDurations.SelectGetNormalBounds(iqrRadius, iqrLowerCoef, iqrUpperCoef)
				***REMOVED***

				for _, trail := range httpTrails ***REMOVED***
					if trail.ConnDuration < minConnDur ||
						trail.ConnDuration > maxConnDur ||
						trail.Duration < minReqDur ||
						trail.Duration > maxReqDur ***REMOVED***
						// Seems like an outlier, add it as a standalone metric
						newSamples = append(newSamples, NewSampleFromTrail(trail))
					***REMOVED*** else ***REMOVED***
						// Aggregate the trail
						aggrData.Add(trail)
					***REMOVED***
				***REMOVED***
			***REMOVED***

			aggrData.CalcAverages()

			if aggrData.Count > 0 ***REMOVED***
				log.WithFields(log.Fields***REMOVED***
					"http_samples": aggrData.Count,
				***REMOVED***).Debug("Aggregated HTTP metrics")
				newSamples = append(newSamples, &Sample***REMOVED***
					Type:   DataTypeAggregatedHTTPReqs,
					Metric: "http_req_li_all",
					Data:   aggrData,
				***REMOVED***)
			***REMOVED***
		***REMOVED***
		delete(c.aggrBuckets, bucketID)
	***REMOVED***

	if len(newSamples) > 0 ***REMOVED***
		c.bufferMutex.Lock()
		c.bufferSamples = append(c.bufferSamples, newSamples...)
		c.bufferMutex.Unlock()
	***REMOVED***
***REMOVED***

func (c *Collector) flushHTTPTrails() ***REMOVED***
	c.bufferMutex.Lock()
	defer c.bufferMutex.Unlock()

	newSamples := []*Sample***REMOVED******REMOVED***
	for _, trail := range c.bufferHTTPTrails ***REMOVED***
		newSamples = append(newSamples, NewSampleFromTrail(trail))
	***REMOVED***
	for _, bucket := range c.aggrBuckets ***REMOVED***
		for _, trails := range bucket ***REMOVED***
			for _, trail := range trails ***REMOVED***
				newSamples = append(newSamples, NewSampleFromTrail(trail))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	c.bufferHTTPTrails = nil
	c.aggrBuckets = map[int64]aggregationBucket***REMOVED******REMOVED***
	c.bufferSamples = append(c.bufferSamples, newSamples...)
***REMOVED***
func (c *Collector) pushMetrics() ***REMOVED***
	c.bufferMutex.Lock()
	if len(c.bufferSamples) == 0 ***REMOVED***
		c.bufferMutex.Unlock()
		return
	***REMOVED***
	buffer := c.bufferSamples
	c.bufferSamples = nil
	c.bufferMutex.Unlock()

	log.WithFields(log.Fields***REMOVED***
		"samples": len(buffer),
	***REMOVED***).Debug("Pushing metrics to cloud")

	for len(buffer) > 0 ***REMOVED***
		var size = len(buffer)
		if size > int(c.config.MaxMetricSamplesPerPackage.Int64) ***REMOVED***
			size = int(c.config.MaxMetricSamplesPerPackage.Int64)
		***REMOVED***
		err := c.client.PushMetric(c.referenceID, c.config.NoCompress.Bool, buffer[:size])
		if err != nil ***REMOVED***
			log.WithFields(log.Fields***REMOVED***
				"error": err,
			***REMOVED***).Warn("Failed to send metrics to cloud")
		***REMOVED***
		buffer = buffer[size:]
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
			thresholdResults[name][t.Source] = t.LastFailed
			if t.LastFailed ***REMOVED***
				testTainted = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"ref":     c.referenceID,
		"tainted": testTainted,
	***REMOVED***).Debug("Sending test finished")

	runStatus := lib.RunStatusFinished
	if c.runStatus != 0 ***REMOVED***
		runStatus = c.runStatus
	***REMOVED***

	err := c.client.TestFinished(c.referenceID, thresholdResults, testTainted, runStatus)
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

// SetRunStatus Set run status
func (c *Collector) SetRunStatus(status lib.RunStatus) ***REMOVED***
	c.runStatus = status
***REMOVED***
