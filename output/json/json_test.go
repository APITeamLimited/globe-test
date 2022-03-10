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
	"go.k6.io/k6/output"
	"go.k6.io/k6/stats"
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
			assert.Equal(t, expected[i-1], string(s.Bytes()))
		***REMOVED***
		assert.NoError(t, s.Err())
		assert.Equal(t, len(expected), i)
	***REMOVED***
***REMOVED***

func generateTestMetricSamples(t testing.TB) ([]stats.SampleContainer, func(io.Reader)) ***REMOVED***
	metric1 := stats.New("my_metric1", stats.Gauge)
	metric2 := stats.New("my_metric2", stats.Counter, stats.Data)
	time1 := time.Date(2021, time.February, 24, 13, 37, 10, 0, time.UTC)
	time2 := time1.Add(10 * time.Second)
	time3 := time2.Add(10 * time.Second)
	connTags := stats.NewSampleTags(map[string]string***REMOVED***"key": "val"***REMOVED***)

	samples := []stats.SampleContainer***REMOVED***
		stats.Sample***REMOVED***Time: time1, Metric: metric1, Value: float64(1), Tags: stats.NewSampleTags(map[string]string***REMOVED***"tag1": "val1"***REMOVED***)***REMOVED***,
		stats.Sample***REMOVED***Time: time1, Metric: metric1, Value: float64(2), Tags: stats.NewSampleTags(map[string]string***REMOVED***"tag2": "val2"***REMOVED***)***REMOVED***,
		stats.ConnectedSamples***REMOVED***Samples: []stats.Sample***REMOVED***
			***REMOVED***Time: time2, Metric: metric2, Value: float64(3), Tags: connTags***REMOVED***,
			***REMOVED***Time: time2, Metric: metric1, Value: float64(4), Tags: connTags***REMOVED***,
		***REMOVED***, Time: time2, Tags: connTags***REMOVED***,
		stats.Sample***REMOVED***Time: time3, Metric: metric2, Value: float64(5), Tags: stats.NewSampleTags(map[string]string***REMOVED***"tag3": "val3"***REMOVED***)***REMOVED***,
	***REMOVED***
	expected := []string***REMOVED***
		`***REMOVED***"type":"Metric","data":***REMOVED***"name":"my_metric1","type":"gauge","contains":"default","tainted":null,"thresholds":["rate<0.01","p(99)<250"],"submetrics":null,"sub":***REMOVED***"name":"","parent":"","suffix":"","tags":null***REMOVED******REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:10Z","value":1,"tags":***REMOVED***"tag1":"val1"***REMOVED******REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Point","data":***REMOVED***"time":"2021-02-24T13:37:10Z","value":2,"tags":***REMOVED***"tag2":"val2"***REMOVED******REMOVED***,"metric":"my_metric1"***REMOVED***`,
		`***REMOVED***"type":"Metric","data":***REMOVED***"name":"my_metric2","type":"counter","contains":"data","tainted":null,"thresholds":[],"submetrics":null,"sub":***REMOVED***"name":"","parent":"","suffix":"","tags":null***REMOVED******REMOVED***,"metric":"my_metric2"***REMOVED***`,
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
	out := WrapSample(stats.Sample***REMOVED***
		Metric: &stats.Metric***REMOVED******REMOVED***,
	***REMOVED***)
	assert.NotEqual(t, out, (*Envelope)(nil))
***REMOVED***

func TestWrapMetricWithMetricPointer(t *testing.T) ***REMOVED***
	t.Parallel()
	out := wrapMetric(&stats.Metric***REMOVED******REMOVED***)
	assert.NotEqual(t, out, (*Envelope)(nil))
***REMOVED***

func setThresholds(t *testing.T, out output.Output) ***REMOVED***
	t.Helper()

	jout, ok := out.(*Output)
	require.True(t, ok)

	ts := stats.NewThresholds([]string***REMOVED***"rate<0.01", "p(99)<250"***REMOVED***)

	jout.SetThresholds(map[string]stats.Thresholds***REMOVED***"my_metric1": ts***REMOVED***)
***REMOVED***
