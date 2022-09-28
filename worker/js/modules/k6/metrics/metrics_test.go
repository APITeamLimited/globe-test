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

	"github.com/APITeamLimited/k6-worker/js/common"
	"github.com/APITeamLimited/k6-worker/js/modulestest"
	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/testutils"
	"github.com/APITeamLimited/k6-worker/metrics"
)

type addTestValue struct ***REMOVED***
	JS     string
	Float  float64
	errStr string
	noTags bool
***REMOVED***

type addTest struct ***REMOVED***
	val          addTestValue
	rt           *goja.Runtime
	hook         *testutils.SimpleLogrusHook
	samples      chan metrics.SampleContainer
	isThrow      bool
	mtyp         metrics.MetricType
	valueType    metrics.ValueType
	js           string
	expectedTags map[string]string
***REMOVED***

func (a addTest) run(t *testing.T) ***REMOVED***
	_, err := a.rt.RunString(a.js)
	if len(a.val.errStr) != 0 && a.isThrow ***REMOVED***
		if assert.Error(t, err) ***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		assert.NoError(t, err)
		if len(a.val.errStr) != 0 && !a.isThrow ***REMOVED***
			lines := a.hook.Drain()
			require.Len(t, lines, 1)
			assert.Contains(t, lines[0].Message, a.val.errStr)
			return
		***REMOVED***
	***REMOVED***
	bufSamples := metrics.GetBufferedSamples(a.samples)
	if assert.Len(t, bufSamples, 1) ***REMOVED***
		sample, ok := bufSamples[0].(metrics.Sample)
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
	types := map[string]metrics.MetricType***REMOVED***
		"Counter": metrics.Counter,
		"Gauge":   metrics.Gauge,
		"Trend":   metrics.Trend,
		"Rate":    metrics.Rate,
	***REMOVED***
	values := map[string]addTestValue***REMOVED***
		"Float":                 ***REMOVED***JS: `2.5`, Float: 2.5***REMOVED***,
		"Int":                   ***REMOVED***JS: `5`, Float: 5.0***REMOVED***,
		"True":                  ***REMOVED***JS: `true`, Float: 1.0***REMOVED***,
		"False":                 ***REMOVED***JS: `false`, Float: 0.0***REMOVED***,
		"null":                  ***REMOVED***JS: `null`, errStr: "is an invalid value for metric"***REMOVED***,
		"undefined":             ***REMOVED***JS: `undefined`, errStr: "is an invalid value for metric"***REMOVED***,
		"NaN":                   ***REMOVED***JS: `NaN`, errStr: "is an invalid value for metric"***REMOVED***,
		"string":                ***REMOVED***JS: `"string"`, errStr: "is an invalid value for metric"***REMOVED***,
		"string 5":              ***REMOVED***JS: `"5.3"`, Float: 5.3***REMOVED***,
		"some object":           ***REMOVED***JS: `***REMOVED***something: 3***REMOVED***`, errStr: "is an invalid value for metric"***REMOVED***,
		"another metric object": ***REMOVED***JS: `m`, errStr: "is an invalid value for metric"***REMOVED***,
		"no argument":           ***REMOVED***JS: ``, errStr: "no value was provided", noTags: true***REMOVED***,
	***REMOVED***
	for fn, mtyp := range types ***REMOVED***
		fn, mtyp := fn, mtyp
		t.Run(fn, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for isTime, valueType := range map[bool]metrics.ValueType***REMOVED***false: metrics.Default, true: metrics.Time***REMOVED*** ***REMOVED***
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
					test.samples = make(chan metrics.SampleContainer, 1000)
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
							if !val.noTags ***REMOVED***
								t.Run(fmt.Sprintf("%s/isThrow=%v/Tags", name, isThrow), func(t *testing.T) ***REMOVED***
									test.js = fmt.Sprintf(`m.add(%v, ***REMOVED***a:1***REMOVED***)`, val.JS)
									test.expectedTags = map[string]string***REMOVED***"key": "value", "a": "1"***REMOVED***
									test.run(t)
								***REMOVED***)
							***REMOVED***
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
