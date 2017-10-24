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
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/stats"
)

type Metric struct ***REMOVED***
	metric *stats.Metric
***REMOVED***

func newMetric(ctxPtr *context.Context, name string, t stats.MetricType, isTime []bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if common.GetState(*ctxPtr) != nil ***REMOVED***
		return nil, errors.New("Metrics must be declared in the init context")
	***REMOVED***

	valueType := stats.Default
	if len(isTime) > 0 && isTime[0] ***REMOVED***
		valueType = stats.Time
	***REMOVED***

	rt := common.GetRuntime(*ctxPtr)
	return common.Bind(rt, Metric***REMOVED***stats.New(name, t, valueType)***REMOVED***, ctxPtr), nil
***REMOVED***

func (m Metric) Add(ctx context.Context, v goja.Value, addTags ...map[string]string) ***REMOVED***
	state := common.GetState(ctx)

	tags := map[string]string***REMOVED***
		"group": state.Group.Path,
	***REMOVED***
	for _, ts := range addTags ***REMOVED***
		for k, v := range ts ***REMOVED***
			tags[k] = v
		***REMOVED***
	***REMOVED***

	vfloat := v.ToFloat()
	if vfloat == 0 && v.ToBoolean() ***REMOVED***
		vfloat = 1.0
	***REMOVED***

	state.Samples = append(state.Samples,
		stats.Sample***REMOVED***Time: time.Now(), Metric: m.metric, Value: vfloat, Tags: tags***REMOVED***,
	)
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
