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

// SystemTagMap is a string to bool map (for lookup efficiency) that is used to keep track
// which system tags should be included with with metrics.
type SystemTagMap map[string]bool

// DefaultSystemTagList includes all of the system tags emitted with metrics by default.
// Other tags that are not enabled by default include: iter, vu, ocsp_status, ip
var DefaultSystemTagList = []string***REMOVED***
	TagProto.String(),
	TagSubProto.String(),
	TagStatus.String(),
	TagMethod.String(),
	TagURL.String(),
	TagName.String(),
	TagGroup.String(),
	TagCheck.String(),
	TagCheck.String(),
	TagError.String(),
	TagErrorCode.String(),
	TagTLSVersion.String(),
***REMOVED***

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

// Map returns the SystemTagMap with current value from SystemTagSet
func (ts *SystemTagSet) Map() SystemTagMap ***REMOVED***
	m := SystemTagMap***REMOVED******REMOVED***
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
