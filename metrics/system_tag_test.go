package metrics

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemTagSetMarshalJSON(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		tagset   SystemTagSet
		expected string
	***REMOVED******REMOVED***
		***REMOVED***TagIP, `["ip"]`***REMOVED***,
		***REMOVED***TagIP | TagProto | TagGroup, `["group","ip","proto"]`***REMOVED***,
		***REMOVED***0, `null`***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := &tc.tagset
		got, err := json.Marshal(ts)
		require.Nil(t, err)
		require.Equal(t, tc.expected, string(got))
	***REMOVED***
***REMOVED***

func TestSystemTagSet_UnmarshalJSON(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		tags []byte
		sets []SystemTagSet
	***REMOVED******REMOVED***
		***REMOVED***[]byte(`[]`), []SystemTagSet***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte(`["ip", "proto"]`), []SystemTagSet***REMOVED***TagIP, TagProto***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := new(SystemTagSet)
		require.Nil(t, json.Unmarshal(tc.tags, ts))
		for _, tag := range tc.sets ***REMOVED***
			assert.True(t, ts.Has(tag))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSystemTagSetTextUnmarshal(t *testing.T) ***REMOVED***
	t.Parallel()

	testMatrix := map[string]SystemTagSet***REMOVED***
		"":                      0,
		"ip":                    TagIP,
		"ip,proto":              TagIP | TagProto,
		"   ip  ,  proto  ":     TagIP | TagProto,
		"   ip  ,   ,  proto  ": TagIP | TagProto,
		"   ip  ,,  proto  ,,":  TagIP | TagProto,
	***REMOVED***

	for input, expected := range testMatrix ***REMOVED***
		set := new(SystemTagSet)
		err := set.UnmarshalText([]byte(input))
		require.NoError(t, err)
		require.Equal(t, expected, *set)
	***REMOVED***
***REMOVED***

func TestTagSetMarshalJSON(t *testing.T) ***REMOVED***
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

func TestTagSet_UnmarshalJSON(t *testing.T) ***REMOVED***
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

func TestTagSetTextUnmarshal(t *testing.T) ***REMOVED***
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
