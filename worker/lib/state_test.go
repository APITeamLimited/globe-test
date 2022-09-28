package lib

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTagMapSet(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("Sync", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tm := NewTagMap(nil)
		tm.Set("mytag", "42")
		v, found := tm.Get("mytag")
		assert.True(t, found)
		assert.Equal(t, "42", v)
	***REMOVED***)

	t.Run("Safe-Concurrent", func(t *testing.T) ***REMOVED***
		t.Parallel()
		tm := NewTagMap(nil)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() ***REMOVED***
			count := 0
			for ***REMOVED***
				select ***REMOVED***
				case <-time.Tick(1 * time.Millisecond):
					count++
					tm.Set("mytag", strconv.Itoa(count))

				case <-ctx.Done():
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()

		go func() ***REMOVED***
			for ***REMOVED***
				select ***REMOVED***
				case <-time.Tick(1 * time.Millisecond):
					tm.Get("mytag")

				case <-ctx.Done():
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***()

		time.Sleep(100 * time.Millisecond)
	***REMOVED***)
***REMOVED***

func TestTagMapGet(t *testing.T) ***REMOVED***
	t.Parallel()
	tm := NewTagMap(map[string]string***REMOVED***
		"key1": "value1",
	***REMOVED***)
	v, ok := tm.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", v)
***REMOVED***

func TestTagMapLen(t *testing.T) ***REMOVED***
	t.Parallel()
	tm := NewTagMap(map[string]string***REMOVED***
		"key1": "value1",
		"key2": "value2",
	***REMOVED***)
	assert.Equal(t, 2, tm.Len())
***REMOVED***

func TestTagMapDelete(t *testing.T) ***REMOVED***
	t.Parallel()
	m := map[string]string***REMOVED***
		"key1": "value1",
		"key2": "value2",
	***REMOVED***
	tm := NewTagMap(m)
	tm.Delete("key1")
	_, ok := m["key1"]
	assert.False(t, ok)
***REMOVED***

func TestTagMapClone(t *testing.T) ***REMOVED***
	t.Parallel()
	tm := NewTagMap(map[string]string***REMOVED***
		"key1": "value1",
		"key2": "value2",
	***REMOVED***)
	m := tm.Clone()
	assert.Equal(t, map[string]string***REMOVED***
		"key1": "value1",
		"key2": "value2",
	***REMOVED***, m)
***REMOVED***
