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

package metrics

import (
	"fmt"
	"regexp"
	"sync"

	"go.k6.io/k6/stats"
)

// Registry is what can create metrics
type Registry struct ***REMOVED***
	metrics map[string]*stats.Metric
	l       sync.RWMutex
***REMOVED***

// NewRegistry returns a new registry
func NewRegistry() *Registry ***REMOVED***
	return &Registry***REMOVED***
		metrics: make(map[string]*stats.Metric),
	***REMOVED***
***REMOVED***

const nameRegexString = "^[\\p***REMOVED***L***REMOVED***\\p***REMOVED***N***REMOVED***\\._ !\\?/&#\\(\\)<>%-]***REMOVED***1,128***REMOVED***$"

var compileNameRegex = regexp.MustCompile(nameRegexString)

func checkName(name string) bool ***REMOVED***
	return compileNameRegex.Match([]byte(name))
***REMOVED***

// NewMetric returns new metric registered to this registry
// TODO have multiple versions returning specific metric types when we have such things
func (r *Registry) NewMetric(name string, typ stats.MetricType, t ...stats.ValueType) (*stats.Metric, error) ***REMOVED***
	r.l.Lock()
	defer r.l.Unlock()

	if !checkName(name) ***REMOVED***
		return nil, fmt.Errorf("Invalid metric name: '%s'", name) //nolint:golint,stylecheck
	***REMOVED***
	oldMetric, ok := r.metrics[name]

	if !ok ***REMOVED***
		m := newMetric(name, typ, t...)
		r.metrics[name] = m
		return m, nil
	***REMOVED***
	if oldMetric.Type != typ ***REMOVED***
		return nil, fmt.Errorf("metric '%s' already exists but with type %s, instead of %s", name, oldMetric.Type, typ)
	***REMOVED***
	if len(t) > 0 ***REMOVED***
		if t[0] != oldMetric.Contains ***REMOVED***
			return nil, fmt.Errorf("metric '%s' already exists but with a value type %s, instead of %s",
				name, oldMetric.Contains, t[0])
		***REMOVED***
	***REMOVED***
	return oldMetric, nil
***REMOVED***

// MustNewMetric is like NewMetric, but will panic if there is an error
func (r *Registry) MustNewMetric(name string, typ stats.MetricType, t ...stats.ValueType) *stats.Metric ***REMOVED***
	m, err := r.NewMetric(name, typ, t...)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return m
***REMOVED***

// Get returns the Metric with the given name. If that metric doesn't exist,
// Get() will return a nil value.
func (r *Registry) Get(name string) *stats.Metric ***REMOVED***
	return r.metrics[name]
***REMOVED***

func newMetric(name string, mt stats.MetricType, vt ...stats.ValueType) *stats.Metric ***REMOVED***
	valueType := stats.Default
	if len(vt) > 0 ***REMOVED***
		valueType = vt[0]
	***REMOVED***
	var sink stats.Sink
	switch mt ***REMOVED***
	case stats.Counter:
		sink = &stats.CounterSink***REMOVED******REMOVED***
	case stats.Gauge:
		sink = &stats.GaugeSink***REMOVED******REMOVED***
	case stats.Trend:
		sink = &stats.TrendSink***REMOVED******REMOVED***
	case stats.Rate:
		sink = &stats.RateSink***REMOVED******REMOVED***
	default:
		return nil
	***REMOVED***
	return &stats.Metric***REMOVED***
		Name:     name,
		Type:     mt,
		Contains: valueType,
		Sink:     sink,
	***REMOVED***
***REMOVED***
