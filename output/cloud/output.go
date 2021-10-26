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
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/cloudapi"
	"go.k6.io/k6/output"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/netext"
	"go.k6.io/k6/lib/netext/httpext"
	"go.k6.io/k6/stats"
)

// TestName is the default Load Impact Cloud test name
const TestName = "k6 test"

// Output sends result data to the Load Impact cloud service.
type Output struct ***REMOVED***
	config      cloudapi.Config
	referenceID string

	executionPlan []lib.ExecutionStep
	duration      int64 // in seconds
	thresholds    map[string][]*stats.Threshold
	client        *MetricsClient

	runStatus lib.RunStatus

	bufferMutex      sync.Mutex
	bufferHTTPTrails []*httpext.Trail
	bufferSamples    []*Sample

	logger logrus.FieldLogger
	opts   lib.Options

	// TODO: optimize this
	//
	// Since the real-time metrics refactoring (https://github.com/k6io/k6/pull/678),
	// we should no longer have to handle metrics that have times long in the past. So instead of a
	// map, we can probably use a simple slice (or even an array!) as a ring buffer to store the
	// aggregation buckets. This should save us a some time, since it would make the lookups and WaitPeriod
	// checks basically O(1). And even if for some reason there are occasional metrics with past times that
	// don't fit in the chosen ring buffer size, we could just send them along to the buffer unaggregated
	aggrBuckets map[int64]map[[3]string]aggregationBucket

	stopSendingMetrics chan struct***REMOVED******REMOVED***
	stopAggregation    chan struct***REMOVED******REMOVED***
	aggregationDone    *sync.WaitGroup
	stopOutput         chan struct***REMOVED******REMOVED***
	outputDone         *sync.WaitGroup
	engineStopFunc     func(error)
***REMOVED***

// Verify that Output implements the wanted interfaces
var _ interface ***REMOVED***
	output.WithRunStatusUpdates
	output.WithThresholds
	output.WithTestRunStop
***REMOVED*** = &Output***REMOVED******REMOVED***

// New creates a new cloud output.
func New(params output.Params) (output.Output, error) ***REMOVED***
	return newOutput(params)
***REMOVED***

// New creates a new cloud output.
func newOutput(params output.Params) (*Output, error) ***REMOVED***
	conf, err := cloudapi.GetConsolidatedConfig(
		params.JSONConfig, params.Environment, params.ConfigArgument, params.ScriptOptions.External)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := validateRequiredSystemTags(params.ScriptOptions.SystemTags); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	logger := params.Logger.WithFields(logrus.Fields***REMOVED***"output": "cloud"***REMOVED***)

	if conf.AggregationPeriod.Duration > 0 &&
		(params.ScriptOptions.SystemTags.Has(stats.TagVU) || params.ScriptOptions.SystemTags.Has(stats.TagIter)) ***REMOVED***
		return nil, errors.New("aggregation cannot be enabled if the 'vu' or 'iter' system tag is also enabled")
	***REMOVED***

	if !conf.Name.Valid || conf.Name.String == "" ***REMOVED***
		scriptPath := params.ScriptPath.String()
		if scriptPath == "" ***REMOVED***
			// Script from stdin without a name, likely from stdin
			return nil, errors.New("script name not set, please specify K6_CLOUD_NAME or options.ext.loadimpact.name")
		***REMOVED***

		conf.Name = null.StringFrom(filepath.Base(scriptPath))
	***REMOVED***
	if conf.Name.String == "-" ***REMOVED***
		conf.Name = null.StringFrom(TestName)
	***REMOVED***

	duration, testEnds := lib.GetEndOffset(params.ExecutionPlan)
	if !testEnds ***REMOVED***
		return nil, errors.New("tests with unspecified duration are not allowed when outputting data to k6 cloud")
	***REMOVED***

	if !(conf.MetricPushConcurrency.Int64 > 0) ***REMOVED***
		return nil, fmt.Errorf("metrics push concurrency must be a positive number but is %d",
			conf.MetricPushConcurrency.Int64)
	***REMOVED***

	if !(conf.MaxMetricSamplesPerPackage.Int64 > 0) ***REMOVED***
		return nil, fmt.Errorf("metric samples per package must be a positive number but is %d",
			conf.MaxMetricSamplesPerPackage.Int64)
	***REMOVED***

	apiClient := cloudapi.NewClient(
		logger, conf.Token.String, conf.Host.String, consts.Version, conf.Timeout.TimeDuration())

	return &Output***REMOVED***
		config:        conf,
		client:        NewMetricsClient(apiClient, logger, conf.Host.String, conf.NoCompress.Bool),
		executionPlan: params.ExecutionPlan,
		duration:      int64(duration / time.Second),
		opts:          params.ScriptOptions,
		aggrBuckets:   map[int64]map[[3]string]aggregationBucket***REMOVED******REMOVED***,
		logger:        logger,

		stopSendingMetrics: make(chan struct***REMOVED******REMOVED***),
		stopAggregation:    make(chan struct***REMOVED******REMOVED***),
		aggregationDone:    &sync.WaitGroup***REMOVED******REMOVED***,
		stopOutput:         make(chan struct***REMOVED******REMOVED***),
		outputDone:         &sync.WaitGroup***REMOVED******REMOVED***,
	***REMOVED***, nil
***REMOVED***

// validateRequiredSystemTags checks if all required tags are present.
func validateRequiredSystemTags(scriptTags *stats.SystemTagSet) error ***REMOVED***
	missingRequiredTags := []string***REMOVED******REMOVED***
	requiredTags := stats.TagName | stats.TagMethod | stats.TagStatus | stats.TagError | stats.TagCheck | stats.TagGroup
	for _, tag := range stats.SystemTagSetValues() ***REMOVED***
		if requiredTags.Has(tag) && !scriptTags.Has(tag) ***REMOVED***
			missingRequiredTags = append(missingRequiredTags, tag.String())
		***REMOVED***
	***REMOVED***
	if len(missingRequiredTags) > 0 ***REMOVED***
		return fmt.Errorf(
			"the cloud output needs the following system tags enabled: %s",
			strings.Join(missingRequiredTags, ", "),
		)
	***REMOVED***
	return nil
***REMOVED***

// Start calls the k6 Cloud API to initialize the test run, and then starts the
// goroutine that would listen for metric samples and send them to the cloud.
func (out *Output) Start() error ***REMOVED***
	if out.config.PushRefID.Valid ***REMOVED***
		out.referenceID = out.config.PushRefID.String
		out.logger.WithField("referenceId", out.referenceID).Debug("directly pushing metrics without init")
		out.startBackgroundProcesses()
		return nil
	***REMOVED***

	thresholds := make(map[string][]string)

	for name, t := range out.thresholds ***REMOVED***
		for _, threshold := range t ***REMOVED***
			thresholds[name] = append(thresholds[name], threshold.Source)
		***REMOVED***
	***REMOVED***
	maxVUs := lib.GetMaxPossibleVUs(out.executionPlan)

	testRun := &cloudapi.TestRun***REMOVED***
		Name:       out.config.Name.String,
		ProjectID:  out.config.ProjectID.Int64,
		VUsMax:     int64(maxVUs),
		Thresholds: thresholds,
		Duration:   out.duration,
	***REMOVED***

	response, err := out.client.CreateTestRun(testRun)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	out.referenceID = response.ReferenceID

	if response.ConfigOverride != nil ***REMOVED***
		out.logger.WithFields(logrus.Fields***REMOVED***
			"override": response.ConfigOverride,
		***REMOVED***).Debug("overriding config options")
		out.config = out.config.Apply(*response.ConfigOverride)
	***REMOVED***

	out.startBackgroundProcesses()

	out.logger.WithFields(logrus.Fields***REMOVED***
		"name":        out.config.Name,
		"projectId":   out.config.ProjectID,
		"duration":    out.duration,
		"referenceId": out.referenceID,
	***REMOVED***).Debug("Started!")
	return nil
***REMOVED***

func (out *Output) startBackgroundProcesses() ***REMOVED***
	aggregationPeriod := time.Duration(out.config.AggregationPeriod.Duration)
	// If enabled, start periodically aggregating the collected HTTP trails
	if aggregationPeriod > 0 ***REMOVED***
		out.aggregationDone.Add(1)
		go func() ***REMOVED***
			defer out.aggregationDone.Done()
			aggregationWaitPeriod := time.Duration(out.config.AggregationWaitPeriod.Duration)
			aggregationTicker := time.NewTicker(aggregationPeriod)
			defer aggregationTicker.Stop()

			for ***REMOVED***
				select ***REMOVED***
				case <-out.stopSendingMetrics:
					return
				case <-aggregationTicker.C:
					out.aggregateHTTPTrails(aggregationWaitPeriod)
				case <-out.stopAggregation:
					out.aggregateHTTPTrails(0)
					out.flushHTTPTrails()
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	out.outputDone.Add(1)
	go func() ***REMOVED***
		defer out.outputDone.Done()
		pushTicker := time.NewTicker(time.Duration(out.config.MetricPushInterval.Duration))
		defer pushTicker.Stop()
		for ***REMOVED***
			select ***REMOVED***
			case <-out.stopSendingMetrics:
				return
			default:
			***REMOVED***
			select ***REMOVED***
			case <-out.stopOutput:
				out.pushMetrics()
				return
			case <-pushTicker.C:
				out.pushMetrics()
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

// Stop gracefully stops all metric emission from the output and when all metric
// samples are emitted, it sends an API to the cloud to finish the test run.
func (out *Output) Stop() error ***REMOVED***
	out.logger.Debug("Stopping the cloud output...")
	close(out.stopAggregation)
	out.aggregationDone.Wait() // could be a no-op, if we have never started the aggregation
	out.logger.Debug("Aggregation stopped, stopping metric emission...")
	close(out.stopOutput)
	out.outputDone.Wait()
	out.logger.Debug("Metric emission stopped, calling cloud API...")
	err := out.testFinished()
	if err != nil ***REMOVED***
		out.logger.WithFields(logrus.Fields***REMOVED***"error": err***REMOVED***).Warn("Failed to send test finished to the cloud")
	***REMOVED*** else ***REMOVED***
		out.logger.Debug("Cloud output successfully stopped!")
	***REMOVED***
	return err
***REMOVED***

// Description returns the URL with the test run results.
func (out *Output) Description() string ***REMOVED***
	return fmt.Sprintf("cloud (%s)", cloudapi.URLForResults(out.referenceID, out.config))
***REMOVED***

// SetRunStatus receives the latest run status from the Engine.
func (out *Output) SetRunStatus(status lib.RunStatus) ***REMOVED***
	out.runStatus = status
***REMOVED***

// SetThresholds receives the thresholds before the output is Start()-ed.
func (out *Output) SetThresholds(scriptThresholds map[string]stats.Thresholds) ***REMOVED***
	thresholds := make(map[string][]*stats.Threshold)
	for name, t := range scriptThresholds ***REMOVED***
		thresholds[name] = append(thresholds[name], t.Thresholds...)
	***REMOVED***
	out.thresholds = thresholds
***REMOVED***

// SetTestRunStopCallback receives the function that stops the engine on error
func (out *Output) SetTestRunStopCallback(stopFunc func(error)) ***REMOVED***
	out.engineStopFunc = stopFunc
***REMOVED***

func useCloudTags(source *httpext.Trail) *httpext.Trail ***REMOVED***
	name, nameExist := source.Tags.Get("name")
	url, urlExist := source.Tags.Get("url")
	if !nameExist || !urlExist || name == url ***REMOVED***
		return source
	***REMOVED***

	newTags := source.Tags.CloneTags()
	newTags["url"] = name

	dest := new(httpext.Trail)
	*dest = *source
	dest.Tags = stats.IntoSampleTags(&newTags)
	dest.Samples = nil

	return dest
***REMOVED***

// AddMetricSamples receives a set of metric samples. This method is never
// called concurrently, so it defers as much of the work as possible to the
// asynchronous goroutines initialized in Start().
func (out *Output) AddMetricSamples(sampleContainers []stats.SampleContainer) ***REMOVED***
	select ***REMOVED***
	case <-out.stopSendingMetrics:
		return
	default:
	***REMOVED***

	if out.referenceID == "" ***REMOVED***
		return
	***REMOVED***

	newSamples := []*Sample***REMOVED******REMOVED***
	newHTTPTrails := []*httpext.Trail***REMOVED******REMOVED***

	for _, sampleContainer := range sampleContainers ***REMOVED***
		switch sc := sampleContainer.(type) ***REMOVED***
		case *httpext.Trail:
			sc = useCloudTags(sc)
			// Check if aggregation is enabled,
			if out.config.AggregationPeriod.Duration > 0 ***REMOVED***
				newHTTPTrails = append(newHTTPTrails, sc)
			***REMOVED*** else ***REMOVED***
				newSamples = append(newSamples, NewSampleFromTrail(sc))
			***REMOVED***
		case *netext.NetTrail:
			// TODO: aggregate?
			values := map[string]float64***REMOVED***
				metrics.DataSentName:     float64(sc.BytesWritten),
				metrics.DataReceivedName: float64(sc.BytesRead),
			***REMOVED***

			if sc.FullIteration ***REMOVED***
				values[metrics.IterationDurationName] = stats.D(sc.EndTime.Sub(sc.StartTime))
				values[metrics.IterationsName] = 1
			***REMOVED***

			newSamples = append(newSamples, &Sample***REMOVED***
				Type:   DataTypeMap,
				Metric: "iter_li_all",
				Data: &SampleDataMap***REMOVED***
					Time:   toMicroSecond(sc.GetTime()),
					Tags:   sc.GetTags(),
					Values: values,
				***REMOVED***,
			***REMOVED***)
		default:
			for _, sample := range sampleContainer.GetSamples() ***REMOVED***
				newSamples = append(newSamples, &Sample***REMOVED***
					Type:   DataTypeSingle,
					Metric: sample.Metric.Name,
					Data: &SampleDataSingle***REMOVED***
						Type:  sample.Metric.Type,
						Time:  toMicroSecond(sample.Time),
						Tags:  sample.Tags,
						Value: sample.Value,
					***REMOVED***,
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(newSamples) > 0 || len(newHTTPTrails) > 0 ***REMOVED***
		out.bufferMutex.Lock()
		out.bufferSamples = append(out.bufferSamples, newSamples...)
		out.bufferHTTPTrails = append(out.bufferHTTPTrails, newHTTPTrails...)
		out.bufferMutex.Unlock()
	***REMOVED***
***REMOVED***

//nolint:funlen,nestif,gocognit
func (out *Output) aggregateHTTPTrails(waitPeriod time.Duration) ***REMOVED***
	out.bufferMutex.Lock()
	newHTTPTrails := out.bufferHTTPTrails
	out.bufferHTTPTrails = nil
	out.bufferMutex.Unlock()

	aggrPeriod := int64(out.config.AggregationPeriod.Duration)

	// Distribute all newly buffered HTTP trails into buckets and sub-buckets

	// this key is here specifically to not incur more allocations then necessary
	// if you change this code please run the benchmarks and add the results to the commit message
	var subBucketKey [3]string
	for _, trail := range newHTTPTrails ***REMOVED***
		trailTags := trail.GetTags()
		bucketID := trail.GetTime().UnixNano() / aggrPeriod

		// Get or create a time bucket for that trail period
		bucket, ok := out.aggrBuckets[bucketID]
		if !ok ***REMOVED***
			bucket = make(map[[3]string]aggregationBucket)
			out.aggrBuckets[bucketID] = bucket
		***REMOVED***
		subBucketKey[0], _ = trailTags.Get("name")
		subBucketKey[1], _ = trailTags.Get("group")
		subBucketKey[2], _ = trailTags.Get("status")

		subBucket, ok := bucket[subBucketKey]
		if !ok ***REMOVED***
			subBucket = aggregationBucket***REMOVED******REMOVED***
			bucket[subBucketKey] = subBucket
		***REMOVED***
		// Either use an existing subbucket key or use the trail tags as a new one
		subSubBucketKey := trailTags
		subSubBucket, ok := subBucket[subSubBucketKey]
		if !ok ***REMOVED***
			for sbTags, sb := range subBucket ***REMOVED***
				if trailTags.IsEqual(sbTags) ***REMOVED***
					subSubBucketKey = sbTags
					subSubBucket = sb
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		subBucket[subSubBucketKey] = append(subSubBucket, trail)
	***REMOVED***

	// Which buckets are still new and we'll wait for trails to accumulate before aggregating
	bucketCutoffID := time.Now().Add(-waitPeriod).UnixNano() / aggrPeriod
	iqrRadius := out.config.AggregationOutlierIqrRadius.Float64
	iqrLowerCoef := out.config.AggregationOutlierIqrCoefLower.Float64
	iqrUpperCoef := out.config.AggregationOutlierIqrCoefUpper.Float64
	newSamples := []*Sample***REMOVED******REMOVED***

	// Handle all aggregation buckets older than bucketCutoffID
	for bucketID, subBuckets := range out.aggrBuckets ***REMOVED***
		if bucketID > bucketCutoffID ***REMOVED***
			continue
		***REMOVED***

		for _, subBucket := range subBuckets ***REMOVED***
			for tags, httpTrails := range subBucket ***REMOVED***
				// start := time.Now() // this is in a combination with the log at the end
				trailCount := int64(len(httpTrails))
				if trailCount < out.config.AggregationMinSamples.Int64 ***REMOVED***
					for _, trail := range httpTrails ***REMOVED***
						newSamples = append(newSamples, NewSampleFromTrail(trail))
					***REMOVED***
					continue
				***REMOVED***

				aggrData := &SampleDataAggregatedHTTPReqs***REMOVED***
					Time: toMicroSecond(time.Unix(0, bucketID*aggrPeriod+aggrPeriod/2)),
					Type: "aggregated_trend",
					Tags: tags,
				***REMOVED***

				if out.config.AggregationSkipOutlierDetection.Bool ***REMOVED***
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
					if trailCount < out.config.AggregationOutlierAlgoThreshold.Int64 ***REMOVED***
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
					/*
						out.logger.WithFields(logrus.Fields***REMOVED***
							"http_samples": aggrData.Count,
							"ratio":        fmt.Sprintf("%.2f", float64(aggrData.Count)/float64(trailCount)),
							"t":            time.Since(start),
						***REMOVED***).Debug("Aggregated HTTP metrics")
					//*/
					newSamples = append(newSamples, &Sample***REMOVED***
						Type:   DataTypeAggregatedHTTPReqs,
						Metric: "http_req_li_all",
						Data:   aggrData,
					***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		delete(out.aggrBuckets, bucketID)
	***REMOVED***

	if len(newSamples) > 0 ***REMOVED***
		out.bufferMutex.Lock()
		out.bufferSamples = append(out.bufferSamples, newSamples...)
		out.bufferMutex.Unlock()
	***REMOVED***
***REMOVED***

func (out *Output) flushHTTPTrails() ***REMOVED***
	out.bufferMutex.Lock()
	defer out.bufferMutex.Unlock()

	newSamples := []*Sample***REMOVED******REMOVED***
	for _, trail := range out.bufferHTTPTrails ***REMOVED***
		newSamples = append(newSamples, NewSampleFromTrail(trail))
	***REMOVED***
	for _, bucket := range out.aggrBuckets ***REMOVED***
		for _, subBucket := range bucket ***REMOVED***
			for _, trails := range subBucket ***REMOVED***
				for _, trail := range trails ***REMOVED***
					newSamples = append(newSamples, NewSampleFromTrail(trail))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	out.bufferHTTPTrails = nil
	out.aggrBuckets = map[int64]map[[3]string]aggregationBucket***REMOVED******REMOVED***
	out.bufferSamples = append(out.bufferSamples, newSamples...)
***REMOVED***

func (out *Output) shouldStopSendingMetrics(err error) bool ***REMOVED***
	if err == nil ***REMOVED***
		return false
	***REMOVED***

	if errResp, ok := err.(cloudapi.ErrorResponse); ok && errResp.Response != nil ***REMOVED***
		return errResp.Response.StatusCode == http.StatusForbidden && errResp.Code == 4
	***REMOVED***

	return false
***REMOVED***

type pushJob struct ***REMOVED***
	done    chan error
	samples []*Sample
***REMOVED***

// ceil(a/b)
func ceilDiv(a, b int) int ***REMOVED***
	r := a / b
	if a%b != 0 ***REMOVED***
		r++
	***REMOVED***
	return r
***REMOVED***

func (out *Output) pushMetrics() ***REMOVED***
	out.bufferMutex.Lock()
	if len(out.bufferSamples) == 0 ***REMOVED***
		out.bufferMutex.Unlock()
		return
	***REMOVED***
	buffer := out.bufferSamples
	out.bufferSamples = nil
	out.bufferMutex.Unlock()

	count := len(buffer)
	out.logger.WithFields(logrus.Fields***REMOVED***
		"samples": count,
	***REMOVED***).Debug("Pushing metrics to cloud")
	start := time.Now()

	numberOfPackages := ceilDiv(len(buffer), int(out.config.MaxMetricSamplesPerPackage.Int64))
	numberOfWorkers := int(out.config.MetricPushConcurrency.Int64)
	if numberOfWorkers > numberOfPackages ***REMOVED***
		numberOfWorkers = numberOfPackages
	***REMOVED***

	ch := make(chan pushJob, numberOfPackages)
	for i := 0; i < numberOfWorkers; i++ ***REMOVED***
		go func() ***REMOVED***
			for job := range ch ***REMOVED***
				err := out.client.PushMetric(out.referenceID, job.samples)
				job.done <- err
				if out.shouldStopSendingMetrics(err) ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	jobs := make([]pushJob, 0, numberOfPackages)

	for len(buffer) > 0 ***REMOVED***
		size := len(buffer)
		if size > int(out.config.MaxMetricSamplesPerPackage.Int64) ***REMOVED***
			size = int(out.config.MaxMetricSamplesPerPackage.Int64)
		***REMOVED***
		job := pushJob***REMOVED***done: make(chan error, 1), samples: buffer[:size]***REMOVED***
		ch <- job
		jobs = append(jobs, job)
		buffer = buffer[size:]
	***REMOVED***

	close(ch)

	for _, job := range jobs ***REMOVED***
		err := <-job.done
		if err != nil ***REMOVED***
			if out.shouldStopSendingMetrics(err) ***REMOVED***
				out.logger.WithError(err).Warn("Stopped sending metrics to cloud due to an error")
				if out.config.StopOnError.Bool ***REMOVED***
					out.engineStopFunc(err)
				***REMOVED***
				close(out.stopSendingMetrics)
				break
			***REMOVED***
			out.logger.WithError(err).Warn("Failed to send metrics to cloud")
		***REMOVED***
	***REMOVED***
	out.logger.WithFields(logrus.Fields***REMOVED***
		"samples": count,
		"t":       time.Since(start),
	***REMOVED***).Debug("Pushing metrics to cloud finished")
***REMOVED***

func (out *Output) testFinished() error ***REMOVED***
	if out.referenceID == "" || out.config.PushRefID.Valid ***REMOVED***
		return nil
	***REMOVED***

	testTainted := false
	thresholdResults := make(cloudapi.ThresholdResult)
	for name, thresholds := range out.thresholds ***REMOVED***
		thresholdResults[name] = make(map[string]bool)
		for _, t := range thresholds ***REMOVED***
			thresholdResults[name][t.Source] = t.LastFailed
			if t.LastFailed ***REMOVED***
				testTainted = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	out.logger.WithFields(logrus.Fields***REMOVED***
		"ref":     out.referenceID,
		"tainted": testTainted,
	***REMOVED***).Debug("Sending test finished")

	runStatus := lib.RunStatusFinished
	if out.runStatus != lib.RunStatusQueued ***REMOVED***
		runStatus = out.runStatus
	***REMOVED***

	return out.client.TestFinished(out.referenceID, thresholdResults, testTainted, runStatus)
***REMOVED***

const expectedGzipRatio = 6 // based on test it is around 6.8, but we don't need to be that accurate
