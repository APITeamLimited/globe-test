package accumulate

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNonexistent(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	assert.Nil(t, b.Get(&stat, "value"))
***REMOVED***

func TestGet(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)

	assert.NotNil(t, b.Get(&stat, "value"))
***REMOVED***

func TestSubmitInternsNames(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	assert.Len(t, b.interned, 1)
	assert.Len(t, b.Data, 1)
	assert.Len(t, b.Data[&stat], 1)
	assert.Contains(t, b.Data[&stat], b.interned["value"])
***REMOVED***

func TestSubmitSortsValues(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)

	dim := b.Get(&stat, "value")
	assert.EqualValues(t, []float64***REMOVED***1, 2, 3***REMOVED***, dim.Values)
	assert.False(t, dim.dirty)
***REMOVED***

func TestSubmitSortsValuesContinously(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 6***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 5***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 4***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)

	dim := b.Get(&stat, "value")
	assert.EqualValues(t, []float64***REMOVED***1, 2, 3, 4, 5, 6***REMOVED***, dim.Values)
	assert.False(t, dim.dirty)
***REMOVED***

func TestSubmitKeepsLast(t *testing.T) ***REMOVED***
	b := New()
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	assert.Equal(t, float64(2), b.Get(&stat, "value").Last)
***REMOVED***

func TestSubmitIgnoresExcluded(t *testing.T) ***REMOVED***
	b := New()
	stat1 := stats.Stat***REMOVED***Name: "test"***REMOVED***
	stat2 := stats.Stat***REMOVED***Name: "test2"***REMOVED***
	b.Exclude["test2"] = true
	b.Submit([][]stats.Point***REMOVED***
		[]stats.Point***REMOVED***
			stats.Point***REMOVED***Stat: &stat1, Values: stats.Values***REMOVED***"value": 3***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat1, Values: stats.Values***REMOVED***"value": 1***REMOVED******REMOVED***,
			stats.Point***REMOVED***Stat: &stat2, Values: stats.Values***REMOVED***"value": 2***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***)
	assert.Len(t, b.Data, 1)
***REMOVED***
