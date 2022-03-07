/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package json

import (
	"time"

	"go.k6.io/k6/stats"
)

//go:generate easyjson -pkg -no_std_marshalers -gen_build_flags -mod=mod .

//easyjson:json
type sampleEnvelope struct ***REMOVED***
	Type string `json:"type"`
	Data struct ***REMOVED***
		Time  time.Time         `json:"time"`
		Value float64           `json:"value"`
		Tags  *stats.SampleTags `json:"tags"`
	***REMOVED*** `json:"data"`
	Metric string `json:"metric"`
***REMOVED***

// wrapSample is used to package a metric sample in a way that's nice to export
// to JSON.
func wrapSample(sample stats.Sample) sampleEnvelope ***REMOVED***
	s := sampleEnvelope***REMOVED***
		Type:   "Point",
		Metric: sample.Metric.Name,
	***REMOVED***
	s.Data.Time = sample.Time
	s.Data.Value = sample.Value
	s.Data.Tags = sample.Tags
	return s
***REMOVED***

//easyjson:json
type metricEnvelope struct ***REMOVED***
	Type   string        `json:"type"`
	Data   *stats.Metric `json:"data"`
	Metric string        `json:"metric"`
***REMOVED***

func wrapMetric(metric *stats.Metric) metricEnvelope ***REMOVED***
	return metricEnvelope***REMOVED***
		Type:   "Metric",
		Metric: metric.Name,
		Data:   metric,
	***REMOVED***
***REMOVED***
