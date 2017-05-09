package cloud

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/mitchellh/mapstructure"
)

type loadimpactConfig struct ***REMOVED***
	ProjectId int    `mapstructure:"project_id"`
	Name      string `mapstructure:"name"`
***REMOVED***

// Collector sends result data to the Load Impact cloud service.
type Collector struct ***REMOVED***
	referenceID string
	initErr     error // Possible error from init call to cloud API

	name       string
	project_id int

	duration   int64
	thresholds map[string][]*stats.Threshold
	client     *Client
***REMOVED***

// New creates a new cloud collector
func New(fname string, src *lib.SourceData, opts lib.Options, version string) (*Collector, error) ***REMOVED***
	token := os.Getenv("K6CLOUD_TOKEN")

	var extConfig loadimpactConfig
	if val, ok := opts.External["loadimpact"]; ok ***REMOVED***
		err := mapstructure.Decode(val, &extConfig)
		if err != nil ***REMOVED***
			log.Warn("Malformed loadimpact settings in script options")
		***REMOVED***
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
		client:     NewClient(token, "", version),
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

	response, err := c.client.CreateTestRun(testRun)

	if err != nil ***REMOVED***
		c.initErr = err
		log.WithFields(log.Fields***REMOVED***
			"error": err,
		***REMOVED***).Error("Cloud collector failed to init")
		return
	***REMOVED***
	c.referenceID = response.ReferenceID

	log.WithFields(log.Fields***REMOVED***
		"name":        c.name,
		"projectId":   c.project_id,
		"duration":    c.duration,
		"referenceId": c.referenceID,
	***REMOVED***).Debug("Cloud collector init successful")
***REMOVED***

func (c *Collector) String() string ***REMOVED***
	if c.initErr == nil ***REMOVED***
		return fmt.Sprintf("Load Impact (https://app.loadimpact.com/k6/runs/%s)", c.referenceID)
	***REMOVED***

	switch c.initErr ***REMOVED***
	case ErrNotAuthorized:
	case ErrNotAuthorized:
		return c.initErr.Error()
	***REMOVED***
	return fmt.Sprintf("Failed to create test in Load Impact cloud")
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

	if c.referenceID != "" ***REMOVED***
		err := c.client.TestFinished(c.referenceID, thresholdResults, testTainted)
		if err != nil ***REMOVED***
			log.WithFields(log.Fields***REMOVED***
				"error": err,
			***REMOVED***).Warn("Failed to send test finished to cloud")
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	if c.referenceID == "" ***REMOVED***
		return
	***REMOVED***

	var cloudSamples []*sample
	for _, samp := range samples ***REMOVED***
		sampleJSON := &sample***REMOVED***
			Type:   "Point",
			Metric: samp.Metric.Name,
			Data: sampleData***REMOVED***
				Type:  samp.Metric.Type,
				Time:  samp.Time,
				Value: samp.Value,
				Tags:  samp.Tags,
			***REMOVED***,
		***REMOVED***
		cloudSamples = append(cloudSamples, sampleJSON)
	***REMOVED***

	if len(cloudSamples) > 0 ***REMOVED***
		err := c.client.PushMetric(c.referenceID, cloudSamples)
		if err != nil ***REMOVED***
			log.WithFields(log.Fields***REMOVED***
				"error":   err,
				"samples": cloudSamples,
			***REMOVED***).Warn("Failed to send metrics to cloud")
		***REMOVED***
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
