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

package csv

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/stats"
)

func TestMakeHeader(t *testing.T) ***REMOVED***
	testdata := map[string][]string***REMOVED***
		"One tag": ***REMOVED***
			"tag1",
		***REMOVED***,
		"Two tags": ***REMOVED***
			"tag1", "tag2",
		***REMOVED***,
	***REMOVED***

	for testname, tags := range testdata ***REMOVED***
		testname, tags := testname, tags
		t.Run(testname, func(t *testing.T) ***REMOVED***
			header := MakeHeader(tags)
			assert.Equal(t, len(tags)+4, len(header))
			assert.Equal(t, "metric_name", header[0])
			assert.Equal(t, "timestamp", header[1])
			assert.Equal(t, "metric_value", header[2])
			assert.Equal(t, "extra_tags", header[len(header)-1])
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSampleToRow(t *testing.T) ***REMOVED***
	testData := []struct ***REMOVED***
		testname    string
		sample      *stats.Sample
		resTags     []string
		ignoredTags []string
	***REMOVED******REMOVED***
		***REMOVED***
			testname: "One res tag, one ignored tag, one extra tag",
			sample: &stats.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: stats.New("my_metric", stats.Gauge),
				Value:  1,
				Tags: stats.NewSampleTags(map[string]string***REMOVED***
					"tag1": "val1",
					"tag2": "val2",
					"tag3": "val3",
				***REMOVED***),
			***REMOVED***,
			resTags:     []string***REMOVED***"tag1"***REMOVED***,
			ignoredTags: []string***REMOVED***"tag2"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			testname: "Two res tags, three extra tags",
			sample: &stats.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: stats.New("my_metric", stats.Gauge),
				Value:  1,
				Tags: stats.NewSampleTags(map[string]string***REMOVED***
					"tag1": "val1",
					"tag2": "val2",
					"tag3": "val3",
					"tag4": "val4",
					"tag5": "val5",
				***REMOVED***),
			***REMOVED***,
			resTags:     []string***REMOVED***"tag1", "tag2"***REMOVED***,
			ignoredTags: []string***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			testname: "Two res tags, two ignored",
			sample: &stats.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: stats.New("my_metric", stats.Gauge),
				Value:  1,
				Tags: stats.NewSampleTags(map[string]string***REMOVED***
					"tag1": "val1",
					"tag2": "val2",
					"tag3": "val3",
					"tag4": "val4",
					"tag5": "val5",
					"tag6": "val6",
				***REMOVED***),
			***REMOVED***,
			resTags:     []string***REMOVED***"tag1", "tag3"***REMOVED***,
			ignoredTags: []string***REMOVED***"tag4", "tag6"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	expected := []struct ***REMOVED***
		baseRow  []string
		extraRow []string
	***REMOVED******REMOVED***
		***REMOVED***
			baseRow: []string***REMOVED***
				"my_metric",
				"1562324644",
				"1.000000",
				"val1",
			***REMOVED***,
			extraRow: []string***REMOVED***
				"tag3=val3",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			baseRow: []string***REMOVED***
				"my_metric",
				"1562324644",
				"1.000000",
				"val1",
				"val2",
			***REMOVED***,
			extraRow: []string***REMOVED***
				"tag3=val3",
				"tag4=val4",
				"tag5=val5",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			baseRow: []string***REMOVED***
				"my_metric",
				"1562324644",
				"1.000000",
				"val1",
				"val3",
			***REMOVED***,
			extraRow: []string***REMOVED***
				"tag2=val2",
				"tag5=val5",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for i := range testData ***REMOVED***
		testname, sample := testData[i].testname, testData[i].sample
		resTags, ignoredTags := testData[i].resTags, testData[i].ignoredTags
		expectedRow := expected[i]

		t.Run(testname, func(t *testing.T) ***REMOVED***
			row := SampleToRow(sample, resTags, ignoredTags, make([]string, 3+len(resTags)+1))
			for ind, cell := range expectedRow.baseRow ***REMOVED***
				assert.Equal(t, cell, row[ind])
			***REMOVED***
			for _, cell := range expectedRow.extraRow ***REMOVED***
				assert.Contains(t, row[len(row)-1], cell)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func readUnCompressedFile(fileName string, fs afero.Fs) string ***REMOVED***
	csvbytes, err := afero.ReadFile(fs, fileName)
	if err != nil ***REMOVED***
		return err.Error()
	***REMOVED***

	return fmt.Sprintf("%s", csvbytes)
***REMOVED***

func readCompressedFile(fileName string, fs afero.Fs) string ***REMOVED***
	file, err := fs.Open(fileName)
	if err != nil ***REMOVED***
		return err.Error()
	***REMOVED***

	gzf, err := gzip.NewReader(file)
	if err != nil ***REMOVED***
		return err.Error()
	***REMOVED***

	csvbytes, err := ioutil.ReadAll(gzf)
	if err != nil ***REMOVED***
		return err.Error()
	***REMOVED***

	return fmt.Sprintf("%s", csvbytes)
***REMOVED***

func TestRun(t *testing.T) ***REMOVED***
	t.Parallel()
	testData := []struct ***REMOVED***
		samples        []stats.SampleContainer
		fileName       string
		fileReaderFunc func(fileName string, fs afero.Fs) string
		outputContent  string
	***REMOVED******REMOVED***
		***REMOVED***
			samples: []stats.SampleContainer***REMOVED***
				stats.Sample***REMOVED***
					Time:   time.Unix(1562324643, 0),
					Metric: stats.New("my_metric", stats.Gauge),
					Value:  1,
					Tags: stats.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
					***REMOVED***),
				***REMOVED***,
				stats.Sample***REMOVED***
					Time:   time.Unix(1562324644, 0),
					Metric: stats.New("my_metric", stats.Gauge),
					Value:  1,
					Tags: stats.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
						"tag4":  "val4",
					***REMOVED***),
				***REMOVED***,
			***REMOVED***,
			fileName:       "test",
			fileReaderFunc: readUnCompressedFile,
			outputContent:  "metric_name,timestamp,metric_value,check,error,extra_tags\n" + "my_metric,1562324643,1.000000,val1,val3,url=val2\n" + "my_metric,1562324644,1.000000,val1,val3,tag4=val4&url=val2\n",
		***REMOVED***,
		***REMOVED***
			samples: []stats.SampleContainer***REMOVED***
				stats.Sample***REMOVED***
					Time:   time.Unix(1562324643, 0),
					Metric: stats.New("my_metric", stats.Gauge),
					Value:  1,
					Tags: stats.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
					***REMOVED***),
				***REMOVED***,
				stats.Sample***REMOVED***
					Time:   time.Unix(1562324644, 0),
					Metric: stats.New("my_metric", stats.Gauge),
					Value:  1,
					Tags: stats.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
						"name":  "val4",
					***REMOVED***),
				***REMOVED***,
			***REMOVED***,
			fileName:       "test.gz",
			fileReaderFunc: readCompressedFile,
			outputContent:  "metric_name,timestamp,metric_value,check,error,extra_tags\n" + "my_metric,1562324643,1.000000,val1,val3,url=val2\n" + "my_metric,1562324644,1.000000,val1,val3,name=val4&url=val2\n",
		***REMOVED***,
	***REMOVED***

	for _, data := range testData ***REMOVED***
		mem := afero.NewMemMapFs()
		output, err := newOutput(output.Params***REMOVED***
			Logger:         testutils.NewLogger(t),
			FS:             mem,
			ConfigArgument: data.fileName,
			ScriptOptions: lib.Options***REMOVED***
				SystemTags: stats.NewSystemTagSet(stats.TagError | stats.TagCheck),
			***REMOVED***,
		***REMOVED***)
		require.NoError(t, err)
		require.NotNil(t, output)

		require.NoError(t, output.Start())
		output.AddMetricSamples(data.samples)
		time.Sleep(1 * time.Second)
		require.NoError(t, output.Stop())

		finalOutput := data.fileReaderFunc(data.fileName, mem)
		assert.Equal(t, data.outputContent, sortExtraTagsForTest(t, finalOutput))
	***REMOVED***
***REMOVED***

func sortExtraTagsForTest(t *testing.T, input string) string ***REMOVED***
	t.Helper()
	r := csv.NewReader(strings.NewReader(input))
	lines, err := r.ReadAll()
	require.NoError(t, err)
	for i, line := range lines[1:] ***REMOVED***
		extraTags := strings.Split(line[len(line)-1], "&")
		sort.Strings(extraTags)
		lines[i+1][len(line)-1] = strings.Join(extraTags, "&")
	***REMOVED***
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	require.NoError(t, w.WriteAll(lines))
	w.Flush()
	return b.String()
***REMOVED***
