package cloud

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"

	"github.com/mitchellh/mapstructure"
)

type loadimpactConfig struct ***REMOVED***
	ProjectId int    `mapstructure:"project_id"`
	Name      string `mapstructure:"name"`
***REMOVED***

// Collector sends results data to the Load Impact cloud service.
type Collector struct ***REMOVED***
	referenceID string

	name       string
	project_id int

	duration   int64
	thresholds map[string][]*stats.Threshold
	client     *Client
***REMOVED***

// New creates a new cloud collector
func New(fname string, src *lib.SourceData, opts lib.Options) (*Collector, error) ***REMOVED***
	token := os.Getenv("K6CLOUD_TOKEN")

	var extConfig loadimpactConfig
	if val, ok := opts.External["loadimpact"]; ok ***REMOVED***
		err := mapstructure.Decode(val, &extConfig)
		if err != nil ***REMOVED***
			// For now we ignore if loadimpact section is malformed
		***REMOVED***
	***REMOVED***

	thresholds := make(map[string][]*stats.Threshold)

	for name, t := range opts.Thresholds ***REMOVED***
		for _, threshold := range t.Thresholds ***REMOVED***
			thresholds[name] = append(thresholds[name], threshold)
		***REMOVED***
	***REMOVED***

	// Sum test duration from options. -1 for unknown duration.
	var duration int64 = -1
	if len(opts.Stages) > 0 ***REMOVED***
		duration = sumStages(opts.Stages)
	***REMOVED*** else if opts.Duration.Valid ***REMOVED***
		// Parse duration if no stages found
		dur, err := time.ParseDuration(opts.Duration.String)
		// ignore error and keep default -1 value
		if err == nil ***REMOVED***
			duration = int64(dur.Seconds())
		***REMOVED***
	***REMOVED***

	return &Collector***REMOVED***
		name:       getName(src, extConfig),
		project_id: getProjectId(extConfig),
		thresholds: thresholds,
		client:     NewClient(token),
		duration:   duration,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) Init() ***REMOVED***

	thresholds := make(map[string][]string)

	for name, t := range c.thresholds ***REMOVED***
		for _, threshold := range t ***REMOVED***
			thresholds[name] = append(thresholds[name], threshold.Source)
		***REMOVED***
	***REMOVED***

	testRun := &TestRun***REMOVED***
		Name:       c.name,
		Thresholds: thresholds,
		Duration:   c.duration,
		ProjectID:  c.project_id,
	***REMOVED***

	// TODO fix this and add proper error handling
	response := c.client.CreateTestRun(testRun)
	if response != nil ***REMOVED***
		c.referenceID = response.ReferenceID
	***REMOVED*** else ***REMOVED***
		log.Warn("Failed to create test in Load Impact cloud")
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"name":        c.name,
		"projectId":   c.project_id,
		"duration":    c.duration,
		"referenceId": c.referenceID,
	***REMOVED***).Debug("Cloud collector init")
***REMOVED***

func (c *Collector) String() string ***REMOVED***
	return fmt.Sprintf("Load Impact (https://app.staging.loadimpact.com/k6/runs/%s)", c.referenceID)
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	<-ctx.Done()

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

	if c.referenceID == "" ***REMOVED***
		c.client.TestFinished(c.referenceID, thresholdResults, testTainted)
	***REMOVED***
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	if c.referenceID == "" ***REMOVED***
		return
	***REMOVED***

	var cloudSamples []*Sample
	for _, sample := range samples ***REMOVED***
		sampleJSON := &Sample***REMOVED***
			Type:   "Point",
			Metric: sample.Metric.Name,
			Data: SampleData***REMOVED***
				Type:  sample.Metric.Type,
				Time:  sample.Time,
				Value: sample.Value,
				Tags:  sample.Tags,
			***REMOVED***,
		***REMOVED***
		cloudSamples = append(cloudSamples, sampleJSON)
	***REMOVED***

	if len(cloudSamples) > 0 ***REMOVED***
		c.client.PushMetric(c.referenceID, cloudSamples)
	***REMOVED***
***REMOVED***

func sumStages(stages []lib.Stage) int64 ***REMOVED***
	var total time.Duration
	for _, stage := range stages ***REMOVED***
		total += stage.Duration
	***REMOVED***

	return int64(total.Seconds())
***REMOVED***

func getProjectId(extConfig loadimpactConfig) int ***REMOVED***
	env := os.Getenv("K6CLOUD_PROJECTID")
	if env != "" ***REMOVED***
		id, err := strconv.Atoi(env)
		if err == nil && id > 0 ***REMOVED***
			return id
		***REMOVED***
	***REMOVED***

	if extConfig.ProjectId > 0 ***REMOVED***
		return extConfig.ProjectId
	***REMOVED***

	return 0
***REMOVED***

func getName(src *lib.SourceData, extConfig loadimpactConfig) string ***REMOVED***
	envName := os.Getenv("K6CLOUD_NAME")
	if envName != "" ***REMOVED***
		return envName
	***REMOVED***

	if extConfig.Name != "" ***REMOVED***
		return extConfig.Name
	***REMOVED***

	if src.Filename != "" ***REMOVED***
		name := filepath.Base(src.Filename)
		if name != "" || name != "." ***REMOVED***
			return name
		***REMOVED***
	***REMOVED***

	return "k6 test"
***REMOVED***
