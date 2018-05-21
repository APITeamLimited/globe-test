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
)

func TestMetricHumanizeValue(t *testing.T) ***REMOVED***
	t.Parallel()
	data := map[*Metric]map[float64][]string***REMOVED***
		***REMOVED***Type: Counter, Contains: Default***REMOVED***: ***REMOVED***
			1.0:     ***REMOVED***"1", "1", "1", "1"***REMOVED***,
			1.5:     ***REMOVED***"1.5", "1.5", "1.5", "1.5"***REMOVED***,
			1.54321: ***REMOVED***"1.54321", "1.54321", "1.54321", "1.54321"***REMOVED***,
		***REMOVED***,
		***REMOVED***Type: Gauge, Contains: Default***REMOVED***: ***REMOVED***
			1.0:     ***REMOVED***"1", "1", "1", "1"***REMOVED***,
			1.5:     ***REMOVED***"1.5", "1.5", "1.5", "1.5"***REMOVED***,
			1.54321: ***REMOVED***"1.54321", "1.54321", "1.54321", "1.54321"***REMOVED***,
		***REMOVED***,
		***REMOVED***Type: Trend, Contains: Default***REMOVED***: ***REMOVED***
			1.0:     ***REMOVED***"1", "1", "1", "1"***REMOVED***,
			1.5:     ***REMOVED***"1.5", "1.5", "1.5", "1.5"***REMOVED***,
			1.54321: ***REMOVED***"1.54321", "1.54321", "1.54321", "1.54321"***REMOVED***,
		***REMOVED***,
		***REMOVED***Type: Counter, Contains: Time***REMOVED***: ***REMOVED***
			D(1):               ***REMOVED***"1ns", "0.00s", "0.00ms", "0.00µs"***REMOVED***,
			D(12):              ***REMOVED***"12ns", "0.00s", "0.00ms", "0.01µs"***REMOVED***,
			D(123):             ***REMOVED***"123ns", "0.00s", "0.00ms", "0.12µs"***REMOVED***,
			D(1234):            ***REMOVED***"1.23µs", "0.00s", "0.00ms", "1.23µs"***REMOVED***,
			D(12345):           ***REMOVED***"12.34µs", "0.00s", "0.01ms", "12.35µs"***REMOVED***,
			D(123456):          ***REMOVED***"123.45µs", "0.00s", "0.12ms", "123.46µs"***REMOVED***,
			D(1234567):         ***REMOVED***"1.23ms", "0.00s", "1.23ms", "1234.57µs"***REMOVED***,
			D(12345678):        ***REMOVED***"12.34ms", "0.01s", "12.35ms", "12345.68µs"***REMOVED***,
			D(123456789):       ***REMOVED***"123.45ms", "0.12s", "123.46ms", "123456.79µs"***REMOVED***,
			D(1234567890):      ***REMOVED***"1.23s", "1.23s", "1234.57ms", "1234567.89µs"***REMOVED***,
			D(12345678901):     ***REMOVED***"12.34s", "12.35s", "12345.68ms", "12345678.90µs"***REMOVED***,
			D(123456789012):    ***REMOVED***"2m3s", "123.46s", "123456.79ms", "123456789.01µs"***REMOVED***,
			D(1234567890123):   ***REMOVED***"20m34s", "1234.57s", "1234567.89ms", "1234567890.12µs"***REMOVED***,
			D(12345678901234):  ***REMOVED***"3h25m45s", "12345.68s", "12345678.90ms", "12345678901.23µs"***REMOVED***,
			D(123456789012345): ***REMOVED***"34h17m36s", "123456.79s", "123456789.01ms", "123456789012.35µs"***REMOVED***,
		***REMOVED***,
		***REMOVED***Type: Gauge, Contains: Time***REMOVED***: ***REMOVED***
			D(1):               ***REMOVED***"1ns", "0.00s", "0.00ms", "0.00µs"***REMOVED***,
			D(12):              ***REMOVED***"12ns", "0.00s", "0.00ms", "0.01µs"***REMOVED***,
			D(123):             ***REMOVED***"123ns", "0.00s", "0.00ms", "0.12µs"***REMOVED***,
			D(1234):            ***REMOVED***"1.23µs", "0.00s", "0.00ms", "1.23µs"***REMOVED***,
			D(12345):           ***REMOVED***"12.34µs", "0.00s", "0.01ms", "12.35µs"***REMOVED***,
			D(123456):          ***REMOVED***"123.45µs", "0.00s", "0.12ms", "123.46µs"***REMOVED***,
			D(1234567):         ***REMOVED***"1.23ms", "0.00s", "1.23ms", "1234.57µs"***REMOVED***,
			D(12345678):        ***REMOVED***"12.34ms", "0.01s", "12.35ms", "12345.68µs"***REMOVED***,
			D(123456789):       ***REMOVED***"123.45ms", "0.12s", "123.46ms", "123456.79µs"***REMOVED***,
			D(1234567890):      ***REMOVED***"1.23s", "1.23s", "1234.57ms", "1234567.89µs"***REMOVED***,
			D(12345678901):     ***REMOVED***"12.34s", "12.35s", "12345.68ms", "12345678.90µs"***REMOVED***,
			D(123456789012):    ***REMOVED***"2m3s", "123.46s", "123456.79ms", "123456789.01µs"***REMOVED***,
			D(1234567890123):   ***REMOVED***"20m34s", "1234.57s", "1234567.89ms", "1234567890.12µs"***REMOVED***,
			D(12345678901234):  ***REMOVED***"3h25m45s", "12345.68s", "12345678.90ms", "12345678901.23µs"***REMOVED***,
			D(123456789012345): ***REMOVED***"34h17m36s", "123456.79s", "123456789.01ms", "123456789012.35µs"***REMOVED***,
		***REMOVED***,
		***REMOVED***Type: Trend, Contains: Time***REMOVED***: ***REMOVED***
			D(1):               ***REMOVED***"1ns", "0.00s", "0.00ms", "0.00µs"***REMOVED***,
			D(12):              ***REMOVED***"12ns", "0.00s", "0.00ms", "0.01µs"***REMOVED***,
			D(123):             ***REMOVED***"123ns", "0.00s", "0.00ms", "0.12µs"***REMOVED***,
			D(1234):            ***REMOVED***"1.23µs", "0.00s", "0.00ms", "1.23µs"***REMOVED***,
			D(12345):           ***REMOVED***"12.34µs", "0.00s", "0.01ms", "12.35µs"***REMOVED***,
			D(123456):          ***REMOVED***"123.45µs", "0.00s", "0.12ms", "123.46µs"***REMOVED***,
			D(1234567):         ***REMOVED***"1.23ms", "0.00s", "1.23ms", "1234.57µs"***REMOVED***,
			D(12345678):        ***REMOVED***"12.34ms", "0.01s", "12.35ms", "12345.68µs"***REMOVED***,
			D(123456789):       ***REMOVED***"123.45ms", "0.12s", "123.46ms", "123456.79µs"***REMOVED***,
			D(1234567890):      ***REMOVED***"1.23s", "1.23s", "1234.57ms", "1234567.89µs"***REMOVED***,
			D(12345678901):     ***REMOVED***"12.34s", "12.35s", "12345.68ms", "12345678.90µs"***REMOVED***,
			D(123456789012):    ***REMOVED***"2m3s", "123.46s", "123456.79ms", "123456789.01µs"***REMOVED***,
			D(1234567890123):   ***REMOVED***"20m34s", "1234.57s", "1234567.89ms", "1234567890.12µs"***REMOVED***,
			D(12345678901234):  ***REMOVED***"3h25m45s", "12345.68s", "12345678.90ms", "12345678901.23µs"***REMOVED***,
			D(123456789012345): ***REMOVED***"34h17m36s", "123456.79s", "123456789.01ms", "123456789012.35µs"***REMOVED***,
		***REMOVED***,
		***REMOVED***Type: Rate, Contains: Default***REMOVED***: ***REMOVED***
			0.0:       ***REMOVED***"0.00%", "0.00%", "0.00%", "0.00%"***REMOVED***,
			0.01:      ***REMOVED***"1.00%", "1.00%", "1.00%", "1.00%"***REMOVED***,
			0.02:      ***REMOVED***"2.00%", "2.00%", "2.00%", "2.00%"***REMOVED***,
			0.022:     ***REMOVED***"2.19%", "2.19%", "2.19%", "2.19%"***REMOVED***, // caused by float truncation
			0.0222:    ***REMOVED***"2.22%", "2.22%", "2.22%", "2.22%"***REMOVED***,
			0.02222:   ***REMOVED***"2.22%", "2.22%", "2.22%", "2.22%"***REMOVED***,
			0.022222:  ***REMOVED***"2.22%", "2.22%", "2.22%", "2.22%"***REMOVED***,
			1.0 / 3.0: ***REMOVED***"33.33%", "33.33%", "33.33%", "33.33%"***REMOVED***,
			0.5:       ***REMOVED***"50.00%", "50.00%", "50.00%", "50.00%"***REMOVED***,
			0.55:      ***REMOVED***"55.00%", "55.00%", "55.00%", "55.00%"***REMOVED***,
			0.555:     ***REMOVED***"55.50%", "55.50%", "55.50%", "55.50%"***REMOVED***,
			0.5555:    ***REMOVED***"55.55%", "55.55%", "55.55%", "55.55%"***REMOVED***,
			0.55555:   ***REMOVED***"55.55%", "55.55%", "55.55%", "55.55%"***REMOVED***,
			0.75:      ***REMOVED***"75.00%", "75.00%", "75.00%", "75.00%"***REMOVED***,
			0.999995:  ***REMOVED***"99.99%", "99.99%", "99.99%", "99.99%"***REMOVED***,
			1.0:       ***REMOVED***"100.00%", "100.00%", "100.00%", "100.00%"***REMOVED***,
			1.5:       ***REMOVED***"150.00%", "150.00%", "150.00%", "150.00%"***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for m, values := range data ***REMOVED***
		m, values := m, values
		t.Run(fmt.Sprintf("type=%s,contains=%s", m.Type.String(), m.Contains.String()), func(t *testing.T) ***REMOVED***
			t.Parallel()
			for v, s := range values ***REMOVED***
				t.Run(fmt.Sprintf("v=%f", v), func(t *testing.T) ***REMOVED***
					t.Run("no fixed unit", func(t *testing.T) ***REMOVED***
						assert.Equal(t, s[0], m.HumanizeValue(v, ""))
					***REMOVED***)
					t.Run("unit fixed in s", func(t *testing.T) ***REMOVED***
						assert.Equal(t, s[1], m.HumanizeValue(v, "s"))
					***REMOVED***)
					t.Run("unit fixed in ms", func(t *testing.T) ***REMOVED***
						assert.Equal(t, s[2], m.HumanizeValue(v, "ms"))
					***REMOVED***)
					t.Run("unit fixed in µs", func(t *testing.T) ***REMOVED***
						assert.Equal(t, s[3], m.HumanizeValue(v, "us"))
					***REMOVED***)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

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

func TestNewSubmetric(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		parent string
		tags   map[string]string
	***REMOVED******REMOVED***
		"my_metric":                 ***REMOVED***"my_metric", nil***REMOVED***,
		"my_metric***REMOVED******REMOVED***":               ***REMOVED***"my_metric", nil***REMOVED***,
		"my_metric***REMOVED***a***REMOVED***":              ***REMOVED***"my_metric", map[string]string***REMOVED***"a": ""***REMOVED******REMOVED***,
		"my_metric***REMOVED***a:1***REMOVED***":            ***REMOVED***"my_metric", map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		"my_metric***REMOVED*** a : 1 ***REMOVED***":        ***REMOVED***"my_metric", map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		"my_metric***REMOVED***a,b***REMOVED***":            ***REMOVED***"my_metric", map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		"my_metric***REMOVED***a:1,b:2***REMOVED***":        ***REMOVED***"my_metric", map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		"my_metric***REMOVED*** a : 1, b : 2 ***REMOVED***": ***REMOVED***"my_metric", map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		name, data := name, data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			parent, sm := NewSubmetric(name)
			assert.Equal(t, data.parent, parent)
			if data.tags != nil ***REMOVED***
				assert.EqualValues(t, data.tags, sm.Tags.tags)
			***REMOVED*** else ***REMOVED***
				assert.Nil(t, sm.Tags)
			***REMOVED***
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
