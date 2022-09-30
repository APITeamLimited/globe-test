package workerMetrics

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnabledTagsMarshalJSON(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		tagset   EnabledTags
		expected string
	***REMOVED******REMOVED***
		***REMOVED***tagset: EnabledTags***REMOVED***"ip": true, "proto": true, "group": true, "custom": true***REMOVED***, expected: `["custom","group","ip","proto"]`***REMOVED***,
		***REMOVED***tagset: EnabledTags***REMOVED******REMOVED***, expected: `[]`***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := &tc.tagset
		got, err := json.Marshal(ts)
		require.Nil(t, err)
		require.Equal(t, tc.expected, string(got))
	***REMOVED***
***REMOVED***

func TestEnabledTagsUnmarshalJSON(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		tags []byte
		sets EnabledTags
	***REMOVED******REMOVED***
		***REMOVED***[]byte(`[]`), EnabledTags***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte(`["ip","custom", "proto"]`), EnabledTags***REMOVED***"ip": true, "proto": true, "custom": true***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := new(EnabledTags)
		require.Nil(t, json.Unmarshal(tc.tags, ts))
		for tag := range tc.sets ***REMOVED***
			assert.True(t, (*ts)[tag])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEnabledTagsTextUnmarshal(t *testing.T) ***REMOVED***
	t.Parallel()

	testMatrix := map[string]EnabledTags***REMOVED***
		"":                           make(EnabledTags),
		"ip":                         ***REMOVED***"ip": true***REMOVED***,
		"ip,proto":                   ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,  proto  ":          ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,   ,  proto  ":      ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,,  proto  ,,":       ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,custom,  proto  ,,": ***REMOVED***"ip": true, "custom": true, "proto": true***REMOVED***,
	***REMOVED***

	for input, expected := range testMatrix ***REMOVED***
		set := new(EnabledTags)
		err := set.UnmarshalText([]byte(input))
		require.NoError(t, err)
		require.Equal(t, expected, *set)
	***REMOVED***
***REMOVED***
