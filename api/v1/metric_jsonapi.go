/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package v1

import (
	"time"

	"go.k6.io/k6/stats"
)

// MetricsJSONAPI is JSON API envelop for metrics
type MetricsJSONAPI struct ***REMOVED***
	Data []metricData `json:"data"`
***REMOVED***

type metricJSONAPI struct ***REMOVED***
	Data metricData `json:"data"`
***REMOVED***

type metricData struct ***REMOVED***
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Metric `json:"attributes"`
***REMOVED***

func newMetricEnvelope(m *stats.Metric, t time.Duration) metricJSONAPI ***REMOVED***
	return metricJSONAPI***REMOVED***
		Data: newMetricData(m, t),
	***REMOVED***
***REMOVED***

func newMetricsJSONAPI(list map[string]*stats.Metric, t time.Duration) MetricsJSONAPI ***REMOVED***
	metrics := make([]metricData, 0, len(list))

	for _, m := range list ***REMOVED***
		metrics = append(metrics, newMetricData(m, t))
	***REMOVED***

	return MetricsJSONAPI***REMOVED***
		Data: metrics,
	***REMOVED***
***REMOVED***

func newMetricData(m *stats.Metric, t time.Duration) metricData ***REMOVED***
	metric := NewMetric(m, t)

	return metricData***REMOVED***
		Type:       "metrics",
		ID:         metric.Name,
		Attributes: metric,
	***REMOVED***
***REMOVED***

// Metrics extract the []v1.Metric from the JSON API envelop
func (m MetricsJSONAPI) Metrics() []Metric ***REMOVED***
	list := make([]Metric, 0, len(m.Data))

	for _, metric := range m.Data ***REMOVED***
		m := metric.Attributes
		m.Name = metric.ID
		list = append(list, m)
	***REMOVED***

	return list
***REMOVED***
