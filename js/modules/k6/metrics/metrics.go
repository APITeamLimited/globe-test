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
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type Metric struct ***REMOVED***
	metric *metrics.Metric
	vu     modules.VU
***REMOVED***

// ErrMetricsAddInInitContext is error returned when adding to metric is done in the init context
var ErrMetricsAddInInitContext = common.NewInitContextError("Adding to metrics in the init context is not supported")

func (mi *ModuleInstance) newMetric(call goja.ConstructorCall, t metrics.MetricType) (*goja.Object, error) ***REMOVED***
	initEnv := mi.vu.InitEnv()
	if initEnv == nil ***REMOVED***
		return nil, errors.New("metrics must be declared in the init context")
	***REMOVED***
	rt := mi.vu.Runtime()
	c, _ := goja.AssertFunction(rt.ToValue(func(name string, isTime ...bool) (*goja.Object, error) ***REMOVED***
		valueType := metrics.Default
		if len(isTime) > 0 && isTime[0] ***REMOVED***
			valueType = metrics.Time
		***REMOVED***
		m, err := initEnv.Registry.NewMetric(name, t, valueType)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		metric := &Metric***REMOVED***metric: m, vu: mi.vu***REMOVED***
		o := rt.NewObject()
		err = o.DefineDataProperty("name", rt.ToValue(name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if err = o.Set("add", rt.ToValue(metric.add)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return o, nil
	***REMOVED***))
	v, err := c(call.This, call.Arguments...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return v.ToObject(rt), nil
***REMOVED***

const warnMessageValueMaxSize = 100

func limitValue(v string) string ***REMOVED***
	vRunes := []rune(v)
	if len(vRunes) < warnMessageValueMaxSize ***REMOVED***
		return v
	***REMOVED***
	difference := int64(len(vRunes) - warnMessageValueMaxSize)
	omitMsg := append(strconv.AppendInt([]byte("... omitting "), difference, 10), " characters ..."...)
	return strings.Join([]string***REMOVED***
		string(vRunes[:warnMessageValueMaxSize/2]),
		string(vRunes[len(vRunes)-warnMessageValueMaxSize/2:]),
	***REMOVED***, string(omitMsg))
***REMOVED***

func (m Metric) add(v goja.Value, addTags ...map[string]string) (bool, error) ***REMOVED***
	state := m.vu.State()
	if state == nil ***REMOVED***
		return false, ErrMetricsAddInInitContext
	***REMOVED***

	// return/throw exception if throw enabled, otherwise just log
	raiseErr := func(err error) (bool, error) ***REMOVED*** //nolint:unparam // we want to just do `return raiseErr(...)`
		if state.Options.Throw.Bool ***REMOVED***
			return false, err
		***REMOVED***
		state.Logger.Warn(err)
		return false, nil
	***REMOVED***
	raiseNan := func() (bool, error) ***REMOVED***
		return raiseErr(fmt.Errorf("'%s' is an invalid value for metric '%s', a number or a boolean value is expected",
			limitValue(v.String()), m.metric.Name))
	***REMOVED***

	if v == nil ***REMOVED***
		return raiseErr(fmt.Errorf("no value was provided for metric '%s', a number or a boolean value is expected",
			m.metric.Name))
	***REMOVED***
	if goja.IsNull(v) ***REMOVED***
		return raiseNan()
	***REMOVED***

	vfloat := v.ToFloat()
	if vfloat == 0 && v.ToBoolean() ***REMOVED***
		vfloat = 1.0
	***REMOVED***

	if math.IsNaN(vfloat) ***REMOVED***
		return raiseNan()
	***REMOVED***

	tags := state.CloneTags()
	for _, ts := range addTags ***REMOVED***
		for k, v := range ts ***REMOVED***
			tags[k] = v
		***REMOVED***
	***REMOVED***

	sample := metrics.Sample***REMOVED***Time: time.Now(), Metric: m.metric, Value: vfloat, Tags: metrics.IntoSampleTags(&tags)***REMOVED***
	metrics.PushIfNotDone(m.vu.Context(), state.Samples, sample)
	return true, nil
***REMOVED***

type (
	// RootModule is the root metrics module
	RootModule struct***REMOVED******REMOVED***
	// ModuleInstance represents an instance of the metrics module
	ModuleInstance struct ***REMOVED***
		vu modules.VU
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &ModuleInstance***REMOVED******REMOVED***
)

// NewModuleInstance implements modules.Module interface
func (*RootModule) NewModuleInstance(m modules.VU) modules.Instance ***REMOVED***
	return &ModuleInstance***REMOVED***vu: m***REMOVED***
***REMOVED***

// New returns a new RootModule.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// Exports returns the exports of the metrics module
func (mi *ModuleInstance) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"Counter": mi.XCounter,
			"Gauge":   mi.XGauge,
			"Trend":   mi.XTrend,
			"Rate":    mi.XRate,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// XCounter is a counter constructor
func (mi *ModuleInstance) XCounter(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, metrics.Counter)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***

// XGauge is a gauge constructor
func (mi *ModuleInstance) XGauge(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, metrics.Gauge)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***

// XTrend is a trend constructor
func (mi *ModuleInstance) XTrend(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, metrics.Trend)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***

// XRate is a rate constructor
func (mi *ModuleInstance) XRate(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, metrics.Rate)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***
