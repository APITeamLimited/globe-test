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
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

					ctxPtr := new(context.Context)
					*ctxPtr = common.WithRuntime(context.Background(), rt)
					rt.Set("metrics", common.Bind(rt, New(), ctxPtr))

					root, _ := lib.NewGroup("", nil)
					child, _ := root.Group("child")
					samples := make(chan stats.SampleContainer, 1000)
					state := &lib.State***REMOVED***
						Options: lib.Options***REMOVED***SystemTags: stats.NewSystemTagSet(stats.TagGroup)***REMOVED***,
						Group:   root,
						Samples: samples,
					***REMOVED***

					isTimeString := ""
					if isTime ***REMOVED***
						isTimeString = `, true`
					***REMOVED***
					_, err := common.RunString(rt,
						fmt.Sprintf(`let m = new metrics.%s("my_metric"%s)`, fn, isTimeString),
					)
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***

					t.Run("ExitInit", func(t *testing.T) ***REMOVED***
						*ctxPtr = lib.WithState(*ctxPtr, state)
						_, err := common.RunString(rt, fmt.Sprintf(`new metrics.%s("my_metric")`, fn))
						assert.EqualError(t, err, "GoError: metrics must be declared in the init context at apply (native)")
					***REMOVED***)

					groups := map[string]*lib.Group***REMOVED***
						"Root":  root,
						"Child": child,
					***REMOVED***
					for name, g := range groups ***REMOVED***
						t.Run(name, func(t *testing.T) ***REMOVED***
							state.Group = g
							for name, val := range values ***REMOVED***
								t.Run(name, func(t *testing.T) ***REMOVED***
									t.Run("Simple", func(t *testing.T) ***REMOVED***
										_, err := common.RunString(rt, fmt.Sprintf(`m.add(%v)`, val.JS))
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
										_, err := common.RunString(rt, fmt.Sprintf(`m.add(%v, ***REMOVED***a:1***REMOVED***)`, val.JS))
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

func TestMetricNames(t *testing.T) ***REMOVED***
	t.Parallel()
	var testMap = map[string]bool***REMOVED***
		"simple":       true,
		"still_simple": true,
		"":             false,
		"@":            false,
		"a":            true,
		"special\n\t":  false,
		// this has both hangul and japanese numerals .
		"hello.World_in_한글一안녕一세상": true,
		// too long
		"tooolooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooog": false,
	***REMOVED***

	for key, value := range testMap ***REMOVED***
		t.Run(key, func(t *testing.T) ***REMOVED***
			assert.Equal(t, value, checkName(key), key)
		***REMOVED***)
	***REMOVED***
***REMOVED***
