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
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// NullHostnameTrie is a nullable HostnameTrie, in the same vein as the nullable types provided by
// package gopkg.in/guregu/null.v3
type NullHostnameTrie struct ***REMOVED***
	Trie  *HostnameTrie
	Valid bool
***REMOVED***

// UnmarshalText converts text data to a valid NullHostnameTrie
func (d *NullHostnameTrie) UnmarshalText(data []byte) error ***REMOVED***
	if len(data) == 0 ***REMOVED***
		*d = NullHostnameTrie***REMOVED******REMOVED***
		return nil
	***REMOVED***
	var err error
	d.Trie, err = NewHostnameTrie(strings.Split(string(data), ","))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Valid = true
	return nil
***REMOVED***

// UnmarshalJSON converts JSON data to a valid NullHostnameTrie
func (d *NullHostnameTrie) UnmarshalJSON(data []byte) error ***REMOVED***
	if bytes.Equal(data, []byte(`null`)) ***REMOVED***
		d.Valid = false
		return nil
	***REMOVED***

	var m []string
	var err error
	if err = json.Unmarshal(data, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Trie, err = NewHostnameTrie(m)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	d.Valid = true
	return nil
***REMOVED***

// MarshalJSON implements json.Marshaler interface
func (d NullHostnameTrie) MarshalJSON() ([]byte, error) ***REMOVED***
	if !d.Valid ***REMOVED***
		return []byte(`null`), nil
	***REMOVED***
	return json.Marshal(d.Trie.source)
***REMOVED***

// HostnameTrie is a tree-structured list of hostname matches with support
// for wildcards exclusively at the start of the pattern. Items may only
// be inserted and searched. Internationalized hostnames are valid.
type HostnameTrie struct ***REMOVED***
	source []string

	children map[rune]*HostnameTrie
***REMOVED***

// NewNullHostnameTrie returns a NullHostnameTrie encapsulating HostnameTrie or an error if the
// input is incorrect
func NewNullHostnameTrie(source []string) (NullHostnameTrie, error) ***REMOVED***
	h, err := NewHostnameTrie(source)
	if err != nil ***REMOVED***
		return NullHostnameTrie***REMOVED******REMOVED***, err
	***REMOVED***
	return NullHostnameTrie***REMOVED***
		Valid: true,
		Trie:  h,
	***REMOVED***, nil
***REMOVED***

// NewHostnameTrie returns a pointer to a new HostnameTrie or an error if the input is incorrect
func NewHostnameTrie(source []string) (*HostnameTrie, error) ***REMOVED***
	h := &HostnameTrie***REMOVED***
		source: source,
	***REMOVED***
	for _, s := range h.source ***REMOVED***
		if err := h.insert(s); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return h, nil
***REMOVED***

// Regex description of hostname pattern to enforce blocks by. Global var
// to avoid compilation penalty at runtime.
// based on regex from https://stackoverflow.com/a/106223/5427244
//nolint:gochecknoglobals,lll
var validHostnamePattern *regexp.Regexp = regexp.MustCompile(`^(\*\.?)?((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))?$`)

func isValidHostnamePattern(s string) error ***REMOVED***
	if len(validHostnamePattern.FindString(s)) != len(s) ***REMOVED***
		return errors.Errorf("invalid hostname pattern %s", s)
	***REMOVED***
	return nil
***REMOVED***

// insert a hostname pattern into the given HostnameTrie. Returns an error
// if hostname pattern is invalid.
func (t *HostnameTrie) insert(s string) error ***REMOVED***
	s = strings.ToLower(s)
	if err := isValidHostnamePattern(s); err != nil ***REMOVED***
		return err
	***REMOVED***

	return t.childInsert(s)
***REMOVED***

func (t *HostnameTrie) childInsert(s string) error ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return nil
	***REMOVED***

	// mask creation of the trie by initializing the root here
	if t.children == nil ***REMOVED***
		t.children = make(map[rune]*HostnameTrie)
	***REMOVED***

	rStr := []rune(s) // need to iterate by runes for intl' names
	last := len(rStr) - 1
	if c, ok := t.children[rStr[last]]; ok ***REMOVED***
		return c.childInsert(string(rStr[:last]))
	***REMOVED***

	t.children[rStr[last]] = &HostnameTrie***REMOVED***children: make(map[rune]*HostnameTrie)***REMOVED***
	return t.children[rStr[last]].childInsert(string(rStr[:last]))
***REMOVED***

// Contains returns whether s matches a pattern in the HostnameTrie
// along with the matching pattern, if one was found.
func (t *HostnameTrie) Contains(s string) (matchedPattern string, matchFound bool) ***REMOVED***
	s = strings.ToLower(s)
	if len(s) == 0 ***REMOVED***
		if len(t.children) == 0 ***REMOVED***
			return "", true
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rStr := []rune(s)
		last := len(rStr) - 1
		if c, ok := t.children[rStr[last]]; ok ***REMOVED***
			if match, matched := c.Contains(string(rStr[:last])); matched ***REMOVED***
				return match + string(rStr[last]), true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if _, wild := t.children['*']; wild ***REMOVED***
		return "*", true
	***REMOVED***

	return "", false
***REMOVED***
