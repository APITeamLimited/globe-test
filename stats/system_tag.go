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
	"bytes"
	"encoding/json"
	"strings"
)

// SystemTagSet is a bitmask that is used to keep track
// which system tags should be included with which metrics.
//go:generate enumer -type=SystemTagSet -transform=snake -trimprefix=Tag -output system_tag_set_gen.go
type SystemTagSet uint32

// TagSet is a string to bool map (for lookup efficiency) that is used to keep track
// which system tags and non-system tags should be included with with metrics.
type TagSet map[string]bool

//nolint: golint
const (
	// Default system tags includes all of the system tags emitted with metrics by default.
	TagProto SystemTagSet = 1 << iota
	TagSubproto
	TagStatus
	TagMethod
	TagURL
	TagName
	TagGroup
	TagCheck
	TagError
	TagErrorCode
	TagTLSVersion
	TagScenario

	// System tags not enabled by default.
	TagIter
	TagVU
	TagOCSPStatus
	TagIP
)

// DefaultSystemTagSet includes all of the system tags emitted with metrics by default.
// Other tags that are not enabled by default include: iter, vu, ocsp_status, ip
//nolint:gochecknoglobals
var DefaultSystemTagSet = TagProto | TagSubproto | TagStatus | TagMethod | TagURL | TagName | TagGroup |
	TagCheck | TagCheck | TagError | TagErrorCode | TagTLSVersion | TagScenario

// Add adds a tag to tag set.
func (i *SystemTagSet) Add(tag SystemTagSet) ***REMOVED***
	if i == nil ***REMOVED***
		i = new(SystemTagSet)
	***REMOVED***
	*i |= tag
***REMOVED***

// Has checks a tag included in tag set.
func (i *SystemTagSet) Has(tag SystemTagSet) bool ***REMOVED***
	if i == nil ***REMOVED***
		return false
	***REMOVED***
	return *i&tag != 0
***REMOVED***

// Map returns the TagSet with current value from SystemTagSet
func (i SystemTagSet) Map() TagSet ***REMOVED***
	m := TagSet***REMOVED******REMOVED***
	for _, tag := range SystemTagSetValues() ***REMOVED***
		if i.Has(tag) ***REMOVED***
			m[tag.String()] = true
		***REMOVED***
	***REMOVED***
	return m
***REMOVED***

// SetString returns comma separated list of the string representation of all values in the set
func (i SystemTagSet) SetString() string ***REMOVED***
	var keys []string
	for _, tag := range SystemTagSetValues() ***REMOVED***
		if i.Has(tag) ***REMOVED***
			keys = append(keys, tag.String())
		***REMOVED***
	***REMOVED***
	return strings.Join(keys, ",")
***REMOVED***

// ToSystemTagSet converts list of tags to SystemTagSet
// TODO: emit error instead of discarding invalid values.
func ToSystemTagSet(tags []string) *SystemTagSet ***REMOVED***
	ts := new(SystemTagSet)
	for _, tag := range tags ***REMOVED***
		if v, err := SystemTagSetString(tag); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return ts
***REMOVED***

// NewSystemTagSet returns a SystemTagSet from input.
func NewSystemTagSet(tags ...SystemTagSet) *SystemTagSet ***REMOVED***
	ts := new(SystemTagSet)
	for _, tag := range tags ***REMOVED***
		ts.Add(tag)
	***REMOVED***
	return ts
***REMOVED***

// MarshalJSON converts the SystemTagSet to a list (JS array).
func (i *SystemTagSet) MarshalJSON() ([]byte, error) ***REMOVED***
	var tags []string
	for _, tag := range SystemTagSetValues() ***REMOVED***
		if i.Has(tag) ***REMOVED***
			tags = append(tags, tag.String())
		***REMOVED***
	***REMOVED***
	return json.Marshal(tags)
***REMOVED***

// UnmarshalJSON converts the tag list back to expected tag set.
func (i *SystemTagSet) UnmarshalJSON(data []byte) error ***REMOVED***
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(tags) != 0 ***REMOVED***
		*i = *ToSystemTagSet(tags)
	***REMOVED***

	return nil
***REMOVED***

// UnmarshalText converts the tag list to SystemTagSet.
func (i *SystemTagSet) UnmarshalText(data []byte) error ***REMOVED***
	var list = bytes.Split(data, []byte(","))

	for _, key := range list ***REMOVED***
		key := strings.TrimSpace(string(key))
		if key == "" ***REMOVED***
			continue
		***REMOVED***
		if v, err := SystemTagSetString(key); err == nil ***REMOVED***
			i.Add(v)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
