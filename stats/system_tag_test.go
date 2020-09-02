/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package stats

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemTagSetMarshalJSON(t *testing.T) ***REMOVED***
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
	tests := []struct ***REMOVED***
		tagset   TagSet
		expected string
	***REMOVED******REMOVED***
		***REMOVED***tagset: TagSet***REMOVED***"ip": true, "proto": true, "group": true, "custom": true***REMOVED***, expected: `["custom","group","ip","proto"]`***REMOVED***,
		***REMOVED***tagset: TagSet***REMOVED******REMOVED***, expected: `[]`***REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := &tc.tagset
		got, err := json.Marshal(ts)
		require.Nil(t, err)
		require.Equal(t, tc.expected, string(got))
	***REMOVED***
***REMOVED***

func TestTagSet_UnmarshalJSON(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		tags []byte
		sets TagSet
	***REMOVED******REMOVED***
		***REMOVED***[]byte(`[]`), TagSet***REMOVED******REMOVED******REMOVED***,
		***REMOVED***[]byte(`["ip","custom", "proto"]`), TagSet***REMOVED***"ip": true, "proto": true, "custom": true***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range tests ***REMOVED***
		ts := new(TagSet)
		require.Nil(t, json.Unmarshal(tc.tags, ts))
		for tag := range tc.sets ***REMOVED***
			assert.True(t, (*ts)[tag])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTagSetTextUnmarshal(t *testing.T) ***REMOVED***
	testMatrix := map[string]TagSet***REMOVED***
		"":                           make(TagSet),
		"ip":                         ***REMOVED***"ip": true***REMOVED***,
		"ip,proto":                   ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,  proto  ":          ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,   ,  proto  ":      ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,,  proto  ,,":       ***REMOVED***"ip": true, "proto": true***REMOVED***,
		"   ip  ,custom,  proto  ,,": ***REMOVED***"ip": true, "custom": true, "proto": true***REMOVED***,
	***REMOVED***

	for input, expected := range testMatrix ***REMOVED***
		set := new(TagSet)
		err := set.UnmarshalText([]byte(input))
		require.NoError(t, err)
		require.Equal(t, expected, *set)
	***REMOVED***
***REMOVED***
