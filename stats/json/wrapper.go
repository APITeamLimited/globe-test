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

func WrapSample(sample *stats.Sample) *Envelope ***REMOVED***
	if sample == nil ***REMOVED***
		return nil
	***REMOVED***
	return &Envelope***REMOVED***
		Type:   "Point",
		Metric: sample.Metric.Name,
		Data:   NewJSONSample(sample),
	***REMOVED***
***REMOVED***

func WrapMetric(metric *stats.Metric) *Envelope ***REMOVED***
	if metric == nil ***REMOVED***
		return nil
	***REMOVED***

	return &Envelope***REMOVED***
		Type:   "Metric",
		Metric: metric.Name,
		Data:   metric,
	***REMOVED***
***REMOVED***
