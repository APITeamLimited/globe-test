package accumulate

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNonexistent(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	assert.Nil(t, b.Data[&stat]["value"])
***REMOVED***

func TestGet(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)

	assert.NotNil(t, b.Data[&stat]["value"])
***REMOVED***

func TestSubmitSortsValues(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)

	dim := b.Data[&stat]["value"]
	assert.EqualValues(t, []float64***REMOVED***1, 2, 3***REMOVED***, dim.Values)
	assert.False(t, dim.dirty)
***REMOVED***

func TestSubmitSortsValuesContinously(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 6***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 5***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 4***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)

	dim := b.Data[&stat]["value"]
	assert.EqualValues(t, []float64***REMOVED***1, 2, 3, 4, 5, 6***REMOVED***, dim.Values)
	assert.False(t, dim.dirty)
***REMOVED***

func TestSubmitKeepsLast(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	assert.Equal(t, float64(2), b.Data[&stat]["value"].Last)
***REMOVED***

func TestSubmitIgnoresExcluded(t *testing.T) ***REMOVED***
	b := New()
	stat1 := stats.Stat***REMOVED***Name: "test"***REMOVED***
	stat2 := stats.Stat***REMOVED***Name: "test2"***REMOVED***
	b.Exclude["test2"] = true
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat1, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat1, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat2, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	assert.Len(t, b.Data, 1)
***REMOVED***

func TestSubmitIgnoresNotInOnly(t *testing.T) ***REMOVED***
	b := New()
	stat1 := stats.Stat***REMOVED***Name: "test"***REMOVED***
	stat2 := stats.Stat***REMOVED***Name: "test2"***REMOVED***
	b.Only["test2"] = true
	b.Submit([][]stats.Sample***REMOVED***
		[]stats.Sample***REMOVED***
			stats.Sample***REMOVED***Stat: &stat1, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat1, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Sample***REMOVED***Stat: &stat2, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	assert.Len(t, b.Data, 1)
***REMOVED***

func TestGetVStatDefault(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	assert.Equal(t, &stat, b.getVStat(&stat, stats.Tags***REMOVED******REMOVED***))
***REMOVED***

func TestGetVStatNoMatch(t *testing.T) ***REMOVED***
	b := New()
	b.GroupBy = []string***REMOVED***"no-match"***REMOVED***
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	assert.Equal(t, &stat, b.getVStat(&stat, stats.Tags***REMOVED******REMOVED***))
***REMOVED***

func TestGetVStatOneTag(t *testing.T) ***REMOVED***
	b := New()
	b.GroupBy = []string***REMOVED***"tag"***REMOVED***
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	vstat := b.getVStat(&stat, stats.Tags***REMOVED***"tag": "value"***REMOVED***)
	assert.NotNil(t, vstat)
	assert.Equal(t, "test***REMOVED***tag: value***REMOVED***", vstat.Name)
***REMOVED***

func TestGetVStatTwoTags(t *testing.T) ***REMOVED***
	b := New()
	b.GroupBy = []string***REMOVED***"tag", "blah"***REMOVED***
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	vstat := b.getVStat(&stat, stats.Tags***REMOVED***"tag": "value", "blah": 12345***REMOVED***)
	assert.NotNil(t, vstat)
	assert.Equal(t, "test***REMOVED***tag: value, blah: 12345***REMOVED***", vstat.Name)
***REMOVED***

func TestGetVStatTwoTagsOneMiss(t *testing.T) ***REMOVED***
	b := New()
	b.GroupBy = []string***REMOVED***"tag", "weh", "blah"***REMOVED***
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	vstat := b.getVStat(&stat, stats.Tags***REMOVED***"tag": "value", "blah": 12345***REMOVED***)
	assert.NotNil(t, vstat)
	assert.Equal(t, "test***REMOVED***tag: value, blah: 12345***REMOVED***", vstat.Name)
***REMOVED***
