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

package stats

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		Type     MetricType
		SinkType Sink
	***REMOVED******REMOVED***
		"Counter": ***REMOVED***Counter, &CounterSink***REMOVED******REMOVED******REMOVED***,
		"Gauge":   ***REMOVED***Gauge, &GaugeSink***REMOVED******REMOVED******REMOVED***,
		"Trend":   ***REMOVED***Trend, &TrendSink***REMOVED******REMOVED******REMOVED***,
		"Rate":    ***REMOVED***Rate, &RateSink***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		name, data := name, data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			m := New("my_metric", data.Type)
			assert.Equal(t, "my_metric", m.Name)
			assert.IsType(t, data.SinkType, m.Sink)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestAddSubmetric(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		err  bool
		tags map[string]string
	***REMOVED******REMOVED***
		"":                        ***REMOVED***true, nil***REMOVED***,
		"  ":                      ***REMOVED***true, nil***REMOVED***,
		"a":                       ***REMOVED***false, map[string]string***REMOVED***"a": ""***REMOVED******REMOVED***,
		"a:1":                     ***REMOVED***false, map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		" a : 1 ":                 ***REMOVED***false, map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		"a,b":                     ***REMOVED***false, map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		` a:"",b: ''`:             ***REMOVED***false, map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		`a:1,b:2`:                 ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		` a : 1, b : 2 `:          ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		`a : '1' , b : "2"`:       ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		`" a" : ' 1' , b : "2 " `: ***REMOVED***false, map[string]string***REMOVED***" a": " 1", "b": "2 "***REMOVED******REMOVED***, //nolint:gocritic
	***REMOVED***

	for name, expected := range testdata ***REMOVED***
		name, expected := name, expected
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			m := New("metric", Trend)
			sm, err := m.AddSubmetric(name)
			if expected.err ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			require.NotNil(t, sm)
			assert.EqualValues(t, expected.tags, sm.Tags.tags)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSampleTags(t *testing.T) ***REMOVED***
	t.Parallel()

	// Nil pointer to SampleTags
	var nilTags *SampleTags
	assert.True(t, nilTags.IsEqual(nilTags))
	assert.Equal(t, map[string]string***REMOVED******REMOVED***, nilTags.CloneTags())

	nilJSON, err := json.Marshal(nilTags)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(nilJSON))

	// Empty SampleTags
	emptyTagMap := map[string]string***REMOVED******REMOVED***
	emptyTags := NewSampleTags(emptyTagMap)
	assert.Nil(t, emptyTags)
	assert.True(t, emptyTags.IsEqual(emptyTags))
	assert.True(t, emptyTags.IsEqual(nilTags))
	assert.Equal(t, emptyTagMap, emptyTags.CloneTags())

	emptyJSON, err := json.Marshal(emptyTags)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(emptyJSON))

	var emptyTagsUnmarshaled *SampleTags
	err = json.Unmarshal(emptyJSON, &emptyTagsUnmarshaled)
	assert.NoError(t, err)
	assert.Nil(t, emptyTagsUnmarshaled)
	assert.True(t, emptyTagsUnmarshaled.IsEqual(emptyTags))
	assert.True(t, emptyTagsUnmarshaled.IsEqual(nilTags))
	assert.Equal(t, emptyTagMap, emptyTagsUnmarshaled.CloneTags())

	// SampleTags with keys and values
	tagMap := map[string]string***REMOVED***"key1": "val1", "key2": "val2"***REMOVED***
	tags := NewSampleTags(tagMap)
	assert.NotNil(t, tags)
	assert.True(t, tags.IsEqual(tags))
	assert.False(t, tags.IsEqual(nilTags))
	assert.False(t, tags.IsEqual(emptyTags))
	assert.False(t, tags.IsEqual(IntoSampleTags(&map[string]string***REMOVED***"key1": "val1", "key2": "val3"***REMOVED***)))
	assert.True(t, tags.Contains(IntoSampleTags(&map[string]string***REMOVED***"key1": "val1"***REMOVED***)))
	assert.False(t, tags.Contains(IntoSampleTags(&map[string]string***REMOVED***"key3": "val1"***REMOVED***)))
	assert.Equal(t, tagMap, tags.CloneTags())

	assert.Nil(t, tags.json) // No cache
	tagsJSON, err := json.Marshal(tags)
	expJSON := `***REMOVED***"key1":"val1","key2":"val2"***REMOVED***`
	assert.NoError(t, err)
	assert.JSONEq(t, expJSON, string(tagsJSON))
	assert.JSONEq(t, expJSON, string(tags.json)) // Populated cache

	var tagsUnmarshaled *SampleTags
	err = json.Unmarshal(tagsJSON, &tagsUnmarshaled)
	assert.NoError(t, err)
	assert.NotNil(t, tagsUnmarshaled)
	assert.True(t, tagsUnmarshaled.IsEqual(tags))
	assert.False(t, tagsUnmarshaled.IsEqual(nilTags))
	assert.Equal(t, tagMap, tagsUnmarshaled.CloneTags())
***REMOVED***

func TestSampleImplementations(t *testing.T) ***REMOVED***
	tagMap := map[string]string***REMOVED***"key1": "val1", "key2": "val2"***REMOVED***
	now := time.Now()

	sample := Sample***REMOVED***
		Metric: New("test_metric", Counter),
		Time:   now,
		Tags:   NewSampleTags(tagMap),
		Value:  1.0,
	***REMOVED***
	samples := Samples(sample.GetSamples())
	cSamples := ConnectedSamples***REMOVED***
		Samples: []Sample***REMOVED***sample***REMOVED***,
		Time:    now,
		Tags:    NewSampleTags(tagMap),
	***REMOVED***
	exp := []Sample***REMOVED***sample***REMOVED***
	assert.Equal(t, exp, sample.GetSamples())
	assert.Equal(t, exp, samples.GetSamples())
	assert.Equal(t, exp, cSamples.GetSamples())
	assert.Equal(t, now, sample.GetTime())
	assert.Equal(t, now, cSamples.GetTime())
	assert.Equal(t, sample.GetTags(), sample.GetTags())
***REMOVED***

func TestGetResolversForTrendColumnsValidation(t *testing.T) ***REMOVED***
	validateTests := []struct ***REMOVED***
		stats  []string
		expErr bool
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED******REMOVED***, false***REMOVED***,
		***REMOVED***[]string***REMOVED***"avg", "min", "med", "max", "p(0)", "p(99)", "p(99.999)", "count"***REMOVED***, false***REMOVED***,
		***REMOVED***[]string***REMOVED***"avg", "p(err)"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"nil", "p(err)"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"p90"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"p(90"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***" avg"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"avg "***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"", "avg "***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"p(-1)"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"p(101)"***REMOVED***, true***REMOVED***,
		***REMOVED***[]string***REMOVED***"p(1)"***REMOVED***, false***REMOVED***,
	***REMOVED***

	for _, tc := range validateTests ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%v", tc.stats), func(t *testing.T) ***REMOVED***
			_, err := GetResolversForTrendColumns(tc.stats)
			if tc.expErr ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func createTestTrendSink(count int) *TrendSink ***REMOVED***
	sink := TrendSink***REMOVED******REMOVED***

	for i := 0; i < count; i++ ***REMOVED***
		sink.Add(Sample***REMOVED***Value: float64(i)***REMOVED***)
	***REMOVED***

	return &sink
***REMOVED***

func TestResolversForTrendColumnsCalculation(t *testing.T) ***REMOVED***
	customResolversTests := []struct ***REMOVED***
		stats      string
		percentile float64
	***REMOVED******REMOVED***
		***REMOVED***"p(50)", 0.5***REMOVED***,
		***REMOVED***"p(99)", 0.99***REMOVED***,
		***REMOVED***"p(99.9)", 0.999***REMOVED***,
		***REMOVED***"p(99.99)", 0.9999***REMOVED***,
		***REMOVED***"p(99.999)", 0.99999***REMOVED***,
	***REMOVED***

	sink := createTestTrendSink(100)

	for _, tc := range customResolversTests ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%v", tc.stats), func(t *testing.T) ***REMOVED***
			res, err := GetResolversForTrendColumns([]string***REMOVED***tc.stats***REMOVED***)
			assert.NoError(t, err)
			assert.Len(t, res, 1)
			for k := range res ***REMOVED***
				assert.InDelta(t, sink.P(tc.percentile), res[k](sink), 0.000001)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
