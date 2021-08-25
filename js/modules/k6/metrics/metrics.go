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
	"regexp"
	"time"

	"github.com/dop251/goja"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/stats"
)

var nameRegexString = "^[\\p***REMOVED***L***REMOVED***\\p***REMOVED***N***REMOVED***\\._ !\\?/&#\\(\\)<>%-]***REMOVED***1,128***REMOVED***$"

var compileNameRegex = regexp.MustCompile(nameRegexString)

func checkName(name string) bool ***REMOVED***
	return compileNameRegex.Match([]byte(name))
***REMOVED***

type Metric struct ***REMOVED***
	metric *stats.Metric
	core   modules.InstanceCore
***REMOVED***

// ErrMetricsAddInInitContext is error returned when adding to metric is done in the init context
var ErrMetricsAddInInitContext = common.NewInitContextError("Adding to metrics in the init context is not supported")

func (mi *ModuleInstance) newMetric(call goja.ConstructorCall, t stats.MetricType) (*goja.Object, error) ***REMOVED***
	if mi.GetInitEnv() == nil ***REMOVED***
		return nil, errors.New("metrics must be declared in the init context")
	***REMOVED***
	rt := mi.GetRuntime()
	c, _ := goja.AssertFunction(rt.ToValue(func(name string, isTime ...bool) (*goja.Object, error) ***REMOVED***
		// TODO: move verification outside the JS
		if !checkName(name) ***REMOVED***
			return nil, common.NewInitContextError(fmt.Sprintf("Invalid metric name: '%s'", name))
		***REMOVED***

		valueType := stats.Default
		if len(isTime) > 0 && isTime[0] ***REMOVED***
			valueType = stats.Time
		***REMOVED***
		m := stats.New(name, t, valueType)

		metric := &Metric***REMOVED***metric: m, core: mi.InstanceCore***REMOVED***
		o := rt.NewObject()
		err := o.DefineDataProperty("name", rt.ToValue(name), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
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

func (m Metric) add(v goja.Value, addTags ...map[string]string) (bool, error) ***REMOVED***
	state := m.core.GetState()
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
	stats.PushIfNotDone(m.core.GetContext(), state.Samples, sample)
	return true, nil
***REMOVED***

type (
	// RootModule is the root metrics module
	RootModule struct***REMOVED******REMOVED***
	// ModuleInstance represents an instance of the metrics module
	ModuleInstance struct ***REMOVED***
		modules.InstanceCore
	***REMOVED***
)

var (
	_ modules.IsModuleV2 = &RootModule***REMOVED******REMOVED***
	_ modules.Instance   = &ModuleInstance***REMOVED******REMOVED***
)

// NewModuleInstance implements modules.IsModuleV2 interface
func (*RootModule) NewModuleInstance(m modules.InstanceCore) modules.Instance ***REMOVED***
	return &ModuleInstance***REMOVED***InstanceCore: m***REMOVED***
***REMOVED***

// New returns a new RootModule.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// GetExports returns the exports of the metrics module
func (mi *ModuleInstance) GetExports() modules.Exports ***REMOVED***
	return modules.GenerateExports(mi)
***REMOVED***

// XCounter is a counter constructor
func (mi *ModuleInstance) XCounter(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, stats.Counter)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***

// XGauge is a gauge constructor
func (mi *ModuleInstance) XGauge(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, stats.Gauge)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***

// XTrend is a trend constructor
func (mi *ModuleInstance) XTrend(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, stats.Trend)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***

// XRate is a rate constructor
func (mi *ModuleInstance) XRate(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object ***REMOVED***
	v, err := mi.newMetric(call, stats.Rate)
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***
	return v
***REMOVED***
