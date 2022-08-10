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

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
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
	testMetric, err := metrics.NewRegistry().NewMetric("my_metric", metrics.Gauge)
	require.NoError(t, err)

	testData := []struct ***REMOVED***
		testname    string
		sample      *metrics.Sample
		resTags     []string
		ignoredTags []string
		timeFormat  string
	***REMOVED******REMOVED***
		***REMOVED***
			testname: "One res tag, one ignored tag, one extra tag",
			sample: &metrics.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: testMetric,
				Value:  1,
				Tags: metrics.NewSampleTags(map[string]string***REMOVED***
					"tag1": "val1",
					"tag2": "val2",
					"tag3": "val3",
				***REMOVED***),
			***REMOVED***,
			resTags:     []string***REMOVED***"tag1"***REMOVED***,
			ignoredTags: []string***REMOVED***"tag2"***REMOVED***,
			timeFormat:  "unix",
		***REMOVED***,
		***REMOVED***
			testname: "Two res tags, three extra tags",
			sample: &metrics.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: testMetric,
				Value:  1,
				Tags: metrics.NewSampleTags(map[string]string***REMOVED***
					"tag1": "val1",
					"tag2": "val2",
					"tag3": "val3",
					"tag4": "val4",
					"tag5": "val5",
				***REMOVED***),
			***REMOVED***,
			resTags:     []string***REMOVED***"tag1", "tag2"***REMOVED***,
			ignoredTags: []string***REMOVED******REMOVED***,
			timeFormat:  "unix",
		***REMOVED***,
		***REMOVED***
			testname: "Two res tags, two ignored, with RFC3339 timestamp",
			sample: &metrics.Sample***REMOVED***
				Time:   time.Unix(1562324644, 0),
				Metric: testMetric,
				Value:  1,
				Tags: metrics.NewSampleTags(map[string]string***REMOVED***
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
			timeFormat:  "rfc3339",
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
				time.Unix(1562324644, 0).Format(time.RFC3339),
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
		timeFormat, err := TimeFormatString(testData[i].timeFormat)
		require.NoError(t, err)
		expectedRow := expected[i]

		t.Run(testname, func(t *testing.T) ***REMOVED***
			row := SampleToRow(sample, resTags, ignoredTags, make([]string, 3+len(resTags)+1), timeFormat)
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

	testMetric, err := metrics.NewRegistry().NewMetric("my_metric", metrics.Gauge)
	require.NoError(t, err)

	testData := []struct ***REMOVED***
		samples        []metrics.SampleContainer
		fileName       string
		fileReaderFunc func(fileName string, fs afero.Fs) string
		timeFormat     string
		outputContent  string
	***REMOVED******REMOVED***
		***REMOVED***
			samples: []metrics.SampleContainer***REMOVED***
				metrics.Sample***REMOVED***
					Time:   time.Unix(1562324643, 0),
					Metric: testMetric,
					Value:  1,
					Tags: metrics.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
					***REMOVED***),
				***REMOVED***,
				metrics.Sample***REMOVED***
					Time:   time.Unix(1562324644, 0),
					Metric: testMetric,
					Value:  1,
					Tags: metrics.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
						"tag4":  "val4",
					***REMOVED***),
				***REMOVED***,
			***REMOVED***,
			fileName:       "test",
			fileReaderFunc: readUnCompressedFile,
			timeFormat:     "",
			outputContent:  "metric_name,timestamp,metric_value,check,error,extra_tags\n" + "my_metric,1562324643,1.000000,val1,val3,url=val2\n" + "my_metric,1562324644,1.000000,val1,val3,tag4=val4&url=val2\n",
		***REMOVED***,
		***REMOVED***
			samples: []metrics.SampleContainer***REMOVED***
				metrics.Sample***REMOVED***
					Time:   time.Unix(1562324643, 0),
					Metric: testMetric,
					Value:  1,
					Tags: metrics.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
					***REMOVED***),
				***REMOVED***,
				metrics.Sample***REMOVED***
					Time:   time.Unix(1562324644, 0),
					Metric: testMetric,
					Value:  1,
					Tags: metrics.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
						"name":  "val4",
					***REMOVED***),
				***REMOVED***,
			***REMOVED***,
			fileName:       "test.gz",
			fileReaderFunc: readCompressedFile,
			timeFormat:     "unix",
			outputContent:  "metric_name,timestamp,metric_value,check,error,extra_tags\n" + "my_metric,1562324643,1.000000,val1,val3,url=val2\n" + "my_metric,1562324644,1.000000,val1,val3,name=val4&url=val2\n",
		***REMOVED***,
		***REMOVED***
			samples: []metrics.SampleContainer***REMOVED***
				metrics.Sample***REMOVED***
					Time:   time.Unix(1562324644, 0),
					Metric: testMetric,
					Value:  1,
					Tags: metrics.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
					***REMOVED***),
				***REMOVED***,
				metrics.Sample***REMOVED***
					Time:   time.Unix(1562324644, 0),
					Metric: testMetric,
					Value:  1,
					Tags: metrics.NewSampleTags(map[string]string***REMOVED***
						"check": "val1",
						"url":   "val2",
						"error": "val3",
						"name":  "val4",
					***REMOVED***),
				***REMOVED***,
			***REMOVED***,
			fileName:       "test",
			fileReaderFunc: readUnCompressedFile,
			timeFormat:     "rfc3339",
			outputContent: "metric_name,timestamp,metric_value,check,error,extra_tags\n" +
				"my_metric," + time.Unix(1562324644, 0).Format(time.RFC3339) + ",1.000000,val1,val3,url=val2\n" +
				"my_metric," + time.Unix(1562324644, 0).Format(time.RFC3339) + ",1.000000,val1,val3,name=val4&url=val2\n",
		***REMOVED***,
	***REMOVED***

	for i, data := range testData ***REMOVED***
		name := fmt.Sprint(i)
		data := data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			mem := afero.NewMemMapFs()
			env := make(map[string]string)
			if data.timeFormat != "" ***REMOVED***
				env["K6_CSV_TIME_FORMAT"] = data.timeFormat
			***REMOVED***

			output, err := newOutput(output.Params***REMOVED***
				Logger:         testutils.NewLogger(t),
				FS:             mem,
				Environment:    env,
				ConfigArgument: data.fileName,
				ScriptOptions: lib.Options***REMOVED***
					SystemTags: metrics.NewSystemTagSet(metrics.TagError | metrics.TagCheck),
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
		***REMOVED***)
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
