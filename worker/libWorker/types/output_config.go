package types

import (
	"bytes"
	"encoding/json"
)

func DefaultOutputConfig() NullOutputConfig {
	return NullOutputConfig{
		OutputConfig{
			Graphs: []MetricGraph{
				{
					Name: "Overview",
					Series: []MetricGraphSeries{
						{
							LoadZone: "global",
							Metric:   "vus",
							Kind:     AreaGraphSeriesType,
							Color:    "#808080",
						},
						{
							LoadZone: "global",
							Metric:   "http_reqs",
							Kind:     LineGraphSeriesType,
							Color:    "#0096FF",
						},
						{
							LoadZone: "global",
							Metric:   "http_req_duration",
							Kind:     LineGraphSeriesType,
							Color:    "#FF00FF",
						},
						{
							LoadZone: "global",
							Metric:   "http_req_failed",
							Kind:     LineGraphSeriesType,
							Color:    "#FF0000",
						},
					},
					DesiredWidth: 3,
				},
			},
		},
		true,
	}
}

const (
	AreaGraphSeriesType   = "area"
	LineGraphSeriesType   = "line"
	ColumnGraphSeriesType = "column"
)

type MetricGraphSeries struct {
	LoadZone string `json:"loadZone"`
	Metric   string `json:"metric"`
	Kind     string `json:"kind"`
	Color    string `json:"color"`
}

type MetricGraph struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Series       []MetricGraphSeries `json:"series"`
	DesiredWidth int                 `json:"desiredWidth"`
}

type OutputConfig struct {
	Graphs []MetricGraph `json:"graphs"`
}

type NullOutputConfig struct {
	Value OutputConfig
	Valid bool
}

func NewNullOutputConfig(outputConfig OutputConfig, valid bool) NullOutputConfig {
	return NullOutputConfig{outputConfig, valid}
}

func NullOutputConfigFrom(outputConfig OutputConfig) NullOutputConfig {
	return NullOutputConfig{outputConfig, true}
}

func (oc NullOutputConfig) MarshalJSON() ([]byte, error) {
	if !oc.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(oc.Value)
}

func (oc *NullOutputConfig) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`null`)) {
		oc.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &oc.Value); err != nil {
		return err
	}
	oc.Valid = true
	return nil
}

func IsValidSeriesKind(kind string) bool {
	return kind == AreaGraphSeriesType || kind == LineGraphSeriesType || kind == ColumnGraphSeriesType
}

func IsBuiltinMetric(metricName string) bool {
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
}

// Ensures color is a valid hex color
func ValidSeriesColor(color string) bool {
	if len(color) != 7 {
		return false
	}

	if color[0] != '#' {
		return false
	}

	for _, c := range color[1:] {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}

	return true
}
