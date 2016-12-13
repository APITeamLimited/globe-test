package json

import (
	"github.com/loadimpact/k6/stats"
	"time"
)

type Envelope struct ***REMOVED***
	Type   string      `json:"type"`
	Data   interface***REMOVED******REMOVED*** `json:"data"`
	Metric string      `json:"metric,omitempty"`
***REMOVED***

type JSONSample struct ***REMOVED***
	Time  time.Time         `json:"time"`
	Value float64           `json:"value"`
	Tags  map[string]string `json:"tags"`
***REMOVED***

func NewJSONSample(sample *stats.Sample) *JSONSample ***REMOVED***
	return &JSONSample***REMOVED***
		Time:  sample.Time,
		Value: sample.Value,
		Tags:  sample.Tags,
	***REMOVED***
***REMOVED***

func Wrap(t interface***REMOVED******REMOVED***) *Envelope ***REMOVED***
	switch data := t.(type) ***REMOVED***
	case stats.Sample:
		return &Envelope***REMOVED***
			Type:   "Point",
			Metric: data.Metric.Name,
			Data:   NewJSONSample(&data),
		***REMOVED***
	case *stats.Metric:
		return &Envelope***REMOVED***
			Type:   "Metric",
			Metric: data.Name,
			Data:   data,
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
