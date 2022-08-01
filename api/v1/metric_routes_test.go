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

package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/metrics"
)

func TestGetMetrics(t *testing.T) ***REMOVED***
	t.Parallel()

	piState := getTestPreInitState(t)
	testMetric, err := piState.Registry.NewMetric("my_metric", metrics.Trend, metrics.Time)
	require.NoError(t, err)
	execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED******REMOVED***, piState)
	require.NoError(t, err)
	engine, err := core.NewEngine(execScheduler, lib.Options***REMOVED******REMOVED***, piState.RuntimeOptions, nil, piState.Logger, piState.Registry)
	require.NoError(t, err)

	engine.MetricsEngine.ObservedMetrics = map[string]*metrics.Metric***REMOVED***
		"my_metric": testMetric,
	***REMOVED***
	engine.MetricsEngine.ObservedMetrics["my_metric"].Tainted = null.BoolFrom(true)

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/metrics", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var doc MetricsJSONAPI
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		if !assert.NotNil(t, doc.Data) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, "metrics", doc.Data[0].Type)
	***REMOVED***)

	t.Run("metrics", func(t *testing.T) ***REMOVED***
		t.Parallel()

		var envelop MetricsJSONAPI
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &envelop))

		metricsData := envelop.Data
		if !assert.Len(t, metricsData, 1) ***REMOVED***
			return
		***REMOVED***

		metric := metricsData[0].Attributes

		assert.Equal(t, "my_metric", metricsData[0].ID)
		assert.True(t, metric.Type.Valid)
		assert.Equal(t, metrics.Trend, metric.Type.Type)
		assert.True(t, metric.Contains.Valid)
		assert.Equal(t, metrics.Time, metric.Contains.Type)
		assert.True(t, metric.Tainted.Valid)
		assert.True(t, metric.Tainted.Bool)

		resMetrics := envelop.Metrics()
		assert.Len(t, resMetrics, 1)
		assert.Equal(t, resMetrics[0].Name, "my_metric")
	***REMOVED***)
***REMOVED***

func TestGetMetric(t *testing.T) ***REMOVED***
	t.Parallel()

	piState := getTestPreInitState(t)
	testMetric, err := piState.Registry.NewMetric("my_metric", metrics.Trend, metrics.Time)
	require.NoError(t, err)
	execScheduler, err := local.NewExecutionScheduler(&minirunner.MiniRunner***REMOVED******REMOVED***, piState)
	require.NoError(t, err)
	engine, err := core.NewEngine(execScheduler, lib.Options***REMOVED******REMOVED***, piState.RuntimeOptions, nil, piState.Logger, piState.Registry)
	require.NoError(t, err)

	engine.MetricsEngine.ObservedMetrics = map[string]*metrics.Metric***REMOVED***
		"my_metric": testMetric,
	***REMOVED***
	engine.MetricsEngine.ObservedMetrics["my_metric"].Tainted = null.BoolFrom(true)

	t.Run("nonexistent", func(t *testing.T) ***REMOVED***
		t.Parallel()

		rw := httptest.NewRecorder()
		NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/metrics/notreal", nil))
		res := rw.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	***REMOVED***)

	t.Run("real", func(t *testing.T) ***REMOVED***
		t.Parallel()

		rw := httptest.NewRecorder()
		NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/metrics/my_metric", nil))
		res := rw.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		t.Run("document", func(t *testing.T) ***REMOVED***
			t.Parallel()

			var doc metricJSONAPI
			assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))

			assert.Equal(t, "metrics", doc.Data.Type)
		***REMOVED***)

		t.Run("metric", func(t *testing.T) ***REMOVED***
			var envelop metricJSONAPI

			assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &envelop))

			metric := envelop.Data.Attributes

			assert.Equal(t, "my_metric", envelop.Data.ID)
			assert.True(t, metric.Type.Valid)
			assert.Equal(t, metrics.Trend, metric.Type.Type)
			assert.True(t, metric.Contains.Valid)
			assert.Equal(t, metrics.Time, metric.Contains.Type)
			assert.True(t, metric.Tainted.Valid)
			assert.True(t, metric.Tainted.Bool)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
