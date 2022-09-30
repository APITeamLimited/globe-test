package workerMetrics

import "errors"

// A MetricType specifies the type of a metric.
type MetricType int

// Possible values for MetricType.
const (
	Counter = MetricType(iota) // A counter that sums its data points
	Gauge                      // A gauge that displays the latest value
	Trend                      // A trend, min/max/avg/med are interesting
	Rate                       // A rate, displays % of values that aren't 0
)

// ErrInvalidMetricType indicates the serialized metric type is invalid.
var ErrInvalidMetricType = errors.New("invalid metric type")

const (
	counterString = "counter"
	gaugeString   = "gauge"
	trendString   = "trend"
	rateString    = "rate"

	defaultString = "default"
	timeString    = "time"
	dataString    = "data"
)

// MarshalJSON serializes a MetricType as a human readable string.
func (t MetricType) MarshalJSON() ([]byte, error) ***REMOVED***
	txt, err := t.MarshalText()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return []byte(`"` + string(txt) + `"`), nil
***REMOVED***

// MarshalText serializes a MetricType as a human readable string.
func (t MetricType) MarshalText() ([]byte, error) ***REMOVED***
	switch t ***REMOVED***
	case Counter:
		return []byte(counterString), nil
	case Gauge:
		return []byte(gaugeString), nil
	case Trend:
		return []byte(trendString), nil
	case Rate:
		return []byte(rateString), nil
	default:
		return nil, ErrInvalidMetricType
	***REMOVED***
***REMOVED***

// UnmarshalText deserializes a MetricType from a string representation.
func (t *MetricType) UnmarshalText(data []byte) error ***REMOVED***
	switch string(data) ***REMOVED***
	case counterString:
		*t = Counter
	case gaugeString:
		*t = Gauge
	case trendString:
		*t = Trend
	case rateString:
		*t = Rate
	default:
		return ErrInvalidMetricType
	***REMOVED***

	return nil
***REMOVED***

func (t MetricType) String() string ***REMOVED***
	switch t ***REMOVED***
	case Counter:
		return counterString
	case Gauge:
		return gaugeString
	case Trend:
		return trendString
	case Rate:
		return rateString
	default:
		return "[INVALID]"
	***REMOVED***
***REMOVED***

// supportedAggregationMethods returns the list of threshold aggregation methods
// that can be used against this MetricType.
func (t MetricType) supportedAggregationMethods() []string ***REMOVED***
	switch t ***REMOVED***
	case Counter:
		return []string***REMOVED***tokenCount, tokenRate***REMOVED***
	case Gauge:
		return []string***REMOVED***tokenValue***REMOVED***
	case Rate:
		return []string***REMOVED***tokenRate***REMOVED***
	case Trend:
		return []string***REMOVED***
			tokenAvg,
			tokenMin,
			tokenMax,
			tokenMed,
			tokenPercentile,
		***REMOVED***
	default:
		// unreachable!
		panic("unreachable")
	***REMOVED***
***REMOVED***

// supportsAggregationMethod returns whether the MetricType supports a
// given threshold aggregation method or not.
func (t MetricType) supportsAggregationMethod(aggregationMethod string) bool ***REMOVED***
	for _, m := range t.supportedAggregationMethods() ***REMOVED***
		if aggregationMethod == m ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***
