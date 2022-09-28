package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO: update test
/*
func TestSumStages(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		Time   types.NullDuration
		Stages []Stage
	***REMOVED******REMOVED***
		"Blank":    ***REMOVED***types.NullDuration***REMOVED******REMOVED***, []Stage***REMOVED******REMOVED******REMOVED***,
		"Infinite": ***REMOVED***types.NullDuration***REMOVED******REMOVED***, []Stage***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		"Limit": ***REMOVED***
			types.NullDurationFrom(10 * time.Second),
			[]Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"InfiniteTail": ***REMOVED***
			types.NullDuration***REMOVED***Duration: types.Duration(10 * time.Second), Valid: false***REMOVED***,
			[]Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***,
				***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, data.Time, SumStages(data.Stages))
		***REMOVED***)
	***REMOVED***
***REMOVED***
*/

func TestMin(t *testing.T) ***REMOVED***
	t.Parallel()
	assert.Equal(t, int64(10), Min(10, 100))
	assert.Equal(t, int64(10), Min(100, 10))
***REMOVED***

func TestMax(t *testing.T) ***REMOVED***
	t.Parallel()
	assert.Equal(t, int64(100), Max(10, 100))
	assert.Equal(t, int64(100), Max(100, 10))
***REMOVED***
