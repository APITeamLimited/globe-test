package jsonlog

import (
	"context"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/stats"
	"net/url"
	"os"
	"time"
)

type Collector struct ***REMOVED***
	f     *os.File
	types map[string]stats.Metric
***REMOVED***

func New(u *url.URL) (*Collector, error) ***REMOVED***
	var fname string

	if u.Path == "" ***REMOVED***
		fname = u.String()
	***REMOVED*** else ***REMOVED***
		fname = u.Path
	***REMOVED***

	logfile, err := os.Create(fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Collector***REMOVED***
		f:     logfile,
		types: map[string]stats.Metric***REMOVED******REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) String() string ***REMOVED***
	return "jsonlog"
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	log.Debug("Writing metrics as JSON to ", c.f.Name())
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			c.writeTypes()
			c.f.Close()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) writeTypes() ***REMOVED***
	types, err := json.Marshal(c.types)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	c.f.WriteString(string(types) + "\n")
***REMOVED***

func (c *Collector) Collect(samples []stats.Sample) ***REMOVED***
	for _, sample := range samples ***REMOVED***
		if _, present := c.types[sample.Metric.Name]; !present ***REMOVED***
			c.types[sample.Metric.Name] = *sample.Metric
		***REMOVED***

		row, err := json.Marshal(NewJSONPoint(&sample))
		if err != nil ***REMOVED***
			// Skip metric if it can't be made into JSON.
			continue
		***REMOVED***
		c.f.WriteString(string(row) + "\n")
	***REMOVED***
***REMOVED***

type JSONPoint struct ***REMOVED***
	Type  string            `json:"type"`
	Time  time.Time         `json:"timestamp"`
	Value float64           `json:"value"`
	Tags  map[string]string `json:"tags"`
***REMOVED***

func NewJSONPoint(sample *stats.Sample) *JSONPoint ***REMOVED***
	return &JSONPoint***REMOVED***
		Type:  sample.Metric.Name,
		Time:  sample.Time,
		Value: sample.Value,
		Tags:  sample.Tags,
	***REMOVED***
***REMOVED***
