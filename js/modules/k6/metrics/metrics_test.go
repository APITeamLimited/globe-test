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
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/stats"
)

func TestMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	types := map[string]stats.MetricType***REMOVED***
		"Counter": stats.Counter,
		"Gauge":   stats.Gauge,
		"Trend":   stats.Trend,
		"Rate":    stats.Rate,
	***REMOVED***
	values := map[string]struct ***REMOVED***
		JS    string
		Float float64
	***REMOVED******REMOVED***
		"Float": ***REMOVED***`2.5`, 2.5***REMOVED***,
		"Int":   ***REMOVED***`5`, 5.0***REMOVED***,
		"True":  ***REMOVED***`true`, 1.0***REMOVED***,
		"False": ***REMOVED***`false`, 0.0***REMOVED***,
	***REMOVED***
	for fn, mtyp := range types ***REMOVED***
		fn, mtyp := fn, mtyp
		t.Run(fn, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for isTime, valueType := range map[bool]stats.ValueType***REMOVED***false: stats.Default, true: stats.Time***REMOVED*** ***REMOVED***
				isTime, valueType := isTime, valueType
				t.Run(fmt.Sprintf("isTime=%v", isTime), func(t *testing.T) ***REMOVED***
					t.Parallel()
					rt := goja.New()
					rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
					mii := &modulestest.InstanceCore***REMOVED***
						Runtime: rt,
						InitEnv: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
						Ctx:     context.Background(),
					***REMOVED***
					m, ok := New().NewModuleInstance(mii).(*ModuleInstance)
					require.True(t, ok)
					require.NoError(t, rt.Set("metrics", m.GetExports().Named))
					root, _ := lib.NewGroup("", nil)
					child, _ := root.Group("child")
					samples := make(chan stats.SampleContainer, 1000)
					state := &lib.State***REMOVED***
						Options: lib.Options***REMOVED***SystemTags: stats.NewSystemTagSet(stats.TagGroup), Throw: null.BoolFrom(true)***REMOVED***,
						Group:   root,
						Samples: samples,
						Tags:    map[string]string***REMOVED***"group": root.Path***REMOVED***,
					***REMOVED***

					isTimeString := ""
					if isTime ***REMOVED***
						isTimeString = `, true`
					***REMOVED***
					_, err := rt.RunString(fmt.Sprintf(`var m = new metrics.%s("my_metric"%s)`, fn, isTimeString))
					require.NoError(t, err)

					t.Run("ExitInit", func(t *testing.T) ***REMOVED***
						mii.State = state
						mii.InitEnv = nil
						_, err := rt.RunString(fmt.Sprintf(`new metrics.%s("my_metric")`, fn))
						assert.Contains(t, err.Error(), "metrics must be declared in the init context")
					***REMOVED***)

					groups := map[string]*lib.Group***REMOVED***
						"Root":  root,
						"Child": child,
					***REMOVED***
					for name, g := range groups ***REMOVED***
						name, g := name, g
						t.Run(name, func(t *testing.T) ***REMOVED***
							state.Group = g
							state.Tags["group"] = g.Path
							for name, val := range values ***REMOVED***
								name, val := name, val
								t.Run(name, func(t *testing.T) ***REMOVED***
									t.Run("Simple", func(t *testing.T) ***REMOVED***
										_, err := rt.RunString(fmt.Sprintf(`m.add(%v)`, val.JS))
										assert.NoError(t, err)
										bufSamples := stats.GetBufferedSamples(samples)
										if assert.Len(t, bufSamples, 1) ***REMOVED***
											sample, ok := bufSamples[0].(stats.Sample)
											require.True(t, ok)

											assert.NotZero(t, sample.Time)
											assert.Equal(t, sample.Value, val.Float)
											assert.Equal(t, map[string]string***REMOVED***
												"group": g.Path,
											***REMOVED***, sample.Tags.CloneTags())
											assert.Equal(t, "my_metric", sample.Metric.Name)
											assert.Equal(t, mtyp, sample.Metric.Type)
											assert.Equal(t, valueType, sample.Metric.Contains)
										***REMOVED***
									***REMOVED***)
									t.Run("Tags", func(t *testing.T) ***REMOVED***
										_, err := rt.RunString(fmt.Sprintf(`m.add(%v, ***REMOVED***a:1***REMOVED***)`, val.JS))
										assert.NoError(t, err)
										bufSamples := stats.GetBufferedSamples(samples)
										if assert.Len(t, bufSamples, 1) ***REMOVED***
											sample, ok := bufSamples[0].(stats.Sample)
											require.True(t, ok)

											assert.NotZero(t, sample.Time)
											assert.Equal(t, sample.Value, val.Float)
											assert.Equal(t, map[string]string***REMOVED***
												"group": g.Path,
												"a":     "1",
											***REMOVED***, sample.Tags.CloneTags())
											assert.Equal(t, "my_metric", sample.Metric.Name)
											assert.Equal(t, mtyp, sample.Metric.Type)
											assert.Equal(t, valueType, sample.Metric.Contains)
										***REMOVED***
									***REMOVED***)
								***REMOVED***)
							***REMOVED***
						***REMOVED***)
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

	mii := &modulestest.InstanceCore***REMOVED***
		Runtime: rt,
		InitEnv: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
		Ctx:     context.Background(),
	***REMOVED***
	m, ok := New().NewModuleInstance(mii).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("metrics", m.GetExports().Named))
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

	mii := &modulestest.InstanceCore***REMOVED***
		Runtime: rt,
		InitEnv: &common.InitEnvironment***REMOVED***Registry: metrics.NewRegistry()***REMOVED***,
		Ctx:     context.Background(),
	***REMOVED***
	m, ok := New().NewModuleInstance(mii).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("metrics", m.GetExports().Named))
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
