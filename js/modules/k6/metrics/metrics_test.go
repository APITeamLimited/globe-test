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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/stats"
)

type addTestValue struct ***REMOVED***
	JS      string
	Float   float64
	isError bool
***REMOVED***

type addTest struct ***REMOVED***
	val          addTestValue
	rt           *goja.Runtime
	hook         *testutils.SimpleLogrusHook
	samples      chan stats.SampleContainer
	isThrow      bool
	mtyp         stats.MetricType
	valueType    stats.ValueType
	js           string
	expectedTags map[string]string
***REMOVED***

func (a addTest) run(t *testing.T) ***REMOVED***
	_, err := a.rt.RunString(a.js)
	if a.val.isError && a.isThrow ***REMOVED***
		if assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		assert.NoError(t, err)
		if a.val.isError && !a.isThrow ***REMOVED***
			lines := a.hook.Drain()
			require.Len(t, lines, 1)
			assert.Contains(t, lines[0].Message, "is an invalid value for metric")
			return
		***REMOVED***
	***REMOVED***
	bufSamples := stats.GetBufferedSamples(a.samples)
	if assert.Len(t, bufSamples, 1) ***REMOVED***
		sample, ok := bufSamples[0].(stats.Sample)
		require.True(t, ok)

		assert.NotZero(t, sample.Time)
		assert.Equal(t, a.val.Float, sample.Value)
		assert.Equal(t, a.expectedTags, sample.Tags.CloneTags())
		assert.Equal(t, "my_metric", sample.Metric.Name)
		assert.Equal(t, a.mtyp, sample.Metric.Type)
		assert.Equal(t, a.valueType, sample.Metric.Contains)
	***REMOVED***
***REMOVED***

func TestMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	types := map[string]stats.MetricType***REMOVED***
		"Counter": stats.Counter,
		"Gauge":   stats.Gauge,
		"Trend":   stats.Trend,
		"Rate":    stats.Rate,
	***REMOVED***
	values := map[string]addTestValue***REMOVED***
		"Float":                 ***REMOVED***JS: `2.5`, Float: 2.5***REMOVED***,
		"Int":                   ***REMOVED***JS: `5`, Float: 5.0***REMOVED***,
		"True":                  ***REMOVED***JS: `true`, Float: 1.0***REMOVED***,
		"False":                 ***REMOVED***JS: `false`, Float: 0.0***REMOVED***,
		"null":                  ***REMOVED***JS: `null`, isError: true***REMOVED***,
		"undefined":             ***REMOVED***JS: `undefined`, isError: true***REMOVED***,
		"NaN":                   ***REMOVED***JS: `NaN`, isError: true***REMOVED***,
		"string":                ***REMOVED***JS: `"string"`, isError: true***REMOVED***,
		"string 5":              ***REMOVED***JS: `"5.3"`, Float: 5.3***REMOVED***,
		"some object":           ***REMOVED***JS: `***REMOVED***something: 3***REMOVED***`, isError: true***REMOVED***,
		"another metric object": ***REMOVED***JS: `m`, isError: true***REMOVED***,
	***REMOVED***
	for fn, mtyp := range types ***REMOVED***
		fn, mtyp := fn, mtyp
		t.Run(fn, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for isTime, valueType := range map[bool]stats.ValueType***REMOVED***false: stats.Default, true: stats.Time***REMOVED*** ***REMOVED***
				isTime, valueType := isTime, valueType
				t.Run(fmt.Sprintf("isTime=%v", isTime), func(t *testing.T) ***REMOVED***
					t.Parallel()
					test := addTest***REMOVED***
						mtyp:      mtyp,
						valueType: valueType,
					***REMOVED***
					test.rt = goja.New()
					test.rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
					mii := &modulestest.VU***REMOVED***
						RuntimeField: test.rt,
						InitEnvField: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
						CtxField:     context.Background(),
					***REMOVED***
					m, ok := New().NewModuleInstance(mii).(*ModuleInstance)
					require.True(t, ok)
					require.NoError(t, test.rt.Set("metrics", m.Exports().Named))
					test.samples = make(chan stats.SampleContainer, 1000)
					state := &lib.State***REMOVED***
						Options: lib.Options***REMOVED******REMOVED***,
						Samples: test.samples,
						Tags: lib.NewTagMap(map[string]string***REMOVED***
							"key": "value",
						***REMOVED***),
					***REMOVED***

					isTimeString := ""
					if isTime ***REMOVED***
						isTimeString = `, true`
					***REMOVED***
					_, err := test.rt.RunString(fmt.Sprintf(`var m = new metrics.%s("my_metric"%s)`, fn, isTimeString))
					require.NoError(t, err)

					t.Run("ExitInit", func(t *testing.T) ***REMOVED***
						mii.StateField = state
						mii.InitEnvField = nil
						_, err := test.rt.RunString(fmt.Sprintf(`new metrics.%s("my_metric")`, fn))
						assert.Contains(t, err.Error(), "metrics must be declared in the init context")
					***REMOVED***)
					mii.StateField = state
					logger := logrus.New()
					logger.Out = ioutil.Discard
					test.hook = &testutils.SimpleLogrusHook***REMOVED***HookedLevels: logrus.AllLevels***REMOVED***
					logger.AddHook(test.hook)
					state.Logger = logger

					for name, val := range values ***REMOVED***
						test.val = val
						for _, isThrow := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
							state.Options.Throw.Bool = isThrow
							test.isThrow = isThrow
							t.Run(fmt.Sprintf("%s/isThrow=%v/Simple", name, isThrow), func(t *testing.T) ***REMOVED***
								test.js = fmt.Sprintf(`m.add(%v)`, val.JS)
								test.expectedTags = map[string]string***REMOVED***"key": "value"***REMOVED***
								test.run(t)
							***REMOVED***)
							t.Run(fmt.Sprintf("%s/isThrow=%v/Tags", name, isThrow), func(t *testing.T) ***REMOVED***
								test.js = fmt.Sprintf(`m.add(%v, ***REMOVED***a:1***REMOVED***)`, val.JS)
								test.expectedTags = map[string]string***REMOVED***"key": "value", "a": "1"***REMOVED***
								test.run(t)
							***REMOVED***)
						***REMOVED***
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMetricGetName(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	mii := &modulestest.VU***REMOVED***
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
		CtxField:     context.Background(),
	***REMOVED***
	m, ok := New().NewModuleInstance(mii).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("metrics", m.Exports().Named))
	v, err := rt.RunString(`
		var m = new metrics.Counter("my_metric")
		m.name
	`)
	require.NoError(t, err)
	require.Equal(t, "my_metric", v.String())

	_, err = rt.RunString(`
		"use strict";
		m.name = "something"
	`)
	require.Error(t, err)
	require.Contains(t, err.Error(), "TypeError: Cannot assign to read only property 'name'")
***REMOVED***

func TestMetricDuplicates(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	mii := &modulestest.VU***REMOVED***
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
		CtxField:     context.Background(),
	***REMOVED***
	m, ok := New().NewModuleInstance(mii).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("metrics", m.Exports().Named))
	_, err := rt.RunString(`
		var m = new metrics.Counter("my_metric")
	`)
	require.NoError(t, err)

	_, err = rt.RunString(`
		var m2 = new metrics.Counter("my_metric")
	`)
	require.NoError(t, err)

	_, err = rt.RunString(`
		var m3 = new metrics.Gauge("my_metric")
	`)
	require.Error(t, err)

	_, err = rt.RunString(`
		var m4 = new metrics.Counter("my_metric", true)
	`)
	require.Error(t, err)

	v, err := rt.RunString(`
		m.name == m2.name && m.name == "my_metric" && m3 === undefined && m4 === undefined
	`)
	require.NoError(t, err)

	require.True(t, v.ToBoolean())
***REMOVED***
