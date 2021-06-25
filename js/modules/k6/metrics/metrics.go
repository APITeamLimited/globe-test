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

package metrics

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/stats"
)

var nameRegexString = "^[\\p***REMOVED***L***REMOVED***\\p***REMOVED***N***REMOVED***\\._ !\\?/&#\\(\\)<>%-]***REMOVED***1,128***REMOVED***$"

var compileNameRegex = regexp.MustCompile(nameRegexString)

func checkName(name string) bool ***REMOVED***
	return compileNameRegex.Match([]byte(name))
***REMOVED***

type Metric struct ***REMOVED***
	metric *stats.Metric
***REMOVED***

// ErrMetricsAddInInitContext is error returned when adding to metric is done in the init context
var ErrMetricsAddInInitContext = common.NewInitContextError("Adding to metrics in the init context is not supported")

func newMetric(ctxPtr *context.Context, name string, t stats.MetricType, isTime []bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if lib.GetState(*ctxPtr) != nil ***REMOVED***
		return nil, errors.New("metrics must be declared in the init context")
	***REMOVED***

	// TODO: move verification outside the JS
	if !checkName(name) ***REMOVED***
		return nil, common.NewInitContextError(fmt.Sprintf("Invalid metric name: '%s'", name))
	***REMOVED***

	valueType := stats.Default
	if len(isTime) > 0 && isTime[0] ***REMOVED***
		valueType = stats.Time
	***REMOVED***

	rt := common.GetRuntime(*ctxPtr)
	bound := common.Bind(rt, Metric***REMOVED***stats.New(name, t, valueType)***REMOVED***, ctxPtr)
	o := rt.NewObject()
	err := o.DefineDataProperty("name", rt.ToValue(name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = o.Set("add", rt.ToValue(bound["add"])); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return o, nil
***REMOVED***

func (m Metric) Add(ctx context.Context, v goja.Value, addTags ...map[string]string) (bool, error) ***REMOVED***
	state := lib.GetState(ctx)
	if state == nil ***REMOVED***
		return false, ErrMetricsAddInInitContext
	***REMOVED***

	tags := state.CloneTags()
	for _, ts := range addTags ***REMOVED***
		for k, v := range ts ***REMOVED***
			tags[k] = v
		***REMOVED***
	***REMOVED***

	vfloat := v.ToFloat()
	if vfloat == 0 && v.ToBoolean() ***REMOVED***
		vfloat = 1.0
	***REMOVED***

	sample := stats.Sample***REMOVED***Time: time.Now(), Metric: m.metric, Value: vfloat, Tags: stats.IntoSampleTags(&tags)***REMOVED***
	stats.PushIfNotDone(ctx, state.Samples, sample)
	return true, nil
***REMOVED***

type Metrics struct***REMOVED******REMOVED***

func New() *Metrics ***REMOVED***
	return &Metrics***REMOVED******REMOVED***
***REMOVED***

func (*Metrics) XCounter(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Counter, isTime)
***REMOVED***

func (*Metrics) XGauge(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Gauge, isTime)
***REMOVED***

func (*Metrics) XTrend(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Trend, isTime)
***REMOVED***

func (*Metrics) XRate(ctx *context.Context, name string, isTime ...bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return newMetric(ctx, name, stats.Rate, isTime)
***REMOVED***
