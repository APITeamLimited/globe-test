package cloud

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

// Collector sends results data to the Load Impact cloud service.
type Collector struct ***REMOVED***
	referenceID string

	duration   int64
	thresholds map[string][]string
	client     *Client
***REMOVED***

func New(fname string, opts lib.Options) (*Collector, error) ***REMOVED***
	referenceID := os.Getenv("K6CLOUD_REFERENCEID")
	token := os.Getenv("K6CLOUD_TOKEN")

	thresholds := make(map[string][]string)

	for name, t := range opts.Thresholds ***REMOVED***
		for _, threshold := range t.Thresholds ***REMOVED***
			thresholds[name] = append(thresholds[name], threshold.Source)
		***REMOVED***
	***REMOVED***

	// Sum test duration from options. -1 for unknown duration.
	var duration int64 = -1
	if len(opts.Stages) > 0 ***REMOVED***
		duration = sumStages(opts.Stages)
	***REMOVED***

	return &Collector***REMOVED***
		referenceID: referenceID,
		thresholds:  thresholds,
		client:      NewClient(token),
		duration:    duration,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) Init() ***REMOVED***
	name := os.Getenv("K6CLOUD_NAME")
	if name == "" ***REMOVED***
		name = "k6 test"
	***REMOVED***

	// TODO fix this and add proper error handling
	if c.referenceID == "" ***REMOVED***
		response := c.client.CreateTestRun(name, c.thresholds, c.duration)
		if response != nil ***REMOVED***
			c.referenceID = response.ReferenceID
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) String() string ***REMOVED***
	return fmt.Sprintf("Load Impact (https://app.staging.loadimpact.com/k6/runs/%s)", c.referenceID)
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	t := time.Now()
	<-ctx.Done()
	s := time.Now()

	c.client.TestFinished(c.referenceID)

	log.Debug(fmt.Sprintf("http://localhost:5000/v1/metrics/%s/%d000/%d000\n", c.referenceID, t.Unix(), s.Unix()))
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***

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

	if len(cloudSamples) > 0 && c.referenceID != "" ***REMOVED***
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
