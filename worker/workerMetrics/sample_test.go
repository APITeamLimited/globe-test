package workerMetrics

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	assert.False(t, tags.Contains(IntoSampleTags(&map[string]string***REMOVED***"nonexistent_key": ""***REMOVED***)))
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
		Metric: newMetric("test_metric", Counter),
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

func TestGetResolversForTrendColumnsCalculation(t *testing.T) ***REMOVED***
	t.Parallel()

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

	for _, tc := range customResolversTests ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%v", tc.stats), func(t *testing.T) ***REMOVED***
			t.Parallel()
			sink := createTestTrendSink(100)

			res, err := GetResolversForTrendColumns([]string***REMOVED***tc.stats***REMOVED***)
			assert.NoError(t, err)
			assert.Len(t, res, 1)
			for k := range res ***REMOVED***
				assert.InDelta(t, sink.P(tc.percentile), res[k](sink), 0.000001)
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
