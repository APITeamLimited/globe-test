/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHostnameTrieInsert(t *testing.T) ***REMOVED***
	t.Parallel()

	hostnames, err := NewHostnameTrie([]string***REMOVED***"foo.bar"***REMOVED***)
	assert.NoError(t, err)

	assert.NoError(t, hostnames.insert("test.k6.io"))
	assert.Error(t, hostnames.insert("inval*d.pattern"))
	assert.NoError(t, hostnames.insert("*valid.pattern"))
***REMOVED***

func TestHostnameTrieContains(t *testing.T) ***REMOVED***
	t.Parallel()

	trie, err := NewHostnameTrie([]string***REMOVED***"sub.test.k6.io", "test.k6.io", "*valid.pattern", "sub.valid.pattern"***REMOVED***)
	require.NoError(t, err)
	cases := map[string]string***REMOVED***
		"K6.Io":                 "",
		"tEsT.k6.Io":            "test.k6.io",
		"TESt.K6.IO":            "test.k6.io",
		"sub.test.k6.io":        "sub.test.k6.io",
		"sub.sub.test.k6.io":    "",
		"blocked.valId.paTtern": "*valid.pattern",
		"valId.paTtern":         "*valid.pattern",
		"sub.valid.pattern":     "sub.valid.pattern", // use the most specific blocker
		"www.sub.valid.pattern": "*valid.pattern",
		"example.test.k6.io":    "",
	***REMOVED***
	for key, value := range cases ***REMOVED***
		host, pattern := key, value
		t.Run(host, func(t *testing.T) ***REMOVED***
			t.Parallel()

			match, matches := trie.Contains(host)
			if pattern == "" ***REMOVED***
				assert.False(t, matches)
				assert.Empty(t, match)
			***REMOVED*** else ***REMOVED***
				assert.True(t, matches)
				assert.Equal(t, pattern, match)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestNullHostnameTrieSource(t *testing.T) ***REMOVED***
	t.Parallel()

	trie, err := NewNullHostnameTrie([]string***REMOVED***"sub.test.k6.io", "test.k6.io", "*valid.pattern", "sub.valid.pattern"***REMOVED***)

	require.NoError(t, err)

	assert.Equal(t, []string***REMOVED***"sub.test.k6.io", "test.k6.io", "*valid.pattern", "sub.valid.pattern"***REMOVED***, trie.Source())
***REMOVED***
