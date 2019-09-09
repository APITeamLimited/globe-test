package tagset

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagSetMarshalJSON(t *testing.T) ***REMOVED***
	var tests = []struct ***REMOVED***
		tagset   TagSet
		expected string
	***REMOVED******REMOVED***
		***REMOVED***IP, `["ip"]`***REMOVED***,
		***REMOVED***0, `null`***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := &tc.tagset
		got, err := json.Marshal(ts)
		require.Nil(t, err)
		require.Equal(t, tc.expected, string(got))
	***REMOVED***

***REMOVED***

func TestTagSet_UnmarshalJSON(t *testing.T) ***REMOVED***
	var tests = []struct ***REMOVED***
		tags []byte
		sets []TagSet
	***REMOVED******REMOVED***
		***REMOVED***[]byte(`[]`), []TagSet***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte(`["ip", "proto"]`), []TagSet***REMOVED***IP, Proto***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := new(TagSet)
		require.Nil(t, json.Unmarshal(tc.tags, ts))
		for _, tag := range tc.sets ***REMOVED***
			assert.True(t, ts.Has(tag))
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestTagSetTextUnmarshal(t *testing.T) ***REMOVED***
	var testMatrix = map[string]TagSet***REMOVED***
		"":                      0,
		"ip":                    IP,
		"ip,proto":              IP | Proto,
		"   ip  ,  proto  ":     IP | Proto,
		"   ip  ,   ,  proto  ": IP | Proto,
		"   ip  ,,  proto  ,,":  IP | Proto,
	***REMOVED***

	for input, expected := range testMatrix ***REMOVED***
		var set = new(TagSet)
		err := set.UnmarshalText([]byte(input))
		require.NoError(t, err)
		require.Equal(t, expected, *set)
	***REMOVED***
***REMOVED***
