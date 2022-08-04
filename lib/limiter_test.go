package lib

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlotLimiterSingleSlot(t *testing.T) ***REMOVED***
	t.Parallel()
	l := NewSlotLimiter(1)
	l.Begin()
	done := false
	go func() ***REMOVED***
		done = true
		l.End()
	***REMOVED***()
	l.Begin()
	assert.True(t, done)
	l.End()
***REMOVED***

func TestSlotLimiterUnlimited(t *testing.T) ***REMOVED***
	t.Parallel()
	l := NewSlotLimiter(0)
	l.Begin()
	l.Begin()
	l.Begin()
***REMOVED***

func TestSlotLimiters(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct***REMOVED*** limit, launches, expMid int ***REMOVED******REMOVED***
		***REMOVED***0, 0, 0***REMOVED***,
		***REMOVED***0, 1, 1***REMOVED***,
		***REMOVED***0, 5, 5***REMOVED***,
		***REMOVED***1, 5, 1***REMOVED***,
		***REMOVED***2, 5, 2***REMOVED***,
		***REMOVED***5, 6, 5***REMOVED***,
		***REMOVED***6, 5, 5***REMOVED***,
		***REMOVED***10, 7, 7***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("limit=%d,launches=%d", tc.limit, tc.launches), func(t *testing.T) ***REMOVED***
			t.Parallel()
			l := NewSlotLimiter(tc.limit)
			wg := sync.WaitGroup***REMOVED******REMOVED***

			switch ***REMOVED***
			case tc.limit == 0:
				wg.Add(tc.launches)
			case tc.launches < tc.limit:
				wg.Add(tc.launches)
			default:
				wg.Add(tc.limit)
			***REMOVED***

			var counter uint32

			for i := 0; i < tc.launches; i++ ***REMOVED***
				go func() ***REMOVED***
					l.Begin()
					atomic.AddUint32(&counter, 1)
					wg.Done()
				***REMOVED***()
			***REMOVED***
			wg.Wait()
			assert.Equal(t, uint32(tc.expMid), atomic.LoadUint32(&counter))

			if tc.limit != 0 && tc.limit < tc.launches ***REMOVED***
				wg.Add(tc.launches - tc.limit)
				for i := 0; i < tc.launches; i++ ***REMOVED***
					l.End()
				***REMOVED***
				wg.Wait()
				assert.Equal(t, uint32(tc.launches), atomic.LoadUint32(&counter))
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMultiSlotLimiter(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("0", func(t *testing.T) ***REMOVED***
		t.Parallel()
		l := NewMultiSlotLimiter(0)
		assert.Nil(t, l.Slot("test"))
	***REMOVED***)
	t.Run("1", func(t *testing.T) ***REMOVED***
		t.Parallel()
		l := NewMultiSlotLimiter(1)
		assert.Equal(t, l.Slot("test"), l.Slot("test"))
		assert.NotNil(t, l.Slot("test"))
	***REMOVED***)
	t.Run("2", func(t *testing.T) ***REMOVED***
		t.Parallel()
		l := NewMultiSlotLimiter(1)
		wg := sync.WaitGroup***REMOVED******REMOVED***
		wg.Add(2)

		var s1, s2 SlotLimiter
		go func() ***REMOVED***
			s1 = l.Slot("ctest")
			wg.Done()
		***REMOVED***()
		go func() ***REMOVED***
			s2 = l.Slot("ctest")
			wg.Done()
		***REMOVED***()
		wg.Wait()

		assert.NotNil(t, s1)
		assert.Equal(t, s1, s2)
		assert.Equal(t, s1, l.Slot("ctest"))
		assert.NotEqual(t, s1, l.Slot("dtest"))
		assert.NotNil(t, l.Slot("dtest"))
	***REMOVED***)
***REMOVED***
