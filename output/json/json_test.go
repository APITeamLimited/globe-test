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

package json

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

func getValidator(t testing.TB, expected []string) func(io.Reader) ***REMOVED***
	return func(rawJSONLines io.Reader) ***REMOVED***
		s := bufio.NewScanner(rawJSONLines)
		i := 0
		for s.Scan() ***REMOVED***
			i++
			if i > len(expected) ***REMOVED***
				t.Errorf("Read unexpected line number %d, expected only %d entries", i, len(expected))
				continue
			***REMOVED***
			assert.JSONEq(t, expected[i-1], string(s.Bytes()))
		***REMOVED***
		assert.NoError(t, s.Err())
		assert.Equal(t, len(expected), i)
	***REMOVED***
***REMOVED***

func generateTestMetricSamples(t testing.TB) ([]metrics.SampleContainer, func(io.Reader)) ***REMOVED***
	registry := metrics.NewRegistry()

	metric1, err := registry.NewMetric("my_metric1", metrics.Gauge)
	require.NoError(t, err)

	_, err = metric1.AddSubmetric("a:1,b:2")
	require.NoError(t, err)

	metric2, err := registry.NewMetric("my_metric2", metrics.Counter, metrics.Data)
	require.NoError(t, err)

	time1 := time.Date(2021, time.February, 24, 13, 37, 10, 0, time.UTC)
	time2 := time1.Add(10 * time.Second)
	time3 := time2.Add(10 * time.Second)

	connTags := metrics.NewSampleTags(map[string]string***REMOVED***"key": "val"***REMOVED***)

	samples := []metrics.SampleContainer***REMOVED***
		metrics.Sample***REMOVED***Time: time1, Metric: metric1, Value: float64(1), Tags: metrics.NewSampleTags(map[string]string***REMOVED***"tag1": "val1"***REMOVED***)***REMOVED***,
		metrics.Sample***REMOVED***Time: time1, Metric: metric1, Value: float64(2), Tags: metrics.NewSampleTags(map[string]string***REMOVED***"tag2": "val2"***REMOVED***)***REMOVED***,
		metrics.ConnectedSamples***REMOVED***Samples: []metrics.Sample***REMOVED***
			***REMOVED***Time: time2, Metric: metric2, Value: float64(3), Tags: connTags***REMOVED***,
			***REMOVED***Time: time2, Metric: metric1, Value: float64(4), Tags: connTags***REMOVED***,
		***REMOVED***, Time: time2, Tags: connTags***REMOVED***,
		metrics.Sample***REMOVED***Time: time3, Metric: metric2, Value: float64(5), Tags: metrics.NewSampleTags(map[string]string***REMOVED***"tag3": "val3"***REMOVED***)***REMOVED***,
	***REMOVED***
	expected := []string***REMOVED***
		`***REMOVED***"type":"Metric","data":***REMOVED***"name":"my_metric1","type":"gauge","contains":"default","tainted":null,"thresholds":["rate<0.01","p(99)<250"],"submetrics":[***REMOVED***"name":"my_metric1***REMOVED***a:1,b:2***REMOVED***","suffix":"a:1,b:2","tags":***REMOVED***"a":"1","b":"2"***REMOVED******REMOVED***]***REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:10Z","value":1,"tags":***REMOVED***"tag1":"val1"***REMOVED******REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:10Z","value":2,"tags":***REMOVED***"tag2":"val2"***REMOVED******REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Metric","data":***REMOVED***"name":"my_metric2","type":"counter","contains":"data","tainted":null,"thresholds":[],"submetrics":null***REMOVED***,"metric":"my_metric2"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:20Z","value":3,"tags":***REMOVED***"key":"val"***REMOVED******REMOVED***,"metric":"my_metric2"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:20Z","value":4,"tags":***REMOVED***"key":"val"***REMOVED******REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:30Z","value":5,"tags":***REMOVED***"tag3":"val3"***REMOVED******REMOVED***,"metric":"my_metric2"***REMOVED***`,
	***REMOVED***

	return samples, getValidator(t, expected)
***REMOVED***

func TestJsonOutputStdout(t *testing.T) ***REMOVED***
	t.Parallel()

	stdout := new(bytes.Buffer)
	out, err := New(output.Params***REMOVED***
		Logger: testutils.NewLogger(t),
		StdOut: stdout,
	***REMOVED***)
	require.NoError(t, err)

	setThresholds(t, out)
	require.NoError(t, out.Start())

	samples, validateResults := generateTestMetricSamples(t)
	out.AddMetricSamples(samples[:2])
	out.AddMetricSamples(samples[2:])
	require.NoError(t, out.Stop())
	validateResults(stdout)
***REMOVED***

func TestJsonOutputFileError(t *testing.T) ***REMOVED***
	t.Parallel()

	stdout := new(bytes.Buffer)
	fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
	out, err := New(output.Params***REMOVED***
		Logger:         testutils.NewLogger(t),
		StdOut:         stdout,
		FS:             fs,
		ConfigArgument: "/json-output",
	***REMOVED***)
	require.NoError(t, err)
	assert.Error(t, out.Start())
***REMOVED***

func TestJsonOutputFile(t *testing.T) ***REMOVED***
	t.Parallel()

	stdout := new(bytes.Buffer)
	fs := afero.NewMemMapFs()
	out, err := New(output.Params***REMOVED***
		Logger:         testutils.NewLogger(t),
		StdOut:         stdout,
		FS:             fs,
		ConfigArgument: "/json-output",
	***REMOVED***)
	require.NoError(t, err)

	setThresholds(t, out)
	require.NoError(t, out.Start())

	samples, validateResults := generateTestMetricSamples(t)
	out.AddMetricSamples(samples[:2])
	out.AddMetricSamples(samples[2:])
	require.NoError(t, out.Stop())

	assert.Empty(t, stdout.Bytes())
	file, err := fs.Open("/json-output")
	require.NoError(t, err)
	validateResults(file)
	assert.NoError(t, file.Close())
***REMOVED***

func TestJsonOutputFileGzipped(t *testing.T) ***REMOVED***
	t.Parallel()

	stdout := new(bytes.Buffer)
	fs := afero.NewMemMapFs()
	out, err := New(output.Params***REMOVED***
		Logger:         testutils.NewLogger(t),
		StdOut:         stdout,
		FS:             fs,
		ConfigArgument: "/json-output.gz",
	***REMOVED***)
	require.NoError(t, err)

	setThresholds(t, out)
	require.NoError(t, out.Start())

	samples, validateResults := generateTestMetricSamples(t)
	out.AddMetricSamples(samples[:2])
	out.AddMetricSamples(samples[2:])
	require.NoError(t, out.Stop())

	assert.Empty(t, stdout.Bytes())
	file, err := fs.Open("/json-output.gz")
	require.NoError(t, err)
	reader, err := gzip.NewReader(file)
	require.NoError(t, err)
	validateResults(reader)
	assert.NoError(t, file.Close())
***REMOVED***

func TestWrapSampleWithSamplePointer(t *testing.T) ***REMOVED***
	t.Parallel()
	out := wrapSample(metrics.Sample***REMOVED***
		Metric: &metrics.Metric***REMOVED******REMOVED***,
	***REMOVED***)
	assert.NotEqual(t, out, (*sampleEnvelope)(nil))
***REMOVED***

func setThresholds(t *testing.T, out output.Output) ***REMOVED***
	t.Helper()

	jout, ok := out.(*Output)
	require.True(t, ok)

	ts := metrics.NewThresholds([]string***REMOVED***"rate<0.01", "p(99)<250"***REMOVED***)
	jout.SetThresholds(map[string]metrics.Thresholds***REMOVED***"my_metric1": ts***REMOVED***)
***REMOVED***
