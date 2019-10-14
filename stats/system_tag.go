package stats

import (
	"bytes"
	"encoding/json"
	"strings"
)

// SystemTagSet is a bitmask that is used to keep track
// which system tags should be included with which metrics.
type SystemTagSet uint32

// TagSet is a string to bool map (for lookup efficiency) that is used to keep track
// which system tags and non-system tags should be included with with metrics.
type TagSet map[string]bool

//nolint: golint
const (
	// Default system tags includes all of the system tags emitted with metrics by default.
	TagProto SystemTagSet = 1 << iota
	TagSubProto
	TagStatus
	TagMethod
	TagURL
	TagName
	TagGroup
	TagCheck
	TagError
	TagErrorCode
	TagTLSVersion

	// System tags not enabled by default.
	TagIter
	TagVU
	TagOCSPStatus
	TagIP
)

// DefaultSystemTagSet includes all of the system tags emitted with metrics by default.
// Other tags that are not enabled by default include: iter, vu, ocsp_status, ip
//nolint:gochecknoglobals
var DefaultSystemTagSet = TagProto | TagSubProto | TagStatus | TagMethod | TagURL | TagName | TagGroup |
	TagCheck | TagCheck | TagError | TagErrorCode | TagTLSVersion

// ToSystemTagSet converts a tag set to tag set
// Non-system tag will be ignored.
func (ts TagSet) ToSystemTagSet() SystemTagSet ***REMOVED***
	sts := SystemTagSet(0)
	for tag, ok := range ts ***REMOVED***
		if v, err := SystemTagSetString(tag); err == nil && ok ***REMOVED***
			sts.Add(v)
		***REMOVED***
	***REMOVED***
	return sts
***REMOVED***

// Add adds a tag to tag set.
func (ts *SystemTagSet) Add(tag SystemTagSet) ***REMOVED***
	if ts == nil ***REMOVED***
		ts = new(SystemTagSet)
	***REMOVED***
	*ts |= tag
***REMOVED***

// Has checks a tag included in tag set.
func (ts *SystemTagSet) Has(tag SystemTagSet) bool ***REMOVED***
	if ts == nil ***REMOVED***
		return false
	***REMOVED***
	return *ts&tag != 0
***REMOVED***

// Map returns the TagSet with current value from SystemTagSet
func (ts *SystemTagSet) Map() TagSet ***REMOVED***
	m := TagSet***REMOVED******REMOVED***
	for _, tag := range SystemTagSetValues() ***REMOVED***
		if ts.Has(tag) ***REMOVED***
			m[tag.String()] = true
		***REMOVED***
	***REMOVED***
	return m
***REMOVED***

// ToSystemTagSet converts list of tags to SystemTagSet
func ToSystemTagSet(tags []string) *SystemTagSet ***REMOVED***
	ts := SystemTagSet(0)
	for _, tag := range tags ***REMOVED***
		if v, err := SystemTagSetString(tag); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return &ts
***REMOVED***

// MarshalJSON converts the SystemTagSet to a list (JS array).
func (ts *SystemTagSet) MarshalJSON() ([]byte, error) ***REMOVED***
	var tags []string
	for _, tag := range SystemTagSetValues() ***REMOVED***
		if ts.Has(tag) ***REMOVED***
			tags = append(tags, tag.String())
		***REMOVED***
	***REMOVED***
	return json.Marshal(tags)
***REMOVED***

// UnmarshalJSON converts the tag list back to expected tag set.
func (ts *SystemTagSet) UnmarshalJSON(data []byte) error ***REMOVED***
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(tags) != 0 ***REMOVED***
		*ts = *ToSystemTagSet(tags)
	***REMOVED***

	return nil
***REMOVED***

// UnmarshalText converts the tag list to SystemTagSet.
func (ts *SystemTagSet) UnmarshalText(data []byte) error ***REMOVED***
	var list = bytes.Split(data, []byte(","))

	for _, key := range list ***REMOVED***
		key := strings.TrimSpace(string(key))
		if key == "" ***REMOVED***
			continue
		***REMOVED***
		if v, err := SystemTagSetString(key); err == nil ***REMOVED***
			ts.Add(v)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
