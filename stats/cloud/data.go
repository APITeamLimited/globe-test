package cloud

import (
	"time"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
)

// Sample is the generic struct that contains all types of data that we send to the cloud.
type Sample struct ***REMOVED***
	Type   string      `json:"type"`
	Metric string      `json:"metric"`
	Data   interface***REMOVED******REMOVED*** `json:"data"`
***REMOVED***

// SampleDataSingle is used in all simple un-aggregated single-value samples.
type SampleDataSingle struct ***REMOVED***
	Time  time.Time         `json:"time"`
	Type  stats.MetricType  `json:"type"`
	Tags  *stats.SampleTags `json:"tags,omitempty"`
	Value float64           `json:"value"`
***REMOVED***

// SampleDataMap is used by samples that contain multiple values, currently
// that's only iteration metrics (`iter_li_all`) and unaggregated HTTP
// requests (`http_req_li_all`).
type SampleDataMap struct ***REMOVED***
	Time   time.Time          `json:"time"`
	Tags   *stats.SampleTags  `json:"tags,omitempty"`
	Values map[string]float64 `json:"values,omitempty"`
***REMOVED***

// NewSampleFromTrail just creates a ready-to-send Sample instance
// directly from a netext.Trail.
func NewSampleFromTrail(trail *netext.Trail) *Sample ***REMOVED***
	return &Sample***REMOVED***
		Type:   "Points",
		Metric: "http_req_li_all",
		Data: SampleDataMap***REMOVED***
			Time: trail.GetTime(),
			Tags: trail.GetTags(),
			Values: map[string]float64***REMOVED***
				metrics.HTTPReqs.Name:        1,
				metrics.HTTPReqDuration.Name: stats.D(trail.Duration),

				metrics.HTTPReqBlocked.Name:        stats.D(trail.Blocked),
				metrics.HTTPReqConnecting.Name:     stats.D(trail.Connecting),
				metrics.HTTPReqTLSHandshaking.Name: stats.D(trail.TLSHandshaking),
				metrics.HTTPReqSending.Name:        stats.D(trail.Sending),
				metrics.HTTPReqWaiting.Name:        stats.D(trail.Waiting),
				metrics.HTTPReqReceiving.Name:      stats.D(trail.Receiving),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// SampleDataAggregatedHTTPReqs is used in aggregated samples for HTTP requests.
type SampleDataAggregatedHTTPReqs struct ***REMOVED***
	Time   time.Time         `json:"time"`
	Type   string            `json:"type"`
	Count  uint64            `json:"count"`
	Tags   *stats.SampleTags `json:"tags,omitempty"`
	Values struct ***REMOVED***
		Duration       AggregatedMetric `json:"http_req_duration"`
		Blocked        AggregatedMetric `json:"http_req_blocked"`
		Connecting     AggregatedMetric `json:"http_req_connecting"`
		TLSHandshaking AggregatedMetric `json:"http_req_tls_handshaking"`
		Sending        AggregatedMetric `json:"http_req_sending"`
		Waiting        AggregatedMetric `json:"http_req_waiting"`
		Receiving      AggregatedMetric `json:"http_req_receiving"`
	***REMOVED*** `json:"values"`
***REMOVED***

// CalcAverages calculates and sets all `Avg` properties in the `Values` struct
func (sdagg *SampleDataAggregatedHTTPReqs) CalcAverages() ***REMOVED***
	count := float64(sdagg.Count)
	sdagg.Values.Duration.Avg = float64(sdagg.Values.Duration.sum) / count
	sdagg.Values.Blocked.Avg = float64(sdagg.Values.Blocked.sum) / count
	sdagg.Values.Connecting.Avg = float64(sdagg.Values.Connecting.sum) / count
	sdagg.Values.TLSHandshaking.Avg = float64(sdagg.Values.TLSHandshaking.sum) / count
	sdagg.Values.Sending.Avg = float64(sdagg.Values.Sending.sum) / count
	sdagg.Values.Waiting.Avg = float64(sdagg.Values.Waiting.sum) / count
	sdagg.Values.Receiving.Avg = float64(sdagg.Values.Receiving.sum) / count
***REMOVED***

// AggregatedMetric is used to store aggregated information for a
// particular metric in an SampleDataAggregatedMap.
type AggregatedMetric struct ***REMOVED***
	Min time.Duration `json:"min"`
	Max time.Duration `json:"max"`
	sum time.Duration `json:"-"`   // ignored in JSON output because of SampleDataAggregatedHTTPReqs.Count
	Avg float64       `json:"avg"` // not updated automatically, has to be set externally
***REMOVED***

// Add the new duration to the internal sum and update Min and Max if necessary
func (am *AggregatedMetric) Add(t time.Duration) ***REMOVED***
	if am.sum == 0 || am.Min > t ***REMOVED***
		am.Min = t
	***REMOVED***
	if am.Max < t ***REMOVED***
		am.Max = t
	***REMOVED***
	am.sum += t
***REMOVED***

type aggregationBucket map[*stats.SampleTags][]*netext.Trail

type durations []time.Duration

func (d durations) Len() int           ***REMOVED*** return len(d) ***REMOVED***
func (d durations) Swap(i, j int)      ***REMOVED*** d[i], d[j] = d[j], d[i] ***REMOVED***
func (d durations) Less(i, j int) bool ***REMOVED*** return d[i] < d[j] ***REMOVED***
func (d durations) GetNormalBounds(iqrCoef float64) (min, max time.Duration) ***REMOVED***
	l := len(d)
	if l == 0 ***REMOVED***
		return
	***REMOVED***

	var q1, q3 time.Duration
	if l%4 == 0 ***REMOVED***
		q1 = d[l/4]
		q3 = d[(l/4)*3]
	***REMOVED*** else ***REMOVED***
		q1 = (d[l/4] + d[(l/4)+1]) / 2
		q3 = (d[(l/4)*3] + d[(l/4)*3+1]) / 2
	***REMOVED***

	iqr := float64(q3 - q1)
	min = q1 - time.Duration(iqrCoef*iqr)
	max = q3 + time.Duration(iqrCoef*iqr)
	return
***REMOVED***
