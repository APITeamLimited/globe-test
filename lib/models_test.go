package lib

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/k6-worker/lib/types"
)

func TestStageJSON(t *testing.T) ***REMOVED***
	t.Parallel()
	s := Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***

	data, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, `***REMOVED***"duration":"10s","target":10***REMOVED***`, string(data))

	var s2 Stage
	assert.NoError(t, json.Unmarshal(data, &s2))
	assert.Equal(t, s, s2)
***REMOVED***

// Suggested by @nkovacs in https://github.com/k6io/k6/issues/207#issuecomment-330545467
func TestDataRaces(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Check race", func(t *testing.T) ***REMOVED***
		t.Parallel()
		group, err := NewGroup("test", nil)
		assert.Nil(t, err, "NewGroup")
		wg := sync.WaitGroup***REMOVED******REMOVED***
		wg.Add(2)
		var check1, check2 *Check
		go func() ***REMOVED***
			var err error // using the outer err would result in a data race
			check1, err = group.Check("race")
			assert.Nil(t, err, "Check 1")
			wg.Done()
		***REMOVED***()
		go func() ***REMOVED***
			var err error
			check2, err = group.Check("race")
			assert.Nil(t, err, "Check 2")
			wg.Done()
		***REMOVED***()
		wg.Wait()
		assert.Equal(t, check1, check2, "Checks are the same")
	***REMOVED***)
	t.Run("Group race", func(t *testing.T) ***REMOVED***
		t.Parallel()
		group, err := NewGroup("test", nil)
		assert.Nil(t, err, "NewGroup")
		wg := sync.WaitGroup***REMOVED******REMOVED***
		wg.Add(2)
		var group1, group2 *Group
		go func() ***REMOVED***
			var err error
			group1, err = group.Group("race")
			assert.Nil(t, err, "Group 1")
			wg.Done()
		***REMOVED***()
		go func() ***REMOVED***
			var err error
			group2, err = group.Group("race")
			assert.Nil(t, err, "Group 2")
			wg.Done()
		***REMOVED***()
		wg.Wait()
		assert.Equal(t, group1, group2, "Groups are the same")
	***REMOVED***)
***REMOVED***
