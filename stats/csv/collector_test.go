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
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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
func TestCollect(t *testing.T) ***REMOVED***
	testSamples := []stats.SampleContainer***REMOVED***
		stats.Sample***REMOVED***
			Time:   time.Unix(1562324643, 0),
			Metric: stats.New("my_metric", stats.Gauge),
			Value:  1,
			Tags: stats.NewSampleTags(map[string]string***REMOVED***
				"tag1": "val1",
				"tag2": "val2",
				"tag3": "val3",
			***REMOVED***),
		***REMOVED***,
		stats.Sample***REMOVED***
			Time:   time.Unix(1562324644, 0),
			Metric: stats.New("my_metric", stats.Gauge),
			Value:  1,
			Tags: stats.NewSampleTags(map[string]string***REMOVED***
				"tag1": "val1",
				"tag2": "val2",
				"tag3": "val3",
				"tag4": "val4",
			***REMOVED***),
		***REMOVED***,
	***REMOVED***

	mem := afero.NewMemMapFs()
	collector, err := New(
		mem,
		stats.SystemTagMap***REMOVED***"tag1": true, "tag2": false, "tag3": true***REMOVED***,
		Config***REMOVED***FileName: null.StringFrom("name"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
	)
	assert.NoError(t, err)
	assert.NotNil(t, collector)

	collector.Collect(testSamples)

	assert.Equal(t, len(testSamples), len(collector.buffer))
***REMOVED***

func TestRun(t *testing.T) ***REMOVED***
	collector, err := New(
		afero.NewMemMapFs(),
		stats.SystemTagMap***REMOVED***"tag1": true, "tag2": false, "tag3": true***REMOVED***,
		Config***REMOVED***FileName: null.StringFrom("name"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
	)
	assert.NoError(t, err)
	assert.NotNil(t, collector)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		defer wg.Done()
		err := collector.Init()
		assert.NoError(t, err)
		collector.Run(ctx)
	***REMOVED***()
	cancel()
	wg.Wait()
***REMOVED***

func TestRunCollect(t *testing.T) ***REMOVED***
	testSamples := []stats.SampleContainer***REMOVED***
		stats.Sample***REMOVED***
			Time:   time.Unix(1562324643, 0),
			Metric: stats.New("my_metric", stats.Gauge),
			Value:  1,
			Tags: stats.NewSampleTags(map[string]string***REMOVED***
				"tag1": "val1",
				"tag2": "val2",
				"tag3": "val3",
			***REMOVED***),
		***REMOVED***,
		stats.Sample***REMOVED***
			Time:   time.Unix(1562324644, 0),
			Metric: stats.New("my_metric", stats.Gauge),
			Value:  1,
			Tags: stats.NewSampleTags(map[string]string***REMOVED***
				"tag1": "val1",
				"tag2": "val2",
				"tag3": "val3",
				"tag4": "val4",
			***REMOVED***),
		***REMOVED***,
	***REMOVED***

	mem := afero.NewMemMapFs()
	collector, err := New(
		mem,
		stats.SystemTagMap***REMOVED***"tag1": true, "tag2": false, "tag3": true***REMOVED***,
		Config***REMOVED***FileName: null.StringFrom("path"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
	)
	assert.NoError(t, err)
	assert.NotNil(t, collector)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		collector.Run(ctx)
		wg.Done()
	***REMOVED***()
	err = collector.Init()
	assert.NoError(t, err)
	collector.Collect(testSamples)
	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()
	csvbytes, _ := afero.ReadFile(mem, "path")
	csvstr := fmt.Sprintf("%s", csvbytes)
	assert.Equal(t,
		"metric_name,timestamp,metric_value,tag1,tag3,extra_tags\n"+
			"my_metric,1562324643,1.000000,val1,val3,\n"+
			"my_metric,1562324644,1.000000,val1,val3,tag4=val4\n",
		csvstr)
***REMOVED***

func TestNew(t *testing.T) ***REMOVED***
	configs := []struct ***REMOVED***
		cfg  Config
		tags stats.SystemTagMap
	***REMOVED******REMOVED***
		***REMOVED***
			cfg: Config***REMOVED***FileName: null.StringFrom("name"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
			tags: stats.SystemTagMap***REMOVED***
				"tag1": true,
				"tag2": false,
				"tag3": true,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			cfg: Config***REMOVED***FileName: null.StringFrom("-"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
			tags: stats.SystemTagMap***REMOVED***
				"tag1": true,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			cfg: Config***REMOVED***FileName: null.StringFrom(""), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
			tags: stats.SystemTagMap***REMOVED***
				"tag1": false,
				"tag2": false,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	expected := []struct ***REMOVED***
		fname       string
		resTags     []string
		ignoredTags []string
	***REMOVED******REMOVED***
		***REMOVED***
			fname: "name",
			resTags: []string***REMOVED***
				"tag1", "tag3",
			***REMOVED***,
			ignoredTags: []string***REMOVED***
				"tag2",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			fname: "-",
			resTags: []string***REMOVED***
				"tag1",
			***REMOVED***,
			ignoredTags: []string***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			fname:   "-",
			resTags: []string***REMOVED******REMOVED***,
			ignoredTags: []string***REMOVED***
				"tag1", "tag2",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for i := range configs ***REMOVED***
		config, expected := configs[i], expected[i]
		t.Run(config.cfg.FileName.String, func(t *testing.T) ***REMOVED***
			collector, err := New(afero.NewMemMapFs(), config.tags, config.cfg)
			assert.NoError(t, err)
			assert.NotNil(t, collector)
			assert.Equal(t, expected.fname, collector.fname)
			sort.Strings(expected.resTags)
			sort.Strings(collector.resTags)
			assert.Equal(t, expected.resTags, collector.resTags)
			sort.Strings(expected.ignoredTags)
			sort.Strings(collector.ignoredTags)
			assert.Equal(t, expected.ignoredTags, collector.ignoredTags)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestGetRequiredSystemTags(t *testing.T) ***REMOVED***
	collector, err := New(
		afero.NewMemMapFs(),
		stats.SystemTagMap***REMOVED***"tag1": true, "tag2": false, "tag3": true***REMOVED***,
		Config***REMOVED***FileName: null.StringFrom("name"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
	)
	assert.NoError(t, err)
	assert.NotNil(t, collector)
	assert.Equal(t, stats.SystemTagSet(0), collector.GetRequiredSystemTags())
***REMOVED***

func TestLink(t *testing.T) ***REMOVED***
	collector, err := New(
		afero.NewMemMapFs(),
		stats.SystemTagMap***REMOVED***"tag1": true, "tag2": false, "tag3": true***REMOVED***,
		Config***REMOVED***FileName: null.StringFrom("path"), SaveInterval: types.NewNullDuration(time.Duration(1), true)***REMOVED***,
	)
	assert.NoError(t, err)
	assert.NotNil(t, collector)
	assert.Equal(t, "path", collector.Link())
***REMOVED***
