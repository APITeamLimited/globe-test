package api

import (
	"errors"
	"github.com/loadimpact/speedboat/stats"
)

var (
	ErrInvalidMetricType = errors.New("Invalid metric type")
	ErrInvalidValueType  = errors.New("Invalid value type")
)

const (
	counterString = `"counter"`
	gaugeString   = `"gauge"`
	trendString   = `"trend"`

	defaultString = `"default"`
	timeString    = `"time"`
)

type Error struct ***REMOVED***
	Title string `json:"title"`
***REMOVED***

type ErrorResponse struct ***REMOVED***
	Errors []Error `json:"errors"`
***REMOVED***

type MetricType stats.MetricType

func (t *MetricType) UnmarshalJSON(data []byte) error ***REMOVED***
	switch string(data) ***REMOVED***
	case counterString:
		*t = MetricType(stats.Counter)
	case gaugeString:
		*t = MetricType(stats.Gauge)
	case trendString:
		*t = MetricType(stats.Trend)
	default:
		return ErrInvalidMetricType
	***REMOVED***

	return nil
***REMOVED***

func (t MetricType) MarshalJSON() ([]byte, error) ***REMOVED***
	switch stats.MetricType(t) ***REMOVED***
	case stats.Counter:
		return []byte(counterString), nil
	case stats.Gauge:
		return []byte(gaugeString), nil
	case stats.Trend:
		return []byte(trendString), nil
	default:
		return nil, ErrInvalidMetricType
	***REMOVED***
***REMOVED***

type ValueType stats.ValueType

func (t *ValueType) UnmarshalJSON(data []byte) error ***REMOVED***
	switch string(data) ***REMOVED***
	case defaultString:
		*t = ValueType(stats.Default)
	case timeString:
		*t = ValueType(stats.Time)
	default:
		return ErrInvalidValueType
	***REMOVED***

	return nil
***REMOVED***

func (t ValueType) MarshalJSON() ([]byte, error) ***REMOVED***
	switch stats.ValueType(t) ***REMOVED***
	case stats.Default:
		return []byte(defaultString), nil
	case stats.Time:
		return []byte(timeString), nil
	default:
		return nil, ErrInvalidValueType
	***REMOVED***
***REMOVED***

type Metric struct ***REMOVED***
	Name     string             `json:"name"`
	Type     MetricType         `json:"type"`
	Contains ValueType          `json:"contains"`
	Data     map[string]float64 `json:"data"`
***REMOVED***

type MetricSet map[string]Metric
