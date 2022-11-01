package types

import (
	"bytes"
	"encoding/json"
)

func DefaultOutputConfig() NullOutputConfig ***REMOVED***
	return NullOutputConfig***REMOVED***
		OutputConfig***REMOVED***
			Graphs: []MetricGraph***REMOVED***
				***REMOVED***
					Name: "Overview",
					Series: []MetricGraphSeries***REMOVED***
						***REMOVED***
							Name:   "VUs",
							Metric: "vus",
							Kind:   AreaGraphSeriesType,
							Color:  "#808080",
						***REMOVED***,
						***REMOVED***
							Name:   "Request Rate",
							Metric: "http_reqs",
							Kind:   LineGraphSeriesType,
							Color:  "#0096FF",
						***REMOVED***,
						***REMOVED***
							Name:   "Request Duration",
							Metric: "http_req_duration",
							Kind:   LineGraphSeriesType,
							Color:  "#FF00FF",
						***REMOVED***,
						***REMOVED***
							Name:   "Failure Rate",
							Metric: "http_req_failed",
							Kind:   LineGraphSeriesType,
							Color:  "#FF0000",
						***REMOVED***,
					***REMOVED***,
					DesiredWidth: 3,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		true,
	***REMOVED***
***REMOVED***

const (
	AreaGraphSeriesType   = "area"
	LineGraphSeriesType   = "line"
	ColumnGraphSeriesType = "column"
)

type MetricGraphSeries struct ***REMOVED***
	Name   string `json:"name"`
	Metric string `json:"metric"`
	Kind   string `json:"kind"`
	Color  string `json:"color"`
***REMOVED***

type MetricGraph struct ***REMOVED***
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Series       []MetricGraphSeries `json:"series"`
	DesiredWidth int                 `json:"desiredWidth"`
***REMOVED***

type OutputConfig struct ***REMOVED***
	Graphs []MetricGraph `json:"graphs"`
***REMOVED***

type NullOutputConfig struct ***REMOVED***
	Value OutputConfig
	Valid bool
***REMOVED***

func NewNullOutputConfig(outputConfig OutputConfig, valid bool) NullOutputConfig ***REMOVED***
	return NullOutputConfig***REMOVED***outputConfig, valid***REMOVED***
***REMOVED***

func NullOutputConfigFrom(outputConfig OutputConfig) NullOutputConfig ***REMOVED***
	return NullOutputConfig***REMOVED***outputConfig, true***REMOVED***
***REMOVED***

func (oc NullOutputConfig) MarshalJSON() ([]byte, error) ***REMOVED***
	if !oc.Valid ***REMOVED***
		return []byte(`null`), nil
	***REMOVED***
	return json.Marshal(oc.Value)
***REMOVED***

func (oc *NullOutputConfig) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, []byte(`null`)) ***REMOVED***
		oc.Valid = false
		return nil
	***REMOVED***

	if err := json.Unmarshal(data, &oc.Value); err != nil ***REMOVED***
		return err
	***REMOVED***
	oc.Valid = true
	return nil
***REMOVED***

func IsValidSeriesKind(kind string) bool ***REMOVED***
	return kind == AreaGraphSeriesType || kind == LineGraphSeriesType || kind == ColumnGraphSeriesType
***REMOVED***

func IsBuiltinMetric(metricName string) bool ***REMOVED***
	return metricName == "data_received" ||
		metricName == "data_sent" ||
		metricName == "http_req_blocked" ||
		metricName == "http_req_connecting" ||
		metricName == "http_req_duration" ||
		metricName == "http_req_failed" ||
		metricName == "http_req_receiving" ||
		metricName == "http_req_sending" ||
		metricName == "http_req_tls_handshaking" ||
		metricName == "http_req_waiting" ||
		metricName == "http_reqs" ||
		metricName == "iteration_duration" ||
		metricName == "iterations" ||
		metricName == "vus" ||
		metricName == "vus_max"
***REMOVED***

// Ensures color is a valid hex color
func ValidSeriesColor(color string) bool ***REMOVED***
	if len(color) != 7 ***REMOVED***
		return false
	***REMOVED***

	if color[0] != '#' ***REMOVED***
		return false
	***REMOVED***

	for _, c := range color[1:] ***REMOVED***
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***
